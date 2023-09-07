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
	//alertmode.RequestDingding("https://oapi.dingtalk.com/robot/send?access_token=5ee47bfc6293f6bb41cb8f3bebb098a6187d3b6b307280b3420c02aeac8fc0dd")
}
