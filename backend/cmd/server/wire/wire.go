//go:build wireinject
// +build wireinject

package wire

import (
	"backend/internal/handler"
	"backend/internal/repository"
	"backend/internal/server"
	"backend/internal/service"
	"backend/pkg/aws"
	"backend/pkg/helper/sid"
	"backend/pkg/jwt"
	"backend/pkg/log"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

var HandlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewUserHandler,
	handler.NewMusicHandler,
)

var ServiceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,
	service.NewMusicService,
)

var RepositorySet = wire.NewSet(
	repository.NewDB,
	repository.NewRedis,
	repository.NewRepository,
	repository.NewUserRepository,
	repository.NewMusicRepository,
)

func NewApp(*viper.Viper, *log.Logger) (*server.Server, func(), error) {
	panic(wire.Build(
		RepositorySet,
		ServiceSet,
		HandlerSet,
		server.NewServer,
		server.NewServerHTTP,
		sid.NewSid,
		jwt.NewJwt,
		aws.NewS3,
	))
}
