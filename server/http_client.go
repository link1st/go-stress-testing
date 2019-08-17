/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-15
* Time: 21:03
 */

package server

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"
)

// http get
func HttpGetResp(url string) (resp *http.Response, err error) {
	resp, err = http.Get(url)
	if err != nil {
		fmt.Println("HttpGet err:", err)

		return
	}

	return
}

// 请求
func HttpRequest(method, url string, body io.Reader, headers map[string]string, timeout time.Duration) (resp *http.Response, err error) {

	// 跳过证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// http cookie接口
	cookieJar, _ := cookiejar.New(nil)

	client := &http.Client{
		Transport: tr,
		Jar:       cookieJar,
		Timeout:   timeout,
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {

		return
	}

	if _, ok := headers["Content-Type"]; !ok {
		if headers == nil {
			headers = make(map[string]string)
		}
		headers["Content-Type"] = "application/x-www-form-urlencoded; charset=utf-8"
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err = client.Do(req)
	if err != nil {

		return
	}

	return
}
