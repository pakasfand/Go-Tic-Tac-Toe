package main

import (
	"encoding/json"
	"log"
	"net/http"
	
	. "server/models"
	. "server/handlers"
	. "shared_types"

	"github.com/gorilla/websocket"
)

func main() {
	var upgrader = websocket.Upgrader{}
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			http.Error(w, "Could not upgrade connection", http.StatusInternalServerError)
			return
		}
		defer c.Close()

		msgChan := make(chan []byte, 10)
		errChan := make(chan error, 10)

		go readMessages(c, msgChan, errChan)

		for {
			select {
			case message := <-msgChan:
				// handle message...
				var msg InboundMessage
				if err := json.Unmarshal(message, &msg); err != nil {
					log.Println("Could not unmarshal JSON:", err)
					continue
				}
				switch msg.Type {
				case MessageTypeCreateGame:
					HandleCreateGame(c)
					log.Println("Client initialized")
				case MessageTypeJoinGame:
					HandleJoinGame(c, msg.GameID)
				case MessageTypeMakeMove:
					if msg.Move != nil {
						HandleMakeMove(c, msg.GameID, *msg.Move)
					}
				case MessageTypeRequestPlayAgain:
					HandleRequestPlayAgain(c, msg.GameID)
				}
			case err := <-errChan:
				log.Println("Error:", err)
				return
			}
		}
	})

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func readMessages(c *websocket.Conn, msgChan chan []byte, errChan chan error) {
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			errChan <- err
			return
		}
		msgChan <- message
	}
}