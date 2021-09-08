{{ define "fruit" }}

WITH base AS (
    SELECT
        "orange" AS fruit
    UNION ALL
    SELECT
        "apple"
)
SELECT
    fruit
FROM {{ .Base }}

{{ end }}