package user

import (
	"context"
	"eventdrivensystem/configs"
	models "eventdrivensystem/internal/models/user"
	"eventdrivensystem/pkg/logger"
	"eventdrivensystem/pkg/util"

	"gorm.io/gorm"
)

type UserDomain struct {
	cfg *configs.AppConfig
	db  *gorm.DB
	log logger.Logger
}

type UserDomainHandler interface {
	BeginTx(ctx context.Context) *gorm.DB
	CreateUser(ctx context.Context, user *models.User, opts ...util.DbOptions) (*models.User, error)
}

func NewUserDomain(cfg *configs.AppConfig, log logger.Logger, db *gorm.DB) UserDomainHandler {
	return &UserDomain{
		cfg: cfg,
		db:  db,
		log: log,
	}
}
