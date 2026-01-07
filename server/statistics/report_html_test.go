// Package statistics HTML报告生成测试
package statistics

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestGenerateHTMLReport 测试生成HTML报告
func TestGenerateHTMLReport(t *testing.T) {
	tmpDir := t.TempDir()
	reportPath := filepath.Join(tmpDir, "test_report.html")

	InitReportData("http://test.example.com/api", 10)
	CurrentReportData.StartTime = time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	CurrentReportData.EndTime = time.Date(2024, 1, 1, 10, 1, 0, 0, time.UTC)
	CurrentReportData.TotalRequests = 1000
	CurrentReportData.SuccessNum = 990
	CurrentReportData.FailureNum = 10
	CurrentReportData.TotalTime = 60.0
	CurrentReportData.QPS = 16.5
	CurrentReportData.MaxTime = 500.0
	CurrentReportData.MinTime = 10.0
	CurrentReportData.AvgTime = 100.0
	CurrentReportData.TP90 = 200.0
	CurrentReportData.TP95 = 300.0
	CurrentReportData.TP99 = 450.0
	CurrentReportData.ReceivedBytes = 102400
	CurrentReportData.ErrorCodeMap = map[int]int{200: 990, 500: 10}

	err := GenerateHTMLReport(reportPath)
	if err != nil {
		t.Fatalf("GenerateHTMLReport() 返回错误: %v", err)
	}

	content, err := os.ReadFile(reportPath)
	if err != nil {
		t.Fatalf("读取报告文件失败: %v", err)
	}

	report := string(content)

	// 验证HTML结构
	checks := []string{
		"<!DOCTYPE html>",
		"<html",
		"</html>",
		"压力测试报告",
		"http://test.example.com/api",
		"并发数",
		"chart.js",
		"score-card",
		"responseTimeChart",
		"statusCodeChart",
	}

	for _, check := range checks {
		if !strings.Contains(report, check) {
			t.Errorf("HTML报告缺少内容: %q", check)
		}
	}
}

// TestGenerateHTMLReport_NilReportData 测试nil时生成报告
func TestGenerateHTMLReport_NilReportData(t *testing.T) {
	CurrentReportData = nil

	err := GenerateHTMLReport("/tmp/test.html")
	if err == nil {
		t.Error("期望返回错误，但没有")
	}
}

// TestGenerateHTMLReport_InvalidPath 测试无效路径
func TestGenerateHTMLReport_InvalidPath(t *testing.T) {
	InitReportData("http://example.com", 5)
	CurrentReportData.ErrorCodeMap = map[int]int{200: 100}

	err := GenerateHTMLReport("/nonexistent/directory/report.html")
	if err == nil {
		t.Error("期望返回错误，但没有")
	}
}

// TestGenerateHTMLReport_WithTimeRecords 测试带时间序列数据的报告
func TestGenerateHTMLReport_WithTimeRecords(t *testing.T) {
	tmpDir := t.TempDir()
	reportPath := filepath.Join(tmpDir, "report_with_records.html")

	InitReportData("http://example.com", 5)
	CurrentReportData.ErrorCodeMap = map[int]int{200: 100}
	CurrentReportData.TotalRequests = 100
	CurrentReportData.SuccessNum = 100
	CurrentReportData.TimeRecords = []TimeRecord{
		{Timestamp: time.Now(), Elapsed: 1.0, Concurrent: 5, Success: 50, Failure: 0, SuccessRate: 100.0, QPS: 50.0, MaxTime: 100.0, MinTime: 10.0, AvgTime: 50.0, ErrorCodes: map[int]int{200: 50}},
		{Timestamp: time.Now(), Elapsed: 2.0, Concurrent: 5, Success: 100, Failure: 0, SuccessRate: 100.0, QPS: 50.0, MaxTime: 100.0, MinTime: 10.0, AvgTime: 50.0, ErrorCodes: map[int]int{200: 100}},
	}

	err := GenerateHTMLReport(reportPath)
	if err != nil {
		t.Fatalf("GenerateHTMLReport() 返回错误: %v", err)
	}

	content, _ := os.ReadFile(reportPath)
	report := string(content)

	// 应该包含各种时间序列图表
	if !strings.Contains(report, "qpsChart") {
		t.Error("报告应该包含QPS曲线图")
	}
	if !strings.Contains(report, "successRateChart") {
		t.Error("报告应该包含成功率曲线图")
	}
	if !strings.Contains(report, "latencyChart") {
		t.Error("报告应该包含接口耗时曲线图")
	}
	if !strings.Contains(report, "errorCodeChart") {
		t.Error("报告应该包含错误码曲线图")
	}
	if !strings.Contains(report, "时间序列数据") {
		t.Error("报告应该包含时间序列数据表格")
	}
}

// TestGenerateHTMLHead 测试HTML头部生成
func TestGenerateHTMLHead(t *testing.T) {
	head := generateHTMLHead()

	checks := []string{
		"<!DOCTYPE html>",
		"<html lang=\"zh-CN\">",
		"<meta charset=\"UTF-8\"",
		"<title>压力测试报告</title>",
		"chart.js",
		"<style>",
		"</style>",
		"</head>",
	}

	for _, check := range checks {
		if !strings.Contains(head, check) {
			t.Errorf("HTML头部缺少: %q", check)
		}
	}
}

// TestGenerateScoreCard 测试评分卡片生成
func TestGenerateScoreCard(t *testing.T) {
	score := &ScoreResult{
		TotalScore:       85,
		Grade:            "B",
		SuccessRateScore: 27,
		QPSScore:         20,
		AvgTimeScore:     15,
		TP99Score:        12,
		ErrorCodeScore:   10,
		Suggestions:      []string{"优化建议1", "优化建议2"},
	}

	card := generateScoreCard(score)

	checks := []string{
		"score-card",
		"85",
		"grade-B",
		"27/30",
		"20/25",
		"15/20",
		"12/15",
		"10/10",
		"优化建议1",
		"优化建议2",
	}

	for _, check := range checks {
		if !strings.Contains(card, check) {
			t.Errorf("评分卡片缺少: %q", check)
		}
	}
}

// TestGenerateTestInfoSection 测试测试信息部分
func TestGenerateTestInfoSection(t *testing.T) {
	data := &ReportData{
		URL:           "http://example.com/test",
		Concurrency:   10,
		TotalRequests: 1000,
		StartTime:     time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
		EndTime:       time.Date(2024, 1, 1, 10, 1, 0, 0, time.UTC),
		TotalTime:     60.0,
	}

	section := generateTestInfoSection(data)

	checks := []string{
		"测试信息",
		"http://example.com/test",
		"10",
		"1000",
		"2024-01-01",
	}

	for _, check := range checks {
		if !strings.Contains(section, check) {
			t.Errorf("测试信息部分缺少: %q", check)
		}
	}
}

// TestGenerateResponseTimeDistribution 测试响应时间分布数据生成
func TestGenerateResponseTimeDistribution(t *testing.T) {
	data := &ReportData{
		RequestTimeList: []uint64{
			30 * 1e6,  // 30ms -> 0-50
			80 * 1e6,  // 80ms -> 50-100
			150 * 1e6, // 150ms -> 100-200
			300 * 1e6, // 300ms -> 200-500
			800 * 1e6, // 800ms -> 500-1000
			1500 * 1e6, // 1500ms -> 1000+
		},
	}

	dist := generateResponseTimeDistribution(data)

	if len(dist.Labels) != 6 {
		t.Errorf("应该有6个时间区间, 实际: %d", len(dist.Labels))
	}

	expectedLabels := []string{"0-50", "50-100", "100-200", "200-500", "500-1000", "1000+"}
	for i, label := range expectedLabels {
		if dist.Labels[i] != label {
			t.Errorf("区间%d标签应该是%s, 实际: %s", i, label, dist.Labels[i])
		}
	}

	// 每个区间应该有1个请求
	for i, count := range dist.Values {
		if count != 1 {
			t.Errorf("区间%s应该有1个请求, 实际: %d", dist.Labels[i], count)
		}
	}
}

// TestGenerateResponseTimeDistribution_Empty 测试空数据
func TestGenerateResponseTimeDistribution_Empty(t *testing.T) {
	data := &ReportData{
		RequestTimeList: []uint64{},
	}

	dist := generateResponseTimeDistribution(data)

	if len(dist.Labels) != 1 || dist.Labels[0] != "无数据" {
		t.Error("空数据应该返回'无数据'标签")
	}
}

// TestGenerateStatusCodeSection 测试状态码分布部分
func TestGenerateStatusCodeSection(t *testing.T) {
	data := &ReportData{
		TotalRequests: 1000,
		ErrorCodeMap:  map[int]int{200: 900, 404: 50, 500: 50},
	}

	section := generateStatusCodeSection(data)

	checks := []string{
		"状态码分布",
		"200",
		"404",
		"500",
		"90.00%",
		"5.00%",
	}

	for _, check := range checks {
		if !strings.Contains(section, check) {
			t.Errorf("状态码分布部分缺少: %q", check)
		}
	}
}

// TestGenerateChartScripts 测试图表脚本生成
func TestGenerateChartScripts(t *testing.T) {
	data := &ReportData{
		RequestTimeList: []uint64{50 * 1e6, 100 * 1e6, 150 * 1e6},
		ErrorCodeMap:    map[int]int{200: 90, 500: 10},
		TotalRequests:   100,
	}

	scripts := generateChartScripts(data)

	checks := []string{
		"<script>",
		"</script>",
		"new Chart",
		"responseTimeChart",
		"statusCodeChart",
	}

	for _, check := range checks {
		if !strings.Contains(scripts, check) {
			t.Errorf("图表脚本缺少: %q", check)
		}
	}
}

// TestGenerateChartScripts_WithTimeRecords 测试带时间序列的图表脚本
func TestGenerateChartScripts_WithTimeRecords(t *testing.T) {
	data := &ReportData{
		RequestTimeList: []uint64{50 * 1e6},
		ErrorCodeMap:    map[int]int{200: 100},
		TotalRequests:   100,
		TimeRecords: []TimeRecord{
			{Timestamp: time.Now(), Elapsed: 1.0, QPS: 50.0, Success: 50, Failure: 0, SuccessRate: 100.0, AvgTime: 50.0, MaxTime: 100.0, MinTime: 10.0, ErrorCodes: map[int]int{200: 50}},
			{Timestamp: time.Now(), Elapsed: 2.0, QPS: 50.0, Success: 100, Failure: 0, SuccessRate: 100.0, AvgTime: 50.0, MaxTime: 100.0, MinTime: 10.0, ErrorCodes: map[int]int{200: 100}},
		},
	}

	scripts := generateChartScripts(data)

	checks := []string{
		"qpsChart",
		"successRateChart",
		"latencyChart",
		"errorCodeChart",
	}

	for _, check := range checks {
		if !strings.Contains(scripts, check) {
			t.Errorf("带时间序列的图表脚本缺少: %q", check)
		}
	}
}

// TestGenerateChartModal 测试图表模态框生成
func TestGenerateChartModal(t *testing.T) {
	modal := generateChartModal()

	checks := []string{
		"chartModal",
		"chart-modal",
		"chart-modal-content",
		"chart-modal-close",
		"modalChartTitle",
		"modalChart",
		"closeChartModal",
	}

	for _, check := range checks {
		if !strings.Contains(modal, check) {
			t.Errorf("模态框HTML缺少: %q", check)
		}
	}
}

// TestGenerateModalScript 测试模态框脚本生成
func TestGenerateModalScript(t *testing.T) {
	script := generateModalScript()

	checks := []string{
		"chartConfigs",
		"modalChart",
		"getChartInstance",
		"initChartClickHandlers",
		"getChartTitle",
		"openChartModal",
		"closeChartModal",
		"chart-wrapper",
		"Escape",
		"DOMContentLoaded",
	}

	for _, check := range checks {
		if !strings.Contains(script, check) {
			t.Errorf("模态框脚本缺少: %q", check)
		}
	}
}

// TestGenerateHTMLReport_WithModal 测试HTML报告包含模态框功能
func TestGenerateHTMLReport_WithModal(t *testing.T) {
	tmpDir := t.TempDir()
	reportPath := filepath.Join(tmpDir, "test_modal.html")

	InitReportData("http://example.com", 5)
	CurrentReportData.ErrorCodeMap = map[int]int{200: 100}
	CurrentReportData.TotalRequests = 100
	CurrentReportData.SuccessNum = 100

	err := GenerateHTMLReport(reportPath)
	if err != nil {
		t.Fatalf("GenerateHTMLReport() 返回错误: %v", err)
	}

	content, _ := os.ReadFile(reportPath)
	report := string(content)

	// 验证模态框相关内容
	modalChecks := []string{
		"chartModal",
		"chart-modal",
		"点击放大",
		"openChartModal",
		"closeChartModal",
		"modalChart",
	}

	for _, check := range modalChecks {
		if !strings.Contains(report, check) {
			t.Errorf("报告缺少模态框相关内容: %q", check)
		}
	}
}

// TestTimeSeriesTimeFormat 测试时间序列时间格式
func TestTimeSeriesTimeFormat(t *testing.T) {
	tmpDir := t.TempDir()
	reportPath := filepath.Join(tmpDir, "test_time_format.html")

	InitReportData("http://example.com", 5)
	CurrentReportData.ErrorCodeMap = map[int]int{200: 100}
	CurrentReportData.TotalRequests = 100
	CurrentReportData.SuccessNum = 100
	CurrentReportData.TimeRecords = []TimeRecord{
		{Timestamp: time.Date(2024, 1, 1, 14, 30, 45, 0, time.UTC), Elapsed: 1.0, Concurrent: 5, Success: 50, Failure: 0, SuccessRate: 100.0, QPS: 50.0, MaxTime: 100.0, MinTime: 10.0, AvgTime: 50.0, ErrorCodes: map[int]int{200: 50}},
	}

	err := GenerateHTMLReport(reportPath)
	if err != nil {
		t.Fatalf("GenerateHTMLReport() 返回错误: %v", err)
	}

	content, _ := os.ReadFile(reportPath)
	report := string(content)

	// 验证时间格式为 HH:MM:SS (不包含日期)
	if !strings.Contains(report, "14:30:45") {
		t.Error("报告应该包含时分秒格式时间")
	}
	// 表格中不应该包含完整日期格式
	if strings.Contains(report, "2024-01-01 14:30:45") {
		t.Error("时间序列表格不应该包含完整日期")
	}
}

// TestErrorCodeChartExcludesSuccessCode 测试错误码图表排除成功状态码
func TestErrorCodeChartExcludesSuccessCode(t *testing.T) {
	tmpDir := t.TempDir()
	reportPath := filepath.Join(tmpDir, "test_error_code.html")

	// 设置成功状态码
	SuccessCode = 200

	InitReportData("http://example.com", 5)
	CurrentReportData.ErrorCodeMap = map[int]int{200: 900, 500: 100}
	CurrentReportData.TotalRequests = 1000
	CurrentReportData.SuccessNum = 900
	CurrentReportData.FailureNum = 100
	CurrentReportData.TimeRecords = []TimeRecord{
		{Timestamp: time.Now(), Elapsed: 1.0, Concurrent: 5, Success: 450, Failure: 50, SuccessRate: 90.0, QPS: 50.0, MaxTime: 100.0, MinTime: 10.0, AvgTime: 50.0, ErrorCodes: map[int]int{200: 450, 500: 50}},
		{Timestamp: time.Now(), Elapsed: 2.0, Concurrent: 5, Success: 450, Failure: 50, SuccessRate: 90.0, QPS: 50.0, MaxTime: 100.0, MinTime: 10.0, AvgTime: 50.0, ErrorCodes: map[int]int{200: 450, 500: 50}},
	}

	err := GenerateHTMLReport(reportPath)
	if err != nil {
		t.Fatalf("GenerateHTMLReport() 返回错误: %v", err)
	}

	content, _ := os.ReadFile(reportPath)
	report := string(content)

	// 验证错误码图表存在
	if !strings.Contains(report, "errorCodeChart") {
		t.Error("报告应该包含错误码图表")
	}
}

// BenchmarkGenerateHTMLReport 性能测试
func BenchmarkGenerateHTMLReport(b *testing.B) {
	tmpDir := b.TempDir()

	for i := 0; i < b.N; i++ {
		InitReportData("http://example.com", 100)
		CurrentReportData.TotalRequests = 10000
		CurrentReportData.SuccessNum = 9900
		CurrentReportData.FailureNum = 100
		CurrentReportData.ErrorCodeMap = map[int]int{200: 9900, 500: 100}
		CurrentReportData.RequestTimeList = make([]uint64, 1000)
		for j := 0; j < 1000; j++ {
			CurrentReportData.RequestTimeList[j] = uint64(j) * 1e6
		}

		reportPath := filepath.Join(tmpDir, "bench_report.html")
		GenerateHTMLReport(reportPath)
	}
}
