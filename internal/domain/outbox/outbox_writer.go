package outbox

import (
	"context"
	models "eventdrivensystem/internal/models/outbox"
	"eventdrivensystem/pkg/util"
)

type OutboxDomainWriter interface {
	CreateOutbox(ctx context.Context, Outbox *models.OutboxEvent, opts ...util.DbOptions) error
}

func (u *OutboxDomain) CreateOutbox(ctx context.Context, outbox *models.OutboxEvent, opts ...util.DbOptions) error {
	return u.createOutboxSql(ctx, outbox, opts...)
}
