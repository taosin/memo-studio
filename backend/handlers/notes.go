package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"memo-studio/backend/database"
	"memo-studio/backend/models"

	"github.com/gin-gonic/gin"
)

type CreateNoteRequest struct {
	Title       interface{} `json:"title" binding:"required"`
	Content     interface{} `json:"content"`
	Tags        []string    `json:"tags"`
	NotebookIDs []int       `json:"notebook_ids"`
}

type UpdateNoteRequest struct {
	Title       interface{} `json:"title"`
	Content     interface{} `json:"content"`
	Tags        []string    `json:"tags"`
	NotebookIDs []int       `json:"notebook_ids"`
}

// normalizeString 将 interface{} 转换为字符串，处理各种类型
func normalizeString(v interface{}) string {
	if v == nil {
		return ""
	}

	switch val := v.(type) {
	case string:
		if val == "[object Object]" || val == "[object object]" {
			return ""
		}
		return val
	case []byte:
		str := string(val)
		if str == "[object Object]" || str == "[object object]" {
			return ""
		}
		return str
	case map[string]interface{}:
		if content, ok := val["content"].(string); ok && content != "[object Object]" && content != "[object object]" {
			return content
		}
		if text, ok := val["text"].(string); ok && text != "[object Object]" && text != "[object object]" {
			return text
		}
		if value, ok := val["value"].(string); ok && value != "[object Object]" && value != "[object object]" {
			return value
		}
		jsonBytes, _ := json.Marshal(val)
		jsonStr := string(jsonBytes)
		if jsonStr == `"[object Object]"` || jsonStr == `"[object object]"` {
			return ""
		}
		return jsonStr
	default:
		if str, ok := val.(fmt.Stringer); ok {
			result := str.String()
			if result == "[object Object]" || result == "[object object]" {
				return ""
			}
			return result
		}
		jsonBytes, _ := json.Marshal(val)
		jsonStr := string(jsonBytes)
		if len(jsonStr) > 2 && jsonStr[0] == '"' && jsonStr[len(jsonStr)-1] == '"' {
			jsonStr = jsonStr[1 : len(jsonStr)-1]
		}
		if jsonStr == "[object Object]" || jsonStr == "[object object]" {
			return ""
		}
		return jsonStr
	}
}

type BatchDeleteRequest struct {
	IDs []int `json:"ids" binding:"required"`
}

type UpdateTagRequest struct {
	Name  string `json:"name" binding:"required,min=1,max=50"`
	Color string `json:"color"`
}

type MergeTagsRequest struct {
	SourceID int `json:"sourceId" binding:"required"`
	TargetID int `json:"targetId" binding:"required"`
}

type CreateTagRequest struct {
	Name  string `json:"name" binding:"required,min=1,max=50"`
	Color string `json:"color"`
}

func mustUserID(c *gin.Context) (int, bool) {
	uidAny, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return 0, false
	}
	return uidAny.(int), true
}

func ensureNoteOwned(c *gin.Context, noteID int, userID int) bool {
	var owner sql.NullInt64
	err := database.DB.QueryRow("SELECT user_id FROM notes WHERE id = ?", noteID).Scan(&owner)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "笔记不存在"})
			return false
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return false
	}
	if owner.Valid && int(owner.Int64) != userID {
		c.JSON(http.StatusNotFound, gin.H{"error": "笔记不存在"})
		return false
	}
	return true
}

// GetNotes 获取所有笔记
func GetNotes(c *gin.Context) {
	userID, ok := mustUserID(c)
	if !ok {
		return
	}
	notes, err := models.ListMemos(models.MemoQuery{
		Limit:  200,
		Offset: 0,
		UserID: &userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取笔记列表失败"})
		return
	}
	if notes == nil {
		notes = []models.Note{}
	}
	c.JSON(http.StatusOK, notes)
}

// GetNote 获取单个笔记
func GetNote(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的笔记ID"})
		return
	}
	userID, ok := mustUserID(c)
	if !ok {
		return
	}
	if !ensureNoteOwned(c, id, userID) {
		return
	}

	note, err := models.GetNote(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "笔记不存在"})
		return
	}
	c.JSON(http.StatusOK, note)
}

// CreateNote 创建笔记
func CreateNote(c *gin.Context) {
	var req CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	userID, ok := mustUserID(c)
	if !ok {
		return
	}

	title := normalizeString(req.Title)
	content := normalizeString(req.Content)

	if title == "" && content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标题和内容不能同时为空"})
		return
	}

	var tagIDs []int
	for _, tagName := range req.Tags {
		if tagName == "" {
			continue
		}
		tag, err := models.CreateTagIfNotExists(tagName, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建标签失败"})
			return
		}
		tagIDs = append(tagIDs, tag.ID)
	}

	note, err := models.CreateNote(title, content, tagIDs, false, "markdown", nil, &userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建笔记失败"})
		return
	}

	if len(req.NotebookIDs) > 0 {
		var validIDs []int
		for _, nid := range req.NotebookIDs {
			if nid <= 0 {
				continue
			}
			nb, _ := models.GetNotebook(nid, userID)
			if nb != nil {
				validIDs = append(validIDs, nid)
			}
		}
		_ = models.SetNoteNotebooks(note.ID, validIDs)
	}

	c.JSON(http.StatusCreated, note)
}

// UpdateNote 更新笔记
func UpdateNote(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的笔记ID"})
		return
	}

	userID, ok := mustUserID(c)
	if !ok {
		return
	}
	if !ensureNoteOwned(c, id, userID) {
		return
	}

	var req UpdateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	title := normalizeString(req.Title)
	content := normalizeString(req.Content)

	if content == "" && title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标题和内容不能同时为空"})
		return
	}

	var tagIDs []int
	for _, tagName := range req.Tags {
		if tagName == "" {
			continue
		}
		tag, err := models.CreateTagIfNotExists(tagName, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建标签失败"})
			return
		}
		tagIDs = append(tagIDs, tag.ID)
	}

	note, err := models.UpdateNote(id, title, content, tagIDs, false, "markdown", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新笔记失败"})
		return
	}

	if req.NotebookIDs != nil {
		var validIDs []int
		for _, nid := range req.NotebookIDs {
			if nid <= 0 {
				continue
			}
			nb, _ := models.GetNotebook(nid, userID)
			if nb != nil {
				validIDs = append(validIDs, nid)
			}
		}
		_ = models.SetNoteNotebooks(id, validIDs)
	}

	c.JSON(http.StatusOK, note)
}

// DeleteNote 删除笔记
func DeleteNote(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的笔记ID"})
		return
	}
	userID, ok := mustUserID(c)
	if !ok {
		return
	}
	if !ensureNoteOwned(c, id, userID) {
		return
	}

	err = models.DeleteNote(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除笔记失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "笔记已删除"})
}

// DeleteNotes 批量删除笔记
func DeleteNotes(c *gin.Context) {
	var req BatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	if len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择要删除的笔记"})
		return
	}
	userID, ok := mustUserID(c)
	if !ok {
		return
	}

	placeholders := strings.Repeat("?,", len(req.IDs))
	placeholders = strings.TrimSuffix(placeholders, ",")
	args := make([]interface{}, 0, len(req.IDs)+1)
	args = append(args, userID)
	for _, id := range req.IDs {
		args = append(args, id)
	}
	_, err := database.DB.Exec("DELETE FROM notes WHERE user_id = ? AND id IN ("+placeholders+")", args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "批量删除笔记失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "deleted": len(req.IDs)})
}

// GetTags 获取所有标签
func GetTags(c *gin.Context) {
	uidAny, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	userID := uidAny.(int)

	if c.Query("withCount") == "1" {
		tags, err := models.GetTagsWithCount(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取标签列表失败"})
			return
		}
		if tags == nil {
			tags = []models.TagWithCount{}
		}
		c.JSON(http.StatusOK, tags)
		return
	}

	tags, err := models.GetAllTags(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取标签列表失败"})
		return
	}
	if tags == nil {
		tags = []models.Tag{}
	}
	c.JSON(http.StatusOK, tags)
}

// CreateTag 创建标签
func CreateTag(c *gin.Context) {
	uidAny, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	userID := uidAny.(int)

	var req CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标签名称不能为空"})
		return
	}

	tag, err := models.CreateTagIfNotExists(strings.TrimSpace(req.Name), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建标签失败"})
		return
	}

	if strings.TrimSpace(req.Color) != "" {
		updated, err := models.UpdateTag(tag.ID, tag.Name, strings.TrimSpace(req.Color))
		if err == nil {
			c.JSON(http.StatusCreated, updated)
			return
		}
	}

	c.JSON(http.StatusCreated, tag)
}

// UpdateTag 更新标签
func UpdateTag(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的标签ID"})
		return
	}

	var req UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标签名称不能为空"})
		return
	}

	_, err = models.GetTagByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "标签不存在"})
		return
	}

	tag, err := models.UpdateTag(id, req.Name, req.Color)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新标签失败"})
		return
	}

	c.JSON(http.StatusOK, tag)
}

// DeleteTag 删除标签
func DeleteTag(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的标签ID"})
		return
	}

	_, err = models.GetTagByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "标签不存在"})
		return
	}

	err = models.DeleteTag(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除标签失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "标签已删除"})
}

// MergeTags 合并标签
func MergeTags(c *gin.Context) {
	var req MergeTagsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	if req.SourceID == req.TargetID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能合并相同的标签"})
		return
	}

	_, err := models.GetTagByID(req.SourceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "源标签不存在"})
		return
	}

	_, err = models.GetTagByID(req.TargetID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "目标标签不存在"})
		return
	}

	err = models.MergeTags(req.SourceID, req.TargetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "合并标签失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "标签合并成功"})
}
