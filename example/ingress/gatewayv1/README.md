# Gateway API v1

By the end of this example you should have a better understanding of how the Layer7 Operator manages Kubernetes Gateway API resources. This example builds on the Advanced example with HTTPS termination, backend TLS, a shared Gateway, and TLS passthrough.

> [!WARNING]
> Gateway API CRDs are cluster-scoped. Use a dedicated or disposable cluster, such as the Kind cluster below, unless you own the cluster's Gateway API installation.
>
> - **Installing:** `make envoy-gateway` installs Envoy Gateway with Helm, which brings in the Gateway API CRDs. These apply to every namespace and can conflict with CRDs already installed by another Gateway API implementation or managed by your Kubernetes platform. Check for existing CRDs first with `kubectl get crd | grep gateway.networking.k8s.io`. If they exist, do not replace them. Use the cluster's existing implementation and GatewayClass instead, and check the installed version with `kubectl get crd gateways.gateway.networking.k8s.io -o jsonpath='{.metadata.annotations.gateway\.networking\.k8s\.io/bundle-version}'`.
> - **Removing:** Deleting the Gateway API CRDs deletes every Gateway, HTTPRoute, TLSRoute, and BackendTLSPolicy in the cluster, not only the resources from this example. Do not remove them from a shared cluster.
> - **Layer7 Operator:** The Operator registers its Gateway API watches at startup. Restart it after installing or removing the CRDs.

### Getting started

1. Use Layer7 Operator v1.3.1+ and Kubernetes v1.34+. The Gateway API CRDs must be the standard channel at v1.4+ for BackendTLSPolicy, or v1.5+ for TLSRoute (TLS passthrough). Your Gateway API implementation must also support the resources you enable.
2. Place a Gateway v11 license in [base/resources/secrets/license/](../../base/resources/secrets/license/) called `license.xml`.
3. Accept the Gateway License. Set `license.accept` to `true` in [gateway-api.yaml](../../gateway/gateway-api.yaml) before proceeding.

   ```yaml
   license:
     accept: true
     secretName: gateway-license
   ```

4. Set up [Kubernetes Gateway API](https://gateway-api.sigs.k8s.io/guides/). This example uses [Envoy Gateway](https://gateway.envoyproxy.io/docs/tasks/quickstart/) and GatewayClass `eg`.
5. The Makefile generates the example TLS certificates using OpenSSL. To use existing certificates, see [Create Certificates](#create-certificates).

### Quickstart

The example folder includes a Makefile that covers the deployment steps in this guide, **meaning you can move directly onto testing your Gateway Deployment.**

You can optionally use [Kind](https://kind.sigs.k8s.io/) (Kubernetes in Docker) with the supplied [kind-config.yaml](../../kind-config.yaml).

1. Go into the example directory.

   ```bash
   cd example
   ```

2. Optionally deploy a Kind Cluster.

   ```bash
   make kind-cluster
   ```

3. Deploy the Gateway API controller.

   ```bash
   make envoy-gateway
   ```

   For Kind, use `make envoy-gateway-kind` instead to expose ports 80 and 443 through the supplied Kind port mappings.

4. Deploy the example.

   ```bash
   make install kubernetes-gateway-example
   ```

5. Wait for all components to be ready.

   ```bash
   watch kubectl get pods
   ```

## If you used the Makefile proceed to [Test your Gateway Deployment](#test-your-gateway-deployment)

### Guide

Run the individual steps below from the repository root. Like the Basic and Advanced examples, these commands use the default namespace.

- [Deploy the Operator](#deploy-the-layer7-operator)
- [Create Repositories](#create-repositories)
- [Create Certificates](#create-certificates)
- [Create a Gateway](#create-a-gateway)
- [Test Gateway Deployment](#test-your-gateway-deployment)
- [Additional Configuration options](#additional-configuration-options)
- [Remove Kind Cluster](#remove-kind-cluster)
- [Remove Custom Resources](#remove-custom-resources)
- [Uninstall the Operator](#uninstall-the-operator)

### Deploy the Layer7 Operator

This step deploys the Layer7 Operator and its resources in namespaced mode. It manages Gateway and Repository Custom Resources in the namespace where it is deployed.

```bash
kubectl apply -f https://github.com/CAAPIM/layer7-operator/releases/download/v1.3.1/bundle.yaml
```

##### Verify the Operator is up and running

```bash
kubectl get pods -l app.kubernetes.io/name=layer7-operator
```

If the Operator was already running when the Gateway API CRDs were installed, restart it to register the Gateway API watches.

### Create Repositories

This example uses the same four repository references as the Advanced example. Create the repositories before the Gateway that references them. Applying the repositories also applies the base resources, including the license Secret.

- [l7-gw-myframework](../../repositories/framework-repository.yaml) — static
- [l7-gw-myapis](../../repositories/apis-repository.yaml) — dynamic
- [l7-gw-mysubscriptions](../../repositories/subscriptions-repository.yaml) — dynamic
- [local-reference-repository](../../repositories/local-reference-repository.yaml) — dynamic

```bash
kubectl apply -k ./example/repositories
```

##### Repository CR

```bash
kubectl get repositories
```

##### Operator Logs

```bash
kubectl logs -f -l app.kubernetes.io/name=layer7-operator
```

### Create Certificates

Generate and apply the example certificates:

```bash
cd example
make pki-gatewayv1-listener pki-gatewayv1-backend
cd ..
```

This creates the `gateway-listener-tls` and `gateway-backend-tls` TLS Secrets and the `gateway-backend-ca` ConfigMap. Both certificates cover `gateway.brcmlabs.com` and `gateway-management.brcmlabs.com`.

To use your own certificates, create Secrets with those names and place the issuing CA certificate in the ConfigMap under `ca.crt`. The backend certificate must cover `backendTLSPolicy.validation.hostname`.

### Create a Gateway

The [Gateway Custom Resource](../../gateway/gateway-api.yaml) uses HTTPS termination with backend TLS by default.

```bash
kubectl apply -f ./example/gateway/gateway-api.yaml
```

The Operator creates:

- A Kubernetes Gateway named `ssg-gateway` with HTTPS listeners on port 443.
- An HTTPRoute for `gateway.brcmlabs.com` to the main Service on port 8443.
- An HTTPRoute for `gateway-management.brcmlabs.com` to the management Service on port 9443.
- A BackendTLSPolicy using `gateway-backend-ca` to validate the backend certificate.

The backend key is applied through `spec.app.externalKeys` as the Gateway's default SSL key.

##### View your new Gateway

```bash
kubectl get pods
kubectl get gateways.gateway.networking.k8s.io
kubectl get httproutes.gateway.networking.k8s.io
kubectl get backendtlspolicies.gateway.networking.k8s.io
```

###### Gateway CR

```bash
kubectl get gateways.security.brcmlabs.com ssg -o yaml
```

### Test your Gateway Deployment

Find the address of the Kubernetes Gateway:

```bash
kubectl get gateways.gateway.networking.k8s.io ssg-gateway
```

Point `gateway.brcmlabs.com` and `gateway-management.brcmlabs.com` at that address using DNS or your hosts file. With `make envoy-gateway-kind`, use `127.0.0.1` when Docker runs locally.

Curl:

```bash
curl https://gateway.brcmlabs.com/api1 \
  -H "client-id: D63FA04C8447" --cacert /tmp/gatewayv1-ca.crt
```

If you used your own certificates, replace `/tmp/gatewayv1-ca.crt` with your CA certificate. To skip verification for a quick connectivity check, use `-k` instead.

Example response once the repositories have synchronised:

```json
{
  "client": "D63FA04C8447",
  "plan": "plan_a",
  "service": "hello api 1",
  "myDemoConfigVal": "suspiciousLlama"
}
```

#### Additional Configuration options

The alternative modes are commented sections in [gateway-api.yaml](../../gateway/gateway-api.yaml). Comment out the active `kubernetesGateway` block and uncomment the mode you want to use. Keep one block active at a time.

##### Shared Gateway

Create the example shared Gateway in the same namespace:

```bash
kubectl apply -f ./example/ingress/gatewayv1/shared-gateway.yaml
```

Enable the **Shared Gateway** block and apply the Layer7 Gateway:

```bash
kubectl apply -f ./example/gateway/gateway-api.yaml
```

The Operator creates the HTTPRoute and BackendTLSPolicy, references `shared-gateway`, and removes its previous managed Gateway. Update DNS to the shared Gateway's address.

##### TLS passthrough

Enable the **TLS passthrough** block and apply the Layer7 Gateway:

```bash
kubectl apply -f ./example/gateway/gateway-api.yaml
```

The Operator replaces its HTTPRoutes with a TLSRoute and removes the BackendTLSPolicy. The client receives the Layer7 Gateway's backend certificate directly.

##### Ingress

`spec.app.ingress` and `spec.app.kubernetesGateway` can both be enabled. Disabling either removes its resources without affecting the other.

##### View the Operator logs

```bash
kubectl logs -f -l app.kubernetes.io/name=layer7-operator
```

Invalid HTTPRoute matches or filters are reported during reconciliation. Check the Operator logs and the route's `Accepted` and `ResolvedRefs` conditions if routing fails.

### Remove Kind Cluster

If you deployed a dedicated Kind cluster for this example, remove it from the example directory:

```bash
cd example
make uninstall-kind
```

### Remove Custom Resources

From the repository root, remove the shared Gateway if you created it:

```bash
kubectl delete -f ./example/ingress/gatewayv1/shared-gateway.yaml
```

Remove the Layer7 Gateway and example certificates:

```bash
kubectl delete -f ./example/gateway/gateway-api.yaml
kubectl delete secret gateway-listener-tls gateway-backend-tls
kubectl delete configmap gateway-backend-ca
```

Remove the repositories and base resources if no other examples use them:

```bash
kubectl delete -k ./example/repositories
```

These steps leave Envoy Gateway and the Gateway API CRDs installed. Only remove the CRDs from a cluster dedicated to this example. Deleting them removes all Gateway API resources in every namespace.

### Uninstall the Operator

```bash
kubectl delete -f https://github.com/CAAPIM/layer7-operator/releases/download/v1.3.1/bundle.yaml
```
