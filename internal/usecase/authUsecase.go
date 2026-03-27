package usecase

import "github.com/HanThamarat/SMART-CONTROL-BACKEND/internal/domain"

type authUsecase struct {
	repo domain.AuthRepository
}

func NewAuthUsecase(r domain.AuthRepository) domain.AuthUsecase {
	return &authUsecase{repo: r}
}

func (u *authUsecase) CredentialAuth(dto domain.AuthDTO) (*domain.AuthResponse, error) {
	return u.repo.CredentialAuth(dto)
}
