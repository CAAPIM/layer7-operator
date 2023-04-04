
1. Place a Gateway license as license.xml into the example folder
2. If you would like to create a TLS secret then add tls.crt and tls.key, then uncomment lines 15-19 in example/kustomization.yaml.
3. Update example/security_v1_gateway.yaml with any changes you would like to make (eg. ingress configuration)

The default external traffic exposure method for Operator Managed Gateways is via Kubernetes Ingress Controller. This can be disabled in example/security_v1_gateway.yaml if you'd like to use a L4 Loadbalancer.

```
$ kubectl apply -k example/
serviceaccount/ssg-serviceaccount created
secret/gateway-license created
secret/gateway-secret created
gateway.security.brcmlabs.com/ssg created
gateway.security.brcmlabs.com/sample-graphman-repo created
```

### Uninstall


#### Remove the Gateway Example 
```
$ kubectl delete -k example/
serviceaccount/ssg-serviceaccount deleted
secret "gateway-license" deleted
secret "gateway-secret" deleted
gateway.security.brcmlabs.com "ssg" deleted
```