DROP keyspace otk_db;
CREATE KEYSPACE otk_db WITH REPLICATION = { 'class': 'SimpleStrategy',
'replication_factor': 1 };
USE otk_db;
CREATE TABLE oauth_client (
    client_ident text PRIMARY KEY,
    name text,
    type text,
    description text,
    organization text,
    registered_by text,
    registered_by_lower text,
    created bigint,
    custom text,
    software_statement text
);
CREATE TABLE oauth_client_name_v2 (
    name text,
    client_ident text,
    organization text,
    PRIMARY KEY (name, organization, client_ident)
);
CREATE TABLE oauth_client_org_v2 (
    organization text,
    client_ident text,
    PRIMARY KEY (organization, client_ident)
);
CREATE TABLE oauth_client_key (
    client_key text,
    secret text,
    scope text,
    callback text,
    environment text,
    expiration bigint,
    secret_expiration bigint,
    status text,
    created bigint,
    created_by text,
    client_ident text,
    custom text,
    serviceIds text,
    accountPlanMappingIds text,
    master int,
    hashing_metadata text,
    resources text,
    PRIMARY KEY(client_key, client_ident)
);
CREATE TABLE oauth_client_key_user (
    created_by text,
    client_ident text,
    client_key text,
    PRIMARY KEY(created_by, client_ident, client_key)
);
CREATE TABLE oauth_client_key_created_by (
    client_key text,
    secret text,
    scope text,
    callback text,
    environment text,
    expiration bigint,
    secret_expiration bigint,
    status text,
    created bigint,
    created_by text,
    client_ident text,
    custom text,
    serviceIds text,
    accountPlanMappingIds text,
    master int,
    resources text,
    PRIMARY KEY((client_ident, created_by), client_key)
);
CREATE INDEX oauth_client_registered_by_lower_idx ON oauth_client (registered_by_lower);
CREATE INDEX ON oauth_client_key (client_ident);
CREATE TABLE oauth_client_key_view_app_instances(
    client_key text,
    client_ident text,
    master int,
    PRIMARY KEY(client_ident, master, client_key)
);
CREATE TABLE oauth_access_token (
    otk_token text PRIMARY KEY,
    secret text,
    expiration bigint,
    scope text,
    resource_owner text,
    created bigint,
    otoken text,
    status text,
    client_key text,
    client_name text,
    client_ident text,
    custom text,
    resources text,
    token_rank int
);
CREATE INDEX ON oauth_access_token (client_ident);
CREATE TABLE oauth_access_token_view_client_key (
    otk_token text,
    resource_owner text,
    client_key text,
    client_name text,
    PRIMARY KEY (client_key, resource_owner, otk_token)
);
CREATE TABLE oauth_access_token_view_resource_owner (
    otk_token text,
    resource_owner text,
    client_key text,
    status text,
    client_name text,
    PRIMARY KEY (resource_owner, client_key, otk_token)
);
CREATE TABLE oauth_access_token_view_client_ident (
    otk_token text,
    client_ident text,
    resource_owner text,
    token_rank int,
    created bigint,
    PRIMARY KEY(
        (client_ident, resource_owner),
        token_rank,
        created,
        otk_token
    )
);
CREATE TABLE oauth_access_token_view_app_active_users (
    otk_token text,
    client_ident text,
    PRIMARY KEY(client_ident, otk_token)
);
CREATE INDEX ON oauth_access_token_view_resource_owner (status);
CREATE TABLE oauth_access_token_pagination (
    otk_token text,
    client_name_otk_token text,
    secret text,
    expiration bigint,
    scope text,
    resource_owner text,
    created bigint,
    otoken text,
    status text,
    client_key text,
    client_name text,
    custom text,
    resources text,
    PRIMARY KEY (client_name_otk_token, resource_owner, otk_token)
);
CREATE INDEX ON oauth_access_token_pagination (status);
CREATE INDEX ON oauth_access_token_pagination (client_key);
CREATE TABLE oauth_access_token_by_resource_owner (
    resource_owner text,
    otk_token text,
    client_key text,
    client_name text,
    created bigint,
    custom text,
    expiration bigint,
    oexpiration bigint,
    otoken text,
    scope text,
    secret text,
    status text,
    resources text,
    PRIMARY KEY (resource_owner, client_key, otk_token)
);
CREATE TABLE oauth_access_token_by_resource_owner_pagination (
    client_name_otk_token text,
    resource_owner text,
    otk_token text,
    client_key text,
    client_name text,
    created bigint,
    custom text,
    expiration bigint,
    oexpiration bigint,
    otoken text,
    scope text,
    secret text,
    status text,
    resources text,
    PRIMARY KEY (client_name_otk_token, resource_owner)
);
CREATE TABLE oauth_refresh_token (
    otk_token text PRIMARY KEY,
    secret text,
    expiration bigint,
    scope text,
    resource_owner text,
    created bigint,
    otoken text,
    // other token = access token oexpiration bigint,
    status text,
    client_key text,
    client_name text,
    custom text,
    resources text
);
CREATE TABLE oauth_refresh_token_view_client_key (
    otk_token text,
    resource_owner text,
    client_key text,
    client_name text,
    PRIMARY KEY (client_key, resource_owner, otk_token)
);
CREATE TABLE oauth_refresh_token_view_resource_owner (
    otk_token text,
    resource_owner text,
    client_key text,
    status text,
    client_name text,
    PRIMARY KEY (resource_owner, client_key, otk_token)
);
CREATE TABLE oauth_initiate (
    otk_token text PRIMARY KEY,
    secret text,
    expiration bigint,
    scope text,
    resource_owner text,
    created bigint,
    verifier text,
    callback text,
    client_key text,
    client_name text,
    custom text,
    resources text
);
CREATE TABLE oauth_initiate_view_verifier (
    otk_token text,
    verifier text,
    PRIMARY KEY (otk_token, verifier)
);
CREATE TABLE oauth_initiate_view_client_key (
    otk_token text,
    client_key text,
    client_name text,
    resource_owner text,
    PRIMARY KEY (client_key, otk_token)
);
CREATE TABLE oauth_initiate_view_resource_owner (
    otk_token text,
    resource_owner text,
    client_key text,
    PRIMARY KEY (resource_owner, client_key, otk_token)
);
CREATE TABLE oauth_session (
    session_key text,
    session_group text,
    expiration bigint,
    value text,
    PRIMARY KEY(session_key, session_group)
);

CREATE TABLE otk_version (
current_version text PRIMARY KEY
);
INSERT INTO otk_version (current_version)
VALUES ('otk4.6.3');
CREATE TABLE oauth_id_token (
    resource_owner text,
    azp text,
    sub text,
    jwt_id text,
    jwt text,
    salt text,
    shared_secret text,
    shared_secret_type text,
    iss text,
    expiration bigint,
    PRIMARY KEY(resource_owner, azp)
);
CREATE TABLE oauth_id_token_view_sub (
    resource_owner text,
    azp text,
    sub text,
    jwt_id text,
    jwt text,
    salt text,
    shared_secret text,
    shared_secret_type text,
    iss text,
    expiration bigint,
    PRIMARY KEY(sub, azp)
);
CREATE TABLE oauth_id_token_view_azp (
    resource_owner text,
    azp text,
    sub text,
    jwt_id text,
    jwt text,
    salt text,
    shared_secret text,
    shared_secret_type text,
    iss text,
    expiration bigint,
    PRIMARY KEY(azp)
);
CREATE TABLE oauth_id_token_view_jwt_id (
    resource_owner text,
    azp text,
    sub text,
    jwt_id text,
    jwt text,
    salt text,
    shared_secret text,
    shared_secret_type text,
    iss text,
    expiration bigint,
    PRIMARY KEY(jwt_id)
);
CREATE TABLE oauth_id_token_view_jwt (
    resource_owner text,
    azp text,
    sub text,
    jwt_id text,
    jwt text,
    salt text,
    shared_secret text,
    shared_secret_type text,
    iss text,
    expiration bigint,
    PRIMARY KEY(jwt)
);
CREATE TABLE IF NOT EXISTS portal_apikey_by_apikey_pk (
    apikey_pk text PRIMARY KEY,
    apikey text,
    apikey_secret text,
    apis text,
    created bigint,
    created_by text,
    label text,
    modified_by text,
    organization text,
    organization_id text,
    status text,
    updated bigint,
    value_xml text
);
CREATE TABLE IF NOT EXISTS portal_apikey_by_apikey (
    apikey text PRIMARY KEY,
    apikey_pk text,
    apikey_secret text,
    apis text,
    created bigint,
    created_by text,
    label text,
    modified_by text,
    organization text,
    organization_id text,
    hashing_metadata text,
    status text,
    updated bigint,
    value_xml text
);
CREATE TABLE IF NOT EXISTS portal_apikey_by_org_and_apikey_pk (
    apikey_pk text,
    apikey text,
    apikey_secret text,
    apis text,
    created bigint,
    created_by text,
    label text,
    modified_by text,
    organization text,
    organization_id text,
    status text,
    updated bigint,
    value_xml text,
    PRIMARY KEY (organization_id, apikey_pk)
);
CREATE TABLE oauth_token_revocation_list (otk_token text PRIMARY KEY);
CREATE TABLE oauth_device_initiate (
    device_code text PRIMARY KEY,
    user_code text,
    scope text,
    resources text,
    client_key text,
    client_name text,
    resource_owner text,
    status text,
    expiration bigint,
    created bigint,
    last_polled_time bigint,
    custom text
);
CREATE TABLE oauth_device_initiate_by_user_code (
    user_code text PRIMARY KEY,
    device_code text,
    scope text,
    resources text,
    client_key text,
    client_name text,
    resource_owner text,
    status text,
    expiration bigint,
    created bigint,
    last_polled_time bigint,
    custom text
);
CREATE TABLE oauth_par_initiate (
    request_uri text PRIMARY KEY,
    request_value text,
    expiration bigint,
    created bigint,
    custom text
);
CREATE TABLE oauth_ciba_initiate (
    auth_req_id text PRIMARY KEY,
    user_info text,
    scope text,
    resources text,
    client_key text,
    client_name text,
    status text,
    expiration bigint,
    created bigint,
    last_polled_time bigint,
    auth_time bigint,
    acr_values text,
    custom text
);