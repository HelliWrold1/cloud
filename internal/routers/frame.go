package routers

import (
	"github.com/HelliWrold1/cloud/internal/handler"

	"github.com/gin-gonic/gin"
)

func init() {
	routerFns = append(routerFns, func(group *gin.RouterGroup) {
		frameRouter(group, handler.NewFrameHandler())
	})
}

func frameRouter(group *gin.RouterGroup, h handler.FrameHandler) {
	group.POST("/frame", h.Create)
	group.DELETE("/frame/:id", h.DeleteByID)
	group.POST("/frames/delete/ids", h.DeleteByIDs)
	group.PUT("/frame/:id", h.UpdateByID)
	group.GET("/frame/:id", h.GetByID)
	group.POST("/frames/ids", h.ListByIDs)
	group.POST("/frames", h.List)
}
