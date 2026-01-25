package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"memo-studio/backend/models"

	"github.com/gin-gonic/gin"
)

type CreateNoteRequest struct {
	Title   interface{} `json:"title"`
	Content interface{}   `json:"content"`
	Tags    []string    `json:"tags"`
}

type UpdateNoteRequest struct {
	Title   interface{} `json:"title"`
	Content interface{}   `json:"content"`
	Tags    []string    `json:"tags"`
}

// normalizeString 将 interface{} 转换为字符串，处理各种类型
func normalizeString(v interface{}) string {
	if v == nil {
		return ""
	}
	
	switch val := v.(type) {
	case string:
		// 检查是否是错误的 "[object Object]" 字符串（完全匹配）
		if val == "[object Object]" || val == "[object object]" {
			return ""
		}
		// 正常字符串直接返回
		return val
	case []byte:
		str := string(val)
		if str == "[object Object]" || str == "[object object]" {
			return ""
		}
		return str
	case map[string]interface{}:
		// 如果是对象，尝试提取常见字段（允许空字符串）
		if content, ok := val["content"].(string); ok {
			if content != "[object Object]" && content != "[object object]" {
				return content
			}
		}
		if text, ok := val["text"].(string); ok {
			if text != "[object Object]" && text != "[object object]" {
				return text
			}
		}
		if value, ok := val["value"].(string); ok {
			if value != "[object Object]" && value != "[object object]" {
				return value
			}
		}
		// 如果对象中没有这些字段，尝试转换为 JSON 字符串
		jsonBytes, err := json.Marshal(val)
		if err == nil {
			jsonStr := string(jsonBytes)
			// 如果转换后的 JSON 是 "[object Object]"，返回空
			if jsonStr == `"[object Object]"` || jsonStr == `"[object object]"` {
				return ""
			}
			return jsonStr
		}
		return ""
	default:
		// 其他类型，尝试转换为字符串
		// 对于实现了 Stringer 接口的类型
		if str, ok := val.(fmt.Stringer); ok {
			result := str.String()
			if result == "[object Object]" || result == "[object object]" {
				return ""
			}
			return result
		}
		// 尝试 JSON 序列化
		jsonBytes, err := json.Marshal(val)
		if err == nil {
			jsonStr := string(jsonBytes)
			// 移除 JSON 字符串的引号（如果是字符串类型）
			if len(jsonStr) > 2 && jsonStr[0] == '"' && jsonStr[len(jsonStr)-1] == '"' {
				jsonStr = jsonStr[1 : len(jsonStr)-1]
			}
			if jsonStr == "[object Object]" || jsonStr == "[object object]" {
				return ""
			}
			return jsonStr
		}
		// 如果都失败了，使用 fmt.Sprintf 作为最后手段
		result := fmt.Sprintf("%v", val)
		if result == "[object Object]" || result == "[object object]" {
			return ""
		}
		return result
	}
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

type CreateTagRequest struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color"`
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

	uidAny, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	userID := uidAny.(int)

	// 规范化 title 和 content 为字符串
	title := normalizeString(req.Title)
	content := normalizeString(req.Content)

	// 调试日志
	log.Printf("[CreateNote] 接收 - title: %v (type: %T), content: %v (type: %T)", req.Title, req.Title, req.Content, req.Content)
	log.Printf("[CreateNote] 规范化后 - title: %q (len: %d), content: %q (len: %d)", title, len(title), content, len(content))

	// 验证必填字段
	if title == "" && content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标题和内容不能同时为空"})
		return
	}

	// 创建或获取标签
	var tagIDs []int
	for _, tagName := range req.Tags {
		if tagName == "" {
			continue
		}
		tag, err := models.CreateTagIfNotExists(tagName, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建标签失败: " + err.Error()})
			return
		}
		tagIDs = append(tagIDs, tag.ID)
	}

	// 创建笔记（旧接口：默认 markdown、不置顶、无附件、无 user_id）
	note, err := models.CreateNote(title, content, tagIDs, false, "markdown", nil, &userID)
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

	uidAny, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	userID := uidAny.(int)

	var req UpdateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 规范化 title 和 content 为字符串
	title := normalizeString(req.Title)
	content := normalizeString(req.Content)

	// 调试日志
	log.Printf("[UpdateNote] 接收 - title: %v (type: %T), content: %v (type: %T)", req.Title, req.Title, req.Content, req.Content)
	log.Printf("[UpdateNote] 规范化后 - title: %q (len: %d), content: %q (len: %d)", title, len(title), content, len(content))

	// 验证必填字段
	if content == "" && title == "" {
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
		tag, err := models.CreateTagIfNotExists(tagName, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建标签失败: " + err.Error()})
			return
		}
		tagIDs = append(tagIDs, tag.ID)
	}

	// 更新笔记（旧接口：默认 markdown、不置顶、无附件）
	note, err := models.UpdateNote(id, title, content, tagIDs, false, "markdown", nil)
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
	uidAny, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	userID := uidAny.(int)
	// withCount=1 时返回包含计数的标签列表（给侧边栏用）
	if c.Query("withCount") == "1" {
		tags, err := models.GetTagsWithCount(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取标签列表失败: " + err.Error()})
			return
		}
		if tags == nil {
			tags = []models.TagWithCount{}
		}
		c.JSON(http.StatusOK, tags)
		return
	}

	tags, err := models.GetAllTags(userID)
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

// CreateTag 创建标签
func CreateTag(c *gin.Context) {
	uidAny, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	userID := uidAny.(int)
	var req CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "标签名称不能为空"})
		return
	}

	tag, err := models.CreateTagIfNotExists(strings.TrimSpace(req.Name), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建标签失败: " + err.Error()})
		return
	}

	// 如果传了 color，则更新一次
	if strings.TrimSpace(req.Color) != "" {
		updated, err := models.UpdateTag(tag.ID, tag.Name, strings.TrimSpace(req.Color))
		if err == nil {
			c.JSON(http.StatusCreated, updated)
			return
		}
	}

	c.JSON(http.StatusCreated, tag)
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
