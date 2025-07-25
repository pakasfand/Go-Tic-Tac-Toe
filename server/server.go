package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	. "server/handlers"
	. "server/models"
	. "shared_types"

	"github.com/gorilla/websocket"
)

func main() {
	// Serve static files (client WASM build output)
	fs := http.FileServer(http.Dir("../client"))
	    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // Set Content-Type explicitly for .js files
        if strings.HasSuffix(r.URL.Path, ".js") {
            w.Header().Set("Content-Type", "application/javascript")
        }
        if strings.HasSuffix(r.URL.Path, ".wasm") {
            w.Header().Set("Content-Type", "application/wasm")
        }

        fs.ServeHTTP(w, r)
    })

	var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true // ⚠️ Allow all origins for dev; lock this down for production
	}}
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
