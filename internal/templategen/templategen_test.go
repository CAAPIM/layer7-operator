package templategen

import (
	b64 "encoding/base64"
	"testing"
	"time"
)

func TestBuildTemplate(t *testing.T) {

	template := PolicyTemplate{
		Uuid:                       "72093738-871a-45bd-b114-ad3a61893ac0",
		ApiPolicyTemplateArguments: []PolicyTemplateArg{{"ptName", "ptValue"}},
	}

	b64LocationUrl := b64.StdEncoding.EncodeToString([]byte("https://localhost:9443/stubbed"))
	b64SsgUrl := b64.StdEncoding.EncodeToString([]byte("bookings"))

	portalApi := PortalAPI{
		TenantId:        "T1",
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

	iterations := 20000

	t.Run("Template Test", func(t *testing.T) {
		if got := BuildTemplate(portalApi); got == "" {
			t.Errorf("BuildTemplate() = %v, want %v", got, "xml string")
		}
	})

	t.Run("Performance Test", func(t *testing.T) {
		duration := perfTest(portalApi, iterations)
		timeout := time.Duration(750 * time.Millisecond)

		if min(duration, timeout) == timeout {
			t.Errorf("perfTest(portalApi,iterations) = %v, want %v", duration, timeout)
		}
	})
}

func perfTest(portalApi PortalAPI, iterations int) time.Duration {
	start := time.Now()
	for a := 0; a < iterations; a++ {
		BuildTemplate(portalApi)
	}
	duration := time.Since(start)

	return duration
}
