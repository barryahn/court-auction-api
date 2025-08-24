package main

import (
	"net/http"

	"github.com/barryx002/court-auction-api/internal/queue"

	"github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    q := queue.NewQueue("localhost:6379")

    // 작업 요청 → 큐에 URL 추가
    r.POST("/crawl", func(c *gin.Context) {
        url := c.Query("url")
        if url == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
            return
        }
        q.PushJob(url)
        c.JSON(http.StatusOK, gin.H{"status": "queued", "url": url})
    })

    // (나중에 DB 붙이면) 크롤링된 결과 반환
    r.GET("/auctions", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"data": "여기에 DB 결과 반환"})
    })

    r.Run(":8080")
}
