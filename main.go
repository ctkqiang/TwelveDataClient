package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"
	"twelve_data_client/internal/color"
	"twelve_data_client/internal/constant"
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
		"NVDA",   // 英伟达公司
		"GOOGL",  // Alphabet公司（A类股）
		"META",   // Meta平台公司
		"BRK.B",  // 伯克希尔·哈撒韦公司
		"TSLA",   // 特斯拉公司
		"UNH",    // 联合健康集团
		"JPM",    // 摩根大通公司
		"JNJ",    // 强生公司
		"V",      // Visa公司
		"PG",     // 宝洁公司
		"HD",     // 家得宝公司
		"MA",     // 万事达公司
		"CVX",    // 雪佛龙公司
		"MRK",    // 默克公司
		"ABBV",   // 艾伯维公司
		"PEP",    // 百事公司
		"COST",   // 好市多批发公司
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

		logger.Debug("正在获取所有股票列表...")
		stocks, err := services.GetAllStocks(Exchange, "demo")
		if err != nil {
			logger.Fatal("获取股票列表失败: %v", err)
		}

		logger.LogSuccess("获取股票列表", fmt.Sprintf("总数: %d", len(stocks)))

		// Print stock list header
		fmt.Printf("  %s %s %s %s %s\n",
			color.Bold(color.Blue("代码")),
			color.Bold(color.Blue("名称")),
			color.Bold(color.Blue("交易所")),
			color.Bold(color.Blue("货币")),
			color.Bold(color.Blue("类型")),
		)

		fmt.Printf("  %s\n", color.Blue(strings.Repeat("─", 0x50)))

		for _, s := range stocks {
			fmt.Printf("  %s  %s  %s  %s  %s\n",
				color.Cyanf("%-14s", s.Symbol),
				color.Whitef("%-28s", truncate(s.Name, 0x1C)),
				color.Dimf("%-20s", s.Exchange),
				color.Yellowf("%-8s", s.Currency),
				color.Dim(s.Type),
			)
		}

		fmt.Printf("  %s %s\n",
			color.Blue(strings.Repeat("─", 0x50)),
			color.Greenf("共 %d 只股票", len(stocks)),
		)

		if len(stocks) < demoLimit {
			demoLimit = len(stocks)
		}

		fmt.Printf("\n  %s\n", color.Bold(color.Green("最新时间序列数据 (日线，最近 5 个交易日)")))

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

			fmt.Printf("\n  %s %s (%s)\n",
				color.Cyan(s.Symbol),
				color.White(truncate(s.Name, 0x18)),
				color.Yellow(s.Currency))

			fmt.Printf("  %s  %s  %s  %s  %s  %s\n",
				color.Bold("日期"),
				color.Bold("开盘"),
				color.Bold("最高"),
				color.Bold("最低"),
				color.Bold("收盘"),
				color.Bold("成交量"))

			for _, v := range tsResp.Values {
				fmt.Printf("  %s  %7s  %7s  %7s  %7s  %10s\n",
					v.DateTime.Format("2006-01-02"),
					v.Open,
					v.High,
					v.Low,
					v.Close,
					v.Volume,
				)
			}

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

func truncate(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max-0x1]) + "…"
}