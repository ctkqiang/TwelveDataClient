package services

import (
	"log"
	"twelve_data_client/internal/constant"
	"twelve_data_client/internal/model"

	"github.com/gorilla/websocket"
)

func GetTwelveDataWebSocket(apiKey string, subscription model.Subscription) (*websocket.Conn, error) {
	url := constant.TWELVED_DATA_WEBSOCKET_URL + "?apiKey=" + apiKey

	connection, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	subscriptionMesage := model.Subscription {
		Action: model.ActionSubscribe,
		Params: model.SubscriptionParams{Symbols: subscription.Params.Symbols},
	}
	log.Println("订阅消息:", subscriptionMesage)

	return connection, nil
}