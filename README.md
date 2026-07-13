# TaskForge

TaskForge is a lightweight **Cloud Tasks-inspired task queue** built in Go.

It is a portfolio project focused on backend fundamentals: **gRPC**, **Protocol Buffers**, **goroutines**, **channels**, **worker pools**, scheduling, retries and dead-letter handling.

> Current storage: in-memory. Tasks are lost when the server stops.

---

## Features

- gRPC API with Protobuf contracts
- CLI commands for creating, listing, watching and cancelling tasks
- Interactive terminal UI with dashboard and task navigation
- Scheduled tasks with delay
- Concurrent workers using goroutines
- Internal queue using channels
- Retry flow with simple backoff
- Dead-letter state after max attempts
- Pretty terminal output and server startup banner
- Graceful shutdown

---

## Architecture

```txt
CLI / TUI
   ↓
gRPC Server
   ↓
Task Store
   ↓
Scheduler
   ↓
Channel Queue
   ↓
Worker Pool
   ↓
Completed / Retry / Dead Letter
```

---

## Project Structure

```txt
.
├── cmd
│   ├── server        # gRPC server entrypoint
│   └── taskforge     # CLI/TUI entrypoint
├── internal
│   ├── app           # app wiring
│   ├── cliui         # terminal output helpers
│   ├── config        # server flags/config
│   ├── grpcserver    # gRPC service implementation
│   ├── queue         # channel-based queue
│   ├── scheduler     # scheduled task dispatcher
│   ├── store         # store interface + memory store
│   ├── task          # task domain model
│   ├── tui           # interactive terminal UI
│   └── worker        # worker pool + handler
├── proto             # Protobuf contract and generated files
├── docs
├── Makefile
├── go.mod
└── go.sum
```

---

## Requirements

- Go
- Make
- protoc
- protoc-gen-go
- protoc-gen-go-grpc

Install Protobuf generators:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Make sure Go binaries are available:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

---

## Setup

```bash
git clone https://github.com/pablo-banker/taskforge.git
cd taskforge

go mod tidy
make proto
```

---

## Run the Server

```bash
make run-server
```

Default server address:

```txt
localhost:50051
```

Run with custom args:

```bash
make run-server ARGS="--addr :50051 --workers 4 --task-timeout 2s"
```

Server flags:

| Flag | Default | Description |
|---|---:|---|
| `--addr` | `:50051` | gRPC server address |
| `--workers` | `4` | Number of workers |
| `--queue-buffer` | `1024` | Internal channel buffer |
| `--task-timeout` | `2s` | Timeout per task |
| `--scheduler-interval` | `1s` | Scheduler tick interval |
| `--store` | `in-memory` | Store type |

---

## CLI Usage

Create a task:

```bash
make create-task ARGS="--name welcome --payload 'hello world'"
```

Create a scheduled task:

```bash
make create-task ARGS="--name email --payload 'send email' --delay 10s"
```

Create a failing task to test retries and dead-letter:

```bash
make create-task ARGS="--name broken --payload 'fail' --max-attempts 3"
```

List tasks:

```bash
make list-tasks
```

List by status:

```bash
make list-tasks ARGS="--status completed"
```

Get a task:

```bash
make get-task ARGS="--id <task-id>"
```

Watch task updates in real time:

```bash
make watch-task ARGS="--id <task-id>"
```

Cancel a task:

```bash
make cancel-task ARGS="--id <task-id>"
```

You can also run the CLI directly:

```bash
go run ./cmd/taskforge <command>
```

---

## Interactive TUI

TaskForge also includes an interactive terminal UI.

Start the server first:

```bash
make run-server
```

Then open the TUI in another terminal:

```bash
make tui
```

Or run directly:

```bash
go run ./cmd/taskforge tui
```

Use a custom address:

```bash
make tui ARGS="--addr localhost:50051"
```

### TUI shortcuts

| Key | Action |
|---|---|
| `q` | Quit |
| `tab` | Next page |
| `1` | Dashboard |
| `2` | Tasks |
| `3` / `c` | Create task |
| `4` | Dead Letter |
| `r` | Refresh |
| `↑/↓` or `j/k` | Move selection |
| `enter` | Open selected task |
| `b` | Back |
| `x` | Cancel selected task |
| `ctrl+s` | Submit create form |
| `esc` | Leave create form |

---

## Task Lifecycle

Successful task:

```txt
scheduled -> pending -> running -> completed
```

Failing task:

```txt
pending -> running -> scheduled -> running -> dead_letter
```

Cancelled task:

```txt
scheduled -> cancelled
```

---

## Makefile Commands

| Command | Description |
|---|---|
| `make proto` | Generate Go files from Protobuf |
| `make run-server` | Run the gRPC server |
| `make create-task ARGS="..."` | Create a task |
| `make list-tasks ARGS="..."` | List tasks |
| `make get-task ARGS="..."` | Get a task by ID |
| `make watch-task ARGS="..."` | Watch task updates |
| `make cancel-task ARGS="..."` | Cancel a task |
| `make tui ARGS="..."` | Open interactive terminal UI |
| `make build` | Build server and CLI binaries |
| `make test` | Run tests |
| `make race` | Run tests with race detector |
| `make tidy` | Run `go mod tidy` |

---

## Build

```bash
make build
```

This creates:

```txt
bin/taskforge-server
bin/taskforge
```

---

## Development

```bash
go fmt ./...
make tidy
make test
```

Regenerate Protobuf files after editing the `.proto`:

```bash
make proto
```

---

## Troubleshooting

### Connection refused

The server is probably not running.

```bash
make run-server
```

Then run the CLI or TUI in another terminal.

### Protobuf generator not found

Install the generators:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### Disable terminal colors

```bash
NO_COLOR=1 make run-server
```

---

## Current Limitations

- In-memory storage only
- No authentication
- No distributed locking
- No persistence layer yet
- Task handler is simulated with payload keywords like `fail` and `slow`

---

## Roadmap

- SQLite or PostgreSQL persistence
- Priority queue
- Configurable retry policies
- Queue-level rate limits
- Structured logs
- Metrics endpoint
- Dockerfile
- CI with GitHub Actions
- More unit and integration tests

---

## License

MIT