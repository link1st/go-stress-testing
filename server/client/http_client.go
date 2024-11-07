// Package client http 客户端
package client

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/http2"

	"github.com/link1st/go-stress-testing/helper"
	"github.com/link1st/go-stress-testing/model"
	httplongclinet "github.com/link1st/go-stress-testing/server/client/http_longclinet"
)

// logErr err
var logErr = log.New(os.Stderr, "", 0)

// HTTPRequest HTTP 请求
// method 方法 GET POST
// url 请求的url
// body 请求的body
// headers 请求头信息
// timeout 请求超时时间
func HTTPRequest(chanID uint64, request *model.Request) (resp *http.Response, requestTime uint64, err error) {
	method := request.Method
	url := request.URL
	body := request.GetBody()
	timeout := request.Timeout
	headers := request.CopyHeaders()

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		logErr.Println("Error creating request:", err)
		return
	}

	// Set Host if provided in headers
	if host, ok := headers["Host"]; ok {
		req.Host = host
	}

	// Ensure Content-Type is set
	if _, ok := headers["Content-Type"]; !ok {
		headers["Content-Type"] = "application/json; charset=utf-8"
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	var client *http.Client
	if request.Keepalive {
		client = httplongclinet.NewClient(chanID, request)
	} else {
		req.Close = true
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		if request.HTTP2 {
			tr = &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
			}
			if err = http2.ConfigureTransport(tr); err != nil {
				logErr.Println("Error configuring HTTP2 transport:", err)
				return
			}
		}
		client = &http.Client{
			Transport: tr,
			Timeout:   timeout,
		}
		if !request.Redirect {
			client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}
		}
	}

	startTime := time.Now()
	resp, err = client.Do(req)
	requestTime = uint64(helper.DiffNano(startTime))
	if err != nil {
		logErr.Println("Request failed:", err)
		return
	}

	return
}
