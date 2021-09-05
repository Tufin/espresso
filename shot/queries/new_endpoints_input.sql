{{ define "new_endpoints_input" }}

(
    SELECT
        "GET" AS request_method,
        "/api/rome/conf/master/tenants" AS path,
        200 AS status_code,
        2 AS hit_count_yesterday,
    UNION ALL
    SELECT
        "GET" AS request_method,
        "/api/rome/conf/master/tenants" AS path,
        500 AS status_code,
        1 AS hit_count_yesterday,
)

{{ end }}
