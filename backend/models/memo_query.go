package models

import (
	"database/sql"
	"fmt"
	"memo-studio/backend/database"
	"strconv"
	"strings"
	"time"
)

type MemoQuery struct {
	Limit       int
	Offset      int
	Q           string
	Tags        []string
	From        *time.Time
	To          *time.Time
	Pinned      *bool
	ContentType string
	UserID      *int
}

func ListMemos(q MemoQuery) ([]Note, error) {
	limit := q.Limit
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	offset := q.Offset
	if offset < 0 {
		offset = 0
	}

	args := []interface{}{}

	// 基础 FROM
	from := " FROM notes n "
	where := " WHERE 1=1 "

	// FTS
	fts := strings.TrimSpace(q.Q)
	if fts != "" {
		from = " FROM notes_fts f JOIN notes n ON n.id = f.rowid "
		where += " AND notes_fts MATCH ? "
		args = append(args, fts)
	}

	// user_id
	if q.UserID != nil && *q.UserID > 0 {
		// 兼容旧数据：user_id 为空的历史 notes 也能被看到
		where += " AND (n.user_id = ? OR n.user_id IS NULL) "
		args = append(args, *q.UserID)
	}

	// pinned
	if q.Pinned != nil {
		if *q.Pinned {
			where += " AND n.pinned = 1 "
		} else {
			where += " AND n.pinned = 0 "
		}
	}

	// content_type
	if strings.TrimSpace(q.ContentType) != "" {
		where += " AND n.content_type = ? "
		args = append(args, strings.TrimSpace(q.ContentType))
	}

	// date range
	if q.From != nil {
		where += " AND n.created_at >= ? "
		args = append(args, q.From.Format(time.RFC3339))
	}
	if q.To != nil {
		where += " AND n.created_at <= ? "
		args = append(args, q.To.Format(time.RFC3339))
	}

	// tags（任意匹配）
	if len(q.Tags) > 0 {
		placeholders := make([]string, 0, len(q.Tags))
		for range q.Tags {
			placeholders = append(placeholders, "?")
		}
		from += " JOIN note_tags nt ON nt.note_id = n.id JOIN tags t ON t.id = nt.tag_id "
		where += " AND t.name IN (" + strings.Join(placeholders, ",") + ") "
		for _, t := range q.Tags {
			args = append(args, t)
		}
	}

	order := " ORDER BY n.pinned DESC, n.created_at DESC, n.id DESC "
	limitOffset := " LIMIT ? OFFSET ? "
	args = append(args, limit, offset)

	// SELECT（tags JOIN 时会产生重复行，所以用 DISTINCT）
	selectPrefix := "SELECT n.id, n.user_id, n.title, n.content, n.pinned, n.content_type, n.created_at, n.updated_at"
	if len(q.Tags) > 0 {
		selectPrefix = "SELECT DISTINCT n.id, n.user_id, n.title, n.content, n.pinned, n.content_type, n.created_at, n.updated_at"
	}
	sqlStr := selectPrefix + from + where
	if fts != "" {
		// 有搜索时按 bm25 排序更合理，但 pinned 仍优先
		order = " ORDER BY n.pinned DESC, bm25(notes_fts), n.created_at DESC "
	}
	sqlStr += order + limitOffset

	rows, err := database.DB.Query(sqlStr, args...)
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

		notebookIDs, _ := GetNotebookIDsByNoteID(note.ID)
		note.NotebookIDs = notebookIDs

		notes = append(notes, note)
	}
	return notes, rows.Err()
}

func ParseTagsParam(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parts := strings.FieldsFunc(s, func(r rune) bool { return r == ',' || r == ' ' || r == '，' || r == ';' })
	out := make([]string, 0, len(parts))
	seen := map[string]bool{}
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if len(p) > 64 {
			p = p[:64]
		}
		if !seen[p] {
			seen[p] = true
			out = append(out, p)
		}
	}
	return out
}

func ParseBoolParam(s string) (*bool, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return nil, nil
	}
	switch s {
	case "1", "true", "yes", "y":
		v := true
		return &v, nil
	case "0", "false", "no", "n":
		v := false
		return &v, nil
	default:
		return nil, fmt.Errorf("invalid bool: %s", s)
	}
}

func ParseLimitOffset(limitStr, offsetStr string) (limit int, offset int) {
	limit = 50
	offset = 0
	if strings.TrimSpace(limitStr) != "" {
		if v, err := strconv.Atoi(limitStr); err == nil {
			limit = v
		}
	}
	if strings.TrimSpace(offsetStr) != "" {
		if v, err := strconv.Atoi(offsetStr); err == nil {
			offset = v
		}
	}
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return
}

