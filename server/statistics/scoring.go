// Package statistics AI评分系统
package statistics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ScoreResult 评分结果
type ScoreResult struct {
	TotalScore      int            `json:"total_score"`      // 总分 (0-100)
	Grade           string         `json:"grade"`            // 评级 (A/B/C/D/F)
	SuccessRateScore int           `json:"success_rate_score"` // 成功率得分
	QPSScore        int            `json:"qps_score"`        // QPS得分
	AvgTimeScore    int            `json:"avg_time_score"`   // 平均响应时间得分
	TP99Score       int            `json:"tp99_score"`       // TP99稳定性得分
	ErrorCodeScore  int            `json:"error_code_score"` // 错误码得分
	Suggestions     []string       `json:"suggestions"`      // 改进建议
	Details         map[string]string `json:"details"`       // 详细说明
}

// AIConfig AI API配置
type AIConfig struct {
	APIEndpoint string // API地址
	APIKey      string // API Key
	Model       string // 模型名称
}

// AI配置全局变量
var (
	AIAPIEndpoint string
	AIAPIKey      string
	AIModel       string = "gpt-3.5-turbo"
)

// SuccessCode 成功状态码（由 -code 参数指定）
var SuccessCode int = 200

// CalculateScore 计算内置评分
func CalculateScore(data *ReportData) *ScoreResult {
	if data == nil {
		return &ScoreResult{
			TotalScore:  0,
			Grade:       "F",
			Suggestions: []string{"没有测试数据"},
		}
	}

	result := &ScoreResult{
		Details:     make(map[string]string),
		Suggestions: make([]string, 0),
	}

	// 1. 成功率评分 (满分30分)
	result.SuccessRateScore = calculateSuccessRateScore(data, result)

	// 2. QPS评分 (满分25分)
	result.QPSScore = calculateQPSScore(data, result)

	// 3. 平均响应时间评分 (满分20分)
	result.AvgTimeScore = calculateAvgTimeScore(data, result)

	// 4. TP99稳定性评分 (满分15分)
	result.TP99Score = calculateTP99Score(data, result)

	// 5. 错误码评分 (满分10分)
	result.ErrorCodeScore = calculateErrorCodeScore(data, result)

	// 计算总分
	result.TotalScore = result.SuccessRateScore + result.QPSScore +
		result.AvgTimeScore + result.TP99Score + result.ErrorCodeScore

	// 确保总分在0-100范围内
	if result.TotalScore < 0 {
		result.TotalScore = 0
	} else if result.TotalScore > 100 {
		result.TotalScore = 100
	}

	// 评级
	result.Grade = calculateGrade(result.TotalScore)

	return result
}

// calculateSuccessRateScore 计算成功率得分
func calculateSuccessRateScore(data *ReportData, result *ScoreResult) int {
	if data.TotalRequests == 0 {
		result.Details["success_rate"] = "无请求数据"
		return 0
	}

	successRate := float64(data.SuccessNum) / float64(data.TotalRequests) * 100
	var score int

	switch {
	case successRate >= 99.99:
		score = 30
		result.Details["success_rate"] = fmt.Sprintf("成功率 %.2f%% - 完美", successRate)
	case successRate >= 99.9:
		score = 28
		result.Details["success_rate"] = fmt.Sprintf("成功率 %.2f%% - 极优", successRate)
	case successRate >= 99.5:
		score = 26
		result.Details["success_rate"] = fmt.Sprintf("成功率 %.2f%% - 优秀", successRate)
	case successRate >= 99:
		score = 24
		result.Details["success_rate"] = fmt.Sprintf("成功率 %.2f%% - 良好", successRate)
	case successRate >= 98:
		score = 21
		result.Details["success_rate"] = fmt.Sprintf("成功率 %.2f%% - 较好", successRate)
		result.Suggestions = append(result.Suggestions, "建议排查失败请求原因，提升成功率至99%以上")
	case successRate >= 95:
		score = 16
		result.Details["success_rate"] = fmt.Sprintf("成功率 %.2f%% - 一般", successRate)
		result.Suggestions = append(result.Suggestions, "成功率一般，需要关注失败原因")
	case successRate >= 90:
		score = 10
		result.Details["success_rate"] = fmt.Sprintf("成功率 %.2f%% - 较差", successRate)
		result.Suggestions = append(result.Suggestions, "成功率较低，需要重点排查失败原因")
	case successRate >= 80:
		score = 5
		result.Details["success_rate"] = fmt.Sprintf("成功率 %.2f%% - 差", successRate)
		result.Suggestions = append(result.Suggestions, "成功率过低，系统可能存在问题")
	default:
		score = 0
		result.Details["success_rate"] = fmt.Sprintf("成功率 %.2f%% - 极差", successRate)
		result.Suggestions = append(result.Suggestions, "成功率严重不足，请立即排查系统问题")
	}

	return score
}

// calculateQPSScore 计算QPS得分
func calculateQPSScore(data *ReportData, result *ScoreResult) int {
	if data.Concurrency == 0 {
		result.Details["qps"] = "无并发数据"
		return 0
	}

	// 计算每并发QPS效率
	qpsPerConcurrency := data.QPS / float64(data.Concurrency)
	var score int

	switch {
	case qpsPerConcurrency >= 100:
		score = 25
		result.Details["qps"] = fmt.Sprintf("QPS %.2f (%.2f/并发) - 极高吞吐", data.QPS, qpsPerConcurrency)
	case qpsPerConcurrency >= 50:
		score = 22
		result.Details["qps"] = fmt.Sprintf("QPS %.2f (%.2f/并发) - 高吞吐", data.QPS, qpsPerConcurrency)
	case qpsPerConcurrency >= 20:
		score = 18
		result.Details["qps"] = fmt.Sprintf("QPS %.2f (%.2f/并发) - 中等吞吐", data.QPS, qpsPerConcurrency)
	case qpsPerConcurrency >= 10:
		score = 12
		result.Details["qps"] = fmt.Sprintf("QPS %.2f (%.2f/并发) - 一般吞吐", data.QPS, qpsPerConcurrency)
		result.Suggestions = append(result.Suggestions, "QPS效率一般，建议检查服务端性能瓶颈")
	case qpsPerConcurrency >= 1:
		score = 5
		result.Details["qps"] = fmt.Sprintf("QPS %.2f (%.2f/并发) - 较低吞吐", data.QPS, qpsPerConcurrency)
		result.Suggestions = append(result.Suggestions, "QPS较低，可能存在性能问题，建议进行性能分析")
	default:
		score = 0
		result.Details["qps"] = fmt.Sprintf("QPS %.2f (%.2f/并发) - 极低吞吐", data.QPS, qpsPerConcurrency)
		result.Suggestions = append(result.Suggestions, "QPS过低，系统吞吐能力严重不足")
	}

	return score
}

// calculateAvgTimeScore 计算平均响应时间得分
func calculateAvgTimeScore(data *ReportData, result *ScoreResult) int {
	avgTime := data.AvgTime // 单位: ms
	var score int

	switch {
	case avgTime <= 50:
		score = 20
		result.Details["avg_time"] = fmt.Sprintf("平均响应 %.2fms - 极快", avgTime)
	case avgTime <= 100:
		score = 19
		result.Details["avg_time"] = fmt.Sprintf("平均响应 %.2fms - 很快", avgTime)
	case avgTime <= 200:
		score = 18
		result.Details["avg_time"] = fmt.Sprintf("平均响应 %.2fms - 快速", avgTime)
	case avgTime <= 300:
		score = 17
		result.Details["avg_time"] = fmt.Sprintf("平均响应 %.2fms - 较快", avgTime)
	case avgTime <= 500:
		score = 15
		result.Details["avg_time"] = fmt.Sprintf("平均响应 %.2fms - 良好", avgTime)
	case avgTime <= 800:
		score = 13
		result.Details["avg_time"] = fmt.Sprintf("平均响应 %.2fms - 一般", avgTime)
		result.Suggestions = append(result.Suggestions, "响应时间略长，建议优化接口性能")
	case avgTime <= 1000:
		score = 11
		result.Details["avg_time"] = fmt.Sprintf("平均响应 %.2fms - 可接受", avgTime)
		result.Suggestions = append(result.Suggestions, "响应时间略长，建议优化接口性能")
	case avgTime <= 1500:
		score = 9
		result.Details["avg_time"] = fmt.Sprintf("平均响应 %.2fms - 略慢", avgTime)
		result.Suggestions = append(result.Suggestions, "响应时间较长，建议检查慢查询或外部调用")
	case avgTime <= 2000:
		score = 7
		result.Details["avg_time"] = fmt.Sprintf("平均响应 %.2fms - 较慢", avgTime)
		result.Suggestions = append(result.Suggestions, "响应时间较长，建议检查慢查询或外部调用")
	case avgTime <= 3000:
		score = 5
		result.Details["avg_time"] = fmt.Sprintf("平均响应 %.2fms - 慢", avgTime)
		result.Suggestions = append(result.Suggestions, "响应时间过长，需要进行性能优化")
	case avgTime <= 5000:
		score = 3
		result.Details["avg_time"] = fmt.Sprintf("平均响应 %.2fms - 很慢", avgTime)
		result.Suggestions = append(result.Suggestions, "响应时间过长，需要进行性能优化")
	default:
		score = 0
		result.Details["avg_time"] = fmt.Sprintf("平均响应 %.2fms - 极慢", avgTime)
		result.Suggestions = append(result.Suggestions, "响应时间严重超标，请立即优化")
	}

	return score
}

// calculateTP99Score 计算TP99稳定性得分
func calculateTP99Score(data *ReportData, result *ScoreResult) int {
	if data.AvgTime == 0 {
		result.Details["tp99"] = "无响应时间数据"
		return 0
	}

	// TP99与平均时间的比值反映稳定性
	ratio := data.TP99 / data.AvgTime
	var score int

	switch {
	case ratio <= 1.5:
		score = 15
		result.Details["tp99"] = fmt.Sprintf("TP99 %.2fms (%.1fx平均) - 非常稳定", data.TP99, ratio)
	case ratio <= 2.0:
		score = 12
		result.Details["tp99"] = fmt.Sprintf("TP99 %.2fms (%.1fx平均) - 稳定", data.TP99, ratio)
	case ratio <= 3.0:
		score = 9
		result.Details["tp99"] = fmt.Sprintf("TP99 %.2fms (%.1fx平均) - 一般", data.TP99, ratio)
		result.Suggestions = append(result.Suggestions, "响应时间波动较大，建议检查是否存在偶发慢请求")
	case ratio <= 5.0:
		score = 5
		result.Details["tp99"] = fmt.Sprintf("TP99 %.2fms (%.1fx平均) - 不稳定", data.TP99, ratio)
		result.Suggestions = append(result.Suggestions, "存在明显的长尾请求，需要排查原因")
	default:
		score = 0
		result.Details["tp99"] = fmt.Sprintf("TP99 %.2fms (%.1fx平均) - 非常不稳定", data.TP99, ratio)
		result.Suggestions = append(result.Suggestions, "响应时间极不稳定，可能存在资源竞争或GC问题")
	}

	return score
}

// calculateErrorCodeScore 计算错误码得分
func calculateErrorCodeScore(data *ReportData, result *ScoreResult) int {
	score := 10 // 基础分

	has4xx := false
	has5xx := false
	onlySuccess := true

	for code := range data.ErrorCodeMap {
		// 使用 SuccessCode 判断成功状态码（支持 -code 参数）
		if code != SuccessCode && code != 0 {
			onlySuccess = false
		}
		if code >= 400 && code < 500 {
			has4xx = true
		}
		if code >= 500 {
			has5xx = true
		}
	}

	if onlySuccess {
		result.Details["error_code"] = fmt.Sprintf("仅有成功状态码(%d) - 完美", SuccessCode)
		return score
	}

	if has5xx {
		score -= 5
		result.Suggestions = append(result.Suggestions, "存在5xx服务端错误，请检查服务端日志")
	}

	if has4xx {
		score -= 3
		result.Suggestions = append(result.Suggestions, "存在4xx客户端错误，请检查请求参数")
	}

	if score < 0 {
		score = 0
	}

	var codeList []string
	for code, count := range data.ErrorCodeMap {
		codeList = append(codeList, fmt.Sprintf("%d:%d", code, count))
	}
	result.Details["error_code"] = fmt.Sprintf("状态码分布: %s", strings.Join(codeList, ", "))

	return score
}

// calculateGrade 计算评级
func calculateGrade(score int) string {
	switch {
	case score >= 90:
		return "A"
	case score >= 80:
		return "B"
	case score >= 70:
		return "C"
	case score >= 60:
		return "D"
	default:
		return "F"
	}
}

// CalculateScoreWithAI 使用外部AI API计算评分
func CalculateScoreWithAI(data *ReportData) (*ScoreResult, error) {
	if AIAPIEndpoint == "" || AIAPIKey == "" {
		return nil, fmt.Errorf("AI API未配置")
	}

	// 先计算内置评分作为基础
	baseResult := CalculateScore(data)

	// 构建AI请求
	prompt := buildAIPrompt(data, baseResult)

	// 调用AI API
	aiResponse, err := callAIAPI(prompt)
	if err != nil {
		// AI调用失败，返回内置评分
		baseResult.Suggestions = append(baseResult.Suggestions, fmt.Sprintf("AI分析失败: %v", err))
		return baseResult, nil
	}

	// 解析AI响应并合并建议
	parseAIResponse(aiResponse, baseResult)

	return baseResult, nil
}

// buildAIPrompt 构建AI提示词
func buildAIPrompt(data *ReportData, baseResult *ScoreResult) string {
	return fmt.Sprintf(`作为性能测试专家，请分析以下压力测试结果并给出专业建议：

## 测试数据
- URL: %s
- 并发数: %d
- 总请求数: %d
- 成功数: %d, 失败数: %d
- 成功率: %.2f%%
- QPS: %.2f
- 平均响应时间: %.2fms
- 最小响应时间: %.2fms
- 最大响应时间: %.2fms
- TP90: %.2fms, TP95: %.2fms, TP99: %.2fms

## 内置评分结果
- 总分: %d/100 (评级: %s)
- 成功率得分: %d/30
- QPS得分: %d/25
- 响应时间得分: %d/20
- 稳定性得分: %d/15
- 错误码得分: %d/10

请用中文给出3-5条具体的性能优化建议，每条建议一行，以"-"开头。`,
		data.URL, data.Concurrency, data.TotalRequests,
		data.SuccessNum, data.FailureNum,
		float64(data.SuccessNum)/float64(data.TotalRequests)*100,
		data.QPS, data.AvgTime, data.MinTime, data.MaxTime,
		data.TP90, data.TP95, data.TP99,
		baseResult.TotalScore, baseResult.Grade,
		baseResult.SuccessRateScore, baseResult.QPSScore,
		baseResult.AvgTimeScore, baseResult.TP99Score, baseResult.ErrorCodeScore)
}

// AIRequest AI API请求结构
type AIRequest struct {
	Model    string       `json:"model"`
	Messages []AIMessage  `json:"messages"`
}

// AIMessage AI消息结构
type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AIResponse AI API响应结构
type AIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// callAIAPI 调用AI API
func callAIAPI(prompt string) (string, error) {
	reqBody := AIRequest{
		Model: AIModel,
		Messages: []AIMessage{
			{Role: "user", Content: prompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("JSON编码失败: %v", err)
	}

	req, err := http.NewRequest("POST", AIAPIEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+AIAPIKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	var aiResp AIResponse
	if err := json.Unmarshal(body, &aiResp); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	if aiResp.Error != nil {
		return "", fmt.Errorf("AI API错误: %s", aiResp.Error.Message)
	}

	if len(aiResp.Choices) == 0 {
		return "", fmt.Errorf("AI无返回结果")
	}

	return aiResp.Choices[0].Message.Content, nil
}

// parseAIResponse 解析AI响应
func parseAIResponse(response string, result *ScoreResult) {
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "-") || strings.HasPrefix(line, "•") {
			suggestion := strings.TrimPrefix(line, "-")
			suggestion = strings.TrimPrefix(suggestion, "•")
			suggestion = strings.TrimSpace(suggestion)
			if suggestion != "" {
				result.Suggestions = append(result.Suggestions, "[AI] "+suggestion)
			}
		}
	}
}
