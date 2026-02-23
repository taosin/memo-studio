package handlers

import (
	"net/http"
	"os"
	"time"

	"memo-studio/backend/services"

	"github.com/gin-gonic/gin"
)

// InsightType 洞察视角类型
type InsightType string

const (
	InsightOverview    InsightType = "overview"     // 概览
	InsightTime       InsightType = "time"          // 时间视角
	InsightTopic      InsightType = "topic"         // 主题视角
	InsightSentiment  InsightType = "sentiment"      // 情感视角
	InsightAction     InsightType = "action"         // 行动视角
	InsightConnection InsightType = "connection"      // 关联视角
	InsightFrequency  InsightType = "frequency"      // 频率视角
	InsightAll        InsightType = "all"           // 全部视角
)

// InsightRequest 洞察请求
type InsightRequest struct {
	Notes       []string     `json:"notes"`
	TimeRange   string       `json:"time_range"`
	Perspectives []InsightType `json:"perspectives"`
}

// InsightResponse 洞察响应
type InsightResponse struct {
	Summary      string             `json:"summary"`
	Perspectives []PerspectiveInsight `json:"perspectives"`
	Highlights   []string          `json:"highlights"`
	ActionItems  []string          `json:"action_items"`
	UpdateTime   string            `json:"update_time"`
}

// PerspectiveInsight 单个视角的洞察
type PerspectiveInsight struct {
	Type      InsightType `json:"type"`
	Name     string     `json:"name"`
	Summary  string     `json:"summary"`
	Details  []DetailItem `json:"details"`
	Highlights []string  `json:"highlights"`
	Score    int        `json:"score"`
}

// DetailItem 详细分析项
type DetailItem struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Icon    string `json:"icon"`
	Count   int    `json:"count"`
}

// SummarizeResponse 总结响应
type SummarizeResponse struct {
	Summary     string   `json:"summary"`
	Highlights  []string `json:"highlights"`
	ActionItems []string `json:"action_items"`
}

// GetInsight 获取笔记洞察（多视角）
// POST /api/insights
func GetInsight(c *gin.Context) {
	var req struct {
		Notes       []string     `json:"notes"`
		TimeRange   string       `json:"time_range"`
		Perspectives []InsightType `json:"perspectives"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误"})
		return
	}

	if req.TimeRange == "" {
		req.TimeRange = "30d"
	}
	if len(req.Perspectives) == 0 {
		req.Perspectives = []InsightType{InsightAll}
	}

	// 检查是否配置了 LLM
	hasAPIKey := os.Getenv("OPENAI_API_KEY") != "" ||
		os.Getenv("LLM_API_KEY") != "" ||
		os.Getenv("ANTHROPIC_API_KEY") != "" ||
		os.Getenv("DEEPSEEK_API_KEY") != "" ||
		os.Getenv("ZHIPU_API_KEY") != ""

	var response InsightResponse

	if hasAPIKey && len(req.Notes) > 0 {
		// 使用 LLM 生成洞察
		llmService := services.NewLLMService()
		aiInsight, err := llmService.GenerateInsight(services.InsightRequest{
			Notes:     req.Notes,
			TimeRange: req.TimeRange,
		})
		
		if err == nil {
			// 转换为多视角格式
			response = convertToMultiPerspective(aiInsight, req)
		} else {
			response = generateBasicInsight(req.Notes, req.TimeRange)
		}
	} else {
		// 使用基础分析
		response = generateBasicInsight(req.Notes, req.TimeRange)
	}

	response.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	c.JSON(http.StatusOK, response)
}

// GetInsightByType 获取特定视角的洞察
// POST /api/insights/:type
func GetInsightByType(c *gin.Context) {
	insightType := InsightType(c.Param("type"))

	var req struct {
		Notes     []string `json:"notes"`
		TimeRange string   `json:"time_range"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误"})
		return
	}

	if req.TimeRange == "" {
		req.TimeRange = "30d"
	}

	response := generatePerspective(insightType, req.Notes, req.TimeRange)
	c.JSON(http.StatusOK, response)
}

// CompareInsights 对比分析
// POST /api/insights/compare
func CompareInsights(c *gin.Context) {
	var req struct {
		Notes1 []string `json:"notes1"`
		Notes2 []string `json:"notes2"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误"})
		return
	}

	insight1 := generateBasicInsight(req.Notes1, "period1")
	insight2 := generateBasicInsight(req.Notes2, "period2")

	c.JSON(http.StatusOK, gin.H{
		"period1":  insight1,
		"period2":  insight2,
		"changes":  generateChanges(insight1, insight2),
	})
}

// SummarizeNote 总结单条笔记
// POST /api/summarize
func SummarizeNote(c *gin.Context) {
	var req struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误"})
		return
	}

	if req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "内容不能为空"})
		return
	}

	// 检查是否配置了 LLM
	hasAPIKey := os.Getenv("OPENAI_API_KEY") != "" ||
		os.Getenv("LLM_API_KEY") != "" ||
		os.Getenv("ANTHROPIC_API_KEY") != ""

	if hasAPIKey {
		llmService := services.NewLLMService()
		summary, err := llmService.GenerateSummary(services.SummarizeRequest{
			Content: req.Content,
		})
		
		if err == nil {
			c.JSON(http.StatusOK, summary)
			return
		}
	}

	// 返回基础总结
	c.JSON(http.StatusOK, SummarizeResponse{
		Summary:    "（请配置 LLM_API_KEY 启用 AI 总结）",
		Highlights: []string{},
		ActionItems: []string{},
	})
}

// BatchSummarize 批量总结
// POST /api/summarize/batch
func BatchSummarize(c *gin.Context) {
	var req struct {
		Notes []string `json:"notes"`
		Limit int      `json:"limit"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误"})
		return
	}

	if len(req.Notes) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "笔记列表不能为空"})
		return
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}

	hasAPIKey := os.Getenv("OPENAI_API_KEY") != "" ||
		os.Getenv("LLM_API_KEY") != ""

	results := make([]SummarizeResponse, 0, len(req.Notes))
	llmService := services.NewLLMService()

	for i, note := range req.Notes {
		if i >= req.Limit {
			break
		}

		var summary SummarizeResponse
		if hasAPIKey {
			s, err := llmService.GenerateSummary(services.SummarizeRequest{Content: note})
			if err == nil {
				summary = SummarizeResponse{
					Summary:     s.Summary,
					Highlights:  s.Highlights,
					ActionItems: s.ActionItems,
				}
			} else {
				summary = SummarizeResponse{Summary: truncate(note, 100)}
			}
		} else {
			summary = SummarizeResponse{Summary: truncate(note, 100)}
		}
		results = append(results, summary)
	}

	c.JSON(http.StatusOK, gin.H{
		"total":   len(req.Notes),
		"limited": len(results),
		"results": results,
	})
}

// ========== 辅助函数 ==========

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func convertToMultiPerspective(aiInsight *services.InsightResponse, req InsightRequest) InsightResponse {
	return InsightResponse{
		Summary: aiInsight.Summary,
		Perspectives: []PerspectiveInsight{
			{
				Type:      InsightTopic,
				Name:     "🏷️ 主题视角",
				Summary:  "基于 AI 分析",
				Highlights: aiInsight.Keywords,
				Score:    85,
			},
			{
				Type:      InsightSentiment,
				Name:     "💭 情感视角",
				Summary:  aiInsight.Sentiment,
				Score:    80,
			},
		},
		Highlights:  aiInsight.Tips,
		ActionItems: aiInsight.Tips,
	}
}

func generateBasicInsight(notes []string, timeRange string) InsightResponse {
	count := len(notes)
	charCount := 0
	for _, n := range notes {
		charCount += len(n)
	}

	return InsightResponse{
		Summary: formatTimeRange(timeRange) + "共记录 " + itoa(count) + " 条笔记",
		Perspectives: []PerspectiveInsight{
			generatePerspective(InsightOverview, notes, timeRange),
			generatePerspective(InsightTopic, notes, timeRange),
			generatePerspective(InsightSentiment, notes, timeRange),
		},
		Highlights:  []string{"继续保持记录习惯"},
		ActionItems: []string{},
	}
}

func generatePerspective(pType InsightType, notes []string, timeRange string) PerspectiveInsight {
	perspective := PerspectiveInsight{
		Type:      pType,
		Highlights: []string{},
		Score:     50,
	}

	switch pType {
	case InsightOverview:
		perspective.Name = "📊 概览"
		perspective.Summary = "共 " + itoa(len(notes)) + " 条笔记"
		perspective.Details = []DetailItem{
			{Title: "笔记数", Content: itoa(len(notes)), Icon: "📝"},
			{Title: "总字数", Content: itoa(len(notes)) + " 字", Icon: "📏"},
		}
		perspective.Score = 70

	case InsightTopic:
		perspective.Name = "🏷️ 主题视角"
		topicStats := analyzeTopics(notes)
		perspective.Summary = topicStats.Summary
		perspective.Details = topicStats.Details
		perspective.Score = topicStats.Score

	case InsightSentiment:
		perspective.Name = "💭 情感视角"
		sentimentStats := analyzeSentiment(notes)
		perspective.Summary = sentimentStats.Summary
		perspective.Highlights = sentimentStats.Highlights
		perspective.Score = sentimentStats.Score

	case InsightAction:
		perspective.Name = "✅ 行动视角"
		actionStats := analyzeActions(notes)
		perspective.Summary = actionStats.Summary
		perspective.Details = actionStats.Details
		perspective.Score = actionStats.Score

	default:
		perspective.Name = "📊 综合"
		perspective.Summary = "记录良好"
	}

	return perspective
}

func formatTimeRange(tr string) string {
	switch tr {
	case "7d":
		return "最近 7 天"
	case "30d":
		return "最近 30 天"
	case "90d":
		return "最近 3 个月"
	default:
		return ""
	}
}

func itoa(n int) string {
	return string(rune('0'+n/1000%10)) + string(rune('0'+n/100%10)) + string(rune('0'+n/10%10)) + string(rune('0'+n%10))
}

type statsResult struct {
	Summary   string
	Details   []DetailItem
	Score     int
	Highlights []string
}

func analyzeTopics(notes []string) statsResult {
	topics := map[string][]string{
		"💻 工作": {"工作", "项目", "任务", "会议"},
		"📚 学习": {"学习", "读书", "课程", "知识"},
		"🏃 健康": {"健康", "运动", "锻炼"},
		"💰 财务": {"钱", "消费", "收入", "理财"},
	}

	countMap := make(map[string]int)
	for _, note := range notes {
		for topic, keywords := range topics {
			for _, kw := range keywords {
				if contains(note, kw) {
					countMap[topic]++
				}
			}
		}
	}

	var details []DetailItem
	var maxCount int
	var topTopic string
	for topic, count := range countMap {
		details = append(details, DetailItem{
			Title:  topic,
			Content: itoa(count) + " 条",
			Icon:   string([]byte(topic)[0]),
		})
		if count > maxCount {
			maxCount = count
			topTopic = topic
		}
	}

	summary := "关注领域较广"
	score := 50
	if topTopic != "" {
		summary = "最关注 " + topTopic[2:] + " 方面"
		score = 70
	}

	return statsResult{Summary: summary, Details: details, Score: score}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func analyzeSentiment(notes []string) statsResult {
	positive := []string{"开心", "高兴", "成功", "收获", "进步"}
	negative := []string{"难过", "焦虑", "失败", "困难"}

	posCount, negCount := 0, 0
	for _, note := range notes {
		for _, w := range positive {
			if contains(note, w) {
				posCount++
				break
			}
		}
		for _, w := range negative {
			if contains(note, w) {
				negCount++
				break
			}
		}
	}

	summary := "情绪平稳"
	score := 50
	highlights := []string{}

	if posCount > negCount {
		summary = "😊 整体积极"
		score = 75
		highlights = append(highlights, "正面情绪占主导")
	} else if negCount > posCount {
		summary = "😔 有些负面情绪"
		score = 40
		highlights = append(highlights, "建议适当放松")
	}

	return statsResult{Summary: summary, Score: score, Highlights: highlights}
}

func analyzeActions(notes []string) statsResult {
	todoWords := []string{"待办", "计划", "要", "应该"}
	doneWords := []string{"完成", "解决", "搞定"}

	todoCount, doneCount := 0, 0
	for _, note := range notes {
		for _, w := range todoWords {
			if contains(note, w) {
				todoCount++
				break
			}
		}
		for _, w := range doneWords {
			if contains(note, w) {
				doneCount++
				break
			}
		}
	}

	summary := "有一定行动记录"
	score := 50
	if todoCount > 0 {
		rate := doneCount * 100 / todoCount
		summary = "完成率 " + itoa(rate) + "%"
		if rate > 70 {
			score = 85
		} else if rate > 40 {
			score = 60
		}
	}

	return statsResult{
		Summary: summary,
		Details: []DetailItem{
			{Title: "待办", Content: itoa(todoCount), Icon: "📋"},
			{Title: "完成", Content: itoa(doneCount), Icon: "✅"},
		},
		Score: score,
	}
}

func generateChanges(insight1, insight2 InsightResponse) []map[string]string {
	return []map[string]string{
		{"category": "记录数量", "before": insight1.Summary, "after": insight2.Summary},
	}
}
