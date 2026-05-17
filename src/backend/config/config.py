import os
from pathlib import Path
from typing import Optional
from pydantic import Field
from pydantic_settings import BaseSettings, SettingsConfigDict

BASE_DIR = Path(os.getcwd())
ENV_FILE_PATH = BASE_DIR / ".env"

class Settings(BaseSettings):
    
    # DB CONFIG
    DB_HOST: str = Field(default="localhost", validation_alias="DB_HOST")
    DB_PORT: int = Field(default=5432, validation_alias="DB_PORT")
    DB_USER: str = Field(default="postgres", validation_alias="DB_USER")
    DB_PASSWORD: str = Field(default="", validation_alias="DB_PASSWORD")
    DB_NAME: str = Field(default="exe101", validation_alias="DB_NAME")
    DB_URL: Optional[str] = Field(default=None, validation_alias="DB_URL")

    # REDIS CACHE
    REDIS_HOST: str = Field(default="localhost", validation_alias="REDIS_HOST")
    REDIS_PORT: int = Field(default=6379, validation_alias="REDIS_PORT")
    REDIS_PASSWORD: str = Field(default="", validation_alias="REDIS_PASSWORD")
    DB_REDIS: int = Field(default=0, validation_alias="DB_REDIS")

    # OPENROUTER
    OPENROUTER_API_KEY: str = Field(..., validation_alias="OPENROUTER_API_KEY")
    OPENROUTER_BASE_URL: str = Field(default="https://openrouter.ai/api/v1", validation_alias="OPENROUTER_BASE_URL")
    LLM_MODEL_NAME: str = Field(..., validation_alias="LLM_MODEL_NAME")
    LLM_TEMPERATURE: float = Field(default=0.2, validation_alias="LLM_TEMPERATURE", ge=0.0, le=1.0)
    LLM_MAX_TOKENS: int = Field(default=2048, validation_alias="LLM_MAX_TOKENS", gt=0)

    # HUGGINGFACE
    HF_TOKEN: str = Field(..., validation_alias="HF_TOKEN")
    EMBEDDING_MODEL_NAME: str = Field(default="BAAI/bge-base-en-v1.5", validation_alias="EMBEDDING_MODEL_NAME")

    # APP CONFIG
    APP_ENV: str = Field(default="development", validation_alias="APP_ENV")
    APP_PORT: int = Field(default=8080, validation_alias="APP_PORT")
    PYTHON_LLM_SERVICE_URL: str = Field(default="http://localhost:5000", validation_alias="PYTHON_LLM_SERVICE_URL")
    CACHE_TTL_SECONDS: int = Field(default=3600, validation_alias="CACHE_TTL_SECONDS", gt=0)

    model_config = SettingsConfigDict(
        env_file=ENV_FILE_PATH,
        env_file_encoding="utf-8",
        extra="ignore"
    )

settings = Settings()