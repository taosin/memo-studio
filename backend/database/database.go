package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// Init 初始化数据库连接和表结构
func Init() error {
	var err error
	DB, err = sql.Open("sqlite3", "./notes.db")
	if err != nil {
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
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

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

	if _, err := DB.Exec(notesTable); err != nil {
		return err
	}

	if _, err := DB.Exec(tagsTable); err != nil {
		return err
	}

	if _, err := DB.Exec(noteTagsTable); err != nil {
		return err
	}

	return nil
}
