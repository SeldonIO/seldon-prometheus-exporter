package query

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"text/template"
)

type QueryInputData struct {
	Range string `json:",omitempty"`
	OffsetExp string `json:",omitempty"`
}

type MetricInstance struct {
	Model string
	Namespace string
	ModelType string
	Value float64
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

func ObtainMetricValues(queryTemplate string, inputData QueryInputData) ([]MetricInstance, error){
	query, err := BuildQueryFromTemplate(inputData, queryTemplate)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	resp, err := QueryPrometheus(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	fmt.Println(query)

	data := resp.(map[string]interface{})["data"]
	result := data.(map[string]interface{})["result"]

	var metricInstances []MetricInstance

	for _, res := range result.([]interface{}) {
		metric := res.(map[string]interface{})["metric"]
		value := res.(map[string]interface{})["value"]

		model := metric.(map[string]interface{})["label_seldon_app"]
		namespace := metric.(map[string]interface{})["namespace"]
		metricVal := value.([]interface{})[1]

		f, err := strconv.ParseFloat(metricVal.(string),64)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		//TODO: "SeldonDeployment" is hardcoded
		metricInstances = append(metricInstances, MetricInstance{Model: model.(string),Namespace: namespace.(string),ModelType: "SeldonDeployment",Value: f})
		fmt.Println(model)
		fmt.Println(namespace)
		fmt.Println(metricVal)
	}
	return metricInstances, nil
}

func QueryPrometheus(query string) (interface{}, error){
	params := url.Values{}
	params.Add("query", query)

	promUrl := os.Getenv("PROMETHEUS_URL")
	if promUrl == "" {
		promUrl = "http://localhost:8080/api/v1/"
	}
	queryType := "query" //not doing range
	promUrl = promUrl + queryType
	baseURL, err := url.Parse(promUrl)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	baseURL.RawQuery = params.Encode()
	queryObj := strings.NewReader(params.Encode())

	req, err := http.NewRequest("GET", baseURL.String(), queryObj)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	promToken := os.Getenv("PROMETHEUS_SELDON_TOKEN")

	if promToken != "" {
		req.Header.Add("Authorization", "Bearer "+promToken)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: tr,
	}

	resp, err := client.Do(req)
	//log.Println("QUERY AGAINST " + baseURL.String() + " of " + inputData.QueryTemplate)
	//log.Println(query)
	//log.Println(params.Encode())
	//log.Println(resp)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	var queryData interface{}
	json.NewDecoder(resp.Body).Decode(&queryData)
	return queryData, nil
}