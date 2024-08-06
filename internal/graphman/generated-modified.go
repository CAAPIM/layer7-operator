package graphman

import (
	"context"
	"time"

	"github.com/Khan/genqlient/graphql"
)

type ActiveConnectorInput struct {
	// The internal entity unique identifier
	Goid string `json:"goid"`
	// The active connector name
	Name string `json:"name"`
	// Whether this active connector is enabled
	Enabled bool `json:"enabled"`
	// The active connector type Examples:- KAFKA, SFTP_POLLING_LISTENER, MQ_NATIVE
	ConnectorType string `json:"connectorType"`
	// The name of the published service hardwired to the active connector
	HardwiredServiceName string `json:"hardwiredServiceName"`
	// The active connector properties
	Properties []*EntityPropertyInput `json:"properties,omitempty"`
	// The advanced properties for active connector
	AdvancedProperties []*EntityPropertyInput `json:"advancedProperties,omitempty"`
	// Ignored at creation time but can be used to compare bundle with gw state
	Checksum string `json:"checksum"`
}

// GetGoid returns ActiveConnectorInput.Goid, and is useful for accessing the field via an interface.
func (v *ActiveConnectorInput) GetGoid() string { return v.Goid }

// GetName returns ActiveConnectorInput.Name, and is useful for accessing the field via an interface.
func (v *ActiveConnectorInput) GetName() string { return v.Name }

// GetEnabled returns ActiveConnectorInput.Enabled, and is useful for accessing the field via an interface.
func (v *ActiveConnectorInput) GetEnabled() bool { return v.Enabled }

// GetConnectorType returns ActiveConnectorInput.ConnectorType, and is useful for accessing the field via an interface.
func (v *ActiveConnectorInput) GetConnectorType() string { return v.ConnectorType }

// GetHardwiredServiceName returns ActiveConnectorInput.HardwiredServiceName, and is useful for accessing the field via an interface.
func (v *ActiveConnectorInput) GetHardwiredServiceName() string { return v.HardwiredServiceName }

// GetProperties returns ActiveConnectorInput.Properties, and is useful for accessing the field via an interface.
func (v *ActiveConnectorInput) GetProperties() []*EntityPropertyInput { return v.Properties }

// GetAdvancedProperties returns ActiveConnectorInput.AdvancedProperties, and is useful for accessing the field via an interface.
func (v *ActiveConnectorInput) GetAdvancedProperties() []*EntityPropertyInput {
	return v.AdvancedProperties
}

// GetChecksum returns ActiveConnectorInput.Checksum, and is useful for accessing the field via an interface.
func (v *ActiveConnectorInput) GetChecksum() string { return v.Checksum }

// The inputs sent with the setClusterProperty Mutation
type AdministrativeUserAccountPropertyInput struct {
	// The administrative user account property unique identifier
	Goid string `json:"goid"`
	// The name of administrative user account property
	Name string `json:"name"`
	// The value of the administrative user account property
	Value string `json:"value"`
	// Ignored at creation time but can be used to compare bundle with gw state
	Checksum string `json:"checksum"`
}

// GetGoid returns AdministrativeUserAccountPropertyInput.Goid, and is useful for accessing the field via an interface.
func (v *AdministrativeUserAccountPropertyInput) GetGoid() string { return v.Goid }

// GetName returns AdministrativeUserAccountPropertyInput.Name, and is useful for accessing the field via an interface.
func (v *AdministrativeUserAccountPropertyInput) GetName() string { return v.Name }

// GetValue returns AdministrativeUserAccountPropertyInput.Value, and is useful for accessing the field via an interface.
func (v *AdministrativeUserAccountPropertyInput) GetValue() string { return v.Value }

// GetChecksum returns AdministrativeUserAccountPropertyInput.Checksum, and is useful for accessing the field via an interface.
func (v *AdministrativeUserAccountPropertyInput) GetChecksum() string { return v.Checksum }

type AuditConfigurationInput struct {
	// The internal entity unique identifier
	Goid string `json:"goid"`
	// log sink unique name
	Name string `json:"name"`
	// Lookup Policy Name
	LookupPolicyName string `json:"lookupPolicyName"`
	// The configuration checksum
	Checksum string `json:"checksum"`
	// Sink Policy Name
	SinkPolicyName string `json:"sinkPolicyName"`
	// FTP Client Configuration
	FtpConfig *AuditFtpConfigurationInput `json:"ftpConfig,omitempty"`
}

// GetGoid returns AuditConfigurationInput.Goid, and is useful for accessing the field via an interface.
func (v *AuditConfigurationInput) GetGoid() string { return v.Goid }

// GetName returns AuditConfigurationInput.Name, and is useful for accessing the field via an interface.
func (v *AuditConfigurationInput) GetName() string { return v.Name }

// GetLookupPolicyName returns AuditConfigurationInput.LookupPolicyName, and is useful for accessing the field via an interface.
func (v *AuditConfigurationInput) GetLookupPolicyName() string { return v.LookupPolicyName }

// GetChecksum returns AuditConfigurationInput.Checksum, and is useful for accessing the field via an interface.
func (v *AuditConfigurationInput) GetChecksum() string { return v.Checksum }

// GetSinkPolicyName returns AuditConfigurationInput.SinkPolicyName, and is useful for accessing the field via an interface.
func (v *AuditConfigurationInput) GetSinkPolicyName() string { return v.SinkPolicyName }

// GetFtpConfig returns AuditConfigurationInput.FtpConfig, and is useful for accessing the field via an interface.
func (v *AuditConfigurationInput) GetFtpConfig() *AuditFtpConfigurationInput { return v.FtpConfig }

type AuditFtpConfigurationInput struct {
	// Host of FTP Server
	Host string `json:"host"`
	// Port of FTP Server
	Port int `json:"port"`
	// FTP connection timeout
	Timeout int `json:"timeout"`
	// FTP user
	User string `json:"user"`
	// FTP password
	Password string `json:"password"`
	// Directory in FTP Server
	Directory string `json:"directory"`
	// To verify server certification
	VerifyServerCert bool `json:"verifyServerCert"`
	// Security Type
	Security AuditFtpSecurityType `json:"security"`
	// Whether this Audit Configuration is enabled
	Enabled bool `json:"enabled"`
}

// GetHost returns AuditFtpConfigurationInput.Host, and is useful for accessing the field via an interface.
func (v *AuditFtpConfigurationInput) GetHost() string { return v.Host }

// GetPort returns AuditFtpConfigurationInput.Port, and is useful for accessing the field via an interface.
func (v *AuditFtpConfigurationInput) GetPort() int { return v.Port }

// GetTimeout returns AuditFtpConfigurationInput.Timeout, and is useful for accessing the field via an interface.
func (v *AuditFtpConfigurationInput) GetTimeout() int { return v.Timeout }

// GetUser returns AuditFtpConfigurationInput.User, and is useful for accessing the field via an interface.
func (v *AuditFtpConfigurationInput) GetUser() string { return v.User }

// GetPassword returns AuditFtpConfigurationInput.Password, and is useful for accessing the field via an interface.
func (v *AuditFtpConfigurationInput) GetPassword() string { return v.Password }

// GetDirectory returns AuditFtpConfigurationInput.Directory, and is useful for accessing the field via an interface.
func (v *AuditFtpConfigurationInput) GetDirectory() string { return v.Directory }

// GetVerifyServerCert returns AuditFtpConfigurationInput.VerifyServerCert, and is useful for accessing the field via an interface.
func (v *AuditFtpConfigurationInput) GetVerifyServerCert() bool { return v.VerifyServerCert }

// GetSecurity returns AuditFtpConfigurationInput.Security, and is useful for accessing the field via an interface.
func (v *AuditFtpConfigurationInput) GetSecurity() AuditFtpSecurityType { return v.Security }

// GetEnabled returns AuditFtpConfigurationInput.Enabled, and is useful for accessing the field via an interface.
func (v *AuditFtpConfigurationInput) GetEnabled() bool { return v.Enabled }

// Indicates the Sink Category
type AuditFtpSecurityType string

const (
	AuditFtpSecurityTypeFtpUnsecured AuditFtpSecurityType = "FTP_UNSECURED"
	AuditFtpSecurityTypeFtpsExplicit AuditFtpSecurityType = "FTPS_EXPLICIT"
	AuditFtpSecurityTypeFtpsImplicit AuditFtpSecurityType = "FTPS_IMPLICIT"
)

type BackgroundTaskPolicyInput struct {
	// The internal entity unique identifier
	Goid string `json:"goid"`
	// The name of the background task policy
	Name string `json:"name"`
	// The internal entity unique identifier
	Guid string `json:"guid"`
	// The folder path background task policy
	FolderPath string `json:"folderPath"`
	// The background task policy
	Policy *PolicyInput `json:"policy,omitempty"`
	Soap   bool         `json:"soap"`
	// The configuration checksum
	Checksum string `json:"checksum"`
}

// GetGoid returns BackgroundTaskPolicyInput.Goid, and is useful for accessing the field via an interface.
func (v *BackgroundTaskPolicyInput) GetGoid() string { return v.Goid }

// GetName returns BackgroundTaskPolicyInput.Name, and is useful for accessing the field via an interface.
func (v *BackgroundTaskPolicyInput) GetName() string { return v.Name }

// GetGuid returns BackgroundTaskPolicyInput.Guid, and is useful for accessing the field via an interface.
func (v *BackgroundTaskPolicyInput) GetGuid() string { return v.Guid }

// GetFolderPath returns BackgroundTaskPolicyInput.FolderPath, and is useful for accessing the field via an interface.
func (v *BackgroundTaskPolicyInput) GetFolderPath() string { return v.FolderPath }

// GetPolicy returns BackgroundTaskPolicyInput.Policy, and is useful for accessing the field via an interface.
func (v *BackgroundTaskPolicyInput) GetPolicy() *PolicyInput { return v.Policy }

// GetSoap returns BackgroundTaskPolicyInput.Soap, and is useful for accessing the field via an interface.
func (v *BackgroundTaskPolicyInput) GetSoap() bool { return v.Soap }

// GetChecksum returns BackgroundTaskPolicyInput.Checksum, and is useful for accessing the field via an interface.
func (v *BackgroundTaskPolicyInput) GetChecksum() string { return v.Checksum }

type CassandraCompression string

const (
	CassandraCompressionNone CassandraCompression = "NONE"
	CassandraCompressionLz4  CassandraCompression = "LZ4"
)

type CassandraConnectionInput struct {
	// The internal entity unique identifier
	Goid string `json:"goid"`
	// The Cassandra Connection name
	Name string `json:"name"`
	// The Cassandra keyspace name
	Keyspace string `json:"keyspace"`
	// The Cassandra connection points
	ContactPoints []string `json:"contactPoints"`
	// The Cassandra server port
	Port int `json:"port"`
	// The username
	Username string `json:"username"`
	// The secure password reference.
	SecurePasswordName string `json:"securePasswordName"`
	// The Cassandra connection compression type
	Compression CassandraCompression `json:"compression"`
	// Whether this Cassandra connection is SSL enabled
	SslEnabled bool `json:"sslEnabled"`
	// Cipher suites used for SSL connection
	CipherSuites []string `json:"cipherSuites"`
	// Whether this Cassandra connection is enabled
	Enabled bool `json:"enabled"`
	// The Cassandra connection properties
	Properties []*EntityPropertyInput `json:"properties,omitempty"`
	// Ignored at creation time but can be used to compare bundle with gw state
	Checksum string `json:"checksum"`
}

// GetGoid returns CassandraConnectionInput.Goid, and is useful for accessing the field via an interface.
func (v *CassandraConnectionInput) GetGoid() string { return v.Goid }

// GetName returns CassandraConnectionInput.Name, and is useful for accessing the field via an interface.
func (v *CassandraConnectionInput) GetName() string { return v.Name }

// GetKeyspace returns CassandraConnectionInput.Keyspace, and is useful for accessing the field via an interface.
func (v *CassandraConnectionInput) GetKeyspace() string { return v.Keyspace }

// GetContactPoints returns CassandraConnectionInput.ContactPoints, and is useful for accessing the field via an interface.
func (v *CassandraConnectionInput) GetContactPoints() []string { return v.ContactPoints }

// GetPort returns CassandraConnectionInput.Port, and is useful for accessing the field via an interface.
func (v *CassandraConnectionInput) GetPort() int { return v.Port }

// GetUsername returns CassandraConnectionInput.Username, and is useful for accessing the field via an interface.
func (v *CassandraConnectionInput) GetUsername() string { return v.Username }

// GetSecurePasswordName returns CassandraConnectionInput.SecurePasswordName, and is useful for accessing the field via an interface.
func (v *CassandraConnectionInput) GetSecurePasswordName() string { return v.SecurePasswordName }

// GetCompression returns CassandraConnectionInput.Compression, and is useful for accessing the field via an interface.
func (v *CassandraConnectionInput) GetCompression() CassandraCompression { return v.Compression }

// GetSslEnabled returns CassandraConnectionInput.SslEnabled, and is useful for accessing the field via an interface.
func (v *CassandraConnectionInput) GetSslEnabled() bool { return v.SslEnabled }

// GetCipherSuites returns CassandraConnectionInput.CipherSuites, and is useful for accessing the field via an interface.
func (v *CassandraConnectionInput) GetCipherSuites() []string { return v.CipherSuites }

// GetEnabled returns CassandraConnectionInput.Enabled, and is useful for accessing the field via an interface.
func (v *CassandraConnectionInput) GetEnabled() bool { return v.Enabled }

// GetProperties returns CassandraConnectionInput.Properties, and is useful for accessing the field via an interface.
func (v *CassandraConnectionInput) GetProperties() []*EntityPropertyInput { return v.Properties }

// GetChecksum returns CassandraConnectionInput.Checksum, and is useful for accessing the field via an interface.
func (v *CassandraConnectionInput) GetChecksum() string { return v.Checksum }

type CertRevocationCheckPropertyType string

const (
	// Type for checking against a CRL from a URL contained in an X.509 Certificate
	CertRevocationCheckPropertyTypeCrlFromCertificate CertRevocationCheckPropertyType = "CRL_FROM_CERTIFICATE"
	// Type for checking against a CRL from a specified URL
	CertRevocationCheckPropertyTypeCrlFromUrl CertRevocationCheckPropertyType = "CRL_FROM_URL"
	// Type for OCSP check against a responder URL contained in an X.509 Certificate
	CertRevocationCheckPropertyTypeOcspFromCertificate CertRevocationCheckPropertyType = "OCSP_FROM_CERTIFICATE"
	// Type for OCSP check against a specified responder URL
	CertRevocationCheckPropertyTypeOcspFromUrl CertRevocationCheckPropertyType = "OCSP_FROM_URL"
)

type CertValidationType string

const (
	CertValidationTypeUseDefault      CertValidationType = "USE_DEFAULT"
	CertValidationTypeCertificateOnly CertValidationType = "CERTIFICATE_ONLY"
	CertValidationTypePathValidation  CertValidationType = "PATH_VALIDATION"
	CertValidationTypeRevocation      CertValidationType = "REVOCATION"
)

type CertificateValidationType string

const (
	CertificateValidationTypeCertificateOnly CertificateValidationType = "CERTIFICATE_ONLY"
	CertificateValidationTypePathValidation  CertificateValidationType = "PATH_VALIDATION"
	CertificateValidationTypeRevocation      CertificateValidationType = "REVOCATION"
)

// The inputs sent with the setClusterProperty Mutation
type ClusterPropertyInput struct {
	// The internal entity unique identifier
	Goid string `json:"goid"`
	// The name of the cluster property
	Name string `json:"name"`
	// The value of the cluster property to set
	Value string `json:"value"`
	// The cluster property description
	Description string `json:"description"`
	// Whether this is a hidden property. (Note that, this field has no effect on the mutation)
	HiddenProperty bool `json:"hiddenProperty"`
	// Ignored at creation time but can be used to compare bundle with gw state
	Checksum string `json:"checksum"`
}

// GetGoid returns ClusterPropertyInput.Goid, and is useful for accessing the field via an interface.
func (v *ClusterPropertyInput) GetGoid() string { return v.Goid }

// GetName returns ClusterPropertyInput.Name, and is useful for accessing the field via an interface.
func (v *ClusterPropertyInput) GetName() string { return v.Name }

// GetValue returns ClusterPropertyInput.Value, and is useful for accessing the field via an interface.
func (v *ClusterPropertyInput) GetValue() string { return v.Value }

// GetDescription returns ClusterPropertyInput.Description, and is useful for accessing the field via an interface.
func (v *ClusterPropertyInput) GetDescription() string { return v.Description }

// GetHiddenProperty returns ClusterPropertyInput.HiddenProperty, and is useful for accessing the field via an interface.
func (v *ClusterPropertyInput) GetHiddenProperty() bool { return v.HiddenProperty }

// GetChecksum returns ClusterPropertyInput.Checksum, and is useful for accessing the field via an interface.
func (v *ClusterPropertyInput) GetChecksum() string { return v.Checksum }

// The inputs sent with the setCustomKeyValue Mutation
type CustomKeyValueInput struct {
	// The goid for the custom key value
	Goid string `json:"goid"`
	// The custom key
	Key string `json:"key"`
	// The custom value in Base64 encoded format
	Value string `json:"value"`
	// Ignored at creation time but can be used to compare bundle with gw state
	Checksum string `json:"checksum"`
}

// GetGoid returns CustomKeyValueInput.Goid, and is useful for accessing the field via an interface.
func (v *CustomKeyValueInput) GetGoid() string { return v.Goid }

// GetKey returns CustomKeyValueInput.Key, and is useful for accessing the field via an interface.
func (v *CustomKeyValueInput) GetKey() string { return v.Key }

// GetValue returns CustomKeyValueInput.Value, and is useful for accessing the field via an interface.
func (v *CustomKeyValueInput) GetValue() string { return v.Value }

// GetChecksum returns CustomKeyValueInput.Checksum, and is useful for accessing the field via an interface.
func (v *CustomKeyValueInput) GetChecksum() string { return v.Checksum }

type DataType string

const (
	DataTypeString      DataType = "STRING"
	DataTypeCertificate DataType = "CERTIFICATE"
	DataTypeInteger     DataType = "INTEGER"
	DataTypeDecimal     DataType = "DECIMAL"
	DataTypeFloat       DataType = "FLOAT"
	DataTypeElement     DataType = "ELEMENT"
	DataTypeBoolean     DataType = "BOOLEAN"
	DataTypeBinary      DataType = "BINARY"
	DataTypeDateTime    DataType = "DATE_TIME"
	DataTypeMessage     DataType = "MESSAGE"
	DataTypeBlob        DataType = "BLOB"
	DataTypeClob        DataType = "CLOB"
	DataTypeUnknown     DataType = "UNKNOWN"
)

type DtdInput struct {
	// The internal entity unique identifier
	Goid string `json:"goid"`
	// A reference to the dtd. This id is what is referred to in policy and is often mirror of the target namespace
	SystemId string `json:"systemId"`
	// The public id for the dtd
	PublicId string `json:"publicId"`
	// An optional description
	Description string `json:"description"`
	// The actual dtd itself
	Content string `json:"content"`
	// Ignored at creation time but can be used to compare bundle with gw state
	Checksum string `json:"checksum"`
}

// GetGoid returns DtdInput.Goid, and is useful for accessing the field via an interface.
func (v *DtdInput) GetGoid() string { return v.Goid }

// GetSystemId returns DtdInput.SystemId, and is useful for accessing the field via an interface.
func (v *DtdInput) GetSystemId() string { return v.SystemId }

// GetPublicId returns DtdInput.PublicId, and is useful for accessing the field via an interface.
func (v *DtdInput) GetPublicId() string { return v.PublicId }

// GetDescription returns DtdInput.Description, and is useful for accessing the field via an interface.
func (v *DtdInput) GetDescription() string { return v.Description }

// GetContent returns DtdInput.Content, and is useful for accessing the field via an interface.
func (v *DtdInput) GetContent() string { return v.Content }

// GetChecksum returns DtdInput.Checksum, and is useful for accessing the field via an interface.
func (v *DtdInput) GetChecksum() string { return v.Checksum }

type EmailListenerInput struct {
	// The internal entity unique identifier
	Goid string `json:"goid"`
	// The name of the email listener. If you are creating several listeners, make sure the name is descriptive
	Name string `json:"name"`
	// Whether this email listener is enabled(active)
	Enabled bool `json:"enabled"`
	// The hostname of the email server. This name is verified against the X.509 certificate
	Hostname string `json:"hostname"`
	// The port number to monitor
	Port int `json:"port"`
	// The type of email server (IMAP or POP3)
	ServerType EmailServerType `json:"serverType"`
	// Whether email server connection (POP3S or IMAPS) is SSL enabled
	SslEnabled bool `json:"sslEnabled"`
	// Whether delete the messages on the mail server after retrieving
	DeleteOnReceive bool `json:"deleteOnReceive"`
	// The folder name to check for emails (Only for IMAP)
	Folder string `json:"folder"`
	// The listener will check for email after the specified number of seconds
	PollInterval int `json:"pollInterval"`
	// Email account name
	Username string `json:"username"`
	// Email account password. The password could be in plain text or secure password reference
	Password string `json:"password"`
	// The name of the published service hardwired to the email listener
	HardwiredServiceName string `json:"hardwiredServiceName"`
	// Permitted maximum size of the message
	SizeLimit int `json:"sizeLimit"`
	// [Optional] The Email listener Properties excluding sizeLimit and
	// HardwiredServiceName. When specified, will replace existing properties
	Properties []*EntityPropertyInput `json:"properties,omitempty"`
	// Ignored at creation time but can be used to compare bundle with gw state
	Checksum string `json:"checksum"`
}

// GetGoid returns EmailListenerInput.Goid, and is useful for accessing the field via an interface.
func (v *EmailListenerInput) GetGoid() string { return v.Goid }

// GetName returns EmailListenerInput.Name, and is useful for accessing the field via an interface.
func (v *EmailListenerInput) GetName() string { return v.Name }

// GetEnabled returns EmailListenerInput.Enabled, and is useful for accessing the field via an interface.
func (v *EmailListenerInput) GetEnabled() bool { return v.Enabled }

// GetHostname returns EmailListenerInput.Hostname, and is useful for accessing the field via an interface.
func (v *EmailListenerInput) GetHostname() string { return v.Hostname }

// GetPort returns EmailListenerInput.Port, and is useful for accessing the field via an interface.
func (v *EmailListenerInput) GetPort() int { return v.Port }

// GetServerType returns EmailListenerInput.ServerType, and is useful for accessing the field via an interface.
func (v *EmailListenerInput) GetServerType() EmailServerType { return v.ServerType }

// GetSslEnabled returns EmailListenerInput.SslEnabled, and is useful for accessing the field via an interface.
func (v *EmailListenerInput) GetSslEnabled() bool { return v.SslEnabled }

// GetDeleteOnReceive returns EmailListenerInput.DeleteOnReceive, and is useful for accessing the field via an interface.
func (v *EmailListenerInput) GetDeleteOnReceive() bool { return v.DeleteOnReceive }

// GetFolder returns EmailListenerInput.Folder, and is useful for accessing the field via an interface.
func (v *EmailListenerInput) GetFolder() string { return v.Folder }

// GetPollInterval returns EmailListenerInput.PollInterval, and is useful for accessing the field via an interface.
func (v *EmailListenerInput) GetPollInterval() int { return v.PollInterval }

// GetUsername returns EmailListenerInput.Username, and is useful for accessing the field via an interface.
func (v *EmailListenerInput) GetUsername() string { return v.Username }

// GetPassword returns EmailListenerInput.Password, and is useful for accessing the field via an interface.
func (v *EmailListenerInput) GetPassword() string { return v.Password }

// GetHardwiredServiceName returns EmailListenerInput.HardwiredServiceName, and is useful for accessing the field via an interface.
func (v *EmailListenerInput) GetHardwiredServiceName() string { return v.HardwiredServiceName }

// GetSizeLimit returns EmailListenerInput.SizeLimit, and is useful for accessing the field via an interface.
func (v *EmailListenerInput) GetSizeLimit() int { return v.SizeLimit }

// GetProperties returns EmailListenerInput.Properties, and is useful for accessing the field via an interface.
func (v *EmailListenerInput) GetProperties() []*EntityPropertyInput { return v.Properties }

// GetChecksum returns EmailListenerInput.Checksum, and is useful for accessing the field via an interface.
func (v *EmailListenerInput) GetChecksum() string { return v.Checksum }

type EmailServerType string

const (
	EmailServerTypeImap EmailServerType = "IMAP"
	EmailServerTypePop3 EmailServerType = "POP3"
)

// The description of an input argument for an encapsulated assertion for use when
// creating or updating an existing encass config
type EncassArgInput struct {
	// The name of the input
	Name string `json:"name"`
	// The type of input
	Type DataType `json:"type"`
	// The order of the argument in the admin gui
	Ordinal int `json:"ordinal"`
	// The prompt in the admin gui for this encass argument
	GuiPrompt bool `json:"guiPrompt"`
	// The label in the admin gui associated with this encass argument
	GuiLabel string `json:"guiLabel"`
}

// GetName returns EncassArgInput.Name, and is useful for accessing the field via an interface.
func (v *EncassArgInput) GetName() string { return v.Name }

// GetType returns EncassArgInput.Type, and is useful for accessing the field via an interface.
func (v *EncassArgInput) GetType() DataType { return v.Type }

// GetOrdinal returns EncassArgInput.Ordinal, and is useful for accessing the field via an interface.
func (v *EncassArgInput) GetOrdinal() int { return v.Ordinal }

// GetGuiPrompt returns EncassArgInput.GuiPrompt, and is useful for accessing the field via an interface.
func (v *EncassArgInput) GetGuiPrompt() bool { return v.GuiPrompt }

// GetGuiLabel returns EncassArgInput.GuiLabel, and is useful for accessing the field via an interface.
func (v *EncassArgInput) GetGuiLabel() string { return v.GuiLabel }

// The description of a new encapsulated assertion configuration being created
type EncassConfigInput struct {
	// The internal entity unique identifier
	Goid string `json:"goid"`
	// The guid for this encass config, can be omitted and a new one is assigned
	Guid string `json:"guid"`
	// The name of the encass config
	Name        string `json:"name"`
	Description string `json:"description"`
	// The policy it points to and its dependencies
	PolicyName string `json:"policyName"`
	// the input argument descriptions for this encass
	EncassArgs []*EncassArgInput `json:"encassArgs,omitempty"`
	// the output descriptions
	EncassResults []*EncassResultInput   `json:"encassResults,omitempty"`
	Properties    []*EntityPropertyInput `json:"properties,omitempty"`
	// Ignored at creation time but can be used to compare bundle with gw state
	Checksum string `json:"checksum"`
}

// GetGoid returns EncassConfigInput.Goid, and is useful for accessing the field via an interface.
func (v *EncassConfigInput) GetGoid() string { return v.Goid }

// GetGuid returns EncassConfigInput.Guid, and is useful for accessing the field via an interface.
func (v *EncassConfigInput) GetGuid() string { return v.Guid }

// GetName returns EncassConfigInput.Name, and is useful for accessing the field via an interface.
func (v *EncassConfigInput) GetName() string { return v.Name }

// GetDescription returns EncassConfigInput.Description, and is useful for accessing the field via an interface.
func (v *EncassConfigInput) GetDescription() string { return v.Description }

// GetPolicyName returns EncassConfigInput.PolicyName, and is useful for accessing the field via an interface.
func (v *EncassConfigInput) GetPolicyName() string { return v.PolicyName }

// GetEncassArgs returns EncassConfigInput.EncassArgs, and is useful for accessing the field via an interface.
func (v *EncassConfigInput) GetEncassArgs() []*EncassArgInput { return v.EncassArgs }

// GetEncassResults returns EncassConfigInput.EncassResults, and is useful for accessing the field via an interface.
func (v *EncassConfigInput) GetEncassResults() []*EncassResultInput { return v.EncassResults }

// GetProperties returns EncassConfigInput.Properties, and is useful for accessing the field via an interface.
func (v *EncassConfigInput) GetProperties() []*EntityPropertyInput { return v.Properties }

// GetChecksum returns EncassConfigInput.Checksum, and is useful for accessing the field via an interface.
func (v *EncassConfigInput) GetChecksum() string { return v.Checksum }

// The description of an output from the encapsulated assertion for use when creating or updating an existing encass config
type EncassResultInput struct {
	// The name of the output
	Name string `json:"name"`
	// The type of the output
	Type DataType `json:"type"`
}

// GetName returns EncassResultInput.Name, and is useful for accessing the field via an interface.
func (v *EncassResultInput) GetName() string { return v.Name }

// GetType returns EncassResultInput.Type, and is useful for accessing the field via an interface.
func (v *EncassResultInput) GetType() DataType { return v.Type }

type EntityFieldOption string

const (
	EntityFieldOptionDefault EntityFieldOption = "DEFAULT"
	EntityFieldOptionNone    EntityFieldOption = "NONE"
	EntityFieldOptionCustom  EntityFieldOption = "CUSTOM"
)

type EntityMutationAction string

const (
	EntityMutationActionNewOrUpdate     EntityMutationAction = "NEW_OR_UPDATE"
	EntityMutationActionNewOrExisting   EntityMutationAction = "NEW_OR_EXISTING"
	EntityMutationActionAlwaysCreateNew EntityMutationAction = "ALWAYS_CREATE_NEW"
	EntityMutationActionIgnore          EntityMutationAction = "IGNORE"
	EntityMutationActionDelete          EntityMutationAction = "DELETE"
)

type EntityMutationStatus string

const (
	EntityMutationStatusNone         EntityMutationStatus = "NONE"
	EntityMutationStatusCreated      EntityMutationStatus = "CREATED"
	EntityMutationStatusUpdated      EntityMutationStatus = "UPDATED"
	EntityMutationStatusDeleted      EntityMutationStatus = "DELETED"
	EntityMutationStatusUsedExisting EntityMutationStatus = "USED_EXISTING"
	EntityMutationStatusIgnored      EntityMutationStatus = "IGNORED"
	EntityMutationStatusError        EntityMutationStatus = "ERROR"
)

type EntityPropertyInput struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// GetName returns EntityPropertyInput.Name, and is useful for accessing the field via an interface.
func (v *EntityPropertyInput) GetName() string { return v.Name }

// GetValue returns EntityPropertyInput.Value, and is useful for accessing the field via an interface.
func (v *EntityPropertyInput) GetValue() string { return v.Value }

type FederatedGroupInput struct {
	Name string `json:"name"`
	// If provided, will try to honour at creation time
	Goid string `json:"goid"`
	// The name of the FiP this group is defined in
	ProviderName string `json:"providerName"`
	Description  string `json:"description"`
	// Ignored at creation time but can be used to compare bundle with gw state
	Checksum string `json:"checksum"`
}

// GetName returns FederatedGroupInput.Name, and is useful for accessing the field via an interface.
func (v *FederatedGroupInput) GetName() string { return v.Name }

// GetGoid returns FederatedGroupInput.Goid, and is useful for accessing the field via an interface.
func (v *FederatedGroupInput) GetGoid() string { return v.Goid }

// GetProviderName returns FederatedGroupInput.ProviderName, and is useful for accessing the field via an interface.
func (v *FederatedGroupInput) GetProviderName() string { return v.ProviderName }

// GetDescription returns FederatedGroupInput.Description, and is useful for accessing the field via an interface.
func (v *FederatedGroupInput) GetDescription() string { return v.Description }

// GetChecksum returns FederatedGroupInput.Checksum, and is useful for accessing the field via an interface.
func (v *FederatedGroupInput) GetChecksum() string { return v.Checksum }

type FederatedIdpInput struct {
	Name string `json:"name"`
	// Will try to match goid if provided
	Goid           string             `json:"goid"`
	SupportsSAML   bool               `json:"supportsSAML"`
	SupportsX509   bool               `json:"supportsX509"`
	CertValidation CertValidationType `json:"certValidation"`
	// The certificates in the trusted certificate table that establish the trust for this FIP
	TrustedCerts []*TrustedCertPartialInput `json:"trustedCerts,omitempty"`
	// The optional checksum is ignored during the mutation but can be used to compare bundle content
	Checksum string `json:"checksum"`
}

// GetName returns FederatedIdpInput.Name, and is useful for accessing the field via an interface.
func (v *FederatedIdpInput) GetName() string { return v.Name }

// GetGoid returns FederatedIdpInput.Goid, and is useful for accessing the field via an interface.
func (v *FederatedIdpInput) GetGoid() string { return v.Goid }

// GetSupportsSAML returns FederatedIdpInput.SupportsSAML, and is useful for accessing the field via an interface.
func (v *FederatedIdpInput) GetSupportsSAML() bool { return v.SupportsSAML }

// GetSupportsX509 returns FederatedIdpInput.SupportsX509, and is useful for accessing the field via an interface.
func (v *FederatedIdpInput) GetSupportsX509() bool { return v.SupportsX509 }

// GetCertValidation returns FederatedIdpInput.CertValidation, and is useful for accessing the field via an interface.
func (v *FederatedIdpInput) GetCertValidation() CertValidationType { return v.CertValidation }

// GetTrustedCerts returns FederatedIdpInput.TrustedCerts, and is useful for accessing the field via an interface.
func (v *FederatedIdpInput) GetTrustedCerts() []*TrustedCertPartialInput { return v.TrustedCerts }

// GetChecksum returns FederatedIdpInput.Checksum, and is useful for accessing the field via an interface.
func (v *FederatedIdpInput) GetChecksum() string { return v.Checksum }

type FederatedUserInput struct {
	Name string `json:"name"`
	// If provided, will try to honour at creation time
	Goid string `json:"goid"`
	// The name of the FiP this user is defined as part of
	ProviderName string `json:"providerName"`
	// Whether to replace existing group memberships or not
	ReplaceGroupMemberships bool `json:"replaceGroupMemberships"`
	// The list of fip group details (names) that this user is member of. If you pass
	// empty array, will reset memberships. If absent, does not affect memberships
	// for current user.
	MemberOf  []*MembershipInput `json:"memberOf,omitempty"`
	Login     string             `json:"login"`
	SubjectDn string             `json:"subjectDn"`
	// A client-side certificate associated with this user to use for pki type authentication
	CertBase64 string `json:"certBase64"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	Checksum   string `json:"checksum"`
}

// GetName returns FederatedUserInput.Name, and is useful for accessing the field via an interface.
func (v *FederatedUserInput) GetName() string { return v.Name }

// GetGoid returns FederatedUserInput.Goid, and is useful for accessing the field via an interface.
func (v *FederatedUserInput) GetGoid() string { return v.Goid }

// GetProviderName returns FederatedUserInput.ProviderName, and is useful for accessing the field via an interface.
func (v *FederatedUserInput) GetProviderName() string { return v.ProviderName }

// GetReplaceGroupMemberships returns FederatedUserInput.ReplaceGroupMemberships, and is useful for accessing the field via an interface.
func (v *FederatedUserInput) GetReplaceGroupMemberships() bool { return v.ReplaceGroupMemberships }

// GetMemberOf returns FederatedUserInput.MemberOf, and is useful for accessing the field via an interface.
func (v *FederatedUserInput) GetMemberOf() []*MembershipInput { return v.MemberOf }

// GetLogin returns FederatedUserInput.Login, and is useful for accessing the field via an interface.
func (v *FederatedUserInput) GetLogin() string { return v.Login }

// GetSubjectDn returns FederatedUserInput.SubjectDn, and is useful for accessing the field via an interface.
func (v *FederatedUserInput) GetSubjectDn() string { return v.SubjectDn }

// GetCertBase64 returns FederatedUserInput.CertBase64, and is useful for accessing the field via an interface.
func (v *FederatedUserInput) GetCertBase64() string { return v.CertBase64 }

// GetFirstName returns FederatedUserInput.FirstName, and is useful for accessing the field via an interface.
func (v *FederatedUserInput) GetFirstName() string { return v.FirstName }

// GetLastName returns FederatedUserInput.LastName, and is useful for accessing the field via an interface.
func (v *FederatedUserInput) GetLastName() string { return v.LastName }

// GetEmail returns FederatedUserInput.Email, and is useful for accessing the field via an interface.
func (v *FederatedUserInput) GetEmail() string { return v.Email }

// GetChecksum returns FederatedUserInput.Checksum, and is useful for accessing the field via an interface.
func (v *FederatedUserInput) GetChecksum() string { return v.Checksum }

type FipCertInput struct {
	// The thumbprint of the cert to use as trust for a federated identity provider
	ThumbprintSha1 string `json:"thumbprintSha1"`
	// The internal entity unique identifier. (Note that, this field has no effect on the mutation)
	Goid string `json:"goid"`
	// The name of the trusted certificate. (Note that, this field has no effect on the mutation)
	Name string `json:"name"`
	// The base 64 encoded string of the certificate. (Note that, this field has no effect on the mutation)
	CertBase64 string `json:"certBase64"`
	// Whether to perform hostname verification with this certificate. (Note that, this field has no effect on the mutation)
	VerifyHostname bool `json:"verifyHostname"`
	// Whether this certificate is a trust anchor. (Note that, this field has no effect on the mutation)
	TrustAnchor bool `json:"trustAnchor"`
	// What the certificate is trusted for. (Note that, this field has no effect on the mutation)
	TrustedFor []TrustedForType `json:"trustedFor"`
	// The revocation check policy type. (Note that, this field has no effect on the mutation)
	RevocationCheckPolicyType PolicyUsageType `json:"revocationCheckPolicyType"`
	// The name of revocation policy. (Note that, this field has no effect on the mutation)
	RevocationCheckPolicyName string `json:"revocationCheckPolicyName"`
	// The Subject DN of this certificate. (Note that, this field has no effect on the mutation)
	SubjectDn string `json:"subjectDn"`
	// The start date of the validity period. (Note that, this field has no effect on the mutation)
	NotBefore string `json:"notBefore"`
	// the end date of the validity period. (Note that, this field has no effect on the mutation)
	NotAfter string `json:"notAfter"`
	// Ignored at creation time but can be used to compare bundle with gw state
	Checksum string `json:"checksum"`
}

// GetThumbprintSha1 returns FipCertInput.ThumbprintSha1, and is useful for accessing the field via an interface.
func (v *FipCertInput) GetThumbprintSha1() string { return v.ThumbprintSha1 }

// GetGoid returns FipCertInput.Goid, and is useful for accessing the field via an interface.
func (v *FipCertInput) GetGoid() string { return v.Goid }

// GetName returns FipCertInput.Name, and is useful for accessing the field via an interface.
func (v *FipCertInput) GetName() string { return v.Name }

// GetCertBase64 returns FipCertInput.CertBase64, and is useful for accessing the field via an interface.
func (v *FipCertInput) GetCertBase64() string { return v.CertBase64 }

// GetVerifyHostname returns FipCertInput.VerifyHostname, and is useful for accessing the field via an interface.
func (v *FipCertInput) GetVerifyHostname() bool { return v.VerifyHostname }

// GetTrustAnchor returns FipCertInput.TrustAnchor, and is useful for accessing the field via an interface.
func (v *FipCertInput) GetTrustAnchor() bool { return v.TrustAnchor }

// GetTrustedFor returns FipCertInput.TrustedFor, and is useful for accessing the field via an interface.
func (v *FipCertInput) GetTrustedFor() []TrustedForType { return v.TrustedFor }

// GetRevocationCheckPolicyType returns FipCertInput.RevocationCheckPolicyType, and is useful for accessing the field via an interface.
func (v *FipCertInput) GetRevocationCheckPolicyType() PolicyUsageType {
	return v.RevocationCheckPolicyType
}

// GetRevocationCheckPolicyName returns FipCertInput.RevocationCheckPolicyName, and is useful for accessing the field via an interface.
func (v *FipCertInput) GetRevocationCheckPolicyName() string { return v.RevocationCheckPolicyName }

// GetSubjectDn returns FipCertInput.SubjectDn, and is useful for accessing the field via an interface.
func (v *FipCertInput) GetSubjectDn() string { return v.SubjectDn }

// GetNotBefore returns FipCertInput.NotBefore, and is useful for accessing the field via an interface.
func (v *FipCertInput) GetNotBefore() string { return v.NotBefore }

// GetNotAfter returns FipCertInput.NotAfter, and is useful for accessing the field via an interface.
func (v *FipCertInput) GetNotAfter() string { return v.NotAfter }

// GetChecksum returns FipCertInput.Checksum, and is useful for accessing the field via an interface.
func (v *FipCertInput) GetChecksum() string { return v.Checksum }

type FipGroupInput struct {
	Name string `json:"name"`
	// If provided, will try to honour at creation time
	Goid string `json:"goid"`
	// The name of the FiP this group is defined in
	ProviderName string `json:"providerName"`
	Description  string `json:"description"`
	// Ignored at creation time but can be used to compare bundle with gw state
	Checksum string `json:"checksum"`
}

// GetName returns FipGroupInput.Name, and is useful for accessing the field via an interface.
func (v *FipGroupInput) GetName() string { return v.Name }

// GetGoid returns FipGroupInput.Goid, and is useful for accessing the field via an interface.
func (v *FipGroupInput) GetGoid() string { return v.Goid }

// GetProviderName returns FipGroupInput.ProviderName, and is useful for accessing the field via an interface.
func (v *FipGroupInput) GetProviderName() string { return v.ProviderName }

// GetDescription returns FipGroupInput.Description, and is useful for accessing the field via an interface.
func (v *FipGroupInput) GetDescription() string { return v.Description }

// GetChecksum returns FipGroupInput.Checksum, and is useful for accessing the field via an interface.
func (v *FipGroupInput) GetChecksum() string { return v.Checksum }

type FipInput struct {
	Name string `json:"name"`
	// Will try to match goid if provided
	Goid                     string                    `json:"goid"`
	EnableCredentialTypeSaml bool                      `json:"enableCredentialTypeSaml"`
	EnableCredentialTypeX509 bool                      `json:"enableCredentialTypeX509"`
	CertificateValidation    CertificateValidationType `json:"certificateValidation"`
	// The certificates in the trusted certificate table that establish the trust for this FIP
	CertificateReferences []*FipCertInput `json:"certificateReferences,omitempty"`
	// The optional checksum is ignored during the mutation but can be used to compare bundle content
	Checksum string `json:"checksum"`
}

// GetName returns FipInput.Name, and is useful for accessing the field via an interface.
func (v *FipInput) GetName() string { return v.Name }

// GetGoid returns FipInput.Goid, and is useful for accessing the field via an interface.
func (v *FipInput) GetGoid() string { return v.Goid }

// GetEnableCredentialTypeSaml returns FipInput.EnableCredentialTypeSaml, and is useful for accessing the field via an interface.
func (v *FipInput) GetEnableCredentialTypeSaml() bool { return v.EnableCredentialTypeSaml }

// GetEnableCredentialTypeX509 returns FipInput.EnableCredentialTypeX509, and is useful for accessing the field via an interface.
func (v *FipInput) GetEnableCredentialTypeX509() bool { return v.EnableCredentialTypeX509 }

// GetCertificateValidation returns FipInput.CertificateValidation, and is useful for accessing the field via an interface.
func (v *FipInput) GetCertificateValidation() CertificateValidationType {
	return v.CertificateValidation
}

// GetCertificateReferences returns FipInput.CertificateReferences, and is useful for accessing the field via an interface.
func (v *FipInput) GetCertificateReferences() []*FipCertInput { return v.CertificateReferences }

// GetChecksum returns FipInput.Checksum, and is useful for accessing the field via an interface.
func (v *FipInput) GetChecksum() string { return v.Checksum }

type FipUserInput struct {
	Name string `json:"name"`
	// If provided, will try to honour at creation time
	Goid string `json:"goid"`
	// The name of the FiP this user is defined as part of
	ProviderName string `json:"providerName"`
	// The list of fip group details (names) that this user is member of. If you pass
	// empty array, will reset memberships. If absent, does not affect memberships
	// for current user.
	MemberOf  []*MembershipInput `json:"memberOf,omitempty"`
	Login     string             `json:"login"`
	SubjectDn string             `json:"subjectDn"`
	// A client-side certificate associated with this user to use for pki type authentication
	CertBase64 string `json:"certBase64"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	Checksum   string `json:"checksum"`
}

// GetName returns FipUserInput.Name, and is useful for accessing the field via an interface.
func (v *FipUserInput) GetName() string { return v.Name }

// GetGoid returns FipUserInput.Goid, and is useful for accessing the field via an interface.
func (v *FipUserInput) GetGoid() string { return v.Goid }

// GetProviderName returns FipUserInput.ProviderName, and is useful for accessing the field via an interface.
func (v *FipUserInput) GetProviderName() string { return v.ProviderName }

// GetMemberOf returns FipUserInput.MemberOf, and is useful for accessing the field via an interface.
func (v *FipUserInput) GetMemberOf() []*MembershipInput { return v.MemberOf }

// GetLogin returns FipUserInput.Login, and is useful for accessing the field via an interface.
func (v *FipUserInput) GetLogin() string { return v.Login }

// GetSubjectDn returns FipUserInput.SubjectDn, and is useful for accessing the field via an interface.
func (v *FipUserInput) GetSubjectDn() string { return v.SubjectDn }

// GetCertBase64 returns FipUserInput.CertBase64, and is useful for accessing the field via an interface.
func (v *FipUserInput) GetCertBase64() string { return v.CertBase64 }

// GetFirstName returns FipUserInput.FirstName, and is useful for accessing the field via an interface.
func (v *FipUserInput) GetFirstName() string { return v.FirstName }

// GetLastName returns FipUserInput.LastName, and is useful for accessing the field via an interface.
func (v *FipUserInput) GetLastName() string { return v.LastName }

// GetEmail returns FipUserInput.Email, and is useful for accessing the field via an interface.
func (v *FipUserInput) GetEmail() string { return v.Email }

// GetChecksum returns FipUserInput.Checksum, and is useful for accessing the field via an interface.
func (v *FipUserInput) GetChecksum() string { return v.Checksum }

type FolderInput struct {
	// The goid for the folder
	Goid string `json:"goid"`
	// The folder name
	Name string `json:"name"`
	// The folder Path
	Path string `json:"path"`
	// The configuration checksum of this folder
	Checksum string `json:"checksum"`
}

// GetGoid returns FolderInput.Goid, and is useful for accessing the field via an interface.
func (v *FolderInput) GetGoid() string { return v.Goid }

// GetName returns FolderInput.Name, and is useful for accessing the field via an interface.
func (v *FolderInput) GetName() string { return v.Name }

// GetPath returns FolderInput.Path, and is useful for accessing the field via an interface.
func (v *FolderInput) GetPath() string { return v.Path }

// GetChecksum returns FolderInput.Checksum, and is useful for accessing the field via an interface.
func (v *FolderInput) GetChecksum() string { return v.Checksum }

type GenericEntityInput struct {
	Goid string `json:"goid"`
	// unique name
	Name string `json:"name"`
	// The configuration checksum
	Checksum string `json:"checksum"`
	// description
	Description string `json:"description"`
	// XML representation of underlying entity details
	ValueXml string `json:"valueXml"`
	// Whether this Generic entity is enabled
	Enabled bool `json:"enabled"`
	// Absolute entity class name of Generic Entity
	EntityClassName string `json:"entityClassName"`
}

// GetGoid returns GenericEntityInput.Goid, and is useful for accessing the field via an interface.
func (v *GenericEntityInput) GetGoid() string { return v.Goid }

// GetName returns GenericEntityInput.Name, and is useful for accessing the field via an interface.
func (v *GenericEntityInput) GetName() string { return v.Name }

// GetChecksum returns GenericEntityInput.Checksum, and is useful for accessing the field via an interface.
func (v *GenericEntityInput) GetChecksum() string { return v.Checksum }

// GetDescription returns GenericEntityInput.Description, and is useful for accessing the field via an interface.
func (v *GenericEntityInput) GetDescription() string { return v.Description }

// GetValueXml returns GenericEntityInput.ValueXml, and is useful for accessing the field via an interface.
func (v *GenericEntityInput) GetValueXml() string { return v.ValueXml }

// GetEnabled returns GenericEntityInput.Enabled, and is useful for accessing the field via an interface.
func (v *GenericEntityInput) GetEnabled() bool { return v.Enabled }

// GetEntityClassName returns GenericEntityInput.EntityClassName, and is useful for accessing the field via an interface.
func (v *GenericEntityInput) GetEntityClassName() string { return v.EntityClassName }

type GlobalPolicyInput struct {
	// The name of the policy. Policies are unique by name.
	Name string `json:"name"`
	// The folder path where to create this policy.  If the path does not exist, it will be created
	FolderPath string `json:"folderPath"`
	// The goid for this policy
	Goid string `json:"goid"`
	// The guid for this service, if none provided, assigned at creation
	Guid string `json:"guid"`
	// Global policy tag. Possible values are :
	// message-completed
	// message-received
	// post-security
	// post-service
	// pre-security
	// pre-service
	Tag string `json:"tag,omitempty"`
	// The policy
	Policy *PolicyInput `json:"policy,omitempty"`
	Soap   bool         `json:"soap"`
	// Ignored at creation time but can be used to compare bundle with gw state
	Checksum string `json:"checksum"`
}

// GetName returns GlobalPolicyInput.Name, and is useful for accessing the field via an interface.
func (v *GlobalPolicyInput) GetName() string { return v.Name }

// GetFolderPath returns GlobalPolicyInput.FolderPath, and is useful for accessing the field via an interface.
func (v *GlobalPolicyInput) GetFolderPath() string { return v.FolderPath }

// GetGoid returns GlobalPolicyInput.Goid, and is useful for accessing the field via an interface.
func (v *GlobalPolicyInput) GetGoid() string { return v.Goid }

// GetGuid returns GlobalPolicyInput.Guid, and is useful for accessing the field via an interface.
func (v *GlobalPolicyInput) GetGuid() string { return v.Guid }

// GetTag returns GlobalPolicyInput.Tag, and is useful for accessing the field via an interface.
func (v *GlobalPolicyInput) GetTag() string { return v.Tag }

// GetPolicy returns GlobalPolicyInput.Policy, and is useful for accessing the field via an interface.
func (v *GlobalPolicyInput) GetPolicy() *PolicyInput { return v.Policy }

// GetSoap returns GlobalPolicyInput.Soap, and is useful for accessing the field via an interface.
func (v *GlobalPolicyInput) GetSoap() bool { return v.Soap }

// GetChecksum returns GlobalPolicyInput.Checksum, and is useful for accessing the field via an interface.
func (v *GlobalPolicyInput) GetChecksum() string { return v.Checksum }

type GroupMappingInput struct {
	ObjClass       string               `json:"objClass"`
	NameAttrName   string               `json:"nameAttrName"`
	MemberAttrName string               `json:"memberAttrName"`
	MemberStrategy *MemberStrategyInput `json:"memberStrategy,omitempty"`
}

// GetObjClass returns GroupMappingInput.ObjClass, and is useful for accessing the field via an interface.
func (v *GroupMappingInput) GetObjClass() string { return v.ObjClass }

// GetNameAttrName returns GroupMappingInput.NameAttrName, and is useful for accessing the field via an interface.
func (v *GroupMappingInput) GetNameAttrName() string { return v.NameAttrName }

// GetMemberAttrName returns GroupMappingInput.MemberAttrName, and is useful for accessing the field via an interface.
func (v *GroupMappingInput) GetMemberAttrName() string { return v.MemberAttrName }

// GetMemberStrategy returns GroupMappingInput.MemberStrategy, and is useful for accessing the field via an interface.
func (v *GroupMappingInput) GetMemberStrategy() *MemberStrategyInput { return v.MemberStrategy }

// IDP Group Reference input
type GroupRefInput struct {
	// The name of group
	Name string `json:"name"`
	// The subjectDn of group
	SubjectDn string `json:"subjectDn"`
	// The name of identity provider that the group belongs to
	ProviderName string `json:"providerName"`
	// The type of identity provider that the group belongs to
	ProviderType IdpType `json:"providerType"`
}

// GetName returns GroupRefInput.Name, and is useful for accessing the field via an interface.
func (v *GroupRefInput) GetName() string { return v.Name }

// GetSubjectDn returns GroupRefInput.SubjectDn, and is useful for accessing the field via an interface.
func (v *GroupRefInput) GetSubjectDn() string { return v.SubjectDn }

// GetProviderName returns GroupRefInput.ProviderName, and is useful for accessing the field via an interface.
func (v *GroupRefInput) GetProviderName() string { return v.ProviderName }

// GetProviderType returns GroupRefInput.ProviderType, and is useful for accessing the field via an interface.
func (v *GroupRefInput) GetProviderType() IdpType { return v.ProviderType }

type HttpConfigurationInput struct {
	// The goid for the http configuration
	Goid string `json:"goid"`
	// The host of the http configuration
	Host string `json:"host"`
	// The port of the http configuration
	Port int `json:"port"`
	// The protocol of the http configuration
	Protocol HttpScheme `json:"protocol"`
	// The path of the http configuration
	Path string `json:"path"`
	// The username of the http configuration
	Username string `json:"username"`
	// The securePasswordName of the http configuration
	SecurePasswordName string `json:"securePasswordName"`
	// The ntlmHost of the http configuration
	NtlmHost string `json:"ntlmHost"`
	// The ntlmDomain of the http configuration
	NtlmDomain string `json:"ntlmDomain"`
	// The tlsVersion of the http configuration
	TlsVersion string `json:"tlsVersion"`
	// The tlsKeyUse of the http configuration
	TlsKeyUse EntityFieldOption `json:"tlsKeyUse"`
	// The tlsKeystoreId of the http configuration
	TlsKeystoreId string `json:"tlsKeystoreId"`
	// The tlsKeyAlias of the http configuration
	TlsKeyAlias string `json:"tlsKeyAlias"`
	// The tlsCipherSuites of the http configuration
	TlsCipherSuites []string `json:"tlsCipherSuites"`
	// The connectTimeout of the http configuration
	ConnectTimeout int `json:"connectTimeout"`
	// The readTimeout of the http configuration
	ReadTimeout int `json:"readTimeout"`
	// The followRedirects of the http configuration
	FollowRedirects bool `json:"followRedirects"`
	// The proxyUse of the http configuration
	ProxyUse EntityFieldOption `json:"proxyUse"`
	// The HttpProxyConfiguration of the http configuration
	ProxyConfiguration *HttpProxyConfigurationInput `json:"proxyConfiguration,omitempty"`
	// Ignored at creation time but can be used to compare bundle with gw state
	Checksum string `json:"checksum"`
}

// GetGoid returns HttpConfigurationInput.Goid, and is useful for accessing the field via an interface.
func (v *HttpConfigurationInput) GetGoid() string { return v.Goid }

// GetHost returns HttpConfigurationInput.Host, and is useful for accessing the field via an interface.
func (v *HttpConfigurationInput) GetHost() string { return v.Host }

// GetPort returns HttpConfigurationInput.Port, and is useful for accessing the field via an interface.
func (v *HttpConfigurationInput) GetPort() int { return v.Port }

// GetProtocol returns HttpConfigurationInput.Protocol, and is useful for accessing the field via an interface.
func (v *HttpConfigurationInput) GetProtocol() HttpScheme { return v.Protocol }

// GetPath returns HttpConfigurationInput.Path, and is useful for accessing the field via an interface.
func (v *HttpConfigurationInput) GetPath() string { return v.Path }

// GetUsername returns HttpConfigurationInput.Username, and is useful for accessing the field via an interface.
func (v *HttpConfigurationInput) GetUsername() string { return v.Username }

// GetSecurePasswordName returns HttpConfigurationInput.SecurePasswordName, and is useful for accessing the field via an interface.
func (v *HttpConfigurationInput) GetSecurePasswordName() string { return v.SecurePasswordName }

// GetNtlmHost returns HttpConfigurationInput.NtlmHost, and is useful for accessing the field via an interface.
func (v *HttpConfigurationInput) GetNtlmHost() string { return v.NtlmHost }

// GetNtlmDomain returns HttpConfigurationInput.NtlmDomain, and is useful for accessing the field via an interface.
func (v *HttpConfigurationInput) GetNtlmDomain() string { return v.NtlmDomain }

// GetTlsVersion returns HttpConfigurationInput.TlsVersion, and is useful for accessing the field via an interface.
func (v *HttpConfigurationInput) GetTlsVersion() string { return v.TlsVersion }

// GetTlsKeyUse returns HttpConfigurationInput.TlsKeyUse, and is useful for accessing the field via an interface.
func (v *HttpConfigurationInput) GetTlsKeyUse() EntityFieldOption { return v.TlsKeyUse }

// GetTlsKeystoreId returns HttpConfigurationInput.TlsKeystoreId, and is useful for accessing the field via an interface.
func (v *HttpConfigurationInput) GetTlsKeystoreId() string { return v.TlsKeystoreId }

// GetTlsKeyAlias returns HttpConfigurationInput.TlsKeyAlias, and is useful for accessing the field via an interface.
func (v *HttpConfigurationInput) GetTlsKeyAlias() string { return v.TlsKeyAlias }

// GetTlsCipherSuites returns HttpConfigurationInput.TlsCipherSuites, and is useful for accessing the field via an interface.
func (v *HttpConfigurationInput) GetTlsCipherSuites() []string { return v.TlsCipherSuites }

// GetConnectTimeout returns HttpConfigurationInput.ConnectTimeout, and is useful for accessing the field via an interface.
func (v *HttpConfigurationInput) GetConnectTimeout() int { return v.ConnectTimeout }

// GetReadTimeout returns HttpConfigurationInput.ReadTimeout, and is useful for accessing the field via an interface.
func (v *HttpConfigurationInput) GetReadTimeout() int { return v.ReadTimeout }

// GetFollowRedirects returns HttpConfigurationInput.FollowRedirects, and is useful for accessing the field via an interface.
func (v *HttpConfigurationInput) GetFollowRedirects() bool { return v.FollowRedirects }

// GetProxyUse returns HttpConfigurationInput.ProxyUse, and is useful for accessing the field via an interface.
func (v *HttpConfigurationInput) GetProxyUse() EntityFieldOption { return v.ProxyUse }

// GetProxyConfiguration returns HttpConfigurationInput.ProxyConfiguration, and is useful for accessing the field via an interface.
func (v *HttpConfigurationInput) GetProxyConfiguration() *HttpProxyConfigurationInput {
	return v.ProxyConfiguration
}

// GetChecksum returns HttpConfigurationInput.Checksum, and is useful for accessing the field via an interface.
func (v *HttpConfigurationInput) GetChecksum() string { return v.Checksum }

// Support Http methods for Web API Service
type HttpMethod string

const (
	HttpMethodDelete  HttpMethod = "DELETE"
	HttpMethodHead    HttpMethod = "HEAD"
	HttpMethodGet     HttpMethod = "GET"
	HttpMethodPost    HttpMethod = "POST"
	HttpMethodPut     HttpMethod = "PUT"
	HttpMethodOptions HttpMethod = "OPTIONS"
	HttpMethodPatch   HttpMethod = "PATCH"
	HttpMethodOther   HttpMethod = "OTHER"
)

type HttpProxyConfigurationInput struct {
	// The proxyHost of the http proxy configuration
	Host string `json:"host"`
	// The proxyPort of the http proxy configuration
	Port int `json:"port"`
	// The proxyUsername of the http proxy configuration
	Username string `json:"username"`
	// The securePasswordName of the http proxy configuration
	SecurePasswordName string `json:"securePasswordName"`
}

// GetHost returns HttpProxyConfigurationInput.Host, and is useful for accessing the field via an interface.
func (v *HttpProxyConfigurationInput) GetHost() string { return v.Host }

// GetPort returns HttpProxyConfigurationInput.Port, and is useful for accessing the field via an interface.
func (v *HttpProxyConfigurationInput) GetPort() int { return v.Port }

// GetUsername returns HttpProxyConfigurationInput.Username, and is useful for accessing the field via an interface.
func (v *HttpProxyConfigurationInput) GetUsername() string { return v.Username }

// GetSecurePasswordName returns HttpProxyConfigurationInput.SecurePasswordName, and is useful for accessing the field via an interface.
func (v *HttpProxyConfigurationInput) GetSecurePasswordName() string { return v.SecurePasswordName }

type HttpScheme string

const (
	HttpSchemeHttp  HttpScheme = "HTTP"
	HttpSchemeHttps HttpScheme = "HTTPS"
	HttpSchemeAny   HttpScheme = "ANY"
)

type IdpType string

const (
	IdpTypeInternal     IdpType = "INTERNAL"
	IdpTypeFederated    IdpType = "FEDERATED"
	IdpTypeLdap         IdpType = "LDAP"
	IdpTypeSimpleLdap   IdpType = "SIMPLE_LDAP"
	IdpTypePolicyBacked IdpType = "POLICY_BACKED"
)

type InternalGroupInput struct {
	Name string `json:"name"`
	// If provided, will try to honour at creation time
	Goid        string `json:"goid"`
	Description string `json:"description"`
	// Ignored at creation time but can be used to compare bundle with gw state
	Checksum string `json:"checksum"`
}

// GetName returns InternalGroupInput.Name, and is useful for accessing the field via an interface.
func (v *InternalGroupInput) GetName() string { return v.Name }

// GetGoid returns InternalGroupInput.Goid, and is useful for accessing the field via an interface.
func (v *InternalGroupInput) GetGoid() string { return v.Goid }

// GetDescription returns InternalGroupInput.Description, and is useful for accessing the field via an interface.
func (v *InternalGroupInput) GetDescription() string { return v.Description }

// GetChecksum returns InternalGroupInput.Checksum, and is useful for accessing the field via an interface.
func (v *InternalGroupInput) GetChecksum() string { return v.Checksum }

type InternalIdpInput struct {
	Goid           string             `json:"goid"`
	Name           string             `json:"name"`
	Checksum       string             `json:"checksum"`
	CertValidation CertValidationType `json:"certValidation"`
}

// GetGoid returns InternalIdpInput.Goid, and is useful for accessing the field via an interface.
func (v *InternalIdpInput) GetGoid() string { return v.Goid }

// GetName returns InternalIdpInput.Name, and is useful for accessing the field via an interface.
func (v *InternalIdpInput) GetName() string { return v.Name }

// GetChecksum returns InternalIdpInput.Checksum, and is useful for accessing the field via an interface.
func (v *InternalIdpInput) GetChecksum() string { return v.Checksum }

// GetCertValidation returns InternalIdpInput.CertValidation, and is useful for accessing the field via an interface.
func (v *InternalIdpInput) GetCertValidation() CertValidationType { return v.CertValidation }

type InternalUserInput struct {
	Name string `json:"name"`
	// If provided, will try to honour at creation time
	Goid string `json:"goid"`
	// Whether to replace existing group memberships or not
	ReplaceGroupMemberships bool `json:"replaceGroupMemberships"`
	// The list of internal group details (names) that this user is member of. If you
	// pass empty array, will reset memberships. If absent, does not affect
	// memberships for current users.
	MemberOf []*MembershipInput `json:"memberOf,omitempty"`
	Login    string             `json:"login"`
	// You can either pass in the hashed password which comes back in queries or the raw passwd directly
	Password string `json:"password"`
	// A client-side certificate associated with this user to use for pki type authentication
	CertBase64 string `json:"certBase64"`
	// SSH public key
	SshPublicKey string `json:"sshPublicKey"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	// Is user enabled or not!
	Enabled bool `json:"enabled"`
	// Ignored at creation time but can be used to compare bundle with gw state
	Checksum string `json:"checksum"`
}

// GetName returns InternalUserInput.Name, and is useful for accessing the field via an interface.
func (v *InternalUserInput) GetName() string { return v.Name }

// GetGoid returns InternalUserInput.Goid, and is useful for accessing the field via an interface.
func (v *InternalUserInput) GetGoid() string { return v.Goid }

// GetReplaceGroupMemberships returns InternalUserInput.ReplaceGroupMemberships, and is useful for accessing the field via an interface.
func (v *InternalUserInput) GetReplaceGroupMemberships() bool { return v.ReplaceGroupMemberships }

// GetMemberOf returns InternalUserInput.MemberOf, and is useful for accessing the field via an interface.
func (v *InternalUserInput) GetMemberOf() []*MembershipInput { return v.MemberOf }

// GetLogin returns InternalUserInput.Login, and is useful for accessing the field via an interface.
func (v *InternalUserInput) GetLogin() string { return v.Login }

// GetPassword returns InternalUserInput.Password, and is useful for accessing the field via an interface.
func (v *InternalUserInput) GetPassword() string { return v.Password }

// GetCertBase64 returns InternalUserInput.CertBase64, and is useful for accessing the field via an interface.
func (v *InternalUserInput) GetCertBase64() string { return v.CertBase64 }

// GetSshPublicKey returns InternalUserInput.SshPublicKey, and is useful for accessing the field via an interface.
func (v *InternalUserInput) GetSshPublicKey() string { return v.SshPublicKey }

// GetFirstName returns InternalUserInput.FirstName, and is useful for accessing the field via an interface.
func (v *InternalUserInput) GetFirstName() string { return v.FirstName }

// GetLastName returns InternalUserInput.LastName, and is useful for accessing the field via an interface.
func (v *InternalUserInput) GetLastName() string { return v.LastName }

// GetEmail returns InternalUserInput.Email, and is useful for accessing the field via an interface.
func (v *InternalUserInput) GetEmail() string { return v.Email }

// GetEnabled returns InternalUserInput.Enabled, and is useful for accessing the field via an interface.
func (v *InternalUserInput) GetEnabled() bool { return v.Enabled }

// GetChecksum returns InternalUserInput.Checksum, and is useful for accessing the field via an interface.
func (v *InternalUserInput) GetChecksum() string { return v.Checksum }

type JdbcConnectionInput struct {
	// The internal entity unique identifier
	Goid string `json:"goid"`
	// The JDBC Connection name
	Name string `json:"name"`
	// The JDBC driver class name
	DriverClass string `json:"driverClass"`
	// The JDBC url
	JdbcUrl string `json:"jdbcUrl"`
	// Whether this JDBC connection is enabled
	Enabled bool `json:"enabled"`
	// The username
	Username string `json:"username"`
	// The password or the secured password reference
	Password string `json:"password"`
	// The minimum connection pool size
	MinPoolSize int `json:"minPoolSize"`
	// The maximum connection pool size
	MaxPoolSize int `json:"maxPoolSize"`
	// The JDBC connection properties excluding 'user' and 'password'
	Properties []*EntityPropertyInput `json:"properties,omitempty"`
	// Ignored at creation time but can be used to compare bundle with gw state
	Checksum string `json:"checksum"`
}

// GetGoid returns JdbcConnectionInput.Goid, and is useful for accessing the field via an interface.
func (v *JdbcConnectionInput) GetGoid() string { return v.Goid }

// GetName returns JdbcConnectionInput.Name, and is useful for accessing the field via an interface.
func (v *JdbcConnectionInput) GetName() string { return v.Name }

// GetDriverClass returns JdbcConnectionInput.DriverClass, and is useful for accessing the field via an interface.
func (v *JdbcConnectionInput) GetDriverClass() string { return v.DriverClass }

// GetJdbcUrl returns JdbcConnectionInput.JdbcUrl, and is useful for accessing the field via an interface.
func (v *JdbcConnectionInput) GetJdbcUrl() string { return v.JdbcUrl }

// GetEnabled returns JdbcConnectionInput.Enabled, and is useful for accessing the field via an interface.
func (v *JdbcConnectionInput) GetEnabled() bool { return v.Enabled }

// GetUsername returns JdbcConnectionInput.Username, and is useful for accessing the field via an interface.
func (v *JdbcConnectionInput) GetUsername() string { return v.Username }

// GetPassword returns JdbcConnectionInput.Password, and is useful for accessing the field via an interface.
func (v *JdbcConnectionInput) GetPassword() string { return v.Password }

// GetMinPoolSize returns JdbcConnectionInput.MinPoolSize, and is useful for accessing the field via an interface.
func (v *JdbcConnectionInput) GetMinPoolSize() int { return v.MinPoolSize }

// GetMaxPoolSize returns JdbcConnectionInput.MaxPoolSize, and is useful for accessing the field via an interface.
func (v *JdbcConnectionInput) GetMaxPoolSize() int { return v.MaxPoolSize }

// GetProperties returns JdbcConnectionInput.Properties, and is useful for accessing the field via an interface.
func (v *JdbcConnectionInput) GetProperties() []*EntityPropertyInput { return v.Properties }

// GetChecksum returns JdbcConnectionInput.Checksum, and is useful for accessing the field via an interface.
func (v *JdbcConnectionInput) GetChecksum() string { return v.Checksum }

type JmsDestinationInput struct {
	// The internal entity unique identifier
	Goid           string `json:"goid"`
	ConnectionGoid string `json:"connectionGoid"`
	// The JMS Destination name
	Name string `json:"name"`
	// The JMS Destination direction (inbound or outbound)
	Direction string `json:"direction"`
	// The JMS provider type
	ProviderType string `json:"providerType"`
	// The initial context factory class name
	InitialContextFactoryClassname string `json:"initialContextFactoryClassname"`
	// The connection factory name
	ConnectionFactoryName string `json:"connectionFactoryName"`
	// The JNDI URL
	JndiUrl string `json:"jndiUrl"`
	// The JNDI username
	JndiUsername string `json:"jndiUsername"`
	// The JNDI password
	JndiPassword string `json:"jndiPassword"`
	// The JNDI SSL details
	JndiSslDetails *JmsSslDetailsInput `json:"jndiSslDetails,omitempty"`
	// The destination type
	DestinationType string `json:"destinationType"`
	// The destination name
	DestinationName string `json:"destinationName"`
	// The username for destination connection
	DestinationUsername string `json:"destinationUsername"`
	// The password for destination connection
	DestinationPassword string `json:"destinationPassword"`
	// The destination SSL details
	DestinationSslDetails *JmsSslDetailsInput `json:"destinationSslDetails,omitempty"`
	// Whether this JMS destination is template
	Template bool `json:"template"`
	// Whether this JMS destination is enabled
	Enabled bool `json:"enabled"`
	// The remaining JMS Destination properties that include inbound options or outbound options or additional properties
	Properties []*EntityPropertyInput `json:"properties,omitempty"`
	// Ignored at creation time but can be used to compare bundle with gw state
	Checksum string `json:"checksum"`
}

// GetGoid returns JmsDestinationInput.Goid, and is useful for accessing the field via an interface.
func (v *JmsDestinationInput) GetGoid() string { return v.Goid }

// GetConnectionGoid returns JmsDestinationInput.ConnectionGoid, and is useful for accessing the field via an interface.
func (v *JmsDestinationInput) GetConnectionGoid() string { return v.ConnectionGoid }

// GetName returns JmsDestinationInput.Name, and is useful for accessing the field via an interface.
func (v *JmsDestinationInput) GetName() string { return v.Name }

// GetDirection returns JmsDestinationInput.Direction, and is useful for accessing the field via an interface.
func (v *JmsDestinationInput) GetDirection() string { return v.Direction }

// GetProviderType returns JmsDestinationInput.ProviderType, and is useful for accessing the field via an interface.
func (v *JmsDestinationInput) GetProviderType() string { return v.ProviderType }

// GetInitialContextFactoryClassname returns JmsDestinationInput.InitialContextFactoryClassname, and is useful for accessing the field via an interface.
func (v *JmsDestinationInput) GetInitialContextFactoryClassname() string {
	return v.InitialContextFactoryClassname
}

// GetConnectionFactoryName returns JmsDestinationInput.ConnectionFactoryName, and is useful for accessing the field via an interface.
func (v *JmsDestinationInput) GetConnectionFactoryName() string { return v.ConnectionFactoryName }

// GetJndiUrl returns JmsDestinationInput.JndiUrl, and is useful for accessing the field via an interface.
func (v *JmsDestinationInput) GetJndiUrl() string { return v.JndiUrl }

// GetJndiUsername returns JmsDestinationInput.JndiUsername, and is useful for accessing the field via an interface.
func (v *JmsDestinationInput) GetJndiUsername() string { return v.JndiUsername }

// GetJndiPassword returns JmsDestinationInput.JndiPassword, and is useful for accessing the field via an interface.
func (v *JmsDestinationInput) GetJndiPassword() string { return v.JndiPassword }

// GetJndiSslDetails returns JmsDestinationInput.JndiSslDetails, and is useful for accessing the field via an interface.
func (v *JmsDestinationInput) GetJndiSslDetails() *JmsSslDetailsInput { return v.JndiSslDetails }

// GetDestinationType returns JmsDestinationInput.DestinationType, and is useful for accessing the field via an interface.
func (v *JmsDestinationInput) GetDestinationType() string { return v.DestinationType }

// GetDestinationName returns JmsDestinationInput.DestinationName, and is useful for accessing the field via an interface.
func (v *JmsDestinationInput) GetDestinationName() string { return v.DestinationName }

// GetDestinationUsername returns JmsDestinationInput.DestinationUsername, and is useful for accessing the field via an interface.
func (v *JmsDestinationInput) GetDestinationUsername() string { return v.DestinationUsername }

// GetDestinationPassword returns JmsDestinationInput.DestinationPassword, and is useful for accessing the field via an interface.
func (v *JmsDestinationInput) GetDestinationPassword() string { return v.DestinationPassword }

// GetDestinationSslDetails returns JmsDestinationInput.DestinationSslDetails, and is useful for accessing the field via an interface.
func (v *JmsDestinationInput) GetDestinationSslDetails() *JmsSslDetailsInput {
	return v.DestinationSslDetails
}

// GetTemplate returns JmsDestinationInput.Template, and is useful for accessing the field via an interface.
func (v *JmsDestinationInput) GetTemplate() bool { return v.Template }

// GetEnabled returns JmsDestinationInput.Enabled, and is useful for accessing the field via an interface.
func (v *JmsDestinationInput) GetEnabled() bool { return v.Enabled }

// GetProperties returns JmsDestinationInput.Properties, and is useful for accessing the field via an interface.
func (v *JmsDestinationInput) GetProperties() []*EntityPropertyInput { return v.Properties }

// GetChecksum returns JmsDestinationInput.Checksum, and is useful for accessing the field via an interface.
func (v *JmsDestinationInput) GetChecksum() string { return v.Checksum }

type JmsSslDetailsInput struct {
	// Whether SSL is enabled
	SslEnabled bool `json:"sslEnabled"`
	// Whether SSL is used for Authentication only
	SslForAuthenticationOnly bool `json:"sslForAuthenticationOnly"`
	// Whether SSL Server Certificate is to be verified
	SslVerifyServerCertificate bool `json:"sslVerifyServerCertificate"`
	// Whether SSL Server Hostname is to be verified
	SslVerifyServerHostname bool `json:"sslVerifyServerHostname"`
	// Private Key Alias for SSL Client Authentication
	SslClientKeyAlias string `json:"sslClientKeyAlias"`
}

// GetSslEnabled returns JmsSslDetailsInput.SslEnabled, and is useful for accessing the field via an interface.
func (v *JmsSslDetailsInput) GetSslEnabled() bool { return v.SslEnabled }

// GetSslForAuthenticationOnly returns JmsSslDetailsInput.SslForAuthenticationOnly, and is useful for accessing the field via an interface.
func (v *JmsSslDetailsInput) GetSslForAuthenticationOnly() bool { return v.SslForAuthenticationOnly }

// GetSslVerifyServerCertificate returns JmsSslDetailsInput.SslVerifyServerCertificate, and is useful for accessing the field via an interface.
func (v *JmsSslDetailsInput) GetSslVerifyServerCertificate() bool {
	return v.SslVerifyServerCertificate
}

// GetSslVerifyServerHostname returns JmsSslDetailsInput.SslVerifyServerHostname, and is useful for accessing the field via an interface.
func (v *JmsSslDetailsInput) GetSslVerifyServerHostname() bool { return v.SslVerifyServerHostname }

// GetSslClientKeyAlias returns JmsSslDetailsInput.SslClientKeyAlias, and is useful for accessing the field via an interface.
func (v *JmsSslDetailsInput) GetSslClientKeyAlias() string { return v.SslClientKeyAlias }

// Defines a current status of a given scheduled task
type JobStatus string

const (
	JobStatusScheduled JobStatus = "SCHEDULED"
	JobStatusCompleted JobStatus = "COMPLETED"
	JobStatusDisabled  JobStatus = "DISABLED"
)

// Defines a scheduled task type
type JobType string

const (
	JobTypeOneTime   JobType = "ONE_TIME"
	JobTypeRecurring JobType = "RECURRING"
)

type KeyInput struct {
	KeystoreId string `json:"keystoreId"`
	Alias      string `json:"alias"`
	// Base64 encoded PKCS12 keystore containing the private key and cert chain for the key entry.
	// The keystore is password-protected using the transaction encryption passphrase provided.
	P12 string `json:"p12"`
	// The private key data in PEM format
	Pem string `json:"pem"`
	// Will try to match at creation time is specified
	Goid string `json:"goid"`
	// SubjectDn of the certificate associated with the key. (Note that, this field has no effect on the mutation)
	SubjectDn string `json:"subjectDn"`
	// Key type. (Note that, this field has no effect on the mutation)
	KeyType string `json:"keyType"`
	// The Key usage types. (Note that, the key usage will not be reset when this field is not specified)
	UsageTypes []KeyUsageType `json:"usageTypes"`
	// The certificate chain in PEM format. (Note that, this field has no effect on the mutation)
	CertChain interface{} `json:"certChain"`
	// Ignored at entity creation time but declared here so you can embed checksums in graphman bundles
	Checksum string `json:"checksum"`
}

// GetKeystoreId returns KeyInput.KeystoreId, and is useful for accessing the field via an interface.
func (v *KeyInput) GetKeystoreId() string { return v.KeystoreId }

// GetAlias returns KeyInput.Alias, and is useful for accessing the field via an interface.
func (v *KeyInput) GetAlias() string { return v.Alias }

// GetP12 returns KeyInput.P12, and is useful for accessing the field via an interface.
func (v *KeyInput) GetP12() string { return v.P12 }

// GetPem returns KeyInput.Pem, and is useful for accessing the field via an interface.
func (v *KeyInput) GetPem() string { return v.Pem }

// GetGoid returns KeyInput.Goid, and is useful for accessing the field via an interface.
func (v *KeyInput) GetGoid() string { return v.Goid }

// GetSubjectDn returns KeyInput.SubjectDn, and is useful for accessing the field via an interface.
func (v *KeyInput) GetSubjectDn() string { return v.SubjectDn }

// GetKeyType returns KeyInput.KeyType, and is useful for accessing the field via an interface.
func (v *KeyInput) GetKeyType() string { return v.KeyType }

// GetUsageTypes returns KeyInput.UsageTypes, and is useful for accessing the field via an interface.
func (v *KeyInput) GetUsageTypes() []KeyUsageType { return v.UsageTypes }

// GetCertChain returns KeyInput.CertChain, and is useful for accessing the field via an interface.
func (v *KeyInput) GetCertChain() interface{} { return v.CertChain }

// GetChecksum returns KeyInput.Checksum, and is useful for accessing the field via an interface.
func (v *KeyInput) GetChecksum() string { return v.Checksum }

type KeyUsageType string

const (
	// Represents a key marked as the default SSL key
	KeyUsageTypeSsl KeyUsageType = "SSL"
	// Represents a key marked as the default CA key
	KeyUsageTypeCa KeyUsageType = "CA"
	// Represents a key marked as the default audit viewer/decryption key
	KeyUsageTypeAuditViewer KeyUsageType = "AUDIT_VIEWER"
	// Represents a key marked as the default audit signing key
	KeyUsageTypeAuditSigning KeyUsageType = "AUDIT_SIGNING"
)

type L7PolicyInput struct {
	// The internal entity unique identifier
	Goid string `json:"goid"`
	// The folder path where to create this policy.  If the path does not exist, it will be created
	FolderPath string `json:"folderPath"`
	// The name of the policy. Policies are unique by name.
	Name string `json:"name"`
	// The guid for this policy, if none provided, assigned at creation
	Guid string `json:"guid"`
	// The policy
	Policy          *PolicyInput           `json:"policy,omitempty"`
	PolicyRevision  *PolicyRevisionInput   `json:"policyRevision,omitempty"`
	PolicyRevisions []*PolicyRevisionInput `json:"policyRevisions,omitempty"`
	Soap            bool                   `json:"soap"`
	PolicyType      L7PolicyType           `json:"policyType"`
	Tag             string                 `json:"tag,omitempty"`
	SubTag          string                 `json:"subTag,omitempty"`
	// Ignored at creation time but can be used to compare bundle with gw state
	Checksum string `json:"checksum,omitempty"`
}

// GetGoid returns L7PolicyInput.Goid, and is useful for accessing the field via an interface.
func (v *L7PolicyInput) GetGoid() string { return v.Goid }

// GetFolderPath returns L7PolicyInput.FolderPath, and is useful for accessing the field via an interface.
func (v *L7PolicyInput) GetFolderPath() string { return v.FolderPath }

// GetName returns L7PolicyInput.Name, and is useful for accessing the field via an interface.
func (v *L7PolicyInput) GetName() string { return v.Name }

// GetGuid returns L7PolicyInput.Guid, and is useful for accessing the field via an interface.
func (v *L7PolicyInput) GetGuid() string { return v.Guid }

// GetPolicy returns L7PolicyInput.Policy, and is useful for accessing the field via an interface.
func (v *L7PolicyInput) GetPolicy() *PolicyInput { return v.Policy }

// GetPolicyRevision returns L7PolicyInput.PolicyRevision, and is useful for accessing the field via an interface.
func (v *L7PolicyInput) GetPolicyRevision() *PolicyRevisionInput { return v.PolicyRevision }

// GetPolicyRevisions returns L7PolicyInput.PolicyRevisions, and is useful for accessing the field via an interface.
func (v *L7PolicyInput) GetPolicyRevisions() []*PolicyRevisionInput { return v.PolicyRevisions }

// GetSoap returns L7PolicyInput.Soap, and is useful for accessing the field via an interface.
func (v *L7PolicyInput) GetSoap() bool { return v.Soap }

// GetPolicyType returns L7PolicyInput.PolicyType, and is useful for accessing the field via an interface.
func (v *L7PolicyInput) GetPolicyType() L7PolicyType { return v.PolicyType }

// GetTag returns L7PolicyInput.Tag, and is useful for accessing the field via an interface.
func (v *L7PolicyInput) GetTag() string { return v.Tag }

// GetSubTag returns L7PolicyInput.SubTag, and is useful for accessing the field via an interface.
func (v *L7PolicyInput) GetSubTag() string { return v.SubTag }

// GetChecksum returns L7PolicyInput.Checksum, and is useful for accessing the field via an interface.
func (v *L7PolicyInput) GetChecksum() string { return v.Checksum }

type L7PolicyType string

const (
	L7PolicyTypeFragment                      L7PolicyType = "FRAGMENT"
	L7PolicyTypePreRoutingFragment            L7PolicyType = "PRE_ROUTING_FRAGMENT"
	L7PolicyTypeSuccessfulRoutingFragment     L7PolicyType = "SUCCESSFUL_ROUTING_FRAGMENT"
	L7PolicyTypeFailedRoutingFragment         L7PolicyType = "FAILED_ROUTING_FRAGMENT"
	L7PolicyTypeAuthenticationSuccessFragment L7PolicyType = "AUTHENTICATION_SUCCESS_FRAGMENT"
	L7PolicyTypeAuthenticationFailureFragment L7PolicyType = "AUTHENTICATION_FAILURE_FRAGMENT"
	L7PolicyTypeAuthorizationSuccessFragment  L7PolicyType = "AUTHORIZATION_SUCCESS_FRAGMENT"
	L7PolicyTypeAuthorizationFailureFragment  L7PolicyType = "AUTHORIZATION_FAILURE_FRAGMENT"
	L7PolicyTypeGlobal                        L7PolicyType = "GLOBAL"
	L7PolicyTypeInternal                      L7PolicyType = "INTERNAL"
	L7PolicyTypePolicyBackedIdp               L7PolicyType = "POLICY_BACKED_IDP"
	L7PolicyTypePolicyBackedOperation         L7PolicyType = "POLICY_BACKED_OPERATION"
	L7PolicyTypePolicyBackedBackgroundTask    L7PolicyType = "POLICY_BACKED_BACKGROUND_TASK"
	L7PolicyTypePolicyBackedServiceMetrics    L7PolicyType = "POLICY_BACKED_SERVICE_METRICS"
)

type L7ServiceInput struct {
	// The goid for this service
	Goid string `json:"goid"`
	// The guid for this service
	Guid string `json:"guid"`
	// The name of the service
	Name string `json:"name"`
	// The resolution path to the service
	ResolutionPath string `json:"resolutionPath"`
	// The service resolvers. They can be used to identify services.
	Resolvers *ServiceResolversInput `json:"resolvers,omitempty"`
	// The type of service
	ServiceType L7ServiceType `json:"serviceType"`
	// The configuration checksum
	Checksum string `json:"checksum"`
	// Whether or not the published service is enabled
	Enabled bool `json:"enabled"`
	// The folder path where to create this service.  If the path does not exist, it will be created
	FolderPath string `json:"folderPath"`
	// Which SOAP version
	SoapVersion SoapVersion `json:"soapVersion,omitempty"`
	// Which http methods are allowed
	MethodsAllowed       []HttpMethod `json:"methodsAllowed"`
	TracingEnabled       bool         `json:"tracingEnabled"`
	WssProcessingEnabled bool         `json:"wssProcessingEnabled,omitempty"`
	// Allow requests intended for operations not supported by the WSDL
	LaxResolution bool                   `json:"laxResolution,omitempty"`
	Properties    []*EntityPropertyInput `json:"properties,omitempty"`
	// The WSDL of the soap service
	Wsdl string `json:"wsdl,omitempty"`
	// URL for the protected service WSDL document
	WsdlUrl string `json:"wsdlUrl"`
	// One or more additional WSDL resources
	WsdlResources []*ServiceResourceInput `json:"wsdlResources,omitempty"`
	// The service policy
	Policy *PolicyInput `json:"policy,omitempty"`
	// This will be ignored during the mutation
	PolicyRevision *PolicyRevisionInput `json:"policyRevision,omitempty"`
	// This will be ignored during the mutation
	PolicyRevisions []*PolicyRevisionInput `json:"policyRevisions,omitempty"`
}

// GetGoid returns L7ServiceInput.Goid, and is useful for accessing the field via an interface.
func (v *L7ServiceInput) GetGoid() string { return v.Goid }

// GetGuid returns L7ServiceInput.Guid, and is useful for accessing the field via an interface.
func (v *L7ServiceInput) GetGuid() string { return v.Guid }

// GetName returns L7ServiceInput.Name, and is useful for accessing the field via an interface.
func (v *L7ServiceInput) GetName() string { return v.Name }

// GetResolutionPath returns L7ServiceInput.ResolutionPath, and is useful for accessing the field via an interface.
func (v *L7ServiceInput) GetResolutionPath() string { return v.ResolutionPath }

// GetResolvers returns L7ServiceInput.Resolvers, and is useful for accessing the field via an interface.
func (v *L7ServiceInput) GetResolvers() *ServiceResolversInput { return v.Resolvers }

// GetServiceType returns L7ServiceInput.ServiceType, and is useful for accessing the field via an interface.
func (v *L7ServiceInput) GetServiceType() L7ServiceType { return v.ServiceType }

// GetChecksum returns L7ServiceInput.Checksum, and is useful for accessing the field via an interface.
func (v *L7ServiceInput) GetChecksum() string { return v.Checksum }

// GetEnabled returns L7ServiceInput.Enabled, and is useful for accessing the field via an interface.
func (v *L7ServiceInput) GetEnabled() bool { return v.Enabled }

// GetFolderPath returns L7ServiceInput.FolderPath, and is useful for accessing the field via an interface.
func (v *L7ServiceInput) GetFolderPath() string { return v.FolderPath }

// GetSoapVersion returns L7ServiceInput.SoapVersion, and is useful for accessing the field via an interface.
func (v *L7ServiceInput) GetSoapVersion() SoapVersion { return v.SoapVersion }

// GetMethodsAllowed returns L7ServiceInput.MethodsAllowed, and is useful for accessing the field via an interface.
func (v *L7ServiceInput) GetMethodsAllowed() []HttpMethod { return v.MethodsAllowed }

// GetTracingEnabled returns L7ServiceInput.TracingEnabled, and is useful for accessing the field via an interface.
func (v *L7ServiceInput) GetTracingEnabled() bool { return v.TracingEnabled }

// GetWssProcessingEnabled returns L7ServiceInput.WssProcessingEnabled, and is useful for accessing the field via an interface.
func (v *L7ServiceInput) GetWssProcessingEnabled() bool { return v.WssProcessingEnabled }

// GetLaxResolution returns L7ServiceInput.LaxResolution, and is useful for accessing the field via an interface.
func (v *L7ServiceInput) GetLaxResolution() bool { return v.LaxResolution }

// GetProperties returns L7ServiceInput.Properties, and is useful for accessing the field via an interface.
func (v *L7ServiceInput) GetProperties() []*EntityPropertyInput { return v.Properties }

// GetWsdl returns L7ServiceInput.Wsdl, and is useful for accessing the field via an interface.
func (v *L7ServiceInput) GetWsdl() string { return v.Wsdl }

// GetWsdlUrl returns L7ServiceInput.WsdlUrl, and is useful for accessing the field via an interface.
func (v *L7ServiceInput) GetWsdlUrl() string { return v.WsdlUrl }

// GetWsdlResources returns L7ServiceInput.WsdlResources, and is useful for accessing the field via an interface.
func (v *L7ServiceInput) GetWsdlResources() []*ServiceResourceInput { return v.WsdlResources }

// GetPolicy returns L7ServiceInput.Policy, and is useful for accessing the field via an interface.
func (v *L7ServiceInput) GetPolicy() *PolicyInput { return v.Policy }

// GetPolicyRevision returns L7ServiceInput.PolicyRevision, and is useful for accessing the field via an interface.
func (v *L7ServiceInput) GetPolicyRevision() *PolicyRevisionInput { return v.PolicyRevision }

// GetPolicyRevisions returns L7ServiceInput.PolicyRevisions, and is useful for accessing the field via an interface.
func (v *L7ServiceInput) GetPolicyRevisions() []*PolicyRevisionInput { return v.PolicyRevisions }

type L7ServiceType string

const (
	L7ServiceTypeWebApi         L7ServiceType = "WEB_API"
	L7ServiceTypeSoap           L7ServiceType = "SOAP"
	L7ServiceTypeInternalWebApi L7ServiceType = "INTERNAL_WEB_API"
	L7ServiceTypeInternalSoap   L7ServiceType = "INTERNAL_SOAP"
)

type LdapIdpInput struct {
	Goid     string `json:"goid"`
	Name     string `json:"name"`
	Checksum string `json:"checksum"`
	// Ldap type
	LdapType string `json:"ldapType"`
	// Ldap server urls
	ServerUrls []string `json:"serverUrls"`
	// Whether or not the gateway presents a client cert when connecting at those ldap urls (only relevant when ldaps url)
	UseSslClientAuth bool `json:"useSslClientAuth"`
	// The alias of the key in the gateway keystore that is used when doing ldaps client cert authentication
	SslClientKeyAlias   string                 `json:"sslClientKeyAlias"`
	SearchBase          string                 `json:"searchBase"`
	BindDn              string                 `json:"bindDn"`
	BindPassword        string                 `json:"bindPassword"`
	Writable            bool                   `json:"writable"`
	WriteBase           string                 `json:"writeBase"`
	SpecifiedAttributes []string               `json:"specifiedAttributes"`
	UserMappings        []*UserMappingInput    `json:"userMappings,omitempty"`
	GroupMappings       []*GroupMappingInput   `json:"groupMappings,omitempty"`
	NtlmProperties      []*EntityPropertyInput `json:"ntlmProperties,omitempty"`
	// Additional properties
	Properties []*EntityPropertyInput `json:"properties,omitempty"`
}

// GetGoid returns LdapIdpInput.Goid, and is useful for accessing the field via an interface.
func (v *LdapIdpInput) GetGoid() string { return v.Goid }

// GetName returns LdapIdpInput.Name, and is useful for accessing the field via an interface.
func (v *LdapIdpInput) GetName() string { return v.Name }

// GetChecksum returns LdapIdpInput.Checksum, and is useful for accessing the field via an interface.
func (v *LdapIdpInput) GetChecksum() string { return v.Checksum }

// GetLdapType returns LdapIdpInput.LdapType, and is useful for accessing the field via an interface.
func (v *LdapIdpInput) GetLdapType() string { return v.LdapType }

// GetServerUrls returns LdapIdpInput.ServerUrls, and is useful for accessing the field via an interface.
func (v *LdapIdpInput) GetServerUrls() []string { return v.ServerUrls }

// GetUseSslClientAuth returns LdapIdpInput.UseSslClientAuth, and is useful for accessing the field via an interface.
func (v *LdapIdpInput) GetUseSslClientAuth() bool { return v.UseSslClientAuth }

// GetSslClientKeyAlias returns LdapIdpInput.SslClientKeyAlias, and is useful for accessing the field via an interface.
func (v *LdapIdpInput) GetSslClientKeyAlias() string { return v.SslClientKeyAlias }

// GetSearchBase returns LdapIdpInput.SearchBase, and is useful for accessing the field via an interface.
func (v *LdapIdpInput) GetSearchBase() string { return v.SearchBase }

// GetBindDn returns LdapIdpInput.BindDn, and is useful for accessing the field via an interface.
func (v *LdapIdpInput) GetBindDn() string { return v.BindDn }

// GetBindPassword returns LdapIdpInput.BindPassword, and is useful for accessing the field via an interface.
func (v *LdapIdpInput) GetBindPassword() string { return v.BindPassword }

// GetWritable returns LdapIdpInput.Writable, and is useful for accessing the field via an interface.
func (v *LdapIdpInput) GetWritable() bool { return v.Writable }

// GetWriteBase returns LdapIdpInput.WriteBase, and is useful for accessing the field via an interface.
func (v *LdapIdpInput) GetWriteBase() string { return v.WriteBase }

// GetSpecifiedAttributes returns LdapIdpInput.SpecifiedAttributes, and is useful for accessing the field via an interface.
func (v *LdapIdpInput) GetSpecifiedAttributes() []string { return v.SpecifiedAttributes }

// GetUserMappings returns LdapIdpInput.UserMappings, and is useful for accessing the field via an interface.
func (v *LdapIdpInput) GetUserMappings() []*UserMappingInput { return v.UserMappings }

// GetGroupMappings returns LdapIdpInput.GroupMappings, and is useful for accessing the field via an interface.
func (v *LdapIdpInput) GetGroupMappings() []*GroupMappingInput { return v.GroupMappings }

// GetNtlmProperties returns LdapIdpInput.NtlmProperties, and is useful for accessing the field via an interface.
func (v *LdapIdpInput) GetNtlmProperties() []*EntityPropertyInput { return v.NtlmProperties }

// GetProperties returns LdapIdpInput.Properties, and is useful for accessing the field via an interface.
func (v *LdapIdpInput) GetProperties() []*EntityPropertyInput { return v.Properties }

type LdapInput struct {
	Name string `json:"name"`
	// Will try to match goid if provided
	Goid     string   `json:"goid"`
	LdapUrls []string `json:"ldapUrls"`
	// Whether or not the gateway presents a client cert when connecting at those ldap urls (only relevant when ldaps url)
	LdapsClientAuthEnabled bool `json:"ldapsClientAuthEnabled"`
	// The ID of the gateway keystore where the key is located
	LdapsClientKeystoreId string `json:"ldapsClientKeystoreId"`
	// The alias of the key in the gateway keystore that is used when doing ldaps client cert authentication
	LdapsClientKeyAlias string               `json:"ldapsClientKeyAlias"`
	SearchBase          string               `json:"searchBase"`
	Writable            bool                 `json:"writable"`
	BindDn              string               `json:"bindDn"`
	BindPassword        string               `json:"bindPassword"`
	UserMappings        []*UserMappingInput  `json:"userMappings,omitempty"`
	GroupMappings       []*GroupMappingInput `json:"groupMappings,omitempty"`
	// The optional checksum is ignored during the mutation but can be used to compare bundle content
	Checksum string `json:"checksum"`
}

// GetName returns LdapInput.Name, and is useful for accessing the field via an interface.
func (v *LdapInput) GetName() string { return v.Name }

// GetGoid returns LdapInput.Goid, and is useful for accessing the field via an interface.
func (v *LdapInput) GetGoid() string { return v.Goid }

// GetLdapUrls returns LdapInput.LdapUrls, and is useful for accessing the field via an interface.
func (v *LdapInput) GetLdapUrls() []string { return v.LdapUrls }

// GetLdapsClientAuthEnabled returns LdapInput.LdapsClientAuthEnabled, and is useful for accessing the field via an interface.
func (v *LdapInput) GetLdapsClientAuthEnabled() bool { return v.LdapsClientAuthEnabled }

// GetLdapsClientKeystoreId returns LdapInput.LdapsClientKeystoreId, and is useful for accessing the field via an interface.
func (v *LdapInput) GetLdapsClientKeystoreId() string { return v.LdapsClientKeystoreId }

// GetLdapsClientKeyAlias returns LdapInput.LdapsClientKeyAlias, and is useful for accessing the field via an interface.
func (v *LdapInput) GetLdapsClientKeyAlias() string { return v.LdapsClientKeyAlias }

// GetSearchBase returns LdapInput.SearchBase, and is useful for accessing the field via an interface.
func (v *LdapInput) GetSearchBase() string { return v.SearchBase }

// GetWritable returns LdapInput.Writable, and is useful for accessing the field via an interface.
func (v *LdapInput) GetWritable() bool { return v.Writable }

// GetBindDn returns LdapInput.BindDn, and is useful for accessing the field via an interface.
func (v *LdapInput) GetBindDn() string { return v.BindDn }

// GetBindPassword returns LdapInput.BindPassword, and is useful for accessing the field via an interface.
func (v *LdapInput) GetBindPassword() string { return v.BindPassword }

// GetUserMappings returns LdapInput.UserMappings, and is useful for accessing the field via an interface.
func (v *LdapInput) GetUserMappings() []*UserMappingInput { return v.UserMappings }

// GetGroupMappings returns LdapInput.GroupMappings, and is useful for accessing the field via an interface.
func (v *LdapInput) GetGroupMappings() []*GroupMappingInput { return v.GroupMappings }

// GetChecksum returns LdapInput.Checksum, and is useful for accessing the field via an interface.
func (v *LdapInput) GetChecksum() string { return v.Checksum }

type ListenPortClientAuth string

const (
	ListenPortClientAuthNone     ListenPortClientAuth = "NONE"
	ListenPortClientAuthOptional ListenPortClientAuth = "OPTIONAL"
	ListenPortClientAuthRequired ListenPortClientAuth = "REQUIRED"
)

type ListenPortFeature string

const (
	ListenPortFeaturePublishedServiceMessageInput ListenPortFeature = "PUBLISHED_SERVICE_MESSAGE_INPUT"
	ListenPortFeaturePolicyManagerAccess          ListenPortFeature = "POLICY_MANAGER_ACCESS"
	ListenPortFeatureEnterpriseManagerAccess      ListenPortFeature = "ENTERPRISE_MANAGER_ACCESS"
	ListenPortFeatureAdministrativeAccess         ListenPortFeature = "ADMINISTRATIVE_ACCESS"
	ListenPortFeatureBrowserBasedAdministration   ListenPortFeature = "BROWSER_BASED_ADMINISTRATION"
	ListenPortFeaturePolicyDownloadService        ListenPortFeature = "POLICY_DOWNLOAD_SERVICE"
	ListenPortFeaturePingService                  ListenPortFeature = "PING_SERVICE"
	ListenPortFeatureWsTrustSecurityTokenService  ListenPortFeature = "WS_TRUST_SECURITY_TOKEN_SERVICE"
	ListenPortFeatureCertificateSigningService    ListenPortFeature = "CERTIFICATE_SIGNING_SERVICE"
	ListenPortFeaturePasswordChangingService      ListenPortFeature = "PASSWORD_CHANGING_SERVICE"
	ListenPortFeatureWsdlDownloadService          ListenPortFeature = "WSDL_DOWNLOAD_SERVICE"
	ListenPortFeatureSnmpQueryService             ListenPortFeature = "SNMP_QUERY_SERVICE"
	ListenPortFeatureBuiltInServices              ListenPortFeature = "BUILT_IN_SERVICES"
	ListenPortFeatureNodeControl                  ListenPortFeature = "NODE_CONTROL"
	ListenPortFeatureInterNodeCommunication       ListenPortFeature = "INTER_NODE_COMMUNICATION"
)

type ListenPortInput struct {
	// The internal entity unique identifier
	Goid string `json:"goid"`
	// The listen port configuration name
	Name string `json:"name"`
	// Whether this listen port configuration is enabled to listen for traffic on the specified port
	Enabled bool `json:"enabled"`
	// Protocol (scheme). Possible values are:
	// HTTP
	// HTTPS
	// HTTP2
	// HTTP2 (Secure)
	// FTP
	// FTPS
	// l7.raw.tcp
	// SSH2
	Protocol string `json:"protocol"`
	// The ListenPort's port number
	// Note: If the listen port is using the SSH2 protocol, avoid using port 22, as
	// it may conflict with the default SSH port 22 on Linux or Unix systems.
	Port int `json:"port"`
	// The name of the published service hardwired to the listen port
	HardwiredServiceName string `json:"hardwiredServiceName"`
	// Which Gateway services can be accessed through this listen port
	EnabledFeatures []ListenPortFeature `json:"enabledFeatures"`
	// The listen port tls settings
	TlsSettings *ListenPortTlsSettingsInput `json:"tlsSettings,omitempty"`
	// The listen port properties
	Properties []*EntityPropertyInput `json:"properties,omitempty"`
	// Ignored at creation time but can be used to compare bundle with gw state
	Checksum string `json:"checksum"`
}

// GetGoid returns ListenPortInput.Goid, and is useful for accessing the field via an interface.
func (v *ListenPortInput) GetGoid() string { return v.Goid }

// GetName returns ListenPortInput.Name, and is useful for accessing the field via an interface.
func (v *ListenPortInput) GetName() string { return v.Name }

// GetEnabled returns ListenPortInput.Enabled, and is useful for accessing the field via an interface.
func (v *ListenPortInput) GetEnabled() bool { return v.Enabled }

// GetProtocol returns ListenPortInput.Protocol, and is useful for accessing the field via an interface.
func (v *ListenPortInput) GetProtocol() string { return v.Protocol }

// GetPort returns ListenPortInput.Port, and is useful for accessing the field via an interface.
func (v *ListenPortInput) GetPort() int { return v.Port }

// GetHardwiredServiceName returns ListenPortInput.HardwiredServiceName, and is useful for accessing the field via an interface.
func (v *ListenPortInput) GetHardwiredServiceName() string { return v.HardwiredServiceName }

// GetEnabledFeatures returns ListenPortInput.EnabledFeatures, and is useful for accessing the field via an interface.
func (v *ListenPortInput) GetEnabledFeatures() []ListenPortFeature { return v.EnabledFeatures }

// GetTlsSettings returns ListenPortInput.TlsSettings, and is useful for accessing the field via an interface.
func (v *ListenPortInput) GetTlsSettings() *ListenPortTlsSettingsInput { return v.TlsSettings }

// GetProperties returns ListenPortInput.Properties, and is useful for accessing the field via an interface.
func (v *ListenPortInput) GetProperties() []*EntityPropertyInput { return v.Properties }

// GetChecksum returns ListenPortInput.Checksum, and is useful for accessing the field via an interface.
func (v *ListenPortInput) GetChecksum() string { return v.Checksum }

type ListenPortTlsSettingsInput struct {
	// Specify whether the client must present a certificate to authenticate: NONE/OPTIONAL/REQUIRED
	ClientAuthentication ListenPortClientAuth `json:"clientAuthentication"`
	// Keystore ID
	KeystoreId string `json:"keystoreId"`
	// Key alias configured for listen port
	KeyAlias string `json:"keyAlias"`
	// TLS versions to be enabled for the listen port
	TlsVersions []string `json:"tlsVersions"`
	// Cipher suites that will be enabled on the SSL listen port
	CipherSuites []string `json:"cipherSuites"`
	// Enforces cipher suites usage in the order of preference
	UseCipherSuitesOrder bool `json:"useCipherSuitesOrder"`
}

// GetClientAuthentication returns ListenPortTlsSettingsInput.ClientAuthentication, and is useful for accessing the field via an interface.
func (v *ListenPortTlsSettingsInput) GetClientAuthentication() ListenPortClientAuth {
	return v.ClientAuthentication
}

// GetKeystoreId returns ListenPortTlsSettingsInput.KeystoreId, and is useful for accessing the field via an interface.
func (v *ListenPortTlsSettingsInput) GetKeystoreId() string { return v.KeystoreId }

// GetKeyAlias returns ListenPortTlsSettingsInput.KeyAlias, and is useful for accessing the field via an interface.
func (v *ListenPortTlsSettingsInput) GetKeyAlias() string { return v.KeyAlias }

// GetTlsVersions returns ListenPortTlsSettingsInput.TlsVersions, and is useful for accessing the field via an interface.
func (v *ListenPortTlsSettingsInput) GetTlsVersions() []string { return v.TlsVersions }

// GetCipherSuites returns ListenPortTlsSettingsInput.CipherSuites, and is useful for accessing the field via an interface.
func (v *ListenPortTlsSettingsInput) GetCipherSuites() []string { return v.CipherSuites }

// GetUseCipherSuitesOrder returns ListenPortTlsSettingsInput.UseCipherSuitesOrder, and is useful for accessing the field via an interface.
func (v *ListenPortTlsSettingsInput) GetUseCipherSuitesOrder() bool { return v.UseCipherSuitesOrder }

// Indicates severity threshold of the log sink
type LogSeverityThreshold string

const (
	LogSeverityThresholdAll     LogSeverityThreshold = "ALL"
	LogSeverityThresholdFinest  LogSeverityThreshold = "FINEST"
	LogSeverityThresholdFiner   LogSeverityThreshold = "FINER"
	LogSeverityThresholdFine    LogSeverityThreshold = "FINE"
	LogSeverityThresholdConfig  LogSeverityThreshold = "CONFIG"
	LogSeverityThresholdInfo    LogSeverityThreshold = "INFO"
	LogSeverityThresholdWarning LogSeverityThreshold = "WARNING"
	LogSeverityThresholdSevere  LogSeverityThreshold = "SEVERE"
)

// Indicates the Sink Category
type LogSinkCategory string

const (
	LogSinkCategoryLog     LogSinkCategory = "LOG"
	LogSinkCategoryTraffic LogSinkCategory = "TRAFFIC"
	LogSinkCategoryAudit   LogSinkCategory = "AUDIT"
	LogSinkCategorySspc    LogSinkCategory = "SSPC"
)

// Indicate the long sink filter, consist of a type and list of values
type LogSinkFilterInput struct {
	// defines the type of log sink
	Type string `json:"type"`
	// defines the list of values
	Values []string `json:"values"`
}

// GetType returns LogSinkFilterInput.Type, and is useful for accessing the field via an interface.
func (v *LogSinkFilterInput) GetType() string { return v.Type }

// GetValues returns LogSinkFilterInput.Values, and is useful for accessing the field via an interface.
func (v *LogSinkFilterInput) GetValues() []string { return v.Values }

type LogSinkInput struct {
	// The internal entity unique identifier
	Goid string `json:"goid"`
	// log sink unique name
	Name string `json:"name"`
	// description of log sink
	Description string `json:"description"`
	// defines whether its a file based log or sysLog
	Type LogSinkType `json:"type"`
	// Whether this log sink is enabled
	Enabled bool `json:"enabled"`
	// defines the severity threshold of log Sink
	Severity LogSeverityThreshold `json:"severity"`
	// defines list of categories
	Categories []LogSinkCategory `json:"categories"`
	// defines syslog host list
	SyslogHosts []string `json:"syslogHosts"`
	// defines list of Log sink filters
	Filters []*LogSinkFilterInput `json:"filters,omitempty"`
	// defines list of log Sink properties
	Properties []*EntityPropertyInput `json:"properties,omitempty"`
	// Ignored at creation time but can be used to compare bundle with gw state
	Checksum string `json:"checksum"`
}

// GetGoid returns LogSinkInput.Goid, and is useful for accessing the field via an interface.
func (v *LogSinkInput) GetGoid() string { return v.Goid }

// GetName returns LogSinkInput.Name, and is useful for accessing the field via an interface.
func (v *LogSinkInput) GetName() string { return v.Name }

// GetDescription returns LogSinkInput.Description, and is useful for accessing the field via an interface.
func (v *LogSinkInput) GetDescription() string { return v.Description }

// GetType returns LogSinkInput.Type, and is useful for accessing the field via an interface.
func (v *LogSinkInput) GetType() LogSinkType { return v.Type }

// GetEnabled returns LogSinkInput.Enabled, and is useful for accessing the field via an interface.
func (v *LogSinkInput) GetEnabled() bool { return v.Enabled }

// GetSeverity returns LogSinkInput.Severity, and is useful for accessing the field via an interface.
func (v *LogSinkInput) GetSeverity() LogSeverityThreshold { return v.Severity }

// GetCategories returns LogSinkInput.Categories, and is useful for accessing the field via an interface.
func (v *LogSinkInput) GetCategories() []LogSinkCategory { return v.Categories }

// GetSyslogHosts returns LogSinkInput.SyslogHosts, and is useful for accessing the field via an interface.
func (v *LogSinkInput) GetSyslogHosts() []string { return v.SyslogHosts }

// GetFilters returns LogSinkInput.Filters, and is useful for accessing the field via an interface.
func (v *LogSinkInput) GetFilters() []*LogSinkFilterInput { return v.Filters }

// GetProperties returns LogSinkInput.Properties, and is useful for accessing the field via an interface.
func (v *LogSinkInput) GetProperties() []*EntityPropertyInput { return v.Properties }

// GetChecksum returns LogSinkInput.Checksum, and is useful for accessing the field via an interface.
func (v *LogSinkInput) GetChecksum() string { return v.Checksum }

// Indicates the type of sink . File Based Or SYSLOG based
type LogSinkType string

const (
	LogSinkTypeFile   LogSinkType = "FILE"
	LogSinkTypeSyslog LogSinkType = "SYSLOG"
)

type MemberStrategyInput struct {
	// Possible values are 0 for MEMBERS_ARE_DN, 1 MEMBERS_ARE_LOGIN, 2 MEMBERS_ARE_NVPAIR, 3 MEMBERS_BY_OU
	Val int `json:"val"`
}

// GetVal returns MemberStrategyInput.Val, and is useful for accessing the field via an interface.
func (v *MemberStrategyInput) GetVal() int { return v.Val }

type MembershipInput struct {
	// The name of group to which the membership is defined
	Name         string `json:"name"`
	Goid         string `json:"goid"`
	Description  string `json:"description"`
	ProviderName string `json:"providerName"`
	Checksum     string `json:"checksum"`
}

// GetName returns MembershipInput.Name, and is useful for accessing the field via an interface.
func (v *MembershipInput) GetName() string { return v.Name }

// GetGoid returns MembershipInput.Goid, and is useful for accessing the field via an interface.
func (v *MembershipInput) GetGoid() string { return v.Goid }

// GetDescription returns MembershipInput.Description, and is useful for accessing the field via an interface.
func (v *MembershipInput) GetDescription() string { return v.Description }

// GetProviderName returns MembershipInput.ProviderName, and is useful for accessing the field via an interface.
func (v *MembershipInput) GetProviderName() string { return v.ProviderName }

// GetChecksum returns MembershipInput.Checksum, and is useful for accessing the field via an interface.
func (v *MembershipInput) GetChecksum() string { return v.Checksum }

type ModuleType string

const (
	ModuleTypeModularAssertion ModuleType = "MODULAR_ASSERTION"
	ModuleTypeCustomAssertion  ModuleType = "CUSTOM_ASSERTION"
)

type OcspNonceUsage string

const (
	// To include nonce in OCSP requests
	OcspNonceUsageIncludeNonce OcspNonceUsage = "INCLUDE_NONCE"
	// Do not include nonce in OCSP requests
	OcspNonceUsageExcludeNonce OcspNonceUsage = "EXCLUDE_NONCE"
	// Let pkix.ocsp.useNonce cluster wide property decide
	OcspNonceUsageUseNonceConditionally OcspNonceUsage = "USE_NONCE_CONDITIONALLY"
)

type PasswdStrategyInput struct {
	// Possible values are 0 for CLEAR, 1 for HASHED
	Val int `json:"val"`
}

// GetVal returns PasswdStrategyInput.Val, and is useful for accessing the field via an interface.
func (v *PasswdStrategyInput) GetVal() int { return v.Val }

type PasswordPolicyInput struct {
	// The internal entity unique identifier
	Goid string `json:"goid"`
	// Ignored at creation time but can be used to compare bundle with gw state
	Checksum string `json:"checksum"`
	// Force password change for new user and reset
	ForcePasswordChangeNewUser bool `json:"forcePasswordChangeNewUser"`
	// To enable/disable no repeating characters
	NoRepeatingCharacters bool `json:"noRepeatingCharacters"`
	// Minimum Password Length - Enter the minimum number of characters ranging from 3 to 128 required for the password.
	MinPasswordLength int `json:"minPasswordLength"`
	// Maximum Password Length - Enter the maximum number of characters ranging from 3 to 128 required for the password.
	MaxPasswordLength int `json:"maxPasswordLength"`
	// Set the number of uppercase letters that are required for the password. ranging from 1 to 128
	UpperMinimum int `json:"upperMinimum"`
	// Set the number of lowercase letters that are required for the password. ranging from 1 to 128
	LowerMinimum int `json:"lowerMinimum"`
	// Sets how many numbers (0-9) are required for the password. ranging from 1 to 128
	NumberMinimum int `json:"numberMinimum"`
	// Sets how many symbol characters are required for the password. ranging from 1 to 128
	SymbolMinimum int `json:"symbolMinimum"`
	// Sets how many non numeric characters are required for the password. ranging from 1 to 128
	NonNumericMinimum int `json:"nonNumericMinimum"`
	// Sets how many characters are required for the password. ranging from 1 to 128
	CharDiffMinimum int `json:"charDiffMinimum"`
	// Enter the number of times, between 1 and 50, that a new password must be different from the current password
	RepeatFrequency int `json:"repeatFrequency"`
	// Days required for the password to be expired. ranging from 1 to 1825
	PasswordExpiry int `json:"passwordExpiry"`
	// Allow One Password Change Per 24 Hours
	AllowableChangesPerDay bool `json:"allowableChangesPerDay"`
}

// GetGoid returns PasswordPolicyInput.Goid, and is useful for accessing the field via an interface.
func (v *PasswordPolicyInput) GetGoid() string { return v.Goid }

// GetChecksum returns PasswordPolicyInput.Checksum, and is useful for accessing the field via an interface.
func (v *PasswordPolicyInput) GetChecksum() string { return v.Checksum }

// GetForcePasswordChangeNewUser returns PasswordPolicyInput.ForcePasswordChangeNewUser, and is useful for accessing the field via an interface.
func (v *PasswordPolicyInput) GetForcePasswordChangeNewUser() bool {
	return v.ForcePasswordChangeNewUser
}

// GetNoRepeatingCharacters returns PasswordPolicyInput.NoRepeatingCharacters, and is useful for accessing the field via an interface.
func (v *PasswordPolicyInput) GetNoRepeatingCharacters() bool { return v.NoRepeatingCharacters }

// GetMinPasswordLength returns PasswordPolicyInput.MinPasswordLength, and is useful for accessing the field via an interface.
func (v *PasswordPolicyInput) GetMinPasswordLength() int { return v.MinPasswordLength }

// GetMaxPasswordLength returns PasswordPolicyInput.MaxPasswordLength, and is useful for accessing the field via an interface.
func (v *PasswordPolicyInput) GetMaxPasswordLength() int { return v.MaxPasswordLength }

// GetUpperMinimum returns PasswordPolicyInput.UpperMinimum, and is useful for accessing the field via an interface.
func (v *PasswordPolicyInput) GetUpperMinimum() int { return v.UpperMinimum }

// GetLowerMinimum returns PasswordPolicyInput.LowerMinimum, and is useful for accessing the field via an interface.
func (v *PasswordPolicyInput) GetLowerMinimum() int { return v.LowerMinimum }

// GetNumberMinimum returns PasswordPolicyInput.NumberMinimum, and is useful for accessing the field via an interface.
func (v *PasswordPolicyInput) GetNumberMinimum() int { return v.NumberMinimum }

// GetSymbolMinimum returns PasswordPolicyInput.SymbolMinimum, and is useful for accessing the field via an interface.
func (v *PasswordPolicyInput) GetSymbolMinimum() int { return v.SymbolMinimum }

// GetNonNumericMinimum returns PasswordPolicyInput.NonNumericMinimum, and is useful for accessing the field via an interface.
func (v *PasswordPolicyInput) GetNonNumericMinimum() int { return v.NonNumericMinimum }

// GetCharDiffMinimum returns PasswordPolicyInput.CharDiffMinimum, and is useful for accessing the field via an interface.
func (v *PasswordPolicyInput) GetCharDiffMinimum() int { return v.CharDiffMinimum }

// GetRepeatFrequency returns PasswordPolicyInput.RepeatFrequency, and is useful for accessing the field via an interface.
func (v *PasswordPolicyInput) GetRepeatFrequency() int { return v.RepeatFrequency }

// GetPasswordExpiry returns PasswordPolicyInput.PasswordExpiry, and is useful for accessing the field via an interface.
func (v *PasswordPolicyInput) GetPasswordExpiry() int { return v.PasswordExpiry }

// GetAllowableChangesPerDay returns PasswordPolicyInput.AllowableChangesPerDay, and is useful for accessing the field via an interface.
func (v *PasswordPolicyInput) GetAllowableChangesPerDay() bool { return v.AllowableChangesPerDay }

type PolicyBackedIdpInput struct {
	// The internal entity unique identifier
	Goid string `json:"goid"`
	// Name of the simple ldap identity provider
	Name string `json:"name"`
	// A checksum of the properties
	Checksum string `json:"checksum"`
	// Authentication Policy Name
	AuthPolicyName string `json:"authPolicyName"`
	// Default Role
	DefaultRoleName string `json:"defaultRoleName"`
	// Additional properties
	Properties []*EntityPropertyInput `json:"properties,omitempty"`
}

// GetGoid returns PolicyBackedIdpInput.Goid, and is useful for accessing the field via an interface.
func (v *PolicyBackedIdpInput) GetGoid() string { return v.Goid }

// GetName returns PolicyBackedIdpInput.Name, and is useful for accessing the field via an interface.
func (v *PolicyBackedIdpInput) GetName() string { return v.Name }

// GetChecksum returns PolicyBackedIdpInput.Checksum, and is useful for accessing the field via an interface.
func (v *PolicyBackedIdpInput) GetChecksum() string { return v.Checksum }

// GetAuthPolicyName returns PolicyBackedIdpInput.AuthPolicyName, and is useful for accessing the field via an interface.
func (v *PolicyBackedIdpInput) GetAuthPolicyName() string { return v.AuthPolicyName }

// GetDefaultRoleName returns PolicyBackedIdpInput.DefaultRoleName, and is useful for accessing the field via an interface.
func (v *PolicyBackedIdpInput) GetDefaultRoleName() string { return v.DefaultRoleName }

// GetProperties returns PolicyBackedIdpInput.Properties, and is useful for accessing the field via an interface.
func (v *PolicyBackedIdpInput) GetProperties() []*EntityPropertyInput { return v.Properties }

type PolicyFragmentInput struct {
	// The internal entity unique identifier
	Goid string `json:"goid"`
	// The folder path where to create this policy.  If the path does not exist, it will be created
	FolderPath string `json:"folderPath"`
	// The name of the policy. Policies are unique by name.
	Name string `json:"name"`
	// The guid for this policy, if none provided, assigned at creation
	Guid string `json:"guid"`
	// The policy
	Policy *PolicyInput `json:"policy,omitempty"`
	Soap   bool         `json:"soap"`
	// Ignored at creation time but can be used to compare bundle with gw state
	Checksum string `json:"checksum"`
}

// GetGoid returns PolicyFragmentInput.Goid, and is useful for accessing the field via an interface.
func (v *PolicyFragmentInput) GetGoid() string { return v.Goid }

// GetFolderPath returns PolicyFragmentInput.FolderPath, and is useful for accessing the field via an interface.
func (v *PolicyFragmentInput) GetFolderPath() string { return v.FolderPath }

// GetName returns PolicyFragmentInput.Name, and is useful for accessing the field via an interface.
func (v *PolicyFragmentInput) GetName() string { return v.Name }

// GetGuid returns PolicyFragmentInput.Guid, and is useful for accessing the field via an interface.
func (v *PolicyFragmentInput) GetGuid() string { return v.Guid }

// GetPolicy returns PolicyFragmentInput.Policy, and is useful for accessing the field via an interface.
func (v *PolicyFragmentInput) GetPolicy() *PolicyInput { return v.Policy }

// GetSoap returns PolicyFragmentInput.Soap, and is useful for accessing the field via an interface.
func (v *PolicyFragmentInput) GetSoap() bool { return v.Soap }

// GetChecksum returns PolicyFragmentInput.Checksum, and is useful for accessing the field via an interface.
func (v *PolicyFragmentInput) GetChecksum() string { return v.Checksum }

type PolicyInput struct {
	// The policy xml
	Xml string `json:"xml,omitempty"`
	// The policy JSON
	Json string `json:"json,omitempty"`
	// The policy YAML
	Yaml string `json:"yaml,omitempty"`
	// The policy code
	Code interface{} `json:"code,omitempty"`
}

// GetXml returns PolicyInput.Xml, and is useful for accessing the field via an interface.
func (v *PolicyInput) GetXml() string { return v.Xml }

// GetJson returns PolicyInput.Json, and is useful for accessing the field via an interface.
func (v *PolicyInput) GetJson() string { return v.Json }

// GetYaml returns PolicyInput.Yaml, and is useful for accessing the field via an interface.
func (v *PolicyInput) GetYaml() string { return v.Yaml }

// GetCode returns PolicyInput.Code, and is useful for accessing the field via an interface.
func (v *PolicyInput) GetCode() interface{} { return v.Code }

type PolicyRevisionInput struct {
	Goid    string    `json:"goid"`
	Ordinal int64     `json:"ordinal"`
	Active  bool      `json:"active"`
	Comment string    `json:"comment"`
	Author  string    `json:"author"`
	Time    time.Time `json:"time"`
	// The policy XML
	Xml string `json:"xml,omitempty"`
	// The policy JSON
	Json string `json:"json,omitempty"`
	// The policy YAML
	Yaml string `json:"yaml,omitempty"`
	// The policy code
	Code interface{} `json:"code,omitempty"`
}

// GetGoid returns PolicyRevisionInput.Goid, and is useful for accessing the field via an interface.
func (v *PolicyRevisionInput) GetGoid() string { return v.Goid }

// GetOrdinal returns PolicyRevisionInput.Ordinal, and is useful for accessing the field via an interface.
func (v *PolicyRevisionInput) GetOrdinal() int64 { return v.Ordinal }

// GetActive returns PolicyRevisionInput.Active, and is useful for accessing the field via an interface.
func (v *PolicyRevisionInput) GetActive() bool { return v.Active }

// GetComment returns PolicyRevisionInput.Comment, and is useful for accessing the field via an interface.
func (v *PolicyRevisionInput) GetComment() string { return v.Comment }

// GetAuthor returns PolicyRevisionInput.Author, and is useful for accessing the field via an interface.
func (v *PolicyRevisionInput) GetAuthor() string { return v.Author }

// GetTime returns PolicyRevisionInput.Time, and is useful for accessing the field via an interface.
func (v *PolicyRevisionInput) GetTime() time.Time { return v.Time }

// GetXml returns PolicyRevisionInput.Xml, and is useful for accessing the field via an interface.
func (v *PolicyRevisionInput) GetXml() string { return v.Xml }

// GetJson returns PolicyRevisionInput.Json, and is useful for accessing the field via an interface.
func (v *PolicyRevisionInput) GetJson() string { return v.Json }

// GetYaml returns PolicyRevisionInput.Yaml, and is useful for accessing the field via an interface.
func (v *PolicyRevisionInput) GetYaml() string { return v.Yaml }

// GetCode returns PolicyRevisionInput.Code, and is useful for accessing the field via an interface.
func (v *PolicyRevisionInput) GetCode() interface{} { return v.Code }

type PolicyUsageType string

const (
	// Do not perform revocation check
	PolicyUsageTypeNone PolicyUsageType = "NONE"
	// Use the default revocation check policy
	PolicyUsageTypeUseDefault PolicyUsageType = "USE_DEFAULT"
	// Use the specified revocation check policy
	PolicyUsageTypeSpecified PolicyUsageType = "SPECIFIED"
)

type RevocationCheckPolicyInput struct {
	// The goid for this revocation check policy
	Goid string `json:"goid"`
	// Name that describes the revocation checking policy
	Name string `json:"name"`
	// Ignored at creation time but can be used to compare bundle with gw state
	Checksum string `json:"checksum"`
	// Use as default revocation check policy
	DefaultPolicy bool `json:"defaultPolicy"`
	// Succeed if revocation status is unknown
	DefaultSuccess bool `json:"defaultSuccess"`
	// Continue processing if server is unavailable
	ContinueOnServerUnavailable bool `json:"continueOnServerUnavailable"`
	// Certificate revocation check properties
	RevocationCheckPolicyItems []*RevocationCheckPolicyItemInput `json:"revocationCheckPolicyItems,omitempty"`
}

// GetGoid returns RevocationCheckPolicyInput.Goid, and is useful for accessing the field via an interface.
func (v *RevocationCheckPolicyInput) GetGoid() string { return v.Goid }

// GetName returns RevocationCheckPolicyInput.Name, and is useful for accessing the field via an interface.
func (v *RevocationCheckPolicyInput) GetName() string { return v.Name }

// GetChecksum returns RevocationCheckPolicyInput.Checksum, and is useful for accessing the field via an interface.
func (v *RevocationCheckPolicyInput) GetChecksum() string { return v.Checksum }

// GetDefaultPolicy returns RevocationCheckPolicyInput.DefaultPolicy, and is useful for accessing the field via an interface.
func (v *RevocationCheckPolicyInput) GetDefaultPolicy() bool { return v.DefaultPolicy }

// GetDefaultSuccess returns RevocationCheckPolicyInput.DefaultSuccess, and is useful for accessing the field via an interface.
func (v *RevocationCheckPolicyInput) GetDefaultSuccess() bool { return v.DefaultSuccess }

// GetContinueOnServerUnavailable returns RevocationCheckPolicyInput.ContinueOnServerUnavailable, and is useful for accessing the field via an interface.
func (v *RevocationCheckPolicyInput) GetContinueOnServerUnavailable() bool {
	return v.ContinueOnServerUnavailable
}

// GetRevocationCheckPolicyItems returns RevocationCheckPolicyInput.RevocationCheckPolicyItems, and is useful for accessing the field via an interface.
func (v *RevocationCheckPolicyInput) GetRevocationCheckPolicyItems() []*RevocationCheckPolicyItemInput {
	return v.RevocationCheckPolicyItems
}

type RevocationCheckPolicyItemInput struct {
	// Type for Checking OCSP or CRL
	Type CertRevocationCheckPropertyType `json:"type"`
	// If the CRL from URL or OCSP from URL option was selected, enter the URL Otherwise provide regex.
	// CRL_FROM_CERTIFICATE &  OCSP_FROM_CERTIFICATE options uses URL Regex &
	// CRL_FROM_URL & OCSP_FROM_URL options uses URLs.
	// This is caller's responsibility to provide valid URL or Regex, Graphman won't validate it.
	Url string `json:"url"`
	// If user permitting the entity that issued the certificate
	AllowIssuerSignature bool `json:"allowIssuerSignature"`
	// Whether to include a nonce in OCSP request, default is to set INCLUDE_NONCE
	NonceUsage OcspNonceUsage `json:"nonceUsage"`
	// The sha1 thumbprint of the certificate
	SignerThumbprintSha1s []string `json:"signerThumbprintSha1s"`
}

// GetType returns RevocationCheckPolicyItemInput.Type, and is useful for accessing the field via an interface.
func (v *RevocationCheckPolicyItemInput) GetType() CertRevocationCheckPropertyType { return v.Type }

// GetUrl returns RevocationCheckPolicyItemInput.Url, and is useful for accessing the field via an interface.
func (v *RevocationCheckPolicyItemInput) GetUrl() string { return v.Url }

// GetAllowIssuerSignature returns RevocationCheckPolicyItemInput.AllowIssuerSignature, and is useful for accessing the field via an interface.
func (v *RevocationCheckPolicyItemInput) GetAllowIssuerSignature() bool {
	return v.AllowIssuerSignature
}

// GetNonceUsage returns RevocationCheckPolicyItemInput.NonceUsage, and is useful for accessing the field via an interface.
func (v *RevocationCheckPolicyItemInput) GetNonceUsage() OcspNonceUsage { return v.NonceUsage }

// GetSignerThumbprintSha1s returns RevocationCheckPolicyItemInput.SignerThumbprintSha1s, and is useful for accessing the field via an interface.
func (v *RevocationCheckPolicyItemInput) GetSignerThumbprintSha1s() []string {
	return v.SignerThumbprintSha1s
}

// Role configuration
type RoleInput struct {
	// The goid for the Role
	Goid string `json:"goid"`
	// Name of a Role
	Name string `json:"name"`
	// The configuration checksum
	Checksum string `json:"checksum"`
	// Type of a role
	RoleType RoleType `json:"roleType"`
	// Description of the role. This is optional
	Description string `json:"description"`
	// Tag: Either Admin or Null
	Tag Tag `json:"tag,omitempty"`
	// Whether to replace the existing assignees with the specified users/groups
	ReplaceAssignees bool `json:"replaceAssignees"`
	// One or more users assigned to the role
	UserAssignees []*UserRefInput `json:"userAssignees,omitempty"`
	// One or more groups assigned to the role
	GroupAssignees []*GroupRefInput `json:"groupAssignees,omitempty"`
}

// GetGoid returns RoleInput.Goid, and is useful for accessing the field via an interface.
func (v *RoleInput) GetGoid() string { return v.Goid }

// GetName returns RoleInput.Name, and is useful for accessing the field via an interface.
func (v *RoleInput) GetName() string { return v.Name }

// GetChecksum returns RoleInput.Checksum, and is useful for accessing the field via an interface.
func (v *RoleInput) GetChecksum() string { return v.Checksum }

// GetRoleType returns RoleInput.RoleType, and is useful for accessing the field via an interface.
func (v *RoleInput) GetRoleType() RoleType { return v.RoleType }

// GetDescription returns RoleInput.Description, and is useful for accessing the field via an interface.
func (v *RoleInput) GetDescription() string { return v.Description }

// GetTag returns RoleInput.Tag, and is useful for accessing the field via an interface.
func (v *RoleInput) GetTag() Tag { return v.Tag }

// GetReplaceAssignees returns RoleInput.ReplaceAssignees, and is useful for accessing the field via an interface.
func (v *RoleInput) GetReplaceAssignees() bool { return v.ReplaceAssignees }

// GetUserAssignees returns RoleInput.UserAssignees, and is useful for accessing the field via an interface.
func (v *RoleInput) GetUserAssignees() []*UserRefInput { return v.UserAssignees }

// GetGroupAssignees returns RoleInput.GroupAssignees, and is useful for accessing the field via an interface.
func (v *RoleInput) GetGroupAssignees() []*GroupRefInput { return v.GroupAssignees }

type RoleType string

const (
	RoleTypeSystem RoleType = "SYSTEM"
	RoleTypeCustom RoleType = "CUSTOM"
)

type SMConfigInput struct {
	// The goid for the CA SSO connection
	Goid string `json:"goid"`
	// Name of the CA SSO configuration
	Name string `json:"name"`
	// Indicates whether the specified configuration is currently enabled or disabled
	Enabled bool `json:"enabled"`
	// Name of the host registered with the CA SSO Policy Server
	AgentHost string `json:"agentHost"`
	// The IP address of the CA SSO agent. This field is required if the Check IP check box is selected
	AgentIP string `json:"agentIP"`
	// CA SSO Policy Server host configuration used by the agent
	AgentHostConfig string `json:"agentHostConfig"`
	// CA SSO shared secret used by the agent to establish communication with the Policy Server
	AgentSecret string `json:"agentSecret"`
	// Choose the FIPS mode supported by the CA SSO Policy Server. The available values are: COMPAT(default)/MIGRATE/ONLY
	CryptoMode SMCryptoMode `json:"cryptoMode"`
	// The CA SSO Policy Server compare the client IP against the address stored in the SSO Token
	IpCheckEnabled bool `json:"ipCheckEnabled"`
	// Whether to update the SSO Token after successful authentication/authorization
	UpdateSSOToken bool `json:"updateSSOToken"`
	// The percentage of servers within a cluster that must be available for Policy Server requests
	ClusterFailoverThreshold int  `json:"clusterFailoverThreshold"`
	NonClusterFailover       bool `json:"nonClusterFailover"`
	// User name of the CA SSO administrator
	Username string `json:"username"`
	// The secure password reference
	SecurePasswordName string `json:"securePasswordName"`
	// The Siteminder configuration properties
	Properties []*EntityPropertyInput `json:"properties,omitempty"`
	// The configuration checksum
	Checksum string `json:"checksum"`
}

// GetGoid returns SMConfigInput.Goid, and is useful for accessing the field via an interface.
func (v *SMConfigInput) GetGoid() string { return v.Goid }

// GetName returns SMConfigInput.Name, and is useful for accessing the field via an interface.
func (v *SMConfigInput) GetName() string { return v.Name }

// GetEnabled returns SMConfigInput.Enabled, and is useful for accessing the field via an interface.
func (v *SMConfigInput) GetEnabled() bool { return v.Enabled }

// GetAgentHost returns SMConfigInput.AgentHost, and is useful for accessing the field via an interface.
func (v *SMConfigInput) GetAgentHost() string { return v.AgentHost }

// GetAgentIP returns SMConfigInput.AgentIP, and is useful for accessing the field via an interface.
func (v *SMConfigInput) GetAgentIP() string { return v.AgentIP }

// GetAgentHostConfig returns SMConfigInput.AgentHostConfig, and is useful for accessing the field via an interface.
func (v *SMConfigInput) GetAgentHostConfig() string { return v.AgentHostConfig }

// GetAgentSecret returns SMConfigInput.AgentSecret, and is useful for accessing the field via an interface.
func (v *SMConfigInput) GetAgentSecret() string { return v.AgentSecret }

// GetCryptoMode returns SMConfigInput.CryptoMode, and is useful for accessing the field via an interface.
func (v *SMConfigInput) GetCryptoMode() SMCryptoMode { return v.CryptoMode }

// GetIpCheckEnabled returns SMConfigInput.IpCheckEnabled, and is useful for accessing the field via an interface.
func (v *SMConfigInput) GetIpCheckEnabled() bool { return v.IpCheckEnabled }

// GetUpdateSSOToken returns SMConfigInput.UpdateSSOToken, and is useful for accessing the field via an interface.
func (v *SMConfigInput) GetUpdateSSOToken() bool { return v.UpdateSSOToken }

// GetClusterFailoverThreshold returns SMConfigInput.ClusterFailoverThreshold, and is useful for accessing the field via an interface.
func (v *SMConfigInput) GetClusterFailoverThreshold() int { return v.ClusterFailoverThreshold }

// GetNonClusterFailover returns SMConfigInput.NonClusterFailover, and is useful for accessing the field via an interface.
func (v *SMConfigInput) GetNonClusterFailover() bool { return v.NonClusterFailover }

// GetUsername returns SMConfigInput.Username, and is useful for accessing the field via an interface.
func (v *SMConfigInput) GetUsername() string { return v.Username }

// GetSecurePasswordName returns SMConfigInput.SecurePasswordName, and is useful for accessing the field via an interface.
func (v *SMConfigInput) GetSecurePasswordName() string { return v.SecurePasswordName }

// GetProperties returns SMConfigInput.Properties, and is useful for accessing the field via an interface.
func (v *SMConfigInput) GetProperties() []*EntityPropertyInput { return v.Properties }

// GetChecksum returns SMConfigInput.Checksum, and is useful for accessing the field via an interface.
func (v *SMConfigInput) GetChecksum() string { return v.Checksum }

type SMCryptoMode string

const (
	SMCryptoModeCompat  SMCryptoMode = "COMPAT"
	SMCryptoModeMigrate SMCryptoMode = "MIGRATE"
	SMCryptoModeFips    SMCryptoMode = "FIPS"
)

type ScheduledTaskInput struct {
	// The internal entity unique identifier
	Goid string `json:"goid"`
	// The name of the scheduled task
	Name string `json:"name"`
	// The name of the policy for scheduled task
	PolicyName string  `json:"policyName"`
	JobType    JobType `json:"jobType"`
	// The cron job expression
	CronExpression string `json:"cronExpression"`
	// Whether to execute on single node
	ExecuteOnSingleNode bool `json:"executeOnSingleNode"`
	// Whether to execute the RECURRING task now?
	ExecuteOnCreation bool `json:"executeOnCreation"`
	// Specify a future execution date for a ONE_TIME task
	ExecutionDate time.Time `json:"executionDate"`
	// The scheduled task status
	Status                JobStatus `json:"status"`
	RunAsUser             string    `json:"runAsUser"`
	RunAsUserProviderName string    `json:"runAsUserProviderName"`
	// The configuration checksum
	Checksum string `json:"checksum"`
}

// GetGoid returns ScheduledTaskInput.Goid, and is useful for accessing the field via an interface.
func (v *ScheduledTaskInput) GetGoid() string { return v.Goid }

// GetName returns ScheduledTaskInput.Name, and is useful for accessing the field via an interface.
func (v *ScheduledTaskInput) GetName() string { return v.Name }

// GetPolicyName returns ScheduledTaskInput.PolicyName, and is useful for accessing the field via an interface.
func (v *ScheduledTaskInput) GetPolicyName() string { return v.PolicyName }

// GetJobType returns ScheduledTaskInput.JobType, and is useful for accessing the field via an interface.
func (v *ScheduledTaskInput) GetJobType() JobType { return v.JobType }

// GetCronExpression returns ScheduledTaskInput.CronExpression, and is useful for accessing the field via an interface.
func (v *ScheduledTaskInput) GetCronExpression() string { return v.CronExpression }

// GetExecuteOnSingleNode returns ScheduledTaskInput.ExecuteOnSingleNode, and is useful for accessing the field via an interface.
func (v *ScheduledTaskInput) GetExecuteOnSingleNode() bool { return v.ExecuteOnSingleNode }

// GetExecuteOnCreation returns ScheduledTaskInput.ExecuteOnCreation, and is useful for accessing the field via an interface.
func (v *ScheduledTaskInput) GetExecuteOnCreation() bool { return v.ExecuteOnCreation }

// GetExecutionDate returns ScheduledTaskInput.ExecutionDate, and is useful for accessing the field via an interface.
func (v *ScheduledTaskInput) GetExecutionDate() time.Time { return v.ExecutionDate }

// GetStatus returns ScheduledTaskInput.Status, and is useful for accessing the field via an interface.
func (v *ScheduledTaskInput) GetStatus() JobStatus { return v.Status }

// GetRunAsUser returns ScheduledTaskInput.RunAsUser, and is useful for accessing the field via an interface.
func (v *ScheduledTaskInput) GetRunAsUser() string { return v.RunAsUser }

// GetRunAsUserProviderName returns ScheduledTaskInput.RunAsUserProviderName, and is useful for accessing the field via an interface.
func (v *ScheduledTaskInput) GetRunAsUserProviderName() string { return v.RunAsUserProviderName }

// GetChecksum returns ScheduledTaskInput.Checksum, and is useful for accessing the field via an interface.
func (v *ScheduledTaskInput) GetChecksum() string { return v.Checksum }

type SchemaInput struct {
	// The internal entity unique identifier
	Goid string `json:"goid"`
	// A reference to the schema. This id is what is referred to in policy and is often mirror of the target namespace
	SystemId string `json:"systemId"`
	// The target namespace in the XML schema
	TargetNs string `json:"targetNs"`
	// An optional description for the schema
	Description string `json:"description"`
	// The content of XML schema
	Content string `json:"content"`
	// Ignored at creation time but can be used to compare bundle with gw state
	Checksum string `json:"checksum"`
}

// GetGoid returns SchemaInput.Goid, and is useful for accessing the field via an interface.
func (v *SchemaInput) GetGoid() string { return v.Goid }

// GetSystemId returns SchemaInput.SystemId, and is useful for accessing the field via an interface.
func (v *SchemaInput) GetSystemId() string { return v.SystemId }

// GetTargetNs returns SchemaInput.TargetNs, and is useful for accessing the field via an interface.
func (v *SchemaInput) GetTargetNs() string { return v.TargetNs }

// GetDescription returns SchemaInput.Description, and is useful for accessing the field via an interface.
func (v *SchemaInput) GetDescription() string { return v.Description }

// GetContent returns SchemaInput.Content, and is useful for accessing the field via an interface.
func (v *SchemaInput) GetContent() string { return v.Content }

// GetChecksum returns SchemaInput.Checksum, and is useful for accessing the field via an interface.
func (v *SchemaInput) GetChecksum() string { return v.Checksum }

type SecretInput struct {
	// Identify the password being stored. You may use letters, numbers, dashes, and underscores.
	// Names that contain spaces or periods are valid, but the resulting stored
	// password cannot be referenced via context variable.
	// Names that contain @ or $ are valid, but the resulting stored password cannot be referenced via context variable.
	Name string `json:"name"`
	// Password or PEM Private Key
	SecretType SecretType `json:"secretType"`
	// The goid for the Secret
	Goid string `json:"goid"`
	// Ignored at entity creation time but declared here so you can embed checksums in graphman bundles
	Checksum string `json:"checksum"`
	// Whether this secret can be referred to in policy via context variable ${secpass...
	VariableReferencable bool `json:"variableReferencable"`
	// Base64 encrypted secret. The encryption is compatible with openssl secret encryption
	// using cypher AES/CBC/PKCS5Padding. You can create this value at command line:
	// > echo -n "<clear text secret>" | openssl enc -aes-256-cbc -md sha256 -pass pass:<password> -a
	Secret string `json:"secret"`
	// Description of the password. This is optional
	Description string `json:"description"`
}

// GetName returns SecretInput.Name, and is useful for accessing the field via an interface.
func (v *SecretInput) GetName() string { return v.Name }

// GetSecretType returns SecretInput.SecretType, and is useful for accessing the field via an interface.
func (v *SecretInput) GetSecretType() SecretType { return v.SecretType }

// GetGoid returns SecretInput.Goid, and is useful for accessing the field via an interface.
func (v *SecretInput) GetGoid() string { return v.Goid }

// GetChecksum returns SecretInput.Checksum, and is useful for accessing the field via an interface.
func (v *SecretInput) GetChecksum() string { return v.Checksum }

// GetVariableReferencable returns SecretInput.VariableReferencable, and is useful for accessing the field via an interface.
func (v *SecretInput) GetVariableReferencable() bool { return v.VariableReferencable }

// GetSecret returns SecretInput.Secret, and is useful for accessing the field via an interface.
func (v *SecretInput) GetSecret() string { return v.Secret }

// GetDescription returns SecretInput.Description, and is useful for accessing the field via an interface.
func (v *SecretInput) GetDescription() string { return v.Description }

type SecretType string

const (
	// Stored password for example used in the jdbc connection
	SecretTypePassword SecretType = "PASSWORD"
	// Secure pem key for example used in the route via ssh assertion
	SecretTypePemPrivateKey SecretType = "PEM_PRIVATE_KEY"
)

type ServerModuleFileInput struct {
	// The internal entity unique identifier
	Goid string `json:"goid"`
	// The Server module name
	Name string `json:"name"`
	// The Server module type
	ModuleType ModuleType `json:"moduleType"`
	// The Server module SHA256 digest value
	ModuleSha256 string `json:"moduleSha256"`
	// The Server module signature
	Signature string `json:"signature"`
	// The base64 encoded signer certificate
	SignerCertBase64 string `json:"signerCertBase64"`
	// The Server module file properties
	Properties []*EntityPropertyInput `json:"properties,omitempty"`
	// Ignored at creation time but can be used to compare bundle with gw state
	Checksum string `json:"checksum"`
}

// GetGoid returns ServerModuleFileInput.Goid, and is useful for accessing the field via an interface.
func (v *ServerModuleFileInput) GetGoid() string { return v.Goid }

// GetName returns ServerModuleFileInput.Name, and is useful for accessing the field via an interface.
func (v *ServerModuleFileInput) GetName() string { return v.Name }

// GetModuleType returns ServerModuleFileInput.ModuleType, and is useful for accessing the field via an interface.
func (v *ServerModuleFileInput) GetModuleType() ModuleType { return v.ModuleType }

// GetModuleSha256 returns ServerModuleFileInput.ModuleSha256, and is useful for accessing the field via an interface.
func (v *ServerModuleFileInput) GetModuleSha256() string { return v.ModuleSha256 }

// GetSignature returns ServerModuleFileInput.Signature, and is useful for accessing the field via an interface.
func (v *ServerModuleFileInput) GetSignature() string { return v.Signature }

// GetSignerCertBase64 returns ServerModuleFileInput.SignerCertBase64, and is useful for accessing the field via an interface.
func (v *ServerModuleFileInput) GetSignerCertBase64() string { return v.SignerCertBase64 }

// GetProperties returns ServerModuleFileInput.Properties, and is useful for accessing the field via an interface.
func (v *ServerModuleFileInput) GetProperties() []*EntityPropertyInput { return v.Properties }

// GetChecksum returns ServerModuleFileInput.Checksum, and is useful for accessing the field via an interface.
func (v *ServerModuleFileInput) GetChecksum() string { return v.Checksum }

type ServiceResolutionConfigInput struct {
	// The internal entity unique identifier
	Goid string `json:"goid"`
	// Ignored at creation time but can be used to compare bundle with gw state
	Checksum string `json:"checksum"`
	// Only Services with a resolution path are accessible
	ResolutionPathRequired bool `json:"resolutionPathRequired"`
	// Resolution paths are case sensitive
	ResolutionPathCaseSensitive bool `json:"resolutionPathCaseSensitive"`
	// Allow resolution by L7-Original-URL header
	UseL7OriginalUrl bool `json:"useL7OriginalUrl"`
	// Allow resolution by Service GOID/OID in URLs
	UseServiceGoid bool `json:"useServiceGoid"`
	// Use SOAP action
	UseSoapAction bool `json:"useSoapAction"`
	// Use SOAP body child namespace
	UseSoapBodyChildNamespace bool `json:"useSoapBodyChildNamespace"`
}

// GetGoid returns ServiceResolutionConfigInput.Goid, and is useful for accessing the field via an interface.
func (v *ServiceResolutionConfigInput) GetGoid() string { return v.Goid }

// GetChecksum returns ServiceResolutionConfigInput.Checksum, and is useful for accessing the field via an interface.
func (v *ServiceResolutionConfigInput) GetChecksum() string { return v.Checksum }

// GetResolutionPathRequired returns ServiceResolutionConfigInput.ResolutionPathRequired, and is useful for accessing the field via an interface.
func (v *ServiceResolutionConfigInput) GetResolutionPathRequired() bool {
	return v.ResolutionPathRequired
}

// GetResolutionPathCaseSensitive returns ServiceResolutionConfigInput.ResolutionPathCaseSensitive, and is useful for accessing the field via an interface.
func (v *ServiceResolutionConfigInput) GetResolutionPathCaseSensitive() bool {
	return v.ResolutionPathCaseSensitive
}

// GetUseL7OriginalUrl returns ServiceResolutionConfigInput.UseL7OriginalUrl, and is useful for accessing the field via an interface.
func (v *ServiceResolutionConfigInput) GetUseL7OriginalUrl() bool { return v.UseL7OriginalUrl }

// GetUseServiceGoid returns ServiceResolutionConfigInput.UseServiceGoid, and is useful for accessing the field via an interface.
func (v *ServiceResolutionConfigInput) GetUseServiceGoid() bool { return v.UseServiceGoid }

// GetUseSoapAction returns ServiceResolutionConfigInput.UseSoapAction, and is useful for accessing the field via an interface.
func (v *ServiceResolutionConfigInput) GetUseSoapAction() bool { return v.UseSoapAction }

// GetUseSoapBodyChildNamespace returns ServiceResolutionConfigInput.UseSoapBodyChildNamespace, and is useful for accessing the field via an interface.
func (v *ServiceResolutionConfigInput) GetUseSoapBodyChildNamespace() bool {
	return v.UseSoapBodyChildNamespace
}

type ServiceResolversInput struct {
	// The soap action referred to in the wsdl
	SoapAction string `json:"soapAction"`
	// The soap actions referred to in the wsdl
	SoapActions []string `json:"soapActions"`
	// Base uri from the wsdl of the service. This is used for service resolution
	BaseUri string `json:"baseUri"`
	// The resolution path to the service.
	ResolutionPath string `json:"resolutionPath"`
}

// GetSoapAction returns ServiceResolversInput.SoapAction, and is useful for accessing the field via an interface.
func (v *ServiceResolversInput) GetSoapAction() string { return v.SoapAction }

// GetSoapActions returns ServiceResolversInput.SoapActions, and is useful for accessing the field via an interface.
func (v *ServiceResolversInput) GetSoapActions() []string { return v.SoapActions }

// GetBaseUri returns ServiceResolversInput.BaseUri, and is useful for accessing the field via an interface.
func (v *ServiceResolversInput) GetBaseUri() string { return v.BaseUri }

// GetResolutionPath returns ServiceResolversInput.ResolutionPath, and is useful for accessing the field via an interface.
func (v *ServiceResolversInput) GetResolutionPath() string { return v.ResolutionPath }

type ServiceResourceInput struct {
	Uri     string `json:"uri"`
	Content string `json:"content"`
}

// GetUri returns ServiceResourceInput.Uri, and is useful for accessing the field via an interface.
func (v *ServiceResourceInput) GetUri() string { return v.Uri }

// GetContent returns ServiceResourceInput.Content, and is useful for accessing the field via an interface.
func (v *ServiceResourceInput) GetContent() string { return v.Content }

type SimpleLdapIdpInput struct {
	// The internal entity unique identifier
	Goid string `json:"goid"`
	// Name of the simple ldap identity provider
	Name string `json:"name"`
	// A checksum of the properties
	Checksum string `json:"checksum"`
	// simple ldap server urls
	ServerUrls []string `json:"serverUrls"`
	// Whether to use client certificate authentication
	UseSslClientAuth bool `json:"useSslClientAuth"`
	// Client key Alias
	SslClientKeyAlias string `json:"sslClientKeyAlias"`
	// Bind DN prefix
	BindDnPatternPrefix string `json:"bindDnPatternPrefix"`
	// Bind DN suffix
	BindDnPatternSuffix string `json:"bindDnPatternSuffix"`
	// Simple Ldap properties
	Properties []*EntityPropertyInput `json:"properties,omitempty"`
}

// GetGoid returns SimpleLdapIdpInput.Goid, and is useful for accessing the field via an interface.
func (v *SimpleLdapIdpInput) GetGoid() string { return v.Goid }

// GetName returns SimpleLdapIdpInput.Name, and is useful for accessing the field via an interface.
func (v *SimpleLdapIdpInput) GetName() string { return v.Name }

// GetChecksum returns SimpleLdapIdpInput.Checksum, and is useful for accessing the field via an interface.
func (v *SimpleLdapIdpInput) GetChecksum() string { return v.Checksum }

// GetServerUrls returns SimpleLdapIdpInput.ServerUrls, and is useful for accessing the field via an interface.
func (v *SimpleLdapIdpInput) GetServerUrls() []string { return v.ServerUrls }

// GetUseSslClientAuth returns SimpleLdapIdpInput.UseSslClientAuth, and is useful for accessing the field via an interface.
func (v *SimpleLdapIdpInput) GetUseSslClientAuth() bool { return v.UseSslClientAuth }

// GetSslClientKeyAlias returns SimpleLdapIdpInput.SslClientKeyAlias, and is useful for accessing the field via an interface.
func (v *SimpleLdapIdpInput) GetSslClientKeyAlias() string { return v.SslClientKeyAlias }

// GetBindDnPatternPrefix returns SimpleLdapIdpInput.BindDnPatternPrefix, and is useful for accessing the field via an interface.
func (v *SimpleLdapIdpInput) GetBindDnPatternPrefix() string { return v.BindDnPatternPrefix }

// GetBindDnPatternSuffix returns SimpleLdapIdpInput.BindDnPatternSuffix, and is useful for accessing the field via an interface.
func (v *SimpleLdapIdpInput) GetBindDnPatternSuffix() string { return v.BindDnPatternSuffix }

// GetProperties returns SimpleLdapIdpInput.Properties, and is useful for accessing the field via an interface.
func (v *SimpleLdapIdpInput) GetProperties() []*EntityPropertyInput { return v.Properties }

type SoapServiceInput struct {
	// The internal entity unique identifier
	Goid string `json:"goid"`
	// The guid for this service, if none provided, assigned at creation
	Guid string `json:"guid"`
	// The folder path where to create this service.  If the path does not exist, it will be created
	FolderPath string `json:"folderPath"`
	// The name of the service
	Name string `json:"name"`
	// The WSDL of the soap service
	Wsdl string `json:"wsdl"`
	// URL for the protected service WSDL document
	WsdlUrl string `json:"wsdlUrl"`
	// One or more additional WSDL resources
	WsdlResources []*ServiceResourceInput `json:"wsdlResources,omitempty"`
	// The resolution path of the service
	ResolutionPath string `json:"resolutionPath"`
	// Soap service resolvers
	Resolvers *SoapServiceResolverInput `json:"resolvers,omitempty"`
	// The policy
	Policy *PolicyInput `json:"policy,omitempty"`
	// Whether the service is enabled (optional, default true)
	Enabled bool `json:"enabled"`
	// The http methods allowed for this service
	MethodsAllowed []HttpMethod `json:"methodsAllowed"`
	// Which SOAP version
	SoapVersion SoapVersion `json:"soapVersion"`
	// Whether or not the gateway should process incoming ws-security soap headers
	WssProcessingEnabled bool `json:"wssProcessingEnabled"`
	TracingEnabled       bool `json:"tracingEnabled"`
	// Allow requests intended for operations not supported by the WSDL
	LaxResolution bool                   `json:"laxResolution"`
	Properties    []*EntityPropertyInput `json:"properties,omitempty"`
	// Ignored at creation time but can be used to compare bundle with gw state
	Checksum string `json:"checksum"`
}

// GetGoid returns SoapServiceInput.Goid, and is useful for accessing the field via an interface.
func (v *SoapServiceInput) GetGoid() string { return v.Goid }

// GetGuid returns SoapServiceInput.Guid, and is useful for accessing the field via an interface.
func (v *SoapServiceInput) GetGuid() string { return v.Guid }

// GetFolderPath returns SoapServiceInput.FolderPath, and is useful for accessing the field via an interface.
func (v *SoapServiceInput) GetFolderPath() string { return v.FolderPath }

// GetName returns SoapServiceInput.Name, and is useful for accessing the field via an interface.
func (v *SoapServiceInput) GetName() string { return v.Name }

// GetWsdl returns SoapServiceInput.Wsdl, and is useful for accessing the field via an interface.
func (v *SoapServiceInput) GetWsdl() string { return v.Wsdl }

// GetWsdlUrl returns SoapServiceInput.WsdlUrl, and is useful for accessing the field via an interface.
func (v *SoapServiceInput) GetWsdlUrl() string { return v.WsdlUrl }

// GetWsdlResources returns SoapServiceInput.WsdlResources, and is useful for accessing the field via an interface.
func (v *SoapServiceInput) GetWsdlResources() []*ServiceResourceInput { return v.WsdlResources }

// GetResolutionPath returns SoapServiceInput.ResolutionPath, and is useful for accessing the field via an interface.
func (v *SoapServiceInput) GetResolutionPath() string { return v.ResolutionPath }

// GetResolvers returns SoapServiceInput.Resolvers, and is useful for accessing the field via an interface.
func (v *SoapServiceInput) GetResolvers() *SoapServiceResolverInput { return v.Resolvers }

// GetPolicy returns SoapServiceInput.Policy, and is useful for accessing the field via an interface.
func (v *SoapServiceInput) GetPolicy() *PolicyInput { return v.Policy }

// GetEnabled returns SoapServiceInput.Enabled, and is useful for accessing the field via an interface.
func (v *SoapServiceInput) GetEnabled() bool { return v.Enabled }

// GetMethodsAllowed returns SoapServiceInput.MethodsAllowed, and is useful for accessing the field via an interface.
func (v *SoapServiceInput) GetMethodsAllowed() []HttpMethod { return v.MethodsAllowed }

// GetSoapVersion returns SoapServiceInput.SoapVersion, and is useful for accessing the field via an interface.
func (v *SoapServiceInput) GetSoapVersion() SoapVersion { return v.SoapVersion }

// GetWssProcessingEnabled returns SoapServiceInput.WssProcessingEnabled, and is useful for accessing the field via an interface.
func (v *SoapServiceInput) GetWssProcessingEnabled() bool { return v.WssProcessingEnabled }

// GetTracingEnabled returns SoapServiceInput.TracingEnabled, and is useful for accessing the field via an interface.
func (v *SoapServiceInput) GetTracingEnabled() bool { return v.TracingEnabled }

// GetLaxResolution returns SoapServiceInput.LaxResolution, and is useful for accessing the field via an interface.
func (v *SoapServiceInput) GetLaxResolution() bool { return v.LaxResolution }

// GetProperties returns SoapServiceInput.Properties, and is useful for accessing the field via an interface.
func (v *SoapServiceInput) GetProperties() []*EntityPropertyInput { return v.Properties }

// GetChecksum returns SoapServiceInput.Checksum, and is useful for accessing the field via an interface.
func (v *SoapServiceInput) GetChecksum() string { return v.Checksum }

// Must have minimum (1 soapAction + baseUri) OR resolutionPath. You can have both too.
type SoapServiceResolverInput struct {
	// One of the SoapAction of the service to resolved. This must be specified along with a base ns from the WSDL
	SoapAction string `json:"soapAction"`
	// One or more soap actions of the service to resolved. This must be specified in the absence of soapAction field.
	SoapActions []string `json:"soapActions"`
	// Base uri from the wsdl of the service. Use this alongside the soapaction
	// property to resolve a soap service without resolutionUri
	BaseUri string `json:"baseUri"`
	// The resolution path of the service if that is how the soap service is resolved
	ResolutionPath string `json:"resolutionPath"`
}

// GetSoapAction returns SoapServiceResolverInput.SoapAction, and is useful for accessing the field via an interface.
func (v *SoapServiceResolverInput) GetSoapAction() string { return v.SoapAction }

// GetSoapActions returns SoapServiceResolverInput.SoapActions, and is useful for accessing the field via an interface.
func (v *SoapServiceResolverInput) GetSoapActions() []string { return v.SoapActions }

// GetBaseUri returns SoapServiceResolverInput.BaseUri, and is useful for accessing the field via an interface.
func (v *SoapServiceResolverInput) GetBaseUri() string { return v.BaseUri }

// GetResolutionPath returns SoapServiceResolverInput.ResolutionPath, and is useful for accessing the field via an interface.
func (v *SoapServiceResolverInput) GetResolutionPath() string { return v.ResolutionPath }

type SoapVersion string

const (
	SoapVersionSoap11  SoapVersion = "SOAP_1_1"
	SoapVersionSoap12  SoapVersion = "SOAP_1_2"
	SoapVersionUnknown SoapVersion = "UNKNOWN"
)

type Tag string

const (
	TagAdmin Tag = "ADMIN"
)

// Input sent with createTrustedCert mutation
type TrustedCertInput struct {
	// The internal entity unique identifier
	Goid string `json:"goid"`
	// The name of the trusted certificate
	Name string `json:"name"`
	// The base 64 encoded string of the certificate
	CertBase64 string `json:"certBase64"`
	// Whether to perform hostname verification with this certificate
	VerifyHostname bool `json:"verifyHostname"`
	// Whether this certificate is a trust anchor
	TrustAnchor bool `json:"trustAnchor"`
	// What the certificate is trusted for
	TrustedFor []TrustedForType `json:"trustedFor"`
	// The revocation check policy type
	RevocationCheckPolicyType PolicyUsageType `json:"revocationCheckPolicyType"`
	// The name of revocation policy.  Required if revocationCheckPolicyType is PolicyUsageType.SPECIFIED
	RevocationCheckPolicyName string `json:"revocationCheckPolicyName"`
	// The Subject DN of this certificate. (Note that, this field has no effect on the mutation)
	SubjectDn string `json:"subjectDn"`
	// The start date of the validity period. (Note that, this field has no effect on the mutation)
	NotBefore string `json:"notBefore"`
	// the end date of the validity period. (Note that, this field has no effect on the mutation)
	NotAfter string `json:"notAfter"`
	// The sha1 thumbprint of the certificate. (Note that, this field has no effect on the mutation)
	ThumbprintSha1 string `json:"thumbprintSha1"`
	// Ignored at creation time but can be used to compare bundle with gw state
	Checksum string `json:"checksum"`
}

// GetGoid returns TrustedCertInput.Goid, and is useful for accessing the field via an interface.
func (v *TrustedCertInput) GetGoid() string { return v.Goid }

// GetName returns TrustedCertInput.Name, and is useful for accessing the field via an interface.
func (v *TrustedCertInput) GetName() string { return v.Name }

// GetCertBase64 returns TrustedCertInput.CertBase64, and is useful for accessing the field via an interface.
func (v *TrustedCertInput) GetCertBase64() string { return v.CertBase64 }

// GetVerifyHostname returns TrustedCertInput.VerifyHostname, and is useful for accessing the field via an interface.
func (v *TrustedCertInput) GetVerifyHostname() bool { return v.VerifyHostname }

// GetTrustAnchor returns TrustedCertInput.TrustAnchor, and is useful for accessing the field via an interface.
func (v *TrustedCertInput) GetTrustAnchor() bool { return v.TrustAnchor }

// GetTrustedFor returns TrustedCertInput.TrustedFor, and is useful for accessing the field via an interface.
func (v *TrustedCertInput) GetTrustedFor() []TrustedForType { return v.TrustedFor }

// GetRevocationCheckPolicyType returns TrustedCertInput.RevocationCheckPolicyType, and is useful for accessing the field via an interface.
func (v *TrustedCertInput) GetRevocationCheckPolicyType() PolicyUsageType {
	return v.RevocationCheckPolicyType
}

// GetRevocationCheckPolicyName returns TrustedCertInput.RevocationCheckPolicyName, and is useful for accessing the field via an interface.
func (v *TrustedCertInput) GetRevocationCheckPolicyName() string { return v.RevocationCheckPolicyName }

// GetSubjectDn returns TrustedCertInput.SubjectDn, and is useful for accessing the field via an interface.
func (v *TrustedCertInput) GetSubjectDn() string { return v.SubjectDn }

// GetNotBefore returns TrustedCertInput.NotBefore, and is useful for accessing the field via an interface.
func (v *TrustedCertInput) GetNotBefore() string { return v.NotBefore }

// GetNotAfter returns TrustedCertInput.NotAfter, and is useful for accessing the field via an interface.
func (v *TrustedCertInput) GetNotAfter() string { return v.NotAfter }

// GetThumbprintSha1 returns TrustedCertInput.ThumbprintSha1, and is useful for accessing the field via an interface.
func (v *TrustedCertInput) GetThumbprintSha1() string { return v.ThumbprintSha1 }

// GetChecksum returns TrustedCertInput.Checksum, and is useful for accessing the field via an interface.
func (v *TrustedCertInput) GetChecksum() string { return v.Checksum }

// Partial TrustedCert input for updates
type TrustedCertPartialInput struct {
	// The internal entity unique identifier
	Goid string `json:"goid"`
	// The name of the trusted certificate
	Name string `json:"name"`
	// The base 64 encoded string of the certificate
	CertBase64 string `json:"certBase64"`
	// Whether to perform hostname verification with this certificate
	VerifyHostname bool `json:"verifyHostname"`
	// Whether this certificate is a trust anchor
	TrustAnchor bool `json:"trustAnchor"`
	// What the certificate is trusted for
	TrustedFor []TrustedForType `json:"trustedFor"`
	// The revocation check policy type
	RevocationCheckPolicyType PolicyUsageType `json:"revocationCheckPolicyType"`
	// The name of revocation policy.  Required if revocationCheckPolicyType is PolicyUsageType.SPECIFIED
	RevocationCheckPolicyName string `json:"revocationCheckPolicyName"`
	// The Subject DN of this certificate. (Note that, this field has no effect on the mutation)
	SubjectDn string `json:"subjectDn"`
	// The start date of the validity period. (Note that, this field has no effect on the mutation)
	NotBefore string `json:"notBefore"`
	// the end date of the validity period. (Note that, this field has no effect on the mutation)
	NotAfter string `json:"notAfter"`
	// The sha1 thumbprint of the certificate. This field is used to find the existing record.
	ThumbprintSha1 string `json:"thumbprintSha1"`
	// Ignored at creation time but can be used to compare bundle with gw state
	Checksum string `json:"checksum"`
}

// GetGoid returns TrustedCertPartialInput.Goid, and is useful for accessing the field via an interface.
func (v *TrustedCertPartialInput) GetGoid() string { return v.Goid }

// GetName returns TrustedCertPartialInput.Name, and is useful for accessing the field via an interface.
func (v *TrustedCertPartialInput) GetName() string { return v.Name }

// GetCertBase64 returns TrustedCertPartialInput.CertBase64, and is useful for accessing the field via an interface.
func (v *TrustedCertPartialInput) GetCertBase64() string { return v.CertBase64 }

// GetVerifyHostname returns TrustedCertPartialInput.VerifyHostname, and is useful for accessing the field via an interface.
func (v *TrustedCertPartialInput) GetVerifyHostname() bool { return v.VerifyHostname }

// GetTrustAnchor returns TrustedCertPartialInput.TrustAnchor, and is useful for accessing the field via an interface.
func (v *TrustedCertPartialInput) GetTrustAnchor() bool { return v.TrustAnchor }

// GetTrustedFor returns TrustedCertPartialInput.TrustedFor, and is useful for accessing the field via an interface.
func (v *TrustedCertPartialInput) GetTrustedFor() []TrustedForType { return v.TrustedFor }

// GetRevocationCheckPolicyType returns TrustedCertPartialInput.RevocationCheckPolicyType, and is useful for accessing the field via an interface.
func (v *TrustedCertPartialInput) GetRevocationCheckPolicyType() PolicyUsageType {
	return v.RevocationCheckPolicyType
}

// GetRevocationCheckPolicyName returns TrustedCertPartialInput.RevocationCheckPolicyName, and is useful for accessing the field via an interface.
func (v *TrustedCertPartialInput) GetRevocationCheckPolicyName() string {
	return v.RevocationCheckPolicyName
}

// GetSubjectDn returns TrustedCertPartialInput.SubjectDn, and is useful for accessing the field via an interface.
func (v *TrustedCertPartialInput) GetSubjectDn() string { return v.SubjectDn }

// GetNotBefore returns TrustedCertPartialInput.NotBefore, and is useful for accessing the field via an interface.
func (v *TrustedCertPartialInput) GetNotBefore() string { return v.NotBefore }

// GetNotAfter returns TrustedCertPartialInput.NotAfter, and is useful for accessing the field via an interface.
func (v *TrustedCertPartialInput) GetNotAfter() string { return v.NotAfter }

// GetThumbprintSha1 returns TrustedCertPartialInput.ThumbprintSha1, and is useful for accessing the field via an interface.
func (v *TrustedCertPartialInput) GetThumbprintSha1() string { return v.ThumbprintSha1 }

// GetChecksum returns TrustedCertPartialInput.Checksum, and is useful for accessing the field via an interface.
func (v *TrustedCertPartialInput) GetChecksum() string { return v.Checksum }

// Defines what a certificate is trusted for
type TrustedForType string

const (
	// Is trusted as an SSL server cert
	TrustedForTypeSsl TrustedForType = "SSL"
	// Is trusted as a CA that signs SSL server certs
	TrustedForTypeSigningServerCerts TrustedForType = "SIGNING_SERVER_CERTS"
	// Is trusted as a CA that signs SSL client certs
	TrustedForTypeSigningClientCerts TrustedForType = "SIGNING_CLIENT_CERTS"
	// Is trusted to sign SAML tokens
	TrustedForTypeSamlIssuer TrustedForType = "SAML_ISSUER"
	// Is trusted as a SAML attesting entity
	TrustedForTypeSamlAttestingEntity TrustedForType = "SAML_ATTESTING_ENTITY"
)

type UserMappingInput struct {
	ObjClass                   string               `json:"objClass"`
	NameAttrName               string               `json:"nameAttrName"`
	LoginAttrName              string               `json:"loginAttrName"`
	PasswdAttrName             string               `json:"passwdAttrName"`
	FirstNameAttrName          string               `json:"firstNameAttrName"`
	LastNameAttrName           string               `json:"lastNameAttrName"`
	EmailNameAttrName          string               `json:"emailNameAttrName"`
	KerberosAttrName           string               `json:"kerberosAttrName"`
	KerberosEnterpriseAttrName string               `json:"kerberosEnterpriseAttrName"`
	UserCertAttrName           string               `json:"userCertAttrName"`
	PasswdType                 *PasswdStrategyInput `json:"passwdType,omitempty"`
}

// GetObjClass returns UserMappingInput.ObjClass, and is useful for accessing the field via an interface.
func (v *UserMappingInput) GetObjClass() string { return v.ObjClass }

// GetNameAttrName returns UserMappingInput.NameAttrName, and is useful for accessing the field via an interface.
func (v *UserMappingInput) GetNameAttrName() string { return v.NameAttrName }

// GetLoginAttrName returns UserMappingInput.LoginAttrName, and is useful for accessing the field via an interface.
func (v *UserMappingInput) GetLoginAttrName() string { return v.LoginAttrName }

// GetPasswdAttrName returns UserMappingInput.PasswdAttrName, and is useful for accessing the field via an interface.
func (v *UserMappingInput) GetPasswdAttrName() string { return v.PasswdAttrName }

// GetFirstNameAttrName returns UserMappingInput.FirstNameAttrName, and is useful for accessing the field via an interface.
func (v *UserMappingInput) GetFirstNameAttrName() string { return v.FirstNameAttrName }

// GetLastNameAttrName returns UserMappingInput.LastNameAttrName, and is useful for accessing the field via an interface.
func (v *UserMappingInput) GetLastNameAttrName() string { return v.LastNameAttrName }

// GetEmailNameAttrName returns UserMappingInput.EmailNameAttrName, and is useful for accessing the field via an interface.
func (v *UserMappingInput) GetEmailNameAttrName() string { return v.EmailNameAttrName }

// GetKerberosAttrName returns UserMappingInput.KerberosAttrName, and is useful for accessing the field via an interface.
func (v *UserMappingInput) GetKerberosAttrName() string { return v.KerberosAttrName }

// GetKerberosEnterpriseAttrName returns UserMappingInput.KerberosEnterpriseAttrName, and is useful for accessing the field via an interface.
func (v *UserMappingInput) GetKerberosEnterpriseAttrName() string {
	return v.KerberosEnterpriseAttrName
}

// GetUserCertAttrName returns UserMappingInput.UserCertAttrName, and is useful for accessing the field via an interface.
func (v *UserMappingInput) GetUserCertAttrName() string { return v.UserCertAttrName }

// GetPasswdType returns UserMappingInput.PasswdType, and is useful for accessing the field via an interface.
func (v *UserMappingInput) GetPasswdType() *PasswdStrategyInput { return v.PasswdType }

// IDP User Reference input
type UserRefInput struct {
	// The name of user
	Name string `json:"name"`
	// The login identity of user
	Login string `json:"login"`
	// The DN of user
	SubjectDn string `json:"subjectDn"`
	// The name of identity provider that the user belongs to
	ProviderName string `json:"providerName"`
	// The type of identity provider that the user belongs to
	ProviderType IdpType `json:"providerType"`
}

// GetName returns UserRefInput.Name, and is useful for accessing the field via an interface.
func (v *UserRefInput) GetName() string { return v.Name }

// GetLogin returns UserRefInput.Login, and is useful for accessing the field via an interface.
func (v *UserRefInput) GetLogin() string { return v.Login }

// GetSubjectDn returns UserRefInput.SubjectDn, and is useful for accessing the field via an interface.
func (v *UserRefInput) GetSubjectDn() string { return v.SubjectDn }

// GetProviderName returns UserRefInput.ProviderName, and is useful for accessing the field via an interface.
func (v *UserRefInput) GetProviderName() string { return v.ProviderName }

// GetProviderType returns UserRefInput.ProviderType, and is useful for accessing the field via an interface.
func (v *UserRefInput) GetProviderType() IdpType { return v.ProviderType }

type WebApiServiceInput struct {
	// The internal entity unique identifier
	Goid string `json:"goid"`
	// The guid for this service, if none provided, assigned at creation
	Guid string `json:"guid"`
	// The folder path where to create this service.  If the path does not exist, it will be created
	FolderPath string `json:"folderPath"`
	// The name of the service
	Name string `json:"name"`
	// The resolution path of the service
	ResolutionPath string `json:"resolutionPath"`
	// The policy
	Policy *PolicyInput `json:"policy,omitempty"`
	// Whether the service is enabled (optional, default to true)
	Enabled bool `json:"enabled"`
	// The http methods allowed for this service
	MethodsAllowed       []HttpMethod           `json:"methodsAllowed"`
	TracingEnabled       bool                   `json:"tracingEnabled"`
	WssProcessingEnabled bool                   `json:"wssProcessingEnabled"`
	Properties           []*EntityPropertyInput `json:"properties,omitempty"`
	// Ignored at creation time but can be used to compare bundle with gw state
	Checksum string `json:"checksum"`
}

// GetGoid returns WebApiServiceInput.Goid, and is useful for accessing the field via an interface.
func (v *WebApiServiceInput) GetGoid() string { return v.Goid }

// GetGuid returns WebApiServiceInput.Guid, and is useful for accessing the field via an interface.
func (v *WebApiServiceInput) GetGuid() string { return v.Guid }

// GetFolderPath returns WebApiServiceInput.FolderPath, and is useful for accessing the field via an interface.
func (v *WebApiServiceInput) GetFolderPath() string { return v.FolderPath }

// GetName returns WebApiServiceInput.Name, and is useful for accessing the field via an interface.
func (v *WebApiServiceInput) GetName() string { return v.Name }

// GetResolutionPath returns WebApiServiceInput.ResolutionPath, and is useful for accessing the field via an interface.
func (v *WebApiServiceInput) GetResolutionPath() string { return v.ResolutionPath }

// GetPolicy returns WebApiServiceInput.Policy, and is useful for accessing the field via an interface.
func (v *WebApiServiceInput) GetPolicy() *PolicyInput { return v.Policy }

// GetEnabled returns WebApiServiceInput.Enabled, and is useful for accessing the field via an interface.
func (v *WebApiServiceInput) GetEnabled() bool { return v.Enabled }

// GetMethodsAllowed returns WebApiServiceInput.MethodsAllowed, and is useful for accessing the field via an interface.
func (v *WebApiServiceInput) GetMethodsAllowed() []HttpMethod { return v.MethodsAllowed }

// GetTracingEnabled returns WebApiServiceInput.TracingEnabled, and is useful for accessing the field via an interface.
func (v *WebApiServiceInput) GetTracingEnabled() bool { return v.TracingEnabled }

// GetWssProcessingEnabled returns WebApiServiceInput.WssProcessingEnabled, and is useful for accessing the field via an interface.
func (v *WebApiServiceInput) GetWssProcessingEnabled() bool { return v.WssProcessingEnabled }

// GetProperties returns WebApiServiceInput.Properties, and is useful for accessing the field via an interface.
func (v *WebApiServiceInput) GetProperties() []*EntityPropertyInput { return v.Properties }

// GetChecksum returns WebApiServiceInput.Checksum, and is useful for accessing the field via an interface.
func (v *WebApiServiceInput) GetChecksum() string { return v.Checksum }

// __deleteKeysInput is used internally by genqlient
type __deleteKeysInput struct {
	Keys []string `json:"keys"`
}

// GetKeys returns __deleteKeysInput.Keys, and is useful for accessing the field via an interface.
func (v *__deleteKeysInput) GetKeys() []string { return v.Keys }

// __deleteL7PortalApiInput is used internally by genqlient
type __deleteL7PortalApiInput struct {
	WebApiServiceResolutionPaths []string `json:"webApiServiceResolutionPaths"`
	PolicyFragmentNames          []string `json:"policyFragmentNames"`
}

// GetWebApiServiceResolutionPaths returns __deleteL7PortalApiInput.WebApiServiceResolutionPaths, and is useful for accessing the field via an interface.
func (v *__deleteL7PortalApiInput) GetWebApiServiceResolutionPaths() []string {
	return v.WebApiServiceResolutionPaths
}

// GetPolicyFragmentNames returns __deleteL7PortalApiInput.PolicyFragmentNames, and is useful for accessing the field via an interface.
func (v *__deleteL7PortalApiInput) GetPolicyFragmentNames() []string { return v.PolicyFragmentNames }

// __deleteSecretsInput is used internally by genqlient
type __deleteSecretsInput struct {
	Secrets []string `json:"secrets"`
}

// GetSecrets returns __deleteSecretsInput.Secrets, and is useful for accessing the field via an interface.
func (v *__deleteSecretsInput) GetSecrets() []string { return v.Secrets }

// __installBundleInput is used internally by genqlient
type __installBundleInput struct {
	ActiveConnectors                    []*ActiveConnectorInput                   `json:"activeConnectors,omitempty"`
	AdministrativeUserAccountProperties []*AdministrativeUserAccountPropertyInput `json:"administrativeUserAccountProperties,omitempty"`
	BackgroundTaskPolicies              []*BackgroundTaskPolicyInput              `json:"backgroundTaskPolicies,omitempty"`
	CassandraConnections                []*CassandraConnectionInput               `json:"cassandraConnections,omitempty"`
	ClusterProperties                   []*ClusterPropertyInput                   `json:"clusterProperties,omitempty"`
	Dtds                                []*DtdInput                               `json:"dtds,omitempty"`
	EmailListeners                      []*EmailListenerInput                     `json:"emailListeners,omitempty"`
	EncassConfigs                       []*EncassConfigInput                      `json:"encassConfigs,omitempty"`
	FipGroups                           []*FipGroupInput                          `json:"fipGroups,omitempty"`
	FipUsers                            []*FipUserInput                           `json:"fipUsers,omitempty"`
	Fips                                []*FipInput                               `json:"fips,omitempty"`
	FederatedGroups                     []*FederatedGroupInput                    `json:"federatedGroups,omitempty"`
	FederatedUsers                      []*FederatedUserInput                     `json:"federatedUsers,omitempty"`
	InternalIdps                        []*InternalIdpInput                       `json:"internalIdps,omitempty"`
	FederatedIdps                       []*FederatedIdpInput                      `json:"federatedIdps,omitempty"`
	LdapIdps                            []*LdapIdpInput                           `json:"ldapIdps,omitempty"`
	SimpleLdapIdps                      []*SimpleLdapIdpInput                     `json:"simpleLdapIdps,omitempty"`
	PolicyBackedIdps                    []*PolicyBackedIdpInput                   `json:"policyBackedIdps,omitempty"`
	GlobalPolicies                      []*GlobalPolicyInput                      `json:"globalPolicies,omitempty"`
	InternalGroups                      []*InternalGroupInput                     `json:"internalGroups,omitempty"`
	InternalSoapServices                []*SoapServiceInput                       `json:"internalSoapServices,omitempty"`
	InternalUsers                       []*InternalUserInput                      `json:"internalUsers,omitempty"`
	InternalWebApiServices              []*WebApiServiceInput                     `json:"internalWebApiServices,omitempty"`
	JdbcConnections                     []*JdbcConnectionInput                    `json:"jdbcConnections,omitempty"`
	JmsDestinations                     []*JmsDestinationInput                    `json:"jmsDestinations,omitempty"`
	Keys                                []*KeyInput                               `json:"keys,omitempty"`
	Ldaps                               []*LdapInput                              `json:"ldaps,omitempty"`
	Roles                               []*RoleInput                              `json:"roles,omitempty"`
	ListenPorts                         []*ListenPortInput                        `json:"listenPorts,omitempty"`
	PasswordPolicies                    []*PasswordPolicyInput                    `json:"passwordPolicies,omitempty"`
	Policies                            []*L7PolicyInput                          `json:"policies,omitempty"`
	PolicyFragments                     []*PolicyFragmentInput                    `json:"policyFragments,omitempty"`
	RevocationCheckPolicies             []*RevocationCheckPolicyInput             `json:"revocationCheckPolicies,omitempty"`
	ScheduledTasks                      []*ScheduledTaskInput                     `json:"scheduledTasks,omitempty"`
	LogSinks                            []*LogSinkInput                           `json:"logSinks,omitempty"`
	Schemas                             []*SchemaInput                            `json:"schemas,omitempty"`
	Secrets                             []*SecretInput                            `json:"secrets,omitempty"`
	HttpConfigurations                  []*HttpConfigurationInput                 `json:"httpConfigurations,omitempty"`
	CustomKeyValues                     []*CustomKeyValueInput                    `json:"customKeyValues,omitempty"`
	ServerModuleFiles                   []*ServerModuleFileInput                  `json:"serverModuleFiles,omitempty"`
	ServiceResolutionConfigs            []*ServiceResolutionConfigInput           `json:"serviceResolutionConfigs,omitempty"`
	Folders                             []*FolderInput                            `json:"folders,omitempty"`
	SmConfigs                           []*SMConfigInput                          `json:"smConfigs,omitempty"`
	Services                            []*L7ServiceInput                         `json:"services,omitempty"`
	SoapServices                        []*SoapServiceInput                       `json:"soapServices,omitempty"`
	TrustedCerts                        []*TrustedCertInput                       `json:"trustedCerts,omitempty"`
	WebApiServices                      []*WebApiServiceInput                     `json:"webApiServices,omitempty"`
	GenericEntities                     []*GenericEntityInput                     `json:"genericEntities,omitempty"`
	AuditConfigurations                 []*AuditConfigurationInput                `json:"auditConfigurations,omitempty"`
}

// GetActiveConnectors returns __installBundleInput.ActiveConnectors, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetActiveConnectors() []*ActiveConnectorInput {
	return v.ActiveConnectors
}

// GetAdministrativeUserAccountProperties returns __installBundleInput.AdministrativeUserAccountProperties, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetAdministrativeUserAccountProperties() []*AdministrativeUserAccountPropertyInput {
	return v.AdministrativeUserAccountProperties
}

// GetBackgroundTaskPolicies returns __installBundleInput.BackgroundTaskPolicies, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetBackgroundTaskPolicies() []*BackgroundTaskPolicyInput {
	return v.BackgroundTaskPolicies
}

// GetCassandraConnections returns __installBundleInput.CassandraConnections, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetCassandraConnections() []*CassandraConnectionInput {
	return v.CassandraConnections
}

// GetClusterProperties returns __installBundleInput.ClusterProperties, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetClusterProperties() []*ClusterPropertyInput {
	return v.ClusterProperties
}

// GetDtds returns __installBundleInput.Dtds, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetDtds() []*DtdInput { return v.Dtds }

// GetEmailListeners returns __installBundleInput.EmailListeners, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetEmailListeners() []*EmailListenerInput { return v.EmailListeners }

// GetEncassConfigs returns __installBundleInput.EncassConfigs, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetEncassConfigs() []*EncassConfigInput { return v.EncassConfigs }

// GetFipGroups returns __installBundleInput.FipGroups, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetFipGroups() []*FipGroupInput { return v.FipGroups }

// GetFipUsers returns __installBundleInput.FipUsers, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetFipUsers() []*FipUserInput { return v.FipUsers }

// GetFips returns __installBundleInput.Fips, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetFips() []*FipInput { return v.Fips }

// GetFederatedGroups returns __installBundleInput.FederatedGroups, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetFederatedGroups() []*FederatedGroupInput { return v.FederatedGroups }

// GetFederatedUsers returns __installBundleInput.FederatedUsers, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetFederatedUsers() []*FederatedUserInput { return v.FederatedUsers }

// GetInternalIdps returns __installBundleInput.InternalIdps, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetInternalIdps() []*InternalIdpInput { return v.InternalIdps }

// GetFederatedIdps returns __installBundleInput.FederatedIdps, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetFederatedIdps() []*FederatedIdpInput { return v.FederatedIdps }

// GetLdapIdps returns __installBundleInput.LdapIdps, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetLdapIdps() []*LdapIdpInput { return v.LdapIdps }

// GetSimpleLdapIdps returns __installBundleInput.SimpleLdapIdps, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetSimpleLdapIdps() []*SimpleLdapIdpInput { return v.SimpleLdapIdps }

// GetPolicyBackedIdps returns __installBundleInput.PolicyBackedIdps, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetPolicyBackedIdps() []*PolicyBackedIdpInput {
	return v.PolicyBackedIdps
}

// GetGlobalPolicies returns __installBundleInput.GlobalPolicies, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetGlobalPolicies() []*GlobalPolicyInput { return v.GlobalPolicies }

// GetInternalGroups returns __installBundleInput.InternalGroups, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetInternalGroups() []*InternalGroupInput { return v.InternalGroups }

// GetInternalSoapServices returns __installBundleInput.InternalSoapServices, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetInternalSoapServices() []*SoapServiceInput {
	return v.InternalSoapServices
}

// GetInternalUsers returns __installBundleInput.InternalUsers, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetInternalUsers() []*InternalUserInput { return v.InternalUsers }

// GetInternalWebApiServices returns __installBundleInput.InternalWebApiServices, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetInternalWebApiServices() []*WebApiServiceInput {
	return v.InternalWebApiServices
}

// GetJdbcConnections returns __installBundleInput.JdbcConnections, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetJdbcConnections() []*JdbcConnectionInput { return v.JdbcConnections }

// GetJmsDestinations returns __installBundleInput.JmsDestinations, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetJmsDestinations() []*JmsDestinationInput { return v.JmsDestinations }

// GetKeys returns __installBundleInput.Keys, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetKeys() []*KeyInput { return v.Keys }

// GetLdaps returns __installBundleInput.Ldaps, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetLdaps() []*LdapInput { return v.Ldaps }

// GetRoles returns __installBundleInput.Roles, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetRoles() []*RoleInput { return v.Roles }

// GetListenPorts returns __installBundleInput.ListenPorts, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetListenPorts() []*ListenPortInput { return v.ListenPorts }

// GetPasswordPolicies returns __installBundleInput.PasswordPolicies, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetPasswordPolicies() []*PasswordPolicyInput {
	return v.PasswordPolicies
}

// GetPolicies returns __installBundleInput.Policies, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetPolicies() []*L7PolicyInput { return v.Policies }

// GetPolicyFragments returns __installBundleInput.PolicyFragments, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetPolicyFragments() []*PolicyFragmentInput { return v.PolicyFragments }

// GetRevocationCheckPolicies returns __installBundleInput.RevocationCheckPolicies, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetRevocationCheckPolicies() []*RevocationCheckPolicyInput {
	return v.RevocationCheckPolicies
}

// GetScheduledTasks returns __installBundleInput.ScheduledTasks, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetScheduledTasks() []*ScheduledTaskInput { return v.ScheduledTasks }

// GetLogSinks returns __installBundleInput.LogSinks, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetLogSinks() []*LogSinkInput { return v.LogSinks }

// GetSchemas returns __installBundleInput.Schemas, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetSchemas() []*SchemaInput { return v.Schemas }

// GetSecrets returns __installBundleInput.Secrets, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetSecrets() []*SecretInput { return v.Secrets }

// GetHttpConfigurations returns __installBundleInput.HttpConfigurations, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetHttpConfigurations() []*HttpConfigurationInput {
	return v.HttpConfigurations
}

// GetCustomKeyValues returns __installBundleInput.CustomKeyValues, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetCustomKeyValues() []*CustomKeyValueInput { return v.CustomKeyValues }

// GetServerModuleFiles returns __installBundleInput.ServerModuleFiles, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetServerModuleFiles() []*ServerModuleFileInput {
	return v.ServerModuleFiles
}

// GetServiceResolutionConfigs returns __installBundleInput.ServiceResolutionConfigs, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetServiceResolutionConfigs() []*ServiceResolutionConfigInput {
	return v.ServiceResolutionConfigs
}

// GetFolders returns __installBundleInput.Folders, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetFolders() []*FolderInput { return v.Folders }

// GetSmConfigs returns __installBundleInput.SmConfigs, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetSmConfigs() []*SMConfigInput { return v.SmConfigs }

// GetServices returns __installBundleInput.Services, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetServices() []*L7ServiceInput { return v.Services }

// GetSoapServices returns __installBundleInput.SoapServices, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetSoapServices() []*SoapServiceInput { return v.SoapServices }

// GetTrustedCerts returns __installBundleInput.TrustedCerts, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetTrustedCerts() []*TrustedCertInput { return v.TrustedCerts }

// GetWebApiServices returns __installBundleInput.WebApiServices, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetWebApiServices() []*WebApiServiceInput { return v.WebApiServices }

// GetGenericEntities returns __installBundleInput.GenericEntities, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetGenericEntities() []*GenericEntityInput { return v.GenericEntities }

// GetAuditConfigurations returns __installBundleInput.AuditConfigurations, and is useful for accessing the field via an interface.
func (v *__installBundleInput) GetAuditConfigurations() []*AuditConfigurationInput {
	return v.AuditConfigurations
}

// deleteKeysDeleteKeysKeysPayload includes the requested fields of the GraphQL type KeysPayload.
type deleteKeysDeleteKeysKeysPayload struct {
	DetailedStatus []*deleteKeysDeleteKeysKeysPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
	Keys           []*deleteKeysDeleteKeysKeysPayloadKeysKey                                    `json:"keys"`
}

// GetDetailedStatus returns deleteKeysDeleteKeysKeysPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *deleteKeysDeleteKeysKeysPayload) GetDetailedStatus() []*deleteKeysDeleteKeysKeysPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// GetKeys returns deleteKeysDeleteKeysKeysPayload.Keys, and is useful for accessing the field via an interface.
func (v *deleteKeysDeleteKeysKeysPayload) GetKeys() []*deleteKeysDeleteKeysKeysPayloadKeysKey {
	return v.Keys
}

// deleteKeysDeleteKeysKeysPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type deleteKeysDeleteKeysKeysPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Status      EntityMutationStatus `json:"status"`
	Description string               `json:"description"`
}

// GetStatus returns deleteKeysDeleteKeysKeysPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *deleteKeysDeleteKeysKeysPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns deleteKeysDeleteKeysKeysPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *deleteKeysDeleteKeysKeysPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// deleteKeysDeleteKeysKeysPayloadKeysKey includes the requested fields of the GraphQL type Key.
// The GraphQL type's documentation follows.
//
// A key entry in the gateway keystore. These entries combine a private
// key and associated certificate and are used for example by listener ports.
// > @l7-entity key|keys
// > @l7-identity-fields alias,keystoreId
// > @l7-summary-fields goid,keystoreId,alias,checksum
// > @l7-excluded-fields pem
type deleteKeysDeleteKeysKeysPayloadKeysKey struct {
	// The internal entity unique identifier
	Goid string `json:"goid"`
	// The gateway keystore identifier
	KeystoreId string `json:"keystoreId"`
	// The name assigned to the key
	Alias string `json:"alias"`
}

// GetGoid returns deleteKeysDeleteKeysKeysPayloadKeysKey.Goid, and is useful for accessing the field via an interface.
func (v *deleteKeysDeleteKeysKeysPayloadKeysKey) GetGoid() string { return v.Goid }

// GetKeystoreId returns deleteKeysDeleteKeysKeysPayloadKeysKey.KeystoreId, and is useful for accessing the field via an interface.
func (v *deleteKeysDeleteKeysKeysPayloadKeysKey) GetKeystoreId() string { return v.KeystoreId }

// GetAlias returns deleteKeysDeleteKeysKeysPayloadKeysKey.Alias, and is useful for accessing the field via an interface.
func (v *deleteKeysDeleteKeysKeysPayloadKeysKey) GetAlias() string { return v.Alias }

// deleteKeysResponse is returned by deleteKeys on success.
type deleteKeysResponse struct {
	// Deletes one or more existing keys
	DeleteKeys *deleteKeysDeleteKeysKeysPayload `json:"deleteKeys"`
}

// GetDeleteKeys returns deleteKeysResponse.DeleteKeys, and is useful for accessing the field via an interface.
func (v *deleteKeysResponse) GetDeleteKeys() *deleteKeysDeleteKeysKeysPayload { return v.DeleteKeys }

// deleteL7PortalApiDeletePolicyFragmentsPolicyFragmentsPayload includes the requested fields of the GraphQL type PolicyFragmentsPayload.
type deleteL7PortalApiDeletePolicyFragmentsPolicyFragmentsPayload struct {
	DetailedStatus []*deleteL7PortalApiDeletePolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns deleteL7PortalApiDeletePolicyFragmentsPolicyFragmentsPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *deleteL7PortalApiDeletePolicyFragmentsPolicyFragmentsPayload) GetDetailedStatus() []*deleteL7PortalApiDeletePolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// deleteL7PortalApiDeletePolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type deleteL7PortalApiDeletePolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Status      EntityMutationStatus `json:"status"`
	Description string               `json:"description"`
}

// GetStatus returns deleteL7PortalApiDeletePolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *deleteL7PortalApiDeletePolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns deleteL7PortalApiDeletePolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *deleteL7PortalApiDeletePolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// deleteL7PortalApiDeleteWebApiServicesWebApiServicesPayload includes the requested fields of the GraphQL type WebApiServicesPayload.
type deleteL7PortalApiDeleteWebApiServicesWebApiServicesPayload struct {
	DetailedStatus []*deleteL7PortalApiDeleteWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns deleteL7PortalApiDeleteWebApiServicesWebApiServicesPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *deleteL7PortalApiDeleteWebApiServicesWebApiServicesPayload) GetDetailedStatus() []*deleteL7PortalApiDeleteWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// deleteL7PortalApiDeleteWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type deleteL7PortalApiDeleteWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Status      EntityMutationStatus `json:"status"`
	Description string               `json:"description"`
}

// GetStatus returns deleteL7PortalApiDeleteWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *deleteL7PortalApiDeleteWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns deleteL7PortalApiDeleteWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *deleteL7PortalApiDeleteWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// deleteL7PortalApiResponse is returned by deleteL7PortalApi on success.
type deleteL7PortalApiResponse struct {
	// Delete existing web api services given their resolution paths
	DeleteWebApiServices *deleteL7PortalApiDeleteWebApiServicesWebApiServicesPayload `json:"deleteWebApiServices"`
	// Delete policy fragments
	DeletePolicyFragments *deleteL7PortalApiDeletePolicyFragmentsPolicyFragmentsPayload `json:"deletePolicyFragments"`
}

// GetDeleteWebApiServices returns deleteL7PortalApiResponse.DeleteWebApiServices, and is useful for accessing the field via an interface.
func (v *deleteL7PortalApiResponse) GetDeleteWebApiServices() *deleteL7PortalApiDeleteWebApiServicesWebApiServicesPayload {
	return v.DeleteWebApiServices
}

// GetDeletePolicyFragments returns deleteL7PortalApiResponse.DeletePolicyFragments, and is useful for accessing the field via an interface.
func (v *deleteL7PortalApiResponse) GetDeletePolicyFragments() *deleteL7PortalApiDeletePolicyFragmentsPolicyFragmentsPayload {
	return v.DeletePolicyFragments
}

// deleteSecretsDeleteSecretsSecretsPayload includes the requested fields of the GraphQL type SecretsPayload.
type deleteSecretsDeleteSecretsSecretsPayload struct {
	DetailedStatus []*deleteSecretsDeleteSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
	Secrets        []*deleteSecretsDeleteSecretsSecretsPayloadSecretsSecret                              `json:"secrets"`
}

// GetDetailedStatus returns deleteSecretsDeleteSecretsSecretsPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *deleteSecretsDeleteSecretsSecretsPayload) GetDetailedStatus() []*deleteSecretsDeleteSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// GetSecrets returns deleteSecretsDeleteSecretsSecretsPayload.Secrets, and is useful for accessing the field via an interface.
func (v *deleteSecretsDeleteSecretsSecretsPayload) GetSecrets() []*deleteSecretsDeleteSecretsSecretsPayloadSecretsSecret {
	return v.Secrets
}

// deleteSecretsDeleteSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type deleteSecretsDeleteSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Status      EntityMutationStatus `json:"status"`
	Description string               `json:"description"`
}

// GetStatus returns deleteSecretsDeleteSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *deleteSecretsDeleteSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns deleteSecretsDeleteSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *deleteSecretsDeleteSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// deleteSecretsDeleteSecretsSecretsPayloadSecretsSecret includes the requested fields of the GraphQL type Secret.
// The GraphQL type's documentation follows.
//
// A secret (password or private key) which is used by gateway policies and other configurations.
// > @l7-entity
// > @l7-identity-fields name
// > @l7-summary-fields goid,name,checksum
// > @l7-excluded-fields
type deleteSecretsDeleteSecretsSecretsPayloadSecretsSecret struct {
	// The goid for the Secret
	Goid string `json:"goid"`
	// Identify the password being stored. You may use letters, numbers, dashes, and underscores.
	// Names that contain spaces or periods are valid, but the resulting stored
	// password cannot be referenced via context variable.
	// Names that contain @ or $ are valid, but the resulting stored password cannot be referenced via context variable.
	Name string `json:"name"`
}

// GetGoid returns deleteSecretsDeleteSecretsSecretsPayloadSecretsSecret.Goid, and is useful for accessing the field via an interface.
func (v *deleteSecretsDeleteSecretsSecretsPayloadSecretsSecret) GetGoid() string { return v.Goid }

// GetName returns deleteSecretsDeleteSecretsSecretsPayloadSecretsSecret.Name, and is useful for accessing the field via an interface.
func (v *deleteSecretsDeleteSecretsSecretsPayloadSecretsSecret) GetName() string { return v.Name }

// deleteSecretsResponse is returned by deleteSecrets on success.
type deleteSecretsResponse struct {
	// Deletes one or more existing secrets
	DeleteSecrets *deleteSecretsDeleteSecretsSecretsPayload `json:"deleteSecrets"`
}

// GetDeleteSecrets returns deleteSecretsResponse.DeleteSecrets, and is useful for accessing the field via an interface.
func (v *deleteSecretsResponse) GetDeleteSecrets() *deleteSecretsDeleteSecretsSecretsPayload {
	return v.DeleteSecrets
}

// installBundleGenericInstallBundleEntitiesBundleEntitiesPayload includes the requested fields of the GraphQL type BundleEntitiesPayload.
type installBundleGenericInstallBundleEntitiesBundleEntitiesPayload struct {
	Summary bool `json:"summary"`
}

// GetSummary returns installBundleGenericInstallBundleEntitiesBundleEntitiesPayload.Summary, and is useful for accessing the field via an interface.
func (v *installBundleGenericInstallBundleEntitiesBundleEntitiesPayload) GetSummary() bool {
	return v.Summary
}

// installBundleGenericResponse is returned by installBundleGeneric on success.
type installBundleGenericResponse struct {
	// Installs bundle of entities using set-based mutation operations
	InstallBundleEntities *installBundleGenericInstallBundleEntitiesBundleEntitiesPayload `json:"installBundleEntities"`
}

// GetInstallBundleEntities returns installBundleGenericResponse.InstallBundleEntities, and is useful for accessing the field via an interface.
func (v *installBundleGenericResponse) GetInstallBundleEntities() *installBundleGenericInstallBundleEntitiesBundleEntitiesPayload {
	return v.InstallBundleEntities
}

// installBundleResponse is returned by installBundle on success.
type installBundleResponse struct {
	// Sets Server module files. Updating the existing server module file is unsupported.
	SetServerModuleFiles *installBundleSetServerModuleFilesServerModuleFilesPayload `json:"setServerModuleFiles"`
	// Create or update existing cluster properties.  If a cluster property with the given name does not
	// exist, one will be created, otherwise the existing one will be updated. This returns the list of
	// entities created and/or updated
	SetClusterProperties *installBundleSetClusterPropertiesClusterPropertiesPayload `json:"setClusterProperties"`
	// Update Service Resolution Configs
	SetServiceResolutionConfigs *installBundleSetServiceResolutionConfigsServiceResolutionConfigsPayLoad `json:"setServiceResolutionConfigs"`
	// Set/Update the Password Policies
	SetPasswordPolicies *installBundleSetPasswordPoliciesPasswordPoliciesPayLoad `json:"setPasswordPolicies"`
	// Create or update existing Administrative User Account Minimum cluster properties.
	// If Administrative User Account Minimum cluster property with the given name
	// does not exist, one will be created, otherwise the existing one will be updated.
	// This returns the list of entities created and/or updated.
	// Below are the allowed Administrative User Account Minimum properties
	// logonMaxAllowableAttempts : Logon attempts must be between 1 and 20
	// logonLockoutTime : Lockout period must be between 1 and 86400 seconds
	// logonSessionExpiry : Expiry period must be between 1 and 86400 seconds
	// logonInactivityPeriod : Inactivity period must be between 1 and 365 days
	SetAdministrativeUserAccountProperties *installBundleSetAdministrativeUserAccountPropertiesAdministrativeUserAccountPropertiesPayload `json:"setAdministrativeUserAccountProperties"`
	// Set the Folders
	SetFolders *installBundleSetFoldersFoldersPayload `json:"setFolders"`
	// Create or update existing revocation check policies.
	// Match is carried by name. If match is found, it will be updated. Otherwise, it will be created.
	SetRevocationCheckPolicies *installBundleSetRevocationCheckPoliciesRevocationCheckPoliciesPayload `json:"setRevocationCheckPolicies"`
	// Create or update trusted certificates.
	// If a certificate with the same sha1 thumbprint already exist, it will be updated.
	SetTrustedCerts *installBundleSetTrustedCertsTrustedCertsPayload `json:"setTrustedCerts"`
	// Creates or updates one or more secrets
	SetSecrets *installBundleSetSecretsSecretsPayload `json:"setSecrets"`
	// Create or update existing http configuration.
	SetHttpConfigurations *installBundleSetHttpConfigurationsHttpConfigurationsPayload `json:"setHttpConfigurations"`
	// Create or update existing custom key values data.  If a custom key value with the given key does not
	// exist, one will be created, otherwise the existing one will be updated. This returns the list of
	// entities created and/or updated
	SetCustomKeyValues *installBundleSetCustomKeyValuesCustomKeyValuePayload `json:"setCustomKeyValues"`
	// Create or Update multiple XML schemas
	SetSchemas *installBundleSetSchemasSchemasPayload `json:"setSchemas"`
	// Create or Update multiple DTD resources
	SetDtds *installBundleSetDtdsDtdsPayload `json:"setDtds"`
	// Create or update JDBC connections.
	// If JDBC connection with the same name exist, the JDBC connection will be updated.
	// If no JDBC connection with the name exist, a new JDBC connection will be created.
	SetJdbcConnections *installBundleSetJdbcConnectionsJdbcConnectionsPayload `json:"setJdbcConnections"`
	// Creates or updates one ore more internal IDP configurations
	SetInternalIdps *installBundleSetInternalIdpsInternalIdpsPayload `json:"setInternalIdps"`
	// Creates or updates one or more fips
	SetFederatedIdps *installBundleSetFederatedIdpsFederatedIdpsPayload `json:"setFederatedIdps"`
	// Creates or updates one or more ldaps
	SetLdapIdps *installBundleSetLdapIdpsLdapIdpsPayload `json:"setLdapIdps"`
	// Creates or updates one or more simple ldaps
	SetSimpleLdapIdps *installBundleSetSimpleLdapIdpsSimpleLdapIdpsPayload `json:"setSimpleLdapIdps"`
	// Creates or updates one or more fips
	SetFips *installBundleSetFipsFipsPayload `json:"setFips"`
	// Creates or updates one or more ldaps
	SetLdaps *installBundleSetLdapsLdapsPayload `json:"setLdaps"`
	// Creates or updates one or more fip groups
	SetFederatedGroups *installBundleSetFederatedGroupsFederatedGroupsPayload `json:"setFederatedGroups"`
	// Creates or updates one or more fip groups
	SetFipGroups *installBundleSetFipGroupsFipGroupsPayload `json:"setFipGroups"`
	// Creates or updates one or more internal groups
	SetInternalGroups *installBundleSetInternalGroupsInternalGroupsPayload `json:"setInternalGroups"`
	// Creates or updates one or more fip users.
	// NOTE: Existing user will be found by either login or subjectDn or name.
	SetFederatedUsers *installBundleSetFederatedUsersFederatedUsersPayload `json:"setFederatedUsers"`
	// Creates or updates one or more fip users.
	// NOTE: Existing user will be found by either login or subjectDn or name.
	SetFipUsers *installBundleSetFipUsersFipUsersPayload `json:"setFipUsers"`
	// Creates or updates one or more internal users
	SetInternalUsers *installBundleSetInternalUsersInternalUsersPayload `json:"setInternalUsers"`
	// Create or update Cassandra connections.
	// If Cassandra connection with the same name exist, the Cassandra connection will be updated.
	// If no Cassandra connection with the name exist, a new Cassandra connection will be created.
	SetCassandraConnections *installBundleSetCassandraConnectionsCassandraConnectionsPayload `json:"setCassandraConnections"`
	// Create or update existing siteminder configurations.
	// Match is carried by name. If match is found, it will be updated. Otherwise, it will be created
	SetSMConfigs *installBundleSetSMConfigsSMConfigsPayload `json:"setSMConfigs"`
	// Create or update policies
	SetPolicies *installBundleSetPoliciesL7PoliciesPayload `json:"setPolicies"`
	// Create or update policy fragments
	SetPolicyFragments *installBundleSetPolicyFragmentsPolicyFragmentsPayload `json:"setPolicyFragments"`
	// Create or update Encapsulated Assertion Configurations
	SetEncassConfigs *installBundleSetEncassConfigsEncassConfigsPayload `json:"setEncassConfigs"`
	// Create or update global policies
	SetGlobalPolicies *installBundleSetGlobalPoliciesGlobalPoliciesPayload `json:"setGlobalPolicies"`
	// Creates or updates one or more background task policies
	SetBackgroundTaskPolicies *installBundleSetBackgroundTaskPoliciesBackgroundTaskPoliciesPayload `json:"setBackgroundTaskPolicies"`
	// Create or update services
	SetServices *installBundleSetServicesL7ServicesPayload `json:"setServices"`
	// Create or update web api services
	SetWebApiServices *installBundleSetWebApiServicesWebApiServicesPayload `json:"setWebApiServices"`
	// Create or update soap services
	SetSoapServices *installBundleSetSoapServicesSoapServicesPayload `json:"setSoapServices"`
	// Create or update Internal web api services
	SetInternalWebApiServices *installBundleSetInternalWebApiServicesInternalWebApiServicesPayload `json:"setInternalWebApiServices"`
	// Create or update Internal soap services
	SetInternalSoapServices *installBundleSetInternalSoapServicesInternalSoapServicesPayload `json:"setInternalSoapServices"`
	// Creates or updates one or more policy backed ldaps
	SetPolicyBackedIdps *installBundleSetPolicyBackedIdpsPolicyBackedIdpsPayload `json:"setPolicyBackedIdps"`
	// Create or update JMS destinations.
	// If JMS destination exists, the JMS destination will be updated.
	// If no JMS destination with given name, direction, providerType exist, a new JMS destination will be created.
	SetJmsDestinations *installBundleSetJmsDestinationsJmsDestinationsPayload `json:"setJmsDestinations"`
	// Create or update existing email listeners.
	// Match is carried by name. If match is found, it will be updated. Otherwise, it will be created.
	SetEmailListeners *installBundleSetEmailListenersEmailListenersPayload `json:"setEmailListeners"`
	// Create or update Listen Ports.
	// If Listen Port with the same name exist, the Listen Port will be updated.
	// If no Listen Port with the name exist, a new Listen Port will be created.
	SetListenPorts *installBundleSetListenPortsListenPortsPayload `json:"setListenPorts"`
	// Create or update existing active connector.
	// Match is carried by name. If match is found, it will be updated. Otherwise, it will be created.
	SetActiveConnectors *installBundleSetActiveConnectorsActiveConnectorsPayload `json:"setActiveConnectors"`
	// Creates or updates one or more scheduled tasks
	SetScheduledTasks *installBundleSetScheduledTasksScheduledTasksPayload `json:"setScheduledTasks"`
	// Create or update Log Sinks.
	// If Log Sink with the same name exist, the Log Sink will be updated.
	// If no Log Sink with the name exist, a new Log Sink will be created.
	SetLogSinks *installBundleSetLogSinksLogSinksPayload `json:"setLogSinks"`
	// Create or update existing generic entities.
	// Match is carried by name. If match is found, it will be updated. Otherwise, it will be created.
	SetGenericEntities *installBundleSetGenericEntitiesGenericEntitiesPayload `json:"setGenericEntities"`
	// Update Roles with user/group assignees.
	// Note: Creating a role is unsupported.
	SetRoles               *installBundleSetRolesRolesPayload                             `json:"setRoles"`
	SetAuditConfigurations *installBundleSetAuditConfigurationsAuditConfigurationsPayload `json:"setAuditConfigurations"`
	// Creates or updates one or more keys
	SetKeys *installBundleSetKeysKeysPayload `json:"setKeys"`
}

// GetSetServerModuleFiles returns installBundleResponse.SetServerModuleFiles, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetServerModuleFiles() *installBundleSetServerModuleFilesServerModuleFilesPayload {
	return v.SetServerModuleFiles
}

// GetSetClusterProperties returns installBundleResponse.SetClusterProperties, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetClusterProperties() *installBundleSetClusterPropertiesClusterPropertiesPayload {
	return v.SetClusterProperties
}

// GetSetServiceResolutionConfigs returns installBundleResponse.SetServiceResolutionConfigs, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetServiceResolutionConfigs() *installBundleSetServiceResolutionConfigsServiceResolutionConfigsPayLoad {
	return v.SetServiceResolutionConfigs
}

// GetSetPasswordPolicies returns installBundleResponse.SetPasswordPolicies, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetPasswordPolicies() *installBundleSetPasswordPoliciesPasswordPoliciesPayLoad {
	return v.SetPasswordPolicies
}

// GetSetAdministrativeUserAccountProperties returns installBundleResponse.SetAdministrativeUserAccountProperties, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetAdministrativeUserAccountProperties() *installBundleSetAdministrativeUserAccountPropertiesAdministrativeUserAccountPropertiesPayload {
	return v.SetAdministrativeUserAccountProperties
}

// GetSetFolders returns installBundleResponse.SetFolders, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetFolders() *installBundleSetFoldersFoldersPayload {
	return v.SetFolders
}

// GetSetRevocationCheckPolicies returns installBundleResponse.SetRevocationCheckPolicies, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetRevocationCheckPolicies() *installBundleSetRevocationCheckPoliciesRevocationCheckPoliciesPayload {
	return v.SetRevocationCheckPolicies
}

// GetSetTrustedCerts returns installBundleResponse.SetTrustedCerts, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetTrustedCerts() *installBundleSetTrustedCertsTrustedCertsPayload {
	return v.SetTrustedCerts
}

// GetSetSecrets returns installBundleResponse.SetSecrets, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetSecrets() *installBundleSetSecretsSecretsPayload {
	return v.SetSecrets
}

// GetSetHttpConfigurations returns installBundleResponse.SetHttpConfigurations, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetHttpConfigurations() *installBundleSetHttpConfigurationsHttpConfigurationsPayload {
	return v.SetHttpConfigurations
}

// GetSetCustomKeyValues returns installBundleResponse.SetCustomKeyValues, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetCustomKeyValues() *installBundleSetCustomKeyValuesCustomKeyValuePayload {
	return v.SetCustomKeyValues
}

// GetSetSchemas returns installBundleResponse.SetSchemas, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetSchemas() *installBundleSetSchemasSchemasPayload {
	return v.SetSchemas
}

// GetSetDtds returns installBundleResponse.SetDtds, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetDtds() *installBundleSetDtdsDtdsPayload { return v.SetDtds }

// GetSetJdbcConnections returns installBundleResponse.SetJdbcConnections, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetJdbcConnections() *installBundleSetJdbcConnectionsJdbcConnectionsPayload {
	return v.SetJdbcConnections
}

// GetSetInternalIdps returns installBundleResponse.SetInternalIdps, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetInternalIdps() *installBundleSetInternalIdpsInternalIdpsPayload {
	return v.SetInternalIdps
}

// GetSetFederatedIdps returns installBundleResponse.SetFederatedIdps, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetFederatedIdps() *installBundleSetFederatedIdpsFederatedIdpsPayload {
	return v.SetFederatedIdps
}

// GetSetLdapIdps returns installBundleResponse.SetLdapIdps, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetLdapIdps() *installBundleSetLdapIdpsLdapIdpsPayload {
	return v.SetLdapIdps
}

// GetSetSimpleLdapIdps returns installBundleResponse.SetSimpleLdapIdps, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetSimpleLdapIdps() *installBundleSetSimpleLdapIdpsSimpleLdapIdpsPayload {
	return v.SetSimpleLdapIdps
}

// GetSetFips returns installBundleResponse.SetFips, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetFips() *installBundleSetFipsFipsPayload { return v.SetFips }

// GetSetLdaps returns installBundleResponse.SetLdaps, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetLdaps() *installBundleSetLdapsLdapsPayload { return v.SetLdaps }

// GetSetFederatedGroups returns installBundleResponse.SetFederatedGroups, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetFederatedGroups() *installBundleSetFederatedGroupsFederatedGroupsPayload {
	return v.SetFederatedGroups
}

// GetSetFipGroups returns installBundleResponse.SetFipGroups, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetFipGroups() *installBundleSetFipGroupsFipGroupsPayload {
	return v.SetFipGroups
}

// GetSetInternalGroups returns installBundleResponse.SetInternalGroups, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetInternalGroups() *installBundleSetInternalGroupsInternalGroupsPayload {
	return v.SetInternalGroups
}

// GetSetFederatedUsers returns installBundleResponse.SetFederatedUsers, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetFederatedUsers() *installBundleSetFederatedUsersFederatedUsersPayload {
	return v.SetFederatedUsers
}

// GetSetFipUsers returns installBundleResponse.SetFipUsers, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetFipUsers() *installBundleSetFipUsersFipUsersPayload {
	return v.SetFipUsers
}

// GetSetInternalUsers returns installBundleResponse.SetInternalUsers, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetInternalUsers() *installBundleSetInternalUsersInternalUsersPayload {
	return v.SetInternalUsers
}

// GetSetCassandraConnections returns installBundleResponse.SetCassandraConnections, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetCassandraConnections() *installBundleSetCassandraConnectionsCassandraConnectionsPayload {
	return v.SetCassandraConnections
}

// GetSetSMConfigs returns installBundleResponse.SetSMConfigs, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetSMConfigs() *installBundleSetSMConfigsSMConfigsPayload {
	return v.SetSMConfigs
}

// GetSetPolicies returns installBundleResponse.SetPolicies, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetPolicies() *installBundleSetPoliciesL7PoliciesPayload {
	return v.SetPolicies
}

// GetSetPolicyFragments returns installBundleResponse.SetPolicyFragments, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetPolicyFragments() *installBundleSetPolicyFragmentsPolicyFragmentsPayload {
	return v.SetPolicyFragments
}

// GetSetEncassConfigs returns installBundleResponse.SetEncassConfigs, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetEncassConfigs() *installBundleSetEncassConfigsEncassConfigsPayload {
	return v.SetEncassConfigs
}

// GetSetGlobalPolicies returns installBundleResponse.SetGlobalPolicies, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetGlobalPolicies() *installBundleSetGlobalPoliciesGlobalPoliciesPayload {
	return v.SetGlobalPolicies
}

// GetSetBackgroundTaskPolicies returns installBundleResponse.SetBackgroundTaskPolicies, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetBackgroundTaskPolicies() *installBundleSetBackgroundTaskPoliciesBackgroundTaskPoliciesPayload {
	return v.SetBackgroundTaskPolicies
}

// GetSetServices returns installBundleResponse.SetServices, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetServices() *installBundleSetServicesL7ServicesPayload {
	return v.SetServices
}

// GetSetWebApiServices returns installBundleResponse.SetWebApiServices, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetWebApiServices() *installBundleSetWebApiServicesWebApiServicesPayload {
	return v.SetWebApiServices
}

// GetSetSoapServices returns installBundleResponse.SetSoapServices, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetSoapServices() *installBundleSetSoapServicesSoapServicesPayload {
	return v.SetSoapServices
}

// GetSetInternalWebApiServices returns installBundleResponse.SetInternalWebApiServices, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetInternalWebApiServices() *installBundleSetInternalWebApiServicesInternalWebApiServicesPayload {
	return v.SetInternalWebApiServices
}

// GetSetInternalSoapServices returns installBundleResponse.SetInternalSoapServices, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetInternalSoapServices() *installBundleSetInternalSoapServicesInternalSoapServicesPayload {
	return v.SetInternalSoapServices
}

// GetSetPolicyBackedIdps returns installBundleResponse.SetPolicyBackedIdps, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetPolicyBackedIdps() *installBundleSetPolicyBackedIdpsPolicyBackedIdpsPayload {
	return v.SetPolicyBackedIdps
}

// GetSetJmsDestinations returns installBundleResponse.SetJmsDestinations, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetJmsDestinations() *installBundleSetJmsDestinationsJmsDestinationsPayload {
	return v.SetJmsDestinations
}

// GetSetEmailListeners returns installBundleResponse.SetEmailListeners, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetEmailListeners() *installBundleSetEmailListenersEmailListenersPayload {
	return v.SetEmailListeners
}

// GetSetListenPorts returns installBundleResponse.SetListenPorts, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetListenPorts() *installBundleSetListenPortsListenPortsPayload {
	return v.SetListenPorts
}

// GetSetActiveConnectors returns installBundleResponse.SetActiveConnectors, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetActiveConnectors() *installBundleSetActiveConnectorsActiveConnectorsPayload {
	return v.SetActiveConnectors
}

// GetSetScheduledTasks returns installBundleResponse.SetScheduledTasks, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetScheduledTasks() *installBundleSetScheduledTasksScheduledTasksPayload {
	return v.SetScheduledTasks
}

// GetSetLogSinks returns installBundleResponse.SetLogSinks, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetLogSinks() *installBundleSetLogSinksLogSinksPayload {
	return v.SetLogSinks
}

// GetSetGenericEntities returns installBundleResponse.SetGenericEntities, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetGenericEntities() *installBundleSetGenericEntitiesGenericEntitiesPayload {
	return v.SetGenericEntities
}

// GetSetRoles returns installBundleResponse.SetRoles, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetRoles() *installBundleSetRolesRolesPayload { return v.SetRoles }

// GetSetAuditConfigurations returns installBundleResponse.SetAuditConfigurations, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetAuditConfigurations() *installBundleSetAuditConfigurationsAuditConfigurationsPayload {
	return v.SetAuditConfigurations
}

// GetSetKeys returns installBundleResponse.SetKeys, and is useful for accessing the field via an interface.
func (v *installBundleResponse) GetSetKeys() *installBundleSetKeysKeysPayload { return v.SetKeys }

// installBundleSetActiveConnectorsActiveConnectorsPayload includes the requested fields of the GraphQL type ActiveConnectorsPayload.
type installBundleSetActiveConnectorsActiveConnectorsPayload struct {
	DetailedStatus []*installBundleSetActiveConnectorsActiveConnectorsPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetActiveConnectorsActiveConnectorsPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetActiveConnectorsActiveConnectorsPayload) GetDetailedStatus() []*installBundleSetActiveConnectorsActiveConnectorsPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetActiveConnectorsActiveConnectorsPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetActiveConnectorsActiveConnectorsPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                                  `json:"action"`
	Status      EntityMutationStatus                                                                                                  `json:"status"`
	Description string                                                                                                                `json:"description"`
	Source      []*installBundleSetActiveConnectorsActiveConnectorsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetActiveConnectorsActiveConnectorsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetActiveConnectorsActiveConnectorsPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetActiveConnectorsActiveConnectorsPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetActiveConnectorsActiveConnectorsPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetActiveConnectorsActiveConnectorsPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetActiveConnectorsActiveConnectorsPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetActiveConnectorsActiveConnectorsPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetActiveConnectorsActiveConnectorsPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetActiveConnectorsActiveConnectorsPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetActiveConnectorsActiveConnectorsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetActiveConnectorsActiveConnectorsPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetActiveConnectorsActiveConnectorsPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetActiveConnectorsActiveConnectorsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetActiveConnectorsActiveConnectorsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetActiveConnectorsActiveConnectorsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetActiveConnectorsActiveConnectorsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetActiveConnectorsActiveConnectorsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetActiveConnectorsActiveConnectorsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetActiveConnectorsActiveConnectorsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetActiveConnectorsActiveConnectorsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetActiveConnectorsActiveConnectorsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetActiveConnectorsActiveConnectorsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetActiveConnectorsActiveConnectorsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetActiveConnectorsActiveConnectorsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetActiveConnectorsActiveConnectorsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetAdministrativeUserAccountPropertiesAdministrativeUserAccountPropertiesPayload includes the requested fields of the GraphQL type AdministrativeUserAccountPropertiesPayload.
type installBundleSetAdministrativeUserAccountPropertiesAdministrativeUserAccountPropertiesPayload struct {
	DetailedStatus []*installBundleSetAdministrativeUserAccountPropertiesAdministrativeUserAccountPropertiesPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetAdministrativeUserAccountPropertiesAdministrativeUserAccountPropertiesPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetAdministrativeUserAccountPropertiesAdministrativeUserAccountPropertiesPayload) GetDetailedStatus() []*installBundleSetAdministrativeUserAccountPropertiesAdministrativeUserAccountPropertiesPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetAdministrativeUserAccountPropertiesAdministrativeUserAccountPropertiesPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetAdministrativeUserAccountPropertiesAdministrativeUserAccountPropertiesPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                                                                        `json:"action"`
	Status      EntityMutationStatus                                                                                                                                        `json:"status"`
	Description string                                                                                                                                                      `json:"description"`
	Source      []*installBundleSetAdministrativeUserAccountPropertiesAdministrativeUserAccountPropertiesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetAdministrativeUserAccountPropertiesAdministrativeUserAccountPropertiesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetAdministrativeUserAccountPropertiesAdministrativeUserAccountPropertiesPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetAdministrativeUserAccountPropertiesAdministrativeUserAccountPropertiesPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetAdministrativeUserAccountPropertiesAdministrativeUserAccountPropertiesPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetAdministrativeUserAccountPropertiesAdministrativeUserAccountPropertiesPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetAdministrativeUserAccountPropertiesAdministrativeUserAccountPropertiesPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetAdministrativeUserAccountPropertiesAdministrativeUserAccountPropertiesPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetAdministrativeUserAccountPropertiesAdministrativeUserAccountPropertiesPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetAdministrativeUserAccountPropertiesAdministrativeUserAccountPropertiesPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetAdministrativeUserAccountPropertiesAdministrativeUserAccountPropertiesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetAdministrativeUserAccountPropertiesAdministrativeUserAccountPropertiesPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetAdministrativeUserAccountPropertiesAdministrativeUserAccountPropertiesPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetAdministrativeUserAccountPropertiesAdministrativeUserAccountPropertiesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetAdministrativeUserAccountPropertiesAdministrativeUserAccountPropertiesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetAdministrativeUserAccountPropertiesAdministrativeUserAccountPropertiesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetAdministrativeUserAccountPropertiesAdministrativeUserAccountPropertiesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetAdministrativeUserAccountPropertiesAdministrativeUserAccountPropertiesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetAdministrativeUserAccountPropertiesAdministrativeUserAccountPropertiesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetAdministrativeUserAccountPropertiesAdministrativeUserAccountPropertiesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetAdministrativeUserAccountPropertiesAdministrativeUserAccountPropertiesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetAdministrativeUserAccountPropertiesAdministrativeUserAccountPropertiesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetAdministrativeUserAccountPropertiesAdministrativeUserAccountPropertiesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetAdministrativeUserAccountPropertiesAdministrativeUserAccountPropertiesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetAdministrativeUserAccountPropertiesAdministrativeUserAccountPropertiesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetAdministrativeUserAccountPropertiesAdministrativeUserAccountPropertiesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetAuditConfigurationsAuditConfigurationsPayload includes the requested fields of the GraphQL type AuditConfigurationsPayload.
type installBundleSetAuditConfigurationsAuditConfigurationsPayload struct {
	DetailedStatus []*installBundleSetAuditConfigurationsAuditConfigurationsPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetAuditConfigurationsAuditConfigurationsPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetAuditConfigurationsAuditConfigurationsPayload) GetDetailedStatus() []*installBundleSetAuditConfigurationsAuditConfigurationsPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetAuditConfigurationsAuditConfigurationsPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetAuditConfigurationsAuditConfigurationsPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                                        `json:"action"`
	Status      EntityMutationStatus                                                                                                        `json:"status"`
	Description string                                                                                                                      `json:"description"`
	Source      []*installBundleSetAuditConfigurationsAuditConfigurationsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetAuditConfigurationsAuditConfigurationsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetAuditConfigurationsAuditConfigurationsPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetAuditConfigurationsAuditConfigurationsPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetAuditConfigurationsAuditConfigurationsPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetAuditConfigurationsAuditConfigurationsPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetAuditConfigurationsAuditConfigurationsPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetAuditConfigurationsAuditConfigurationsPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetAuditConfigurationsAuditConfigurationsPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetAuditConfigurationsAuditConfigurationsPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetAuditConfigurationsAuditConfigurationsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetAuditConfigurationsAuditConfigurationsPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetAuditConfigurationsAuditConfigurationsPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetAuditConfigurationsAuditConfigurationsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetAuditConfigurationsAuditConfigurationsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetAuditConfigurationsAuditConfigurationsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetAuditConfigurationsAuditConfigurationsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetAuditConfigurationsAuditConfigurationsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetAuditConfigurationsAuditConfigurationsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetAuditConfigurationsAuditConfigurationsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetAuditConfigurationsAuditConfigurationsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetAuditConfigurationsAuditConfigurationsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetAuditConfigurationsAuditConfigurationsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetAuditConfigurationsAuditConfigurationsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetAuditConfigurationsAuditConfigurationsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetAuditConfigurationsAuditConfigurationsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetBackgroundTaskPoliciesBackgroundTaskPoliciesPayload includes the requested fields of the GraphQL type BackgroundTaskPoliciesPayload.
type installBundleSetBackgroundTaskPoliciesBackgroundTaskPoliciesPayload struct {
	DetailedStatus []*installBundleSetBackgroundTaskPoliciesBackgroundTaskPoliciesPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetBackgroundTaskPoliciesBackgroundTaskPoliciesPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetBackgroundTaskPoliciesBackgroundTaskPoliciesPayload) GetDetailedStatus() []*installBundleSetBackgroundTaskPoliciesBackgroundTaskPoliciesPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetBackgroundTaskPoliciesBackgroundTaskPoliciesPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetBackgroundTaskPoliciesBackgroundTaskPoliciesPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                                              `json:"action"`
	Status      EntityMutationStatus                                                                                                              `json:"status"`
	Description string                                                                                                                            `json:"description"`
	Source      []*installBundleSetBackgroundTaskPoliciesBackgroundTaskPoliciesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetBackgroundTaskPoliciesBackgroundTaskPoliciesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetBackgroundTaskPoliciesBackgroundTaskPoliciesPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetBackgroundTaskPoliciesBackgroundTaskPoliciesPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetBackgroundTaskPoliciesBackgroundTaskPoliciesPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetBackgroundTaskPoliciesBackgroundTaskPoliciesPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetBackgroundTaskPoliciesBackgroundTaskPoliciesPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetBackgroundTaskPoliciesBackgroundTaskPoliciesPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetBackgroundTaskPoliciesBackgroundTaskPoliciesPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetBackgroundTaskPoliciesBackgroundTaskPoliciesPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetBackgroundTaskPoliciesBackgroundTaskPoliciesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetBackgroundTaskPoliciesBackgroundTaskPoliciesPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetBackgroundTaskPoliciesBackgroundTaskPoliciesPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetBackgroundTaskPoliciesBackgroundTaskPoliciesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetBackgroundTaskPoliciesBackgroundTaskPoliciesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetBackgroundTaskPoliciesBackgroundTaskPoliciesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetBackgroundTaskPoliciesBackgroundTaskPoliciesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetBackgroundTaskPoliciesBackgroundTaskPoliciesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetBackgroundTaskPoliciesBackgroundTaskPoliciesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetBackgroundTaskPoliciesBackgroundTaskPoliciesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetBackgroundTaskPoliciesBackgroundTaskPoliciesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetBackgroundTaskPoliciesBackgroundTaskPoliciesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetBackgroundTaskPoliciesBackgroundTaskPoliciesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetBackgroundTaskPoliciesBackgroundTaskPoliciesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetBackgroundTaskPoliciesBackgroundTaskPoliciesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetBackgroundTaskPoliciesBackgroundTaskPoliciesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetCassandraConnectionsCassandraConnectionsPayload includes the requested fields of the GraphQL type CassandraConnectionsPayload.
type installBundleSetCassandraConnectionsCassandraConnectionsPayload struct {
	DetailedStatus []*installBundleSetCassandraConnectionsCassandraConnectionsPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetCassandraConnectionsCassandraConnectionsPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetCassandraConnectionsCassandraConnectionsPayload) GetDetailedStatus() []*installBundleSetCassandraConnectionsCassandraConnectionsPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetCassandraConnectionsCassandraConnectionsPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetCassandraConnectionsCassandraConnectionsPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                                          `json:"action"`
	Status      EntityMutationStatus                                                                                                          `json:"status"`
	Description string                                                                                                                        `json:"description"`
	Source      []*installBundleSetCassandraConnectionsCassandraConnectionsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetCassandraConnectionsCassandraConnectionsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetCassandraConnectionsCassandraConnectionsPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetCassandraConnectionsCassandraConnectionsPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetCassandraConnectionsCassandraConnectionsPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetCassandraConnectionsCassandraConnectionsPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetCassandraConnectionsCassandraConnectionsPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetCassandraConnectionsCassandraConnectionsPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetCassandraConnectionsCassandraConnectionsPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetCassandraConnectionsCassandraConnectionsPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetCassandraConnectionsCassandraConnectionsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetCassandraConnectionsCassandraConnectionsPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetCassandraConnectionsCassandraConnectionsPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetCassandraConnectionsCassandraConnectionsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetCassandraConnectionsCassandraConnectionsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetCassandraConnectionsCassandraConnectionsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetCassandraConnectionsCassandraConnectionsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetCassandraConnectionsCassandraConnectionsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetCassandraConnectionsCassandraConnectionsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetCassandraConnectionsCassandraConnectionsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetCassandraConnectionsCassandraConnectionsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetCassandraConnectionsCassandraConnectionsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetCassandraConnectionsCassandraConnectionsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetCassandraConnectionsCassandraConnectionsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetCassandraConnectionsCassandraConnectionsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetCassandraConnectionsCassandraConnectionsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetClusterPropertiesClusterPropertiesPayload includes the requested fields of the GraphQL type ClusterPropertiesPayload.
type installBundleSetClusterPropertiesClusterPropertiesPayload struct {
	DetailedStatus []*installBundleSetClusterPropertiesClusterPropertiesPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetClusterPropertiesClusterPropertiesPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetClusterPropertiesClusterPropertiesPayload) GetDetailedStatus() []*installBundleSetClusterPropertiesClusterPropertiesPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetClusterPropertiesClusterPropertiesPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetClusterPropertiesClusterPropertiesPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                                    `json:"action"`
	Status      EntityMutationStatus                                                                                                    `json:"status"`
	Description string                                                                                                                  `json:"description"`
	Source      []*installBundleSetClusterPropertiesClusterPropertiesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetClusterPropertiesClusterPropertiesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetClusterPropertiesClusterPropertiesPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetClusterPropertiesClusterPropertiesPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetClusterPropertiesClusterPropertiesPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetClusterPropertiesClusterPropertiesPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetClusterPropertiesClusterPropertiesPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetClusterPropertiesClusterPropertiesPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetClusterPropertiesClusterPropertiesPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetClusterPropertiesClusterPropertiesPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetClusterPropertiesClusterPropertiesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetClusterPropertiesClusterPropertiesPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetClusterPropertiesClusterPropertiesPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetClusterPropertiesClusterPropertiesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetClusterPropertiesClusterPropertiesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetClusterPropertiesClusterPropertiesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetClusterPropertiesClusterPropertiesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetClusterPropertiesClusterPropertiesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetClusterPropertiesClusterPropertiesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetClusterPropertiesClusterPropertiesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetClusterPropertiesClusterPropertiesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetClusterPropertiesClusterPropertiesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetClusterPropertiesClusterPropertiesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetClusterPropertiesClusterPropertiesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetClusterPropertiesClusterPropertiesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetClusterPropertiesClusterPropertiesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetCustomKeyValuesCustomKeyValuePayload includes the requested fields of the GraphQL type CustomKeyValuePayload.
type installBundleSetCustomKeyValuesCustomKeyValuePayload struct {
	DetailedStatus []*installBundleSetCustomKeyValuesCustomKeyValuePayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetCustomKeyValuesCustomKeyValuePayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetCustomKeyValuesCustomKeyValuePayload) GetDetailedStatus() []*installBundleSetCustomKeyValuesCustomKeyValuePayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetCustomKeyValuesCustomKeyValuePayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetCustomKeyValuesCustomKeyValuePayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                               `json:"action"`
	Status      EntityMutationStatus                                                                                               `json:"status"`
	Description string                                                                                                             `json:"description"`
	Source      []*installBundleSetCustomKeyValuesCustomKeyValuePayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetCustomKeyValuesCustomKeyValuePayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetCustomKeyValuesCustomKeyValuePayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetCustomKeyValuesCustomKeyValuePayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetCustomKeyValuesCustomKeyValuePayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetCustomKeyValuesCustomKeyValuePayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetCustomKeyValuesCustomKeyValuePayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetCustomKeyValuesCustomKeyValuePayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetCustomKeyValuesCustomKeyValuePayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetCustomKeyValuesCustomKeyValuePayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetCustomKeyValuesCustomKeyValuePayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetCustomKeyValuesCustomKeyValuePayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetCustomKeyValuesCustomKeyValuePayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetCustomKeyValuesCustomKeyValuePayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetCustomKeyValuesCustomKeyValuePayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetCustomKeyValuesCustomKeyValuePayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetCustomKeyValuesCustomKeyValuePayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetCustomKeyValuesCustomKeyValuePayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetCustomKeyValuesCustomKeyValuePayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetCustomKeyValuesCustomKeyValuePayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetCustomKeyValuesCustomKeyValuePayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetCustomKeyValuesCustomKeyValuePayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetCustomKeyValuesCustomKeyValuePayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetCustomKeyValuesCustomKeyValuePayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetCustomKeyValuesCustomKeyValuePayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetCustomKeyValuesCustomKeyValuePayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetDtdsDtdsPayload includes the requested fields of the GraphQL type DtdsPayload.
type installBundleSetDtdsDtdsPayload struct {
	DetailedStatus []*installBundleSetDtdsDtdsPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetDtdsDtdsPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetDtdsDtdsPayload) GetDetailedStatus() []*installBundleSetDtdsDtdsPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetDtdsDtdsPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetDtdsDtdsPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                          `json:"action"`
	Status      EntityMutationStatus                                                                          `json:"status"`
	Description string                                                                                        `json:"description"`
	Source      []*installBundleSetDtdsDtdsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetDtdsDtdsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetDtdsDtdsPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetDtdsDtdsPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetDtdsDtdsPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetDtdsDtdsPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetDtdsDtdsPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetDtdsDtdsPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetDtdsDtdsPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetDtdsDtdsPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetDtdsDtdsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetDtdsDtdsPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetDtdsDtdsPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetDtdsDtdsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetDtdsDtdsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetDtdsDtdsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetDtdsDtdsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetDtdsDtdsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetDtdsDtdsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetDtdsDtdsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetDtdsDtdsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetDtdsDtdsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetDtdsDtdsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetDtdsDtdsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetDtdsDtdsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetDtdsDtdsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetEmailListenersEmailListenersPayload includes the requested fields of the GraphQL type EmailListenersPayload.
type installBundleSetEmailListenersEmailListenersPayload struct {
	DetailedStatus []*installBundleSetEmailListenersEmailListenersPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetEmailListenersEmailListenersPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetEmailListenersEmailListenersPayload) GetDetailedStatus() []*installBundleSetEmailListenersEmailListenersPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetEmailListenersEmailListenersPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetEmailListenersEmailListenersPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                              `json:"action"`
	Status      EntityMutationStatus                                                                                              `json:"status"`
	Description string                                                                                                            `json:"description"`
	Source      []*installBundleSetEmailListenersEmailListenersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetEmailListenersEmailListenersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetEmailListenersEmailListenersPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetEmailListenersEmailListenersPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetEmailListenersEmailListenersPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetEmailListenersEmailListenersPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetEmailListenersEmailListenersPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetEmailListenersEmailListenersPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetEmailListenersEmailListenersPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetEmailListenersEmailListenersPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetEmailListenersEmailListenersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetEmailListenersEmailListenersPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetEmailListenersEmailListenersPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetEmailListenersEmailListenersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetEmailListenersEmailListenersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetEmailListenersEmailListenersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetEmailListenersEmailListenersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetEmailListenersEmailListenersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetEmailListenersEmailListenersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetEmailListenersEmailListenersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetEmailListenersEmailListenersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetEmailListenersEmailListenersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetEmailListenersEmailListenersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetEmailListenersEmailListenersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetEmailListenersEmailListenersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetEmailListenersEmailListenersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetEncassConfigsEncassConfigsPayload includes the requested fields of the GraphQL type EncassConfigsPayload.
type installBundleSetEncassConfigsEncassConfigsPayload struct {
	DetailedStatus []*installBundleSetEncassConfigsEncassConfigsPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetEncassConfigsEncassConfigsPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetEncassConfigsEncassConfigsPayload) GetDetailedStatus() []*installBundleSetEncassConfigsEncassConfigsPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetEncassConfigsEncassConfigsPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetEncassConfigsEncassConfigsPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                            `json:"action"`
	Status      EntityMutationStatus                                                                                            `json:"status"`
	Description string                                                                                                          `json:"description"`
	Source      []*installBundleSetEncassConfigsEncassConfigsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetEncassConfigsEncassConfigsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetEncassConfigsEncassConfigsPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetEncassConfigsEncassConfigsPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetEncassConfigsEncassConfigsPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetEncassConfigsEncassConfigsPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetEncassConfigsEncassConfigsPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetEncassConfigsEncassConfigsPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetEncassConfigsEncassConfigsPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetEncassConfigsEncassConfigsPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetEncassConfigsEncassConfigsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetEncassConfigsEncassConfigsPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetEncassConfigsEncassConfigsPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetEncassConfigsEncassConfigsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetEncassConfigsEncassConfigsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetEncassConfigsEncassConfigsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetEncassConfigsEncassConfigsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetEncassConfigsEncassConfigsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetEncassConfigsEncassConfigsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetEncassConfigsEncassConfigsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetEncassConfigsEncassConfigsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetEncassConfigsEncassConfigsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetEncassConfigsEncassConfigsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetEncassConfigsEncassConfigsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetEncassConfigsEncassConfigsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetEncassConfigsEncassConfigsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetFederatedGroupsFederatedGroupsPayload includes the requested fields of the GraphQL type FederatedGroupsPayload.
type installBundleSetFederatedGroupsFederatedGroupsPayload struct {
	DetailedStatus []*installBundleSetFederatedGroupsFederatedGroupsPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetFederatedGroupsFederatedGroupsPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetFederatedGroupsFederatedGroupsPayload) GetDetailedStatus() []*installBundleSetFederatedGroupsFederatedGroupsPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetFederatedGroupsFederatedGroupsPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetFederatedGroupsFederatedGroupsPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                                `json:"action"`
	Status      EntityMutationStatus                                                                                                `json:"status"`
	Description string                                                                                                              `json:"description"`
	Source      []*installBundleSetFederatedGroupsFederatedGroupsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetFederatedGroupsFederatedGroupsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetFederatedGroupsFederatedGroupsPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetFederatedGroupsFederatedGroupsPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetFederatedGroupsFederatedGroupsPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetFederatedGroupsFederatedGroupsPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetFederatedGroupsFederatedGroupsPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetFederatedGroupsFederatedGroupsPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetFederatedGroupsFederatedGroupsPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetFederatedGroupsFederatedGroupsPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetFederatedGroupsFederatedGroupsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetFederatedGroupsFederatedGroupsPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetFederatedGroupsFederatedGroupsPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetFederatedGroupsFederatedGroupsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetFederatedGroupsFederatedGroupsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetFederatedGroupsFederatedGroupsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetFederatedGroupsFederatedGroupsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetFederatedGroupsFederatedGroupsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetFederatedGroupsFederatedGroupsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetFederatedGroupsFederatedGroupsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetFederatedGroupsFederatedGroupsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetFederatedGroupsFederatedGroupsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetFederatedGroupsFederatedGroupsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetFederatedGroupsFederatedGroupsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetFederatedGroupsFederatedGroupsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetFederatedGroupsFederatedGroupsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetFederatedIdpsFederatedIdpsPayload includes the requested fields of the GraphQL type FederatedIdpsPayload.
type installBundleSetFederatedIdpsFederatedIdpsPayload struct {
	DetailedStatus []*installBundleSetFederatedIdpsFederatedIdpsPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetFederatedIdpsFederatedIdpsPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetFederatedIdpsFederatedIdpsPayload) GetDetailedStatus() []*installBundleSetFederatedIdpsFederatedIdpsPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetFederatedIdpsFederatedIdpsPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetFederatedIdpsFederatedIdpsPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                            `json:"action"`
	Status      EntityMutationStatus                                                                                            `json:"status"`
	Description string                                                                                                          `json:"description"`
	Source      []*installBundleSetFederatedIdpsFederatedIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetFederatedIdpsFederatedIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetFederatedIdpsFederatedIdpsPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetFederatedIdpsFederatedIdpsPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetFederatedIdpsFederatedIdpsPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetFederatedIdpsFederatedIdpsPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetFederatedIdpsFederatedIdpsPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetFederatedIdpsFederatedIdpsPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetFederatedIdpsFederatedIdpsPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetFederatedIdpsFederatedIdpsPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetFederatedIdpsFederatedIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetFederatedIdpsFederatedIdpsPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetFederatedIdpsFederatedIdpsPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetFederatedIdpsFederatedIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetFederatedIdpsFederatedIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetFederatedIdpsFederatedIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetFederatedIdpsFederatedIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetFederatedIdpsFederatedIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetFederatedIdpsFederatedIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetFederatedIdpsFederatedIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetFederatedIdpsFederatedIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetFederatedIdpsFederatedIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetFederatedIdpsFederatedIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetFederatedIdpsFederatedIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetFederatedIdpsFederatedIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetFederatedIdpsFederatedIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetFederatedUsersFederatedUsersPayload includes the requested fields of the GraphQL type FederatedUsersPayload.
type installBundleSetFederatedUsersFederatedUsersPayload struct {
	DetailedStatus []*installBundleSetFederatedUsersFederatedUsersPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetFederatedUsersFederatedUsersPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetFederatedUsersFederatedUsersPayload) GetDetailedStatus() []*installBundleSetFederatedUsersFederatedUsersPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetFederatedUsersFederatedUsersPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetFederatedUsersFederatedUsersPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                              `json:"action"`
	Status      EntityMutationStatus                                                                                              `json:"status"`
	Description string                                                                                                            `json:"description"`
	Source      []*installBundleSetFederatedUsersFederatedUsersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetFederatedUsersFederatedUsersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetFederatedUsersFederatedUsersPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetFederatedUsersFederatedUsersPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetFederatedUsersFederatedUsersPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetFederatedUsersFederatedUsersPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetFederatedUsersFederatedUsersPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetFederatedUsersFederatedUsersPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetFederatedUsersFederatedUsersPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetFederatedUsersFederatedUsersPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetFederatedUsersFederatedUsersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetFederatedUsersFederatedUsersPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetFederatedUsersFederatedUsersPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetFederatedUsersFederatedUsersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetFederatedUsersFederatedUsersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetFederatedUsersFederatedUsersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetFederatedUsersFederatedUsersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetFederatedUsersFederatedUsersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetFederatedUsersFederatedUsersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetFederatedUsersFederatedUsersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetFederatedUsersFederatedUsersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetFederatedUsersFederatedUsersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetFederatedUsersFederatedUsersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetFederatedUsersFederatedUsersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetFederatedUsersFederatedUsersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetFederatedUsersFederatedUsersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetFipGroupsFipGroupsPayload includes the requested fields of the GraphQL type FipGroupsPayload.
type installBundleSetFipGroupsFipGroupsPayload struct {
	DetailedStatus []*installBundleSetFipGroupsFipGroupsPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetFipGroupsFipGroupsPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetFipGroupsFipGroupsPayload) GetDetailedStatus() []*installBundleSetFipGroupsFipGroupsPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetFipGroupsFipGroupsPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetFipGroupsFipGroupsPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                    `json:"action"`
	Status      EntityMutationStatus                                                                                    `json:"status"`
	Description string                                                                                                  `json:"description"`
	Source      []*installBundleSetFipGroupsFipGroupsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetFipGroupsFipGroupsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetFipGroupsFipGroupsPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetFipGroupsFipGroupsPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetFipGroupsFipGroupsPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetFipGroupsFipGroupsPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetFipGroupsFipGroupsPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetFipGroupsFipGroupsPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetFipGroupsFipGroupsPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetFipGroupsFipGroupsPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetFipGroupsFipGroupsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetFipGroupsFipGroupsPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetFipGroupsFipGroupsPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetFipGroupsFipGroupsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetFipGroupsFipGroupsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetFipGroupsFipGroupsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetFipGroupsFipGroupsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetFipGroupsFipGroupsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetFipGroupsFipGroupsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetFipGroupsFipGroupsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetFipGroupsFipGroupsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetFipGroupsFipGroupsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetFipGroupsFipGroupsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetFipGroupsFipGroupsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetFipGroupsFipGroupsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetFipGroupsFipGroupsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetFipUsersFipUsersPayload includes the requested fields of the GraphQL type FipUsersPayload.
type installBundleSetFipUsersFipUsersPayload struct {
	DetailedStatus []*installBundleSetFipUsersFipUsersPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetFipUsersFipUsersPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetFipUsersFipUsersPayload) GetDetailedStatus() []*installBundleSetFipUsersFipUsersPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetFipUsersFipUsersPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetFipUsersFipUsersPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                  `json:"action"`
	Status      EntityMutationStatus                                                                                  `json:"status"`
	Description string                                                                                                `json:"description"`
	Source      []*installBundleSetFipUsersFipUsersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetFipUsersFipUsersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetFipUsersFipUsersPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetFipUsersFipUsersPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetFipUsersFipUsersPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetFipUsersFipUsersPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetFipUsersFipUsersPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetFipUsersFipUsersPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetFipUsersFipUsersPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetFipUsersFipUsersPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetFipUsersFipUsersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetFipUsersFipUsersPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetFipUsersFipUsersPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetFipUsersFipUsersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetFipUsersFipUsersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetFipUsersFipUsersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetFipUsersFipUsersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetFipUsersFipUsersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetFipUsersFipUsersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetFipUsersFipUsersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetFipUsersFipUsersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetFipUsersFipUsersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetFipUsersFipUsersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetFipUsersFipUsersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetFipUsersFipUsersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetFipUsersFipUsersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetFipsFipsPayload includes the requested fields of the GraphQL type FipsPayload.
type installBundleSetFipsFipsPayload struct {
	DetailedStatus []*installBundleSetFipsFipsPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetFipsFipsPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetFipsFipsPayload) GetDetailedStatus() []*installBundleSetFipsFipsPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetFipsFipsPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetFipsFipsPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                          `json:"action"`
	Status      EntityMutationStatus                                                                          `json:"status"`
	Description string                                                                                        `json:"description"`
	Source      []*installBundleSetFipsFipsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetFipsFipsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetFipsFipsPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetFipsFipsPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetFipsFipsPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetFipsFipsPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetFipsFipsPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetFipsFipsPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetFipsFipsPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetFipsFipsPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetFipsFipsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetFipsFipsPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetFipsFipsPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetFipsFipsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetFipsFipsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetFipsFipsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetFipsFipsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetFipsFipsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetFipsFipsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetFipsFipsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetFipsFipsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetFipsFipsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetFipsFipsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetFipsFipsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetFipsFipsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetFipsFipsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetFoldersFoldersPayload includes the requested fields of the GraphQL type FoldersPayload.
type installBundleSetFoldersFoldersPayload struct {
	DetailedStatus []*installBundleSetFoldersFoldersPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetFoldersFoldersPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetFoldersFoldersPayload) GetDetailedStatus() []*installBundleSetFoldersFoldersPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetFoldersFoldersPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetFoldersFoldersPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                `json:"action"`
	Status      EntityMutationStatus                                                                                `json:"status"`
	Description string                                                                                              `json:"description"`
	Source      []*installBundleSetFoldersFoldersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetFoldersFoldersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetFoldersFoldersPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetFoldersFoldersPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetFoldersFoldersPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetFoldersFoldersPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetFoldersFoldersPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetFoldersFoldersPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetFoldersFoldersPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetFoldersFoldersPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetFoldersFoldersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetFoldersFoldersPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetFoldersFoldersPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetFoldersFoldersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetFoldersFoldersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetFoldersFoldersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetFoldersFoldersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetFoldersFoldersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetFoldersFoldersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetFoldersFoldersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetFoldersFoldersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetFoldersFoldersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetFoldersFoldersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetFoldersFoldersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetFoldersFoldersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetFoldersFoldersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetGenericEntitiesGenericEntitiesPayload includes the requested fields of the GraphQL type GenericEntitiesPayload.
type installBundleSetGenericEntitiesGenericEntitiesPayload struct {
	DetailedStatus []*installBundleSetGenericEntitiesGenericEntitiesPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetGenericEntitiesGenericEntitiesPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetGenericEntitiesGenericEntitiesPayload) GetDetailedStatus() []*installBundleSetGenericEntitiesGenericEntitiesPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetGenericEntitiesGenericEntitiesPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetGenericEntitiesGenericEntitiesPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                                `json:"action"`
	Status      EntityMutationStatus                                                                                                `json:"status"`
	Description string                                                                                                              `json:"description"`
	Source      []*installBundleSetGenericEntitiesGenericEntitiesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetGenericEntitiesGenericEntitiesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetGenericEntitiesGenericEntitiesPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetGenericEntitiesGenericEntitiesPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetGenericEntitiesGenericEntitiesPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetGenericEntitiesGenericEntitiesPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetGenericEntitiesGenericEntitiesPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetGenericEntitiesGenericEntitiesPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetGenericEntitiesGenericEntitiesPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetGenericEntitiesGenericEntitiesPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetGenericEntitiesGenericEntitiesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetGenericEntitiesGenericEntitiesPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetGenericEntitiesGenericEntitiesPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetGenericEntitiesGenericEntitiesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetGenericEntitiesGenericEntitiesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetGenericEntitiesGenericEntitiesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetGenericEntitiesGenericEntitiesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetGenericEntitiesGenericEntitiesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetGenericEntitiesGenericEntitiesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetGenericEntitiesGenericEntitiesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetGenericEntitiesGenericEntitiesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetGenericEntitiesGenericEntitiesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetGenericEntitiesGenericEntitiesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetGenericEntitiesGenericEntitiesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetGenericEntitiesGenericEntitiesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetGenericEntitiesGenericEntitiesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetGlobalPoliciesGlobalPoliciesPayload includes the requested fields of the GraphQL type GlobalPoliciesPayload.
type installBundleSetGlobalPoliciesGlobalPoliciesPayload struct {
	DetailedStatus []*installBundleSetGlobalPoliciesGlobalPoliciesPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetGlobalPoliciesGlobalPoliciesPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetGlobalPoliciesGlobalPoliciesPayload) GetDetailedStatus() []*installBundleSetGlobalPoliciesGlobalPoliciesPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetGlobalPoliciesGlobalPoliciesPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetGlobalPoliciesGlobalPoliciesPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                              `json:"action"`
	Status      EntityMutationStatus                                                                                              `json:"status"`
	Description string                                                                                                            `json:"description"`
	Source      []*installBundleSetGlobalPoliciesGlobalPoliciesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetGlobalPoliciesGlobalPoliciesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetGlobalPoliciesGlobalPoliciesPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetGlobalPoliciesGlobalPoliciesPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetGlobalPoliciesGlobalPoliciesPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetGlobalPoliciesGlobalPoliciesPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetGlobalPoliciesGlobalPoliciesPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetGlobalPoliciesGlobalPoliciesPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetGlobalPoliciesGlobalPoliciesPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetGlobalPoliciesGlobalPoliciesPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetGlobalPoliciesGlobalPoliciesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetGlobalPoliciesGlobalPoliciesPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetGlobalPoliciesGlobalPoliciesPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetGlobalPoliciesGlobalPoliciesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetGlobalPoliciesGlobalPoliciesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetGlobalPoliciesGlobalPoliciesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetGlobalPoliciesGlobalPoliciesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetGlobalPoliciesGlobalPoliciesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetGlobalPoliciesGlobalPoliciesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetGlobalPoliciesGlobalPoliciesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetGlobalPoliciesGlobalPoliciesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetGlobalPoliciesGlobalPoliciesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetGlobalPoliciesGlobalPoliciesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetGlobalPoliciesGlobalPoliciesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetGlobalPoliciesGlobalPoliciesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetGlobalPoliciesGlobalPoliciesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetHttpConfigurationsHttpConfigurationsPayload includes the requested fields of the GraphQL type HttpConfigurationsPayload.
type installBundleSetHttpConfigurationsHttpConfigurationsPayload struct {
	DetailedStatus []*installBundleSetHttpConfigurationsHttpConfigurationsPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetHttpConfigurationsHttpConfigurationsPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetHttpConfigurationsHttpConfigurationsPayload) GetDetailedStatus() []*installBundleSetHttpConfigurationsHttpConfigurationsPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetHttpConfigurationsHttpConfigurationsPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetHttpConfigurationsHttpConfigurationsPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                                      `json:"action"`
	Status      EntityMutationStatus                                                                                                      `json:"status"`
	Description string                                                                                                                    `json:"description"`
	Source      []*installBundleSetHttpConfigurationsHttpConfigurationsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetHttpConfigurationsHttpConfigurationsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetHttpConfigurationsHttpConfigurationsPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetHttpConfigurationsHttpConfigurationsPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetHttpConfigurationsHttpConfigurationsPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetHttpConfigurationsHttpConfigurationsPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetHttpConfigurationsHttpConfigurationsPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetHttpConfigurationsHttpConfigurationsPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetHttpConfigurationsHttpConfigurationsPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetHttpConfigurationsHttpConfigurationsPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetHttpConfigurationsHttpConfigurationsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetHttpConfigurationsHttpConfigurationsPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetHttpConfigurationsHttpConfigurationsPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetHttpConfigurationsHttpConfigurationsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetHttpConfigurationsHttpConfigurationsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetHttpConfigurationsHttpConfigurationsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetHttpConfigurationsHttpConfigurationsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetHttpConfigurationsHttpConfigurationsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetHttpConfigurationsHttpConfigurationsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetHttpConfigurationsHttpConfigurationsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetHttpConfigurationsHttpConfigurationsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetHttpConfigurationsHttpConfigurationsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetHttpConfigurationsHttpConfigurationsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetHttpConfigurationsHttpConfigurationsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetHttpConfigurationsHttpConfigurationsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetHttpConfigurationsHttpConfigurationsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetInternalGroupsInternalGroupsPayload includes the requested fields of the GraphQL type InternalGroupsPayload.
type installBundleSetInternalGroupsInternalGroupsPayload struct {
	DetailedStatus []*installBundleSetInternalGroupsInternalGroupsPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetInternalGroupsInternalGroupsPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalGroupsInternalGroupsPayload) GetDetailedStatus() []*installBundleSetInternalGroupsInternalGroupsPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetInternalGroupsInternalGroupsPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetInternalGroupsInternalGroupsPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                              `json:"action"`
	Status      EntityMutationStatus                                                                                              `json:"status"`
	Description string                                                                                                            `json:"description"`
	Source      []*installBundleSetInternalGroupsInternalGroupsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetInternalGroupsInternalGroupsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetInternalGroupsInternalGroupsPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalGroupsInternalGroupsPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetInternalGroupsInternalGroupsPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalGroupsInternalGroupsPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetInternalGroupsInternalGroupsPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalGroupsInternalGroupsPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetInternalGroupsInternalGroupsPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalGroupsInternalGroupsPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetInternalGroupsInternalGroupsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetInternalGroupsInternalGroupsPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalGroupsInternalGroupsPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetInternalGroupsInternalGroupsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetInternalGroupsInternalGroupsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetInternalGroupsInternalGroupsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetInternalGroupsInternalGroupsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalGroupsInternalGroupsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetInternalGroupsInternalGroupsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalGroupsInternalGroupsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetInternalGroupsInternalGroupsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetInternalGroupsInternalGroupsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetInternalGroupsInternalGroupsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalGroupsInternalGroupsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetInternalGroupsInternalGroupsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalGroupsInternalGroupsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetInternalIdpsInternalIdpsPayload includes the requested fields of the GraphQL type InternalIdpsPayload.
type installBundleSetInternalIdpsInternalIdpsPayload struct {
	DetailedStatus []*installBundleSetInternalIdpsInternalIdpsPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetInternalIdpsInternalIdpsPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalIdpsInternalIdpsPayload) GetDetailedStatus() []*installBundleSetInternalIdpsInternalIdpsPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetInternalIdpsInternalIdpsPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetInternalIdpsInternalIdpsPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                          `json:"action"`
	Status      EntityMutationStatus                                                                                          `json:"status"`
	Description string                                                                                                        `json:"description"`
	Source      []*installBundleSetInternalIdpsInternalIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetInternalIdpsInternalIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetInternalIdpsInternalIdpsPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalIdpsInternalIdpsPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetInternalIdpsInternalIdpsPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalIdpsInternalIdpsPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetInternalIdpsInternalIdpsPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalIdpsInternalIdpsPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetInternalIdpsInternalIdpsPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalIdpsInternalIdpsPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetInternalIdpsInternalIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetInternalIdpsInternalIdpsPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalIdpsInternalIdpsPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetInternalIdpsInternalIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetInternalIdpsInternalIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetInternalIdpsInternalIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetInternalIdpsInternalIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalIdpsInternalIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetInternalIdpsInternalIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalIdpsInternalIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetInternalIdpsInternalIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetInternalIdpsInternalIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetInternalIdpsInternalIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalIdpsInternalIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetInternalIdpsInternalIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalIdpsInternalIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetInternalSoapServicesInternalSoapServicesPayload includes the requested fields of the GraphQL type InternalSoapServicesPayload.
type installBundleSetInternalSoapServicesInternalSoapServicesPayload struct {
	DetailedStatus []*installBundleSetInternalSoapServicesInternalSoapServicesPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetInternalSoapServicesInternalSoapServicesPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalSoapServicesInternalSoapServicesPayload) GetDetailedStatus() []*installBundleSetInternalSoapServicesInternalSoapServicesPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetInternalSoapServicesInternalSoapServicesPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetInternalSoapServicesInternalSoapServicesPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                                          `json:"action"`
	Status      EntityMutationStatus                                                                                                          `json:"status"`
	Description string                                                                                                                        `json:"description"`
	Source      []*installBundleSetInternalSoapServicesInternalSoapServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetInternalSoapServicesInternalSoapServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetInternalSoapServicesInternalSoapServicesPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalSoapServicesInternalSoapServicesPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetInternalSoapServicesInternalSoapServicesPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalSoapServicesInternalSoapServicesPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetInternalSoapServicesInternalSoapServicesPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalSoapServicesInternalSoapServicesPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetInternalSoapServicesInternalSoapServicesPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalSoapServicesInternalSoapServicesPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetInternalSoapServicesInternalSoapServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetInternalSoapServicesInternalSoapServicesPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalSoapServicesInternalSoapServicesPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetInternalSoapServicesInternalSoapServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetInternalSoapServicesInternalSoapServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetInternalSoapServicesInternalSoapServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetInternalSoapServicesInternalSoapServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalSoapServicesInternalSoapServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetInternalSoapServicesInternalSoapServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalSoapServicesInternalSoapServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetInternalSoapServicesInternalSoapServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetInternalSoapServicesInternalSoapServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetInternalSoapServicesInternalSoapServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalSoapServicesInternalSoapServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetInternalSoapServicesInternalSoapServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalSoapServicesInternalSoapServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetInternalUsersInternalUsersPayload includes the requested fields of the GraphQL type InternalUsersPayload.
type installBundleSetInternalUsersInternalUsersPayload struct {
	DetailedStatus []*installBundleSetInternalUsersInternalUsersPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetInternalUsersInternalUsersPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalUsersInternalUsersPayload) GetDetailedStatus() []*installBundleSetInternalUsersInternalUsersPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetInternalUsersInternalUsersPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetInternalUsersInternalUsersPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                            `json:"action"`
	Status      EntityMutationStatus                                                                                            `json:"status"`
	Description string                                                                                                          `json:"description"`
	Source      []*installBundleSetInternalUsersInternalUsersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetInternalUsersInternalUsersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetInternalUsersInternalUsersPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalUsersInternalUsersPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetInternalUsersInternalUsersPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalUsersInternalUsersPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetInternalUsersInternalUsersPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalUsersInternalUsersPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetInternalUsersInternalUsersPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalUsersInternalUsersPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetInternalUsersInternalUsersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetInternalUsersInternalUsersPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalUsersInternalUsersPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetInternalUsersInternalUsersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetInternalUsersInternalUsersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetInternalUsersInternalUsersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetInternalUsersInternalUsersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalUsersInternalUsersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetInternalUsersInternalUsersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalUsersInternalUsersPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetInternalUsersInternalUsersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetInternalUsersInternalUsersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetInternalUsersInternalUsersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalUsersInternalUsersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetInternalUsersInternalUsersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalUsersInternalUsersPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetInternalWebApiServicesInternalWebApiServicesPayload includes the requested fields of the GraphQL type InternalWebApiServicesPayload.
type installBundleSetInternalWebApiServicesInternalWebApiServicesPayload struct {
	DetailedStatus []*installBundleSetInternalWebApiServicesInternalWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetInternalWebApiServicesInternalWebApiServicesPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalWebApiServicesInternalWebApiServicesPayload) GetDetailedStatus() []*installBundleSetInternalWebApiServicesInternalWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetInternalWebApiServicesInternalWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetInternalWebApiServicesInternalWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                                              `json:"action"`
	Status      EntityMutationStatus                                                                                                              `json:"status"`
	Description string                                                                                                                            `json:"description"`
	Source      []*installBundleSetInternalWebApiServicesInternalWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetInternalWebApiServicesInternalWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetInternalWebApiServicesInternalWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalWebApiServicesInternalWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetInternalWebApiServicesInternalWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalWebApiServicesInternalWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetInternalWebApiServicesInternalWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalWebApiServicesInternalWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetInternalWebApiServicesInternalWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalWebApiServicesInternalWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetInternalWebApiServicesInternalWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetInternalWebApiServicesInternalWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalWebApiServicesInternalWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetInternalWebApiServicesInternalWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetInternalWebApiServicesInternalWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetInternalWebApiServicesInternalWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetInternalWebApiServicesInternalWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalWebApiServicesInternalWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetInternalWebApiServicesInternalWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalWebApiServicesInternalWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetInternalWebApiServicesInternalWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetInternalWebApiServicesInternalWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetInternalWebApiServicesInternalWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalWebApiServicesInternalWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetInternalWebApiServicesInternalWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetInternalWebApiServicesInternalWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetJdbcConnectionsJdbcConnectionsPayload includes the requested fields of the GraphQL type JdbcConnectionsPayload.
type installBundleSetJdbcConnectionsJdbcConnectionsPayload struct {
	DetailedStatus []*installBundleSetJdbcConnectionsJdbcConnectionsPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetJdbcConnectionsJdbcConnectionsPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetJdbcConnectionsJdbcConnectionsPayload) GetDetailedStatus() []*installBundleSetJdbcConnectionsJdbcConnectionsPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetJdbcConnectionsJdbcConnectionsPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetJdbcConnectionsJdbcConnectionsPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                                `json:"action"`
	Status      EntityMutationStatus                                                                                                `json:"status"`
	Description string                                                                                                              `json:"description"`
	Source      []*installBundleSetJdbcConnectionsJdbcConnectionsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetJdbcConnectionsJdbcConnectionsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetJdbcConnectionsJdbcConnectionsPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetJdbcConnectionsJdbcConnectionsPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetJdbcConnectionsJdbcConnectionsPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetJdbcConnectionsJdbcConnectionsPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetJdbcConnectionsJdbcConnectionsPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetJdbcConnectionsJdbcConnectionsPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetJdbcConnectionsJdbcConnectionsPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetJdbcConnectionsJdbcConnectionsPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetJdbcConnectionsJdbcConnectionsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetJdbcConnectionsJdbcConnectionsPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetJdbcConnectionsJdbcConnectionsPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetJdbcConnectionsJdbcConnectionsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetJdbcConnectionsJdbcConnectionsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetJdbcConnectionsJdbcConnectionsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetJdbcConnectionsJdbcConnectionsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetJdbcConnectionsJdbcConnectionsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetJdbcConnectionsJdbcConnectionsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetJdbcConnectionsJdbcConnectionsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetJdbcConnectionsJdbcConnectionsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetJdbcConnectionsJdbcConnectionsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetJdbcConnectionsJdbcConnectionsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetJdbcConnectionsJdbcConnectionsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetJdbcConnectionsJdbcConnectionsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetJdbcConnectionsJdbcConnectionsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetJmsDestinationsJmsDestinationsPayload includes the requested fields of the GraphQL type JmsDestinationsPayload.
type installBundleSetJmsDestinationsJmsDestinationsPayload struct {
	DetailedStatus []*installBundleSetJmsDestinationsJmsDestinationsPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetJmsDestinationsJmsDestinationsPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetJmsDestinationsJmsDestinationsPayload) GetDetailedStatus() []*installBundleSetJmsDestinationsJmsDestinationsPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetJmsDestinationsJmsDestinationsPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetJmsDestinationsJmsDestinationsPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                                `json:"action"`
	Status      EntityMutationStatus                                                                                                `json:"status"`
	Description string                                                                                                              `json:"description"`
	Source      []*installBundleSetJmsDestinationsJmsDestinationsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetJmsDestinationsJmsDestinationsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetJmsDestinationsJmsDestinationsPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetJmsDestinationsJmsDestinationsPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetJmsDestinationsJmsDestinationsPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetJmsDestinationsJmsDestinationsPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetJmsDestinationsJmsDestinationsPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetJmsDestinationsJmsDestinationsPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetJmsDestinationsJmsDestinationsPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetJmsDestinationsJmsDestinationsPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetJmsDestinationsJmsDestinationsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetJmsDestinationsJmsDestinationsPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetJmsDestinationsJmsDestinationsPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetJmsDestinationsJmsDestinationsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetJmsDestinationsJmsDestinationsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetJmsDestinationsJmsDestinationsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetJmsDestinationsJmsDestinationsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetJmsDestinationsJmsDestinationsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetJmsDestinationsJmsDestinationsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetJmsDestinationsJmsDestinationsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetJmsDestinationsJmsDestinationsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetJmsDestinationsJmsDestinationsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetJmsDestinationsJmsDestinationsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetJmsDestinationsJmsDestinationsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetJmsDestinationsJmsDestinationsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetJmsDestinationsJmsDestinationsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetKeysKeysPayload includes the requested fields of the GraphQL type KeysPayload.
type installBundleSetKeysKeysPayload struct {
	DetailedStatus []*installBundleSetKeysKeysPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetKeysKeysPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetKeysKeysPayload) GetDetailedStatus() []*installBundleSetKeysKeysPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetKeysKeysPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetKeysKeysPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                          `json:"action"`
	Status      EntityMutationStatus                                                                          `json:"status"`
	Description string                                                                                        `json:"description"`
	Source      []*installBundleSetKeysKeysPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetKeysKeysPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetKeysKeysPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetKeysKeysPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetKeysKeysPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetKeysKeysPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetKeysKeysPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetKeysKeysPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetKeysKeysPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetKeysKeysPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetKeysKeysPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetKeysKeysPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetKeysKeysPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetKeysKeysPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetKeysKeysPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetKeysKeysPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetKeysKeysPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetKeysKeysPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetKeysKeysPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetKeysKeysPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetKeysKeysPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetKeysKeysPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetKeysKeysPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetKeysKeysPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetKeysKeysPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetKeysKeysPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetLdapIdpsLdapIdpsPayload includes the requested fields of the GraphQL type LdapIdpsPayload.
type installBundleSetLdapIdpsLdapIdpsPayload struct {
	DetailedStatus []*installBundleSetLdapIdpsLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetLdapIdpsLdapIdpsPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetLdapIdpsLdapIdpsPayload) GetDetailedStatus() []*installBundleSetLdapIdpsLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetLdapIdpsLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetLdapIdpsLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                  `json:"action"`
	Status      EntityMutationStatus                                                                                  `json:"status"`
	Description string                                                                                                `json:"description"`
	Source      []*installBundleSetLdapIdpsLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetLdapIdpsLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetLdapIdpsLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetLdapIdpsLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetLdapIdpsLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetLdapIdpsLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetLdapIdpsLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetLdapIdpsLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetLdapIdpsLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetLdapIdpsLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetLdapIdpsLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetLdapIdpsLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetLdapIdpsLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetLdapIdpsLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetLdapIdpsLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetLdapIdpsLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetLdapIdpsLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetLdapIdpsLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetLdapIdpsLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetLdapIdpsLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetLdapIdpsLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetLdapIdpsLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetLdapIdpsLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetLdapIdpsLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetLdapIdpsLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetLdapIdpsLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetLdapsLdapsPayload includes the requested fields of the GraphQL type LdapsPayload.
type installBundleSetLdapsLdapsPayload struct {
	DetailedStatus []*installBundleSetLdapsLdapsPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetLdapsLdapsPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetLdapsLdapsPayload) GetDetailedStatus() []*installBundleSetLdapsLdapsPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetLdapsLdapsPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetLdapsLdapsPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                            `json:"action"`
	Status      EntityMutationStatus                                                                            `json:"status"`
	Description string                                                                                          `json:"description"`
	Source      []*installBundleSetLdapsLdapsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetLdapsLdapsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetLdapsLdapsPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetLdapsLdapsPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetLdapsLdapsPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetLdapsLdapsPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetLdapsLdapsPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetLdapsLdapsPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetLdapsLdapsPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetLdapsLdapsPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetLdapsLdapsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetLdapsLdapsPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetLdapsLdapsPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetLdapsLdapsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetLdapsLdapsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetLdapsLdapsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetLdapsLdapsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetLdapsLdapsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetLdapsLdapsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetLdapsLdapsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetLdapsLdapsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetLdapsLdapsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetLdapsLdapsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetLdapsLdapsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetLdapsLdapsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetLdapsLdapsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetListenPortsListenPortsPayload includes the requested fields of the GraphQL type ListenPortsPayload.
type installBundleSetListenPortsListenPortsPayload struct {
	DetailedStatus []*installBundleSetListenPortsListenPortsPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetListenPortsListenPortsPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetListenPortsListenPortsPayload) GetDetailedStatus() []*installBundleSetListenPortsListenPortsPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetListenPortsListenPortsPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetListenPortsListenPortsPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                        `json:"action"`
	Status      EntityMutationStatus                                                                                        `json:"status"`
	Description string                                                                                                      `json:"description"`
	Source      []*installBundleSetListenPortsListenPortsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetListenPortsListenPortsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetListenPortsListenPortsPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetListenPortsListenPortsPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetListenPortsListenPortsPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetListenPortsListenPortsPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetListenPortsListenPortsPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetListenPortsListenPortsPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetListenPortsListenPortsPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetListenPortsListenPortsPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetListenPortsListenPortsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetListenPortsListenPortsPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetListenPortsListenPortsPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetListenPortsListenPortsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetListenPortsListenPortsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetListenPortsListenPortsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetListenPortsListenPortsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetListenPortsListenPortsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetListenPortsListenPortsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetListenPortsListenPortsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetListenPortsListenPortsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetListenPortsListenPortsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetListenPortsListenPortsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetListenPortsListenPortsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetListenPortsListenPortsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetListenPortsListenPortsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetLogSinksLogSinksPayload includes the requested fields of the GraphQL type LogSinksPayload.
type installBundleSetLogSinksLogSinksPayload struct {
	DetailedStatus []*installBundleSetLogSinksLogSinksPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetLogSinksLogSinksPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetLogSinksLogSinksPayload) GetDetailedStatus() []*installBundleSetLogSinksLogSinksPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetLogSinksLogSinksPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetLogSinksLogSinksPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                  `json:"action"`
	Status      EntityMutationStatus                                                                                  `json:"status"`
	Description string                                                                                                `json:"description"`
	Source      []*installBundleSetLogSinksLogSinksPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetLogSinksLogSinksPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetLogSinksLogSinksPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetLogSinksLogSinksPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetLogSinksLogSinksPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetLogSinksLogSinksPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetLogSinksLogSinksPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetLogSinksLogSinksPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetLogSinksLogSinksPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetLogSinksLogSinksPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetLogSinksLogSinksPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetLogSinksLogSinksPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetLogSinksLogSinksPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetLogSinksLogSinksPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetLogSinksLogSinksPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetLogSinksLogSinksPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetLogSinksLogSinksPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetLogSinksLogSinksPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetLogSinksLogSinksPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetLogSinksLogSinksPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetLogSinksLogSinksPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetLogSinksLogSinksPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetLogSinksLogSinksPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetLogSinksLogSinksPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetLogSinksLogSinksPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetLogSinksLogSinksPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetPasswordPoliciesPasswordPoliciesPayLoad includes the requested fields of the GraphQL type PasswordPoliciesPayLoad.
type installBundleSetPasswordPoliciesPasswordPoliciesPayLoad struct {
	DetailedStatus []*installBundleSetPasswordPoliciesPasswordPoliciesPayLoadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetPasswordPoliciesPasswordPoliciesPayLoad.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetPasswordPoliciesPasswordPoliciesPayLoad) GetDetailedStatus() []*installBundleSetPasswordPoliciesPasswordPoliciesPayLoadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetPasswordPoliciesPasswordPoliciesPayLoadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetPasswordPoliciesPasswordPoliciesPayLoadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                                  `json:"action"`
	Status      EntityMutationStatus                                                                                                  `json:"status"`
	Description string                                                                                                                `json:"description"`
	Source      []*installBundleSetPasswordPoliciesPasswordPoliciesPayLoadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetPasswordPoliciesPasswordPoliciesPayLoadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetPasswordPoliciesPasswordPoliciesPayLoadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetPasswordPoliciesPasswordPoliciesPayLoadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetPasswordPoliciesPasswordPoliciesPayLoadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetPasswordPoliciesPasswordPoliciesPayLoadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetPasswordPoliciesPasswordPoliciesPayLoadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetPasswordPoliciesPasswordPoliciesPayLoadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetPasswordPoliciesPasswordPoliciesPayLoadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetPasswordPoliciesPasswordPoliciesPayLoadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetPasswordPoliciesPasswordPoliciesPayLoadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetPasswordPoliciesPasswordPoliciesPayLoadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetPasswordPoliciesPasswordPoliciesPayLoadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetPasswordPoliciesPasswordPoliciesPayLoadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetPasswordPoliciesPasswordPoliciesPayLoadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetPasswordPoliciesPasswordPoliciesPayLoadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetPasswordPoliciesPasswordPoliciesPayLoadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetPasswordPoliciesPasswordPoliciesPayLoadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetPasswordPoliciesPasswordPoliciesPayLoadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetPasswordPoliciesPasswordPoliciesPayLoadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetPasswordPoliciesPasswordPoliciesPayLoadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetPasswordPoliciesPasswordPoliciesPayLoadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetPasswordPoliciesPasswordPoliciesPayLoadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetPasswordPoliciesPasswordPoliciesPayLoadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetPasswordPoliciesPasswordPoliciesPayLoadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetPasswordPoliciesPasswordPoliciesPayLoadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetPoliciesL7PoliciesPayload includes the requested fields of the GraphQL type L7PoliciesPayload.
type installBundleSetPoliciesL7PoliciesPayload struct {
	DetailedStatus []*installBundleSetPoliciesL7PoliciesPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetPoliciesL7PoliciesPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetPoliciesL7PoliciesPayload) GetDetailedStatus() []*installBundleSetPoliciesL7PoliciesPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetPoliciesL7PoliciesPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetPoliciesL7PoliciesPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                    `json:"action"`
	Status      EntityMutationStatus                                                                                    `json:"status"`
	Description string                                                                                                  `json:"description"`
	Source      []*installBundleSetPoliciesL7PoliciesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetPoliciesL7PoliciesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetPoliciesL7PoliciesPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetPoliciesL7PoliciesPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetPoliciesL7PoliciesPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetPoliciesL7PoliciesPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetPoliciesL7PoliciesPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetPoliciesL7PoliciesPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetPoliciesL7PoliciesPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetPoliciesL7PoliciesPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetPoliciesL7PoliciesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetPoliciesL7PoliciesPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetPoliciesL7PoliciesPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetPoliciesL7PoliciesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetPoliciesL7PoliciesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetPoliciesL7PoliciesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetPoliciesL7PoliciesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetPoliciesL7PoliciesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetPoliciesL7PoliciesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetPoliciesL7PoliciesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetPoliciesL7PoliciesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetPoliciesL7PoliciesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetPoliciesL7PoliciesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetPoliciesL7PoliciesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetPoliciesL7PoliciesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetPoliciesL7PoliciesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetPolicyBackedIdpsPolicyBackedIdpsPayload includes the requested fields of the GraphQL type PolicyBackedIdpsPayload.
type installBundleSetPolicyBackedIdpsPolicyBackedIdpsPayload struct {
	DetailedStatus []*installBundleSetPolicyBackedIdpsPolicyBackedIdpsPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetPolicyBackedIdpsPolicyBackedIdpsPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetPolicyBackedIdpsPolicyBackedIdpsPayload) GetDetailedStatus() []*installBundleSetPolicyBackedIdpsPolicyBackedIdpsPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetPolicyBackedIdpsPolicyBackedIdpsPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetPolicyBackedIdpsPolicyBackedIdpsPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                                  `json:"action"`
	Status      EntityMutationStatus                                                                                                  `json:"status"`
	Description string                                                                                                                `json:"description"`
	Source      []*installBundleSetPolicyBackedIdpsPolicyBackedIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetPolicyBackedIdpsPolicyBackedIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetPolicyBackedIdpsPolicyBackedIdpsPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetPolicyBackedIdpsPolicyBackedIdpsPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetPolicyBackedIdpsPolicyBackedIdpsPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetPolicyBackedIdpsPolicyBackedIdpsPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetPolicyBackedIdpsPolicyBackedIdpsPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetPolicyBackedIdpsPolicyBackedIdpsPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetPolicyBackedIdpsPolicyBackedIdpsPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetPolicyBackedIdpsPolicyBackedIdpsPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetPolicyBackedIdpsPolicyBackedIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetPolicyBackedIdpsPolicyBackedIdpsPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetPolicyBackedIdpsPolicyBackedIdpsPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetPolicyBackedIdpsPolicyBackedIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetPolicyBackedIdpsPolicyBackedIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetPolicyBackedIdpsPolicyBackedIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetPolicyBackedIdpsPolicyBackedIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetPolicyBackedIdpsPolicyBackedIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetPolicyBackedIdpsPolicyBackedIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetPolicyBackedIdpsPolicyBackedIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetPolicyBackedIdpsPolicyBackedIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetPolicyBackedIdpsPolicyBackedIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetPolicyBackedIdpsPolicyBackedIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetPolicyBackedIdpsPolicyBackedIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetPolicyBackedIdpsPolicyBackedIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetPolicyBackedIdpsPolicyBackedIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetPolicyFragmentsPolicyFragmentsPayload includes the requested fields of the GraphQL type PolicyFragmentsPayload.
type installBundleSetPolicyFragmentsPolicyFragmentsPayload struct {
	DetailedStatus []*installBundleSetPolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetPolicyFragmentsPolicyFragmentsPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetPolicyFragmentsPolicyFragmentsPayload) GetDetailedStatus() []*installBundleSetPolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetPolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetPolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                                `json:"action"`
	Status      EntityMutationStatus                                                                                                `json:"status"`
	Description string                                                                                                              `json:"description"`
	Source      []*installBundleSetPolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetPolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetPolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetPolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetPolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetPolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetPolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetPolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetPolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetPolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetPolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetPolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetPolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetPolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetPolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetPolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetPolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetPolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetPolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetPolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetPolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetPolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetPolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetPolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetPolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetPolicyFragmentsPolicyFragmentsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetRevocationCheckPoliciesRevocationCheckPoliciesPayload includes the requested fields of the GraphQL type RevocationCheckPoliciesPayload.
type installBundleSetRevocationCheckPoliciesRevocationCheckPoliciesPayload struct {
	DetailedStatus []*installBundleSetRevocationCheckPoliciesRevocationCheckPoliciesPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetRevocationCheckPoliciesRevocationCheckPoliciesPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetRevocationCheckPoliciesRevocationCheckPoliciesPayload) GetDetailedStatus() []*installBundleSetRevocationCheckPoliciesRevocationCheckPoliciesPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetRevocationCheckPoliciesRevocationCheckPoliciesPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetRevocationCheckPoliciesRevocationCheckPoliciesPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                                                `json:"action"`
	Status      EntityMutationStatus                                                                                                                `json:"status"`
	Description string                                                                                                                              `json:"description"`
	Source      []*installBundleSetRevocationCheckPoliciesRevocationCheckPoliciesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetRevocationCheckPoliciesRevocationCheckPoliciesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetRevocationCheckPoliciesRevocationCheckPoliciesPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetRevocationCheckPoliciesRevocationCheckPoliciesPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetRevocationCheckPoliciesRevocationCheckPoliciesPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetRevocationCheckPoliciesRevocationCheckPoliciesPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetRevocationCheckPoliciesRevocationCheckPoliciesPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetRevocationCheckPoliciesRevocationCheckPoliciesPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetRevocationCheckPoliciesRevocationCheckPoliciesPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetRevocationCheckPoliciesRevocationCheckPoliciesPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetRevocationCheckPoliciesRevocationCheckPoliciesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetRevocationCheckPoliciesRevocationCheckPoliciesPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetRevocationCheckPoliciesRevocationCheckPoliciesPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetRevocationCheckPoliciesRevocationCheckPoliciesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetRevocationCheckPoliciesRevocationCheckPoliciesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetRevocationCheckPoliciesRevocationCheckPoliciesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetRevocationCheckPoliciesRevocationCheckPoliciesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetRevocationCheckPoliciesRevocationCheckPoliciesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetRevocationCheckPoliciesRevocationCheckPoliciesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetRevocationCheckPoliciesRevocationCheckPoliciesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetRevocationCheckPoliciesRevocationCheckPoliciesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetRevocationCheckPoliciesRevocationCheckPoliciesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetRevocationCheckPoliciesRevocationCheckPoliciesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetRevocationCheckPoliciesRevocationCheckPoliciesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetRevocationCheckPoliciesRevocationCheckPoliciesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetRevocationCheckPoliciesRevocationCheckPoliciesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetRolesRolesPayload includes the requested fields of the GraphQL type RolesPayload.
type installBundleSetRolesRolesPayload struct {
	DetailedStatus []*installBundleSetRolesRolesPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetRolesRolesPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetRolesRolesPayload) GetDetailedStatus() []*installBundleSetRolesRolesPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetRolesRolesPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetRolesRolesPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                            `json:"action"`
	Status      EntityMutationStatus                                                                            `json:"status"`
	Description string                                                                                          `json:"description"`
	Source      []*installBundleSetRolesRolesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetRolesRolesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetRolesRolesPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetRolesRolesPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetRolesRolesPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetRolesRolesPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetRolesRolesPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetRolesRolesPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetRolesRolesPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetRolesRolesPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetRolesRolesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetRolesRolesPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetRolesRolesPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetRolesRolesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetRolesRolesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetRolesRolesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetRolesRolesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetRolesRolesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetRolesRolesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetRolesRolesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetRolesRolesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetRolesRolesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetRolesRolesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetRolesRolesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetRolesRolesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetRolesRolesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetSMConfigsSMConfigsPayload includes the requested fields of the GraphQL type SMConfigsPayload.
type installBundleSetSMConfigsSMConfigsPayload struct {
	DetailedStatus []*installBundleSetSMConfigsSMConfigsPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetSMConfigsSMConfigsPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetSMConfigsSMConfigsPayload) GetDetailedStatus() []*installBundleSetSMConfigsSMConfigsPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetSMConfigsSMConfigsPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetSMConfigsSMConfigsPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                    `json:"action"`
	Status      EntityMutationStatus                                                                                    `json:"status"`
	Description string                                                                                                  `json:"description"`
	Source      []*installBundleSetSMConfigsSMConfigsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetSMConfigsSMConfigsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetSMConfigsSMConfigsPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetSMConfigsSMConfigsPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetSMConfigsSMConfigsPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetSMConfigsSMConfigsPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetSMConfigsSMConfigsPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetSMConfigsSMConfigsPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetSMConfigsSMConfigsPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetSMConfigsSMConfigsPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetSMConfigsSMConfigsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetSMConfigsSMConfigsPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetSMConfigsSMConfigsPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetSMConfigsSMConfigsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetSMConfigsSMConfigsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetSMConfigsSMConfigsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetSMConfigsSMConfigsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetSMConfigsSMConfigsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetSMConfigsSMConfigsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetSMConfigsSMConfigsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetSMConfigsSMConfigsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetSMConfigsSMConfigsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetSMConfigsSMConfigsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetSMConfigsSMConfigsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetSMConfigsSMConfigsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetSMConfigsSMConfigsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetScheduledTasksScheduledTasksPayload includes the requested fields of the GraphQL type ScheduledTasksPayload.
type installBundleSetScheduledTasksScheduledTasksPayload struct {
	DetailedStatus []*installBundleSetScheduledTasksScheduledTasksPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetScheduledTasksScheduledTasksPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetScheduledTasksScheduledTasksPayload) GetDetailedStatus() []*installBundleSetScheduledTasksScheduledTasksPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetScheduledTasksScheduledTasksPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetScheduledTasksScheduledTasksPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                              `json:"action"`
	Status      EntityMutationStatus                                                                                              `json:"status"`
	Description string                                                                                                            `json:"description"`
	Source      []*installBundleSetScheduledTasksScheduledTasksPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetScheduledTasksScheduledTasksPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetScheduledTasksScheduledTasksPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetScheduledTasksScheduledTasksPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetScheduledTasksScheduledTasksPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetScheduledTasksScheduledTasksPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetScheduledTasksScheduledTasksPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetScheduledTasksScheduledTasksPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetScheduledTasksScheduledTasksPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetScheduledTasksScheduledTasksPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetScheduledTasksScheduledTasksPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetScheduledTasksScheduledTasksPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetScheduledTasksScheduledTasksPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetScheduledTasksScheduledTasksPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetScheduledTasksScheduledTasksPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetScheduledTasksScheduledTasksPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetScheduledTasksScheduledTasksPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetScheduledTasksScheduledTasksPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetScheduledTasksScheduledTasksPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetScheduledTasksScheduledTasksPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetScheduledTasksScheduledTasksPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetScheduledTasksScheduledTasksPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetScheduledTasksScheduledTasksPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetScheduledTasksScheduledTasksPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetScheduledTasksScheduledTasksPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetScheduledTasksScheduledTasksPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetSchemasSchemasPayload includes the requested fields of the GraphQL type SchemasPayload.
type installBundleSetSchemasSchemasPayload struct {
	DetailedStatus []*installBundleSetSchemasSchemasPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetSchemasSchemasPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetSchemasSchemasPayload) GetDetailedStatus() []*installBundleSetSchemasSchemasPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetSchemasSchemasPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetSchemasSchemasPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                `json:"action"`
	Status      EntityMutationStatus                                                                                `json:"status"`
	Description string                                                                                              `json:"description"`
	Source      []*installBundleSetSchemasSchemasPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetSchemasSchemasPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetSchemasSchemasPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetSchemasSchemasPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetSchemasSchemasPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetSchemasSchemasPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetSchemasSchemasPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetSchemasSchemasPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetSchemasSchemasPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetSchemasSchemasPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetSchemasSchemasPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetSchemasSchemasPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetSchemasSchemasPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetSchemasSchemasPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetSchemasSchemasPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetSchemasSchemasPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetSchemasSchemasPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetSchemasSchemasPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetSchemasSchemasPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetSchemasSchemasPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetSchemasSchemasPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetSchemasSchemasPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetSchemasSchemasPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetSchemasSchemasPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetSchemasSchemasPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetSchemasSchemasPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetSecretsSecretsPayload includes the requested fields of the GraphQL type SecretsPayload.
type installBundleSetSecretsSecretsPayload struct {
	DetailedStatus []*installBundleSetSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetSecretsSecretsPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetSecretsSecretsPayload) GetDetailedStatus() []*installBundleSetSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                `json:"action"`
	Status      EntityMutationStatus                                                                                `json:"status"`
	Description string                                                                                              `json:"description"`
	Source      []*installBundleSetSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetSecretsSecretsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetServerModuleFilesServerModuleFilesPayload includes the requested fields of the GraphQL type ServerModuleFilesPayload.
type installBundleSetServerModuleFilesServerModuleFilesPayload struct {
	DetailedStatus []*installBundleSetServerModuleFilesServerModuleFilesPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetServerModuleFilesServerModuleFilesPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetServerModuleFilesServerModuleFilesPayload) GetDetailedStatus() []*installBundleSetServerModuleFilesServerModuleFilesPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetServerModuleFilesServerModuleFilesPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetServerModuleFilesServerModuleFilesPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                                    `json:"action"`
	Status      EntityMutationStatus                                                                                                    `json:"status"`
	Description string                                                                                                                  `json:"description"`
	Source      []*installBundleSetServerModuleFilesServerModuleFilesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetServerModuleFilesServerModuleFilesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetServerModuleFilesServerModuleFilesPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetServerModuleFilesServerModuleFilesPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetServerModuleFilesServerModuleFilesPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetServerModuleFilesServerModuleFilesPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetServerModuleFilesServerModuleFilesPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetServerModuleFilesServerModuleFilesPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetServerModuleFilesServerModuleFilesPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetServerModuleFilesServerModuleFilesPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetServerModuleFilesServerModuleFilesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetServerModuleFilesServerModuleFilesPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetServerModuleFilesServerModuleFilesPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetServerModuleFilesServerModuleFilesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetServerModuleFilesServerModuleFilesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetServerModuleFilesServerModuleFilesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetServerModuleFilesServerModuleFilesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetServerModuleFilesServerModuleFilesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetServerModuleFilesServerModuleFilesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetServerModuleFilesServerModuleFilesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetServerModuleFilesServerModuleFilesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetServerModuleFilesServerModuleFilesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetServerModuleFilesServerModuleFilesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetServerModuleFilesServerModuleFilesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetServerModuleFilesServerModuleFilesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetServerModuleFilesServerModuleFilesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetServiceResolutionConfigsServiceResolutionConfigsPayLoad includes the requested fields of the GraphQL type ServiceResolutionConfigsPayLoad.
type installBundleSetServiceResolutionConfigsServiceResolutionConfigsPayLoad struct {
	DetailedStatus []*installBundleSetServiceResolutionConfigsServiceResolutionConfigsPayLoadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetServiceResolutionConfigsServiceResolutionConfigsPayLoad.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetServiceResolutionConfigsServiceResolutionConfigsPayLoad) GetDetailedStatus() []*installBundleSetServiceResolutionConfigsServiceResolutionConfigsPayLoadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetServiceResolutionConfigsServiceResolutionConfigsPayLoadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetServiceResolutionConfigsServiceResolutionConfigsPayLoadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                                                  `json:"action"`
	Status      EntityMutationStatus                                                                                                                  `json:"status"`
	Description string                                                                                                                                `json:"description"`
	Source      []*installBundleSetServiceResolutionConfigsServiceResolutionConfigsPayLoadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetServiceResolutionConfigsServiceResolutionConfigsPayLoadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetServiceResolutionConfigsServiceResolutionConfigsPayLoadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetServiceResolutionConfigsServiceResolutionConfigsPayLoadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetServiceResolutionConfigsServiceResolutionConfigsPayLoadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetServiceResolutionConfigsServiceResolutionConfigsPayLoadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetServiceResolutionConfigsServiceResolutionConfigsPayLoadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetServiceResolutionConfigsServiceResolutionConfigsPayLoadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetServiceResolutionConfigsServiceResolutionConfigsPayLoadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetServiceResolutionConfigsServiceResolutionConfigsPayLoadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetServiceResolutionConfigsServiceResolutionConfigsPayLoadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetServiceResolutionConfigsServiceResolutionConfigsPayLoadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetServiceResolutionConfigsServiceResolutionConfigsPayLoadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetServiceResolutionConfigsServiceResolutionConfigsPayLoadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetServiceResolutionConfigsServiceResolutionConfigsPayLoadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetServiceResolutionConfigsServiceResolutionConfigsPayLoadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetServiceResolutionConfigsServiceResolutionConfigsPayLoadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetServiceResolutionConfigsServiceResolutionConfigsPayLoadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetServiceResolutionConfigsServiceResolutionConfigsPayLoadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetServiceResolutionConfigsServiceResolutionConfigsPayLoadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetServiceResolutionConfigsServiceResolutionConfigsPayLoadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetServiceResolutionConfigsServiceResolutionConfigsPayLoadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetServiceResolutionConfigsServiceResolutionConfigsPayLoadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetServiceResolutionConfigsServiceResolutionConfigsPayLoadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetServiceResolutionConfigsServiceResolutionConfigsPayLoadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetServiceResolutionConfigsServiceResolutionConfigsPayLoadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetServicesL7ServicesPayload includes the requested fields of the GraphQL type L7ServicesPayload.
type installBundleSetServicesL7ServicesPayload struct {
	DetailedStatus []*installBundleSetServicesL7ServicesPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetServicesL7ServicesPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetServicesL7ServicesPayload) GetDetailedStatus() []*installBundleSetServicesL7ServicesPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetServicesL7ServicesPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetServicesL7ServicesPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                    `json:"action"`
	Status      EntityMutationStatus                                                                                    `json:"status"`
	Description string                                                                                                  `json:"description"`
	Source      []*installBundleSetServicesL7ServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetServicesL7ServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetServicesL7ServicesPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetServicesL7ServicesPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetServicesL7ServicesPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetServicesL7ServicesPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetServicesL7ServicesPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetServicesL7ServicesPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetServicesL7ServicesPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetServicesL7ServicesPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetServicesL7ServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetServicesL7ServicesPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetServicesL7ServicesPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetServicesL7ServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetServicesL7ServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetServicesL7ServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetServicesL7ServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetServicesL7ServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetServicesL7ServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetServicesL7ServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetServicesL7ServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetServicesL7ServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetServicesL7ServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetServicesL7ServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetServicesL7ServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetServicesL7ServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetSimpleLdapIdpsSimpleLdapIdpsPayload includes the requested fields of the GraphQL type SimpleLdapIdpsPayload.
type installBundleSetSimpleLdapIdpsSimpleLdapIdpsPayload struct {
	DetailedStatus []*installBundleSetSimpleLdapIdpsSimpleLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetSimpleLdapIdpsSimpleLdapIdpsPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetSimpleLdapIdpsSimpleLdapIdpsPayload) GetDetailedStatus() []*installBundleSetSimpleLdapIdpsSimpleLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetSimpleLdapIdpsSimpleLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetSimpleLdapIdpsSimpleLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                              `json:"action"`
	Status      EntityMutationStatus                                                                                              `json:"status"`
	Description string                                                                                                            `json:"description"`
	Source      []*installBundleSetSimpleLdapIdpsSimpleLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetSimpleLdapIdpsSimpleLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetSimpleLdapIdpsSimpleLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetSimpleLdapIdpsSimpleLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetSimpleLdapIdpsSimpleLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetSimpleLdapIdpsSimpleLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetSimpleLdapIdpsSimpleLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetSimpleLdapIdpsSimpleLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetSimpleLdapIdpsSimpleLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetSimpleLdapIdpsSimpleLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetSimpleLdapIdpsSimpleLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetSimpleLdapIdpsSimpleLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetSimpleLdapIdpsSimpleLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetSimpleLdapIdpsSimpleLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetSimpleLdapIdpsSimpleLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetSimpleLdapIdpsSimpleLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetSimpleLdapIdpsSimpleLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetSimpleLdapIdpsSimpleLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetSimpleLdapIdpsSimpleLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetSimpleLdapIdpsSimpleLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetSimpleLdapIdpsSimpleLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetSimpleLdapIdpsSimpleLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetSimpleLdapIdpsSimpleLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetSimpleLdapIdpsSimpleLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetSimpleLdapIdpsSimpleLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetSimpleLdapIdpsSimpleLdapIdpsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetSoapServicesSoapServicesPayload includes the requested fields of the GraphQL type SoapServicesPayload.
type installBundleSetSoapServicesSoapServicesPayload struct {
	DetailedStatus []*installBundleSetSoapServicesSoapServicesPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetSoapServicesSoapServicesPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetSoapServicesSoapServicesPayload) GetDetailedStatus() []*installBundleSetSoapServicesSoapServicesPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetSoapServicesSoapServicesPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetSoapServicesSoapServicesPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                          `json:"action"`
	Status      EntityMutationStatus                                                                                          `json:"status"`
	Description string                                                                                                        `json:"description"`
	Source      []*installBundleSetSoapServicesSoapServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetSoapServicesSoapServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetSoapServicesSoapServicesPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetSoapServicesSoapServicesPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetSoapServicesSoapServicesPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetSoapServicesSoapServicesPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetSoapServicesSoapServicesPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetSoapServicesSoapServicesPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetSoapServicesSoapServicesPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetSoapServicesSoapServicesPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetSoapServicesSoapServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetSoapServicesSoapServicesPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetSoapServicesSoapServicesPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetSoapServicesSoapServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetSoapServicesSoapServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetSoapServicesSoapServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetSoapServicesSoapServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetSoapServicesSoapServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetSoapServicesSoapServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetSoapServicesSoapServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetSoapServicesSoapServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetSoapServicesSoapServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetSoapServicesSoapServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetSoapServicesSoapServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetSoapServicesSoapServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetSoapServicesSoapServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetTrustedCertsTrustedCertsPayload includes the requested fields of the GraphQL type TrustedCertsPayload.
type installBundleSetTrustedCertsTrustedCertsPayload struct {
	DetailedStatus []*installBundleSetTrustedCertsTrustedCertsPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetTrustedCertsTrustedCertsPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetTrustedCertsTrustedCertsPayload) GetDetailedStatus() []*installBundleSetTrustedCertsTrustedCertsPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetTrustedCertsTrustedCertsPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetTrustedCertsTrustedCertsPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                          `json:"action"`
	Status      EntityMutationStatus                                                                                          `json:"status"`
	Description string                                                                                                        `json:"description"`
	Source      []*installBundleSetTrustedCertsTrustedCertsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetTrustedCertsTrustedCertsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetTrustedCertsTrustedCertsPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetTrustedCertsTrustedCertsPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetTrustedCertsTrustedCertsPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetTrustedCertsTrustedCertsPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetTrustedCertsTrustedCertsPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetTrustedCertsTrustedCertsPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetTrustedCertsTrustedCertsPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetTrustedCertsTrustedCertsPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetTrustedCertsTrustedCertsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetTrustedCertsTrustedCertsPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetTrustedCertsTrustedCertsPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetTrustedCertsTrustedCertsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetTrustedCertsTrustedCertsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetTrustedCertsTrustedCertsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetTrustedCertsTrustedCertsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetTrustedCertsTrustedCertsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetTrustedCertsTrustedCertsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetTrustedCertsTrustedCertsPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetTrustedCertsTrustedCertsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetTrustedCertsTrustedCertsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetTrustedCertsTrustedCertsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetTrustedCertsTrustedCertsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetTrustedCertsTrustedCertsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetTrustedCertsTrustedCertsPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetWebApiServicesWebApiServicesPayload includes the requested fields of the GraphQL type WebApiServicesPayload.
type installBundleSetWebApiServicesWebApiServicesPayload struct {
	DetailedStatus []*installBundleSetWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatus `json:"detailedStatus"`
}

// GetDetailedStatus returns installBundleSetWebApiServicesWebApiServicesPayload.DetailedStatus, and is useful for accessing the field via an interface.
func (v *installBundleSetWebApiServicesWebApiServicesPayload) GetDetailedStatus() []*installBundleSetWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatus {
	return v.DetailedStatus
}

// installBundleSetWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatus includes the requested fields of the GraphQL type EntityMutationDetailedStatus.
type installBundleSetWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatus struct {
	Action      EntityMutationAction                                                                                              `json:"action"`
	Status      EntityMutationStatus                                                                                              `json:"status"`
	Description string                                                                                                            `json:"description"`
	Source      []*installBundleSetWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty `json:"source"`
	Target      []*installBundleSetWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty `json:"target"`
}

// GetAction returns installBundleSetWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatus.Action, and is useful for accessing the field via an interface.
func (v *installBundleSetWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatus) GetAction() EntityMutationAction {
	return v.Action
}

// GetStatus returns installBundleSetWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatus.Status, and is useful for accessing the field via an interface.
func (v *installBundleSetWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatus) GetStatus() EntityMutationStatus {
	return v.Status
}

// GetDescription returns installBundleSetWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatus.Description, and is useful for accessing the field via an interface.
func (v *installBundleSetWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatus) GetDescription() string {
	return v.Description
}

// GetSource returns installBundleSetWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatus.Source, and is useful for accessing the field via an interface.
func (v *installBundleSetWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatus) GetSource() []*installBundleSetWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty {
	return v.Source
}

// GetTarget returns installBundleSetWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatus.Target, and is useful for accessing the field via an interface.
func (v *installBundleSetWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatus) GetTarget() []*installBundleSetWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty {
	return v.Target
}

// installBundleSetWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatusSourceAnyProperty) GetValue() interface{} {
	return v.Value
}

// installBundleSetWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty includes the requested fields of the GraphQL type AnyProperty.
type installBundleSetWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty struct {
	// The name of property
	Name string `json:"name"`
	// The value of the property
	Value interface{} `json:"value"`
}

// GetName returns installBundleSetWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Name, and is useful for accessing the field via an interface.
func (v *installBundleSetWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetName() string {
	return v.Name
}

// GetValue returns installBundleSetWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty.Value, and is useful for accessing the field via an interface.
func (v *installBundleSetWebApiServicesWebApiServicesPayloadDetailedStatusEntityMutationDetailedStatusTargetAnyProperty) GetValue() interface{} {
	return v.Value
}

// The query or mutation executed by deleteKeys.
const deleteKeys_Operation = `
mutation deleteKeys ($keys: [String!]!) {
	deleteKeys(aliases: $keys) {
		detailedStatus {
			status
			description
		}
		keys {
			goid
			keystoreId
			alias
		}
	}
}
`

func deleteKeys(
	ctx_ context.Context,
	client_ graphql.Client,
	keys []string,
) (*deleteKeysResponse, error) {
	req_ := &graphql.Request{
		OpName: "deleteKeys",
		Query:  deleteKeys_Operation,
		Variables: &__deleteKeysInput{
			Keys: keys,
		},
	}
	var err_ error

	var data_ deleteKeysResponse
	resp_ := &graphql.Response{Data: &data_}

	err_ = client_.MakeRequest(
		ctx_,
		req_,
		resp_,
	)

	return &data_, err_
}

// The query or mutation executed by deleteL7PortalApi.
const deleteL7PortalApi_Operation = `
mutation deleteL7PortalApi ($webApiServiceResolutionPaths: [String!]!, $policyFragmentNames: [String!]!) {
	deleteWebApiServices(resolutionPaths: $webApiServiceResolutionPaths) {
		detailedStatus {
			status
			description
		}
	}
	deletePolicyFragments(names: $policyFragmentNames) {
		detailedStatus {
			status
			description
		}
	}
}
`

func deleteL7PortalApi(
	ctx_ context.Context,
	client_ graphql.Client,
	webApiServiceResolutionPaths []string,
	policyFragmentNames []string,
) (*deleteL7PortalApiResponse, error) {
	req_ := &graphql.Request{
		OpName: "deleteL7PortalApi",
		Query:  deleteL7PortalApi_Operation,
		Variables: &__deleteL7PortalApiInput{
			WebApiServiceResolutionPaths: webApiServiceResolutionPaths,
			PolicyFragmentNames:          policyFragmentNames,
		},
	}
	var err_ error

	var data_ deleteL7PortalApiResponse
	resp_ := &graphql.Response{Data: &data_}

	err_ = client_.MakeRequest(
		ctx_,
		req_,
		resp_,
	)

	return &data_, err_
}

// The query or mutation executed by deleteSecrets.
const deleteSecrets_Operation = `
mutation deleteSecrets ($secrets: [String!]!) {
	deleteSecrets(names: $secrets) {
		detailedStatus {
			status
			description
		}
		secrets {
			goid
			name
		}
	}
}
`

func deleteSecrets(
	ctx_ context.Context,
	client_ graphql.Client,
	secrets []string,
) (*deleteSecretsResponse, error) {
	req_ := &graphql.Request{
		OpName: "deleteSecrets",
		Query:  deleteSecrets_Operation,
		Variables: &__deleteSecretsInput{
			Secrets: secrets,
		},
	}
	var err_ error

	var data_ deleteSecretsResponse
	resp_ := &graphql.Response{Data: &data_}

	err_ = client_.MakeRequest(
		ctx_,
		req_,
		resp_,
	)

	return &data_, err_
}

// The query or mutation executed by installBundle.
const installBundle_Operation = `
mutation installBundle ($activeConnectors: [ActiveConnectorInput!]! = [], $administrativeUserAccountProperties: [AdministrativeUserAccountPropertyInput!]! = [], $backgroundTaskPolicies: [BackgroundTaskPolicyInput!]! = [], $cassandraConnections: [CassandraConnectionInput!]! = [], $clusterProperties: [ClusterPropertyInput!]! = [], $dtds: [DtdInput!]! = [], $emailListeners: [EmailListenerInput!]! = [], $encassConfigs: [EncassConfigInput!]! = [], $fipGroups: [FipGroupInput!]! = [], $fipUsers: [FipUserInput!]! = [], $fips: [FipInput!]! = [], $federatedGroups: [FederatedGroupInput!]! = [], $federatedUsers: [FederatedUserInput!]! = [], $internalIdps: [InternalIdpInput!] = [], $federatedIdps: [FederatedIdpInput!]! = [], $ldapIdps: [LdapIdpInput!] = [], $simpleLdapIdps: [SimpleLdapIdpInput!] = [], $policyBackedIdps: [PolicyBackedIdpInput!] = [], $globalPolicies: [GlobalPolicyInput!]! = [], $internalGroups: [InternalGroupInput!]! = [], $internalSoapServices: [SoapServiceInput!]! = [], $internalUsers: [InternalUserInput!]! = [], $internalWebApiServices: [WebApiServiceInput!]! = [], $jdbcConnections: [JdbcConnectionInput!]! = [], $jmsDestinations: [JmsDestinationInput!]! = [], $keys: [KeyInput!]! = [], $ldaps: [LdapInput!]! = [], $roles: [RoleInput!]! = [], $listenPorts: [ListenPortInput!]! = [], $passwordPolicies: [PasswordPolicyInput!]! = [], $policies: [L7PolicyInput!]! = [], $policyFragments: [PolicyFragmentInput!]! = [], $revocationCheckPolicies: [RevocationCheckPolicyInput!]! = [], $scheduledTasks: [ScheduledTaskInput!]! = [], $logSinks: [LogSinkInput!]! = [], $schemas: [SchemaInput!]! = [], $secrets: [SecretInput!]! = [], $httpConfigurations: [HttpConfigurationInput!]! = [], $customKeyValues: [CustomKeyValueInput!]! = [], $serverModuleFiles: [ServerModuleFileInput!]! = [], $serviceResolutionConfigs: [ServiceResolutionConfigInput!]! = [], $folders: [FolderInput!]! = [], $smConfigs: [SMConfigInput!]! = [], $services: [L7ServiceInput!]! = [], $soapServices: [SoapServiceInput!]! = [], $trustedCerts: [TrustedCertInput!]! = [], $webApiServices: [WebApiServiceInput!]! = [], $genericEntities: [GenericEntityInput!]! = [], $auditConfigurations: [AuditConfigurationInput!]! = []) {
	setServerModuleFiles(input: $serverModuleFiles) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setClusterProperties(input: $clusterProperties) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setServiceResolutionConfigs(input: $serviceResolutionConfigs) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setPasswordPolicies(input: $passwordPolicies) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setAdministrativeUserAccountProperties(input: $administrativeUserAccountProperties) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setFolders(input: $folders) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setRevocationCheckPolicies(input: $revocationCheckPolicies) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setTrustedCerts(input: $trustedCerts) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setSecrets(input: $secrets) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setHttpConfigurations(input: $httpConfigurations) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setCustomKeyValues(input: $customKeyValues) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setSchemas(input: $schemas) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setDtds(input: $dtds) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setJdbcConnections(input: $jdbcConnections) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setInternalIdps(input: $internalIdps) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setFederatedIdps(input: $federatedIdps) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setLdapIdps(input: $ldapIdps) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setSimpleLdapIdps(input: $simpleLdapIdps) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setFips(input: $fips) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setLdaps(input: $ldaps) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setFederatedGroups(input: $federatedGroups) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setFipGroups(input: $fipGroups) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setInternalGroups(input: $internalGroups) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setFederatedUsers(input: $federatedUsers) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setFipUsers(input: $fipUsers) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setInternalUsers(input: $internalUsers) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setCassandraConnections(input: $cassandraConnections) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setSMConfigs(input: $smConfigs) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setPolicies(input: $policies) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setPolicyFragments(input: $policyFragments) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setEncassConfigs(input: $encassConfigs) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setGlobalPolicies(input: $globalPolicies) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setBackgroundTaskPolicies(input: $backgroundTaskPolicies) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setServices(input: $services) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setWebApiServices(input: $webApiServices) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setSoapServices(input: $soapServices) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setInternalWebApiServices(input: $internalWebApiServices) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setInternalSoapServices(input: $internalSoapServices) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setPolicyBackedIdps(input: $policyBackedIdps) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setJmsDestinations(input: $jmsDestinations) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setEmailListeners(input: $emailListeners) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setListenPorts(input: $listenPorts) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setActiveConnectors(input: $activeConnectors) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setScheduledTasks(input: $scheduledTasks) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setLogSinks(input: $logSinks) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setGenericEntities(input: $genericEntities) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setRoles(input: $roles) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setAuditConfigurations(input: $auditConfigurations) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
	setKeys(input: $keys) {
		detailedStatus {
			action
			status
			description
			source {
				name
				value
			}
			target {
				name
				value
			}
		}
	}
}
`

func installBundle(
	ctx_ context.Context,
	client_ graphql.Client,
	activeConnectors []*ActiveConnectorInput,
	administrativeUserAccountProperties []*AdministrativeUserAccountPropertyInput,
	backgroundTaskPolicies []*BackgroundTaskPolicyInput,
	cassandraConnections []*CassandraConnectionInput,
	clusterProperties []*ClusterPropertyInput,
	dtds []*DtdInput,
	emailListeners []*EmailListenerInput,
	encassConfigs []*EncassConfigInput,
	fipGroups []*FipGroupInput,
	fipUsers []*FipUserInput,
	fips []*FipInput,
	federatedGroups []*FederatedGroupInput,
	federatedUsers []*FederatedUserInput,
	internalIdps []*InternalIdpInput,
	federatedIdps []*FederatedIdpInput,
	ldapIdps []*LdapIdpInput,
	simpleLdapIdps []*SimpleLdapIdpInput,
	policyBackedIdps []*PolicyBackedIdpInput,
	globalPolicies []*GlobalPolicyInput,
	internalGroups []*InternalGroupInput,
	internalSoapServices []*SoapServiceInput,
	internalUsers []*InternalUserInput,
	internalWebApiServices []*WebApiServiceInput,
	jdbcConnections []*JdbcConnectionInput,
	jmsDestinations []*JmsDestinationInput,
	keys []*KeyInput,
	ldaps []*LdapInput,
	roles []*RoleInput,
	listenPorts []*ListenPortInput,
	passwordPolicies []*PasswordPolicyInput,
	policies []*L7PolicyInput,
	policyFragments []*PolicyFragmentInput,
	revocationCheckPolicies []*RevocationCheckPolicyInput,
	scheduledTasks []*ScheduledTaskInput,
	logSinks []*LogSinkInput,
	schemas []*SchemaInput,
	secrets []*SecretInput,
	httpConfigurations []*HttpConfigurationInput,
	customKeyValues []*CustomKeyValueInput,
	serverModuleFiles []*ServerModuleFileInput,
	serviceResolutionConfigs []*ServiceResolutionConfigInput,
	folders []*FolderInput,
	smConfigs []*SMConfigInput,
	services []*L7ServiceInput,
	soapServices []*SoapServiceInput,
	trustedCerts []*TrustedCertInput,
	webApiServices []*WebApiServiceInput,
	genericEntities []*GenericEntityInput,
	auditConfigurations []*AuditConfigurationInput,
) (*installBundleResponse, error) {
	req_ := &graphql.Request{
		OpName: "installBundle",
		Query:  installBundle_Operation,
		Variables: &__installBundleInput{
			ActiveConnectors:                    activeConnectors,
			AdministrativeUserAccountProperties: administrativeUserAccountProperties,
			BackgroundTaskPolicies:              backgroundTaskPolicies,
			CassandraConnections:                cassandraConnections,
			ClusterProperties:                   clusterProperties,
			Dtds:                                dtds,
			EmailListeners:                      emailListeners,
			EncassConfigs:                       encassConfigs,
			FipGroups:                           fipGroups,
			FipUsers:                            fipUsers,
			Fips:                                fips,
			FederatedGroups:                     federatedGroups,
			FederatedUsers:                      federatedUsers,
			InternalIdps:                        internalIdps,
			FederatedIdps:                       federatedIdps,
			LdapIdps:                            ldapIdps,
			SimpleLdapIdps:                      simpleLdapIdps,
			PolicyBackedIdps:                    policyBackedIdps,
			GlobalPolicies:                      globalPolicies,
			InternalGroups:                      internalGroups,
			InternalSoapServices:                internalSoapServices,
			InternalUsers:                       internalUsers,
			InternalWebApiServices:              internalWebApiServices,
			JdbcConnections:                     jdbcConnections,
			JmsDestinations:                     jmsDestinations,
			Keys:                                keys,
			Ldaps:                               ldaps,
			Roles:                               roles,
			ListenPorts:                         listenPorts,
			PasswordPolicies:                    passwordPolicies,
			Policies:                            policies,
			PolicyFragments:                     policyFragments,
			RevocationCheckPolicies:             revocationCheckPolicies,
			ScheduledTasks:                      scheduledTasks,
			LogSinks:                            logSinks,
			Schemas:                             schemas,
			Secrets:                             secrets,
			HttpConfigurations:                  httpConfigurations,
			CustomKeyValues:                     customKeyValues,
			ServerModuleFiles:                   serverModuleFiles,
			ServiceResolutionConfigs:            serviceResolutionConfigs,
			Folders:                             folders,
			SmConfigs:                           smConfigs,
			Services:                            services,
			SoapServices:                        soapServices,
			TrustedCerts:                        trustedCerts,
			WebApiServices:                      webApiServices,
			GenericEntities:                     genericEntities,
			AuditConfigurations:                 auditConfigurations,
		},
	}
	var err_ error

	var data_ installBundleResponse
	resp_ := &graphql.Response{Data: &data_}

	err_ = client_.MakeRequest(
		ctx_,
		req_,
		resp_,
	)

	return &data_, err_
}

// The query or mutation executed by installBundleGeneric.
const installBundleGeneric_Operation = `
mutation installBundleGeneric {
	installBundleEntities {
		summary
	}
}
`

func installBundleGeneric(
	ctx_ context.Context,
	client_ graphql.Client,
) (*installBundleGenericResponse, error) {
	req_ := &graphql.Request{
		OpName: "installBundleGeneric",
		Query:  installBundleGeneric_Operation,
	}
	var err_ error

	var data_ installBundleGenericResponse
	resp_ := &graphql.Response{Data: &data_}

	err_ = client_.MakeRequest(
		ctx_,
		req_,
		resp_,
	)

	return &data_, err_
}
