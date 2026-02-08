package models

import (
	"memo-studio/backend/database"
)

// UserStats 当前用户统计
type UserStats struct {
	NotesCount       int `json:"notes_count"`
	TagsCount        int `json:"tags_count"`
	ResourcesCount   int `json:"resources_count"`
	NotebooksCount   int `json:"notebooks_count"`
	PinnedCount      int `json:"pinned_count"`
	NotesCreated7d   int `json:"notes_created_7d"`
	NotesUpdated7d  int `json:"notes_updated_7d"`
}

// GetUserStats 获取指定用户的统计信息
func GetUserStats(userID int) (*UserStats, error) {
	s := &UserStats{}
	// notes（含 user_id 为空的历史数据）
	err := database.DB.QueryRow(`
		SELECT COUNT(*) FROM notes WHERE user_id = ? OR user_id IS NULL
	`, userID).Scan(&s.NotesCount)
	if err != nil {
		return nil, err
	}
	err = database.DB.QueryRow(`
		SELECT COUNT(*) FROM tags WHERE user_id = ?
	`, userID).Scan(&s.TagsCount)
	if err != nil {
		return nil, err
	}
	err = database.DB.QueryRow(`
		SELECT COUNT(*) FROM resources WHERE user_id = ?
	`, userID).Scan(&s.ResourcesCount)
	if err != nil {
		return nil, err
	}
	err = database.DB.QueryRow(`
		SELECT COUNT(*) FROM notebooks WHERE user_id = ?
	`, userID).Scan(&s.NotebooksCount)
	if err != nil {
		return nil, err
	}
	err = database.DB.QueryRow(`
		SELECT COUNT(*) FROM notes WHERE (user_id = ? OR user_id IS NULL) AND pinned = 1
	`, userID).Scan(&s.PinnedCount)
	if err != nil {
		return nil, err
	}
	err = database.DB.QueryRow(`
		SELECT COUNT(*) FROM notes WHERE (user_id = ? OR user_id IS NULL) AND created_at >= datetime('now', '-7 days')
	`, userID).Scan(&s.NotesCreated7d)
	if err != nil {
		return nil, err
	}
	err = database.DB.QueryRow(`
		SELECT COUNT(*) FROM notes WHERE (user_id = ? OR user_id IS NULL) AND updated_at >= datetime('now', '-7 days')
	`, userID).Scan(&s.NotesUpdated7d)
	if err != nil {
		return nil, err
	}
	return s, nil
}
