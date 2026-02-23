package database

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

// Init 初始化数据库连接和表结构
func Init() error {
	var err error
	dbPath := os.Getenv("MEMO_DB_PATH")
	if dbPath == "" {
		dbPath = "./notes.db"
	}
	
	// 确保数据库文件所在目录存在
if dbDir := filepath.Dir(dbPath); dbDir != "." {
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return fmt.Errorf("无法创建数据库目录 %s: %w", dbDir, err)
		}
	}
	
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	// 测试连接
	if err = DB.Ping(); err != nil {
		return err
	}

	// 推荐的 SQLite pragma（不影响兼容性）
	// 注意：FTS5 需要通过 go build tag sqlite_fts5 启用
	if _, err := DB.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return err
	}
	_, _ = DB.Exec(`PRAGMA journal_mode = WAL;`)
	_, _ = DB.Exec(`PRAGMA busy_timeout = 5000;`)

	// 创建表 & 迁移
	if err := runMigrations(); err != nil {
		return err
	}

	log.Println("数据库初始化成功")
	return nil
}

// runMigrations 创建数据库表并执行迁移
func runMigrations() error {
	// 关键：DDL/Schema 操作强制走同一个连接，避免连接池导致的 schema 可见性问题
	conn, err := DB.Conn(context.Background())
	if err != nil {
		return err
	}
	defer conn.Close()
	ctx := context.Background()

	// 获取当前 user_version
	var ver int
	if err := conn.QueryRowContext(ctx, `PRAGMA user_version;`).Scan(&ver); err != nil {
		return err
	}

	// v1：基础 schema（notes/tags/users + FTS5）
	if ver < 1 {
		if err := ensureSchemaV1(ctx, conn); err != nil {
			return err
		}
		if _, err := conn.ExecContext(ctx, `PRAGMA user_version = 1;`); err != nil {
			return err
		}
		ver = 1
	}

	// v2：notes 扩展字段（pinned/content_type/user_id）
	if ver < 2 {
		if err := ensureNotesColumnsV2(ctx, conn); err != nil {
			return err
		}
		if _, err := conn.ExecContext(ctx, `PRAGMA user_version = 2;`); err != nil {
			return err
		}
		ver = 2
	}

	// v3：resources（附件）表 + note_resources 关联表
	if ver < 3 {
		if err := ensureResourcesSchemaV3(ctx, conn); err != nil {
			return err
		}
		if _, err := conn.ExecContext(ctx, `PRAGMA user_version = 3;`); err != nil {
			return err
		}
		ver = 3
	}

	// v4：users 增加 is_admin，并初始化默认管理员
	if ver < 4 {
		if err := ensureUsersAdminV4(ctx, conn); err != nil {
			return err
		}
		if _, err := conn.ExecContext(ctx, `PRAGMA user_version = 4;`); err != nil {
			return err
		}
		ver = 4
	}

	// v5：users 增加 must_change_password，并更新/修复默认管理员策略
	if ver < 5 {
		if err := ensureUsersBootstrapV5(ctx, conn); err != nil {
			return err
		}
		if _, err := conn.ExecContext(ctx, `PRAGMA user_version = 5;`); err != nil {
			return err
		}
		ver = 5
	}

	// v6：tags 增加 user_id 并迁移历史数据；同时把历史 notes.user_id 迁移到管理员
	if ver < 6 {
		if err := ensureMultiUserIsolationV6(ctx, conn); err != nil {
			return err
		}
		if _, err := conn.ExecContext(ctx, `PRAGMA user_version = 6;`); err != nil {
			return err
		}
		ver = 6
	}

	// v7：移除 tags.name 全局 UNIQUE，改为 (user_id, name) 唯一
	if ver < 7 {
		if err := ensureTagsUniquePerUserV7(ctx, conn); err != nil {
			return err
		}
		if _, err := conn.ExecContext(ctx, `PRAGMA user_version = 7;`); err != nil {
			return err
		}
		ver = 7
	}

	// v8：笔记本（notebooks）表 + note_notebooks 关联表
	if ver < 8 {
		if err := ensureNotebooksV8(ctx, conn); err != nil {
			return err
		}
		if _, err := conn.ExecContext(ctx, `PRAGMA user_version = 8;`); err != nil {
			return err
		}
		ver = 8
	}

	// v9：笔记位置（location）字段
	if ver < 9 {
		if err := ensureLocationV9(ctx, conn); err != nil {
			return err
		}
		if _, err := conn.ExecContext(ctx, `PRAGMA user_version = 9;`); err != nil {
			return err
		}
		ver = 9
	}

	return nil
}

func ensureNotebooksV8(ctx context.Context, conn *sql.Conn) error {
	notebooksTable := `
	CREATE TABLE IF NOT EXISTS notebooks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		color TEXT,
		sort_order INTEGER NOT NULL DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`
	noteNotebooksTable := `
	CREATE TABLE IF NOT EXISTS note_notebooks (
		note_id INTEGER NOT NULL,
		notebook_id INTEGER NOT NULL,
		PRIMARY KEY (note_id, notebook_id),
		FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE,
		FOREIGN KEY (notebook_id) REFERENCES notebooks(id) ON DELETE CASCADE
	);`
	if _, err := conn.ExecContext(ctx, notebooksTable); err != nil {
		return err
	}
	if _, err := conn.ExecContext(ctx, noteNotebooksTable); err != nil {
		return err
	}
	_, _ = conn.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS idx_notebooks_user_id ON notebooks(user_id);`)
	_, _ = conn.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS idx_note_notebooks_notebook_id ON note_notebooks(notebook_id);`)
	return nil
}

// v9：笔记位置字段
func ensureLocationV9(ctx context.Context, conn *sql.Conn) error {
	// location 字段：存储地点名称
	if ok, err := columnExists(ctx, conn, "notes", "location"); err != nil {
		return err
	} else if !ok {
		if _, err := conn.ExecContext(ctx, `ALTER TABLE notes ADD COLUMN location TEXT;`); err != nil {
			return err
		}
	}

	// latitude 字段：纬度
	if ok, err := columnExists(ctx, conn, "notes", "latitude"); err != nil {
		return err
	} else if !ok {
		if _, err := conn.ExecContext(ctx, `ALTER TABLE notes ADD COLUMN latitude REAL;`); err != nil {
			return err
		}
	}

	// longitude 字段：经度
	if ok, err := columnExists(ctx, conn, "notes", "longitude"); err != nil {
		return err
	} else if !ok {
		if _, err := conn.ExecContext(ctx, `ALTER TABLE notes ADD COLUMN longitude REAL;`); err != nil {
			return err
		}
	}

	return nil
}

func ensureSchemaV1(ctx context.Context, conn *sql.Conn) error {
	// 创建笔记表
	notesTable := `
	CREATE TABLE IF NOT EXISTS notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		content TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// 全文检索表（FTS5）
	//
	// 注意：FTS5 在 mattn/go-sqlite3 需要以 build tag 启用（sqlite_fts5）。
	// 这里使用 rowid 与 notes.id 对齐，并通过触发器维护一致性。
	notesFTSTable := `
	CREATE VIRTUAL TABLE IF NOT EXISTS notes_fts
	USING fts5(content, note_id UNINDEXED, tokenize='unicode61');`

	notesFTSTriggers := []string{
		// 新增
		`CREATE TRIGGER IF NOT EXISTS notes_ai AFTER INSERT ON notes BEGIN
			INSERT INTO notes_fts(rowid, content, note_id) VALUES (new.id, COALESCE(new.content, ''), new.id);
		END;`,
		// 删除（直接 DELETE，避免特殊语法兼容问题）
		`CREATE TRIGGER IF NOT EXISTS notes_ad AFTER DELETE ON notes BEGIN
			DELETE FROM notes_fts WHERE rowid = old.id;
		END;`,
		// 更新：先删后插
		`CREATE TRIGGER IF NOT EXISTS notes_au AFTER UPDATE ON notes BEGIN
			DELETE FROM notes_fts WHERE rowid = old.id;
			INSERT INTO notes_fts(rowid, content, note_id) VALUES (new.id, COALESCE(new.content, ''), new.id);
		END;`,
	}

	// 创建标签表
	tagsTable := `
	CREATE TABLE IF NOT EXISTS tags (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		name TEXT NOT NULL,
		color TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
	);`

	// 创建笔记标签关联表
	noteTagsTable := `
	CREATE TABLE IF NOT EXISTS note_tags (
		note_id INTEGER NOT NULL,
		tag_id INTEGER NOT NULL,
		PRIMARY KEY (note_id, tag_id),
		FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE,
		FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
	);`

	// 创建用户表
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		email TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := conn.ExecContext(ctx, notesTable); err != nil {
		return err
	}

	// 创建 FTS5 虚表（如果编译未启用 FTS5，这里会报错）
	if _, err := conn.ExecContext(ctx, notesFTSTable); err != nil {
		return err
	}
	// 触发器可能会更新实现：启动时尽力 drop 再重建
	// 注意：在部分环境里 schema 前缀 drop 可能失效，因此通过 sqlite_master/sqlite_temp_master 查找实际存在的触发器名并 drop。
	dropExisting := func(masterTable string) {
		rows, e := conn.QueryContext(ctx,
			fmt.Sprintf(`SELECT name FROM %s WHERE type='trigger' AND name IN ('notes_ai','notes_ad','notes_au')`, masterTable),
		)
		if e != nil {
			return
		}
		defer rows.Close()
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				continue
			}
			// 安全转义双引号
			escaped := `"` + string([]rune(name)) + `"`
			_, _ = conn.ExecContext(ctx, fmt.Sprintf(`DROP TRIGGER IF EXISTS %s;`, escaped))
		}
	}
	dropExisting("sqlite_master")
	dropExisting("sqlite_temp_master")

	// 兜底：再尝试不带 schema 的 drop（避免一些奇怪的命名空间问题）
	_, _ = conn.ExecContext(ctx, `DROP TRIGGER IF EXISTS notes_ai;`)
	_, _ = conn.ExecContext(ctx, `DROP TRIGGER IF EXISTS notes_ad;`)
	_, _ = conn.ExecContext(ctx, `DROP TRIGGER IF EXISTS notes_au;`)

	for _, trg := range notesFTSTriggers {
		if _, err := conn.ExecContext(ctx, trg); err != nil {
			return err
		}
	}

	if _, err := conn.ExecContext(ctx, tagsTable); err != nil {
		return err
	}

	if _, err := conn.ExecContext(ctx, noteTagsTable); err != nil {
		return err
	}

	if _, err := conn.ExecContext(ctx, usersTable); err != nil {
		return err
	}

	// 首次创建/升级后，尽量把历史 notes 同步进 FTS 表（避免空索引）
	// - 使用 NOT EXISTS 防止重复插入
	// - COALESCE 避免 NULL
	_, _ = conn.ExecContext(ctx, `
		INSERT INTO notes_fts(rowid, content, note_id)
		SELECT n.id, COALESCE(n.content, ''), n.id
		FROM notes n
		WHERE NOT EXISTS (SELECT 1 FROM notes_fts f WHERE f.rowid = n.id);
	`)

	return nil
}

func ensureNotesColumnsV2(ctx context.Context, conn *sql.Conn) error {
	// pinned
	if ok, err := columnExists(ctx, conn, "notes", "pinned"); err != nil {
		return err
	} else if !ok {
		if _, err := conn.ExecContext(ctx, `ALTER TABLE notes ADD COLUMN pinned INTEGER NOT NULL DEFAULT 0;`); err != nil {
			return err
		}
	}

	// content_type
	if ok, err := columnExists(ctx, conn, "notes", "content_type"); err != nil {
		return err
	} else if !ok {
		// default: markdown（方便前端做渲染策略）
		if _, err := conn.ExecContext(ctx, `ALTER TABLE notes ADD COLUMN content_type TEXT NOT NULL DEFAULT 'markdown';`); err != nil {
			return err
		}
	}

	// user_id（可为空，兼容旧数据）
	if ok, err := columnExists(ctx, conn, "notes", "user_id"); err != nil {
		return err
	} else if !ok {
		if _, err := conn.ExecContext(ctx, `ALTER TABLE notes ADD COLUMN user_id INTEGER;`); err != nil {
			return err
		}
	}

	return nil
}

func ensureResourcesSchemaV3(ctx context.Context, conn *sql.Conn) error {
	resourcesTable := `
	CREATE TABLE IF NOT EXISTS resources (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		filename TEXT NOT NULL,
		storage_path TEXT NOT NULL,
		mime_type TEXT,
		size INTEGER,
		sha256 TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
	);`

	noteResourcesTable := `
	CREATE TABLE IF NOT EXISTS note_resources (
		note_id INTEGER NOT NULL,
		resource_id INTEGER NOT NULL,
		PRIMARY KEY (note_id, resource_id),
		FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE,
		FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE
	);`

	if _, err := conn.ExecContext(ctx, resourcesTable); err != nil {
		return err
	}
	if _, err := conn.ExecContext(ctx, noteResourcesTable); err != nil {
		return err
	}
	return nil
}

func ensureUsersAdminV4(ctx context.Context, conn *sql.Conn) error {
	// is_admin 列
	if ok, err := columnExists(ctx, conn, "users", "is_admin"); err != nil {
		return err
	} else if !ok {
		if _, err := conn.ExecContext(ctx, `ALTER TABLE users ADD COLUMN is_admin INTEGER NOT NULL DEFAULT 0;`); err != nil {
			return err
		}
	}

	// v4 只做 schema，不再在这里写入默认管理员（避免固定密码）
	return nil
}

func ensureUsersBootstrapV5(ctx context.Context, conn *sql.Conn) error {
	// must_change_password 列
	if ok, err := columnExists(ctx, conn, "users", "must_change_password"); err != nil {
		return err
	} else if !ok {
		if _, err := conn.ExecContext(ctx, `ALTER TABLE users ADD COLUMN must_change_password INTEGER NOT NULL DEFAULT 0;`); err != nil {
			return err
		}
	}

	// 规则：
	// 1) 如果设置了 MEMO_ADMIN_PASSWORD：确保 admin 存在并重置为该密码，同时强制 must_change_password=1
	// 2) 如果没有任何用户：创建 admin，并随机生成密码（打印到日志），强制 must_change_password=1
	// 3) 如果已存在 admin 且其密码仍等于旧默认 admin123：记录警告并强制 must_change_password=1

	adminPassword := strings.TrimSpace(os.Getenv("MEMO_ADMIN_PASSWORD"))

	// 是否有用户
	var userCount int
	if err := conn.QueryRowContext(ctx, `SELECT COUNT(1) FROM users`).Scan(&userCount); err != nil {
		return err
	}

	// 查 admin
	var adminID int
	var adminHash string
	err := conn.QueryRowContext(ctx, `SELECT id, password FROM users WHERE username = 'admin' LIMIT 1`).Scan(&adminID, &adminHash)
	adminExists := err == nil
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	setAdminPassword := func(pw string) (string, error) {
		hashed, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
		if err != nil {
			return "", err
		}
		return string(hashed), nil
	}

	// 1) env 指定：重置/创建 admin
	if adminPassword != "" {
		hashed, err := setAdminPassword(adminPassword)
		if err != nil {
			return err
		}
		if adminExists {
			_, err = conn.ExecContext(ctx,
				`UPDATE users SET password = ?, is_admin = 1, must_change_password = 1 WHERE id = ?`,
				hashed, adminID,
			)
			return err
		}
		_, err = conn.ExecContext(ctx,
			`INSERT INTO users (username, password, email, created_at, is_admin, must_change_password) VALUES (?, ?, ?, ?, 1, 1)`,
			"admin", hashed, "", time.Now().Format("2006-01-02 15:04:05"),
		)
		return err
	}

	// 2) 没有任何用户：创建 admin + 随机密码
	if userCount == 0 && !adminExists {
		pw := randomPassword(16)
		hashed, err := setAdminPassword(pw)
		if err != nil {
			return err
		}
		if _, err := conn.ExecContext(ctx,
			`INSERT INTO users (username, password, email, created_at, is_admin, must_change_password) VALUES (?, ?, ?, ?, 1, 1)`,
			"admin", hashed, "", time.Now().Format("2006-01-02 15:04:05"),
		); err != nil {
			return err
		}
		log.Printf("[BOOTSTRAP] 已创建默认管理员 admin，初始密码：%s（请登录后立即修改；或设置 MEMO_ADMIN_PASSWORD 覆盖）", pw)
		return nil
	}

	// 3) 仍为旧默认：admin123
	if adminExists {
		if bcrypt.CompareHashAndPassword([]byte(adminHash), []byte("admin123")) == nil {
			log.Printf("[SECURITY] 检测到 admin 仍使用旧默认密码 admin123，建议设置 MEMO_ADMIN_PASSWORD 并登录后修改密码")
			_, _ = conn.ExecContext(ctx, `UPDATE users SET must_change_password = 1 WHERE id = ?`, adminID)
		}
	}

	return nil
}

func randomPassword(n int) string {
	const alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz23456789!@#$%^&*"
	if n <= 0 {
		n = 16
	}
	b := make([]byte, n)
	
	// 使用 crypto/rand 生成安全随机密码
	if _, err := rand.Read(b); err != nil {
		// crypto/rand 失败时使用备用方案（不应该发生）
		log.Printf("[SECURITY] crypto/rand 不可用，使用备用随机源")
		for i := range b {
			b[i] = alphabet[int(time.Now().UnixNano())%len(alphabet)]
		}
	} else {
		for i := range b {
			b[i] = alphabet[int(b[i])%len(alphabet)]
		}
	}
	return string(b)
}

func ensureMultiUserIsolationV6(ctx context.Context, conn *sql.Conn) error {
	// tags.user_id
	if ok, err := columnExists(ctx, conn, "tags", "user_id"); err != nil {
		return err
	} else if !ok {
		if _, err := conn.ExecContext(ctx, `ALTER TABLE tags ADD COLUMN user_id INTEGER;`); err != nil {
			return err
		}
	}
	// per-user unique（SQLite 允许多个 NULL，因此 public tags（user_id NULL）仍可存在）
	_, _ = conn.ExecContext(ctx, `CREATE UNIQUE INDEX IF NOT EXISTS idx_tags_user_name ON tags(user_id, name);`)

	// 找到一个“主用户”用于迁移旧数据（优先管理员）
	var primaryUserID int
	err := conn.QueryRowContext(ctx, `SELECT id FROM users WHERE is_admin = 1 ORDER BY id ASC LIMIT 1`).Scan(&primaryUserID)
	if err == sql.ErrNoRows {
		err = conn.QueryRowContext(ctx, `SELECT id FROM users ORDER BY id ASC LIMIT 1`).Scan(&primaryUserID)
	}
	if err != nil {
		// 没有用户：跳过迁移（后续注册后会有新数据）
		return nil
	}

	// 把历史 notes/tags 的 NULL user_id 迁移到主用户（避免“所有人共享数据”）
	_, _ = conn.ExecContext(ctx, `UPDATE notes SET user_id = ? WHERE user_id IS NULL;`, primaryUserID)
	_, _ = conn.ExecContext(ctx, `UPDATE tags SET user_id = ? WHERE user_id IS NULL;`, primaryUserID)
	return nil
}

// SQLite 无法直接 DROP 列级 UNIQUE 约束，因此通过重建表移除 tags.name 的全局 UNIQUE。
func ensureTagsUniquePerUserV7(ctx context.Context, conn *sql.Conn) error {
	// 如果旧库存在 tags.name UNIQUE，会生成 sqlite_autoindex_tags_1；我们统一重建 tags 表。
	_, _ = conn.ExecContext(ctx, `PRAGMA foreign_keys = OFF;`)
	if _, err := conn.ExecContext(ctx, `BEGIN;`); err != nil {
		return err
	}

	// 新表：不含 name UNIQUE，只保留 (user_id,name) 唯一索引
	if _, err := conn.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS tags_new (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			name TEXT NOT NULL,
			color TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
		);
	`); err != nil {
		_, _ = conn.ExecContext(ctx, `ROLLBACK;`)
		return err
	}

	// 复制数据（保留 id，保证 note_tags 引用不变）
	if _, err := conn.ExecContext(ctx, `
		INSERT INTO tags_new (id, user_id, name, color, created_at)
		SELECT id, user_id, name, color, created_at FROM tags;
	`); err != nil {
		_, _ = conn.ExecContext(ctx, `ROLLBACK;`)
		return err
	}

	// 替换表
	if _, err := conn.ExecContext(ctx, `DROP TABLE tags;`); err != nil {
		_, _ = conn.ExecContext(ctx, `ROLLBACK;`)
		return err
	}
	if _, err := conn.ExecContext(ctx, `ALTER TABLE tags_new RENAME TO tags;`); err != nil {
		_, _ = conn.ExecContext(ctx, `ROLLBACK;`)
		return err
	}

	// 重建 per-user unique index
	if _, err := conn.ExecContext(ctx, `CREATE UNIQUE INDEX IF NOT EXISTS idx_tags_user_name ON tags(user_id, name);`); err != nil {
		_, _ = conn.ExecContext(ctx, `ROLLBACK;`)
		return err
	}

	if _, err := conn.ExecContext(ctx, `COMMIT;`); err != nil {
		_, _ = conn.ExecContext(ctx, `ROLLBACK;`)
		return err
	}
	_, _ = conn.ExecContext(ctx, `PRAGMA foreign_keys = ON;`)
	return nil
}

func columnExists(ctx context.Context, conn *sql.Conn, table, column string) (bool, error) {
	// PRAGMA table_info 不支持占位符绑定 table 名，只能拼接；这里做最小转义校验
	if strings.TrimSpace(table) == "" || strings.TrimSpace(column) == "" {
		return false, fmt.Errorf("invalid table/column")
	}
	table = strings.ReplaceAll(table, `"`, `""`)
	rows, err := conn.QueryContext(ctx, fmt.Sprintf(`PRAGMA table_info("%s");`, table))
	if err != nil {
		return false, err
	}
	defer rows.Close()

	// table_info 输出：cid,name,type,notnull,dflt_value,pk
	for rows.Next() {
		var cid int
		var name, typ string
		var notnull int
		var dflt sql.NullString
		var pk int
		if err := rows.Scan(&cid, &name, &typ, &notnull, &dflt, &pk); err != nil {
			return false, err
		}
		if name == column {
			return true, nil
		}
	}
	return false, rows.Err()
}
