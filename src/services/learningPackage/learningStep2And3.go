package learningPackage

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	"project0/src/actions/mapsActions"
	"project0/src/actions/mapsActions/userGetBoxAction"
	"project0/src/controllers/userMapController"
	"project0/src/models"
	"project0/src/repositories"
	"project0/src/services/helpers"
	"strings"
)

func learningStep2And3(data string, user models.User) (text string, buttons tg.InlineKeyboardMarkup) {
	charData := strings.Fields(data)

	info1 := "*Шаг 2:*\nЗдесь ты научишься добывать ресурсы!\nИсследуй местность, не бойся нажимать на кнопки и не забудь взять подарочки 🎁 сверху, там ты получишь инструменты!"
	info2 := "Отлично!\nВозьми второй подарок 🎁, и я объясню тебе, что с этим делать!"
	infoNextStep := "*Шаг 3:*\nПоздравляю! Ты получил новые предметы!\nЧтобы использовать их, нажми на кнопку «Вещи 🧦»"

	switch true {
	case strings.Contains(data, "goodsMoving"),
		strings.Contains(data, "Меню"),
		strings.Contains(data, "category"),
		strings.Contains(data, "move 44208"):
		text, buttons = userMapController.GetMyMap(user)
		if strings.Contains(user.MenuLocation, "step2") {
			text = fmt.Sprintf("%s%s%s%s❗️ Пока еще рано это нажимать 🤫", info1, v.GetString("msg_separator"), text, v.GetString("msg_separator"))
		} else {
			text = fmt.Sprintf("%s%s%s%s❗️ Пока еще рано это нажимать 🤫", info2, v.GetString("msg_separator"), text, v.GetString("msg_separator"))
		}

	case strings.Contains(data, "box"):
		cell := models.Cell{ID: uint(helpers.ToInt(charData[1]))}.GetCell()
		text, buttons = userGetBoxAction.UserGetBox(user, cell)

		if strings.Contains(user.MenuLocation, "step2") {
			text = fmt.Sprintf("%s%s%s", info2, v.GetString("msg_separator"), text)
			user.MenuLocation = "learning step3"
			repositories.UpdateUser(user)
		} else {
			text = fmt.Sprintf("%s%s%s", infoNextStep, v.GetString("msg_separator"), text)
			user.MenuLocation = "learning step4"
			repositories.UpdateUser(user)
		}

	default:
		text, buttons = mapsActions.MapsActions(user, data)
		if strings.Contains(user.MenuLocation, "step2") {
			text = fmt.Sprintf("%s%s%s", info1, v.GetString("msg_separator"), text)
		} else {
			text = fmt.Sprintf("%s%s%s", info2, v.GetString("msg_separator"), text)
		}
	}

	return text, buttons
}
