package repositories

import (
	"encoding/json"
	"time"

	"github.com/HanThamarat/SMART-CONTROL-BACKEND/internal/domain"
	"gorm.io/gorm"
)

type gormWidgetRepository struct {
	db *gorm.DB
}

func NewGormWidgetRepository(db *gorm.DB) domain.WidgetRepository {
	return &gormWidgetRepository{db: db}
}

func (r *gormWidgetRepository) CreateNewWidget(dto domain.WidgetDTO) (*domain.Widget, error) {
	var dataPrepare domain.Widget

	byte, err := json.Marshal(dto.WidgetSetting)
	if err != nil {
		return nil, err
	}

	widgetSettingjsonString := string(byte)

	dataPrepare.WidgetName = dto.WidgetName
	dataPrepare.WidgetType = dto.WidgetType
	dataPrepare.WidgetSetting = widgetSettingjsonString
	dataPrepare.Status = dto.Status

	if err := r.db.Create(&dataPrepare).Error; err != nil {
		return nil, err
	}

	return &dataPrepare, nil
}

func (r *gormWidgetRepository) FindAllWidget() (*[]domain.WidgetResponse, error) {
	var formatModel []domain.WidgetFormat

	if err := r.db.Table("widgets as w").
		Joins("inner join widget_types wt on w.widget_type = wt.id").
		Select(`w.id as id, w.widget_name as widget_name, wt."name" as widget_type, w.widget_setting as widget_setting, w.status as status`).
		Where("w.deleted_at IS NULL").
		Scan(&formatModel).Error; err != nil {
		return nil, err
	}

	responseModel := []domain.WidgetResponse{}

	for _, item := range formatModel {
		var ws map[string]any

		err := json.Unmarshal([]byte(item.WidgetSetting), &ws)
		if err != nil {
			return nil, err
		}

		newItem := domain.WidgetResponse{
			ID:            item.ID,
			WidgetName:    item.WidgetName,
			WidgetType:    item.WidgetType,
			WidgetSetting: ws,
			Status:        item.Status,
		}

		responseModel = append(responseModel, newItem)
	}

	return &responseModel, nil
}

func (r *gormWidgetRepository) FindWidget(id int) (*domain.WidgetResponse, error) {
	var formatModel domain.WidgetFormat

	if err := r.db.Table("widgets as w").
		Joins("inner join widget_types wt on w.widget_type = wt.id").
		Select(`w.id as id, w.widget_name as widget_name, wt."name" as widget_type, w.widget_setting as widget_setting, w.status as status`).
		Where("w.id = ? and w.deleted_at IS NULL", id).
		Find(&formatModel).Error; err != nil {
		return nil, err
	}

	var ws map[string]any

	err := json.Unmarshal([]byte(formatModel.WidgetSetting), &ws)
	if err != nil {
		return nil, err
	}

	responseModel := domain.WidgetResponse{
		ID:            formatModel.ID,
		WidgetName:    formatModel.WidgetName,
		WidgetType:    formatModel.WidgetType,
		WidgetSetting: ws,
		Status:        formatModel.Status,
	}

	return &responseModel, nil
}

func (r *gormWidgetRepository) UpdateWidget(id int, dto domain.WidgetUpdateDTO) (*domain.WidgetResponse, error) {
	var formatModel domain.WidgetFormat
	var model domain.Widget

	if err := r.db.Model(&model).Where("id = ?", id).Updates(map[string]interface{}{
		"widget_name": dto.WidgetName,
		"status":      dto.Status,
	}).Error; err != nil {
		return nil, err
	}

	if err := r.db.Table("widgets as w").
		Joins("inner join widget_types wt on w.widget_type = wt.id").
		Select(`w.id as id, w.widget_name as widget_name, wt."name" as widget_type, w.widget_setting as widget_setting, w.status as status`).
		Where("w.id = ? and w.deleted_at IS NULL", id).
		Find(&formatModel).Error; err != nil {
		return nil, err
	}

	var ws map[string]any

	err := json.Unmarshal([]byte(formatModel.WidgetSetting), &ws)
	if err != nil {
		return nil, err
	}

	responseModel := domain.WidgetResponse{
		ID:            formatModel.ID,
		WidgetName:    formatModel.WidgetName,
		WidgetType:    formatModel.WidgetType,
		WidgetSetting: ws,
		Status:        formatModel.Status,
	}

	return &responseModel, nil
}

func (r *gormWidgetRepository) DeleteWidget(id int) (*domain.WidgetResponse, error) {
	var formatModel domain.WidgetFormat
	var model domain.Widget

	if err := r.db.Model(&model).Where("id = ? and deleted_at IS NULL", id).Update("deleted_at", time.Now()).Error; err != nil {
		return nil, err
	}

	if err := r.db.Table("widgets as w").
		Joins("inner join widget_types wt on w.widget_type = wt.id").
		Select(`w.id as id, w.widget_name as widget_name, wt."name" as widget_type, w.widget_setting as widget_setting, w.status as status`).
		Where("w.id = ?", id).
		Find(&formatModel).Error; err != nil {
		return nil, err
	}

	var ws map[string]any

	err := json.Unmarshal([]byte(formatModel.WidgetSetting), &ws)
	if err != nil {
		return nil, err
	}

	responseModel := domain.WidgetResponse{
		ID:            formatModel.ID,
		WidgetName:    formatModel.WidgetName,
		WidgetType:    formatModel.WidgetType,
		WidgetSetting: ws,
		Status:        formatModel.Status,
	}

	return &responseModel, nil
}
