# go实现的压测工具【单台机器100w连接压测实战】

本文介绍压测是什么，解释压测的专属名词，教大家如何压测。介绍市面上的常见压测工具(ab、locust、Jmeter、go实现的压测工具、云压测)，对比这些压测工具，教大家如何选择一款适合自己的压测工具，本文还有两个压测实战项目：

- 单台机器对 HTTP 短连接 QPS 1W+ 的压测实战
- 单台机器 100W 长连接的压测实战
- 对 grpc 接口进行压测
> 简单扩展即可支持 私有协议

## 目录
- [1、项目说明](#1项目说明)
    - [1.1 go-stress-testing](#11-go-stress-testing)
    - [1.2 项目体验](#12-项目体验)
- [2、压测](#2压测)
    - [2.1 压测是什么](#21-压测是什么)
    - [2.2 为什么要压测](#22-为什么要压测)
    - [2.3 压测名词解释](#23-压测名词解释)
        - [2.3.1 压测类型解释](#231-压测类型解释)
        - [2.3.2 压测名词解释](#232-压测名词解释)
        - [2.3.3 机器性能指标解释](#233-机器性能指标解释)
        - [2.3.4 访问指标解释](#234-访问指标解释)
    - [3.4 如何计算压测指标](#24-如何计算压测指标)
- [3、常见的压测工具](#3常见的压测工具)
    - [3.1 ab](#31-ab)
    - [3.2 locust](#32-locust)
    - [3.3 JMeter](#33-JMeter)
    - [3.4 云压测](#34-云压测)
        - [3.4.1 云压测介绍](#341-云压测介绍)
        - [3.4.2 阿里云 性能测试 PTS](#342-阿里云-性能测试-PTS)
        - [3.4.3 腾讯云 压测大师 LM](#343-腾讯云-压测大师-LM)
- [4、go-stress-testing go语言实现的压测工具](#4go-stress-testing-go语言实现的压测工具)
    - [4.1 介绍](#41-介绍)
    - [4.2 用法](#42-用法)
    - [4.3 实现](#43-实现)
    - [4.4 go-stress-testing 对 Golang web 压测](#44-go-stress-testing-对-golang-web-压测)
    - [4.5 grpc压测](#45-grpc压测)
- [5、压测工具的比较](#5压测工具的比较)
    - [5.1 比较](#51-比较)
    - [5.2 如何选择压测工具](#52-如何选择压测工具)
- [6、单台机器100w连接压测实战](#6单台机器100w连接压测实战)
    - [6.1 说明](#61-说明)
    - [6.2 内核优化](#62-内核优化)
    - [6.3 客户端配置](#63-客户端配置)
    - [6.4 准备](#64-准备)
    - [6.5 压测数据](#65-压测数据)
- [7、常见问题](#7常见问题)
- [8、总结](#8总结)
- [9、参考文献](#9参考文献)


## 1、项目说明
### 1.1 go-stress-testing

go 实现的压测工具，每个用户用一个协程的方式模拟，最大限度的利用 CPU 资源

### 1.2 项目体验

- 可以在 mac/linux/windows 不同平台下执行的命令

- [go-stress-testing](https://github.com/link1st/go-stress-testing/releases) 压测工具下载地址

参数说明:

`-c` 表示并发数

`-n` 每个并发执行请求的次数，总请求的次数 = 并发数 `*` 每个并发执行请求的次数

`-u` 需要压测的地址

```shell

# 运行 以mac为示例
./go-stress-testing-mac -c 1 -n 100 -u https://www.baidu.com/

```

- 压测结果展示

执行以后，终端每秒钟都会输出一次结果，压测完成以后输出执行的压测结果

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

参数解释:

**耗时**: 程序运行耗时。程序每秒钟输出一次压测结果

**并发数**: 并发数，启动的协程数

**成功数**: 压测中，请求成功的数量

**失败数**: 压测中，请求失败的数量

**qps**: 当前压测的QPS(每秒钟处理请求数量)

**最长耗时**: 压测中，单个请求最长的响应时长

**最短耗时**: 压测中，单个请求最短的响应时长

**平均耗时**: 压测中，单个请求平均的响应时长

**错误码**: 压测中，接口返回的 code码:返回次数的集合

## 2、压测
### 2.1 压测是什么

压测，即压力测试，是确立系统稳定性的一种测试方法，通常在系统正常运作范围之外进行，以考察其功能极限和隐患。

主要检测服务器的承受能力，包括用户承受能力（多少用户同时玩基本不影响质量）、流量承受等。

### 2.2 为什么要压测

- 压测的目的就是通过压测(模拟真实用户的行为)，测算出机器的性能(单台机器的QPS)，从而推算出系统在承受指定用户数(100W)时，需要多少机器能支撑得住
- 压测是在上线前为了应对未来可能达到的用户数量的一次预估(提前演练)，压测以后通过优化程序的性能或准备充足的机器，来保证用户的体验。

### 2.3 压测名词解释
#### 2.3.1 压测类型解释

| 压测类型 |   解释  |
| :----   | :---- |
| 压力测试(Stress Testing)          |  也称之为强度测试，测试一个系统的最大抗压能力，在强负载(大数据、高并发)的情况下，测试系统所能承受的最大压力，预估系统的瓶颈    |
| 并发测试(Concurrency Testing)     |  通过模拟很多用户同一时刻访问系统或对系统某一个功能进行操作，来测试系统的性能，从中发现问题(并发读写、线程控制、资源争抢)      |
| 耐久性测试(Configuration Testing) |  通过对系统在大负荷的条件下长时间运行，测试系统、机器的长时间运行下的状况,从中发现问题(内存泄漏、数据库连接池不释放、资源不回收)     |


#### 2.3.2 压测名词解释

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
| 最大响应时间(Max Response Time) | 在一次压测中，从发出请求或指令系统做出的反映(响应)的最大时间  |
| 最少响应时间(Mininum Response Time) | 在一次压测中，从发出请求或指令系统做出的反映(响应)的最少时间  |
| 平均响应时间(Average Response Time) | 在一次压测中，从发出请求或指令系统做出的反映(响应)的平均时间  |

#### 2.3.3 机器性能指标解释

| 机器性能 |   解释  |
| :----   | :---- |
| CUP利用率(CPU Usage)       |  CUP 利用率分用户态、系统态和空闲态，CPU利用率是指:CPU执行非系统空闲进程的时间与CPU总执行时间的比率      |
| 内存使用率(Memory usage)    |  内存使用率指的是此进程所开销的内存。      |
| IO(Disk input/ output)    |  磁盘的读写包速率       |
| 网卡负载(Network Load)      |  网卡的进出带宽,包量       |

#### 2.3.4 访问指标解释

| 访问 |   解释  |
| :----   | :---- |
| PV(页面浏览量 Page View)           |  用户每打开1个网站页面，记录1个PV。用户多次打开同一页面，PV值累计多次      |
| UV(网站独立访客 Unique Visitor)    |  通过互联网访问、流量网站的自然人。1天内相同访客多次访问网站，只计算为1个独立访客       |

### 2.4 如何计算压测指标

- 压测我们需要有目的性的压测，这次压测我们需要达到什么目标(如:单台机器的性能为 100QPS?网站能同时满足100W人同时在线)
- 可以通过以下计算方法来进行计算:
- 压测原则:每天80%的访问量集中在20%的时间里，这20%的时间就叫做峰值
- 公式: ( 总PV数`*`80% ) / ( 每天的秒数`*`20% ) = 峰值时间每秒钟请求数(QPS)
- 机器: 峰值时间每秒钟请求数(QPS) / 单台机器的QPS = 需要的机器的数量

- 假设:网站每天的用户数(100W)，每天的用户的访问量约为3000W PV，这台机器的需要多少QPS?
> ( 30000000\*0.8 ) / (86400 * 0.2) ≈ 1389 (QPS)

- 假设:单台机器的的QPS是69，需要需要多少台机器来支撑？
> 1389 / 69 ≈ 20

## 3、常见的压测工具
### 3.1 ab
- 简介

ApacheBench 是 Apache 服务器自带的一个web压力测试工具，简称 ab。ab 又是一个命令行工具，对发起负载的本机要求很低，根据 ab 命令可以创建很多的并发访问线程，模拟多个访问者同时对某一 URL 地址进行访问，因此可以用来测试目标服务器的负载压力。总的来说 ab 工具小巧简单，上手学习较快，可以提供需要的基本性能指标，但是没有图形化结果，不能监控。

ab 属于一个轻量级的压测工具，结果不会特别准确，可以用作参考。

- 安装

```shell
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

```shell
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

- `Requests per second` 吞吐量，指的是某个并发用户下单位时间内处理的请求数。等效于 QPS，其实可以看作同一个统计方式，只是叫法不同而已。

- `Time per request` 用户平均请求等待时间

- `Time per request` 服务器处理时间

### 3.2 Locust

- 简介

是非常简单易用、分布式、python 开发的压力测试工具。有图形化界面，支持将压测数据导出。

- 安装

```shell
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
from locust import HttpUser, TaskSet, task

# 定义用户行为
class UserBehavior(TaskSet):

    @task
    def baidu_index(self):
        self.client.get("/")

class WebsiteUser(HttpUser):
    task = [UserBehavior] # 指向一个定义的用户行为类
    min_wait = 3000 # 执行事务之间用户等待时间的下界（单位：毫秒）
    max_wait = 6000 # 执行事务之间用户等待时间的上界（单位：毫秒）
```

- 启动压测

```shell
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

### 3.3 JMeter

- 简介

Apache JMeter是Apache组织开发的基于Java的压力测试工具。用于对软件做压力测试，它最初被设计用于Web应用测试，但后来扩展到其他测试领域。
JMeter能够对应用程序做功能/回归测试，通过创建带有断言的脚本来验证你的程序返回了你期望的结果。

- 安装

访问 https://jmeter-plugins.org/install/Install/ 下载解压以后即可使用

- 用法

JMeter的功能过于强大，这里暂时不介绍用法，可以查询相关文档使用(参考文献中有推荐的教程文档)


### 3.4 云压测

#### 3.4.1 云压测介绍

顾名思义就是将压测脚本部署在云端，通过云端对对我们的应用进行全方位压测，只需要配置压测的参数，无需准备实体机，云端自动给我们分配需要压测的云主机，对被压测目标进行压测。

云压测的优势:

1. 轻易的实现分布式部署
2. 能够模拟海量用户的访问
3. 流量可以从全国各地发起，更加真实的反映用户的体验
4. 全方位的监控压测指标
5. 文档比较完善

当然了云压测是一款商业产品，在使用的时候自然还是需要收费的，而且价格还是比较昂贵的~

#### 3.4.2 阿里云 性能测试 PTS

PTS（Performance Testing Service）是面向所有技术背景人员的云化测试工具。有别于传统工具的繁复，PTS以互联网化的交互，提供性能测试、API调试和监测等多种能力。自研和适配开源的功能都可以轻松模拟任意体量的用户访问业务的场景，任务随时发起，免去繁琐的搭建和维护成本。更是紧密结合监控、流控等兄弟产品提供一站式高可用能力，高效检验和管理业务性能。

阿里云同样还是支持渗透测试，通过模拟黑客对业务系统进行全面深入的安全测试。


#### 3.4.3 腾讯云 压测大师 LM

通过创建虚拟机器人模拟多用户的并发场景，提供一整套完整的服务器压测解决方案


## 4、go-stress-testing go语言实现的压测工具

### 4.1 介绍

- go-stress-testing 是go语言实现的简单压测工具，源码开源、支持二次开发，可以压测http、webSocket请求、私有rpc调用，使用协程模拟单个用户，可以更高效的利用CPU资源。

- 项目地址 [https://github.com/link1st/go-stress-testing](https://github.com/link1st/go-stress-testing)

### 4.2 用法

- [go-stress-testing](https://github.com/link1st/go-stress-testing/releases) 下载地址
- clone 项目源码运行的时候，需要将项目 clone 到 **$GOPATH** 目录下
- 支持参数:

```
Usage of ./go-stress-testing-mac:
  -c uint
      并发数 (default 1)
  -n uint
      请求数(单个并发/协程) (default 1)
  -u string
      压测地址
  -d string
      调试模式 (default "false")
  -H value
      自定义头信息传递给服务器 示例:-H 'Content-Type: application/json'
  -data string
      HTTP POST方式传送数据
  -v string
      验证方法 http 支持:statusCode、json webSocket支持:json
  -p string
      curl文件路径
```

- `-n` 是单个用户请求的次数，请求总次数 = `-c`* `-n`， 这里考虑的是模拟用户行为，所以这个是每个用户请求的次数

- 下载以后执行下面命令即可压测

- 使用示例:

```
# 查看用法
./go-stress-testing-mac

# 使用请求百度页面
./go-stress-testing-mac -c 1 -n 100 -u https://www.baidu.com/

# 使用debug模式请求百度页面
./go-stress-testing-mac -c 1 -n 1 -d true -u https://www.baidu.com/

# 使用 curl文件(文件在curl目录下) 的方式请求
./go-stress-testing-mac -c 1 -n 1 -p curl/baidu.curl.txt

# 压测webSocket连接
./go-stress-testing-mac -c 10 -n 10 -u ws://127.0.0.1:8089/acc
```

- 完整压测命令示例
```shell script
# 更多参数 支持 header、post body
go run main.go -c 1 -n 1 -d true -u 'https://page.aliyun.com/delivery/plan/list' \
  -H 'authority: page.aliyun.com' \
  -H 'accept: application/json, text/plain, */*' \
  -H 'content-type: application/x-www-form-urlencoded' \
  -H 'origin: https://cn.aliyun.com' \
  -H 'sec-fetch-site: same-site' \
  -H 'sec-fetch-mode: cors' \
  -H 'sec-fetch-dest: empty' \
  -H 'referer: https://cn.aliyun.com/' \
  -H 'accept-language: zh-CN,zh;q=0.9' \
  -H 'cookie: aliyun_choice=CN; JSESSIONID=J8866281-CKCFJ4BUZ7GDO9V89YBW1-KJ3J5V9K-GYUW7; maliyun_temporary_console0=1AbLByOMHeZe3G41KYd5WWZvrM%2BGErkaLcWfBbgveKA9ifboArprPASvFUUfhwHtt44qsDwVqMk8Wkdr1F5LccYk2mPCZJiXb0q%2Bllj5u3SQGQurtyPqnG489y%2FkoA%2FEvOwsXJTvXTFQPK%2BGJD4FJg%3D%3D; cna=L3Q5F8cHDGgCAXL3r8fEZtdU; isg=BFNThsmSCcgX-sUcc5Jo2s2T4tF9COfKYi8g9wVwr3KphHMmjdh3GrHFvPTqJD_C; l=eBaceXLnQGBjstRJBOfwPurza77OSIRAguPzaNbMiT5POw1B5WAlWZbqyNY6C3GVh6lwR37EODnaBeYBc3K-nxvOu9eFfGMmn' \
  -data 'adPlanQueryParam=%7B%22adZone%22%3A%7B%22positionList%22%3A%5B%7B%22positionId%22%3A83%7D%5D%7D%2C%22requestId%22%3A%2217958651-f205-44c7-ad5d-f8af92a6217a%22%7D'
```

- 使用 curl文件进行压测

curl是Linux在命令行下的工作的文件传输工具，是一款很强大的http命令行工具。

使用curl文件可以压测使用非GET的请求，支持设置http请求的 method、cookies、header、body等参数


**I:** chrome 浏览器生成 curl文件，打开开发者模式(快捷键F12)，如图所示，生成 curl 在终端执行命令
![chrome cURL](https://img.mukewang.com/5d60eddd0001f4b016661114.png)

**II:** postman 生成 curl 命令
![postman cURL](https://img.mukewang.com/5ed79b590001837120581530.png)

生成内容粘贴到项目目录下的**curl/baidu.curl.txt**文件中，执行下面命令就可以从curl.txt文件中读取需要压测的内容进行压测了

```
# 使用 curl文件(文件在curl目录下) 的方式请求
go run main.go -c 1 -n 1 -p curl/baidu.curl.txt
```


### 4.3 实现

- 具体需求可以查看项目源码

- 项目目录结构

```
|____main.go                      // main函数，获取命令行参数
|____server                       // 处理程序目录
| |____dispose.go                 // 压测启动，注册验证器、启动统计函数、启动协程进行压测
| |____statistics                 // 统计目录
| | |____statistics.go            // 接收压测统计结果并处理
| |____golink                     // 建立连接目录
| | |____http_link.go             // http建立连接
| | |____websocket_link.go        // webSocket建立连接
| |____client                     // 请求数据客户端目录
| | |____http_client.go           // http客户端
| | |____websocket_client.go      // webSocket客户端
| |____verify                     // 对返回数据校验目录
| | |____http_verify.go           // http返回数据校验
| | |____websokcet_verify.go      // webSocket返回数据校验
|____heper                        // 通用函数目录
| |____heper.go                   // 通用函数
|____model                        // 模型目录
| |____request_model.go           // 请求数据模型
| |____curl_model.go              // curl文件解析
|____vendor                       // 项目依赖目录
```


### 4.4 go-stress-testing 对 Golang web 压测


这里使用go-stress-testing对go server进行压测(部署在同一台机器上)，并统计压测结果

- 申请的服务器配置

CPU: 4核 (Intel Xeon(Cascade Lake) Platinum 8269  2.5 GHz/3.2 GHz)

内存: 16G
硬盘: 20G SSD
系统: CentOS 7.6

go version: go1.12.9 linux/amd64

![go-stress-testing01](https://img.mukewang.com/5d64a48e0001bb8421170573.png)

- go server

```golang
package main

import (
    "log"
    "net/http"
    "runtime"
)

const (
    httpPort = "8088"
)

func main() {

    runtime.GOMAXPROCS(runtime.NumCPU() - 1)

    hello := func(w http.ResponseWriter, req *http.Request) {
        data := "Hello, go-stress-testing! \n"

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

- go_stress_testing 压测命令

```
./go-stress-testing-linux -c 100 -n 10000 -u http://127.0.0.1:8088/
```


- 压测结果
- [压测结果 示例](https://github.com/link1st/go-stress-testing/issues/32)

| 并发数  |  go_stress_testing QPS  |
| :----: |  :----:  |
|   1    | 6394.86  |
|   4    | 16909.36 |
|   10   | 18456.81 |
|   20   | 19490.50 |
|   30   | 19947.47 |
|   50   | 19922.56 |
|   80   | 19155.33 |
|   100  | 18336.46 |

从压测的结果上看：效果还不错，压测QPS有接近2W

### 4.5 grpc压测
- 介绍如何压测 grpc 接口
> [添加对 grpc 接口压测 commit](https://github.com/link1st/go-stress-testing/commit/2b4b14aaf026d08276531cf76f42de90efd3bc61)
- 1. 启动Server
```shell script
# 进入 grpc server 目录
cd tests/grpc

# 启动 grpc server
go run main.go
```

- 2. 对 grpc server 协议进行压测
```
# 回到项目根目录
go run main.go -c 300 -n 1000 -u grpc://127.0.0.1:8099 -data world

开始启动  并发数:300 请求数:1000 请求参数:
request:
 form:grpc
 url:grpc://127.0.0.1:8099
 method:POST
 headers:map[Content-Type:application/x-www-form-urlencoded; charset=utf-8]
 data:world
 verify:
 timeout:30s
 debug:false

─────┬───────┬───────┬───────┬────────┬────────┬────────┬────────┬────────┬────────┬────────
 耗时 │ 并发数 │ 成功数 │ 失败数 │   qps  │最长耗时  │最短耗时 │平均耗时  │下载字节 │字节每秒  │ 错误码
─────┼───────┼───────┼───────┼────────┼────────┼────────┼────────┼────────┼────────┼────────
   1s│    186│  14086│      0│34177.69│   22.40│    0.63│    8.78│        │        │200:14086
   2s│    265│  30408│      0│26005.09│   32.68│    0.63│   11.54│        │        │200:30408
   3s│    300│  46747│      0│21890.46│   40.84│    0.63│   13.70│        │        │200:46747
   4s│    300│  62837│      0│20057.06│   45.81│    0.63│   14.96│        │        │200:62837
   5s│    300│  79119│      0│19134.52│   45.81│    0.63│   15.68│        │        │200:79119
```

- 如何扩展其它私有协议
> 由于私有协议、grpc 协议 都涉及到代码的书写，所以需要 编写go 的代码才能完成
> 参考 [添加对 grpc 接口压测 commit](https://github.com/link1st/go-stress-testing/commit/2b4b14aaf026d08276531cf76f42de90efd3bc61)

## 5、压测工具的比较
### 5.1 比较

| -         |  ab     | locust  | Jmeter  | go-stress-testing  | 云压测  |
| :----     |  :----  |  :----  |  :----  |  :----             |  :---- |
|   实现语言 |    C    |  Python |  Java   |      Golang        |  -     |
|   UI界面  |    无   |   有     |    有   |        无          |    无   | 
|   优势    |  使用简单，上手简单   |  支持分布式、压测数据支持导出   | 插件丰富，支持生成HTML报告   |      项目开源，使用简单，没有依赖，支持webSocket压测         |   更加真实的模拟用户，支持更高的压测力度  |


### 5.2 如何选择压测工具

这个世界上**没有最好的，只有最适合的**，工具千千万，选择一款适合你的才是最重要的

在实际使用中有各种场景，选择工具的时候就需要考虑这些:

* 明确你的目的，需要做什么压测、压测的目标是什么？

* 使用的工具你是否熟悉，你愿意花多大的成本了解它？

* 你是为了测试还是想了解其中的原理？

* 工具是否能支持你需要压测的场景


## 6、单台机器100w连接压测实战
### 6.1 说明

之前写了一篇文章，[基于websocket单台机器支持百万连接分布式聊天(IM)系统](https://github.com/link1st/gowebsocket)(不了解这个项目可以查看上一篇或搜索一下文章)，这里我们要实现单台机器支持100W连接的压测

目标:

* 单台机器能保持100W个长连接
* 机器的CPU、内存、网络、I/O 状态都正常

说明:

gowebsocket 分布式聊天(IM)系统:

* 之前用户连接以后有个全员广播，这里需要将用户连接、退出等事件关闭


- 服务器准备:
> 由于自己手上没有自己的服务器，所以需要临时购买的云服务器

压测服务器:

16台(稍后解释为什么需要16台机器)

CPU: 2核
内存: 8G
硬盘: 20G
系统: CentOS 7.6


![webSocket压测服务器](https://img.mukewang.com/5d64ce2d000126cd19970588.png)

被压测服务:

1台

CPU: 4核
内存: 32G
硬盘: 20G SSD
系统: CentOS 7.6

![webSocket被压测服务器](https://img.mukewang.com/5d64cdfd00013d9a19890606.png)


### 6.2 内核优化

- 修改程序最大打开文件数

被压测服务器需要保持100W长连接，客户和服务器端是通过socket通讯的，每个连接需要建立一个socket，程序需要保持100W长连接就需要单个程序能打开100W个文件句柄


```
# 查看系统默认的值
ulimit -n
# 设置最大打开文件数
ulimit -n 1040000
```

这里设置的要超过100W，程序除了有100W连接还有其它资源连接(数据库、资源等连接)，这里设置为 104W

centOS 7.6 上述设置不生效，需要手动修改配置文件

`vim /etc/security/limits.conf`

这里需要把硬限制和软限制、root用户和所有用户都设置为 1040000

core 是限制内核文件的大小，这里设置为 unlimited

```
# 添加以下参数
root soft nofile 1040000
root hard nofile 1040000

root soft nofile 1040000
root hard nproc 1040000

root soft core unlimited
root hard core unlimited

* soft nofile 1040000
* hard nofile 1040000

* soft nofile 1040000
* hard nproc 1040000

* soft core unlimited
* hard core unlimited
```

注意:

`/proc/sys/fs/file-max` 表示系统级别的能够打开的文件句柄的数量，不能小于limits中设置的值

如果file-max的值小于limits设置的值会导致系统重启以后无法登录

```
# file-max 设置的值参考
cat /proc/sys/fs/file-max
12553500
```

修改以后重启服务器，`ulimit -n` 查看配置是否生效


### 6.3 客户端配置

由于linux端口的范围是 `0~65535(2^16-1)`这个和操作系统无关，不管linux是32位的还是64位的

这个数字是由于tcp协议决定的，tcp协议头部表示端口只有16位，所以最大值只有65535(如果每台机器多几个虚拟ip就能突破这个限制)

1024以下是系统保留端口，所以能使用的1024到65535

如果需要100W长连接，每台机器有 65535-1024 个端口， 100W / (65535-1024) ≈ 15.5，所以这里需要16台服务器

- `vim /etc/sysctl.conf` 在文件末尾添加

```
net.ipv4.ip_local_port_range = 1024 65000
net.ipv4.tcp_mem = 786432 2097152 3145728
net.ipv4.tcp_rmem = 4096 4096 16777216
net.ipv4.tcp_wmem = 4096 4096 16777216
```

`sysctl -p` 修改配置以后使得配置生效命令

配置解释:

- `ip_local_port_range` 表示TCP/UDP协议允许使用的本地端口号 范围:1024~65000
- `tcp_mem` 确定TCP栈应该如何反映内存使用，每个值的单位都是内存页（通常是4KB）。第一个值是内存使用的下限；第二个值是内存压力模式开始对缓冲区使用应用压力的上限；第三个值是内存使用的上限。在这个层次上可以将报文丢弃，从而减少对内存的使用。对于较大的BDP可以增大这些值（注意，其单位是内存页而不是字节）
- `tcp_rmem` 为自动调优定义socket使用的内存。第一个值是为socket接收缓冲区分配的最少字节数；第二个值是默认值（该值会被rmem_default覆盖），缓冲区在系统负载不重的情况下可以增长到这个值；第三个值是接收缓冲区空间的最大字节数（该值会被rmem_max覆盖）。
- `tcp_wmem` 为自动调优定义socket使用的内存。第一个值是为socket发送缓冲区分配的最少字节数；第二个值是默认值（该值会被wmem_default覆盖），缓冲区在系统负载不重的情况下可以增长到这个值；第三个值是发送缓冲区空间的最大字节数（该值会被wmem_max覆盖）。

### 6.4 准备


1. 在被压测服务器上启动Server服务(gowebsocket)

2. 查看被压测服务器的内网端口

3. 登录上16台压测服务器，这里我提前把需要优化的系统做成了镜像，申请机器的时候就可以直接使用这个镜像(参数已经调好)

![压测服务器16台准备](https://img.mukewang.com/5d64cb130001f50912630962.png)

4. 启动压测

```
 ./go_stress_testing_linux -c 62500 -n 1  -u ws://192.168.0.74:443/acc
```

`62500*16 = 100W `正好可以达到我们的要求

建立连接以后，`-n 1`发送一个**ping**的消息给服务器，收到响应以后保持连接不中断

5. 通过 gowebsocket服务器的http接口，实时查询连接数和项目启动的协程数

6. 压测过程中查看系统状态

```
# linux 命令
ps      # 查看进程内存、cup使用情况
iostat  # 查看系统IO情况
nload   # 查看网络流量情况
/proc/pid/status # 查看进程状态
```

### 6.5 压测数据

- 压测以后，查看连接数到100W，然后保持10分钟观察系统是否正常

- 观察以后，系统运行正常、CPU、内存、I/O 都正常，打开页面都正常

- 压测完成以后的数据

查看goWebSocket连接数统计，可以看到 **clientsLen**连接数为100W，**goroutine**数量2000008个，每个连接两个goroutine加上项目启动默认的8个。这里可以看到连接数满足了100W

![查看goWebSocket连接数统计](https://img.mukewang.com/5d64ca86000119ad10080892.png)

从压测服务上查看连接数是否达到了要求，压测完成的统计数据并发数为62500，是每个客户端连接的数量,总连接数： `62500*16=100W`，

![压测服务16台 压测完成](https://img.mukewang.com/5d64ca1d00015a1412630962.png)

- 记录内存使用情况，分别记录了1W到100W连接数内存使用情况

| 连接数      |  内存 |
| :----:     | :----:|
|   10000    | 281M  |
|   100000   | 2.7g  |
|   200000   | 5.4g  |
|   500000   | 13.1g |
|   1000000  | 25.8g |


100W连接时的查看内存详细数据:

```
cat /proc/pid/status
VmSize: 27133804 kB
```

`27133804/1000000≈27.1` 100W连接，占用了25.8g的内存，粗略计算了一下，一个连接占用了27.1Kb的内存，由于goWebSocket项目每个用户连接起了两个协程处理用户的读写事件，所以内存占用稍微多一点

如果需要如何减少内存使用可以参考 **@Roy11568780** 大佬给的解决方案
> 传统的golang中是采用的一个goroutine循环read的方法对应每一个socket。实际百万链路场景中这是巨大的资源浪费，优化的原理也不是什么新东西，golang中一样也可以使用epoll的，把fd拿到epoll中，检测到事件然后在协程池里面去读就行了，看情况读写分别10-20的协程goroutine池应该就足够了

至此，压测已经全部完成，单台机器支持100W连接已经满足~

## 7.常见问题
- **Q:** 压测过程中会出现大量 **TIME_WAIT**

 A: 参考TCP四次挥手原理，主动关闭连接的一方会出现 **TIME_WAIT** 状态，等待的时长为 2MSL(约1分钟左右)

 原因是：主动断开的一方回复 ACK 消息可能丢失，TCP 是可靠的传输协议，在没有收到 ACK 消息的另一端会重试，重新发送FIN消息，所以主动关闭的一方会等待 2MSL 时间，防止对方重试，这就出现了大量 **TIME_WAIT** 状态（参考: 四次挥手的最后两次）

TCP 握手：
<img border="0" src="https://img.mukewang.com/5ec504300001aa7b08301233.png" width="830"/>

## 8、总结

到这里压测总算完成，本次压测花费16元巨款。

单台机器支持100W连接是实测是满足的，但是实际业务比较复杂，还是需要持续优化~

本文通过介绍什么是压测，在什么情况下需要压测，通过单台机器100W长连接的压测实战了解Linux内核的参数的调优。如果觉得现有的压测工具不适用，可以自己实现或者是改造成属于自己的自己的工具。


## 9、参考文献

[性能测试工具](https://testerhome.com/topics/17068)

[性能测试常见名词解释](https://blog.csdn.net/r455678/article/details/53063989)

[性能测试名词解释](https://codeigniter.org.cn/forums/blog-39678-2456.html)

[PV、TPS、QPS是怎么计算出来的？](https://www.zhihu.com/question/21556347)

[超实用压力测试工具－ab工具](https://www.jianshu.com/p/43d04d8baaf7)

[Locust 介绍](http://www.testclass.net/locust/introduce)

[Jmeter性能测试 入门](https://www.cnblogs.com/TankXiao/p/4045439.html)

[基于websocket单台机器支持百万连接分布式聊天(IM)系统](https://github.com/link1st/gowebsocket)

[https://github.com/link1st/go-stress-testing](https://github.com/link1st/go-stress-testing)

github 搜:link1st 查看项目 go-stress-testing

### 意见反馈

- 在项目中遇到问题可以直接在这里找找答案或者提问 [issues](https://github.com/link1st/go-stress-testing/issues)
- 也可以添加我的微信(申请信息填写:公司、姓名，我好备注下)，直接反馈给我
<br/>
<p align="center">
     <img border="0" src="https://img.mukewang.com/5eb376b60001ddc208300832.png" alt="添加link1st的微信" width="200"/>
</p>
