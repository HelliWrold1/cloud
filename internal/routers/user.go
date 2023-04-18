package routers

import (
	"github.com/HelliWrold1/cloud/internal/handler"
	"github.com/zhufuyi/sponge/pkg/gin/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	routerFns = append(routerFns, func(group *gin.RouterGroup) {
		userRouter(group, handler.NewUserHandler())
	})
}

func userRouter(group *gin.RouterGroup, h handler.UserHandler) {
	group.POST("/user/register", h.Create)
	group.DELETE("/user/:id", middleware.Auth(), h.DeleteByID)
	group.POST("/users/delete/ids", middleware.Auth(), h.DeleteByIDs)
	group.PUT("/user/:id", middleware.Auth(), h.UpdateByID)
	group.GET("/user/:id", middleware.Auth(), h.GetByID)
	group.POST("/users/ids", middleware.Auth(), h.ListByIDs)
	group.POST("/users", middleware.Auth(), h.List)
	group.POST("/user/login", h.LoginUser)
	group.PUT("/user/update", middleware.Auth(), h.UpdateByUserPasswordToNew)
}
