package usecase

import "github.com/HanThamarat/SMART-CONTROL-BACKEND/internal/domain"

type widgetUsecase struct {
	repo domain.WidgetRepository
}

func NewWidgetUsecase(r domain.WidgetRepository) domain.WidgetUsecase {
	return &widgetUsecase{repo: r}
}

func (u *widgetUsecase) CreateNewWidget(dto domain.WidgetDTO) (*domain.Widget, error) {
	return u.repo.CreateNewWidget(dto)
}

func (u *widgetUsecase) FindAllWidget() (*[]domain.WidgetResponse, error) {
	return u.repo.FindAllWidget()
}

func (u *widgetUsecase) FindWidget(id int) (*domain.WidgetResponse, error) {
	return u.repo.FindWidget(id)
}

func (u *widgetUsecase) UpdateWidget(id int, dto domain.WidgetUpdateDTO) (*domain.WidgetResponse, error) {
	return u.repo.UpdateWidget(id, dto)
}

func (u *widgetUsecase) DeleteWidget(id int) (*domain.WidgetResponse, error) {
	return u.repo.DeleteWidget(id)
}
