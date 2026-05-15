package services

import (
	"twelve_data_client/internal/constant"
	"twelve_data_client/internal/logger"
	"twelve_data_client/internal/model"

	"github.com/gorilla/websocket"
)

func GetTwelveDataWebSocket(apiKey string, subscription model.Subscription) (*websocket.Conn, error) {
	url := constant.TWELVED_DATA_WEBSOCKET_URL + "?apikey=" + apiKey

	logger.Debug("正在连接到 WebSocket: %s", constant.TWELVED_DATA_WEBSOCKET_URL)

	connection, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		logger.LogError("WebSocket 连接", err)
		return nil, err
	}

	logger.LogSuccess("WebSocket 连接已建立")

	subscriptionMesage := model.Subscription{
		Action: model.ActionSubscribe,
		Params: model.SubscriptionParams{Symbols: subscription.Params.Symbols},
	}

	logger.Debug("准备发送订阅消息: %d 个符号", len(subscriptionMesage.Params.Symbols))

	if err := connection.WriteJSON(subscriptionMesage); err != nil {
		logger.LogError("发送订阅消息", err)
		connection.Close()
		return nil, err
	}

	logger.LogSuccess("订阅消息已发送")

	return connection, nil
}