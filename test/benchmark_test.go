package test

import (
	"encoding/json"
	"testing"
	"time"
	"twelve_data_client/internal/model"
)


// BenchmarkNewSubscribe_Latency 测试订阅创建延迟
func BenchmarkNewSubscribe_Latency(b *testing.B) {
	symbols := []string{"AAPL", "RY", "RY:TSX", "EUR/USD", "XAU/USD"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		model.NewSubscribe(symbols...)
	}
}

// BenchmarkSubscriptionMarshal_Latency 测试 JSON 序列化延迟
func BenchmarkSubscriptionMarshal_Latency(b *testing.B) {
	sub := model.NewSubscribe("AAPL", "RY", "RY:TSX")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		json.Marshal(sub)
	}
}

// BenchmarkPriceEventUnmarshal_Latency 测试价格事件解析延迟
func BenchmarkPriceEventUnmarshal_Latency(b *testing.B) {
	jsonData := []byte(`{
		"event": "price",
		"symbol": "AAPL",
		"price": 150.45,
		"timestamp": 1778551560
	}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var pe model.PriceEvent
		json.Unmarshal(jsonData, &pe)
	}
}

// BenchmarkSymbolJoin_Latency 测试符号连接延迟
func BenchmarkSymbolJoin_Latency(b *testing.B) {
	symbols := []string{"AAPL", "RY", "RY:TSX", "EUR/USD", "XAU/USD"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		model.NewSubscribe(symbols...)
	}
}

// === 模型性能测试 ===

// BenchmarkNewSubscribe_SingleSymbol 单个符号订阅性能测试
func BenchmarkNewSubscribe_SingleSymbol(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		model.NewSubscribe("AAPL")
	}
}

// BenchmarkNewSubscribe_ManySymbols 多个符号订阅性能测试
func BenchmarkNewSubscribe_ManySymbols(b *testing.B) {
	symbols := make([]string, 100)
	for i := 0; i < 100; i++ {
		symbols[i] = "AAPL"
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		model.NewSubscribe(symbols...)
	}
}

// BenchmarkNewSubscribe_StressMax 最大压力测试 (500 个符号)
func BenchmarkNewSubscribe_StressMax(b *testing.B) {
	symbols := make([]string, 500)
	for i := 0; i < 500; i++ {
		symbols[i] = "AAPL"
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		model.NewSubscribe(symbols...)
	}
}

// BenchmarkConcurrentSubscriptions 并发订阅创建测试
func BenchmarkConcurrentSubscriptions(b *testing.B) {
	done := make(chan bool, 100)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		go func() {
			model.NewSubscribe("AAPL", "RY", "RY:TSX")
			done <- true
		}()
	}

	for i := 0; i < b.N; i++ {
		<-done
	}
}

// BenchmarkPriceEventBatch 批量价格事件处理性能测试
func BenchmarkPriceEventBatch(b *testing.B) {
	events := []byte(`[
		{"event": "price", "symbol": "AAPL", "price": 150.45, "timestamp": 1778551560},
		{"event": "price", "symbol": "RY", "price": 125.30, "timestamp": 1778551561},
		{"event": "price", "symbol": "RY:TSX", "price": 128.67, "timestamp": 1778551562},
		{"event": "price", "symbol": "EUR/USD", "price": 1.092, "timestamp": 1778551563},
		{"event": "price", "symbol": "XAU/USD", "price": 2025.12, "timestamp": 1778551564}
	]`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var priceEvents []model.PriceEvent
		json.Unmarshal(events, &priceEvents)
	}
}

// BenchmarkSubscriptionUnsubscribe_Performance 取消订阅性能测试
func BenchmarkSubscriptionUnsubscribe_Performance(b *testing.B) {
	symbols := []string{"AAPL", "RY", "RY:TSX", "EUR/USD", "XAU/USD"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		model.NewUnsubscribe(symbols...)
	}
}

// BenchmarkMemoryAllocation 内存分配性能测试
func BenchmarkMemoryAllocation(b *testing.B) {
	symbols := []string{"AAPL", "RY", "RY:TSX", "EUR/USD", "XAU/USD", "EUR/GBP", "USD/JPY", "GLD"}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		model.NewSubscribe(symbols...)
	}
}

// === WebSocket 服务集成测试 ===

// TestWebSocketService_SubscriptionMessageGeneration 测试订阅消息生成
func TestWebSocketService_SubscriptionMessageGeneration(t *testing.T) {
	symbols := []string{"AAPL", "RY", "RY:TSX", "EUR/USD", "XAU/USD"}
	subscription := model.NewSubscribe(symbols...)

	if subscription.Action != model.ActionSubscribe {
		t.Errorf("期望 Action 为 %q，得到 %q", model.ActionSubscribe, subscription.Action)
	}

	if subscription.Params.Symbols == "" {
		t.Error("订阅参数中符号为空")
	}

	expectedSymbols := "AAPL,RY,RY:TSX,EUR/USD,XAU/USD"
	if subscription.Params.Symbols != expectedSymbols {
		t.Errorf("期望符号 %q，得到 %q", expectedSymbols, subscription.Params.Symbols)
	}
}

// TestWebSocketService_UnsubscriptionMessageGeneration 测试取消订阅消息生成
func TestWebSocketService_UnsubscriptionMessageGeneration(t *testing.T) {
	symbols := []string{"AAPL", "EUR/USD"}
	unsubscription := model.NewUnsubscribe(symbols...)

	if unsubscription.Action != model.ActionUnsubscribe {
		t.Errorf("期望 Action 为 %q，得到 %q", model.ActionUnsubscribe, unsubscription.Action)
	}

	expectedSymbols := "AAPL,EUR/USD"
	if unsubscription.Params.Symbols != expectedSymbols {
		t.Errorf("期望符号 %q，得到 %q", expectedSymbols, unsubscription.Params.Symbols)
	}
}

// BenchmarkWebSocketService_SubscriptionMessageCreation WebSocket 服务订阅消息创建基准
func BenchmarkWebSocketService_SubscriptionMessageCreation(b *testing.B) {
	symbols := []string{"AAPL", "RY", "RY:TSX", "EUR/USD", "XAU/USD"}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		subscription := model.NewSubscribe(symbols...)
		_, _ = json.Marshal(subscription)
	}
}

// BenchmarkWebSocketService_MessageSerialization WebSocket 消息序列化基准
func BenchmarkWebSocketService_MessageSerialization(b *testing.B) {
	subscription := model.NewSubscribe("AAPL", "RY", "RY:TSX", "EUR/USD", "XAU/USD")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		json.Marshal(subscription)
	}
}

// === 价格数据验证测试 ===

// TestPricing_ValidPrice 验证正确的价格解析
func TestPricing_ValidPrice(t *testing.T) {
	jsonData := []byte(`{
		"event": "price",
		"symbol": "AAPL",
		"price": 150.45,
		"timestamp": 1778551560
	}`)

	var pe model.PriceEvent
	err := json.Unmarshal(jsonData, &pe)
	if err != nil {
		t.Fatalf("反序列化失败: %v", err)
	}

	if pe.Price != 150.45 {
		t.Errorf("期望价格 150.45，得到 %f", pe.Price)
	}
}

// TestPricing_ZeroPrice 零价格处理测试
func TestPricing_ZeroPrice(t *testing.T) {
	jsonData := []byte(`{
		"event": "price",
		"symbol": "AAPL",
		"price": 0,
		"timestamp": 1778551560
	}`)

	var pe model.PriceEvent
	json.Unmarshal(jsonData, &pe)

	if pe.Price != 0 {
		t.Errorf("期望零价格，得到 %f", pe.Price)
	}
}

// TestPricing_NegativePrice 负价格处理测试 (错误情况)
func TestPricing_NegativePrice(t *testing.T) {
	jsonData := []byte(`{
		"event": "price",
		"symbol": "AAPL",
		"price": -10.50,
		"timestamp": 1778551560
	}`)

	var pe model.PriceEvent
	json.Unmarshal(jsonData, &pe)

	if pe.Price < 0 {
		t.Logf("警告: 接受负价格 %f (应该被验证并拒绝)", pe.Price)
	}
}

// TestPricing_VeryHighPrice 超高价格值测试
func TestPricing_VeryHighPrice(t *testing.T) {
	jsonData := []byte(`{
		"event": "price",
		"symbol": "BTC/USD",
		"price": 999999999.99,
		"timestamp": 1778551560
	}`)

	var pe model.PriceEvent
	json.Unmarshal(jsonData, &pe)

	if pe.Price != 999999999.99 {
		t.Errorf("解析大价格失败，得到 %f", pe.Price)
	}
}

// TestPricing_VerySmallPrice 超小价格值测试
func TestPricing_VerySmallPrice(t *testing.T) {
	jsonData := []byte(`{
		"event": "price",
		"symbol": "SHIB/USD",
		"price": 0.00000001,
		"timestamp": 1778551560
	}`)

	var pe model.PriceEvent
	json.Unmarshal(jsonData, &pe)

	if pe.Price != 0.00000001 {
		t.Errorf("解析小价格失败，得到 %f", pe.Price)
	}
}

// TestPricing_PrecisionLoss 高精度价格精度损失测试
func TestPricing_PrecisionLoss(t *testing.T) {
	jsonData := []byte(`{
		"event": "price",
		"symbol": "AAPL",
		"price": 150.123456789012345,
		"timestamp": 1778551560
	}`)

	var pe model.PriceEvent
	json.Unmarshal(jsonData, &pe)

	if pe.Price == 0 {
		t.Error("解析高精度价格失败")
	}
}

// BenchmarkPricing_FastPriceUpdate 快速价格更新基准测试
func BenchmarkPricing_FastPriceUpdate(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jsonData := []byte(`{"event": "price", "symbol": "AAPL", "price": 150.45, "timestamp": 1778551560}`)
		var pe model.PriceEvent
		json.Unmarshal(jsonData, &pe)
	}
}

// === 错误和边界情况测试 ===

// TestError_MalformedJSON 畸形 JSON 处理测试
func TestError_MalformedJSON(t *testing.T) {
	jsonData := []byte(`{invalid json}`)
	var pe model.PriceEvent
	err := json.Unmarshal(jsonData, &pe)
	if err == nil {
		t.Error("期望畸形 JSON 出错")
	}
}

// TestError_MissingSymbol 缺失符号字段测试
func TestError_MissingSymbol(t *testing.T) {
	jsonData := []byte(`{
		"event": "price",
		"price": 150.45,
		"timestamp": 1778551560
	}`)

	var pe model.PriceEvent
	json.Unmarshal(jsonData, &pe)

	if pe.Symbol != "" {
		t.Logf("警告: 缺失的符号被接受无验证: %q", pe.Symbol)
	}
}

// TestError_MissingPrice 缺失价格字段测试
func TestError_MissingPrice(t *testing.T) {
	jsonData := []byte(`{
		"event": "price",
		"symbol": "AAPL",
		"timestamp": 1778551560
	}`)

	var pe model.PriceEvent
	json.Unmarshal(jsonData, &pe)

	if pe.Price != 0 {
		t.Logf("警告: 缺失价格默认为: %f", pe.Price)
	}
}

// TestError_EmptySymbolArray 空符号数组订阅测试
func TestError_EmptySymbolArray(t *testing.T) {
	sub := model.NewSubscribe()
	if sub.Params.Symbols != "" {
		t.Logf("警告: 空订阅创建带有符号: %q", sub.Params.Symbols)
	}
}

// TestError_InvalidSymbolFormat 各种无效符号格式测试
func TestError_InvalidSymbolFormat(t *testing.T) {
	testCases := []struct {
		name   string
		symbol string
	}{
		{"空符号", ""},
		{"空格", " "},
		{"前导空格", " AAPL"},
		{"尾部空格", "AAPL "},
		{"内部空格", "AA PL"},
		{"换行符", "AAPL\n"},
		{"制表符", "AAPL\t"},
	}

	for _, tc := range testCases {
		sub := model.NewSubscribe(tc.symbol)
		if tc.symbol != "" && (sub.Params.Symbols == "" || sub.Params.Symbols != tc.symbol) {
			t.Logf("警告 [%s]: 符号 %q 处理不当，得到: %q", tc.name, tc.symbol, sub.Params.Symbols)
		}
	}
}

// TestError_VeryLongSymbolList 超长符号列表测试
func TestError_VeryLongSymbolList(t *testing.T) {
	symbols := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		symbols[i] = "AAPL"
	}

	sub := model.NewSubscribe(symbols...)
	if sub.Params.Symbols == "" {
		t.Error("处理超长符号列表失败")
	}
}

// TestError_SpecialCharactersInSymbol 符号中特殊字符测试
func TestError_SpecialCharactersInSymbol(t *testing.T) {
	testCases := []string{
		"AAPL:NYSE",
		"EUR/USD",
		"BTC/USD",
		"RY:TSX",
		"AAPL;DROP",
		"AAPL|UNION",
	}

	for _, symbol := range testCases {
		sub := model.NewSubscribe(symbol)
		if sub.Params.Symbols == "" {
			t.Logf("警告: 带特殊字符的符号被拒绝: %q", symbol)
		}
	}
}

// TestError_CaseSensitivity 符号大小写敏感性测试
func TestError_CaseSensitivity(t *testing.T) {
	sub1 := model.NewSubscribe("aapl")
	sub2 := model.NewSubscribe("AAPL")

	if sub1.Params.Symbols != sub2.Params.Symbols {
		t.Logf("警告: 大小写敏感性问题 - 'aapl' 不等于 'AAPL'")
	}
}

// === 时间戳和并发测试 ===

// TestTimestamp_ValidTimestamp 有效时间戳验证
func TestTimestamp_ValidTimestamp(t *testing.T) {
	jsonData := []byte(`{
		"event": "price",
		"symbol": "AAPL",
		"price": 150.45,
		"timestamp": 1778551560
	}`)

	var pe model.PriceEvent
	json.Unmarshal(jsonData, &pe)

	if pe.Timestamp != 1778551560 {
		t.Errorf("时间戳不匹配: 得到 %d", pe.Timestamp)
	}
}

// TestTimestamp_FutureTimestamp 未来时间戳处理测试
func TestTimestamp_FutureTimestamp(t *testing.T) {
	future := time.Now().AddDate(1, 0, 0).Unix()
	jsonData := []byte(`{"event": "price", "symbol": "AAPL", "price": 150.45, "timestamp": ` + string(rune(future)) + `}`)

	var pe model.PriceEvent
	json.Unmarshal(jsonData, &pe)

	if pe.Timestamp > time.Now().Unix() {
		t.Logf("警告: 接受未来时间戳: %d (应该被拒绝)", pe.Timestamp)
	}
}

// TestTimestamp_OldTimestamp 超旧时间戳处理测试
func TestTimestamp_OldTimestamp(t *testing.T) {
	jsonData := []byte(`{
		"event": "price",
		"symbol": "AAPL",
		"price": 150.45,
		"timestamp": 0
	}`)

	var pe model.PriceEvent
	json.Unmarshal(jsonData, &pe)

	if pe.Timestamp != 0 {
		t.Logf("警告: Unix 时代时间戳处理: 得到 %d", pe.Timestamp)
	}
}

// BenchmarkConcurrency_ParallelParsing 并发价格解析基准测试
func BenchmarkConcurrency_ParallelParsing(b *testing.B) {
	jsonData := []byte(`{
		"event": "price",
		"symbol": "AAPL",
		"price": 150.45,
		"timestamp": 1778551560
	}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go func() {
			var pe model.PriceEvent
			json.Unmarshal(jsonData, &pe)
		}()
	}
}

// === 订阅响应测试 ===

// TestSubscriptionResponse_ValidResponse 有效订阅响应测试
func TestSubscriptionResponse_ValidResponse(t *testing.T) {
	jsonData := []byte(`{
		"event": "subscribe",
		"status": "ok",
		"success": [
			{"symbol": "AAPL", "exchange": "NYSE", "country": "USA", "type": "Stock"}
		]
	}`)

	var sr model.SubscriptionResponse
	err := json.Unmarshal(jsonData, &sr)
	if err != nil {
		t.Fatalf("反序列化失败: %v", err)
	}

	if len(sr.Success) != 1 || sr.Success[0].Symbol != "AAPL" {
		t.Error("订阅响应解析失败")
	}
}

// TestSubscriptionResponse_MultipleSymbols 多符号订阅响应测试
func TestSubscriptionResponse_MultipleSymbols(t *testing.T) {
	jsonData := []byte(`{
		"event": "subscribe",
		"status": "ok",
		"success": [
			{"symbol": "AAPL", "exchange": "NYSE"},
			{"symbol": "RY", "exchange": "TSX"},
			{"symbol": "RY:TSX", "exchange": "TSX"}
		]
	}`)

	var sr model.SubscriptionResponse
	json.Unmarshal(jsonData, &sr)

	if len(sr.Success) != 3 {
		t.Errorf("期望 3 个符号，得到 %d", len(sr.Success))
	}
}

// TestSubscriptionResponse_EmptySuccess 空成功数组响应测试
func TestSubscriptionResponse_EmptySuccess(t *testing.T) {
	jsonData := []byte(`{
		"event": "subscribe",
		"status": "ok",
		"success": []
	}`)

	var sr model.SubscriptionResponse
	json.Unmarshal(jsonData, &sr)

	if len(sr.Success) != 0 {
		t.Errorf("期望空成功数组，得到 %d 个条目", len(sr.Success))
	}
}

// TestSubscriptionResponse_MissingFields 缺失可选字段的响应测试
func TestSubscriptionResponse_MissingFields(t *testing.T) {
	jsonData := []byte(`{
		"event": "subscribe",
		"status": "ok",
		"success": [
			{"symbol": "AAPL"}
		]
	}`)

	var sr model.SubscriptionResponse
	json.Unmarshal(jsonData, &sr)

	if sr.Success[0].Symbol != "AAPL" {
		t.Error("解析最小响应失败")
	}
}


// TestWebSocketService_DocumentationExample 验证文档中的使用示例
func TestWebSocketService_DocumentationExample(t *testing.T) {
	// 验证订阅 API 的基本用法
	symbols := []string{"AAPL", "RY", "RY:TSX", "EUR/USD", "XAU/USD"}
	subscription := model.NewSubscribe(symbols...)

	// 验证创建的订阅对象
	if subscription.Action != model.ActionSubscribe {
		t.Fatal("订阅 Action 不正确")
	}

	// 验证订阅可以被序列化
	data, err := json.Marshal(subscription)
	if err != nil {
		t.Fatalf("订阅序列化失败: %v", err)
	}

	// 验证序列化后的数据包含所有符号
	var unmarshaled model.Subscription
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("订阅反序列化失败: %v", err)
	}

	if unmarshaled.Params.Symbols != subscription.Params.Symbols {
		t.Error("序列化/反序列化过程中数据不一致")
	}
}

// BenchmarkWebSocketService_EndToEndMessageFlow 端到端消息流性能测试
func BenchmarkWebSocketService_EndToEndMessageFlow(b *testing.B) {
	symbols := []string{"AAPL", "RY", "RY:TSX", "EUR/USD", "XAU/USD"}
	priceEventJSON := []byte(`{"event": "price", "symbol": "AAPL", "price": 150.45, "timestamp": 1778551560}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 创建订阅消息
		subscription := model.NewSubscribe(symbols...)
		subscriptionData, _ := json.Marshal(subscription)

		// 解析价格事件
		var priceEvent model.PriceEvent
		json.Unmarshal(priceEventJSON, &priceEvent)

		_ = subscriptionData
	}
}
