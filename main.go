package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Welcome struct {
	ModelID       string         `json:"model_id"`
	Name          string         `json:"name"`
	PhotoSHA      string         `json:"photo_sha"`
	DealerExist   []string       `json:"dealer_exist"`
	Modifications []Modification `json:"modifications"`
	PhotoSha666   string         `json:"photo_sha666"`
}

type Modification struct {
	ModificationID string       `json:"modification_id"`
	Name           string       `json:"name"`
	Producing      string       `json:"producing"`
	Price          string       `json:"price"`
	Options        []string     `json:"options"`
	OptionsObj     []OptionsObj `json:"options_obj"`
	Colors         []Color      `json:"colors"`
}

type Color struct {
	ColorID     string       `json:"color_id"`
	Name        string       `json:"name"`
	HexValue    string       `json:"hex_value"`
	QueueNo     string       `json:"queue_no"`
	ExpectDate  string       `json:"expect_date"`
	PhotoSHA    string       `json:"photo_sha"`
	StockData   []StockDatum `json:"stock_data"`
	PhotoSha666 string       `json:"photo_sha666"`
}

type StockDatum struct {
	RegionID string `json:"region_id"`
	Stock    string `json:"stock"`
}

type OptionsObj struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Imagesha    string `json:"imagesha"`
}

type Client struct {
	UserID    int64
	FirstName string
	Subscribe bool
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

var db *sql.DB
var bot *tgbotapi.BotAPI

func isIdInTheTable(update tgbotapi.Update) bool {
	// Ma'lumotlarni o'qib olish
	sqlSelect := `SELECT chat_id FROM users`
	rows, err := db.Query(sqlSelect)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var client Client
		err = rows.Scan(&client.UserID)
		if err != nil {
			log.Fatal(err)
		}
		if client.UserID == update.Message.From.ID {
			return false
		}
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return true
}

func sendTimeMessage(msg string) {
	// Ma'lumotlarni o'qib olish
	sqlSelect := `SELECT chat_id, subscribe FROM users`
	rows, err := db.Query(sqlSelect)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var client Client
		err = rows.Scan(&client.UserID, &client.Subscribe)
		if err != nil {
			log.Fatal(err)
		}
		if client.Subscribe {
			msgConfig := tgbotapi.NewMessage(client.UserID, msg)
			bot.Send(msgConfig)
		}
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func scheduleAtTwoPM() {
	msg := ""
	cnt := 0

	for {
		if time.Now().Hour() == 22 && time.Now().Minute() == 22 {
			msgEveryDay := ""

			if cnt == 0 {
				msgEveryDay = "Сегодня не было новостей"
			} else {
				msgEveryDay = "Сегодня были новости"
				cnt = 0
			}

			sendTimeMessage(msgEveryDay)
		}

		time.Sleep(time.Minute)
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

		var welcome []Welcome
		if err = json.Unmarshal(body, &welcome); err != nil {
			fmt.Println("Error:", err)
			return
		}

		newMsg := ""

		for _, w := range welcome {
			newMsg += w.ModelID + " " + w.Name + "\n"
		}

		if msg == "" {
			msg = newMsg
		}

		if msg != newMsg {
			msg = newMsg
			sendTimeMessage(msg)
			cnt++
		}
	}
}

func main() {
	// Postgresga ulanish uchun DSN (Data Source Name) yaratish
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, _ = sql.Open("postgres", psqlInfo)
	defer db.Close()

	err := db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("connected to the database")

	// .env faylidan tokenni o'qib olish
	godotenv.Load()
	token := os.Getenv("botToken")

	bot, _ = tgbotapi.NewBotAPI(token)
	fmt.Println("bot online")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	go scheduleAtTwoPM() // Funksiyani fon jarayoni sifatida ishga tushirish

	for update := range updates {
		if update.Message != nil {
			switch update.Message.Text {
			case "/start":
				if isIdInTheTable(update) {
					sqlStatement := `
						INSERT INTO users (chat_id, first_name, subscribe)
						VALUES ($1, $2, $3)
						RETURNING id`
					id := 0
					err = db.QueryRow(sqlStatement, update.Message.Chat.ID, update.Message.From.FirstName, false).Scan(&id)
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
		} else if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err = bot.Request(callback); err != nil {
				log.Panic(err)
			}

			client := Client{
				UserID:    update.CallbackQuery.From.ID,
				FirstName: update.CallbackQuery.From.FirstName,
			}

			if update.CallbackQuery.Data == "Subscribe" {
				client.Subscribe = true
			} else {
				client.Subscribe = false
			}

			sqlUpdate := `
				UPDATE users
				SET subscribe = $1
				WHERE chat_id = $2`
			res, _ := db.Exec(sqlUpdate, client.Subscribe, client.UserID)

			_, err = res.RowsAffected()
			if err != nil {
				log.Fatal(err)
			}

			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Вы выбрали: "+update.CallbackQuery.Data)
			bot.Send(msg)
		}
	}
}
