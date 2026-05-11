package services

import (
	"twelve_data_client/internal/constant"

	"github.com/gorilla/websocket"
)

func GetTwelveDataWebSocket(apiKey string) (*websocket.Conn, error) {
	url := constant.TWELVED_DATA_WEBSOCKET_URL + "?apiKey=" + apiKey

	connection, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	
	return connection, nil
}