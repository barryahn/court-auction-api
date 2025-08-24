package repository

import (
	"context"

	"github.com/barryx002/court-auction-api/internal/models"
)

// EnsureSchema creates required tables and indexes if they do not exist.
func (db *DB) EnsureSchema(ctx context.Context) error {
    _, err := db.Pool.Exec(ctx, `
        CREATE TABLE IF NOT EXISTS auctions (
            id BIGSERIAL PRIMARY KEY,
            case_number TEXT NOT NULL,
            court TEXT NOT NULL,
            auction_date TEXT,
            address TEXT,
            item_type TEXT,
            min_bid_price TEXT,
            crawled_at TIMESTAMPTZ DEFAULT NOW(),
            UNIQUE(case_number, court)
        );

        CREATE INDEX IF NOT EXISTS idx_auctions_crawled_at ON auctions (crawled_at DESC);
    `)
    return err
}

// UpsertAuctionItems inserts or updates multiple auction items.
func (db *DB) UpsertAuctionItems(ctx context.Context, items []models.AuctionItem) error {
    batch := &Batch{}
    batch.init()

    for _, it := range items {
        batch.queue(`
            INSERT INTO auctions (case_number, court, auction_date, address, item_type, min_bid_price)
            VALUES ($1, $2, $3, $4, $5, $6)
            ON CONFLICT (case_number, court)
            DO UPDATE SET
                auction_date = EXCLUDED.auction_date,
                address = EXCLUDED.address,
                item_type = EXCLUDED.item_type,
                min_bid_price = EXCLUDED.min_bid_price,
                crawled_at = NOW();
        `, it.CaseNumber, it.Court, it.AuctionDate, it.Address, it.ItemType, it.MinBidPrice)
    }

    return batch.send(ctx, db)
}

// ListAuctions returns a page of auctions.
func (db *DB) ListAuctions(ctx context.Context, limit, offset int) ([]models.AuctionItem, error) {
    rows, err := db.Pool.Query(ctx, `
        SELECT case_number, court, auction_date, address, item_type, min_bid_price
        FROM auctions
        ORDER BY crawled_at DESC
        LIMIT $1 OFFSET $2
    `, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    list := make([]models.AuctionItem, 0)
    for rows.Next() {
        var it models.AuctionItem
        if err := rows.Scan(&it.CaseNumber, &it.Court, &it.AuctionDate, &it.Address, &it.ItemType, &it.MinBidPrice); err != nil {
            return nil, err
        }
        list = append(list, it)
    }
    return list, rows.Err()
}

// --- lightweight batch helper using pgxpool ---

type Batch struct {
    statements []stmt
}

type stmt struct {
    sql  string
    args []any
}

func (b *Batch) init() { b.statements = make([]stmt, 0) }

func (b *Batch) queue(sql string, args ...any) {
    b.statements = append(b.statements, stmt{sql: sql, args: args})
}

func (b *Batch) send(ctx context.Context, db *DB) error {
    if len(b.statements) == 0 {
        return nil
    }
    tx, err := db.Pool.Begin(ctx)
    if err != nil {
        return err
    }
    defer func() { _ = tx.Rollback(ctx) }()

    for _, s := range b.statements {
        if _, err := tx.Exec(ctx, s.sql, s.args...); err != nil {
            return err
        }
    }
    return tx.Commit(ctx)
}


