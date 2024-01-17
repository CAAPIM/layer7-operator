# API Reference

Packages:

- [security.brcmlabs.com/v1](#securitybrcmlabscomv1)

# security.brcmlabs.com/v1

Resource Types:

- [Repository](#repository)




## Repository
<sup><sup>[↩ Parent](#securitybrcmlabscomv1 )</sup></sup>






Repository is the Schema for the repositories API

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
      <td>security.brcmlabs.com/v1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>Repository</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#repositoryspec">spec</a></b></td>
        <td>object</td>
        <td>
          Spec - Repository Spec<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#repositorystatus">status</a></b></td>
        <td>object</td>
        <td>
          Status - Repository Status<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Repository.spec
<sup><sup>[↩ Parent](#repository)</sup></sup>



Spec - Repository Spec

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
        <td><b>endpoint</b></td>
        <td>string</td>
        <td>
          Endoint - Git repository endpoint<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>annotations</b></td>
        <td>map[string]string</td>
        <td>
          Annotations - Custom Annotations<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#repositoryspecauth">auth</a></b></td>
        <td>object</td>
        <td>
          Auth contains a reference to the credentials required to connect to your Git repository<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>branch</b></td>
        <td>string</td>
        <td>
          Branch - specify which branch to clone if branch and tag are both specified branch will take precedence and tag will be ignored if branch and tag are both missing the entire repository will be cloned<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Enabled - if enabled this repository will be synced<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>labels</b></td>
        <td>map[string]string</td>
        <td>
          Name Repository name Name string `json:"name"` Labels - Custom Labels<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#repositoryspeclocalreference">localReference</a></b></td>
        <td>object</td>
        <td>
          LocalReference lets the Repository controller use a local Kubernetes Configmap/Secret as a repository source This is not currently implemented<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>remoteName</b></td>
        <td>string</td>
        <td>
          Remote Name - defaults to "origin"<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#repositoryspecsync">sync</a></b></td>
        <td>object</td>
        <td>
          RepositorySyncConfig<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>tag</b></td>
        <td>string</td>
        <td>
          Tag - clone a specific tag. tags do not change, once cloned this will not be checked for updates<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Repository.spec.auth
<sup><sup>[↩ Parent](#repositoryspec)</sup></sup>



Auth contains a reference to the credentials required to connect to your Git repository

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
        <td><b>existingSecretName</b></td>
        <td>string</td>
        <td>
          ExistingSecretName reference an existing secret<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>knownHosts</b></td>
        <td>string</td>
        <td>
          KnownHosts is required for SSH Auth<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>password</b></td>
        <td>string</td>
        <td>
          Password repository Password password or token are acceptable<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>sshKey</b></td>
        <td>string</td>
        <td>
          SSHKey for Git SSH Authentication<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>sshKeyPass</b></td>
        <td>string</td>
        <td>
          SSHKeyPass<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>token</b></td>
        <td>string</td>
        <td>
          Token repository Access Token<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Auth Type defaults to basic, possible options are none, basic or ssh<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>username</b></td>
        <td>string</td>
        <td>
          Username repository username<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>vendor</b></td>
        <td>string</td>
        <td>
          Vendor i.e. Github, Gitlab, BitBucket<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Repository.spec.localReference
<sup><sup>[↩ Parent](#repositoryspec)</sup></sup>



LocalReference lets the Repository controller use a local Kubernetes Configmap/Secret as a repository source This is not currently implemented

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
        <td><b>configMapName</b></td>
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


### Repository.spec.sync
<sup><sup>[↩ Parent](#repositoryspec)</sup></sup>



RepositorySyncConfig

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
        <td><b>interval</b></td>
        <td>integer</td>
        <td>
          Configure how frequently the remote is checked for new commits<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Repository.status
<sup><sup>[↩ Parent](#repository)</sup></sup>



Status - Repository Status

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
        <td><b>commit</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lastAppliedSummary</b></td>
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
        <td><b>ready</b></td>
        <td>boolean</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>storageSecretName</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>summary</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>updated</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>vendor</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>
