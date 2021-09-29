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
- Write tests in your own development stack: programming language, IDE and CI/CD pipeline

## Writing Your Own SQL Tests
1. Write an SQL query using Go Text Template notation, for example:
   ```
   {{ define "fruit" }}
 
   WITH base AS (
       {{ .Base }}
   )
   SELECT
       fruit
   FROM base

   {{ end }}
   ```

   The query may contain parameters, like {{ .Base }}
2. Add additional SQL queries to pass as arguments to the main query, for example:  
   ```
   {{ define "base" }}

   SELECT
       "orange" AS fruit
   UNION ALL
   SELECT
       "apple"

   {{ end }}
   ```
   
3. Write your expected result query, for example:
   ```
   {{ define "fruit_result" }}

   SELECT
       "orange" AS fruit
   UNION ALL
   SELECT
       "apple"

   {{ end }}
   ```
4. Create a query definition file describing your query and one or more tests, for example:
   ```
   Name: fruit
   Requires:
   - Base
   Tests:
     Test1:
       Args:
       - Name: Base
         Source: base
       Result:
         Source: fruit_result
   ```
5. Put all files together in a directory

## Query Definitions
The query definition file specifies how to construct an SQL query from the SQL templates.  
A definition file contains Tests, each specifying how to construct an SQL query and an optional expected result.  
To create a query from its components, espresso looks for an SQL template with the same name as the definition file itself, and then interprets the Args.  
Each Arg must have a name field that corresponds to an argument in the corresponding SQL template and one of the following fields:
1. Source - Another SQL template which will be parsed and injected into the containing query.
2. Table - a table name which will be combined with the Google project and BigQuery dataset (defined in Shot) and injected into the containing query.
3. Const - a string that will be injected into the containing query.

Source may contain its own args, or, if it doesn't, espresso will look for a corresponding template file with a test of the same name and parse the args from there.

Result has the same semantics as an Arg, except it has no Name field.

## Access To BigQuery
The tests require access to BigQuery API. 
Please set the following environment variables to grant espresso access to BigQuery:
- `export GCLOUD_PROJECT_ID=<your GCP project ID>`
- `export BIGQUERY_KEY=<a service account with permissions to use BigQuery>`

## Running Tests From The Command-line
```
go build
./espresso -dir="./shot/queries/fruit/" -query="fruit" -test="Test1"
```

## Running Tests From Docker
```
docker run --rm -t -e GCLOUD_PROJECT_ID=$GCLOUD_PROJECT_ID -e BIGQUERY_KEY=$BIGQUERY_KEY -v $(pwd)/shot:/shot:ro tufin/espresso -dir="/shot" -query="fruit" -test="Test1"
```

## Running Tests From Golang
1. Embed your tests directory
2. Create an "Espresso Shot" and run it
3. Use standard Go assertions to check the expected result against the actual output
```
func TestEspressoShot_Filesystem(t *testing.T) {
	queryValues, resultValues, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", os.DirFS("./queries/fruit")).RunTest("fruit", "Test1", []bigquery.QueryParameter{}, &map[string]bigquery.Value{})
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}
```

You can also embed the SQL templates directory into the code:
```
//go:embed queries/fruit
var templates embed.FS

func TestEspressoShot_Embed(t *testing.T) {
	queryValues, resultValues, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", templates).RunTest("fruit", "Test1", []bigquery.QueryParameter{}, &map[string]bigquery.Value{})
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}
```

## Running Tests In Other Programming Languages
Currently only Go is supported.
If you'd like to contribute additional language support, please start a dicssussion.

## Current Status
- This is an initial proof-of-concept and request-for-comments
- Please submit your feedback as pull requests, issues or discussions.
