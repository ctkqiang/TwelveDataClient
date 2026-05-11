package services

import (
	"twelve_data_client/internal/constant"

	"github.com/gorilla/websocket"
)

func GetTwelveDataWebSocket() (*websocket.Conn, error) {
	url := constant.TWELVED_DATA_WEBSOCKET_URL  + ""

	connection, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	
	return connection, nil
}