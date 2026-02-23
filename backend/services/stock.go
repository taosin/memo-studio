package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// StockInfo 股票信息
type StockInfo struct {
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	Market       string  `json:"market"`       // 深圳/上海
	Price        float64 `json:"price"`        // 当前价格
	Change       float64 `json:"change"`       // 涨跌额
	ChangePercent float64 `json:"change_percent"` // 涨跌幅
	Open         float64 `json:"open"`         // 开盘价
	PreClose     float64 `json:"pre_close"`    // 昨收价
	High         float64 `json:"high"`         // 最高价
	Low          float64 `json:"low"`          // 最低价
	Volume       int64   `json:"volume"`       // 成交量（手）
	Turnover     float64 `json:"turnover"`     // 成交额（万）
	PE           float64 `json:"pe"`           // 市盈率
	PB           float64 `json:"pb"`           // 市净率
	Dividend     float64 `json:"dividend"`     // 股息率
	MarketCap    float64 `json:"market_cap"`  // 总市值（亿）
	MarketCapStr string  `json:"market_cap_str"` // 格式化市值
	UpdateTime   string  `json:"update_time"`
}

// StockFundFlow 资金流向
type StockFundFlow struct {
	Code           string  `json:"code"`
	MainNetInflow  float64 `json:"main_net_inflow"`  // 主力净流入
	MainNetInflowRate float64 `json:"main_net_inflow_rate"` // 主力净流入占比
	SuperNetInflow float64 `json:"super_net_inflow"` // 超大单净流入
	LargeNetInflow float64 `json:"large_net_inflow"` // 大单净流入
	MediumNetInflow float64 `json:"medium_net_inflow"` // 中单净流入
	SmallNetInflow float64 `json:"small_net_inflow"` // 小单净流入
	UpdateTime     string  `json:"update_time"`
}

// StockHolder 股东信息
type StockHolder struct {
	HolderName string  `json:"holder_name"` // 股东名称
	HolderType string  `json:"holder_type"` // 股东类型
	Shares     float64 `json:"shares"`      // 持股数（万股）
	Ratio      float64 `json:"ratio"`       // 持股比例
	Change     float64 `json:"change"`      // 持股变动
	ChangeRatio float64 `json:"change_ratio"` // 变动比例
	ReportDate string  `json:"report_date"` // 公告日期
}

// StockFinance 财务指标
type StockFinance struct {
	Code          string  `json:"code"`
	Revenue       float64 `json:"revenue"`       // 营业收入（亿）
	RevenueYoY    float64 `json:"revenue_yoy"`  // 营收同比
	Profit        float64 `json:"profit"`        // 净利润（亿）
	ProfitYoY     float64 `json:"profit_yoy"`   // 净利润同比
	EPS           float64 `json:"eps"`           // 每股收益
	ROE           float64 `json:"roe"`           // 净资产收益率
	DebtRatio     float64 `json:"debt_ratio"`   // 资产负债率
	CashFlow      float64 `json:"cash_flow"`    // 经营现金流（亿）
	GrossMargin   float64 `json:"gross_margin"` // 毛利率
	NetMargin     float64 `json:"net_margin"`   // 净利率
	ReportDate    string  `json:"report_date"`  // 报告期
}

// StockSearch 股票搜索结果
type StockSearch struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Market string `json:"market"`
}

// GetStockInfo 获取股票实时信息
func GetStockInfo(stockCode string) (*StockInfo, error) {
	// 转换股票代码格式
	code := formatStockCode(stockCode)
	if code == "" {
		return nil, fmt.Errorf("无效的股票代码")
	}

	// 使用新浪财经 API
	apiURL := fmt.Sprintf("https://hq.sinajs.cn/list=%s", code)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("Referer", "http://finance.sina.com.cn")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取数据失败: %v", err)
	}

	// 解析新浪财经数据格式
	// 格式: var hq_str_sh600519="贵州茅台,1680.50,1665.00,1685.00,1690.00,1660.00,1670.00,23456789,3.89亿,0.14%,28.5,1680.50,2026-02-17 10:30:00";
	stockInfo, err := parseSinaResponse(string(body), code)
	if err != nil {
		return nil, err
	}

	return stockInfo, nil
}

// formatStockCode 格式化股票代码
func formatStockCode(code string) string {
	code = strings.TrimSpace(code)

	// 如果已经是正确格式，直接返回
	if strings.HasPrefix(code, "sh") || strings.HasPrefix(code, "sz") {
		return code
	}

	// 判断市场
	if len(code) == 6 {
		switch code[0] {
		case '0', '3':
			return "sz" + code // 深圳
		case '5', '6':
			return "sh" + code // 上海
		}
	}

	return ""
}

// GetStockFundFlow 获取股票资金流向
func GetStockFundFlow(stockCode string) (*StockFundFlow, error) {
	code := formatStockCode(stockCode)
	if code == "" {
		return nil, fmt.Errorf("无效的股票代码")
	}

	// 使用东方财富资金流向 API
	apiURL := fmt.Sprintf("http://push2.eastmoney.com/api/qt/stock/get?secid=%s&fields=f43,f50,f51,f52,f53,f54,f55,f57,f58,f59,f60,f61,f62,f63,f64,f65,f66,f67,f68,f69,f70",
		getSecID(code))

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取数据失败: %v", err)
	}

	var result struct {
		Data struct {
			F43 float64 `json:"f43"` // 主力净流入
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析数据失败: %v", err)
	}

	return &StockFundFlow{
		Code:          stockCode,
		MainNetInflow: result.Data.F43,
		UpdateTime:    time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

// getSecID 获取东方财富 secid
func getSecID(code string) string {
	if strings.HasPrefix(code, "sh") {
		return "1." + code[2:]
	}
	return "0." + code[2:]
}

// GetStockFinance 获取股票财务数据
func GetStockFinance(stockCode string) (*StockFinance, error) {
	code := formatStockCode(stockCode)
	if code == "" {
		return nil, fmt.Errorf("无效的股票代码")
	}

	// 使用新浪财经财务数据
	apiURL := fmt.Sprintf("https://hq.sinajs.cn/list=%s", code)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("Referer", "http://finance.sina.com.cn")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 解析财务数据（简化版）
	return &StockFinance{
		Code:       stockCode,
		ReportDate: time.Now().Format("2006-01-02"),
	}, nil
}

// GetStockHolders 获取主要股东信息
func GetStockHolders(stockCode string) ([]StockHolder, error) {
	code := formatStockCode(stockCode)
	if code == "" {
		return nil, fmt.Errorf("无效的股票代码")
	}

	// 返回空数据（实际需要从东方财富等API获取）
	return []StockHolder{}, nil
}

// parseSinaResponse 解析新浪财经响应
func parseSinaResponse(response, code string) (*StockInfo, error) {
	// 提取数据部分
	start := strings.Index(response, "=")
	if start == -1 {
		return nil, fmt.Errorf("数据格式错误")
	}

	dataStr := response[start+2 : len(response)-2]
	parts := strings.Split(dataStr, ",")

	if len(parts) < 32 {
		return nil, fmt.Errorf("数据不完整")
	}

	// 解析数据
	stockInfo := &StockInfo{
		Code:       code,
		Name:       parts[0],
		Open:       parseFloat(parts[1]),
		PreClose:   parseFloat(parts[2]),
		Price:      parseFloat(parts[3]), // 当前价
		High:       parseFloat(parts[4]),
		Low:        parseFloat(parts[5]),
		Volume:     parseInt64(parts[8]),
		Turnover:   parseFloat(parts[9]),
		Change:     parseFloat(parts[31]),
		UpdateTime: parts[30],
	}

	// 计算涨跌幅
	if stockInfo.PreClose > 0 {
		stockInfo.ChangePercent = (stockInfo.Change / stockInfo.PreClose) * 100
	}

	// 设置市场
	if strings.HasPrefix(code, "sh") {
		stockInfo.Market = "上海"
	} else {
		stockInfo.Market = "深圳"
	}

	// 估算 PE 和市值（需要额外 API）
	stockInfo.PE = 0
	stockInfo.MarketCap = 0

	return stockInfo, nil
}

// GetStockList 获取股票列表
func GetStockList(keyword string) ([]StockSearch, error) {
	// 使用同花顺股票API
	apiURL := fmt.Sprintf("http://search.tianyancha.com/api/v4/stock/search?keyword=%s",
		url.QueryEscape(keyword))

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var result struct {
		Data []struct {
			Code   string `json:"code"`
			Name   string `json:"name"`
			Market string `json:"market"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		// 如果 API 失败，返回空列表
		return []StockSearch{}, nil
	}

	stockList := make([]StockSearch, 0, len(result.Data))
	for _, s := range result.Data {
		stockList = append(stockList, StockSearch{
			Code:   s.Code,
			Name:   s.Name,
			Market: s.Market,
		})
	}

	return stockList, nil
}

// parseFloat 解析浮点数
func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}

// parseInt64 解析整数
func parseInt64(s string) int64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return int64(f)
}

// GetStockHistory 获取股票历史数据
func GetStockHistory(stockCode string, days int) ([]StockHistory, error) {
	code := formatStockCode(stockCode)
	if code == "" {
		return nil, fmt.Errorf("无效的股票代码")
	}

	// 使用新浪财经历史数据 API
	// 格式: https://finance.sina.com.cn/realstock/company/sh600519/nc.shtml
	apiURL := fmt.Sprintf("https://quotes.sina.cn/cn/api/json.php/KL_MarketDataService.getKLineData?symbol=%s&scale=240&ma=no&datalen=%d",
		code, days)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析 JSON
	var data []struct {
		D string  `json:"d"` // 日期
		O float64 `json:"o"` // 开盘
		C float64 `json:"c"` // 收盘
		H float64 `json:"h"` // 最高
		L float64 `json:"l"` // 最低
		V int64   `json:"v"` // 成交量
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	history := make([]StockHistory, 0, len(data))
	for _, d := range data {
		history = append(history, StockHistory{
			Date:   d.D,
			Open:   d.O,
			Close:  d.C,
			High:   d.H,
			Low:    d.L,
			Volume: d.V,
		})
	}

	return history, nil
}

// StockHistory 股票历史数据
type StockHistory struct {
	Date   string  `json:"date"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Volume int64   `json:"volume"`
}

// StockAnalysis 股票分析结果
type StockAnalysis struct {
	Summary     string        `json:"summary"`
	Signals     []AnalysisSignal `json:"signals"`
	Suggestion  string        `json:"suggestion"`
	Risks       []string     `json:"risks"`
	Tips        []string     `json:"tips"`
}

// AnalysisSignal 分析信号
type AnalysisSignal struct {
	Type string `json:"type"` // technical, volume, valuation
	Icon string `json:"icon"`
	Text string `json:"text"`
}

// AnalyzeStock 分析股票
func AnalyzeStock(stock *StockInfo) *StockAnalysis {
	analysis := &StockAnalysis{
		Summary: fmt.Sprintf("%s（%s）今日%s%.2f元（%.2f%%），当前价格¥%.2f",
			stock.Name,
			stock.Code,
			getChangeText(stock.Change),
			stock.Change,
			stock.ChangePercent,
			stock.Price,
		),
		Signals: make([]AnalysisSignal, 0),
		Risks:   []string{},
		Tips:    []string{},
	}

	// 技术面分析
	if stock.Change > 0 {
		analysis.Signals = append(analysis.Signals, AnalysisSignal{
			Type: "technical",
			Icon: "📈",
			Text: "价格上涨，技术面偏强",
		})
	} else if stock.Change < 0 {
		analysis.Signals = append(analysis.Signals, AnalysisSignal{
			Type: "technical",
			Icon: "📉",
			Text: "价格下跌，需关注支撑位",
		})
	}

	// 成交量分析
	if stock.Volume > 10000000 { // 1000万手以上
		analysis.Signals = append(analysis.Signals, AnalysisSignal{
			Type: "volume",
			Icon: "📊",
			Text: fmt.Sprintf("成交量 %d 万手，较为活跃", stock.Volume/10000),
		})
	}

	// 估值分析
	if stock.PE > 0 {
		if stock.PE > 50 {
			analysis.Signals = append(analysis.Signals, AnalysisSignal{
				Type: "valuation",
				Icon: "💰",
				Text: fmt.Sprintf("市盈率 %.2f 倍，估值偏高", stock.PE),
			})
		} else if stock.PE > 0 && stock.PE < 20 {
			analysis.Signals = append(analysis.Signals, AnalysisSignal{
				Type: "valuation",
				Icon: "💵",
				Text: fmt.Sprintf("市盈率 %.2f 倍，估值合理", stock.PE),
			})
		}
	}

	// 建议
	if stock.Change > 3 {
		analysis.Suggestion = "涨幅较大，建议减仓或观望"
	} else if stock.Change < -3 {
		analysis.Suggestion = "跌幅较大，关注支撑位，可适当补仓"
	} else if stock.Change > 0 {
		analysis.Suggestion = "可持有，关注上方压力位"
	} else {
		analysis.Suggestion = "建议观望，注意止损"
	}

	// 风险提示
	analysis.Risks = []string{
		"市场整体回调风险",
		"行业政策变化影响",
		"公司业绩不及预期",
		"大盘系统性风险",
	}

	// 投资建议
	analysis.Tips = []string{
		"分散投资，不要满仓一只股票",
		"设置止损位，控制风险",
		"关注公司基本面变化",
		"保持长期投资心态",
		"不要追涨杀跌",
	}

	return analysis
}

// getChangeText 获取涨跌描述
func getChangeText(change float64) string {
	switch {
	case change > 0:
		return "上涨"
	case change < 0:
		return "下跌"
	default:
		return "持平"
	}
}

// GetHotStocks 获取热门股票列表
func GetHotStocks() []StockSearch {
	return []StockSearch{
		{Code: "000001", Name: "平安银行", Market: "深圳"},
		{Code: "600519", Name: "贵州茅台", Market: "上海"},
		{Code: "600036", Name: "招商银行", Market: "上海"},
		{Code: "000002", Name: "万 科Ａ", Market: "深圳"},
		{Code: "601398", Name: "工商银行", Market: "上海"},
		{Code: "601857", Name: "中国石油", Market: "上海"},
		{Code: "601988", Name: "中国银行", Market: "上海"},
		{Code: "600000", Name: "浦发银行", Market: "上海"},
		{Code: "000725", Name: "京东方A", Market: "深圳"},
		{Code: "002594", Name: "比亚迪", Market: "深圳"},
	}
}
