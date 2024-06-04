# Portal Integration
By the end of this example you should have a better understanding of the Layer7 Operator <==> Portal Integration. This example introduces the L7Portal and L7Api custom resources.

**NOTE** this example is experimental only! It **should not** be used in Portal deployments, and does not carry any official Broadcom support.


### Getting started
1. Place a gateway v10 or v11 license in [base/resources/secrets/license/](../base/resources/secrets/license/) called license.xml.
2. Place the registry credential for the Portal images in [./portal-integration/secrets](../portal-integration/secrets/) called docker-secret.yaml available on the [CA API Developer Portal Solutions & Patches](https://techdocs.broadcom.com/us/product-content/recommended-reading/technical-document-index/ca-api-developer-portal-solutions-and-patches.html) page.
3. Accept the Gateway License
  - license.accept defaults to false in [Gateway examples](../gateway/advanced-gateway.yaml)
  - update license.accept to true before proceeding
  ```
  license:
    accept: true
    secretName: gateway-license
  ```
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
    - **NOTE:** the Portal requires an ingress controller that supports ssl/tls passthrough for mutual ssl/tls.
        - This will add the following [command line argument](https://kubernetes.github.io/ingress-nginx/user-guide/cli-arguments/) to nginx, in the ingress-nginx namespace. **It has only been tested** deployed with the above commands
        - ```--enable-ssl-passthrough```
        - The following command will edit your nginx deployment
            - ```make configure-nginx-ssl-passthrough```
5. Resources
   - You will need a machine that is capable of running the Portal Core stack and the Gateway
     - At a minimum you should have 8(v)cpu and 16GB RAM allocated to your Kind instance or Kubernetes node.
6. DNS/Host file configuration
   - You will need the following entries in your hosts file or local DNS
     If you're running Kind locally
     ```
     127.0.0.1 gateway.brcmlabs.com portal.brcmlabs.com apim-dev-portal.brcmlabs.com dev-portal-ssg.brcmlabs.com dev-portal-enroll.brcmlabs.com dev-portal-sync.brcmlabs.com dev-portal-sso.brcmlabs.com dev-portal-analytics.brcmlabs.com dev-portal-broker.brcmlabs.com
     ```
     If you're running Kind on a remote VM
     ```
     <VIRTUAL-MACHINE-IP> gateway.brcmlabs.com portal.brcmlabs.com apim-dev-portal.brcmlabs.com dev-portal-ssg.brcmlabs.com dev-portal-enroll.brcmlabs.com dev-portal-sync.brcmlabs.com dev-portal-sso.brcmlabs.com dev-portal-analytics.brcmlabs.com dev-portal-broker.brcmlabs.com
     i.e.
     192.168.1.40 gateway.brcmlabs.com portal.brcmlabs.com apim-dev-portal.brcmlabs.com dev-portal-ssg.brcmlabs.com dev-portal-enroll.brcmlabs.com dev-portal-sync.brcmlabs.com dev-portal-sso.brcmlabs.com dev-portal-analytics.brcmlabs.com dev-portal-broker.brcmlabs.com
     ```
     ***NOTE*** If you are using an existing Kubernetes Cluster you can retrieve the correct address after the Prometheus Stack has been deployed

     ```
     kubectl get ingress
     ```
     output
     ```
     NAME                 CLASS   HOSTS                                                                                                    ADDRESS        PORTS     AGE
     portal-ingress       nginx   apim-dev-portal.brcmlabs.com,dev-portal-ssg.brcmlabs.com,dev-portal-analytics.brcmlabs.com + 5 more...   <ip-address>   80, 443   57m
     ```
     In your hosts file - the ingress address will be the same for the Gateway Ingress record
     ```
     <ip-address> gateway.brcmlabs.com portal.brcmlabs.com apim-dev-portal.brcmlabs.com dev-portal-ssg.brcmlabs.com dev-portal-enroll.brcmlabs.com dev-portal-sync.brcmlabs.com dev-portal-sso.brcmlabs.com dev-portal-analytics.brcmlabs.com dev-portal-broker.brcmlabs.com
     ```

     - We recommend sticking with the defaults to try out this experimental example as they are used to provision a Portal Tenant
       - If you wish to change the default you can do so in [portal-values.yaml](../portal-integration/portal-values.yaml)
         - set portal.domain
         ```yaml
         ...
         portal:
           domain: yourportaldomain.com
         ...
         ```
         - Your hosts file will need to use yourportaldomain.com in place of brcmlabs.com
         - The [enroll-payload](./enroll-payload.json) will also need to be updated
         - Finally set the following environment variable
           - export PORTAL_DOMAIN=yourportaldomain.com

### Guide
* [Quickstart](#quickstart)
    * [Using an existing Kubernetes Cluster](#existing-kubernetes-cluster)
    * [Using Kind](#kind)
- [Deploy the Operator](#deploy-the-layer7-operator)
- [Create a Gateway](#create-a-gateway)
- [Configure the Developer Portal](#test-your-gateway-deployment)
- [Remove Kind Cluster](#remove-kind-cluster)
- [Remove Custom Resources](#remove-custom-resources)
- [Uninstall the Operator CRs](#uninstall-the-operator)

## Quickstart
A Makefile is included in the example directory that makes deploying this example a one step process. If you have access to a Docker Machine you can use [Kind](https://kind.sigs.k8s.io/) (Kubernetes in Docker). This example can optionally deploy a Kind Cluster for you (you just need to make sure that you've [installed Kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation))

The kind configuration is in the base of the example folder. If your docker machine is remote (you are using a VM or remote machine) then uncomment the network section and set the apiServerAddress to the address of your VM/Remote machine
```
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
# networking:
#   apiServerAddress: "192.168.1.64"
#   apiServerPort: 6443
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
```

### Kind
```
make kind-cluster nginx-kind configure-nginx-ssl-passthrough portal-example
```

### Existing Kubernetes Cluster
```
make portal-example
```
if you don't have an ingress controller you can deploy nginx with the following
```
make nginx configure-nginx-ssl-passthrough
```
if you are using kind
```
make nginx-kind configure-nginx-ssl-passthrough
```

## If you used the Makefile proceed to [Test your Gateway Deployment](#test-your-gateway-deployment)

### Create a Gateway
The [Gateway Custom Resource](../gateway/portal-gateway.yaml) in this example has the following pre-configured

```
kubectl apply -f ./example/gateway/portal-gateway.yaml
```

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
kubectl logs -f $(kubectl get pods -oname | grep layer7-operator-controller-manager) manager
```

###### Portal CR


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

### Uninstall the Operator
```
kubectl delete -f 
```