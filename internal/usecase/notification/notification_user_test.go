package notification_test

import (
	"context"
	"eventdrivensystem/internal/models/notification"
	"eventdrivensystem/pkg/errors"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSendNotificationUserRegistration(t *testing.T) {
	ctx := context.TODO()
	defer func() {
		ctx.Done()
	}()

	testCases := []struct {
		name      string
		id        strfmt.UUID4
		mockSetup func(*testPrep)
		wantError error
	}{
		{
			name: "success",
			id:   strfmt.UUID4("123e4567-e89b-12d3-a456-426614174000"),
			mockSetup: func(prep *testPrep) {

				expectedNotif := notification.Notification{
					ID:      strfmt.UUID4("123e4567-e89b-12d3-a456-426614174000"),
					UserID:  strfmt.UUID4("123e4567-e89b-12d3-a456-426614174000"),
					Status:  notification.NotificationStatusPending,
					Type:    notification.NotificationTypeUserRegistration,
					Message: notification.NotificationMessageUserRegistration,
				}

				prep.mockNotif.EXPECT().
					GetNotification(ctx, gomock.Any()).
					Return(&expectedNotif, nil)

				updateNotif := expectedNotif
				updateNotif.Status = notification.NotificationStatusSent
				prep.mockNotif.EXPECT().
					UpdateNotification(ctx, &updateNotif).
					Return(&updateNotif, nil)
			},
			wantError: nil,
		},
		{
			name: "failed_get_notification",
			id:   strfmt.UUID4("123e4567-e89b-12d3-a456-426614174000"),
			mockSetup: func(prep *testPrep) {
				prep.mockNotif.EXPECT().
					GetNotification(ctx, gomock.Any()).
					Return(nil, errors.ErrSQLGet)
			},
			wantError: errors.ErrSQLGet,
		},
		{
			name: "notification_already_sent",
			id:   strfmt.UUID4("123e4567-e89b-12d3-a456-426614174000"),
			mockSetup: func(prep *testPrep) {
				expectedNotif := &notification.Notification{
					ID:      strfmt.UUID4("123e4567-e89b-12d3-a456-426614174000"),
					UserID:  strfmt.UUID4("123e4567-e89b-12d3-a456-426614174000"),
					Status:  notification.NotificationStatusSent,
					Type:    notification.NotificationTypeUserRegistration,
					Message: notification.NotificationMessageUserRegistration,
				}

				prep.mockNotif.EXPECT().
					GetNotification(ctx, gomock.Any()).
					Return(expectedNotif, nil)
			},
			wantError: nil,
		},
		{
			name: "failed_update_notification",
			id:   strfmt.UUID4("123e4567-e89b-12d3-a456-426614174000"),
			mockSetup: func(prep *testPrep) {
				expectedNotif := notification.Notification{
					ID:      strfmt.UUID4("123e4567-e89b-12d3-a456-426614174000"),
					UserID:  strfmt.UUID4("123e4567-e89b-12d3-a456-426614174000"),
					Status:  notification.NotificationStatusPending,
					Type:    notification.NotificationTypeUserRegistration,
					Message: notification.NotificationMessageUserRegistration,
				}

				prep.mockNotif.EXPECT().
					GetNotification(ctx, gomock.Any()).
					Return(&expectedNotif, nil)

				updateNotif := expectedNotif
				updateNotif.Status = notification.NotificationStatusSent
				prep.mockNotif.EXPECT().
					UpdateNotification(ctx, &updateNotif).
					Return(nil, errors.ErrSQLCreate)
			},
			wantError: errors.ErrSQLCreate,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prep := setupTest(t)
			defer prep.ctrl.Finish()
			defer prep.mockDB.DB.Close()

			// Setup mocks
			tc.mockSetup(prep)

			// Execute
			err := prep.notifUseCase.SendNotificationUserRegistration(ctx, tc.id)

			// Assert
			if tc.wantError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.wantError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
