# This Makefile drives the quickstart sections of all examples.
#SHELL := /bin/bash

NAMESPACE ?= default

.PHONY: basic advanced otel-gateway wait elastic metrics-server prepare-collector test grafana-stack prometheus otel otk

wait:
ifeq ($(OS),Windows_NT)
	TIMEOUT /T $(t)
else
	sleep $(t)
endif

tmp:=$(shell mktemp -d)

create-collector:
	echo ${t}
	cp ./otel-elastic/collector.yaml $(tmp)/collector.yaml
	sed -i -e 's/APM_AUTH_TOKEN/'${t}'/g' $(tmp)/collector.yaml
	cat $(tmp)/collector.yaml
	kubectl apply -f $(tmp)/collector.yaml
	rm -rf $(tmp)

install:
	kubectl apply -f https://github.com/CAAPIM/layer7-operator/releases/download/v1.2.2/bundle.yaml
	@$(MAKE) --silent t=10 wait
	kubectl wait --for=condition=ready --timeout=600s pod -l app.kubernetes.io/name=layer7-operator

basic:
	kubectl apply -k ./repositories
	@$(MAKE) --silent t=10 wait 
	kubectl apply -f ./gateway/basic-gateway.yaml

advanced:
	kubectl apply -k ./repositories
	@$(MAKE) --silent t=10 wait
	kubectl apply -f ./gateway/advanced-gateway.yaml

otk-single:
	kubectl apply -k ./repositories
	@$(MAKE) --silent t=10 wait
	kubectl apply -f ./gateway/otk/otk-single.yaml

portal-example: redis
	kubectl apply -f https://github.com/CAAPIM/layer7-operator/releases/download/v1.2.2/bundle.yaml -n ${NAMESPACE}
	@$(MAKE) --silent t=15 wait
	kubectl -n ${NAMESPACE} wait --for=condition=ready --timeout=600s pod -l app.kubernetes.io/name=layer7-operator
	kubectl apply -k ./base -n ${NAMESPACE}
	kubectl apply -f ./gateway/portal-gateway.yaml -n ${NAMESPACE}
	helm repo add layer7 https://caapim.github.io/apim-charts/
	helm repo update
	@echo "####################################################################################################################################"
	@echo "Deploying the Layer7 API Developer Portal, this will take several minutes and complete when all components are ready."
	@echo "####################################################################################################################################"
	helm upgrade --install portal --set ingress.create=false --set-file portal.registryCredentials=./portal-integration/secrets/docker-secret.yaml -f ./portal-integration/portal-values.yaml layer7/portal -n ${NAMESPACE}
	kubectl -n ${NAMESPACE} wait --for=condition=ready --timeout=600s pod -l app=apim
	@$(MAKE) --silent t=15 wait
	@echo "####################################################################################################################################"
	@echo "Provisioning Portal Tenant"
	@echo "####################################################################################################################################"
	./portal-integration/create-tenant.sh -d ./portal-integration/enroll-payload.json -n ${NAMESPACE}

g2c-agent:
#sed -e 's/ip:.*/ip: "$(shell kubectl get svc ingress-nginx-controller -n ingress-nginx -ojsonpath="{.spec.clusterIP}")"/g' ./portal-integration/g2c-agent/deployment-template.yaml > ./portal-integration/g2c-agent/deployment.yaml
	kubectl apply -k ./portal-integration/g2c-agent
	echo "deploy g2cagent"

# add function to generate agent config
redis:
	helm upgrade -i standalone -f ./portal-integration/redis/redis-values.yaml oci://registry-1.docker.io/bitnamicharts/redis -n ${NAMESPACE}
	@$(MAKE) --silent t=15 wait
	kubectl wait --for=condition=ready --timeout=600s pod -l app.kubernetes.io/name=redis -n ${NAMESPACE}

default-redis:
	helm upgrade -i standalone -f ./portal-integration/redis/redis-values.yaml oci://registry-1.docker.io/bitnamicharts/redis
	@$(MAKE) --silent t=15 wait
	kubectl wait --for=condition=ready --timeout=600s pod -l app.kubernetes.io/name=redis

grafana-stack:
	helm repo add grafana https://grafana-community.github.io/helm-charts
	helm upgrade --install --values ./otel-lgtm/grafana-stack/loki-overrides.yaml loki grafana/loki -n grafana-loki --create-namespace
	helm upgrade --install --values ./otel-lgtm/grafana-stack/promtail-overrides.yaml promtail grafana/promtail -n grafana-loki
	helm upgrade --install --values ./otel-lgtm/grafana-stack/tempo-overrides.yaml tempo grafana/tempo -n grafana-loki
	helm upgrade --install --values ./otel-lgtm/grafana-stack/mimir-distributed-overrides.yaml mimir grafana/mimir-distributed -n grafana-loki

otel-lgtm-example-kind: install cert-manager prometheus-lgtm open-telemetry grafana-stack contour-kind
	kubectl apply -f ./otel-lgtm/collector.yaml
	kubectl apply -f ./otel-lgtm/instrumentation.yaml
	kubectl apply -k ./repositories
	@$(MAKE) --silent t=10 wait
	kubectl apply -f ./gateway/otel-lgtm-gateway.yaml


otel-lgtm-example: install cert-manager prometheus-lgtm open-telemetry grafana-stack
	kubectl apply -f ./otel-lgtm/collector.yaml
	kubectl apply -f ./otel-lgtm/instrumentation.yaml
	kubectl apply -k ./repositories
	@$(MAKE) --silent t=10 wait
	kubectl apply -f ./gateway/otel-lgtm-gateway.yaml

# otel-prometheus-example-kind: install cert-manager prometheus open-telemetry jaeger contour-kind
# 	kubectl apply -f ./otel-prometheus/servicemonitor.yaml
# 	kubectl apply -f ./otel-prometheus/collector.yaml
# 	kubectl apply -f ./otel-prometheus/instrumentation.yaml
# 	kubectl apply -f ./otel-prometheus/observability/jaeger/jaeger.yaml
# 	kubectl apply -f ./otel-prometheus/observability/jaeger/ingress.yaml
# 	kubectl apply -k ./repositories
# 	@$(MAKE) --silent t=10 wait
# 	kubectl apply -f ./gateway/otel-prometheus-gateway.yaml

# otel-prometheus-example: install cert-manager prometheus open-telemetry jaeger
# 	kubectl apply -f ./otel-prometheus/servicemonitor.yaml
# 	kubectl apply -f ./otel-prometheus/collector.yaml
# 	kubectl apply -f ./otel-prometheus/instrumentation.yaml
# 	kubectl apply -f ./otel-prometheus/observability/jaeger/jaeger.yaml
# 	kubectl apply -f ./otel-prometheus/observability/jaeger/ingress.yaml
# 	kubectl apply -k ./repositories
# 	@$(MAKE) --silent t=10 wait
# 	kubectl apply -f ./gateway/otel-prometheus-gateway.yaml

otel-elastic-example-kind: install cert-manager open-telemetry elastic contour-kind
	@$(MAKE) --silent t=$(shell kubectl get secret/apm-server-quickstart-apm-token -o go-template='{{index .data "secret-token" | base64decode}}') create-collector
	kubectl apply -f ./otel-elastic/instrumentation.yaml
	@$(MAKE) --silent t=10 wait
	kubectl apply -k ./repositories
	kubectl apply -f ./gateway/otel-elastic-gateway.yaml
	@echo "Creating Layer7 Dashboard"
	elasticPass=$(shell kubectl get secret quickstart-es-elastic-user -o go-template='{{.data.elastic | base64decode}}')
	@curl -XPOST -s --write-out "%{http_code}" --output /dev/null https://kibana.brcmlabs.com/api/saved_objects/_import?createNewCopies=false -H "kbn-xsrf: true" -k -uelastic:$(shell kubectl get secret quickstart-es-elastic-user -o go-template='{{.data.elastic | base64decode}}') -F "file=@./otel-elastic/dashboard/apim-dashboard.ndjson"
	@echo "\n#####################################################\nComplete\n#####################################################"
	@echo "1. Open kibana.brcmlabs.com in a browser\n2. Complete the Post Installation tasks\n3. Test your Gateway\n4. View Logs, Metrics and Traces\n\nusername: elastic\npassword: $(shell kubectl get secret quickstart-es-elastic-user -o go-template='{{.data.elastic | base64decode}}')"
	@echo "#####################################################\n"

otel-elastic-example: install cert-manager open-telemetry elastic
	@$(MAKE) --silent t=$(shell kubectl get secret/apm-server-quickstart-apm-token -o go-template='{{index .data "secret-token" | base64decode}}') create-collector
	kubectl apply -f ./otel-elastic/instrumentation.yaml
	@$(MAKE) --silent t=10 wait
	kubectl apply -k ./repositories
	kubectl apply -f ./gateway/otel-elastic-gateway.yaml
	@echo "Creating Layer7 Dashboard"
	elasticPass=$(shell kubectl get secret quickstart-es-elastic-user -o go-template='{{.data.elastic | base64decode}}')
	@curl -XPOST -s --write-out "%{http_code}" --output /dev/null https://kibana.brcmlabs.com/api/saved_objects/_import?createNewCopies=false -H "kbn-xsrf: true" -k -uelastic:$(shell kubectl get secret quickstart-es-elastic-user -o go-template='{{.data.elastic | base64decode}}') -F "file=@./otel-elastic/dashboard/apim-dashboard.ndjson"
	@echo "\n#####################################################\nComplete\n#####################################################"
	@echo "1. Open kibana.brcmlabs.com in a browser\n2. Complete the Post Installation tasks\n3. Test your Gateway\n4. View Logs, Metrics and Traces\n\nusername: elastic\npassword: $(shell kubectl get secret quickstart-es-elastic-user -o go-template='{{.data.elastic | base64decode}}')"
	@echo "#####################################################\n"

cert-manager:
	kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.20.0/cert-manager.yaml
	@$(MAKE) --silent t=10 wait
	kubectl wait --for=condition=ready --timeout=600s pod -l app=cert-manager -n cert-manager
	kubectl wait --for=condition=ready --timeout=600s pod -l app=cainjector -n cert-manager
	kubectl wait --for=condition=ready --timeout=600s pod -l app=webhook -n cert-manager

elastic:
	kubectl create -f https://download.elastic.co/downloads/eck/3.3.1/crds.yaml
	kubectl apply -f https://download.elastic.co/downloads/eck/3.3.1/operator.yaml
	@$(MAKE) --silent t=10 wait
	kubectl wait --for=condition=ready --timeout=600s pod -l control-plane=elastic-operator -n elastic-system
	kubectl apply -f ./otel-elastic/components
	@$(MAKE) --silent t=10 wait
	kubectl wait --for=condition=ready --timeout=600s pod -l elasticsearch.k8s.elastic.co/statefulset-name=quickstart-es-default
	kubectl wait --for=condition=ready --timeout=600s pod -l apm.k8s.elastic.co/name=apm-server-quickstart
	kubectl wait --for=condition=ready --timeout=600s pod -l kibana.k8s.elastic.co/name=quickstart
	kubectl wait --for=condition=ready --timeout=600s pod -l agent.k8s.elastic.co/name=elastic-agent
	kubectl wait --for=condition=ready --timeout=600s pod -l beat.k8s.elastic.co/name=filebeat
	kubectl wait --for=condition=ready --timeout=600s pod -l beat.k8s.elastic.co/name=metricbeat

metrics-server:
	kubectl apply -f ./metrics-server/metrics-server-0-6-3.yaml

open-telemetry:
	kubectl apply -f https://github.com/open-telemetry/opentelemetry-operator/releases/download/v0.146.0/opentelemetry-operator.yaml
	@$(MAKE) --silent t=10 wait
	kubectl wait --for=condition=ready --timeout=600s pod -l app.kubernetes.io/name=opentelemetry-operator -n opentelemetry-operator-system

# prometheus:
# 	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
# 	helm repo update
# 	-kubectl create ns monitoring
# 	kubectl apply -k ./otel-prometheus/monitoring/grafana/
# 	helm upgrade -i prometheus -f ./otel-prometheus/monitoring/prometheus/prometheus-values.yaml prometheus-community/kube-prometheus-stack -n monitoring

prometheus-lgtm:
	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	helm repo update
	-kubectl create ns monitoring
	kubectl apply -k ./otel-lgtm/prometheus/grafana-dashboard/
	@$(MAKE) NAMESPACE=monitoring pki
	helm upgrade -i prometheus -f ./otel-lgtm/prometheus/prometheus-values.yaml prometheus-community/kube-prometheus-stack -n monitoring

jaeger:
	-kubectl create namespace observability
	kubectl apply -f https://github.com/jaegertracing/jaeger-operator/releases/download/v1.65.0/jaeger-operator.yaml -n observability
	@$(MAKE) --silent t=10 wait
	kubectl wait --for=condition=ready --timeout=600s pod -l name=jaeger-operator -n observability

contour: pki
	helm repo add contour https://projectcontour.github.io/helm-charts/
	helm upgrade --install contour contour/contour --namespace projectcontour --create-namespace
	@$(MAKE) --silent t=10 wait
	kubectl wait --for=condition=ready --timeout=600s pod -l app.kubernetes.io/name=contour,app.kubernetes.io/component=contour -n projectcontour
	kubectl wait --for=condition=ready --timeout=600s pod -l app.kubernetes.io/name=contour,app.kubernetes.io/component=envoy -n projectcontour

contour-kind: pki
	helm repo add contour https://projectcontour.github.io/helm-charts/
	helm upgrade --install contour contour/contour --namespace projectcontour --set envoy.useHostPort.http=true --set envoy.useHostPort.https=true  --create-namespace
	@$(MAKE) --silent t=10 wait
	kubectl wait --for=condition=ready --timeout=600s pod -l app.kubernetes.io/name=contour,app.kubernetes.io/component=contour -n projectcontour
	kubectl wait --for=condition=ready --timeout=600s pod -l app.kubernetes.io/name=contour,app.kubernetes.io/component=envoy -n projectcontour

uninstall-contour:
	helm uninstall contour -n projectcontour

pki:
	-openssl req -x509 -newkey rsa:2048 -nodes -days 365 -subj "/CN=*.brcmlabs.com" \
	  -addext "subjectAltName=DNS:*.brcmlabs.com,DNS:brcmlabs.com" \
	  -keyout /tmp/tls.key -out /tmp/tls.crt
	-kubectl create secret tls brcmlabs \
	  --cert=/tmp/tls.crt --key=/tmp/tls.key \
	  -n $(NAMESPACE) --dry-run=client -o yaml | kubectl apply -f -

pki-otk:
	@test -f /tmp/otk-internal-tls.crt || \
	  openssl req -x509 -newkey rsa:2048 -nodes -days 365 \
	  -subj "/CN=gateway-internal-otk.brcmlabs.com" \
	  -addext "subjectAltName=DNS:gateway-internal-otk.brcmlabs.com,DNS:gateway-internal-otk-pm.brcmlabs.com" \
	  -keyout /tmp/otk-internal-tls.key -out /tmp/otk-internal-tls.crt
	@test -f /tmp/otk-dmz-tls.crt || \
	  openssl req -x509 -newkey rsa:2048 -nodes -days 365 \
	  -subj "/CN=gateway-dmz-otk.brcmlabs.com" \
	  -addext "subjectAltName=DNS:gateway-dmz-otk.brcmlabs.com,DNS:gateway-dmz-otk-pm.brcmlabs.com" \
	  -keyout /tmp/otk-dmz-tls.key -out /tmp/otk-dmz-tls.crt
	@test -f /tmp/otk-dmz-mtls.crt || { \
	  openssl req -x509 -newkey rsa:2048 -nodes -days 365 \
	  -subj "/CN=otk-dmz-mtls-ca" \
	  -keyout /tmp/otk-dmz-mtls-ca.key -out /tmp/otk-dmz-mtls-ca.crt && \
	  openssl req -newkey rsa:2048 -nodes \
	  -subj "/CN=dmz-gateway" \
	  -keyout /tmp/otk-dmz-mtls.key -out /tmp/otk-dmz-mtls.csr && \
	  openssl x509 -req -in /tmp/otk-dmz-mtls.csr \
	  -CA /tmp/otk-dmz-mtls-ca.crt -CAkey /tmp/otk-dmz-mtls-ca.key \
	  -CAcreateserial -days 365 -out /tmp/otk-dmz-mtls.crt; \
	}

pki-otk-dmz: pki-otk
	-kubectl create secret tls otk-dmz-secret \
	  --cert=/tmp/otk-dmz-tls.crt --key=/tmp/otk-dmz-tls.key \
	  -n $(NAMESPACE) --dry-run=client -o yaml | kubectl apply -f -
	-kubectl create secret tls otk-dmz-mtls-secret \
	  --cert=/tmp/otk-dmz-mtls.crt --key=/tmp/otk-dmz-mtls.key \
	  -n $(NAMESPACE) --dry-run=client -o yaml | kubectl apply -f -
	-kubectl create secret generic otk-internal-ca-cert \
	  --from-file=ca.crt=/tmp/otk-internal-tls.crt \
	  -n $(NAMESPACE) --dry-run=client -o yaml | kubectl apply -f -

pki-otk-internal: pki-otk
	-kubectl create secret tls otk-internal-secret \
	  --cert=/tmp/otk-internal-tls.crt --key=/tmp/otk-internal-tls.key \
	  -n $(NAMESPACE) --dry-run=client -o yaml | kubectl apply -f -
	-kubectl create secret generic otk-dmz-mtls-public-cert \
	  --from-file=dmz-gateway.crt=/tmp/otk-dmz-mtls.crt \
	  -n $(NAMESPACE) --dry-run=client -o yaml | kubectl apply -f -

cassandra:
	helm upgrade --install cassandra -f ./otk/database/cassandra/cassandra-values.yaml oci://registry-1.docker.io/bitnamicharts/cassandra

mysql:
	helm upgrade --install mysql -f ./otk/database/mysql/mysql-values.yaml oci://registry-1.docker.io/bitnamicharts/mysql
	kubectl wait --for=condition=ready --timeout=600s pod -l app.kubernetes.io/name=mysql

enable-otel-l7operator:
	kubectl apply -f ./otel-lgtm/layer7-operator/collector.yaml
	kubectl patch deployment layer7-operator-controller-manager -p '{"spec":{"template":{"spec":{"containers":[{"env":[{"name":"ENABLE_OTEL","value":"true"},{"name":"OTEL_EXPORTER_OTLP_ENDPOINT","value":"layer7-operator-collector:4317"},{"name": "OTEL_METRIC_PREFIX","value":"layer7_"}], "name": "manager"}]}}}}'

uninstall:
	-kubectl delete -f ./gateway
	-kubectl delete -k ./repositories/
	-kubectl delete -f ./otel-prometheus/collector.yaml
	-kubectl delete -f ./otel-elastic/collector.yaml
	-kubectl delete -f ./otel-prometheus/instrumentation.yaml
	-kubectl delete -f ./otel-elastic/instrumentation.yaml
	-kubectl delete -f ./otel-prometheus/observability/jaeger
	-kubectl delete -f https://github.com/open-telemetry/opentelemetry-operator/releases/download/v0.146.0/opentelemetry-operator.yaml
	-helm uninstall prometheus -n monitoring
	-kubectl delete -k ./otel-prometheus/monitoring/grafana/
	-kubectl delete -f https://github.com/cert-manager/cert-manager/releases/download/v1.20.0/cert-manager.yaml
	-kubectl delete -f https://github.com/jaegertracing/jaeger-operator/releases/download/v1.65.0/jaeger-operator.yaml -n observability
	-kubectl delete ns observability
	-kubectl delete ns monitoring
	-kubectl delete -f https://github.com/CAAPIM/layer7-operator/releases/download/v1.2.2/bundle.yaml
	-helm uninstall mysql
	-helm uninstall cassandra
	-helm uninstall contour -n projectcontour

uninstall-kind:
	kind delete cluster --name layer7

kind-cluster:
	kind create cluster --name layer7 --config ./kind-config.yaml