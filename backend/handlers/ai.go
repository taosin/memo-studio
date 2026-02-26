package handlers

import (
	"net/http"
	"os"

	"memo-studio/backend/services"

	"github.com/gin-gonic/gin"
)

// AIAssistRequest AI 辅助请求
type AIAssistRequest struct {
	Action  string `json:"action"`  // polish, continue, summarize, translate, expand, simplify
	Content string `json:"content"` // 输入内容
}

// AIAssistResponse AI 辅助响应
type AIAssistResponse struct {
	Result string `json:"result"`
	Action string `json:"action"`
}

// getPromptForAction 根据 action 获取提示词
func getPromptForAction(action, content string) string {
	switch action {
	case "polish":
		return "请润色以下文本，使其表达更流畅、更专业，保持原意不变：\n\n" + content
	case "continue":
		return "请根据以下内容继续写作，保持风格一致，续写约100-200字：\n\n" + content
	case "summarize":
		return "请总结以下内容的要点，生成简洁的摘要（约50-100字）：\n\n" + content
	case "translate":
		return "请将以下内容翻译成英文（如果是英文则翻译成中文）：\n\n" + content
	case "expand":
		return "请扩展以下内容，增加更多细节和说明，使其更加丰富完整：\n\n" + content
	case "simplify":
		return "请简化以下内容，使其更加简洁易懂，去除冗余表达：\n\n" + content
	default:
		return "请处理以下内容：\n\n" + content
	}
}

// AIAssist AI 辅助接口
// POST /api/ai/assist
func AIAssist(c *gin.Context) {
	var req AIAssistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误"})
		return
	}

	if req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "内容不能为空"})
		return
	}

	if req.Action == "" {
		req.Action = "polish"
	}

	// 检查是否配置了 LLM API Key
	hasAPIKey := os.Getenv("OPENAI_API_KEY") != "" ||
		os.Getenv("LLM_API_KEY") != "" ||
		os.Getenv("ANTHROPIC_API_KEY") != "" ||
		os.Getenv("DEEPSEEK_API_KEY") != "" ||
		os.Getenv("ZHIPU_API_KEY") != ""

	if !hasAPIKey {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "AI 功能未配置，请在环境变量中设置 LLM API Key",
		})
		return
	}

	// 创建 LLM 服务
	llmService := services.NewLLMService()

	// 构建提示词
	prompt := getPromptForAction(req.Action, req.Content)

	// 调用 LLM
	messages := []services.ChatMessage{
		{Role: "system", Content: "你是一个专业的写作助手，帮助用户优化和处理文本内容。请直接输出处理结果，不要添加任何解释或说明。"},
		{Role: "user", Content: prompt},
	}

	result, err := llmService.Chat(messages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "AI 处理失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, AIAssistResponse{
		Result: result,
		Action: req.Action,
	})
}
