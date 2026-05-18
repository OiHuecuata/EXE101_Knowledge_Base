package handler

import (
	"context"
	"io"
	"net/http"

	"backend/service/chat"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChatHandler struct {
	chatService chat.ChatService
}

func NewChatHandler(chatService chat.ChatService) *ChatHandler {
	return &ChatHandler{chatService: chatService}
}

// POST /chats
func (h *ChatHandler) CreateSession(c *gin.Context) {
	var req struct {
		Title string `json:"title" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}

	session, err := h.chatService.CreateNewSession(c.Request.Context(), req.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, session)
}

// GET /chats
func (h *ChatHandler) GetSessions(c *gin.Context) {
	sessions, err := h.chatService.GetAllSessions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

// GET /chats/:id/messages
func (h *ChatHandler) GetMessages(c *gin.Context) {
	idStr := c.Param("id")
	sessionID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID format"})
		return
	}

	messages, err := h.chatService.GetSessionHistory(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

// POST /chats/:id/messages (Stream SSE)
func (h *ChatHandler) StreamMessage(c *gin.Context) {
	idStr := c.Param("id")
	sessionID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID format"})
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content is required"})
		return
	}

	outChan := make(chan string, 10)
	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	go func() {
		if err := h.chatService.StreamChatMessage(ctx, sessionID, req.Content, outChan); err != nil {
			cancel()
		}
	}()

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	c.Stream(func(w io.Writer) bool {
		select {
		case <-ctx.Done():
			return false
		case chunk, ok := <-outChan:
			if !ok {
				return false
			}
			c.SSEvent("message", chunk)
			return true
		}
	})
}
