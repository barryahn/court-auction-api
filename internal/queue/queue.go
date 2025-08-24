package queue

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Queue struct {
    Client *redis.Client
}

func NewQueue(addr string) *Queue {
    rdb := redis.NewClient(&redis.Options{Addr: addr})
    return &Queue{Client: rdb}
}

// 작업 추가
func (q *Queue) PushJob(url string) error {
    return q.Client.RPush(ctx, "crawl_queue", url).Err()
}

// 작업 꺼내기 (블록킹)
func (q *Queue) PopJob() (string, error) {
    res, err := q.Client.BLPop(ctx, 0, "crawl_queue").Result()
    if err != nil {
        return "", err
    }
    return res[1], nil
}
