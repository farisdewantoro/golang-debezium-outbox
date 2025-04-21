package notification_test

import (
	"eventdrivensystem/configs"
	"eventdrivensystem/internal/domain"
	domainMocks "eventdrivensystem/internal/mocks/domain"
	"eventdrivensystem/internal/usecase/notification"
	"eventdrivensystem/pkg/databases"
	"eventdrivensystem/pkg/logger"
	"testing"

	"github.com/golang/mock/gomock"
)

type testPrep struct {
	ctrl         *gomock.Controller
	mockNotif    *domainMocks.MockNotificationDomainHandler
	mockOutbox   *domainMocks.MockOutboxDomainHandler
	log          logger.Logger
	cfg          *configs.AppConfig
	notifUseCase notification.NotificationUsecaseHandler
	mockDB       *databases.MockDB
}

func setupTest(t *testing.T) *testPrep {
	ctrl := gomock.NewController(t)
	mockNotif := domainMocks.NewMockNotificationDomainHandler(ctrl)
	mockOutbox := domainMocks.NewMockOutboxDomainHandler(ctrl)
	log := logger.Init(logger.Options{
		Output:    logger.OutputStdout,
		Formatter: logger.FormatJSON,
		Level:     logger.LevelInfo,
	})
	cfg := &configs.AppConfig{}
	mockDB := databases.SetupMockDB(t)

	dom := &domain.Domain{
		Notification: mockNotif,
		Outbox:       mockOutbox,
	}

	notifUseCase := notification.NewNotificationUsecase(cfg, log, dom)

	return &testPrep{
		ctrl:         ctrl,
		mockNotif:    mockNotif,
		mockOutbox:   mockOutbox,
		log:          log,
		cfg:          cfg,
		notifUseCase: notifUseCase,
		mockDB:       mockDB,
	}
}
