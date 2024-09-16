package container

import (
	"database/sql"

	"github.com/sebastianreh/user-balance-api/internal/app/services"
	"github.com/sebastianreh/user-balance-api/internal/domain/balance"
	"github.com/sebastianreh/user-balance-api/internal/infrastructure/config"
	"github.com/sebastianreh/user-balance-api/internal/infrastructure/postgresql"
	"github.com/sebastianreh/user-balance-api/internal/interfaces/http"
	"github.com/sebastianreh/user-balance-api/pkg/csv"
	"github.com/sebastianreh/user-balance-api/pkg/email"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
)

type Dependencies struct {
	Config             config.Config
	Logs               logger.Logger
	SQL                *sql.DB
	PingHandler        *http.PingHandler
	UserHandler        *http.UserHandler
	TransactionHandler *http.TransactionHandler
	BalanceHandler     *http.BalanceHandler
	MigrationHandler   *http.MigrationHandler
}

func Build() Dependencies {
	var dependencies Dependencies
	dependencies.Config = config.NewConfig()
	logs := logger.NewLogger()
	dependencies.Logs = logs
	dependencies.PingHandler = http.NewPingHandler(dependencies.Config)

	pgDB, err := postgresql.NewPostgresDB(dependencies.Config, dependencies.Logs)
	if err != nil {
		logs.Fatal("Database initialization error, shutting down server")
	}

	sqlMigrations := postgresql.NewSQLMigrations(logs, pgDB)
	err = sqlMigrations.RunMigrations()
	if err != nil {
		logs.Fatal("Database migration error, shutting down server")
	}

	dependencies.SQL = pgDB

	userSQLRepository := postgresql.NewSQLUserRepository(dependencies.Logs, dependencies.SQL)
	transactionSQLRepository := postgresql.NewSQLTransactionRepository(dependencies.Logs, dependencies.SQL)

	balanceCalculator := balance.NewBalanceCalculator()

	smtpConfig := dependencies.Config.SMTP
	csvProcessor := csv.NewCsvProcessor()
	emailService := email.NewSMTPEmailService(smtpConfig.Username, smtpConfig.Password, smtpConfig.From, smtpConfig.SendTo,
		smtpConfig.Host, smtpConfig.Port)
	userService := services.NewUserService(dependencies.Logs, userSQLRepository)
	transactionService := services.NewTransactionService(dependencies.Logs, transactionSQLRepository)
	balanceService := services.NewBalanceService(dependencies.Logs, userSQLRepository,
		transactionSQLRepository, balanceCalculator)
	migrationService := services.NewMigrationService(dependencies.Config, dependencies.Logs, userSQLRepository,
		transactionSQLRepository, csvProcessor)
	migrationsReportService := services.NewMigrationReportService(dependencies.Logs, emailService)

	dependencies.UserHandler = http.NewUserHandler(dependencies.Logs, userService)
	dependencies.TransactionHandler = http.NewTransactionHandler(dependencies.Logs, transactionService)
	dependencies.BalanceHandler = http.NewBalanceHandler(dependencies.Logs, balanceService)
	dependencies.MigrationHandler = http.NewMigrationHandler(dependencies.Logs, migrationService, migrationsReportService)

	return dependencies
}
