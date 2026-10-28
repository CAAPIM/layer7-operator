# GemFire Shared State Example
This example stands up a standalone GemFire cluster that a Gateway can use as a shared-state backend,
either as a general-purpose shared state client or as an alternative to Redis for API key storage in
the [Portal Integration example](../portal-integration/).

GemFire is licensed Broadcom/VMware software — there is no public registry for it. You need
access to Broadcom's Support Portal registry and must
tell the Makefile which one to use: `GEMFIRE_REGISTRY`, `REGISTRY_USER`, and `REGISTRY_TOKEN` have
no defaults, and `make gemfire-operator`/`gemfire-cluster` fail fast with a named-variable error
if any of them are unset.

**Broadcom Support Portal registry** (`GEMFIRE_CHART_REPO`/`GEMFIRE_CONTROLLER_IMAGE_REPO`/
`GEMFIRE_CLUSTER_IMAGE` default to this registry's path layout, so only the registry and
credentials need to be set):
```
GEMFIRE_REGISTRY=registry.packages.broadcom.com \
REGISTRY_USER=<support-portal-email> \
REGISTRY_TOKEN=<registry-token> \
make gemfire-operator
```
Generate a Registry Token from the Support Portal: log in at [support.broadcom.com](https://support.broadcom.com),
go to **My Downloads**, click **Registry Tokens** (only visible if your account has an active
GemFire/Tanzu entitlement), then **Generate Token**. Your username is the email you sign in with.

Once the registry variables above are set (export them, or prefix every `make` call below), stand
up the example GemFire cluster and bootstrap it in three steps:
```
make gemfire-operator
make gemfire-cluster
make gemfire-configure FUNCTIONS_JAR=/path/to/layer7-gemfire-functions-*.jar SECURITY_JAR=/path/to/layer7-gemfire-security-*.jar
```
- `gemfire-operator` installs VMware's GemFire operator via Helm.
- `gemfire-cluster` applies [`./cluster1.yaml`](./cluster1.yaml), a 2 locator/2 server `GemFireCluster`.
- `gemfire-configure` runs a one-time `gfsh` bootstrap ([`./configure-cluster.sh`](./configure-cluster.sh))
  that deploys the Layer7 function and security jars and creates the regions
  (`layer7gw_keyvalue`, `layer7gw_session`, `layer7gw_sortedset`, `layer7gw_ratelimiter`, `layer7gw_counter`) the
  Gateway's GemFire shared-state client expects. `FUNCTIONS_JAR` and `SECURITY_JAR` must point to local paths for
  the two jars — there is no default, and the script fails with a clear error if either is unset or not found.
  Both jars ship together inside `Layer7_API_Gateway_Gemfire_Extension_11.#.#.zip`, available from the
  [Broadcom Support Download Center](https://techdocs.broadcom.com/us/en/ca-enterprise-software/layer7-api-management/api-gateway/11-2/release-notes/list-of-update-files.html)
  (the Gateway release's "List of Update Files" page) — unzip it and point `FUNCTIONS_JAR`/`SECURITY_JAR` at the
  extracted `layer7-gemfire-functions-xxx.jar`/`layer7-gemfire-security-xxx.jar`. See
  [Connect to an External GemFire Datastore](https://techdocs.broadcom.com/us/en/ca-enterprise-software/layer7-api-management/api-gateway/11-2/install-configure-upgrade/connect-to-a-gemfire-datastore/connect-to-an-external-gemfire-datastore.html)
  for full background.

## Configuring the Gateway CR
Once the cluster is bootstrapped, point a Gateway CR's `spec.app.gemfire` config at
`cluster1-locator-clusterip:10334`. The minimal, no-auth/no-TLS configuration looks like this
(matches [`../gateway/portal-gateway.yaml`](../gateway/portal-gateway.yaml)):
```yaml
spec:
  app:
    gemfire:
      enabled: true
      testOnStart: true
      locators:
      - host: cluster1-locator-clusterip
        port: 10334
      auth:
        enabled: false
      ssl:
        enabled: false
```

This cluster can also be used as an alternative backend for API key storage in the
[Portal Integration example](../portal-integration/) — Redis remains required there for the
Portal-specific state store, but GemFire can optionally replace Redis as the store for API keys.
This is decided when configuring the Gateway CR, see
[Using GemFire instead of Redis for API key storage](../portal-integration/readme.md#using-gemfire-instead-of-redis-for-api-key-storage)
for details.

### Full field reference (`spec.app.gemfire`)
- `enabled` — turns the GemFire shared-state client on.
- `testOnStart` — verifies connectivity to the locators when the Gateway node starts.
- `locators` — list of `{host, port}` pairs, e.g. `cluster1-locator-clusterip:10334`.
- `existingSecret` — mounts an existing secret containing GemFire configuration instead of letting the
  operator manage one. The secret must contain a key called `sharedstate_client.yaml`. **If Redis is
  also enabled on the same Gateway, this must reference the same secret as `redis.existingSecret`**
  (or both must be left empty for the operator to manage it) — the two backends share one
  `sharedstate_client.yaml`.
- `certs` — `[]{enabled, secretName, key}`; `key` must match the file referenced in `existingSecret`.
- Region name overrides — only needed if you changed the region names in
  [`configure-cluster.sh`](./configure-cluster.sh) from their defaults:
  - `gwKeyValueRegionName` (default `layer7gw_keyvalue`)
  - `gwCounterRegionName` (default `layer7gw_counter`)
  - `gwRateLimiterRegionName` (default `layer7gw_ratelimiter`)
  - `gwSortedSetRegionName` (default `layer7gw_sortedset`)
- `auth` — `{enabled, username, passwordEncoded, passwordPlaintext}`.
- `dynamicProperties` — a `map[string]string` of additional GemFire client properties.
- `ssl` — `{enabled, enabledComponents (defaults to all), keystoreType, truststoreType, keystore,
  truststore}`. `keystoreType`/`truststoreType` default to `JKS`, and also support `PKCS12`.
  `keystore`/`truststore` are each `{existingSecretName, existingSecretKey, passwordEncoded,
  passwordPlaintext}` — `existingSecretKey` defaults to `keystore.jks`/`truststore.jks` respectively.

### Enabling auth
```yaml
    gemfire:
      enabled: true
      testOnStart: true
      locators:
      - host: cluster1-locator-clusterip
        port: 10334
      auth:
        enabled: true
        username: gemfireadmin
        passwordPlaintext: 7layer
      ssl:
        enabled: false
```
This assumes the cluster's management/peer security is backed by a `gemfireadmin`/`7layer` credential
pair (e.g. via a `mgmtSvcCredentialsSecretName` secret on the `GemFireCluster`) — adjust to match
whatever credentials your cluster is configured with.

### Enabling SSL
```yaml
    gemfire:
      enabled: true
      testOnStart: true
      locators:
      - host: cluster1-locator-clusterip
        port: 10334
      auth:
        enabled: false
      ssl:
        enabled: true
        keystoreType: PKCS12
        truststoreType: PKCS12
        keystore:
          existingSecretName: cluster1-cert
          existingSecretKey: keystore.p12
          passwordPlaintext: <cluster1-cert-password>
        truststore:
          existingSecretName: cluster1-cert
          existingSecretKey: truststore.p12
          passwordPlaintext: <cluster1-cert-password>
```
`<cluster1-cert-password>` is the auto-generated password for the cluster's TLS cert secret, retrieved with:
```
kubectl get secret cluster1-cert -o jsonpath='{.data.password}' | base64 -d
```
Auth and SSL can be combined (both blocks enabled at once).
