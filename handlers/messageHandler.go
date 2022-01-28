package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	user "project0/repository"
)

var msg tgbotapi.MessageConfig

func messageResolver(update tgbotapi.Update) tgbotapi.MessageConfig {

	newMessage := update.Message.Text

	switch newMessage {
	case "🗺\nКарта":
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "🏔⛰🗻⬜️\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\U0001F7E9\U0001F7E9\U0001F7E9🚪\n⛰🗻⬜️⬜️\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9⛪️\U0001F7E8🏪\U0001F7E9\n☃️⬜️⬜️⬜️\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\n⬜️⬜️⬜️🔥\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\n\U0001F7E9\U0001F7E9\U0001FAB5\U0001F7E9\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9🏥\U0001F7E8🏦\U0001F7E9\n\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\n\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8🕦\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\n\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\n\U0001FAA8\U0001FAA8🐚\U0001F7E9\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\U0001F7E9\U0001F7E9\U0001F7E9🍄\n\U0001F7E6\U0001F7E6\U0001F7E6\U0001F7E6\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\U0001F7EB🍅\U0001F7EB🥔\n\U0001F7E6\U0001F7E6\U0001F7E6\U0001F7E6\U0001F7E9\U0001F7E8🐱\U0001F7E8\U0001F7E9🥕\U0001F7EB🌽\U0001F7EB\n\U0001F9CA\U0001F9CA\U0001F7E6\U0001F7E6\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\U0001F7EB🍎\U0001F7EB🍓\n\U0001F9CA\U0001F9CA\U0001F7E6\U0001F7E6\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9🌳🌿🌱🌵")
		msg.ReplyMarkup = moveKeyboard
	case "/start":
		res := user.GetOrCreateUser(update)
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Приветствую тебя,  "+res.Username)
	case "/menu":
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Меню")
		msg.ReplyMarkup = mainKeyboard
	case "👤👔\nПрофиль":
		res := user.GetOrCreateUser(update)
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "\bПрофиль:\b\nТы жмых из под пня, и звать тебя "+res.Username+"!\nНо я буду звать ультра-мышь!")
		msg.ReplyMarkup = profileKeyboard
	case "📝 Изменить имя на Писю? 📝":
		res := user.UpdateUsername(update)
		fmt.Println(res)
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Введите новое имя..., если не хочешь быть ультра мышью)")
		msg.ReplyMarkup = profileKeyboard
	case "👜\nИнвентарь":
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Инвентарь")
		msg.ReplyMarkup = backpackKeyboard
	default:
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Сам ты "+newMessage)
	}

	return msg
}
