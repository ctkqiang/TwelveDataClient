package model

// Subscription 是发送给 Twelve Data WebSocket 的订阅/取消订阅消息。
//
// 示例请求体：
//
//	{
//	  "action": "subscribe",
//	  "params": {
//	    "symbols": "AAPL,RY,RY:TSX,EUR/USD,BTC/USD"
//	  }
//	}
type Subscription struct {
	Action string             `json:"action"`
	Params SubscriptionParams `json:"params"`
}

// SubscriptionParams 携带订阅消息的参数。
// Symbols 是以逗号分隔的标的代码列表（可附加交易所后缀，例如 "RY:TSX"）。
type SubscriptionParams struct {
	Symbols string `json:"symbols"`
}

// Twelve Data WebSocket 协议的操作类型常量。
const (
	ActionSubscribe   = "subscribe"
	ActionUnsubscribe = "unsubscribe"
	ActionHeartbeat   = "heartbeat"
	ActionReset       = "reset"
)

// NewSubscribe 构建针对指定标的的 "subscribe" 消息。
func NewSubscribe(symbols ...string) Subscription {
	return Subscription{
		Action: ActionSubscribe,
		Params: SubscriptionParams{Symbols: joinSymbols(symbols)},
	}
}

// NewUnsubscribe 构建针对指定标的的 "unsubscribe" 消息。
func NewUnsubscribe(symbols ...string) Subscription {
	return Subscription{
		Action: ActionUnsubscribe,
		Params: SubscriptionParams{Symbols: joinSymbols(symbols)},
	}
}

func joinSymbols(symbols []string) string {
	out := ""
	for i, s := range symbols {
		if i > 0 {
			out += ","
		}
		out += s
	}
	return out
}

// SubscriptionResponse 表示顶级的 JSON 结构
type SubscriptionResponse struct {
	Event   string               `json:"event"`
	Status  string               `json:"status"`
	Success []SubscriptionDetail `json:"success"`
}

// SubscriptionDetail 表示 success 数组中的单个资产条目
type SubscriptionDetail struct {
	Symbol   string `json:"symbol"`
	Exchange string `json:"exchange"`
	Country  string `json:"country"`
	Type     string `json:"type"`
}

// PriceEvent 表示 Twelve Data WebSocket 推送的实时价格事件
//
//	{
//	  "event": "price",
//	  "symbol": "BTC/USD",
//	  "currency_base": "Bitcoin",
//	  "currency_quote": "US Dollar",
//	  "exchange": "Coinbase Pro",
//	  "type": "Digital Currency",
//	  "timestamp": 1778551560,
//	  "price": 81276.71
//	}
type PriceEvent struct {
	Event        string  `json:"event"`
	Symbol       string  `json:"symbol"`
	CurrencyBase string  `json:"currency_base"`
	CurrencyQuote string `json:"currency_quote"`
	Exchange     string  `json:"exchange"`
	Type         string  `json:"type"`
	Timestamp    int64   `json:"timestamp"`
	Price        float64 `json:"price"`
}