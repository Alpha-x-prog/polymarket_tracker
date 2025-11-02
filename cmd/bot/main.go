package main

import (
	"log"
	"os"
	"polymarket_tg_bot/internal/bot"
	"polymarket_tg_bot/internal/storage"

	"github.com/joho/godotenv"
)

func mustEnv(key string) string {

	_ = godotenv.Load(".env")
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("environment variable %s is not set", key)
	}
	return v
}

func main() {
	token := mustEnv("TG_BOT_TOKEN")

	st, err := storage.New("polymarket_bot.db")
	if err != nil {
		log.Fatalf("init storage: %v", err)
	}
	defer st.Close()

	b, err := bot.New(token, st)
	if err != nil {
		log.Fatalf("init bot: %v", err)
	}

	if err := b.Run(); err != nil {
		log.Fatalf("run bot: %v", err)
	}
}
