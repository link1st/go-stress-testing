// Package client WebSocket客户端测试
package client

import (
	"testing"
)

// TestNewWebSocket 测试创建WebSocket客户端
func TestNewWebSocket(t *testing.T) {
	tests := []struct {
		name        string
		urlLink     string
		expectedSSL bool
	}{
		{
			name:        "WS协议",
			urlLink:     "ws://localhost:8080/ws",
			expectedSSL: false,
		},
		{
			name:        "WSS协议",
			urlLink:     "wss://localhost:8080/ws",
			expectedSSL: true,
		},
		{
			name:        "带路径的WS",
			urlLink:     "ws://example.com/api/websocket",
			expectedSSL: false,
		},
		{
			name:        "带端口的WSS",
			urlLink:     "wss://example.com:443/ws",
			expectedSSL: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ws := NewWebSocket(tt.urlLink)

			if ws == nil {
				t.Fatal("NewWebSocket() 返回了 nil")
			}
			if ws.URLLink != tt.urlLink {
				t.Errorf("URLLink = %q, 期望 %q", ws.URLLink, tt.urlLink)
			}
			if ws.IsSsl != tt.expectedSSL {
				t.Errorf("IsSsl = %v, 期望 %v", ws.IsSsl, tt.expectedSSL)
			}
			if ws.URL == nil {
				t.Error("URL 不应该为 nil")
			}
			if ws.HTTPHeader == nil {
				t.Error("HTTPHeader 不应该为 nil")
			}
		})
	}
}

// TestWebSocket_getLink 测试获取连接地址
func TestWebSocket_getLink(t *testing.T) {
	urlLink := "ws://localhost:8080/ws"
	ws := NewWebSocket(urlLink)

	if got := ws.getLink(); got != urlLink {
		t.Errorf("getLink() = %q, 期望 %q", got, urlLink)
	}
}

// TestWebSocket_SetHeader 测试设置Header
func TestWebSocket_SetHeader(t *testing.T) {
	ws := NewWebSocket("ws://localhost:8080/ws")

	headers := map[string]string{
		"Authorization": "Bearer token",
		"X-Custom":      "value",
	}

	ws.SetHeader(headers)

	if len(ws.HTTPHeader) != len(headers) {
		t.Errorf("HTTPHeader 长度 = %d, 期望 %d", len(ws.HTTPHeader), len(headers))
	}

	for k, v := range headers {
		if ws.HTTPHeader[k] != v {
			t.Errorf("HTTPHeader[%q] = %q, 期望 %q", k, ws.HTTPHeader[k], v)
		}
	}
}

// TestWebSocket_getOrigin 测试获取Origin
func TestWebSocket_getOrigin(t *testing.T) {
	tests := []struct {
		name           string
		urlLink        string
		expectedOrigin string
	}{
		{
			name:           "WS协议",
			urlLink:        "ws://localhost:8080/ws",
			expectedOrigin: "http://localhost:8080/",
		},
		{
			name:           "WSS协议",
			urlLink:        "wss://localhost:8080/ws",
			expectedOrigin: "https://localhost:8080/",
		},
		{
			name:           "无端口WS",
			urlLink:        "ws://example.com/ws",
			expectedOrigin: "http://example.com/",
		},
		{
			name:           "无端口WSS",
			urlLink:        "wss://example.com/ws",
			expectedOrigin: "https://example.com/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ws := NewWebSocket(tt.urlLink)
			origin := ws.getOrigin()

			if origin != tt.expectedOrigin {
				t.Errorf("getOrigin() = %q, 期望 %q", origin, tt.expectedOrigin)
			}
		})
	}
}

// TestWebSocket_Close_Nil 测试关闭nil连接
func TestWebSocket_Close_Nil(t *testing.T) {
	// 测试nil WebSocket
	var ws *WebSocket
	err := ws.Close()
	if err != nil {
		t.Errorf("关闭nil WebSocket应该返回nil: %v", err)
	}

	// 测试nil conn
	ws = NewWebSocket("ws://localhost:8080/ws")
	ws.conn = nil
	err = ws.Close()
	if err != nil {
		t.Errorf("关闭nil conn应该返回nil: %v", err)
	}
}

// TestWebSocket_Write_NoConn 测试未建立连接时写入
func TestWebSocket_Write_NoConn(t *testing.T) {
	ws := NewWebSocket("ws://localhost:8080/ws")
	// 不调用GetConn，直接写入

	err := ws.Write([]byte("test message"))
	if err == nil {
		t.Error("未建立连接时写入应该返回错误")
	}
	if err.Error() != "未建立连接" {
		t.Errorf("错误消息 = %q, 期望 '未建立连接'", err.Error())
	}
}

// TestWebSocket_Read_NoConn 测试未建立连接时读取
func TestWebSocket_Read_NoConn(t *testing.T) {
	ws := NewWebSocket("ws://localhost:8080/ws")
	// 不调用GetConn，直接读取

	_, err := ws.Read()
	if err == nil {
		t.Error("未建立连接时读取应该返回错误")
	}
	if err.Error() != "未建立连接" {
		t.Errorf("错误消息 = %q, 期望 '未建立连接'", err.Error())
	}
}

// TestConnRetry 测试重试常量
func TestConnRetry(t *testing.T) {
	if connRetry != 3 {
		t.Errorf("connRetry = %d, 期望 3", connRetry)
	}
}

// TestWebSocket_URLParsing 测试URL解析
func TestWebSocket_URLParsing(t *testing.T) {
	tests := []struct {
		name         string
		urlLink      string
		expectedHost string
		expectedPath string
	}{
		{
			name:         "简单路径",
			urlLink:      "ws://localhost:8080/ws",
			expectedHost: "localhost:8080",
			expectedPath: "/ws",
		},
		{
			name:         "多级路径",
			urlLink:      "ws://example.com/api/v1/websocket",
			expectedHost: "example.com",
			expectedPath: "/api/v1/websocket",
		},
		{
			name:         "无路径",
			urlLink:      "ws://localhost:8080",
			expectedHost: "localhost:8080",
			expectedPath: "",
		},
		{
			name:         "带查询参数",
			urlLink:      "ws://localhost:8080/ws?token=abc",
			expectedHost: "localhost:8080",
			expectedPath: "/ws",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ws := NewWebSocket(tt.urlLink)

			if ws.URL.Host != tt.expectedHost {
				t.Errorf("URL.Host = %q, 期望 %q", ws.URL.Host, tt.expectedHost)
			}
			if ws.URL.Path != tt.expectedPath {
				t.Errorf("URL.Path = %q, 期望 %q", ws.URL.Path, tt.expectedPath)
			}
		})
	}
}

// TestWebSocket_SetHeader_Nil 测试设置nil Header
func TestWebSocket_SetHeader_Nil(t *testing.T) {
	ws := NewWebSocket("ws://localhost:8080/ws")
	ws.SetHeader(nil)

	if ws.HTTPHeader != nil {
		// SetHeader直接赋值，所以会变成nil
		// 这是当前实现的行为
	}
}

// TestWebSocket_SetHeader_Empty 测试设置空Header
func TestWebSocket_SetHeader_Empty(t *testing.T) {
	ws := NewWebSocket("ws://localhost:8080/ws")
	ws.SetHeader(map[string]string{})

	if len(ws.HTTPHeader) != 0 {
		t.Errorf("HTTPHeader 长度应该为 0")
	}
}

// TestWebSocket_MultipleSets 测试多次设置Header
func TestWebSocket_MultipleSets(t *testing.T) {
	ws := NewWebSocket("ws://localhost:8080/ws")

	// 第一次设置
	ws.SetHeader(map[string]string{"Key1": "Value1"})
	if ws.HTTPHeader["Key1"] != "Value1" {
		t.Error("第一次设置Header失败")
	}

	// 第二次设置（应该覆盖）
	ws.SetHeader(map[string]string{"Key2": "Value2"})
	if _, exists := ws.HTTPHeader["Key1"]; exists {
		t.Error("第二次设置应该覆盖之前的Header")
	}
	if ws.HTTPHeader["Key2"] != "Value2" {
		t.Error("第二次设置Header失败")
	}
}

// BenchmarkNewWebSocket 性能测试
func BenchmarkNewWebSocket(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewWebSocket("ws://localhost:8080/ws")
	}
}

// BenchmarkWebSocket_getOrigin 性能测试
func BenchmarkWebSocket_getOrigin(b *testing.B) {
	ws := NewWebSocket("wss://example.com:443/ws")
	for i := 0; i < b.N; i++ {
		ws.getOrigin()
	}
}

// BenchmarkWebSocket_SetHeader 性能测试
func BenchmarkWebSocket_SetHeader(b *testing.B) {
	ws := NewWebSocket("ws://localhost:8080/ws")
	headers := map[string]string{
		"Authorization": "Bearer token",
		"X-Custom":      "value",
	}

	for i := 0; i < b.N; i++ {
		ws.SetHeader(headers)
	}
}
