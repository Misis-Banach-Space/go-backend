debug:
	docker compose -f deployment/docker-compose.yml up db rabbit --build -d 
	go run cmd/server/main.go

deploy:
	docker compose -f deployment/docker-compose.yml up --build -d
