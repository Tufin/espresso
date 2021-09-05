# espresso - a framework for testing BigQuery queries

## Goals
- Componentization: compose complex queries from smaller, reusable components
- Test driven development: write tests for each query and run them like unit tests (except for the fact that they make calls to BigQuery)
- Data as code: input and required output for tests can be defined as part of the code (as well as in real database tables)
- No new languages to learn
- Allow the user to run tests in their own development stack incl. programming language, IDE and CI/CD pipeline

## Using
1. Write an SQL query using Go Text Template notation, for example [report_summary.sql](queries/report_summary.sql).
   The query can contain parameters.
2. Add additional SQL queries to pass as paramaters to the main query. These can be data files like [new_endpoints_input.sql](queries/new_endpoints_input.sql) or additional sub-queries like [get_new_endpoints.sql](queries/get_new_endpoints.sql)
3. Write your result query like [report_summary_result.sql](queries/report_summary_result.sql)
4. Create a test definition YAML file decribing your query and tests, for example: [report_summary.yaml](queries/report_summary.yaml)
5. Put all files in a directory, [create an Espresso Shot and run the test](shot_test.go)

## Current Status
- This is an initial proof-of-concept and request-for-comments
- Please submit your feedback as pull requests
