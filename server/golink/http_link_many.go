/**
* Created by GoLand.
* User: link1st
* Date: 2020/7/31
* Time: 8:36 下午
 */

package golink

import (
	"time"

	"go-stress-testing/model"
)

// 接口分步压测
type ReqListMany struct {
	list []*model.Request
}

func (r *ReqListMany) getCount() int {
	return len(r.list)
}

var (
	clientList *ReqListMany
)

// 接口分步压测示例
func init() {

	clientList = &ReqListMany{}

	// TODO::接口分步压测示例
	// 需要压测的接口参数
	clients := make([]*model.Request, 0)

	// 压测第一步
	clients = append(clients, &model.Request{
		Url:    "https://page.aliyun.com/delivery/plan/list", // 请求url
		Form:   "http",                                       // 请求方式 示例参数:http/webSocket/tcp
		Method: "POST",                                       // 请求方法 示例参数:GET/POST/PUT
		Headers: map[string]string{
			"referer": "https://cn.aliyun.com/",
			"cookie":  "aliyun_choice=CN; JSESSIONID=J8866281-CKCFJ4BUZ7GDO9V89YBW1-KJ3J5V9K-GYUW7; maliyun_temporary_console0=1AbLByOMHeZe3G41KYd5WWZvrM%2BGErkaLcWfBbgveKA9ifboArprPASvFUUfhwHtt44qsDwVqMk8Wkdr1F5LccYk2mPCZJiXb0q%2Bllj5u3SQGQurtyPqnG489y%2FkoA%2FEvOwsXJTvXTFQPK%2BGJD4FJg%3D%3D; cna=L3Q5F8cHDGgCAXL3r8fEZtdU; isg=BFNThsmSCcgX-sUcc5Jo2s2T4tF9COfKYi8g9wVwr3KphHMmjdh3GrHFvPTqJD_C; l=eBaceXLnQGBjstRJBOfwPurza77OSIRAguPzaNbMiT5POw1B5WAlWZbqyNY6C3GVh6lwR37EODnaBeYBc3K-nxvOu9eFfGMmn",
		},                                                                                                                                                                            // headers 头信息
		Body:    "adPlanQueryParam=%7B%22adZone%22%3A%7B%22positionList%22%3A%5B%7B%22positionId%22%3A83%7D%5D%7D%2C%22requestId%22%3A%2217958651-f205-44c7-ad5d-f8af92a6217a%22%7D", // 消息体
		Verify:  "statusCode",                                                                                                                                                        // 验证的方法 示例参数:statusCode、json
		Timeout: 30 * time.Second,                                                                                                                                                    // 是否开启Debug模式
		Debug:   false,                                                                                                                                                               // 是否开启Debug模式
	})

	// 压测第二步
	clients = append(clients, &model.Request{
		Url:    "https://page.aliyun.com/delivery/plan/list", // 请求url
		Form:   "http",                                       // 请求方式 示例参数:http/webSocket/tcp
		Method: "POST",                                       // 请求方法 示例参数:GET/POST/PUT
		Headers: map[string]string{
			"referer": "https://cn.aliyun.com/",
			"cookie":  "aliyun_choice=CN; JSESSIONID=J8866281-CKCFJ4BUZ7GDO9V89YBW1-KJ3J5V9K-GYUW7; maliyun_temporary_console0=1AbLByOMHeZe3G41KYd5WWZvrM%2BGErkaLcWfBbgveKA9ifboArprPASvFUUfhwHtt44qsDwVqMk8Wkdr1F5LccYk2mPCZJiXb0q%2Bllj5u3SQGQurtyPqnG489y%2FkoA%2FEvOwsXJTvXTFQPK%2BGJD4FJg%3D%3D; cna=L3Q5F8cHDGgCAXL3r8fEZtdU; isg=BFNThsmSCcgX-sUcc5Jo2s2T4tF9COfKYi8g9wVwr3KphHMmjdh3GrHFvPTqJD_C; l=eBaceXLnQGBjstRJBOfwPurza77OSIRAguPzaNbMiT5POw1B5WAlWZbqyNY6C3GVh6lwR37EODnaBeYBc3K-nxvOu9eFfGMmn",
		},                                                                                                                                                                            // headers 头信息
		Body:    "adPlanQueryParam=%7B%22adZone%22%3A%7B%22positionList%22%3A%5B%7B%22positionId%22%3A83%7D%5D%7D%2C%22requestId%22%3A%2217958651-f205-44c7-ad5d-f8af92a6217a%22%7D", // 消息体
		Verify:  "statusCode",                                                                                                                                                        // 验证的方法 示例参数:statusCode、json
		Timeout: 30 * time.Second,                                                                                                                                                    // 是否开启Debug模式
		Debug:   false,                                                                                                                                                               // 是否开启Debug模式
	})

	clientList.list = clients

	// TODO::注释下面一行代码
	clientList.list = nil
}

func getRequestList(request *model.Request) []*model.Request {

	if len(clientList.list) <= 0 {

		return []*model.Request{request}
	}

	return clientList.list
}
