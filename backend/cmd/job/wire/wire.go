//go:build wireinject
// +build wireinject

package wire

import (
	"backend/internal/job"
	"backend/pkg/log"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

var JobSet = wire.NewSet(job.NewJob)

func NewApp(*viper.Viper, *log.Logger) (*job.Job, func(), error) {
	panic(wire.Build(
		JobSet,
	))
}
