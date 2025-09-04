{{/*
Expand the name of the chart.
*/}}
{{- define "slurm.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "slurm.fullname" -}}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "slurm.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "slurm.labels" -}}
helm.sh/chart: {{ include "slurm.chart" . }}
{{ include "slurm.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "slurm.selectorLabels" -}}
app.kubernetes.io/name: {{ include "slurm.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "slurm.serviceAccountName" -}}
{{ default (include "slurm.fullname" .) .Values.serviceAccount.name }}
{{- end }}

{{/*
Create the name of the role to use
*/}}
{{- define "slurm.roleName" -}}
{{ default (include "slurm.fullname" .) .Values.serviceAccount.role.name }}
{{- end }}

{{/*
Create the name of the rolebinding to use
*/}}
{{- define "slurm.roleBindingName" -}}
{{ default (include "slurm.fullname" .) .Values.serviceAccount.roleBinding.name }}
{{- end }}

{{/*
Return the proper Docker Image Registry Secret Names
*/}}
{{- define "slurm.imagePullSecrets" -}}
{{ include "common.images.renderPullSecrets" (dict "images" (list .Values.slurmctld.image .Values.slurmdCPU.image .Values.slurmdGPU.image ) "context" $) }}
{{- end -}}
