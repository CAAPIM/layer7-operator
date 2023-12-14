package templategen

import (
	b64 "encoding/base64"
	"fmt"
	"strings"
	"time"
	"unicode"
)

type PortalAPI struct {
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

	// Fragment: func FromApiFragmentTemplate(portalApi v1.PortalAPI)
	// Service: func FromApiServiceTemplate(portalApi v1.PortalAPI)
	// Restman: func FromRestmamTemplate(portalApi v1.PortalAPI, apiServiceXml string, apiFragmentXml string, isSoapApi string, wsdl string)

	fragment := FromApiFragmentTemplate(portalApi)
	service := FromApiServiceTemplate(portalApi)
	restman := FromRestmamTemplate(portalApi, service, fragment, "false", "")

	//sanitizedString := strings.ReplaceAll(restman, "\n", "")
	//sanitizedString = strings.Trim(sanitizedString, " ")
	//sanitizedString = stripSpaces(sanitizedString)
	//fmt.Println(sanitizedString)

	return restman
}

func test() {

	template := PolicyTemplate{
		Uuid:                       "72093738-871a-45bd-b114-ad3a61893ac0",
		ApiPolicyTemplateArguments: []PolicyTemplateArg{{"ptName", "ptValue"}},
	}

	b64LocationUrl := b64.StdEncoding.EncodeToString([]byte("https://localhost:9443/stubbed"))
	b64SsgUrl := b64.StdEncoding.EncodeToString([]byte("bookings"))

	portalApi := PortalAPI{
		Name:            "Booking",
		Uuid:            "17b0fb67-03d3-4340-ae68-b489e1835075",
		UuidStripped:    "17b0fb6703d34340ae68b489e1835075", // calculated
		ServiceId:       "0dd8af1599c43b74a7acb743aa3b3836", // calculated
		SsgUrl:          "bookings",                         // used in service
		SsgUrlBase64:    b64SsgUrl,                          // used in fragment
		LocationUrl:     b64LocationUrl,
		ApiEnabled:      true,
		CustomFields:    []CustomField{{Name: "Custom Field 1", Value: "three"}},
		PolicyTemplates: []PolicyTemplate{template, {Uuid: "92092f24-6ca1-3f19-b29e-70287c64a369"}},
		ModifyTs:        1694490707365,
	}

	restman := BuildTemplate(portalApi)
	perfTest(portalApi, 20000)
	fmt.Printf("%s\n", restman)
}

func perfTest(portalApi PortalAPI, iterations int) time.Duration {
	start := time.Now()
	for a := 0; a < iterations; a++ {
		BuildTemplate(portalApi)
	}
	duration := time.Since(start)
	fmt.Printf("Performance Results: %s millis for %d iterations\n", duration, iterations)

	return duration
}

func stripSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			// if the character is a space, drop it
			return -1
		}
		// else keep it in the string
		return r
	}, str)
}
