package main

import (
	"github.com/go-chassis/go-chassis"
	"github.com/go-mesh/openlogging"
	_ "github.com/huaweicse/auth/adaptor/gochassis"
	_ "github.com/huaweicse/cse-collector"
	"github.com/huaweicse/cse-collector/example/restaurant/resource"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/rpc/server/

func main() {
	chassis.RegisterSchema("rest", &resource.RestaurantResource{})
	if err := chassis.Init(); err != nil {
		openlogging.Error("Init failed." + err.Error())
		return
	}
	chassis.Run()
}
