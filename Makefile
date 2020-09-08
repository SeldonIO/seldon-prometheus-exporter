VERSION=0.0.1
IMG=seldonio/seldon-prometheus-exporter:${VERSION}

go-run:
	go run ./.

go-build:
	go build .

docker-build:
	docker build . -t ${IMG}

docker-push:
	docker push ${IMG}

helm-install:
	helm upgrade seldon-prometheus-exporter ./helm/seldon-prometheus-exporter --namespace=seldon-system --install --recreate-pods

helm-delete:
	helm delete seldon-prometheus-exporter --namespace=seldon-system