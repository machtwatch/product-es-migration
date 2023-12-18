package postgres

import (
	"context"
	"fmt"
	"product-es-migration/config"

	"github.com/machtwatch/catalystdk/go/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
)

// dsnConfig represent data source name (DSN) configuration
type dsnConfig struct {
	host     string
	user     string
	password string
	db       string
	port     int
	ssl      string
}

func OpenPostgres() *gorm.DB {
	pgConfig := postgres.Config{
		DSN:                  getPostgresSlaveDSN(),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}
	db, err := gorm.Open(postgres.New(pgConfig), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.StdFatal(context.Background(), nil, err, "gorm.Open() got error on connecting to database - database.openPostgres()")
	}

	if err := db.Use(tracing.NewPlugin(tracing.WithoutQueryVariables())); err != nil {
		log.StdFatal(context.Background(), nil, err, "db.Use(tracing.NewPlugin(tracing.WithoutQueryVariables())) got error  - database.openPostgres()")
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.StdFatal(context.Background(), nil, err, "db.DB() got error on opening database - database.openPostgres()")
	}

	if err = sqlDB.Ping(); err != nil {
		log.StdFatal(context.Background(), nil, err, "sqlDB.Ping() got error on ping database - database.openPostgres()")
	}

	log.Info("successfully connected to postgres database")

	return db
}

func ClosePostgres(conn *gorm.DB) {
	sqlDB, err := conn.DB()
	if err != nil {
		log.StdFatal(context.Background(), nil, err, "conn.DB()Error occurred while closing a DB connection - database.closePostgres()")
	}

	sqlDB.Close()
}

// GetPostgresMasterDSN get Postgres's slave data source name
func getPostgresSlaveDSN() string {
	return writePostgreDSNString(dsnConfig{
		host:     config.POSTGRES_HOST_SLAVE,
		user:     config.POSTGRES_USERNAME_SLAVE,
		password: config.POSTGRES_PASSWORD_SLAVE,
		db:       config.POSTGRES_DATABASE_SLAVE,
		port:     config.POSTGRES_PORT_SLAVE,
		ssl:      config.POSTGRES_SSL_MODE_SLAVE,
	})
}

// writePostgreDSNString write Postgres DSN format string
func writePostgreDSNString(dsn dsnConfig) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%v sslmode=%s", dsn.host, dsn.user, dsn.password, dsn.db, dsn.port, dsn.ssl)
}
