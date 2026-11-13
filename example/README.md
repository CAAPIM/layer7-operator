# Layer7 Operator Examples
These examples cover a broader set of the features that the Layer7 Operator provides and serve as a starting point for implementing these in your own environments or just trying them out.

## Prerequisites
- Kubernetes v1.34+
- Gateway v11.x License
- Ingress or Kubernetes Gateway v1 (you can also expose Gateway Services as L4 LoadBalancers)
- Accept the Gateway License
  - license.accept defaults to false in all of the [Gateway examples](./gateway/)
  - update license.accept to true before proceeding

> **Breaking change for users upgrading from v1.2.x or earlier:** The Kubernetes Ingress examples and their `Makefile` quickstart use **[Contour](https://projectcontour.io/)** as the ingress controller (`ingressClassName: contour`, `projectcontour.io/upstream-protocol.tls` annotations). If you were using the previous nginx-based examples, your ingress configuration is **incompatible** and must be updated. Replace `ingressClassName: nginx` with `ingressClassName: contour` and substitute any `nginx.ingress.kubernetes.io/*` annotations with their Contour equivalents. Run `make contour` (or `make contour-kind` for Kind clusters) to install Contour before applying those examples.

The [Gateway API v1 examples](./ingress/gatewayv1/) use [Envoy Gateway](https://gateway.envoyproxy.io/docs/tasks/quickstart/).
Refer to the [Kubernetes Gateway API documentation](https://gateway-api.sigs.k8s.io/guides/)
for setup, then follow the example's [Quickstart](./ingress/gatewayv1/README.md#quickstart)
to deploy with `make kubernetes-gateway-example`.

The basic and advanced examples can be run in a single namespace, The OTel Examples require multiple namespaces for the additional components. Your Kubernetes user or service account must have sufficient privileges to create namespaces, deployments, configmaps, secrets, service accounts, roles, etc..

Each example also includes a [Kind](https://kind.sigs.k8s.io/) (Kubernetes in Docker) Quickstart which you can utilise if you have access to a Docker Machine.

#### Getting started
1. Place a gateway v11 license in [./base/resources/secrets/license](./base/resources/secrets/license).
2. If you would like to create a TLS secret for your ingress controller then add tls.crt and tls.key to [./base/resources/secrets/tls](./base/resources/secrets/tls)
    - these will be referenced later on.

#### Examples
Most examples use [kustomize](https://kustomize.io/). The basic example covers a simple deployment with a single Gateway 3 Repository CRs (custom resource) configured for static and dynamic updates resepectively with a focus on repositories. Building on the basic example, the advanced example focuses on gateway configuration.

Gateways
- [Basic](./basic)
- [Advanced](./advanced)
- [Ingress and Gateway v1](./ingress/)
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
