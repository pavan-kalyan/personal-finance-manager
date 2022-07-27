package app

import "github.com/gin-gonic/gin"

func (s *Server) InitRoutes() *gin.Engine {
	router := s.router

	v1 := router.Group("/v1/api")
	{
		//accounts
		v1.GET("/accounts/:id", s.GetAccount())
		v1.GET("/accounts", s.ApiStatus())
		v1.DELETE("/accounts/:id", s.DeleteAccount())
		v1.POST("/accounts", s.AddAccount())
		v1.PUT("/accounts/:id", s.UpdateAccount())

		//organizations
		v1.GET("/organizations", s.ApiStatus())

		//health check
		v1.GET("/ping", s.ApiStatus())
	}

	return router
}
