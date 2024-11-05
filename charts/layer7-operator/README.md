# Layer7 Operator Helm Chart
This Helm Chart installs the Layer7 Operator

## Requirements
- This Chart requires cluster privileges to install
- [Certmanager](https://cert-manager.io/docs/installation/kubectl/)
- Helm >= 3.7

## Deploy Cert-Manager
This chart depends on cert-manager. If you do not already have cert-manager deployed, run the following command

Please check [Certmanager](https://cert-manager.io/docs/installation/kubectl/) for any updated installation instructions
```
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.3/cert-manager.yaml
```

## Getting Started
Add the layer7-operator repository
```
helm repo add layer7-operator https://caapim.github.io/layer7-operator/
``` 
Update your repositories
```
helm repo update
```
Install the Operator
```
helm upgrade -i layer7-operator layer7-operator/layer7-operator -n layer7-operator-system --create-namespace
```

## Parameters

### Common Parameters

| Name                | Description                                   | Value           |
| ------------------- | --------------------------------------------- | --------------- |
| `nameOverride`      | String to partially override fullname         | `""`            |
| `fullnameOverride`  | String to fully override fullname             | `""`            |
| `clusterDomain`     | Kubernetes cluster domain name                | `cluster.local` |
| `commonLabels`      | Labels to add to all deployed objects         | `{}`            |
| `commonAnnotations` | Annotations to add to all deployed objects    | `{}`            |
| `podLabels`         | Labels to add to the Layer7 Operator Pod      | `{}`            |
| `podAnnotations`    | Annotations to add to the Layer7 Operator Pod | `{}`            |

### RBAC Parameters

| Name                         | Description                                          | Value  |
| ---------------------------- | ---------------------------------------------------- | ------ |
| `serviceAccount.create`      | Specifies whether a ServiceAccount should be created | `true` |
| `serviceAccount.annotations` | Additional custom annotations for the ServiceAccount | `{}`   |
| `serviceAccount.name`        | The name of the ServiceAccount to use.               | `""`   |
| `rbac.create`                | Specifies whether RBAC resources should be created   | `true` |

### Layer7 Operator Parameters

| Name                              | Description                                                                                                    | Value                                                                                                                                                     |
| --------------------------------- | -------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `managedNamespaces`               | Namespaces that the Operator will manage. By default it will watch all namespaces.                             | `[""]`                                                                                                                                                    |
| `replicas`                        | Number of Layer7 Operator replicas. This value should not be changed\                                          | `1`                                                                                                                                                       |
| `webhook.enabled`                 | This creates Validating and Mutating Webhook configurations                                                    | `true`                                                                                                                                                    |
| `webhook.tls.certmanager.enabled` | This creates a self-signed issuer and cert-manager certificate, cert-manager is required if this is true       | `true`                                                                                                                                                    |
| `webhook.tls.existingTlsSecret`   | This allows you use an existing secret of type kubernetes.io/tls                                               | `webhook-cert-secret`                                                                                                                                     |
| `podSecurityContext`              | Layer7 Operator Pod Security Context                                                                           | `{}`                                                                                                                                                      |
| `containerSecurityContext`        | Layer7 Operator Container Security Context                                                                     | `{}`                                                                                                                                                      |
| `image.registry`                  | Layer7 Operator image registry                                                                                 | `REGISTRY_NAME`                                                                                                                                           |
| `image.repository`                | Layer7 Operator image repository                                                                               | `REPOSITORY_NAME/layer7api`                                                                                                                               |
| `image.tag`                       | Layer7 Operator image tag                                                                                      | `v1.1.0`                                                                                                                                                  |
| `image.pullPolicy`                | Layer7 Operator image pull policy                                                                              | `IfNotPresent`                                                                                                                                            |
| `image.pullSecrets`               | Layer7 Operator image pull secrets                                                                             | `[]`                                                                                                                                                      |
| `resources.limits.cpu`            | The cpu limits for the Layer7 Operator container                                                               | `500m`                                                                                                                                                    |
| `resources.limits.memory`         | The memory limits for the Layer7 Operator container                                                            | `100Mi`                                                                                                                                                   |
| `resources.requests.cpu`          | The cpu requests for the Layer7 Operator container                                                             | `100m`                                                                                                                                                    |
| `resources.requests.memory`       | The memory requests for Layer7 Operator container                                                              | `20Mi`                                                                                                                                                    |
| `args`                            | The arguments to pass to the Layer7 Operator Container. Setting --zap-log-level=10 will increase log verbosity | `["--health-probe-bind-address=:8081","--metrics-bind-address=127.0.0.1:8080","--leader-elect","--zap-log-level=info","--zap-time-encoding=rfc3339nano"]` |
| `otel.enabled`                    | Enable OpenTelemetry Metrics for the Layer7 Operator                                                           | `false`                                                                                                                                                   |
| `otel.otlpEndpoint`               | OTel Collector GRPC endpoint                                                                                   | `localhost:4317`                                                                                                                                          |
| `otel.metricPrefix`               | OTel metric prefix that will be prepended to each metric emitted from the Layer7 Operator                      | `layer7_`                                                                                                                                                 |

### RBAC Proxy Parameters

| Name                                                          | Description                                                    | Value                                                                                                       |
| ------------------------------------------------------------- | -------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------- |
| `rbacProxy.image.registry`                                    | Layer7 Operator image registry                                 | `REGISTRY_NAME`                                                                                             |
| `rbacProxy.image.repository`                                  | Layer7 Operator image repository                               | `REPOSITORY_NAME/layer7api`                                                                                 |
| `rbacProxy.image.tag`                                         | Layer7 Operator image tag                                      | `v0.13.1`                                                                                                   |
| `rbacProxy.image.pullPolicy`                                  | Redis&reg; image pull policy                                   | `IfNotPresent`                                                                                              |
| `rbacProxy.image.pullSecrets`                                 | Redis&reg; image pull secrets                                  | `[]`                                                                                                        |
| `rbacProxy.args`                                              | The arguments to pass to the RBAC Proxy Container.             | `["--secure-listen-address=0.0.0.0:8443","--upstream=http://127.0.0.1:8080/","--logtostderr=true","--v=0"]` |
| `rbacProxy.containerSecurityContext.allowPrivilegeEscalation` | Is it possible to escalate the RBAC Proxy container privileges | `false`                                                                                                     |
| `rbacProxy.containerSecurityContext.capabilities.drop`        | Set RBAC Proxy container Security Context capabilities to drop | `["ALL"]`                                                                                                   |
| `rbacProxy.resources.limits.cpu`                              | The cpu limits for the RBAC Proxy container                    | `500m`                                                                                                      |
| `rbacProxy.resources.limits.memory`                           | The memory limits for the RBAC Proxy container                 | `128Mi`                                                                                                     |
| `rbacProxy.resources.requests.cpu`                            | The cpu requests for the RBAC Proxy container                  | `5m`                                                                                                        |
| `rbacProxy.resources.requests.memory`                         | The memory requests for RBAC Proxy container                   | `64Mi`                                                                                                      |

### Proxy Configuration

| Name                               | Description                                                                                                                                                                                       | Value   |
| ---------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------- |
| `proxy.httpProxy`                  | HTTP proxy                                                                                                                                                                                        | `nil`   |
| `proxy.httpsProxy`                 | HTTPS proxy                                                                                                                                                                                       | `nil`   |
| `proxy.noProxy`                    | Proxy exclusion                                                                                                                                                                                   | `nil`   |
| `proxy.caBundle.enabled`           | Mount a configmap to the Layer7 Operator Container with a Trusted CA bundle                                                                                                                       | `false` |
| `proxy.caBundle.existingConfigmap` | Existing configmap containing a ca bundle                                                                                                                                                         | `nil`   |
| `proxy.caBundle.create`            | Create the ca bundle                                                                                                                                                                              | `false` |
| `proxy.caBundle.key`               | Existing configmap key that has the ca-bundle. Set this if you are specifying your own configmap or if you are using a label to inject a trusted ca bundle into the configMap this Chart creates. | `nil`   |
| `proxy.caBundle.labels`            | Labels to add to the created ca bundle                                                                                                                                                            | `{}`    |
| `proxy.caBundle.annotations`       | to add the created ca bundle                                                                                                                                                                      | `{}`    |
| `proxy.caBundle.pem`               | optional even when create is true given that certain labels will automatically inject the ca contents.                                                                                            | `nil`   |

