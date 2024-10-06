.PHONY: run

run: 
	go run cmd/songs/main.go cmd/songs/logger.go -path_to_config .env