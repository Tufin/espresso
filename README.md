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
        Result: report_summary_result
   ```
5. Put all files together in a directory

## Running Tests From The Command-line
```
go build
./espresso -dir="shot/queries" -def="report_summary.yaml" -query="report_summary" -test="Test1"
````

## Running Tests From Golang
1. Embed your tests directory
2. Create an "Espresso Shot" and run it
3. Use standard Go assertions to check the expected result against the actual output
```
//go:embed queries
var sqlTemplates embed.FS

func TestEspressoShot_Embed(t *testing.T) {
	queryValues, resultValues, err := shot.NewShot(env.GetGCPProjectID(), sqlTemplates).RunTest("queries/report_summary.yaml", "report_summary", "Test1", []bigquery.QueryParameter{})
	require.NoError(t, err)
	require.Equal(t, queryValues, resultValues)
}
```

You can also pass the tests directory without embedding it:
```
func TestEspressoShot_Filesystem(t *testing.T) {
	fileSystem := os.DirFS(".")
	queryValues, resultValues, err := shot.NewShot(env.GetGCPProjectID(), fileSystem).RunTest("queries/report_summary.yaml", "report_summary", "Test1", []bigquery.QueryParameter{})
	require.NoError(t, err)
	require.Equal(t, queryValues, resultValues)
}
```

## Access To BigQuery
Please set the following environment variables to grant espresso access to BigQuery:
- `export GCLOUD_PROJECT_ID=<your GCP project ID>`
- `export BIGQUERY_KEY=<a service account with permissions to use BigQuery>`

## Current Status
- This is an initial proof-of-concept and request-for-comments
- Please submit your feedback as pull requests
