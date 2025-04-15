package order

import (
	"context"
	models "eventdrivensystem/internal/models/order"
	"eventdrivensystem/pkg/util"

	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (o *OrderDomain) createOrderSql(ctx context.Context, order *models.Order, opts ...util.DbOptions) (*models.Order, error) {
	var (
		db  *gorm.DB
		opt util.DbOptions
	)

	if len(opts) > 0 {
		opt = opts[0]
	}

	db = opt.Extract(ctx, o.db)

	// Generate UUID if not already set
	if order.ID == "" {
		order.ID = strfmt.UUID4(uuid.New().String())
	}

	return order, db.Create(order).Error
}
