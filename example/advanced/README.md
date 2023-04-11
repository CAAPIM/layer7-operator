# Advanced
By the end of this example you should have a better understanding of the Layer7 Operator and the Custom Resources it manages. This example builds on the basic example with a stronger focus on the Gateway Custom Resource.

### Getting started
1. Place a gateway v10 or v11 license in [base/resources/secrets/license/](../base/resources/secrets/license/) called license.xml.
2. If you would like to create a TLS secret for your ingress controller then add tls.crt and tls.key to [base/resources/secrets/tls](../base/resources/secrets/tls)
    - these will be referenced later on.

### Guide
- [Deploy the Operator](#deploy-the-layer7-operator)
- [Create Repositories](#create-repository-custom-resources)
- [Create a Gateway](#create-a-gateway-custom-resource)
- [Test Gateway Deployment](#test-your-gateway-deployment)
- [Remove Custom Resources](#remove-custom-resources)
- [Uninstall the Operator CRs](#uninstall-the-operator)

#### Deploy the Layer7 Operator
This step will deploy the Layer7 Operator and all of its resources in namespaced mode. This means that it will only manage Gateway and Repository Custom Resources in the Kubernetes Namespace that it's deployed in.

```
$ kubectl apply -f deploy/bundle.yaml

customresourcedefinition.apiextensions.k8s.io/gateways.security.brcmlabs.com created
customresourcedefinition.apiextensions.k8s.io/repositories.security.brcmlabs.com created
serviceaccount/layer7-operator-controller-manager created
role.rbac.authorization.k8s.io/layer7-operator-leader-election-role created
role.rbac.authorization.k8s.io/layer7-operator-manager-role created
role.rbac.authorization.k8s.io/layer7-operator-proxy-role created
rolebinding.rbac.authorization.k8s.io/layer7-operator-leader-election-rolebinding created
rolebinding.rbac.authorization.k8s.io/layer7-operator-manager-rolebinding created
rolebinding.rbac.authorization.k8s.io/layer7-operator-proxy-rolebinding created
configmap/layer7-operator-manager-config created
service/layer7-operator-controller-manager-metrics-service created
deployment.apps/layer7-operator-controller-manager created
```

##### Verify the Operator is up and running
```
$ kubectl get pods

NAME                                                  READY   STATUS    RESTARTS   AGE
layer7-operator-controller-manager-7647b58697-qd9vg   2/2     Running   0          27s
```

#### Create Repository Custom Resources
This example ships with 3 pre-configured Graphman repositories. The repository controller is responsible for synchronising these with the Operator and should always be created before Gateway resources that reference them to avoid race conditions. ***race conditions will be resolved automatically.***

- [l7-gw-myframework](https://github.com/Gazza7205/l7GWMyFramework)
- [l7-gw-mysubscriptions](https://github.com/Gazza7205/l7GWMySubscriptions)
- [l7-gw-myapis](https://github.com/Gazza7205/l7GWMyAPIs)

```
$ kubectl apply -k example/repositories

secret/gateway-license configured
secret/gateway-secret unchanged
secret/graphman-encryption-secret unchanged
secret/graphman-repository-secret configured
secret/harbor-reg-cred configured
repository.security.brcmlabs.com/l7-gw-myapis created
repository.security.brcmlabs.com/l7-gw-myframework created
repository.security.brcmlabs.com/l7-gw-mysubscriptions created
```

#### Create a Gateway Custom Resource
The [Gateway Custom Resource](./ssg-gateway.yaml) in this example has the following pre-configured

#### Configured in this example
- InitContainers
This initContainer has a helloworld service and some very basic scripts that use echo.
```
initContainers:
- name: gateway-init
  image: harbor.sutraone.com/operator/gateway-init:1.0.4
  imagePullPolicy: IfNotPresent
  volumeMounts:
  - name: config-directory
    mountPath: /opt/docker/custom
```

The bootstrap script works with the initContainer by moving everything in the shared volume (/opt/docker/custom) to the correct locations on the Gateway for bootstrap. This functionality is separate from Operator managed Repositories meaning the Operator will not automatically sync these and both Restman and Graphman bundles can be used.

```
bootstrap:
  script:
    enabled: true
```
- Bundles (cluster properties)
  - 1 Restman bundle
  - 1 Graphman bundle
```
bundle:
  - type: restman
    source: secret
    name: restman-bootstrap-bundle
  - type: graphman
    source: secret
    name: graphman-bootstrap-bundle
```
- Custom Gateway Ports
```
listenPorts:
  harden: false
  custom:
    enabled: true
  ports:
  - name: Default HTTPS (8443)
    port: "8443"
    enabled: true
    protocol: HTTPS
    managementFeatures:
    - Published service message input
    # - Administrative access
    # - Browser-based administration
    # - Built-in services
    properties:
    - name: server
      value: A
    tls:
      enabled: true
      #privateKey: 00000000000000000000000000000002:ssl
      clientAuthentication: Optional
      versions:
      - TLSv1.2
      - TLSv1.3
      useCipherSuitesOrder: true
      cipherSuites:
      - TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
      - TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384
      - TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384
      - TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA384
      - TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA
      - TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA
      - TLS_DHE_RSA_WITH_AES_256_GCM_SHA384
      - TLS_DHE_RSA_WITH_AES_256_CBC_SHA256
      - TLS_DHE_RSA_WITH_AES_256_CBC_SHA
      - TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
      - TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256
      - TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256
      - TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256
      - TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA
      - TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA
      - TLS_DHE_RSA_WITH_AES_128_GCM_SHA256
      - TLS_DHE_RSA_WITH_AES_128_CBC_SHA256
      - TLS_DHE_RSA_WITH_AES_128_CBC_SHA
      - TLS_AES_256_GCM_SHA384
      - TLS_AES_128_GCM_SHA256
...
```
- Autoscaling
```
autoscaling:
  enabled: true
  hpa:
    minReplicas: 2
    maxReplicas: 3
    metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 60
    behavior:
      scaleDown:
        stabilizationWindowSeconds: 300
        policies:
        - type: Pods
          value: 1
          periodSeconds: 60
      scaleUp:
        stabilizationWindowSeconds: 0
        policies:
        - type: Percent
          value: 100
          periodSeconds: 15
```
- Security Context
Allows setting container and pod security contexts.
Note: Pod and Security context options are currently limited to Gateway version 10.1.00_CR3
```
containerSecurityContext:
  runAsNonRoot: true
  runAsUser: 3000
  capabilities:
    drop:
    - ALL
  allowPrivilegeEscalation: false
podSecurityContext:
  runAsUser: 3000
  runAsGroup: 3000
  fsGroup: 3000
```
- Pod Disruption Budget
```
pdb:
  enabled: true
  minAvailable: 1
  maxUnavailable: 0
```

#### Additional Configuration options
- Tolerations
- TopologySpreadConstraints
- Affinity
- NodeSelector

#### Deploy the Gateway CR
```
$ kubectl apply -k example/advanced/

serviceaccount/ssg-serviceaccount created
secret/gateway-license configured
secret/gateway-secret configured
secret/graphman-bootstrap-bundle configured
secret/graphman-encryption-secret configured
secret/graphman-repository-secret configured
secret/restman-bootstrap-bundle configured
gateway.security.brcmlabs.com/ssg created
```

##### View your new Gateway
In this example we're using an Autoscaler, 1 node will be present while the autoscaler is initially configured.
```
$ kubectl get pods
NAME                   READY   STATUS    RESTARTS   AGE
...
ssg-7698bc565b-qrz5g   1/1     Running   0          2m45s
ssg-7698bc565b-szrbj   1/1     Running   0          2m45s
```
The Operator also keeps status of any given Gateway CR up-to-date so you also run the following
```
$ kubectl get gateways
NAME   AGE
ssg    2m45s


$ kubectl get gateway ssg -oyaml

...
status:
  conditions:
  ...
  gateway:
  - name: ssg-8589cbc944-r6292
    phase: Running
    ready: true
    startTime: 2023-04-06 15:05:41 +0000 UTC
  - name: ssg-8589cbc944-q4n67
    phase: Running
    ready: true
    startTime: 2023-04-06 15:05:56 +0000 UTC
  host: gateway.brcmlabs.com
  image: docker.io/caapim/gateway:10.1.00_CR3
  ready: 2
  replicas: 2
  repositoryStatus:
  - branch: main
    commit: 9175df94cc51e650d97a40c06b7855496a1612e2
    enabled: true
    endpoint: https://github.com/Gazza7205/l7GWMyFramework
    name: l7-gw-myframework
    secretName: graphman-repository-secret
    storageSecretName: l7-gw-myframework-repository
    type: static
  ...
  state: Ready
  version: 10.1.00_CR3

```


##### View the Operator logs
```
$ kubectl logs layer7-operator-controller-manager-7647b58697-qd9vg manager

...
1.6805472375519047e+09  INFO    Starting workers        {"controller": "gateway", "controllerGroup": "security.brcmlabs.com", "controllerKind": "Gateway", "worker count": 1}
1.6805472375519912e+09  INFO    Starting workers        {"controller": "repository", "controllerGroup": "security.brcmlabs.com", "controllerKind": "Repository", "worker count": 1}
1.6805480463029926e+09  INFO    controllers.Gateway     Creating ConfigMap      {"Name": "ssg", "Namespace": "layer7"}
1.680548046309193e+09   INFO    controllers.Gateway     Creating ConfigMap      {"Name": "ssg-system", "Namespace": "layer7"}
1.6805480463136642e+09  INFO    controllers.Gateway     Creating ConfigMap      {"Name": "ssg-cwp-bundle", "Namespace": "layer7"}
1.6805480463188894e+09  INFO    controllers.Gateway     Creating ConfigMap      {"Name": "ssg-listen-port-bundle", "Namespace": "layer7"}
1.680548046426919e+09   INFO    controllers.Gateway     Creating Service        {"Name": "ssg", "Namespace": "layer7"}
1.6805480465468638e+09  INFO    controllers.Gateway     Deployment hasn't been created yet      {"Name": "ssg", "Namespace": "layer7"}
1.6805480466609669e+09  INFO    controllers.Gateway     Creating ConfigMap      {"Name": "ssg-repository-init-config", "Namespace": "layer7"}
1.6805480466660128e+09  INFO    controllers.Gateway     Creating Deployment     {"Name": "ssg", "Namespace": "layer7"}
1.6805480472615528e+09  INFO    controllers.Repository  Creating Storage Secret {"Name": "l7-gw-myframework", "Namespace": "layer7"}
1.680548047275876e+09   INFO    controllers.Repository  Reconciled      {"Name": "l7-gw-myframework", "Namespace": "layer7", "Commit": "4b6c3ff1f174e4095ceadb31153392084fbaa61b"}
1.6805786502375867e+09  INFO    controllers.Gateway     Applying Latest Commit  {"Repo": "l7-gw-myapis", "Directory": "/", "Commit": "3791f11c9b588b383ce87535f46d4fc1526ae83b", "Pod": "ssg-57d96567cb-n24g9", "Name": "ssg", "Namespace": "layer7"}
1.6805786509813132e+09  INFO    controllers.Gateway     Applying Latest Commit  {"Repo": "l7-gw-mysubscriptions", "Directory": "/", "Commit": "fd6b225159fcd8fccf4bd61e31f40cdac64eccfa", "Pod": "ssg-57d96567cb-n24g9", "Name": "ssg", "Namespace": "layer7"}
...

```

##### Inspect the Status of your Custom Resources

###### Gateway CR
The Gateway Controller tracks gateway pods and the repositories that have been applied to the deployment
```
$ kubectl get gateway ssg -oyaml

status:
 ...
  gateway:
  - name: ssg-6b7d7fd999-n5bsj
    phase: Running
    ready: true
    startTime: 2023-04-03 18:57:24 +0000 UTC
  host: gateway.brcmlabs.com
  image: caapim/gateway:10.1.00_CR3
  ready: 1
  replicas: 1
repositoryStatus:
- branch: main
  commit: c93028b807cf1b62bce0142a80ad4f6203207e8d
  enabled: true
  endpoint: https://github.com/Gazza7205/l7GWMyFramework
  name: l7-gw-myframework
  secretName: graphman-repository-secret
  storageSecretName: l7-gw-myframework-repository
  type: static
- branch: main
  commit: 3791f11c9b588b383ce87535f46d4fc1526ae83b
  enabled: true
  endpoint: https://github.com/Gazza7205/l7GWMyAPIs
  name: l7-gw-myapis
  secretName: graphman-repository-secret
  storageSecretName: l7-gw-myapis-repository
  type: dynamic
- branch: main
  commit: fd6b225159fcd8fccf4bd61e31f40cdac64eccfa
  enabled: true
  endpoint: https://github.com/Gazza7205/l7GWMySubscriptions
  name: l7-gw-mysubscriptions
  secretName: graphman-repository-secret
  storageSecretName: l7-gw-mysubscriptions-repository
  type: dynamic
state: Ready
version: 10.1.00_CR3
```

###### Repository CR
The Repository Controller keeps tracks the latest available commit, where it's stored (if it's less than 1mb we create a Kubernetes secret) and when it was last updated.
```
$ kubectl get repositories

NAME                    AGE
l7-gw-myapis            10s
l7-gw-myframework       10s
l7-gw-mysubscriptions   10s

$ kubectl get repository l7-gw-myapis -oyaml
...
status:
  commit: 3791f11c9b588b383ce87535f46d4fc1526ae83b
  name: l7-gw-myapis
  storageSecretName: l7-gw-myapis-repository
  updated: 2023-04-04 02:53:53.298060678 +0000 UTC m=+752.481758238
  vendor: Github
```

##### Test your Gateway Deployment
```
$ kubectl get ingress

NAME   CLASS   HOSTS                  ADDRESS        PORTS     AGE
ssg    nginx   gateway.brcmlabs.com   34.89.126.80   80, 443   54m
```

Add the following to your hosts file for DNS resolution
```
Format
$ADDRESS $HOST

example
34.89.126.80 gateway.brcmlabs.com
```
Curl
```
$ curl https://gateway.brcmlabs.com/api1 -H "client-id: D63FA04C8447" -k
```
##### Sign into Policy Manager
Policy Manager access is less relevant in a deployment like this because we haven't specified an external MySQL database, any changes that we make will only apply to the Gateway that we're connected to and won't survive a restart. It is still useful to check what's been applied. We configured custom ports where we disabled Policy Manager access on 8443, we're also using an ingress controller meaning that port 9443 is not accessible without port forwarding.

Port-Forward
```
$ kubectl get pods
NAME                   READY   STATUS    RESTARTS   AGE
...
ssg-7698bc565b-qrz5g   1/1     Running   0          54m
ssg-7698bc565b-szrbj   1/1     Running   0          54m

$ kubectl port-forward ssg-7698bc565b-szrbj 9443:9443
```
Policy Manager
```
username: admin
password: 7layer
gateway: localhost:9443
```

#### Remove Custom Resources
```
$ kubectl delete -k example/basic/
$ kubectl delete -k example/repositories/

secret "gateway-license" deleted
secret "gateway-secret" deleted
secret "graphman-encryption-secret" deleted
secret "graphman-repository-secret" deleted
secret "harbor-reg-cred" deleted
gateway.security.brcmlabs.com "ssg" deleted
repository.security.brcmlabs.com "l7-gw-myapis" deleted
repository.security.brcmlabs.com "l7-gw-myframework" deleted
repository.security.brcmlabs.com "l7-gw-mysubscriptions" deleted
```

### Uninstall the Operator
```
$ kubectl delete -f deploy/bundle.yaml

customresourcedefinition.apiextensions.k8s.io "gateways.security.brcmlabs.com" deleted
customresourcedefinition.apiextensions.k8s.io "repositories.security.brcmlabs.com" deleted
serviceaccount "layer7-operator-controller-manager" deleted
role.rbac.authorization.k8s.io "layer7-operator-leader-election-role" deleted
role.rbac.authorization.k8s.io "layer7-operator-manager-role" deleted
role.rbac.authorization.k8s.io "layer7-operator-proxy-role" deleted
rolebinding.rbac.authorization.k8s.io "layer7-operator-leader-election-rolebinding" deleted
rolebinding.rbac.authorization.k8s.io "layer7-operator-manager-rolebinding" deleted
rolebinding.rbac.authorization.k8s.io "layer7-operator-proxy-rolebinding" deleted
configmap "layer7-operator-manager-config" deleted
service "layer7-operator-controller-manager-metrics-service" deleted
deployment.apps "layer7-operator-controller-manager" deleted
```