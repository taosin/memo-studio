package handlers

import (
	"net/http"

	"memo-studio/backend/models"

	"github.com/gin-gonic/gin"
)

// GetStats GET /api/stats
func GetStats(c *gin.Context) {
	userID, ok := mustUserID(c)
	if !ok {
		return
	}
	stats, err := models.GetUserStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}
