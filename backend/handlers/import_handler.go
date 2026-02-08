package handlers

import (
	"net/http"
	"strings"

	"memo-studio/backend/models"

	"github.com/gin-gonic/gin"
)

// ImportNoteItem 单条导入笔记结构
type ImportNoteItem struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

// ImportRequest 导入请求体
type ImportRequest struct {
	Notes []ImportNoteItem `json:"notes"`
}

// ImportNotes POST /api/import
func ImportNotes(c *gin.Context) {
	userID, ok := mustUserID(c)
	if !ok {
		return
	}
	var req ImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}
	if req.Notes == nil {
		req.Notes = []ImportNoteItem{}
	}
	if len(req.Notes) > 500 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "单次导入最多 500 条笔记"})
		return
	}
	created := 0
	failed := 0
	for _, item := range req.Notes {
		title := strings.TrimSpace(item.Title)
		content := strings.TrimSpace(item.Content)
		if title == "" && content == "" {
			failed++
			continue
		}
		if title == "" && len(content) > 0 {
			if len(content) > 77 {
				title = content[:77] + "..."
			} else {
				title = content
			}
		}
		if title == "" {
			title = "未命名"
		}
		var tagIDs []int
		for _, tagName := range item.Tags {
			if tagName == "" {
				continue
			}
			tag, err := models.CreateTagIfNotExists(tagName, userID)
			if err != nil {
				continue
			}
			tagIDs = append(tagIDs, tag.ID)
		}
		_, err := models.CreateNote(title, content, tagIDs, false, "markdown", nil, &userID)
		if err != nil {
			failed++
			continue
		}
		created++
	}
	c.JSON(http.StatusOK, gin.H{
		"created": created,
		"failed":  failed,
		"total":   len(req.Notes),
	})
}
