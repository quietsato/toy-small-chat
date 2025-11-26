package routes

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/quietsato/toy-small-chat/api/internal/applications/account/controller"
	authserviceimpl "github.com/quietsato/toy-small-chat/api/internal/applications/account/infrastructure/serviceimpl"
	"github.com/quietsato/toy-small-chat/api/internal/di"
)

type ctxKeyAccountID struct{}

func getAccountIDFromContext(ctx context.Context) *string {
	if accountID, ok := ctx.Value(ctxKeyAccountID{}).(string); ok {
		return &accountID
	}
	return nil
}

func accountCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		_, claims, _ := jwtauth.FromContext(ctx)
		accountId, ok := claims[authserviceimpl.AccountIDKey]
		if ok {
			ctx = context.WithValue(ctx, ctxKeyAccountID{}, accountId)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func createAccount(dic *di.Container) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		c := controller.NewCreateAccountController(dic.Account.Repo, dic.Auth.Service)

		bytes, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		inp := controller.CreateAccountInput{}
		if err := json.Unmarshal(bytes, &inp); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		res, err := c.CreateAccount(ctx, inp)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		resBytes, err := json.Marshal(res)
		if err != nil {
			slog.ErrorContext(ctx, "failed to login", slog.Any("err", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(resBytes); err != nil {
			slog.ErrorContext(ctx, "failed to write response", slog.Any("err", err))
		}
	})
}

func login(dic *di.Container) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		bytes, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		inp := controller.LoginInput{}
		if err := json.Unmarshal(bytes, &inp); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		res, err := controller.NewLoginController(dic.Account.Query).Login(ctx, inp)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		resBytes, err := json.Marshal(res)
		if err != nil {
			slog.ErrorContext(ctx, "failed to login", slog.Any("err", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(resBytes); err != nil {
			slog.ErrorContext(ctx, "failed to write response", slog.Any("err", err))
		}
	})
}
