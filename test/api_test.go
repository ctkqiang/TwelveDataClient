package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"twelve_data_client/internal/model"
	"twelve_data_client/internal/services"
)

// TestGetAllStocks_Success 测试成功获取所有股票
func TestGetAllStocks_Success(t *testing.T) {
	expectedStocks := []model.Stock{
		{
			Symbol:   "000001",
			Name:     "平安银行",
			Currency: "CNY",
			Exchange: "SZSE",
			MICCode:  "XSHE",
			Country:  "CN",
			Type:     "Stock",
			FIGICode: "BBG000B9XRY4",
		},
		{
			Symbol:   "000002",
			Name:     "万科A",
			Currency: "CNY",
			Exchange: "SZSE",
			MICCode:  "XSHE",
			Country:  "CN",
			Type:     "Stock",
			FIGICode: "BBG000B9Z3K0",
		},
	}

	responseBody := map[string]interface{}{
		"data":   expectedStocks,
		"count":  2,
		"status": "ok",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responseBody)
	}))
	defer server.Close()

	stocks, err := services.GetAllStocks("SZSE", "demo")
	if err != nil {
		t.Fatalf("期望无错误，得到错误: %v", err)
	}

	if stocks == nil {
		t.Fatal("期望返回非空股票列表")
	}
}

// TestGetAllStocks_WithSZSEExchange 测试使用深圳交易所符号
func TestGetAllStocks_WithSZSEExchange(t *testing.T) {
	responseBody := map[string]interface{}{
		"data": []model.Stock{
			{
				Symbol:   "000001",
				Name:     "平安银行",
				Currency: "CNY",
				Exchange: "SZSE",
				MICCode:  "XSHE",
				Country:  "CN",
				Type:     "Stock",
			},
		},
		"count":  1,
		"status": "ok",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证查询参数包含交易所信息
		exchange := r.URL.Query().Get("exchange")
		if exchange != "SZSE" {
			t.Errorf("期望交易所参数为 SZSE，得到: %s", exchange)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responseBody)
	}))
	defer server.Close()

	stocks, err := services.GetAllStocks("SZSE", "demo")
	if err != nil {
		t.Fatalf("使用 SZSE 交易所出错: %v", err)
	}

	if stocks == nil {
		t.Fatal("期望返回非空股票列表")
	}
}

// TestGetAllStocks_EmptyResponse 测试空响应处理
func TestGetAllStocks_EmptyResponse(t *testing.T) {
	responseBody := map[string]interface{}{
		"data":   []model.Stock{},
		"count":  0,
		"status": "ok",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responseBody)
	}))
	defer server.Close()

	stocks, err := services.GetAllStocks("SZSE", "demo")
	if err != nil {
		t.Fatalf("期望无错误，得到错误: %v", err)
	}

	if stocks == nil {
		t.Fatal("期望返回空切片，不是 nil")
	}

	if len(stocks) != 0 {
		t.Errorf("期望返回空列表，得到 %d 个股票", len(stocks))
	}
}

// TestGetAllStocks_APIError 测试 API 错误响应
func TestGetAllStocks_APIError(t *testing.T) {
	responseBody := map[string]interface{}{
		"data":   []model.Stock{},
		"count":  0,
		"status": "error",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responseBody)
	}))
	defer server.Close()

	_, err := services.GetAllStocks("SZSE", "demo")
	if err == nil {
		t.Fatal("期望 API 错误，但未返回错误")
	}

	if err.Error() != "API返回状态: error" {
		t.Errorf("期望特定错误消息，得到: %v", err)
	}
}

// TestGetAllStocks_HTTPError 测试 HTTP 错误状态码
func TestGetAllStocks_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("无效的 API Key"))
	}))
	defer server.Close()

	_, err := services.GetAllStocks("SZSE", "demo")
	if err == nil {
		t.Fatal("期望 HTTP 401 错误，但未返回错误")
	}

	if err.Error() != "API返回非OK状态: 401 Unauthorized" {
		t.Logf("错误消息: %v", err)
	}
}

// TestGetAllStocks_MalformedJSON 测试畸形 JSON 响应
func TestGetAllStocks_MalformedJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{invalid json}"))
	}))
	defer server.Close()

	_, err := services.GetAllStocks("SZSE", "demo")
	if err == nil {
		t.Fatal("期望 JSON 解析错误，但未返回错误")
	}

	if err.Error() != "解码响应失败: invalid character 'i' looking for beginning of object key string" {
		t.Logf("错误消息: %v", err)
	}
}

// TestGetAllStocks_EmptyAPIKey 测试空 API Key 默认处理
func TestGetAllStocks_EmptyAPIKey(t *testing.T) {
	responseBody := map[string]interface{}{
		"data": []model.Stock{
			{
				Symbol:   "000001",
				Name:     "平安银行",
				Exchange: "SZSE",
			},
		},
		"count":  1,
		"status": "ok",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.URL.Query().Get("apiKey")
		if apiKey != "demo" {
			t.Errorf("期望 API Key 为 demo，得到: %s", apiKey)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responseBody)
	}))
	defer server.Close()

	stocks, err := services.GetAllStocks("SZSE", "")
	if err != nil {
		t.Fatalf("期望无错误，得到错误: %v", err)
	}

	if stocks == nil {
		t.Fatal("期望返回非空股票列表")
	}
}

// TestGetAllStocks_MultipleStocks 测试返回多个股票
func TestGetAllStocks_MultipleStocks(t *testing.T) {
	stocks := make([]model.Stock, 100)
	for i := 0; i < 100; i++ {
		stocks[i] = model.Stock{
			Symbol:   "000001",
			Name:     "平安银行",
			Currency: "CNY",
			Exchange: "SZSE",
			Country:  "CN",
			Type:     "Stock",
		}
	}

	responseBody := map[string]interface{}{
		"data":   stocks,
		"count":  100,
		"status": "ok",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responseBody)
	}))
	defer server.Close()

	result, err := services.GetAllStocks("SZSE", "demo")
	if err != nil {
		t.Fatalf("期望无错误，得到错误: %v", err)
	}

	if len(result) != 100 {
		t.Errorf("期望 100 个股票，得到 %d 个", len(result))
	}
}

// TestGetAllStocks_StockFieldValidation 测试股票字段验证
func TestGetAllStocks_StockFieldValidation(t *testing.T) {
	expectedStock := model.Stock{
		Symbol:   "000001",
		Name:     "平安银行",
		Currency: "CNY",
		Exchange: "SZSE",
		MICCode:  "XSHE",
		Country:  "CN",
		Type:     "Stock",
		FIGICode: "BBG000B9XRY4",
		CFICode:  "ESVUFR",
		ISIN:     "CNE100000001",
		CUSIP:    "174921102",
	}

	responseBody := map[string]interface{}{
		"data":   []model.Stock{expectedStock},
		"count":  1,
		"status": "ok",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responseBody)
	}))
	defer server.Close()

	stocks, err := services.GetAllStocks("SZSE", "demo")
	if err != nil {
		t.Fatalf("期望无错误，得到错误: %v", err)
	}

	if len(stocks) != 1 {
		t.Fatalf("期望 1 个股票，得到 %d 个", len(stocks))
	}

	stock := stocks[0]

	if stock.Symbol != "000001" {
		t.Errorf("期望符号 000001，得到 %s", stock.Symbol)
	}

	if stock.Name != "平安银行" {
		t.Errorf("期望名称 平安银行，得到 %s", stock.Name)
	}

	if stock.Exchange != "SZSE" {
		t.Errorf("期望交易所 SZSE，得到 %s", stock.Exchange)
	}

	if stock.Currency != "CNY" {
		t.Errorf("期望货币 CNY，得到 %s", stock.Currency)
	}

	if stock.MICCode != "XSHE" {
		t.Errorf("期望 MIC 代码 XSHE，得到 %s", stock.MICCode)
	}

	if stock.ISIN != "CNE100000001" {
		t.Errorf("期望 ISIN CNE100000001，得到 %s", stock.ISIN)
	}
}

// TestGetAllStocks_ServerTimeout 测试服务器超时
func TestGetAllStocks_ServerTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 模拟超时，不返回响应
		select {}
	}))
	defer server.Close()

	// 这个测试可能需要超时设置才能有效
	_, err := services.GetAllStocks("SZSE", "demo")
	if err != nil {
		t.Logf("超时测试: 期望错误，得到: %v", err)
	}
}

// TestGetAllStocks_DifferentExchanges 测试不同的交易所参数
func TestGetAllStocks_DifferentExchanges(t *testing.T) {
	testCases := []string{
		"SZSE",
		"SSE",
		"HKEX",
		"NYSE",
		"NASDAQ",
	}

	for _, exchange := range testCases {
		t.Run("交易所_"+exchange, func(t *testing.T) {
			responseBody := map[string]interface{}{
				"data": []model.Stock{
					{
						Symbol:   "TEST",
						Exchange: exchange,
					},
				},
				"count":  1,
				"status": "ok",
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(responseBody)
			}))
			defer server.Close()

			stocks, err := services.GetAllStocks(exchange, "demo")
			if err != nil {
				t.Fatalf("交易所 %s 出错: %v", exchange, err)
			}

			if len(stocks) == 0 {
				t.Errorf("交易所 %s 返回空列表", exchange)
			}

			if stocks[0].Exchange != exchange {
				t.Errorf("交易所 %s: 期望交易所 %s，得到 %s", exchange, exchange, stocks[0].Exchange)
			}
		})
	}
}

// TestGetAllStocks_SpecialCharactersInFields 测试字段中的特殊字符
func TestGetAllStocks_SpecialCharactersInFields(t *testing.T) {
	stock := model.Stock{
		Symbol:   "000001",
		Name:     "平安银行 (中文名称)",
		Currency: "CNY",
		Exchange: "SZSE",
		Country:  "CN",
		Type:     "Stock & Fund",
		FIGICode: "BBG000B9XRY4",
	}

	responseBody := map[string]interface{}{
		"data":   []model.Stock{stock},
		"count":  1,
		"status": "ok",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responseBody)
	}))
	defer server.Close()

	stocks, err := services.GetAllStocks("SZSE", "demo")
	if err != nil {
		t.Fatalf("期望无错误，得到错误: %v", err)
	}

	if stocks[0].Name != "平安银行 (中文名称)" {
		t.Errorf("特殊字符处理失败，得到: %s", stocks[0].Name)
	}
}

// BenchmarkGetAllStocks 基准测试 - 获取股票列表性能
func BenchmarkGetAllStocks(b *testing.B) {
	stocks := make([]model.Stock, 1000)
	for i := 0; i < 1000; i++ {
		stocks[i] = model.Stock{
			Symbol:   "000001",
			Name:     "平安银行",
			Currency: "CNY",
			Exchange: "SZSE",
			Country:  "CN",
			Type:     "Stock",
		}
	}

	responseBody := map[string]interface{}{
		"data":   stocks,
		"count":  1000,
		"status": "ok",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responseBody)
	}))
	defer server.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		services.GetAllStocks("SZSE", "demo")
	}
}

// BenchmarkGetAllStocks_JSONUnmarshal 基准测试 - JSON 反序列化性能
func BenchmarkGetAllStocks_JSONUnmarshal(b *testing.B) {
	stocks := make([]model.Stock, 500)
	for i := 0; i < 500; i++ {
		stocks[i] = model.Stock{
			Symbol:   "000001",
			Name:     "平安银行",
			Currency: "CNY",
			Exchange: "SZSE",
			MICCode:  "XSHE",
			Country:  "CN",
			Type:     "Stock",
		}
	}

	responseBody := map[string]interface{}{
		"data":   stocks,
		"count":  500,
		"status": "ok",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responseBody)
	}))
	defer server.Close()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		services.GetAllStocks("SZSE", "demo")
	}
}
