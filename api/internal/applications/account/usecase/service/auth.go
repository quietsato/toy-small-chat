package service

type AuthService interface {
	GenerateToken(id string) string
}
