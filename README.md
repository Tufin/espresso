[![CI](https://github.com/Tufin/espresso/workflows/go/badge.svg)](https://github.com/Tufin/espresso/actions)
[![codecov](https://codecov.io/gh/tufin/espresso/branch/main/graph/badge.svg?token=4neEgts50n)](https://codecov.io/gh/tufin/espresso)
[![Go Report Card](https://goreportcard.com/badge/github.com/tufin/espresso)](https://goreportcard.com/report/github.com/tufin/espresso)
[![GoDoc](https://godoc.org/github.com/tufin/espresso?status.svg)](https://godoc.org/github.com/tufin/espresso)
[![Docker Image Version](https://img.shields.io/docker/v/tufin/espresso?sort=semver)](https://hub.docker.com/r/tufin/espresso/tags)

# espresso - a framework for testing BigQuery queries

## Goals
- Componentization: compose complex queries from smaller, reusable components
- Test driven development: write tests for each query and run them like unit tests (except for the fact that they make calls to BigQuery)
- Data as code: input and required output for tests can be defined as part of the code (as well as in real database tables)
- No new languages to learn
- Run tests in your own development stack: programming language, IDE and CI/CD pipeline

## Writing Your Own SQL Tests
1. Write an SQL query using Go Text Template notation, for example [report_summary.sql](shot/queries/report_summary.sql):
   ```
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
   ```
   The query may contain parameters.
2. Add additional SQL queries to pass as paramaters to the main query.  
   These can be data files like [new_endpoints_input.sql](shot/queries/new_endpoints_input.sql) or additional sub-queries like [get_new_endpoints.sql](shot/queries/get_new_endpoints.sql)
3. Write your result query like [report_summary_result.sql](shot/queries/report_summary_result.sql):
   ```
   {{ define "report_summary_result" }}

    (
        SELECT
            2 AS hit_count_yesterday,
            1 AS total_endpoints,
            1 AS total_endpoints_yesterday,
    )

    {{ end }}
   ```
    The test will expect the result of the test to be equal to this.
4. Create a test definition YAML file decribing your query and tests like [report_summary.yaml](shot/queries/report_summary.yaml):
   ```
   Name: report_summary
   Requires:
   - Endpoints
   Tests:
     Test1:
       Args:
       - Name: Endpoints
         Source: get_new_endpoints
         Args:
         - Name: NewEndpoints
           Source: new_endpoints_input
       Result: 
         Source: report_summary_result
   ```
5. Put all files together in a directory

## Access To BigQuery
The tests require access to BigQuery API. 
Please set the following environment variables to grant espresso access to BigQuery:
- `export GCLOUD_PROJECT_ID=<your GCP project ID>`
- `export BIGQUERY_KEY=<a service account with permissions to use BigQuery>`

## Running Tests From The Command-line
```
go build
./espresso -dir="./shot/queries/endpoints/" -def="report_summary.yaml" -test="Test1"
````

## Running Tests From Docker
```
docker run --rm -t -e GCLOUD_PROJECT_ID=$GCLOUD_PROJECT_ID -e BIGQUERY_KEY=$BIGQUERY_KEY -v $(pwd)/shot:/shot:ro tufin/espresso -dir="/shot" -def="queries/endpoints/report_summary.yaml" -test="Test1"
```

## Running Tests From Golang
1. Embed your tests directory
2. Create an "Espresso Shot" and run it
3. Use standard Go assertions to check the expected result against the actual output
```
//go:embed queries/endpoints
var endpointTemplates embed.FS

func TestEspressoShot_Embed(t *testing.T) {
	queryValues, resultValues, err := shot.NewShot(env.GetGCPProjectID(), endpointTemplates).RunTest("queries/endpoints/report_summary.yaml", "Test1", []bigquery.QueryParameter{})
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}
```

You can also pass the tests directory without embedding it:
```
func TestEspressoShot_Filesystem(t *testing.T) {
	queryValues, resultValues, err := shot.NewShot(env.GetGCPProjectID(), os.DirFS("./queries/endpoints")).RunTest("report_summary.yaml", "Test1", []bigquery.QueryParameter{})
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}
```

## Current Status
- This is an initial proof-of-concept and request-for-comments
- Please submit your feedback as pull requests
