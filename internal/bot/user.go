package bot

import (
	"context"
	"fmt"
	"strings"
)

func (b *Bot) handleUserProfile(chatID int64, addr string) {
	// 1. адрес: если не передали — берём сохранённый из БД
	if addr == "" {
		saved, err := b.store.GetDefaultWallet(chatID)
		if err != nil {
			b.Send(chatID, "Ошибка чтения кошелька: "+err.Error())
			return
		}
		if saved == "" {
			b.Send(chatID, "Кошелёк не задан. Сначала сделай: /setwallet 0x...")
			return
		}
		addr = saved
	}

	ctx := context.Background()

	// 2. тянем данные с Polymarket
	value, err := b.pm.GetUserTotalValue(ctx, addr)
	if err != nil {
		b.Send(chatID, "Не смог получить total value: "+err.Error())
		return
	}

	openPos, _ := b.pm.GetUserPositions(ctx, addr)
	closedPos, _ := b.pm.GetUserClosedPositions(ctx, addr, 50)
	traded, _ := b.pm.GetUserTraded(ctx, addr)

	// 3. посчитаем PnL и biggest win
	var realizedTotal float64
	var biggestWin float64
	var biggestWinTitle string

	for _, cp := range closedPos {
		realizedTotal += cp.RealizedPnL
		if cp.RealizedPnL > biggestWin {
			biggestWin = cp.RealizedPnL
			if cp.Title != "" {
				biggestWinTitle = cp.Title
			} else if cp.Slug != "" {
				biggestWinTitle = cp.Slug
			} else {
				biggestWinTitle = cp.ConditionID
			}
		}
	}

	// нереализованный — по открытым
	var unrealized float64
	for _, op := range openPos {
		unrealized += op.CurrentValue - op.InitialValue
	}

	// 4. собираем ответ
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("👤 User: %s\n", addr))
	sb.WriteString(fmt.Sprintf("Positions value: $%.2f\n", value))
	sb.WriteString(fmt.Sprintf("Open positions: %d\n", len(openPos)))
	sb.WriteString(fmt.Sprintf("Closed positions: %d\n", len(closedPos)))
	sb.WriteString(fmt.Sprintf("Predictions (markets touched): %d\n", traded))
	sb.WriteString(fmt.Sprintf("Realized PnL: $%.2f\n", realizedTotal))
	sb.WriteString(fmt.Sprintf("Unrealized PnL: $%.2f\n", unrealized))

	if biggestWin > 0 {
		sb.WriteString(fmt.Sprintf("Biggest win: $%.2f — %s\n", biggestWin, biggestWinTitle))
	}

	b.Send(chatID, sb.String())
}
