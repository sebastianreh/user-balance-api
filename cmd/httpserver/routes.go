package httpserver

func (s *Server) Routes() {
	root := s.Server.Group(s.dependencies.Config.Prefix)
	s.Server.GET("/ping", s.dependencies.PingHandler.Ping)

	root.POST("/migrate", s.dependencies.MigrationHandler.UploadMigrationCSV)
	usersGroup := root.Group("/users")
	usersGroup.GET("/:user_id/balance", s.dependencies.BalanceHandler.GetUserBalanceWithOptions)
}
