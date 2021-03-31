package router

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/sckacr/calltaxi/handler"
)

func New() *gin.Engine {
	router := gin.New()

	router.POST("/api/hook/passengers", handler.PassengerHandlerFunc)
	router.POST("/api/hook/operators", handler.OperatorHandlerFunc)

	return router
}
