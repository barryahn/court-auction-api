package main

import (
	"context"
	"net/http"
	"time"

	"github.com/barryx002/court-auction-api/internal/config"
	"github.com/barryx002/court-auction-api/internal/queue"
	"github.com/barryx002/court-auction-api/internal/repository"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	db, err := repository.NewDB(ctx, cfg.PostgresDSN)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r := gin.Default()

	// 헬스체크: DB 연결 확인
	r.GET("/healthz", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := db.Pool.Ping(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	q := queue.NewQueue(cfg.RedisAddr)

	// 작업 요청 → 큐에 URL 추가
	r.POST("/crawl", func(c *gin.Context) {
		url := c.Query("url")
		if url == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
			return
		}
		if err := q.PushJob(url); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "queued", "url": url})
	})

	// (나중에 DB 붙이면) 크롤링된 결과 반환
	r.GET("/auctions", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "여기에 DB 결과 반환"})
	})

	r.Run(cfg.APIServerAddr)
}
