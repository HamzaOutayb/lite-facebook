package handler

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
	"sync"

	"social-network/internal/models"
	utils "social-network/pkg"
	"social-network/pkg/middlewares"

	"github.com/gorilla/websocket"
)

var (
	userConnections = make(map[int][]*websocket.Conn)
	userConnMu      sync.RWMutex

	convSubscribers = make(map[int][]int)
	convSubMu       sync.RWMutex
)

// handle messages of ws
func (h *Handler) MessagesHandler(upgrader websocket.Upgrader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(middlewares.UserIDKey).(models.User)
		if !ok {
			utils.WriteJson(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("WebSocket upgrade error:", err)
			return
		}
		defer conn.Close()

		userConnMu.Lock()
		userConnections[user.ID] = append(userConnections[user.ID], conn)
		userConnMu.Unlock()

		defer removeConnection(user.ID, conn)

		conversations, err := h.Service.FetchConversations(user.ID)
		if err != nil {
			fmt.Println("Fetch conversations error:", err)
			return
		}

		addUserToConversations(user.ID, conversations)

		// Send initial data
		initialMsg := models.WSMessage{
			Type:          "conversations",
			Conversations: conversations,
			OnlineUsers:   getOnlineUsers(conversations),
		}
		if err := conn.WriteJSON(initialMsg); err != nil {
			fmt.Println("Initial message send error:", err)
			return
		}

		notifyUserStatus(user.ID, "online", conversations)

		for {
			var msg models.WSMessage
			if err := conn.ReadJSON(&msg); err != nil {
				if websocket.IsUnexpectedCloseError(err) {
					fmt.Printf("WebSocket closed: %v\n", err)
				}
				break
			}
			msg.Message.SenderID = user.ID
			handleMessage(msg, h)
		}

		notifyUserStatus(user.ID, "offline", conversations)
	}
}

// add user to conversations
func addUserToConversations(userID int, conversations []models.ConversationsInfo) {
	convSubMu.Lock()
	defer convSubMu.Unlock()

	for _, conv := range conversations {
		convID := conv.Conversation.ID
		subscribers := convSubscribers[convID]

		if !slices.Contains(subscribers, userID) {
			convSubscribers[convID] = append(subscribers, userID)
		}
	}
}

// Remove connection from user's connections
func removeConnection(userID int, conn *websocket.Conn) {
	userConnMu.Lock()
	defer userConnMu.Unlock()

	conns := userConnections[userID]
	for i, c := range conns {
		if c == conn {
			userConnections[userID] = slices.Delete(conns, i, i+1)
			break
		}
	}

	if len(userConnections[userID]) == 0 {
		delete(userConnections, userID)
		cleanupConversationSubscriptions(userID)
	}
}

// Remove user from conversations if last connection
func cleanupConversationSubscriptions(userID int) {
	convSubMu.Lock()
	defer convSubMu.Unlock()

	for convID, subscribers := range convSubscribers {
		if index := slices.Index(subscribers, userID); index != -1 {
			Tabupdated := slices.Delete(subscribers, index, index+1)
			if len(Tabupdated) == 0 {
				delete(convSubscribers, convID)
			} else {
				convSubscribers[convID] = Tabupdated
			}
		}
	}
}

// handle the message by type
func handleMessage(msg models.WSMessage, h *Handler) {
	switch msg.Type {
	case "new_message":
		msg.Message.Content = strings.TrimSpace(msg.Message.Content)
		if len(msg.Message.Content) == 0 || len(msg.Message.Content) > 500 {
			sendError(msg.Message.SenderID, "Invalid message content")
			return
		}

		convSubMu.RLock()
		subscribers, ok := convSubscribers[msg.Message.ConversationID]
		convSubMu.RUnlock()

		if !ok || !slices.Contains(subscribers, msg.Message.SenderID) {
			sendError(msg.Message.SenderID, "Not authorized for this conversation")
			return
		}

		if err := h.Service.CreateMessage(&msg.Message); err != nil {
			fmt.Println("Create message error:", err)
			sendError(msg.Message.SenderID, "Failed to send message")
			return
		}

		distributeMessage(msg, subscribers)
	}
}

// Notify status
func notifyUserStatus(userID int, status string, conversations []models.ConversationsInfo) {
	msg := models.WSMessage{
		Type:    status,
		Message: models.Message{SenderID: userID},
	}

	var allSubscribers []int
	convSubMu.RLock()
	for _, conv := range conversations {
		allSubscribers = append(allSubscribers, convSubscribers[conv.Conversation.ID]...)
	}
	convSubMu.RUnlock()

	subscribers := uniqueInts(allSubscribers)
	distributeMessage(msg, subscribers)
}

// destribute message for all
func distributeMessage(msg models.WSMessage, receivers []int) {
	userConnMu.RLock()
	defer userConnMu.RUnlock()

	for _, userID := range receivers {
		if conns, ok := userConnections[userID]; ok {
			for _, conn := range conns {
				if err := conn.WriteJSON(msg); err != nil {
					fmt.Println("Message distribution error:", err)
				}
			}
		}
	}
}

// get online users
func getOnlineUsers(conversations []models.ConversationsInfo) []int {
	userConnMu.RLock()
	defer userConnMu.RUnlock()

	var online []int
	for _, conv := range conversations {
		if conv.Conversation.Type == "private" {
			if _, ok := userConnections[conv.UserInfo.ID]; ok {
				online = append(online, conv.UserInfo.ID)
			}
		}
	}
	return online
}

// for errors
func sendError(userID int, message string) {
	userConnMu.RLock()
	defer userConnMu.RUnlock()

	if conns, ok := userConnections[userID]; ok {
		errMsg := models.WSMessage{
			Type:    "error",
			Message: models.Message{Content: message},
		}
		for _, conn := range conns {
			conn.WriteJSON(errMsg)
		}
	}
}

// Deduplicate subscribers
func uniqueInts(slice []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
