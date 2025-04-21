package user_test

import (
	"eventdrivensystem/configs"
	"eventdrivensystem/internal/domain"
	domainMocks "eventdrivensystem/internal/mocks/domain"
	"eventdrivensystem/internal/usecase/user"
	"eventdrivensystem/pkg/databases"
	"eventdrivensystem/pkg/logger"
	"testing"

	"github.com/golang/mock/gomock"
)

type testPrep struct {
	ctrl        *gomock.Controller
	mockUser    *domainMocks.MockUserDomainHandler
	mockOutbox  *domainMocks.MockOutboxDomainHandler
	mockNotif   *domainMocks.MockNotificationDomainHandler
	log         logger.Logger
	cfg         *configs.AppConfig
	userUseCase user.UserUsecaseHandler
	mockDB      *databases.MockDB
}

func setupTest(t *testing.T) *testPrep {
	ctrl := gomock.NewController(t)
	mockUser := domainMocks.NewMockUserDomainHandler(ctrl)
	mockOutbox := domainMocks.NewMockOutboxDomainHandler(ctrl)
	mockNotif := domainMocks.NewMockNotificationDomainHandler(ctrl)
	log := logger.Init(logger.Options{
		Output:    logger.OutputStdout,
		Formatter: logger.FormatJSON,
		Level:     logger.LevelInfo,
	})
	cfg := &configs.AppConfig{}
	mockDB := databases.SetupMockDB(t)

	dom := &domain.Domain{
		User:         mockUser,
		Outbox:       mockOutbox,
		Notification: mockNotif,
	}

	userUseCase := user.NewUserUsecase(cfg, log, dom)

	return &testPrep{
		ctrl:        ctrl,
		mockUser:    mockUser,
		mockOutbox:  mockOutbox,
		mockNotif:   mockNotif,
		log:         log,
		cfg:         cfg,
		userUseCase: userUseCase,
		mockDB:      mockDB,
	}
}
