.PHONY: help

help:
	@echo "Hello"

clean:
	find . -type d -name "__pycache__" -exec rm -rf {} +
	rm -rf .ruff_cache .black_cache .mypy_cache .pytest_cache build dist *.egg-info
	find . -type f -name "*.log" -delete
	@echo "All cache cleared"