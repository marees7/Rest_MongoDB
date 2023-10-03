package controllers

func (s *Server) InitializeRoutes() {
	s.Router.POST("/user", s.CreateUser)
	s.Router.POST("/login", s.UserLogin)
	s.Router.GET("/user/me", s.GetLoggedInUser)
}
