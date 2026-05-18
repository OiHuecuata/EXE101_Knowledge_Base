.DEFAULT_GOAL := help

.PHONY: help tree clean packages parse pipeline migrate-down migrate-up ui

export PYTHONPATH := $(shell pwd)/src

help:
	@echo "Available commands:"
	@echo "  make help          : Show commands and how to use"
	@echo "  make tree          : Display project structure"
	@echo "  make clean         : Remove cache files, logs"
	@echo "  make packages      : Install all dependencies"
	@echo "  make parse         : Convert PPTX/DOCX to Markdown"
	@echo "  make pipeline      : Run the main pipeline"
	@echo "  make migrate-up    : Run database migrations up"
	@echo "  make migrate-down  : Rollback database migrations"
	@echo "  make ui            : Run frontend development server"

packages:
	poetry install --no-root

tree:
	tree -I "EXE101|processed|node_modules"

clean:
	find . -type d -name "__pycache__" -exec rm -rf {} +
	rm -rf .ruff_cache .black_cache .mypy_cache .pytest_cache build dist *.egg-info
	find . -type f -name "*.log" -delete

parse:
	poetry run python src/pipeline/parser.py

pipeline:
	poetry run python src/pipeline/main.py

migrate-down:
	psql "$(DB_URL)" -f database/migration/init_schema.down.sql

migrate-up:
	psql "$(DB_URL)" -f database/migration/init_schema.up.sql

ui:
	cd src/frontend && npm run dev