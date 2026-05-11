package main

import (
	"flag"
	"fmt"
)

var userArguement string

func init() {
	flag.StringVar(&userArguement, "apiKey", "{YOUR_API_KEY_HERE}", "go run . --apiKey=\"{YOUR_API_KEY}\"")
}

func main() {
	flag.Parse()

	if userArguement == "{YOUR_API_KEY_HERE}" {
		panic("呜呜... 找不到 TwelveData 的 API Key，人家的脑子转不动了 。快去这里救救我: https://twelvedata.com/")
	}
	fmt.Println("你正在使用 API Key 运行: " + userArguement)

	
}