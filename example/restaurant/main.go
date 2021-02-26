package main

import (
	"github.com/go-chassis/go-chassis/v2"
	"github.com/go-chassis/openlog"
	_ "github.com/go-chassis/go-chassis-cloud/provider/huawei/engine"
	_ "github.com/huaweicse/cse-collector"
	"github.com/huaweicse/cse-collector/example/restaurant/resource"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/rpc/server/

func main() {
	chassis.RegisterSchema("rest", &resource.RestaurantResource{})
	if err := chassis.Init(); err != nil {
		openlog.Error("Init failed." + err.Error())
		return
	}
	chassis.Run()
}
