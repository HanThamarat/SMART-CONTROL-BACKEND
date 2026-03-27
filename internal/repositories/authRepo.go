package repositories

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/HanThamarat/SMART-CONTROL-BACKEND/internal/domain"
	"github.com/HanThamarat/SMART-CONTROL-BACKEND/pkg/encrypt"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type gormAuthRepository struct {
	db *gorm.DB
}

func NewGormAuthRepository(db *gorm.DB) domain.AuthRepository {
	return &gormAuthRepository{db}
}

func (r *gormAuthRepository) CredentialAuth(dto domain.AuthDTO) (*domain.AuthResponse, error) {
	var userInfo domain.User

	if err := r.db.Where("(LOWER(email) = ? OR LOWER(username) = ?) AND status = ?", strings.ToLower(dto.Credential), strings.ToLower(dto.Credential), true).First(&userInfo).Error; err != nil {
		return nil, err
	}

	verifyPassword := encrypt.VerifyPassword(dto.Password, *userInfo.Password)

	if !verifyPassword {
		return nil, errors.New("Invalid your password.")
	}

	claims := jwt.MapClaims{
		"userId": userInfo.ID,
		"name":   userInfo.Name,
		"email":  userInfo.Email,
		"status": userInfo.Status,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	authResult := domain.AuthResponse{
		ID:        int(userInfo.ID),
		Email:     userInfo.Email,
		Name:      userInfo.Name,
		AuthToken: t,
	}

	return &authResult, nil
}
