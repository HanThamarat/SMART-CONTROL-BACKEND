package initial

import (
	"fmt"
	"os"

	"github.com/HanThamarat/SMART-CONTROL-BACKEND/internal/domain"
	"github.com/HanThamarat/SMART-CONTROL-BACKEND/pkg/encrypt"
	"gorm.io/gorm"
)

func UserInit(db *gorm.DB) {
	var user domain.User
	var count int64

	recheck := db.Model(&user).Count(&count)

	if recheck.Error != nil {
		fmt.Println("Have something wrong in initialize progress : ", recheck.Error)
		return
	}

	if count != 0 {
		fmt.Println("✅ User initialize success.")
		return
	}

	password := os.Getenv("PASSWORD")
	passwordHashing, err := encrypt.HashPassword(password)

	if err != nil {
		fmt.Println("Have something wrong in initialize progress : ", err.Error())
		return
	}

	user.Email = os.Getenv("EMAIL")
	user.Username = os.Getenv("USERNAME")
	user.Name = os.Getenv("NAME")
	user.Password = &passwordHashing

	createNewUser := db.Create(&user)

	if createNewUser.Error != nil {
		fmt.Println("Have something wrong in initialize progress : ", createNewUser.Error)
		return
	}

	fmt.Println("✅ User initialize success.")
}
