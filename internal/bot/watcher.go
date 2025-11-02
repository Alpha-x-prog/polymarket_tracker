package bot

import (
	"context"
	"fmt"
	"log"
	"time"
)

// lastActivityByWallet — чтобы не слать одно и то же
var lastActivityByWallet = make(map[string]string)

func (b *Bot) startWatcher() {
	ticker := time.NewTicker(15 * time.Second)
	for range ticker.C {
		b.checkAllWallets()
	}
}

func (b *Bot) checkAllWallets() {
	// достаём из БД всех подписанных
	subs, err := b.store.GetAllSubs()
	if err != nil {
		log.Println("get subs:", err)
		return
	}

	ctx := context.Background()

	for chatID, wallets := range subs {
		for _, w := range wallets {
			b.checkWalletAndNotify(ctx, chatID, w)
		}
	}
}

func (b *Bot) checkWalletAndNotify(ctx context.Context, chatID int64, wallet string) {
	act, err := b.pm.GetUserLastActivity(ctx, wallet)
	if err != nil {
		log.Printf("Error checking wallet %s: %v", wallet, err)
		return
	}
	if act == nil {
		return
	}

	// у некоторых активностей может не быть id — надо подстраховаться
	id := act.ID
	if id == "" {
		// соберём суррогатный id
		id = act.Type + "|" + act.ConditionID + "|" + act.CreatedAt
	}

	// если уже слали это событие — выходим
	if lastActivityByWallet[wallet] == id {
		return
	}

	// запоминаем
	lastActivityByWallet[wallet] = id

	// можно фильтрануть только трейды
	if act.Type != "" && act.Type != "TRADE" {
		return
	}

	title := act.MarketTitle
	if title == "" {
		title = act.Title
	}
	if title == "" {
		title = act.ConditionID
	}

	msg := fmt.Sprintf("🟦 %s сделал %s в рынке: %s", wallet, act.Type, title)
	b.Send(chatID, msg)
}
