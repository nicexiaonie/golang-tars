package main

import (
	"fmt"

	"github.com/TarsCloud/TarsGo/tars"

	"demo-user/tars-protocol/Demo"
)

func main() {
	comm := tars.NewCommunicator()
	obj := fmt.Sprintf("Demo.UserServer.UserServerObjObj@tcp -h 8.141.7.34 -p 17191 -t 60000")
	app := new(Demo.UserServerObj)
	comm.StringToProxy(obj, app)
	var out, i int32
	i = 123
	ret, err := app.Add(i, i*2, &out)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ret, out)
}
