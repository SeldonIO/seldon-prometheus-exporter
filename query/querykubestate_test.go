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
	if query != `sum by (label_seldon_deployment_id,namespace) ((sum_over_time(kube_pod_labels{label_app_kubernetes_io_managed_by=~"seldon-core"}[1d] ) / scalar(max(sum_over_time(kube_pod_labels[1d] )))) * on(pod,namespace) group_right(label_seldon_deployment_id) max by (namespace,pod,container,namespace) (avg_over_time(kube_pod_container_info[1d] )))` {
		t.Errorf("Bad seldon sum prom query: %s", query)
	}

}

func TestContainerUsageAgainstPrometheus(t *testing.T) {

	resp, err := QueryPrometheus(`sum by (label_seldon_deployment_id,namespace) ((sum_over_time(kube_pod_labels{label_app_kubernetes_io_managed_by=~"seldon-core"}[1d] ) / scalar(max(sum_over_time(kube_pod_labels[1d] )))) * on(pod,namespace) group_right(label_seldon_deployment_id) max by (namespace,pod,container,namespace) (avg_over_time(kube_pod_container_info[1d] )))`)
	if err != nil {
		t.Errorf(err.Error())
	}
	data := resp.(map[string]interface{})["data"]
	result := data.(map[string]interface{})["result"]
	for _, res := range result.([]interface{}) {
		metric := res.(map[string]interface{})["metric"]
		value := res.(map[string]interface{})["value"]

		model := metric.(map[string]interface{})["label_seldon_deployment_id"]
		namespace := metric.(map[string]interface{})["namespace"]
		metricVal := value.([]interface{})[0]
		fmt.Println(model)
		fmt.Println(namespace)
		fmt.Println(metricVal)
	}

	//TODO: need to register metrics and create entries for different labels
	// above is a metric for containers and each iteration of loop is a different set of labels
	// see https://github.com/SeldonIO/seldon-deploy/issues/1294#issuecomment-675010586
	// example from https://github.com/prometheus/client_golang/issues/364 may be good reference
	// or https://github.com/SeldonIO/seldon-core/blob/3b92c880cd115af386e3cb5e1152adec4c1f65cb/executor/api/metric/client.go#L15
	// should we respond to scrapes... which could be problematic as the prom client library provides the handler?
	// or just have a cycle process on the same frequence as the scrape interval that updates the metrics
	// could alternatively implement endpont that returns data in prom format without using prom client.. but that feels awkward
}