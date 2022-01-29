package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	user "project0/repository"
)

var msg tgbotapi.MessageConfig

func messageResolver(update tgbotapi.Update) tgbotapi.MessageConfig {
	res := user.GetOrCreateUser(update)
	username := res.Username

	newMessage := update.Message.Text
	fmt.Println(username, newMessage)

	if username == "Пися" && newMessage == "👤👔\nПрофиль" || newMessage == "📝 Изменить имя обратно? 📝" {
		switch newMessage {
		case "👤👔\nПрофиль":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "\bПрофиль:\b\nТы жмых из под пня, и звать тебя "+res.Username+"!\nНо я буду звать ультра-мышь!")
			msg.ReplyMarkup = profileKeyboardBackUsername
		case "📝 Изменить имя обратно? 📝":
			res := user.UpdateUsername(update, true)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "\bПрофиль:\b\nТы жмых из под пня, и звать тебя "+res.Username+"!\nНо я буду звать ультра-мышь!")
			msg.ReplyMarkup = profileKeyboard
		}
	} else {
		switch newMessage {
		case "/start":
			res := user.GetOrCreateUser(update)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Приветствую тебя,  "+res.Username)
			user.GetOrCreateLocation(update)
		case "/menu":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Меню")
			msg.ReplyMarkup = mainKeyboard
		case "🗺\nКарта":
			res := user.GetOrCreateLocation(update)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Карта: "+res.Maps+"\n🏔⛰🗻⬜️\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\U0001F7E9\U0001F7E9\U0001F7E9🚪\n⛰🗻⬜️⬜️\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9⛪️\U0001F7E8🏪\U0001F7E9\n☃️⬜️⬜️⬜️\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\n⬜️⬜️⬜️🔥\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\n\U0001F7E9\U0001F7E9\U0001FAB5\U0001F7E9\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9🏥\U0001F7E8🏦\U0001F7E9\n\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\n\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8🕦\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\n\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E8\n\U0001FAA8\U0001FAA8🐚\U0001F7E9\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\U0001F7E9\U0001F7E9\U0001F7E9🍄\n\U0001F7E6\U0001F7E6\U0001F7E6\U0001F7E6\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\U0001F7EB🍅\U0001F7EB🥔\n\U0001F7E6\U0001F7E6\U0001F7E6\U0001F7E6\U0001F7E9\U0001F7E8🐱\U0001F7E8\U0001F7E9🥕\U0001F7EB🌽\U0001F7EB\n\U0001F9CA\U0001F9CA\U0001F7E6\U0001F7E6\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9\U0001F7EB🍎\U0001F7EB🍓\n\U0001F9CA\U0001F9CA\U0001F7E6\U0001F7E6\U0001F7E9\U0001F7E8\U0001F7E8\U0001F7E8\U0001F7E9🌳🌿🌱🌵")
			msg.ReplyMarkup = moveKeyboard

		case "👤👔\nПрофиль":
			res := user.GetOrCreateUser(update)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "\bПрофиль:\b\nТы жмых из под пня, и звать тебя "+res.Username+"!\nНо я буду звать ультра-мышь!")
			msg.ReplyMarkup = profileKeyboard
		case "📝 Изменить имя на Писю? 📝":
			user.UpdateUsername(update, false)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Эх, Пися, для меня ты все равно ультра-мышь :D")
			msg.ReplyMarkup = profileKeyboardBackUsername
		case "👜\nИнвентарь":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Инвентарь")
			msg.ReplyMarkup = backpackKeyboard
		default:
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Сам ты "+newMessage)
		}
	}

	return msg
}
