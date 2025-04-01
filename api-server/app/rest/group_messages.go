package rest

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"ride_sharing_api/app/assert"
	"ride_sharing_api/app/sqlc"
	"ride_sharing_api/app/utils"
	"sync"

	websocket "golang.org/x/net/websocket"
)

func groupMessageHandlers(h *http.ServeMux) {
	chat := groupChatHandler{
		chats: make(map[string]*groupChatWs),
	}

	wsChatServer := websocket.Server{
		Handler: chat.handleWs,
		Config: websocket.Config{
			Origin: nil,
		},
	}

	h.HandleFunc("POST /groups/by-id/{id}/send-message", handle(chat.groupMessageCreate).with(bearerAuth(false)).build())
	h.Handle("/groups/messages/{groupId}", wsChatServer)
}

type groupChatHandler struct {
	mutex sync.Mutex
	chats map[string]*groupChatWs
}

type GroupMessageData struct {
	MessageId   string  `json:"messageId"`
	GroupId     string  `json:"groupId"`
	Content     string  `json:"content"`
	SentBy      string  `json:"sentBy"`
	SentByEmail string  `json:"sentByEmail"`
	RepliesTo   *string `json:"repliesTo"`
	CreatedAt   string  `json:"createdAt"`
}

type createGroupMessageData struct {
	GroupId   *string `json:"groupId" validate:"required"`
	Content   *string `json:"content" validate:"required"`
	RepliesTo *string `json:"repliesTo"`
}

type groupChatWs struct {
	clients []*groupChatClient
}

type groupChatClient struct {
	conn *websocket.Conn
	done chan struct{}
}

func (c *groupChatHandler) groupMessageCreate(w http.ResponseWriter, r *http.Request) {
	user := getMiddlewareData[sqlc.User](r, "user")

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error: Invalid request body.", "error:", err)
		httpWriteErr(w, http.StatusBadRequest, "Invalid request body.")
		return
	}

	var createParams createGroupMessageData
	err = json.Unmarshal(data, &createParams)
	if err != nil {
		log.Println("Error: Invalid JSON in request body.", "error:", err)
		httpWriteErr(w, http.StatusBadRequest, "Invalid JSON in request body.", err.Error())
		return
	}

	err = utils.Validate.Struct(createParams)
	if err != nil {
		log.Println("Error: Invalid JSON in request body.", "error:", err)
		httpWriteErr(w, http.StatusBadRequest, "Missing/Invalid fields in request body.", err.Error())
		return
	}

	argsCreateGroup := sqlc.GroupMessagesCreateParams{
		Content:   *createParams.Content,
		GroupID:   *createParams.GroupId,
		SentBy:    user.ID,
		RepliesTo: utils.SqlNullStr(createParams.RepliesTo),
	}

	msg, err := state.queries.GroupMessagesCreate(r.Context(), argsCreateGroup)
	assert.Nil(err)

	var repliesTo *string = nil
	if msg.RepliesTo.Valid {
		repliesTo = &msg.RepliesTo.String
	}

	msgData := GroupMessageData{
		MessageId:   msg.ID,
		GroupId:     msg.GroupID,
		Content:     msg.Content,
		SentBy:      user.ID,
		SentByEmail: user.Email,
		RepliesTo:   repliesTo,
		CreatedAt:   msg.CreatedAt,
	}

	resp, err := json.Marshal(msgData)

	chat, exists := c.chats[msg.GroupID]
	if exists {
		for _, client := range chat.clients {
			if client.conn == nil {
				continue
			}

			client.Write(resp)
		}
	}

	assert.Nil(err, "Failed to serialize message.")
	w.WriteHeader(201)
	w.Write(resp)
}

func (c *groupChatHandler) handleWs(conn *websocket.Conn) {
	r := conn.Request()
	if r == nil {
		conn.Write([]byte("Missing required HTTP request."))
		conn.WriteClose(http.StatusBadRequest)
		return
	}

	groupId := r.PathValue("groupId")
	if groupId == "" {
		conn.Write([]byte("Missing 'groupId' parameter in route."))
		conn.WriteClose(http.StatusBadRequest)
		return
	}

	locked := true
	c.mutex.Lock()
	defer func() {
		if locked {
			c.mutex.Unlock()
		}
	}()

	done := make(chan struct{})
	client := groupChatClient{
		conn: conn,
		done: done,
	}

	chat, exists := c.chats[groupId]
	if exists {
		chat.clients = append(chat.clients, &client)
	} else {
		chat = &groupChatWs{
			clients: []*groupChatClient{&client},
		}

		c.chats[groupId] = chat
	}

	msgs, err := state.queries.GroupMessagesGetMany(r.Context(), groupId)
	assert.Nil(err)

	for _, msg := range msgs {
		var repliesTo *string = nil
		if msg.RepliesTo.Valid {
			repliesTo = &msg.RepliesTo.String
		}
		msgData := GroupMessageData{
			MessageId:   msg.ID,
			GroupId:     msg.GroupID,
			Content:     msg.Content,
			SentBy:      msg.SentBy,
			SentByEmail: msg.SentByEmail,
			RepliesTo:   repliesTo,
			CreatedAt:   msg.CreatedAt,
		}

		resp, err := json.Marshal(msgData)
		assert.Nil(err)

		_, err = client.Write(resp)
		assert.Nil(err)
	}

	locked = false
	c.mutex.Unlock()
	<-done // wait for websocket to close
	client.conn.WriteClose(http.StatusOK)
	client.conn = nil
}

func (ws *groupChatClient) Write(p []byte) (int, error) {
	if ws.conn == nil {
		return 0, nil
	}

	n, err := ws.conn.Write(p)
	if err != nil {
		return n, err
	}

	if n != len(p) {
		return n, io.ErrShortWrite
	}

	return n, nil
}
