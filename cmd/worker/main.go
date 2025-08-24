package main

import (
	"fmt"

	"github.com/barryx002/court-auction-api/internal/crawler"
	"github.com/barryx002/court-auction-api/internal/queue"
)

func main() {
    q := queue.NewQueue("localhost:6379")

    for {
        url, err := q.PopJob()
        if err != nil {
            fmt.Println("큐에서 작업 꺼내기 실패:", err)
            continue
        }

        fmt.Println("크롤링 시작:", url)
        items, err := crawler.CrawlAuctionList(url)
        if err != nil {
            fmt.Println("크롤링 실패:", err)
            continue
        }

        // TODO: DB 저장 로직 추가
        fmt.Printf("크롤링 완료: %d 건 수집\n", len(items))
    }
}
