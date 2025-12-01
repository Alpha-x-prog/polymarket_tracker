package bot

import (
	"context"
	"fmt"
	"strings"
	"time"
)

func (b *Bot) HandleMarketInfo(chatID int64, arg string) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	slug := strings.TrimSpace(arg)
	if slug == "" {
		b.Send(chatID, "Usage: /market <slug>")
		return
	}

	m, err := b.pm.GetMarketBySlug(ctx, slug)
	if err != nil || m == nil {
		b.Send(chatID, "❌ Market not found")
		return
	}
	// выбираем категорию
	cat := m.Category
	if cat == "" && len(m.Categories) > 0 && m.Categories[0].Label != "" {
		cat = m.Categories[0].Label
	}
	if cat == "" {
		cat = "—"
	}

	desc := strings.TrimSpace(m.Description)
	if desc == "" {
		desc = "No description"
	} else if len(desc) > 400 {
		desc = desc[:397] + "..."
	}

	res := strings.TrimSpace(m.ResolutionSource)

	spread := m.Spread
	if spread == 0 && m.BestBid > 0 && m.BestAsk > 0 {
		spread = m.BestAsk - m.BestBid
	}

	var resLine string
	if res != "" {
		resLine = fmt.Sprintf("Resolution: %s\n", res)
	}

	msg := fmt.Sprintf(
		"📌 %s\n"+
			"Category: %s\n"+
			"Slug: %s\n"+
			"ID: %s\n"+
			"24h Volume: %.2f USDC\n"+
			"Total Volume: %.2f USDC\n"+
			"Liquidity: %.2f USDC\n"+
			"Spread: %.4f (bid %.4f / ask %.4f)\n"+
			"%s"+
			"\n%s\n\n"+
			"https://polymarket.com/event/%s",
		m.Question,
		cat,
		m.Slug,
		m.ID,
		m.Volume24hr,
		m.VolumeNum,
		m.LiquidityNum,
		spread,
		m.BestBid,
		m.BestAsk,
		resLine,
		desc,
		m.Slug,
	)

	b.Send(chatID, msg)

}
