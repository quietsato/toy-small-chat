package serviceimpl_test

import (
	"testing"

	"github.com/go-chi/jwtauth/v5"
	"github.com/quietsato/toy-small-chat/api/internal/applications/account/infrastructure/serviceimpl"
	"github.com/stretchr/testify/require"
)

func TestAuthServiceImpl_GenerateToken(t *testing.T) {
	t.Parallel()

	t.Run("トークンが正常に生成される", func(t *testing.T) {
		t.Parallel()

		secretKey := []byte("test-secret-key")
		auth := serviceimpl.NewAuthService(secretKey)

		accountID := "test-account-id"
		token := auth.GenerateToken(accountID)

		require.NotEmpty(t, token)
	})

	t.Run("生成されたトークンが検証可能", func(t *testing.T) {
		t.Parallel()

		secretKey := []byte("test-secret-key")
		auth := serviceimpl.NewAuthService(secretKey)

		accountID := "test-account-id"
		token := auth.GenerateToken(accountID)

		// トークンを検証
		jwt := jwtauth.New("HS256", secretKey, nil)
		parsedToken, err := jwt.Decode(token)
		require.NoError(t, err)

		claims := parsedToken.PrivateClaims()
		require.Equal(t, accountID, claims[serviceimpl.AccountIDKey])
	})

	t.Run("異なるアカウントIDで異なるトークンが生成される", func(t *testing.T) {
		t.Parallel()

		secretKey := []byte("test-secret-key")
		auth := serviceimpl.NewAuthService(secretKey)

		token1 := auth.GenerateToken("account-1")
		token2 := auth.GenerateToken("account-2")

		require.NotEqual(t, token1, token2, "expected different tokens for different account IDs")
	})
}

func TestNewAuthService(t *testing.T) {
	t.Parallel()

	t.Run("正しく初期化される", func(t *testing.T) {
		t.Parallel()
		secretKey := []byte("test-secret-key")
		auth := serviceimpl.NewAuthService(secretKey)

		require.NotNil(t, auth)
	})
}

func TestAuthServiceImpl_GetTokenAuthForMiddleware(t *testing.T) {
	t.Parallel()

	t.Run("JWTAuthが取得できる", func(t *testing.T) {
		t.Parallel()
		secretKey := []byte("test-secret-key")
		auth := serviceimpl.NewAuthService(secretKey)

		tokenAuth := auth.GetTokenAuthForMiddleware()

		require.NotNil(t, tokenAuth)
	})
}
