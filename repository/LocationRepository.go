package repository

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project0/config"
)

type Location struct {
	ID       uint `gorm:"primaryKey"`
	UserTgId uint
	UserID   uint
	AxisX    *int
	AxisY    *int
	Map      string `gorm:"embedded"`
}

func GetOrCreateMyLocation(update tgbotapi.Update) Location {
	resUser := GetOrCreateUser(update)

	AsX := 7
	AsY := 2

	result := Location{
		UserTgId: uint(update.Message.From.ID),
		AxisX:    &AsX,
		AxisY:    &AsY,
		Map:      "Main Place",
	}

	err := config.Db.Where(&Location{UserID: resUser.ID}).FirstOrCreate(&result).Error
	if err != nil {
		panic(err)
	}

	return result
}

func UpdateLocation(update tgbotapi.Update, LocationStruct Location) Location {
	myLocation := GetOrCreateMyLocation(update)
	//location := LocationStruct

	if myLocation.Map != LocationStruct.Map {
		resCellule := GetCellule(Cellule{Map: myLocation.Map, AxisX: *LocationStruct.AxisX, AxisY: *LocationStruct.AxisY})
		resTeleport := GetTeleport(Teleport{ID: resCellule.TeleportId})
		LocationStruct = Location{
			AxisX: &resTeleport.StartX,
			AxisY: &resTeleport.StartY,
			Map:   resTeleport.Map,
		}
	}

	fmt.Println(LocationStruct)
	var result Cellule
	var err error

	err = config.Db.First(&result, &Cellule{Map: LocationStruct.Map, AxisX: *LocationStruct.AxisX, AxisY: *LocationStruct.AxisY}).Error
	if err != nil {
		if err.Error() == "record not found" {
			return myLocation
		}
		panic(err)
	}

	fmt.Println("Result location: ", result, "LocaStruct: ", LocationStruct)

	if !result.CanStep && result.Type == "teleport" {
		err = config.Db.Where(&Location{UserTgId: uint(update.Message.From.ID)}).Updates(LocationStruct).Error
		if err != nil {
			panic(err)
		}
	} else if !result.CanStep {
		return myLocation
	} else {
		err = config.Db.Where(&Location{UserTgId: uint(update.Message.From.ID)}).Updates(LocationStruct).Error
		if err != nil {
			panic(err)
		}
	}

	myLocation = GetOrCreateMyLocation(update)
	return myLocation
}
