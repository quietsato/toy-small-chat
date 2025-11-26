package routes

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/quietsato/toy-small-chat/api/internal/applications/room/controller"
	"github.com/quietsato/toy-small-chat/api/internal/di"
)

type ctxKeyRoomID struct{}

func getRoomIDFromContext(ctx context.Context) *string {
	if roomID, ok := ctx.Value(ctxKeyRoomID{}).(string); ok {
		return &roomID
	}
	return nil
}

func roomCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		roomID := chi.URLParam(r, "roomID")
		slog.InfoContext(ctx, "RoomCtx", slog.String("roomID", roomID))
		ctx = context.WithValue(ctx, ctxKeyRoomID{}, roomID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getRooms(dic *di.Container) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		c := controller.NewGetRoomsController(dic.Room.Query)
		rooms, err := c.GetRooms(ctx, controller.GetRoomsInput{})
		if err != nil {
			slog.Error("failed to get rooms", slog.Any("err", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		res, err := json.Marshal(rooms)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(res)
		if err != nil {
			slog.Error("failed to write", slog.Any("err", err))
		}
	})
}

func createRoom(dic *di.Container) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		accountID := getAccountIDFromContext(ctx)
		if accountID == nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		defer r.Body.Close()
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		inp := controller.CreateRoomInput{}
		inp.CreatedBy = *accountID
		if err := json.Unmarshal(bytes, &inp); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		c := controller.NewCreateRoomController(dic.Room.Repo)
		rooms, err := c.CreateRoom(ctx, inp)
		if err != nil {
			slog.Error("failed to create room", slog.Any("err", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		res, err := json.Marshal(rooms)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(res)
		if err != nil {
			slog.Error("failed to write", slog.Any("err", err))
		}
	})
}
