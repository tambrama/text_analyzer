CMD_A = cmd/receiver
CMD_B = cmd/analyzer
BUILD_DIR = build
BIN_A = $(BUILD_DIR)/receiver
BIN_B = $(BUILD_DIR)/service-b

DOCKER_FILE = docker-compose.yml

.PHONY: all build run clean fmt mod docker-up docker-down docker-rebuild docker-logs docs test

all: docker-up


build: 
	@mkdir -p $(BUILD_DIR)
	go build -o $(BIN_A) ./$(CMD_A)
	go build -o $(BIN_B) ./$(CMD_B)

run: build
	@echo "Starting Service A..."
	./$(BIN_A) &
	@echo "Starting Service B..."
	./$(BIN_B) &
	@echo "Services running. Press Ctrl+C to stop."

clean:
	rm -rf $(BUILD_DIR)

fmt:
	go fmt ./...

mod:
	go mod tidy

docker-up:
	docker-compose -f $(DOCKER_FILE) up --build -d

docker-down:
	docker-compose -f $(DOCKER_FILE) down

docker-rebuild:
	docker-compose -f $(DOCKER_FILE) down -v && \
	docker-compose -f $(DOCKER_FILE) up --build -d


docker-logs:
	docker-compose -f $(DOCKER_FILE) logs -f

docs:
	swag init --parseDependency --parseInternal --generalInfo $(CMD_A)/main.go --output ./docs
	swag init --parseDependency --parseInternal --generalInfo $(CMD_B)/main.go --output ./docs

test:
	go test -v -tags=integration ./test