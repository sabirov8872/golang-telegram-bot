package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
)

type DealerExist []string

// Option - Options for modification
type Option struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageSha    string `json:"imagesha"`
}

// Modification - Car modifications
type Modification struct {
	ModificationID string   `json:"modification_id"`
	Name           string   `json:"name"`
	Producing      string   `json:"producing"`
	Price          string   `json:"price"`
	Options        []string `json:"options"`
	OptionsObj     []Option `json:"options_obj"`
	Colors         []Color  `json:"colors"`
}

// StockData - Stock information for each region
type StockData struct {
	RegionID string `json:"region_id"`
	Stock    string `json:"stock"`
}

// Color - Color options for the car
type Color struct {
	ColorID     string      `json:"color_id"`
	Name        string      `json:"name"`
	HexValue    string      `json:"hex_value"`
	QueueNo     string      `json:"queue_no"`
	ExpectDate  string      `json:"expect_date"`
	PhotoSha    string      `json:"photo_sha"`
	StockData   []StockData `json:"stock_data"`
	PhotoSha666 string      `json:"photo_sha666"`
}

// CarModel - Main model structure
type CarModel struct {
	ModelID       string         `json:"model_id"`
	Name          string         `json:"name"`
	PhotoSha      string         `json:"photo_sha"`
	DealerExist   DealerExist    `json:"dealer_exist"`
	Modifications []Modification `json:"modifications"`
	As666         int            `json:"as666"`
	PhotoSha666   string         `json:"photo_sha666"`
}

// So'rov uchun struktura
type RequestData struct {
	IsWeb    string `json:"is_web"`
	FilialID int    `json:"filial_id"`
}

var bot *tgbotapi.BotAPI

func sendMessage(chatID int64, msg string) {
	msgConfig := tgbotapi.NewMessage(chatID, msg)
	bot.Send(msgConfig)
}

func main() {
	// .env faylidan tokenni o'qib olish
	godotenv.Load()
	token := os.Getenv("botToken")

	bot, _ = tgbotapi.NewBotAPI(token)
	fmt.Println("bot online")
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatId := update.Message.Chat.ID

		if update.Message.Text == "/start" {
			// URL manzilini aniqlash
			url := "https://savdo.uzavtosanoat.uz/t/ap/stream/ph&models"

			// So'rov ma'lumotlarini tayyorlash
			requestData := RequestData{
				IsWeb:    "Y",
				FilialID: 100,
			}

			// Ma'lumotlarni JSON formatiga o'tkazish
			jsonData, err := json.Marshal(requestData)
			if err != nil {
				log.Fatalf("Error marshalling JSON: %v", err)
			}

			// HTTP POST so'rovini tayyorlash
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
			if err != nil {
				log.Fatalf("Error creating request: %v", err)
			}

			// Header qo'shish
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("rcode", "SuiQuAy4qEkJhv51JS0n8jW1KeGxaxnr")

			// So'rovni amalga oshirish
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Fatalf("Error sending request: %v", err)
			}
			defer resp.Body.Close()

			// Javobni o'qish
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatalf("Error reading response body: %v", err)
			}

			var carModelList []CarModel

			if err := json.Unmarshal(body, &carModelList); err != nil {
				log.Fatalf("Error unmarshalling JSON: %v", err)
			}

			//fmt.Println("Response Status:", resp.Status)
			msg := ""

			for _, carModel := range carModelList {
				//fmt.Printf("Model ID: %s, Name: %s\n", carModel.ModelID, carModel.Name)
				msg += carModel.Name + "\n"
			}

			sendMessage(chatId, msg)
		}
	}
}
