/*
Requires:
- NewEndpoints
*/

{{ define "get_new_endpoints" }}

(
    SELECT * FROM {{ .NewEndpoints }}
)

{{ end }}
