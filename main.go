package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math/rand"
	"strconv"
)

const BotToken = "7368281324:AAGEPlq6znDSXDkEJ-penPuq4nKTfX5RJRg"

var bot *tgbotapi.BotAPI
var symbols = [5]string{"+", "-", "*", "/", "%"}

func questions() (string, string) {
	symbol := rand.Intn(len(symbols))
	rnd1 := rand.Intn(21)
	rnd2 := rand.Intn(21)
	sum := 0

	if rnd2 == 0 {
		rnd2++
	}

	if symbols[symbol] == "+" {
		sum = rnd1 + rnd2
	} else if symbols[symbol] == "-" {
		sum = rnd1 - rnd2
	} else if symbols[symbol] == "*" {
		sum = rnd1 * rnd2
	} else if symbols[symbol] == "/" {
		sum = rnd1 / rnd2
	} else {
		sum = rnd1 % rnd2
	}

	return strconv.Itoa(rnd1) + " " + symbols[symbol] + " " + strconv.Itoa(rnd2) + " = ?", strconv.Itoa(sum)
}

func sendMessage(chatID int64, msg string) {
	msgConfig := tgbotapi.NewMessage(chatID, msg)
	bot.Send(msgConfig)
}

func main() {
	bot, _ = tgbotapi.NewBotAPI(BotToken)
	fmt.Println("Bot ishlayapti!!")

	updateConfig := tgbotapi.NewUpdate(0)
	updates := bot.GetUpdatesChan(updateConfig)

	isStart := false
	cnt := 0
	isQuestion := false
	isTrueQuestion := false
	rightAnswer := 0
	var answer, question string

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatId := update.Message.Chat.ID

		if isTrueQuestion {
			natija := ""

			if update.Message.Text == answer {
				natija = "to'g'ri ✅"
				rightAnswer++
			} else {
				natija = "noto'g'ri ❌\nTo'g'ri javob✅:  " + answer
			}

			sendMessage(chatId, "Javob "+natija)
			isTrueQuestion = false
		}

		if cnt == 10 {
			sendMessage(chatId, "To'gri javoblar✅:  "+strconv.Itoa(rightAnswer)+"\nNoto'gri javoblar❌:  "+strconv.Itoa(10-rightAnswer))
			sendMessage(chatId, "Sizga hisoblashga oid 10 ta savol beriladi va oxirida sizning natijangiz aytiladi.\n\n/ - sonning butun qismini oling\n% - sonning qoldiq qismini oling\n\nMasalan:\n7 / 3 = 2\n7 % 3 = 1\n\nBoshlash uchun boshlashni ustidan bosing.\n👉  /boshlash  👈")
			isStart = false
			cnt = 0
			isQuestion = false
			isTrueQuestion = false
			rightAnswer = 0
		}

		if !isQuestion && update.Message.Text == "/boshlash" {
			isStart = true
			isQuestion = true
		}

		if !isTrueQuestion && isQuestion {
			cnt++
			question, answer = questions()
			sendMessage(chatId, strconv.Itoa(cnt)+" - savol:\n"+question)
			isTrueQuestion = true
		}

		if !isStart && update.Message.Text == "/start" {
			sendMessage(chatId, "Salom "+update.Message.From.FirstName+". Mening ismim my_bot. Men hozircha sizga hisoblashga oid savollar bera olaman!!")
			sendMessage(chatId, "Sizga hisoblashga oid 10 ta savol beriladi va oxirida sizning natijangiz aytiladi.\n\n/ - sonning butun qismini oling\n% - sonning qoldiq qismini oling\n\nMasalan:\n7 / 3 = 2\n7 % 3 = 1\n\nBoshlash uchun boshlashni ustidan bosing.\n👉  /boshlash  👈")
		}
	}
}
