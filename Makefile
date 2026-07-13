PROTO_FILE=proto/taskforge/v1/taskforge.proto

.PHONY: proto
proto:
	protoc \
		-I proto \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_FILE)

.PHONY: run-server
run-server:
	go run ./cmd/server $(ARGS)

.PHONY: build
build:
	mkdir -p bin
	go build -o bin/taskforge-server ./cmd/server
	go build -o bin/taskforge ./cmd/taskforge

.PHONY: test
test:
	go test ./...

.PHONY: race
race:
	go test -race ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: run-cli
run-cli:
	go run ./cmd/taskforge $(ARGS)

.PHONY: create-task
create-task:
	go run ./cmd/taskforge create $(ARGS)

.PHONY: list-tasks
list-tasks:
	go run ./cmd/taskforge list $(ARGS)

.PHONY: get-task
get-task:
	go run ./cmd/taskforge get $(ARGS)

.PHONY: watch-task
watch-task:
	go run ./cmd/taskforge watch $(ARGS)

.PHONY: cancel-task
cancel-task:
	go run ./cmd/taskforge cancel $(ARGS)

.PHONY: tui
tui:
	go run ./cmd/taskforge tui $(ARGS)