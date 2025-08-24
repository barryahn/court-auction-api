package main

import (
	"context"
	"fmt"

	"github.com/barryx002/court-auction-api/internal/config"
	"github.com/barryx002/court-auction-api/internal/crawler"
	"github.com/barryx002/court-auction-api/internal/queue"
	"github.com/barryx002/court-auction-api/internal/repository"
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

	// 스키마 보장
	if err := db.EnsureSchema(ctx); err != nil {
		panic(err)
	}

	q := queue.NewQueue(cfg.RedisAddr)

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

		if len(items) > 0 {
			if err := db.UpsertAuctionItems(ctx, items); err != nil {
				fmt.Println("DB 저장 실패:", err)
				continue
			}
		}
		fmt.Printf("크롤링 완료 및 저장: %d 건\n", len(items))
	}
}
