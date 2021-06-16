package httplongclinet

import (
	"crypto/tls"
	"go-stress-testing/model"
	"golang.org/x/net/http2"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	LangHttpClient *http.Client
	once           sync.Once
)

//初始化长连接客户端参数
func CreateLangHttpClient(request *model.Request) {
	once.Do(func() {
		tr := &http.Transport{}
		if request.HTTP2 {
			//使用真实证书 验证证书 模拟真实请求
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
			http2.ConfigureTransport(tr)
		} else {
			// 跳过证书验证
			tr = &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				MaxIdleConns:        0,                // 最大连接数,默认0无穷大
				MaxIdleConnsPerHost: request.MaxCon,   // 对每个host的最大连接数量(MaxIdleConnsPerHost<=MaxIdleConns)
				IdleConnTimeout:     90 * time.Second, // 多长时间未使用自动关闭连接
				TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			}
		}

		LangHttpClient = &http.Client{
			Transport: tr,
		}
	})
}
