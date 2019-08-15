/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-15
* Time: 21:03
 */

package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// http get
func HttpGet(url string) (data string, err error) {

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("HttpGet err:", err)
		return
	}

	defer resp.Body.Close()

	// 读取数据
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("HttpGet err:", err)

		return data, err
	}

	data = string(body)
	return
}