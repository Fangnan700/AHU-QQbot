package main

import (
	"Kira-qbot/initialize"
	"Kira-qbot/server"
)

func main() {
	initialize.Init()

	engine := server.CreateRouterEngine()
	engine.Run("0.0.0.0:5701")
}
