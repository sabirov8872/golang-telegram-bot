package bot

import (
	"database/sql"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Start(db *sql.DB) {
	token := os.Getenv("botToken")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("%s online", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	go TimeMessage(db, bot)

	for update := range updates {
		if update.Message != nil {
			HandleMessage(update, db, bot)
		} else if update.CallbackQuery != nil {
			HandleCallback(update, db, bot)
		}
	}
}
