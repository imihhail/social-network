package urlHandlers

import (
	"backend/validators"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == "http://localhost:3000"
	},
}

var clients = make(map[*Client]bool)
var broadcast = make(chan SocketMessage)

type SocketMessage struct {
	Type   string `json:"type"`
	Status string `json:"status"`
	// From             string   `json:"fromuser"`
	// FromId           string   `json:"fromuserid"`
	Message          string   `json:"message"`
	To               string   `json:"touser"`
	ConnectedClients []string `json:"connectedclients"`
}

type Client struct {
	connection  *websocket.Conn
	send        chan []byte
	connOwnerId string
	mu          sync.Mutex
	lastActive  time.Time
}

func periodicUserPresenceCheck() {
	for {
		time.Sleep(time.Minute)
		// Iterate through clients and update their online status based on lastActive
		currentTimestamp := time.Now()
		for client := range clients {
			client.mu.Lock()
			if currentTimestamp.Sub(client.lastActive) > 3*time.Minute {
				client.connection.Close()
				delete(clients, client)
			}
			client.mu.Unlock()
		}
	}
}

func handleMessages(userId string) {
	for {
		msg := <-broadcast
		switch msg.Type {
		case "message":
			// set new message into db
			validators.ValidateSetNewMessage(userId, msg.Message, msg.To)
			// send to user directly
			for client := range clients {
				if msg.To == client.connOwnerId {
					client.mu.Lock()
					err := client.connection.WriteJSON(msg)
					if err != nil {
						fmt.Println(err)
						client.connection.Close()
						delete(clients, client)
					}
					client.mu.Unlock()
				}
			}
		case "onlineStatus":
			users := []string{}
			for key := range clients {
				users = append(users, key.connOwnerId)
			}
			allUsers := SocketMessage{
				Type:             "onlineStatus",
				Status:           "online",
				ConnectedClients: users,
			}
			// broadcast everyone that you are online/offline
			for client := range clients {
				client.mu.Lock()
				client.connection.WriteJSON(allUsers)
				client.mu.Unlock()
			}
		}
	}
}

func HandleSocket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Websocket attempt!")

	cookie, err := r.Cookie("socialNetworkSession")
	userId := validators.ValidateUserSession(cookie.Value)
	if userId == "0" {
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &Client{
		connection:  conn,
		connOwnerId: userId,
		send:        make(chan []byte, 256),
	}

	clients[client] = true

	defer func() {
		client.connection.Close()
	}()

	go handleMessages(userId)
	go periodicUserPresenceCheck()

	for {
		var msg SocketMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			client.mu.Lock()
			delete(clients, client)
			client.mu.Unlock()
			return
		}
		client.mu.Lock()
		client.lastActive = time.Now()
		client.mu.Unlock()
		broadcast <- msg
	}
}
