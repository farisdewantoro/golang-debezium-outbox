package user_test

import (
	"context"
	"eventdrivensystem/internal/models/notification"
	user "eventdrivensystem/internal/models/user"
	"eventdrivensystem/pkg/errors"
	"eventdrivensystem/pkg/util"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	testCases := []struct {
		name      string
		param     *user.CreateUserParam
		mockSetup func(*testPrep)
		wantError error
	}{
		{
			name: "success",
			param: &user.CreateUserParam{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(prep *testPrep) {
				ctx := context.Background()
				expectedUser := &user.User{
					ID:       strfmt.UUID4("123e4567-e89b-12d3-a456-426614174000"),
					Email:    "test@example.com",
					Password: "password123",
				}
				expectedNotif := &notification.Notification{
					UserID:  expectedUser.ID,
					Status:  notification.NotificationStatusPending,
					Type:    notification.NotificationTypeUserRegistration,
					Message: notification.NotificationMessageUserRegistration,
				}

				prep.mockUser.EXPECT().BeginTx(ctx).Return(prep.mockDB.Gorm)
				prep.mockUser.EXPECT().
					CreateUser(ctx, gomock.Any(), util.DbOptions{Transaction: prep.mockDB.Gorm}).
					Return(expectedUser, nil)
				prep.mockNotif.EXPECT().
					CreateNotification(ctx, gomock.Any(), util.DbOptions{Transaction: prep.mockDB.Gorm}).
					Return(expectedNotif, nil)
				prep.mockOutbox.EXPECT().
					CreateOutbox(ctx, gomock.Any(), util.DbOptions{Transaction: prep.mockDB.Gorm}).
					Return(nil)
			},
			wantError: nil,
		},
		{
			name: "failed_create_user",
			param: &user.CreateUserParam{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(prep *testPrep) {
				ctx := context.Background()
				prep.mockUser.EXPECT().BeginTx(ctx).Return(prep.mockDB.Gorm)
				prep.mockUser.EXPECT().
					CreateUser(ctx, gomock.Any(), util.DbOptions{Transaction: prep.mockDB.Gorm}).
					Return(nil, errors.ErrSQLCreate)
			},
			wantError: errors.ErrSQLCreate,
		},
		{
			name: "failed_create_notification",
			param: &user.CreateUserParam{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(prep *testPrep) {
				ctx := context.Background()
				expectedUser := &user.User{
					ID:       strfmt.UUID4("123e4567-e89b-12d3-a456-426614174000"),
					Email:    "test@example.com",
					Password: "password123",
				}

				prep.mockUser.EXPECT().BeginTx(ctx).Return(prep.mockDB.Gorm)
				prep.mockUser.EXPECT().
					CreateUser(ctx, gomock.Any(), util.DbOptions{Transaction: prep.mockDB.Gorm}).
					Return(expectedUser, nil)
				prep.mockNotif.EXPECT().
					CreateNotification(ctx, gomock.Any(), util.DbOptions{Transaction: prep.mockDB.Gorm}).
					Return(nil, errors.ErrSQLCreate)
			},
			wantError: errors.ErrSQLCreate,
		},
		{
			name: "failed_create_outbox",
			param: &user.CreateUserParam{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(prep *testPrep) {
				ctx := context.Background()
				expectedUser := &user.User{
					ID:       strfmt.UUID4("123e4567-e89b-12d3-a456-426614174000"),
					Email:    "test@example.com",
					Password: "password123",
				}
				expectedNotif := &notification.Notification{
					UserID:  expectedUser.ID,
					Status:  notification.NotificationStatusPending,
					Type:    notification.NotificationTypeUserRegistration,
					Message: notification.NotificationMessageUserRegistration,
				}

				prep.mockUser.EXPECT().BeginTx(ctx).Return(prep.mockDB.Gorm)
				prep.mockUser.EXPECT().
					CreateUser(ctx, gomock.Any(), util.DbOptions{Transaction: prep.mockDB.Gorm}).
					Return(expectedUser, nil)
				prep.mockNotif.EXPECT().
					CreateNotification(ctx, gomock.Any(), util.DbOptions{Transaction: prep.mockDB.Gorm}).
					Return(expectedNotif, nil)
				prep.mockOutbox.EXPECT().
					CreateOutbox(ctx, gomock.Any(), util.DbOptions{Transaction: prep.mockDB.Gorm}).
					Return(errors.ErrSQLCreate)
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
			err := prep.userUseCase.CreateUser(context.Background(), tc.param)

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
