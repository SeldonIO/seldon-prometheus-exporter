package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/seldonIO/seldon-prometheus-exporter/query"
	str2duration "github.com/xhit/go-str2duration/v2"
	"log"
	"os"
	"time"
)

//Define the metrics we wish to expose
var modelCpuUsageMetric = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "model_cpu_usage_seconds_total",
		Help: "cpu usage for an ML deployment",
	},
	[]string{"namespace","type","name"},
)

var modelCpuRequestsMetric = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "model_cpu_requests",
		Help: "cpu requests for an ML deployment",
	},
	[]string{"namespace","type","name"},
)

var modelCpuLimitsMetric = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "model_cpu_limits",
		Help: "cpu limits for an ML deployment",
	},
	[]string{"namespace","type","name"},
)

var modelMemoryUsageMetric = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "model_memory_usage_bytes",
		Help: "memory usage for an ML deployment",
	},
	[]string{"namespace","type","name"},
)

var modelMemoryRequestsMetric = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "model_memory_requests_bytes",
		Help: "memory requests for an ML deployment",
	},
	[]string{"namespace","type","name"},
)

var modelMemoryLimitsMetric = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "model_memory_limits_bytes",
		Help: "memory limits for an ML deployment",
	},
	[]string{"namespace","type","name"},
)

var modelContainersMetric = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "model_containers_average",
		Help: "Average running containers for an ML deployment",
	},
	[]string{"namespace","type","name"},
)


func init() {
	//Register metrics with prometheus
	prometheus.MustRegister(modelCpuUsageMetric)
	prometheus.MustRegister(modelCpuRequestsMetric)
	prometheus.MustRegister(modelCpuLimitsMetric)

	prometheus.MustRegister(modelMemoryUsageMetric)
	prometheus.MustRegister(modelMemoryRequestsMetric)
	prometheus.MustRegister(modelMemoryLimitsMetric)

	prometheus.MustRegister(modelContainersMetric)

	// periodically collect metrics
	go func() {
		timePeriod := os.Getenv("TIME_PERIOD")
		if timePeriod == "" {
			timePeriod = "2m"
		}
		duration, err := str2duration.ParseDuration(timePeriod)
		if err != nil {
			panic(err)
		}

		for {
			collectMetricsUsingQueryTemplate(modelCpuUsageMetric,query.CpuUsageSumSeldonTemplate)
			collectMetricsUsingQueryTemplate(modelCpuRequestsMetric,query.CpuRequestSeldonTemplate)
			collectMetricsUsingQueryTemplate(modelCpuLimitsMetric,query.CpuLimitSeldonTemplate)

			collectMetricsUsingQueryTemplate(modelMemoryUsageMetric,query.MemUsageSumSeldonTemplate)
			collectMetricsUsingQueryTemplate(modelMemoryRequestsMetric,query.MemRequestSeldonTemplate)
			collectMetricsUsingQueryTemplate(modelMemoryLimitsMetric,query.MemLimitSeldonTemplate)

			collectMetricsUsingQueryTemplate(modelContainersMetric,query.ContainersUsageSumSeldonTemplate)

			time.Sleep(time.Duration(duration))
		}
	}()
}

func collectMetricsUsingQueryTemplate(metric *prometheus.GaugeVec, queryTemplate string) {

	timePeriod := os.Getenv("TIME_PERIOD")
	if timePeriod == "" {
		timePeriod = "2m"
	}

	var inputData = query.DefaultInputData
	inputData.Range = timePeriod

	queryReturnValues, err := query.ObtainMetricValues(queryTemplate, inputData)
	if err != nil {
		log.Println(err)
	}
	metric.Reset()
	for _, v := range queryReturnValues {
		metric.WithLabelValues(v.Namespace, v.ModelType, v.Model).Set(v.Value)
	}
}