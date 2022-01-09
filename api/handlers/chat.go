package handlers

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/brokeyourbike/lets-go-chat/db"
	"github.com/brokeyourbike/lets-go-chat/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type MessagesRepo interface {
	Create(msg models.Message) (models.Message, error)
	GetAfterDateExcludingUserId(after time.Time, userId uuid.UUID) ([]models.Message, error)
}

type Chat struct {
	chatHub         *Hub
	activeUsersRepo ActiveUsersRepo
	tokensRepo      TokensRepo
	messagesRepo    MessagesRepo
}

func NewChat(h *Hub, a ActiveUsersRepo, t TokensRepo, m MessagesRepo) *Chat {
	return &Chat{chatHub: h, activeUsersRepo: a, tokensRepo: t, messagesRepo: m}
}

func (c *Chat) HandleChat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := uuid.Parse(r.URL.Query().Get("token"))
		if err != nil {
			http.Error(w, "Token format invalid", http.StatusBadRequest)
			return
		}

		token, err := c.tokensRepo.Get(t)
		if errors.Is(err, db.ErrTokenNotFound) {
			http.Error(w, "Token invalid", http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, "Token cannot be validated", http.StatusInternalServerError)
			return
		}

		if token.ExpiresAt.Before(time.Now()) {
			http.Error(w, "Token expired", http.StatusBadRequest)
			return
		}

		c.tokensRepo.InvalidateByUserId(token.UserID)
		c.activeUsersRepo.Add(token.UserID)
		defer c.activeUsersRepo.Delete(token.UserID)

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Cannot upgrade request to websocket protocol", http.StatusInternalServerError)
			return
		}

		client := &Client{hub: c.chatHub, conn: conn, send: make(chan []byte, 256)}
		client.hub.register <- client

		messages, err := c.messagesRepo.GetAfterDateExcludingUserId(time.Now(), token.UserID)
		if err != nil {
			http.Error(w, "Cannot fetch previous messages", http.StatusInternalServerError)
			return
		}

		for _, msg := range messages {
			client.send <- []byte(msg.Text)
		}

		go client.write()
		go client.read(c.messagesRepo)
	}
}

type Message struct {
	client  *Client
	content []byte
}

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				if client == message.client {
					continue
				}
				select {
				case client.send <- message.content:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	userID uuid.UUID
}

func (c *Client) read(messagesRepo MessagesRepo) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		messagesRepo.Create(models.Message{ID: uuid.New(), UserID: c.userID, Text: string(message), CreatedAt: time.Now()})
		c.hub.broadcast <- Message{client: c, content: message}
	}
}

func (c *Client) write() {
	defer func() {
		c.conn.Close()
	}()
	for {
		message, ok := <-c.send
		c.conn.SetWriteDeadline(time.Now().Add(writeWait))
		if !ok {
			// The hub closed the channel.
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		w.Write(message)

		if err := w.Close(); err != nil {
			return
		}
	}
}
