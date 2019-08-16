/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-15
* Time: 21:03
 */

package server

import (
	"fmt"
	"net/http"
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
