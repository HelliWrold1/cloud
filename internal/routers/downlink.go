package routers

import (
	"github.com/HelliWrold1/cloud/internal/handler"

	"github.com/gin-gonic/gin"
)

func init() {
	routerFns = append(routerFns, func(group *gin.RouterGroup) {
		downlinkRouter(group, handler.NewDownlinkHandler())
	})
}

func downlinkRouter(group *gin.RouterGroup, h handler.DownlinkHandler) {
	group.POST("/downlink", h.Create)
	group.DELETE("/downlink/:id", h.DeleteByID)
	group.POST("/downlinks/delete/ids", h.DeleteByIDs)
	group.PUT("/downlink/:id", h.UpdateByID)
	group.GET("/downlink/:id", h.GetByID)
	group.POST("/downlinks/ids", h.ListByIDs)
	group.POST("/downlinks", h.List)
}
