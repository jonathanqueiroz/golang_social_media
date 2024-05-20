package websocket

import (
	"log"
	"net/http"
	"project01/src/auth"
	"project01/src/models"

	"github.com/gorilla/websocket"
)

var userChannels = make(map[uint64]chan models.Notification) // user channels
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// HandleConnections handles websocket connections
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	userID, err := auth.ExtractUserID(r)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	if _, ok := userChannels[userID]; !ok {
		userChannels[userID] = make(chan models.Notification)
	}

	for notification := range userChannels[userID] {
		err := ws.WriteJSON(notification)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}
	}
}

// SendNotification sends a notification to a user
func SendNotification(userID uint64, notification models.Notification) {
	userChannels[userID] <- notification
}
