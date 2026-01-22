package handlers

import (
	"net/http"
	"strconv"

	"memo-studio/backend/models"

	"github.com/gin-gonic/gin"
)

type CreateNoteRequest struct {
	Title   string   `json:"title" binding:"required"`
	Content string   `json:"content" binding:"required"`
	Tags    []string `json:"tags"`
}

// GetNotes 获取所有笔记
func GetNotes(c *gin.Context) {
	notes, err := models.GetAllNotes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建或获取标签
	var tagIDs []int
	for _, tagName := range req.Tags {
		if tagName == "" {
			continue
		}
		tag, err := models.CreateTagIfNotExists(tagName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		tagIDs = append(tagIDs, tag.ID)
	}

	// 创建笔记
	note, err := models.CreateNote(req.Title, req.Content, tagIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, note)
}

// GetTags 获取所有标签
func GetTags(c *gin.Context) {
	tags, err := models.GetAllTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tags)
}
