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
	)

	flag.Uint64Var(&concurrency, "c", 1, "并发数")
	flag.Uint64Var(&totalNumber, "n", 1, "请求总数")

	// 解析参数
	flag.Parse()

	if concurrency == 0 || totalNumber == 0 {
		flag.Usage()

		return
	}

	fmt.Printf("开始启动  并发数:%d 请求数:%d\n", concurrency, totalNumber)

	dispose(concurrency, totalNumber)

	return
}

func dispose(concurrency, totalNumber uint64) {

	// 设置接收数据缓存
	ch := make(chan *model.RequestResults, 1000)
	var (
		// TODO::容易丢数据 或不及时返回
		wg sync.WaitGroup
	)

	go server.ReceivingResults(concurrency, ch)

	for i := uint64(0); i < concurrency; i++ {
		wg.Add(1)
		go goLink(i, ch, totalNumber, &wg)
	}

	wg.Wait()

	close(ch)

	time.Sleep(1 * time.Second)


	fmt.Println("完成")

}

func goLink(id uint64, ch chan<- *model.RequestResults, totalNumber uint64, wg *sync.WaitGroup) {

	defer func() {
		wg.Done()
	}()

	fmt.Printf("启动协程 编号:%05d \n", id)

	for i := uint64(0); i < totalNumber; i++ {
		requestResults := &model.RequestResults{
			Time:      4,
			IsSucceed: true,
			ErrCode:   0,
		}

		ch <- requestResults

		time.Sleep(10 * time.Millisecond)
	}

	return
}
