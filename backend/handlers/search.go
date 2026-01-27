package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"memo-studio/backend/models"

	"github.com/gin-gonic/gin"
)

// SearchNotes 全文搜索（兼容旧接口）
// GET /api/search?q=...&limit=50&offset=0
// 注意：现在统一走 /api/memos?q=...，这里保留是为了旧前端不报 404
func SearchNotes(c *gin.Context) {
	uidAny, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	userID := uidAny.(int)

	q := strings.TrimSpace(c.Query("q"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	notes, err := models.ListMemos(models.MemoQuery{
		Limit:  limit,
		Offset: offset,
		Q:      q,
		UserID: &userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "全文搜索失败: " + err.Error()})
		return
	}

	if notes == nil {
		notes = []models.Note{}
	}
	c.JSON(http.StatusOK, notes)
}

