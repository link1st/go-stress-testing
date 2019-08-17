# go-stress-testing
go 实现的压测工具


过程
先写代码 优化 书写文档 优化文档 绘制图标 发布

## 压测名词解释
- 并发
- 并行
- QPS
- TPS
- 耗时
- 响应时间
- 并发数
- 总请求量


## 常见的压测工具
- 常用的压测工具实现的语言、使用方法、比较、说明
- Jmeter

## 网络压测
- 阿里云等

## 为什么要实现一个压测工具
- 为什么

## go语言实现压测
- 怎么保持连接
- 实现过程
- 耗时怎么计算，算不算连接事件
- tcp
- webSocket
- http常见的压测

## 新机器部署进行压测
- 申请机器
- 部署环境
- 启动压测 
- 被压测的机器收集数据

## 压测以后
- 数据汇聚成图表
- 百度 ECharts

## 注意事项

## 工作
- 通过文件的方式设置 body Headers
- 解析curl的参数

```
curl 'https://api.h5game.g.mi.com/task/getTaskList' -H 'Accept: application/json, text/plain, */*' -H 'Referer: https://static.g.mi.com/game/platform/linkGame.html' -H 'User-Agent: Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1' -H 'Sec-Fetch-Mode: cors' -H 'Content-Type: application/x-www-form-urlencoded' --data 'fromApp=&serviceToken=_SW01_qyWxdVoBW2Qlv3G1yk9aObtnXVj4fnZ9cC9TceXkUaTnH3%2FCzuucnCi2M0gfaVy2&channel=1' --compressed

```

## 完善的


## 反思

[性能测试工具](https://testerhome.com/topics/17068)


