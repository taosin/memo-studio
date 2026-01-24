package handlers

import (
	"net/http"
	"strconv"

	"memo-studio/backend/models"

	"github.com/gin-gonic/gin"
)

type CreateNoteRequest struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type UpdateNoteRequest struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type BatchDeleteRequest struct {
	IDs []int `json:"ids" binding:"required"`
}

type UpdateTagRequest struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color"`
}

type MergeTagsRequest struct {
	SourceID int `json:"sourceId" binding:"required"`
	TargetID int `json:"targetId" binding:"required"`
}

// GetNotes 获取所有笔记
func GetNotes(c *gin.Context) {
	notes, err := models.GetAllNotes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取笔记列表失败: " + err.Error()})
		return
	}

	// 确保返回数组（即使为空）
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 验证必填字段
	if req.Title == "" && req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标题和内容不能同时为空"})
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建标签失败: " + err.Error()})
			return
		}
		tagIDs = append(tagIDs, tag.ID)
	}

	// 创建笔记
	note, err := models.CreateNote(req.Title, req.Content, tagIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建笔记失败: " + err.Error()})
		return
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

	var req UpdateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 验证必填字段
	if req.Content == "" && req.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标题和内容不能同时为空"})
		return
	}

	// 检查笔记是否存在
	_, err = models.GetNote(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "笔记不存在"})
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建标签失败: " + err.Error()})
			return
		}
		tagIDs = append(tagIDs, tag.ID)
	}

	// 更新笔记
	note, err := models.UpdateNote(id, req.Title, req.Content, tagIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新笔记失败: " + err.Error()})
		return
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

	// 检查笔记是否存在
	_, err = models.GetNote(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "笔记不存在"})
		return
	}

	err = models.DeleteNote(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除笔记失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "笔记已删除"})
}

// DeleteNotes 批量删除笔记
func DeleteNotes(c *gin.Context) {
	var req BatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	if len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择要删除的笔记"})
		return
	}

	err := models.DeleteNotes(req.IDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "批量删除笔记失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "deleted": len(req.IDs), "message": "已删除 " + strconv.Itoa(len(req.IDs)) + " 条笔记"})
}

// GetTags 获取所有标签
func GetTags(c *gin.Context) {
	tags, err := models.GetAllTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取标签列表失败: " + err.Error()})
		return
	}

	// 确保返回数组（即使为空）
	if tags == nil {
		tags = []models.Tag{}
	}

	c.JSON(http.StatusOK, tags)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标签名称不能为空"})
		return
	}

	// 检查标签是否存在
	_, err = models.GetTagByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "标签不存在"})
		return
	}

	tag, err := models.UpdateTag(id, req.Name, req.Color)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新标签失败: " + err.Error()})
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

	// 检查标签是否存在
	_, err = models.GetTagByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "标签不存在"})
		return
	}

	err = models.DeleteTag(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除标签失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "标签已删除"})
}

// MergeTags 合并标签
func MergeTags(c *gin.Context) {
	var req MergeTagsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	if req.SourceID == req.TargetID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能合并相同的标签"})
		return
	}

	// 检查源标签是否存在
	_, err := models.GetTagByID(req.SourceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "源标签不存在"})
		return
	}

	// 检查目标标签是否存在
	_, err = models.GetTagByID(req.TargetID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "目标标签不存在"})
		return
	}

	err = models.MergeTags(req.SourceID, req.TargetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "合并标签失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "标签合并成功"})
}
