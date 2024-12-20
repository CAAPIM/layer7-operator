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
