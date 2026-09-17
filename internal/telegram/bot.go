package telegram

import (
	"WEBSITE/internal/database"
	"context"
	"fmt"
	"log"
	"time"

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
		chatId := update.Message.Chat.ID

		switch update.Message.Text {
		case "/start":
			text = "Добро пожаловать в телеграм бота Online-Shop. Для помощи используйте /help"
		case "/help":
			text = `Доступные команды:\n
				/start - Начать работу\n
				/catalog - Каталог товаров"`
		case "/catalog":
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			rows, err := database.DB.QueryContext(ctx, "SELECT id, name, price, stock FROM products")
			if err != nil {
				log.Println("Ошибка получения товаров из БД для бота: ", err)
				text = "Не удалось загрузить каталог"
				break
			}

			text = "Каталог товаров:\n"

			for rows.Next() {
				var (
					id    int
					name  string
					price float64
					stock int
				)

				if err := rows.Scan(&id, &name, &price, &stock); err != nil {
					log.Println("Ошибка чтения товара: ", err)
					break
				}

				text += fmt.Sprintf("🔹 <b>%s</b>\n Цена: <code>%.2f</code> \n В наличии: %d шт.\n", name, price, stock)
			}
			rows.Close()

		default:
			text = "Такой команды нет. Используйте /help"
		}

		msg := tgbotapi.NewMessage(chatId, text)
		msg.ParseMode = "HTML"

		_, err := bot.Send(msg)
		if err != nil {
			log.Println("Ошибка отправки сообщения: ", err)
		}
	}
}
