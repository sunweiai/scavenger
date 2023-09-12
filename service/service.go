package service

import (
	"k8s.io/client-go/kubernetes"
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
	//dingTalk         string
	intervalTime time.Duration
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
	//envVar.dingTalk = os.Getenv("DINGTALK")

	envVar.namespaceExclude = strings.Split(namespaceExlude, ";")
	envVar.sourceType = strings.Split(sourceType, ";")
	return &envVar
}

func GetMetrics(envVar *varEnv) []utils.MetricsInfo {
	var kubeClient = &utils.SourceLimit{
		MemLimit:   envVar.memLimit,
		CpuLimit:   envVar.cpuLimit,
		NameSpace:  envVar.namespaceExclude,
		SourceType: envVar.sourceType,
		Job:        envVar.job,
		//Dingtalk:   envVar.dingTalk,
	}

	// 获取超过阈值的metrics列表
	ctx, v1api, cancel := kubeClient.ClientProm(envVar.prometheusUrl)
	//mem := utils.MetricsMemValue("k8s-test", "default", "whoami-6cdf669df7-mqjwx", ctx, v1api)
	sourceMetricsList := kubeClient.MetricsCPUValue(ctx, v1api, cancel)

	return sourceMetricsList
}
func SendMes(dingtalkUrl, bodydata string) {
	if len(bodydata) > 0 {
		alertmode.RequestDingding(dingtalkUrl, bodydata)
	}
}

func HandMes(metrics []utils.MetricsInfo) (string, string) {
	dingtalk := alertmode.UnmarshalJson("E:\\gitee_code\\go_project\\scavenger\\dingtalk-mes.json")
	dingbody, url := alertmode.MarkdownBody(dingtalk, metrics)
	return dingbody, url
}

func HandService(clientSet *kubernetes.Clientset, sourceMetrics []utils.MetricsInfo) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Runtime err caught: %v", r)
		}
	}()
	kindInfo := new(KindInfo)
	if len(sourceMetrics) > 0 {
		for _, metrics := range sourceMetrics {
			sourceType := kindInfo.GetPodType(clientSet, metrics.Namespace, metrics.Pod)
			mesData, url := HandMes(sourceMetrics)
			SendMes(url, mesData)
			DeleteSource(clientSet, sourceType.sourceName, sourceType.kind, sourceType.nameSpace)
		}
	}
}
