package databases

import (
	"database/sql"
	"eventdrivensystem/configs"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Required for postgres driver
	_ "github.com/golang-migrate/migrate/v4/source/file"       // Required for file source
)

func NewSqlDb(cfg *configs.AppConfig) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.SQL.Username,
		cfg.SQL.Password,
		cfg.SQL.Host,
		cfg.SQL.Port,
		cfg.SQL.Database)

	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return
	}

	db.Use(tracing.NewPlugin())

	sqlDB, err := db.DB()
	if err != nil {
		return
	}

	// Create the connection pool
	sqlDB.SetConnMaxIdleTime(time.Minute * 5)

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(5)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(30)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(620 * time.Second)

	log.Info("Database connected successfully", cfg.SQL.Database)
	return db, nil
}

func NewMigrate(cfg *configs.AppConfig) (*migrate.Migrate, error) {
	migrationsPath := "file://./docs/db/migrations"
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.SQL.Username,
		cfg.SQL.Password,
		cfg.SQL.Host,
		cfg.SQL.Port,
		cfg.SQL.Database)
	// Create the migrator instance
	migrator, err := migrate.New(migrationsPath, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator instance: %w", err)
	}

	return migrator, nil

}

type MockDB struct {
	DB   *sql.DB
	Mock sqlmock.Sqlmock
	Gorm *gorm.DB
}

func SetupMockDB(t *testing.T) *MockDB {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	return &MockDB{
		DB:   db,
		Mock: mock,
		Gorm: gormDB,
	}
}
