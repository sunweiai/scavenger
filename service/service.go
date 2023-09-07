package service

import (
	"fmt"
	"log"
	"os"
	"scavenger/alertmode"
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
	dingTalk         string
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
	envVar.dingTalk = os.Getenv("DINGTALK")

	envVar.namespaceExclude = strings.Split(namespaceExlude, ";")
	envVar.sourceType = strings.Split(sourceType, ";")
	return &envVar
}

func HandMetrics(envVar *varEnv) {
	//clientSet := client.RestClient()
	//kindinfo := new(KindInfo)
	var kubeClient = &utils.SourceLimit{
		MemLimit:   envVar.memLimit,
		CpuLimit:   envVar.cpuLimit,
		NameSpace:  envVar.namespaceExclude,
		SourceType: envVar.sourceType,
		Job:        envVar.job,
	}

	// 获取超过阈值的metrics列表
	ctx, v1api, cancel := kubeClient.ClientProm(envVar.prometheusUrl)
	//mem := utils.MetricsMemValue("k8s-test", "default", "whoami-6cdf669df7-mqjwx", ctx, v1api)
	podMetricsList := kubeClient.MetricsCPUValue(ctx, v1api, cancel)

	fmt.Println(podMetricsList)
	// 根据metrics列表获取到要删除的pod
	if len(podMetricsList) > 0 {
		//delayTime := time.Now()
		for _, podMetrics := range podMetricsList {
			jsonData := `{"msgtype": "markdown","markdown": {"title":"====高负载pod预警====","text":"====高负载pod警告==== \n\n 即将删除资源: \n\n
namespace: ` + podMetrics.Namespace + ` \n\n podname: ` + podMetrics.Pod + `"}}`
			fmt.Printf(jsonData)
			alertmode.RequestDingding(envVar.dingTalk, jsonData)
			//sourceType := kindinfo.GetPodType(clientSet, podMetrics.Namespace, podMetrics.Pod)
			//DeleteSource(clientSet, sourceType.sourceName, sourceType.kind, sourceType.nameSpace)
		}
	}

}
