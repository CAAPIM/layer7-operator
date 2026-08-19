// Copyright (c) 2026 Broadcom Inc. and its subsidiaries. All Rights Reserved.
// AI assistance has been used to generate some or all contents of this file. That includes, but is not limited to, new code, modifying existing code, stylistic edits.

package templategen

import (
	b64 "encoding/base64"
	"slices"
	"strings"
	"testing"
	"time"
)

// allVerbs is the hardcoded fallback list emitted when HttpMethods is unset for a REST API.
var allVerbs = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"}

// buildRestman renders the restman template with an explicit isSoapApi value. BuildTemplate
// hardcodes isSoapApi to "false", so the SOAP branch of the <l7:Verbs> conditional is only
// reachable this way.
func buildRestman(portalApi PortalAPI, isSoapApi string) string {
	fragment := FromApiFragmentTemplate(portalApi)
	service := FromApiServiceTemplate(portalApi)
	return FromRestmamTemplate(portalApi, service, fragment, isSoapApi, "")
}

// assertVerbs checks the rendered output contains exactly the expected verbs and nothing else.
func assertVerbs(t *testing.T, got string, want []string, context string) {
	t.Helper()

	if count := strings.Count(got, "<l7:Verb>"); count != len(want) {
		t.Errorf("%s produced %d <l7:Verb> entries, want %d", context, count, len(want))
	}

	for _, verb := range want {
		if !strings.Contains(got, "<l7:Verb>"+verb+"</l7:Verb>") {
			t.Errorf("%s did not contain <l7:Verb>%s</l7:Verb>", context, verb)
		}
	}

	for _, verb := range allVerbs {
		if slices.Contains(want, verb) {
			continue
		}
		if strings.Contains(got, "<l7:Verb>"+verb+"</l7:Verb>") {
			t.Errorf("%s unexpectedly contained <l7:Verb>%s</l7:Verb>", context, verb)
		}
	}
}

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
		timeout := time.Duration(3000 * time.Millisecond)

		if min(duration, timeout) == timeout {
			t.Errorf("perfTest(portalApi,iterations) = %v, want %v", duration, timeout)
		}
	})

	t.Run("HttpMethods Test", func(t *testing.T) {
		restrictedApi := portalApi
		restrictedApi.HttpMethods = []string{"GET"}

		got := BuildTemplate(restrictedApi)

		if got == "" {
			t.Errorf("BuildTemplate() = %v, want %v", got, "xml string")
		}

		if strings.Count(got, "<l7:Verb>") != 1 {
			t.Errorf("BuildTemplate() with HttpMethods=[GET] produced %d <l7:Verb> entries, want 1",
				strings.Count(got, "<l7:Verb>"))
		}

		if !strings.Contains(got, "<l7:Verb>GET</l7:Verb>") {
			t.Errorf("BuildTemplate() with HttpMethods=[GET] did not contain <l7:Verb>GET</l7:Verb>")
		}

		for _, verb := range []string{"POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"} {
			if strings.Contains(got, "<l7:Verb>"+verb+"</l7:Verb>") {
				t.Errorf("BuildTemplate() with HttpMethods=[GET] unexpectedly contained <l7:Verb>%s</l7:Verb>", verb)
			}
		}
	})

	// Backward-compatibility guard: every pre-existing Portal API leaves HttpMethods unset and
	// must keep getting the full 7-verb list. Without this, a change to the shared <l7:Verbs>
	// conditional could silently restrict every existing API.
	t.Run("HttpMethods Unset Falls Back To All Verbs", func(t *testing.T) {
		assertVerbs(t, BuildTemplate(portalApi), allVerbs, "BuildTemplate() with HttpMethods unset")
	})

	// An explicitly empty list must behave identically to unset. This matters because Portal
	// serializes PortalMeta with @JsonInclude(NON_NULL), so an API with no methods selected
	// ships "httpMethods": [] (not an omitted field) in its L7Api CR. Both paths rely on the
	// template's len(...) > 0 check, so lock that equivalence in.
	t.Run("HttpMethods Empty Falls Back To All Verbs", func(t *testing.T) {
		emptyApi := portalApi
		emptyApi.HttpMethods = []string{}

		assertVerbs(t, BuildTemplate(emptyApi), allVerbs, "BuildTemplate() with HttpMethods=[]")
	})

	t.Run("HttpMethods Multiple Verbs", func(t *testing.T) {
		restrictedApi := portalApi
		restrictedApi.HttpMethods = []string{"GET", "POST"}

		assertVerbs(t, BuildTemplate(restrictedApi), []string{"GET", "POST"},
			"BuildTemplate() with HttpMethods=[GET,POST]")
	})

	// OTHER is a Gateway-native value in Graphman's HttpMethod enum and is deliberately not
	// filtered out before emission, so it must survive into the bundle verbatim.
	t.Run("HttpMethods Passes Through OTHER", func(t *testing.T) {
		otherApi := portalApi
		otherApi.HttpMethods = []string{"OTHER"}

		got := BuildTemplate(otherApi)

		if count := strings.Count(got, "<l7:Verb>"); count != 1 {
			t.Errorf("BuildTemplate() with HttpMethods=[OTHER] produced %d <l7:Verb> entries, want 1", count)
		}

		if !strings.Contains(got, "<l7:Verb>OTHER</l7:Verb>") {
			t.Errorf("BuildTemplate() with HttpMethods=[OTHER] did not contain <l7:Verb>OTHER</l7:Verb>")
		}
	})

	// The SOAP branch is currently unreachable in production (BuildTemplate hardcodes
	// isSoapApi to "false"), but the <l7:Verbs> conditional still routes through it, so pin
	// its behaviour to catch an accidental regression of the shared conditional.
	t.Run("Soap Unset Falls Back To Get And Post", func(t *testing.T) {
		assertVerbs(t, buildRestman(portalApi, "true"), []string{"GET", "POST"},
			"FromRestmamTemplate() with isSoapApi=true and HttpMethods unset")
	})

	// Branch ordering is intentional: HttpMethods takes priority over the SOAP fallback.
	t.Run("HttpMethods Overrides Soap Fallback", func(t *testing.T) {
		restrictedApi := portalApi
		restrictedApi.HttpMethods = []string{"GET"}

		assertVerbs(t, buildRestman(restrictedApi, "true"), []string{"GET"},
			"FromRestmamTemplate() with isSoapApi=true and HttpMethods=[GET]")
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
