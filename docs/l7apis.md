# API Reference

Packages:

- [security.brcmlabs.com/v1alpha1](#securitybrcmlabscomv1alpha1)

# security.brcmlabs.com/v1alpha1

Resource Types:

- [L7Api](#l7api)




## L7Api
<sup><sup>[↩ Parent](#securitybrcmlabscomv1alpha1 )</sup></sup>






L7Api is the Schema for the l7apis API

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>security.brcmlabs.com/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>L7Api</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#l7apispec">spec</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#l7apistatus">status</a></b></td>
        <td>object</td>
        <td>
          L7ApiStatus defines the observed state of L7Api<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### L7Api.spec
<sup><sup>[↩ Parent](#l7api)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>deploymentTags</b></td>
        <td>[]string</td>
        <td>
          DeploymentTags target Gateway deployments that this API should be published to<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>graphmanBundle</b></td>
        <td>string</td>
        <td>
          GraphmanBundle associated with this API currently limited to Service and Fragments auto generated when PortalMeta is set and PortalPublished is true<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>l7Portal</b></td>
        <td>string</td>
        <td>
          L7Portal is the L7Portal that this API is associated with when Portal Published is true<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#l7apispecportalmeta">portalMeta</a></b></td>
        <td>object</td>
        <td>
          PortalMeta is reserved for the API Developer Portal<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>portalPublished</b></td>
        <td>boolean</td>
        <td>
          PortalPublished<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>serviceUrl</b></td>
        <td>string</td>
        <td>
          ServiceUrl on the API Gateway<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### L7Api.spec.portalMeta
<sup><sup>[↩ Parent](#l7apispec)</sup></sup>



PortalMeta is reserved for the API Developer Portal

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>apiId</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>apiUuid</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>checksum</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>createTs</b></td>
        <td>integer</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#l7apispecportalmetacustomfieldvaluesindex">customFieldValues</a></b></td>
        <td>[]object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>locationUrl</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>modifyTs</b></td>
        <td>integer</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#l7apispecportalmetapolicyentitiesindex">policyEntities</a></b></td>
        <td>[]object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>publishedTs</b></td>
        <td>integer</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>serviceId</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ssgServiceType</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ssgUrl</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ssgUrlEncoded</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>tenantId</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### L7Api.spec.portalMeta.customFieldValues[index]
<sup><sup>[↩ Parent](#l7apispecportalmeta)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### L7Api.spec.portalMeta.policyEntities[index]
<sup><sup>[↩ Parent](#l7apispecportalmeta)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>policyEntityUuid</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#l7apispecportalmetapolicyentitiesindexpolicytemplateargumentsindex">policyTemplateArguments</a></b></td>
        <td>[]object</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### L7Api.spec.portalMeta.policyEntities[index].policyTemplateArguments[index]
<sup><sup>[↩ Parent](#l7apispecportalmetapolicyentitiesindex)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### L7Api.status
<sup><sup>[↩ Parent](#l7api)</sup></sup>



L7ApiStatus defines the observed state of L7Api

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>checksum</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#l7apistatusgatewaysindex">gateways</a></b></td>
        <td>[]object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ready</b></td>
        <td>boolean</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### L7Api.status.gateways[index]
<sup><sup>[↩ Parent](#l7apistatus)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>checksum</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>deployment</b></td>
        <td>string</td>
        <td>
          Phase       corev1.PodPhase `json:"phase,omitempty"` Reason: Success/Failed to sync because of x<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastUpdated</b></td>
        <td>string</td>
        <td>
          Ready       bool            `json:"ready,omitempty"`<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>
