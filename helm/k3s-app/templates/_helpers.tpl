{{/*
Expand the name of the chart.
*/}}
{{- define "k3s-app.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "k3s-app.fullname" -}}
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
{{- define "k3s-app.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "k3s-app.labels" -}}
helm.sh/chart: {{ include "k3s-app.chart" . }}
{{ include "k3s-app.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "k3s-app.selectorLabels" -}}
app.kubernetes.io/name: {{ include "k3s-app.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Frontend name
*/}}
{{- define "k3s-app.frontend.name" -}}
{{- printf "%s-frontend" (include "k3s-app.fullname" .) }}
{{- end }}

{{/*
API name
*/}}
{{- define "k3s-app.api.name" -}}
{{- printf "%s-api" (include "k3s-app.fullname" .) }}
{{- end }}

{{/*
RPC name
*/}}
{{- define "k3s-app.rpc.name" -}}
{{- printf "%s-rpc" (include "k3s-app.fullname" .) }}
{{- end }}

{{/*
MySQL name
*/}}
{{- define "k3s-app.mysql.name" -}}
{{- printf "%s-mysql" (include "k3s-app.fullname" .) }}
{{- end }}

{{/*
Redis name
*/}}
{{- define "k3s-app.redis.name" -}}
{{- printf "%s-redis" (include "k3s-app.fullname" .) }}
{{- end }}

{{/*
Service Account name
*/}}
{{- define "k3s-app.serviceAccountName" -}}
{{- if .Values.rbac.create }}
{{- default (printf "%s-sa" (include "k3s-app.fullname" .)) .Values.rbac.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.rbac.serviceAccount.name }}
{{- end }}
{{- end }}


