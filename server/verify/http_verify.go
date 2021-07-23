// Package verify 校验
package verify

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"go-stress-testing/model"
)

// getZipData 处理gzip压缩
func getZipData(response *http.Response) (body []byte, err error) {
	var reader io.ReadCloser
	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(response.Body)
		defer func() {
			_ = reader.Close()
		}()
	default:
		reader = response.Body
	}
	body, err = ioutil.ReadAll(reader)
	response.Body = ioutil.NopCloser(bytes.NewReader(body))
	return
}

// HTTPStatusCode 通过 HTTP 状态码判断是否请求成功
func HTTPStatusCode(request *model.Request, response *http.Response) (code int, isSucceed bool) {
	defer func() {
		_ = response.Body.Close()
	}()
	code = response.StatusCode
	if code == request.Code {
		isSucceed = true
	}
	// 开启调试模式
	if request.GetDebug() {
		body, err := getZipData(response)
		fmt.Printf("请求结果 httpCode:%d body:%s err:%v \n", response.StatusCode, string(body), err)
	}
	io.Copy(ioutil.Discard, response.Body)
	return
}

/***************************  返回值为json  ********************************/

// ResponseJSON 返回数据结构体
type ResponseJSON struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// HTTPJson  通过返回的Body 判断
// 返回示例: {"code":200,"msg":"Success","data":{}}
// code 默认将http code作为返回码，http code 为200时 取body中的返回code
func HTTPJson(request *model.Request, response *http.Response) (code int, isSucceed bool) {
	defer func() {
		_ = response.Body.Close()
	}()
	code = response.StatusCode
	if code == http.StatusOK {
		body, err := getZipData(response)
		if err != nil {
			code = model.ParseError
			fmt.Printf("请求结果 ioutil.ReadAll err:%v", err)
		} else {
			responseJSON := &ResponseJSON{}
			err = json.Unmarshal(body, responseJSON)
			if err != nil {
				code = model.ParseError
				fmt.Printf("请求结果 json.Unmarshal err:%v", err)
			} else {
				code = responseJSON.Code
				// body 中code返回200为返回数据成功
				if responseJSON.Code == request.Code {
					isSucceed = true
				}
			}
		}
		// 开启调试模式
		if request.GetDebug() {
			fmt.Printf("请求结果 httpCode:%d body:%s err:%v \n", response.StatusCode, string(body), err)
		}
	}
	io.Copy(ioutil.Discard, response.Body)
	return
}
