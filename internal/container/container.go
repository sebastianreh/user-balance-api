package container

import (
	"database/sql"
	"github.com/sebastianreh/user-balance-api/internal/app/services"
	"github.com/sebastianreh/user-balance-api/internal/domain/balance"
	"github.com/sebastianreh/user-balance-api/internal/infrastructure/config"
	"github.com/sebastianreh/user-balance-api/internal/infrastructure/postgre_sql"
	"github.com/sebastianreh/user-balance-api/internal/interfaces/http"
	"github.com/sebastianreh/user-balance-api/pkg/csv"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
)

type Dependencies struct {
	Config           config.Config
	Logs             logger.Logger
	SQL              *sql.DB
	PingHandler      *http.PingHandler
	BalanceHandler   *http.BalanceHandler
	MigrationHandler *http.MigrationHandler
}

func Build() Dependencies {
	dependencies := Dependencies{}
	dependencies.Config = config.NewConfig()
	logs := logger.NewLogger()
	dependencies.Logs = logs
	dependencies.PingHandler = http.NewPingHandler(dependencies.Config)

	pgDb, err := postgre_sql.NewPostgresDB(dependencies.Config, dependencies.Logs)
	if err != nil {
		logs.Fatal("Database initialization error, shutting down server")
	}

	sqlMigrations := postgre_sql.NewSqlMigrations(logs, pgDb)
	err = sqlMigrations.RunMigrations()
	if err != nil {
		logs.Fatal("Database migration error, shutting down server")
	}

	dependencies.SQL = pgDb

	userSQLRepository := postgre_sql.NewSqlUserRepository(dependencies.Logs, dependencies.SQL)
	transactionSQLRepository := postgre_sql.NewSqlTransactionRepository(dependencies.Logs, dependencies.SQL)

	balanceCalculator := balance.NewBalanceCalculator()

	csvProcessor := csv.NewCsvProcessor()

	balanceService := services.NewBalanceService(dependencies.Logs, userSQLRepository,
		transactionSQLRepository, balanceCalculator)
	migrationService := services.NewMigrationService(dependencies.Logs, userSQLRepository,
		transactionSQLRepository, csvProcessor)

	dependencies.BalanceHandler = http.NewBalanceHandler(dependencies.Logs, balanceService)
	dependencies.MigrationHandler = http.NewMigrationHandler(dependencies.Logs, migrationService)

	return dependencies
}
