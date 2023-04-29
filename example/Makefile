# This Makefile makes installing these examples more simple

.PHONY: basic advanced otel-gateway


install: 
	kubectl apply -f ../deploy/bundle.yaml

basic:
	kubectl apply -k ./repositories
	kubectl apply -k ./basic

advanced:
	kubectl apply -k ./repositories
	kubectl apply -k ./advanced

otel-example: install cert-manager prometheus open-telemetry jaeger nginx
	sleep 30
	kubectl apply -f ./otel/collector.yaml
	kubectl apply -f ./otel/observability/jaeger/jaeger.yaml
	sleep 10
	kubectl apply -k ./repositories
	kubectl apply -k ./otel-gateway

cert-manager:
	kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.11.0/cert-manager.yaml

open-telemetry:
	sleep 30
	kubectl apply -f https://github.com/open-telemetry/opentelemetry-operator/releases/latest/download/opentelemetry-operator.yaml
prometheus:
	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	helm repo update
	kubectl create ns monitoring
	kubectl apply -k ./otel/monitoring/grafana/
	helm upgrade -i prometheus -f ./otel/monitoring/prometheus/prometheus-values.yaml prometheus-community/kube-prometheus-stack -n monitoring

jaeger:
	kubectl create namespace observability
	kubectl create -f https://github.com/jaegertracing/jaeger-operator/releases/download/v1.44.0/jaeger-operator.yaml -n observability

nginx:
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml

uninstall:
	kind delete cluster --name layer7

kind-cluster:
	kind create cluster --name layer7 --config ./kind-config.yaml