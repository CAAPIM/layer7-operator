package graphman

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Bundle struct {
	WebApiServices                      []*WebApiServiceInput                     `json:"webApiServices,omitempty"`
	InternalWebApiServices              []*WebApiServiceInput                     `json:"internalWebApiServices,omitempty"`
	SoapServices                        []*SoapServiceInput                       `json:"soapServices,omitempty"`
	InternalSoapServices                []*SoapServiceInput                       `json:"internalSoapServices,omitempty"`
	PolicyFragments                     []*PolicyFragmentInput                    `json:"policyFragments,omitempty"`
	EncassConfigs                       []*EncassConfigInput                      `json:"encassConfigs,omitempty"`
	ClusterProperties                   []*ClusterPropertyInput                   `json:"clusterProperties,omitempty"`
	JdbcConnections                     []*JdbcConnectionInput                    `json:"jdbcConnections,omitempty"`
	TrustedCerts                        []*TrustedCertInput                       `json:"trustedCerts,omitempty"`
	Schemas                             []*SchemaInput                            `json:"schemas,omitempty"`
	Dtds                                []*DtdInput                               `json:"dtds,omitempty"`
	Fips                                []*FipInput                               `json:"fips,omitempty"`
	Ldaps                               []*LdapInput                              `json:"ldaps,omitempty"`
	InternalGroups                      []*InternalGroupInput                     `json:"internalGroups,omitempty"`
	FipGroups                           []*FipGroupInput                          `json:"fipGroups,omitempty"`
	InternalUsers                       []*InternalUserInput                      `json:"internalUsers,omitempty"`
	FipUsers                            []*FipUserInput                           `json:"fipUsers,omitempty"`
	Secrets                             []*SecretInput                            `json:"secrets,omitempty"`
	Keys                                []*KeyInput                               `json:"keys,omitempty"`
	CassandraConnections                []*CassandraConnectionInput               `json:"cassandraConnections,omitempty"`
	JmsDestinations                     []*JmsDestinationInput                    `json:"jmsDestinations,omitempty"`
	GlobalPolicies                      []*GlobalPolicyInput                      `json:"globalPolicies,omitempty"`
	BackgroundTasks                     []*BackgroundTaskPolicyInput              `json:"backgroundTaskPolicies,omitempty"`
	ScheduledTasks                      []*ScheduledTaskInput                     `json:"scheduledTasks,omitempty"`
	ServerModuleFiles                   []*ServerModuleFileInput                  `json:"serverModuleFiles,omitempty"`
	SiteMinderConfigs                   []*SMConfigInput                          `json:"smConfigs,omitempty"`
	ActiveConnectors                    []*ActiveConnectorInput                   `json:"activeConnectors,omitempty"`
	EmailListeners                      []*EmailListenerInput                     `json:"emailListeners,omitempty"`
	ListenPorts                         []*ListenPortInput                        `json:"listenPorts,omitempty"`
	AdministrativeUserAccountProperties []*AdministrativeUserAccountPropertyInput `json:"administrativeUserAccountProperties,omitempty"`
	PasswordPolicies                    []*PasswordPolicyInput                    `json:"passwordPolicies,omitempty"`
	RevocationCheckPolicies             []*RevocationCheckPolicyInput             `json:"revocationCheckPolicies,omitempty"`
	LogSinks                            []*LogSinkInput                           `json:"logSinks,omitempty"`
	HttpConfigurations                  []*HttpConfigurationInput                 `json:"httpConfigurations,omitempty"`
	CustomKeyValues                     []*CustomKeyValueInput                    `json:"customKeyValues,omitempty"`
	ServiceResolutionConfigs            []*ServiceResolutionConfigInput           `json:"serviceResolutionConfigs,omitempty"`
	Folders                             []*FolderInput                            `json:"folders,omitempty"`
	FederatedIdps                       []*FederatedIdpInput                      `json:"federatedIdps,omitempty"`
	FederatedGroups                     []*FederatedGroupInput                    `json:"federatedGroups,omitempty"`
	FederatedUsers                      []*FederatedUserInput                     `json:"federatedUsers,omitempty"`
	InternalIdps                        []*InternalIdpInput                       `json:"internalIdps,omitempty"`
	LdapIdps                            []*LdapIdpInput                           `json:"ldapIdps,omitempty"`
	SimpleLdapIdps                      []*SimpleLdapIdpInput                     `json:"simpleLdapIdps,omitempty"`
	PolicyBackedIdps                    []*PolicyBackedIdpInput                   `json:"policyBackedIdps,omitempty"`
	Policies                            []*L7PolicyInput                          `json:"policies,omitempty"`
	Services                            []*L7ServiceInput                         `json:"services,omitempty"`
}

var entities = []string{
	"clusterProperties",
	"encassConfigs",
	"jdbcConnections",
	"cassandraConnections",
	"trustedCerts",
	"schemas",
	"dtds",
	"fips",
	"ldaps",
	"internalGroups",
	"fipGroups",
	"internalUsers",
	"fipUsers",
	"scheduledTasks",
	"jmsDestinations",
	"secrets",
	"keys",
	"listenPorts",
	"activeConnectors",
	"serverModuleFiles",
	"emailListeners",
	"smConfigs",
	"administrativeUserAccountProperties",
	"passwordPolicies",
	"revocationCheckPolicies",
	"logSinks",
	"httpConfigurations",
	"customKeyValues",
	"serviceResolutionConfigs",
	"folders",
	".internalwebapi",
	".webapi",
	".soap",
	".internalsoap",
	".global",
	".policy",
	"federatedIdps",
	"federatedGroups",
	"federatedUsers",
	"internalIdps",
	"ldapIdps",
	"simpleLdapIdps",
	"policyBackedIdps",
	"policies",
	"services",
	".service",
}

type BundleApplyErrors struct {
	Errors []BundleApplyError `json:"errors,omitempty"`
}

type BundleApplyError struct {
	Name   string `json:"name,omitempty"`
	Status string `json:"status,omitempty"`
	Entity string `json:"entity,omitempty"`
	Detail string `json:"detail,omitempty"`
}

func ConcatBundle(src []byte, dest []byte) ([]byte, error) {
	srcBundle := Bundle{}
	destBundle := Bundle{}

	err := json.Unmarshal(dest, &destBundle)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(src, &srcBundle)
	if err != nil {
		return nil, err
	}

	destBundle.ActiveConnectors = append(destBundle.ActiveConnectors, srcBundle.ActiveConnectors...)
	destBundle.BackgroundTasks = append(destBundle.BackgroundTasks, srcBundle.BackgroundTasks...)
	destBundle.CassandraConnections = append(destBundle.CassandraConnections, srcBundle.CassandraConnections...)
	destBundle.ClusterProperties = append(destBundle.ClusterProperties, srcBundle.ClusterProperties...)
	destBundle.Dtds = append(destBundle.Dtds, srcBundle.Dtds...)
	destBundle.EmailListeners = append(destBundle.EmailListeners, srcBundle.EmailListeners...)
	destBundle.EncassConfigs = append(destBundle.EncassConfigs, srcBundle.EncassConfigs...)
	destBundle.FipGroups = append(destBundle.FipGroups, srcBundle.FipGroups...)
	destBundle.FipUsers = append(destBundle.FipUsers, srcBundle.FipUsers...)
	destBundle.Fips = append(destBundle.Fips, srcBundle.Fips...)
	destBundle.GlobalPolicies = append(destBundle.GlobalPolicies, srcBundle.GlobalPolicies...)
	destBundle.InternalGroups = append(destBundle.InternalGroups, srcBundle.InternalGroups...)
	destBundle.InternalSoapServices = append(destBundle.InternalSoapServices, srcBundle.InternalSoapServices...)
	destBundle.InternalUsers = append(destBundle.InternalUsers, srcBundle.InternalUsers...)
	destBundle.InternalWebApiServices = append(destBundle.InternalWebApiServices, srcBundle.InternalWebApiServices...)
	destBundle.JdbcConnections = append(destBundle.JdbcConnections, srcBundle.JdbcConnections...)
	destBundle.JmsDestinations = append(destBundle.JmsDestinations, srcBundle.JmsDestinations...)
	destBundle.Keys = append(destBundle.Keys, srcBundle.Keys...)
	destBundle.LdapIdps = append(destBundle.LdapIdps, srcBundle.LdapIdps...)
	destBundle.ListenPorts = append(destBundle.ListenPorts, srcBundle.ListenPorts...)
	destBundle.PolicyFragments = append(destBundle.PolicyFragments, srcBundle.PolicyFragments...)
	destBundle.ScheduledTasks = append(destBundle.ScheduledTasks, srcBundle.ScheduledTasks...)
	destBundle.Schemas = append(destBundle.Schemas, srcBundle.Schemas...)
	destBundle.Secrets = append(destBundle.Secrets, srcBundle.Secrets...)
	destBundle.ServerModuleFiles = append(destBundle.ServerModuleFiles, srcBundle.ServerModuleFiles...)
	destBundle.SiteMinderConfigs = append(destBundle.SiteMinderConfigs, srcBundle.SiteMinderConfigs...)
	destBundle.SoapServices = append(destBundle.SoapServices, srcBundle.InternalSoapServices...)
	destBundle.TrustedCerts = append(destBundle.TrustedCerts, srcBundle.TrustedCerts...)
	destBundle.WebApiServices = append(destBundle.WebApiServices, srcBundle.WebApiServices...)
	destBundle.AdministrativeUserAccountProperties = append(destBundle.AdministrativeUserAccountProperties, srcBundle.AdministrativeUserAccountProperties...)
	destBundle.PasswordPolicies = append(destBundle.PasswordPolicies, srcBundle.PasswordPolicies...)
	destBundle.RevocationCheckPolicies = append(destBundle.RevocationCheckPolicies, srcBundle.RevocationCheckPolicies...)
	destBundle.LogSinks = append(destBundle.LogSinks, srcBundle.LogSinks...)
	destBundle.HttpConfigurations = append(destBundle.HttpConfigurations, srcBundle.HttpConfigurations...)
	destBundle.CustomKeyValues = append(destBundle.CustomKeyValues, srcBundle.CustomKeyValues...)
	destBundle.ServiceResolutionConfigs = append(destBundle.ServiceResolutionConfigs, srcBundle.ServiceResolutionConfigs...)
	destBundle.Folders = append(destBundle.Folders, srcBundle.Folders...)
	destBundle.FederatedIdps = append(destBundle.FederatedIdps, srcBundle.FederatedIdps...)
	destBundle.FederatedGroups = append(destBundle.FederatedGroups, srcBundle.FederatedGroups...)
	destBundle.FederatedUsers = append(destBundle.FederatedUsers, srcBundle.FederatedUsers...)
	destBundle.InternalIdps = append(destBundle.InternalIdps, srcBundle.InternalIdps...)
	destBundle.LdapIdps = append(destBundle.LdapIdps, srcBundle.LdapIdps...)
	destBundle.SimpleLdapIdps = append(destBundle.SimpleLdapIdps, srcBundle.SimpleLdapIdps...)
	destBundle.PolicyBackedIdps = append(destBundle.PolicyBackedIdps, srcBundle.PolicyBackedIdps...)
	destBundle.Policies = append(destBundle.Policies, srcBundle.Policies...)
	destBundle.Services = append(destBundle.Services, srcBundle.Services...)

	bundleBytes, err := json.Marshal(destBundle)
	if err != nil {
		return nil, err
	}

	return bundleBytes, nil

}

// Implode - convert an exploded Graphman directory into a single JSON file.
func Implode(path string) ([]byte, error) {
	bundle, err := implodeBundle(path)
	if err != nil {
		return nil, err
	}

	bundleBytes, err := json.Marshal(bundle)
	if err != nil {
		return nil, err
	}
	return bundleBytes, nil
}

func RemoveL7PortalApi(username string, password string, target string, apiName string, policyFragmentName string) ([]byte, error) {
	resp, err := deleteL7PortalApi(context.Background(), gqlClient(username, password, target, ""), []string{apiName}, []string{policyFragmentName})
	if err != nil {
		return nil, err
	}
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	return respBytes, nil
}

func ApplyDynamicBundle(username string, password string, target string, encpass string, bundleBytes []byte) ([]byte, error) {
	bundle := Bundle{}
	err := json.Unmarshal(bundleBytes, &bundle)
	if err != nil {
		return nil, err
	}

	resp, err := installBundle(context.Background(), gqlClient(username, password, target, encpass), bundle.ActiveConnectors, bundle.AdministrativeUserAccountProperties, bundle.BackgroundTasks, bundle.CassandraConnections, bundle.ClusterProperties, bundle.Dtds, bundle.EmailListeners, bundle.EncassConfigs, bundle.FipGroups, bundle.FipUsers, bundle.Fips, bundle.FederatedGroups, bundle.FederatedUsers, bundle.InternalIdps, bundle.FederatedIdps, bundle.LdapIdps, bundle.SimpleLdapIdps, bundle.PolicyBackedIdps, bundle.GlobalPolicies, bundle.InternalGroups, bundle.InternalSoapServices, bundle.InternalUsers, bundle.InternalWebApiServices, bundle.JdbcConnections, bundle.JmsDestinations, bundle.Keys, bundle.Ldaps, bundle.ListenPorts, bundle.PasswordPolicies, bundle.Policies, bundle.PolicyFragments, bundle.RevocationCheckPolicies, bundle.ScheduledTasks, bundle.LogSinks, bundle.Schemas, bundle.Secrets, bundle.HttpConfigurations, bundle.CustomKeyValues, bundle.ServerModuleFiles, bundle.ServiceResolutionConfigs, bundle.Folders, bundle.SiteMinderConfigs, bundle.Services, bundle.SoapServices, bundle.TrustedCerts, bundle.WebApiServices)

	detailedErr := CheckDetailedStatus(bundle, resp)
	if detailedErr != nil {
		return nil, detailedErr
	}

	if err != nil {
		return nil, err
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return respBytes, nil
}

// parseEntity determines which entity the file from a Graphman directory belongs to
// this works with a static list of globally defined entities
func parseEntity(path string) (string, bool) {
	for _, e := range entities {
		if strings.Contains(path, e) {
			return e, true
		}
	}

	return "", false
}

// readBundle unmarshals a JSON file in the specified Graphman directory into the working Bundle object.
func readBundle(entityType string, file string, bundle *Bundle) (Bundle, error) {
	f, _ := os.ReadFile(file)
	switch entityType {
	case ".webapi":
		webApiService := WebApiServiceInput{}
		err := json.Unmarshal(f, &webApiService)
		if err != nil {
			return *bundle, err
		}
		bundle.WebApiServices = append(bundle.WebApiServices, &webApiService)
	case ".internalwebapi":
		internalWebApiService := WebApiServiceInput{}
		err := json.Unmarshal(f, &internalWebApiService)
		if err != nil {
			return *bundle, err
		}
		bundle.WebApiServices = append(bundle.InternalWebApiServices, &internalWebApiService)
	case ".soap":
		soapService := SoapServiceInput{}
		err := json.Unmarshal(f, &soapService)
		if err != nil {
			return *bundle, err
		}
		bundle.SoapServices = append(bundle.SoapServices, &soapService)
	case ".internalsoap":
		internalSoapService := SoapServiceInput{}
		err := json.Unmarshal(f, &internalSoapService)
		if err != nil {
			return *bundle, err
		}
		bundle.SoapServices = append(bundle.InternalSoapServices, &internalSoapService)
	case ".global":
		globalPolicy := GlobalPolicyInput{}
		err := json.Unmarshal(f, &globalPolicy)
		if err != nil {
			return *bundle, err
		}
		bundle.GlobalPolicies = append(bundle.GlobalPolicies, &globalPolicy)
	case ".policy":
		policyFragment := L7PolicyInput{}
		err := json.Unmarshal(f, &policyFragment)
		if err != nil {
			return *bundle, err
		}
		bundle.Policies = append(bundle.Policies, &policyFragment)
	case ".service":
		service := L7ServiceInput{}
		err := json.Unmarshal(f, &service)
		if err != nil {
			return *bundle, nil
		}
		bundle.Services = append(bundle.Services, &service)
	case "clusterProperties":
		clusterProperty := ClusterPropertyInput{}
		err := json.Unmarshal(f, &clusterProperty)
		if err != nil {
			return *bundle, err
		}
		bundle.ClusterProperties = append(bundle.ClusterProperties, &clusterProperty)
	case "scheduledTasks":
		scheduledTask := ScheduledTaskInput{}
		err := json.Unmarshal(f, &scheduledTask)
		if err != nil {
			return *bundle, err
		}
		bundle.ScheduledTasks = append(bundle.ScheduledTasks, &scheduledTask)
	case "encassConfigs":
		encassConfig := EncassConfigInput{}
		err := json.Unmarshal(f, &encassConfig)
		if err != nil {
			return *bundle, err
		}
		bundle.EncassConfigs = append(bundle.EncassConfigs, &encassConfig)
	case "jdbcConnections":
		jdbcConnection := JdbcConnectionInput{}
		err := json.Unmarshal(f, &jdbcConnection)
		if err != nil {
			return *bundle, err
		}
		bundle.JdbcConnections = append(bundle.JdbcConnections, &jdbcConnection)
	case "trustedCerts":
		trustedCert := TrustedCertInput{}
		err := json.Unmarshal(f, &trustedCert)
		if err != nil {
			return *bundle, err
		}
		bundle.TrustedCerts = append(bundle.TrustedCerts, &trustedCert)
	case "schemas":
		schema := SchemaInput{}
		err := json.Unmarshal(f, &schema)
		if err != nil {
			return *bundle, err
		}
		bundle.Schemas = append(bundle.Schemas, &schema)
	case "dtds":
		dtd := DtdInput{}
		err := json.Unmarshal(f, &dtd)
		if err != nil {
			return *bundle, err
		}
		bundle.Dtds = append(bundle.Dtds, &dtd)
	case "fips":
		fip := FipInput{}
		err := json.Unmarshal(f, &fip)
		if err != nil {
			return *bundle, err
		}
		bundle.Fips = append(bundle.Fips, &fip)
	case "ldaps":
		ldap := LdapInput{}
		err := json.Unmarshal(f, &ldap)
		if err != nil {
			return *bundle, err
		}
		bundle.Ldaps = append(bundle.Ldaps, &ldap)
	case "internalGroups":
		internalGroup := InternalGroupInput{}
		err := json.Unmarshal(f, &internalGroup)
		if err != nil {
			return *bundle, err
		}
		bundle.InternalGroups = append(bundle.InternalGroups, &internalGroup)
	case "fipGroups":
		fipsGroup := FipGroupInput{}
		err := json.Unmarshal(f, &fipsGroup)
		if err != nil {
			return *bundle, err
		}
		bundle.FipGroups = append(bundle.FipGroups, &fipsGroup)
	case "internalUsers":
		internalUser := InternalUserInput{}
		err := json.Unmarshal(f, &internalUser)
		if err != nil {
			return *bundle, err
		}
		bundle.InternalUsers = append(bundle.InternalUsers, &internalUser)
	case "fipUsers":
		fipsUser := FipUserInput{}
		err := json.Unmarshal(f, &fipsUser)
		if err != nil {
			return *bundle, err
		}
		bundle.FipUsers = append(bundle.FipUsers, &fipsUser)
	case "secrets":
		secret := SecretInput{}
		err := json.Unmarshal(f, &secret)
		if err != nil {
			return *bundle, err
		}
		bundle.Secrets = append(bundle.Secrets, &secret)
	case "keys":
		key := KeyInput{}
		err := json.Unmarshal(f, &key)
		if err != nil {
			return *bundle, err
		}
		bundle.Keys = append(bundle.Keys, &key)
	case "jmsDestinations":
		jmsDestination := JmsDestinationInput{}
		err := json.Unmarshal(f, &jmsDestination)
		if err != nil {
			return *bundle, err
		}
		bundle.JmsDestinations = append(bundle.JmsDestinations, &jmsDestination)

	case "activeConnectors":
		activeConnector := ActiveConnectorInput{}
		err := json.Unmarshal(f, &activeConnector)
		if err != nil {
			return *bundle, nil
		}
		bundle.ActiveConnectors = append(bundle.ActiveConnectors, &activeConnector)
	case "listenPorts":
		listenPort := ListenPortInput{}
		err := json.Unmarshal(f, &listenPort)
		if err != nil {
			return *bundle, nil
		}
		bundle.ListenPorts = append(bundle.ListenPorts, &listenPort)
	case "emailListeners":
		emailListener := EmailListenerInput{}
		err := json.Unmarshal(f, &emailListener)
		if err != nil {
			return *bundle, nil
		}
		bundle.EmailListeners = append(bundle.EmailListeners, &emailListener)
	case "serverModuleFiles":
		serverModuleFile := ServerModuleFileInput{}
		err := json.Unmarshal(f, &serverModuleFile)
		if err != nil {
			return *bundle, nil
		}
		bundle.ServerModuleFiles = append(bundle.ServerModuleFiles, &serverModuleFile)
	case "smConfigs":
		smConfig := SMConfigInput{}
		err := json.Unmarshal(f, &smConfig)
		if err != nil {
			return *bundle, nil
		}
		bundle.SiteMinderConfigs = append(bundle.SiteMinderConfigs, &smConfig)

	case "administrativeUserAccountProperties":
		administrativeUserAccountProperty := AdministrativeUserAccountPropertyInput{}
		err := json.Unmarshal(f, &administrativeUserAccountProperty)
		if err != nil {
			return *bundle, nil
		}
		bundle.AdministrativeUserAccountProperties = append(bundle.AdministrativeUserAccountProperties, &administrativeUserAccountProperty)
	case "passwordPolicies":
		passwordPolicy := PasswordPolicyInput{}
		err := json.Unmarshal(f, &passwordPolicy)
		if err != nil {
			return *bundle, nil
		}
		bundle.PasswordPolicies = append(bundle.PasswordPolicies, &passwordPolicy)
	case "revocationCheckPolicies":
		revocationCheckPolicy := RevocationCheckPolicyInput{}
		err := json.Unmarshal(f, &revocationCheckPolicy)
		if err != nil {
			return *bundle, nil
		}
		bundle.RevocationCheckPolicies = append(bundle.RevocationCheckPolicies, &revocationCheckPolicy)
	case "logSinks":
		logSink := LogSinkInput{}
		err := json.Unmarshal(f, &logSink)
		if err != nil {
			return *bundle, nil
		}
		bundle.LogSinks = append(bundle.LogSinks, &logSink)
	case "httpConfigurations":
		httpConfiguration := HttpConfigurationInput{}
		err := json.Unmarshal(f, &httpConfiguration)
		if err != nil {
			return *bundle, nil
		}
		bundle.HttpConfigurations = append(bundle.HttpConfigurations, &httpConfiguration)
	case "customKeyValues":
		customKeyValue := CustomKeyValueInput{}
		err := json.Unmarshal(f, &customKeyValue)
		if err != nil {
			return *bundle, nil
		}
		bundle.CustomKeyValues = append(bundle.CustomKeyValues, &customKeyValue)
	case "serviceResolutionConfigs":
		serviceResolutionConfig := ServiceResolutionConfigInput{}
		err := json.Unmarshal(f, &serviceResolutionConfig)
		if err != nil {
			return *bundle, nil
		}
		bundle.ServiceResolutionConfigs = append(bundle.ServiceResolutionConfigs, &serviceResolutionConfig)
	case "folders":
		folder := FolderInput{}
		err := json.Unmarshal(f, &folder)
		if err != nil {
			return *bundle, nil
		}
		bundle.Folders = append(bundle.Folders, &folder)
	case "federatedIdps":
		federatedIdp := FederatedIdpInput{}
		err := json.Unmarshal(f, &federatedIdp)
		if err != nil {
			return *bundle, nil
		}
		bundle.FederatedIdps = append(bundle.FederatedIdps, &federatedIdp)
	case "federatedGroups":
		federatedGroup := FederatedGroupInput{}
		err := json.Unmarshal(f, &federatedGroup)
		if err != nil {
			return *bundle, nil
		}
		bundle.FederatedGroups = append(bundle.FederatedGroups, &federatedGroup)
	case "federatedUsers":
		federatedUser := FederatedUserInput{}
		err := json.Unmarshal(f, &federatedUser)
		if err != nil {
			return *bundle, nil
		}
		bundle.FederatedUsers = append(bundle.FederatedUsers, &federatedUser)
	case "internalIdps":
		internalIdp := InternalIdpInput{}
		err := json.Unmarshal(f, &internalIdp)
		if err != nil {
			return *bundle, nil
		}
		bundle.InternalIdps = append(bundle.InternalIdps, &internalIdp)
	case "ldapIdps":
		ldapIdp := LdapIdpInput{}
		err := json.Unmarshal(f, &ldapIdp)
		if err != nil {
			return *bundle, nil
		}
		bundle.LdapIdps = append(bundle.LdapIdps, &ldapIdp)
	case "simpleLdapIdps":
		simpleLdapIdp := SimpleLdapIdpInput{}
		err := json.Unmarshal(f, &simpleLdapIdp)
		if err != nil {
			return *bundle, nil
		}
		bundle.SimpleLdapIdps = append(bundle.SimpleLdapIdps, &simpleLdapIdp)
	case "policyBackedIdps":
		policyBackedIdp := PolicyBackedIdpInput{}
		err := json.Unmarshal(f, &policyBackedIdp)
		if err != nil {
			return *bundle, nil
		}
		bundle.PolicyBackedIdps = append(bundle.PolicyBackedIdps, &policyBackedIdp)
	case "policies":
		policy := L7PolicyInput{}
		err := json.Unmarshal(f, &policy)
		if err != nil {
			return *bundle, nil
		}
		bundle.Policies = append(bundle.Policies, &policy)
	case "services":
		service := L7ServiceInput{}
		err := json.Unmarshal(f, &service)
		if err != nil {
			return *bundle, nil
		}
		bundle.Services = append(bundle.Services, &service)
	}

	return *bundle, nil
}

// implodeBundle takes a Graphman directory and returns a Bundle Object.
func implodeBundle(path string) (Bundle, error) {
	bundle := Bundle{}
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			name, isEntity := parseEntity(path)
			if isEntity {
				bundle, err = readBundle(name, path, &bundle)
				if err != nil {
					return err
				}
				return nil
			}
		}
		return nil
	})

	if err != nil {
		return bundle, err
	}
	return bundle, nil
}

// Not used - reserved for future use
// func parseEntities(bundle Bundle) {
// 	v := reflect.ValueOf(bundle)
// 	typeOfS := v.Type()

// 	for i := 0; i < v.NumField(); i++ {
// 		fmt.Printf("%s %v\n", typeOfS.Field(i).Name, v.Field(i).Interface())
// 	}
// }

func CheckDetailedStatus(bundle Bundle, resp *installBundleResponse) error {
	var bundleApplyErrors BundleApplyErrors

	if resp.SetActiveConnectors != nil {
		for i, r := range resp.SetActiveConnectors.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "ActiveConnectors", Name: bundle.ActiveConnectors[i].Name})
			}
		}
	}
	if resp.SetBackgroundTaskPolicies != nil {
		for i, r := range resp.SetBackgroundTaskPolicies.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "BackgroundTaskPolicies", Name: bundle.BackgroundTasks[i].Name})
			}
		}
	}
	if resp.SetCassandraConnections != nil {
		for i, r := range resp.SetCassandraConnections.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "CassandraConnections", Name: bundle.CassandraConnections[i].Name})
			}
		}
	}
	if resp.SetClusterProperties != nil {
		for i, r := range resp.SetClusterProperties.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "ClusterProperties", Name: bundle.ClusterProperties[i].Name})
			}
		}
	}
	if resp.SetDtds != nil {
		for i, r := range resp.SetDtds.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "Dtds", Name: bundle.Dtds[i].Description})
			}
		}
	}
	if resp.SetEmailListeners != nil {
		for i, r := range resp.SetEmailListeners.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "EmailListeners", Name: bundle.EmailListeners[i].Name})
			}
		}
	}
	if resp.SetEncassConfigs != nil {
		for i, r := range resp.SetEncassConfigs.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "EncassConfigs", Name: bundle.EncassConfigs[i].Name})
			}
		}
	}
	if resp.SetFipGroups != nil {
		for i, r := range resp.SetFipGroups.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "FipGroups", Name: bundle.FipGroups[i].Name})
			}
		}
	}
	if resp.SetFipUsers != nil {
		for i, r := range resp.SetFipUsers.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "FipUsers", Name: bundle.FipUsers[i].Name})
			}
		}
	}
	if resp.SetFips != nil {
		for i, r := range resp.SetFips.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "Fips", Name: bundle.Fips[i].Name})
			}
		}
	}
	if resp.SetGlobalPolicies != nil {
		for i, r := range resp.SetGlobalPolicies.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "GlobalPolicies", Name: bundle.GlobalPolicies[i].Name})
			}
		}
	}
	if resp.SetInternalGroups != nil {
		for i, r := range resp.SetInternalGroups.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "InternalGroups", Name: bundle.InternalGroups[i].Name})
			}
		}
	}
	if resp.SetInternalSoapServices != nil {
		for i, r := range resp.SetInternalSoapServices.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "InternalSoapServices", Name: bundle.InternalSoapServices[i].Name})
			}
		}
	}
	if resp.SetInternalUsers != nil {
		for i, r := range resp.SetInternalUsers.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "InternalUsers", Name: bundle.InternalUsers[i].Name})
			}
		}
	}
	if resp.SetInternalWebApiServices != nil {
		for i, r := range resp.SetInternalWebApiServices.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "InternalWebApiServices", Name: bundle.InternalWebApiServices[i].Name})
			}
		}
	}
	if resp.SetJdbcConnections != nil {
		for i, r := range resp.SetJdbcConnections.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "JdbcConnections", Name: bundle.JdbcConnections[i].Name})
			}
		}
	}
	if resp.SetJmsDestinations != nil {
		for i, r := range resp.SetJmsDestinations.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "JmsDestinations", Name: bundle.JmsDestinations[i].Name})
			}
		}
	}
	if resp.SetKeys != nil {
		for i, r := range resp.SetKeys.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "Keys", Name: bundle.Keys[i].Alias})
			}
		}
	}
	if resp.SetLdaps != nil {
		for i, r := range resp.SetLdaps.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "Ldaps", Name: bundle.LdapIdps[i].Name})
			}
		}
	}
	if resp.SetListenPorts != nil {
		for i, r := range resp.SetListenPorts.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "ListenPorts", Name: bundle.ListenPorts[i].Name})
			}
		}
	}
	if resp.SetPolicyFragments != nil {
		for i, r := range resp.SetPolicyFragments.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "PolicyFragments", Name: bundle.PolicyFragments[i].Name})
			}
		}
	}
	if resp.SetSMConfigs != nil {
		for i, r := range resp.SetSMConfigs.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "SMConfigs", Name: bundle.SiteMinderConfigs[i].Name})
			}
		}
	}
	if resp.SetScheduledTasks != nil {
		for i, r := range resp.SetScheduledTasks.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "ScheduledTasks", Name: bundle.ScheduledTasks[i].Name})
			}
		}
	}
	if resp.SetSchemas != nil {
		for i, r := range resp.SetSchemas.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "Schemas", Name: bundle.Schemas[i].Description})
			}
		}
	}
	if resp.SetSecrets != nil {
		for i, r := range resp.SetSecrets.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "Secrets", Name: bundle.Secrets[i].Name})
			}
		}
	}
	if resp.SetServerModuleFiles != nil {
		for i, r := range resp.SetServerModuleFiles.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "ServerModuleFiles", Name: bundle.ServerModuleFiles[i].Name})
			}
		}
	}
	if resp.SetSoapServices != nil {
		for i, r := range resp.SetSoapServices.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "SoapServices", Name: bundle.SoapServices[i].Name})
			}
		}
	}
	if resp.SetTrustedCerts != nil {
		for i, r := range resp.SetTrustedCerts.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "TrustedCerts", Name: bundle.TrustedCerts[i].Name})
			}
		}
	}
	if resp.SetWebApiServices != nil {
		for i, r := range resp.SetWebApiServices.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "WebApiServices", Name: bundle.WebApiServices[i].Name})
			}
		}
	}

	if resp.SetAdministrativeUserAccountProperties != nil {
		for i, r := range resp.SetAdministrativeUserAccountProperties.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "AdministrativeUserAccountProperties", Name: bundle.AdministrativeUserAccountProperties[i].Name})
			}
		}
	}
	if resp.SetPasswordPolicies != nil {
		for _, r := range resp.SetPasswordPolicies.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "PasswordPolicies", Name: "Password Policy"})
			}
		}
	}
	if resp.SetRevocationCheckPolicies != nil {
		for i, r := range resp.SetRevocationCheckPolicies.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "RevocationCheckPolicies", Name: bundle.RevocationCheckPolicies[i].Name})
			}
		}
	}
	if resp.SetLogSinks != nil {
		for i, r := range resp.SetLogSinks.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "LogSinks", Name: bundle.LogSinks[i].Name})
			}
		}
	}
	if resp.SetHttpConfigurations != nil {
		for i, r := range resp.SetHttpConfigurations.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "HttpConfigurations", Name: bundle.HttpConfigurations[i].Host})
			}
		}
	}
	if resp.SetCustomKeyValues != nil {
		for i, r := range resp.SetCustomKeyValues.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "CustomKeyValues", Name: bundle.CustomKeyValues[i].Key})
			}
		}
	}
	if resp.SetServiceResolutionConfigs != nil {
		for _, r := range resp.SetServiceResolutionConfigs.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "ServiceResolutionConfigs", Name: "Service Resolution Config"})
			}
		}
	}
	if resp.SetFolders != nil {
		for i, r := range resp.SetFolders.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "Folders", Name: bundle.Folders[i].Name})
			}
		}
	}

	if resp.SetFederatedIdps != nil {
		for i, r := range resp.SetFederatedIdps.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "FederatedIdps", Name: bundle.FederatedIdps[i].Name})
			}
		}
	}
	if resp.SetFederatedGroups != nil {
		for i, r := range resp.SetFederatedGroups.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "FederatedGroups", Name: bundle.FederatedGroups[i].Name})
			}
		}
	}
	if resp.SetFederatedUsers != nil {
		for i, r := range resp.SetFederatedUsers.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "FederatedUsers", Name: bundle.FederatedUsers[i].Name})
			}
		}
	}
	if resp.SetInternalIdps != nil {
		for i, r := range resp.SetInternalIdps.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "InternalIdps", Name: bundle.InternalIdps[i].Name})
			}
		}
	}
	if resp.SetLdapIdps != nil {
		for i, r := range resp.SetLdapIdps.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "LdapIdps", Name: bundle.LdapIdps[i].Name})
			}
		}
	}
	if resp.SetSimpleLdapIdps != nil {
		for i, r := range resp.SetSimpleLdapIdps.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "SimpleLdapIdps", Name: bundle.SimpleLdapIdps[i].Name})
			}
		}
	}
	if resp.SetPolicyBackedIdps != nil {
		for i, r := range resp.SetPolicyBackedIdps.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "PolicyBackedIdps", Name: bundle.PolicyBackedIdps[i].Name})
			}
		}
	}
	if resp.SetPolicies != nil {
		for i, r := range resp.SetPolicies.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "Policies", Name: bundle.Policies[i].Name})
			}
		}
	}
	if resp.SetServices != nil {
		for i, r := range resp.SetServices.DetailedStatus {
			if r.Status == "ERROR" {
				bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "Services", Name: bundle.Services[i].Name})
			}
		}
	}

	if len(bundleApplyErrors.Errors) > 0 {
		errorBytes, _ := json.Marshal(bundleApplyErrors)
		return fmt.Errorf("errors: %s", string(errorBytes))
	}

	return nil
}
