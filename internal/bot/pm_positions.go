package bot

import (
	"context"
	"fmt"
	"polymarket_tg_bot/internal/polymarket"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// сколько позиций показываем за раз
const positionsPerPage = 5

// /pm positions <0x...>
func (b *Bot) handlePMPositions(chatID int64, addr string) {
	if addr == "" {
		// пытаемся взять из БД
		saved, err := b.store.GetDefaultWallet(chatID)
		if err != nil {
			b.Send(chatID, "Ошибка чтения кошелька: "+err.Error())
			return
		}
		if saved == "" {
			b.Send(chatID, "У тебя не сохранён кошелёк. Сделай: /setwallet 0x....")
			return
		}
		addr = saved
	}

	ctx := context.Background()
	positions, err := b.pm.GetUserPositions(ctx, addr)
	if err != nil {
		b.Send(chatID, "❌ error: "+err.Error())
		return
	}
	if len(positions) == 0 {
		b.Send(chatID, "Нет открытых позиций у этого пользователя.")
		return
	}

	// положили в кэш
	if b.positionsCache == nil {
		b.positionsCache = make(map[int64][]polymarket.UserPosition)
	}
	b.positionsCache[chatID] = positions

	// показали первую страницу
	b.sendPositionsPage(chatID, addr, 0)
}

func (b *Bot) handlePMValue(chatID int64, addr string) {
	// 1. если адрес не передали в команду – пробуем взять из БД
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

	// 2. запрашиваем у Polymarket total value
	ctx := context.Background()
	value, err := b.pm.GetUserTotalValue(ctx, addr)
	if err != nil {
		b.Send(chatID, "Не смог получить value у Polymarket: "+err.Error())
		return
	}

	// 3. красиво выводим
	b.Send(chatID, fmt.Sprintf("💰 Total value for %s: %.2f USD", addr, value))
}

func (b *Bot) sendPositionsPage(chatID int64, addr string, page int) {
	positions, ok := b.positionsCache[chatID]
	if !ok || len(positions) == 0 {
		b.Send(chatID, "Пока нет позиций в кэше. Сначала вызови /pm positions <addr>.")
		return
	}

	total := len(positions)
	maxPage := (total - 1) / positionsPerPage
	if page < 0 {
		page = 0
	}
	if page > maxPage {
		page = maxPage
	}

	start := page * positionsPerPage
	end := start + positionsPerPage
	if end > total {
		end = total
	}
	pagePositions := positions[start:end]

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📊 Positions for %s\n", addr))
	sb.WriteString(fmt.Sprintf("Page %d/%d\n\n", page+1, maxPage+1))

	for _, p := range pagePositions {
		title := p.Title
		if title == "" {
			title = p.Slug
		}
		if title == "" {
			title = p.ConditionID
		}

		flags := ""
		if p.Redeemable {
			flags += " ✅redeem"
		}
		if p.Mergeable {
			flags += " 🧩merge"
		}
		if p.NegativeRisk {
			flags += " ⚠️neg"
		}

		sb.WriteString(fmt.Sprintf(
			"• %s\n  outcome: %s\n  cur price: %.3f\n  value: %.2f$  pnl: %.2f$\n%s\n\n",
			title,
			p.Outcome,
			p.CurPrice,
			p.CurrentValue,
			p.CashPnL,
			flags,
		))
	}

	kb := positionsKeyboard(addr, page, maxPage)

	msg := tgbotapi.NewMessage(chatID, sb.String())
	msg.ReplyMarkup = kb
	b.api.Send(msg)
}

// редактирование при листании
func (b *Bot) editPositionsPage(msg *tgbotapi.Message, addr string, page int) {
	positions, ok := b.positionsCache[msg.Chat.ID]
	if !ok || len(positions) == 0 {
		return
	}

	total := len(positions)
	maxPage := (total - 1) / positionsPerPage
	if page < 0 {
		page = 0
	}
	if page > maxPage {
		page = maxPage
	}

	start := page * positionsPerPage
	end := start + positionsPerPage
	if end > total {
		end = total
	}
	pagePositions := positions[start:end]

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📊 Positions for %s\n", addr))
	sb.WriteString(fmt.Sprintf("Page %d/%d\n\n", page+1, maxPage+1))

	for _, p := range pagePositions {
		title := p.Title
		if title == "" {
			title = p.Slug
		}
		if title == "" {
			title = p.ConditionID
		}

		flags := ""
		if p.Redeemable {
			flags += " ✅redeem"
		}
		if p.Mergeable {
			flags += " 🧩merge"
		}
		if p.NegativeRisk {
			flags += " ⚠️neg"
		}

		sb.WriteString(fmt.Sprintf(
			"• %s\n  outcome: %s\n  cur price: %.3f\n  value: %.2f$  pnl: %.2f$\n%s\n\n",
			title,
			p.Outcome,
			p.CurPrice,
			p.CurrentValue,
			p.CashPnL,
			flags,
		))
	}

	kb := positionsKeyboard(addr, page, maxPage)

	edit := tgbotapi.NewEditMessageTextAndMarkup(
		msg.Chat.ID,
		msg.MessageID,
		sb.String(),
		kb,
	)
	b.api.Send(edit)
}

func positionsKeyboard(addr string, page, maxPage int) tgbotapi.InlineKeyboardMarkup {
	prevPage := page - 1
	if prevPage < 0 {
		prevPage = 0
	}
	nextPage := page + 1
	if nextPage > maxPage {
		nextPage = maxPage
	}

	row := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("« 1", fmt.Sprintf("pos:0:%s", addr)),
		tgbotapi.NewInlineKeyboardButtonData("‹", fmt.Sprintf("pos:%d:%s", prevPage, addr)),
		tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d/%d", page+1, maxPage+1), "noop"),
		tgbotapi.NewInlineKeyboardButtonData("›", fmt.Sprintf("pos:%d:%s", nextPage, addr)),
		tgbotapi.NewInlineKeyboardButtonData("»", fmt.Sprintf("pos:%d:%s", maxPage, addr)),
	}

	return tgbotapi.NewInlineKeyboardMarkup(row)
}
