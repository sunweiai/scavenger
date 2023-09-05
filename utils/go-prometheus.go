package utils

import (
	"context"
	"fmt"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"os"
	"strconv"
	"time"

	api "github.com/prometheus/client_golang/api"
)

type MetricsInfo struct {
	job       string
	pod       string
	namespace string
	cpuUsage  float64
	memUsage  int64
}

type SourceLimit struct {
	CpuLimit   float64
	MemLimit   int64
	NameSpace  []string
	SourceType []string
	Job        string
}

func InArray(value string, arrays []string) bool {
	for _, array := range arrays {
		if value == array {
			return true
		}
	}
	return false
}

func (sl *SourceLimit) ClientProm(prometheusURL string) (context.Context, v1.API) {
	fmt.Println("创建prometheus的连接")
	client, err := api.NewClient(api.Config{Address: prometheusURL})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return ctx, v1api
}

func (sl *SourceLimit) MetricsCPUValue(ctx context.Context, v1api v1.API) []MetricsInfo {
	//cpuch := make(chan metricsInfo)
	podMetrics := make([]MetricsInfo, 0)
	var mi MetricsInfo

	r := v1.Range{
		Start: time.Now().Add(-time.Minute * 5),
		End:   time.Now(),
		Step:  time.Minute,
	}
	fmt.Printf(sl.Job)
	result, warnings, err := v1api.QueryRange(ctx, "sum(irate(container_cpu_usage_seconds_total{container!=\"POD\",job=\"k8s-test\"}[5m])) by (namespace, pod)", r, v1.WithTimeout(5*time.Second))
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}
	metrix, ok := result.(model.Matrix)
	if !ok {
		fmt.Printf("查询不是矩阵类型")
	}
	//fmt.Printf("matrix:%v\n", metrix.String())
	for i := 0; i < len(metrix); i++ {
		//fmt.Printf("value:%v,metric:%v\n", metrix[i].Values[0], metrix[i].Metric)

		cpuUsage, err := strconv.ParseFloat(metrix[i].Values[0].Value.String(), 64)
		memUsage := MetricsMemValue(sl.Job, string(metrix[i].Metric["namespace"]), string(metrix[i].Metric["pod"]), ctx, v1api)

		if err != nil {
			fmt.Printf("Error convert cpu or memory value ")
		}
		fmt.Printf("cpu:%f,mem:%d\n", cpuUsage, memUsage)
		if InArray(string(metrix[i].Metric["namespace"]), sl.NameSpace) {
			if cpuUsage >= sl.CpuLimit || memUsage/1024/1024 >= sl.MemLimit {
				mi.namespace = string(metrix[i].Metric["namespace"])
				mi.pod = string(metrix[i].Metric["pod"])
				mi.cpuUsage, err = strconv.ParseFloat(metrix[i].Values[0].Value.String(), 64)
				mi.memUsage = memUsage
				podMetrics = append(podMetrics, mi)

			}
		}

	}
	return podMetrics
}

func MetricsMemValue(job, namespace, podname string, ctx context.Context, v1api v1.API) int64 {
	memResult, warnings, err := v1api.Query(ctx, "sum(container_memory_working_set_bytes{job=\""+job+"\",namespace=\""+namespace+"\",container!=\"\",container!=\"POD\",pod=\""+podname+"\"}) by (namespace,pod)", time.Now(), v1.WithTimeout(5*time.Second))
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}
	memUsage, err := strconv.ParseInt(memResult.(model.Vector)[0].Value.String(), 10, 64)
	return memUsage
}
