package bot

import (
	"log"
	"polymarket_tg_bot/internal/polymarket"
	"polymarket_tg_bot/internal/storage"
	"strconv"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api            *tgbotapi.BotAPI
	pm             *polymarket.Client
	r              *Router
	store          *storage.Storage
	mu             sync.Mutex
	lastSeen       map[string]string                   // wallet -> last activity ID
	positionsCache map[int64][]polymarket.UserPosition // chatID -> positions
}

func New(token string, store *storage.Storage) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	b := &Bot{
		api:            api,
		pm:             polymarket.NewClient(),
		r:              NewRouter(),
		store:          store,
		lastSeen:       make(map[string]string),
		positionsCache: make(map[int64][]polymarket.UserPosition),
	}

	b.registerHandlers()

	// Start background watchers
	go b.startWatcher()       // for user wallets
	go b.startMarketWatcher() // for markets by ID

	return b, nil
}

func (b *Bot) Run() error {
	log.Printf("Authorized on @%s", b.api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates := b.api.GetUpdatesChan(u)

	for upd := range updates {
		switch {
		// обычные сообщения
		case upd.Message != nil:
			b.r.Dispatch(b, upd.Message.Chat.ID, upd.Message.Text)

		// нажатия на inline-кнопки
		case upd.CallbackQuery != nil:
			b.handleCallback(upd.CallbackQuery)
		}
	}
	return nil
}

func (b *Bot) Send(chatID int64, text string) {
	b.send(chatID, text)
}

func (b *Bot) handleCallback(q *tgbotapi.CallbackQuery) {
	data := q.Data

	// pos:<page>:<addr>
	if strings.HasPrefix(data, "pos:") {
		parts := strings.SplitN(data, ":", 3)
		if len(parts) == 3 {
			pageStr := parts[1]
			addr := parts[2]

			page, err := strconv.Atoi(pageStr)
			if err == nil {
				// редактируем сообщение
				b.editPositionsPage(q.Message, addr, page)
			}
		}
	}

	// обязательно ответить на callback, чтобы "часики" исчезли
	b.api.Request(tgbotapi.NewCallback(q.ID, ""))
}
