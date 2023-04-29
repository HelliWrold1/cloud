package routers

import (
	"github.com/HelliWrold1/cloud/internal/handler"
	"github.com/zhufuyi/sponge/pkg/gin/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	routerFns = append(routerFns, func(group *gin.RouterGroup) {
		frameRouter(group, handler.NewFrameHandler())
	})
}

func frameRouter(group *gin.RouterGroup, h handler.FrameHandler) {
	group.POST("/frame", middleware.AuthAdmin(), h.Create)
	group.DELETE("/frame/:id", middleware.AuthAdmin(), h.DeleteByID)
	group.POST("/frames/delete/ids", middleware.AuthAdmin(), h.DeleteByIDs)
	group.PUT("/frame/:id", middleware.AuthAdmin(), h.UpdateByID)
	group.GET("/frame/:id", middleware.Auth(), h.GetByID)
	group.POST("/frames/ids", middleware.Auth(), h.ListByIDs)
	group.POST("/frames", middleware.Auth(), h.List)
}
