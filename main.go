/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-15
* Time: 13:44
 */

package main

import (
	"flag"
	"fmt"
	"go-stress-testing/model"
	"go-stress-testing/server"
	"sync"
	"time"
)

func main() {

	var (
		concurrency uint64
		totalNumber uint64
		debugStr    string
		requestUrl  string
		path        string
	)

	flag.Uint64Var(&concurrency, "c", 1, "并发数")
	flag.Uint64Var(&totalNumber, "n", 1, "请求总数")
	flag.StringVar(&debugStr, "d", "false", "调试模式")
	flag.StringVar(&requestUrl, "u", "", "请求地址")
	flag.StringVar(&path, "p", "", "curl文件")

	// 解析参数
	flag.Parse()
	if concurrency == 0 || totalNumber == 0 || (requestUrl == "" && path == "") {
		fmt.Printf("示例: go run main.go -c 1 -n 1 -u https://www.baidu.com/ \n")
		fmt.Printf("-c %d -n %d -d %v -u %s \n", concurrency, totalNumber, debugStr, requestUrl)

		flag.Usage()

		return
	}

	debug := debugStr == "true"
	request, err := model.NewRequest(requestUrl, "", 0, debug, path)
	if err != nil {
		fmt.Printf("参数不合法 %v \n", err)

		return
	}

	fmt.Printf("\n 开始启动  并发数:%d 请求数:%d 请求参数: \n", concurrency, totalNumber)
	request.Print()

	dispose(concurrency, totalNumber, request)

	return
}

func dispose(concurrency, totalNumber uint64, request *model.Request) {

	// 设置接收数据缓存
	ch := make(chan *model.RequestResults, 1000)
	var (
		// TODO::容易丢数据 或不及时返回
		wg sync.WaitGroup
	)

	go server.ReceivingResults(concurrency, ch)

	for i := uint64(0); i < concurrency; i++ {
		wg.Add(1)
		go goLink(i, ch, totalNumber, &wg, request)
	}

	wg.Wait()

	close(ch)

	time.Sleep(50 * time.Microsecond)

}

// 请求时间
// diff 纳秒
func forHowLong(startTime time.Time) (diff uint64) {
	startTimeStamp := startTime.UnixNano()
	endTimeStamp := time.Now().UnixNano()
	diff = uint64(endTimeStamp - startTimeStamp)

	return
}

// http go link
func goLink(chanId uint64, ch chan<- *model.RequestResults, totalNumber uint64, wg *sync.WaitGroup, request *model.Request) {

	defer func() {
		wg.Done()
	}()

	// fmt.Printf("启动协程 编号:%05d \n", chanId)

	for i := uint64(0); i < totalNumber; i++ {

		var (
			startTime = time.Now()
			isSucceed = false
			errCode   = model.HttpOk
		)

		resp, err := server.HttpRequest(request.Method, request.Url, request.Body, request.Headers, request.Timeout)
		// resp, err := server.HttpGetResp(request.Url)
		if err != nil {
			errCode = model.RequestErr // 请求错误
		} else {
			// 验证请求是否成功
			errCode, isSucceed = request.VerifyHttp(request, resp)
		}

		requestTime := forHowLong(startTime)

		requestResults := &model.RequestResults{
			Time:      requestTime,
			IsSucceed: isSucceed,
			ErrCode:   errCode,
		}

		requestResults.SetId(chanId, i)

		ch <- requestResults
	}

	return
}
