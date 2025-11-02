package bot

import (
	"context"
	"fmt"
	"log"
	"math"
	"strings"
	"sync"
	"time"
)

func (b *Bot) startMarketWatcher() {
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	// remember last price per (chatID, marketID)
	var mu sync.Mutex
	last := make(map[string]float64)

	for range ticker.C {
		subs, err := b.store.GetAllMarkets()
		if err != nil {
			log.Printf("Error getting market subscriptions: %v", err)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		for chatID, marketIDs := range subs {
			for _, mID := range marketIDs {
				m, err := b.pm.GetMarketByID(ctx, mID)
				if err != nil || m == nil || len(m.Outcomes) == 0 {
					continue
				}

				// pick YES price (first outcome or outcome with Name == "Yes")
				var price float64
				for _, o := range m.Outcomes {
					if strings.EqualFold(o.Name, "Yes") {
						price = o.Price
						break
					}
				}
				if price == 0 {
					// fallback: first outcome
					price = m.Outcomes[0].Price
				}

				key := fmt.Sprintf("%d:%s", chatID, mID)

				mu.Lock()
				prev, ok := last[key]
				if !ok {
					last[key] = price
					mu.Unlock()
					continue
				}

				// notify only on change
				if math.Abs(price-prev) >= 0.01 {
					last[key] = price
					mu.Unlock()

					// human-friendly message
					msg := fmt.Sprintf(
						"📈 Market %s price changed: %.2f → %.2f\n\n%s",
						mID[:8]+"...", prev, price, m.Question,
					)
					b.Send(chatID, msg)
				} else {
					mu.Unlock()
				}
			}
		}

		cancel()
	}
}

