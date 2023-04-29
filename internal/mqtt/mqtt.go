package mqtt

import (
	"context"
	"fmt"
	"github.com/HelliWrold1/cloud/internal/cache"
	"github.com/HelliWrold1/cloud/internal/dao"
	"github.com/HelliWrold1/cloud/internal/model"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	jsoniter "github.com/json-iterator/go"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/mysql"
	"time"
)

var frameDao dao.FrameDao
var downlinkDao dao.DownlinkDao

var client MQTT.Client
var uplinkToCloud string = "uplinkToCloud"
var rulesFromCloud string = "rulesFromCloud"
var downlinkToNode string = "downlinkToNode"

type SubTopicInfo struct {
	Topic string
	Qos   byte
}

type PubMsgInfo struct {
	Topic   string `json:"topic"`
	Qos     byte   `json:"qos"`
	Payload string `json:"payload"`
	Retain  bool   `json:"retain"`
}

func Init() error {

	// 获取数据库对象
	frameDao = dao.NewFrameDao(
		model.GetDB(),
		cache.NewFrameCache(model.GetCacheType()),
	)

	downlinkDao = dao.NewDownlinkDao(
		model.GetDB(),
		cache.NewDownlinkCache(model.GetCacheType()),
	)

	opts := MQTT.NewClientOptions().AddBroker("tcp://localhost:1883")

	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(5 * time.Second)
	//opts.SetReconnectingHandler();
	opts.SetClientID("go-mqtt-client")

	opts.SetOnConnectHandler(func(c MQTT.Client) {
		fmt.Println("MQTT connected")

		if token := client.Subscribe(uplinkToCloud, 1, messageReceivedCallback); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
		}
	})

	opts.SetConnectionLostHandler(func(c MQTT.Client, err error) {
		fmt.Println("MQTT connection lost")
	})

	// 创建MQTT客户端
	client = MQTT.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// Publish 发布消息
func Publish(p PubMsgInfo) error {
	if token := client.Publish(p.Topic, p.Qos, p.Retain, p.Payload); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	fmt.Println("Publish topic: " + p.Topic + " Message: " + p.Payload)
	return nil
}

// Subscribe 订阅主题消息
func Subscribe(s SubTopicInfo) error {
	if token := client.Subscribe(s.Topic, s.Qos, messageReceivedCallback); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	fmt.Println("Subscribe topic " + s.Topic)
	return nil
}

// Unsubscribe 取消订阅主题消息
func Unsubscribe(s SubTopicInfo) error {
	if token := client.Unsubscribe(s.Topic); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	fmt.Println("Unsubscribe topic " + s.Topic)
	return nil
}

func Close() {
	client.Disconnect(250)
}

// messageReceivedCallback MQTT消息收到回调
func messageReceivedCallback(client MQTT.Client, message MQTT.Message) {
	fmt.Printf("Received message on topic [%s]: \n%s\n", message.Topic(), message.Payload())

	topic := message.Topic()
	payloadStr := string(message.Payload())
	if topic == uplinkToCloud {
		var obj map[string]interface{}

		err := jsoniter.Unmarshal([]byte(payloadStr), &obj)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		dataType, exist := obj["datatype"].(float64)
		devAddr, _ := obj["devaddr"].(string)
		gatewayMac, _ := obj["mac"].(string)
		utcTimeStr, _ := obj["datetime"].(string)
		utcTime, _ := time.Parse("2006-01-02T15:04:05Z", utcTimeStr)
		if !exist {
			logger.Info("datatype not exist")
			return
		}
		err = frameDao.Create(context.Background(), &model.Frame{
			Model: mysql.Model{
				CreatedAt: utcTime, // 插入的是utcTime, 框架会自动把UTC时间转换为localtime存入数据库
			},
			Frame:      payloadStr,
			DevAddr:    devAddr,
			DataType:   int(dataType),
			GatewayMac: gatewayMac,
		})
		if err != nil {
			return
		}
	}
}
