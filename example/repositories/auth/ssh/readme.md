# SSH Authentication with Git
This folder contains an example structure of what's required to use SSH Authentication with a Git repository.

## Important
known hosts are written to a temporary file in this integration and are shared between repositories

## Usage
While credentials can be provided in [plaintext](./apis-repository-plaintext-auth.yaml) it is recommended that you use Kubernetes Secrets.

Your OpenSSH Key can be in plaintext or encrypted format.
The default in [kustomization.yaml](./kustomization.yaml) is encrypted where a [passphrase](./sshpass.txt) is required.

## Using this folder as a starting point
1. Add your private key contents to [ssh_key_encrypted.key](./ssh_key_encrypted.key)
2. Add your private key passphrase to [sshpass.txt](./sshpass.txt)
Your key can be encrypted with the following command
```
ssh-keygen -p -f ssh_key_encrypted.key
```

3. Update [known_hosts](./known_hosts)
The following defaults have been added to this example
```
github.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOMqqnkVzrm0SdG6UOoqKLsabgH5C9okWi0dh2l9GKJl
gitlab.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIAfuCHKVTjquxvt6CM6tdG4SLp1Btn/nOeHHE5UOzRdf
bitbucket.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIIazEu89wgQZ4bqs3d63QSMzYVa0MuJ2e2gKTKqu+UUO
```

4. Update the endpoint, names and branch in [apis-repository.yaml](./apis-repository.yaml)
```
apiVersion: security.brcmlabs.com/v1
kind: Repository
metadata:
  name: l7-gw-myapis
spec:
  name: l7-gw-myapis
  enabled: true
  endpoint: ssh://git@[yourgitserver.com]/[username]/[reponame]
  #endpoint: ssh://git@[yourgitserver.com]:[port]/[username]/[reponame]
  branch: main
  auth:
    type: ssh
    existingSecretName: myapis-sshkey
```
5. Apply using Kustomize
```
kubectl apply -k ./example/repositories/auth/ssh
```

## Troubleshooting
The Layer7 Operator will log an error if there is an issue with your key, key passphrase or known_hosts

You can view the operator log with the following command
```
kubectl logs -f -l app.kubernetes.io/name=layer7-operator
```