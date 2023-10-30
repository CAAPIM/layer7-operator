# Layer7 Gateway Operator
The Layer7 Gateway Operator, built using the [Operator SDK](https://github.com/operator-framework/operator-sdk) covers all aspects of deploying, maintaining and upgrading API Gateways in Kubernetes.

##### Note: The Operator examples currently use ***Gateway 11.0.00_CR1*** as a base. The OTel examples use a custom 11.0.00_CR1 image.

## About
The Operator is currently in an Alpha state and therefore does not currently carry a support/maintenance statement. Please check out the [Gateway Helm Chart](https://github.com/CAAPIM/apim-charts/tree/stable/charts/gateway) for supported container Gateway Deployment options

The initial release gives basic coverage to instantiating one or more Operator managed Gateway Deployments with or without an external MySQL Database.

The Layer7 Operator is restricted to manage the namespace it's deployed in by default. There is also a cluster-wide option available where the operator can watch all or multiple namespaces.


## Documentation
- [API (Custom Resource) Docs](./docs/readme.md)
- [Wiki](https://github.com/CAAPIM/layer7-operator/wiki)
- Techdocs (coming soon)

### Deployment Types

#### Application Level
- Database Backed Gateway Clusters
- Embedded Derby Database (Ephemeral) Gateways

#### Custom Resources
- Gateway
- Repository
- L7Api
- L7Portal


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
- External Secrets
- OpenTelemetry Integration (check out the examples!)

#### External Secrets
A new configuration option for external secrets has been created. This allows you to reference existing Kubernetes secrets which are synced with the Gateway's Stored Passwords for use in things like JDBC connections or policy.

External providers can be configured with the [external secrets operator](https://external-secrets.io).

```
app:
...
  externalSecrets:
    - name: database-credentials-gcp
      enabled: true
      description: GCP Database credentials
      variableReferencable: true
      auth:
        encryption:
          passphrase: 7layer
          existingSecret: ""
    - name: local-secret
      enabled: true
      description: local secret
      variableReferencable: true
      auth:
        encryption:
          passphrase: 7layer
          existingSecret: ""
```

#### Portal Integration (alpha)
The L7Portal controller is responsible for synchronizing Portal Managed APIs, the API Controller is responsible for applying those to target Gateway Deployments.

This integration handles CRUD on an individual API basis, the relationship between APIs and Gateways is deploymentTag.

Update operator image in deploy/bundle.yaml to harbor.sutraone.com/operator/layer7-operator:portal_integration and set imagePullPolicy to always.

There is a pre-configured example with an existing Portal Tenant containing 3 Portal APIs. To create your own update the [l7Portal CR]((./config/samples/security_v1alpha1_l7portal.yaml))
1. Create the Portal Enrolment bundle
  - kubectl create secret generic graphman-portal-bootstrap-bundle --from-file=./example/base/resources/secrets/bundles/portal-integration.json 
2. Apply the portal gateway [here](./example/gateway/portal-gateway.yaml)
  - kubectl apply -f ./example/gateway/portal-gateway.yaml
3. Apply the l7portal CR [here](./config/samples/security_v1alpha1_l7portal.yaml)
  - kubectl apply -f ./config/samples/security_v1alpha1_l7portal.yaml


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

## Installation
There are currently two ways to deploy the Layer7 Gateway Operator. A Helm Chart will be available in the future.

### Install with kubectl
Clone this repository to get started
```
git clone https://github.com/caapim/layer7-operator.git
cd layer7-operator
```
bundle.yaml contains all of the manifests that the Layer7 Operator requires.

#### OwnNamespace
By default the Operator manages the namespace that it is deployed into and does not create any cluster roles/role bindings. 

```
kubectl apply -f deploy/bundle.yaml
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
kubectl apply -f deploy/cw-bundle.yaml
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
  version: "11.0.00_CR1"
  license:
    accept: true
    secretName: gateway-license
  app:
    replicas: 1
    image: docker.io/caapim/gateway:11.0.00_CR1
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
kubectl get svc
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
kubectl port-forward svc/ssg 8443

username: admin
password: 7layer
gateway: localhost
```

5. Remove the Gateway Resource
```
kubectl delete gateway ssg
```

#### Remove the Operator
if you installed the operator using kubectl
```
kubectl delete -k deploy/bundle.yaml|cw-bundle.yaml
```

if you installed the operator in Openshift

``` 
operator-sdk cleanup <operatorPackageName> [flags]
```

