package utils

import (
	"context"
	"fmt"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"log"
	"os"
	"strconv"
	"time"

	api "github.com/prometheus/client_golang/api"
)

type MetricsInfo struct {
	Job       string
	Pod       string
	Namespace string
	CpuUsage  float64
	MemUsage  int64
	TimeStap  time.Time
}

type SourceLimit struct {
	CpuLimit   float64
	MemLimit   int64
	NameSpace  []string
	SourceType []string
	Job        string
	Dingtalk   string
}

func InArray(value string, arrays []string) bool {
	for _, array := range arrays {
		if value == array {
			return true
		}
	}
	return false
}

// 创建到Prometheus的连接
func (sl *SourceLimit) ClientProm(prometheusURL string) (context.Context, v1.API, context.CancelFunc) {
	//fmt.Println("创建prometheus的连接")
	client, err := api.NewClient(api.Config{Address: prometheusURL})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	return ctx, v1api, cancel
}

// 使用queryrange语句查询所有pod的CPU占用，并根据label调用内存的查询
func (sl *SourceLimit) MetricsCPUValue(ctx context.Context, v1api v1.API, cancel context.CancelFunc) []MetricsInfo {
	defer cancel()
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Runtime err caught: %v", r)
		}
	}()
	podMetrics := make([]MetricsInfo, 0)
	var mi MetricsInfo

	r := v1.Range{
		Start: time.Now().Add(-time.Minute * 5),
		End:   time.Now(),
		Step:  time.Minute,
	}
	result, warnings, err := v1api.QueryRange(ctx, "sum(irate(container_cpu_usage_seconds_total{container!=\"\",container!=\"POD\",job=\""+sl.Job+"\",pod!=\"\"}[5m])) by (namespace,pod)", r, v1.WithTimeout(5*time.Second))
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
	if len(metrix) > 0 {
		for i := 0; i < len(metrix); i++ {
			cpuUsage, err := strconv.ParseFloat(metrix[i].Values[0].Value.String(), 64)
			memUsage := MetricsMemValue(sl.Job, string(metrix[i].Metric["namespace"]), string(metrix[i].Metric["pod"]), ctx, v1api)

			if err != nil {
				fmt.Printf("Error convert cpu or memory value ")
			}
			//fmt.Printf("namespace:%s,cpu:%f,mem:%d\n", string(metrix[i].Metric["namespace"]), cpuUsage, memUsage)
			// 排除掉集群namespace等，可自定义
			if !InArray(string(metrix[i].Metric["namespace"]), sl.NameSpace) {
				if cpuUsage >= sl.CpuLimit || memUsage/1024/1024 >= sl.MemLimit {
					mi.Namespace = string(metrix[i].Metric["namespace"])
					mi.Pod = string(metrix[i].Metric["pod"])
					mi.CpuUsage, err = strconv.ParseFloat(metrix[i].Values[0].Value.String(), 64)
					mi.TimeStap = metrix[i].Values[0].Timestamp.Time()
					mi.MemUsage = memUsage
					podMetrics = append(podMetrics, mi)
				}
			}

		}
	}
	return podMetrics
}

// 使用query语句查询pod的内存占用
func MetricsMemValue(job, namespace, podname string, ctx context.Context, v1api v1.API) int64 {
	memResult, warnings, err := v1api.Query(ctx, "sum(container_memory_working_set_bytes{job=\""+job+"\",namespace=\""+namespace+"\",container!=\"\",container!=\"POD\",pod=\""+podname+"\"}) by (namespace,pod)", time.Now(), v1.WithTimeout(5*time.Second))
	var memUsage int64
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}
	//fmt.Println(job, namespace, podname)
	if len(memResult.(model.Vector)) > 0 {
		memUsage, err = strconv.ParseInt(memResult.(model.Vector)[0].Value.String(), 10, 64)
		if err != nil {
			fmt.Printf("err:%v\n", err)
		}
	}

	return memUsage
}
