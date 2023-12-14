// Code generated by qtc from "portal-api-restman-template.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

// ----------------------------------------------
//
//
// Variables Used:
// - portalApi.Name
// - portalApi.UuidStripped
// - portalApi.Uuid
// - portalApi.ServiceId
// - portalApi.SsgUrl
// - portalApi.ApiEnabled
// - portalApi.ModifyTs
// - portalApi.CustomFields
// - apiServiceXml
// - apiFragmentXml
// - isSoapAPI
// - wsdl (empty string for REST)
//
// -------------------------------------------------
//

//line portal-api-restman-template.qtpl:20
package templategen

//line portal-api-restman-template.qtpl:20
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line portal-api-restman-template.qtpl:20
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line portal-api-restman-template.qtpl:20
func StreamFromRestmamTemplate(qw422016 *qt422016.Writer, portalApi PortalAPI, apiServiceXml string, apiFragmentXml string, isSoapApi string, wsdl string) {
//line portal-api-restman-template.qtpl:20
	qw422016.N().S(`
<l7:Bundle xmlns:l7="http://ns.l7tech.com/2010/04/gateway-management">
  <l7:References>
    <l7:Item xmlns:l7="http://ns.l7tech.com/2010/04/gateway-management">
      <l7:Name>`)
//line portal-api-restman-template.qtpl:24
	qw422016.E().S(portalApi.Name)
//line portal-api-restman-template.qtpl:24
	qw422016.N().S(`-fragment</l7:Name>
      <l7:Id>`)
//line portal-api-restman-template.qtpl:25
	qw422016.E().S(portalApi.UuidStripped)
//line portal-api-restman-template.qtpl:25
	qw422016.N().S(`</l7:Id>
      <l7:Type>POLICY</l7:Type>
      <l7:Resource>
        <l7:Policy guid="`)
//line portal-api-restman-template.qtpl:28
	qw422016.E().S(portalApi.Uuid)
//line portal-api-restman-template.qtpl:28
	qw422016.N().S(`" id="`)
//line portal-api-restman-template.qtpl:28
	qw422016.E().S(portalApi.UuidStripped)
//line portal-api-restman-template.qtpl:28
	qw422016.N().S(`" version="0">
          <l7:PolicyDetail folderId="ddb84c6f397d7dbd3cca71d3043f019c" guid="`)
//line portal-api-restman-template.qtpl:29
	qw422016.E().S(portalApi.Uuid)
//line portal-api-restman-template.qtpl:29
	qw422016.N().S(`" id="`)
//line portal-api-restman-template.qtpl:29
	qw422016.E().S(portalApi.UuidStripped)
//line portal-api-restman-template.qtpl:29
	qw422016.N().S(`">
            <l7:Name>`)
//line portal-api-restman-template.qtpl:30
	qw422016.E().S(portalApi.Name)
//line portal-api-restman-template.qtpl:30
	qw422016.N().S(`-fragment</l7:Name>
            <l7:PolicyType>Include</l7:PolicyType>
            <l7:Properties>
              <l7:Property key="revision">
                <l7:LongValue>1</l7:LongValue>
              </l7:Property>
              <l7:Property key="soap">
                <l7:BooleanValue>`)
//line portal-api-restman-template.qtpl:37
	qw422016.E().S(isSoapApi)
//line portal-api-restman-template.qtpl:37
	qw422016.N().S(`</l7:BooleanValue>
              </l7:Property>
            </l7:Properties>
          </l7:PolicyDetail>
          <l7:Resources>
            <l7:ResourceSet tag="policy">
              <l7:Resource type="policy">`)
//line portal-api-restman-template.qtpl:43
	qw422016.E().S(apiFragmentXml)
//line portal-api-restman-template.qtpl:43
	qw422016.N().S(`</l7:Resource>
            </l7:ResourceSet>
          </l7:Resources>
        </l7:Policy>
      </l7:Resource>
    </l7:Item>
    <l7:Item xmlns:l7="http://ns.l7tech.com/2010/04/gateway-management">
      <l7:Name>`)
//line portal-api-restman-template.qtpl:50
	qw422016.E().S(portalApi.Name)
//line portal-api-restman-template.qtpl:50
	qw422016.N().S(`</l7:Name>
      <l7:Id>`)
//line portal-api-restman-template.qtpl:51
	qw422016.E().S(portalApi.ServiceId)
//line portal-api-restman-template.qtpl:51
	qw422016.N().S(`</l7:Id>
      <l7:Type>SERVICE</l7:Type>
      <l7:Resource>
        <l7:Service id="`)
//line portal-api-restman-template.qtpl:54
	qw422016.E().S(portalApi.ServiceId)
//line portal-api-restman-template.qtpl:54
	qw422016.N().S(`" xmlns:l7="http://ns.l7tech.com/2010/04/gateway-management">
          <l7:ServiceDetail folderId="ddb84c6f397d7dbd3cca71d3043f019c" id="`)
//line portal-api-restman-template.qtpl:55
	qw422016.E().S(portalApi.ServiceId)
//line portal-api-restman-template.qtpl:55
	qw422016.N().S(`">
            <l7:Name>`)
//line portal-api-restman-template.qtpl:56
	qw422016.E().S(portalApi.Name)
//line portal-api-restman-template.qtpl:56
	qw422016.N().S(`</l7:Name>
            <l7:Enabled>true</l7:Enabled>
            <l7:ServiceMappings>
              <l7:HttpMapping>
                <l7:UrlPattern>/`)
//line portal-api-restman-template.qtpl:60
	qw422016.E().S(portalApi.SsgUrl)
//line portal-api-restman-template.qtpl:60
	qw422016.N().S(`*</l7:UrlPattern>
                <l7:Verbs>
                  `)
//line portal-api-restman-template.qtpl:62
	if isSoapApi == "true" {
//line portal-api-restman-template.qtpl:62
		qw422016.N().S(`
                    <l7:Verb>GET</l7:Verb>
                    <l7:Verb>POST</l7:Verb>
                  `)
//line portal-api-restman-template.qtpl:65
	} else {
//line portal-api-restman-template.qtpl:65
		qw422016.N().S(`
                    <l7:Verb>GET</l7:Verb>
                    <l7:Verb>POST</l7:Verb>
                    <l7:Verb>PUT</l7:Verb>
                    <l7:Verb>DELETE</l7:Verb>
                    <l7:Verb>OPTIONS</l7:Verb>
                    <l7:Verb>PATCH</l7:Verb>
                    <l7:Verb>HEAD</l7:Verb>
                  `)
//line portal-api-restman-template.qtpl:73
	}
//line portal-api-restman-template.qtpl:73
	qw422016.N().S(`
                </l7:Verbs>
              </l7:HttpMapping>
            </l7:ServiceMappings>
            <l7:Properties>
              <l7:Property key="internal">
                <l7:BooleanValue>false</l7:BooleanValue>
              </l7:Property>
              <l7:Property key="soap">
                <l7:BooleanValue>`)
//line portal-api-restman-template.qtpl:82
	qw422016.E().S(isSoapApi)
//line portal-api-restman-template.qtpl:82
	qw422016.N().S(`</l7:BooleanValue>
              </l7:Property>
              <l7:Property key="tracingEnabled">
                <l7:BooleanValue>false</l7:BooleanValue>
              </l7:Property>
              <l7:Property key="wssProcessingEnabled">
                <l7:BooleanValue>false</l7:BooleanValue>
              </l7:Property>
              <l7:Property key="property.portalID">
                <l7:StringValue>`)
//line portal-api-restman-template.qtpl:91
	qw422016.E().S(portalApi.Uuid)
//line portal-api-restman-template.qtpl:91
	qw422016.N().S(`</l7:StringValue>
              </l7:Property>
              <l7:Property key="property.internal.portalAPIEnabled">
                <l7:StringValue>`)
//line portal-api-restman-template.qtpl:94
	qw422016.E().V(portalApi.ApiEnabled)
//line portal-api-restman-template.qtpl:94
	qw422016.N().S(`</l7:StringValue>
              </l7:Property>
              <l7:Property key="property.portalModifyTS">
                <l7:StringValue>`)
//line portal-api-restman-template.qtpl:97
	qw422016.N().D(portalApi.ModifyTs)
//line portal-api-restman-template.qtpl:97
	qw422016.N().S(`</l7:StringValue>
              </l7:Property>
              `)
//line portal-api-restman-template.qtpl:99
	for _, customField := range portalApi.CustomFields {
//line portal-api-restman-template.qtpl:99
		qw422016.N().S(`
              <l7:Property key="property.`)
//line portal-api-restman-template.qtpl:100
		qw422016.E().S(customField.Name)
//line portal-api-restman-template.qtpl:100
		qw422016.N().S(`" xmlns:l7="http://ns.l7tech.com/2010/04/gateway-management">
                  <l7:StringValue>`)
//line portal-api-restman-template.qtpl:101
		qw422016.E().S(customField.Value)
//line portal-api-restman-template.qtpl:101
		qw422016.N().S(`</l7:StringValue>
              </l7:Property>
              `)
//line portal-api-restman-template.qtpl:103
	}
//line portal-api-restman-template.qtpl:103
	qw422016.N().S(`
            </l7:Properties>
          </l7:ServiceDetail>
          <l7:Resources>
            <l7:ResourceSet tag="policy">
              <l7:Resource type="policy">`)
//line portal-api-restman-template.qtpl:108
	qw422016.E().S(apiServiceXml)
//line portal-api-restman-template.qtpl:108
	qw422016.N().S(`</l7:Resource>
            </l7:ResourceSet>
            `)
//line portal-api-restman-template.qtpl:110
	if isSoapApi == "true" {
//line portal-api-restman-template.qtpl:110
		qw422016.N().S(`
              <l7:ResourceSet tag="wsdl">
                <l7:Resource type="wsdl">`)
//line portal-api-restman-template.qtpl:112
		qw422016.E().S(wsdl)
//line portal-api-restman-template.qtpl:112
		qw422016.N().S(`</l7:Resource>
              </l7:ResourceSet>
            `)
//line portal-api-restman-template.qtpl:114
	}
//line portal-api-restman-template.qtpl:114
	qw422016.N().S(`
          </l7:Resources>
        </l7:Service>
      </l7:Resource>
    </l7:Item>
  </l7:References>
  <l7:Mappings>
    <l7:Mapping action="NewOrUpdate" srcId="`)
//line portal-api-restman-template.qtpl:121
	qw422016.E().S(portalApi.UuidStripped)
//line portal-api-restman-template.qtpl:121
	qw422016.N().S(`" type="POLICY" xmlns:l7="http://ns.l7tech.com/2010/04/gateway-management"/>
    <l7:Mapping action="NewOrUpdate" srcId="`)
//line portal-api-restman-template.qtpl:122
	qw422016.E().S(portalApi.ServiceId)
//line portal-api-restman-template.qtpl:122
	qw422016.N().S(`" type="SERVICE" xmlns:l7="http://ns.l7tech.com/2010/04/gateway-management"/>
  </l7:Mappings>
</l7:Bundle>
`)
//line portal-api-restman-template.qtpl:125
}

//line portal-api-restman-template.qtpl:125
func WriteFromRestmamTemplate(qq422016 qtio422016.Writer, portalApi PortalAPI, apiServiceXml string, apiFragmentXml string, isSoapApi string, wsdl string) {
//line portal-api-restman-template.qtpl:125
	qw422016 := qt422016.AcquireWriter(qq422016)
//line portal-api-restman-template.qtpl:125
	StreamFromRestmamTemplate(qw422016, portalApi, apiServiceXml, apiFragmentXml, isSoapApi, wsdl)
//line portal-api-restman-template.qtpl:125
	qt422016.ReleaseWriter(qw422016)
//line portal-api-restman-template.qtpl:125
}

//line portal-api-restman-template.qtpl:125
func FromRestmamTemplate(portalApi PortalAPI, apiServiceXml string, apiFragmentXml string, isSoapApi string, wsdl string) string {
//line portal-api-restman-template.qtpl:125
	qb422016 := qt422016.AcquireByteBuffer()
//line portal-api-restman-template.qtpl:125
	WriteFromRestmamTemplate(qb422016, portalApi, apiServiceXml, apiFragmentXml, isSoapApi, wsdl)
//line portal-api-restman-template.qtpl:125
	qs422016 := string(qb422016.B)
//line portal-api-restman-template.qtpl:125
	qt422016.ReleaseByteBuffer(qb422016)
//line portal-api-restman-template.qtpl:125
	return qs422016
//line portal-api-restman-template.qtpl:125
}
