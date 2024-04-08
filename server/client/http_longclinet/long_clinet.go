// Package httplongclinet Keepalive
package httplongclinet

import (
	"crypto/tls"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/http2"

	"github.com/link1st/go-stress-testing/model"
)

var (
	mutex   sync.RWMutex
	clients = make(map[uint64]*http.Client, 0)
)

// NewClient new
func NewClient(i uint64, request *model.Request) *http.Client {
	client := getClient(i)
	if client != nil {
		return client
	}
	return setClient(i, request)
}

func getClient(i uint64) *http.Client {
	mutex.RLock()
	defer mutex.RUnlock()
	return clients[i]
}

func setClient(i uint64, request *model.Request) *http.Client {
	mutex.Lock()
	defer mutex.Unlock()
	client := createLangHTTPClient(request)
	clients[i] = client
	return client
}

// createLangHTTPClient 初始化长连接客户端参数
func createLangHTTPClient(request *model.Request) *http.Client {
	tr := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:        0,                // 最大连接数,默认0无穷大
		MaxIdleConnsPerHost: request.MaxCon,   // 对每个host的最大连接数量(MaxIdleConnsPerHost<=MaxIdleConns)
		IdleConnTimeout:     90 * time.Second, // 多长时间未使用自动关闭连接
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}
	if request.HTTP2 {
		// 使用真实证书 验证证书 模拟真实请求
		tr = &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:        0,                // 最大连接数,默认0无穷大
			MaxIdleConnsPerHost: request.MaxCon,   // 对每个host的最大连接数量(MaxIdleConnsPerHost<=MaxIdleConns)
			IdleConnTimeout:     90 * time.Second, // 多长时间未使用自动关闭连接
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: false},
		}
		_ = http2.ConfigureTransport(tr)
	}
	return &http.Client{
		Transport: tr,
	}
}
