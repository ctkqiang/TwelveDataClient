package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
	"twelve_data_client/internal/constant"
	"twelve_data_client/internal/model"
)

const defaultTimeout = 30 * time.Second

type apiResponse struct {
	Data   []model.Stock `json:"data"`
	Count  int           `json:"count"`
	Status string        `json:"status"`
}

// GetAllStocks 从 Twelve Data API 获取股票列表。
//
// exchange 为空时返回全部标的；非空时作为 query 参数传递给 API 端进行服务端过滤，
// 同时也会在客户端做二次校验。
//
// apiKey 为空时自动降级为 "demo"。
func GetAllStocks(exchange string, apiKey string) ([]model.Stock, error) {
	var r apiResponse

	if apiKey == "" {
		apiKey = "demo"
	}

	u, err := url.Parse(constant.TWELVED_DATA_API_URL + "/stocks")
	if err != nil {
		return nil, fmt.Errorf("解析 API URL 失败: %w", err)
	}

	q := u.Query()
	q.Set("apikey", apiKey)
	if exchange != "" {
		q.Set("exchange", exchange)
	}
	u.RawQuery = q.Encode()

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("构建请求失败: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求 API 失败: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 返回 HTTP %d: %s", resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, &r); err != nil {
		return nil, fmt.Errorf("解码 JSON 失败: %w", err)
	}

	if r.Status != "ok" {
		return nil, fmt.Errorf("API 状态异常: %s", r.Status)
	}

	return r.Data, nil
}

func GetTimeSeries(apiKey string, params *model.TimeSeriesParams) (*model.TimeSeriesResponse, error) {
	var tsResp model.TimeSeriesResponse
	
	if apiKey == "" {
		return nil, fmt.Errorf("API密钥不能为空")
	}

	q := url.Values{}
	q.Add("symbol", params.Symbol)
	q.Add("interval", params.Interval)
	q.Add("apikey", apiKey)

	if !params.StartDate.IsZero() && !params.EndDate.IsZero() {
		q.Add("start_date", params.StartDate.Format("2006-01-02"))
		q.Add("end_date", params.EndDate.Format("2006-01-02"))
	} else if params.OutputSize > 0 {
		q.Add("outputsize", fmt.Sprintf("%d", params.OutputSize))
	}

	if params.PrePost {
		q.Add("prepost", "true")
	}
	if params.TimeZone != "" {
		q.Add("timezone", params.TimeZone)
	}

	baseURL := constant.TWELVED_DATA_API_URL + "/time_series"
	fullURL := fmt.Sprintf("%s?%s", baseURL, q.Encode())

	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("请求时间序列数据失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 返回非 OK 状态: %s", resp.Status)
	}


	if err := json.NewDecoder(resp.Body).Decode(&tsResp); err != nil {
		return nil, fmt.Errorf("解码响应失败: %w", err)
	}

	if tsResp.Status != "ok" {
		return nil, fmt.Errorf("API 返回错误状态: %s", tsResp.Status)
	}

	return &tsResp, nil
}