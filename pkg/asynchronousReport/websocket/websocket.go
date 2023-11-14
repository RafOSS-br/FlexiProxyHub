package websocket

import (
	utils "FlexiProxyHub/internal/utils/generic"
	"FlexiProxyHub/pkg/database"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var DB *database.Database
var Log *zap.Logger

func SetLog(log *zap.Logger) {
	Log = log
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to upgrade connection:", err)
		return
	}
	filename, err := utils.GetCookie(r, "filename")
	if err != nil {
		fmt.Println("Failed to get cookie:", err)
	}
	// Read messages from the WebSocket connection
	database.DB = DB
	for {
		// Write a message back to the WebSocket connection

		if database.CheckIfTheHashHasCompleted(filename) {
			err = conn.WriteMessage(websocket.TextMessage, []byte("File downloaded!"))
			if err != nil {
				fmt.Println("Failed to write message:", err)
				break
			}
			break
		}
		time.Sleep(time.Second * 5)
	}
}
