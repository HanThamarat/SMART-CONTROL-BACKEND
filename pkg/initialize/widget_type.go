package initial

import (
	"encoding/json"
	"fmt"

	"github.com/HanThamarat/SMART-CONTROL-BACKEND/internal/domain"
	"github.com/HanThamarat/SMART-CONTROL-BACKEND/internal/types"
	"gorm.io/gorm"
)

func WidgetTypeInit(db *gorm.DB) {
	var widgetType domain.WidgetType
	var count int64

	recheck := db.Model(&widgetType).Count(&count)

	if recheck.Error != nil {
		fmt.Println("Have something wrong in initialize progress : ", recheck.Error)
		return
	}

	if count != 0 {
		fmt.Println("✅ Widget Type initialize success.")
		return
	}

	format := types.PercentageSwitch{
		MinValue:     0,
		MaxValue:     0,
		CurrentValue: 0,
	}

	formatByte, err := json.Marshal(format)

	if err != nil {
		fmt.Println("Have something wrong in initialize progress : ", recheck.Error)
		return
	}

	convertToString := string(formatByte)

	widgetType.Name = "Ligth Control"
	widgetType.Status = true
	widgetType.WidgetFormat = convertToString

	if err := db.Create(&widgetType).Error; err != nil {
		fmt.Println("Have something wrong in initialize progress : ", recheck.Error)
		return
	}

	fmt.Println("✅ Widget type initialize success.")
}
