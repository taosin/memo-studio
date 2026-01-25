package models

import (
	"database/sql"
	"log"
	"memo-studio/backend/database"
	"strconv"
	"strings"
	"time"
)

type Note struct {
	ID          int        `json:"id"`
	UserID      *int       `json:"user_id,omitempty"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	ContentType string     `json:"content_type"`
	Pinned      bool       `json:"pinned"`
	Tags        []Tag      `json:"tags"`
	Resources   []Resource `json:"resources"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Tag struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
}

type TagWithCount struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
	NoteCount int       `json:"note_count"`
}

// CreateNote 创建笔记
func CreateNote(title, content string, tagIDs []int, pinned bool, contentType string, resourceIDs []int, userID *int) (*Note, error) {
	if strings.TrimSpace(contentType) == "" {
		contentType = "markdown"
	}
	tx, err := database.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var userParam interface{} = nil
	if userID != nil {
		userParam = *userID
	}

	result, err := tx.Exec(
		"INSERT INTO notes (title, content, pinned, content_type, user_id) VALUES (?, ?, ?, ?, ?)",
		title, content, pinned, contentType, userParam,
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

	// 关联附件
	for _, rid := range resourceIDs {
		if rid <= 0 {
			continue
		}
		_, err = tx.Exec(
			"INSERT INTO note_resources (note_id, resource_id) VALUES (?, ?)",
			noteID, rid,
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
func UpdateNote(id int, title, content string, tagIDs []int, pinned bool, contentType string, resourceIDs []int) (*Note, error) {
	if strings.TrimSpace(contentType) == "" {
		contentType = "markdown"
	}
	tx, err := database.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// 更新笔记
	_, err = tx.Exec(
		"UPDATE notes SET title = ?, content = ?, pinned = ?, content_type = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		title, content, pinned, contentType, id,
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

	// 删除旧的附件关联
	_, err = tx.Exec("DELETE FROM note_resources WHERE note_id = ?", id)
	if err != nil {
		return nil, err
	}
	// 添加新的附件关联
	for _, rid := range resourceIDs {
		if rid <= 0 {
			continue
		}
		_, err = tx.Exec(
			"INSERT INTO note_resources (note_id, resource_id) VALUES (?, ?)",
			id, rid,
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
	var userID sql.NullInt64
	var pinnedInt int
	var contentType string
	err := database.DB.QueryRow(
		"SELECT id, user_id, title, content, pinned, content_type, created_at, updated_at FROM notes WHERE id = ?",
		id,
	).Scan(&note.ID, &userID, &note.Title, &note.Content, &pinnedInt, &contentType, &note.CreatedAt, &note.UpdatedAt)

	if err != nil {
		return nil, err
	}

	if userID.Valid {
		v := int(userID.Int64)
		note.UserID = &v
	}
	note.Pinned = pinnedInt != 0
	note.ContentType = contentType

	// 清理 content 字段（只清理 [object Object] 字符串）
	note.Content = cleanContent(note.Content)
	note.Title = cleanContent(note.Title)
	
	// 调试日志
	log.Printf("[GetNote] 读取 - ID: %d, Title: %q (len: %d), Content: %q (len: %d)", 
		note.ID, note.Title, len(note.Title), note.Content, len(note.Content))

	// 获取标签
	tags, err := GetTagsByNoteID(id)
	if err != nil {
		return nil, err
	}
	note.Tags = tags

	// 获取附件
	resources, err := GetResourcesByNoteID(id)
	if err != nil {
		return nil, err
	}
	note.Resources = resources

	return note, nil
}

// GetAllNotes 获取所有笔记
func GetAllNotes() ([]Note, error) {
	rows, err := database.DB.Query(
		"SELECT id, user_id, title, content, pinned, content_type, created_at, updated_at FROM notes ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		var userID sql.NullInt64
		var pinnedInt int
		var contentType string
		err := rows.Scan(&note.ID, &userID, &note.Title, &note.Content, &pinnedInt, &contentType, &note.CreatedAt, &note.UpdatedAt)
		if err != nil {
			return nil, err
		}

		if userID.Valid {
			v := int(userID.Int64)
			note.UserID = &v
		}
		note.Pinned = pinnedInt != 0
		note.ContentType = contentType

		// 清理 content 和 title 字段（只清理 [object Object] 字符串）
		note.Content = cleanContent(note.Content)
		note.Title = cleanContent(note.Title)
		
		// 调试日志（只记录前几条）
		if len(notes) < 3 {
			log.Printf("[GetAllNotes] 读取 - ID: %d, Title: %q (len: %d), Content: %q (len: %d)", 
				note.ID, note.Title, len(note.Title), note.Content, len(note.Content))
		}

		// 获取标签
		tags, err := GetTagsByNoteID(note.ID)
		if err != nil {
			return nil, err
		}
		note.Tags = tags

		// 获取附件（列表场景：N+1，当前规模可接受；后续可做聚合优化）
		resources, err := GetResourcesByNoteID(note.ID)
		if err != nil {
			return nil, err
		}
		note.Resources = resources

		notes = append(notes, note)
	}

	return notes, nil
}

// SearchNotes FTS5 全文搜索（按内容）
func SearchNotes(q string, limit, offset int) ([]Note, error) {
	q = strings.TrimSpace(q)
	if q == "" {
		// 空查询：退化为最新列表
		return GetAllNotes()
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	rows, err := database.DB.Query(
		`SELECT n.id, n.user_id, n.title, n.content, n.pinned, n.content_type, n.created_at, n.updated_at
		 FROM notes_fts f
		 JOIN notes n ON n.id = f.rowid
		 WHERE notes_fts MATCH ?
		 ORDER BY bm25(notes_fts)
		 LIMIT ? OFFSET ?`,
		q, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		var userID sql.NullInt64
		var pinnedInt int
		var contentType string
		if err := rows.Scan(&note.ID, &userID, &note.Title, &note.Content, &pinnedInt, &contentType, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, err
		}

		if userID.Valid {
			v := int(userID.Int64)
			note.UserID = &v
		}
		note.Pinned = pinnedInt != 0
		note.ContentType = contentType

		note.Content = cleanContent(note.Content)
		note.Title = cleanContent(note.Title)

		tags, err := GetTagsByNoteID(note.ID)
		if err != nil {
			return nil, err
		}
		note.Tags = tags

		resources, err := GetResourcesByNoteID(note.ID)
		if err != nil {
			return nil, err
		}
		note.Resources = resources
		notes = append(notes, note)
	}

	return notes, nil
}

// GetTagsWithCount 获取标签列表（包含笔记计数）
func GetTagsWithCount() ([]TagWithCount, error) {
	rows, err := database.DB.Query(
		`SELECT t.id, t.name, t.color, t.created_at, COUNT(nt.note_id) AS note_count
		 FROM tags t
		 LEFT JOIN note_tags nt ON t.id = nt.tag_id
		 GROUP BY t.id
		 ORDER BY note_count DESC, t.created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []TagWithCount
	for rows.Next() {
		var t TagWithCount
		if err := rows.Scan(&t.ID, &t.Name, &t.Color, &t.CreatedAt, &t.NoteCount); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, nil
}

// GetTagByName 根据名称获取标签
func GetTagByName(name string) (*Tag, error) {
	var tag Tag
	err := database.DB.QueryRow(
		"SELECT id, name, color, created_at FROM tags WHERE name = ?",
		name,
	).Scan(&tag.ID, &tag.Name, &tag.Color, &tag.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// RandomNotes 随机回顾笔记（可按标签过滤，可按天数过滤）
func RandomNotes(limit int, tagName string, withinDays int) ([]Note, error) {
	if limit <= 0 || limit > 20 {
		limit = 1
	}

	args := []interface{}{}
	query := `
		SELECT n.id, n.user_id, n.title, n.content, n.pinned, n.content_type, n.created_at, n.updated_at
		FROM notes n
	`

	where := " WHERE 1=1 "

	if strings.TrimSpace(tagName) != "" {
		// tagName -> tag_id
		tag, err := GetTagByName(strings.TrimSpace(tagName))
		if err != nil {
			// 没有该标签：直接返回空
			if err == sql.ErrNoRows {
				return []Note{}, nil
			}
			return nil, err
		}
		query += " JOIN note_tags nt ON nt.note_id = n.id "
		where += " AND nt.tag_id = ? "
		args = append(args, tag.ID)
	}

	if withinDays > 0 {
		where += " AND n.created_at >= datetime('now', ?) "
		args = append(args, "-"+strconv.Itoa(withinDays)+" days")
	}

	query += where + " ORDER BY RANDOM() LIMIT ? "
	args = append(args, limit)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		var userID sql.NullInt64
		var pinnedInt int
		var contentType string
		if err := rows.Scan(&note.ID, &userID, &note.Title, &note.Content, &pinnedInt, &contentType, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, err
		}
		if userID.Valid {
			v := int(userID.Int64)
			note.UserID = &v
		}
		note.Pinned = pinnedInt != 0
		note.ContentType = contentType
		note.Content = cleanContent(note.Content)
		note.Title = cleanContent(note.Title)
		tags, err := GetTagsByNoteID(note.ID)
		if err != nil {
			return nil, err
		}
		note.Tags = tags

		resources, err := GetResourcesByNoteID(note.ID)
		if err != nil {
			return nil, err
		}
		note.Resources = resources
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
