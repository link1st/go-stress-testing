// Package statistics 报告生成测试
package statistics

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestInitReportData 测试初始化报告数据
func TestInitReportData(t *testing.T) {
	InitReportData("http://example.com", 10)

	if CurrentReportData == nil {
		t.Fatal("InitReportData() CurrentReportData 为 nil")
	}
	if CurrentReportData.URL != "http://example.com" {
		t.Errorf("URL = %q, 期望 'http://example.com'", CurrentReportData.URL)
	}
	if CurrentReportData.Concurrency != 10 {
		t.Errorf("Concurrency = %d, 期望 10", CurrentReportData.Concurrency)
	}
	if CurrentReportData.ErrorCodeMap == nil {
		t.Error("ErrorCodeMap 不应为 nil")
	}
	if CurrentReportData.TimeRecords == nil {
		t.Error("TimeRecords 不应为 nil")
	}
}

// TestAddTimeRecord 测试添加时间记录
func TestAddTimeRecord(t *testing.T) {
	InitReportData("http://example.com", 5)

	record := TimeRecord{
		Timestamp:   time.Now(),
		Elapsed:     1.0,
		Concurrent:  5,
		Success:     100,
		Failure:     0,
		SuccessRate: 100.0,
		QPS:         100.0,
		MaxTime:     50.0,
		MinTime:     10.0,
		AvgTime:     25.0,
		ErrorCodes:  map[int]int{200: 100},
	}

	AddTimeRecord(record)

	if len(CurrentReportData.TimeRecords) != 1 {
		t.Errorf("TimeRecords 长度 = %d, 期望 1", len(CurrentReportData.TimeRecords))
	}
	if CurrentReportData.TimeRecords[0].QPS != 100.0 {
		t.Errorf("QPS = %f, 期望 100.0", CurrentReportData.TimeRecords[0].QPS)
	}
	if CurrentReportData.TimeRecords[0].SuccessRate != 100.0 {
		t.Errorf("SuccessRate = %f, 期望 100.0", CurrentReportData.TimeRecords[0].SuccessRate)
	}
}

// TestAddTimeRecord_NilReportData 测试nil时添加记录
func TestAddTimeRecord_NilReportData(t *testing.T) {
	CurrentReportData = nil

	// 不应该panic
	record := TimeRecord{Elapsed: 1.0}
	AddTimeRecord(record)
}

// TestFinalizeReport 测试完成报告
func TestFinalizeReport(t *testing.T) {
	InitReportData("http://example.com", 5)

	errCode := &sync.Map{}
	errCode.Store(200, 95)
	errCode.Store(500, 5)

	requestTimeList := []uint64{
		100000000, 150000000, 200000000, 120000000, 180000000,
		250000000, 90000000, 300000000, 110000000, 160000000,
	}

	// 不设置 OutputPath，不生成文件
	OutputPath = ""
	FinalizeReport(95, 5, 10.0, 950.0, 300.0, 90.0, 150.0, 10240, errCode, requestTimeList)

	if CurrentReportData.SuccessNum != 95 {
		t.Errorf("SuccessNum = %d, 期望 95", CurrentReportData.SuccessNum)
	}
	if CurrentReportData.FailureNum != 5 {
		t.Errorf("FailureNum = %d, 期望 5", CurrentReportData.FailureNum)
	}
	if CurrentReportData.TotalRequests != 100 {
		t.Errorf("TotalRequests = %d, 期望 100", CurrentReportData.TotalRequests)
	}
}

// TestFinalizeReport_NilReportData 测试nil时完成报告
func TestFinalizeReport_NilReportData(t *testing.T) {
	CurrentReportData = nil
	errCode := &sync.Map{}

	// 不应该panic
	FinalizeReport(100, 0, 10.0, 1000.0, 100.0, 10.0, 50.0, 1024, errCode, nil)
}

// TestGenerateMarkdownReport 测试生成Markdown报告
func TestGenerateMarkdownReport(t *testing.T) {
	tmpDir := t.TempDir()
	reportPath := filepath.Join(tmpDir, "test_report.md")

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

	err := GenerateMarkdownReport(reportPath)
	if err != nil {
		t.Fatalf("GenerateMarkdownReport() 返回错误: %v", err)
	}

	// 读取生成的文件
	content, err := os.ReadFile(reportPath)
	if err != nil {
		t.Fatalf("读取报告文件失败: %v", err)
	}

	report := string(content)

	// 验证报告内容
	checks := []string{
		"# 压力测试报告",
		"http://test.example.com/api",
		"并发数 | 10",
		"总请求数 | 1000",
		"成功请求数 | 990",
		"失败请求数 | 10",
		"99.00%",
		"TP90",
		"TP95",
		"TP99",
		"状态码分布",
		"| 200 |",
		"| 500 |",
	}

	for _, check := range checks {
		if !strings.Contains(report, check) {
			t.Errorf("报告缺少内容: %q", check)
		}
	}
}

// TestGenerateMarkdownReport_NilReportData 测试nil时生成报告
func TestGenerateMarkdownReport_NilReportData(t *testing.T) {
	CurrentReportData = nil

	err := GenerateMarkdownReport("/tmp/test.md")
	if err == nil {
		t.Error("期望返回错误，但没有")
	}
}

// TestGenerateMarkdownReport_InvalidPath 测试无效路径
func TestGenerateMarkdownReport_InvalidPath(t *testing.T) {
	InitReportData("http://example.com", 5)
	CurrentReportData.ErrorCodeMap = map[int]int{200: 100}

	err := GenerateMarkdownReport("/nonexistent/directory/report.md")
	if err == nil {
		t.Error("期望返回错误，但没有")
	}
}

// TestGenerateMarkdownReport_WithTimeRecords 测试带时间记录的报告
func TestGenerateMarkdownReport_WithTimeRecords(t *testing.T) {
	tmpDir := t.TempDir()
	reportPath := filepath.Join(tmpDir, "report_with_records.md")

	InitReportData("http://example.com", 5)
	CurrentReportData.ErrorCodeMap = map[int]int{200: 100}
	CurrentReportData.TimeRecords = []TimeRecord{
		{Timestamp: time.Now(), Elapsed: 1.0, Concurrent: 5, Success: 50, Failure: 0, SuccessRate: 100.0, QPS: 50.0, MaxTime: 100.0, MinTime: 10.0, AvgTime: 50.0, ErrorCodes: map[int]int{200: 50}},
		{Timestamp: time.Now(), Elapsed: 2.0, Concurrent: 5, Success: 100, Failure: 0, SuccessRate: 100.0, QPS: 50.0, MaxTime: 100.0, MinTime: 10.0, AvgTime: 50.0, ErrorCodes: map[int]int{200: 100}},
	}

	err := GenerateMarkdownReport(reportPath)
	if err != nil {
		t.Fatalf("GenerateMarkdownReport() 返回错误: %v", err)
	}

	content, _ := os.ReadFile(reportPath)
	report := string(content)

	if !strings.Contains(report, "时间序列数据") {
		t.Error("报告应该包含时间序列数据")
	}
}

// TestReportData_Struct 测试ReportData结构体
func TestReportData_Struct(t *testing.T) {
	data := ReportData{
		URL:           "http://example.com",
		Concurrency:   10,
		TotalRequests: 1000,
		SuccessNum:    990,
		FailureNum:    10,
		TotalTime:     60.0,
		QPS:           16.5,
		MaxTime:       500.0,
		MinTime:       10.0,
		AvgTime:       100.0,
		TP90:          200.0,
		TP95:          300.0,
		TP99:          450.0,
		ReceivedBytes: 102400,
		ErrorCodeMap:  map[int]int{200: 990},
	}

	if data.URL != "http://example.com" {
		t.Errorf("URL = %q, 期望 'http://example.com'", data.URL)
	}
	if data.TotalRequests != 1000 {
		t.Errorf("TotalRequests = %d, 期望 1000", data.TotalRequests)
	}
}

// TestTimeRecord_Struct 测试TimeRecord结构体
func TestTimeRecord_Struct(t *testing.T) {
	now := time.Now()
	record := TimeRecord{
		Timestamp:   now,
		Elapsed:     1.5,
		Concurrent:  10,
		Success:     100,
		Failure:     5,
		SuccessRate: 95.24,
		QPS:         66.67,
		MaxTime:     200.0,
		MinTime:     20.0,
		AvgTime:     100.0,
		ErrorCodes:  map[int]int{200: 100, 500: 5},
	}

	if record.Elapsed != 1.5 {
		t.Errorf("Elapsed = %f, 期望 1.5", record.Elapsed)
	}
	if record.Concurrent != 10 {
		t.Errorf("Concurrent = %d, 期望 10", record.Concurrent)
	}
	if record.SuccessRate != 95.24 {
		t.Errorf("SuccessRate = %f, 期望 95.24", record.SuccessRate)
	}
	if len(record.ErrorCodes) != 2 {
		t.Errorf("ErrorCodes长度 = %d, 期望 2", len(record.ErrorCodes))
	}
}

// BenchmarkGenerateMarkdownReport 性能测试
func BenchmarkGenerateMarkdownReport(b *testing.B) {
	tmpDir := b.TempDir()

	for i := 0; i < b.N; i++ {
		InitReportData("http://example.com", 100)
		CurrentReportData.TotalRequests = 10000
		CurrentReportData.SuccessNum = 9900
		CurrentReportData.FailureNum = 100
		CurrentReportData.ErrorCodeMap = map[int]int{200: 9900, 500: 100}

		reportPath := filepath.Join(tmpDir, "bench_report.md")
		GenerateMarkdownReport(reportPath)
	}
}
