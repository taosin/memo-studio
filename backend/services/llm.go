package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// ModelType 模型类型
type ModelType string

const (
	// 云端模型
	ModelOpenAI    ModelType = "openai"     // OpenAI (GPT-4/GPT-3.5)
	ModelClaude    ModelType = "claude"     // Anthropic Claude
	ModelDeepSeek  ModelType = "deepseek"   // DeepSeek
	ModelGLM       ModelType = "glm"        // 智谱 GLM
	ModelYi        ModelType = "yi"         // 零一万物 Yi
	ModelQwen      ModelType = "qwen"       // 阿里通义千问
	ModelKimi      ModelType = "kimi"       // 月之暗面 Kimi
	ModelSpark     ModelType = "spark"      // 讯飞星火

	// 本地模型
	ModelOllama    ModelType = "ollama"     // Ollama 本地模型
	ModelLocalAI   ModelType = "localai"    // LocalAI
	ModelLMStudio  ModelType = "lmstudio"   // LM Studio
	ModelAnything  ModelType = "anything"    // AnythingLLM
)

// ModelCategory 模型分类
type ModelCategory string

const (
	CategoryCloud  ModelCategory = "cloud"   // 云端模型
	CategoryLocal  ModelCategory = "local"   // 本地模型
)

// ModelConfig 模型配置
type ModelConfig struct {
	Type      ModelType      `json:"type"`
	Name      string         `json:"name"`        // 显示名称
	Category  ModelCategory  `json:"category"`    // 分类
	APIKey    string         `json:"api_key"`     // API Key
	BaseURL   string         `json:"base_url"`    // API 地址
	Model     string         `json:"model"`       // 模型名称
	MaxTokens int            `json:"max_tokens"`  // 最大 token
	Context   int            `json:"context"`     // 上下文长度（本地模型）
	GPU       bool           `json:"gpu"`         // 是否支持 GPU
}

// LocalModelInfo 本地模型信息
type LocalModelInfo struct {
	Name       string `json:"name"`        // 模型名称
	Size       string `json:"size"`        // 模型大小
	Format     string `json:"format"`      // 格式 (gguf/ggml)
	Parameters string `json:"parameters"`  // 参数 (7B/13B/70B)
	Status     string `json:"status"`      // 状态 (ready/loading/error)
}

// ModelListResponse 模型列表响应
type ModelListResponse struct {
	Category ModelCategory    `json:"category"`
	Models   []ModelConfig    `json:"models"`
	LocalModels []LocalModelInfo `json:"local_models,omitempty"`
	Active   string           `json:"active"`
}

// DefaultModels 默认模型配置
func DefaultModels() []ModelConfig {
	return []ModelConfig{
		// ===== 云端模型 =====
		{
			Type:     ModelOpenAI,
			Name:     "OpenAI GPT-4",
			Category: CategoryCloud,
			BaseURL:  "https://api.openai.com/v1",
			Model:    "gpt-4",
			MaxTokens: 8192,
		},
		{
			Type:     ModelOpenAI,
			Name:     "OpenAI GPT-3.5",
			Category: CategoryCloud,
			BaseURL:  "https://api.openai.com/v1",
			Model:    "gpt-3.5-turbo",
			MaxTokens: 4096,
		},
		{
			Type:     ModelClaude,
			Name:     "Anthropic Claude-3",
			Category: CategoryCloud,
			BaseURL:  "https://api.anthropic.com/v1",
			Model:    "claude-3-sonnet-20240229",
			MaxTokens: 4096,
		},
		{
			Type:     ModelDeepSeek,
			Name:     "DeepSeek Chat",
			Category: CategoryCloud,
			BaseURL:  "https://api.deepseek.com/chat",
			Model:    "deepseek-chat",
			MaxTokens: 4096,
		},
		{
			Type:     ModelGLM,
			Name:     "智谱 GLM-4",
			Category: CategoryCloud,
			BaseURL:  "https://open.bigmodel.cn/api/paas/v4",
			Model:    "glm-4",
			MaxTokens: 4096,
		},
		{
			Type:     ModelYi,
			Name:     "零一万物 Yi-34B",
			Category: CategoryCloud,
			BaseURL:  "https://api.lingyiwanwu.com/v1",
			Model:    "yi-34b-chat",
			MaxTokens: 4096,
		},
		{
			Type:     ModelQwen,
			Name:     "阿里通义千问",
			Category: CategoryCloud,
			BaseURL:  "https://dashscope.aliyuncs.com/compatible-mode/v1",
			Model:    "qwen-turbo",
			MaxTokens: 4096,
		},
		{
			Type:     ModelKimi,
			Name:     "月之暗面 Kimi",
			Category: CategoryCloud,
			BaseURL:  "https://api.moonshot.cn/v1",
			Model:    "moonshot-v1-8k",
			MaxTokens: 8192,
		},
		{
			Type:     ModelSpark,
			Name:     "讯飞星火",
			Category: CategoryCloud,
			BaseURL:  "https://spark-api.xf-yun.com/v1",
			Model:    "general",
			MaxTokens: 4096,
		},

		// ===== 本地模型 =====
		{
			Type:     ModelOllama,
			Name:     "Ollama (本地)",
			Category: CategoryLocal,
			BaseURL:  "http://localhost:11434/v1",
			Model:    "llama2",
			MaxTokens: 2048,
			Context:   4096,
			GPU:       true,
		},
		{
			Type:     ModelLocalAI,
			Name:     "LocalAI (本地)",
			Category: CategoryLocal,
			BaseURL:  "http://localhost:8080/v1",
			Model:    "llama-2-7b",
			MaxTokens: 2048,
			Context:   4096,
			GPU:       false,
		},
		{
			Type:     ModelLMStudio,
			Name:     "LM Studio (本地)",
			Category: CategoryLocal,
			BaseURL:  "http://localhost:1234/v1",
			Model:    "llama-2-7b-chat",
			MaxTokens: 2048,
			Context:   4096,
			GPU:       false,
		},
		{
			Type:     ModelAnything,
			Name:     "AnythingLLM (本地)",
			Category: CategoryLocal,
			BaseURL:  "http://localhost:3001/v1",
			Model:    "mistral-7b",
			MaxTokens: 2048,
			Context:   4096,
			GPU:       false,
		},
	}
}

// GetModelsByCategory 按分类获取模型
func GetModelsByCategory(category ModelCategory) []ModelConfig {
	var result []ModelConfig
	for _, m := range DefaultModels() {
		if m.Category == category {
			result = append(result, m)
		}
	}
	return result
}

// LocalModels 本地模型配置
var LocalModels = []ModelConfig{
	{
		Type:     ModelOllama,
		Name:     "Ollama - Llama 2",
		Category: CategoryLocal,
		BaseURL:  "http://localhost:11434/v1",
		Model:    "llama2",
		Context:  4096,
	},
	{
		Type:     ModelOllama,
		Name:     "Ollama - CodeLlama",
		Category: CategoryLocal,
		BaseURL:  "http://localhost:11434/v1",
		Model:    "codellama",
		Context:  4096,
	},
	{
		Type:     ModelOllama,
		Name:     "Ollama - Mistral",
		Category: CategoryLocal,
		BaseURL:  "http://localhost:11434/v1",
		Model:    "mistral",
		Context:  4096,
	},
	{
		Type:     ModelOllama,
		Name:     "Ollama - Qwen",
		Category: CategoryLocal,
		BaseURL:  "http://localhost:11434/v1",
		Model:    "qwen",
		Context:  4096,
	},
	{
		Type:     ModelOllama,
		Name:     "Ollama - DeepSeek",
		Category: CategoryLocal,
		BaseURL:  "http://localhost:11434/v1",
		Model:    "deepseek-coder",
		Context:  4096,
	},
	{
		Type:     ModelLMStudio,
		Name:     "LM Studio - Llama 2",
		Category: CategoryLocal,
		BaseURL:  "http://localhost:1234/v1",
		Model:    "llama-2-7b-chat",
		Context:  4096,
	},
	{
		Type:     ModelLMStudio,
		Name:     "LM Studio - Mistral",
		Category: CategoryLocal,
		BaseURL:  "http://localhost:1234/v1",
		Model:    "mistral-7b-instruct",
		Context:  4096,
	},
	{
		Type:     ModelLMStudio,
		Name:     "LM Studio - Yi-6B",
		Category: CategoryLocal,
		BaseURL:  "http://localhost:1234/v1",
		Model:    "yi-6b-chat",
		Context:  4096,
	},
	{
		Type:     ModelLocalAI,
		Name:     "LocalAI - Llama 2",
		Category: CategoryLocal,
		BaseURL:  "http://localhost:8080/v1",
		Model:    "llama-2-7b",
		Context:  4096,
	},
	{
		Type:     ModelAnything,
		Name:     "AnythingLLM - Mistral",
		Category: CategoryLocal,
		BaseURL:  "http://localhost:3001/v1",
		Model:    "mistral-7b",
		Context:  4096,
	},
}

// GetActiveModel 获取当前配置的模型
func GetActiveModel() ModelConfig {
	// 1. 检查环境变量 LLM_MODEL_TYPE
	modelType := ModelType(os.Getenv("LLM_MODEL_TYPE"))
	
	for _, m := range DefaultModels() {
		if m.Type == modelType {
			// 使用环境变量覆盖配置
			return overrideFromEnv(m)
		}
	}
	
	// 2. 检查云端 API Key 配置
	if os.Getenv("OPENAI_API_KEY") != "" {
		return overrideFromEnv(DefaultModels()[0])
	}
	if os.Getenv("ANTHROPIC_API_KEY") != "" {
		return overrideFromEnv(DefaultModels()[2])
	}
	if os.Getenv("DEEPSEEK_API_KEY") != "" {
		return overrideFromEnv(DefaultModels()[3])
	}
	if os.Getenv("ZHIPU_API_KEY") != "" {
		return overrideFromEnv(DefaultModels()[4])
	}
	
	// 3. 检查统一 API Key
	if os.Getenv("LLM_API_KEY") != "" {
		return overrideFromEnv(DefaultModels()[0])
	}
	
	// 4. 默认使用第一个（GPT-4）
	return overrideFromEnv(DefaultModels()[0])
}

// overrideFromEnv 从环境变量覆盖配置
func overrideFromEnv(m ModelConfig) ModelConfig {
	if apiKey := os.Getenv("LLM_API_KEY"); apiKey != "" {
		m.APIKey = apiKey
	}
	if baseURL := os.Getenv("LLM_BASE_URL"); baseURL != "" {
		m.BaseURL = baseURL
	}
	if model := os.Getenv("LLM_MODEL"); model != "" {
		m.Model = model
	}
	return m
}

// ChatMessage 聊天消息
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest 聊天请求
type ChatRequest struct {
	Model       string       `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int         `json:"max_tokens,omitempty"`
	Temperature float64      `json:"temperature,omitempty"`
	Stream      bool        `json:"stream,omitempty"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice 选择
type Choice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

// Usage 使用量
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens     int `json:"total_tokens"`
}

// LLMService 大模型服务
type LLMService struct {
	Model ModelConfig
}

// NewLLMService 创建 LLM 服务
func NewLLMService() *LLMService {
	return &LLMService{
		Model: GetActiveModel(),
	}
}

// SetModel 设置模型
func (s *LLMService) SetModel(modelType ModelType) {
	for _, m := range DefaultModels() {
		if m.Type == modelType {
			s.Model = m
			return
		}
	}
	// 本地模型
	for _, m := range LocalModels {
		if m.Type == modelType {
			s.Model = m
			return
		}
	}
}

// SetLocalModel 设置本地模型
func (s *LLMService) SetLocalModel(baseURL, modelName string) {
	s.Model = ModelConfig{
		Type:     ModelOllama,
		Name:     "自定义本地模型",
		Category: CategoryLocal,
		BaseURL:  baseURL,
		Model:    modelName,
		Context:  4096,
	}
}

// Chat 聊天
func (s *LLMService) Chat(messages []ChatMessage) (string, error) {
	req := ChatRequest{
		Model:       s.Model.Model,
		Messages:    messages,
		MaxTokens:   s.Model.MaxTokens,
		Temperature: 0.7,
		Stream:     false,
	}

	body := s.buildRequestBody(req)
	resp, err := s.sendRequest(body)
	if err != nil {
		return "", err
	}

	return s.parseResponse(resp)
}

// buildRequestBody 构建请求体
func (s *LLMService) buildRequestBody(req ChatRequest) map[string]interface{} {
	body := map[string]interface{}{
		"model":    req.Model,
		"messages": req.Messages,
		"stream":   false,
	}
	
	if req.MaxTokens > 0 {
		body["max_tokens"] = req.MaxTokens
	}
	if req.Temperature > 0 {
		body["temperature"] = req.Temperature
	}
	
	return body
}

// sendRequest 发送请求
func (s *LLMService) sendRequest(body map[string]interface{}) ([]byte, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	httpReq, err := http.NewRequest("POST", s.Model.BaseURL+"/chat/completions", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	s.setHeaders(httpReq)

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API 错误 (%d): %s", resp.StatusCode, string(respBody))
	}

	return io.ReadAll(resp.Body)
}

// setHeaders 设置请求头
func (s *LLMService) setHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	
	// 云端模型
	switch s.Model.Type {
	case ModelClaude:
		req.Header.Set("x-api-key", s.Model.APIKey)
		req.Header.Set("anthropic-version", "2023-06-01")
	case ModelOpenAI, ModelKimi, ModelGLM, ModelYi, ModelQwen, ModelDeepSeek, ModelSpark:
		req.Header.Set("Authorization", "Bearer "+s.Model.APIKey)
	case ModelOllama, ModelLocalAI, ModelLMStudio, ModelAnything:
		// 本地模型通常不需要认证
		if s.Model.APIKey != "" {
			req.Header.Set("Authorization", "Bearer "+s.Model.APIKey)
		}
	}
}

// parseResponse 解析响应
func (s *LLMService) parseResponse(respBody []byte) (string, error) {
	var resp ChatResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("无返回结果")
	}

	return resp.Choices[0].Message.Content, nil
}

// CheckLocalHealth 检查本地模型服务健康状态
func (s *LLMService) CheckLocalHealth(baseURL string) (bool, string) {
	client := &http.Client{Timeout: 5 * time.Second}
	
	resp, err := client.Get(baseURL + "/models")
	if err != nil {
		return false, "服务不可达"
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == 200 {
		return true, "运行正常"
	}
	return false, "服务异常"
}

// InsightRequest 洞察请求
type InsightRequest struct {
	Notes     []string `json:"notes"`
	TimeRange string   `json:"time_range"`
}

// InsightResponse 洞察响应
type InsightResponse struct {
	Summary    string   `json:"summary"`
	Keywords   []string `json:"keywords"`
	Categories []string `json:"categories"`
	Sentiment  string   `json:"sentiment"`
	Trends     []string `json:"trends"`
	Tips       []string `json:"tips"`
}

// GenerateInsight 生成洞察
func (s *LLMService) GenerateInsight(req InsightRequest) (*InsightResponse, error) {
	notes := strings.Join(req.Notes, "\n---\n")
	
	prompt := fmt.Sprintf(`分析以下笔记，提供洞察报告（用中文，JSON 格式）：

{
  "summary": "整体总结（1-2句话）",
  "keywords": ["关键词1", "关键词2", "关键词3"],
  "categories": ["分类1", "分类2"],
  "sentiment": "positive/negative/neutral",
  "trends": ["趋势1", "趋势2"],
  "tips": ["建议1", "建议2"]
}

笔记内容：
%s

时间范围：%s`, notes, req.TimeRange)

	messages := []ChatMessage{
		{Role: "system", Content: "你是一个笔记分析助手。请用中文回复严格的 JSON 格式。"},
		{Role: "user", Content: prompt},
	}

	result, err := s.Chat(messages)
	if err != nil {
		return nil, err
	}

	// 清理 markdown 代码块
	result = strings.TrimPrefix(result, "```json")
	result = strings.TrimPrefix(result, "```")
	result = strings.TrimSuffix(result, "```")
	result = strings.TrimSpace(result)

	var insight InsightResponse
	if err := json.Unmarshal([]byte(result), &insight); err != nil {
		return &InsightResponse{Summary: result}, nil
	}

	return &insight, nil
}

// SummarizeRequest 总结请求
type SummarizeRequest struct {
	Content string `json:"content"`
}

// SummarizeResponse 总结响应
type SummarizeResponse struct {
	Summary    string   `json:"summary"`
	Highlights []string `json:"highlights"`
	ActionItems []string `json:"action_items"`
}

// GenerateSummary 生成总结
func (s *LLMService) GenerateSummary(req SummarizeRequest) (*SummarizeResponse, error) {
	prompt := fmt.Sprintf(`请对以下内容进行总结，用中文回复严格的 JSON 格式：

{
  "summary": "内容总结（简洁）",
  "highlights": ["要点1", "要点2"],
  "action_items": ["可执行的任务或建议"]
}

内容：
%s`, req.Content)

	messages := []ChatMessage{
		{Role: "system", Content: "你是一个笔记总结助手。请用中文回复严格的 JSON 格式。"},
		{Role: "user", Content: prompt},
	}

	result, err := s.Chat(messages)
	if err != nil {
		return nil, err
	}

	// 清理 markdown 代码块
	result = strings.TrimPrefix(result, "```json")
	result = strings.TrimPrefix(result, "```")
	result = strings.TrimSuffix(result, "```")
	result = strings.TrimSpace(result)

	var summary SummarizeResponse
	if err := json.Unmarshal([]byte(result), &summary); err != nil {
		return &SummarizeResponse{Summary: result}, nil
	}

	return &summary, nil
}
