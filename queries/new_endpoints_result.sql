{{ define "new_endpoints_result" }}

(
    SELECT
        "GET" AS request_method,
        "/api/rome/conf/master/tenants" AS path,
        200 AS status_code,
        3 AS hit_count,
        CAST(28 AS FLOAT64) AS avg_request_time,
        FALSE AS new_endpoint_seen,
        FALSE AS new_endpoint_unseen,
        FALSE AS new_error,
        FALSE AS ongoing_error,
        DATETIME(TIMESTAMP "2021-08-04T01:33:00.011152") AS first_seen,
        DATETIME(TIMESTAMP "2021-08-05T01:34:00.011152") AS last_seen,
        # avg yesterday
        CAST(33 AS FLOAT64) AS avg_request_time_yesterday,
        2 AS hit_count_yesterday,
        # stats yesterday
        2 AS stats_hit_count_yesterday,
        [STRUCT(2 AS hit_count, "sx-company" AS name)] AS tenants_yesterday,
        [
            STRUCT(1 AS hit_count, "effi" AS name),
            STRUCT(1 AS hit_count, "reuven" AS name)
        ] AS users_yesterday,
        [STRUCT(2 AS hit_count, "Java" AS name)] AS agents_yesterday,
        1 AS total_tenants_yesterday,
        2 AS total_users_yesterday,
        1 AS total_agents_yesterday,
        # stats all
        3 AS stats_hit_count,
        [STRUCT(3 AS hit_count, "sx-company" AS name)] AS tenants,
        [
            STRUCT(2 AS hit_count, "effi" AS name),
            STRUCT(1 AS hit_count, "reuven" AS name)
        ] AS users,
        [STRUCT(3 AS hit_count, "Java" AS name)] AS agents,
        1 AS total_tenants,
        2 AS total_users,
        1 AS total_agents,
    UNION ALL
    SELECT
        "GET" AS request_method,
        "/api/rome/conf/master/tenants" AS path,
        500 AS status_code,
        1 AS hit_count,
        CAST(28 AS FLOAT64) AS avg_request_time,
        FALSE AS new_endpoint_seen,
        FALSE AS new_endpoint_unseen,
        TRUE AS new_error,
        FALSE AS ongoing_error,
        DATETIME(TIMESTAMP "2021-08-05T01:33:00.011152") AS first_seen,
        DATETIME(TIMESTAMP "2021-08-05T01:33:00.011152") AS last_seen,
        # avg yesterday
        CAST(28 AS FLOAT64) AS avg_request_time_yesterday,
        1 AS hit_count_yesterday,
        # stats yesterday
        1 AS stats_hit_count_yesterday,
        [STRUCT(1 AS hit_count, "sx-company" AS name)] AS tenants_yesterday,
        [STRUCT(1 AS hit_count, "effi" AS name)] AS users_yesterday,
        [STRUCT(1 AS hit_count, "Java" AS name)] AS agents_yesterday,
        1 AS total_tenants_yesterday,
        1 AS total_users_yesterday,
        1 AS total_agents_yesterday,
        # stats all
        1 AS stats_hit_count,
        [STRUCT(1 AS hit_count, "sx-company" AS name)] AS tenants,
        [STRUCT(1 AS hit_count, "effi" AS name)] AS users,
        [STRUCT(1 AS hit_count, "Java" AS name)] AS agents,
        1 AS total_tenants,
        1 AS total_users,
        1 AS total_agents,
)

{{ end }}
