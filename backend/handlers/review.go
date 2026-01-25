package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"memo-studio/backend/models"

	"github.com/gin-gonic/gin"
)

// RandomReview 随机回顾
// GET /api/review/random?limit=1&tag=xxx&days=0
func RandomReview(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "1"))
	tag := strings.TrimSpace(c.Query("tag"))
	days, _ := strconv.Atoi(c.DefaultQuery("days", "0"))

	notes, err := models.RandomNotes(limit, tag, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "随机回顾失败: " + err.Error()})
		return
	}
	if notes == nil {
		notes = []models.Note{}
	}
	c.JSON(http.StatusOK, notes)
}

