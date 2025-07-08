package workers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2/log"
	"github.com/hibiken/asynq"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	OTP   string `json:"opt"`
}

func (d *RedisTaskDistributor) DistributeTaskSendVerifyEmail(payload *PayloadSendVerifyEmail, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	info, err := d.Client.EnqueueContext(context.Background(), task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info(fmt.Sprintf("type: %s | payload: %d | queue: %d | max_retry: %d | enqueued task", task.Type(), task.Payload(), info.Queue, info.MaxRetry))
	return nil
}
