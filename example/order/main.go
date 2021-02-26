package main

import (
	"github.com/go-chassis/go-chassis/v2"
	"github.com/go-chassis/openlog"
	_ "github.com/go-chassis/go-chassis-cloud/provider/huawei/engine"
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
		openlog.Error("Init failed." + err.Error())
		return
	}

	err := chassis.Run()
	if err != nil {
		panic(err)
	}
}
