all:
	go build -o bin/subscriber sub/cmd/main.go
docker-push:
	docker-compose up
