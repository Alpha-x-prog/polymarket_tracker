package bot

import (
	"strings"
)

type Router struct{}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Dispatch(b *Bot, chatID int64, text string) {
	text = strings.TrimSpace(text)
	if text == "" {
		b.Send(chatID, "Unknown command")
		return
	}

	switch {
	case text == "/start":
		b.handleStart(chatID)
	case strings.HasPrefix(text, "/track "):
		wallet := strings.TrimSpace(strings.TrimPrefix(text, "/track"))
		if wallet == "" {
			b.Send(chatID, "Usage: /track <wallet_address>")
			return
		}
		b.handleTrack(chatID, wallet)
	case text == "/track-list":
		b.handleTrackList(chatID)
	case strings.HasPrefix(text, "/track-remove "):
		wallet := strings.TrimSpace(strings.TrimPrefix(text, "/track-remove"))
		if wallet == "" {
			b.Send(chatID, "Usage: /track-remove <wallet_address>")
			return
		}
		b.handleTrackRemove(chatID, wallet)
	case strings.HasPrefix(text, "/watch-market "):
		query := strings.TrimSpace(strings.TrimPrefix(text, "/watch-market"))
		b.HandleWatchMarketQuery(chatID, query)
	case strings.HasPrefix(text, "/watch-market-id "):
		id := strings.TrimSpace(strings.TrimPrefix(text, "/watch-market-id"))
		b.HandleWatchMarketID(chatID, id)
	case text == "/markets-list":
		b.HandleMarketsList(chatID)
	case strings.HasPrefix(text, "/pm positions "):
		addr := strings.TrimSpace(strings.TrimPrefix(text, "/pm positions "))
		b.handlePMPositions(chatID, addr)
	case strings.HasPrefix(text, "/pm value "):
		addr := strings.TrimSpace(strings.TrimPrefix(text, "/pm value "))
		b.handlePMPositions(chatID, addr)
	default:
		b.Send(chatID, "Unknown command")
	}
}
