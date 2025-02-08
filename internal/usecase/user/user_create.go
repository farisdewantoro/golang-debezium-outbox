package user

import (
	"context"
	notificationModels "eventdrivensystem/internal/models/notification"
	outboxModels "eventdrivensystem/internal/models/outbox"
	userModels "eventdrivensystem/internal/models/user"

	"eventdrivensystem/pkg/errors"
	"eventdrivensystem/pkg/util"

	"github.com/jackc/pgtype"
)

func (u *UserUsecase) CreateUser(ctx context.Context, param *userModels.CreateUserParam) error {
	dbTx := u.userDomain.BeginTx(ctx)
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
	user, err := u.userDomain.CreateUser(ctx, param.ToDomain(), dbOptions)

	if err != nil {
		return err
	}

	pNotif := notificationModels.Notification{
		Status:  notificationModels.NotificationStatusPending,
		UserID:  user.ID,
		Type:    notificationModels.NotificationTypeUserRegistration,
		Message: notificationModels.NotificationMessageUserRegistration,
	}

	_, err = u.notificationDomain.CreateNotification(ctx, &pNotif, dbOptions)
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
		AggregateID:   pNotif.ID.String(),
		AggregateType: pNotif.TableName(),
		EventType:     outboxModels.OutboxEventTypeNotifUserRegistration,
	}

	err = u.outboxDomain.CreateOutbox(ctx, &outbox, dbOptions)
	if err != nil {
		return err
	}

	return nil
}
