{{/*
Expand the name of the chart.
*/}}
{{- define "layer7-operator.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "layer7-operator.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "layer7-operator.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "layer7-operator.labels" -}}
helm.sh/chart: {{ include "layer7-operator.chart" . }}
{{ include "layer7-operator.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "layer7-operator.selectorLabels" -}}
app.kubernetes.io/name: {{ include "layer7-operator.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "layer7-operator.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "layer7-operator.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Managed Namespaces.
*/}}
{{- define "layer7-operator.managedNamespaces" -}}
{{- $managedNamespaces := len .Values.managedNamespaces }}
{{- if gt $managedNamespaces 0 }}
  {{- join "," .Values.managedNamespaces }}
{{- end -}}
{{- end -}}

{{/*
Trusted CA Configmap name.
*/}}
{{- define "layer7-operator.trusted-ca" -}}
{{- if .Values.proxy.caBundle.enabled }}
{{- if and (not .Values.proxy.caBundle.create) (.Values.proxy.caBundle.existingConfigmap) }}
  {{- print .Values.proxy.caBundle.existingConfigmap }}
{{- else }}
  {{- default (include "layer7-operator.fullname" .) }}-trusted-ca
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Trusted CA Configmap key.
*/}}
{{- define "layer7-operator.trusted-ca-key" -}}
{{- if  and (.Values.proxy.caBundle.enabled) (.Values.proxy.caBundle.key)  }}
  {{- print .Values.proxy.caBundle.existingConfigmap }}
{{- else }}
  {{- default "ca-bundle.crt" }}
{{- end -}}
{{- end -}}


{{/*
Webhook TLS secret.
*/}}
{{- define "layer7-operator.webhook-secret" -}}
{{- if  and (.Values.webhook.enabled) (.Values.webhook.tls.certmanager.enabled)  }}
  {{- default (include "layer7-operator.fullname" .) }}-serving-cert
{{- else }}
  {{- default .Values.webhook.tls.existingTlsSecret }}
{{- end -}}
{{- end -}}

{{ include "layer7-operator.fullname" . }}-serving-cert