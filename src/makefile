BACKEND_IMAGE = nam-observability/backend

TAG = latest

build-backend:
	docker build -t $(BACKEND_IMAGE):$(TAG) ./backend

build: build-backend

run:
	docker compose up -d

down:
	docker compose down
