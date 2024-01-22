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