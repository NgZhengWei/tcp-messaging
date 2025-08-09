.PHONY: server
server:
	go run cmd/server/server.go

.PHONY: client
client:
	go run cmd/client/client.go
