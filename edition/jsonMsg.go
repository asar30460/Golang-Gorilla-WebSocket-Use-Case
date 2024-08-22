package edition

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"encoding/json"
	"log"
)

type Message struct {
	ClientID int    `json:"client_id"`
	Message  string `json:"message"`
}

type ServerJsonMsg struct {
	clients   map[int]*ClientJsonMsg
	broadcast chan Message
}

type ClientJsonMsg struct {
	clientID int // 模擬使用者ID
	conn     *websocket.Conn
	receive  chan Message // 每個使用者只需要有接收通道，而發送是發給伺服器，讓伺服器決定這筆訊息要讓哪些使用者接收
}

func NewGeneralServerJsonMsg() *ServerJsonMsg {
	return &ServerJsonMsg{
		clients:   make(map[int]*ClientJsonMsg),
		broadcast: make(chan Message, 64),
	}
}

func (sjm *ServerJsonMsg) HandleWSJsonMsg(ctx *gin.Context) {
	client_id++
	log.Printf("A new client %d is trying to connect...", client_id)

	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Fatalln("Error upgrading connection:", err)
		return
	}
	defer ws.Close()

	// 每個使用者都有屬於自己的傳送及接收通道
	client := &ClientJsonMsg{
		clientID: client_id,
		conn:     ws,
		receive:  make(chan Message, 64),
	}
	sjm.clients[client_id] = client

	// 啟動 Goroutine 來處理接收訊息
	go func() {
		for {
			msg := <-client.receive
			messageJSON, err := json.Marshal(msg)
			if err != nil {
				log.Println("Error marshaling message:", err)
				continue
			}

			// 寫入接收到的訊息至客戶端
			if err := ws.WriteMessage(websocket.TextMessage, messageJSON); err != nil {
				log.Println("Error writing message:", err)
				return
			}
			log.Printf("Message sent to client %d: %s", client.clientID, msg.Message)
		}
	}()

	// 處理從 WebSocket 接收到的訊息
	for {
		_, message, err := ws.ReadMessage() // 讀取使用者發送的訊息
		if err != nil {
			log.Printf("Error reading message from client %d: %v", client.clientID, err)
			break
		}
		log.Printf("Message received from client %d: %s", client.clientID, message)

		// 將訊息封裝成 Message 結構
		msg := Message{
			ClientID: client.clientID,
			Message:  string(message),
		}

		// 將訊息發送到廣播通道
		sjm.broadcast <- msg
	}
}

func (s *ServerJsonMsg) BroadcastMessageJsonMsg() {
	for {
		msg := <-s.broadcast
		log.Printf("Broadcasting message from client %d: %s", msg.ClientID, msg.Message)
		for _, client := range s.clients {
			client.receive <- msg
		}
	}
}