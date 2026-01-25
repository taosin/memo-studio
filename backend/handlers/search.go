package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"memo-studio/backend/models"

	"github.com/gin-gonic/gin"
)

// SearchNotes FTS5 全文搜索
// GET /api/search?q=...&limit=50&offset=0
func SearchNotes(c *gin.Context) {
	q := strings.TrimSpace(c.Query("q"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	notes, err := models.SearchNotes(q, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "全文搜索失败: " + err.Error()})
		return
	}

	if notes == nil {
		notes = []models.Note{}
	}
	c.JSON(http.StatusOK, notes)
}

