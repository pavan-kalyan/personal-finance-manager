package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-personal-finance/pkg/api"
)

type Server struct {
	router          *gin.Engine
	accountsService *api.AccountsService
}

func NewServer(router *gin.Engine, accountsService *api.AccountsService) *Server {
	return &Server{
		router:          router,
		accountsService: accountsService,
	}
}

func (s *Server) Run() error {

	router := s.InitRoutes()
	err := router.Run()

	if err != nil {
		return fmt.Errorf("Failed to run server with err: %w\n", err)
	}
	return nil
}
