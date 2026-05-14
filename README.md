# TwelveData 客户端 - WebSocket 和 REST API

一个用 Go 编写的高性能 TwelveData 实时行情客户端，同时支持 WebSocket 实时数据流和 REST API 股票查询。支持全球股票、外汇、加密货币、贵金属的实时数据获取和历史数据查询。

## 目录

- [项目特性](#项目特性)
- [快速开始](#快速开始)
- [安装](#安装)
- [项目架构](#项目架构)
- [使用方法](#使用方法)
  - [WebSocket 实时数据](#websocket-实时数据)
  - [REST API 股票查询](#rest-api-股票查询)
- [开发和测试](#开发和测试)
  - [基准测试](#基准测试)
  - [API 测试](#api-测试)
- [项目结构](#项目结构)
- [API 文档](#api-文档)
- [常见问题](#常见问题)

## 项目特性
![](./docs/屏幕截图%202026-05-14%20151935.png)

### WebSocket 功能

- 实时数据流 - 通过 WebSocket 获取即时市场数据
- 多资产支持 - 股票、外汇、加密货币、贵金属
- 高性能 - 优化的消息处理和并发支持
- 自动心跳 - 保持连接活跃
- 灵活订阅 - 动态添加/删除监听符号
- 错误处理 - 完善的异常恢复机制
- 彩色输出 - 终端友好的格式化显示

### REST API 功能

- 股票查询 - 获取特定交易所的股票列表
- 批量操作 - 支持获取数百个股票
- 多语言支持 - 支持中文股票名称
- 完整的字段信息 - ISIN、CUSIP、FIGI 等标准化代码
- 全球交易所 - 支持 SZSE、SSE、HKEX、NYSE、NASDAQ 等

### 测试和质量保证

- 40+ 单元测试
- 完整的基准测试套件
- 性能分析工具集成
- 代码覆盖率报告
- 多场景压力测试

## 快速开始

### 前置要求

- Go 1.19 或更高版本
- TwelveData API Key (获取地址: https://twelvedata.com/)

### 基本用法 (WebSocket)

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

## 项目架构

### 系统架构

```
┌─────────────────────────────────────────────────────────────┐
│                      TwelveDataClient                        │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌─────────────────────┐         ┌────────────────────────┐ │
│  │   WebSocket 层      │         │    REST API 层         │ │
│  ├─────────────────────┤         ├────────────────────────┤ │
│  │ GetTwelveDataWS()   │         │ GetAllStocks()         │ │
│  │ 实时价格流          │         │ 股票列表查询           │ │
│  │ 订阅/取消订阅       │         │ 多交易所支持           │ │
│  └────────┬────────────┘         └────────┬───────────────┘ │
│           │                               │                 │
│           └──────────┬────────────────────┘                 │
│                      │                                       │
│  ┌───────────────────▼────────────────────────────────────┐ │
│  │            数据模型层 (Model)                          │ │
│  ├──────────────────────────────────────────────────────┤ │
│  │ • Subscription (订阅消息)                             │ │
│  │ • PriceEvent (价格事件)                               │ │
│  │ • Stock (股票信息)                                    │ │
│  │ • SubscriptionResponse (订阅响应)                     │ │
│  └───────────────────────────────────────────────────────┘ │
│                                                               │
│  ┌──────────────┐  ┌──────────┐  ┌────────────────────┐    │
│  │ 配置层       │  │ 工具层   │  │ 彩色输出           │    │
│  ├──────────────┤  ├──────────┤  ├────────────────────┤    │
│  │ 常量定义     │  │ 字符串   │  │ 格式化输出         │    │
│  │ API 端点     │  │ 扩展     │  │ 终端友好           │    │
│  └──────────────┘  └──────────┘  └────────────────────┘    │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

### 目录结构

```
TwelveDataClient/
├── main.go                              主入口文件
├── go.mod                               Go 模块定义
├── go.sum                               依赖校验文件
├── README.md                            项目文档
│
├── internal/
│   ├── constant/
│   │   └── config.go                   常量和配置
│   │
│   ├── model/
│   │   ├── subscription.go             WebSocket 数据模型
│   │   └── stocke.go                   REST API 数据模型
│   │
│   ├── services/
│   │   ├── get_websocket.go            WebSocket 服务
│   │   └── get_rest.go                 REST API 服务
│   │
│   ├── color/
│   │   └── color.go                    彩色输出工具
│   │
│   └── extensions/
│       └── strings.go                  字符串扩展
│
└── test/
    ├── benchmark_test.go               WebSocket 基准测试
    ├── api_test.go                     REST API 测试
    └── RUN_TESTS.md                    测试运行指南
```

## 使用方法

### WebSocket 实时数据

#### 订阅单个符号

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

#### 订阅多个符号

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

#### 处理实时数据

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

#### 取消订阅

```go
// 取消订阅特定符号
unsubscription := model.NewUnsubscribe("AAPL", "EUR/USD")

// 发送取消订阅消息
conn.WriteJSON(unsubscription)
```

### REST API 股票查询

#### 获取深交所股票

```go
import (
    "fmt"
    "log"
    "twelve_data_client/internal/services"
)

func main() {
    // 获取深交所 (SZSE) 的股票列表
    stocks, err := services.GetAllStocks("SZSE", "demo")
    if err != nil {
        log.Fatalf("获取股票失败: %v", err)
    }

    // 遍历股票
    for _, stock := range stocks {
        fmt.Printf("符号: %s, 名称: %s\n", stock.Symbol, stock.Name)
    }

    fmt.Printf("总共获取 %d 个股票\n", len(stocks))
}
```

#### 获取其他交易所股票

```go
// 上海交易所 (SSE)
stocks, err := services.GetAllStocks("SSE", "your_api_key")

// 香港交易所 (HKEX)
stocks, err := services.GetAllStocks("HKEX", "your_api_key")

// 纽约证券交易所 (NYSE)
stocks, err := services.GetAllStocks("NYSE", "your_api_key")
```

#### 处理返回数据

```go
stocks, err := services.GetAllStocks("SZSE", "demo")
if err != nil {
    log.Printf("错误: %v", err)
    return
}

// 访问股票字段
for _, stock := range stocks {
    fmt.Printf("代码: %s\n", stock.Symbol)      // 股票代码
    fmt.Printf("名称: %s\n", stock.Name)        // 股票名称
    fmt.Printf("交易所: %s\n", stock.Exchange)  // 交易所
    fmt.Printf("货币: %s\n", stock.Currency)    // 货币代码
    fmt.Printf("国家: %s\n", stock.Country)     // 国家代码
    fmt.Printf("ISIN: %s\n", stock.ISIN)        // ISIN 代码
    fmt.Printf("CUSIP: %s\n", stock.CUSIP)      // CUSIP 代码
    fmt.Printf("FIGI: %s\n", stock.FIGICode)    // FIGI 代码
}
```

## 开发和测试

### 基准测试

本项目包含全面的基准测试套件，覆盖延迟、性能、价格数据和并发场景。

#### 快速测试

最简单的测试方式:

```bash
# 运行所有单元测试
go test -v ./test/

# 运行所有基准测试
go test -v -bench=. -benchmem ./test/
```

#### 完整的综合测试报告

```bash
go test -v -bench=. -benchmem -cover -timeout=30s ./test/ 2>&1 | tee test_report.txt
```

#### 分类运行测试

```bash
# 只运行单元测试
go test -v -run=^Test ./test/

# 只运行基准测试
go test -v -bench=. -run=^$ ./test/

# 运行特定类别 (WebSocket 服务)
go test -v -run=WebSocketService ./test/

# 运行特定类别 (价格数据)
go test -v -run=Pricing ./test/
```

#### 测试覆盖率

```bash
# 显示覆盖率
go test -cover ./test/

# 生成可视化覆盖率报告
go test -coverprofile=coverage.out ./test/
go tool cover -html=coverage.out -o coverage.html
```

#### 性能分析

```bash
# CPU 分析
go test -bench=. -cpuprofile=cpu.prof ./test/
go tool pprof -http=:8080 cpu.prof

# 内存分析
go test -bench=. -memprofile=mem.prof ./test/
go tool pprof -http=:8080 mem.prof
```

### API 测试

REST API 的完整测试套件包含 15+ 个功能测试和 2 个性能基准。

#### 运行所有 API 测试

```bash
# 运行所有 API 测试
go test -v -run=^TestGetAllStocks ./test/

# 运行 SZSE 交易所测试
go test -v -run=TestGetAllStocks_WithSZSEExchange ./test/

# 运行字段验证测试
go test -v -run=TestGetAllStocks_StockFieldValidation ./test/
```

#### 运行 API 基准

```bash
# 运行所有 API 基准
go test -v -bench=GetAllStocks -benchmem ./test/

# 运行 JSON 反序列化基准
go test -v -bench=BenchmarkGetAllStocks_JSONUnmarshal ./test/
```

#### 测试覆盖范围

| 类别 | 测试数 | 覆盖范围 |
|-----|-------|---------|
| 功能测试 | 8 | 成功响应、空列表、多个股票、字段验证 |
| 错误处理 | 4 | API 错误、HTTP 错误、JSON 错误、超时 |
| 数据验证 | 3 | 字段完整性、多交易所、特殊字符 |
| 性能测试 | 2 | 大规模数据获取、JSON 反序列化 |

### 推荐的命令组合

#### 日常开发 (快速检查)

```bash
go test -v ./test/
```

#### 提交前检查 (详细验证)

```bash
go test -v -bench=. -benchmem -cover ./test/
```

#### 性能优化 (详细分析)

```bash
go test -v -bench=. -benchmem -benchtime=10s -cpuprofile=cpu.prof ./test/
go tool pprof -http=:8080 cpu.prof
```

#### 持续集成 (CI/CD)

```bash
go test -v -race -cover -timeout=60s ./test/ && \
go test -v -bench=. -benchmem ./test/
```

## 项目结构

### 关键文件说明

#### main.go
- WebSocket 连接初始化
- 消息接收和处理循环
- 心跳保活机制
- 优雅关闭

#### internal/model/
- **subscription.go** - WebSocket 订阅消息、价格事件、订阅响应数据结构
- **stocke.go** - REST API 股票信息数据结构

#### internal/services/
- **get_websocket.go** - WebSocket 连接和订阅服务 (GetTwelveDataWebSocket)
- **get_rest.go** - REST API 股票查询服务 (GetAllStocks)

#### internal/constant/
- **config.go** - API 端点、WebSocket URL、常量定义

#### internal/color/
- **color.go** - 终端彩色输出工具

#### test/
- **benchmark_test.go** - WebSocket 基准测试 (40+ 测试)
- **api_test.go** - REST API 单元测试和基准 (15+ 测试)
- **RUN_TESTS.md** - 详细的测试运行指南

## API 文档

### WebSocket API

#### Subscription 结构体

```go
type Subscription struct {
    Action string             // "subscribe" | "unsubscribe"
    Params SubscriptionParams
}

type SubscriptionParams struct {
    Symbols string // 逗号分隔的符号列表，如 "AAPL,RY,EUR/USD"
}
```

#### PriceEvent 结构体

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

#### 函数 API

```go
// 建立 WebSocket 连接并发送订阅消息
func GetTwelveDataWebSocket(apiKey string, sub Subscription) (*websocket.Conn, error)

// 创建订阅消息
func NewSubscribe(symbols ...string) Subscription

// 创建取消订阅消息
func NewUnsubscribe(symbols ...string) Subscription
```

### REST API

#### Stock 结构体

```go
type Stock struct {
    Symbol    string // 股票代码
    Name      string // 股票名称
    Currency  string // 货币代码
    Exchange  string // 交易所代码
    MICCode   string // MIC 代码
    Country   string // 国家代码
    Type      string // 资产类型
    FIGICode  string // FIGI 代码
    CFICode   string // CFI 代码
    ISIN      string // ISIN 代码
    CUSIP     string // CUSIP 代码
}
```

#### 函数 API

```go
// 获取特定交易所的股票列表
func GetAllStocks(exchange string, apiKey string) ([]Stock, error)
```

| 参数 | 类型 | 说明 | 示例 |
|-----|------|------|------|
| exchange | string | 交易所代码 | "SZSE", "SSE", "HKEX" |
| apiKey | string | API 密钥 (空值默认为 "demo") | "demo", "your_key" |

## 常见问题

### Q: 如何修改订阅的符号？

A: 创建一个新的订阅对象并通过 WebSocket 发送，或取消订阅旧的再订阅新的：

```go
// 取消订阅旧符号
unsubscribe := model.NewUnsubscribe("AAPL")
conn.WriteJSON(unsubscribe)

// 订阅新符号
subscribe := model.NewSubscribe("AAPL", "EUR/USD")
conn.WriteJSON(subscribe)
```

### Q: API Key 错误怎么办？

A: 确保：
- API Key 有效且未过期
- 检查网络连接
- 验证 TwelveData 服务是否可用
- 使用 "demo" API Key 进行测试

### Q: 如何获得最佳性能？

A: 
- 使用连接池管理多个连接
- 批量处理消息而不是逐个处理
- 增加管道缓冲防止阻塞
- 根据网络状况调整心跳参数

### Q: 支持哪些交易所？

A: 主要支持：
- SZSE (深圳交易所)
- SSE (上海交易所)
- HKEX (香港交易所)
- NYSE (纽约证券交易所)
- NASDAQ (纳斯达克)
- 以及其他 TwelveData 支持的交易所

### Q: 如何实现断线重连？

A: 监控连接错误并重新建立连接：

```go
for {
    conn, err := services.GetTwelveDataWebSocket(apiKey, subscription)
    if err != nil {
        log.Printf("连接失败: %v，10 秒后重试", err)
        time.Sleep(10 * time.Second)
        continue
    }
    defer conn.Close()
    // 处理连接...
}
```

## 贡献指南

欢迎贡献！请遵循以下步骤：

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/your-feature`)
3. 提交更改 (`git commit -am 'Add some feature'`)
4. 推送到分支 (`git push origin feature/your-feature`)
5. 提交 Pull Request

### 代码规范

- 遵循 Go 官方编码规范
- 运行 `go fmt` 格式化代码
- 运行 `go test` 验证测试
- 为新功能添加测试覆盖

## 许可证

MIT License - 详见 LICENSE 文件

## 联系方式

- GitHub Issues: 提交问题
- Email: ctk@grandpine.com
- 项目维护: ctkqiang

## 致谢

感谢 TwelveData 提供优质的 WebSocket 和 REST API 服务。

---

最后更新: 2026-05-12
版本: 1.0.0
状态: 稳定版本
