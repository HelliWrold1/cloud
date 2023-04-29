// nolint

package ecode

import "github.com/zhufuyi/sponge/pkg/errcode"

var (
	MQTTNO       = 83
	MQTTName     = "mqtt"
	MQTTBaseCode = errcode.HCode(MQTTNO)

	ErrPublishMQTT     = errcode.NewError(MQTTBaseCode+1, "failed to publish "+MQTTName)
	ErrJSONMQTT        = errcode.NewError(MQTTBaseCode+2, "error json mqtt payload ")
	ErrSubscribeMQTT   = errcode.NewError(MQTTBaseCode+3, "failed to subscribe")
	ErrUnsubscribeMQTT = errcode.NewError(MQTTBaseCode+4, "failed to unsubscribe")
	// for each error code added, add +1 to the previous error code
)
