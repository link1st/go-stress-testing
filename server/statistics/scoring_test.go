// Package statistics 评分系统测试
package statistics

import (
	"testing"
)

// TestCalculateScore_Perfect 测试满分场景
func TestCalculateScore_Perfect(t *testing.T) {
	data := &ReportData{
		URL:           "http://example.com",
		Concurrency:   10,
		TotalRequests: 1000,
		SuccessNum:    1000,
		FailureNum:    0,
		QPS:           1000.0,
		AvgTime:       50.0,
		MinTime:       10.0,
		MaxTime:       100.0,
		TP90:          60.0,
		TP95:          70.0,
		TP99:          75.0,
		ErrorCodeMap:  map[int]int{200: 1000},
	}

	result := CalculateScore(data)

	if result.TotalScore < 90 {
		t.Errorf("完美数据总分应该>=90, 实际: %d", result.TotalScore)
	}
	if result.Grade != "A" {
		t.Errorf("完美数据应该评级A, 实际: %s", result.Grade)
	}
	if result.SuccessRateScore != 30 {
		t.Errorf("成功率100%%应该得30分, 实际: %d", result.SuccessRateScore)
	}
}

// TestCalculateScore_Poor 测试差评场景
func TestCalculateScore_Poor(t *testing.T) {
	data := &ReportData{
		URL:           "http://example.com",
		Concurrency:   10,
		TotalRequests: 1000,
		SuccessNum:    800,
		FailureNum:    200,
		QPS:           5.0,
		AvgTime:       2000.0,
		MinTime:       100.0,
		MaxTime:       10000.0,
		TP90:          5000.0,
		TP95:          8000.0,
		TP99:          9000.0,
		ErrorCodeMap:  map[int]int{200: 800, 500: 200},
	}

	result := CalculateScore(data)

	if result.TotalScore > 50 {
		t.Errorf("差数据总分应该<=50, 实际: %d", result.TotalScore)
	}
	if result.Grade == "A" || result.Grade == "B" {
		t.Errorf("差数据不应该评级A或B, 实际: %s", result.Grade)
	}
	if len(result.Suggestions) == 0 {
		t.Error("差数据应该有优化建议")
	}
}

// TestCalculateScore_NilData 测试空数据
func TestCalculateScore_NilData(t *testing.T) {
	result := CalculateScore(nil)

	if result.TotalScore != 0 {
		t.Errorf("空数据应该得0分, 实际: %d", result.TotalScore)
	}
	if result.Grade != "F" {
		t.Errorf("空数据应该评级F, 实际: %s", result.Grade)
	}
}

// TestCalculateGrade 测试评级计算
func TestCalculateGrade(t *testing.T) {
	tests := []struct {
		score int
		grade string
	}{
		{100, "A"},
		{95, "A"},
		{90, "A"},
		{89, "B"},
		{80, "B"},
		{79, "C"},
		{70, "C"},
		{69, "D"},
		{60, "D"},
		{59, "F"},
		{0, "F"},
	}

	for _, tt := range tests {
		grade := calculateGrade(tt.score)
		if grade != tt.grade {
			t.Errorf("分数%d应该评级%s, 实际: %s", tt.score, tt.grade, grade)
		}
	}
}

// TestCalculateSuccessRateScore 测试成功率评分
func TestCalculateSuccessRateScore(t *testing.T) {
	tests := []struct {
		name     string
		success  uint64
		total    uint64
		minScore int
		maxScore int
	}{
		{"100%成功率", 1000, 1000, 30, 30},
		{"99%成功率", 990, 1000, 24, 24},
		{"95%成功率", 950, 1000, 16, 16},
		{"90%成功率", 900, 1000, 10, 10},
		{"80%成功率", 800, 1000, 5, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &ReportData{
				TotalRequests: tt.total,
				SuccessNum:    tt.success,
				ErrorCodeMap:  map[int]int{200: int(tt.success)},
			}
			result := &ScoreResult{Details: make(map[string]string)}
			score := calculateSuccessRateScore(data, result)
			if score < tt.minScore || score > tt.maxScore {
				t.Errorf("%s: 期望分数在[%d,%d], 实际: %d", tt.name, tt.minScore, tt.maxScore, score)
			}
		})
	}
}

// TestCalculateQPSScore 测试QPS评分
func TestCalculateQPSScore(t *testing.T) {
	tests := []struct {
		name        string
		qps         float64
		concurrency uint64
		minScore    int
	}{
		{"高QPS", 1000, 10, 20},
		{"中QPS", 100, 10, 10},
		{"低QPS", 10, 10, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &ReportData{
				QPS:         tt.qps,
				Concurrency: tt.concurrency,
			}
			result := &ScoreResult{Details: make(map[string]string)}
			score := calculateQPSScore(data, result)
			if score < tt.minScore {
				t.Errorf("%s: 期望分数>=%d, 实际: %d", tt.name, tt.minScore, score)
			}
		})
	}
}

// TestCalculateAvgTimeScore 测试平均响应时间评分
func TestCalculateAvgTimeScore(t *testing.T) {
	tests := []struct {
		avgTime  float64
		minScore int
	}{
		{30, 20},  // 极快
		{80, 19},  // 很快
		{150, 18}, // 快速
		{250, 17}, // 较快
		{400, 15}, // 良好
		{700, 13}, // 一般
		{900, 11}, // 可接受
		{1200, 9}, // 略慢
		{1800, 7}, // 较慢
		{2500, 5}, // 慢
		{4000, 3}, // 很慢
		{6000, 0}, // 极慢
	}

	for _, tt := range tests {
		data := &ReportData{AvgTime: tt.avgTime}
		result := &ScoreResult{Details: make(map[string]string)}
		score := calculateAvgTimeScore(data, result)
		if score < tt.minScore {
			t.Errorf("平均响应%.0fms: 期望分数>=%d, 实际: %d", tt.avgTime, tt.minScore, score)
		}
	}
}

// TestCalculateTP99Score 测试TP99稳定性评分
func TestCalculateTP99Score(t *testing.T) {
	tests := []struct {
		avgTime  float64
		tp99     float64
		minScore int
	}{
		{100, 120, 12}, // 非常稳定 (1.2x)
		{100, 180, 10}, // 稳定 (1.8x)
		{100, 250, 5},  // 一般 (2.5x)
		{100, 400, 3},  // 不稳定 (4x)
		{100, 600, 0},  // 非常不稳定 (6x)
	}

	for _, tt := range tests {
		data := &ReportData{AvgTime: tt.avgTime, TP99: tt.tp99}
		result := &ScoreResult{Details: make(map[string]string)}
		score := calculateTP99Score(data, result)
		if score < tt.minScore {
			t.Errorf("TP99=%.0fms,Avg=%.0fms: 期望分数>=%d, 实际: %d", tt.tp99, tt.avgTime, tt.minScore, score)
		}
	}
}

// TestCalculateErrorCodeScore 测试错误码评分
func TestCalculateErrorCodeScore(t *testing.T) {
	// 保存原始值
	originalSuccessCode := SuccessCode
	defer func() { SuccessCode = originalSuccessCode }()

	tests := []struct {
		name        string
		codes       map[int]int
		successCode int
		expected    int
	}{
		{"仅200", map[int]int{200: 1000}, 200, 10},
		{"有4xx", map[int]int{200: 900, 404: 100}, 200, 7},
		{"有5xx", map[int]int{200: 900, 500: 100}, 200, 5},
		{"有4xx和5xx", map[int]int{200: 800, 404: 100, 500: 100}, 200, 2},
		{"自定义成功码201", map[int]int{201: 1000}, 201, 10},
		{"自定义成功码204", map[int]int{204: 900, 500: 100}, 204, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SuccessCode = tt.successCode
			data := &ReportData{ErrorCodeMap: tt.codes}
			result := &ScoreResult{Details: make(map[string]string)}
			score := calculateErrorCodeScore(data, result)
			if score != tt.expected {
				t.Errorf("%s: 期望分数=%d, 实际: %d", tt.name, tt.expected, score)
			}
		})
	}
}

// TestScoreResult_Struct 测试ScoreResult结构体
func TestScoreResult_Struct(t *testing.T) {
	result := ScoreResult{
		TotalScore:       85,
		Grade:            "B",
		SuccessRateScore: 27,
		QPSScore:         20,
		AvgTimeScore:     15,
		TP99Score:        12,
		ErrorCodeScore:   10,
		Suggestions:      []string{"建议1", "建议2"},
		Details:          map[string]string{"key": "value"},
	}

	if result.TotalScore != 85 {
		t.Errorf("TotalScore = %d, 期望 85", result.TotalScore)
	}
	if result.Grade != "B" {
		t.Errorf("Grade = %s, 期望 B", result.Grade)
	}
	if len(result.Suggestions) != 2 {
		t.Errorf("Suggestions长度 = %d, 期望 2", len(result.Suggestions))
	}
}

// TestBuildAIPrompt 测试AI提示词构建
func TestBuildAIPrompt(t *testing.T) {
	data := &ReportData{
		URL:           "http://example.com",
		Concurrency:   10,
		TotalRequests: 1000,
		SuccessNum:    990,
		FailureNum:    10,
		QPS:           100.0,
		AvgTime:       50.0,
		MinTime:       10.0,
		MaxTime:       100.0,
		TP90:          60.0,
		TP95:          70.0,
		TP99:          80.0,
	}
	baseResult := &ScoreResult{
		TotalScore:       85,
		Grade:            "B",
		SuccessRateScore: 27,
		QPSScore:         20,
		AvgTimeScore:     15,
		TP99Score:        12,
		ErrorCodeScore:   10,
	}

	prompt := buildAIPrompt(data, baseResult)

	if prompt == "" {
		t.Error("AI提示词不应为空")
	}
	if len(prompt) < 100 {
		t.Error("AI提示词应该包含足够的上下文信息")
	}
}

// TestParseAIResponse 测试AI响应解析
func TestParseAIResponse(t *testing.T) {
	result := &ScoreResult{
		Suggestions: []string{},
	}

	response := `分析结果如下：
- 建议增加服务器资源
- 优化数据库查询
• 使用缓存提升性能
普通文本不会被解析
- 检查网络延迟`

	parseAIResponse(response, result)

	if len(result.Suggestions) != 4 {
		t.Errorf("应该解析出4条建议, 实际: %d", len(result.Suggestions))
	}
}

// TestResolveAIEndpoint tests named-provider alias resolution
func TestResolveAIEndpoint(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		want     string
	}{
		{"orcarouter alias", "orcarouter", "https://api.orcarouter.ai/v1/chat/completions"},
		{"openrouter alias", "openrouter", "https://openrouter.ai/api/v1/chat/completions"},
		{"full URL passes through", "https://api.example.com/v1/chat/completions", "https://api.example.com/v1/chat/completions"},
		{"empty string passes through", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ResolveAIEndpoint(tt.endpoint); got != tt.want {
				t.Errorf("ResolveAIEndpoint(%q) = %q, want %q", tt.endpoint, got, tt.want)
			}
		})
	}
}

// TestResolveAIModel tests the default model of named providers
func TestResolveAIModel(t *testing.T) {
	if got := ResolveAIModel("https://api.orcarouter.ai/v1/chat/completions"); got != "orcarouter/fusion-mini" {
		t.Errorf("orcarouter URL should resolve to orcarouter/fusion-mini, got %q", got)
	}
	if got := ResolveAIModel("orcarouter"); got != "orcarouter/fusion-mini" {
		t.Errorf("orcarouter alias should resolve to orcarouter/fusion-mini, got %q", got)
	}
	if got := ResolveAIModel("https://openrouter.ai/api/v1/chat/completions"); got != "" {
		t.Errorf("openrouter should not have a default model, got %q", got)
	}
	if got := ResolveAIModel("https://api.example.com/v1/chat/completions"); got != "" {
		t.Errorf("custom URL should not have a default model, got %q", got)
	}
}

// BenchmarkCalculateScore 性能测试
func BenchmarkCalculateScore(b *testing.B) {
	data := &ReportData{
		URL:           "http://example.com",
		Concurrency:   100,
		TotalRequests: 10000,
		SuccessNum:    9900,
		FailureNum:    100,
		QPS:           1000.0,
		AvgTime:       100.0,
		MinTime:       10.0,
		MaxTime:       500.0,
		TP90:          150.0,
		TP95:          200.0,
		TP99:          300.0,
		ErrorCodeMap:  map[int]int{200: 9900, 500: 100},
	}

	for i := 0; i < b.N; i++ {
		CalculateScore(data)
	}
}
