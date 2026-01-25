package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
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
		name TEXT NOT NULL UNIQUE,
		color TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
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

	// 默认管理员：admin/admin123
	// - 若已有 admin 用户则不覆盖
	var cnt int
	if err := conn.QueryRowContext(ctx, `SELECT COUNT(1) FROM users WHERE username = 'admin'`).Scan(&cnt); err != nil {
		return err
	}
	if cnt > 0 {
		return nil
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = conn.ExecContext(ctx,
		`INSERT INTO users (username, password, email, created_at, is_admin) VALUES (?, ?, ?, ?, 1)`,
		"admin", string(hashed), "", time.Now().Format("2006-01-02 15:04:05"),
	)
	return err
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
