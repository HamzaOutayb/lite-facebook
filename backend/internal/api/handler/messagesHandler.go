package handler

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strings"
	"sync"

	"social-network/internal/models"
	utils "social-network/pkg"

	"github.com/gofrs/uuid/v5"
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
		user, ok := r.Context().Value(utils.UserIDKey).(models.UserInfo)
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

		notifyUserStatus(user.ID, "online", conversations, h)

		for {
			var msg models.WSMessage
			typeMessage, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err) {
					fmt.Printf("WebSocket closed: %v\n", err)
				}
				break
			}

			if typeMessage == websocket.BinaryMessage {
				fmt.Println(len(message))
				if len(message) < 4 {
					fmt.Println("short")
					continue
				}
				idLen := binary.LittleEndian.Uint32(message[0:4])
				fmt.Println(idLen)
				if len(message) < int(idLen)+4 {
					fmt.Println("short dfdf")
					continue
				}
				fmt.Println("idLen", idLen)
				err = json.Unmarshal(message[4:4+idLen], &msg)
				if err != nil {
					fmt.Printf("error n json: %v\n", err)
					break
				}
				fmt.Println(msg)
				path := HandleImage(msg.Type, message[4+idLen:])
				fmt.Println(path)
				msg.Type = "new_message"
				msg.Message.Image = path

			} else if typeMessage == websocket.TextMessage {
				err = json.Unmarshal(message, &msg)
				if err != nil {
					fmt.Printf("error n json: %v\n", err)
					break
				}
				fmt.Println(message)
				msg.Message.Image = ""
			}
			msg.Message.SenderID = user.ID
			handleMessage(msg, h, conn)
		}

		notifyUserStatus(user.ID, "offline", conversations, h)
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
func handleMessage(msg models.WSMessage, h *Handler, conn *websocket.Conn) {
	switch msg.Type {
	case "new_message":
		if msg.Message.Image == "" {
			msg.Message.Content = strings.TrimSpace(msg.Message.Content)
			if len(msg.Message.Content) == 0 || len(msg.Message.Content) > 500 {
				sendError(msg.Message.SenderID, "Invalid message content")
				return
			}
		}

		convSubMu.RLock()
		subscribers, ok := convSubscribers[msg.Message.ConversationID]
		convSubMu.RUnlock()

		if !ok || !slices.Contains(subscribers, msg.Message.SenderID) {
			sendError(msg.Message.SenderID, "Not authorized for this conversation")
			return
		}
		var err error
		if err = h.Service.CreateMessage(&msg); err != nil {
			fmt.Println("Create message error:", err)
			sendError(msg.Message.SenderID, "Failed to send message")
			return
		}

		if msg.UserInfo, err = h.Service.GetUserByID(msg.Message.SenderID); err != nil {
			fmt.Println("Get user error", err)
			sendError(msg.Message.SenderID, "Failed to send message")
			return
		}

		distributeMessage(msg, subscribers)
	case "conversations":

		conversations, err := h.Service.FetchConversations(msg.Message.SenderID)
		if err != nil {
			fmt.Println("Fetch conversations error:", err)
			return
		}

		initialMsg := models.WSMessage{
			Type:          "conversations",
			Conversations: conversations,
			OnlineUsers:   getOnlineUsers(conversations),
		}
		if err := conn.WriteJSON(initialMsg); err != nil {
			fmt.Println("Initial message send error:", err)
			return
		}
	}
}

// Notify status
func notifyUserStatus(userID int, status string, conversations []models.ConversationsInfo, h *Handler) {
	msg := models.WSMessage{
		Type:    status,
		Message: models.Message{SenderID: userID},
	}
	var err error
	if msg.UserInfo, err = h.Service.GetUserByID(msg.Message.SenderID); err != nil {
		fmt.Println("Get user error", err)
		sendError(msg.Message.SenderID, "Failed to send message")
		return
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

func (Handler *Handler) HandelMessagesHestories(w http.ResponseWriter, r *http.Request) {
	_, data, err := Handler.AfterGet(w, r)
	if err.Err != nil {
		return
	}
	messages, err := Handler.Service.FetchMessagesHestories(data.Before, data.ConversationID)
	if err.Err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	utils.WriteJson(w, http.StatusOK, messages)
}

func HandleImage(filename string, buffer []byte) string {
	extensions := []string{".png", ".jepg", ".gif", ".jpg"}
	extIndex := slices.IndexFunc(extensions, func(ext string) bool {
		return strings.HasSuffix(filename, ext)
	})
	if extIndex == -1 {
		return ""
	}
	imageName, _ := uuid.NewV4()
	err := os.WriteFile("../front-end/public/images/"+imageName.String()+extensions[extIndex], buffer, 0o644)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return imageName.String() + extensions[extIndex]
}
