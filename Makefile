.PHONY: obu
obu:
	@go build -o bin/obu obu/main.go
	@./bin/obu

reveiver:
	@go build -o bin/receiver ./receiver
	@./bin/receiver