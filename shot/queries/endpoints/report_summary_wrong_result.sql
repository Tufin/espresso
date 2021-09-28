{{ define "report_summary_wrong_result" }}

(
    SELECT
        2 AS hit_count_yesterday,
        1 AS total_endpoints,
        2 AS total_endpoints_yesterday,
)

{{ end }}
