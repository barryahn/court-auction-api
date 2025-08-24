package models

type AuctionItem struct {
    CaseNumber   string `json:"case_number"`
    Court        string `json:"court"`
    AuctionDate  string `json:"auction_date"`
    Address      string `json:"address"`
    ItemType     string `json:"item_type"`
    MinBidPrice  string `json:"min_bid_price"`
}
