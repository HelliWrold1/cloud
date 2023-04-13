package main

import (
	"github.com/HelliWrold1/cloud/cmd/frame/initial"
	MQTT "github.com/HelliWrold1/cloud/internal/mqtt"
	"github.com/zhufuyi/sponge/pkg/app"
)

// @title frame api docs
// @description http server api docs
// @schemes http https
// @version v0.0.0
// @host localhost:8080
func main() {
	initial.Config()
	MQTT.Init()
	defer MQTT.Close()
	servers := initial.RegisterServers()
	closes := initial.RegisterClose(servers)

	a := app.New(servers, closes)
	a.Run()
}
