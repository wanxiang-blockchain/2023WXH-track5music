package internal

import (
	"backend/internal/model"
	"backend/pkg/log"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Migrate struct {
	db  *gorm.DB
	log *log.Logger
}

func NewMigrate(db *gorm.DB, log *log.Logger) *Migrate {
	return &Migrate{
		db:  db,
		log: log,
	}
}
func (m *Migrate) Run() {
	if err := m.db.AutoMigrate(&model.Single{}, &model.Track{}, &model.File{}); err != nil {
		m.log.Error("migrate error", zap.Error(err))
		return
	}
	m.log.Info("AutoMigrate end")
}
