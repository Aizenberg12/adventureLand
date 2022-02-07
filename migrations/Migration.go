package migrations

import (
	"fmt"
	"project0/config"
	"project0/repository"
)

func Migrate() bool {
	err := config.Db.AutoMigrate(
		repository.User{},
		repository.Location{},
		repository.Cellule{},
		repository.Map{},
		repository.Teleport{},
		repository.Test{},
		repository.Sectest{},
	)
	if err != nil {
		fmt.Println("Migration failed")
		panic(err)
		return false
	}
	return true
}
