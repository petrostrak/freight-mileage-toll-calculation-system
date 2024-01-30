.PHONY: obu
obu:
	@go build -o bin/obu obu/main.go
	@./bin/obu

reveiver:
	@go build -o bin/receiver ./receiver
	@./bin/receiver

.PHONY: distance_calculator
distance_calculator:
	@go build -o bin/distance_calculator ./distance_calculator
	@./bin/distance_calculator

.PHONY: aggregator
aggregator:
	@go build -o bin/aggregator ./aggregator
	@./bin/aggregator

.PHONY: proto
proto: 
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/ptypes.proto

.PHONY: gateway
gateway:
	@go build -o bin/gateway ./gateway
	@./bin/gateway