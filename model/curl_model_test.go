// Package model 数据模型
package model

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// TestCurl 测试函数
func TestCurl(t *testing.T) {
	// ../curl.txt
	c, err := ParseTheFile("../curl/post.curl.txt")
	fmt.Println(c, err)

	if err != nil {
		return
	}
	fmt.Printf("curl:%s \n", c.String())
	fmt.Printf("url:%s \n", c.GetURL())
	fmt.Printf("method:%s \n", c.GetMethod())
	fmt.Printf("body:%v \n", c.GetBody())
	fmt.Printf("body string:%v \n", c.GetBody())
	fmt.Printf("headers:%s \n", c.GetHeadersStr())
}

// TestBinaryDataSupport 测试二进制数据支持
func TestBinaryDataSupport(t *testing.T) {
	// 创建测试二进制文件
	testDir := t.TempDir()
	testFile := filepath.Join(testDir, "test.bin")
	testData := []byte{0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE, 0xFD}

	err := os.WriteFile(testFile, testData, 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 创建测试curl文件
	curlFile := filepath.Join(testDir, "test.curl.txt")
	curlContent := `curl -X POST http://localhost:8088/upload \
  -H "Content-Type: application/octet-stream" \
  --data-binary @` + testFile

	err = os.WriteFile(curlFile, []byte(curlContent), 0644)
	if err != nil {
		t.Fatalf("创建curl文件失败: %v", err)
	}

	// 解析curl文件
	curl, err := ParseTheFile(curlFile)
	if err != nil {
		t.Fatalf("解析curl文件失败: %v", err)
	}

	// 测试GetBodyBytes方法
	bodyBytes, err := curl.GetBodyBytes()
	if err != nil {
		t.Fatalf("获取二进制数据失败: %v", err)
	}

	// 验证数据
	if len(bodyBytes) != len(testData) {
		t.Errorf("数据长度不匹配: 期望 %d, 实际 %d", len(testData), len(bodyBytes))
	}

	for i, b := range testData {
		if bodyBytes[i] != b {
			t.Errorf("数据不匹配 位置 %d: 期望 0x%02X, 实际 0x%02X", i, b, bodyBytes[i])
		}
	}

	t.Logf("✓ 二进制数据读取成功，长度: %d 字节", len(bodyBytes))
}

// TestBinaryDataWithoutFile 测试不带文件的--data-binary
func TestBinaryDataWithoutFile(t *testing.T) {
	testDir := t.TempDir()
	curlFile := filepath.Join(testDir, "test.curl.txt")
	curlContent := `curl -X POST http://localhost:8088/test \
  -H "Content-Type: application/octet-stream" \
  --data-binary "plain text data"`

	err := os.WriteFile(curlFile, []byte(curlContent), 0644)
	if err != nil {
		t.Fatalf("创建curl文件失败: %v", err)
	}

	curl, err := ParseTheFile(curlFile)
	if err != nil {
		t.Fatalf("解析curl文件失败: %v", err)
	}

	bodyBytes, err := curl.GetBodyBytes()
	if err != nil {
		t.Fatalf("获取数据失败: %v", err)
	}

	expected := "plain text data"
	if string(bodyBytes) != expected {
		t.Errorf("数据不匹配: 期望 %s, 实际 %s", expected, string(bodyBytes))
	}

	t.Logf("✓ 文本数据处理成功: %s", string(bodyBytes))
}

// TestNormalDataStillWorks 测试普通--data参数仍然正常工作
func TestNormalDataStillWorks(t *testing.T) {
	testDir := t.TempDir()
	curlFile := filepath.Join(testDir, "test.curl.txt")
	curlContent := `curl -X POST http://localhost:8088/test \
  -H "Content-Type: application/x-www-form-urlencoded" \
  --data "key=value&foo=bar"`

	err := os.WriteFile(curlFile, []byte(curlContent), 0644)
	if err != nil {
		t.Fatalf("创建curl文件失败: %v", err)
	}

	curl, err := ParseTheFile(curlFile)
	if err != nil {
		t.Fatalf("解析curl文件失败: %v", err)
	}

	// 现在 --data 也会被 GetBodyBytes 处理（不带@时返回字节数组）
	bodyBytes, err := curl.GetBodyBytes()
	if err != nil {
		t.Fatalf("获取数据失败: %v", err)
	}

	expected := "key=value&foo=bar"

	// GetBodyBytes 应该返回字节数组形式的数据
	if string(bodyBytes) != expected {
		t.Errorf("GetBodyBytes 数据不匹配: 期望 %s, 实际 %s", expected, string(bodyBytes))
	}

	// GetBody 也应该正常工作
	body := curl.GetBody()
	if body != expected {
		t.Errorf("GetBody 数据不匹配: 期望 %s, 实际 %s", expected, body)
	}

	t.Logf("✓ 普通data参数仍然正常工作: GetBody=%s, GetBodyBytes=%s", body, string(bodyBytes))
}

// TestDataWithFile 测试 --data @filename 语法
func TestDataWithFile(t *testing.T) {
	testDir := t.TempDir()

	// 创建测试数据文件
	dataFile := filepath.Join(testDir, "data.txt")
	testContent := "name=John&age=30"
	err := os.WriteFile(dataFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("创建数据文件失败: %v", err)
	}

	// 创建 curl 文件，使用 --data @filename
	curlFile := filepath.Join(testDir, "test.curl.txt")
	curlContent := `curl -X POST http://localhost:8088/test \
  -H "Content-Type: application/x-www-form-urlencoded" \
  --data @` + dataFile

	err = os.WriteFile(curlFile, []byte(curlContent), 0644)
	if err != nil {
		t.Fatalf("创建curl文件失败: %v", err)
	}

	curl, err := ParseTheFile(curlFile)
	if err != nil {
		t.Fatalf("解析curl文件失败: %v", err)
	}

	// 测试 GetBodyBytes 能读取文件
	bodyBytes, err := curl.GetBodyBytes()
	if err != nil {
		t.Fatalf("获取数据失败: %v", err)
	}

	if string(bodyBytes) != testContent {
		t.Errorf("数据不匹配: 期望 %s, 实际 %s", testContent, string(bodyBytes))
	}

	t.Logf("✓ --data @filename 语法正常工作: %s", string(bodyBytes))
}

// TestDataRawWithAtSymbol 测试 --data-raw 不处理 @ 符号
func TestDataRawWithAtSymbol(t *testing.T) {
	testDir := t.TempDir()
	curlFile := filepath.Join(testDir, "test.curl.txt")

	// --data-raw 的 @ 应该被当作普通字符
	curlContent := `curl -X POST http://localhost:8088/test \
  --data-raw "@username=test@example.com"`

	err := os.WriteFile(curlFile, []byte(curlContent), 0644)
	if err != nil {
		t.Fatalf("创建curl文件失败: %v", err)
	}

	curl, err := ParseTheFile(curlFile)
	if err != nil {
		t.Fatalf("解析curl文件失败: %v", err)
	}

	// --data-raw 不在 GetBodyBytes 的处理范围内
	bodyBytes, err := curl.GetBodyBytes()
	if err != nil {
		t.Fatalf("获取数据失败: %v", err)
	}

	// 应该返回空，因为 --data-raw 不被 GetBodyBytes 处理
	if len(bodyBytes) != 0 {
		t.Errorf("--data-raw 不应该被 GetBodyBytes 处理")
	}

	// 应该通过 GetBody 获取原始字符串
	body := curl.GetBody()
	expected := "@username=test@example.com"
	if body != expected {
		t.Errorf("数据不匹配: 期望 %s, 实际 %s", expected, body)
	}

	t.Logf("✓ --data-raw 正确处理 @ 符号: %s", body)
}
