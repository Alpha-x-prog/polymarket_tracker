package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api *tgbotapi.BotAPI
}

func New(token string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	b := &Bot{
		api: api,
	}

	return b, nil
}

func (b *Bot) Run() error {
	log.Printf("Authorized on @%s", b.api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates := b.api.GetUpdatesChan(u)

	for upd := range updates {
		if upd.Message == nil {
			continue
		}
		b.send(upd.Message.Chat.ID, "✅ Bot ready. Polymarket integration coming soon.")
	}
	return nil
}
