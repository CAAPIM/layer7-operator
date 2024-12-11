# Basic Authentication with Git
This folder contains an example structure of what's required to use Basic Authentication with a Git repository.

## Usage
While credentials can be provided in [plaintext](./apis-repository-plaintext-auth.yaml) it is recommended that you use Kubernetes Secrets.

The default in [kustomization.yaml](./kustomization.yaml) uses an existing secret that is created from [repository-secret.env](./repository-secret.env).

## Using this folder as a starting point
1. Update [repository-secret.env](./repository-secret.env) with your git username and personal access token
2. Update the endpoint, names and branch in [apis-repository.yaml](./apis-repository.yaml)
```
apiVersion: security.brcmlabs.com/v1
kind: Repository
metadata:
  name: l7-gw-myapis
spec:
  name: l7-gw-myapis
  enabled: true
  endpoint: https://[yourgitserver.com]/[username]/[reponame]
  branch: main
  auth:
    type: basic
    existingSecretName: myapis-basic-auth
```
3. Apply using Kustomize
```
kubectl apply -k ./example/repositories/auth/basic
```

## Troubleshooting
The Layer7 Operator will log an error if there is an issue with credentials.

You can view the operator log with the following command
```
kubectl logs -f -l app.kubernetes.io/name=layer7-operator
```