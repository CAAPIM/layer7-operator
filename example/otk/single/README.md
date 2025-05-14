# OTK Single Mode Install
By the end of this example you should have a better understanding of the Layer7 Operator and ways you can use the OTK with Ephemeral Gateways. This example builds on the advanced example with a stronger focus on the OTK configuration in the Gateway Custom Resource. This example focuses on MySQL as the OTK Token store.

### Getting started
1. Place a gateway v11 license in [base/resources/secrets/license/](../base/resources/secrets/license/) called license.xml.
2. Accept the Gateway License
  - license.accept defaults to false in [Gateway examples](../gateway/advanced-gateway.yaml)
  - update license.accept to true before proceeding
  ```
  license:
    accept: true
    secretName: gateway-license
  ```
3. If you would like to create a TLS secret for your ingress controller then add tls.crt and tls.key to [base/resources/secrets/tls](../base/resources/secrets/tls)
    - these will be referenced later on.
4. You will need an ingress controller like nginx
    - if you do not have one installed already you can use the makefile in the example directory to deploy one
        - ```cd example```
        - Generic Kubernetes
            - ```make nginx```
        - Kind (Kubernetes in Docker)
            - follow the steps in Quickstart
            or
            - ```make nginx-kind```
    - return to the previous folder
        - ```cd ..```

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
3. Deploy Nginx (based on your environment)
```
make nginx
or
make nginx-kind
```
4. Deploy the example
```
make install mysql otk-single
```
5. Wait for all components to be ready
```
watch kubectl get pods
```
output
```
NAME                                                  READY   STATUS    RESTARTS   AGE
layer7-operator-controller-manager-8654fdb66c-ssjfd   2/2     Running   0          107s
mysql-0                                               1/1     Running   0          90s
ssg-7b7694d995-clwkk                                  1/1     Running   0          61s
ssg-7b7694d995-qptbj                                  1/1     Running   0          76s
```

## If you used the Makefile proceed to [Test your Gateway Deployment](#test-your-gateway-deployment)


### Guide
- [Deploy the Operator](#deploy-the-layer7-operator)
- [Create Repositories](#create-repositories)
- [The OTK Customizations Repository](#otk-customizations-repository)
- [Create a Gateway](#create-a-gateway)
- [Create Access Token](#create-access-token)
- [Test Gateway Deployment](#test-your-gateway-deployment)
- [Access OAuth Manager](#access-oauth-manager)
- [Remove Kind Cluster](#remove-kind-cluster)
- [Remove Custom Resources](#remove-custom-resources)
- [Uninstall the Operator CRs](#uninstall-the-operator)

### Deploy the Layer7 Operator
This step will deploy the Layer7 Operator and all of its resources in namespaced mode. This means that it will only manage Gateway and Repository Custom Resources in the Kubernetes Namespace that it's deployed in.

```
kubectl apply -f https://github.com/CAAPIM/layer7-operator/releases/download/v1.2.0/bundle.yaml
```

##### Verify the Operator is up and running
```
kubectl get pods

NAME                                                  READY   STATUS    RESTARTS   AGE
layer7-operator-controller-manager-7647b58697-qd9vg   2/2     Running   0          27s
```

### Create Repositories
This example ships with 2 pre-configured Graphman repository.

- [otk-customizations-single](https://github.com/Gazza7205/l7GWMyFramework)
- [local-reference-repository](../../base/resources/secrets/repository/local-repository.json)

```
kubectl apply -k ./example/repositories
```

##### Operator Logs
```
kubectl logs -f -l app.kubernetes.io/name=layer7-operator
```

##### Repository CR
The Repository Controller keeps tracks the latest available commit, where it's stored (if it's less than 1mb we create a Kubernetes secret) and when it was last updated. The Storage Secret is a gzipped graphman bundle (json) used in the Graphman Init Container to remove dependencies on git during deployment.

***Note: If the repository exceeds 1mb in compressed format each Graphman Init Container will clone it at runtime. This represents a single point of failure if your Git Server is down, we recommended creating your own initContainer with the larger graphman bundle.***
```
kubectl get repositories

NAME                         AGE
l7-gw-myapis                 10s
l7-gw-myframework            10s
l7-gw-mysubscriptions        10s
local-reference-repository   10s

kubectl get repository l7-gw-myapis -oyaml
...
status:
  commit: 3791f11c9b588b383ce87535f46d4fc1526ae83b
  name: l7-gw-myapis
  storageSecretName: l7-gw-myapis-repository
  updated: 2023-04-04 02:53:53.298060678 +0000 UTC m=+752.481758238
  vendor: Github
```

##### Local Reference Repositories
The local-reference-repository above uses a Kubernetes Secret as its source of truth. This removes external dependencies and therefore a potential single point of failure. Kubernetes Secrets are limited to 1MB in size.

The Layer7 Operator watches the secret `local-repo-secret` and updates Gateways that reference this repository if it changes.
```
apiVersion: security.brcmlabs.com/v1
kind: Repository
metadata:
  name: local-reference-repository
spec:
  enabled: true
  localReference:
    secretName: local-repo-secret
  type: local
  auth: {}
```

## OTK Customization Repository
Modifying the default behaviour and configuration of the OTK is important for a variety of reasons. The Layer7 Operator provides an easy way to do this using repositories.

We provide a [sample repo](link to sample repo) that has everthing pre-configured to work with this example. [This](link to community getting started repo) repo provides steps on how to create your own customizations repo using graphman.



### Create a Gateway
The [Gateway Custom Resource](../gateway/advanced-gateway.yaml) in this example has the following pre-configured

```
kubectl apply -f ./example/gateway/advanced-gateway.yaml
```

- InitContainers
This initContainer has a helloworld service and some very basic scripts that use echo.
```
initContainers:
- name: gateway-init
  image: docker.io/layer7api/gateway-init:1.0.0
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
Note: Pod and Security context options are currently limited to Gateway version 10.1.00_CR3,CR4 and 11.0.00_CR1,CR2
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

##### View your new Gateway
In this example we're using an Autoscaler, 1 node will be present while the autoscaler is initially configured.
```
kubectl get pods
NAME                   READY   STATUS    RESTARTS   AGE
...
ssg-7698bc565b-qrz5g   1/1     Running   0          2m45s
ssg-7698bc565b-szrbj   1/1     Running   0          2m45s
```

##### View the Operator logs
```
kubectl logs -f -l app.kubernetes.io/name=layer7-operator
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
  image: caapim/gateway:11.1.2
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
version: 11.1.2
```

### Create Access Token
In order to use of our gateway test endpoint we will need to get an access token

Run the following command
```
curl 'https://gateway.brcmlabs.com/auth/oauth/v2/token' -H 'Content-Type: application/x-www-form-urlencoded' -d 'client_id=54f0c455-4d80-421f-82ca-9194df24859d&client_secret=a0f2742f-31c7-436f-9802-b7015b8fd8e6&grant_type=password&username=admin&password=7layer'
```
output
```
{
  "access_token":"7326970c-f3e5-41f3-80c8-2f35dd05948b-1743426200",
  "token_type":"Bearer",
  "expires_in":3600,
  "scope":"oob",
  "resource": ["https://gateway.brcmlabs.com/*"]
}
```
Take note of your access_token, you will use that in the next step

### Test your Gateway Deployment
```
kubectl get ingress

NAME   CLASS   HOSTS                  ADDRESS              PORTS     AGE
ssg    nginx   gateway.brcmlabs.com,gateway-otk-management.brcmlabs.com   <YOUR-EXTERNAL-IP> or localhost   80, 443   54m
```

Add the following to your hosts file for DNS resolution
```
Format
$ADDRESS $HOST

example
<YOUR-EXTERNAL-IP> or 127.0.0.1 gateway.brcmlabs.com gateway-otk-management.brcmlabs.com
```
Curl
```
curl https://gateway.brcmlabs.com/v1/auth/test -H "Authorization: Bearer 7326970c-f3e5-41f3-80c8-2f35dd05948b-1743426200" -k
```
Response
```
{
  "message": "success",
  "user": "admin",
  "resource": "https://gateway.brcmlabs.com/*",
  "scope": "oob",
  "gateway": "bATu1QUbd7VJVZMqwISI8g+WefIJLL9e"
}
```

### Access OAuth Manager
To interact with OAuth Clients in the browser you can use the built-in OAuth Manager. This has been configured on a different hostname for this example.


1. Open a browser and navigate to https://gateway-otk-management.brcmlabs.com/oauth/manager
2. Accept the certificate warning and proceed
3. Enter the following credentials and click login
```
username: admin
password: 7layer
```
4. On the top menu select 'Clients' to see the OAuth Clients. These are pre-populated as this example loads test data.
5. On the top menu select 'Tokens' to see Tokens that have been issued.

Expired tokens will be removed automatically.


##### Sign into Policy Manager
Policy Manager access is less relevant in a deployment like this because we haven't specified an external MySQL database, any changes that we make will only apply to the Gateway that we're connected to and won't survive a restart. It is still useful to check what's been applied. We configured custom ports where we disabled Policy Manager access on 8443, we're also using an ingress controller meaning that port 9443 is not accessible without port forwarding.

Port-Forward
```
kubectl get pods
NAME                   READY   STATUS    RESTARTS   AGE
...
ssg-7698bc565b-qrz5g   1/1     Running   0          54m
ssg-7698bc565b-szrbj   1/1     Running   0          54m

kubectl port-forward ssg-7698bc565b-szrbj 9443:9443
```
Policy Manager
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
kubectl delete -f ./example/gateway/advanced-gateway.yaml
kubectl delete -k ./example/repositories/
```

### Uninstall the Operator
```
kubectl delete -f https://github.com/CAAPIM/layer7-operator/releases/download/v1.2.0/bundle.yaml
```