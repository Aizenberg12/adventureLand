package repository

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/config"
	"strings"
	"time"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	TgId         uint   `gorm:"embedded"`
	TgChatId     uint   `gorm:"embedded"`
	Username     string `gorm:"embedded"`
	Avatar       string `gorm:"embedded"`
	FirstName    string `gorm:"embedded"`
	LastName     string `gorm:"embedded"`
	Health       uint   `gorm:"embedded"`
	Satiety      uint   `gorm:"embedded"`
	Money        *int   `gorm:"embedded"`
	Head         *Item
	HeadId       *int
	LeftHand     *Item
	LeftHandId   *int
	RightHand    *Item
	RightHandId  *int
	Body         *Item
	BodyId       *int
	Foot         *Item
	FootId       *int
	Shoes        *Item
	ShoesId      *int
	MenuLocation string    `gorm:"embedded"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	Deleted      bool      `gorm:"embedded"`
	OnlineMap    *bool     `gorm:"embedded"`
}

func GetOrCreateUser(update tgbotapi.Update) User {
	userTgId := GetUserTgId(update)
	MoneyUserStart := 10
	UserOnline := false

	replacer := strings.NewReplacer("_", " ", "*", " ")
	var outUsername string
	outUsername = replacer.Replace(update.Message.From.UserName)

	result := User{
		TgId:      userTgId,
		TgChatId:  uint(update.Message.Chat.ID),
		Username:  outUsername,
		FirstName: update.Message.From.FirstName,
		LastName:  update.Message.From.LastName,
		Avatar:    "👤",
		Satiety:   100,
		Health:    100,
		Money:     &MoneyUserStart,
		OnlineMap: &UserOnline,
	}
	err := config.Db.
		Preload("Head").
		Preload("RightHand").
		Preload("LeftHand").
		Preload("Body").
		Preload("Foot").
		Preload("Shoes").
		Where(&User{TgId: userTgId}).
		FirstOrCreate(&result).
		Error

	if err != nil {
		panic(err)
	}

	return result
}

func GetUser(user User) User {
	var result User
	err := config.Db.
		Preload("Head").
		Preload("RightHand").
		Preload("LeftHand").
		Preload("Body").
		Preload("Foot").
		Preload("Shoes").
		Where(user).
		First(&result).Error
	if err != nil {
		panic(err)
	}

	return result
}

func UpdateUser(update tgbotapi.Update, UserStruct User) User {
	var err error
	userTgId := GetUserTgId(update)
	err = config.Db.Where(&User{TgId: userTgId}).Updates(UserStruct).Error
	if err != nil {
		panic(err)
	}

	res := GetUser(User{TgId: userTgId})
	return res
}

func SetNullUserField(update tgbotapi.Update, queryFeild string) {
	var err error
	userTgId := GetUserTgId(update)
	err = config.Db.Model(&User{}).Where(&User{TgId: userTgId}).Update(queryFeild, nil).Error

	if err != nil {
		panic(err)
	}
}

func GetUserInfo(update tgbotapi.Update) string {
	userTgId := GetUserTgId(update)
	user := GetUser(User{TgId: userTgId})
	userIsOnline := "📳 Вкл"

	if !*user.OnlineMap {
		userIsOnline = "📴 Откл"
	}

	messageMap := fmt.Sprintf("🔅 🔆 *Профиль* 🔆 🔅\n\n*Твое имя*: %s\n*Аватар*: %s\n*Золото*: %d 💰\n*Здоровье*: _%d_ ❤️\n*Сытость*: _%d_ 😋️\n*Онлайн*: _%s_",
		user.Username, user.Avatar, *user.Money, user.Health, user.Satiety, userIsOnline)

	return messageMap
}

func IsDressedItem(user User, userItem UserItem) (string, string) {
	dressItem := "Надеть ✅"
	dressItemData := v.GetString("callback_char.dress_good")

	if user.HeadId != nil && userItem.ItemId == *user.HeadId ||
		user.LeftHandId != nil && userItem.ItemId == *user.LeftHandId ||
		user.RightHandId != nil && userItem.ItemId == *user.RightHandId ||
		user.BodyId != nil && userItem.ItemId == *user.BodyId ||
		user.FootId != nil && userItem.ItemId == *user.FootId ||
		user.ShoesId != nil && userItem.ItemId == *user.ShoesId {

		dressItem = "Снять ❎"
		dressItemData = v.GetString("callback_char.take_off_good")
	}

	return dressItem, dressItemData
}

func CheckUserHasInstrument(user User, instrument Instrument) (string, Item) {
	if instrument.Type == "hand" {
		return "Ok", *instrument.Good
	}
	if user.LeftHandId != nil && *user.LeftHandId == *instrument.GoodId {
		return "Ok", *user.LeftHand
	}
	if user.RightHandId != nil && *user.RightHandId == *instrument.GoodId {
		return "Ok", *user.RightHand
	}
	return "User dont have instrument", Item{}
}

func CheckUserHasLighter(update tgbotapi.Update, user User) string {
	if user.LeftHandId != nil && user.LeftHand.Type == "light" {
		_, res := UpdateUserInstrument(update, user, *user.LeftHand)
		return res
	}
	if user.RightHandId != nil && user.RightHand.Type == "light" {
		_, res := UpdateUserInstrument(update, user, *user.RightHand)
		return res
	}
	return "Ok"
}
