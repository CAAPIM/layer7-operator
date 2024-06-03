#!/bin/bash
ca=$(kubectl -n default get secret/portal-sa-token -o jsonpath='{.data.ca\.crt}')
token=$(kubectl -n default get secret/portal-sa-token -o jsonpath='{.data.token}' | base64 --decode)
controlPlane="https://kubernetes.default.svc.cluster.local"

echo "
apiVersion: v1
kind: Config
clusters:
- name: portal-example
  cluster:
    certificate-authority-data: ${ca}
    server: ${controlPlane}
contexts:
- name: portal-example
  context:
    cluster: portal-example
    namespace: default
    user: portal-example
current-context: portal-example
users:
- name: portal-example
  user:
    token: ${token}
" >./portal-integration/portal.kubeconfig

