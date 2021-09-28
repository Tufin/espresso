{{ define "base_with_arg" }}

SELECT
    "orange" AS fruit
UNION ALL
SELECT
    "{{ .Fruit }}"

{{ end }}