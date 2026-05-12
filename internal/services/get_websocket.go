package services

import (
	"log"
	"twelve_data_client/internal/color"
	"twelve_data_client/internal/constant"
	"twelve_data_client/internal/model"

	"github.com/gorilla/websocket"
)

func GetTwelveDataWebSocket(apiKey string, subscription model.Subscription) (*websocket.Conn, error) {
	url := constant.TWELVED_DATA_WEBSOCKET_URL + "?apikey=" + apiKey

	connection, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	subscriptionMesage := model.Subscription{
		Action: model.ActionSubscribe,
		Params: model.SubscriptionParams{Symbols: subscription.Params.Symbols},
	}
	log.Printf("订阅消息: %s\n", color.Greenf("%+v", subscriptionMesage))

	if err := connection.WriteJSON(subscriptionMesage); err != nil {
		connection.Close()
		return nil, err
	}

	return connection, nil
}