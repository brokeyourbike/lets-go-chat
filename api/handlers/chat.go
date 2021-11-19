package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Hub struct {
}

func NewHub() *Hub {
	return &Hub{}
}

func (hub *Hub) HandleChat() http.HandlerFunc {
	type query struct {
		Token string `json:"token"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Cannot upgrade request to websocket protocol", http.StatusInternalServerError)
			return
		}
		defer conn.Close()
		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				http.Error(w, "Cannot read message", http.StatusInternalServerError)
				break
			}
			log.Printf("recv: %s", message)
			err = conn.WriteMessage(mt, message)
			if err != nil {
				http.Error(w, "Cannot write message", http.StatusInternalServerError)
				break
			}
		}
	}
}
