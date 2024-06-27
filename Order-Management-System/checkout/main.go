package main

import (
	"runtime"

	"github.com/sofiukl/oms/oms-checkout/core"
	"github.com/sofiukl/oms/oms-checkout/dispatcher"
)

// numDispatcher denotes no of dispatcher
const numDispatcher = 4

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	app := core.App{}
	app.Initialize()
	dispatcher.StartDispatcher(numDispatcher)
	app.Run(":" + app.Config.ServerPort)
}
