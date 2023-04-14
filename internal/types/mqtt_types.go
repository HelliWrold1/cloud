package types

type MQTTPublishRequest struct {
	Topic   string `json:"topic"`
	Qos     byte   `json:"qos"`
	Payload string `json:"payload"`
	Retain  bool   `json:"retain"`
}

type MQTTPublishResponse struct {
}
