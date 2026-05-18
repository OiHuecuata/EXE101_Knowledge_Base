package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	// DB Config
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBURL      string `mapstructure:"DB_URL"`

	// Redis Config
	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     string `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	DBRedis       int    `mapstructure:"DB_REDIS"`

	// API Key
	OpenRouterAPIKey  string  `mapstructure:"OPENROUTER_API_KEY"`
	OpenRouterBaseURL string  `mapstructure:"OPENROUTER_BASE_URL"`
	LLMModelName      string  `mapstructure:"LLM_MODEL_NAME"`
	LLMTemperature    float64 `mapstructure:"LLM_TEMPERATURE"`
	LLMMaxTokens      int     `mapstructure:"LLM_MAX_TOKENS"`

	// HuggingFace
	HFToken            string `mapstructure:"HF_TOKEN"`
	EmbeddingModelName string `mapstructure:"EMBEDDING_MODEL_NAME"`

	// Application
	AppEnv              string `mapstructure:"APP_ENV"`
	AppPort             string `mapstructure:"APP_PORT"`
	PythonLLMServiceURL string `mapstructure:"PYTHON_LLM_SERVICE_URL"`
	CacheTTLSeconds     int    `mapstructure:"CACHE_TTL_SECONDS"`
}

func LoadConfig() *Config {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./src/backend")

	// Environment Variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Warning: No .env file found, relying on system environment variables")
	}

	// Default
	setDefaults()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	return &config
}

func setDefaults() {
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_PORT", "8080")

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_NAME", "exe101_db")

	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("DB_REDIS", 0)

	viper.SetDefault("OPENROUTER_BASE_URL", "https://openrouter.ai/api/v1")
	viper.SetDefault("LLM_TEMPERATURE", 0.2)
	viper.SetDefault("LLM_MAX_TOKENS", 2048)
	viper.SetDefault("EMBEDDING_MODEL_NAME", "BAAI/bge-base-en-v1.5")
	viper.SetDefault("CACHE_TTL_SECONDS", 3600)
}
