package bot

import (
	"github.com/mounlion/markdownwatcher/parsing"
	"fmt"
	"github.com/mounlion/markdownwatcher/database"
	"gopkg.in/telegram-bot-api.v4"
	"log"
)

const DNSDomain = "https://www.dns-shop.ru"
var BotToken *string

func SetBotToken(value *string)  {
	BotToken = value
}

func SendCatalog(newItems []parsing.Item, updateItems[]parsing.Item)  {
	if len(newItems) > 0 || len(updateItems) > 0 {

		var (
			newItemsString string
			updateItemsString string
		)

		if len(newItems) > 0 {
			newItemsString += "<b>Новые товары</b>\n\n"
			newItemsString += CatalogMessage(newItems)
		}

		if len(updateItems) > 0 {
			updateItemsString += "<b>Обновление цен</b>\n\n"
			updateItemsString += CatalogMessage(updateItems)
		}

		bot, err := tgbotapi.NewBotAPI(*BotToken)
		if err != nil {
			log.Panic(err)
		}

		users := database.GetUsers()

		for _, user := range users {
			if user.IsActive {
				if len(newItemsString) > 0 {
					sendCatalogMessage(bot, &user, &newItemsString)
				}
				if len(updateItemsString) > 0 {
					sendCatalogMessage(bot, &user, &updateItemsString)
				}
			}
		}
	}
}

func sendCatalogMessage(bot *tgbotapi.BotAPI, user *database.User, message *string)  {
	msg := tgbotapi.NewMessage(user.Id, *message)
	msg.ParseMode = "HTML"
	msg.DisableWebPagePreview = true
	_, err := bot.Send(msg)
	if err != nil {
		if err.Error() == "Forbidden: bot was blocked by the user" {
			database.Subscribe(int(user.Id), false)
		}
	}
}

func SendServiceMessage(text string)  {
	users := database.GetUsers()

	bot, err := tgbotapi.NewBotAPI(*BotToken)
	if err != nil {
		log.Panic(err)
	}

	for _, user := range users {
		if user.IsActive && user.IsAdmin {
			msg := tgbotapi.NewMessage(user.Id, text)
			msg.ParseMode = "HTML"
			_, err := bot.Send(msg)
			if err != nil {
				if err.Error() == "Forbidden: bot was blocked by the user" {
					database.Subscribe(int(user.Id), false)
				}
			}
		}
	}
}

func CatalogMessage(items []parsing.Item) string {
	var catalog string

	for _, val := range items {
		catalog += fmt.Sprintf("<a href=\"%s%s\">%s</a>\n", DNSDomain, val.Url, val.Title)
		catalog += fmt.Sprintf("<b>%d₽</b>", val.Price)
		if val.OldPrice != 0 {
			catalog += fmt.Sprintf("    <code>%d₽</code>", val.OldPrice)
		}
		if len(val.Desc) > 0 {
			catalog += fmt.Sprintf("<i>%s</i>", val.Desc)
		}
		catalog += "\n\n"
	}

	return catalog
}