package outbox

import (
	"context"
	models "eventdrivensystem/internal/models/outbox"
	"eventdrivensystem/pkg/util"
)

func (u *OutboxDomain) CreateOutbox(ctx context.Context, outbox *models.OutboxEvent, opts ...util.DbOptions) error {
	return u.createOutboxSql(ctx, outbox, opts...)
}
