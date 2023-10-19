package graphman

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type Bundle struct {
	WebApiServices         []*WebApiServiceInput        `json:"webApiServices,omitempty"`
	InternalWebApiServices []*WebApiServiceInput        `json:"internalWebApiServices,omitempty"`
	SoapServices           []*SoapServiceInput          `json:"soapServices,omitempty"`
	InternalSoapServices   []*SoapServiceInput          `json:"internalSoapServices,omitempty"`
	PolicyFragments        []*PolicyFragmentInput       `json:"policyFragments,omitempty"`
	EncassConfigs          []*EncassConfigInput         `json:"encassConfigs,omitempty"`
	ClusterProperties      []*ClusterPropertyInput      `json:"clusterProperties,omitempty"`
	JdbcConnections        []*JdbcConnectionInput       `json:"jdbcConnections,omitempty"`
	TrustedCerts           []*TrustedCertInput          `json:"trustedCerts,omitempty"`
	Schemas                []*SchemaInput               `json:"schemas,omitempty"`
	Dtds                   []*DtdInput                  `json:"dtds,omitempty"`
	Fips                   []*FipInput                  `json:"fips,omitempty"`
	LdapIdps               []*LdapInput                 `json:"ldaps,omitempty"`
	InternalGroups         []*InternalGroupInput        `json:"internalGroups,omitempty"`
	FipGroups              []*FipGroupInput             `json:"fipGroups,omitempty"`
	InternalUsers          []*InternalUserInput         `json:"internalUsers,omitempty"`
	FipUsers               []*FipUserInput              `json:"fipUsers,omitempty"`
	Secrets                []*SecretInput               `json:"secrets,omitempty"`
	Keys                   []*KeyInput                  `json:"keys,omitempty"`
	CassandraConnections   []*CassandraConnectionInput  `json:"cassandraConnections,omitempty"`
	JmsDestinations        []*JmsDestinationInput       `json:"jmsDestinations,omitempty"`
	GlobalPolicies         []*GlobalPolicyInput         `json:"globalPolicies,omitempty"`
	BackgroundTasks        []*BackgroundTaskPolicyInput `json:"backgroundTaskPolicies,omitempty"`
	ScheduledTasks         []*ScheduledTaskInput        `json:"scheduledTasks,omitempty"`
	ServerModuleFiles      []*ServerModuleFileInput     `json:"serverModuleFiles,omitempty"`
	SiteMinderConfigs      []*SMConfigInput             `json:"smConfigs,omitempty"`
	ActiveConnectors       []*ActiveConnectorInput      `json:"activeConnectors,omitempty"`
	EmailListeners         []*EmailListenerInput        `json:"emailListeners,omitempty"`
	ListenPorts            []*ListenPortInput           `json:"listenPorts,omitempty"`
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
	".internalwebapi",
	".webapi",
	".soap",
	".internalsoap",
	".global",
	".policy",
	".bgpolicy",
}

type BundleApplyErrors struct {
	Errors []BundleApplyError `json:"errors,omitempty"`
}

type BundleApplyError struct {
	Status string `json:"status,omitempty"`
	Entity string `json:"entity,omitempty"`
	Detail string `json:"detail,omitempty"`
}

func Query(username string, password string, target string, encpass string) ([]byte, error) {
	resp, err := everything(context.Background(), gqlClient(username, password, target, encpass))
	if err != nil {
		return nil, err
	}
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return respBytes, nil
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

	// copy(destBundle..., srcBundle...)
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

// Apply - applies a bundle to a target gateway.
func Apply(path string, username string, password string, target string, encpass string) ([]byte, error) {
	bundle, err := implodeBundle(path)
	if err != nil {
		return nil, err
	}

	resp, err := applyBundle(context.Background(), gqlClient(username, password, target, encpass), bundle.ClusterProperties, bundle.WebApiServices, bundle.EncassConfigs, bundle.TrustedCerts, bundle.Dtds, bundle.Schemas, bundle.JdbcConnections, bundle.SoapServices, bundle.PolicyFragments, bundle.Fips, bundle.LdapIdps, bundle.FipGroups, bundle.InternalGroups, bundle.FipUsers, bundle.InternalUsers, bundle.Keys, bundle.Secrets, bundle.CassandraConnections, bundle.JmsDestinations, bundle.InternalWebApiServices, bundle.InternalSoapServices, bundle.EmailListeners, bundle.ListenPorts, bundle.ActiveConnectors, bundle.SiteMinderConfigs, bundle.GlobalPolicies, bundle.BackgroundTasks, bundle.ScheduledTasks, bundle.ServerModuleFiles)
	if err != nil {
		return nil, err
	}

	err = CheckApplyErrors(resp)
	if err != nil {
		return nil, err
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return respBytes, nil
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

	resp, err := applyBundle(context.Background(), gqlClient(username, password, target, encpass), bundle.ClusterProperties, bundle.WebApiServices, bundle.EncassConfigs, bundle.TrustedCerts, bundle.Dtds, bundle.Schemas, bundle.JdbcConnections, bundle.SoapServices, bundle.PolicyFragments, bundle.Fips, bundle.LdapIdps, bundle.FipGroups, bundle.InternalGroups, bundle.FipUsers, bundle.InternalUsers, bundle.Keys, bundle.Secrets, bundle.CassandraConnections, bundle.JmsDestinations, bundle.InternalWebApiServices, bundle.InternalSoapServices, bundle.EmailListeners, bundle.ListenPorts, bundle.ActiveConnectors, bundle.SiteMinderConfigs, bundle.GlobalPolicies, bundle.BackgroundTasks, bundle.ScheduledTasks, bundle.ServerModuleFiles)
	if err != nil {
		return nil, err
	}

	err = CheckApplyErrors(resp)
	if err != nil {
		return nil, err
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return respBytes, nil
}

// parseEntities determines which entity the file from a Graphman directory belongs to
// this works with a static list of globally defined entities
func parseEntity(path string) (string, bool) {
	for _, e := range entities {
		if strings.Contains(path, e) {
			return e, true
		}
	}

	return "", false
}

// Read bundle unmarshals a JSON file in the specified Graphman directory into the working Bundle object.
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
		policyFragment := PolicyFragmentInput{}
		err := json.Unmarshal(f, &policyFragment)
		if err != nil {
			return *bundle, err
		}
		bundle.PolicyFragments = append(bundle.PolicyFragments, &policyFragment)
	case ".bgpolicy":
		backgroundTask := BackgroundTaskPolicyInput{}
		err := json.Unmarshal(f, &backgroundTask)
		if err != nil {
			return *bundle, err
		}
		bundle.BackgroundTasks = append(bundle.BackgroundTasks, &backgroundTask)
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
		ldapIdp := LdapInput{}
		err := json.Unmarshal(f, &ldapIdp)
		if err != nil {
			return *bundle, err
		}
		bundle.LdapIdps = append(bundle.LdapIdps, &ldapIdp)
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
func parseEntities(bundle Bundle) {
	v := reflect.ValueOf(bundle)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fmt.Printf("%s %v\n", typeOfS.Field(i).Name, v.Field(i).Interface())
	}
}

func CheckApplyErrors(resp *applyBundleResponse) error {
	var bundleApplyErrors BundleApplyErrors

	for _, r := range resp.SetActiveConnectors.DetailedStatus {
		if r.Status == "ERROR" {
			bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "ActiveConnectors"})
		}
	}
	for _, r := range resp.SetBackgroundTaskPolicies.DetailedStatus {
		if r.Status == "ERROR" {
			bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "BackgroundTaskPolicies"})
		}
	}
	for _, r := range resp.SetCassandraConnections.DetailedStatus {
		if r.Status == "ERROR" {
			bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "CassandraConnections"})
		}
	}
	for _, r := range resp.SetClusterProperties.DetailedStatus {
		if r.Status == "ERROR" {
			bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "ClusterProperties"})
		}
	}
	for _, r := range resp.SetDtds.DetailedStatus {
		if r.Status == "ERROR" {
			bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "Dtds"})
		}
	}
	for _, r := range resp.SetEmailListeners.DetailedStatus {
		if r.Status == "ERROR" {
			bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "EmailListeners"})
		}
	}
	for _, r := range resp.SetEncassConfigs.DetailedStatus {
		if r.Status == "ERROR" {
			bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "EncassConfigs"})
		}
	}
	for _, r := range resp.SetFipGroups.DetailedStatus {
		if r.Status == "ERROR" {
			bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "FipGroups"})
		}
	}
	for _, r := range resp.SetFipUsers.DetailedStatus {
		if r.Status == "ERROR" {
			bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "FipUsers"})
		}
	}
	for _, r := range resp.SetFips.DetailedStatus {
		if r.Status == "ERROR" {
			bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "Fips"})
		}
	}
	for _, r := range resp.SetGlobalPolicies.DetailedStatus {
		if r.Status == "ERROR" {
			bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "GlobalPolicies"})
		}
	}
	for _, r := range resp.SetInternalGroups.DetailedStatus {
		if r.Status == "ERROR" {
			bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "InternalGroups"})
		}
	}
	for _, r := range resp.SetInternalSoapServices.DetailedStatus {
		if r.Status == "ERROR" {
			bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "InternalSoapServices"})
		}
	}
	for _, r := range resp.SetInternalUsers.DetailedStatus {
		if r.Status == "ERROR" {
			bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "InternalUsers"})
		}
	}
	for _, r := range resp.SetInternalWebApiServices.DetailedStatus {
		if r.Status == "ERROR" {
			bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "InternalWebApiServices"})
		}
	}
	for _, r := range resp.SetJdbcConnections.DetailedStatus {
		if r.Status == "ERROR" {
			bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "JdbcConnections"})
		}
	}
	for _, r := range resp.SetJmsDestinations.DetailedStatus {
		if r.Status == "ERROR" {
			bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "JmsDestinations"})
		}
	}
	for _, r := range resp.SetKeys.DetailedStatus {
		if r.Status == "ERROR" {
			bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "Keys"})
		}
	}
	for _, r := range resp.SetLdaps.DetailedStatus {
		if r.Status == "ERROR" {
			bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "Ldaps"})
		}
	}
	for _, r := range resp.SetListenPorts.DetailedStatus {
		if r.Status == "ERROR" {
			bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "ListenPorts"})
		}
	}
	for _, r := range resp.SetPolicyFragments.DetailedStatus {
		if r.Status == "ERROR" {
			bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "PolicyFragments"})
		}
	}
	for _, r := range resp.SetSMConfigs.DetailedStatus {
		if r.Status == "ERROR" {
			bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "SMConfigs"})
		}
	}
	for _, r := range resp.SetScheduledTasks.DetailedStatus {
		if r.Status == "ERROR" {
			bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "ScheduledTasks"})
		}
	}
	for _, r := range resp.SetSchemas.DetailedStatus {
		if r.Status == "ERROR" {
			bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "Schemas"})
		}
	}
	for _, r := range resp.SetSecrets.DetailedStatus {
		if r.Status == "ERROR" {
			bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "Secrets"})
		}
	}
	for _, r := range resp.SetServerModuleFiles.DetailedStatus {
		if r.Status == "ERROR" {
			bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "ServerModuleFiles"})
		}
	}
	for _, r := range resp.SetSoapServices.DetailedStatus {
		if r.Status == "ERROR" {
			bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "SoapServices"})
		}
	}
	for _, r := range resp.SetTrustedCerts.DetailedStatus {
		if r.Status == "ERROR" {
			bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "TrustedCerts"})
		}
	}
	for _, r := range resp.SetWebApiServices.DetailedStatus {
		if r.Status == "ERROR" {
			bundleApplyErrors.Errors = append(bundleApplyErrors.Errors, BundleApplyError{Status: string(r.Status), Detail: r.Description, Entity: "WebApiServices"})
		}
	}

	if len(bundleApplyErrors.Errors) > 0 {
		errorBytes, _ := json.Marshal(bundleApplyErrors)
		return fmt.Errorf("errors: %s", string(errorBytes))
	}

	return nil
}
