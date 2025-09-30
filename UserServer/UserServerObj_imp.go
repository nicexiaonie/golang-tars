package main

import (
	"context"
	"fmt"
)

// UserServerObjImp servant implementation
type UserServerObjImp struct {
}

func (imp *UserServerObjImp) EchoHello(ctx context.Context, name string, greeting *string) (int32, error) {
	*greeting = "Hello, " + name
	return 0, nil
}

// Init servant init
func (imp *UserServerObjImp) Init() error {
	//initialize servant here:
	//...
	fmt.Println("UserServerObjImp Init")
	return nil
}

// Destroy servant destroy
func (imp *UserServerObjImp) Destroy() {
	//destroy servant here:
	//...
}

func (imp *UserServerObjImp) Add(ctx context.Context, a int32, b int32, c *int32) (int32, error) {
	//Doing something in your function
	//...
	fmt.Println("Add", a, b, c)
	return 0, nil
}
func (imp *UserServerObjImp) Sub(ctx context.Context, a int32, b int32, c *int32) (int32, error) {
	//Doing something in your function
	//...
	return 0, nil
}
