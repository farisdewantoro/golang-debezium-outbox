package order

import (
	"context"
	notificationModels "eventdrivensystem/internal/models/notification"
	orderModels "eventdrivensystem/internal/models/order"
	outboxModels "eventdrivensystem/internal/models/outbox"
	"eventdrivensystem/pkg/errors"
	"eventdrivensystem/pkg/util"
	"fmt"

	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgtype"
)

func (u *OrderUsecase) CreateOrder(ctx context.Context, param *orderModels.CreateOrderParam) error {
	dbTx := u.orderDomain.BeginTx(ctx)
	var (
		err error
	)

	defer func() {
		if tmpErr := util.FirstNotNil(recover(), err); tmpErr != nil {
			if tmpErr != err {
				u.log.ErrorWithContext(ctx, tmpErr)
				return
			}

			if errRollback := dbTx.Rollback().Error; errRollback != nil {
				u.log.ErrorWithContext(ctx, errRollback)
				return
			}
		} else {
			if errCommit := dbTx.Commit().Error; errCommit != nil {
				u.log.ErrorWithContext(ctx, errCommit)
				err = errors.ErrSQLTx
				return
			}
		}
	}()

	dbOptions := util.DbOptions{
		Transaction: dbTx,
	}

	order, err := u.orderDomain.CreateOrder(ctx, param.ToDomain(), dbOptions)
	if err != nil {
		return err
	}

	// Convert uint UserID to strfmt.UUID4
	userIDStr := fmt.Sprintf("%d", order.UserID)
	userUUID := strfmt.UUID4(userIDStr)

	pNotif := notificationModels.Notification{
		Status:  notificationModels.NotificationStatusPending,
		UserID:  userUUID,
		Type:    notificationModels.NotificationTypeOrderCreated,
		Message: "Your order has been created successfully",
	}

	notification, err := u.notificationDomain.CreateNotification(ctx, &pNotif, dbOptions)
	if err != nil {
		return err
	}

	jsonB := pgtype.JSONB{}
	if err := jsonB.Set(pNotif); err != nil {
		u.log.ErrorfWithContext(ctx, "error parse json outbox: %v", err)
		return errors.ErrParseJsonOutbox
	}

	outbox := outboxModels.OutboxEvent{
		Payload:       &jsonB,
		AggregateID:   notification.ID.String(),
		AggregateType: notification.TableName(),
		EventType:     outboxModels.OutboxEventTypeNotifOrderCreated,
	}

	err = u.outboxDomain.CreateOutbox(ctx, &outbox, dbOptions)
	if err != nil {
		return err
	}

	return nil
}
