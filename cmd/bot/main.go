package bot

import (
	"github.com/joho/godotenv"
	"log"
	"telegram_bot/internal/bot"
	"telegram_bot/internal/database"
)

func Run() {
	// .env faylidan o'qish
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Bazaga ulanish
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Botni ishga tushirish
	bot.Start(db)
}
