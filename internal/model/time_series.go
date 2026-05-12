package model

import "time"

type TimeSeriesBatchResponse map[string]TimeSeriesResponse

type TimeSeriesParams struct {
	Symbol     string    // 单个 symbol，例如 "AAPL" 或 "000001:SZSE"
	Interval   string    // 数据间隔: 1min, 5min, 1h, 1day 等
	StartDate  time.Time // 可选：开始日期
	EndDate    time.Time // 可选：结束日期
	OutputSize int       // 可选：返回数据量，与日期参数互斥，默认 30
	PrePost    bool      // 可选：是否包含盘前/盘后数据
	TimeZone   string    // 可选：时区，如 "America/New_York"
}

// TimeSeriesValue 代表单条 OHLCV 数据
type TimeSeriesValue struct {
	DateTime DateOnly `json:"datetime"`
	Open     string   `json:"open"`
	High     string   `json:"high"`
	Low      string   `json:"low"`
	Close    string   `json:"close"`
	Volume   string   `json:"volume"`
}

// TimeSeriesMeta 代表响应中的元数据信息
type TimeSeriesMeta struct {
	Symbol   string `json:"symbol"`
	Interval string `json:"interval"`
	Currency string `json:"currency"`
}

// TimeSeriesResponse 是 /time_series 接口的完整响应结构
type TimeSeriesResponse struct {
	Meta   TimeSeriesMeta    `json:"meta"`
	Values []TimeSeriesValue `json:"values"`
	Status string            `json:"status"`
}