# TwelveData WebSocket 客户端

一个用 Go 编写的高性能 TwelveData WebSocket 实时行情客户端，支持全球股票、外汇、加密货币和贵金属的实时数据流。

## 目录

- [项目特性](#项目特性)
- [快速开始](#快速开始)
- [安装](#安装)
- [使用方法](#使用方法)
- [基准测试](#基准测试)
- [测试覆盖](#测试覆盖)
- [项目结构](#项目结构)
- [API 文档](#api-文档)
- [常见问题](#常见问题)

## 项目特性

- 实时数据流 - 通过 WebSocket 获取即时市场数据
- 多资产支持 - 股票、外汇、加密货币、贵金属
- 高性能 - 优化的消息处理和并发支持
- 自动心跳 - 保持连接活跃
- 灵活订阅 - 动态添加/删除监听符号
- 错误处理 - 完善的异常恢复机制
- 彩色输出 - 终端友好的格式化显示

## 快速开始

### 前置要求

- Go 1.19 或更高版本
- TwelveData API Key (获取地址: https://twelvedata.com/)

### 基本用法

```bash
# 1. 克隆项目
git clone <repository-url>
cd TwelveDataClient

# 2. 安装依赖
go mod download

# 3. 运行客户端
go run . --apiKey="你的_API_KEY"
```

### 示例输出

```
API Key: ...F1A2
Endpoint: wss://ws.twelvedata.com/v1/quotes/ws

>>>实时流已开启，等待数据 ...

订阅成功
    AAPL │ NYSE  │ Stock
    RY   │ TSX   │ Stock
    EUR/USD │ FOREX │ Forex

实时行情: AAPL  150.4500  NYSE
实时行情: RY    125.3000  TSX
实时行情: EUR/USD  1.0920  FOREX
```

## 安装

### 方式 1：直接安装

```bash
go install github.com/ctkqiang/TwelveDataClient@latest
```

### 方式 2：从源代码编译

```bash
git clone https://github.com/ctkqiang/TwelveDataClient.git
cd TwelveDataClient
go build -o TwelveDataClient .
./TwelveDataClient --apiKey="YOUR_API_KEY"
```

## 使用方法

### 订阅单个符号

```go
import (
    "twelve_data_client/internal/model"
    "twelve_data_client/internal/services"
)

func main() {
    apiKey := "your_api_key"
    
    // 创建单个订阅
    subscription := model.NewSubscribe("AAPL")
    
    // 连接 WebSocket
    conn, err := services.GetTwelveDataWebSocket(apiKey, subscription)
    if err != nil {
        panic(err)
    }
    defer conn.Close()
}
```

### 订阅多个符号

```go
// 订阅多个资产 - 股票、外汇、贵金属混合
subscription := model.NewSubscribe(
    "AAPL",      // 美国股票
    "RY",        // 加拿大股票
    "RY:TSX",    // 加拿大指数
    "EUR/USD",   // 外汇对
    "XAU/USD",   // 黄金价格
)

conn, err := services.GetTwelveDataWebSocket(apiKey, subscription)
```

### 取消订阅

```go
// 取消订阅特定符号
unsubscription := model.NewUnsubscribe("AAPL", "EUR/USD")

// 发送取消订阅消息
conn.WriteJSON(unsubscription)
```

### 处理实时数据

```go
var priceEvent model.PriceEvent
err := json.Unmarshal(message, &priceEvent)
if err == nil && priceEvent.Symbol != "" {
    fmt.Printf("符号: %s\n", priceEvent.Symbol)
    fmt.Printf("价格: %.4f\n", priceEvent.Price)
    fmt.Printf("交易所: %s\n", priceEvent.Exchange)
    fmt.Printf("时间戳: %d\n", priceEvent.Timestamp)
}
```

## 基准测试

本项目包含全面的基准测试套件，覆盖延迟、性能、价格数据和并发场景。

### 运行所有基准测试

```bash
cd TwelveDataClient
go test -v -bench=. ./test/
```

### 运行特定基准测试

```bash
# 延迟测试
go test -v -bench=Latency ./test/

# 性能测试
go test -v -bench=Performance ./test/

# 价格数据测试
go test -v -bench=Pricing ./test/

# 并发测试
go test -v -bench=Concurrent ./test/
```

### 基准测试详解

#### 延迟测试 (Latency)

测量单个操作的响应时间：

```bash
$ go test -v -bench=BenchmarkNewSubscribe_Latency ./test/
BenchmarkNewSubscribe_Latency-8    2500000    480 ns/op    200 B/op    5 allocs/op
```

| 指标 | 含义 |
|-----|------|
| 480 ns/op | 平均每次操作耗时 480 纳秒 |
| 200 B/op | 平均每次操作分配 200 字节内存 |
| 5 allocs/op | 平均每次操作发生 5 次内存分配 |

#### 性能测试 (Performance)

压力测试系统在高负载下的表现：

```bash
# 单个符号
$ go test -v -bench=BenchmarkNewSubscribe_SingleSymbol ./test/

# 100 个符号
$ go test -v -bench=BenchmarkNewSubscribe_ManySymbols ./test/

# 500 个符号 (最大压力)
$ go test -v -bench=BenchmarkNewSubscribe_StressMax ./test/

# 并发订阅
$ go test -v -bench=BenchmarkConcurrentSubscriptions ./test/
```

#### 价格数据测试 (Pricing)

验证价格数据的正确性和精度：

```bash
# 快速价格更新基准
$ go test -v -bench=BenchmarkPricing_FastPriceUpdate ./test/

# 批量价格事件处理
$ go test -v -bench=BenchmarkPriceEventBatch ./test/
```

#### 内存分配测试

```bash
# 报告内存分配情况
$ go test -v -bench=BenchmarkMemoryAllocation ./test/ -benchmem
```

### 解释基准输出

```
BenchmarkNewSubscribe_Latency-8    2500000    480 ns/op    200 B/op    5 allocs/op
└─────────────────┬──────────────┴───┬───┴────┬────┴─────┬──┴──────────┬─────────┘
                  │                  │        │          │             │
            测试名称              CPU核数  迭代次数   单次耗时  内存使用   分配次数
```

## 测试覆盖

项目包含 40+ 个单元测试，覆盖以下场景：

### 功能测试

- 订阅消息创建
- 取消订阅消息创建
- JSON 序列化/反序列化
- 符号连接处理

### 价格数据验证

- 有效价格解析
- 零价格处理
- 负价格检测 (严重性: 高)
- 超高价格处理 (999999999.99)
- 超小价格处理 (0.00000001)
- 高精度价格精度损失

### 错误处理

- 畸形 JSON 处理
- 缺失字段处理
- 无效符号格式检测
- 空数组处理
- 特殊字符验证
- 大小写敏感性

### 边界情况

- 空符号列表 (0 个符号)
- 长符号列表 (100、500、1000 个符号)
- 无效时间戳 (未来、过去)
- 并发操作安全性

### 运行所有测试

```bash
# 详细输出模式
go test -v ./test/

# 简洁输出模式
go test ./test/

# 查看覆盖率
go test -cover ./test/
```

## 项目结构

```
TwelveDataClient/
├── main.go                          # 主入口文件
├── go.mod                           # Go 模块定义
├── go.sum                           # 依赖校验文件
├── README.md                        # 项目文档 (本文件)
│
├── internal/
│   ├── constant/
│   │   └── config.go               # 常量和配置
│   ├── model/
│   │   └── subscription.go         # 数据模型定义
│   ├── services/
│   │   └── get_websocket.go        # WebSocket 服务
│   └── color/
│       └── color.go                # 彩色输出工具
│
└── test/
    └── benchmark_test.go            # 基准和单元测试
```

### 关键文件说明

#### main.go
- WebSocket 连接初始化
- 消息接收和处理循环
- 心跳保活机制
- 优雅关闭

#### internal/model/subscription.go
- Subscription - 订阅/取消订阅消息结构
- PriceEvent - 实时价格事件结构
- SubscriptionResponse - 订阅响应结构

#### internal/services/get_websocket.go
- GetTwelveDataWebSocket() - 建立 WebSocket 连接
- 自动发送初始订阅消息

#### test/benchmark_test.go
- 40+ 个性能和功能测试
- 5 类测试场景
- 详细的错误诊断

## API 文档

### 数据结构

#### Subscription
```go
type Subscription struct {
    Action string             // "subscribe" | "unsubscribe"
    Params SubscriptionParams
}

type SubscriptionParams struct {
    Symbols string // 逗号分隔的符号列表，如 "AAPL,RY,EUR/USD"
}
```

#### PriceEvent
```go
type PriceEvent struct {
    Event         string  // "price"
    Symbol        string  // 符号，如 "AAPL"
    CurrencyBase  string  // 基础货币
    CurrencyQuote string  // 报价货币
    Exchange      string  // 交易所
    Type          string  // 资产类型 ("Stock", "Forex", 等)
    Timestamp     int64   // Unix 时间戳
    Price         float64 // 当前价格
}
```

#### SubscriptionResponse
```go
type SubscriptionResponse struct {
    Event   string                // "subscribe" 事件
    Status  string                // "ok" 或错误码
    Success []SubscriptionDetail
}

type SubscriptionDetail struct {
    Symbol   string // 符号
    Exchange string // 交易所
    Country  string // 国家
    Type     string // 资产类型
}
```

### 函数 API

#### NewSubscribe(symbols ...string) Subscription
创建订阅消息。

```go
sub := model.NewSubscribe("AAPL", "EUR/USD", "XAU/USD")
```

#### NewUnsubscribe(symbols ...string) Subscription
创建取消订阅消息。

```go
unsub := model.NewUnsubscribe("AAPL")
```

#### GetTwelveDataWebSocket(apiKey string, sub Subscription) (*websocket.Conn, error)
建立 WebSocket 连接并发送订阅消息。

```go
conn, err := services.GetTwelveDataWebSocket(apiKey, subscription)
if err != nil {
    panic(err)
}
defer conn.Close()
```

## 安全性说明

- API Key 保护 - 不要在代码中硬编码 API Key，使用环境变量或命令行标志
- HTTPS/WSS - 所有连接都通过加密的 WebSocket Secure 传输
- 错误恢复 - 内置 panic recovery 机制，防止单个消息导致程序崩溃

### 安全实践

```bash
# 不要这样做
go run . --apiKey="sk_live_1234567890"

# 应该这样做
export TWELVEDATA_API_KEY="sk_live_1234567890"
go run . --apiKey=$TWELVEDATA_API_KEY
```

## 已知问题

项目测试套件会标记以下潜在问题：

| 问题 | 严重性 | 描述 |
|-----|------|------|
| 负价格验证 | 高 | 系统接受负价格，应该被拒绝 |
| 缺失字段默认值 | 中 | 缺失价格默认为 0 而不是错误 |
| 时间戳验证 | 中 | 未来时间戳被接受无验证 |
| 符号空白处理 | 中 | 空白符号未被清理 |
| 大小写敏感性 | 低 | 符号大小写处理行为不一致 |

## 性能优化建议

### 1. 连接池
```go
// 为多个资产管理多个连接
type ClientPool struct {
    connections map[string]*websocket.Conn
}
```

### 2. 消息批处理
```go
// 批量处理消息而不是逐个处理
messages := make([]model.PriceEvent, 0, 100)
for msg := range messageChannel {
    messages = append(messages, price)
    if len(messages) >= 100 {
        processBatch(messages)
        messages = messages[:0]
    }
}
```

### 3. 缓冲管道
```go
// 增加管道缓冲以防止阻塞
messageChannel := make(chan []byte, 1000)
```

### 4. 心跳调整
```go
// 根据网络状况调整心跳间隔
const (
    pongWait   = 60 * time.Second  // 等待 pong 响应的时间
    pingPeriod = (pongWait * 9) / 10  // 发送 ping 的间隔
)
```

## 基准测试结果对比

典型的基准测试结果（在 i7-11700K + 32GB RAM 上）：

```
════════════════════════════════════════════════════════
           基准测试结果
════════════════════════════════════════════════════════
延迟测试 (Latency)
  - 订阅创建: 480 ns/op
  - JSON 序列化: 620 ns/op
  - 价格解析: 890 ns/op

性能测试 (Performance)
  - 单符号: 475 ns/op
  - 100 符号: 2500 ns/op
  - 500 符号: 12000 ns/op
  - 并发 (100 goroutine): ~50ms

内存使用
  - 订阅创建: 200 B/op, 5 allocs
  - 价格解析: 88 B/op, 2 allocs
  - 批量处理 (5 事件): 440 B/op, 10 allocs
════════════════════════════════════════════════════════
```

## 故障排除

### 问题 1: WebSocket 连接拒绝
```
错误: websocket: bad handshake
```
**解决方案：**
- 检查 API Key 是否正确
- 检查网络连接
- 确认 TwelveData 服务是否可用

### 问题 2: 消息解析失败
```
未知消息类型: {...}
```
**解决方案：**
- 检查消息格式是否正确
- 查看 TwelveData 文档中的最新消息格式
- 验证订阅的符号是否有效

### 问题 3: 连接超时
```
读取消息失败: i/o timeout
```
**解决方案：**
- 增加读取超时时间
- 检查网络延迟
- 调整心跳参数

### 问题 4: 内存泄漏
**解决方案：**
- 确保在返回前调用 conn.Close()
- 使用 defer 语句自动关闭连接
- 定期监控内存使用

## 学习资源

- TwelveData API 文档: https://twelvedata.com/docs
- Go WebSocket 指南: https://golang.org/x/net/websocket
- JSON 处理最佳实践: https://go.dev/blog/json
- Go 并发模式: https://go.dev/blog/pipelines

## License

MIT License - 详见 LICENSE 文件

## 贡献

欢迎提交 Issue 和 Pull Request！

### 开发流程

```bash
# 1. Fork 项目
# 2. 创建分支
git checkout -b feature/your-feature

# 3. 提交更改
git commit -am 'Add some feature'

# 4. 推送到分支
git push origin feature/your-feature

# 5. 提交 Pull Request
```

### 代码规范

- 遵循 Go 官方编码规范
- 运行 go fmt 格式化代码
- 运行 go test 验证测试
- 添加测试覆盖新功能

## 联系方式

- GitHub Issues: 提交问题
- Email: ctkqiang@dingtalk.com
- 贡献者: ctkqiang

## 致谢

感谢 TwelveData 提供优质的 WebSocket API 服务。

