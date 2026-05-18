package chat

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"backend/config"
	"backend/model"
	"backend/repository"

	"github.com/google/uuid"
)

type ChatService interface {
	CreateNewSession(ctx context.Context, title string) (*model.ChatSession, error)
	GetAllSessions(ctx context.Context) ([]model.ChatSession, error)
	GetSessionHistory(ctx context.Context, sessionID uuid.UUID) ([]model.ChatMessage, error)
	StreamChatMessage(ctx context.Context, sessionID uuid.UUID, userContent string, outChan chan<- string) error
}

type chatService struct {
	cfg        *config.Config
	pgRepo     repository.PostgresRepository
	redisRepo  repository.RedisRepository
	httpClient *http.Client
}

func NewChatService(cfg *config.Config, pgRepo repository.PostgresRepository, redisRepo repository.RedisRepository) ChatService {
	return &chatService{
		cfg:        cfg,
		pgRepo:     pgRepo,
		redisRepo:  redisRepo,
		httpClient: &http.Client{Timeout: 5 * time.Minute}, // Timeout dài để giữ kết nối stream với Python
	}
}

func (s *chatService) CreateNewSession(ctx context.Context, title string) (*model.ChatSession, error) {
	session := &model.ChatSession{
		ID:        uuid.New(),
		Title:     title,
		CreatedAt: time.Now(),
	}

	if err := s.pgRepo.CreateSession(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *chatService) GetAllSessions(ctx context.Context) ([]model.ChatSession, error) {
	return s.pgRepo.GetSessions(ctx)
}

func (s *chatService) GetSessionHistory(ctx context.Context, sessionID uuid.UUID) ([]model.ChatMessage, error) {

	messages, err := s.redisRepo.GetMessages(ctx, sessionID)
	if err == nil && len(messages) > 0 {
		return messages, nil
	}

	messages, err = s.pgRepo.GetMessagesBySessionID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (s *chatService) StreamChatMessage(ctx context.Context, sessionID uuid.UUID, userContent string, outChan chan<- string) error {
	defer close(outChan)

	userMsg := &model.ChatMessage{
		SessionID: sessionID,
		Role:      "user",
		Content:   userContent,
		CreatedAt: time.Now(),
	}
	if err := s.pgRepo.CreateMessage(ctx, userMsg); err != nil {
		return fmt.Errorf("failed to save user message: %w", err)
	}

	ttl := time.Duration(s.cfg.CacheTTLSeconds) * time.Second
	_ = s.redisRepo.SaveMessage(ctx, sessionID, userMsg, ttl)

	history, err := s.GetSessionHistory(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to retrieve chat history: %w", err)
	}

	pythonPayload := map[string]interface{}{
		"session_id": sessionID.String(),
		"messages":   history,
	}
	jsonData, err := json.Marshal(pythonPayload)
	if err != nil {
		return err
	}

	pythonURL := fmt.Sprintf("%s/api/v1/chat/stream", s.cfg.PythonLLMServiceURL)
	req, err := http.NewRequestWithContext(ctx, "POST", pythonURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("python service unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("python service returned error status: %d", resp.StatusCode)
	}

	var fullAssistantResponse bytes.Buffer
	reader := bufio.NewReader(resp.Body)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if len(line) > 5 && line[:5] == "data:" {
			chunk := line[5:]
			outChan <- chunk
			fullAssistantResponse.WriteString(chunk)
		}
	}

	assistantMsg := &model.ChatMessage{
		SessionID: sessionID,
		Role:      "assistant",
		Content:   fullAssistantResponse.String(),
		CreatedAt: time.Now(),
	}
	if err := s.pgRepo.CreateMessage(ctx, assistantMsg); err != nil {
		return fmt.Errorf("failed to save assistant message to postgres: %w", err)
	}

	_ = s.redisRepo.SaveMessage(ctx, sessionID, assistantMsg, ttl)

	return nil
}
