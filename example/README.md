# Layer7 Operator Examples
These examples cover a broader set of the features that the Layer7 Operator provides and serve as a starting point for implementing these in your own environments.

## Prerequisites
- Kubernetes v1.25+
- Gateway v10/11.x License
- Ingress Controller (You can also expose Gateway Services as L4 LoadBalancers)

#### Getting started
1. Place a gateway v10 or v11 license in [resources/secrets/license](./resources/secrets/license).
2. If you would like to create a TLS secret for your ingress controller then add tls.crt and tls.key to [resources/secrets/tls](./resources/secrets/tls)
    - these will be referenced later on.

#### Examples
All examples use [kustomize](https://kustomize.io/). The basic example covers a simple deployment with a single Gateway 3 Repository CRs (custom resource) configured for static and dynamic updates resepectively with a focus on repositories. Building on the basic example, the advanced example focuses on gateway configuration.

Gateways
- [Basic](./basic)
- [Advanced](./advanced)

Repositories (used in both examples)
- [Repositories](./repositories/)

#### Coming soon
- OTel
  - Prometheus/Granana
  - ECK (Elastic)