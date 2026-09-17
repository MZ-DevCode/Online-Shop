package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartBot(token string) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Println("Error to started tg bot: ", err)
		return
	}

	bot.Debug = true
	log.Printf("Бот запущен")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		var text string

		if update.Message.Text == "/start" {
			text = "Добро пожаловать в телеграм бота Online-Shop. Для помощи используйте /help"
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)

		_, err := bot.Send(msg)
		if err != nil {
			log.Println("Ошибка отправки сообщения: ", err)
		}
	}
}
