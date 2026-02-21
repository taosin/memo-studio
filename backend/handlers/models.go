package handlers

import (
	"net/http"
	"os"

	"memo-studio/backend/services"

	"github.com/gin-gonic/gin"
)

// ModelResponse 模型响应
type ModelResponse struct {
	Category   string                `json:"category"`
	Models     []ModelInfo          `json:"models"`
	LocalModels []LocalModelInfo     `json:"local_models,omitempty"`
	Active     ActiveModelInfo       `json:"active"`
}

// ModelInfo 模型信息
type ModelInfo struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	Category  string `json:"category"`
	Model     string `json:"model"`
	Available bool   `json:"available"`
	LocalURL  string `json:"local_url,omitempty"`
}

// LocalModelInfo 本地模型信息
type LocalModelInfo struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	URL       string `json:"url"`
	Status    string `json:"status"`
	Available bool   `json:"available"`
}

// ActiveModelInfo 当前激活的模型
type ActiveModelInfo struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Model    string `json:"model"`
	BaseURL  string `json:"base_url"`
}

// SetModelRequest 设置模型请求
type SetModelRequest struct {
	Type    string `json:"type"`              // 模型类型
	BaseURL string `json:"base_url,omitempty"` // 自定义 API 地址（本地模型）
	Model   string `json:"model,omitempty"`   // 模型名称
}

// LocalHealthRequest 本地服务健康检查
type LocalHealthRequest struct {
	BaseURL string `json:"base_url"`
}

// SetActiveModel 设置当前模型
// POST /api/models/active
func SetActiveModel(c *gin.Context) {
	var req SetModelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误"})
		return
	}

	// 验证模型类型
	valid := false
	for _, m := range services.DefaultModels() {
		if string(m.Type) == req.Type {
			valid = true
			break
		}
	}
	// 本地模型
	for _, m := range services.LocalModels {
		if string(m.Type) == req.Type {
			valid = true
			break
		}
	}

	if !valid && req.BaseURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的模型类型"})
		return
	}

	// 设置环境变量
	os.Setenv("LLM_MODEL_TYPE", req.Type)
	if req.BaseURL != "" {
		os.Setenv("LLM_BASE_URL", req.BaseURL)
	}
	if req.Model != "" {
		os.Setenv("LLM_MODEL", req.Model)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "模型切换成功",
		"model":   req.Type,
	})
}

// AddLocalModel 添加自定义本地模型
// POST /api/models/local
func AddLocalModel(c *gin.Context) {
	var req struct {
		Name    string `json:"name"`
		Type    string `json:"type"` // ollama/lmstudio/localai/anything
		BaseURL string `json:"base_url"`
		Model   string `json:"model"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误"})
		return
	}

	if req.BaseURL == "" || req.Model == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "BaseURL 和 Model 不能为空"})
		return
	}

	// 验证连接
	llmService := services.NewLLMService()
	llmService.SetLocalModel(req.BaseURL, req.Model)
	healthy, message := llmService.CheckLocalHealth(req.BaseURL)

	c.JSON(http.StatusOK, gin.H{
		"success":  healthy,
		"message":  message,
		"name":     req.Name,
		"base_url": req.BaseURL,
		"model":    req.Model,
	})
}

// CheckLocalHealth 检查本地模型服务健康状态
// POST /api/models/local/health
func CheckLocalHealth(c *gin.Context) {
	var req LocalHealthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误"})
		return
	}

	if req.BaseURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供服务地址"})
		return
	}

	llmService := services.NewLLMService()
	healthy, message := llmService.CheckLocalHealth(req.BaseURL)

	c.JSON(http.StatusOK, gin.H{
		"healthy": healthy,
		"message": message,
		"url":     req.BaseURL,
	})
}

// GetModels 获取模型列表（分类）
// GET /api/models
func GetModels(c *gin.Context) {
	category := c.Query("category") // cloud/local/all

	response := ModelResponse{
		Active: ActiveModelInfo{},
	}

	// 当前激活的模型
	activeModel := services.GetActiveModel()
	response.Active = ActiveModelInfo{
		Type:     string(activeModel.Type),
		Name:     activeModel.Name,
		Category: string(activeModel.Category),
		Model:    activeModel.Model,
		BaseURL:  activeModel.BaseURL,
	}

	// 检查云端 API Key
	hasGlobalKey := os.Getenv("LLM_API_KEY") != ""
	hasOpenAI := hasGlobalKey || os.Getenv("OPENAI_API_KEY") != ""
	hasClaude := hasGlobalKey || os.Getenv("ANTHROPIC_API_KEY") != ""
	hasDeepSeek := hasGlobalKey || os.Getenv("DEEPSEEK_API_KEY") != ""
	hasZhipu := hasGlobalKey || os.Getenv("ZHIPU_API_KEY") != ""

	// 添加云端模型
	if category == "" || category == "cloud" {
		for _, m := range services.DefaultModels() {
			if m.Category == services.CategoryCloud {
				available := false
				switch m.Type {
				case services.ModelOpenAI:
					available = hasOpenAI
				case services.ModelClaude:
					available = hasClaude
				case services.ModelDeepSeek:
					available = hasDeepSeek
				case services.ModelGLM:
					available = hasZhipu
				default:
					available = hasGlobalKey
				}

				response.Models = append(response.Models, ModelInfo{
					Type:      string(m.Type),
					Name:      m.Name,
					Category:  string(m.Category),
					Model:     m.Model,
					Available: available,
				})
			}
		}
	}

	// 添加本地模型配置
	if category == "" || category == "local" {
		for _, m := range services.LocalModels {
			response.LocalModels = append(response.LocalModels, LocalModelInfo{
				Name:      m.Name,
				Type:      string(m.Type),
				URL:       m.BaseURL,
				Status:    "config", // config/ready/error
				Available: true,
			})
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetCloudModels 获取云端模型列表
// GET /api/models/cloud
func GetCloudModels(c *gin.Context) {
	response := make([]ModelInfo, 0)

	hasGlobalKey := os.Getenv("LLM_API_KEY") != ""
	hasOpenAI := hasGlobalKey || os.Getenv("OPENAI_API_KEY") != ""

	for _, m := range services.DefaultModels() {
		if m.Category == services.CategoryCloud {
			response = append(response, ModelInfo{
				Type:      string(m.Type),
				Name:      m.Name,
				Category:  string(m.Category),
				Model:     m.Model,
				Available: hasOpenAI,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"models": response,
	})
}

// GetLocalModels 获取本地模型列表
// GET /api/models/local
func GetLocalModels(c *gin.Context) {
	response := make([]LocalModelInfo, 0)

	for _, m := range services.LocalModels {
		response = append(response, LocalModelInfo{
			Name:      m.Name,
			Type:      string(m.Type),
			URL:       m.BaseURL,
			Status:    "config",
			Available: true,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"models": response,
	})
}

// GetModelConfig 获取当前模型配置
// GET /api/models/config
func GetModelConfig(c *gin.Context) {
	config := services.GetActiveModel()

	c.JSON(http.StatusOK, gin.H{
		"type":      string(config.Type),
		"name":      config.Name,
		"category":  string(config.Category),
		"model":     config.Model,
		"base_url":  config.BaseURL,
		"max_tokens": config.MaxTokens,
		"context":   config.Context,
	})
}

// GetAvailableModels 获取可用的模型列表（根据环境变量）
// GET /api/models/available
func GetAvailableModels(c *gin.Context) {
	hasGlobalKey := os.Getenv("LLM_API_KEY") != ""
	hasOpenAI := hasGlobalKey || os.Getenv("OPENAI_API_KEY") != ""
	hasClaude := hasGlobalKey || os.Getenv("ANTHROPIC_API_KEY") != ""
	hasDeepSeek := hasGlobalKey || os.Getenv("DEEPSEEK_API_KEY") != ""
	hasZhipu := hasGlobalKey || os.Getenv("ZHIPU_API_KEY") != ""

	available := make([]ModelInfo, 0)

	for _, m := range services.DefaultModels() {
		if m.Category == services.CategoryCloud {
			isAvailable := false
			switch m.Type {
			case services.ModelOpenAI:
				isAvailable = hasOpenAI
			case services.ModelClaude:
				isAvailable = hasClaude
			case services.ModelDeepSeek:
				isAvailable = hasDeepSeek
			case services.ModelGLM:
				isAvailable = hasZhipu
			default:
				isAvailable = hasGlobalKey
			}

			if isAvailable {
				available = append(available, ModelInfo{
					Type:      string(m.Type),
					Name:      m.Name,
					Category:  string(m.Category),
					Model:     m.Model,
					Available: true,
				})
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"models": available,
		"hint":   "配置 LLM_API_KEY 后可用",
	})
}

// TestModelConnection 测试模型连接
// POST /api/models/test
func TestModelConnection(c *gin.Context) {
	var req struct {
		BaseURL string `json:"base_url"`
		Model   string `json:"model"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误"})
		return
	}

	if req.BaseURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供服务地址"})
		return
	}

	llmService := services.NewLLMService()
	llmService.SetLocalModel(req.BaseURL, req.Model)

	healthy, message := llmService.CheckLocalHealth(req.BaseURL)

	c.JSON(http.StatusOK, gin.H{
		"success": healthy,
		"message": message,
		"url":     req.BaseURL,
		"model":   req.Model,
	})
}
