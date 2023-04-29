package handler

import (
	"context"
	"fmt"
	"github.com/HelliWrold1/cloud/internal/cache"
	"github.com/HelliWrold1/cloud/internal/dao"
	"github.com/HelliWrold1/cloud/internal/ecode"
	"github.com/HelliWrold1/cloud/internal/model"
	MQTT "github.com/HelliWrold1/cloud/internal/mqtt"
	"github.com/HelliWrold1/cloud/internal/types"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	jsoniter "github.com/json-iterator/go"
	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/mysql"
	"time"
)

var _ MQTTHandler = (*mqttHandler)(nil)

type MQTTHandler interface {
	Publish(ctx *gin.Context)
	Subscribe(ctx *gin.Context)
	Unsubscribe(ctx *gin.Context)
}

type mqttHandler struct {
	iDao           dao.DownlinkDao
	rulesFromCloud string
	downlinkToNode string
}

func NewMQTTHandler() MQTTHandler {
	return &mqttHandler{
		iDao: dao.NewDownlinkDao(
			model.GetDB(),
			cache.NewDownlinkCache(model.GetCacheType()),
		),
		rulesFromCloud: "rulesFromCloud",
		downlinkToNode: "downlinkToNode",
	}
}

// Publish a MQTT message
// @Summary publish mqtt message
// @Description submit information to publish mqtt message
// @Tags MQTT
// @accept json
// @Produce json
// @Security BearerTokenAuth
// @Param data body types.MQTTPublishRequest true "mqtt information"
// @Success 200 {object} types.Result{}
// @Router /api/v1/mqtt/publish [post]
func (h *mqttHandler) Publish(c *gin.Context) {
	form := &types.MQTTPublishRequest{}

	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	pubMsgInfo := &MQTT.PubMsgInfo{}
	err = copier.Copy(pubMsgInfo, form)

	// 解析要发布的json数据
	var obj map[string]interface{}
	err = jsoniter.Unmarshal([]byte(pubMsgInfo.Payload), &obj)
	if err != nil {
		response.Error(c, ecode.ErrJSONMQTT)
		return
	}

	// 发布MQTT消息
	err = MQTT.Publish(*pubMsgInfo)
	fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05"))
	fmt.Println(time.Now().Local().Format("2006-01-02 15:04:05"))

	// 向数据库插入发布的命令
	if pubMsgInfo.Topic == h.downlinkToNode {
		devAddr, _ := obj["devaddr"].(string)
		err = h.iDao.Create(context.Background(), &model.Downlink{
			Model: mysql.Model{
				CreatedAt: time.Now().Local(),
			},
			DownLink: pubMsgInfo.Payload,
			DevAddr:  devAddr,
		})
	}

	// 向数据库插入当前规则
	if pubMsgInfo.Topic == h.rulesFromCloud {
		err = h.iDao.Create(context.Background(), &model.Downlink{
			Model: mysql.Model{
				CreatedAt: time.Now().UTC(),
			},
			DownLink: pubMsgInfo.Payload,
		})
	}

	if err != nil {
		response.Error(c, ecode.ErrPublishMQTT)
		return
	}
	response.Success(c)
}

// Subscribe a MQTT topic
// @Summary subscribe mqtt topic
// @Description submit information to subscribe mqtt message
// @Tags MQTT
// @accept json
// @Produce json
// @Security BearerTokenAuth
// @Param data body types.MQTTSubscribeRequest true "mqtt information"
// @Success 200 {object} types.Result{}
// @Router /api/v1/mqtt/subscribe [post]
func (h *mqttHandler) Subscribe(c *gin.Context) {
	form := &types.MQTTSubscribeRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	subTopicInfo := &MQTT.SubTopicInfo{}
	err = copier.Copy(subTopicInfo, form)
	err = MQTT.Subscribe(*subTopicInfo)
	if err != nil {
		response.Error(c, ecode.ErrSubscribeMQTT)
		return
	}
	response.Success(c)
}

// Unsubscribe a MQTT topic
// @Summary unsubscribe mqtt topic
// @Description submit information to unsubscribe mqtt message
// @Tags MQTT
// @accept json
// @Produce json
// @Security BearerTokenAuth
// @Param data body types.MQTTUnsubscribeRequest true "mqtt information"
// @Success 200 {object} types.Result{}
// @Router /api/v1/mqtt/unsubscribe [post]
func (h *mqttHandler) Unsubscribe(c *gin.Context) {
	form := &types.MQTTUnsubscribeRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	subTopicInfo := &MQTT.SubTopicInfo{}
	err = copier.Copy(subTopicInfo, form)
	err = MQTT.Unsubscribe(*subTopicInfo)
	if err != nil {
		response.Error(c, ecode.ErrUnsubscribeMQTT)
		return
	}
	response.Success(c)
}
