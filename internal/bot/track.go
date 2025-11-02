package bot

import (
	"fmt"
)

func (b *Bot) handleStart(chatID int64) {
	msg := "👋 Welcome to Polymarket Wallet Tracker!\n\n" +
		"Commands:\n" +
		"/track <wallet> - Subscribe to a wallet address\n" +
		"/track-list - List your subscribed wallets\n" +
		"/track-remove <wallet> - Unsubscribe from a wallet"
	b.Send(chatID, msg)
}

func (b *Bot) handleTrack(chatID int64, wallet string) {
	if err := b.store.AddWallet(chatID, wallet); err != nil {
		b.Send(chatID, fmt.Sprintf("❌ Error: %v", err))
		return
	}
	b.Send(chatID, fmt.Sprintf("✅ Now tracking wallet: %s", wallet))
}

func (b *Bot) handleTrackList(chatID int64) {
	wallets, err := b.store.GetWallets(chatID)
	if err != nil {
		b.Send(chatID, fmt.Sprintf("❌ Error: %v", err))
		return
	}

	if len(wallets) == 0 {
		b.Send(chatID, "📋 No wallets subscribed. Use /track <wallet> to add one.")
		return
	}

	msg := "📋 Your subscribed wallets:\n\n"
	for i, wallet := range wallets {
		msg += fmt.Sprintf("%d. %s\n", i+1, wallet)
	}
	b.Send(chatID, msg)
}

func (b *Bot) handleTrackRemove(chatID int64, wallet string) {
	if err := b.store.RemoveWallet(chatID, wallet); err != nil {
		b.Send(chatID, fmt.Sprintf("❌ Error: %v", err))
		return
	}
	b.Send(chatID, fmt.Sprintf("✅ Removed wallet: %s", wallet))
}

func (b *Bot) registerHandlers() {
	// Handlers are registered in router.go Dispatch method
	// This method exists for future extensibility
}

