package main

import (
	"github.com/HelliWrold1/cloud/cmd/cloud/initial"
	"github.com/zhufuyi/sponge/pkg/app"
	"github.com/zhufuyi/sponge/pkg/jwt"
	"github.com/zhufuyi/sponge/pkg/logger"

	MQTT "github.com/HelliWrold1/cloud/internal/mqtt"
)

// @title cloud api docs
// @description http server api docs
// @schemes http https
// @version v0.0.0
// @host localhost:8080
// @securityDefinitions.apiKey		BearerTokenAuth
// @in 								header
// @name 							Authorization
// @description 					Bearer token authentication
func main() {
	initial.Config()
	jwt.Init()
	err := MQTT.Init()
	if err != nil {
		logger.Debug("MQTT Connection error")
	}
	defer MQTT.Close()
	servers := initial.RegisterServers()
	closes := initial.RegisterClose(servers)

	a := app.New(servers, closes)
	a.Run()
}
