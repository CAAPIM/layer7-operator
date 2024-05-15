# Layer7 Operator Examples
These examples cover a broader set of the features that the Layer7 Operator provides and serve as a starting point for implementing these in your own environments or just trying them out.

## Prerequisites
- Kubernetes v1.25+
- Gateway v10/11.x License
- Ingress Controller (You can also expose Gateway Services as L4 LoadBalancers)
- Accept the Gateway License
  - license.accept defaults to false in all of the [Gateway examples](./gateway/)
  - update license.accept to true before proceeding

The basic and advanced examples can be run in a single namespace, The OTel Examples require multiple namespaces for the additional components. Your Kubernetes user or service account must have sufficient privileges to create namespaces, deployments, configmaps, secrets, service accounts, roles, etc..

Each example also includes a [Kind](https://kind.sigs.k8s.io/) (Kubernetes in Docker) Quickstart which you can utilise if you have access to a Docker Machine.

#### Getting started
1. Place a gateway v10 or v11 license in [./base/resources/secrets/license](./base/resources/secrets/license).
2. If you would like to create a TLS secret for your ingress controller then add tls.crt and tls.key to [./base/resources/secrets/tls](./base/resources/secrets/tls)
    - these will be referenced later on.

#### Examples
All examples use [kustomize](https://kustomize.io/). The basic example covers a simple deployment with a single Gateway 3 Repository CRs (custom resource) configured for static and dynamic updates resepectively with a focus on repositories. Building on the basic example, the advanced example focuses on gateway configuration.

Gateways
- [Basic](./basic)
- [Advanced](./advanced)

Open Telemetry Example
- [Grafana LGTM Stack](./otel-lgtm/)

Other examples
- [Elastic Stack](./otel-elastic)
- [Prometheus](./otel-prometheus)

Repositories (used in both examples)
- [Repositories](./repositories/)