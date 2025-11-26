package routes_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/quietsato/toy-small-chat/api/internal/applications/account/infrastructure/serviceimpl"
	"github.com/quietsato/toy-small-chat/api/internal/applications/room/infrastructure/queryprocessorimpl"
	"github.com/quietsato/toy-small-chat/api/internal/di"
	"github.com/quietsato/toy-small-chat/api/internal/server/routes"
	"github.com/stretchr/testify/require"
)

func TestPublicRoutes(t *testing.T) {

}

func TestProtectedRoutes(t *testing.T) {
	t.Parallel()

	t.Run("JWT がない場合は UnauthorizedError", func(t *testing.T) {
		t.Parallel()

		secretKey := []byte("dummy")
		auth := serviceimpl.NewAuthService(secretKey)

		r := chi.NewRouter()
		routes.Setup(r, &di.Container{
			Auth: di.AuthDeps{
				Service:    auth,
				Middleware: auth,
			},
		})

		body := bytes.NewBufferString(`{}`)
		req := httptest.NewRequest(http.MethodPost, "/rooms", body)
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		require.Equal(t, http.StatusUnauthorized, rr.Result().StatusCode)
	})

	t.Run("無効なJWTの場合はUnauthorizedError", func(t *testing.T) {
		t.Parallel()

		secretKey := []byte("dummy")
		auth := serviceimpl.NewAuthService(secretKey)

		r := chi.NewRouter()
		routes.Setup(r, &di.Container{
			Room: di.RoomDeps{
				Query: queryprocessorimpl.NewRoomQueryProcessorOnDB(nil),
			},
			Auth: di.AuthDeps{
				Service:    auth,
				Middleware: auth,
			},
		})

		// 異なる秘密鍵で作成したトークン
		wrongJwt := jwtauth.New("HS256", []byte("wrong-key"), nil)
		_, token, _ := wrongJwt.Encode(map[string]any{"account_id": "123"})

		body := bytes.NewBufferString(`{}`)
		req := httptest.NewRequest(http.MethodPost, "/rooms", body)
		req.Header.Add("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		require.Equal(t, http.StatusUnauthorized, rr.Result().StatusCode)
	})
}
