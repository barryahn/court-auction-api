package crawler

import (
	"fmt"

	"github.com/barryx002/court-auction-api/internal/models"

	"github.com/gocolly/colly"
)

func CrawlAuctionList(url string) ([]models.AuctionItem, error) {
    c := colly.NewCollector()

    items := []models.AuctionItem{}

    c.OnHTML("table.Ltbl_list tr", func(e *colly.HTMLElement) {
        cols := e.DOM.Find("td")
        if cols.Length() > 0 {
            item := models.AuctionItem{
                CaseNumber:  cols.Eq(0).Text(),
                Court:       cols.Eq(1).Text(),
                AuctionDate: cols.Eq(2).Text(),
                Address:     cols.Eq(3).Text(),
                ItemType:    cols.Eq(4).Text(),
                MinBidPrice: cols.Eq(5).Text(),
            }
            items = append(items, item)
        }
    })

    err := c.Visit(url)
    if err != nil {
        return nil, err
    }

    fmt.Printf("크롤링 완료: %d건\n", len(items))
    return items, nil
}
