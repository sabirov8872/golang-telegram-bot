package bot

import (
	"database/sql"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"log"
	"net/http"
	"time"
)

type Model struct {
	ModelID string `json:"model_id"`
	Name    string `json:"name"`
}

type User struct {
	UserID    int64
	FirstName string
	Subscribe bool
}

func HandleMessage(update tgbotapi.Update, db *sql.DB, bot *tgbotapi.BotAPI) {
	switch update.Message.Text {
	case "/start":
		if isIdInTheTable(update, db) {
			sqlStatement := `
				INSERT INTO users (chat_id, first_name, subscribe)
				VALUES ($1, $2, $3)
				RETURNING id`
			id := 0
			err := db.QueryRow(sqlStatement, update.Message.Chat.ID, update.Message.From.FirstName, false).Scan(&id)
			if err != nil {
				log.Fatal(err)
			}
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет. Я могу рассказать вам новости на этом сайте https://savdo.uzavtosanoat.uz. Выберите Subscribe или Unsubscribe")
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Subscribe", "Subscribe"),
				tgbotapi.NewInlineKeyboardButtonData("Unsubscribe", "Unsubscribe"),
			),
		)
		bot.Send(msg)
	}
}

func HandleCallback(update tgbotapi.Update, db *sql.DB, bot *tgbotapi.BotAPI) {
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
	if _, err := bot.Request(callback); err != nil {
		log.Panic(err)
	}

	user := User{
		UserID:    update.CallbackQuery.From.ID,
		FirstName: update.CallbackQuery.From.FirstName,
	}

	if update.CallbackQuery.Data == "Subscribe" {
		user.Subscribe = true
	} else {
		user.Subscribe = false
	}

	sqlUpdate := `
		UPDATE users
		SET subscribe = $1
		WHERE chat_id = $2`
	res, _ := db.Exec(sqlUpdate, user.Subscribe, user.UserID)

	_, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Вы выбрали: "+update.CallbackQuery.Data)
	bot.Send(msg)
}

func TimeMessage(db *sql.DB, bot *tgbotapi.BotAPI) {
	msg := ""
	cnt := 0

	for {
		if time.Now().Hour() == 11 && time.Now().Minute() == 12 {
			msgEveryDay := ""

			if cnt == 0 {
				msgEveryDay = "Сегодня не было новостей"
			} else {
				msgEveryDay = "Сегодня были новости"
				cnt = 0
			}

			sendTimeMessage(msgEveryDay, db, bot)
		}

		time.Sleep(time.Minute)
		getRequest(msg, db, &cnt, bot)
	}
}

func sendTimeMessage(msg string, db *sql.DB, bot *tgbotapi.BotAPI) {
	// Ma'lumotlarni o'qib olish
	sqlSelect := `SELECT chat_id, subscribe FROM users`
	rows, err := db.Query(sqlSelect)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err = rows.Scan(&user.UserID, &user.Subscribe)
		if err != nil {
			log.Fatal(err)
		}
		if user.Subscribe {
			msgConfig := tgbotapi.NewMessage(user.UserID, msg)
			bot.Send(msgConfig)
		}
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func isIdInTheTable(update tgbotapi.Update, db *sql.DB) bool {
	// Ma'lumotlarni o'qib olish
	sqlSelect := `SELECT chat_id FROM users`
	rows, err := db.Query(sqlSelect)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err = rows.Scan(&user.UserID)
		if err != nil {
			log.Fatal(err)
		}
		if user.UserID == update.Message.From.ID {
			return false
		}
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return true
}

func getRequest(msg string, db *sql.DB, cnt *int, bot *tgbotapi.BotAPI) {
	url := "https://savdo.uzavtosanoat.uz/models.json?1721365266"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var models []Model
	if err = json.Unmarshal(body, &models); err != nil {
		fmt.Println("Error:", err)
		return
	}

	newMsg := ""

	for _, model := range models {
		newMsg += model.ModelID + " " + model.Name + "\n"
	}

	if msg == "" {
		msg = newMsg
	}

	if msg != newMsg {
		msg = newMsg
		sendTimeMessage(msg, db, bot)
		*cnt++
	}
}
