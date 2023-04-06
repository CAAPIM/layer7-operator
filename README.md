# Layer7 Gateway Operator
The Layer7 Gateway Operator, built using the [Operator SDK](https://github.com/operator-framework/operator-sdk) covers all aspects of deploying, maintaining and upgrading API Gateways in Kubernetes.

##### Note: The Operator currently supports ***Gateway 10.1.00_CR3 only***. All examples reference will reference that image for this release.

## About
The Operator is currently in an Alpha state and therefore does not currently carry a support/maintenance statement. Please check out the [Gateway Helm Chart](https://github.com/CAAPIM/apim-charts/tree/stable/charts/gateway) for supported container Gateway Deployment options

The initial release gives basic coverage to instantiating one or more Operator managed Gateway Deployments with or without an external MySQL Database.

The Layer7 Operator is restricted to manage the namespace it's deployed in by default. There is also a cluster-wide option available where the operator can watch all or multiple namespaces.

### Deployment Types

#### Application Level
- Database Backed Gateway Clusters
- Embedded Derby Database (Ephemeral) Gateways

#### Custom Resources
- Gateway
- Repository

#### Features
- Gateway Helm Chart feature parity (no sample mysql/hazelcast/influxdb/grafana deployments)
- Graphman Integration
- Git Integration for Graphman Repositories
  - Dynamic updates
  - Static (bootstrap) updates
- Dynamic Volumes for existing Kubernetes Configmaps/Secrets.
- Dynamic Volumes for CSI Secret Volumes.
- Application level port configuration
- Dedicated Service for access to Policy Manager/Gateway management services (when an external MySQL database is present).

#### Coming Soon
- Monitoring
  - OTel Integration

### Under consideration
- OTK support (operator managed)
- Additional Custom Resources

## Prerequisites
- Kubernetes v1.25+
- Gateway v10/11.x License
- Ingress Controller (You can also expose Gateway Services as L4 LoadBalancers)

## Index
 - [Deploy the Layer7 Operator](#installation)
   - [Using Kubernetes CLI (kubectl)](#install-with-kubectl)
     - [OwnNamespace](#ownnamespace)
     - [All/Multiple Namespaces](#allmultiple-namespaces)
   - [Using OLM (Openshift)](#install-on-openshift)
 - [Create a Simple Gateway](#create-a-simple-gateway)
 - [More Examples](./example)
 <!-- - [Monitoring](#monitoring)
   - [OTel (Open Telemetry)](#open-telemetry)
     - [Install Open Telemetry](#install-open-telemetry)
     - [OTel Collector Sidecar](#otel-collector-sidecar)
     - [Gateway Configuration](#gateway-configuration)
   - [Prometheus](#prometheus)
     - [Install Prometheus](#install-prometheus)
     - [Service Monitor](#service-monitor)
     - [Grafana Dashboard](#grafana-dashboard) -->

## Installation
There are currently two ways to deploy the Layer7 Gateway Operator. A Helm Chart will be available in the future.

### Install with kubectl
Clone this repository to get started
```
$ git clone https://github.com/caapim/layer7-operator.git
$ cd layer7-operator
```
bundle.yaml contains all of the manifests that the Layer7 Operator requires.

#### OwnNamespace
By default the Operator manages the namespace that it is deployed into and does not create any cluster roles/role bindings. 

```
$ kubectl apply -f deploy/bundle.yaml
```

#### All/Multiple Namespaces
You can also configure the Operator to watch all or multiple namespaces. This will create a namespace called <i>layer7-operator-system</i>. The default is all namespaces, you can update this by changing the following in deploy/cw-bundle.yaml

default (watches all namespaces)
```
env:
- name: WATCH_NAMESPACE
  value: ""
```
limit to specific namespaces
```
env:
- name: WATCH_NAMESPACE
  value: "ns1,ns2,ns3"
```

Once you have updated deploy/cw-bundle.yaml run the following command to install the Operator

```
$ kubectl apply -f deploy/cw-bundle.yaml
```

### Install on OpenShift
The Layer7 Operator <b>has not been published</b> to any Operator Catalogs, you can still deploy it using the operator-sdk cli. The only supported install mode in OpenShift is OwnNamespace.

```
operator-sdk run bundle docker.io/layer7api/layer7-operator-bundle:v1.0.1 --install-mode OwnNamespace
```

### Create a Simple Gateway
These steps will deploy a gateway with no additional configuration. This gateway deployment will run in ephemeral mode, without an external MySQL database. This is useful if you're new to the Layer7 API Gateway, Kubernetes Operators or you just want to make sure that the Layer7 Operator is up and running. More advanced configuration can be found [here](link)

1. Create a Kubernetes Secret with your Gateway license file

<i><b>important:</b> the file name should be license.xml</i>
```
kubectl create secret generic gateway-license --from-file=/path/to/license.xml
```
2. Create a Gateway Custom Resource
```
kubectl apply -f - <<EOF
apiVersion: security.brcmlabs.com/v1
kind: Gateway
metadata:
  name: ssg
spec:
  version: "10.1.00_CR3"
  license:
    accept: true
    secretName: gateway-license
  app:
    replicas: 1
    image: docker.io/caapim/gateway:10.1.00_CR3
    management:
      username: admin
      password: 7layer
      cluster:
        password: 7layer
        hostname: gateway.brcmlabs.com
    service:
      # annotations:
      type: LoadBalancer
      ports:
      - name: https
        port: 8443
        targetPort: 8443
        protocol: TCP
      - name: management
        port: 9443
        targetPort: 9443
        protocol: TCP
EOF
```
3. Get the LoadBalancer Address
```
$ kubectl get svc
```
expected output
```
NAME  TYPE           CLUSTER-IP    EXTERNAL-IP     PORT(S)                         AGE
ssg   LoadBalancer   10.68.4.126   35.189.116.20   8443:31384/TCP,9443:31359/TCP   7m39s
```
if your output looks like this that means you don't have an External IP Provisioner in your Kubernetes Cluster. You can still access your Gateway using port-forward.
```
NAME  TYPE           CLUSTER-IP    EXTERNAL-IP     PORT(S)                         AGE
ssg   LoadBalancer   10.68.4.126   <PENDING>   8443:31384/TCP,9443:31359/TCP   7m39s
```

4. Sign into Policy Manager
```
username: admin
password: 7layer
gateway: 35.189.116.20
```
if you don't have an external ip
```
$ kubectl port-forward svc/ssg 8443

username: admin
password: 7layer
gateway: localhost
```

5. Remove the Gateway Resource
```
$ kubectl delete gateway ssg
```

#### Remove the Operator
if you installed the operator using kubectl
```
$ kubectl delete -k deploy/bundle.yaml|cw-bundle.yaml
```

if you installed the operator in Openshift

``` 
$ operator-sdk cleanup <operatorPackageName> [flags]
```

# Watch this space for future updates.

### Monitoring


#### Open Telemetry


##### Install Open Telemetry

##### Otel Collector Sidecar

##### Gateway Configuration


#### Prometheus
##### Install Prometheus

##### Service Monitor

##### Grafana Dashboard


