# Basic
By the end of this example you should have a better understanding of the Layer7 Operator and the Custom Resources it manages. This example has a stronger focus on the Repository Custom Resource with a very basic Gateway Custom Resource that uses a L4 Loadbalancer.


### Getting started
1. Place a gateway v10 or v11 license in [base/resources/secrets/license/](../base/resources/secrets/license/) called license.xml.
2. Accept the Gateway License
  - license.accept defaults to false in [Gateway examples](../gateway/basic-gateway.yaml)
  - update license.accept to true before proceeding
  ```
  license:
    accept: true
    secretName: gateway-license
  ```
3. If you would like to create a TLS secret for your ingress controller then add tls.crt and tls.key to [base/resources/secrets/tls](../base/resources/secrets/tls)
    - these will be referenced later on.

### Quickstart
The example folder includes a makefile that covers the first 3 sections of the guide ***meaning you can move directly onto testing your Gateway Deployment.***

To avoid impacting an existing cluster you can optionally use [Kind](https://kind.sigs.k8s.io/)(Kubernetes in Docker). The default configuration will work if your docker engine is running locally, if you have a remote docker engine open [kind-config.yaml](../kind-config.yaml), uncomment lines 3-5 and set the apiServerAddress on line 4 to the IP Address of your remote docker engine.

1. Go into the example directory
```
cd example
```
2. Optionally deploy a Kind Cluster
```
make kind-cluster
```
3. Deploy the example
```
make install basic
```
4. Wait for all components to be ready
```
watch kubectl get pods
```
output
```
NAME                                                  READY   STATUS    RESTARTS   AGE
layer7-operator-controller-manager-8654fdb66c-zdswp   2/2     Running   0          5m59s
ssg-64ccd9dd48-bqstf                                  1/1     Running   0          61s
```

## If you used the Makefile proceed to [Test your Gateway Deployment](#test-your-gateway-deployment)


### Guide
- [Deploy the Operator](#deploy-the-layer7-operator)
- [Create Repositories](#create-repositories)
- [Create a Gateway](#create-a-gateway)
- [Test Gateway Deployment](#test-your-gateway-deployment)
- [Remove Kind Cluster](#remove-kind-cluster)
- [Remove Custom Resources](#remove-custom-resources)
- [Uninstall the Operator CRs](#uninstall-the-operator)

### Deploy the Layer7 Operator
This step will deploy the Layer7 Operator and all of its resources in namespaced mode. This means that it will only manage Gateway and Repository Custom Resources in the Kubernetes Namespace that it's deployed in.

```
kubectl apply -f https://github.com/CAAPIM/layer7-operator/releases/download/v1.0.8/bundle.yaml
```

##### Verify the Operator is up and running
```
kubectl get pods

NAME                                                  READY   STATUS    RESTARTS   AGE
layer7-operator-controller-manager-7647b58697-qd9vg   2/2     Running   0          27s
```

### Create Repositories
This example ships with 4 pre-configured Graphman repositories. The repository controller is responsible for synchronising these with the Operator and should always be created before Gateway resources that reference them to avoid race conditions. ***race conditions will be resolved automatically.***

- Git Repositories
  - [l7-gw-myframework](https://github.com/Gazza7205/l7GWMyFramework)
  - [l7-gw-mysubscriptions](https://github.com/Gazza7205/l7GWMySubscriptions)
  - [l7-gw-myapis](https://github.com/Gazza7205/l7GWMyAPIs)
- Local Repository

  ***Local Repositories eliminate the need for the repository controller to rely on an external service. This instead works with Kubernetes Secrets that contain graphman bundles.***
  - [local-reference-repository](../base/resources/secrets/repository/local-repository.json)

```
kubectl apply -k ./example/repositories
```

##### Operator Logs
```
kubectl logs -f $(kubectl get pods -oname | grep layer7-operator-controller-manager) manager
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

### Create a Gateway
```
kubectl apply -f ./example/gateway/basic-gateway.yaml
```

##### Referencing the repositories we created
[basic-gateway.yaml](../gateway/basic-gateway.yaml) contains 3 repository references, the 'type' defines how a repository is applied to the Container Gateway.
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
  - name: local-reference-repository
    enabled: true
    type: dynamic
    encryption: {}
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
```

You will see an init container starting with graphman-static-init.
```
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
```

##### View the Operator logs
```
kubectl logs -f $(kubectl get pods -oname | grep layer7-operator-controller-manager) manager
```

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
  image: caapim/gateway:11.1.1
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
version: 11.1.1
```

### Test your Gateway Deployment
```
kubectl get svc

NAME  TYPE           CLUSTER-IP     EXTERNAL-IP                PORT(S)                         AGE
ssg   LoadBalancer   10.68.4.161    ***<YOUR-EXTERNAL-IP>***   8443:31747/TCP,9443:30778/TCP   41m

if your output looks like this that means you don't have an External IP Provisioner in your Kubernetes Cluster. You can still access your Gateway using port-forward.

NAME  TYPE           CLUSTER-IP    EXTERNAL-IP     PORT(S)                         AGE
ssg   LoadBalancer   10.68.4.126   <PENDING>       8443:31384/TCP,9443:31359/TCP   7m39s
```

If EXTERNAL-IP is stuck in \<PENDING> state
- open a second terminal and run the following
```
kubectl port-forward svc/ssg 9443:9443
```

```
curl https://<YOUR-EXTERNAL-IP>:8443/api1 -H "client-id: D63FA04C8447" -k

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
This configuration removes port 2124, disables 8080 (HTTP) and hardens 8443 and 9443 where 9443 is the only port that allows a Policy Manager connection. The [advanced example](../gateway/advanced-gateway.yaml) shows how this can be customised with your own ports.

```
username: admin
password: 7layer
gateway: <YOUR-EXTERNAL-IP>:9443
```
or if you used port-forward
```
username: admin
password: 7layer
gateway: localhost:9443
```

### Remove Kind Cluster
If you used the Quickstart option and deployed Kind, all you will need to do is remove the Kind Cluster.

Make sure that you're in the example folder
```
pwd
```

output
```
/path/to/layer7-operator/example
```

Remove the Kind Cluster
```
make uninstall-kind
```

### Remove Custom Resources
```
kubectl delete -f ./example/gateway/basic-gateway.yaml
kubectl delete -k ./example/repositories/
```

### Uninstall the Operator
```
kubectl delete -f https://github.com/CAAPIM/layer7-operator/releases/download/v1.0.8/bundle.yaml
```