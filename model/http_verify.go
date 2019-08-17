/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-16
* Time: 16:03
 */

package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// 通过Http状态码判断是否请求成功
func HttpStatusCode(request *Request, response *http.Response) (code int, isSucceed bool) {

	defer response.Body.Close()
	code = response.StatusCode
	if code == http.StatusOK {
		isSucceed = true
	}

	// 开启调试模式
	if request.GetDebug() {
		body, err := ioutil.ReadAll(response.Body)
		fmt.Printf("请求结果 httpCode:%d body:%s err:%v \n", response.StatusCode, string(body), err)
	}

	return
}

/***************************  返回值为json  ********************************/

type ResponseJson struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// 通过返回的Body 判断
// 返回示例:: {"code":200,"msg":"Success","data":{}}
// code 默认将http code作为返回码，http code 为200时 取body中的返回code
func HttpJson(request *Request, response *http.Response) (code int, isSucceed bool) {

	defer response.Body.Close()
	code = response.StatusCode
	if code == http.StatusOK {

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			code = ParseError
			fmt.Printf("请求结果 ioutil.ReadAll err:%v", err)
		} else {
			responseJson := &ResponseJson{}
			err = json.Unmarshal(body, responseJson)
			if err != nil {
				code = ParseError
				fmt.Printf("请求结果 json.Unmarshal err:%v", err)
			} else {

				code = responseJson.Code

				// body 中code返回200为返回数据成功
				if responseJson.Code == 200 {
					isSucceed = true
				}
			}
		}

		// 开启调试模式
		if request.GetDebug() {
			fmt.Printf("请求结果 httpCode:%d body:%s err:%v \n", response.StatusCode, string(body), err)
		}
	}

	return
}
