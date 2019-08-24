# go-stress-testing
go 实现的压测工具 未完，待续~



## 目录
- [1、项目说明](#1项目说明)
    - [1.1 go-stress-testing](#11-go-stress-testing)
    - [1.2 项目体验](#12-项目体验)
- [2、介绍go-stress-testing](#2介绍go-stress-testing)
- [3、压测](#3压测)
    - [3.1 压测是什么](#31-压测是什么)
    - [3.2 为什么要压测](#32-为什么要压测)
    - [3.3 压测名词解释](#33-压测名词解释)
        - [3.3.1 压测类型解释](#331-压测类型解释)
        - [3.3.2 压测名词解释](#332-压测名词解释)
        - [3.3.3 机器性能指标解释](#333-机器性能指标解释)
        - [3.3.4 访问指标](#334-访问指标)
    - [3.4 如何计算压测指标](#34-如何计算压测指标)
    
- [4、压测工具](#4压测工具)
    - [4.1 ab](#41-ab)
    - [4.2 locust](#42-locust)
    - [4.3 Jmeter](#43-Jmeter)
    - [4.4 云压测](#44-云压测)
        - [4.4.1 云压测介绍](#441-云压测介绍)
        - [4.4.2 阿里云 性能测试 PTS](#442-阿里云-性能测试-PTS)
        - [4.4.3 腾讯云 压测大师 LM](#443-腾讯云-压测大师-LM)
    - [4.5 比较](#45-比较)
    - [4.6 如何选择](#46-如何选择)


## 1、项目说明
### 1.1 go-stress-testing

go 实现的压测工具，每个用户用一个协程的方式模拟，最大限度的利用CPU资源

### 1.2 项目体验

- 在 mac/linux/windows 不同平台下执行的命令

参数说明:

`-c` 表示并发数

`-n` 每个并发执行请求的次数，总请求的次数 = 并发数 * 每个并发执行请求的次数

`-u` 需要压测的地址

```shell script

# clone 项目
git clone git@github.com:link1st/go-stress-testing.git 

# 进入项目目录
cd go-stress-testing

# mac 
./go_stress_testing_mac -c 1 -n 100 -u https://www.baidu.com/

# linux
./go_stress_testing_linux  -c 1 -n 100 -u https://www.baidu.com/

# windows
go_stress_testing_win.exe -c="1" -n="100" -u="https://www.baidu.com/"

```

- 压测结果展示

执行以后，终端每秒钟都会输出一次结果，执行完成以后输出执行的统计

压测结果展示:

```

─────┬───────┬───────┬───────┬────────┬────────┬────────┬────────┬────────
 耗时│ 并发数 │ 成功数│ 失败数 │   qps  │最长耗时 │最短耗时│平均耗时 │ 错误码
─────┼───────┼───────┼───────┼────────┼────────┼────────┼────────┼────────
   1s│      1│      8│      0│    8.09│  133.16│  110.98│  123.56│200:8
   2s│      1│     15│      0│    8.02│  138.74│  110.98│  124.61│200:15
   3s│      1│     23│      0│    7.80│  220.43│  110.98│  128.18│200:23
   4s│      1│     31│      0│    7.83│  220.43│  110.23│  127.67│200:31
   5s│      1│     39│      0│    7.81│  220.43│  110.23│  128.03│200:39
   6s│      1│     46│      0│    7.72│  220.43│  110.23│  129.59│200:46
   7s│      1│     54│      0│    7.79│  220.43│  110.23│  128.42│200:54
   8s│      1│     62│      0│    7.81│  220.43│  110.23│  128.09│200:62
   9s│      1│     70│      0│    7.79│  220.43│  110.23│  128.33│200:70
  10s│      1│     78│      0│    7.82│  220.43│  106.47│  127.85│200:78
  11s│      1│     84│      0│    7.64│  371.02│  106.47│  130.96│200:84
  12s│      1│     91│      0│    7.63│  371.02│  106.47│  131.02│200:91
  13s│      1│     99│      0│    7.66│  371.02│  106.47│  130.54│200:99
  13s│      1│    100│      0│    7.66│  371.02│  106.47│  130.52│200:100


*************************  结果 stat  ****************************
处理协程数量: 1
请求总数: 100 总请求时间: 13.055 秒 successNum: 100 failureNum: 0
*************************  结果 end   ****************************

```

**耗时**: 程序运行耗时，每秒钟输出一次压测结果，最后一条为执行完成以后统计

**并发数**: 并发数

**成功数**: 压测中，请求成功的数量

**失败数**: 压测中，请求失败的数量

**qps**: 当前压测的QPS(每秒钟处理请求数量)

**最长耗时**: 压测中，单个请求最长的响应时长

**最短耗时**: 压测中，单个请求最短的响应时长

**平均耗时**: 压测中，单个请求平均的响应时长

**错误码**: 压测中，接口返回的 code码:返回次数 的集合


## 2、介绍go-stress-testing

### 2.1 介绍

- go-stress-testing 是go语言实现的简单压测工具，源码开源、支持二次开发，可以压测http、webSocket请求，使用协程模拟单个用户，可以更高效的利用CPU资源。

- 项目地址 [https://github.com/link1st/go-stress-testing](https://github.com/link1st/go-stress-testing)

### 2.2 用法

- 支持参数:

```
Usage of ./go_stress_testing_mac:
  -c uint
        并发数 (default 1)
  -d string
        调试模式 (default "false")
  -n uint
        请求总数 (default 1)
  -p string
        curl文件路径
  -u string
        请求地址
  -v string
        验证方法 http 支持:statusCode、json webSocket支持:json (default "statusCode")
```

- 使用示例:

```
# 查看用法
./go_stress_testing_mac

# 使用debug模式请求百度页面
./go_stress_testing_mac -c 1 -n 1 -d true -u https://www.baidu.com/

# 使用 curl文件(文件在curl目录下) 的方式请求
./go_stress_testing_mac -c 1 -n 1 -p curl/baidu.curl.txt

# 使用json的方式验证返回信息
./go_stress_testing_mac -c 1 -n 1 -u https://www.baidu.com/ -v json

# 压测webSocket连接
./go_stress_testing_mac -c 1 -n 1 -u ws://127.0.0.1:8089/acc
```


- 使用 curl文件进行压测

curl是Linux在命令行山下的工作的文件传输工具，是一款很强大的http命令行工具。

使用curl文件可以压测使用非GET的请求，支持设置http请求的 method、cookies、header、body等参数


chrome 浏览器生成 curl文件，打开开发者模式(快捷键F12)，如图所示，生成 curl 在终端执行命令
![copy cURL](https://img.mukewang.com/5d60eddd0001f4b016661114.png)

生成内容粘贴到项目目录下的**curl/baidu.curl.txt**文件中，执行下面命令就可以从curl.txt文件中读取需要压测的内容进行压测了

```
# 使用 curl文件(文件在curl目录下) 的方式请求
./go_stress_testing_mac -c 1 -n 1 -p curl/baidu.curl.txt
```

## 3、压测
### 3.1 压测是什么

压测，即压力测试，是确立系统稳定性的一种测试方法，通常在系统正常运作范围之外进行，以考察其功能极限和隐患。一般针对网络游戏压力测试从传统的意义来讲是对网络游戏的服务器不断施加“压力”的测试，是通过确定一个系统的瓶颈或者不能接收的性能点，来获得系统能提供的最大服务级别的测试。 一款网络游戏在上市前，游戏研发团队或运营商是会对其进行游戏压力测试的， 目的是为了了解游戏服务器的承受能力。以便更好的有目的的进行运营或研发。 主要检测游戏服务器的承受能力，包括用户承受能力（多少用户同时玩基本不影响质量）、流量承受等。

### 3.2 为什么要压测

- 压测的目的就是通过压测(模拟真实用户的行为)，测算出机器的性能(单台机器的QPS)，从而推算出系统在承受指定用户数(100W)时，需要多少机器能支撑得住
- 压测是在上线前为了应对未来可能达到的用户数量的一次预估，通过优化程序的性能或准备充足的机器，来保证用户的体验。


### 3.3 压测名词解释

#### 3.3.1 压测类型解释

| 压测类型 |   解释  |
| :----   | :---- |
| 压力测试(Stress Testing)          |  也称之为强度测试，测试一个系统的最大抗压能力，在强负载(大数据、高并发)的情况下，测试系统所能承受的最大压力，预估系统的瓶颈    |
| 并发测试(Concurrency Testing)     |  通过模拟很多用户同一时刻访问系统或对系统某一个功能进行操作，来测试系统的性能，从中发现问题(并发读写、线程控制、资源争抢)      |
| 耐久性测试(Configuration Testing) |  通过对系统在大负荷的条件下长时间运行，测试系统、机器的长时间运行下的状况,从中发现问题(内存泄漏、数据库连接池不释放、资源不回收)     |


#### 3.3.2 压测名词解释

| 压测名词 |   解释  |
| :----   | :---- |
| 并发(Concurrency)     |  指一个处理器同时处理多个任务的能力(逻辑上处理的能力)     |
| 并行(Parallel)        |  多个处理器或者是多核的处理器同时处理多个不同的任务(物理上同时执行)     |
| QPS(每秒钟查询数量 Query Per Second) | 服务器每秒钟处理请求数量 (req/sec  请求数/秒  一段时间内总请求数/请求时间)    |
| 事务(Transactions) | 是用户一次或者是几次请求的集合    |
| TPS(每秒钟处理事务数量 Transaction Per Second) | 服务器每秒钟处理事务数量(一个事务可能包括多个请求)    |
| 请求成功数(Request Success Number) | 在一次压测中，请求成功的数量    |
| 请求失败数(Request Failures Number) | 在一次压测中，请求失败的数量    |
| 错误率(Error Rate) | 在压测中，请求成功的数量与请求失败数量的比率  |
| 最大响应时间(Max Response Time) | 在一次事务中，从发出请求或指令系统做出的反映(响应)的最大时间  |
| 最少响应时间(Mininum Response Time) | 在一次事务中，从发出请求或指令系统做出的反映(响应)的最少时间  |
| 平均响应时间(Average Response Time) | 在一次事务中，从发出请求或指令系统做出的反映(响应)的平均时间  |

#### 3.3.3 机器性能指标解释

| 机器性能 |   解释  |
| :----   | :---- |
| CUP利用率(CPU Usage)       |  CUP 利用率分用户态、系统态和空闲态，CPU利用率是指:CPU执行非系统空闲进程的时间与CPU总执行时间的比率      |
| 内存使用率(Memory usage)    |  内存使用率指的是此进程所开销的内存。      |
| IO(Disk input/ output)    |  磁盘的读写包速率       |
| 网卡负载(Network Load)      |  网卡的进出带宽,包量       |

#### 3.3.4 访问指标

| 访问 |   解释  |
| :----   | :---- |
| PV(页面浏览量 Page View)           |  用户每打开1个网站页面，记录1个PV。用户多次打开同一页面，PV值累计多次      |
| UV(网站独立访客 Unique Visitor)    |  通过互联网访问、流量网站的自然人。1天内相同访客多次访问网站，只计算为1个独立访客       |

### 3.4 如何计算压测指标

- 压测我们需要有目的性的压测，这次压测我们需要达到什么目标(如:单台机器的性能为100QPS?网站能同时满足100W人同时在线)
- 可以通过一下计算方法来进行计算:
- 压测原则:每天80%的访问量集中在20%的时间里，这20%的时间就叫做峰值
- 公式: ( 总PV数*80% ) / ( 每天的秒数*20% ) = 峰值时间每秒钟请求数(QPS)
- 机器: 峰值时间每秒钟请求数(QPS) / 单台机器的QPS = 需要的机器的数量

- 假设:网站每天的用户数(100W)，每天的用户的访问量约为3000W PV，这台机器的需要多少QPS?
> ( 30000000*0.8 ) / (86400 * 0.2) ≈ 1389 (QPS)

- 假设:单台机器的的QPS是69，需要需要多少台机器来支撑？
> 1389 / 69 ≈ 20 


## 4、常见的压测工具
### 4.1 ab
- 简介

ApacheBench 是 Apache服务器自带的一个web压力测试工具，简称ab。ab又是一个命令行工具，对发起负载的本机要求很低，根据ab命令可以创建很多的并发访问线程，模拟多个访问者同时对某一URL地址进行访问，因此可以用来测试目标服务器的负载压力。总的来说ab工具小巧简单，上手学习较快，可以提供需要的基本性能指标，但是没有图形化结果，不能监控。

ab属于一个轻量级的压测工具，结果不会特别准确，可以用作参考。

- 安装

```shell script
# 在linux环境安装
sudo yum -y install httpd
```

- 用法

```
Usage: ab [options] [http[s]://]hostname[:port]/path
用法：ab [选项] 地址

选项：
Options are:
    -n requests      #执行的请求数，即一共发起多少请求。
    -c concurrency   #请求并发数。
    -s timeout       #指定每个请求的超时时间，默认是30秒。
    -k               #启用HTTP KeepAlive功能，即在一个HTTP会话中执行多个请求。默认时，不启用KeepAlive功能。
```

- 压测命令

```shell script
# 使用ab压测工具，对百度的链接 请求100次，并发数1
ab -n 100 -c 1 https://www.baidu.com/
```

压测结果

```
~ >ab -n 100 -c 1 https://www.baidu.com/
This is ApacheBench, Version 2.3 <$Revision: 1430300 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking www.baidu.com (be patient).....done


Server Software:        BWS/1.1
Server Hostname:        www.baidu.com
Server Port:            443
SSL/TLS Protocol:       TLSv1.2,ECDHE-RSA-AES128-GCM-SHA256,2048,128

Document Path:          /
Document Length:        227 bytes

Concurrency Level:      1
Time taken for tests:   9.430 seconds
Complete requests:      100
Failed requests:        0
Write errors:           0
Total transferred:      89300 bytes
HTML transferred:       22700 bytes
Requests per second:    10.60 [#/sec] (mean)
Time per request:       94.301 [ms] (mean)
Time per request:       94.301 [ms] (mean, across all concurrent requests)
Transfer rate:          9.25 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:       54   70  16.5     69     180
Processing:    18   24  12.0     23     140
Waiting:       18   24  12.0     23     139
Total:         72   94  20.5     93     203

Percentage of the requests served within a certain time (ms)
  50%     93
  66%     99
  75%    101
  80%    102
  90%    108
  95%    122
  98%    196
  99%    203
 100%    203 (longest request)
```

- 主要关注的测试指标

- `Concurrency Level` 并发请求数

- `Time taken for tests` 整个测试时间

- `Complete requests` 完成请求个数

- `Failed requests` 失败个数

- `Requests per second` 吞吐量，指的是某个并发用户数下单位时间内处理的请求数。等效于QPS，其实可以看作同一个统计方式，只是叫法不同而已。

- `Time per request` 用户平均请求等待时间

- `Time per request` 服务器处理时间

### 4.2 Locust

- 简介

是非常简单易用、分布式、python开发的压力测试工具。有图形化界面，支持将压测数据导出。

- 安装

```shell script
# pip3 安装locust
pip3  install locust
# 查看是否安装成功
locust -h
# 运行 Locust 分布在多个进程/机器库
pip3 install pyzmq
# webSocket 压测库
pip3 install websocket-client
```

- 用法

编写压测脚本 **test.py**

```python
from locust import HttpLocust, TaskSet, task

# 定义用户行为
class UserBehavior(TaskSet):

    @task
    def baidu_index(self):
        self.client.get("/")


class WebsiteUser(HttpLocust):
    task_set = UserBehavior # 指向一个定义的用户行为类
    min_wait = 3000 # 执行事务之间用户等待时间的下界（单位：毫秒）
    max_wait = 6000 # 执行事务之间用户等待时间的上界（单位：毫秒）
```

- 启动压测

```shell script
locust -f  test.py --host=https://www.baidu.com
```

访问 http://localhost:8089 进入压测首页

Number of users to simulate 模拟用户数

Hatch rate (users spawned/second) 每秒钟增加用户数

点击 "Start swarming" 进入压测页面


![locust 首页](https://img.mukewang.com/5d5e4f81000179cd25541372.png)


压测界面右上角有:被压测的地址、当前状态、RPS、失败率、开始或重启按钮

性能测试参数

- `Type` 请求的类型，例如GET/POST

- `Name` 请求的路径

- `Request` 当前请求的数量

- `Fails` 当前请求失败的数量

- `Median` 中间值，单位毫秒，请求响应时间的中间值

- `Average` 平均值，单位毫秒，请求的平均响应时间

- `Min` 请求的最小服务器响应时间，单位毫秒

- `Max` 请求的最大服务器响应时间，单位毫秒

- `Average size` 单个请求的大小，单位字节

- `Current RPS` 代表吞吐量(Requests Per Second的缩写)，指的是某个并发用户数下单位时间内处理的请求数。等效于QPS，其实可以看作同一个统计方式，只是叫法不同而已。

![locust 压测页面](https://img.mukewang.com/5d5e4fad000177e125501368.png)

### 4.3 Jmeter

- 简介

Apache JMeter是Apache组织开发的基于Java的压力测试工具。用于对软件做压力测试，它最初被设计用于Web应用测试，但后来扩展到其他测试领域。
JMeter能够对应用程序做功能/回归测试，通过创建带有断言的脚本来验证你的程序返回了你期望的结果。

- 安装

访问 https://jmeter-plugins.org/install/Install/ 下载解压以后即可使用

- 用法

JMeter的功能过于强大，这里暂时不介绍用法，可以查询相关文档使用(参考文献中有推荐的教程文档)


### 4.4 云压测

#### 4.4.1 云压测介绍

顾名思义就是将压测脚本部署在云端，通过云端对对我们的应用进行全方位压测，只需要配置压测的参数，无需准备实体机，云端自动给我们分配需要压测的云主机，对被压测目标进行压测。

云压测的优势:

1. 轻易的实现分布式部署
2. 能够模拟海量用户的访问
3. 流量可以从全国各地发起，更加真实的反映用户的体验
4. 全方位的监控压测指标


#### 4.4.2 阿里云 性能测试 PTS

PTS（Performance Testing Service）是面向所有技术背景人员的云化测试工具。有别于传统工具的繁复，PTS以互联网化的交互，提供性能测试、API调试和监测等多种能力。自研和适配开源的功能都可以轻松模拟任意体量的用户访问业务的场景，任务随时发起，免去繁琐的搭建和维护成本。更是紧密结合监控、流控等兄弟产品提供一站式高可用能力，高效检验和管理业务性能。

阿里云同样还是支持渗透测试，通过模拟黑客对业务系统进行全面深入的安全测试。


#### 4.4.3 腾讯云 压测大师 LM

通过创建虚拟机器人模拟多用户的并发场景，提供一整套完整的服务器压测解决方案



### 4.5 比较


### 与ab对比 对go http接口压测

申请的服务器情况

- go server 文件

```golang
package main

import (
    "log"
    "net/http"
)

const (
    httpPort = "8088"
)

func main() {
    hello := func(w http.ResponseWriter, req *http.Request) {
        data := "Hello, World!"

        w.Header().Add("Server", "golang")
        w.Write([]byte(data))

        return
    }

    http.HandleFunc("/", hello)
    err := http.ListenAndServe(":"+httpPort, nil)

    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

```

- ab 压测命令

```
ab -c 100 -n 100000  http://127.0.0.1:8088/
```

- go_stress_testing 

```
./go_stress_testing_linux -c 100 -n 10000 -u http://127.0.0.1:8088/
```


### 压测注意事项

- ab `-n` 是总压测次数
- go-stress-testing `-n` 是单个用户请求的次数，总次数= `-c`* `-n`， 这里考虑的是模拟用户行为，所以这个是每个用户请求的次数


### 问题 Q/A


 未完，待续~~~


过程
先写代码 优化 书写文档 优化文档 绘制图标 发布






## 常见的压测工具
- 常用的压测工具实现的语言、使用方法、比较、说明
- Jmeter
- AB

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
- 压测结束条件

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
- webSocket 连接
- 连接以后初始化动作
- 循环事件
- 显示压测数量
- 实现用户数 每秒加几个用户
- 程序结束 优化
- webSocket 建立连接以后，发送消息
- webSocket 首发数据模型调整
- 优化，把测试的go 文件抽离 和 locust类似

## 完善的


## 反思

[性能测试工具](https://testerhome.com/topics/17068)

[性能测试常见名词解释](https://blog.csdn.net/r455678/article/details/53063989)

[性能测试名词解释](https://codeigniter.org.cn/forums/blog-39678-2456.html)

[PV、TPS、QPS是怎么计算出来的？](https://www.zhihu.com/question/21556347)

[超实用压力测试工具－ab工具](https://www.jianshu.com/p/43d04d8baaf7)

[Locust 介绍](http://www.testclass.net/locust/introduce)

[Jmeter性能测试 入门](https://www.cnblogs.com/TankXiao/p/4045439.html)

[阿里云 性能测试 PTS](https://cn.aliyun.com/product/pts)

[腾讯云 压测大师 LM](https://cloud.tencent.com/product/lm/details)


github 搜:link1st 查看项目

[https://github.com/link1st/go-stress-testing](https://github.com/link1st/go-stress-testing)


