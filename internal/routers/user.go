package routers

import (
	"github.com/HelliWrold1/cloud/internal/handler"

	"github.com/gin-gonic/gin"
)

func init() {
	routerFns = append(routerFns, func(group *gin.RouterGroup) {
		userRouter(group, handler.NewUserHandler())
	})
}

func userRouter(group *gin.RouterGroup, h handler.UserHandler) {
	group.POST("/user", h.Create)
	group.DELETE("/user/:id", h.DeleteByID)
	group.POST("/users/delete/ids", h.DeleteByIDs)
	group.PUT("/user/:id", h.UpdateByID)
	group.GET("/user/:id", h.GetByID)
	group.POST("/users/ids", h.ListByIDs)
	group.POST("/users", h.List)
	group.PUT("/user/update", h.UpdateByUsernamePasswordToNew)
}
