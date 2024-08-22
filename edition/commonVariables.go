package edition

import (
	"github.com/gorilla/websocket"
	"net/http"
)

var client_id int = -1

type Req struct {
	ClientID int `json:"client_id"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}