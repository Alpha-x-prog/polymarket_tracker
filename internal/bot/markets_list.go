// internal/bot/markets_list.go
package bot

import (
	"context"
	"fmt"
	"time"
)

func (b *Bot) HandleTrackedMarkets(chatID int64) {
	ids, err := b.store.GetMarkets(chatID)
	if err != nil {
		b.Send(chatID, fmt.Sprintf("❌ Error: %v", err))
		return
	}
	if len(ids) == 0 {
		b.Send(chatID, "You don't track any markets.\nAdd: /track-market <text or slug> or /track-market-id <id>")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	msg := "📋 Your tracked markets:\n\n"
	// тянем названия последовательно (достаточно и так)
	for i, id := range ids {
		m, err := b.pm.GetMarketByID(ctx, id)
		if err != nil || m == nil {
			msg += fmt.Sprintf("%d) %s\n   (fetch failed)\n\n", i+1, id)
			continue
		}
		line := fmt.Sprintf("%d) %s\n   id: %s", i+1, m.Question, m.ID)
		if m.Category != "" {
			line += fmt.Sprintf("  [%s]", m.Category)
		}
		msg += line + "\n\n"
	}
	msg += "— View one: /market <id>\n— Stop tracking: /untrack-market-id <id> (если сделаешь команду)"

	b.Send(chatID, msg)
}
