package routers

import (
	"github.com/HelliWrold1/cloud/internal/handler"

	"github.com/gin-gonic/gin"
)

func init() {
	routerFns = append(routerFns, func(group *gin.RouterGroup) {
		MQTTRouter(group, handler.NewMQTTHandler())
	})
}

func MQTTRouter(group *gin.RouterGroup, h handler.MQTTHandler) {
	group.POST("/mqtt/publish", h.Publish)
}
