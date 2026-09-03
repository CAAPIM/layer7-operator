# Layer7 Operator Examples
These examples cover a broader set of the features that the Layer7 Operator provides and serve as a starting point for implementing these in your own environments or just trying them out.

## Prerequisites
- Kubernetes v1.34+
- Gateway v11.x License
- Ingress Controller (You can also expose Gateway Services as L4 LoadBalancers)
- Accept the Gateway License
  - license.accept defaults to false in all of the [Gateway examples](./gateway/)
  - update license.accept to true before proceeding

> **Breaking change for users upgrading from v1.2.x or earlier:** All example Gateway YAML files and the `Makefile` quickstart now use **[Contour](https://projectcontour.io/)** as the ingress controller (`ingressClassName: contour`, `projectcontour.io/upstream-protocol.tls` annotations). If you were using the previous nginx-based examples, your ingress configuration is **incompatible** and must be updated. Replace `ingressClassName: nginx` with `ingressClassName: contour` and substitute any `nginx.ingress.kubernetes.io/*` annotations with their Contour equivalents. Run `make contour` (or `make contour-kind` for Kind clusters) to install Contour before applying the examples.

The basic and advanced examples can be run in a single namespace, The OTel Examples require multiple namespaces for the additional components. Your Kubernetes user or service account must have sufficient privileges to create namespaces, deployments, configmaps, secrets, service accounts, roles, etc..

Each example also includes a [Kind](https://kind.sigs.k8s.io/) (Kubernetes in Docker) Quickstart which you can utilise if you have access to a Docker Machine.

#### Getting started
1. Place a gateway v11 license in [./base/resources/secrets/license](./base/resources/secrets/license).
2. If you would like to create a TLS secret for your ingress controller then add tls.crt and tls.key to [./base/resources/secrets/tls](./base/resources/secrets/tls)
    - these will be referenced later on.

#### Examples
All examples use [kustomize](https://kustomize.io/). The basic example covers a simple deployment with a single Gateway 3 Repository CRs (custom resource) configured for static and dynamic updates resepectively with a focus on repositories. Building on the basic example, the advanced example focuses on gateway configuration.

Gateways
- [Basic](./basic)
- [Advanced](./advanced)
- OTK
  - [Single](./otk/single)

Open Telemetry Example
- [Grafana LGTM Stack](./otel-lgtm/)

Portal Integration Example
- [Portal Integration](./portal-integration/)

Other examples
- [Elastic Stack](./otel-elastic)

Repositories (used in most of the examples)
- [Repositories](./repositories/)

GemFire Shared State Example
- [GemFire](./gemfire/)

GemFire is licensed Broadcom/VMware software — there is no public registry for it. You need
access to Broadcom's Support Portal registry and must
tell the Makefile which one to use: `GEMFIRE_REGISTRY`, `REGISTRY_USER`, and `REGISTRY_TOKEN` have
no defaults, and `make gemfire-operator`/`gemfire-cluster` fail fast with a named-variable error
if any of them are unset.

**Broadcom Support Portal registry** (`GEMFIRE_CHART_REPO`/`GEMFIRE_CONTROLLER_IMAGE_REPO`/
`GEMFIRE_CLUSTER_IMAGE` default to this registry's path layout, so only the registry and
credentials need to be set):
```
GEMFIRE_REGISTRY=registry.packages.broadcom.com \
REGISTRY_USER=<support-portal-email> \
REGISTRY_TOKEN=<registry-token> \
make gemfire-operator
```
Generate a Registry Token from the Support Portal: log in at [support.broadcom.com](https://support.broadcom.com),
go to **My Downloads**, click **Registry Tokens** (only visible if your account has an active
GemFire/Tanzu entitlement), then **Generate Token**. Your username is the email you sign in with.

Once the registry variables above are set (export them, or prefix every `make` call below), stand
up the example GemFire cluster and bootstrap it in three steps:
```
make gemfire-operator
make gemfire-cluster
make gemfire-configure FUNCTIONS_JAR=/path/to/layer7-gemfire-functions-*.jar SECURITY_JAR=/path/to/layer7-gemfire-security-*.jar
```
- `gemfire-operator` installs VMware's GemFire operator via Helm.
- `gemfire-cluster` applies [`./gemfire/cluster1.yaml`](./gemfire/cluster1.yaml), a 2 locator/2 server `GemFireCluster`.
- `gemfire-configure` runs a one-time `gfsh` bootstrap ([`./gemfire/configure-cluster.sh`](./gemfire/configure-cluster.sh))
  that deploys the Layer7 function and security jars and creates the regions
  (`layer7gw_keyvalue`, `layer7gw_session`, `layer7gw_sortedset`, `layer7gw_ratelimiter`, `layer7gw_counter`) the
  Gateway's GemFire shared-state client expects. `FUNCTIONS_JAR` and `SECURITY_JAR` must point to local paths for
  the two jars — there is no default, and the script fails with a clear error if either is unset or not found.
  Both jars ship together inside `Layer7_API_Gateway_Gemfire_Extension_11.#.#.zip`, available from the
  [Broadcom Support Download Center](https://techdocs.broadcom.com/us/en/ca-enterprise-software/layer7-api-management/api-gateway/11-2/release-notes/list-of-update-files.html)
  (the Gateway release's "List of Update Files" page) — unzip it and point `FUNCTIONS_JAR`/`SECURITY_JAR` at the
  extracted `layer7-gemfire-functions-xxx.jar`/`layer7-gemfire-security-xxx.jar`. See
  [Connect to an External GemFire Datastore](https://techdocs.broadcom.com/us/en/ca-enterprise-software/layer7-api-management/api-gateway/11-2/install-configure-upgrade/connect-to-a-gemfire-datastore/connect-to-an-external-gemfire-datastore.html)
  for full background.

Point a Gateway CR's `spec.app.gemfire` config at `cluster1-locator-clusterip:10334` once this is done.