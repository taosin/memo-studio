package services

import (
	"strings"
	"unicode/utf8"
)

// KnownLocations 已知的地点列表及其变体
var KnownLocations = map[string][]string{
	"北京":     {"北京", "北京市", "beijing", "BJ", "北京城区", "北京市区"},
	"上海":     {"上海", "上海市", "shanghai", "SH"},
	"广州":     {"广州", "广东省广州市", "guangzhou", "GZ"},
	"深圳":     {"深圳", "广东省深圳市", "shenzhen", "SZ"},
	"杭州":     {"杭州", "浙江省杭州市", "hangzhou", "HZ"},
	"成都":     {"成都", "四川省成都市", "chengdu", "CD"},
	"重庆":     {"重庆", "重庆市", "chongqing", "CQ"},
	"武汉":     {"武汉", "湖北省武汉市", "wuhan", "WH"},
	"西安":     {"西安", "陕西省西安市", "xian", "xi'an", "XA"},
	"南京":     {"南京", "江苏省南京市", "nanjing", "NJ"},
	"天津":     {"天津", "天津市", "tianjin", "TJ"},
	"苏州":     {"苏州", "江苏省苏州市", "suzhou", "SZ"},
	"东京":     {"东京", "日本东京", "tokyo", "TKY"},
	"纽约":     {"纽约", "美国纽约", "new york", "NYC", "New York"},
	"伦敦":     {"伦敦", "英国伦敦", "london", "LDN"},
	"巴黎":     {"巴黎", "法国巴黎", "paris"},
	"悉尼":     {"悉尼", "澳大利亚悉尼", "sydney"},
	"洛杉矶":    {"洛杉矶", "美国洛杉矶", "los angeles", "LA"},
	"旧金山":    {"旧金山", "美国旧金山", "san francisco", "SF", "San Francisco"},
	"新加坡":    {"新加坡", "Singapore", "Singapura"},
	"香港":     {"香港", "中国香港", "hong kong", "HK"},
	"澳门":     {"澳门", "中国澳门", "macau"},
	"台北":     {"台北", "台湾台北", "taipei"},
	"哈尔滨":    {"哈尔滨", "黑龙江省哈尔滨市", "harbin"},
	"长春":     {"长春", "吉林省长春市"},
	"沈阳":     {"沈阳", "辽宁省沈阳市"},
	"石家庄":    {"石家庄", "河北省石家庄市"},
	"郑州":     {"郑州", "河南省郑州市"},
	"长沙":     {"长沙", "湖南省长沙市"},
	"南昌":     {"南昌", "江西省南昌市"},
	"合肥":     {"合肥", "安徽省合肥市"},
	"太原":     {"太原", "山西省太原市"},
	"济南":     {"济南", "山东省济南市"},
	"青岛":     {"青岛", "山东省青岛市"},
	"福州":     {"福州", "福建省福州市"},
	"厦门":     {"厦门", "福建省厦门市"},
	"南宁":     {"南宁", "广西南宁市"},
	"贵阳":     {"贵阳", "贵州省贵阳市"},
	"昆明":     {"昆明", "云南省昆明市"},
	"拉萨":     {"拉萨", "西藏拉萨"},
	"兰州":     {"兰州", "甘肃省兰州市"},
	"西宁":     {"西宁", "青海省西宁市"},
	"呼和浩特":  {"呼和浩特", "内蒙古呼和浩特"},
	"乌鲁木齐":  {"乌鲁木齐", "新疆乌鲁木齐"},
	"海口":     {"海口", "海南省海口市"},
	"三亚":     {"三亚", "海南省三亚市"},
	"珠海":     {"珠海", "广东省珠海市"},
	"东莞":     {"东莞", "广东省东莞市"},
	"佛山":     {"佛山", "广东省佛山市"},
	"宁波":     {"宁波", "浙江省宁波市"},
	"无锡":     {"无锡", "江苏省无锡市"},
	"常州":     {"常州", "江苏省常州市"},
	"大连":     {"大连", "辽宁省大连市"},
}

// ExtractLocation 从文本中提取地点
func ExtractLocation(content string) string {
	// 1. 检查已知的地点
	for canonical, variants := range KnownLocations {
		for _, variant := range variants {
			// 检查内容中是否包含该地点变体
			if strings.Contains(content, variant) {
				return canonical
			}
		}
	}

	// 2. 简单的正则匹配（检测 "在 X" 模式）
	patterns := []string{
		`在([^\s，。！？,.!?]{2,10})`,
		`在([^\s，。！？,.!?]{2,10})[的|过|时|候]`,
		`位于([^\s，。！？,.!?]{2,10})`,
		`回到([^\s，。！？,.!?]{2,10})`,
	}

	for _, pattern := range patterns {
		if match := extractByPattern(content, pattern); match != "" {
			// 检查是否是已知的变体
			for canonical, variants := range KnownLocations {
				for _, variant := range variants {
					if match == variant || strings.Contains(variant, match) || strings.Contains(match, variant) {
						return canonical
					}
				}
			}
			return match
		}
	}

	// 3. 检测 "X 年" 等时间模式，忽略
	// 4. 检测省份/城市后缀
	citySuffixes := []string{"市", "省", "区", "县", "镇", "村", "街道"}
	for _, suffix := range citySuffixes {
		// 检测 "X市" 模式
		if idx := strings.Index(content, suffix); idx > 0 {
			// 向前查找可能的地点名（2-5个字符）
			start := idx - 5
			if start < 0 {
				start = 0
			}
			candidate := strings.TrimSpace(content[start:idx])
			if len(candidate) >= 2 {
				return candidate + suffix
			}
		}
	}

	return ""
}

// extractByPattern 简单的正则提取
func extractByPattern(content, pattern string) string {
	// 简化实现：直接查找
	// 在实际应用中可以使用 regexp 包
	searchStr := strings.ReplaceAll(pattern, "\\s", " ")
	searchStr = strings.ReplaceAll(searchStr, "{2,10}", "")

	// 查找 "在" 后面的内容
	if strings.HasPrefix(pattern, "在") {
		idx := strings.Index(content, "在")
		if idx >= 0 && idx < len(content)-1 {
			// 提取 "在" 后面的2-10个非标点字符
			end := idx + 1
			count := 0
			for end < len(content) && count < 10 {
				r, size := utf8.DecodeRuneInString(content[end:])
				if size == 0 {
					break
				}
				if r == ' ' || r == '\n' || r == '，' || r == '。' ||
					r == '！' || r == '？' || r == ',' || r == '.' ||
					r == '!' || r == '?' {
					break
				}
				end += size
				count++
			}
			result := strings.TrimSpace(content[idx+1 : end])
			if len(result) >= 2 {
				return result
			}
		}
	}

	return ""
}

// LocationWithCoords 地点及其坐标
type LocationWithCoords struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// GetLocationCoords 获取地点的坐标（模拟数据）
// 在实际应用中，可以调用地理编码 API（如高德、Google Maps）
func GetLocationCoords(locationName string) *LocationWithCoords {
	// 中国主要城市坐标（示例数据）
	coords := map[string][2]float64{
		"北京":     {39.9042, 116.4074},
		"上海":     {31.2304, 121.4737},
		"广州":     {23.1291, 113.2644},
		"深圳":     {22.5431, 114.0579},
		"杭州":     {30.2741, 120.1551},
		"成都":     {30.5728, 104.0668},
		"重庆":     {29.4316, 106.9123},
		"武汉":     {30.5928, 114.3055},
		"西安":     {34.3416, 108.9398},
		"南京":     {32.0603, 118.7969},
		"天津":     {39.1256, 117.1909},
		"东京":     {35.6762, 139.6503},
		"纽约":     {40.7128, -74.0060},
		"伦敦":     {51.5074, -0.1278},
		"巴黎":     {48.8566, 2.3522},
		"悉尼":     {-33.8688, 151.2093},
		"洛杉矶":    {34.0522, -118.2437},
		"旧金山":    {37.7749, -122.4194},
		"新加坡":    {1.3521, 103.8198},
		"香港":     {22.3193, 114.1694},
		"台北":     {25.0330, 121.5654},
	}

	if coord, ok := coords[locationName]; ok {
		return &LocationWithCoords{
			Name:      locationName,
			Latitude:  coord[0],
			Longitude: coord[1],
		}
	}

	return nil
}

// DetectAndExtractLocation 检测并提取地点，返回标准名称和坐标
func DetectAndExtractLocation(content string) *LocationWithCoords {
	location := ExtractLocation(content)
	if location == "" {
		return nil
	}

	coords := GetLocationCoords(location)
	if coords != nil {
		return coords
	}

	// 如果没有坐标数据，返回地点名称
	return &LocationWithCoords{
		Name:      location,
		Latitude:  0,
		Longitude: 0,
	}
}

// DetectLocationInNotes 检测多条笔记的地点
func DetectLocationInNotes(notes []string) map[int]string {
	results := make(map[int]string)
	for i, note := range notes {
		if loc := ExtractLocation(note); loc != "" {
			results[i] = loc
		}
	}
	return results
}
