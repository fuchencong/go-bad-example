.PHONY: dev backend frontend build test

dev:
	@echo "Run 'make backend' and 'make frontend' in separate terminals"

backend:
	go run ./cmd/server

frontend:
	cd frontend && npm run dev

build:
	go build ./...
	cd frontend && npm run build

test:
	go test ./...
	cd frontend && npm run lint
