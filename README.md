# TaskForge

TaskForge is a lightweight **Cloud Tasks-inspired task queue** built in Go.

It was created as a portfolio project to explore backend concepts using:

- Go
- gRPC
- Protocol Buffers
- Goroutines
- Channels
- Context cancellation
- Worker pools
- Scheduling
- Retries
- Dead-letter queues
- CLI tooling

The project currently runs with an **in-memory store**, which means all tasks are lost when the server stops.

---

## Overview

TaskForge allows you to create tasks, schedule them, execute them with concurrent workers, retry failed tasks, and move tasks to a dead-letter state after the max attempts are reached.

High-level flow:

```txt
Client / CLI
  ↓
gRPC Server
  ↓
Task Store
  ↓
Scheduler
  ↓
Queue Channel
  ↓
Worker Pool
  ↓
Completed / Retry / Dead Letter
```

---

## Features

- Create tasks through a gRPC API
- Interact with the system through a CLI
- Schedule tasks with delay
- Execute tasks with a worker pool
- Process jobs internally using channels
- Retry failed tasks
- Move exhausted tasks to dead-letter
- Watch task updates in real time through gRPC streaming
- Cancel pending or scheduled tasks
- Graceful server shutdown
- Pretty terminal UI for server startup and CLI output

---

## Project Structure

```txt
.
├── cmd
│   ├── server
│   │   └── main.go
│   └── taskforge
│       └── main.go
├── docs
│   └── screenshots
├── internal
│   ├── app
│   │   └── app.go
│   ├── cliui
│   │   ├── server.go
│   │   └── terminal.go
│   ├── config
│   │   └── config.go
│   ├── grpcserver
│   │   ├── convert.go
│   │   └── service.go
│   ├── queue
│   │   └── manager.go
│   ├── scheduler
│   │   └── scheduler.go
│   ├── store
│   │   ├── memory.go
│   │   └── store.go
│   ├── task
│   │   └── task.go
│   └── worker
│       ├── handler.go
│       └── pool.go
├── proto
│   └── taskforge
│       └── v1
│           ├── taskforge.proto
│           ├── taskforge.pb.go
│           └── taskforge_grpc.pb.go
├── Makefile
├── go.mod
└── go.sum
```

### Main folders

| Folder | Purpose |
|---|---|
| `cmd/server` | Starts the TaskForge gRPC server |
| `cmd/taskforge` | CLI entrypoint |
| `internal/app` | Wires the application together |
| `internal/config` | Loads server flags/config |
| `internal/grpcserver` | Implements the gRPC service |
| `internal/queue` | Internal channel-based queue |
| `internal/scheduler` | Dispatches scheduled tasks when they are ready |
| `internal/store` | Store interface and in-memory implementation |
| `internal/task` | Internal task domain model |
| `internal/worker` | Worker pool and task handler |
| `internal/cliui` | Terminal UI helpers |
| `proto` | Protobuf contract and generated gRPC files |

---

## Requirements

You need:

- Go
- Make
- `protoc`
- `protoc-gen-go`
- `protoc-gen-go-grpc`

Install the Protobuf Go generators:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Make sure Go binaries are available in your `PATH`:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

---

## Installation

Clone the repository:

```bash
git clone https://github.com/pablo-banker/taskforge.git
cd taskforge
```

Install dependencies:

```bash
go mod tidy
```

Generate Protobuf files:

```bash
make proto
```

---

## Running the Server

Start the gRPC server:

```bash
make run-server
```

By default, the server runs on:

```txt
localhost:50051
```

You can pass custom server args:

```bash
make run-server ARGS="--addr :50051 --workers 4 --task-timeout 2s"
```

Available server flags:

| Flag | Default | Description |
|---|---:|---|
| `--addr` | `:50051` | gRPC server address |
| `--workers` | `4` | Number of concurrent workers |
| `--queue-buffer` | `1024` | Internal queue channel buffer size |
| `--task-timeout` | `2s` | Timeout per task execution |
| `--scheduler-interval` | `1s` | Interval used by the scheduler |
| `--store` | `in-memory` | Store type used by the app |

Example:

```bash
make run-server ARGS="--addr :50051 --workers 8 --task-timeout 5s --scheduler-interval 500ms"
```

---

## CLI Usage

The CLI can be executed through the Makefile:

```bash
make run-cli ARGS="<command>"
```

Or directly:

```bash
go run ./cmd/taskforge <command>
```

Available commands:

```txt
create
get
list
cancel
watch
```

---

## Create a Task

```bash
make create-task ARGS="--name welcome --payload 'hello world'"
```

Or:

```bash
go run ./cmd/taskforge create --name welcome --payload "hello world"
```

Example output:

```txt
╭─ Task 9d21371d ─────────────────────────────────────────────────────────────╮
│ ID           9d21371d-29f1-48ba-8259-670a1f67cbef                           │
│ Queue        default                                                        │
│ Name         welcome                                                        │
│ Payload      hello world                                                    │
│ Status       ● pending                                                      │
│ Priority     normal                                                         │
│ Attempts     0/3                                                            │
│ Run at       2026-07-13T10:58:31-03:00                                      │
│ Created      2026-07-13T10:58:21-03:00                                      │
│ Updated      2026-07-13T10:58:21-03:00                                      │
╰─────────────────────────────────────────────────────────────────────────────╯
```

---

## Create a Task in a Custom Queue

```bash
make create-task ARGS="--queue emails --name welcome-email --payload 'send welcome email'"
```

Or:

```bash
go run ./cmd/taskforge create \
  --queue emails \
  --name welcome-email \
  --payload "send welcome email"
```

---

## Create a Scheduled Task

```bash
make create-task ARGS="--name email --payload 'send email' --delay 10s"
```

Or:

```bash
go run ./cmd/taskforge create \
  --queue emails \
  --name welcome-email \
  --payload "send welcome email" \
  --delay 10s
```

The task will be created with the `scheduled` status and executed after the delay.

---

## Create a Failing Task

Use a payload containing `fail` to simulate a failure:

```bash
make create-task ARGS="--name broken --payload 'fail' --max-attempts 3"
```

Or:

```bash
go run ./cmd/taskforge create \
  --name broken \
  --payload "fail" \
  --max-attempts 3
```

The worker will retry the task until it reaches the max attempts.  
After that, the task moves to `dead_letter`.

---

## Create a Slow Task

Use a payload containing `slow` to simulate a slow task:

```bash
make create-task ARGS="--name slow-task --payload 'slow job'"
```

If the task takes longer than the configured `--task-timeout`, it fails and enters the retry flow.

---

## List Tasks

```bash
make list-tasks
```

Or:

```bash
go run ./cmd/taskforge list
```

---

## List Tasks by Queue

```bash
make list-tasks ARGS="--queue emails"
```

Or:

```bash
go run ./cmd/taskforge list --queue emails
```

---

## List Tasks by Status

```bash
make list-tasks ARGS="--status completed"
```

Or:

```bash
go run ./cmd/taskforge list --status completed
```

Available statuses:

```txt
pending
scheduled
queued
running
completed
failed
cancelled
dead_letter
```

---

## Get a Task

```bash
make get-task ARGS="--id <task-id>"
```

Or:

```bash
go run ./cmd/taskforge get --id <task-id>
```

Example:

```bash
make get-task ARGS="--id 9d21371d-29f1-48ba-8259-670a1f67cbef"
```

---

## Watch a Task

```bash
make watch-task ARGS="--id <task-id>"
```

Or:

```bash
go run ./cmd/taskforge watch --id <task-id>
```

Example:

```bash
make watch-task ARGS="--id 9d21371d-29f1-48ba-8259-670a1f67cbef"
```

This command streams task updates in real time.

Example successful lifecycle:

```txt
scheduled
running
completed
```

Example failing lifecycle:

```txt
pending
running
scheduled
running
scheduled
running
dead_letter
```

---

## Cancel a Task

```bash
make cancel-task ARGS="--id <task-id>"
```

Or:

```bash
go run ./cmd/taskforge cancel --id <task-id>
```

Only tasks that are not running or already finished can be cancelled.

---

## Makefile Commands

| Command | Description |
|---|---|
| `make proto` | Generate Go files from Protobuf |
| `make run-server` | Run the gRPC server |
| `make run-server ARGS="..."` | Run the server with custom flags |
| `make run-cli ARGS="<command>"` | Run the CLI with custom args |
| `make create-task ARGS="..."` | Create a task |
| `make list-tasks ARGS="..."` | List tasks |
| `make get-task ARGS="..."` | Get a task by ID |
| `make watch-task ARGS="..."` | Watch task updates |
| `make cancel-task ARGS="..."` | Cancel a task |
| `make build` | Build server and CLI binaries |
| `make test` | Run tests |
| `make race` | Run tests with race detector |
| `make tidy` | Run `go mod tidy` |

---

## Example Flow

Start the server:

```bash
make run-server
```

Create a scheduled task:

```bash
make create-task ARGS="--name email --payload 'send email' --delay 10s"
```

Copy the task ID and watch it:

```bash
make watch-task ARGS="--id <task-id>"
```

List all tasks:

```bash
make list-tasks
```

Create a failing task:

```bash
make create-task ARGS="--name broken --payload 'fail' --max-attempts 3"
```

Watch it until it reaches `dead_letter`:

```bash
make watch-task ARGS="--id <task-id>"
```

---

## Task Lifecycle

A task can move through the following statuses:

```txt
scheduled -> pending -> running -> completed
```

If it fails:

```txt
pending -> running -> scheduled -> running -> dead_letter
```

If it is cancelled before execution:

```txt
scheduled -> cancelled
```

---

## Retry Behavior

When a task fails:

1. The worker marks it as failed internally.
2. If attempts are still available, the task is scheduled again.
3. The retry delay uses a simple backoff based on the attempt number.
4. If attempts are exhausted, the task becomes `dead_letter`.

Example:

```txt
attempt 1 failed -> retry later
attempt 2 failed -> retry later
attempt 3 failed -> dead_letter
```

---

## Server Startup UI

When the server starts, TaskForge prints a terminal dashboard with:

- ASCII logo
- Runtime config
- Makefile commands
- Quick examples
- Server status

Example:

```txt
TASKFORGE

A lightweight Cloud Tasks-inspired queue built with Go, gRPC, Protobuf, goroutines and channels.

╭─ Runtime ─────────────────────────────────────╮
│ Address       :50051                           │
│ Workers       4                                │
│ Task timeout  2s                               │
│ Scheduler     1s                               │
│ Queue buffer  1024                             │
│ Store         in-memory                        │
│ Go version    go1.25.3                         │
╰────────────────────────────────────────────────╯

✓ TaskForge server is ready
```

---

## Build Binaries

```bash
make build
```

This creates:

```txt
bin/taskforge-server
bin/taskforge
```

Run the server binary:

```bash
./bin/taskforge-server
```

Run the CLI binary:

```bash
./bin/taskforge create --name welcome --payload "hello"
```

---

## Tests

Run tests:

```bash
make test
```

Run tests with the race detector:

```bash
make race
```

---

## Development

Format all files:

```bash
go fmt ./...
```

Tidy dependencies:

```bash
make tidy
```

Regenerate Protobuf files after editing the `.proto` contract:

```bash
make proto
```

---

## Troubleshooting

### Connection refused

If you see:

```txt
rpc error: code = Unavailable desc = connection error
```

The server is probably not running.

Start it first:

```bash
make run-server
```

Then run the CLI in another terminal.

---

### `protoc-gen-go` not found

If `make proto` fails because `protoc-gen-go` cannot be found, install the generators:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Then ensure your Go bin path is available:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

---

### Terminal colors are not desired

Disable colors with:

```bash
NO_COLOR=1 make run-server
```

Or:

```bash
NO_COLOR=1 make list-tasks
```

---

## Current Limitations

- Tasks are stored in memory
- Tasks are lost when the server stops
- There is no authentication
- There is no distributed locking
- There is no persistence layer yet
- The task handler is simulated through payload keywords like `fail` and `slow`

---


## License

MIT
