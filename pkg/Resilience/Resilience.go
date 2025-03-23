package Resilience

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func Retry(operation func() error, maxRetries int, baseDelay int) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err := operation()
		if err == nil {
			return nil
		}
		delay := time.Duration(baseDelay*(1<<uint(i))) * time.Millisecond
		time.Sleep(delay)
	}
	return fmt.Errorf("after %d attempts, last error: %v", maxRetries, err)
}

func Timeout(operation func() error, timeout int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
	defer cancel()
	done := make(chan error, 1)

	go func() {
		done <- operation()
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return errors.New("operation timed out")
	}

}

type DeadLetterQueue struct {
	message []string
}

func NewDeadLetterQueue() *DeadLetterQueue {
	return &DeadLetterQueue{
		message: make([]string, 0),
	}
}

func (dlq *DeadLetterQueue) AddMessage(msg string) {
	dlq.message = append(dlq.message, msg)
}
func (dlq *DeadLetterQueue) GetMessages() []string {
	return dlq.message
}

func ProcessWithDLQ(messages []string, process func(msg string) error, dlq *DeadLetterQueue) {
	for _, msg := range messages {
		err := process(msg)
		if err != nil {
			dlq.AddMessage(msg)
		}
	}
}
