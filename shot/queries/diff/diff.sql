{{ define "diff" }}

WITH table1 AS (
    {{ .Table1 }}
),
table2 AS (
    {{ .Table2 }}
)
(
  SELECT * FROM table1
  EXCEPT DISTINCT
  SELECT * from table2
)

UNION ALL

(
  SELECT * FROM table2
  EXCEPT DISTINCT
  SELECT * from table1
)

{{ end }}