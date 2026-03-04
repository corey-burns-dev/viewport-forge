package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisQueue struct {
	client       *redis.Client
	queueKey     string
	statusPrefix string
}

type Viewport struct {
	Name   string `json:"name"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type CaptureJob struct {
	ID          string     `json:"id"`
	URL         string     `json:"url"`
	Requested   time.Time  `json:"requestedAt"`
	Viewports   []Viewport `json:"viewports,omitempty"`
	RequestedBy string     `json:"requestedBy,omitempty"`
}

func NewRedisQueue(ctx context.Context, addr, queueKey, statusPrefix string) (*RedisQueue, error) {
	client := redis.NewClient(&redis.Options{Addr: addr})
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return &RedisQueue{
		client:       client,
		queueKey:     queueKey,
		statusPrefix: statusPrefix,
	}, nil
}

func (q *RedisQueue) Enqueue(ctx context.Context, job CaptureJob) error {
	blob, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("marshal job: %w", err)
	}

	if err := q.client.LPush(ctx, q.queueKey, blob).Err(); err != nil {
		return fmt.Errorf("enqueue redis: %w", err)
	}

	statusKey := q.statusPrefix + job.ID
	if err := q.client.HSet(ctx, statusKey, map[string]interface{}{
		"id":           job.ID,
		"url":          job.URL,
		"state":        "queued",
		"requested_at": job.Requested.UTC().Format(time.RFC3339),
	}).Err(); err != nil {
		return fmt.Errorf("seed status: %w", err)
	}

	return nil
}

func (q *RedisQueue) GetStatus(ctx context.Context, jobID string) (map[string]string, error) {
	statusKey := q.statusPrefix + jobID
	result, err := q.client.HGetAll(ctx, statusKey).Result()
	if err != nil {
		return nil, fmt.Errorf("get status: %w", err)
	}
	if len(result) == 0 {
		return nil, errors.New("not found")
	}
	return result, nil
}

func (q *RedisQueue) Close() error {
	return q.client.Close()
}
