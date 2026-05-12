package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
	"twelve_data_client/internal/color"
	"twelve_data_client/internal/constant"
	"twelve_data_client/internal/model"
	"twelve_data_client/internal/services"

	"github.com/gorilla/websocket"
)

var (
	userArguement string
	mode          string
	symbols       = []string{"AAPL", "RY", "RY:TSX", "EUR/USD", "XAU/USD"}
)

func init() {
	flag.StringVar(&userArguement, "apiKey", "{YOUR_API_KEY_HERE}", "go run . --apiKey=\"{YOUR_API_KEY}\"")
	flag.StringVar(&mode, "mode", "api", "go run . --mode=\"api\" 或 \"ws\"" )
}

func main() {
	var (
		messageChannel   = make(chan []byte)
		interruptChannel = make(chan os.Signal, 0x1)
	)

	flag.Parse()

	if userArguement == "{YOUR_API_KEY_HERE}" {
		fmt.Println(color.Red("Missing API Key. Get one at: https://twelvedata.com/"))
		os.Exit(1)
	}
	
	fmt.Println(color.Cyanf("API Key: ...%s", userArguement[len(userArguement)-4:]))
	
	if mode == "api" {
		fmt.Println(color.Dimf("Endpoint: %s", constant.TWELVED_DATA_API_URL))

		
	} else if mode == "ws" {
		fmt.Println(color.Dimf("Endpoint: %s", constant.TWELVED_DATA_WEBSOCKET_URL))
	
		connection, err := services.GetTwelveDataWebSocket(userArguement, model.NewSubscribe(symbols...))
		if err != nil {
			panic(err)
		}

		defer connection.Close()

		const (
			pongWait   = 30 * time.Second
			pingPeriod = (pongWait * 9) / 10
		)

		done := make(chan struct{})

		connection.SetPongHandler(func(string) error {
			connection.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})

		go func() {
			ticker := time.NewTicker(pingPeriod)
			defer ticker.Stop()

			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					connection.SetWriteDeadline(time.Now().Add(10 * time.Second))
					if err := connection.WriteMessage(websocket.PingMessage, nil); err != nil {
						log.Println(color.Redf("心跳发送失败: %v", err))
						return
					}
				}
			}
		}()

		subscriptionConfirmed := false

		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Println("读取协程异常恢复:", r)
				}
			}()

			for {
				_, message, err := connection.ReadMessage()
				if err != nil {
					log.Println(color.Redf("读取消息失败: %v", err))
					close(messageChannel)
					return
				}

				var subResp model.SubscriptionResponse
				if err := json.Unmarshal(message, &subResp); err == nil && len(subResp.Success) > 0 {
					if !subscriptionConfirmed {
						subscriptionConfirmed = true
						fmt.Println(color.Bold(color.Green("  Subscribed")))
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

				var price model.PriceEvent
				if err := json.Unmarshal(message, &price); err == nil && price.Symbol != "" {
					messageChannel <- message
					continue
				}

				log.Println(color.Redf("未知消息类型: %s", string(message)))
			}
		}()

		signal.Notify(interruptChannel, os.Interrupt)
		fmt.Println(color.Blue(">>> 实时流已开启，等待数据 ..."))

		for {
			select {
			case message, ok := <-messageChannel:
				if !ok {
					fmt.Println(color.Yellow("数据流通道已关闭，程序退出。"))
					close(done)
					return
				}

				var price model.PriceEvent
				if err := json.Unmarshal(message, &price); err != nil {
					log.Println(color.Redf("解析价格失败: %v", err))
					continue
				}

				fmt.Printf("  %s  %s  %s\n",
					color.Cyanf("%-12s", price.Symbol),
					color.Yellowf("%12.4f", price.Price),
					color.Dim(price.Exchange))

			case <-interruptChannel:
				fmt.Println(color.Yellow("\n收到停止信号，正在关闭连接..."))
				close(done)
				return
			}
		}
	} else {
		fmt.Println(color.Red("无效模式。请使用 --api 或 --ws。"))
		os.Exit(1)
	}
}