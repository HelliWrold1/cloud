package main

import (
	"github.com/HelliWrold1/cloud/cmd/cloud/initial"
	"github.com/zhufuyi/sponge/pkg/app"

	MQTT "github.com/HelliWrold1/cloud/internal/mqtt"
)

// @title cloud api docs
// @description http server api docs
// @schemes http https
// @version v0.0.0
// @host localhost:8080
func main() {
	initial.Config()
	err := MQTT.Init()
	if err != nil {
		return
	}
	defer MQTT.Close()
	servers := initial.RegisterServers()
	closes := initial.RegisterClose(servers)

	a := app.New(servers, closes)
	a.Run()
}
