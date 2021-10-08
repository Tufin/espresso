{{ define "fruit_result_duplicates" }}

SELECT
    "orange" AS fruit
UNION ALL
SELECT
    "apple"
UNION ALL
SELECT
    "apple"

{{ end }}