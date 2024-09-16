package httpserver

import (
	_ "github.com/sebastianreh/user-balance-api/docs/swagger"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (s *Server) Routes() {
	s.Server.GET("/ping", s.dependencies.PingHandler.Ping)

	root := s.Server.Group(s.dependencies.Config.Prefix)
	root.GET("/swagger/*", echoSwagger.WrapHandler)
	root.POST("/migrate", s.dependencies.MigrationHandler.UploadMigrationCSV)

	usersGroup := root.Group("/users")
	usersGroup.GET("/:user_id/balance", s.dependencies.BalanceHandler.GetUserBalanceWithOptions)
	usersGroup.POST("/create", s.dependencies.UserHandler.CreateUser)
	usersGroup.PUT("/:id", s.dependencies.UserHandler.UpdateUser)
	usersGroup.DELETE("/:id", s.dependencies.UserHandler.DeleteUser)
	usersGroup.GET("/:id", s.dependencies.UserHandler.GetUser)

	transactionsGroup := root.Group("/transactions")
	transactionsGroup.POST("/create", s.dependencies.TransactionHandler.CreateTransaction)
	transactionsGroup.PUT("/:id", s.dependencies.TransactionHandler.UpdateTransaction)
	transactionsGroup.GET("/:id", s.dependencies.TransactionHandler.GetTransaction)
	transactionsGroup.DELETE("/:id", s.dependencies.TransactionHandler.DeleteTransaction)
}
