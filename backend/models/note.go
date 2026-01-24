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

// UpdateNote 更新笔记
func UpdateNote(id int, title, content string, tagIDs []int) (*Note, error) {
	tx, err := database.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// 更新笔记
	_, err = tx.Exec(
		"UPDATE notes SET title = ?, content = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		title, content, id,
	)
	if err != nil {
		return nil, err
	}

	// 删除旧的标签关联
	_, err = tx.Exec("DELETE FROM note_tags WHERE note_id = ?", id)
	if err != nil {
		return nil, err
	}

	// 添加新的标签关联
	for _, tagID := range tagIDs {
		_, err = tx.Exec(
			"INSERT INTO note_tags (note_id, tag_id) VALUES (?, ?)",
			id, tagID,
		)
		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return GetNote(id)
}

// DeleteNote 删除笔记
func DeleteNote(id int) error {
	_, err := database.DB.Exec("DELETE FROM notes WHERE id = ?", id)
	return err
}

// DeleteNotes 批量删除笔记
func DeleteNotes(ids []int) error {
	if len(ids) == 0 {
		return nil
	}

	// 构建占位符
	placeholders := ""
	for i := 0; i < len(ids); i++ {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "?"
	}

	query := "DELETE FROM notes WHERE id IN (" + placeholders + ")"
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	_, err := database.DB.Exec(query, args...)
	return err
}

// cleanContent 清理 content 字段，移除错误的 "[object Object]" 字符串
// 注意：只清理完全匹配的 "[object Object]"，不清理其他内容
func cleanContent(content string) string {
	// 只清理完全匹配的 "[object Object]" 字符串
	if content == "[object Object]" || content == "[object object]" {
		return ""
	}
	// 其他内容（包括空字符串、正常内容等）都原样返回
	return content
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

	// 清理 content 字段
	note.Content = cleanContent(note.Content)
	note.Title = cleanContent(note.Title)

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

		// 清理 content 和 title 字段
		note.Content = cleanContent(note.Content)
		note.Title = cleanContent(note.Title)

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

// GetTagByID 根据ID获取标签
func GetTagByID(id int) (*Tag, error) {
	var tag Tag
	err := database.DB.QueryRow(
		"SELECT id, name, color, created_at FROM tags WHERE id = ?",
		id,
	).Scan(&tag.ID, &tag.Name, &tag.Color, &tag.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &tag, nil
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

// UpdateTag 更新标签
func UpdateTag(id int, name, color string) (*Tag, error) {
	_, err := database.DB.Exec(
		"UPDATE tags SET name = ?, color = ? WHERE id = ?",
		name, color, id,
	)
	if err != nil {
		return nil, err
	}

	// 获取更新后的标签
	var tag Tag
	err = database.DB.QueryRow(
		"SELECT id, name, color, created_at FROM tags WHERE id = ?",
		id,
	).Scan(&tag.ID, &tag.Name, &tag.Color, &tag.CreatedAt)

	if err != nil {
		return nil, err
	}

	// 标签更新后，关联的笔记会自动使用新的标签信息（通过 JOIN 查询）

	return &tag, nil
}

// DeleteTag 删除标签
func DeleteTag(id int) error {
	// 先删除所有笔记标签关联
	_, err := database.DB.Exec("DELETE FROM note_tags WHERE tag_id = ?", id)
	if err != nil {
		return err
	}

	// 删除标签
	_, err = database.DB.Exec("DELETE FROM tags WHERE id = ?", id)
	return err
}

// MergeTags 合并标签
func MergeTags(sourceID, targetID int) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 获取所有使用源标签的笔记
	rows, err := tx.Query("SELECT DISTINCT note_id FROM note_tags WHERE tag_id = ?", sourceID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var noteIDs []int
	for rows.Next() {
		var noteID int
		if err := rows.Scan(&noteID); err != nil {
			return err
		}
		noteIDs = append(noteIDs, noteID)
	}

	// 对于每个笔记，如果已经有目标标签则删除源标签，否则替换
	for _, noteID := range noteIDs {
		var hasTarget bool
		err := tx.QueryRow(
			"SELECT COUNT(*) > 0 FROM note_tags WHERE note_id = ? AND tag_id = ?",
			noteID, targetID,
		).Scan(&hasTarget)

		if err != nil {
			return err
		}

		if hasTarget {
			// 已有目标标签，只删除源标签
			_, err = tx.Exec("DELETE FROM note_tags WHERE note_id = ? AND tag_id = ?", noteID, sourceID)
		} else {
			// 没有目标标签，替换源标签为目标标签
			_, err = tx.Exec(
				"UPDATE note_tags SET tag_id = ? WHERE note_id = ? AND tag_id = ?",
				targetID, noteID, sourceID,
			)
		}

		if err != nil {
			return err
		}
	}

	// 删除源标签
	_, err = tx.Exec("DELETE FROM tags WHERE id = ?", sourceID)
	if err != nil {
		return err
	}

	return tx.Commit()
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
