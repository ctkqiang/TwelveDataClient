package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"twelve_data_client/internal/constant"
	"twelve_data_client/internal/model"
)

type APIResponse struct {
	Data   []model.Stock `json:"data"`
	Count  int           `json:"count"`
	Status string        `json:"status"`
}

func GetAllStocks(exchange string, apiKey string) ([] model.Stock,  error) {
	var apiResp APIResponse
	
	if apiKey == "" {
		apiKey = "demo"
	}

	url := constant.TWELVED_DATA_API_URL + "/stocks?apiKey=" + apiKey
	if exchange != "" {
		url = fmt.Sprintf("%s?exchange=%s", url, exchange)
	}


	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("获取股票数据失败: %w", err)
	}
	
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API返回非OK状态: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("解码响应失败: %w", err)
	}

	if apiResp.Status != "ok" {
		return nil, fmt.Errorf("API返回状态: %s", apiResp.Status)
	}

	return apiResp.Data, nil
}

func GetTimeSeries() error {
	return nil
}