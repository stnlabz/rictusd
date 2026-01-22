package main

import (
	"rictusd/api"
	"rictusd/core"
)

const (
	logPath  = "logs/rictusd.log"
	digitdIP = "192.168.20.102"
)

func main() {
	// 1. Initialize the Environment via Core
	core.InitLogger(logPath)

	// 2. Start the Engine
	enforcer := core.NewEngine(digitdIP)
	
	// 3. Start the API Router
	go api.StartServer(enforcer, "8080")

	enforcer.Log("RictusD Engine and Logger initialized via Core.")

	select {}
}
