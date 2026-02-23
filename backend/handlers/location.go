package handlers

import (
	"net/http"
	"strconv"

	"memo-studio/backend/models"
	"memo-studio/backend/services"

	"github.com/gin-gonic/gin"
)

// UpdateNoteLocation 更新笔记位置
// PUT /api/memos/:id/location
func UpdateNoteLocation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的笔记 ID"})
		return
	}

	var req struct {
		Location  string  `json:"location"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误"})
		return
	}

	// 更新笔记位置
	err = models.UpdateNoteLocation(id, req.Location, req.Latitude, req.Longitude)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新位置失败: " + err.Error()})
		return
	}

	// 获取更新后的笔记
	note, err := models.GetNote(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取笔记失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"note":     note,
		"location": note.Location,
	})
}

// DetectNoteLocation 检测笔记中的位置（AI 识别）
// POST /api/memos/:id/detect-location
func DetectNoteLocation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的笔记 ID"})
		return
	}

	// 获取笔记
	note, err := models.GetNote(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "笔记不存在"})
		return
	}

	// 从内容中检测位置
	locationInfo := services.DetectAndExtractLocation(note.Content)

	if locationInfo == nil {
		c.JSON(http.StatusOK, gin.H{
			"detected": false,
			"message":  "未检测到位置信息",
		})
		return
	}

	// 返回检测到的位置
	c.JSON(http.StatusOK, gin.H{
		"detected":  true,
		"location":  locationInfo.Name,
		"latitude":  locationInfo.Latitude,
		"longitude": locationInfo.Longitude,
		"suggest":   "是否保存此位置到笔记？",
	})
}

// SaveDetectedLocation 检测并保存位置
// POST /api/memos/:id/detect-and-save
func SaveDetectedLocation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的笔记 ID"})
		return
	}

	// 获取笔记
	note, err := models.GetNote(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "笔记不存在"})
		return
	}

	// 检测位置
	locationInfo := services.DetectAndExtractLocation(note.Content)
	if locationInfo == nil {
		c.JSON(http.StatusOK, gin.H{
			"detected": false,
			"message":  "未检测到位置信息",
		})
		return
	}

	// 保存位置
	err = models.UpdateNoteLocation(id, locationInfo.Name, locationInfo.Latitude, locationInfo.Longitude)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存位置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"message":   "位置已保存",
		"location":  locationInfo.Name,
		"latitude":  locationInfo.Latitude,
		"longitude": locationInfo.Longitude,
	})
}

// GetNotesByLocation 按位置筛选笔记
// GET /api/notes?location=北京
func GetNotesByLocation(c *gin.Context) {
	location := c.Query("location")
	if location == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请指定 location 参数"})
		return
	}

	notes, err := models.GetNotesByLocation(location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"location": location,
		"count":    len(notes),
		"notes":    notes,
	})
}

// GetLocationsStats 获取所有位置统计
// GET /api/locations/stats
func GetLocationsStats(c *gin.Context) {
	stats, err := models.GetLocationStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"locations": stats,
	})
}

// BatchDetectLocations 批量检测笔记位置
// POST /api/locations/batch-detect
func BatchDetectLocations(c *gin.Context) {
	var req struct {
		NoteIDs []int `json:"note_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误"})
		return
	}

	results := make(map[int]map[string]interface{})
	for _, id := range req.NoteIDs {
		note, err := models.GetNote(id)
		if err != nil {
			continue
		}

		locationInfo := services.DetectAndExtractLocation(note.Content)
		if locationInfo != nil {
			results[id] = map[string]interface{}{
				"location":  locationInfo.Name,
				"latitude":  locationInfo.Latitude,
				"longitude": locationInfo.Longitude,
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"total":     len(req.NoteIDs),
		"detected":  len(results),
		"locations": results,
	})
}
