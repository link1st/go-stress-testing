// Package statistics HTML报告生成
package statistics

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

// GenerateHTMLReport 生成HTML格式报告
func GenerateHTMLReport(filepath string) error {
	if CurrentReportData == nil {
		return fmt.Errorf("没有报告数据")
	}

	// 计算评分
	var scoreResult *ScoreResult
	if AIAPIEndpoint != "" && AIAPIKey != "" {
		scoreResult, _ = CalculateScoreWithAI(CurrentReportData)
	} else {
		scoreResult = CalculateScore(CurrentReportData)
	}

	data := CurrentReportData
	var sb strings.Builder

	// HTML头部
	sb.WriteString(generateHTMLHead())

	// Body开始
	sb.WriteString("<body>\n")
	sb.WriteString("<div class=\"container\">\n")

	// 标题
	sb.WriteString("<h1>压力测试报告</h1>\n")

	// AI评分卡片
	sb.WriteString(generateScoreCard(scoreResult))

	// 测试信息
	sb.WriteString(generateTestInfoSection(data))

	// 测试结果摘要
	sb.WriteString(generateResultSummarySection(data))

	// 响应时间统计
	sb.WriteString(generateResponseTimeSection(data))

	// 图表区域
	sb.WriteString(generateChartsSection(data))

	// 状态码分布
	sb.WriteString(generateStatusCodeSection(data))

	// 时间序列数据
	if len(data.TimeRecords) > 0 {
		sb.WriteString(generateTimeSeriesSection(data))
	}

	// 页脚
	sb.WriteString("<footer>\n")
	sb.WriteString(fmt.Sprintf("<p>报告生成时间: %s</p>\n", data.EndTime.Format("2006-01-02 15:04:05")))
	sb.WriteString("<p>Powered by <a href=\"https://github.com/link1st/go-stress-testing\">go-stress-testing</a></p>\n")
	sb.WriteString("</footer>\n")

	sb.WriteString("</div>\n")

	// 模态框 HTML
	sb.WriteString(generateChartModal())

	// JavaScript图表代码
	sb.WriteString(generateChartScripts(data))

	sb.WriteString("</body>\n</html>")

	// 写入文件
	err := os.WriteFile(filepath, []byte(sb.String()), 0644)
	if err != nil {
		fmt.Printf("写入HTML报告文件失败: %v\n", err)
		return err
	}

	fmt.Printf("\nHTML报告已导出到: %s\n", filepath)
	return nil
}

// generateHTMLHead 生成HTML头部
func generateHTMLHead() string {
	return `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>压力测试报告</title>
<script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
<style>
:root {
  --primary-color: #3498db;
  --success-color: #27ae60;
  --warning-color: #f39c12;
  --danger-color: #e74c3c;
  --bg-color: #f5f7fa;
  --card-bg: #ffffff;
  --text-color: #2c3e50;
  --border-color: #e1e8ed;
}
* { box-sizing: border-box; margin: 0; padding: 0; }
body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  background: var(--bg-color);
  color: var(--text-color);
  line-height: 1.6;
  padding: 20px;
}
.container { max-width: 1200px; margin: 0 auto; }
h1 {
  text-align: center;
  color: var(--primary-color);
  margin-bottom: 30px;
  font-size: 2em;
}
h2 {
  color: var(--text-color);
  border-bottom: 2px solid var(--primary-color);
  padding-bottom: 10px;
  margin: 30px 0 20px 0;
}
.score-card {
  background: linear-gradient(135deg, var(--primary-color), #2980b9);
  color: white;
  border-radius: 16px;
  padding: 30px;
  margin-bottom: 30px;
  display: flex;
  flex-wrap: wrap;
  gap: 30px;
  box-shadow: 0 10px 30px rgba(52, 152, 219, 0.3);
}
.score-main {
  display: flex;
  align-items: center;
  gap: 20px;
}
.score-circle {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  background: rgba(255,255,255,0.2);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border: 4px solid rgba(255,255,255,0.5);
}
.score-number { font-size: 48px; font-weight: bold; line-height: 1; }
.score-label { font-size: 14px; opacity: 0.9; }
.grade-badge {
  font-size: 64px;
  font-weight: bold;
  text-shadow: 2px 2px 4px rgba(0,0,0,0.2);
}
.grade-A { color: #2ecc71; }
.grade-B { color: #9acd32; }
.grade-C { color: #f1c40f; }
.grade-D { color: #e67e22; }
.grade-F { color: #e74c3c; }
.score-details { flex: 1; min-width: 300px; }
.score-item {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  border-bottom: 1px solid rgba(255,255,255,0.2);
}
.score-item:last-child { border-bottom: none; }
.suggestions {
  width: 100%;
  background: rgba(255,255,255,0.1);
  border-radius: 8px;
  padding: 15px;
  margin-top: 10px;
}
.suggestions h3 { margin-bottom: 10px; font-size: 16px; }
.suggestions ul { padding-left: 20px; }
.suggestions li { margin: 5px 0; font-size: 14px; }
.card {
  background: var(--card-bg);
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 20px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.05);
}
table {
  width: 100%;
  border-collapse: collapse;
  margin: 10px 0;
}
th, td {
  padding: 12px 15px;
  text-align: left;
  border-bottom: 1px solid var(--border-color);
}
th {
  background: var(--bg-color);
  font-weight: 600;
  color: var(--text-color);
}
tr:hover { background: #f8f9fa; }
.charts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 20px;
  margin: 20px 0;
}
.chart-container {
  background: var(--card-bg);
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.05);
}
.chart-container h3 {
  margin-bottom: 15px;
  color: var(--text-color);
  font-size: 16px;
}
.chart-wrapper {
  position: relative;
  height: 300px;
  cursor: pointer;
}
.chart-wrapper:hover {
  opacity: 0.9;
}
.chart-wrapper::after {
  content: '点击放大';
  position: absolute;
  top: 8px;
  right: 8px;
  background: rgba(0,0,0,0.6);
  color: white;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  opacity: 0;
  transition: opacity 0.2s;
}
.chart-wrapper:hover::after {
  opacity: 1;
}
/* 图表放大模态框 */
.chart-modal {
  display: none;
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0,0,0,0.85);
  z-index: 1000;
  justify-content: center;
  align-items: center;
  cursor: pointer;
}
.chart-modal.active {
  display: flex;
}
.chart-modal-content {
  background: white;
  border-radius: 12px;
  padding: 30px;
  width: 85vw;
  height: 80vh;
  max-width: 1400px;
  position: relative;
  display: flex;
  flex-direction: column;
}
.chart-modal-content h3 {
  margin-bottom: 15px;
  color: var(--text-color);
  flex-shrink: 0;
}
.chart-modal-content canvas {
  flex: 1;
  width: 100% !important;
  height: auto !important;
  min-height: 0;
}
.chart-modal-close {
  position: absolute;
  top: -15px;
  right: -15px;
  width: 36px;
  height: 36px;
  background: var(--danger-color);
  color: white;
  border: none;
  border-radius: 50%;
  font-size: 20px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 10px rgba(0,0,0,0.3);
}
.chart-modal-close:hover {
  background: #c0392b;
}
.metric-value { font-weight: bold; color: var(--primary-color); }
.success { color: var(--success-color); }
.warning { color: var(--warning-color); }
.danger { color: var(--danger-color); }
footer {
  text-align: center;
  padding: 30px;
  color: #7f8c8d;
  font-size: 14px;
}
footer a { color: var(--primary-color); text-decoration: none; }
@media (max-width: 768px) {
  .score-card { flex-direction: column; align-items: center; }
  .charts-grid { grid-template-columns: 1fr; }
  .chart-wrapper { height: 250px; }
}
</style>
</head>
`
}

// generateScoreCard 生成评分卡片
func generateScoreCard(score *ScoreResult) string {
	var sb strings.Builder

	sb.WriteString("<div class=\"score-card\">\n")
	sb.WriteString("<div class=\"score-main\">\n")

	// 分数圆圈
	sb.WriteString("<div class=\"score-circle\">\n")
	sb.WriteString(fmt.Sprintf("<span class=\"score-number\">%d</span>\n", score.TotalScore))
	sb.WriteString("<span class=\"score-label\">总分</span>\n")
	sb.WriteString("</div>\n")

	// 评级徽章
	sb.WriteString(fmt.Sprintf("<span class=\"grade-badge grade-%s\">%s</span>\n", score.Grade, score.Grade))

	sb.WriteString("</div>\n")

	// 各项得分
	sb.WriteString("<div class=\"score-details\">\n")
	sb.WriteString(fmt.Sprintf("<div class=\"score-item\"><span>成功率</span><span>%d/30</span></div>\n", score.SuccessRateScore))
	sb.WriteString(fmt.Sprintf("<div class=\"score-item\"><span>QPS吞吐</span><span>%d/25</span></div>\n", score.QPSScore))
	sb.WriteString(fmt.Sprintf("<div class=\"score-item\"><span>响应速度</span><span>%d/20</span></div>\n", score.AvgTimeScore))
	sb.WriteString(fmt.Sprintf("<div class=\"score-item\"><span>稳定性</span><span>%d/15</span></div>\n", score.TP99Score))
	sb.WriteString(fmt.Sprintf("<div class=\"score-item\"><span>错误码</span><span>%d/10</span></div>\n", score.ErrorCodeScore))
	sb.WriteString("</div>\n")

	// 建议
	if len(score.Suggestions) > 0 {
		sb.WriteString("<div class=\"suggestions\">\n")
		sb.WriteString("<h3>优化建议</h3>\n")
		sb.WriteString("<ul>\n")
		for _, suggestion := range score.Suggestions {
			sb.WriteString(fmt.Sprintf("<li>%s</li>\n", suggestion))
		}
		sb.WriteString("</ul>\n")
		sb.WriteString("</div>\n")
	}

	sb.WriteString("</div>\n")
	return sb.String()
}

// generateTestInfoSection 生成测试信息部分
func generateTestInfoSection(data *ReportData) string {
	var sb strings.Builder
	sb.WriteString("<h2>测试信息</h2>\n")
	sb.WriteString("<div class=\"card\">\n")
	sb.WriteString("<table>\n")
	sb.WriteString("<tr><th>项目</th><th>值</th></tr>\n")
	sb.WriteString(fmt.Sprintf("<tr><td>压测地址</td><td><code>%s</code></td></tr>\n", data.URL))
	sb.WriteString(fmt.Sprintf("<tr><td>并发数</td><td class=\"metric-value\">%d</td></tr>\n", data.Concurrency))
	sb.WriteString(fmt.Sprintf("<tr><td>总请求数</td><td class=\"metric-value\">%d</td></tr>\n", data.TotalRequests))
	sb.WriteString(fmt.Sprintf("<tr><td>开始时间</td><td>%s</td></tr>\n", data.StartTime.Format("2006-01-02 15:04:05")))
	sb.WriteString(fmt.Sprintf("<tr><td>结束时间</td><td>%s</td></tr>\n", data.EndTime.Format("2006-01-02 15:04:05")))
	sb.WriteString(fmt.Sprintf("<tr><td>总耗时</td><td class=\"metric-value\">%.3f 秒</td></tr>\n", data.TotalTime))
	sb.WriteString("</table>\n")
	sb.WriteString("</div>\n")
	return sb.String()
}

// generateResultSummarySection 生成结果摘要部分
func generateResultSummarySection(data *ReportData) string {
	var sb strings.Builder
	successRate := float64(0)
	if data.TotalRequests > 0 {
		successRate = float64(data.SuccessNum) / float64(data.TotalRequests) * 100
	}

	successClass := "success"
	if successRate < 99 {
		successClass = "warning"
	}
	if successRate < 95 {
		successClass = "danger"
	}

	sb.WriteString("<h2>测试结果摘要</h2>\n")
	sb.WriteString("<div class=\"card\">\n")
	sb.WriteString("<table>\n")
	sb.WriteString("<tr><th>指标</th><th>值</th></tr>\n")
	sb.WriteString(fmt.Sprintf("<tr><td>成功请求数</td><td class=\"success\">%d</td></tr>\n", data.SuccessNum))
	sb.WriteString(fmt.Sprintf("<tr><td>失败请求数</td><td class=\"%s\">%d</td></tr>\n",
		map[bool]string{true: "", false: "danger"}[data.FailureNum == 0], data.FailureNum))
	sb.WriteString(fmt.Sprintf("<tr><td>成功率</td><td class=\"%s\">%.2f%%</td></tr>\n", successClass, successRate))
	sb.WriteString(fmt.Sprintf("<tr><td>QPS</td><td class=\"metric-value\">%.2f</td></tr>\n", data.QPS))
	sb.WriteString(fmt.Sprintf("<tr><td>下载字节数</td><td>%d</td></tr>\n", data.ReceivedBytes))
	sb.WriteString("</table>\n")
	sb.WriteString("</div>\n")
	return sb.String()
}

// generateResponseTimeSection 生成响应时间统计部分
func generateResponseTimeSection(data *ReportData) string {
	var sb strings.Builder
	sb.WriteString("<h2>响应时间统计</h2>\n")
	sb.WriteString("<div class=\"card\">\n")
	sb.WriteString("<table>\n")
	sb.WriteString("<tr><th>指标</th><th>值 (ms)</th></tr>\n")
	sb.WriteString(fmt.Sprintf("<tr><td>最小响应时间</td><td class=\"success\">%.2f</td></tr>\n", data.MinTime))
	sb.WriteString(fmt.Sprintf("<tr><td>最大响应时间</td><td class=\"warning\">%.2f</td></tr>\n", data.MaxTime))
	sb.WriteString(fmt.Sprintf("<tr><td>平均响应时间</td><td class=\"metric-value\">%.2f</td></tr>\n", data.AvgTime))
	sb.WriteString(fmt.Sprintf("<tr><td>TP90</td><td>%.2f</td></tr>\n", data.TP90))
	sb.WriteString(fmt.Sprintf("<tr><td>TP95</td><td>%.2f</td></tr>\n", data.TP95))
	sb.WriteString(fmt.Sprintf("<tr><td>TP99</td><td>%.2f</td></tr>\n", data.TP99))
	sb.WriteString("</table>\n")
	sb.WriteString("</div>\n")
	return sb.String()
}

// generateChartsSection 生成图表区域
func generateChartsSection(data *ReportData) string {
	var sb strings.Builder
	sb.WriteString("<h2>可视化分析</h2>\n")
	sb.WriteString("<div class=\"charts-grid\">\n")

	// 响应时间分布图
	sb.WriteString("<div class=\"chart-container\">\n")
	sb.WriteString("<h3>响应时间分布</h3>\n")
	sb.WriteString("<div class=\"chart-wrapper\"><canvas id=\"responseTimeChart\"></canvas></div>\n")
	sb.WriteString("</div>\n")

	// 状态码饼图
	sb.WriteString("<div class=\"chart-container\">\n")
	sb.WriteString("<h3>状态码分布</h3>\n")
	sb.WriteString("<div class=\"chart-wrapper\"><canvas id=\"statusCodeChart\"></canvas></div>\n")
	sb.WriteString("</div>\n")

	// 时间序列图表 (如果有时间序列数据)
	if len(data.TimeRecords) > 0 {
		// QPS时间曲线
		sb.WriteString("<div class=\"chart-container\">\n")
		sb.WriteString("<h3>QPS 随时间变化</h3>\n")
		sb.WriteString("<div class=\"chart-wrapper\"><canvas id=\"qpsChart\"></canvas></div>\n")
		sb.WriteString("</div>\n")

		// 成功率时间曲线
		sb.WriteString("<div class=\"chart-container\">\n")
		sb.WriteString("<h3>成功率 随时间变化</h3>\n")
		sb.WriteString("<div class=\"chart-wrapper\"><canvas id=\"successRateChart\"></canvas></div>\n")
		sb.WriteString("</div>\n")

		// 接口耗时时间曲线
		sb.WriteString("<div class=\"chart-container\">\n")
		sb.WriteString("<h3>接口耗时 随时间变化</h3>\n")
		sb.WriteString("<div class=\"chart-wrapper\"><canvas id=\"latencyChart\"></canvas></div>\n")
		sb.WriteString("</div>\n")

		// 错误码时间曲线
		sb.WriteString("<div class=\"chart-container\">\n")
		sb.WriteString("<h3>错误码 随时间变化</h3>\n")
		sb.WriteString("<div class=\"chart-wrapper\"><canvas id=\"errorCodeChart\"></canvas></div>\n")
		sb.WriteString("</div>\n")
	}

	sb.WriteString("</div>\n")
	return sb.String()
}

// generateStatusCodeSection 生成状态码分布部分
func generateStatusCodeSection(data *ReportData) string {
	var sb strings.Builder
	sb.WriteString("<h2>状态码分布</h2>\n")
	sb.WriteString("<div class=\"card\">\n")
	sb.WriteString("<table>\n")
	sb.WriteString("<tr><th>状态码</th><th>次数</th><th>占比</th></tr>\n")

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
		codeClass := ""
		if code >= 400 && code < 500 {
			codeClass = "warning"
		} else if code >= 500 {
			codeClass = "danger"
		} else if code == 200 {
			codeClass = "success"
		}
		sb.WriteString(fmt.Sprintf("<tr><td class=\"%s\">%d</td><td>%d</td><td>%.2f%%</td></tr>\n",
			codeClass, code, count, percentage))
	}

	sb.WriteString("</table>\n")
	sb.WriteString("</div>\n")
	return sb.String()
}

// generateTimeSeriesSection 生成时间序列数据部分
func generateTimeSeriesSection(data *ReportData) string {
	var sb strings.Builder
	sb.WriteString("<h2>时间序列数据</h2>\n")
	sb.WriteString("<div class=\"card\">\n")
	sb.WriteString("<table>\n")
	sb.WriteString("<tr><th>时间</th><th>并发数</th><th>成功数</th><th>失败数</th><th>成功率</th><th>QPS</th><th>平均(ms)</th><th>最长(ms)</th><th>最短(ms)</th></tr>\n")

	for _, record := range data.TimeRecords {
		sb.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%d</td><td>%d</td><td>%d</td><td>%.2f%%</td><td>%.2f</td><td>%.2f</td><td>%.2f</td><td>%.2f</td></tr>\n",
			record.Timestamp.Format("15:04:05"),
			record.Concurrent, record.Success, record.Failure, record.SuccessRate,
			record.QPS, record.AvgTime, record.MaxTime, record.MinTime))
	}

	sb.WriteString("</table>\n")
	sb.WriteString("</div>\n")
	return sb.String()
}

// generateChartModal 生成图表放大模态框
func generateChartModal() string {
	return `<!-- 图表放大模态框 -->
<div id="chartModal" class="chart-modal" onclick="closeChartModal(event)">
  <div class="chart-modal-content" onclick="event.stopPropagation()">
    <button class="chart-modal-close" onclick="closeChartModal()">&times;</button>
    <h3 id="modalChartTitle"></h3>
    <canvas id="modalChart"></canvas>
  </div>
</div>
`
}

// generateModalScript 生成模态框 JavaScript 代码
func generateModalScript() string {
	return `
// 存储所有图表配置
const chartConfigs = {};
let modalChart = null;

// 获取图表实例
function getChartInstance(canvasId) {
  const canvas = document.getElementById(canvasId);
  if (canvas) {
    return Chart.getChart(canvas);
  }
  return null;
}

// 初始化图表点击事件
function initChartClickHandlers() {
  const chartWrappers = document.querySelectorAll('.chart-wrapper');
  chartWrappers.forEach(wrapper => {
    wrapper.addEventListener('click', function() {
      const canvas = this.querySelector('canvas');
      if (canvas) {
        openChartModal(canvas.id);
      }
    });
  });
}

// 获取图表标题
function getChartTitle(canvasId) {
  const titles = {
    'responseTimeChart': '响应时间分布',
    'statusCodeChart': '状态码分布',
    'qpsChart': 'QPS 随时间变化',
    'successRateChart': '成功率 随时间变化',
    'latencyChart': '接口耗时 随时间变化',
    'errorCodeChart': '错误码 随时间变化'
  };
  return titles[canvasId] || '图表';
}

// 深拷贝图表数据（只拷贝可序列化的部分）
function cloneChartData(data) {
  if (!data) return { labels: [], datasets: [] };
  return {
    labels: data.labels ? [...data.labels] : [],
    datasets: data.datasets ? data.datasets.map(ds => ({
      label: ds.label,
      data: [...ds.data],
      backgroundColor: ds.backgroundColor,
      borderColor: ds.borderColor,
      borderWidth: ds.borderWidth,
      fill: ds.fill,
      tension: ds.tension,
      pointRadius: ds.pointRadius,
      borderDash: ds.borderDash
    })) : []
  };
}

// 打开模态框
function openChartModal(canvasId) {
  const originalChart = getChartInstance(canvasId);
  if (!originalChart) return;

  const modal = document.getElementById('chartModal');
  const modalCanvas = document.getElementById('modalChart');
  const modalTitle = document.getElementById('modalChartTitle');

  // 设置标题
  modalTitle.textContent = getChartTitle(canvasId);

  // 销毁旧的模态框图表
  if (modalChart) {
    modalChart.destroy();
    modalChart = null;
  }

  // 手动构建新配置（避免循环引用问题）
  const chartType = originalChart.config.type;
  const chartData = cloneChartData(originalChart.data);

  const config = {
    type: chartType,
    data: chartData,
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: chartType === 'doughnut' ? { position: 'right' } : { display: chartType !== 'bar' }
      }
    }
  };

  // 为非饼图添加坐标轴配置
  if (chartType !== 'doughnut') {
    config.options.scales = {
      x: { beginAtZero: true },
      y: { beginAtZero: true }
    };
  }

  // 清除之前的样式，让 CSS flex 布局控制大小
  modalCanvas.style.width = '';
  modalCanvas.style.height = '';

  modalChart = new Chart(modalCanvas, config);

  // 显示模态框
  modal.classList.add('active');
  document.body.style.overflow = 'hidden';
}

// 关闭模态框
function closeChartModal(event) {
  if (event && event.target !== event.currentTarget) return;

  const modal = document.getElementById('chartModal');
  modal.classList.remove('active');
  document.body.style.overflow = '';

  // 销毁模态框图表
  if (modalChart) {
    modalChart.destroy();
    modalChart = null;
  }
}

// ESC 键关闭模态框
document.addEventListener('keydown', function(e) {
  if (e.key === 'Escape') {
    closeChartModal();
  }
});

// 页面加载后初始化点击事件
document.addEventListener('DOMContentLoaded', initChartClickHandlers);
// 如果 DOM 已经加载完成，立即初始化
if (document.readyState !== 'loading') {
  initChartClickHandlers();
}
`
}

// generateChartScripts 生成图表JavaScript代码
func generateChartScripts(data *ReportData) string {
	var sb strings.Builder
	sb.WriteString("<script>\n")

	// 响应时间分布数据
	responseTimeData := generateResponseTimeDistribution(data)
	responseTimeJSON, _ := json.Marshal(responseTimeData)

	// 状态码数据
	statusCodes := make([]int, 0, len(data.ErrorCodeMap))
	for code := range data.ErrorCodeMap {
		statusCodes = append(statusCodes, code)
	}
	sort.Ints(statusCodes)

	var statusLabels []string
	var statusValues []int
	var statusColors []string
	for _, code := range statusCodes {
		statusLabels = append(statusLabels, fmt.Sprintf("%d", code))
		statusValues = append(statusValues, data.ErrorCodeMap[code])
		if code >= 500 {
			statusColors = append(statusColors, "#e74c3c")
		} else if code >= 400 {
			statusColors = append(statusColors, "#f39c12")
		} else {
			statusColors = append(statusColors, "#27ae60")
		}
	}
	statusLabelsJSON, _ := json.Marshal(statusLabels)
	statusValuesJSON, _ := json.Marshal(statusValues)
	statusColorsJSON, _ := json.Marshal(statusColors)

	// 响应时间分布图
	sb.WriteString(fmt.Sprintf(`
// 响应时间分布图
const rtData = %s;
new Chart(document.getElementById('responseTimeChart'), {
  type: 'bar',
  data: {
    labels: rtData.labels,
    datasets: [{
      label: '请求数',
      data: rtData.values,
      backgroundColor: 'rgba(52, 152, 219, 0.7)',
      borderColor: 'rgba(52, 152, 219, 1)',
      borderWidth: 1
    }]
  },
  options: {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: { display: false }
    },
    scales: {
      x: { title: { display: true, text: '响应时间 (ms)' } },
      y: { title: { display: true, text: '请求数' }, beginAtZero: true }
    }
  }
});
`, string(responseTimeJSON)))

	// 状态码饼图
	sb.WriteString(fmt.Sprintf(`
// 状态码饼图
new Chart(document.getElementById('statusCodeChart'), {
  type: 'doughnut',
  data: {
    labels: %s,
    datasets: [{
      data: %s,
      backgroundColor: %s,
      borderWidth: 2,
      borderColor: '#fff'
    }]
  },
  options: {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: { position: 'right' }
    }
  }
});
`, string(statusLabelsJSON), string(statusValuesJSON), string(statusColorsJSON)))

	// 时间序列图表 (如果有时间序列数据)
	if len(data.TimeRecords) > 0 {
		var timeLabels []string
		var qpsValues []float64
		var successRateValues []float64
		var avgTimeValues []float64
		var maxTimeValues []float64
		var minTimeValues []float64

		// 收集所有错误码（排除成功状态码）
		allErrorCodes := make(map[int]bool)
		for _, record := range data.TimeRecords {
			for code := range record.ErrorCodes {
				// 排除成功状态码，只记录错误码
				if code != SuccessCode && code != 0 {
					allErrorCodes[code] = true
				}
			}
		}

		// 排序错误码
		var errorCodeList []int
		for code := range allErrorCodes {
			errorCodeList = append(errorCodeList, code)
		}
		sort.Ints(errorCodeList)

		// 构建错误码数据集
		errorCodeData := make(map[int][]int)
		for _, code := range errorCodeList {
			errorCodeData[code] = make([]int, 0)
		}

		for _, record := range data.TimeRecords {
			// 时间格式: 时:分:秒
			timeLabels = append(timeLabels, record.Timestamp.Format("15:04:05"))
			qpsValues = append(qpsValues, record.QPS)
			successRateValues = append(successRateValues, record.SuccessRate)
			avgTimeValues = append(avgTimeValues, record.AvgTime)
			maxTimeValues = append(maxTimeValues, record.MaxTime)
			minTimeValues = append(minTimeValues, record.MinTime)

			// 收集每个时间点的错误码数量
			for _, code := range errorCodeList {
				count := 0
				if v, ok := record.ErrorCodes[code]; ok {
					count = v
				}
				errorCodeData[code] = append(errorCodeData[code], count)
			}
		}

		timeLabelsJSON, _ := json.Marshal(timeLabels)
		qpsValuesJSON, _ := json.Marshal(qpsValues)
		successRateValuesJSON, _ := json.Marshal(successRateValues)
		avgTimeValuesJSON, _ := json.Marshal(avgTimeValues)
		maxTimeValuesJSON, _ := json.Marshal(maxTimeValues)
		minTimeValuesJSON, _ := json.Marshal(minTimeValues)

		// QPS时间曲线
		sb.WriteString(fmt.Sprintf(`
// QPS 随时间变化
new Chart(document.getElementById('qpsChart'), {
  type: 'line',
  data: {
    labels: %s,
    datasets: [{
      label: 'QPS',
      data: %s,
      borderColor: '#3498db',
      backgroundColor: 'rgba(52, 152, 219, 0.1)',
      fill: true,
      tension: 0.3,
      pointRadius: 3
    }]
  },
  options: {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: { display: false }
    },
    scales: {
      x: {
        title: { display: true, text: '时间' },
        ticks: { maxRotation: 45, minRotation: 45 }
      },
      y: { title: { display: true, text: 'QPS' }, beginAtZero: true }
    }
  }
});
`, string(timeLabelsJSON), string(qpsValuesJSON)))

		// 成功率时间曲线
		sb.WriteString(fmt.Sprintf(`
// 成功率 随时间变化
new Chart(document.getElementById('successRateChart'), {
  type: 'line',
  data: {
    labels: %s,
    datasets: [{
      label: '成功率(%%)',
      data: %s,
      borderColor: '#27ae60',
      backgroundColor: 'rgba(39, 174, 96, 0.1)',
      fill: true,
      tension: 0.3,
      pointRadius: 3
    }]
  },
  options: {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: { display: false }
    },
    scales: {
      x: {
        title: { display: true, text: '时间' },
        ticks: { maxRotation: 45, minRotation: 45 }
      },
      y: {
        title: { display: true, text: '成功率(%%)' },
        beginAtZero: true,
        max: 100
      }
    }
  }
});
`, string(timeLabelsJSON), string(successRateValuesJSON)))

		// 接口耗时时间曲线
		sb.WriteString(fmt.Sprintf(`
// 接口耗时 随时间变化
new Chart(document.getElementById('latencyChart'), {
  type: 'line',
  data: {
    labels: %s,
    datasets: [
      {
        label: '平均耗时(ms)',
        data: %s,
        borderColor: '#3498db',
        backgroundColor: 'transparent',
        tension: 0.3,
        pointRadius: 3
      },
      {
        label: '最大耗时(ms)',
        data: %s,
        borderColor: '#e74c3c',
        backgroundColor: 'transparent',
        tension: 0.3,
        pointRadius: 3,
        borderDash: [5, 5]
      },
      {
        label: '最小耗时(ms)',
        data: %s,
        borderColor: '#27ae60',
        backgroundColor: 'transparent',
        tension: 0.3,
        pointRadius: 3,
        borderDash: [5, 5]
      }
    ]
  },
  options: {
    responsive: true,
    maintainAspectRatio: false,
    scales: {
      x: {
        title: { display: true, text: '时间' },
        ticks: { maxRotation: 45, minRotation: 45 }
      },
      y: { title: { display: true, text: '耗时(ms)' }, beginAtZero: true }
    }
  }
});
`, string(timeLabelsJSON), string(avgTimeValuesJSON), string(maxTimeValuesJSON), string(minTimeValuesJSON)))

		// 错误码时间曲线
		var errorCodeDatasets []string
		colors := []string{"#27ae60", "#3498db", "#f39c12", "#e74c3c", "#9b59b6", "#1abc9c", "#34495e", "#e91e63"}
		for i, code := range errorCodeList {
			dataJSON, _ := json.Marshal(errorCodeData[code])
			color := colors[i%len(colors)]
			errorCodeDatasets = append(errorCodeDatasets, fmt.Sprintf(`{
        label: '%d',
        data: %s,
        borderColor: '%s',
        backgroundColor: 'transparent',
        tension: 0.3,
        pointRadius: 3
      }`, code, string(dataJSON), color))
		}

		sb.WriteString(fmt.Sprintf(`
// 错误码 随时间变化
new Chart(document.getElementById('errorCodeChart'), {
  type: 'line',
  data: {
    labels: %s,
    datasets: [%s]
  },
  options: {
    responsive: true,
    maintainAspectRatio: false,
    scales: {
      x: {
        title: { display: true, text: '时间' },
        ticks: { maxRotation: 45, minRotation: 45 }
      },
      y: { title: { display: true, text: '请求数' }, beginAtZero: true }
    }
  }
});
`, string(timeLabelsJSON), strings.Join(errorCodeDatasets, ",\n      ")))
	}

	// 添加模态框 JavaScript
	sb.WriteString(generateModalScript())

	sb.WriteString("</script>\n")
	return sb.String()
}

// ResponseTimeDistribution 响应时间分布数据
type ResponseTimeDistribution struct {
	Labels []string `json:"labels"`
	Values []int    `json:"values"`
}

// generateResponseTimeDistribution 生成响应时间分布数据
func generateResponseTimeDistribution(data *ReportData) ResponseTimeDistribution {
	if len(data.RequestTimeList) == 0 {
		return ResponseTimeDistribution{
			Labels: []string{"无数据"},
			Values: []int{0},
		}
	}

	// 定义时间区间 (毫秒)
	buckets := []struct {
		label string
		max   uint64 // 纳秒
	}{
		{"0-50", 50 * 1e6},
		{"50-100", 100 * 1e6},
		{"100-200", 200 * 1e6},
		{"200-500", 500 * 1e6},
		{"500-1000", 1000 * 1e6},
		{"1000+", ^uint64(0)},
	}

	counts := make([]int, len(buckets))

	for _, t := range data.RequestTimeList {
		for i, bucket := range buckets {
			if t <= bucket.max {
				counts[i]++
				break
			}
		}
	}

	labels := make([]string, len(buckets))
	for i, bucket := range buckets {
		labels[i] = bucket.label
	}

	return ResponseTimeDistribution{
		Labels: labels,
		Values: counts,
	}
}
