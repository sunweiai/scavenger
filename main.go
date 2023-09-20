package main

import (
	"scavenger/client"
	"scavenger/service"
	"time"
)

func main() {
	for {
		env := service.GetEnv()
		metricsList := service.GetMetrics(env)
		clientSet := client.RestClient()
		service.HandService(clientSet, metricsList)
		time.Sleep(15 * time.Second)
	}
}
