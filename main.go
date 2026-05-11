package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"twelve_data_client/internal/constant"
	"twelve_data_client/internal/services"
)

var userArguement string

func init() {
	flag.StringVar(&userArguement, "apiKey", "{YOUR_API_KEY_HERE}", "go run . --apiKey=\"{YOUR_API_KEY}\"")
}

func main() {
	var (
		messageChannel = make(chan []byte)
		interruptChannel = make(chan os.Signal, 0x1)
	)

	flag.Parse()

	if userArguement == "{YOUR_API_KEY_HERE}" {
		panic("呜呜... 找不到 TwelveData 的 API Key，人家的脑子转不动了 。快去这里救救我: https://twelvedata.com/")
	}

	fmt.Println("你正在使用 API Key 运行: " + userArguement)
	fmt.Println("当前你正在调用: " + constant.TWELVED_DATA_WEBSOCKET_URL)

	connection, err := services.GetTwelveDataWebSocket(userArguement)
	if err != nil {
		panic(err)
	}

	defer connection.Close()

	go func() {
		for {
				_, message, err := connection.ReadMessage()
				if err != nil {
					close(messageChannel)
					panic(err)
				}

				messageChannel <- message
			}
	}()

	signal.Notify(interruptChannel, os.Interrupt)
	fmt.Println("实时流已开启，等待数据...")

	for {
		select {
		case message, ok := <-messageChannel:
			if !ok {
				fmt.Println("数据流通道已关闭，程序退出。")
				return
			}
			fmt.Printf("实时行情: %s\n", message)

		case <-interruptChannel:
			fmt.Println("收到停止信号，正在关闭连接...")
			return
		}
	}
	
}