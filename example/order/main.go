package main

import (
	"github.com/go-chassis/go-chassis"
	"github.com/go-mesh/openlogging"
	_ "github.com/huaweicse/auth/adaptor/gochassis"
	_ "github.com/huaweicse/cse-collector"
	"github.com/huaweicse/cse-collector/example/order/resource"
	"log"
	"os"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/rest/client/

func main() {
	envs := os.Environ()
	for _, e := range envs {
		log.Println(e)
	}
	chassis.RegisterSchema("rest", &resource.OrderResource{})
	//Init framework
	if err := chassis.Init(); err != nil {
		openlogging.Error("Init failed." + err.Error())
		return
	}

	err := chassis.Run()
	if err != nil {
		panic(err)
	}
}
