package routers

import (
	"github.com/HelliWrold1/cloud/internal/handler"
	"github.com/zhufuyi/sponge/pkg/gin/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	routerFns = append(routerFns, func(group *gin.RouterGroup) {
		downlinkRouter(group, handler.NewDownlinkHandler())
	})
}

func downlinkRouter(group *gin.RouterGroup, h handler.DownlinkHandler) {
	group.POST("/downlink", middleware.AuthAdmin(), h.Create)
	group.DELETE("/downlink/:id", middleware.AuthAdmin(), h.DeleteByID)
	group.POST("/downlinks/delete/ids", middleware.AuthAdmin(), h.DeleteByIDs)
	group.PUT("/downlink/:id", middleware.AuthAdmin(), h.UpdateByID)
	group.GET("/downlink/:id", middleware.Auth(), h.GetByID)
	group.POST("/downlinks/ids", middleware.Auth(), h.ListByIDs)
	group.POST("/downlinks", middleware.Auth(), h.List)
}
