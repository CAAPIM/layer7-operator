# This Makefile drives the quickstart sections of all examples.

.PHONY: basic advanced otel-gateway

install:
	kubectl apply -f ../deploy/bundle.yaml

basic:
	kubectl apply -k ./repositories
	sleep 10
	kubectl apply -k ./basic

advanced:
	kubectl apply -k ./repositories
	sleep 10
	kubectl apply -k ./advanced

otel-example-kind: install cert-manager prometheus open-telemetry jaeger nginx-kind
	kubectl apply -f ./otel/collector.yaml
	kubectl apply -f ./otel/observability/jaeger/jaeger.yaml
	kubectl apply -k ./repositories
	sleep 10
	kubectl apply -k ./otel-gateway

otel-example: install cert-manager prometheus open-telemetry jaeger
	sleep 10
	kubectl apply -k ./repositories
	kubectl apply -k ./otel-gateway

cert-manager:
	kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.11.0/cert-manager.yaml
	sleep 5
	kubectl wait --for=condition=ready --timeout=60s pod -l app=cert-manager -n cert-manager
	kubectl wait --for=condition=ready --timeout=60s pod -l app=cainjector -n cert-manager
	kubectl wait --for=condition=ready --timeout=60s pod -l app=webhook -n cert-manager
	
open-telemetry:
	kubectl apply -f https://github.com/open-telemetry/opentelemetry-operator/releases/latest/download/opentelemetry-operator.yaml
	sleep 5
	kubectl wait --for=condition=ready --timeout=60s pod -l app.kubernetes.io/name=opentelemetry-operator -n opentelemetry-operator-system
	kubectl apply -f ./otel/collector.yaml

prometheus:
	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	helm repo update
	-kubectl create ns monitoring
	kubectl apply -k ./otel/monitoring/grafana/
	helm upgrade -i prometheus -f ./otel/monitoring/prometheus/prometheus-values.yaml prometheus-community/kube-prometheus-stack -n monitoring

jaeger:
	-kubectl create namespace observability
	kubectl apply -f https://github.com/jaegertracing/jaeger-operator/releases/download/v1.44.0/jaeger-operator.yaml -n observability
	sleep 5
	kubectl wait --for=condition=ready --timeout=60s pod -l name=jaeger-operator -n observability
	kubectl apply -f ./otel/observability/jaeger/jaeger.yaml
	kubectl apply -f ./otel/observability/jaeger/ingress.yaml

nginx-kind:
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
	kubectl wait --for=condition=ready --timeout=60s pod -l app.kubernetes.io/name=ingress-nginx,app.kubernetes.io/component=controller -n ingress-nginx

nginx:
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/cloud/deploy.yaml
	kubectl wait --for=condition=ready --timeout=60s pod -l app.kubernetes.io/name=ingress-nginx,app.kubernetes.io/component=controller -n ingress-nginx

uninstall:
	-kubectl delete -k ./otel-gateway
	-kubectl delete -k ./repositories/
	-kubectl delete -f ./otel/collector.yaml
	-kubectl delete -f ./otel/observability/jaeger/jaeger.yaml
	-kubectl delete -f ./otel/observability/jaeger/ingress.yaml
	-kubectl delete -f ../deploy/bundle.yaml
	-kubectl delete -f https://github.com/open-telemetry/opentelemetry-operator/releases/latest/download/opentelemetry-operator.yaml
	-helm uninstall prometheus -n monitoring
	-kubectl delete -k ./otel/monitoring/grafana/
	-kubectl delete -f https://github.com/cert-manager/cert-manager/releases/download/v1.11.0/cert-manager.yaml
	-kubectl delete -f https://github.com/jaegertracing/jaeger-operator/releases/download/v1.44.0/jaeger-operator.yaml -n observability
	-kubectl delete -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
	-kubectl delete -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/cloud/deploy.yaml
	-kubectl delete ns observability
	-kubectl delete ns monitoring

uninstall-kind:
	kind delete cluster --name layer7

kind-cluster:
	kind create cluster --name layer7 --config ./kind-config.yaml