{{ define "report_summary" }}

WITH base AS (
    SELECT
        request_method,
        path,
        SUM(hit_count_yesterday) AS hit_count_yesterday,
    FROM {{ .Endpoints }}
    WHERE status_code<400
    GROUP BY 
        request_method,
        path
)
SELECT
    COUNT(*) AS total_endpoints,
    COUNTIF(hit_count_yesterday>0) AS total_endpoints_yesterday,
    SUM(hit_count_yesterday) AS hit_count_yesterday,
FROM base

{{ end }}