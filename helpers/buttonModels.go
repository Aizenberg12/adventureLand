package helpers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/repository"
	"strings"
)

func BackpackInlineKeyboard(items []repository.UserItem, i int) tgbotapi.InlineKeyboardMarkup {
	if len(items) == 0 {
		return tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Пусто...(", "emptyBackPack"),
			),
		)
	}
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(items[i].Item.View+" "+
				repository.ToString(*items[i].Count)+"шт."+
				"   +"+repository.ToString(*items[i].Item.Healing)+" ♥️️"+
				"   +"+repository.ToString(*items[i].Item.Satiety)+"\U0001F9C3", "description "+repository.ToString(items[i].ID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🍽 1шт", "eatFood "+repository.ToString(items[i].ID)+" "+repository.ToString(i)),
			tgbotapi.NewInlineKeyboardButtonData("🔺", "backpackMoving "+repository.ToString(i-1)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🗑 все!", "throwOutFood "+repository.ToString(items[i].ID)+" "+repository.ToString(i)),
			tgbotapi.NewInlineKeyboardButtonData("🔻", "backpackMoving "+repository.ToString(i+1)),
		),
	)
}

func GoodsInlineKeyboard(user repository.User, userItems []repository.UserItem, i int) tgbotapi.InlineKeyboardMarkup {
	if len(userItems) == 0 {
		return tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Пусто...(", "emptyGoods"),
			),
		)
	}

	text, data := repository.IsDressedItem(user, userItems[i])
	itemDesciption := "Описания нет("
	if userItems[i].Item.Description != nil {
		itemDesciption = *userItems[i].Item.Description
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(userItems[i].Item.View+" "+repository.ToString(*userItems[i].Count)+"шт. "+userItems[i].Item.Name+"  "+itemDesciption,
				"description "+repository.ToString(userItems[i].ID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(text, data+" "+repository.ToString(userItems[i].ID)+" "+repository.ToString(i)),
			tgbotapi.NewInlineKeyboardButtonData("🔺", "goodsMoving "+repository.ToString(i-1)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🗑", "throwOutGood "+repository.ToString(userItems[i].ID)+" "+repository.ToString(i)),
			tgbotapi.NewInlineKeyboardButtonData("🔻", "goodsMoving "+repository.ToString(i+1)),
		),
	)
}

func EmodjiInlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	var buttons [][]tgbotapi.InlineKeyboardButton
	var listOfAvatar []string
	listOfAvatar = strings.Fields("🐶 🐱 🐭 🐹 🐰 🦊 🐻 🐼 ‍️🐨 🐯 🦁 🐮 🐷 🐸 🐵 🐦 🐧 🐔 🐤 🐥 🦆 🐴 🦄 🐺 🐗 🐝 🦋 🐛 🐌 🐞 🪲 🪰 🐜 🕷 🪳 🦖 🦕 🐙 🦀 🐟 🐠 🐡 🦭")

	for x := 0; x < len(listOfAvatar); x = x + 8 {
		var row []tgbotapi.InlineKeyboardButton
		for i := 0; i < 8; i++ {
			sum := x + i
			if len(listOfAvatar) > sum {
				row = append(row, tgbotapi.NewInlineKeyboardButtonData(listOfAvatar[sum], "changeAvatar "+listOfAvatar[sum]))
			}
		}
		buttons = append(buttons, row)
	}

	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

func ProfileKeyboard(user repository.User) tgbotapi.ReplyKeyboardMarkup {
	var onlineButton string
	if *user.OnlineMap {
		onlineButton = "Онлайн (📳♻️📴)"
	} else {
		onlineButton = "Офлайн (📴♻️📳)"
	}
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📝 Изменить имя? 📝"),
			tgbotapi.NewKeyboardButton(user.Avatar+" Изменить аватар? "+user.Avatar),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(onlineButton),
			tgbotapi.NewKeyboardButton("Меню"),
		),
	)
}

func MainKeyboard(user repository.User) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🗺 Карта 🗺"),
			tgbotapi.NewKeyboardButton(user.Avatar+" Профиль 👔"),
		),
	)
}

func ChooseInstrument(char []string, cell repository.Cellule, user repository.User) tgbotapi.InlineKeyboardMarkup {
	instruments := repository.GetInstrumentsUserCanUse(user, cell)

	if len(instruments) != 0 {
		var row []tgbotapi.InlineKeyboardButton

		for _, instrument := range instruments {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(instrument, instrument+" "+char[3]+" "+char[4]))
		}

		return tgbotapi.NewInlineKeyboardMarkup(
			row,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Отмена", "cancel"),
			),
		)
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("На карту?", "cancel"),
		),
	)
}
