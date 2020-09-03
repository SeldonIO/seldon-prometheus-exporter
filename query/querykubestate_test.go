package query

import (
	"fmt"
	"testing"
)

func TestContainerUsage(t *testing.T) {

	query,  err := BuildQueryFromTemplate(DefaultInputData, ContainersUsageSumSeldonTemplate)
	if err != nil {
		t.Errorf(err.Error())
	}
	if query != `sum by (label_seldon_app,namespace) ((sum_over_time(kube_pod_labels{label_app_kubernetes_io_managed_by=~"seldon-core"}[1d] ) / scalar(max(sum_over_time(kube_pod_labels[1d] )))) * on(pod,namespace) group_right(label_seldon_app) max by (namespace,pod,container,namespace) (avg_over_time(kube_pod_container_info[1d] )))` {
		t.Errorf("Bad seldon sum prom query: %s", query)
	}

}

func TestContainerUsageAgainstPrometheus(t *testing.T) {

	resp, err := QueryPrometheus(`sum by (label_seldon_app,namespace) ((sum_over_time(kube_pod_labels{label_app_kubernetes_io_managed_by=~"seldon-core"}[1d] ) / scalar(max(sum_over_time(kube_pod_labels[1d] )))) * on(pod,namespace) group_right(label_seldon_app) max by (namespace,pod,container,namespace) (avg_over_time(kube_pod_container_info[1d] )))`)
	if err != nil {
		t.Errorf(err.Error())
	}
	data := resp.(map[string]interface{})["data"]
	result := data.(map[string]interface{})["result"]
	firstResult := result.([]interface{})[0]
	metric := firstResult.(map[string]interface{})["metric"]
	value := firstResult.(map[string]interface{})["value"]

	model := metric.(map[string]interface{})["label_seldon_app"]
	namespace := metric.(map[string]interface{})["namespace"]
	metricVal := value.([]interface{})[0]
	fmt.Println(model)
	fmt.Println(namespace)
	fmt.Println(metricVal)
}