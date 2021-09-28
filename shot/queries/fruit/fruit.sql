{{ define "fruit" }}

WITH base AS (
    {{ .Base }}
)
SELECT
    fruit
FROM base

{{ end }}