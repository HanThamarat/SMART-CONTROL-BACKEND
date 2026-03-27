package domain

type AuthDTO struct {
	Credential string `json:"credential"`
	Password   string `json:"password"`
}

type AuthResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AuthToken string `json:"authToken"`
}

type AuthRepository interface {
	CredentialAuth(dto AuthDTO) (*AuthResponse, error)
}

type AuthUsecase interface {
	CredentialAuth(dto AuthDTO) (*AuthResponse, error)
}
