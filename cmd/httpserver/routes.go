package httpserver

func (s *Server) Routes() {
	root := s.Server.Group(s.dependencies.Config.Prefix)
	s.Server.GET("/ping", s.dependencies.PingHandler.Ping)

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
