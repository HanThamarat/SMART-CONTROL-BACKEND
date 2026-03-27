package domain

import (
	"time"

	"gorm.io/gorm"
)

type WidgetType struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"column:name" json:"name"`
	Status       bool           `gorm:"default:true" json:"status"`
	WidgetFormat string         `gorm:"column:widget_format" json:"widget_format"`
	CreatedAt    time.Time      `gorm:"column:created_at;" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;" json:"deleted_at"`

	Widgets []Widget `gorm:"foreignKey:WidgetType;references:ID" json:"widgets"`
}

type Widget struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	WidgetName    string         `gorm:"column:widget_name" json:"widget_name"`
	WidgetType    uint           `gorm:"column:widget_type;type:bigint(20) unsigned" json:"widget_type"`
	WidgetSetting string         `gorm:"column:widget_setting" json:"widget_setting"`
	Status        bool           `gorm:"default:true" json:"status"`
	CreatedAt     time.Time      `gorm:"column:created_at;" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;" json:"deleted_at"`
}

type WidgetDTO struct {
	WidgetName    string `json:"widget_name"`
	WidgetType    uint   `json:"widget_type"`
	WidgetSetting any    `json:"widget_setting"`
	Status        bool   `json:"status"`
}

type WidgetUpdateDTO struct {
	WidgetName    string `json:"widget_name"`
	WidgetSetting *any   `json:"widget_setting"`
	Status        bool   `json:"status"`
}

type WidgetFormat struct {
	ID            uint   `gorm:"column:id" json:"id"`
	WidgetName    string `gorm:"column:widget_name" json:"widget_name"`
	WidgetType    string `gorm:"column:widget_type" json:"widget_type"`
	WidgetSetting string `gorm:"column:widget_setting" json:"widget_setting"`
	Status        bool   `gorm:"column:status" json:"status"`
}

type WidgetResponse struct {
	ID            uint   `gorm:"column:id" json:"id"`
	WidgetName    string `gorm:"column:widget_name" json:"widget_name"`
	WidgetType    string `gorm:"column:widget_type" json:"widget_type"`
	WidgetSetting any    `gorm:"column:widget_setting" json:"widget_setting"`
	Status        bool   `gorm:"column:status" json:"status"`
}

type WidgetRepository interface {
	CreateNewWidget(dto WidgetDTO) (*Widget, error)
	FindAllWidget() (*[]WidgetResponse, error)
	FindWidget(id int) (*WidgetResponse, error)
	UpdateWidget(id int, dto WidgetUpdateDTO) (*WidgetResponse, error)
	DeleteWidget(id int) (*WidgetResponse, error)
}

type WidgetUsecase interface {
	CreateNewWidget(dto WidgetDTO) (*Widget, error)
	FindAllWidget() (*[]WidgetResponse, error)
	FindWidget(id int) (*WidgetResponse, error)
	UpdateWidget(id int, dto WidgetUpdateDTO) (*WidgetResponse, error)
	DeleteWidget(id int) (*WidgetResponse, error)
}
