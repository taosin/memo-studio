package handlers

import (
	"net/http"
	"strconv"

	"memo-studio/backend/models"

	"github.com/gin-gonic/gin"
)

// ListNotebooks GET /api/notebooks
func ListNotebooks(c *gin.Context) {
	userID, ok := mustUserID(c)
	if !ok {
		return
	}
	list, err := models.ListNotebooks(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取笔记本列表失败: " + err.Error()})
		return
	}
	if list == nil {
		list = []models.Notebook{}
	}
	c.JSON(http.StatusOK, list)
}

// GetNotebook GET /api/notebooks/:id
func GetNotebook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的笔记本ID"})
		return
	}
	userID, ok := mustUserID(c)
	if !ok {
		return
	}
	nb, err := models.GetNotebook(id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if nb == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "笔记本不存在"})
		return
	}
	c.JSON(http.StatusOK, nb)
}

type CreateNotebookRequest struct {
	Name      string `json:"name"`
	Color     string `json:"color"`
	SortOrder int    `json:"sort_order"`
}

type UpdateNotebookRequest struct {
	Name      string `json:"name"`
	Color     string `json:"color"`
	SortOrder *int   `json:"sort_order"`
}

// CreateNotebook POST /api/notebooks
func CreateNotebook(c *gin.Context) {
	var req CreateNotebookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}
	userID, ok := mustUserID(c)
	if !ok {
		return
	}
	nb, err := models.CreateNotebook(userID, req.Name, req.Color, req.SortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建笔记本失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, nb)
}

// UpdateNotebook PUT /api/notebooks/:id
func UpdateNotebook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的笔记本ID"})
		return
	}
	userID, ok := mustUserID(c)
	if !ok {
		return
	}
	var req UpdateNotebookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}
	nb, err := models.UpdateNotebook(id, userID, req.Name, req.Color, req.SortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新笔记本失败: " + err.Error()})
		return
	}
	if nb == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "笔记本不存在"})
		return
	}
	c.JSON(http.StatusOK, nb)
}

// DeleteNotebook DELETE /api/notebooks/:id
func DeleteNotebook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的笔记本ID"})
		return
	}
	userID, ok := mustUserID(c)
	if !ok {
		return
	}
	if err := models.DeleteNotebook(id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除笔记本失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// ListNotebookNotes GET /api/notebooks/:id/notes?limit=50&offset=0
func ListNotebookNotes(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的笔记本ID"})
		return
	}
	userID, ok := mustUserID(c)
	if !ok {
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	if offset < 0 {
		offset = 0
	}
	notes, err := models.ListNotesByNotebookID(id, userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取笔记列表失败: " + err.Error()})
		return
	}
	if notes == nil {
		notes = []models.Note{}
	}
	c.JSON(http.StatusOK, notes)
}
