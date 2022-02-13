package repository

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/config"
	"project0/helpers"
)

type Item struct {
	ID            uint   `gorm:"primaryKey"`
	Name          string `gorm:"embedded"`
	View          string `gorm:"embedded"`
	Type          string `gorm:"embedded"`
	CanTake       bool   `gorm:"embedded"`
	CanTakeWith   *Item  `gorm:"foreignKey:CanTakeWithId"`
	CanTakeWithId *int   `gorm:"embedded"`
	Healing       *int   `gorm:"embedded"`
	Damage        *int   `gorm:"embedded"`
	Satiety       *int   `gorm:"embedded"`
	Cost          *int   `gorm:"embedded"`
}

func UserGetItem(update tgbotapi.Update, LocationStruct Location, char []string) string {
	var resultCell Cellule
	var err error

	err = config.Db.
		Preload("Item").
		First(&resultCell, &Cellule{Map: LocationStruct.Map, AxisX: *LocationStruct.AxisX, AxisY: *LocationStruct.AxisY}).
		Error
	if err != nil {
		panic(err)
	}

	if resultCell.ItemID != nil {
		switch resultCell.Item.Type {
		case "food", "pick":
			res := UserGetItemUpdateModels(update, resultCell)
			if res != "Ok" {
				return res
			}
		}
	} else {
		return "0"
	}
	return "Ты взял " + char[2] + " 1шт.\nВ ячейке: " + helpers.ToString(*resultCell.CountItem-1) + " шт."
}

func UserGetItemUpdateModels(update tgbotapi.Update, resultCell Cellule) string {
	countAfterUserGetItem := *resultCell.CountItem - 1
	user := GetUser(User{TgId: uint(update.Message.From.ID)})

	if *user.Money >= *resultCell.Item.Cost {
		resUserItem := GetOrCreateUserItem(update, *resultCell.Item)

		if canUserTakeItem(resUserItem) {
			updateUserMoney := *user.Money - *resultCell.Item.Cost
			AddUserItemCount(update, resUserItem, resultCell, updateUserMoney)
			UpdateCellule(resultCell.ID, Cellule{CountItem: &countAfterUserGetItem})
			return "Ok"
		}
		return "У тебя уже есть такой!"
	}
	return "Не хватает деняк!"
}

func canUserTakeItem(item UserItem) bool {
	if item.Item.Type == "pick" && *item.Count < 1 {
		return true
	} else if item.Item.Type == "food" {
		return true
	}
	return false
}
