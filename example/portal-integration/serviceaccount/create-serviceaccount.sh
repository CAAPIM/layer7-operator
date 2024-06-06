#!/bin/bash
namespace=$1

if [[ -z ${namespace} ]]; then
echo "please set namespace"
exit 1
fi

ca=$(kubectl -n ${namespace} get secret/portal-sa-token -o jsonpath='{.data.ca\.crt}')
token=$(kubectl -n ${namespace} get secret/portal-sa-token -o jsonpath='{.data.token}' | base64 --decode)
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
    namespace: ${namespace}
    user: portal-example
current-context: portal-example
users:
- name: portal-example
  user:
    token: ${token}
" >./portal-integration/portal.kubeconfig

echo "###############################################"
echo "portal.kubeconfig written to ./portal-integration/portal.kubeconfig"
echo "###############################################"
echo "###############################################"
echo "############## portal.kubeconfig ##############"
echo "###############################################"
echo "###############################################"
cat ./portal-integration/portal.kubeconfig
