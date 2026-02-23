package handlers

import (
	"net/http"
	"strconv"

	"memo-studio/backend/services"

	"github.com/gin-gonic/gin"
)

// GetStockInfo 获取股票信息
// GET /api/stocks/:code
func GetStockInfo(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供股票代码"})
		return
	}

	stock, err := services.GetStockInfo(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stock": stock,
	})
}

// SearchStocks 搜索股票
// GET /api/stocks/search?q=关键词
func SearchStocks(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供搜索关键词"})
		return
	}

	results, err := services.GetStockList(query)
	if err != nil {
		// 如果搜索失败，返回热门股票
		results = services.GetHotStocks()
	}

	c.JSON(http.StatusOK, gin.H{
		"results": results,
	})
}

// GetHotStocks 获取热门股票
// GET /api/stocks/hot
func GetHotStocks(c *gin.Context) {
	stocks := services.GetHotStocks()
	c.JSON(http.StatusOK, gin.H{
		"stocks": stocks,
	})
}

// GetStockHistory 获取股票历史数据
// GET /api/stocks/:code/history?days=30
func GetStockHistory(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供股票代码"})
		return
	}

	days := 30
	if d := c.Query("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil {
			days = parsed
		}
	}

	history, err := services.GetStockHistory(code, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"days":    days,
		"history": history,
	})
}

// AnalyzeStock 分析股票
// POST /api/stocks/analyze
func AnalyzeStock(c *gin.Context) {
	var req struct {
		Code string `json:"code"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误"})
		return
	}

	if req.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供股票代码"})
		return
	}

	// 获取股票信息
	stock, err := services.GetStockInfo(req.Code)
	if err != nil {
		// 如果获取失败，使用模拟数据
		stock = &services.StockInfo{
			Code:         req.Code,
			Name:         "模拟股票",
			Market:       "上海",
			Price:        100.00,
			Change:       2.50,
			ChangePercent: 2.56,
			Volume:       5000000,
			PE:          25.0,
		}
	}

	// 分析股票
	analysis := services.AnalyzeStock(stock)

	c.JSON(http.StatusOK, gin.H{
		"stock":    stock,
		"analysis": analysis,
	})
}

// StockResponse 股票 API 统一响应格式
type StockResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string     `json:"error,omitempty"`
}
