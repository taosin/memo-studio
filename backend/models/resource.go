package models

import (
	"database/sql"
	"memo-studio/backend/database"
	"strings"
	"time"
)

type Resource struct {
	ID          int       `json:"id"`
	UserID      *int      `json:"user_id,omitempty"`
	Filename    string    `json:"filename"`
	StoragePath string    `json:"storage_path"`
	URL         string    `json:"url"`
	MimeType    string    `json:"mime_type"`
	Size        int64     `json:"size"`
	Sha256      string    `json:"sha256"`
	CreatedAt   time.Time `json:"created_at"`
}

func normalizeStoragePath(p string) string {
	p = strings.TrimSpace(p)
	p = strings.TrimPrefix(p, "/")
	return p
}

func resourceURL(storagePath string) string {
	sp := normalizeStoragePath(storagePath)
	if sp == "" {
		return ""
	}
	return "/uploads/" + sp
}

func CreateResource(userID *int, filename, storagePath, mimeType string, size int64, sha256 string) (*Resource, error) {
	var userParam interface{} = nil
	if userID != nil {
		userParam = *userID
	}
	storagePath = normalizeStoragePath(storagePath)

	res, err := database.DB.Exec(
		`INSERT INTO resources (user_id, filename, storage_path, mime_type, size, sha256)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		userParam, filename, storagePath, mimeType, size, sha256,
	)
	if err != nil {
		return nil, err
	}
	id64, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return GetResource(int(id64))
}

func GetResource(id int) (*Resource, error) {
	var r Resource
	var user sql.NullInt64
	err := database.DB.QueryRow(
		`SELECT id, user_id, filename, storage_path, mime_type, size, sha256, created_at
		 FROM resources WHERE id = ?`,
		id,
	).Scan(&r.ID, &user, &r.Filename, &r.StoragePath, &r.MimeType, &r.Size, &r.Sha256, &r.CreatedAt)
	if err != nil {
		return nil, err
	}
	if user.Valid {
		v := int(user.Int64)
		r.UserID = &v
	}
	r.StoragePath = normalizeStoragePath(r.StoragePath)
	r.URL = resourceURL(r.StoragePath)
	return &r, nil
}

func GetResourcesByNoteID(noteID int) ([]Resource, error) {
	rows, err := database.DB.Query(
		`SELECT r.id, r.user_id, r.filename, r.storage_path, r.mime_type, r.size, r.sha256, r.created_at
		 FROM note_resources nr
		 JOIN resources r ON r.id = nr.resource_id
		 WHERE nr.note_id = ?
		 ORDER BY r.created_at ASC, r.id ASC`,
		noteID,
	)
	if err != nil {
		// 兼容：如果旧库没创建 resources 表，这里会报错；上层会返回 500。
		return nil, err
	}
	defer rows.Close()

	var list []Resource
	for rows.Next() {
		var r Resource
		var user sql.NullInt64
		if err := rows.Scan(&r.ID, &user, &r.Filename, &r.StoragePath, &r.MimeType, &r.Size, &r.Sha256, &r.CreatedAt); err != nil {
			return nil, err
		}
		if user.Valid {
			v := int(user.Int64)
			r.UserID = &v
		}
		r.StoragePath = normalizeStoragePath(r.StoragePath)
		r.URL = resourceURL(r.StoragePath)
		list = append(list, r)
	}
	return list, rows.Err()
}

// ListResourcesResult 资源列表项（含总数）
type ListResourcesResult struct {
	Items []Resource `json:"items"`
	Total int        `json:"total"`
}

// ListResourcesByUserID 分页列出当前用户的资源
func ListResourcesByUserID(userID int, limit, offset int) (*ListResourcesResult, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	var total int
	err := database.DB.QueryRow(
		`SELECT COUNT(*) FROM resources WHERE user_id = ?`,
		userID,
	).Scan(&total)
	if err != nil {
		return nil, err
	}

	rows, err := database.DB.Query(
		`SELECT id, user_id, filename, storage_path, mime_type, size, sha256, created_at
		 FROM resources WHERE user_id = ?
		 ORDER BY created_at DESC, id DESC
		 LIMIT ? OFFSET ?`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Resource
	for rows.Next() {
		var r Resource
		var user sql.NullInt64
		if err := rows.Scan(&r.ID, &user, &r.Filename, &r.StoragePath, &r.MimeType, &r.Size, &r.Sha256, &r.CreatedAt); err != nil {
			return nil, err
		}
		if user.Valid {
			v := int(user.Int64)
			r.UserID = &v
		}
		r.StoragePath = normalizeStoragePath(r.StoragePath)
		r.URL = resourceURL(r.StoragePath)
		list = append(list, r)
	}
	if list == nil {
		list = []Resource{}
	}
	return &ListResourcesResult{Items: list, Total: total}, rows.Err()
}

// DeleteResource 删除资源（仅删除数据库记录；物理文件可由定时任务清理）
func DeleteResource(id int, userID int) error {
	res, err := database.DB.Exec(
		`DELETE FROM resources WHERE id = ? AND user_id = ?`,
		id, userID,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	_, _ = database.DB.Exec(`DELETE FROM note_resources WHERE resource_id = ?`, id)
	return nil
}

