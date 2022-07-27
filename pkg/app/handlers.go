package app

import (
	gin "github.com/gin-gonic/gin"
	"go-personal-finance/pkg/api"
	"net/http"
	"strconv"
)

func (s *Server) ApiStatus() gin.HandlerFunc {
	return func(context *gin.Context) {

		response := map[string]string{
			"status": "success",
			"data":   "result",
		}
		context.JSON(http.StatusOK, response)
	}
}
func (s *Server) DeleteAccount() gin.HandlerFunc {
	return func(context *gin.Context) {

		id, err := strconv.Atoi(context.Param("id"))
		response := gin.H{
			"status": "failure",
		}
		if err != nil {
			context.JSON(http.StatusBadRequest, response)
			return
		}
		err = s.accountsService.DeleteAccount(int64(id))
		response = gin.H{
			"status": "success",
		}
		context.JSON(http.StatusOK, response)
	}
}

func (s *Server) GetAccount() gin.HandlerFunc {
	return func(context *gin.Context) {

		id, err := strconv.Atoi(context.Param("id"))
		response := gin.H{
			"status": "failure",
		}
		if err != nil {
			context.JSON(http.StatusBadRequest, response)
			return
		}
		account, err := s.accountsService.GetAccount(int64(id))
		response = gin.H{
			"status": "success",
			"data":   account,
		}
		context.JSON(http.StatusOK, response)
	}
}
func (s *Server) AddAccount() gin.HandlerFunc {
	return func(context *gin.Context) {

		accountCreationRequest := api.AccountCreationRequest{}
		err := context.BindJSON(&accountCreationRequest)
		response := gin.H{
			"status": "failure",
		}
		if err != nil {
			context.JSON(http.StatusInternalServerError, response)
			return
		}
		err = s.accountsService.CreateAccount(accountCreationRequest)
		if err != nil {
			context.JSON(http.StatusInternalServerError, response)
			return
		}
		response = gin.H{
			"status": "success",
		}
		context.JSON(http.StatusOK, response)
	}
}

func (s *Server) UpdateAccount() gin.HandlerFunc {
	return func(context *gin.Context) {

		id, err := strconv.Atoi(context.Param("id"))
		accountUpdateRequest := api.AccountUpdateRequest{}
		err = context.BindJSON(&accountUpdateRequest)
		response := gin.H{
			"status": "failure",
		}
		if err != nil {
			context.JSON(http.StatusInternalServerError, response)
			return
		}
		err = s.accountsService.UpdateAccount(int64(id), accountUpdateRequest)
		response = gin.H{
			"status": "success",
		}
		context.JSON(http.StatusOK, response)
	}
}
