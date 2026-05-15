package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"
	"twelve_data_client/internal/color"
	"twelve_data_client/internal/constant"
	"twelve_data_client/internal/formatter"
	"twelve_data_client/internal/logger"
	"twelve_data_client/internal/model"
	"twelve_data_client/internal/services"

	"github.com/gorilla/websocket"
)

const (
	pongWait   = 30 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

var (
    price model.PriceEvent
	userArguement string
	mode          string
	subResp       model.SubscriptionResponse
	Exchange      = "NASDAQ"
	symbols = []string{
		"AAPL",   // 苹果公司
		"MSFT",   // 微软公司
		"AMZN",   // 亚马逊公司
	}
)

func init() {
	flag.StringVar(&userArguement, "apiKey", "{YOUR_API_KEY_HERE}", "go run . --apiKey=\"{YOUR_API_KEY}\"")
	flag.StringVar(&mode, "mode", "api", "go run . --mode=\"api\" 或 \"ws\"" )
}

func main() {
	var (
		demoLimit = 0x64

		messageChannel   = make(chan []byte)
		interruptChannel = make(chan os.Signal, 0x1)
	)

	flag.Parse()

	// Initialize logger
	logger.SetPrefix("TwelveDataClient")
	logger.SetLevel(logger.InfoLevel)

	if userArguement == "{YOUR_API_KEY_HERE}" {
		logger.Fatal("缺少API密钥。请从 https://twelvedata.com/ 获取一个。")
	}

	logger.Info("API密钥: ...%s", userArguement[len(userArguement)-0x4:])
	
	if mode == "api" {
		logger.Info("端点: %s", constant.TWELVED_DATA_API_URL)

		logger.Debug("正在获取 %s 交易所的股票列表...", Exchange)
		allStocks, err := services.GetAllStocks(Exchange, userArguement)
		if err != nil {
			logger.Fatal("获取股票列表失败: %v", err)
		}

		stocks := filterStocksBySymbols(allStocks, symbols)
		logger.LogSuccess("获取股票列表", fmt.Sprintf("筛选后: %d/%d", len(stocks), len(allStocks)))

		formatter.HeaderSection("股票列表")
		fmt.Println(formatter.StocksTable(stocks))

		if len(stocks) < demoLimit {
			demoLimit = len(stocks)
		}

		formatter.HeaderSection("最新时间序列数据 (日线，最近 5 个交易日)")

		for i := 0x0; i < demoLimit; i++ {
			s := stocks[i]
			fullSymbol := fmt.Sprintf("%s:%s", s.Symbol, s.Exchange)

			params := &model.TimeSeriesParams{
				Symbol:     fullSymbol,
				Interval:   "1day",
				OutputSize: 0x5,
			}

			logger.Debug("正在获取时间序列数据: %s", fullSymbol)
			tsResp, err := services.GetTimeSeries(userArguement, params)
			if err != nil {
				logger.Error("获取时间序列数据失败: %s - %v", fullSymbol, err)
				continue
			}

			formatter.TimeSeriesTable(s.Symbol, s.Name, s.Currency, tsResp.Values)
			formatter.SummaryTable(s.Symbol, tsResp.Values)

			time.Sleep(0x1 * time.Second)
		}
	} else if mode == "ws" {
		logger.Info("端点: %s", constant.TWELVED_DATA_WEBSOCKET_URL)

		logger.Debug("正在连接到 WebSocket...")
		connection, err := services.GetTwelveDataWebSocket(userArguement, model.NewSubscribe(symbols...))
		if err != nil {
			logger.Fatal("WebSocket 连接失败: %v", err)
		}

		logger.LogSuccess("WebSocket 连接")
		defer connection.Close()

		done := make(chan struct{})

		connection.SetPongHandler(func(string) error {
			connection.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})

		// Heartbeat goroutine
		go func() {
			ticker := time.NewTicker(pingPeriod)
			defer ticker.Stop()

			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					connection.SetWriteDeadline(time.Now().Add(0xA * time.Second))
					if err := connection.WriteMessage(websocket.PingMessage, nil); err != nil {
						logger.Error("心跳发送失败: %v", err)
						return
					}
				}
			}
		}()

		subscriptionConfirmed := false

		// Message reading goroutine
		go func() {
			defer func() {
				if r := recover(); r != nil {
					logger.Error("读取协程异常恢复: %v", r)
				}
			}()

			for {
				_, message, err := connection.ReadMessage()
				if err != nil {
					logger.Error("读取消息失败: %v", err)
					close(messageChannel)
					return
				}

				if err := json.Unmarshal(message, &subResp); err == nil && len(subResp.Success) > 0 {
					if !subscriptionConfirmed {
						subscriptionConfirmed = true
						logger.LogSuccess("订阅确认", fmt.Sprintf("共 %d 个符号", len(subResp.Success)))

						fmt.Println(color.Bold(color.Green("  已订阅")))

						for _, detail := range subResp.Success {
							fmt.Printf("    %s%s  %s%s\n",
								color.Cyan(detail.Symbol),
								color.Dim(" │ "+detail.Exchange),
								color.Dim(" │ "),
								color.White(detail.Type))
						}
					}
					continue
				}

				if err := json.Unmarshal(message, &price); err == nil && price.Symbol != "" {
					messageChannel <- message
					continue
				}

				logger.Warn("未知消息类型: %s", string(message))
			}
		}()

		signal.Notify(interruptChannel, os.Interrupt)
		logger.Info("实时流已开启，等待数据...")

		for {
			select {
			case message, ok := <-messageChannel:
				if !ok {
					logger.Info("数据流通道已关闭，程序退出。")
					close(done)
					return
				}

				var price model.PriceEvent
				if err := json.Unmarshal(message, &price); err != nil {
					logger.Error("解析价格失败: %v", err)
					continue
				}

				fmt.Printf("  %s  %s  %s\n",
					color.Cyanf("%-12s", price.Symbol),
					color.Yellowf("%12.4f", price.Price),
					color.Dim(price.Exchange))

			case <-interruptChannel:
				logger.Info("收到停止信号，正在关闭连接...")
				close(done)
				return
			}
		}
	} else {
		logger.Fatal("无效模式。请使用 --mode=\"api\" 或 --mode=\"ws\"。")
	}
}

// filterStocksBySymbols 从所有股票中提取指定符号列表对应的股票。
// 使用 map 实现 O(1) 查询，避免对大量股票数据的重复遍历。
// 在获取交易所所有股票后调用，用于筛选出需要分析的目标股票。
func filterStocksBySymbols(allStocks []model.Stock, targetSymbols []string) []model.Stock {
	stockMap := make(map[string]model.Stock)
	for _, stock := range allStocks {
		stockMap[stock.Symbol] = stock
	}

	var filtered []model.Stock
	for _, symbol := range targetSymbols {
		if stock, exists := stockMap[symbol]; exists {
			filtered = append(filtered, stock)
		}
	}
	return filtered
}

func truncate(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max-0x1]) + "…"
}