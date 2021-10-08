{{ define "base_duplicates" }}

SELECT
    "orange" AS fruit
UNION ALL
SELECT
    "orange"
UNION ALL
SELECT
    "apple"

{{ end }}