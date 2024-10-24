package templategen

type PortalAPI struct {
	TenantId     string `json:"tenantId"`
	Uuid         string `json:"apiUuid"`
	UuidStripped string `json:"apiId"`     // templates have cases where uuid is stripped of -
	ServiceId    string `json:"serviceId"` // Portal calculates this as UUID.nameUUIDFromBytes(api.getUuid().toString().getBytes());
	Name         string `json:"name"`
	//Description       string           `json:"description"`
	//Type              string           `json:"type"`
	//PortalStatus      string           `json:"portalStatus"`
	ApiEnabled bool `json:"enabled"` // Added as a String to make it easier to use in qtpl
	//AccessStatus      string           `json:"accessStatus"`
	SsgUrl       string `json:"ssgUrl"`
	SsgUrlBase64 string `json:"ssgUrlEncoded"` // added as Fragment wants B64 version
	LocationUrl  string `json:"locationUrl"`   // backend full URL
	//Version      string `json:"version"`
	//ApiEulaUuid       string           `json:"apiEulaUuid"`
	PublishedTs    int    `json:"publishedTs"`
	CreateTs       int    `json:"createTs"`
	ModifyTs       int    `json:"modifyTs"`
	SsgServiceType string `json:"ssgServiceType"`
	//ApplicationUsage  int              `json:"applicationUsage"`
	//Tags              []string         `json:"tags"`
	PolicyTemplates []PolicyTemplate `json:"policyEntities"`              // required by qtpl templates
	CustomFields    []CustomField    `json:"customFieldValues,omitempty"` // required by qtpl templates
	//PublishedByPortal bool             `json:"publishedByPortal"`
	Checksum string `json:"checksum"`
}

type PolicyTemplate struct {
	Uuid                       string              `json:"policyEntityUuid"`
	ApiPolicyTemplateArguments []PolicyTemplateArg `json:"policyTemplateArguments"`
}

type PolicyTemplateArg struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type CustomField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

//go:generate go get -u github.com/valyala/quicktemplate/qtc
//go:generate qtc -dir=
func BuildTemplate(portalApi PortalAPI) string {
	fragment := FromApiFragmentTemplate(portalApi)
	service := FromApiServiceTemplate(portalApi)
	restman := FromRestmamTemplate(portalApi, service, fragment, "false", "")
	return restman
}
