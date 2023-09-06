package main

import (
	"scavenger/service"
	"time"
)

func main() {
	for {
		time.Sleep(15 * time.Second)
		env := service.GetEnv()
		service.HandMetrics(env)
	}
}
