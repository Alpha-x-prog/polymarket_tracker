package bot

import (
	"context"
	"fmt"
	"time"
)

func (b *Bot) HandleWatchMarketQuery(chatID int64, query string) {
	if query == "" {
		b.Send(chatID, "Usage: /watch-market <name or slug>")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	markets, err := b.pm.SearchMarkets(ctx, query, 5)
	if err != nil {
		b.Send(chatID, fmt.Sprintf("❌ Error searching markets: %v", err))
		return
	}

	if len(markets) == 0 {
		b.Send(chatID, "No markets found")
		return
	}

	msg := "Found markets:\n\n"
	for i, m := range markets {
		msg += fmt.Sprintf("%d) %s", i+1, m.Question)
		if m.Category != "" {
			msg += fmt.Sprintf(" [%s]", m.Category)
		}
		msg += fmt.Sprintf("\n   id: %s\n\n", m.ID)
	}
	msg += "Send: /watch-market-id <id>"

	b.Send(chatID, msg)
}

func (b *Bot) HandleWatchMarketID(chatID int64, marketID string) {
	if marketID == "" {
		b.Send(chatID, "Usage: /watch-market-id <market_id>")
		return
	}

	// Optional: fetch and validate market exists
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	market, err := b.pm.GetMarketByID(ctx, marketID)
	if err != nil {
		b.Send(chatID, fmt.Sprintf("❌ Error fetching market: %v", err))
		return
	}

	if market == nil {
		b.Send(chatID, fmt.Sprintf("❌ Market with ID %s not found", marketID))
		return
	}

	// Add to storage
	if err := b.store.AddMarket(chatID, marketID); err != nil {
		b.Send(chatID, fmt.Sprintf("❌ Error: %v", err))
		return
	}

	b.Send(chatID, fmt.Sprintf("✅ Tracking market %s\n\n%s", marketID, market.Question))
}

func (b *Bot) HandleMarketsList(chatID int64) {
	markets, err := b.store.GetMarkets(chatID)
	if err != nil {
		b.Send(chatID, fmt.Sprintf("❌ Error: %v", err))
		return
	}

	if len(markets) == 0 {
		b.Send(chatID, "You don't track any markets")
		return
	}

	msg := "📋 Your tracked markets:\n\n"
	for i, mID := range markets {
		msg += fmt.Sprintf("%d. %s\n", i+1, mID)
	}
	b.Send(chatID, msg)
}

