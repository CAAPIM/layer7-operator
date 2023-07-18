# This Makefile drives the quickstart sections of all examples.
#SHELL := /bin/bash

.PHONY: basic advanced otel-gateway wait elastic metrics-server prepare-collector test

wait:
ifeq ($(OS),Windows_NT)
	TIMEOUT /T $(t)
else
	sleep $(t)
endif

elasticPass=$(shell kubectl get secret quickstart-es-elastic-user -o go-template='{{.data.elastic | base64decode}}')
tmp:=$(shell mktemp -d)

create-collector:
	echo ${t}
	cp ./otel-eck/collector.yaml $(tmp)/collector.yaml
	sed -i -e 's/APM_AUTH_TOKEN/'${t}'/g' $(tmp)/collector.yaml
	cat $(tmp)/collector.yaml
	kubectl apply -f $(tmp)/collector.yaml
	rm -rf $(tmp)

install:
	kubectl apply -f ../deploy/bundle.yaml

basic:
	kubectl apply -k ./repositories
	@$(MAKE) --silent t=10 wait 
	kubectl apply -k ./basic

advanced:
	kubectl apply -k ./repositories
	@$(MAKE) --silent t=10 wait
	kubectl apply -k ./advanced

otel-example-kind: install cert-manager prometheus open-telemetry jaeger nginx-kind
	kubectl apply -f ./otel/collector.yaml
	kubectl apply -f ./otel/observability/jaeger/jaeger.yaml
	kubectl apply -f ./otel/observability/jaeger/ingress.yaml
	kubectl apply -k ./repositories
	@$(MAKE) --silent t=10 wait
	kubectl apply -k ./otel-gateway

otel-example: install cert-manager prometheus open-telemetry jaeger
	kubectl apply -f ./otel/collector.yaml
	kubectl apply -f ./otel/observability/jaeger/jaeger.yaml
	kubectl apply -f ./otel/observability/jaeger/ingress.yaml
	kubectl apply -k ./repositories
	@$(MAKE) --silent t=10 wait
	kubectl apply -k ./otel-gateway

elastic-example-kind: install cert-manager open-telemetry elastic nginx-kind
	@$(MAKE) --silent t=$(shell kubectl get secret/apm-server-quickstart-apm-token -o go-template='{{index .data "secret-token" | base64decode}}') create-collector
	kubectl apply -f ./otel-eck/instrumentation.yaml
	kubectl apply -f ./otel-eck/manifests/ingress.yaml
	@$(MAKE) --silent t=10 wait
	kubectl apply -k ./repositories
	kubectl apply -k ./otel-eck-gateway
	@echo "Creating Layer7 Dashboard"
	@curl -XPOST -s --write-out "%{http_code}" --output /dev/null https://kibana.brcmlabs.com/api/saved_objects/_import?createNewCopies=true -H "kbn-xsrf: true" -k -uelastic:${elasticPass} -F "file=@./otel-eck/dashboard/apim-dashboard.ndjson"
	@echo "\n#####################################################\nComplete\n#####################################################"
	@echo "1. Update your hosts file\n2. Open kibana.brcmlabs.com in a browser\n3. Complete the Post Installation tasks\n4. Test your Gateway\n5. View Logs, Metrics and Traces\n\nusername: elastic\npassword: ${elasticPass}"
	@echo "#####################################################\n"

elastic-example: install cert-manager open-telemetry elastic
	@$(MAKE) --silent t=$(shell kubectl get secret/apm-server-quickstart-apm-token -o go-template='{{index .data "secret-token" | base64decode}}') create-collector
	kubectl apply -f ./otel-eck/instrumentation.yaml
	kubectl apply -f ./otel-eck/manifests/ingress.yaml
	@$(MAKE) --silent t=10 wait
	kubectl apply -k ./repositories
	kubectl apply -k ./otel-eck-gateway
	@echo "Creating Layer7 Dashboard"
	@curl -XPOST -s --write-out "%{http_code}" --output /dev/null https://kibana.brcmlabs.com/api/saved_objects/_import?createNewCopies=true -H "kbn-xsrf: true" -k -uelastic:${elasticPass} -F "file=@./otel-eck/dashboard/apim-dashboard.ndjson"
	@echo "\n#####################################################\nComplete\n#####################################################"
	@echo "1. Update your hosts file\n2. Open kibana.brcmlabs.com in a browser\n3. Complete the Post Installation tasks\n4. Test your Gateway\n5. View Logs, Metrics and Traces\n\nusername: elastic\npassword: ${elasticPass}"
	@echo "#####################################################\n"

cert-manager:
	kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.11.0/cert-manager.yaml
	@$(MAKE) --silent t=10 wait
	kubectl wait --for=condition=ready --timeout=60s pod -l app=cert-manager -n cert-manager
	kubectl wait --for=condition=ready --timeout=60s pod -l app=cainjector -n cert-manager
	kubectl wait --for=condition=ready --timeout=60s pod -l app=webhook -n cert-manager

elastic:
	kubectl create -f https://download.elastic.co/downloads/eck/2.8.0/crds.yaml
	kubectl apply -f https://download.elastic.co/downloads/eck/2.8.0/operator.yaml
	@$(MAKE) --silent t=10 wait
	kubectl wait --for=condition=ready --timeout=60s pod -l control-plane=elastic-operator -n elastic-system
	kubectl apply -k ./elastic
	@$(MAKE) --silent t=10 wait
	kubectl wait --for=condition=ready --timeout=120s pod -l elasticsearch.k8s.elastic.co/statefulset-name=quickstart-es-default
	kubectl wait --for=condition=ready --timeout=120s pod -l apm.k8s.elastic.co/name=apm-server-quickstart
	kubectl wait --for=condition=ready --timeout=120s pod -l kibana.k8s.elastic.co/name=quickstart
	kubectl wait --for=condition=ready --timeout=120s pod -l agent.k8s.elastic.co/name=elastic-agent
	kubectl wait --for=condition=ready --timeout=120s pod -l beat.k8s.elastic.co/name=filebeat
	kubectl wait --for=condition=ready --timeout=120s pod -l beat.k8s.elastic.co/name=metricbeat

metrics-server:
	kubectl apply -f ./metrics-server/metrics-server-0-6-3.yaml

open-telemetry:
	kubectl apply -f https://github.com/open-telemetry/opentelemetry-operator/releases/download/v0.76.1/opentelemetry-operator.yaml
	@$(MAKE) --silent t=10 wait
	kubectl wait --for=condition=ready --timeout=60s pod -l app.kubernetes.io/name=opentelemetry-operator -n opentelemetry-operator-system

prometheus:
	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	helm repo update
	-kubectl create ns monitoring
	kubectl apply -k ./otel/monitoring/grafana/
	helm upgrade -i prometheus -f ./otel/monitoring/prometheus/prometheus-values.yaml prometheus-community/kube-prometheus-stack -n monitoring

jaeger:
	-kubectl create namespace observability
	kubectl apply -f https://github.com/jaegertracing/jaeger-operator/releases/download/v1.44.0/jaeger-operator.yaml -n observability
	@$(MAKE) --silent t=10 wait
	kubectl wait --for=condition=ready --timeout=60s pod -l name=jaeger-operator -n observability

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