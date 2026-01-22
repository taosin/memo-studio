package models

import (
	"database/sql"
	"memo-studio/backend/database"
	"time"
)

type Note struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []Tag     `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Tag struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateNote 创建笔记
func CreateNote(title, content string, tagIDs []int) (*Note, error) {
	tx, err := database.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result, err := tx.Exec(
		"INSERT INTO notes (title, content) VALUES (?, ?)",
		title, content,
	)
	if err != nil {
		return nil, err
	}

	noteID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// 关联标签
	for _, tagID := range tagIDs {
		_, err = tx.Exec(
			"INSERT INTO note_tags (note_id, tag_id) VALUES (?, ?)",
			noteID, tagID,
		)
		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return GetNote(int(noteID))
}

// GetNote 获取单个笔记
func GetNote(id int) (*Note, error) {
	note := &Note{}
	err := database.DB.QueryRow(
		"SELECT id, title, content, created_at, updated_at FROM notes WHERE id = ?",
		id,
	).Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt)

	if err != nil {
		return nil, err
	}

	// 获取标签
	tags, err := GetTagsByNoteID(id)
	if err != nil {
		return nil, err
	}
	note.Tags = tags

	return note, nil
}

// GetAllNotes 获取所有笔记
func GetAllNotes() ([]Note, error) {
	rows, err := database.DB.Query(
		"SELECT id, title, content, created_at, updated_at FROM notes ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt)
		if err != nil {
			return nil, err
		}

		// 获取标签
		tags, err := GetTagsByNoteID(note.ID)
		if err != nil {
			return nil, err
		}
		note.Tags = tags

		notes = append(notes, note)
	}

	return notes, nil
}

// GetTagsByNoteID 获取笔记的标签
func GetTagsByNoteID(noteID int) ([]Tag, error) {
	rows, err := database.DB.Query(
		`SELECT t.id, t.name, t.color, t.created_at 
		 FROM tags t 
		 INNER JOIN note_tags nt ON t.id = nt.tag_id 
		 WHERE nt.note_id = ?`,
		noteID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var tag Tag
		err := rows.Scan(&tag.ID, &tag.Name, &tag.Color, &tag.CreatedAt)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

// GetAllTags 获取所有标签
func GetAllTags() ([]Tag, error) {
	rows, err := database.DB.Query(
		"SELECT id, name, color, created_at FROM tags ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var tag Tag
		err := rows.Scan(&tag.ID, &tag.Name, &tag.Color, &tag.CreatedAt)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

// CreateTagIfNotExists 如果标签不存在则创建
func CreateTagIfNotExists(name string) (*Tag, error) {
	// 先查找是否存在
	var tag Tag
	err := database.DB.QueryRow(
		"SELECT id, name, color, created_at FROM tags WHERE name = ?",
		name,
	).Scan(&tag.ID, &tag.Name, &tag.Color, &tag.CreatedAt)

	if err == sql.ErrNoRows {
		// 标签不存在，创建新标签
		result, err := database.DB.Exec(
			"INSERT INTO tags (name, color) VALUES (?, ?)",
			name, getTagColor(name),
		)
		if err != nil {
			return nil, err
		}

		tagID, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		return &Tag{
			ID:    int(tagID),
			Name:  name,
			Color: getTagColor(name),
		}, nil
	} else if err != nil {
		return nil, err
	}

	return &tag, nil
}

// getTagColor 根据标签名生成颜色
func getTagColor(name string) string {
	colors := []string{
		"#FF6B6B", "#4ECDC4", "#45B7D1", "#FFA07A",
		"#98D8C8", "#F7DC6F", "#BB8FCE", "#85C1E2",
	}
	hash := 0
	for _, char := range name {
		hash = int(char) + ((hash << 5) - hash)
	}
	return colors[abs(hash)%len(colors)]
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
