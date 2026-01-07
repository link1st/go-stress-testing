# go-stress-testing 项目架构文档

## 1. 项目概述

go-stress-testing 是一个用 Go 语言实现的高性能压力测试工具。利用 Go 协程的高并发特性，能够在单台机器上模拟百万级连接，支持 HTTP、WebSocket、gRPC 等多种协议的压测。

### 1.1 核心特性

- **高并发**: 每个用户使用一个协程模拟，充分利用 CPU 资源
- **多协议支持**: HTTP/HTTPS、WebSocket/WSS、gRPC、Radius
- **实时统计**: 每秒输出 QPS、响应时间、错误码等指标
- **灵活验证**: 支持插件式响应验证器扩展
- **跨平台**: 支持 Linux、macOS、Windows

---

## 2. 项目结构

```
go-stress-testing/
│
├── main.go                         # 程序入口，命令行参数解析
│
├── model/                          # 数据模型层
│   ├── request_model.go            # 请求模型定义、验证器注册
│   ├── curl_model.go               # CURL 文件解析器
│   └── curl_model_test.go          # 单元测试
│
├── server/                         # 核心服务层
│   ├── dispose.go                  # 压测调度器（核心入口）
│   │
│   ├── client/                     # 协议客户端
│   │   ├── clienter.go             # 客户端接口定义
│   │   ├── http_client.go          # HTTP 客户端
│   │   ├── websocket_client.go     # WebSocket 客户端
│   │   ├── grpc_client.go          # gRPC 客户端
│   │   └── http_longclinet/        # HTTP 长连接客户端
│   │       └── long_client.go
│   │
│   ├── golink/                     # 连接处理器
│   │   ├── http_link.go            # HTTP 单连接处理
│   │   ├── http_link_many.go       # HTTP 多连接处理
│   │   ├── http_link_weigh.go      # HTTP 加权连接
│   │   ├── websocket_link.go       # WebSocket 连接处理
│   │   ├── grpc_link.go            # gRPC 连接处理
│   │   └── radius_link.go          # Radius 协议处理
│   │
│   ├── verify/                     # 响应验证器
│   │   ├── http_verify.go          # HTTP 响应验证
│   │   └── websokcet_verify.go     # WebSocket 响应验证
│   │
│   └── statistics/                 # 统计模块
│       ├── statistics.go           # 统计收集与输出
│       └── statistics_test.go      # 单元测试
│
├── helper/                         # 工具函数
│   └── helper.go                   # 通用辅助函数
│
├── tools/                          # 工具模块
│   └── sort.go                     # 排序工具
│
├── proto/                          # Protocol Buffers
│   ├── pb.pb.go                    # 生成的 PB 代码
│   └── pb_grpc.pb.go               # 生成的 gRPC 代码
│
├── tests/                          # 测试服务
│   ├── servers.go                  # 测试 HTTP 服务
│   └── grpc/main.go                # 测试 gRPC 服务
│
└── curl/                           # CURL 示例文件
    └── baidu.curl.txt
```

---

## 3. 系统架构图

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              命令行入口 (main.go)                             │
│                                                                             │
│    参数解析: -c(并发数) -n(请求数) -u(URL) -H(Headers) -data(Body) ...         │
└────────────────────────────────────┬────────────────────────────────────────┘
                                     │
                                     ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           请求模型层 (model/)                                 │
│  ┌─────────────────────┐      ┌─────────────────────┐                       │
│  │   Request 结构体     │      │   CURL 解析器        │                       │
│  │  - URL              │      │  - ParseCurlFile()  │                       │
│  │  - Form (协议类型)   │      │  - 支持Header/Body   │                       │
│  │  - Method           │      └─────────────────────┘                       │
│  │  - Headers          │                                                    │
│  │  - Body             │      ┌─────────────────────┐                       │
│  │  - Verify           │      │  验证器注册表         │                       │
│  │  - Timeout          │      │  verifyMapHTTP      │                       │
│  └─────────────────────┘      │  verifyMapWebSocket │                       │
│                               └─────────────────────┘                       │
└────────────────────────────────────┬────────────────────────────────────────┘
                                     │
                                     ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                        压测调度器 (server/dispose.go)                         │
│                                                                             │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │  Dispose(ctx, concurrency, totalNumber, request)                    │   │
│   │                                                                     │   │
│   │  1. 创建结果通道 ch := make(chan *RequestResults, 1000)              │   │
│   │  2. 启动统计协程 go statistics.ReceivingResults(...)                 │   │
│   │  3. 根据协议类型启动 N 个压测协程                                       │   │
│   │  4. 等待所有压测完成 wg.Wait()                                        │   │
│   │  5. 关闭通道，等待统计完成                                             │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
└────────────────────────────────────┬────────────────────────────────────────┘
                                     │
              ┌──────────────────────┼──────────────────────┐
              │                      │                      │
              ▼                      ▼                      ▼
┌──────────────────────┐ ┌──────────────────────┐ ┌──────────────────────┐
│   HTTP 处理流程       │ │  WebSocket 处理流程   │ │   gRPC 处理流程       │
│   (golink/http_*)    │ │(golink/websocket_*)  │ │  (golink/grpc_*)     │
│                      │ │                      │ │                      │
│  for i := 0; i < n   │ │  ws.GetConn()        │ │  grpc.Dial()         │
│    HTTPRequest()     │ │  for i := 0; i < n   │ │  for i := 0; i < n   │
│    验证响应           │ │    SendMessage()     │ │    pb.SayHello()     │
│    发送结果到 ch      │ │    验证响应           │ │    发送结果到 ch      │
│  end                 │ │    发送结果到 ch      │ │  end                 │
└──────────┬───────────┘ └──────────┬───────────┘ └──────────┬───────────┘
           │                        │                        │
           └────────────────────────┼────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                          协议客户端层 (server/client/)                        │
│                                                                             │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐              │
│  │  http_client    │  │websocket_client │  │  grpc_client    │              │
│  │                 │  │                 │  │                 │              │
│  │ HTTPRequest()   │  │ NewWebSocket()  │  │ GrpcRequest()   │              │
│  │ - 创建请求       │  │ GetConn()       │  │ - 建立连接       │              │
│  │ - 设置Header    │  │ Write()         │  │ - 发送请求       │              │
│  │ - 发送/接收      │  │ Read()          │  │ - 接收响应       │              │
│  │ - 返回耗时       │  │ Close()         │  │                 │              │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘              │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                    http_longclinet (HTTP 长连接)                      │    │
│  │    - 连接池管理                                                       │    │
│  │    - Keep-Alive 支持                                                 │    │
│  │    - HTTP/2.0 支持                                                   │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                          响应验证层 (server/verify/)                          │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │  HTTP 验证器                                                         │    │
│  │  - HTTPStatusCode(): 按 HTTP 状态码验证（默认200为成功）                │    │
│  │  - HTTPJson(): 按返回 JSON 中的 code 字段验证                          │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │  WebSocket 验证器                                                    │    │
│  │  - WebSocketJSON(): 按消息 JSON 格式验证                              │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         统计模块 (server/statistics/)                        │
│                                                                             │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │  ReceivingResults(ctx, concurrent, ch, wg)                          │   │
│   │                                                                     │   │
│   │  实时统计:                                                           │   │
│   │  - processingTime (总请求时间)                                       │   │
│   │  - successNum / failureNum (成功/失败数)                             │   │
│   │  - maxTime / minTime (最大/最小响应时间)                              │   │
│   │  - errCode sync.Map (错误码分布)                                     │   │
│   │                                                                     │   │
│   │  每秒输出:                                                           │   │
│   │  ┌─────┬───────┬───────┬───────┬────────┬────────┬────────┐         │   │
│   │  │ 耗时│ 并发数│ 成功数│ 失败数│   qps  │最长耗时│平均耗时│         │   │
│   │  └─────┴───────┴───────┴───────┴────────┴────────┴────────┘         │   │
│   │                                                                     │   │
│   │  最终输出: tp90/tp95/tp99、总请求数、成功率、字节传输等               │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 4. 核心模块详解

### 4.1 请求模型 (model/request_model.go)

#### 4.1.1 Request 结构体

```go
type Request struct {
    URL       string            // 压测目标地址
    Form      string            // 协议类型: http/webSocket/grpc/radius
    Method    string            // HTTP 方法: GET/POST/PUT/DELETE
    Headers   map[string]string // HTTP 请求头
    Body      string            // 请求体
    Verify    string            // 验证方法名称
    Timeout   time.Duration     // 超时时间
    Debug     bool              // 调试模式
    MaxCon    int               // 每连接最大请求数
    HTTP2     bool              // 是否启用 HTTP/2.0
    Keepalive bool              // 是否启用 Keep-Alive
}
```

#### 4.1.2 协议类型常量

```go
const (
    FormTypeHTTP      = "http"       // HTTP/HTTPS 协议
    FormTypeWebSocket = "webSocket"  // WebSocket/WSS 协议
    FormTypeGRPC      = "grpc"       // gRPC 协议
    FormTypeRadius    = "radius"     // Radius 认证协议
)
```

#### 4.1.3 验证器注册机制

```go
// 验证函数类型定义
type VerifyHTTP func(request *Request, response *http.Response,
                     body []byte) (code int, isSucceed bool)

type VerifyWebSocket func(request *Request, seq string,
                          msg []byte) (code int, isSucceed bool)

// 全局验证器映射
var verifyMapHTTP = make(map[string]VerifyHTTP)
var verifyMapWebSocket = make(map[string]VerifyWebSocket)

// 注册函数
func RegisterVerifyHTTP(name string, verify VerifyHTTP) {
    verifyMapHTTP[name] = verify
}
```

### 4.2 压测调度器 (server/dispose.go)

压测调度器是整个系统的核心，负责：
1. 创建结果收集通道
2. 启动统计协程
3. 根据协议类型启动相应数量的压测协程
4. 等待所有压测完成并汇总结果

```go
func Dispose(ctx context.Context, concurrency uint64, totalNumber uint64,
             request *model.Request) {

    // 1. 创建带缓冲的结果通道
    ch := make(chan *model.RequestResults, 1000)

    // 2. 启动统计协程
    var wgReceiving sync.WaitGroup
    wgReceiving.Add(1)
    go statistics.ReceivingResults(ctx, concurrency, ch, &wgReceiving)

    // 3. 启动压测协程
    var wg sync.WaitGroup
    for i := uint64(0); i < concurrency; i++ {
        wg.Add(1)
        switch request.Form {
        case model.FormTypeHTTP:
            go golink.HTTP(ctx, i, ch, totalNumber, &wg, request)
        case model.FormTypeWebSocket:
            ws := client.NewWebSocket(request.URL)
            ws.GetConn()
            go golink.WebSocket(ctx, i, ch, totalNumber, &wg, request, ws)
        case model.FormTypeGRPC:
            go golink.Grpc(ctx, i, ch, totalNumber, &wg, request)
        }
    }

    // 4. 等待压测完成
    wg.Wait()
    close(ch)

    // 5. 等待统计完成
    wgReceiving.Wait()
}
```

### 4.3 连接处理器 (server/golink/)

#### 4.3.1 HTTP 连接处理

```go
// golink/http_link.go
func HTTP(ctx context.Context, chanID uint64, ch chan<- *model.RequestResults,
          totalNumber uint64, wg *sync.WaitGroup, request *model.Request) {
    defer wg.Done()

    for i := uint64(0); i < totalNumber; i++ {
        // 检查上下文是否取消
        select {
        case <-ctx.Done():
            return
        default:
        }

        // 发送请求
        isSucceed, errCode, requestTime, contentLength := send(chanID, request)

        // 构建结果
        result := &model.RequestResults{
            Time:          requestTime,
            IsSucceed:     isSucceed,
            ErrCode:       errCode,
            ReceivedBytes: contentLength,
        }

        // 发送到统计通道
        ch <- result
    }
}
```

#### 4.3.2 WebSocket 连接处理

```go
// golink/websocket_link.go
func WebSocket(ctx context.Context, chanID uint64, ch chan<- *model.RequestResults,
               totalNumber uint64, wg *sync.WaitGroup, request *model.Request,
               ws *client.WebSocket) {
    defer wg.Done()

    for i := uint64(0); i < totalNumber; i++ {
        // 发送消息
        startTime := time.Now()
        ws.Write([]byte(request.Body))

        // 读取响应
        msg, _ := ws.Read()
        requestTime := uint64(time.Since(startTime).Nanoseconds())

        // 验证响应
        errCode, isSucceed := request.GetVerifyWebSocket()(request, "", msg)

        // 发送结果
        ch <- &model.RequestResults{
            Time:      requestTime,
            IsSucceed: isSucceed,
            ErrCode:   errCode,
        }
    }
}
```

### 4.4 协议客户端 (server/client/)

#### 4.4.1 HTTP 客户端

```go
// client/http_client.go
func HTTPRequest(chanID uint64, request *model.Request) (
    *http.Response, uint64, error) {

    // 创建 HTTP 请求
    req, _ := http.NewRequest(request.Method, request.URL,
                              strings.NewReader(request.Body))

    // 设置请求头
    for key, value := range request.Headers {
        req.Header.Set(key, value)
    }

    // 获取客户端（支持长连接/短连接）
    var httpClient *http.Client
    if request.Keepalive {
        httpClient = getKeepAliveClient(chanID, request)
    } else {
        httpClient = getShortClient(request)
    }

    // 发送请求并计时
    startTime := time.Now()
    resp, err := httpClient.Do(req)
    requestTime := uint64(time.Since(startTime).Nanoseconds())

    return resp, requestTime, err
}
```

#### 4.4.2 WebSocket 客户端

```go
// client/websocket_client.go
type WebSocket struct {
    URL     string
    Conn    *websocket.Conn
    Headers http.Header
}

func NewWebSocket(url string) *WebSocket {
    return &WebSocket{URL: url}
}

func (w *WebSocket) GetConn() error {
    // 建立 WebSocket 连接（最多重试3次）
    var err error
    for i := 0; i < 3; i++ {
        w.Conn, _, err = websocket.DefaultDialer.Dial(w.URL, w.Headers)
        if err == nil {
            return nil
        }
    }
    return err
}

func (w *WebSocket) Write(data []byte) error {
    return w.Conn.WriteMessage(websocket.TextMessage, data)
}

func (w *WebSocket) Read() ([]byte, error) {
    _, msg, err := w.Conn.ReadMessage()
    return msg, err
}
```

### 4.5 统计模块 (server/statistics/)

```go
// statistics/statistics.go
func ReceivingResults(ctx context.Context, concurrent uint64,
                      ch <-chan *model.RequestResults, wg *sync.WaitGroup) {
    defer wg.Done()

    var (
        processingTime uint64       // 累计请求时间
        successNum     uint64       // 成功数
        failureNum     uint64       // 失败数
        maxTime        uint64       // 最大响应时间
        minTime        uint64 = math.MaxUint64
        errCode        sync.Map     // 错误码分布
    )

    // 每秒输出统计
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            // 输出实时统计表格
            printTable(processingTime, successNum, failureNum,
                      maxTime, minTime, concurrent)

        case result, ok := <-ch:
            if !ok {
                // 通道关闭，输出最终统计
                printFinalStats(processingTime, successNum, failureNum,
                               maxTime, minTime, errCode)
                return
            }

            // 累计统计
            processingTime += result.Time
            if result.IsSucceed {
                successNum++
            } else {
                failureNum++
            }

            // 更新极值
            if result.Time > maxTime {
                maxTime = result.Time
            }
            if result.Time < minTime {
                minTime = result.Time
            }

            // 统计错误码
            count, _ := errCode.LoadOrStore(result.ErrCode, uint64(0))
            errCode.Store(result.ErrCode, count.(uint64)+1)
        }
    }
}
```

---

## 5. 数据流程图

### 5.1 压测请求数据流

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   用户输入    │────▶│  参数解析     │────▶│ Request对象  │
│  命令行参数   │     │  main.go     │     │  model层     │
└──────────────┘     └──────────────┘     └──────────────┘
                                                 │
                                                 ▼
                     ┌──────────────────────────────────────────┐
                     │           Dispose 调度器                  │
                     │        server/dispose.go                 │
                     └──────────────────────────────────────────┘
                                         │
              ┌──────────────────────────┼──────────────────────────┐
              │                          │                          │
              ▼                          ▼                          ▼
     ┌────────────────┐        ┌────────────────┐        ┌────────────────┐
     │  协程 1         │        │  协程 2         │        │  协程 N         │
     │  HTTP/WS/gRPC  │        │  HTTP/WS/gRPC  │   ...  │  HTTP/WS/gRPC  │
     └────────┬───────┘        └────────┬───────┘        └────────┬───────┘
              │                         │                         │
              │      ┌──────────────────┼──────────────────┐      │
              │      │                  │                  │      │
              ▼      ▼                  ▼                  ▼      ▼
     ┌────────────────────────────────────────────────────────────────┐
     │                    结果通道 (chan *RequestResults)              │
     │                         缓冲区: 1000                           │
     └────────────────────────────────────────────────────────────────┘
                                         │
                                         ▼
                     ┌──────────────────────────────────────────┐
                     │            统计协程                       │
                     │     statistics.ReceivingResults          │
                     │                                          │
                     │  ┌────────────────────────────────────┐  │
                     │  │  每秒输出:                          │  │
                     │  │  - QPS                             │  │
                     │  │  - 成功/失败数                       │  │
                     │  │  - 响应时间 (最大/最小/平均)          │  │
                     │  │  - 错误码分布                        │  │
                     │  └────────────────────────────────────┘  │
                     └──────────────────────────────────────────┘
                                         │
                                         ▼
                     ┌──────────────────────────────────────────┐
                     │              终端输出                     │
                     │  ┌─────┬───────┬───────┬────────┐       │
                     │  │ 耗时│ 并发数│ 成功数│   qps  │ ...   │
                     │  └─────┴───────┴───────┴────────┘       │
                     └──────────────────────────────────────────┘
```

### 5.2 单次请求处理流程

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              单次 HTTP 请求流程                               │
└─────────────────────────────────────────────────────────────────────────────┘

    ┌───────────┐
    │ 开始请求   │
    └─────┬─────┘
          │
          ▼
    ┌───────────────────────────┐
    │ 1. 构建 http.Request 对象  │
    │    - 设置 URL              │
    │    - 设置 Method           │
    │    - 设置 Headers          │
    │    - 设置 Body             │
    └─────────────┬─────────────┘
                  │
                  ▼
    ┌───────────────────────────┐
    │ 2. 获取 HTTP Client        │
    │    - Keepalive? 长连接池    │
    │    - HTTP2? 多路复用        │
    │    - 否则: 短连接           │
    └─────────────┬─────────────┘
                  │
                  ▼
    ┌───────────────────────────┐
    │ 3. 发送请求 (计时开始)      │
    │    client.Do(req)         │
    └─────────────┬─────────────┘
                  │
                  ▼
    ┌───────────────────────────┐
    │ 4. 接收响应 (计时结束)      │
    │    - 读取响应体             │
    │    - 处理 gzip 解压         │
    └─────────────┬─────────────┘
                  │
                  ▼
    ┌───────────────────────────┐
    │ 5. 验证响应                │
    │    - statusCode: 检查状态码 │
    │    - json: 检查返回码       │
    └─────────────┬─────────────┘
                  │
                  ▼
    ┌───────────────────────────┐
    │ 6. 构建结果对象             │
    │    RequestResults {       │
    │      Time: 请求耗时        │
    │      IsSucceed: 是否成功   │
    │      ErrCode: 错误码       │
    │      ReceivedBytes: 字节数 │
    │    }                      │
    └─────────────┬─────────────┘
                  │
                  ▼
    ┌───────────────────────────┐
    │ 7. 发送到统计通道          │
    │    ch <- result           │
    └─────────────┬─────────────┘
                  │
                  ▼
    ┌───────────┐
    │ 请求完成   │
    └───────────┘
```

---

## 6. 扩展开发指南

### 6.1 添加新的验证器

```go
// 1. 在 server/verify/ 目录下创建验证函数
func HTTPCustomVerify(request *model.Request, response *http.Response,
                      body []byte) (code int, isSucceed bool) {
    // 自定义验证逻辑
    // 返回错误码和是否成功
    return http.StatusOK, true
}

// 2. 在 server/dispose.go 的 init() 中注册
func init() {
    model.RegisterVerifyHTTP("custom", verify.HTTPCustomVerify)
}

// 3. 使用: ./go-stress-testing -v custom -u http://...
```

### 6.2 添加新的协议支持

```go
// 1. 在 model/request_model.go 添加协议常量
const FormTypeCustom = "custom"

// 2. 在 server/client/ 创建客户端
// client/custom_client.go
type CustomClient struct {
    // 客户端配置
}

func (c *CustomClient) Request(data []byte) ([]byte, uint64, error) {
    // 实现请求逻辑
}

// 3. 在 server/golink/ 创建连接处理器
// golink/custom_link.go
func Custom(ctx context.Context, chanID uint64, ch chan<- *model.RequestResults,
            totalNumber uint64, wg *sync.WaitGroup, request *model.Request) {
    // 实现压测循环
}

// 4. 在 server/dispose.go 的 Dispose() 中添加分支
case model.FormTypeCustom:
    go golink.Custom(ctx, i, ch, totalNumber, &wg, request)
```

### 6.3 添加新的统计指标

```go
// 在 server/statistics/statistics.go 中修改

// 1. 添加指标变量
var customMetric uint64

// 2. 在结果处理中累计
if result.CustomField > 0 {
    customMetric += result.CustomField
}

// 3. 在输出中展示
fmt.Printf("自定义指标: %d\n", customMetric)
```

---

## 7. 性能优化建议

### 7.1 内核参数优化

```bash
# /etc/sysctl.conf

# 端口范围
net.ipv4.ip_local_port_range = 1024 65000

# TCP 内存
net.ipv4.tcp_mem = 786432 2097152 3145728
net.ipv4.tcp_rmem = 4096 4096 16777216
net.ipv4.tcp_wmem = 4096 4096 16777216

# 文件句柄
fs.file-max = 2000000
```

### 7.2 程序参数建议

| 场景 | 并发数(-c) | 请求数(-n) | 其他参数 |
|------|-----------|-----------|---------|
| HTTP 短连接压测 | 100-500 | 10000 | - |
| HTTP 长连接压测 | 100-500 | 10000 | -k |
| WebSocket 长连接 | 10000-60000 | 1 | - |
| gRPC 压测 | 300-500 | 1000 | - |

### 7.3 资源估算

| 连接数 | 内存占用 | 单连接内存 |
|--------|---------|-----------|
| 10,000 | ~281MB | ~28KB |
| 100,000 | ~2.7GB | ~27KB |
| 1,000,000 | ~25.8GB | ~27KB |

---

## 8. 附录

### 8.1 命令行参数完整列表

| 参数 | 类型 | 默认值 | 说明 |
|------|------|-------|------|
| -c | uint | 1 | 并发数 |
| -n | uint | 1 | 单个并发请求数 |
| -u | string | - | 压测地址 |
| -d | string | "false" | 调试模式 |
| -k | bool | false | 启用 Keep-Alive |
| -http2 | bool | false | 启用 HTTP/2.0 |
| -m | int | 1 | 单个 host 最大连接数 |
| -H | string | - | 自定义 Header (可多个) |
| -data | string | - | POST 请求体 |
| -v | string | - | 验证方法 |
| -p | string | - | CURL 文件路径 |

### 8.2 支持的验证方法

| 验证方法 | 协议 | 说明 |
|---------|------|------|
| statusCode | HTTP | 按 HTTP 状态码验证 (200=成功) |
| json | HTTP | 按返回 JSON 中 code 字段验证 |
| json | WebSocket | 按消息 JSON 格式验证 |

### 8.3 相关文件

- `build.sh` - 跨平台编译脚本
- `Dockerfile` - 容器化部署配置
- `curl/` - CURL 示例文件目录
- `tests/` - 测试服务代码
