package main

import (
	"websocket_basic/edition"
	"github.com/gin-gonic/gin"
)

func main() {
	server := edition.NewGeneralServerJsonMsg()

	go server.BroadcastMessageJsonMsg() // 啟動廣播訊息的 Goroutine

	r := gin.Default()
	r.GET("/ws", func(ctx *gin.Context) {
		server.HandleWSJsonMsg(ctx)
	})
	r.Run(":8080")
}
