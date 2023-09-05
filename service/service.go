package service

import (
	"fmt"
	"log"
	"os"
	"scavenger/utils"
	"strconv"
	"strings"
	"time"
)

type varEnv struct {
	cpuLimit         float64
	memLimit         int64
	namespaceExclude []string
	sourceType       []string
	job              string
	prometheusUrl    string
	intervalTime     time.Duration
}

func GetEnv() *varEnv {
	var err error
	var envVar varEnv
	envVar.cpuLimit, err = strconv.ParseFloat(os.Getenv("CPULIMIT"), 64)
	if err != nil {
		log.Fatal("Get env int error!")
	}
	envVar.memLimit, err = strconv.ParseInt(os.Getenv("MEMLIMIT"), 10, 64)
	if err != nil {
		log.Fatal("Get env int error!")
	}
	namespaceExlude := os.Getenv("NAMESPACE")
	sourceType := os.Getenv("SOURCETYPE")
	convertExpire, err := strconv.ParseInt(os.Getenv("INTERVALTIME"), 10, 64)
	envVar.intervalTime = time.Duration(convertExpire)
	envVar.job = os.Getenv("JOB")
	envVar.prometheusUrl = os.Getenv("URL")

	envVar.namespaceExclude = strings.Split(namespaceExlude, ";")
	envVar.sourceType = strings.Split(sourceType, ";")
	return &envVar
}

func HandMetrics(envVar *varEnv) {
	//clientSet := client.RestClient()
	var kubeClient = &utils.SourceLimit{
		MemLimit:   envVar.memLimit,
		CpuLimit:   envVar.cpuLimit,
		NameSpace:  envVar.namespaceExclude,
		SourceType: envVar.sourceType,
		Job:        envVar.job,
	}

	// 获取超过阈值的metrics列表
	ctx, v1api := kubeClient.ClientProm(envVar.prometheusUrl)
	podMetricsList := kubeClient.MetricsCPUValue(ctx, v1api)

	fmt.Println(podMetricsList)
	// 根据metrics列表获取到要删除的pod
	//for _, podMetrics = range podMetricsList {
	//
	//}
}
