package graphman

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"github.com/Khan/genqlient/graphql"
)

type MappingAction string

const (
	MappingActionNewOrUpdate     MappingAction = "NEW_OR_UPDATE"
	MappingActionNewOrExisting   MappingAction = "NEW_OR_EXISTING"
	MappingActionAlwaysCreateNew MappingAction = "ALWAYS_CREATE_NEW"
	MappingActionDelete          MappingAction = "DELETE"
	MappingActionIgnore          MappingAction = "IGNORE"
)

type MutationStatus string

const (
	MutationStatusNone         MutationStatus = "NONE"
	MutationStatusCreated      MutationStatus = "CREATED"
	MutationStatusUpdated      MutationStatus = "UPDATED"
	MutationStatusDeleted      MutationStatus = "DELETED"
	MutationStatusUsedExisting MutationStatus = "USED_EXISTING"
	MutationStatusIgnored      MutationStatus = "IGNORED"
	MutationStatusError        MutationStatus = "ERROR"
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
	Roles                               []*RoleInput                              `json:"roles,omitempty"`
	GenericEntities                     []*GenericEntityInput                     `json:"genericEntities,omitempty"`
	AuditConfigurations                 []*AuditConfigurationInput                `json:"auditConfigurations,omitempty"`
	Properties                          *BundleProperties                         `json:"properties,omitempty"`
}

type MappingInstructionInput struct {
	Action         MappingAction `json:"action,omitempty"`
	Default        bool          `json:"default,omitempty"`
	FailOnNew      bool          `json:"failOnNew,omitempty"`
	FailOnExisting bool          `json:"failOnExisting,omitempty"`
	Nodef          bool          `json:"nodef,omitempty"`
	Source         interface{}   `json:"source,omitempty"`
}

type BundlePropertyMeta struct {
	Id        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Author    string `json:"author,omitempty"`
	Hostname  string `json:"hostname,omitempty"`
	TimeStamp string `json:"timestamp,omitempty"`
}

type BundleProperties struct {
	Meta          BundlePropertyMeta `json:"meta,omitempty"`
	DefaultAction MappingAction      `json:"defaultAction,omitempty"`
	Mappings      BundleMappings     `json:"mappings,omitempty"`
}

type BundleMappings struct {
	WebApiServices                      []*MappingInstructionInput `json:"webApiServices,omitempty"`
	InternalWebApiServices              []*MappingInstructionInput `json:"internalWebApiServices,omitempty"`
	SoapServices                        []*MappingInstructionInput `json:"soapServices,omitempty"`
	InternalSoapServices                []*MappingInstructionInput `json:"internalSoapServices,omitempty"`
	PolicyFragments                     []*MappingInstructionInput `json:"policyFragments,omitempty"`
	EncassConfigs                       []*MappingInstructionInput `json:"encassConfigs,omitempty"`
	ClusterProperties                   []*MappingInstructionInput `json:"clusterProperties,omitempty"`
	JdbcConnections                     []*MappingInstructionInput `json:"jdbcConnections,omitempty"`
	TrustedCerts                        []*MappingInstructionInput `json:"trustedCerts,omitempty"`
	Schemas                             []*MappingInstructionInput `json:"schemas,omitempty"`
	Dtds                                []*MappingInstructionInput `json:"dtds,omitempty"`
	Fips                                []*MappingInstructionInput `json:"fips,omitempty"`
	Ldaps                               []*MappingInstructionInput `json:"ldaps,omitempty"`
	InternalGroups                      []*MappingInstructionInput `json:"internalGroups,omitempty"`
	FipGroups                           []*MappingInstructionInput `json:"fipGroups,omitempty"`
	InternalUsers                       []*MappingInstructionInput `json:"internalUsers,omitempty"`
	FipUsers                            []*MappingInstructionInput `json:"fipUsers,omitempty"`
	Secrets                             []*MappingInstructionInput `json:"secrets,omitempty"`
	Keys                                []*MappingInstructionInput `json:"keys,omitempty"`
	CassandraConnections                []*MappingInstructionInput `json:"cassandraConnections,omitempty"`
	JmsDestinations                     []*MappingInstructionInput `json:"jmsDestinations,omitempty"`
	GlobalPolicies                      []*MappingInstructionInput `json:"globalPolicies,omitempty"`
	BackgroundTasks                     []*MappingInstructionInput `json:"backgroundTaskPolicies,omitempty"`
	ScheduledTasks                      []*MappingInstructionInput `json:"scheduledTasks,omitempty"`
	ServerModuleFiles                   []*MappingInstructionInput `json:"serverModuleFiles,omitempty"`
	SiteMinderConfigs                   []*MappingInstructionInput `json:"smConfigs,omitempty"`
	ActiveConnectors                    []*MappingInstructionInput `json:"activeConnectors,omitempty"`
	EmailListeners                      []*MappingInstructionInput `json:"emailListeners,omitempty"`
	ListenPorts                         []*MappingInstructionInput `json:"listenPorts,omitempty"`
	AdministrativeUserAccountProperties []*MappingInstructionInput `json:"administrativeUserAccountProperties,omitempty"`
	PasswordPolicies                    []*MappingInstructionInput `json:"passwordPolicies,omitempty"`
	RevocationCheckPolicies             []*MappingInstructionInput `json:"revocationCheckPolicies,omitempty"`
	LogSinks                            []*MappingInstructionInput `json:"logSinks,omitempty"`
	HttpConfigurations                  []*MappingInstructionInput `json:"httpConfigurations,omitempty"`
	CustomKeyValues                     []*MappingInstructionInput `json:"customKeyValues,omitempty"`
	ServiceResolutionConfigs            []*MappingInstructionInput `json:"serviceResolutionConfigs,omitempty"`
	Folders                             []*MappingInstructionInput `json:"folders,omitempty"`
	FederatedIdps                       []*MappingInstructionInput `json:"federatedIdps,omitempty"`
	FederatedGroups                     []*MappingInstructionInput `json:"federatedGroups,omitempty"`
	FederatedUsers                      []*MappingInstructionInput `json:"federatedUsers,omitempty"`
	InternalIdps                        []*MappingInstructionInput `json:"internalIdps,omitempty"`
	LdapIdps                            []*MappingInstructionInput `json:"ldapIdps,omitempty"`
	SimpleLdapIdps                      []*MappingInstructionInput `json:"simpleLdapIdps,omitempty"`
	PolicyBackedIdps                    []*MappingInstructionInput `json:"policyBackedIdps,omitempty"`
	Policies                            []*MappingInstructionInput `json:"policies,omitempty"`
	Services                            []*MappingInstructionInput `json:"services,omitempty"`
	Roles                               []*MappingInstructionInput `json:"roles,omitempty"`
	GenericEntities                     []*MappingInstructionInput `json:"genericEntities,omitempty"`
	AuditConfigurations                 []*MappingInstructionInput `json:"auditConfigurations,omitempty"`
}

type MutationDetailedStatus struct {
	DetailedStatus []DetailedStatus `json:"detailedStatus,omitempty"`
}

type DetailedStatus struct {
	Action      MappingAction  `json:"action,omitempty"`
	Status      MutationStatus `json:"status,omitempty"`
	Description string         `json:"description,omitempty"`
	Source      interface{}    `json:"source,omitempty"`
	Target      interface{}    `json:"target,omitempty"`
}

type BundleResponseDetailedStatus struct {
	WebApiServices                      *MutationDetailedStatus `json:"setWebApiServices,omitempty"`
	InternalWebApiServices              *MutationDetailedStatus `json:"setInternalWebApiServices,omitempty"`
	SoapServices                        *MutationDetailedStatus `json:"setSoapServices,omitempty"`
	InternalSoapServices                *MutationDetailedStatus `json:"setInternalSoapServices,omitempty"`
	PolicyFragments                     *MutationDetailedStatus `json:"setPolicyFragments,omitempty"`
	EncassConfigs                       *MutationDetailedStatus `json:"setEncassConfigs,omitempty"`
	ClusterProperties                   *MutationDetailedStatus `json:"setClusterProperties,omitempty"`
	JdbcConnections                     *MutationDetailedStatus `json:"setJdbcConnections,omitempty"`
	TrustedCerts                        *MutationDetailedStatus `json:"setTrustedCerts,omitempty"`
	Schemas                             *MutationDetailedStatus `json:"setSchemas,omitempty"`
	Dtds                                *MutationDetailedStatus `json:"setDtds,omitempty"`
	Fips                                *MutationDetailedStatus `json:"setFips,omitempty"`
	Ldaps                               *MutationDetailedStatus `json:"setLdaps,omitempty"`
	InternalGroups                      *MutationDetailedStatus `json:"setInternalGroups,omitempty"`
	FipGroups                           *MutationDetailedStatus `json:"setFipGroups,omitempty"`
	InternalUsers                       *MutationDetailedStatus `json:"setInternalUsers,omitempty"`
	FipUsers                            *MutationDetailedStatus `json:"setFipUsers,omitempty"`
	Secrets                             *MutationDetailedStatus `json:"setSecrets,omitempty"`
	Keys                                *MutationDetailedStatus `json:"setKeys,omitempty"`
	CassandraConnections                *MutationDetailedStatus `json:"setCassandraConnections,omitempty"`
	JmsDestinations                     *MutationDetailedStatus `json:"setJmsDestinations,omitempty"`
	GlobalPolicies                      *MutationDetailedStatus `json:"setGlobalPolicies,omitempty"`
	BackgroundTasks                     *MutationDetailedStatus `json:"setBackgroundTaskPolicies,omitempty"`
	ScheduledTasks                      *MutationDetailedStatus `json:"setScheduledTasks,omitempty"`
	ServerModuleFiles                   *MutationDetailedStatus `json:"setServerModuleFiles,omitempty"`
	SiteMinderConfigs                   *MutationDetailedStatus `json:"setSMConfigs,omitempty"`
	ActiveConnectors                    *MutationDetailedStatus `json:"setActiveConnectors,omitempty"`
	EmailListeners                      *MutationDetailedStatus `json:"setEmailListeners,omitempty"`
	ListenPorts                         *MutationDetailedStatus `json:"setListenPorts,omitempty"`
	AdministrativeUserAccountProperties *MutationDetailedStatus `json:"setAdministrativeUserAccountProperties,omitempty"`
	PasswordPolicies                    *MutationDetailedStatus `json:"setPasswordPolicies,omitempty"`
	RevocationCheckPolicies             *MutationDetailedStatus `json:"setRevocationCheckPolicies,omitempty"`
	LogSinks                            *MutationDetailedStatus `json:"setLogSinks,omitempty"`
	HttpConfigurations                  *MutationDetailedStatus `json:"setHttpConfigurations,omitempty"`
	CustomKeyValues                     *MutationDetailedStatus `json:"setCustomKeyValues,omitempty"`
	ServiceResolutionConfigs            *MutationDetailedStatus `json:"setServiceResolutionConfigs,omitempty"`
	Folders                             *MutationDetailedStatus `json:"setFolders,omitempty"`
	FederatedIdps                       *MutationDetailedStatus `json:"setFederatedIdps,omitempty"`
	FederatedGroups                     *MutationDetailedStatus `json:"setFederatedGroups,omitempty"`
	FederatedUsers                      *MutationDetailedStatus `json:"setFederatedUsers,omitempty"`
	InternalIdps                        *MutationDetailedStatus `json:"setInternalIdps,omitempty"`
	LdapIdps                            *MutationDetailedStatus `json:"setLdapIdps,omitempty"`
	SimpleLdapIdps                      *MutationDetailedStatus `json:"setSimpleLdapIdps,omitempty"`
	PolicyBackedIdps                    *MutationDetailedStatus `json:"setPolicyBackedIdps,omitempty"`
	Policies                            *MutationDetailedStatus `json:"setPolicies,omitempty"`
	Services                            *MutationDetailedStatus `json:"setServices,omitempty"`
	Roles                               *MutationDetailedStatus `json:"setRoles,omitempty"`
	GenericEntities                     *MutationDetailedStatus `json:"setGenericEntities,omitempty"`
	AuditConfigurations                 *MutationDetailedStatus `json:"setAuditConfigurations,omitempty"`
}

type MutationError struct {
	Errors []Error `json:"errors,omitempty"`
}

type Error struct {
	Message    string          `json:"message,omitempty"`
	Extensions ErrorExtensions `json:"extensions,omitempty"`
}

type ErrorExtensions struct {
	Classification string `json:"classification,omitempty"`
}

type BundleApplyError struct {
	Entity string         `json:"entity,omitempty"`
	Error  DetailedStatus `json:"error,omitempty"`
}

type MappingSource struct {
	Name           string `json:"name,omitempty"`
	Alias          string `json:"alias,omitempty"`
	KeystoreId     string `json:"keystoreId,omitempty"`
	ThumbprintSha1 string `json:"thumbprintSha1,omitempty"`
	SystemId       string `json:"systemId,omitempty"`
	Port           int    `json:"port,omitempty"`
	Key            string `json:"key,omitempty"`
	ResolutionPath string `json:"resolutionPath,omitempty"`
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
	"roles",
	"genericEntities",
	"auditConfigurations",
}

func ConcatBundle(src []byte, dest []byte) ([]byte, error) {
	srcBundle := Bundle{}
	destBundle := Bundle{}

	if len(src) == 0 {
		return dest, nil
	}

	err := json.Unmarshal(dest, &destBundle)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(src, &srcBundle)
	if err != nil {
		return nil, err
	}

	destBundle = combineBundle(srcBundle, destBundle)

	bundleBytes, err := json.Marshal(destBundle)
	if err != nil {
		return nil, err
	}
	return bundleBytes, nil
}

func AddMappings(src []byte, dest []byte) ([]byte, error) {
	srcBundle := Bundle{}
	destBundle := Bundle{}

	if len(src) == 0 {
		return dest, nil
	}

	err := json.Unmarshal(dest, &destBundle)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(src, &srcBundle)
	if err != nil {
		return nil, err
	}

	destBundle.Properties = srcBundle.Properties

	bundleBytes, err := json.Marshal(destBundle)
	if err != nil {
		return nil, err
	}
	return bundleBytes, nil
}

func matchOptionsLevelFormat(value string) string {
	re := regexp.MustCompile(`{{(.*)}}`)
	match := re.FindStringSubmatch(value)
	if len(match) > 1 {
		return match[1]

	} else {
		re := regexp.MustCompile(`{(.*)}`)
		match = re.FindStringSubmatch(value)
		if len(match) > 1 {
			return match[1]
		}
	}
	return ""
}

func parseCacString(entityType string, configFile string, f interface{}) (string, error) {
	v := reflect.ValueOf(f)
	match := matchOptionsLevelFormat(v.String())
	if match != "" {
		dir, _ := filepath.Split(configFile)
		fBytes, err := os.ReadFile(dir + match)
		if err != nil {
			return "", err
		}
		switch entityType {
		case "trustedCerts":
			return base64.StdEncoding.EncodeToString(fBytes), nil
		case "keys-key":
			return strings.Join(strings.Split(string(fBytes), "\\n"), ""), nil
		case "keys-crt", ".service":
			return string(fBytes), nil
		}

	}

	return "", nil
}

func parsePacCode(p *PolicyInput, file string) (PolicyInput, error) {
	v := reflect.ValueOf(*p)
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).String() != "" {
			if typeOfS.Field(i).Name == "Xml" || typeOfS.Field(i).Name == "Json" || typeOfS.Field(i).Name == "Yaml" {

				match := matchOptionsLevelFormat(v.Field(i).String())
				if match != "" {
					dir, _ := filepath.Split(file)
					fBytes, err := os.ReadFile(dir + match)
					if err != nil {
						return *p, err
					}
					switch typeOfS.Field(i).Name {
					case "Xml":
						p.Xml = strings.Join(strings.Split(string(fBytes), "\\n"), "")
					case "Json":
						p.Json = strings.Join(strings.Split(string(fBytes), "\n"), "")
					case "Yaml":
						p.Yaml = strings.Join(strings.Split(string(fBytes), "\\n"), "")
					}
				}
			}
		}
	}
	return *p, nil
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
	ext := strings.Split(file, ".")[len(strings.Split(file, "."))-1]
	if ext != "json" {
		return *bundle, nil
	}
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
		policy, err := parsePacCode(policyFragment.Policy, file)
		if err != nil {
			return *bundle, err
		}
		policyFragment.Policy = &policy

		bundle.Policies = append(bundle.Policies, &policyFragment)
	case ".service":
		service := L7ServiceInput{}
		err := json.Unmarshal(f, &service)
		if err != nil {
			return *bundle, nil
		}
		if service.ServiceType == "SOAP" {
			wsdl, err := parseCacString(entityType, file, service.Wsdl)
			if err != nil {
				return *bundle, nil
			}
			if wsdl != "" {
				service.Wsdl = wsdl
			}
		}
		policy, err := parsePacCode(service.Policy, file)
		if err != nil {
			return *bundle, err
		}
		service.Policy = &policy
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
		cert, err := parseCacString(entityType, file, trustedCert.CertBase64)
		if err != nil {
			return *bundle, err
		}
		if cert != "" {
			trustedCert.CertBase64 = cert
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
		certChainString, err := parseCacString(entityType+"-crt", file, key.CertChain)
		if err != nil {
			return *bundle, err
		}
		keyField := key.Pem
		if key.P12 != "" {
			keyField = key.P12
		}
		privKey, err := parseCacString(entityType+"-key", file, keyField)
		if err != nil {
			return *bundle, err
		}
		if certChainString != "" {
			certsChain := []string{}
			certStrings := strings.SplitAfter(certChainString, "-----END CERTIFICATE-----")
			for crt := range certStrings {
				if strings.Contains(certStrings[crt], "-----BEGIN CERTIFICATE-----") {
					certsChain = append(certsChain, strings.Join(strings.Split(certStrings[crt], "\n"), "\n"))
				}
			}
			key.CertChain = certsChain
		}
		if privKey != "" {
			if key.P12 != "" {
				key.P12 = privKey
			} else {
				key.Pem = privKey
			}
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
	case "roles":
		role := RoleInput{}
		err := json.Unmarshal(f, &role)
		if err != nil {
			return *bundle, nil
		}
		bundle.Roles = append(bundle.Roles, &role)
	case "genericEntities":
		genericEntity := GenericEntityInput{}
		err := json.Unmarshal(f, &genericEntity)
		if err != nil {
			return *bundle, nil
		}
		bundle.GenericEntities = append(bundle.GenericEntities, &genericEntity)
	case "auditConfigurations":
		auditConfiguration := AuditConfigurationInput{}
		err := json.Unmarshal(f, &auditConfiguration)
		if err != nil {
			return *bundle, nil
		}
		bundle.AuditConfigurations = append(bundle.AuditConfigurations, &auditConfiguration)
	}
	return *bundle, nil
}

func parseBundleProperties(path string) (BundleProperties, error) {
	f, err := os.ReadFile(path + "/bundle-properties.json")
	if err != nil {
		return BundleProperties{}, nil
	}
	bundleProperties := BundleProperties{}
	err = json.Unmarshal(f, &bundleProperties)
	if err != nil {
		return BundleProperties{}, err
	}
	return bundleProperties, nil
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

// SubtractBundle subtracts source from a new bundle by combining the diff and configuring delete mappings
func SubtractBundle(src []byte, new []byte) (delta []byte, combined []byte, err error) {

	var srcBundle Bundle
	var newBundle Bundle
	var deltaBundle Bundle

	deltaBundle.Properties = &BundleProperties{}

	err = json.Unmarshal(src, &srcBundle)
	if err != nil {
		return nil, nil, err
	}

	err = json.Unmarshal(new, &newBundle)
	if err != nil {
		return nil, nil, err
	}

	for _, s := range srcBundle.ActiveConnectors {
		found := false
		for _, d := range newBundle.ActiveConnectors {
			if s.Name == d.Name {
				if !reflect.DeepEqual(s, d) {
					deltaBundle.ActiveConnectors = append(deltaBundle.ActiveConnectors, d)
				}
				s = d
				found = true
			}
		}
		if !found {
			deltaBundle.ActiveConnectors = append(deltaBundle.ActiveConnectors, s)
			newBundle.ActiveConnectors = append(newBundle.ActiveConnectors, s)
			newBundle.Properties.Mappings.ActiveConnectors = append(newBundle.Properties.Mappings.ActiveConnectors, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.BackgroundTasks {
		found := false
		for _, d := range newBundle.BackgroundTasks {
			if s.Name == d.Name {
				if !reflect.DeepEqual(s, d) {
					deltaBundle.BackgroundTasks = append(deltaBundle.BackgroundTasks, d)
				}
				s = d
				found = true
			}
		}
		if !found {
			deltaBundle.BackgroundTasks = append(deltaBundle.BackgroundTasks, s)
			newBundle.BackgroundTasks = append(newBundle.BackgroundTasks, s)
			newBundle.Properties.Mappings.BackgroundTasks = append(newBundle.Properties.Mappings.BackgroundTasks, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.CassandraConnections {
		found := false
		for _, d := range newBundle.CassandraConnections {
			if s.Name == d.Name {
				if !reflect.DeepEqual(s, d) {
					deltaBundle.CassandraConnections = append(deltaBundle.CassandraConnections, d)
				}
				s = d
				found = true
			}
		}
		if !found {
			deltaBundle.CassandraConnections = append(deltaBundle.CassandraConnections, s)
			newBundle.CassandraConnections = append(newBundle.CassandraConnections, s)
			newBundle.Properties.Mappings.CassandraConnections = append(newBundle.Properties.Mappings.CassandraConnections, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.ClusterProperties {
		found := false
		for _, d := range newBundle.ClusterProperties {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.ClusterProperties = append(deltaBundle.ClusterProperties, d)
				}
				s = d

			}
		}
		if !found {
			deltaBundle.ClusterProperties = append(deltaBundle.ClusterProperties, s)
			newBundle.ClusterProperties = append(newBundle.ClusterProperties, s)
			newBundle.Properties.Mappings.ClusterProperties = append(newBundle.Properties.Mappings.ClusterProperties, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.Dtds {
		found := false
		for _, d := range newBundle.Dtds {
			if s.SystemId == d.SystemId {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.Dtds = append(deltaBundle.Dtds, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.Dtds = append(deltaBundle.Dtds, s)
			newBundle.Dtds = append(newBundle.Dtds, s)
			newBundle.Properties.Mappings.Dtds = append(newBundle.Properties.Mappings.Dtds, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					SystemId: s.SystemId,
				},
			})
		}
	}

	for _, s := range srcBundle.EmailListeners {
		found := false
		for _, d := range newBundle.EmailListeners {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.EmailListeners = append(deltaBundle.EmailListeners, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.EmailListeners = append(deltaBundle.EmailListeners, s)
			newBundle.EmailListeners = append(newBundle.EmailListeners, s)
			newBundle.Properties.Mappings.EmailListeners = append(newBundle.Properties.Mappings.EmailListeners, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.EncassConfigs {
		found := false
		for _, d := range newBundle.EncassConfigs {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.EncassConfigs = append(deltaBundle.EncassConfigs, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.EncassConfigs = append(deltaBundle.EncassConfigs, s)
			newBundle.EncassConfigs = append(newBundle.EncassConfigs, s)
			newBundle.Properties.Mappings.EncassConfigs = append(newBundle.Properties.Mappings.EncassConfigs, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.FipGroups {
		found := false
		for _, d := range newBundle.FipGroups {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.FipGroups = append(deltaBundle.FipGroups, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.FipGroups = append(deltaBundle.FipGroups, s)
			newBundle.FipGroups = append(newBundle.FipGroups, s)
			newBundle.Properties.Mappings.FipGroups = append(newBundle.Properties.Mappings.FipGroups, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.FipUsers {
		found := false
		for _, d := range newBundle.FipUsers {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.FipUsers = append(deltaBundle.FipUsers, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.FipUsers = append(deltaBundle.FipUsers, s)
			newBundle.FipUsers = append(newBundle.FipUsers, s)
			newBundle.Properties.Mappings.FipUsers = append(newBundle.Properties.Mappings.FipUsers, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.Fips {
		found := false
		for _, d := range newBundle.Fips {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.Fips = append(deltaBundle.Fips, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.Fips = append(deltaBundle.Fips, s)
			newBundle.Fips = append(newBundle.Fips, s)
			newBundle.Properties.Mappings.Fips = append(newBundle.Properties.Mappings.Fips, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.GlobalPolicies {
		found := false
		for _, d := range newBundle.GlobalPolicies {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.GlobalPolicies = append(deltaBundle.GlobalPolicies, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.GlobalPolicies = append(deltaBundle.GlobalPolicies, s)
			newBundle.GlobalPolicies = append(newBundle.GlobalPolicies, s)
			newBundle.Properties.Mappings.GlobalPolicies = append(newBundle.Properties.Mappings.GlobalPolicies, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.InternalGroups {
		found := false
		for _, d := range newBundle.InternalGroups {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.InternalGroups = append(deltaBundle.InternalGroups, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.InternalGroups = append(deltaBundle.InternalGroups, s)
			newBundle.InternalGroups = append(newBundle.InternalGroups, s)
			newBundle.Properties.Mappings.InternalGroups = append(newBundle.Properties.Mappings.InternalGroups, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.InternalSoapServices {
		found := false
		for _, d := range newBundle.InternalSoapServices {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.InternalSoapServices = append(deltaBundle.InternalSoapServices, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.InternalSoapServices = append(deltaBundle.InternalSoapServices, s)
			newBundle.InternalSoapServices = append(newBundle.InternalSoapServices, s)
			newBundle.Properties.Mappings.InternalSoapServices = append(newBundle.Properties.Mappings.InternalSoapServices, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.InternalUsers {
		found := false
		for _, d := range newBundle.InternalUsers {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.InternalUsers = append(deltaBundle.InternalUsers, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.InternalUsers = append(deltaBundle.InternalUsers, s)
			newBundle.InternalUsers = append(newBundle.InternalUsers, s)
			newBundle.Properties.Mappings.InternalUsers = append(newBundle.Properties.Mappings.InternalUsers, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.InternalWebApiServices {
		found := false
		for _, d := range newBundle.InternalWebApiServices {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.InternalWebApiServices = append(deltaBundle.InternalWebApiServices, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.InternalWebApiServices = append(deltaBundle.InternalWebApiServices, s)
			newBundle.InternalWebApiServices = append(newBundle.InternalWebApiServices, s)
			newBundle.Properties.Mappings.InternalWebApiServices = append(newBundle.Properties.Mappings.InternalWebApiServices, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.JdbcConnections {
		found := false
		for _, d := range newBundle.JdbcConnections {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.JdbcConnections = append(deltaBundle.JdbcConnections, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.JdbcConnections = append(deltaBundle.JdbcConnections, s)
			newBundle.JdbcConnections = append(newBundle.JdbcConnections, s)
			newBundle.Properties.Mappings.JdbcConnections = append(newBundle.Properties.Mappings.JdbcConnections, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.JmsDestinations {
		found := false
		for _, d := range newBundle.JmsDestinations {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.JmsDestinations = append(deltaBundle.JmsDestinations, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.JmsDestinations = append(deltaBundle.JmsDestinations, s)
			newBundle.JmsDestinations = append(newBundle.JmsDestinations, s)
			newBundle.Properties.Mappings.JmsDestinations = append(newBundle.Properties.Mappings.JmsDestinations, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.Keys {
		found := false
		for _, d := range newBundle.Keys {
			if s.Alias == d.Alias {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.Keys = append(deltaBundle.Keys, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.Keys = append(deltaBundle.Keys, s)
			newBundle.Keys = append(newBundle.Keys, s)
			newBundle.Properties.Mappings.Keys = append(newBundle.Properties.Mappings.Keys, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					KeystoreId: "00000000000000000000000000000002",
					Alias:      s.Alias,
				},
			})
		}
	}

	for _, s := range srcBundle.LdapIdps {
		found := false
		for _, d := range newBundle.LdapIdps {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.LdapIdps = append(deltaBundle.LdapIdps, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.LdapIdps = append(deltaBundle.LdapIdps, s)
			newBundle.LdapIdps = append(newBundle.LdapIdps, s)
			newBundle.Properties.Mappings.LdapIdps = append(newBundle.Properties.Mappings.LdapIdps, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.ListenPorts {
		found := false
		for _, d := range newBundle.ListenPorts {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.ListenPorts = append(deltaBundle.ListenPorts, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.ListenPorts = append(deltaBundle.ListenPorts, s)
			newBundle.ListenPorts = append(newBundle.ListenPorts, s)
			newBundle.Properties.Mappings.ListenPorts = append(newBundle.Properties.Mappings.ListenPorts, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.PolicyFragments {
		found := false
		for _, d := range newBundle.PolicyFragments {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.PolicyFragments = append(deltaBundle.PolicyFragments, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.PolicyFragments = append(deltaBundle.PolicyFragments, s)
			newBundle.PolicyFragments = append(newBundle.PolicyFragments, s)
			newBundle.Properties.Mappings.PolicyFragments = append(newBundle.Properties.Mappings.PolicyFragments, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.ScheduledTasks {
		found := false
		for _, d := range newBundle.ScheduledTasks {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.ScheduledTasks = append(deltaBundle.ScheduledTasks, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.ScheduledTasks = append(deltaBundle.ScheduledTasks, s)
			newBundle.ScheduledTasks = append(newBundle.ScheduledTasks, s)
			newBundle.Properties.Mappings.ScheduledTasks = append(newBundle.Properties.Mappings.ScheduledTasks, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.Schemas {
		found := false
		for _, d := range newBundle.Schemas {
			if s.SystemId == d.SystemId {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.Schemas = append(deltaBundle.Schemas, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.Schemas = append(deltaBundle.Schemas, s)
			newBundle.Schemas = append(newBundle.Schemas, s)
			newBundle.Properties.Mappings.Schemas = append(newBundle.Properties.Mappings.Schemas, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					SystemId: s.SystemId,
				},
			})
		}
	}

	for _, s := range srcBundle.Secrets {
		found := false
		for _, d := range newBundle.Secrets {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.Secrets = append(deltaBundle.Secrets, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.Secrets = append(deltaBundle.Secrets, s)
			newBundle.Secrets = append(newBundle.Secrets, s)
			newBundle.Properties.Mappings.Secrets = append(newBundle.Properties.Mappings.Secrets, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.ServerModuleFiles {
		found := false
		for _, d := range newBundle.ServerModuleFiles {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.ServerModuleFiles = append(deltaBundle.ServerModuleFiles, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.ServerModuleFiles = append(deltaBundle.ServerModuleFiles, s)
			newBundle.ServerModuleFiles = append(newBundle.ServerModuleFiles, s)
			newBundle.Properties.Mappings.ServerModuleFiles = append(newBundle.Properties.Mappings.ServerModuleFiles, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.SiteMinderConfigs {
		found := false
		for _, d := range newBundle.SiteMinderConfigs {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.SiteMinderConfigs = append(deltaBundle.SiteMinderConfigs, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.SiteMinderConfigs = append(deltaBundle.SiteMinderConfigs, s)
			newBundle.SiteMinderConfigs = append(newBundle.SiteMinderConfigs, s)
			newBundle.Properties.Mappings.SiteMinderConfigs = append(newBundle.Properties.Mappings.SiteMinderConfigs, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.SoapServices {
		found := false
		for _, d := range newBundle.SoapServices {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.SoapServices = append(deltaBundle.SoapServices, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.SoapServices = append(deltaBundle.SoapServices, s)
			newBundle.SoapServices = append(newBundle.SoapServices, s)
			newBundle.Properties.Mappings.SoapServices = append(newBundle.Properties.Mappings.SoapServices, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.TrustedCerts {
		found := false
		for _, d := range newBundle.TrustedCerts {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.TrustedCerts = append(deltaBundle.TrustedCerts, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.TrustedCerts = append(deltaBundle.TrustedCerts, s)
			newBundle.TrustedCerts = append(newBundle.TrustedCerts, s)
			newBundle.Properties.Mappings.TrustedCerts = append(newBundle.Properties.Mappings.TrustedCerts, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					ThumbprintSha1: s.ThumbprintSha1,
				},
			})
		}
	}

	for _, s := range srcBundle.WebApiServices {
		found := false
		for _, d := range newBundle.WebApiServices {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.WebApiServices = append(deltaBundle.WebApiServices, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.WebApiServices = append(deltaBundle.WebApiServices, s)
			newBundle.WebApiServices = append(newBundle.WebApiServices, s)
			newBundle.Properties.Mappings.WebApiServices = append(newBundle.Properties.Mappings.WebApiServices, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	// AdministrativeUserAccountProperties will be kept if removed as they can't actually be removed
	// future versions may reset to default when not present, this may not represent expected behaviour or cause breaking changes
	// so the default will be to persist even if they aren't present in the new bundle.
	for _, s := range srcBundle.AdministrativeUserAccountProperties {
		found := false
		for _, d := range newBundle.AdministrativeUserAccountProperties {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.AdministrativeUserAccountProperties = append(deltaBundle.AdministrativeUserAccountProperties, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.AdministrativeUserAccountProperties = append(deltaBundle.AdministrativeUserAccountProperties, s)
			newBundle.AdministrativeUserAccountProperties = append(newBundle.AdministrativeUserAccountProperties, s)
		}
	}

	if len(newBundle.PasswordPolicies) == 0 {
		newBundle.PasswordPolicies = append(newBundle.PasswordPolicies, srcBundle.PasswordPolicies...)
	}

	for _, s := range srcBundle.RevocationCheckPolicies {
		found := false
		for _, d := range newBundle.RevocationCheckPolicies {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.RevocationCheckPolicies = append(deltaBundle.RevocationCheckPolicies, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.RevocationCheckPolicies = append(deltaBundle.RevocationCheckPolicies, s)
			newBundle.RevocationCheckPolicies = append(newBundle.RevocationCheckPolicies, s)
			newBundle.Properties.Mappings.RevocationCheckPolicies = append(newBundle.Properties.Mappings.RevocationCheckPolicies, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.LogSinks {
		found := false
		for _, d := range newBundle.LogSinks {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.LogSinks = append(deltaBundle.LogSinks, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.LogSinks = append(deltaBundle.LogSinks, s)
			newBundle.LogSinks = append(newBundle.LogSinks, s)
			newBundle.Properties.Mappings.LogSinks = append(newBundle.Properties.Mappings.LogSinks, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.HttpConfigurations {
		found := false
		for _, d := range newBundle.HttpConfigurations {
			if s.Host == d.Host && s.Port == d.Port {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.HttpConfigurations = append(deltaBundle.HttpConfigurations, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.HttpConfigurations = append(deltaBundle.HttpConfigurations, s)
			newBundle.HttpConfigurations = append(newBundle.HttpConfigurations, s)
			newBundle.Properties.Mappings.HttpConfigurations = append(newBundle.Properties.Mappings.HttpConfigurations, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Host,
					Port: s.Port,
				},
			})
		}
	}

	for _, s := range srcBundle.CustomKeyValues {
		found := false
		for _, d := range newBundle.CustomKeyValues {
			if s.Key == d.Key {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.CustomKeyValues = append(deltaBundle.CustomKeyValues, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.CustomKeyValues = append(deltaBundle.CustomKeyValues, s)
			newBundle.CustomKeyValues = append(newBundle.CustomKeyValues, s)
			newBundle.Properties.Mappings.CustomKeyValues = append(newBundle.Properties.Mappings.CustomKeyValues, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Key: s.Key,
				},
			})
		}
	}

	if len(newBundle.ServiceResolutionConfigs) == 0 {
		newBundle.ServiceResolutionConfigs = append(newBundle.ServiceResolutionConfigs, srcBundle.ServiceResolutionConfigs...)
	}

	for _, s := range srcBundle.Folders {
		found := false
		for _, d := range newBundle.Folders {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.Folders = append(deltaBundle.Folders, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.Folders = append(deltaBundle.Folders, s)
			newBundle.Folders = append(newBundle.Folders, s)
			newBundle.Properties.Mappings.Folders = append(newBundle.Properties.Mappings.Folders, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.FederatedIdps {
		found := false
		for _, d := range newBundle.FederatedIdps {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.FederatedIdps = append(deltaBundle.FederatedIdps, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.FederatedIdps = append(deltaBundle.FederatedIdps, s)
			newBundle.FederatedIdps = append(newBundle.FederatedIdps, s)
			newBundle.Properties.Mappings.FederatedIdps = append(newBundle.Properties.Mappings.FederatedIdps, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.FederatedGroups {
		found := false
		for _, d := range newBundle.FederatedGroups {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.FederatedGroups = append(deltaBundle.FederatedGroups, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.FederatedGroups = append(deltaBundle.FederatedGroups, s)
			newBundle.FederatedGroups = append(newBundle.FederatedGroups, s)
			newBundle.Properties.Mappings.FederatedGroups = append(newBundle.Properties.Mappings.FederatedGroups, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.FederatedUsers {
		found := false
		for _, d := range newBundle.FederatedUsers {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.FederatedUsers = append(deltaBundle.FederatedUsers, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.FederatedUsers = append(deltaBundle.FederatedUsers, s)
			newBundle.FederatedUsers = append(newBundle.FederatedUsers, s)
			newBundle.Properties.Mappings.FederatedUsers = append(newBundle.Properties.Mappings.FederatedUsers, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.FederatedUsers {
		found := false
		for _, d := range newBundle.FederatedUsers {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.FederatedUsers = append(deltaBundle.FederatedUsers, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.FederatedUsers = append(deltaBundle.FederatedUsers, s)
			newBundle.FederatedUsers = append(newBundle.FederatedUsers, s)
			newBundle.Properties.Mappings.FederatedUsers = append(newBundle.Properties.Mappings.FederatedUsers, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.InternalIdps {
		found := false
		for _, d := range newBundle.InternalIdps {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.InternalIdps = append(deltaBundle.InternalIdps, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.InternalIdps = append(deltaBundle.InternalIdps, s)
			newBundle.InternalIdps = append(newBundle.InternalIdps, s)
			newBundle.Properties.Mappings.InternalIdps = append(newBundle.Properties.Mappings.InternalIdps, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.LdapIdps {
		found := false
		for _, d := range newBundle.LdapIdps {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.LdapIdps = append(deltaBundle.LdapIdps, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.LdapIdps = append(deltaBundle.LdapIdps, s)
			newBundle.LdapIdps = append(newBundle.LdapIdps, s)
			newBundle.Properties.Mappings.LdapIdps = append(newBundle.Properties.Mappings.LdapIdps, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.SimpleLdapIdps {
		found := false
		for _, d := range newBundle.SimpleLdapIdps {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.SimpleLdapIdps = append(deltaBundle.SimpleLdapIdps, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.SimpleLdapIdps = append(deltaBundle.SimpleLdapIdps, s)
			newBundle.SimpleLdapIdps = append(newBundle.SimpleLdapIdps, s)
			newBundle.Properties.Mappings.SimpleLdapIdps = append(newBundle.Properties.Mappings.SimpleLdapIdps, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.PolicyBackedIdps {
		found := false
		for _, d := range newBundle.PolicyBackedIdps {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.PolicyBackedIdps = append(deltaBundle.PolicyBackedIdps, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.PolicyBackedIdps = append(deltaBundle.PolicyBackedIdps, s)
			newBundle.PolicyBackedIdps = append(newBundle.PolicyBackedIdps, s)
			newBundle.Properties.Mappings.PolicyBackedIdps = append(newBundle.Properties.Mappings.PolicyBackedIdps, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.Policies {
		found := false
		for _, d := range newBundle.Policies {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.Policies = append(deltaBundle.Policies, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.Policies = append(deltaBundle.Policies, s)
			newBundle.Policies = append(newBundle.Policies, s)
			newBundle.Properties.Mappings.Policies = append(newBundle.Properties.Mappings.Policies, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.Services {
		found := false
		for _, d := range newBundle.Services {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.Services = append(deltaBundle.Services, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.Services = append(deltaBundle.Services, s)
			newBundle.Services = append(newBundle.Services, s)
			newBundle.Properties.Mappings.Services = append(newBundle.Properties.Mappings.Services, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					ResolutionPath: s.ResolutionPath,
				},
			})
		}
	}

	for _, s := range srcBundle.Roles {
		found := false
		for _, d := range newBundle.Roles {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.Roles = append(deltaBundle.Roles, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.Roles = append(deltaBundle.Roles, s)
			newBundle.Roles = append(newBundle.Roles, s)
			newBundle.Properties.Mappings.Roles = append(newBundle.Properties.Mappings.Roles, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.GenericEntities {
		found := false
		for _, d := range newBundle.GenericEntities {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.GenericEntities = append(deltaBundle.GenericEntities, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.GenericEntities = append(deltaBundle.GenericEntities, s)
			newBundle.GenericEntities = append(newBundle.GenericEntities, s)
			newBundle.Properties.Mappings.GenericEntities = append(newBundle.Properties.Mappings.GenericEntities, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	for _, s := range srcBundle.AuditConfigurations {
		found := false
		for _, d := range newBundle.AuditConfigurations {
			if s.Name == d.Name {
				found = true
				if !reflect.DeepEqual(s, d) {
					deltaBundle.AuditConfigurations = append(deltaBundle.AuditConfigurations, d)
				}
				s = d
			}
		}
		if !found {
			deltaBundle.AuditConfigurations = append(deltaBundle.AuditConfigurations, s)
			newBundle.AuditConfigurations = append(newBundle.AuditConfigurations, s)
			newBundle.Properties.Mappings.AuditConfigurations = append(newBundle.Properties.Mappings.AuditConfigurations, &MappingInstructionInput{
				Action: MappingActionDelete,
				Source: MappingSource{
					Name: s.Name,
				},
			})
		}
	}

	deltaBundle.Properties.Mappings = newBundle.Properties.Mappings

	delta, err = json.Marshal(deltaBundle)
	if err != nil {
		return nil, nil, err
	}
	combined, err = json.Marshal(newBundle)
	if err != nil {
		return nil, nil, err
	}

	return delta, combined, nil
}

func combineBundle(srcBundle Bundle, destBundle Bundle) Bundle {
	for _, s := range srcBundle.ActiveConnectors {
		found := false
		for _, d := range destBundle.ActiveConnectors {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.ActiveConnectors = append(destBundle.ActiveConnectors, s)
		}
	}

	for _, s := range srcBundle.BackgroundTasks {
		found := false
		for _, d := range destBundle.BackgroundTasks {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.BackgroundTasks = append(destBundle.BackgroundTasks, s)
		}
	}

	for _, s := range srcBundle.CassandraConnections {
		found := false
		for _, d := range destBundle.CassandraConnections {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.CassandraConnections = append(destBundle.CassandraConnections, s)
		}
	}

	for _, s := range srcBundle.ClusterProperties {
		found := false
		for _, d := range destBundle.ClusterProperties {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.ClusterProperties = append(destBundle.ClusterProperties, s)
		}
	}

	for _, s := range srcBundle.Dtds {
		found := false
		for _, d := range destBundle.Dtds {
			if s.PublicId == d.PublicId {
				found = true
			}
		}
		if !found {
			destBundle.Dtds = append(destBundle.Dtds, s)
		}
	}

	for _, s := range srcBundle.EmailListeners {
		found := false
		for _, d := range destBundle.EmailListeners {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.EmailListeners = append(destBundle.EmailListeners, s)
		}
	}

	for _, s := range srcBundle.EncassConfigs {
		found := false
		for _, d := range destBundle.EncassConfigs {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.EncassConfigs = append(destBundle.EncassConfigs, s)
		}
	}

	for _, s := range srcBundle.FipGroups {
		found := false
		for _, d := range destBundle.FipGroups {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.FipGroups = append(destBundle.FipGroups, s)
		}
	}

	for _, s := range srcBundle.FipUsers {
		found := false
		for _, d := range destBundle.FipUsers {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.FipUsers = append(destBundle.FipUsers, s)
		}
	}

	for _, s := range srcBundle.Fips {
		found := false
		for _, d := range destBundle.Fips {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.Fips = append(destBundle.Fips, s)
		}
	}

	for _, s := range srcBundle.GlobalPolicies {
		found := false
		for _, d := range destBundle.GlobalPolicies {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.GlobalPolicies = append(destBundle.GlobalPolicies, s)
		}
	}

	for _, s := range srcBundle.InternalGroups {
		found := false
		for _, d := range destBundle.InternalGroups {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.InternalGroups = append(destBundle.InternalGroups, s)
		}
	}

	for _, s := range srcBundle.InternalSoapServices {
		found := false
		for _, d := range destBundle.InternalSoapServices {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.InternalSoapServices = append(destBundle.InternalSoapServices, s)
		}
	}

	for _, s := range srcBundle.InternalUsers {
		found := false
		for _, d := range destBundle.InternalUsers {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.InternalUsers = append(destBundle.InternalUsers, s)
		}
	}

	for _, s := range srcBundle.InternalWebApiServices {
		found := false
		for _, d := range destBundle.InternalWebApiServices {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.InternalWebApiServices = append(destBundle.InternalWebApiServices, s)
		}
	}

	for _, s := range srcBundle.JdbcConnections {
		found := false
		for _, d := range destBundle.JdbcConnections {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.JdbcConnections = append(destBundle.JdbcConnections, s)
		}
	}

	for _, s := range srcBundle.JmsDestinations {
		found := false
		for _, d := range destBundle.JmsDestinations {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.JmsDestinations = append(destBundle.JmsDestinations, s)
		}
	}

	for _, s := range srcBundle.Keys {
		found := false
		for _, d := range destBundle.Keys {
			if s.Alias == d.Alias {
				found = true
			}
		}
		if !found {
			destBundle.Keys = append(destBundle.Keys, s)
		}
	}

	for _, s := range srcBundle.LdapIdps {
		found := false
		for _, d := range destBundle.LdapIdps {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.LdapIdps = append(destBundle.LdapIdps, s)
		}
	}

	for _, s := range srcBundle.ListenPorts {
		found := false
		for _, d := range destBundle.ListenPorts {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.ListenPorts = append(destBundle.ListenPorts, s)
		}
	}

	for _, s := range srcBundle.PolicyFragments {
		found := false
		for _, d := range destBundle.PolicyFragments {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.PolicyFragments = append(destBundle.PolicyFragments, s)
		}
	}

	for _, s := range srcBundle.ScheduledTasks {
		found := false
		for _, d := range destBundle.ScheduledTasks {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.ScheduledTasks = append(destBundle.ScheduledTasks, s)
		}
	}

	for _, s := range srcBundle.Schemas {
		found := false
		for _, d := range destBundle.Schemas {
			if s.SystemId == d.SystemId {
				found = true
			}
		}
		if !found {
			destBundle.Schemas = append(destBundle.Schemas, s)
		}
	}

	for _, s := range srcBundle.Secrets {
		found := false
		for _, d := range destBundle.Secrets {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.Secrets = append(destBundle.Secrets, s)
		}
	}

	for _, s := range srcBundle.ServerModuleFiles {
		found := false
		for _, d := range destBundle.ServerModuleFiles {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.ServerModuleFiles = append(destBundle.ServerModuleFiles, s)
		}
	}

	for _, s := range srcBundle.SiteMinderConfigs {
		found := false
		for _, d := range destBundle.SiteMinderConfigs {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.SiteMinderConfigs = append(destBundle.SiteMinderConfigs, s)
		}
	}

	for _, s := range srcBundle.SoapServices {
		found := false
		for _, d := range destBundle.SoapServices {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.SoapServices = append(destBundle.SoapServices, s)
		}
	}

	for _, s := range srcBundle.TrustedCerts {
		found := false
		for _, d := range destBundle.TrustedCerts {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.TrustedCerts = append(destBundle.TrustedCerts, s)
		}
	}

	for _, s := range srcBundle.TrustedCerts {
		found := false
		for _, d := range destBundle.TrustedCerts {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.TrustedCerts = append(destBundle.TrustedCerts, s)
		}
	}

	for _, s := range srcBundle.WebApiServices {
		found := false
		for _, d := range destBundle.WebApiServices {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.WebApiServices = append(destBundle.WebApiServices, s)
		}
	}

	for _, s := range srcBundle.AdministrativeUserAccountProperties {
		found := false
		for _, d := range destBundle.AdministrativeUserAccountProperties {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.AdministrativeUserAccountProperties = append(destBundle.AdministrativeUserAccountProperties, s)
		}
	}

	if len(destBundle.PasswordPolicies) == 0 {
		destBundle.PasswordPolicies = append(destBundle.PasswordPolicies, srcBundle.PasswordPolicies...)
	}

	for _, s := range srcBundle.RevocationCheckPolicies {
		found := false
		for _, d := range destBundle.RevocationCheckPolicies {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.RevocationCheckPolicies = append(destBundle.RevocationCheckPolicies, s)
		}
	}

	for _, s := range srcBundle.LogSinks {
		found := false
		for _, d := range destBundle.LogSinks {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.LogSinks = append(destBundle.LogSinks, s)
		}
	}

	for _, s := range srcBundle.HttpConfigurations {
		found := false
		for _, d := range destBundle.HttpConfigurations {
			if s.Host == d.Host {
				found = true
			}
		}
		if !found {
			destBundle.HttpConfigurations = append(destBundle.HttpConfigurations, s)
		}
	}

	for _, s := range srcBundle.CustomKeyValues {
		found := false
		for _, d := range destBundle.CustomKeyValues {
			if s.Key == d.Key {
				found = true
			}
		}
		if !found {
			destBundle.CustomKeyValues = append(destBundle.CustomKeyValues, s)
		}
	}

	if len(destBundle.ServiceResolutionConfigs) == 0 {
		destBundle.ServiceResolutionConfigs = append(destBundle.ServiceResolutionConfigs, srcBundle.ServiceResolutionConfigs...)
	}

	for _, s := range srcBundle.Folders {
		found := false
		for _, d := range destBundle.Folders {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.Folders = append(destBundle.Folders, s)
		}
	}

	for _, s := range srcBundle.FederatedIdps {
		found := false
		for _, d := range destBundle.FederatedIdps {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.FederatedIdps = append(destBundle.FederatedIdps, s)
		}
	}

	for _, s := range srcBundle.FederatedGroups {
		found := false
		for _, d := range destBundle.FederatedGroups {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.FederatedGroups = append(destBundle.FederatedGroups, s)
		}
	}

	for _, s := range srcBundle.FederatedUsers {
		found := false
		for _, d := range destBundle.FederatedUsers {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.FederatedUsers = append(destBundle.FederatedUsers, s)
		}
	}

	for _, s := range srcBundle.FederatedUsers {
		found := false
		for _, d := range destBundle.FederatedUsers {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.FederatedUsers = append(destBundle.FederatedUsers, s)
		}
	}

	for _, s := range srcBundle.InternalIdps {
		found := false
		for _, d := range destBundle.InternalIdps {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.InternalIdps = append(destBundle.InternalIdps, s)
		}
	}

	for _, s := range srcBundle.LdapIdps {
		found := false
		for _, d := range destBundle.LdapIdps {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.LdapIdps = append(destBundle.LdapIdps, s)
		}
	}

	for _, s := range srcBundle.SimpleLdapIdps {
		found := false
		for _, d := range destBundle.SimpleLdapIdps {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.SimpleLdapIdps = append(destBundle.SimpleLdapIdps, s)
		}
	}

	for _, s := range srcBundle.PolicyBackedIdps {
		found := false
		for _, d := range destBundle.PolicyBackedIdps {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.PolicyBackedIdps = append(destBundle.PolicyBackedIdps, s)
		}
	}

	for _, s := range srcBundle.Policies {
		found := false
		for _, d := range destBundle.Policies {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.Policies = append(destBundle.Policies, s)
		}
	}

	for _, s := range srcBundle.Services {
		found := false
		for _, d := range destBundle.Services {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.Services = append(destBundle.Services, s)
		}
	}

	for _, s := range srcBundle.Roles {
		found := false
		for _, d := range destBundle.Roles {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.Roles = append(destBundle.Roles, s)
		}
	}

	for _, s := range srcBundle.GenericEntities {
		found := false
		for _, d := range destBundle.GenericEntities {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.GenericEntities = append(destBundle.GenericEntities, s)
		}
	}

	for _, s := range srcBundle.AuditConfigurations {
		found := false
		for _, d := range destBundle.AuditConfigurations {
			if s.Name == d.Name {
				found = true
			}
		}
		if !found {
			destBundle.AuditConfigurations = append(destBundle.AuditConfigurations, s)
		}
	}

	return destBundle
}

func installGenericBundle(
	ctx_ context.Context,
	client_ graphql.Client,
	bundle *Bundle,
) (interface{}, error) {
	req_ := &graphql.Request{
		OpName:    "installBundle",
		Query:     installBundleGeneric_Operation,
		Variables: bundle,
	}
	var err_ error

	var data_ BundleResponseDetailedStatus
	resp_ := &graphql.Response{Data: &data_}

	err_ = client_.MakeRequest(
		ctx_,
		req_,
		resp_,
	)
	return &data_, err_
}

func deleteGenericBundle(
	ctx_ context.Context,
	client_ graphql.Client,
	bundle *Bundle,
) (interface{}, error) {
	req_ := &graphql.Request{
		OpName:    "deleteBundle",
		Query:     deleteBundleGeneric_Operation,
		Variables: bundle,
	}
	var err_ error

	var data_ BundleResponseDetailedStatus
	resp_ := &graphql.Response{Data: &data_}

	err_ = client_.MakeRequest(
		ctx_,
		req_,
		resp_,
	)
	return &data_, err_
}
