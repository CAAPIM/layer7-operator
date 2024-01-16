# API Reference

Packages:

- [security.brcmlabs.com/v1alpha1](#securitybrcmlabscomv1alpha1)

# security.brcmlabs.com/v1alpha1

Resource Types:

- [L7Portal](#l7portal)




## L7Portal
<sup><sup>[↩ Parent](#securitybrcmlabscomv1alpha1 )</sup></sup>






L7Portal is the Schema for the l7portals API

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
      <td>L7Portal</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#l7portalspec">spec</a></b></td>
        <td>object</td>
        <td>
          L7PortalSpec defines the desired state of L7Portal<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#l7portalstatus">status</a></b></td>
        <td>object</td>
        <td>
          L7PortalStatus defines the observed state of L7Portal<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### L7Portal.spec
<sup><sup>[↩ Parent](#l7portal)</sup></sup>



L7PortalSpec defines the desired state of L7Portal

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
        <td><b><a href="#l7portalspecauth">auth</a></b></td>
        <td>object</td>
        <td>
          Auth - Portal credentials<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>deploymentTags</b></td>
        <td>[]string</td>
        <td>
          Deployment Tags - determines which Gateway deployments these APIs will be applied to<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Enabled - if enabled this Portal and its APIs will be synced<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>endpoint</b></td>
        <td>string</td>
        <td>
          Endoint - Portal endpoint<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>enrollmentBundle</b></td>
        <td>string</td>
        <td>
          EnrollmentBundle - allows a custom enrollment bundle to be set in the Portal CR<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>labels</b></td>
        <td>map[string]string</td>
        <td>
          Labels - Custom Labels<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>mode</b></td>
        <td>string</td>
        <td>
          Mode determines how or if the Portal is contacted defaults to auto, options are auto, local. Local requires enrollmentBundle to be set.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name Portal name<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>syncIntervalSeconds</b></td>
        <td>integer</td>
        <td>
          SyncIntervalSeconds how often the Portal CR is reconciled. Default is 10 seconds<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### L7Portal.spec.auth
<sup><sup>[↩ Parent](#l7portalspec)</sup></sup>



Auth - Portal credentials

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
        <td><b>clientId</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>clientSecret</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>endpoint</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>existingSecretName</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### L7Portal.status
<sup><sup>[↩ Parent](#l7portal)</sup></sup>



L7PortalStatus defines the observed state of L7Portal

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
        <td><b>apiCount</b></td>
        <td>integer</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>apiSummaryConfigMap</b></td>
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
        <td><b><a href="#l7portalstatusenrollmentbundle">enrollmentBundle</a></b></td>
        <td>object</td>
        <td>
          EnrollmentBundle<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastUpdated</b></td>
        <td>integer</td>
        <td>
          <br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#l7portalstatusproxiesindex">proxies</a></b></td>
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


### L7Portal.status.enrollmentBundle
<sup><sup>[↩ Parent](#l7portalstatus)</sup></sup>



EnrollmentBundle

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
        <td><b>lastUpdated</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>secretName</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### L7Portal.status.proxies[index]
<sup><sup>[↩ Parent](#l7portalstatus)</sup></sup>



GatewayProxy

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
        <td><b><a href="#l7portalstatusproxiesindexgatewaysindex">gateways</a></b></td>
        <td>[]object</td>
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
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type - Ephemeral or DbBacked<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### L7Portal.status.proxies[index].gateways[index]
<sup><sup>[↩ Parent](#l7portalstatusproxiesindex)</sup></sup>





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
        <td><b>lastUpdated</b></td>
        <td>string</td>
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
        <td><b>synchronised</b></td>
        <td>boolean</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>
