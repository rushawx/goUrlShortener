docker-up: 
	docker network create url-short || true
	docker compose up -d --build

docker-down:
	docker compose down --volumes

start-redis: 
	go run main.go

start-memory: 
	go run main.go -d
