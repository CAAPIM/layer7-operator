# Ingress examples

These examples expose the Layer7 Gateway through an ingress or Gateway API controller.

### Getting started

Follow the [example prerequisites](../README.md#prerequisites), then choose a guide
below. Each guide includes a Quickstart and the individual deployment steps.

### Examples

- [Gateway API v1](./gatewayv1/): Envoy Gateway with HTTPS and backend TLS, an existing
  shared Gateway, or TLS passthrough. Uses `spec.app.kubernetesGateway` and the
  `kubernetes-gateway-example` Makefile target.
- [Ingress in the Advanced example](../advanced/): Contour with Kubernetes Ingress.
  Uses `spec.app.ingress` and the `advanced` Makefile target.

Ingress and Gateway API resources can be enabled together. Each has its own
configuration and resource lifecycle. Test the Gateway API endpoint before moving
DNS, then disable Ingress when the migration is complete.

The guides live here; Layer7 Gateway manifests remain in [example/gateway](../gateway/)
alongside the other Gateway examples. Run Makefile commands from `example/`.
