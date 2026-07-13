package worker

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pablo-banker/taskforge/internal/task"
)

type Handler interface {
	Handle(ctx context.Context, t task.Task) error
}

type SimulatedHandler struct{}

func (h SimulatedHandler) Handle(ctx context.Context, t task.Task) error {
	payload := strings.ToLower(strings.TrimSpace(t.Payload))

	if strings.Contains(payload, "fail") {
		return errors.New("simulated task failure")
	}

	duration := 300 * time.Millisecond

	if strings.Contains(payload, "slow") {
		duration = 3 * time.Second
	}

	select {
	case <-time.After(duration):
		fmt.Printf(
			"[handler] executed task=%s queue=%s name=%s payload=%q\n",
			t.ID,
			t.Queue,
			t.Name,
			t.Payload,
		)

		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}
