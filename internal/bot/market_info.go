// internal/bot/market_info.go
package bot

import (
	"context"
	"fmt"
	"strings"
	"time"
)

func (b *Bot) HandleMarketInfo(chatID int64, id string) {
	id = strings.TrimSpace(id)
	if id == "" {
		b.Send(chatID, "Usage: /market <condition_id>")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	m, err := b.pm.GetMarketByID(ctx, id)
	if err != nil {
		b.Send(chatID, "❌ Error fetching market: "+err.Error())
		return
	}
	if m == nil {
		b.Send(chatID, "❌ Market not found")
		return
	}

	// исходы (если есть)
	outcomes := m.Outcomes.Names()
	outStr := "n/a"
	if len(outcomes) > 0 {
		if len(outcomes) > 4 { // чтобы не раздувать сообщение
			outcomes = outcomes[:4]
		}
		outStr = strings.Join(outcomes, " | ")
	}

	text := fmt.Sprintf(
		"🧾 Market info\n\n"+
			"Title: %s\n"+
			"ID: %s\n"+
			"Slug: %s\n"+
			"Category: %s\n"+
			"Outcomes: %s\n"+
			"24h Volume: $%.2f\n"+
			"Open Interest: $%.2f\n",
		m.Question, m.ID, m.Slug, m.Category, outStr, m.Volume24, m.OI,
	)

	b.Send(chatID, text)
}
