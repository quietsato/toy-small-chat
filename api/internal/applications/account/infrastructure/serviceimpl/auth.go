package serviceimpl

import (
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/quietsato/toy-small-chat/api/internal/applications/account/usecase/service"
	authmiddleware "github.com/quietsato/toy-small-chat/api/internal/server/middlewares/auth"
)

const AccountIDKey = "accountId"

type AuthServiceImpl struct {
	tokenAuth *jwtauth.JWTAuth
}

func NewAuthService(secretKey []byte) *AuthServiceImpl {
	tokenAuth := jwtauth.New("HS256", secretKey, nil) // TODO: 公開鍵暗号方式に変更する
	return &AuthServiceImpl{tokenAuth}
}

func (a *AuthServiceImpl) GetTokenAuthForMiddleware() *jwtauth.JWTAuth {
	return a.tokenAuth
}

func (a *AuthServiceImpl) GenerateToken(id string) string {
	claims := make(map[string]any, 2)
	jwtauth.SetExpiryIn(claims, 1*time.Hour) // TODO: リフレッシュトークンに対応させたら有効期限を短くする
	claims[AccountIDKey] = id
	_, tokenString, _ := a.tokenAuth.Encode(claims)
	return tokenString
}

var _ interface {
	service.AuthService
	authmiddleware.Provider // infrastructure 層なので middleware の都合を見るのは問題ないと考える
} = new(AuthServiceImpl)
