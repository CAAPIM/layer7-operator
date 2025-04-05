-- Broadcom Inc.
-- Database schema for OTK
--
-- This script is used to create OAuth default test clients for testing purposes in non-production environments.
--
-- This script inserts test data. It is not necessary but recommended to install this script in non-production environments.
-- It is provided for testing purposes only.
--
-- NOTE: this script will first remove all existing test data!
--
-- run this command on command line to create the Cassandra schema: cqlsh -k <keyspace> -f <this-file-name>.cql
--
--
-- truncate all the data
--
USE otk_db;

truncate oauth_client;
truncate oauth_client_name_v2;
truncate oauth_client_org_v2;
truncate oauth_client_key;
truncate oauth_client_key_user;
truncate oauth_client_key_created_by;
--
-- OpenID Connect Client for the Basic Client Profile specification
--
INSERT INTO oauth_client (client_ident, name, description, organization, registered_by, type, created, custom)
VALUES ('123456800-otk', 'OpenID Connect Basic Client Profile', 'Test for OpenID Connect BCP', 'CA Technologies Inc.', 'admin', 'confidential', 0, '{}');
--
INSERT INTO oauth_client_name_v2 (client_ident, name, organization)
VALUES ('123456800-otk', 'openid connect basic client profile', 'ca technologies inc.');
--
INSERT INTO oauth_client_org_v2 (client_ident, organization)
VALUES ('123456800-otk', 'ca technologies inc.');

INSERT INTO oauth_client_key (client_key, secret, status, created_by, client_ident, callback, scope, expiration, created, environment, custom, master)
VALUES ('5eed868e-7ad0-4172-88f2-704bcf78b61e', '2054e4d7-77f2-46c9-bc4d-11a47255a6ec', 'ENABLED', 'admin', '123456800-otk', 'YOUR_SSG/oauth/v2/client/bcp?auth=done', 'openid email profile phone address', 0, 0, 'ALL', '{}', 0);

INSERT INTO oauth_client_key_created_by (client_key, secret, status, created_by, client_ident, callback, scope, expiration, created, environment, custom, master)
VALUES ('5eed868e-7ad0-4172-88f2-704bcf78b61e', '2054e4d7-77f2-46c9-bc4d-11a47255a6ec', 'ENABLED', 'admin', '123456800-otk', 'YOUR_SSG/oauth/v2/client/bcp?auth=done', 'openid email profile phone address', 0, 0, 'ALL', '{}', 0);

INSERT INTO oauth_client_key_user (client_key, created_by, client_ident)
VALUES ('5eed868e-7ad0-4172-88f2-704bcf78b61e','admin', '123456800-otk');

--
-- OpenID Connect Client for the Implicit Client Profile specification
--
INSERT INTO oauth_client (client_ident, name, description, organization, registered_by, type, created, custom)
VALUES ('123456801-otk', 'OpenID Connect Implicit Client Profile', 'Test for OpenID Connect ICP', 'CA Technologies Inc.', 'admin', 'public', 0, '{}');
--
INSERT INTO oauth_client_name_v2 (client_ident, name, organization)
VALUES ('123456801-otk', 'openid connect basic client profile', 'ca technologies inc.');
--
INSERT INTO oauth_client_org_v2 (client_ident, organization)
VALUES ('123456801-otk', 'ca technologies inc.');

INSERT INTO oauth_client_key (client_key, secret, status, created_by, client_ident, callback, scope, expiration, created, environment, custom, master)
VALUES ('5edc4a38-75ec-4617-8854-1a71ff1e0a2e', '5005a669-0295-4602-be7d-6a75342db6d8', 'ENABLED', 'admin', '123456801-otk', 'YOUR_SSG/oauth/v2/client/icp?auth=done', 'openid email profile phone address', 0, 0, 'ALL', '{}', 0);

INSERT INTO oauth_client_key_created_by (client_key, secret, status, created_by, client_ident, callback, scope, expiration, created, environment, custom, master)
VALUES ('5edc4a38-75ec-4617-8854-1a71ff1e0a2e', '5005a669-0295-4602-be7d-6a75342db6d8', 'ENABLED', 'admin', '123456801-otk', 'YOUR_SSG/oauth/v2/client/icp?auth=done', 'openid email profile phone address', 0, 0, 'ALL', '{}', 0);

INSERT INTO oauth_client_key_user (client_key, created_by, client_ident)
VALUES ('5edc4a38-75ec-4617-8854-1a71ff1e0a2e', 'admin', '123456801-otk');

--
-- Create an OAuth 2.0 client
--
--
INSERT INTO oauth_client (client_ident, name, description, organization, registered_by, type, created, custom)
VALUES ('TestClient2.0', 'OAuth2Client', 'OAuth 2.0 test client hosted on the ssg', 'CA Technologies Inc.', 'admin', 'confidential', 0, '{}');
--
INSERT INTO oauth_client_name_v2 (client_ident, name, organization)
VALUES ('TestClient2.0', 'oauth2client', 'ca technologies inc.');
--
INSERT INTO oauth_client_org_v2 (client_ident, organization)
VALUES ('TestClient2.0', 'ca technologies inc.');
--

INSERT INTO oauth_client_key (client_key, secret, status, created_by, client_ident, callback, expiration, created, environment, scope, custom, master)
VALUES ('54f0c455-4d80-421f-82ca-9194df24859d', 'a0f2742f-31c7-436f-9802-b7015b8fd8e6', 'ENABLED', 'admin', 'TestClient2.0', 'YOUR_SSG/oauth/v2/client/authcode?auth=done,YOUR_SSG/oauth/v2/client/implicit?auth=done', 0, 0, 'ALL', 'oob', '{}', 0);

INSERT INTO oauth_client_key_created_by (client_key, secret, status, created_by, client_ident, callback, expiration, created, environment, scope, custom, master)
VALUES ('54f0c455-4d80-421f-82ca-9194df24859d', 'a0f2742f-31c7-436f-9802-b7015b8fd8e6', 'ENABLED', 'admin', 'TestClient2.0', 'YOUR_SSG/oauth/v2/client/authcode?auth=done,YOUR_SSG/oauth/v2/client/implicit?auth=done', 0, 0, 'ALL', 'oob', '{}', 0);

INSERT INTO oauth_client_key_user (client_key, created_by, client_ident)
VALUES ('54f0c455-4d80-421f-82ca-9194df24859d', 'admin', 'TestClient2.0');


--
-- Swagger OAuth2 Client
--
INSERT INTO oauth_client (client_ident, name, description, organization, registered_by, type, custom)
VALUES ('2dc86f35-773c-47e2-958f-3f4bdfc5ea3a', 'Swagger OAuth2 Client', 'Swagger API testing', 'CA Technologies Inc.', 'admin', 'public', '{}');
--
INSERT INTO oauth_client_name_v2 (client_ident, name, organization)
VALUES ('2dc86f35-773c-47e2-958f-3f4bdfc5ea3a', 'swagger oauth2 client', 'ca technologies inc.');
--
INSERT INTO oauth_client_org_v2 (client_ident, organization)
VALUES ('2dc86f35-773c-47e2-958f-3f4bdfc5ea3a', 'ca technologies inc.');
--

INSERT INTO oauth_client_key (client_key, secret, status, created_by, client_ident, callback, scope, environment, custom)
VALUES ('2dc86f35-773c-47e2-958f-3f4bdfc5ea3a', '6aa6b190-056d-4a80-8604-79d9e44896ef', 'ENABLED', 'admin', '2dc86f35-773c-47e2-958f-3f4bdfc5ea3a', 'YOUR_SWAGGER_SERVER/oauth2-redirect.html', 'openid email profile phone address', 'ALL', '{
    "openid_registration": {
        "request": {},
        "response": {
            "client_id": "2dc86f35-773c-47e2-958f-3f4bdfc5ea3a",
            "client_secret": "6aa6b190-056d-4a80-8604-79d9e44896ef",
            "client_secret_expires_at": "0",
            "client_id_issued_at": 0,
            "registration_access_token": "",
            "registration_client_uri": "",
            "token_endpoint_auth_method": "client_secret_basic",
            "token_endpoint_auth_signing_alg": "",
            "application_type": "",
            "redirect_uris": [
                "YOUR_SWAGGER_SERVER/oauth2-redirect.html"
            ],
            "client_name": "Swagger OAuth2 Client",
            "subject_type": "pairwise",
            "sector_identifier_uri": "",
            "contacts": [],
            "response_types": [
                "code",
                "implicit",
                "token"
            ],
            "grant_types": [],
            "id_token_signed_response_alg": "HS256",
            "userinfo_signed_response_alg": "",
            "environment": "ALL",
            "organization": "CA Technologies Inc.",
            "master": false,
            "description": "Swagger API testing",
            "scope": "openid email profile phone address",
            "jwks": "",
            "jwks_uri": ""
        }
    }
}');
INSERT INTO oauth_client_key_created_by (client_key, secret, status, created_by, client_ident, callback, scope, environment, custom)
VALUES ('2dc86f35-773c-47e2-958f-3f4bdfc5ea3a', '6aa6b190-056d-4a80-8604-79d9e44896ef', 'ENABLED', 'admin', '2dc86f35-773c-47e2-958f-3f4bdfc5ea3a', 'YOUR_SWAGGER_SERVER/oauth2-redirect.html', 'openid email profile phone address', 'ALL', '{
    "openid_registration": {
        "request": {},
        "response": {
            "client_id": "2dc86f35-773c-47e2-958f-3f4bdfc5ea3a",
            "client_secret": "6aa6b190-056d-4a80-8604-79d9e44896ef",
            "client_secret_expires_at": "0",
            "client_id_issued_at": 0,
            "registration_access_token": "",
            "registration_client_uri": "",
            "token_endpoint_auth_method": "client_secret_basic",
            "token_endpoint_auth_signing_alg": "",
            "application_type": "",
            "redirect_uris": [
                "YOUR_SWAGGER_SERVER/oauth2-redirect.html"
            ],
            "client_name": "Swagger OAuth2 Client",
            "subject_type": "pairwise",
            "sector_identifier_uri": "",
            "contacts": [],
            "response_types": [
                "code",
                "implicit",
                "token"
            ],
            "grant_types": [],
            "id_token_signed_response_alg": "HS256",
            "userinfo_signed_response_alg": "",
            "environment": "ALL",
            "organization": "CA Technologies Inc.",
            "master": false,
            "description": "Swagger API testing",
            "scope": "openid email profile phone address",
            "jwks": "",
            "jwks_uri": ""
        }
    }
}');
INSERT INTO oauth_client_key_user (client_key, created_by, client_ident)
VALUES ('2dc86f35-773c-47e2-958f-3f4bdfc5ea3a', 'admin', '2dc86f35-773c-47e2-958f-3f4bdfc5ea3a');