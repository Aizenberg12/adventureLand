package mapsActions

import (
	"errors"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/controllers/questController"
	"project0/src/models"
	"project0/src/services/helpers"
)

func CheckQuestActions(user models.User, charData []string) (msg string, buttons tg.InlineKeyboardMarkup, err error) {
	// Квесты
	switch charData[0] {
	case "quests":
		msg, buttons = listOfQuests(user)
	case "dailyQuests":
		msg, buttons = listOfDailyQuests(user)
	case v.GetString("callback_char.quest"):
		msg, buttons = questController.OpenQuest(uint(helpers.ToInt(charData[1])), user)
	case v.GetString("callback_char.user_get_quest"):
		msg, buttons = userGetQuest(user, charData)
	case v.GetString("callback_char.user_done_quest"):
		msg, buttons = questController.UserDoneQuest(uint(helpers.ToInt(charData[1])), user)
	default:
		err = errors.New("not quest actions")
	}

	return msg, buttons, err
}

func userGetQuest(user models.User, charData []string) (msg string, buttons tg.InlineKeyboardMarkup) {
	models.UserQuest{
		UserId:  user.ID,
		QuestId: uint(helpers.ToInt(charData[1])),
	}.CreateOrUpdateUserQuest()
	msg, buttons = questController.OpenQuest(uint(helpers.ToInt(charData[1])), user)
	return msg, buttons
}

func listOfQuests(user models.User) (msg string, buttons tg.InlineKeyboardMarkup) {
	msg = v.GetString("user_location.tasks_menu_message")
	buttons = questController.AllQuestsMessageKeyboard(user, false)
	return msg, buttons
}

func listOfDailyQuests(user models.User) (msg string, buttons tg.InlineKeyboardMarkup) {
	msg = v.GetString("user_location.daily_tasks_menu_message")
	buttons = questController.AllQuestsMessageKeyboard(user, true)
	return msg, buttons
}
