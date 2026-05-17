package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {

	// DB Config
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBURL      string

	// Redis Config
	RedisHost     string
	RedisPort     string
	RedisPassword string
	DBRedis       string

	// API Key
	OpenRouterAPIKey  string
	OpenRouterBaseURL string
	LLMModelName      string
	LLMTemperature    float64
	LLMMaxTokens      int

	// HuggingFace
	HFToken            string
	EmbeddingModelName string

	// Application
	AppEnv              string
	AppPort             string
	PythonLLMServiceURL string
	CacheTTLSeconds     int
}

func LoadConfig() *Config {

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using system environment variables")
	}

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "exe101_db"),
		DBURL:      getEnv("DB_URL", ""),

		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		DBRedis:       getEnv("DB_REDIS", "0"),

		OpenRouterAPIKey:  getEnv("OPENROUTER_API_KEY", ""),
		OpenRouterBaseURL: getEnv("OPENROUTER_BASE_URL", "https://openrouter.ai/api/v1"),
		LLMModelName:      getEnv("LLM_MODEL_NAME", ""),
		LLMTemperature:    getFloat("LLM_TEMPERATURE", 0.1),
		LLMMaxTokens:      getInt("LLM_MAX_TOKENS", 2048),

		HFToken:            getEnv("HF_TOKEN", ""),
		EmbeddingModelName: getEnv("EMBEDDING_MODEL_NAME", "BAAI/bge-m3"),

		AppEnv:              getEnv("APP_ENV", "development"),
		AppPort:             getEnv("APP_PORT", "8080"),
		PythonLLMServiceURL: getEnv("PYTHON_LLM_SERVICE_URL", "http://localhost:5000"),
		CacheTTLSeconds:     getInt("CACHE_TTL_SECONDS", 3600),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getInt(name string, defaultValue int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getFloat(name string, defaultValue float64) float64 {
	valueStr := getEnv(name, "")
	if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
		return value
	}
	return defaultValue
}
