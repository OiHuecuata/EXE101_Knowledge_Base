package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"backend/model"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisRepository interface {
	SaveMessage(ctx context.Context, sessionID uuid.UUID, message *model.ChatMessage, ttl time.Duration) error
	GetMessages(ctx context.Context, sessionID uuid.UUID) ([]model.ChatMessage, error)
	ClearCache(ctx context.Context, sessionID uuid.UUID) error
}

type redisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) RedisRepository {
	return &redisRepository{client: client}
}

func (r *redisRepository) getSessionKey(sessionID uuid.UUID) string {
	return fmt.Sprintf("chat:session:%s", sessionID.String())
}

func (r *redisRepository) SaveMessage(ctx context.Context, sessionID uuid.UUID, message *model.ChatMessage, ttl time.Duration) error {
	key := r.getSessionKey(sessionID)

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	pipe := r.client.Pipeline()
	pipe.RPush(ctx, key, data)
	pipe.Expire(ctx, key, ttl)

	_, err = pipe.Exec(ctx)
	return err
}

func (r *redisRepository) GetMessages(ctx context.Context, sessionID uuid.UUID) ([]model.ChatMessage, error) {
	key := r.getSessionKey(sessionID)

	results, err := r.client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil
	}

	var messages []model.ChatMessage
	for _, stringData := range results {
		var msg model.ChatMessage
		if err := json.Unmarshal([]byte(stringData), &msg); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func (r *redisRepository) ClearCache(ctx context.Context, sessionID uuid.UUID) error {
	key := r.getSessionKey(sessionID)
	return r.client.Del(ctx, key).Err()
}
