# Basic
By the end of this example you should have a better understanding of the Layer7 Operator and the Custom Resources it manages. This example has a stronger focus on the Repository Custom Resource with a very basic Gateway Custom Resource that uses a L4 Loadbalancer.


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
kubectl apply -f deploy/bundle.yaml

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
kubectl get pods

NAME                                                  READY   STATUS    RESTARTS   AGE
layer7-operator-controller-manager-7647b58697-qd9vg   2/2     Running   0          27s
```

#### Create Repository Custom Resources
This example ships with 3 pre-configured Graphman repositories. The repository controller is responsible for synchronising these with the Operator and should always be created before Gateway resources that reference them to avoid race conditions. ***race conditions will be resolved automatically.***

- [l7-gw-myframework](https://github.com/Gazza7205/l7GWMyFramework)
- [l7-gw-mysubscriptions](https://github.com/Gazza7205/l7GWMySubscriptions)
- [l7-gw-myapis](https://github.com/Gazza7205/l7GWMyAPIs)

```
kubectl apply -k example/repositories

secret/gateway-license configured
secret/gateway-secret unchanged
secret/graphman-encryption-secret unchanged
secret/graphman-repository-secret configured
repository.security.brcmlabs.com/l7-gw-myapis created
repository.security.brcmlabs.com/l7-gw-myframework created
repository.security.brcmlabs.com/l7-gw-mysubscriptions created
```

##### Operator Logs
```
kubectl logs <layer7-operator-pod> manager

...
1.6805762965185595e+09 INFO controllers.Repository Creating Storage Secret {"Name": "l7-gw-myapis-repository", "Namespace": "layer7"}
1.6805762965343177e+09 INFO controllers.Repository Reconciled {"Name": "l7-gw-myapis", "Namespace": "layer7", "Commit": "3791f11c9b588b383ce87535f46d4fc1526ae83b"}
1.680576296929594e+09 INFO controllers.Repository Creating Storage Secret {"Name": "l7-gw-myframework-repository", "Namespace": "layer7"}
1.6805762969402978e+09 INFO controllers.Repository Reconciled {"Name": "l7-gw-myframework", "Namespace": "layer7", "Commit": "c93028b807cf1b62bce0142a80ad4f6203207e8d"}
1.6805762973589563e+09 INFO controllers.Repository Creating Storage Secret {"Name": "l7-gw-mysubscriptions-repository", "Namespace": "layer7"}
1.6805762973709154e+09 INFO controllers.Repository Reconciled {"Name": "l7-gw-mysubscriptions", "Namespace": "layer7", "Commit": "fd6b225159fcd8fccf4bd61e31f40cdac64eccfa"} 
...

```

##### Repository CR
The Repository Controller keeps tracks the latest available commit, where it's stored (if it's less than 1mb we create a Kubernetes secret) and when it was last updated. The Storage Secret is a gzipped graphman bundle (json) used in the Graphman Init Container to remove dependencies on git during deployment.

***Note: If the repository exceeds 1mb in compressed format each Graphman Init Container will clone it at runtime. This represents a single point of failure if your Git Server is down, we recommended creating your own initContainer with the larger graphman bundle.***
```
kubectl get repositories

NAME                    AGE
l7-gw-myapis            10s
l7-gw-myframework       10s
l7-gw-mysubscriptions   10s

kubectl get repository l7-gw-myapis -oyaml
...
status:
  commit: 3791f11c9b588b383ce87535f46d4fc1526ae83b
  name: l7-gw-myapis
  storageSecretName: l7-gw-myapis-repository
  updated: 2023-04-04 02:53:53.298060678 +0000 UTC m=+752.481758238
  vendor: Github
```

#### Create a Gateway Custom Resource
```
kubectl apply -k example/basic/

serviceaccount/ssg-serviceaccount created
secret/gateway-license configured
secret/gateway-secret configured
secret/graphman-bootstrap-bundle configured
secret/graphman-encryption-secret configured
secret/graphman-repository-secret configured
secret/restman-bootstrap-bundle configured
gateway.security.brcmlabs.com/ssg created

```

##### Referencing the repositories we created
[ssg-gateway.yaml](./ssg-gateway.yaml) contains 3 repository references, the 'type' defines how a repository is applied to the Container Gateway.
- Dynamic repositories are applied directly to the Graphman endpoint on the Gateway which does not require a gateway restart
- Static repositories are bootstrapped to the Container Gateway with an initContainer which requires a gateway restart.
```
repositoryReferences:
  - name: l7-gw-myframework
    enabled: true
    type: ***static***
    encryption:
      existingSecret: graphman-encryption-secret
      key: FRAMEWORK_ENCRYPTION_PASSPHRASE
  - name: l7-gw-myapis
    enabled: true
    type: ***dynamic***
    encryption:
      existingSecret: graphman-encryption-secret
      key: APIS_ENCRYPTION_PASSPHRASE
  - name: l7-gw-mysubscriptions
    enabled: true
    type: ***dynamic***
    encryption:
      existingSecret: graphman-encryption-secret
      key: SUBSCRIPTIONS_ENCRYPTION_PASSPHRASE
```

##### View your new Gateway
```
kubectl get pods

NAME                                                  READY   STATUS    RESTARTS   AGE
layer7-operator-controller-manager-7647b58697-qd9vg   2/2     Running   0          15m
ssg-57d96567cb-n24g9                                  1/1     Running   0          73s
```

##### Static Graphman Repositories
Because we created the l7-gw-myframework repository reference with type 'static' the Layer7 Operator automatically injects an initContainer to bootstrap the repository to the Container Gateway.
Note: the suffix here graphman-static-init-***c1b58adb6d*** is generated using all static commit ids, if a static repository changes the Gateway will be updated.
```
kubectl describe pods ssg-57d96567cb-n24g9

...
Init Containers:
  graphman-static-init-c1b58adb6d:
    Container ID:   containerd://21924ae85d25437d3634ea5da1415c9bb58d678600f9fd67d4f0b0360857d7c5
    Image:          docker.io/layer7api/graphman-static-init:1.0.0
    Image ID:       docker.io/layer7api/graphman-static-init@sha256:24189a432c0283845664c6fd54c3e8d9f86ad9d35ef12714bb3a18b7aba85aa4
    Port:           <none>
    Host Port:      <none>
    State:          Terminated
      Reason:       Completed
      Exit Code:    0
      Started:      Tue, 04 Apr 2023 04:11:18 +0100
      Finished:     Tue, 04 Apr 2023 04:11:18 +0100
...
```
##### View the Graphman InitContainer logs
We should see that our static repository l7-gw-myframework has been picked up and moved to the bootstrap folder.
```
kubectl logs ssg-57d96567cb-n24g9 graphman-static-init-c1b58adb6d

l7-gw-myframework with 40kbs written to /opt/SecureSpan/Gateway/node/default/etc/bootstrap/bundle/graphman/0/0_l7-gw-myframework.json
```

##### View the Operator logs
```
kubectl logs layer7-operator-controller-manager-7647b58697-qd9vg manager

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
kubectl get gateway ssg -oyaml

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
kubectl get repository l7-gw-myapis -oyaml
...
status:
  commit: 7332f861e11612a91ca9de6b079826b9377dae6a
  name: l7-gw-myapis
  storageSecretName: l7-gw-myapis-repository
  updated: 2023-04-06 15:00:20.144406434 +0000 UTC m=+21.758241719
  vendor: Github
```

##### Test your Gateway Deployment
```
kubectl get svc

NAME  TYPE           CLUSTER-IP     EXTERNAL-IP         PORT(S)                         AGE
ssg   LoadBalancer   10.68.4.161    ***34.89.84.69***   8443:31747/TCP,9443:30778/TCP   41m

if your output looks like this that means you don't have an External IP Provisioner in your Kubernetes Cluster. You can still access your Gateway using port-forward.

NAME  TYPE           CLUSTER-IP    EXTERNAL-IP     PORT(S)                         AGE
ssg   LoadBalancer   10.68.4.126   <PENDING>       8443:31384/TCP,9443:31359/TCP   7m39s
```

If EXTERNAL-IP is stuck in \<PENDING> state
```
kubectl port-forward svc/ssg 9443:9443
```

```
curl https://34.89.84.69:8443/api1 -H "client-id: D63FA04C8447" -k

or if you used port-forward

curl https://localhost:9443/api1 -H "client-id: D63FA04C8447" -k

```
Response
```
{
  "client" : "D63FA04C8447",
  "plan" : "plan_a",
  "service" : "hello api 1",
  "myDemoConfigVal" : "suspiciousLlama"
}
```

##### Sign into Policy Manager
Policy Manager access is less relevant in a deployment like this because we haven't specified an external MySQL database, any changes that we make will only apply to the Gateway that we're connected to and won't survive a restart. It is still useful to check what's been applied. In our configuration we set the following which overrides the default application port configuration.
```
...
listenPorts:
  harden: true
...
```
This configuration removes port 2124, disables 8080 (HTTP) and hardens 8443 and 9443 where 9443 is the only port that allows a Policy Manager connection. The [advanced example](../advanced/ssg-gateway.yaml) shows how this can be customised with your own ports.

```
username: admin
password: 7layer
gateway: 35.189.116.20:9443
```
or if you used port-forward
```
username: admin
password: 7layer
gateway: localhost:9443
```


#### Remove Custom Resources
```
kubectl delete -k example/basic/
kubectl delete -k example/repositories/

secret "gateway-license" deleted
secret "gateway-secret" deleted
secret "graphman-encryption-secret" deleted
secret "graphman-repository-secret" deleted
gateway.security.brcmlabs.com "ssg" deleted
repository.security.brcmlabs.com "l7-gw-myapis" deleted
repository.security.brcmlabs.com "l7-gw-myframework" deleted
repository.security.brcmlabs.com "l7-gw-mysubscriptions" deleted
```

### Uninstall the Operator
```
kubectl delete -f deploy/bundle.yaml

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