package config

import (
	"flag"
	"time"
)

type Config struct {
	Addr              string
	Workers           int
	QueueBuffer       int
	TaskTimeout       time.Duration
	SchedulerInterval time.Duration
	Store             string
}

func Load() Config {
	addr := flag.String("addr", ":50051", "gRPC server address")
	workers := flag.Int("workers", 4, "number of concurrent workers")
	queueBuffer := flag.Int("queue-buffer", 1024, "internal queue buffer size")
	taskTimeout := flag.Duration("task-timeout", 2*time.Second, "timeout per task execution")
	schedulerInterval := flag.Duration("scheduler-interval", time.Second, "scheduler interval")
	store := flag.String("store", "in-memory", "task store type")

	flag.Parse()

	return Config{
		Addr:              *addr,
		Workers:           *workers,
		QueueBuffer:       *queueBuffer,
		TaskTimeout:       *taskTimeout,
		SchedulerInterval: *schedulerInterval,
		Store:             *store,
	}
}
