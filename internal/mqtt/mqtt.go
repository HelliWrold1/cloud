package mqtt

import (
	"context"
	"fmt"
	"github.com/HelliWrold1/cloud/internal/cache"
	"github.com/HelliWrold1/cloud/internal/dao"
	"github.com/HelliWrold1/cloud/internal/model"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	jsoniter "github.com/json-iterator/go"
	"github.com/zhufuyi/sponge/pkg/mysql"
	"os"
	"time"
)

var client MQTT.Client
var uplinkToCloud string = "uplinkToCloud"
var downlinkToNode string = "downlinkToNode"
var frameDao dao.FrameDao

func Init() {

	// 获取数据库对象
	frameDao = dao.NewFrameDao(
		model.GetDB(),
		cache.NewFrameCache(model.GetCacheType()),
	)

	opts := MQTT.NewClientOptions().AddBroker("tcp://localhost:1883")

	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(5 * time.Second)
	//opts.SetReconnectingHandler();
	opts.SetClientID("go-mqtt-client")

	opts.SetOnConnectHandler(func(c MQTT.Client) {
		fmt.Println("MQTT connected")
	})

	opts.SetConnectionLostHandler(func(c MQTT.Client, err error) {
		fmt.Println("MQTT connection lost")
	})

	// 创建MQTT客户端
	client = MQTT.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := client.Subscribe(uplinkToCloud, 1, messageReceivedCallback); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}

// PublishCmd 发布指令消息
func PublishCmd(payload string) error {
	if token := client.Publish(downlinkToNode, 1, false, payload); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func Close() {
	client.Disconnect(250)
}

// messageReceivedCallback MQTT消息收到回调
func messageReceivedCallback(client MQTT.Client, message MQTT.Message) {
	fmt.Printf("Received message on topic [%s]: \n%s\n", message.Topic(), message.Payload())
	if message.Topic() == uplinkToCloud {
		frame := string(message.Payload())
		var obj map[string]interface{}

		err := jsoniter.Unmarshal([]byte(frame), &obj)
		if err != nil {
			fmt.Println(err.Error())
		}
		dateType, _ := obj["datetype"].(int)
		devAddr, _ := obj["devaddr"].(string)
		gatewayMac, _ := obj["mac"].(string)
		//localTimeStr, _ := obj["localtime"].(string)
		//localTime, _ := time.Parse("2006-01-02 15:04:05", localTimeStr)
		utcTimeStr, _ := obj["datetime"].(string)
		utcTime, _ := time.Parse("2006-01-02T15:04:05Z", utcTimeStr)
		frameDao.Create(context.Background(), &model.Frame{
			Model: mysql.Model{
				//CreatedAt: localTime, // 插入的是localTime
				CreatedAt: utcTime, // 插入的是utcTime
			},
			Frame:      frame,
			DevAddr:    devAddr,
			DataType:   dateType,
			GatewayMac: gatewayMac,
		})
	}
}
