package repository

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/config"
	"strconv"
	"strings"
)

type Map struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"embedded"`
	SizeX int    `gorm:"embedded"`
	SizeY int    `gorm:"embedded"`
}

type MapButtons struct {
	Up    string
	Left  string
	Right string
	Down  string
}

func DefaultButtons() MapButtons {
	return MapButtons{
		Up:    "🔼",
		Left:  "◀️",
		Right: "▶️",
		Down:  "🔽",
	}
}

type UserMap struct {
	leftIndent  int
	rightIndent int
	upperIndent int
	downIndent  int
}

var displayMapSize = 5

func DefaultUserMap(location Location) UserMap {
	return UserMap{
		leftIndent:  location.AxisX - displayMapSize,
		rightIndent: location.AxisX + displayMapSize,
		upperIndent: location.AxisY + displayMapSize,
		downIndent:  location.AxisY - displayMapSize,
	}
}

func GetMap(update tgbotapi.Update) Map {
	resLocation := GetOrCreateMyLocation(update)
	result := Map{}

	err := config.Db.Where(&Map{Name: resLocation.Map}).FirstOrCreate(&result).Error

	if err != nil {
		panic(err)
	}

	return result
}

func GetMyMap(update tgbotapi.Update) (textMessage string, buttons MapButtons) {
	resUser := GetOrCreateUser(update)
	resLocation := GetOrCreateMyLocation(update)
	resMap := GetMap(update)
	buttons = DefaultButtons()
	mapSize := CalculateUserMapBorder(resLocation, resMap)

	var result []Cellule

	err := config.Db.Where("map = '" + resLocation.Map + "' and axis_x >= " + ToString(mapSize.leftIndent) + " and axis_x <= " + ToString(mapSize.rightIndent) + " and axis_y >= " + ToString(mapSize.downIndent) + " and axis_y <= " + ToString(mapSize.upperIndent)).Order("axis_x ASC").Order("axis_y ASC").Find(&result).Error
	if err != nil {
		panic(err)
	}

	var Maps [][]string
	Maps = make([][]string, resMap.SizeY+1)

	type Point = [2]int
	m := map[Point]Cellule{}

	for _, cell := range result {
		m[Point{cell.AxisX, cell.AxisY}] = cell
	}

	if cell := m[Point{resLocation.AxisX, resLocation.AxisY + 1}]; !cell.CanStep {
		if cell.View == "" {
			buttons.Up = "🚫"
		} else {
			buttons.Up = cell.View
		}
	}
	if cell := m[Point{resLocation.AxisX, resLocation.AxisY - 1}]; !cell.CanStep {
		if cell.View == "" {
			buttons.Down = "🚫"
		} else {
			buttons.Down = cell.View
		}
	}
	if cell := m[Point{resLocation.AxisX + 1, resLocation.AxisY}]; !cell.CanStep {
		if cell.View == "" {
			buttons.Right = "🚫"
		} else {
			buttons.Right = cell.View
		}
	}
	if cell := m[Point{resLocation.AxisX - 1, resLocation.AxisY}]; !cell.CanStep {
		if cell.View == "" {
			buttons.Left = "🚫"
		} else {
			buttons.Left = cell.View
		}
	}

	m[Point{resLocation.AxisX, resLocation.AxisY}] = Cellule{View: resUser.Avatar}

	for i := range Maps {
		for x := mapSize.leftIndent; x <= mapSize.rightIndent; x++ {
			if m[Point{x, i}].ID != 0 || m[Point{x, i}] == m[Point{resLocation.AxisX, resLocation.AxisY}] {
				Maps[i] = append(Maps[i], m[Point{x, i}].View)
			} else {
				Maps[i] = append(Maps[i], "\U0001FAA8")
			}
		}
	}

	messageMap := "*Карта*: _" + resLocation.Map + "_ *X*: _" + ToString(resLocation.AxisX) + "_  *Y*: _" + ToString(resLocation.AxisY) + "_"

	for i, row := range Maps {
		if i >= mapSize.downIndent && i <= mapSize.upperIndent {
			messageMap = strings.Join(row, ``) + "\n" + messageMap
		}
	}

	return messageMap, buttons
}

func ToString(int int) string {
	return strconv.FormatInt(int64(int), 10)
}

func CalculateUserMapBorder(resLocation Location, resMap Map) UserMap {
	mapSize := DefaultUserMap(resLocation)

	if resLocation.AxisX < displayMapSize {
		mapSize.leftIndent = 0
		mapSize.rightIndent = displayMapSize * 2
	}
	if resLocation.AxisY < displayMapSize {
		mapSize.downIndent = 0
		mapSize.upperIndent = displayMapSize * 2
	}
	if mapSize.rightIndent > resMap.SizeX && resLocation.AxisX > displayMapSize {
		mapSize.leftIndent = resMap.SizeX - displayMapSize*2
		mapSize.rightIndent = resMap.SizeX
	}
	if mapSize.upperIndent > resMap.SizeY && resLocation.AxisY > displayMapSize {
		mapSize.downIndent = resMap.SizeY - displayMapSize*2
		mapSize.upperIndent = resMap.SizeY
	}

	return mapSize
}
