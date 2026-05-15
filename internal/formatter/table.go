// Package formatter 提供命令行表格格式化输出功能，将结构化数据转换为易读的表格显示。
package formatter

import (
	"fmt"
	"strconv"
	"strings"
	"twelve_data_client/internal/color"
	"twelve_data_client/internal/model"
)

// StocksTable 将股票列表格式化为对齐的表格输出到控制台。
// 包含代码、名称、交易所、货币和类型列，支持颜色高亮显示。
func StocksTable(stocks []model.Stock) string {
	if len(stocks) == 0 {
		return "No stocks found"
	}

	// Header
	header := fmt.Sprintf("  %-14s  %-30s  %-20s  %-10s  %-15s",
		color.Bold(color.Blue("代码")),
		color.Bold(color.Blue("名称")),
		color.Bold(color.Blue("交易所")),
		color.Bold(color.Blue("货币")),
		color.Bold(color.Blue("类型")),
	)

	fmt.Println(header)
	fmt.Println("  " + color.Blue(strings.Repeat("─", 92)))

	// Rows
	for _, stock := range stocks {
		row := fmt.Sprintf("  %-14s  %-30s  %-20s  %-10s  %-15s",
			color.Cyan(stock.Symbol),
			truncateString(stock.Name, 28),
			color.Dim(stock.Exchange),
			color.Yellow(stock.Currency),
			stock.Type,
		)
		fmt.Println(row)
	}

	fmt.Println("  " + color.Blue(strings.Repeat("─", 92)))
	return fmt.Sprintf("共 %d 只股票", len(stocks))
}

// TimeSeriesTable 将历史行情数据（OHLCV）格式化为表格并输出到控制台。
// 显示日期、开盘价、最高价、最低价、收盘价和成交量，使用颜色标记涨跌趋势。
func TimeSeriesTable(symbol string, name string, currency string, values []model.TimeSeriesValue) {
	if len(values) == 0 {
		fmt.Printf("No data available for %s\n", symbol)
		return
	}

	// Print header
	fmt.Printf("\n  %s %s (%s)\n",
		color.Bold(color.Cyan(symbol)),
		color.Bold(color.White(truncateString(name, 30))),
		color.Bold(color.Yellow(currency)))

	// Table header
	header := fmt.Sprintf("  %-12s  %-10s  %-10s  %-10s  %-10s  %-12s",
		color.Bold("日期"),
		color.Bold("开盘"),
		color.Bold("最高"),
		color.Bold("最低"),
		color.Bold("收盘"),
		color.Bold("成交量"),
	)
	fmt.Println(header)
	fmt.Println("  " + color.Dim(strings.Repeat("─", 76)))

	// Data rows
	for _, v := range values {
		row := fmt.Sprintf("  %-12s  %10s  %10s  %10s  %10s  %12s",
			color.Dim(v.DateTime.Format("2006-01-02")),
			formatPrice(v.Open),
			color.Green(formatPrice(v.High)),
			color.Red(formatPrice(v.Low)),
			formatPrice(v.Close),
			formatVolume(v.Volume),
		)
		fmt.Println(row)
	}
}

// SummaryTable 将时间序列数据的统计信息格式化为表格并输出到控制台。
// 包括数据点数、平均价格、最高价格、最低价格、总成交量和平均成交量。
func SummaryTable(symbol string, values []model.TimeSeriesValue) {
	if len(values) == 0 {
		return
	}

	fmt.Printf("\n  %s\n", color.Bold(color.Cyan("统计信息")))

	// Table header
	header := fmt.Sprintf("  %-20s  %-20s",
		color.Bold("指标"),
		color.Bold("数值"),
	)
	fmt.Println(header)
	fmt.Println("  " + color.Dim(strings.Repeat("─", 42)))

	// Calculate statistics
	stats := calculateStats(values)

	rows := []struct {
		label string
		value string
	}{
		{"数据点数", color.Cyan(strconv.Itoa(len(values)))},
		{"平均价格", color.Yellow(fmt.Sprintf("%.5f", stats.AvgPrice))},
		{"最高价格", color.Green(fmt.Sprintf("%.5f", stats.MaxPrice))},
		{"最低价格", color.Red(fmt.Sprintf("%.5f", stats.MinPrice))},
		{"总成交量", color.White(formatVolumeInt(stats.TotalVolume))},
		{"平均成交量", color.White(formatVolumeInt(int64(stats.AvgVolume)))},
	}

	for _, r := range rows {
		row := fmt.Sprintf("  %-20s  %-20s", r.label, r.value)
		fmt.Println(row)
	}
}

// Stats holds calculated statistics
type Stats struct {
	AvgPrice    float64
	MaxPrice    float64
	MinPrice    float64
	TotalVolume int64
	AvgVolume   float64
}

// calculateStats calculates statistics from time series data
func calculateStats(values []model.TimeSeriesValue) Stats {
	if len(values) == 0 {
		return Stats{}
	}

	var (
		totalPrice  float64
		maxPrice    float64
		minPrice    = 999999.0
		totalVolume int64
	)

	for _, v := range values {
		price := parseFloat(v.Close)
		totalPrice += price

		if price > maxPrice {
			maxPrice = price
		}
		if price < minPrice {
			minPrice = price
		}

		volume := parseInt(v.Volume)
		totalVolume += volume
	}

	return Stats{
		AvgPrice:    totalPrice / float64(len(values)),
		MaxPrice:    maxPrice,
		MinPrice:    minPrice,
		TotalVolume: totalVolume,
		AvgVolume:   float64(totalVolume) / float64(len(values)),
	}
}

// ErrorTable displays an error for a symbol
func ErrorTable(symbol string, err error) {
	fmt.Printf("  %s %s - %v\n",
		color.Red("✗ FAILED"),
		color.Cyan(symbol),
		color.Red(err.Error()))
}

// Helper functions

func truncateString(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}

func formatPrice(price string) string {
	value := parseFloat(price)
	return fmt.Sprintf("%.5f", value)
}

func formatVolume(volume string) string {
	val := parseInt(volume)
	return formatVolumeInt(val)
}

func formatVolumeInt(volume int64) string {
	if volume >= 1_000_000 {
		return fmt.Sprintf("%.1fM", float64(volume)/1_000_000)
	}
	if volume >= 1_000 {
		return fmt.Sprintf("%.1fK", float64(volume)/1_000)
	}
	return fmt.Sprintf("%d", volume)
}

func parseFloat(s string) float64 {
	val, _ := strconv.ParseFloat(s, 64)
	return val
}

func parseInt(s string) int64 {
	val, _ := strconv.ParseInt(s, 10, 64)
	return val
}

// HeaderSection prints a formatted section header
func HeaderSection(title string) {
	fmt.Printf("\n%s\n", color.Bold(color.Green(title)))
	fmt.Println(color.Blue(strings.Repeat("═", 80)))
}

// SubHeader prints a formatted sub-header
func SubHeader(title string) {
	fmt.Printf("\n%s\n", color.Bold(color.Cyan(title)))
}
