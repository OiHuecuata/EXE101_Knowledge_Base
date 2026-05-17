.DEFAULT_GOAL := help

.PHONY: help structure clean packages parse

help:
	@echo "Available commands:"
	@echo "  make help      : Show commands and how to use"
	@echo "  make tree      : Display project structure"
	@echo "  make clean     : Remove cache files, logs"
	@echo "  make packages  : Install all dependencies"
	@echo "  make parse     : Convert PPTX/DOCX to Markdown"

packages:
	poetry install --no-root

tree:
	tree -I "EXE101|processed"

clean:
	find . -type d -name "__pycache__" -exec rm -rf {} +
	rm -rf .ruff_cache .black_cache .mypy_cache .pytest_cache build dist *.egg-info
	find . -type f -name "*.log" -delete

parse:
	poetry run python src/parser/main.py