package services

import (
	"errors"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v "github.com/spf13/viper"
	r "project0/repository"
	"strings"
	"time"
)

func StringToArrayCharacters(word string) []rune {
	chars := []rune(word)
	return chars
}

func CheckWordEquals(userWord string) string {
	var resultString string
	activeWord, _ := r.GetActiveWord()

	userArray := StringToArrayCharacters(userWord)
	secretArray := StringToArrayCharacters(activeWord.SecretWord)
	for i, char := range secretArray {
		if char != userArray[i] {
			if strings.ContainsRune(activeWord.SecretWord, userArray[i]) {
				if strings.Count(activeWord.SecretWord, string(userArray[i])) < strings.Count(userWord, string(userArray[i])) {
					resultString += getCharColor(userArray, secretArray, i)
				} else {
					resultString += "🟨"
				}
			} else {
				resultString += "⬜️"
			}
		} else {
			resultString += "✅"
		}
	}

	return resultString
}

func getCharColor(userArray []rune, secretArray []rune, i int) string {
	var secretIndexArray []int
	for x, char := range secretArray {
		if char == userArray[i] {
			secretIndexArray = append(secretIndexArray, x)
		}
	}

	var userIndexArray []int
	for x, char := range userArray {
		if char == userArray[i] {
			userIndexArray = append(userIndexArray, x)
		}
	}

	var Y []int
	var X []int

	for y, sIndex := range secretIndexArray {
		for x, uIndex := range userIndexArray {
			if sIndex == uIndex {
				Y = append(Y, y)
				X = append(X, x)
			}
		}
	}

	resultSecretArray := secretIndexArray
	resultUserArray := userIndexArray

	for y := len(Y) - 1; y >= 0; y-- {
		resultSecretArray = RemoveIndex(resultSecretArray, y)
	}

	for x := len(Y) - 1; x >= 0; x-- {
		resultUserArray = RemoveIndex(resultUserArray, x)
	}

	for y, resultIndex := range resultUserArray {
		if resultIndex == i {
			if y+1 <= len(resultSecretArray) {
				return "🟨"
			}
		}
	}

	return "⬜️"
}

func RemoveIndex(s []int, index int) []int {
	return append(s[:index], s[index+1:]...)
}

func FormattedUserWord(userWord string) string {

	res := strings.ToLower(userWord)
	res = strings.TrimSpace(res)

	return res

}

func WordleMenuButtons() tg.InlineKeyboardMarkup {
	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("📊 Статистика", "wordleUserStatistic"),
		),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("📚 Правила", "wordleRegulations"),
			tg.NewInlineKeyboardButtonData("⚠️ Выйти", "cancel"),
		),
	)
}

func buttonStatistic() tg.InlineKeyboardMarkup {
	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("📊 Статистика", "wordleUserStatistic"),
		),
	)
}

func WordleMap(user r.User) (string, tg.InlineKeyboardMarkup) {
	var msgText string
	var btns tg.InlineKeyboardMarkup

	msgText += "\U0001F9EE *Игра Вуордле!*\U0001F9EE\n"
	countTries := 6

	_, err := r.GetActiveWord()
	if err != nil {
		msgText = fmt.Sprintf("%s\n\n_Соре, сегодня нет слова_ \U0001F97A \n\n_Приходи завтра, мб уже будет...)_", msgText)
		btns = buttonStatistic()

		return msgText, btns
	}

	game := r.GetOrCreateWordleGameProcess(user)
	words := r.GetUserWords(user, time.Now())

	for i, word := range words {
		if i < countTries {
			row := CheckWordEquals(word.Word)
			msgText += fmt.Sprintf("\n%s - *%s*", row, strings.ToUpper(word.Word))
		}
	}

	if len(words) < 6 {
		for x := len(words); x < countTries; x++ {
			msgText += "\n\U0001F7E6\U0001F7E6\U0001F7E6\U0001F7E6\U0001F7E6"
			if x == len(words) && game.Status == "new" {
				msgText += " - _Напиши слово!_"
			}
		}
	}

	var lastText string
	if game.Status == "new" && game.CountTries < countTries {
		lastText = "Только 5 букв! 👉🤚 Ни больше, ни меньше! 👌"
	} else if game.Status == "win" {
		lastText = "🏆 Поздравляю, сегодня ты выйграл! 🏆"
	} else if game.Status == "lose" {
		lastText = "☠️ Ты проиграл :C Ну ничего, попробуй завтра еще раз! 👋"
	}

	btns = WordleMenuButtons()
	msgText = fmt.Sprintf("%s%s_%s_", msgText, v.GetString("msg_separator"), lastText)

	return msgText, btns
}

func CheckUserWordFormat(user r.User, game r.WordleGameProcess, userWord string) (string, error) {
	var msgText string

	if game.Status != "new" {
		msgText = "\U0001FAC2 Игра уже окончена! Приходи завтра) 🤝"
		return msgText, errors.New("game ended")
	}

	if len(strings.Fields(userWord)) != 1 {
		msgText = "Некорректное количество слов"
		return msgText, errors.New("too many words")
	}

	userWord = FormattedUserWord(userWord)

	if chars := []rune(userWord); len(chars) > 5 {
		msgText = "‼️ Слишком много букв ‼️"
		return msgText, errors.New("too many chars")
	} else if len(chars) < 5 {
		msgText = "‼️ Слишком мало букв ‼️"
		return msgText, errors.New("not enough chars")
	}

	if !IsDictionaryHasWord(userWord) {
		msgText = "‼️ Я не нашел в словаре такое слово! Не придумывай)) ‼️"
		return msgText, errors.New("is not word")
	}

	words := r.GetUserWords(user, time.Now())
	for _, word := range words {
		if word.Word == userWord {
			msgText = "‼️ Такое слово уже было ‼️"
			return msgText, errors.New("word duplicate")
		}
	}

	return msgText, nil
}

func UserSendNextWord(user r.User, newMessage string) (string, tg.InlineKeyboardMarkup) {
	game := r.GetOrCreateWordleGameProcess(user)

	msgText, err := CheckUserWordFormat(user, game, newMessage)
	if err != nil {
		msg, btns := WordleMap(user)
		msgText = fmt.Sprintf("%s\n%s", msg, msgText)
		return msgText, btns
	}

	word := FormattedUserWord(newMessage)
	r.CreateUserWord(user, word)
	activeWord, _ := r.GetActiveWord()

	if word == activeWord.SecretWord {
		game.Status = "win"
		game.CountTries++
	}

	game.UpdateWordleGameProcess(user)

	return WordleMap(user)
}
