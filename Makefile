.PHONY: setup redis api worker web

setup:
	cd frontend && npm install
	cd worker && npm install
	cd backend && go mod tidy

redis:
	docker compose up -d redis

api:
	cd backend && go run ./cmd/api

worker:
	cd worker && npm run start

web:
	cd frontend && npm run dev
