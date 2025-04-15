package order

import (
	"context"
	models "eventdrivensystem/internal/models/order"
	"eventdrivensystem/pkg/util"

	"gorm.io/gorm"
)

type OrderDomainWriter interface {
	CreateOrder(ctx context.Context, order *models.Order, opts ...util.DbOptions) (*models.Order, error)
}

func (o *OrderDomain) CreateOrder(ctx context.Context, order *models.Order, opts ...util.DbOptions) (*models.Order, error) {
	return o.createOrderSql(ctx, order, opts...)
}

func (o *OrderDomain) BeginTx(ctx context.Context) *gorm.DB {
	return o.db.Begin()
}
