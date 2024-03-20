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
        <td><b><a href="#l7portalspecbootstrap">bootstrap</a></b></td>
        <td>object</td>
        <td>
          Bootstrap configures an initContainer that is injected into the Gateway Deployment. It contains Portal specific L7Apis that are bulk loaded via bootstrap. this does not currently include the Portal enrollment bundle which is currently loaded via a different initContainer<br/>
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
        <td><b>portalManaged</b></td>
        <td>boolean</td>
        <td>
          PortalManaged if PortalManaged is true the portal controller will not manage APIs and will be only be responsible for maintaining a list of L7Api Portal published Metadata.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>portalTenant</b></td>
        <td>string</td>
        <td>
          PortalTenant is the tenantId of the API Developer Portal<br/>
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


### L7Portal.spec.bootstrap
<sup><sup>[↩ Parent](#l7portalspec)</sup></sup>



Bootstrap configures an initContainer that is injected into the Gateway Deployment. It contains Portal specific L7Apis that are bulk loaded via bootstrap. this does not currently include the Portal enrollment bundle which is currently loaded via a different initContainer

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
        <td><b>initContainerImage</b></td>
        <td>string</td>
        <td>
          InitContainerImage for the initContainer<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>initContainerImagePullPolicy</b></td>
        <td>string</td>
        <td>
          InitContainerImagePullPolicy<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#l7portalspecbootstrapinitcontainersecuritycontext">initContainerSecurityContext</a></b></td>
        <td>object</td>
        <td>
          InitContainerSecurityContext<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### L7Portal.spec.bootstrap.initContainerSecurityContext
<sup><sup>[↩ Parent](#l7portalspecbootstrap)</sup></sup>



InitContainerSecurityContext

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
        <td><b>allowPrivilegeEscalation</b></td>
        <td>boolean</td>
        <td>
          AllowPrivilegeEscalation controls whether a process can gain more privileges than its parent process. This bool directly controls if the no_new_privs flag will be set on the container process. AllowPrivilegeEscalation is true always when the container is: 1) run as Privileged 2) has CAP_SYS_ADMIN Note that this field cannot be set when spec.os.name is windows.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#l7portalspecbootstrapinitcontainersecuritycontextcapabilities">capabilities</a></b></td>
        <td>object</td>
        <td>
          The capabilities to add/drop when running containers. Defaults to the default set of capabilities granted by the container runtime. Note that this field cannot be set when spec.os.name is windows.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>privileged</b></td>
        <td>boolean</td>
        <td>
          Run container in privileged mode. Processes in privileged containers are essentially equivalent to root on the host. Defaults to false. Note that this field cannot be set when spec.os.name is windows.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>procMount</b></td>
        <td>string</td>
        <td>
          procMount denotes the type of proc mount to use for the containers. The default is DefaultProcMount which uses the container runtime defaults for readonly paths and masked paths. This requires the ProcMountType feature flag to be enabled. Note that this field cannot be set when spec.os.name is windows.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnlyRootFilesystem</b></td>
        <td>boolean</td>
        <td>
          Whether this container has a read-only root filesystem. Default is false. Note that this field cannot be set when spec.os.name is windows.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>runAsGroup</b></td>
        <td>integer</td>
        <td>
          The GID to run the entrypoint of the container process. Uses runtime default if unset. May also be set in PodSecurityContext.  If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence. Note that this field cannot be set when spec.os.name is windows.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>runAsNonRoot</b></td>
        <td>boolean</td>
        <td>
          Indicates that the container must run as a non-root user. If true, the Kubelet will validate the image at runtime to ensure that it does not run as UID 0 (root) and fail to start the container if it does. If unset or false, no such validation will be performed. May also be set in PodSecurityContext.  If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>runAsUser</b></td>
        <td>integer</td>
        <td>
          The UID to run the entrypoint of the container process. Defaults to user specified in image metadata if unspecified. May also be set in PodSecurityContext.  If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence. Note that this field cannot be set when spec.os.name is windows.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#l7portalspecbootstrapinitcontainersecuritycontextselinuxoptions">seLinuxOptions</a></b></td>
        <td>object</td>
        <td>
          The SELinux context to be applied to the container. If unspecified, the container runtime will allocate a random SELinux context for each container.  May also be set in PodSecurityContext.  If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence. Note that this field cannot be set when spec.os.name is windows.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#l7portalspecbootstrapinitcontainersecuritycontextseccompprofile">seccompProfile</a></b></td>
        <td>object</td>
        <td>
          The seccomp options to use by this container. If seccomp options are provided at both the pod & container level, the container options override the pod options. Note that this field cannot be set when spec.os.name is windows.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#l7portalspecbootstrapinitcontainersecuritycontextwindowsoptions">windowsOptions</a></b></td>
        <td>object</td>
        <td>
          The Windows specific settings applied to all containers. If unspecified, the options from the PodSecurityContext will be used. If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence. Note that this field cannot be set when spec.os.name is linux.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### L7Portal.spec.bootstrap.initContainerSecurityContext.capabilities
<sup><sup>[↩ Parent](#l7portalspecbootstrapinitcontainersecuritycontext)</sup></sup>



The capabilities to add/drop when running containers. Defaults to the default set of capabilities granted by the container runtime. Note that this field cannot be set when spec.os.name is windows.

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
        <td><b>add</b></td>
        <td>[]string</td>
        <td>
          Added capabilities<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>drop</b></td>
        <td>[]string</td>
        <td>
          Removed capabilities<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### L7Portal.spec.bootstrap.initContainerSecurityContext.seLinuxOptions
<sup><sup>[↩ Parent](#l7portalspecbootstrapinitcontainersecuritycontext)</sup></sup>



The SELinux context to be applied to the container. If unspecified, the container runtime will allocate a random SELinux context for each container.  May also be set in PodSecurityContext.  If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence. Note that this field cannot be set when spec.os.name is windows.

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
        <td><b>level</b></td>
        <td>string</td>
        <td>
          Level is SELinux level label that applies to the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>role</b></td>
        <td>string</td>
        <td>
          Role is a SELinux role label that applies to the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type is a SELinux type label that applies to the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>user</b></td>
        <td>string</td>
        <td>
          User is a SELinux user label that applies to the container.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### L7Portal.spec.bootstrap.initContainerSecurityContext.seccompProfile
<sup><sup>[↩ Parent](#l7portalspecbootstrapinitcontainersecuritycontext)</sup></sup>



The seccomp options to use by this container. If seccomp options are provided at both the pod & container level, the container options override the pod options. Note that this field cannot be set when spec.os.name is windows.

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
        <td><b>type</b></td>
        <td>string</td>
        <td>
          type indicates which kind of seccomp profile will be applied. Valid options are: 
 Localhost - a profile defined in a file on the node should be used. RuntimeDefault - the container runtime default profile should be used. Unconfined - no profile should be applied.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>localhostProfile</b></td>
        <td>string</td>
        <td>
          localhostProfile indicates a profile defined in a file on the node should be used. The profile must be preconfigured on the node to work. Must be a descending path, relative to the kubelet's configured seccomp profile location. Must be set if type is "Localhost". Must NOT be set for any other type.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### L7Portal.spec.bootstrap.initContainerSecurityContext.windowsOptions
<sup><sup>[↩ Parent](#l7portalspecbootstrapinitcontainersecuritycontext)</sup></sup>



The Windows specific settings applied to all containers. If unspecified, the options from the PodSecurityContext will be used. If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence. Note that this field cannot be set when spec.os.name is linux.

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
        <td><b>gmsaCredentialSpec</b></td>
        <td>string</td>
        <td>
          GMSACredentialSpec is where the GMSA admission webhook (https://github.com/kubernetes-sigs/windows-gmsa) inlines the contents of the GMSA credential spec named by the GMSACredentialSpecName field.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>gmsaCredentialSpecName</b></td>
        <td>string</td>
        <td>
          GMSACredentialSpecName is the name of the GMSA credential spec to use.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>hostProcess</b></td>
        <td>boolean</td>
        <td>
          HostProcess determines if a container should be run as a 'Host Process' container. All of a Pod's containers must have the same effective HostProcess value (it is not allowed to have a mix of HostProcess containers and non-HostProcess containers). In addition, if HostProcess is true then HostNetwork must also be set to true.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>runAsUserName</b></td>
        <td>string</td>
        <td>
          The UserName in Windows to run the entrypoint of the container process. Defaults to the user specified in image metadata if unspecified. May also be set in PodSecurityContext. If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence.<br/>
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
