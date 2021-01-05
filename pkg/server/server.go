package server

import (
	apirule "BlankZhu/Go-DNAT-Server/pkg/api/rule"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/rules/:ID", apirule.GetByID)
	r.GET("/rules", apirule.GetAll)
	r.POST("/rules", apirule.AddOne)
	r.DELETE("/rules/:ID", apirule.DeleteByID)
	r.PATCH("/rules/:ID", apirule.UpdateByID)

	return r
}
