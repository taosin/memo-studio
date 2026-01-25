package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
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

	// 创建表
	if err := createTables(); err != nil {
		return err
	}

	log.Println("数据库初始化成功")
	return nil
}

// createTables 创建数据库表
func createTables() error {
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
		// 删除（FTS5 推荐用 delete 命令）
		`CREATE TRIGGER IF NOT EXISTS notes_ad AFTER DELETE ON notes BEGIN
			INSERT INTO notes_fts(notes_fts, rowid, content, note_id) VALUES('delete', old.id, COALESCE(old.content, ''), old.id);
		END;`,
		// 更新：先 delete 再 insert
		`CREATE TRIGGER IF NOT EXISTS notes_au AFTER UPDATE ON notes BEGIN
			INSERT INTO notes_fts(notes_fts, rowid, content, note_id) VALUES('delete', old.id, COALESCE(old.content, ''), old.id);
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

	if _, err := DB.Exec(notesTable); err != nil {
		return err
	}

	// 创建 FTS5 虚表（如果编译未启用 FTS5，这里会报错）
	if _, err := DB.Exec(notesFTSTable); err != nil {
		return err
	}
	for _, trg := range notesFTSTriggers {
		if _, err := DB.Exec(trg); err != nil {
			return err
		}
	}

	if _, err := DB.Exec(tagsTable); err != nil {
		return err
	}

	if _, err := DB.Exec(noteTagsTable); err != nil {
		return err
	}

	if _, err := DB.Exec(usersTable); err != nil {
		return err
	}

	// 首次创建/升级后，尽量把历史 notes 同步进 FTS 表（避免空索引）
	// - 使用 NOT EXISTS 防止重复插入
	// - COALESCE 避免 NULL
	_, _ = DB.Exec(`
		INSERT INTO notes_fts(rowid, content, note_id)
		SELECT n.id, COALESCE(n.content, ''), n.id
		FROM notes n
		WHERE NOT EXISTS (SELECT 1 FROM notes_fts f WHERE f.rowid = n.id);
	`)

	return nil
}
