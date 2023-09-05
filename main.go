package main

import (
	"fmt"
	"scavenger/service"
)

func main() {
	env := service.GetEnv()
	fmt.Println("env:%v\n", env)
	service.HandMetrics(env)
}
