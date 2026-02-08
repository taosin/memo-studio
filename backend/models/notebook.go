package models

import (
	"database/sql"
	"memo-studio/backend/database"
	"strings"
	"time"
)

type Notebook struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	NoteCount int       `json:"note_count,omitempty"`
}

func ListNotebooks(userID int) ([]Notebook, error) {
	rows, err := database.DB.Query(`
		SELECT n.id, n.user_id, n.name, n.color, n.sort_order, n.created_at, n.updated_at,
		       COALESCE(cnt.c, 0) AS note_count
		FROM notebooks n
		LEFT JOIN (SELECT notebook_id, COUNT(*) AS c FROM note_notebooks GROUP BY notebook_id) cnt ON cnt.notebook_id = n.id
		WHERE n.user_id = ?
		ORDER BY n.sort_order ASC, n.id ASC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Notebook
	for rows.Next() {
		var nb Notebook
		var noteCount int
		if err := rows.Scan(&nb.ID, &nb.UserID, &nb.Name, &nb.Color, &nb.SortOrder, &nb.CreatedAt, &nb.UpdatedAt, &noteCount); err != nil {
			return nil, err
		}
		nb.NoteCount = noteCount
		list = append(list, nb)
	}
	return list, rows.Err()
}

func GetNotebook(id int, userID int) (*Notebook, error) {
	var nb Notebook
	var noteCount int
	err := database.DB.QueryRow(`
		SELECT n.id, n.user_id, n.name, n.color, n.sort_order, n.created_at, n.updated_at,
		       COALESCE((SELECT COUNT(*) FROM note_notebooks WHERE notebook_id = n.id), 0)
		FROM notebooks n
		WHERE n.id = ? AND n.user_id = ?
	`, id, userID).Scan(&nb.ID, &nb.UserID, &nb.Name, &nb.Color, &nb.SortOrder, &nb.CreatedAt, &nb.UpdatedAt, &noteCount)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	nb.NoteCount = noteCount
	return &nb, nil
}

func CreateNotebook(userID int, name, color string, sortOrder int) (*Notebook, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		name = "未命名笔记本"
	}
	result, err := database.DB.Exec(`
		INSERT INTO notebooks (user_id, name, color, sort_order, updated_at) VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
	`, userID, name, color, sortOrder)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return GetNotebook(int(id), userID)
}

func UpdateNotebook(id, userID int, name, color string, sortOrder *int) (*Notebook, error) {
	nb, err := GetNotebook(id, userID)
	if err != nil || nb == nil {
		return nil, err
	}
	if name != "" {
		nb.Name = strings.TrimSpace(name)
	}
	if color != "" {
		nb.Color = color
	}
	if sortOrder != nil {
		nb.SortOrder = *sortOrder
	}
	_, err = database.DB.Exec(`
		UPDATE notebooks SET name = ?, color = ?, sort_order = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND user_id = ?
	`, nb.Name, nb.Color, nb.SortOrder, id, userID)
	if err != nil {
		return nil, err
	}
	return GetNotebook(id, userID)
}

func DeleteNotebook(id, userID int) error {
	_, err := database.DB.Exec(`DELETE FROM notebooks WHERE id = ? AND user_id = ?`, id, userID)
	return err
}

func GetNotebookIDsByNoteID(noteID int) ([]int, error) {
	rows, err := database.DB.Query(`SELECT notebook_id FROM note_notebooks WHERE note_id = ?`, noteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func SetNoteNotebooks(noteID int, notebookIDs []int) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec(`DELETE FROM note_notebooks WHERE note_id = ?`, noteID)
	if err != nil {
		return err
	}
	for _, nid := range notebookIDs {
		if nid <= 0 {
			continue
		}
		_, err = tx.Exec(`INSERT OR IGNORE INTO note_notebooks (note_id, notebook_id) VALUES (?, ?)`, noteID, nid)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func ListNotesByNotebookID(notebookID, userID int, limit, offset int) ([]Note, error) {
	nb, err := GetNotebook(notebookID, userID)
	if err != nil || nb == nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 50
	}
	rows, err := database.DB.Query(`
		SELECT n.id, n.user_id, n.title, n.content, n.content_type, n.pinned, n.created_at, n.updated_at
		FROM notes n
		INNER JOIN note_notebooks nn ON nn.note_id = n.id AND nn.notebook_id = ?
		WHERE n.user_id = ?
		ORDER BY n.pinned DESC, n.updated_at DESC
		LIMIT ? OFFSET ?
	`, notebookID, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		var userIDNull sql.NullInt64
		if err := rows.Scan(&note.ID, &userIDNull, &note.Title, &note.Content, &note.ContentType, &note.Pinned, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, err
		}
		if userIDNull.Valid {
			u := int(userIDNull.Int64)
			note.UserID = &u
		}
		notes = append(notes, note)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for i := range notes {
		tags, _ := GetTagsByNoteID(notes[i].ID)
		notes[i].Tags = tags
	}
	return notes, nil
}

func CountNotesByNotebookID(notebookID, userID int) (int, error) {
	var c int
	err := database.DB.QueryRow(`
		SELECT COUNT(*)
		FROM notes n
		INNER JOIN note_notebooks nn ON nn.note_id = n.id AND nn.notebook_id = ?
		WHERE n.user_id = ?
	`, notebookID, userID).Scan(&c)
	return c, err
}
