// Package statistics 测试报告生成
package statistics

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/link1st/go-stress-testing/tools"
)

// ReportData 报告数据结构
type ReportData struct {
	URL             string            // 压测URL
	Concurrency     uint64            // 并发数
	TotalRequests   uint64            // 总请求数
	SuccessNum      uint64            // 成功数
	FailureNum      uint64            // 失败数
	TotalTime       float64           // 总耗时(秒)
	QPS             float64           // QPS
	MaxTime         float64           // 最大响应时间(ms)
	MinTime         float64           // 最小响应时间(ms)
	AvgTime         float64           // 平均响应时间(ms)
	TP90            float64           // TP90(ms)
	TP95            float64           // TP95(ms)
	TP99            float64           // TP99(ms)
	ReceivedBytes   int64             // 下载字节数
	ErrorCodeMap    map[int]int       // 错误码分布
	TimeRecords     []TimeRecord      // 时间序列记录
	RequestTimeList []uint64          // 请求时间列表
	StartTime       time.Time         // 开始时间
	EndTime         time.Time         // 结束时间
}

// TimeRecord 时间序列记录
type TimeRecord struct {
	Timestamp   time.Time // 记录时间点
	Elapsed     float64   // 耗时(秒)
	Concurrent  int       // 并发数
	Success     uint64    // 成功数
	Failure     uint64    // 失败数
	SuccessRate float64   // 成功率(%)
	QPS         float64   // QPS
	MaxTime     float64   // 最大响应时间(ms)
	MinTime     float64   // 最小响应时间(ms)
	AvgTime     float64   // 平均响应时间(ms)
	ErrorCodes  map[int]int // 错误码分布
}

// OutputPath 报告输出路径
var OutputPath string

// OutputFormat 报告格式: html(默认) 或 md
var OutputFormat string = "html"

// CurrentReportData 当前报告数据
var CurrentReportData *ReportData

// InitReportData 初始化报告数据
func InitReportData(url string, concurrency uint64) {
	CurrentReportData = &ReportData{
		URL:           url,
		Concurrency:   concurrency,
		ErrorCodeMap:  make(map[int]int),
		TimeRecords:   make([]TimeRecord, 0),
		StartTime:     time.Now(),
	}
}

// AddTimeRecord 添加时间序列记录
func AddTimeRecord(record TimeRecord) {
	if CurrentReportData != nil {
		CurrentReportData.TimeRecords = append(CurrentReportData.TimeRecords, record)
	}
}

// FinalizeReport 完成报告数据
func FinalizeReport(successNum, failureNum uint64, totalTime, qps, maxTime, minTime, avgTime float64,
	receivedBytes int64, errCode *sync.Map, requestTimeList []uint64) {
	if CurrentReportData == nil {
		return
	}

	CurrentReportData.EndTime = time.Now()
	CurrentReportData.TotalRequests = successNum + failureNum
	CurrentReportData.SuccessNum = successNum
	CurrentReportData.FailureNum = failureNum
	CurrentReportData.TotalTime = totalTime
	CurrentReportData.QPS = qps
	CurrentReportData.MaxTime = maxTime
	CurrentReportData.MinTime = minTime
	CurrentReportData.AvgTime = avgTime
	CurrentReportData.ReceivedBytes = receivedBytes
	CurrentReportData.RequestTimeList = requestTimeList

	// 转换错误码map
	errCode.Range(func(key, value interface{}) bool {
		if k, ok := key.(int); ok {
			if v, ok := value.(int); ok {
				CurrentReportData.ErrorCodeMap[k] = v
			}
		}
		return true
	})

	// 计算TP值
	if len(requestTimeList) > 0 {
		var all tools.Uint64List = requestTimeList
		sort.Sort(all)
		CurrentReportData.TP90 = float64(all[int(float64(len(all))*0.90)]) / 1e6
		CurrentReportData.TP95 = float64(all[int(float64(len(all))*0.95)]) / 1e6
		CurrentReportData.TP99 = float64(all[int(float64(len(all))*0.99)]) / 1e6
	}

	// 如果设置了输出路径，根据格式生成报告
	if OutputPath != "" {
		if OutputFormat == "md" {
			GenerateMarkdownReport(OutputPath)
		} else {
			GenerateHTMLReport(OutputPath)
		}
	}
}

// GenerateMarkdownReport 生成Markdown格式报告
func GenerateMarkdownReport(filepath string) error {
	if CurrentReportData == nil {
		return fmt.Errorf("没有报告数据")
	}

	data := CurrentReportData
	var sb strings.Builder

	// 标题
	sb.WriteString("# 压力测试报告\n\n")

	// 测试信息
	sb.WriteString("## 测试信息\n\n")
	sb.WriteString(fmt.Sprintf("| 项目 | 值 |\n"))
	sb.WriteString(fmt.Sprintf("|------|----|\n"))
	sb.WriteString(fmt.Sprintf("| 压测地址 | `%s` |\n", data.URL))
	sb.WriteString(fmt.Sprintf("| 并发数 | %d |\n", data.Concurrency))
	sb.WriteString(fmt.Sprintf("| 总请求数 | %d |\n", data.TotalRequests))
	sb.WriteString(fmt.Sprintf("| 开始时间 | %s |\n", data.StartTime.Format("2006-01-02 15:04:05")))
	sb.WriteString(fmt.Sprintf("| 结束时间 | %s |\n", data.EndTime.Format("2006-01-02 15:04:05")))
	sb.WriteString(fmt.Sprintf("| 总耗时 | %.3f 秒 |\n", data.TotalTime))
	sb.WriteString("\n")

	// 测试结果摘要
	sb.WriteString("## 测试结果摘要\n\n")
	sb.WriteString(fmt.Sprintf("| 指标 | 值 |\n"))
	sb.WriteString(fmt.Sprintf("|------|----|\n"))
	sb.WriteString(fmt.Sprintf("| 成功请求数 | %d |\n", data.SuccessNum))
	sb.WriteString(fmt.Sprintf("| 失败请求数 | %d |\n", data.FailureNum))
	successRate := float64(0)
	if data.TotalRequests > 0 {
		successRate = float64(data.SuccessNum) / float64(data.TotalRequests) * 100
	}
	sb.WriteString(fmt.Sprintf("| 成功率 | %.2f%% |\n", successRate))
	sb.WriteString(fmt.Sprintf("| QPS | %.2f |\n", data.QPS))
	sb.WriteString(fmt.Sprintf("| 下载字节数 | %d |\n", data.ReceivedBytes))
	sb.WriteString("\n")

	// 响应时间统计
	sb.WriteString("## 响应时间统计\n\n")
	sb.WriteString(fmt.Sprintf("| 指标 | 值 (ms) |\n"))
	sb.WriteString(fmt.Sprintf("|------|--------|\n"))
	sb.WriteString(fmt.Sprintf("| 最小响应时间 | %.2f |\n", data.MinTime))
	sb.WriteString(fmt.Sprintf("| 最大响应时间 | %.2f |\n", data.MaxTime))
	sb.WriteString(fmt.Sprintf("| 平均响应时间 | %.2f |\n", data.AvgTime))
	sb.WriteString(fmt.Sprintf("| TP90 | %.2f |\n", data.TP90))
	sb.WriteString(fmt.Sprintf("| TP95 | %.2f |\n", data.TP95))
	sb.WriteString(fmt.Sprintf("| TP99 | %.2f |\n", data.TP99))
	sb.WriteString("\n")

	// 状态码分布
	sb.WriteString("## 状态码分布\n\n")
	sb.WriteString(fmt.Sprintf("| 状态码 | 次数 | 占比 |\n"))
	sb.WriteString(fmt.Sprintf("|--------|------|------|\n"))

	// 排序状态码
	codes := make([]int, 0, len(data.ErrorCodeMap))
	for code := range data.ErrorCodeMap {
		codes = append(codes, code)
	}
	sort.Ints(codes)

	for _, code := range codes {
		count := data.ErrorCodeMap[code]
		percentage := float64(0)
		if data.TotalRequests > 0 {
			percentage = float64(count) / float64(data.TotalRequests) * 100
		}
		sb.WriteString(fmt.Sprintf("| %d | %d | %.2f%% |\n", code, count, percentage))
	}
	sb.WriteString("\n")

	// 时间序列数据
	if len(data.TimeRecords) > 0 {
		sb.WriteString("## 时间序列数据\n\n")
		sb.WriteString("| 耗时(s) | 并发数 | 成功数 | 失败数 | QPS | 最长耗时(ms) | 最短耗时(ms) | 平均耗时(ms) |\n")
		sb.WriteString("|---------|--------|--------|--------|-----|-------------|-------------|-------------|\n")

		for _, record := range data.TimeRecords {
			sb.WriteString(fmt.Sprintf("| %.0f | %d | %d | %d | %.2f | %.2f | %.2f | %.2f |\n",
				record.Elapsed, record.Concurrent, record.Success, record.Failure,
				record.QPS, record.MaxTime, record.MinTime, record.AvgTime))
		}
		sb.WriteString("\n")
	}

	// 写入文件
	err := os.WriteFile(filepath, []byte(sb.String()), 0644)
	if err != nil {
		fmt.Printf("写入报告文件失败: %v\n", err)
		return err
	}

	fmt.Printf("\n报告已导出到: %s\n", filepath)
	return nil
}
