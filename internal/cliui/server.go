package cliui

import (
	"fmt"
	"runtime"
	"time"
)

type ServerBannerConfig struct {
	Addr              string
	Workers           int
	TaskTimeout       time.Duration
	SchedulerInterval time.Duration
	QueueBuffer       int
	Store             string
}

func PrintServerBanner(cfg ServerBannerConfig) {
	fmt.Println()
	printLogo()
	fmt.Println()

	printTop(" Runtime ", 78)
	printRow("Address", color(cyan, cfg.Addr), 78)
	printRow("Workers", fmt.Sprintf("%d", cfg.Workers), 78)
	printRow("Task timeout", cfg.TaskTimeout.String(), 78)
	printRow("Scheduler", cfg.SchedulerInterval.String(), 78)
	printRow("Queue buffer", fmt.Sprintf("%d", cfg.QueueBuffer), 78)
	printRow("Store", cfg.Store, 78)
	printRow("Go version", runtime.Version(), 78)
	printBottom(78)

	fmt.Println()

	printCommands()

	fmt.Println()
	PrintSuccess("TaskForge server is ready")
	PrintInfo("Keep this terminal open and use another terminal to run the CLI commands.")
	fmt.Println()
}

func printLogo() {
	logo := `
████████╗ █████╗ ███████╗██╗  ██╗███████╗ ██████╗ ██████╗  ██████╗ ███████╗
╚══██╔══╝██╔══██╗██╔════╝██║ ██╔╝██╔════╝██╔═══██╗██╔══██╗██╔════╝ ██╔════╝
   ██║   ███████║███████╗█████╔╝ █████╗  ██║   ██║██████╔╝██║  ███╗█████╗  
   ██║   ██╔══██║╚════██║██╔═██╗ ██╔══╝  ██║   ██║██╔══██╗██║   ██║██╔══╝  
   ██║   ██║  ██║███████║██║  ██╗██║     ╚██████╔╝██║  ██║╚██████╔╝███████╗
   ╚═╝   ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝╚═╝      ╚═════╝ ╚═╝  ╚═╝ ╚═════╝ ╚══════╝
`

	fmt.Print(color(cyan, logo))
	fmt.Println(color(gray, "A lightweight Cloud Tasks-inspired queue built with Go, gRPC, Protobuf, goroutines and channels."))
}

func printCommands() {
	width := 100

	printTop(" Quick Start ", width)
	printCommand("Start server", "make run-server", width)
	printCommand("Custom server", `make run-server ARGS="--addr :50051 --workers 4 --task-timeout 2s"`, width)
	printBottom(width)

	fmt.Println()

	printTop(" Create Tasks ", width)
	printCommand("Simple task", `make create-task ARGS="--name welcome --payload 'hello world'"`, width)
	printCommand("Scheduled task", `make create-task ARGS="--queue emails --name welcome-email --payload 'send email' --delay 10s"`, width)
	printCommand("Failing task", `make create-task ARGS="--name broken-task --payload 'fail' --max-attempts 3"`, width)
	printBottom(width)

	fmt.Println()

	printTop(" Inspect Tasks ", width)
	printCommand("List all", "make list-tasks", width)
	printCommand("List by status", `make list-tasks ARGS="--status completed"`, width)
	printCommand("Get by ID", `make get-task ARGS="--id <task-id>"`, width)
	printCommand("Watch updates", `make watch-task ARGS="--id <task-id>"`, width)
	printCommand("Cancel task", `make cancel-task ARGS="--id <task-id>"`, width)
	printBottom(width)
}

func printCommand(label string, command string, width int) {
	labelLine := fmt.Sprintf("│ %s", color(gray, label))
	fmt.Println(labelLine + spaces(width-visibleLen(labelLine)-1) + "│")

	commandLine := fmt.Sprintf("│   %s", color(green, command))
	fmt.Println(commandLine + spaces(width-visibleLen(commandLine)-1) + "│")
}
