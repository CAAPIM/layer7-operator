# Dual Gateway OTK Configuration Guide

This guide describes how to configure a Dual Gateway OAuth Toolkit (OTK) deployment using the Layer7 Gateway Operator. In a dual gateway setup, one gateway acts as the DMZ (Demilitarized Zone) gateway and another acts as the Internal gateway, providing enhanced security and separation of concerns.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Configuration Overview](#configuration-overview)
- [Step 1: Create Required Secrets](#step-1-create-required-secrets)
- [Step 2: Configure DMZ Gateway](#step-2-configure-dmz-gateway)
- [Step 3: Configure Internal Gateway](#step-3-configure-internal-gateway)
- [Key Configuration Fields](#key-configuration-fields)
- [Deployment](#deployment)
- [Certificate Synchronization](#certificate-synchronization)
- [External Gateway Support](#external-gateway-support)
- [Troubleshooting](#troubleshooting)

## Overview

The Dual Gateway OTK deployment consists of:

- **DMZ Gateway**: Handles external client requests and acts as the OAuth authorization server
- **Internal Gateway**: Handles token validation and resource server operations

The operator automatically synchronizes certificates and keys between the two gateways, ensuring secure communication and proper OAuth flow.

## Configuration Overview

The dual gateway setup requires:

1. **TLS Secrets**: For DMZ and Internal gateway keys/certificates
2. **Auth Secrets**: For gateway authentication credentials
3. **DMZ Gateway Configuration**: With `type: dmz`
4. **Internal Gateway Configuration**: With `type: internal`

## Step 1: Create Required Secrets

### Create TLS Secrets

You need to create TLS secrets for both DMZ and Internal gateways. These secrets must be of type `kubernetes.io/tls` and contain:
- `tls.crt`: The certificate
- `tls.key`: The private key

#### Option 1: Using the Provided Script

A helper script is available to generate self-signed certificates and create all required secrets:

```bash
cd example/gateway/otk/secrets
./create-secrets.sh <namespace>
```

This script creates:
- `otk-dmz-tls-secret`: TLS secret for DMZ gateway
- `otk-internal-tls-secret`: TLS secret for Internal gateway
- `otk-dmz-auth-secret`: Authentication secret for DMZ gateway (username: `admin`, password: `7layer`)
- `otk-internal-auth-secret`: Authentication secret for Internal gateway (username: `admin`, password: `7layer`)

#### Option 2: Manual Secret Creation

**DMZ TLS Secret:**

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: otk-dmz-tls-secret
  namespace: default
type: kubernetes.io/tls
data:
  tls.crt: <base64-encoded-certificate>
  tls.key: <base64-encoded-private-key>
```

**Internal TLS Secret:**

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: otk-internal-tls-secret
  namespace: default
type: kubernetes.io/tls
data:
  tls.crt: <base64-encoded-certificate>
  tls.key: <base64-encoded-private-key>
```

**DMZ Auth Secret:**

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: otk-dmz-auth-secret
  namespace: default
type: Opaque
stringData:
  SSG_ADMIN_USERNAME: admin
  SSG_ADMIN_PASSWORD: 7layer
```

**Internal Auth Secret:**

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: otk-internal-auth-secret
  namespace: default
type: Opaque
stringData:
  SSG_ADMIN_USERNAME: admin
  SSG_ADMIN_PASSWORD: 7layer
```

## Step 2: Configure DMZ Gateway

The DMZ gateway configuration should include:

- `otk.type: dmz`
- Reference to Internal gateway
- DMZ key secret reference
- Internal auth secret for communication with Internal gateway
- Database configuration

### Sample DMZ Gateway Configuration

```yaml
apiVersion: security.brcmlabs.com/v1
kind: Gateway
metadata:
  name: otk-ssg-dmz
spec:
  version: "11.1.3"
  license:
    accept: true
    secretName: gateway-license
  app:
    replicas: 1
    image: docker.io/caapim/gateway:11.1.3
    imagePullPolicy: IfNotPresent
    resources:
      requests:
        memory: 8Gi
        cpu: 3
      limits:
        memory: 8Gi
        cpu: 3
    # ExternalKeys with otk flag set to true for OTK-specific key usage
    externalKeys:
    - name: otk-dmz-tls-secret
      enabled: true
      alias: otk-dmz-key
      keyUsageType: SSL
      otk: true
    otk:
      enabled: true
      initContainerImage: docker.io/caapim/otk-install:4.6.4
      type: dmz
      # Reference to Internal gateway (can be Gateway name or external hostname)
      internalGatewayReference: otk-ssg-internal
      # InternalGatewayPort is used when the Internal gateway is external
      # If not specified, defaults to 9443 or the gateway's graphmanDynamicSync port
      internalGatewayPort: 9443
      # SyncIntervalSeconds determines how often certificates are synchronized
      # Defaults to RuntimeSyncIntervalSeconds if not specified, or 10 seconds if neither is set
      syncIntervalSeconds: 30
      # Reference to the TLS secret for DMZ key
      dmzKeySecret: otk-dmz-tls-secret
      # Auth secret for Internal gateway communication
      internalAuthSecret: otk-internal-auth-secret
      database:
        type: mysql
        create: true
        connectionName: OAuth
        auth:
          gateway:
            username: otk_user
            password: otkUserPass
          readOnly:
            username: readonly_user
            password: readonly_userPass
          admin:
            username: admin
            password: adminPass
        properties:
          minimumPoolSize: 3
          maximumPoolSize: 15
        sql:
          databaseName: otk_db
          jdbcUrl: jdbc:mysql://mysql.brcmlabs.com:3306/otk_db_init
          jdbcDriverClass: com.mysql.cj.jdbc.Driver
          connectionProperties:
            c3p0.maxConnectionAge: "100"
            c3p0.maxIdleTime: "1000"
          manageSchema: true
          databaseWaitTimeout: 60
    management:
      secretName: gateway-secret
      graphman:
        enabled: true
        initContainerImage: docker.io/caapim/graphman-static-init:1.0.4
        dynamicSyncPort: 9443
      cluster:
        hostname: gateway.brcmlabs.com
    service:
      type: ClusterIP
      ports:
      - name: https
        port: 8443
        targetPort: 8443
        protocol: TCP
```

## Step 3: Configure Internal Gateway

The Internal gateway configuration should include:

- `otk.type: internal`
- Reference to DMZ gateway
- Internal key secret reference
- DMZ auth secret for communication with DMZ gateway
- Database configuration (shared with DMZ)

### Sample Internal Gateway Configuration

```yaml
apiVersion: security.brcmlabs.com/v1
kind: Gateway
metadata:
  name: otk-ssg-internal
spec:
  version: "11.1.3"
  license:
    accept: true
    secretName: gateway-license
  app:
    replicas: 1
    image: docker.io/caapim/gateway:11.1.3
    imagePullPolicy: IfNotPresent
    resources:
      requests:
        memory: 8Gi
        cpu: 3
      limits:
        memory: 8Gi
        cpu: 3
    # ExternalKeys with otk flag set to true for OTK-specific key usage
    externalKeys:
    - name: otk-internal-tls-secret
      enabled: true
      alias: otk-internal-key
      keyUsageType: SSL
      otk: true
    otk:
      enabled: true
      initContainerImage: docker.io/caapim/otk-install:4.6.4
      type: internal
      # Reference to DMZ gateway (can be Gateway name or external hostname)
      dmzGatewayReference: otk-ssg-dmz
      # DmzGatewayPort is used when the DMZ gateway is external
      # If not specified, defaults to 9443 or the gateway's graphmanDynamicSync port
      dmzGatewayPort: 9443
      # SyncIntervalSeconds determines how often certificates are synchronized
      # Defaults to RuntimeSyncIntervalSeconds if not specified, or 10 seconds if neither is set
      syncIntervalSeconds: 30
      # Reference to the TLS secret for Internal key
      internalKeySecret: otk-internal-tls-secret
      # Auth secret for DMZ gateway communication
      dmzAuthSecret: otk-dmz-auth-secret
      database:
        type: mysql
        create: true
        connectionName: OAuth
        auth:
          gateway:
            username: otk_user
            password: otkUserPass
          readOnly:
            username: readonly_user
            password: readonly_userPass
          admin:
            username: admin
            password: adminPass
        properties:
          minimumPoolSize: 3
          maximumPoolSize: 15
        sql:
          databaseName: otk_db
          jdbcUrl: jdbc:mysql://mysql.brcmlabs.com:3306/otk_db_init
          jdbcDriverClass: com.mysql.cj.jdbc.Driver
          connectionProperties:
            c3p0.maxConnectionAge: "100"
            c3p0.maxIdleTime: "1000"
          manageSchema: true
          databaseWaitTimeout: 60
    management:
      secretName: gateway-secret
      graphman:
        enabled: true
        initContainerImage: docker.io/caapim/graphman-static-init:1.0.4
      cluster:
        hostname: gateway.brcmlabs.com
    service:
      type: ClusterIP
      ports:
      - name: https
        port: 8443
        targetPort: 8443
        protocol: TCP
      - name: management
        port: 9443
        targetPort: 9443
        protocol: TCP
```

## Key Configuration Fields

### OTK-Specific Fields

| Field | Description | Required | Default |
|-------|-------------|----------|---------|
| `otk.enabled` | Enable OTK installation | Yes | `false` |
| `otk.type` | OTK type: `dmz`, `internal`, or `single` | Yes | - |
| `otk.initContainerImage` | OTK init container image | Yes | - |
| `otk.dmzKeySecret` | Reference to TLS secret containing DMZ key/cert | Yes (DMZ) | - |
| `otk.internalKeySecret` | Reference to TLS secret containing Internal key/cert | Yes (Internal) | - |
| `otk.dmzAuthSecret` | Reference to secret with DMZ gateway credentials | Yes (Internal) | - |
| `otk.internalAuthSecret` | Reference to secret with Internal gateway credentials | Yes (DMZ) | - |
| `otk.dmzGatewayReference` | Reference to DMZ gateway (name or hostname) | Yes (Internal) | - |
| `otk.internalGatewayReference` | Reference to Internal gateway (name or hostname) | Yes (DMZ) | - |
| `otk.dmzGatewayPort` | Port for DMZ gateway (when external) | No | `9443` or `graphmanDynamicSync` port |
| `otk.internalGatewayPort` | Port for Internal gateway (when external) | No | `9443` or `graphmanDynamicSync` port |
| `otk.syncIntervalSeconds` | Certificate sync interval in seconds | No | `RuntimeSyncIntervalSeconds` or `10` |
| `otk.port` | OTK port (defaults to 8443) | No | `8443` |

### External Keys Configuration

Both gateways must have `externalKeys` configured with the `otk: true` flag:

```yaml
externalKeys:
- name: otk-dmz-tls-secret  # or otk-internal-tls-secret
  enabled: true
  alias: otk-dmz-key        # or otk-internal-key
  keyUsageType: SSL
  otk: true                  # Required for OTK key handling
```

## Deployment

### 1. Create Secrets

```bash
cd example/gateway/otk/secrets
./create-secrets.sh default
```

### 2. Deploy DMZ Gateway

```bash
kubectl apply -f example/gateway/otk/otk-ssg-dmz.yaml
```

### 3. Deploy Internal Gateway

```bash
kubectl apply -f example/gateway/otk/otk-ssg-internal.yaml
```

### 4. Verify Deployment

```bash
# Check gateway pods
kubectl get pods -l app=gateway

# Check gateway status
kubectl get gateway otk-ssg-dmz
kubectl get gateway otk-ssg-internal

# Check logs
kubectl logs -l app=gateway,gateway-name=otk-ssg-dmz
kubectl logs -l app=gateway,gateway-name=otk-ssg-internal
```

## Certificate Synchronization

The operator automatically synchronizes certificates between DMZ and Internal gateways:

1. **DMZ Certificate → Internal Gateway**: When the DMZ certificate is updated, it's automatically published to the Internal gateway as a trusted certificate and used for FIP (Federated Identity Provider) user creation.

2. **Internal Certificate → DMZ Gateway**: When the Internal certificate is updated, it's automatically published to the DMZ gateway as a trusted certificate.

3. **Sync Interval**: Controlled by `syncIntervalSeconds` (default: 10 seconds or `RuntimeSyncIntervalSeconds`).

4. **Key Updates**: When DMZ or Internal keys are updated:
   - The key is synchronized to the respective gateway
   - The DMZ private key name is updated in the cluster-wide property `otk.dmz.private_key.name` (DMZ gateway only)
   - Old certificates are removed before new ones are published

## External Gateway Support

The operator supports scenarios where one or both gateways are external (not managed by the operator):

### External DMZ Gateway

If the DMZ gateway is external, configure the Internal gateway with:

```yaml
otk:
  type: internal
  dmzGatewayReference: external-dmz-gateway.example.com
  dmzGatewayPort: 9443  # Port for Graphman API
  dmzAuthSecret: otk-dmz-auth-secret
```

### External Internal Gateway

If the Internal gateway is external, configure the DMZ gateway with:

```yaml
otk:
  type: dmz
  internalGatewayReference: external-internal-gateway.example.com
  internalGatewayPort: 9443  # Port for Graphman API
  internalAuthSecret: otk-internal-auth-secret
```

### External Gateway Requirements

- Graphman API must be enabled and accessible
- Authentication credentials must be provided via auth secrets
- The correct port must be specified if different from default (9443)
- The gateway must be reachable from the operator's network

## Troubleshooting

### Common Issues

1. **Certificates not synchronizing**
   - Verify `syncIntervalSeconds` is set appropriately
   - Check that Graphman is enabled on both gateways
   - Verify auth secrets are correctly configured
   - Check operator logs for errors

2. **Gateway communication failures**
   - Verify gateway references are correct (name or hostname)
   - Check network connectivity between gateways
   - Verify ports are correctly configured
   - Ensure auth secrets contain valid credentials

3. **Key update failures**
   - Verify TLS secrets are of type `kubernetes.io/tls`
   - Check that secrets contain both `tls.crt` and `tls.key`
   - Ensure `externalKeys` have `otk: true` flag
   - Verify alias matches the expected value

4. **Database connection issues**
   - Verify database credentials in `otk.database.auth`
   - Check JDBC URL is correct and accessible
   - Ensure database exists or `create: true` is set
   - Verify database wait timeout is sufficient

### Checking Certificate Sync Status

```bash
# Check annotations for certificate thumbprints
kubectl get gateway otk-ssg-dmz -o jsonpath='{.metadata.annotations}'
kubectl get gateway otk-ssg-internal -o jsonpath='{.metadata.annotations}'

# Check cluster-wide properties
kubectl exec -it <dmz-pod> -- /opt/SecureSpan/Gateway/node/default/bin/ssgconfig \
  get cluster-wide-properties | grep otk.dmz.private_key.name
```

### Operator Logs

```bash
# View operator logs
kubectl logs -n <operator-namespace> -l control-plane=controller-manager

# Filter for OTK-related logs
kubectl logs -n <operator-namespace> -l control-plane=controller-manager | grep -i otk
```

## Additional Resources

- [Layer7 Gateway Operator Documentation](https://github.com/broadcom/layer7-operator)
- [OAuth Toolkit Documentation](https://techdocs.broadcom.com/us/en/ca-enterprise-software/layer7-api-management/api-gateway/11-1.html)
- Example configurations: `example/gateway/otk/`

---

**Note**: This configuration guide assumes you have a working Kubernetes cluster with the Layer7 Gateway Operator installed. Adjust namespaces, hostnames, and other values according to your environment.

