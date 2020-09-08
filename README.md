# Seldon Prometheus Exporter

This exporter should lookup kube-state-metrics for models and downsample relevant metrics to be exported at collected again.

## How it Works

The aim is to get SeldonDeployment metrics with a low enough res to query on. 

We do this by querying metrics from kube-state-metrics and averaging them. So we query over, say 1hr, and put in a new metric that represents the average over that hour. 

This is necessary as prometheus has a max limit on querying. See https://docs.google.com/document/d/1w8rU9gYGQ3fmm6FBI9WuXKic-wZp2Q8OUsJdnHBXwa4/edit?usp=sharing

## Metrics Format

```
model_containers_average{namespace="seldon",type="seldon",name="iris"} 1.21
model_cpu_usage_seconds_total{namespace="seldon",type="seldon",name="iris"} 0.019797927555476498
model_memory_usage_bytes{namespace="seldon",type="seldon",name="iris"} 138986.313
model_containers_average{namespace="seldon",type="seldon",name="income"} 1.15
model_cpu_usage_seconds_total{namespace="seldon",type="seldon",name="income"} 0.028586977665586328
model_memory_usage_bytes{namespace="seldon",type="seldon",name="income"} 235434.334
model_cpu_requests{name="income",namespace="seldon",type="SeldonDeployment"} 1.1
model_cpu_limits{name="income",namespace="seldon",type="SeldonDeployment"} 1
model_memory_requests_bytes{name="income-default",namespace="seldon",type="SeldonDeployment"} 1.073741824e+09
model_memory_limits_bytes{name="income-default",namespace="seldon",type="SeldonDeployment"} 1.073741824e+09
```

These are based on [kube-state-metrics](https://github.com/kubernetes/kube-state-metrics/blob/e43aaa6d6e3554d050ead73b4814566b771377d1/docs/pod-metrics.md) that provide the data as a basis.

## Configuration

The query frequency needs to be configured. This needs to be configured both on the internal querying frequency and with the same frequency on scraping of the exporter.

The exporter is both scraped by prometheus and reads from prometheus. It needs configuration for both.

Data refresh frequency is set with `TIME_PERIOD`. Format is a prom time period.

The prometheus to gather data from is set with `PROMETHEUS_URL`.

If a token is needed then set in `PROMETHEUS_SELDON_TOKEN`.

These are also exposed in the helm chart values file.

The chart exposes a metrics port and adds annotations for scraping by seldon core analytics.

# How to Run

First port-forward to a prometheus in a cluster running Seldon. 
```
kubectl port-forward -n seldon-system svc/seldon-core-analytics-prometheus-seldon 8080:80
```
Then `go run ./.`

Go to `http://localhost:8000/metrics` to see metrics.

Or use helm chart with `make helm-install`

## Notes

Project initially based on https://rsmitty.github.io/Prometheus-Exporters-Revamp/
