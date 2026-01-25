package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"memo-studio/backend/database"
	"memo-studio/backend/models"

	"github.com/gin-gonic/gin"
)

type CreateMemoRequest struct {
	Title       string   `json:"title" binding:"max=200"`
	Content     string   `json:"content" binding:"max=50000"`
	Tags        []string `json:"tags" binding:"max=50"`
	Pinned      bool     `json:"pinned"`
	ContentType string   `json:"content_type" binding:"omitempty,oneof=markdown"`
	ResourceIDs []int    `json:"resource_ids" binding:"max=50"`
}

type UpdateMemoRequest struct {
	Title       string   `json:"title" binding:"max=200"`
	Content     string   `json:"content" binding:"max=50000"`
	Tags        []string `json:"tags" binding:"max=50"`
	Pinned      bool     `json:"pinned"`
	ContentType string   `json:"content_type" binding:"omitempty,oneof=markdown"`
	ResourceIDs []int    `json:"resource_ids" binding:"max=50"`
}

func parseTimeParam(s string) (*time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	// 支持：RFC3339 / YYYY-MM-DD / YYYY-MM-DD HH:MM:SS
	layouts := []string{
		time.RFC3339,
		"2006-01-02",
		"2006-01-02 15:04:05",
	}
	for _, l := range layouts {
		if t, err := time.ParseInLocation(l, s, time.Local); err == nil {
			return &t, nil
		}
	}
	return nil, strconv.ErrSyntax
}

func getAuthUserID(c *gin.Context) (int, bool) {
	v, ok := c.Get("userID")
	if !ok {
		return 0, false
	}
	id, ok := v.(int)
	return id, ok
}

func ensureMemoOwnerOrPublic(noteID int, userID int) error {
	var owner sql.NullInt64
	err := database.DB.QueryRow("SELECT user_id FROM notes WHERE id = ?", noteID).Scan(&owner)
	if err != nil {
		return err
	}
	// 兼容旧数据（user_id 为空）：允许操作
	if !owner.Valid {
		return nil
	}
	if int(owner.Int64) != userID {
		return sql.ErrNoRows
	}
	return nil
}

// ListMemos GET /api/memos
func ListMemos(c *gin.Context) {
	limit, offset := models.ParseLimitOffset(c.Query("limit"), c.Query("offset"))
	tags := models.ParseTagsParam(c.Query("tags"))
	if len(tags) == 0 {
		// 兼容：tag=xxx
		tags = models.ParseTagsParam(c.Query("tag"))
	}

	from, err := parseTimeParam(c.Query("from"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from 参数格式错误"})
		return
	}
	to, err := parseTimeParam(c.Query("to"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to 参数格式错误"})
		return
	}
	pinned, err := models.ParseBoolParam(c.Query("pinned"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pinned 参数格式错误"})
		return
	}

	contentType := strings.TrimSpace(c.Query("content_type"))
	if contentType == "" {
		contentType = strings.TrimSpace(c.Query("type"))
	}
	if contentType != "" && contentType != "markdown" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "目前仅支持 content_type=markdown"})
		return
	}

	userID, ok := getAuthUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	notes, err := models.ListMemos(models.MemoQuery{
		Limit:       limit,
		Offset:      offset,
		Q:           c.Query("q"),
		Tags:        tags,
		From:        from,
		To:          to,
		Pinned:      pinned,
		ContentType: contentType,
		UserID:      &userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取 memos 失败: " + err.Error()})
		return
	}
	if notes == nil {
		notes = []models.Note{}
	}
	c.JSON(http.StatusOK, notes)
}

// CreateMemo POST /api/memos
func CreateMemo(c *gin.Context) {
	userID, ok := getAuthUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	var req CreateMemoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Content = strings.TrimSpace(req.Content)
	if req.Title == "" && req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标题和内容不能同时为空"})
		return
	}
	if req.ContentType == "" {
		req.ContentType = "markdown"
	}
	if req.ContentType != "markdown" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "目前仅支持 content_type=markdown"})
		return
	}

	// 创建或获取标签
	var tagIDs []int
	for _, tagName := range req.Tags {
		tagName = strings.TrimSpace(tagName)
		if tagName == "" {
			continue
		}
		tag, err := models.CreateTagIfNotExists(tagName, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建标签失败: " + err.Error()})
			return
		}
		tagIDs = append(tagIDs, tag.ID)
	}

	uid := userID
	note, err := models.CreateNote(req.Title, req.Content, tagIDs, req.Pinned, req.ContentType, req.ResourceIDs, &uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建 memo 失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, note)
}

// UpdateMemo PUT /api/memos/:id
func UpdateMemo(c *gin.Context) {
	userID, ok := getAuthUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 memo ID"})
		return
	}
	// 权限校验（user_id 为空的旧数据也允许）
	if err := ensureMemoOwnerOrPublic(id, userID); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "memo 不存在"})
			return
		}
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限"})
		return
	}

	var req UpdateMemoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}
	req.Title = strings.TrimSpace(req.Title)
	req.Content = strings.TrimSpace(req.Content)
	if req.Title == "" && req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标题和内容不能同时为空"})
		return
	}
	if req.ContentType == "" {
		req.ContentType = "markdown"
	}
	if req.ContentType != "markdown" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "目前仅支持 content_type=markdown"})
		return
	}

	var tagIDs []int
	for _, tagName := range req.Tags {
		tagName = strings.TrimSpace(tagName)
		if tagName == "" {
			continue
		}
		tag, err := models.CreateTagIfNotExists(tagName, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建标签失败: " + err.Error()})
			return
		}
		tagIDs = append(tagIDs, tag.ID)
	}

	note, err := models.UpdateNote(id, req.Title, req.Content, tagIDs, req.Pinned, req.ContentType, req.ResourceIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新 memo 失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, note)
}

// DeleteMemo DELETE /api/memos/:id
func DeleteMemo(c *gin.Context) {
	userID, ok := getAuthUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 memo ID"})
		return
	}
	if err := ensureMemoOwnerOrPublic(id, userID); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "memo 不存在"})
			return
		}
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限"})
		return
	}
	if err := models.DeleteNote(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除 memo 失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

