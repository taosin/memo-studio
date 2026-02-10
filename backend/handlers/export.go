package handlers

import (
	"bytes"
	"net/http"
	"strconv"
	"strings"
	"time"

	"memo-studio/backend/models"

	"github.com/gin-gonic/gin"
)

// ExportNotes GET /api/export?format=json|markdown&limit=500
func ExportNotes(c *gin.Context) {
	userID, ok := mustUserID(c)
	if !ok {
		return
	}
	format := strings.ToLower(strings.TrimSpace(c.DefaultQuery("format", "json")))
	if format != "json" && format != "markdown" {
		format = "json"
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "500"))
	if limit <= 0 {
		limit = 500
	}
	if limit > 2000 {
		limit = 2000
	}
	notes, err := models.ListMemos(models.MemoQuery{
		Limit:  limit,
		Offset: 0,
		UserID: &userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "导出失败: " + err.Error()})
		return
	}
	if notes == nil {
		notes = []models.Note{}
	}
	if format == "markdown" {
		var buf bytes.Buffer
		buf.WriteString("# Memo Studio 导出\n\n")
		buf.WriteString("导出时间: " + time.Now().Format(time.RFC3339) + "\n\n")
		buf.WriteString("---\n\n")
		for i, n := range notes {
			buf.WriteString("## ")
			buf.WriteString(escapeMarkdownTitle(n.Title))
			buf.WriteString("\n\n")
			if len(n.Tags) > 0 {
				buf.WriteString("标签: ")
				for j, t := range n.Tags {
					if j > 0 {
						buf.WriteString(", ")
					}
					buf.WriteString(t.Name)
				}
				buf.WriteString("\n\n")
			}
			buf.WriteString(n.Content)
			buf.WriteString("\n\n")
			if i < len(notes)-1 {
				buf.WriteString("---\n\n")
			}
		}
		c.Header("Content-Type", "text/markdown; charset=utf-8")
		c.Header("Content-Disposition", "attachment; filename=memo-export-"+time.Now().Format("20060102-150405")+".md")
		c.Data(http.StatusOK, "text/markdown; charset=utf-8", buf.Bytes())
		return
	}
	// json
	c.Header("Content-Disposition", "attachment; filename=memo-export-"+time.Now().Format("20060102-150405")+".json")
	c.JSON(http.StatusOK, gin.H{
		"exported_at": time.Now().Format(time.RFC3339),
		"count":       len(notes),
		"notes":       notes,
	})
}

func escapeMarkdownTitle(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "\r", ""), "\n", " ")
}
