package main

import (
	"fmt"
	"os"

	"github.com/TarsCloud/TarsGo/tars"

	"demo-user/tars-protocol/Demo"
)

func main() {
	// Get server config
	cfg := tars.GetServerConfig()

	// New servant imp
	imp := new(UserServerObjImp)
	err := imp.Init()
	if err != nil {
		fmt.Printf("UserServerObjImp init fail, err:(%s)\n", err)
		os.Exit(-1)
	}
	// New servant
	app := new(Demo.UserServerObj)
	// Register Servant
	app.AddServantWithContext(imp, cfg.App+"."+cfg.Server+".UserServerObj")

	// Run application
	tars.Run()
}
