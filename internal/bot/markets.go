// internal/bot/markets.go
package bot

import (
	"context"
	"fmt"
	"polymarket_tg_bot/internal/polymarket"
	"time"
)

// internal/bot/markets.go
func (b *Bot) HandleTrackMarketQuery(chatID int64, query string) {
	if query == "" {
		b.Send(chatID, "Usage: /track-market <name or slug>")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// если пришёл slug — сразу точный фетч
	if polymarket.LooksLikeSlug(query) { // сделай экспортируемым, если хочешь; либо скопируй логику сюда
		m, err := b.pm.GetMarketBySlug(ctx, query)
		if err != nil || m == nil {
			b.Send(chatID, "Not found by slug. Try text search.")
			return
		}
		b.Send(chatID, fmt.Sprintf("Found market:\n%s\nid: %s\n\nUse: /track-market-id %s", m.Question, m.ID, m.ID))
		return
	}

	// иначе текстовый поиск через /public-search
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
		line := fmt.Sprintf("%d) %s", i+1, m.Question)
		if m.Category != "" {
			line += fmt.Sprintf(" [%s]", m.Category)
		}
		msg += line + fmt.Sprintf("\n   id: %s\n\n", m.ID)
	}
	msg += "Send: /track-market-id <id>"
	b.Send(chatID, msg)
}

func (b *Bot) HandleTrackMarketID(chatID int64, marketID string) {
	if marketID == "" {
		b.Send(chatID, "Usage: /track-market-id <market_id>")
		return
	}

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
