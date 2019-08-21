/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-16
* Time: 16:03
 */

package verify

import (
	"encoding/json"
	"fmt"
	"go-stress-testing/model"
)

/***************************  返回值为json  ********************************/

// 返回数据结构体
type WebSocketResponseJson struct {
	Seq      string `json:"seq"`
	Cmd      string `json:"cmd"`
	Response struct {
		Code    int         `json:"code"`
		CodeMsg string      `json:"codeMsg"`
		Data    interface{} `json:"data"`
	} `json:"response"`
}

// 通过返回的Body 判断
// 返回示例: {"seq":"1566276523281-585638","cmd":"heartbeat","response":{"code":200,"codeMsg":"Success","data":null}}
// code 取body中的返回code
func WebSocketJson(request *model.Request, seq string, msg []byte) (code int, isSucceed bool) {

	responseJson := &WebSocketResponseJson{}
	err := json.Unmarshal(msg, responseJson)
	if err != nil {
		code = model.ParseError
		fmt.Printf("请求结果 json.Unmarshal msg:%s err:%v", string(msg), err)
	} else {

		if seq != responseJson.Seq {
			code = model.ParseError
			fmt.Println("请求和返回seq不一致 ~请求:", seq, responseJson.Seq, string(msg))
		} else {
			code = responseJson.Response.Code
			// body 中code返回200为返回数据成功
			if code == 200 {
				isSucceed = true
			}
		}
	}

	// 开启调试模式
	if request.GetDebug() {
		fmt.Printf("请求结果 seq:%s body:%s \n", seq, string(msg))
	}

	return
}
