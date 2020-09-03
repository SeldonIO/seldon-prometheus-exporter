package query

import (
	"bytes"
	"log"
	"text/template"
)
type QueryInputData struct {
	Range string `json:",omitempty"`
	OffsetExp string `json:",omitempty"`
}

var DefaultInputData = 	QueryInputData{
Range: `1d`,
OffsetExp: ``,
}

	// seldon usage monitoring metrics

var ContainersUsageSumSeldonTemplate = `sum by (label_seldon_app,namespace) ((sum_over_time(kube_pod_labels{label_app_kubernetes_io_managed_by=~"seldon-core"}[{{.Range}}] {{.OffsetExp}}) / scalar(max(sum_over_time(kube_pod_labels[{{.Range}}] {{.OffsetExp}})))) * on(pod,namespace) group_right(label_seldon_app) max by (namespace,pod,container,namespace) (avg_over_time(kube_pod_container_info[{{.Range}}] {{.OffsetExp}})))`
var MemUsageSumSeldonTemplate = `sort_desc(sum by (label_seldon_app,namespace) ((sum_over_time(kube_pod_labels{label_app_kubernetes_io_managed_by=~"seldon-core"}[{{.Range}}] {{.OffsetExp}}) / scalar(max(sum_over_time(kube_pod_labels[{{.Range}}] {{.OffsetExp}})))) * on(pod,namespace) group_right(label_seldon_app) sum by (namespace,pod,container) (rate(container_memory_usage_bytes[{{.Range}}] {{.OffsetExp}}))))`
var CpuUsageSumSeldonTemplate = `sort_desc(sum by (label_seldon_app,namespace) ((sum_over_time(kube_pod_labels{label_app_kubernetes_io_managed_by=~"seldon-core"}[{{.Range}}] {{.OffsetExp}}) / scalar(max(sum_over_time(kube_pod_labels[{{.Range}}] {{.OffsetExp}})))) * on(pod,namespace) group_right(label_seldon_app) sum by (namespace,pod,container) (rate(container_cpu_usage_seconds_total[{{.Range}}] {{.OffsetExp}}))))`


func BuildQueryFromTemplate(inputData QueryInputData, templateStr string) (query string, err error) {
	var result string

	tmpl := template.New("query-template")
	tmpl, err = tmpl.Parse(templateStr)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var resultBuffer bytes.Buffer
	if err := tmpl.Execute(&resultBuffer, inputData); err == nil {
		result = resultBuffer.String()
	} else {
		log.Println(err)
		return "", err
	}

	return result, nil
}
