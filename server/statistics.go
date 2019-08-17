/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-15
* Time: 18:14
 */

package server

import (
	"fmt"
	"go-stress-testing/model"
	"time"
)

// 接收结果
// 统计的时间都是纳秒，显示的时间 都是毫秒
// concurrent 并发数
func ReceivingResults(concurrent uint64, ch <-chan *model.RequestResults) {

	// 时间
	var (
		processingTime uint64 // 处理总时间
		requestTime    uint64 // 请求总时间
		maxTime        uint64 // 最大时长
		minTime        uint64 // 最小时长
		successNum     uint64 // 成功处理数，code为0
		failureNum     uint64 // 处理失败数，code不为0
	)

	statTime := uint64(time.Now().UnixNano())

	// 错误码/错误个数
	var errCode = make(map[int]int)

	// 每个秒时间输出一次 计算结果
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				go calculateData(concurrent, processingTime, requestTime, maxTime, minTime, successNum, failureNum, errCode)
			}
		}
	}()

	for data := range ch {
		// fmt.Println("处理一条数据", data.Id, data.Time, data.IsSucceed, data.ErrCode)
		processingTime = processingTime + data.Time

		if maxTime <= data.Time {
			maxTime = data.Time
		}

		if minTime == 0 {
			minTime = data.Time
		} else if minTime > data.Time {
			minTime = data.Time
		}

		// 是否请求成功
		if data.IsSucceed == true {
			successNum = successNum + 1
		} else {
			failureNum = failureNum + 1
		}

		// 统计错误码
		if value, ok := errCode[data.ErrCode]; ok {
			errCode[data.ErrCode] = value + 1
		} else {
			errCode[data.ErrCode] = 1
		}
	}

	endTime := uint64(time.Now().UnixNano())
	requestTime = endTime - statTime

	fmt.Println("*************************  结果 stat  ****************************")
	fmt.Println("处理协程数量:", concurrent, "程序处理总时长:", fmt.Sprintf("%.3f", float64(processingTime/concurrent)/1e9), "秒")
	fmt.Println("请求总数:", successNum+failureNum, "总请求时间:", fmt.Sprintf("%.3f", float64(requestTime)/1e9),
		"秒", "successNum:", successNum, "failureNum:", failureNum)

	calculateData(concurrent, processingTime, requestTime, maxTime, minTime, successNum, failureNum, errCode)

	fmt.Println("*************************  结果 end  ****************************")
}

// 计算数据
func calculateData(concurrent, processingTime, requestTime, maxTime, minTime, successNum, failureNum uint64, errCode map[int]int) {
	if processingTime == 0 {
		processingTime = 1
	}

	// 平均 每个协程成功数*总协程数据/总耗时 (每秒)
	qps := float64(successNum*1e9*concurrent) / float64(processingTime)

	// 平均时长 总耗时/总请求数 毫秒
	averageTime := float64(processingTime) / float64(successNum*1e6*concurrent)

	maxTimeFloat := float64(maxTime) / 1e6
	minTimeFloat := float64(minTime) / 1e6

	// 打印的时长都为毫秒
	result := fmt.Sprintf("请求总数:%8d|successNum:%8d|failureNum:%8d|qps:%9.3f|maxTime:%9.3f|minTime:%9.3f|平均时长:%9.3f|errCode:%v", successNum+failureNum, successNum, failureNum, qps, maxTimeFloat, minTimeFloat, averageTime, errCode)
	fmt.Println(result)
}

// 打印表头信息
func header() {
	// 打印的时长都为毫秒
	result := fmt.Sprintf("请求时间\t|\t总请求数\t|\t成功数\t|失败数\t|\tQPS\t|\t最长耗时\t|\t最短耗时\t|\t平均耗时\t|\t错误码\t")
	fmt.Println(result)

	return
}

func table(successNum, failureNum uint64, errCode map[uint32]int, qps, averageTime, maxTimeFloat, minTimeFloat float64) {
	// 打印的时长都为毫秒
	result := fmt.Sprintf("请求总数:%8d|successNum:%8d|failureNum:%8d|qps:%9.3f|maxTime:%9.3f|minTime:%9.3f|平均时长:%9.3f|errCode:%v", successNum+failureNum, successNum, failureNum, qps, maxTimeFloat, minTimeFloat, averageTime, errCode)
	fmt.Println(result)

	return
}
