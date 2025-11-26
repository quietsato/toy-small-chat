package di

import (
	"github.com/jackc/pgx/v5/pgxpool"
	accountqueryimpl "github.com/quietsato/toy-small-chat/api/internal/applications/account/infrastructure/queryprocessorimpl"
	accountrepoimpl "github.com/quietsato/toy-small-chat/api/internal/applications/account/infrastructure/repositoryimpl"
	accountserviceimpl "github.com/quietsato/toy-small-chat/api/internal/applications/account/infrastructure/serviceimpl"
	accountquery "github.com/quietsato/toy-small-chat/api/internal/applications/account/usecase/queryprocessor"
	accountrepo "github.com/quietsato/toy-small-chat/api/internal/applications/account/usecase/repository"
	accountservice "github.com/quietsato/toy-small-chat/api/internal/applications/account/usecase/service"
	messagequeryimpl "github.com/quietsato/toy-small-chat/api/internal/applications/message/infrastructure/queryprocessorimpl"
	messagerepoimpl "github.com/quietsato/toy-small-chat/api/internal/applications/message/infrastructure/repositoryimpl"
	messagequery "github.com/quietsato/toy-small-chat/api/internal/applications/message/usecase/queryprocessor"
	messagerepo "github.com/quietsato/toy-small-chat/api/internal/applications/message/usecase/repository"
	roomqueryimpl "github.com/quietsato/toy-small-chat/api/internal/applications/room/infrastructure/queryprocessorimpl"
	roomrepoimpl "github.com/quietsato/toy-small-chat/api/internal/applications/room/infrastructure/repositoryimpl"
	roomquery "github.com/quietsato/toy-small-chat/api/internal/applications/room/usecase/queryprocessor"
	roomrepo "github.com/quietsato/toy-small-chat/api/internal/applications/room/usecase/repository"
	authmiddleware "github.com/quietsato/toy-small-chat/api/internal/server/middlewares/auth"
)

type AccountDeps struct {
	Repo  accountrepo.AccountRepository
	Query accountquery.AccountQueryProcessor
}

type MessageDeps struct {
	Repo  messagerepo.MessageRepository
	Query messagequery.MessageQueryProcessor
}

type RoomDeps struct {
	Repo  roomrepo.RoomRepository
	Query roomquery.RoomQueryProcessor
}

type AuthDeps struct {
	Service    accountservice.AuthService
	Middleware authmiddleware.Provider
}

type Container struct {
	Account AccountDeps
	Message MessageDeps
	Room    RoomDeps
	Auth    AuthDeps
}

func New(pool *pgxpool.Pool, jwtSecretKey string) *Container {
	auth := accountserviceimpl.NewAuthService([]byte(jwtSecretKey))

	return &Container{
		Account: AccountDeps{
			Repo:  accountrepoimpl.NewAccountRepositoryOnDB(pool),
			Query: accountqueryimpl.NewAccountQueryProcessorOnDB(pool),
		},
		Message: MessageDeps{
			Repo:  messagerepoimpl.NewMessageRepositoryOnDB(pool),
			Query: messagequeryimpl.NewMessageQueryProcessorOnDB(pool),
		},
		Room: RoomDeps{
			Repo:  roomrepoimpl.NewRoomRepositoryOnDB(pool),
			Query: roomqueryimpl.NewRoomQueryProcessorOnDB(pool),
		},
		Auth: AuthDeps{
			Service:    auth,
			Middleware: auth,
		},
	}
}
