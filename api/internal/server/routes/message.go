package routes

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/quietsato/toy-small-chat/api/internal/applications/message/controller"
	"github.com/quietsato/toy-small-chat/api/internal/di"
)

func getMessages(dic *di.Container) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		roomID, ok := ctx.Value(ctxKeyRoomID{}).(string)
		if !ok {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		c := controller.NewGetMessagesController(dic.Message.Query)
		msgs, err := c.GetMessages(controller.GetMessagesInput{RoomID: roomID})
		if err != nil {
			slog.Error("failed to get messages", slog.Any("err", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		res, err := json.Marshal(msgs)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(res)
		if err != nil {
			panic(err) // TODO: fix
		}
	})
}

func createMessage(dic *di.Container) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		roomID := getRoomIDFromContext(ctx)
		accountID := getAccountIDFromContext(ctx)
		if roomID == nil || accountID == nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		inp := controller.CreateMessageInput{}
		if err := json.Unmarshal(body, &inp); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		inp.RoomID = *roomID
		inp.AuthorID = *accountID

		c := controller.NewCreateMessageController(dic.Message.Repo)
		err = c.CreateMessage(ctx, inp)
		if err != nil {
			slog.Error("failed to create message", slog.Any("err", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		res, err := json.Marshal(controller.CreateMessageOutput{})
		if err != nil {
			slog.Warn("failed to marshal response", slog.Any("err", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(res)
		if err != nil {
			slog.Warn("failed to write response", slog.Any("err", err))
		}
	})
}
