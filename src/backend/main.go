package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"backend/config"
	"backend/handler"
	"backend/repository"
	"backend/service/chat"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func main() {

	cfg := config.LoadConfig()

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	pgPool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pgPool.Close()

	if err := pgPool.Ping(ctx); err != nil {
		log.Fatalf("Database ping failed: %v\n", err)
	}

	redisAddr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.DBRedis,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Redis ping failed: %v\n", err)
	}
	defer rdb.Close()

	pgRepo := repository.NewPostgresRepository(pgPool)
	redisRepo := repository.NewRedisRepository(rdb)
	chatService := chat.NewChatService(cfg, pgRepo, redisRepo)
	chatHandler := handler.NewChatHandler(chatService)

	r := gin.Default()

	api := r.Group("/chats")
	{
		api.POST("", chatHandler.CreateSession)
		api.GET("", chatHandler.GetSessions)
		api.POST("/:id/messages", chatHandler.StreamMessage)
		api.GET("/:id/messages", chatHandler.GetMessages)
	}

	serverAddr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server đang chạy ổn định tại port %s trong môi trường %s", cfg.AppPort, cfg.AppEnv)
	if err := http.ListenAndServe(serverAddr, r); err != nil {
		log.Fatalf("Server dừng đột ngột: %v\n", err)
	}
}
