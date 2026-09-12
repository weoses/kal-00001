{{/*
Common labels applied to every object.
*/}}
{{- define "memelo.labels" -}}
app.kubernetes.io/part-of: memelo
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Labels identifying a single service's objects, e.g. its Deployment/Service/Secret.
Usage: include "memelo.serviceLabels" (dict "root" $ "name" $name)
*/}}
{{- define "memelo.serviceLabels" -}}
{{ include "memelo.labels" .root }}
app.kubernetes.io/name: {{ .name }}
{{- end }}

{{/*
Full image reference for a service.
Usage: include "memelo.image" (dict "root" $ "svc" $svc)
*/}}
{{- define "memelo.image" -}}
{{ .root.Values.image.registry }}/{{ .svc.imageName }}:{{ .root.Values.image.tag }}
{{- end }}
