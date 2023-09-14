//go:build wireinject
// +build wireinject

package wire

import (
	"backend/cmd/migration/internal"
	"backend/internal/repository"
	"backend/pkg/log"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

var RepositorySet = wire.NewSet(
	repository.NewDB,
	repository.NewRedis,
	repository.NewRepository,
	repository.NewUserRepository,
)

func NewApp(*viper.Viper, *log.Logger) (*internal.Migrate, func(), error) {
	panic(wire.Build(
		RepositorySet,
		internal.NewMigrate,
	))
}
