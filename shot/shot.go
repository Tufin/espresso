package shot

import (
	"fmt"
	"io/fs"

	"cloud.google.com/go/bigquery"
	log "github.com/sirupsen/logrus"
	"github.com/tufin/espresso/shot/bq"
)

// Shot is used to load queries and run tests for them
type Shot struct {
	bqClient  bq.Client
	fsys      fs.FS
	projectID string
	dataset   string
}

func NewShot(client bq.Client, project string, dataset string, fs fs.FS) Shot {

	return Shot{
		bqClient:  client,
		fsys:      fs,
		projectID: project,
		dataset:   dataset,
	}
}

func NewShotWithClient(project string, dataset string, fs fs.FS) Shot {
	return Shot{
		bqClient:  bq.NewClient(project),
		fsys:      fs,
		projectID: project,
		dataset:   dataset,
	}
}

/*
RunQuery runs a single SQL query against BigQuery
query is the name of the query. There must be a correponding yaml definition file and an SQL template in the filesystem.
testName is the name of the test to run, it must appear in the yaml definition file
params are BigQuery paramaters
row is a variable of type to be read from the table (must adhere to https://pkg.go.dev/cloud.google.com/go/bigquery#RowIterator.Next requirements)
The result will be a slice of the same type of 'row' with the result of the query
*/
func (shot Shot) RunQuery(query string, testName string, params []bigquery.QueryParameter, row interface{}) ([]interface{}, error) {
	metadata, err := getMetadata(shot.fsys, query)
	if err != nil {
		log.Errorf("failed to get metadata with %v", err)
		return nil, err
	}

	test, ok := metadata.Tests[testName]
	if !ok {
		err := fmt.Errorf("test %q undefined", testName)
		log.Error(err)
		return nil, err
	}

	queryValues, err := shot.loadAndRun(metadata.Name, testName, test.Args, params, row)
	if err != nil {
		return nil, err
	}

	return queryValues, nil
}

/*
GetTestResults performs a test by running two SQL queries: 
Query 1: The test query
Query 2: The exptected result query
query is the name of the query. There must be a correponding yaml definition file and an SQL template in the filesystem.
testName is the name of the test to run, it must appear in the yaml definition file
params are BigQuery paramaters
row is a variable of type to be read from the table (must adhere to https://pkg.go.dev/cloud.google.com/go/bigquery#RowIterator.Next requirements)
The results will be two slices of the same type of 'row' with the result of each query (test and result)
*/
func (shot Shot) GetTestResults(query string, testName string, params []bigquery.QueryParameter, row interface{}) (interface{}, interface{}, error) {

	metadata, err := getMetadata(shot.fsys, query)
	if err != nil {
		log.Errorf("failed to get metadata with %v", err)
		return nil, nil, err
	}

	test, ok := metadata.Tests[testName]
	if !ok {
		err := fmt.Errorf("test %q undefined", testName)
		log.Error(err)
		return nil, nil, err
	}

	queryValues, err := shot.loadAndRun(metadata.Name, testName, test.Args, params, row)
	if err != nil {
		return nil, nil, err
	}

	if test.Result.Source == "" {
		return nil, nil, fmt.Errorf("result source is missing")
	}

	resultValues, err := shot.loadAndRun(test.Result.Source, testName, test.Result.Args, params, row)
	if err != nil {
		return nil, nil, err
	}

	return queryValues, resultValues, nil
}

/*
RunTest performs a test by running a combined query and checking that the result is empty:
"The test query"
EXCEPT DISTINCT 
"The exptected result query"

query is the name of the query. There must be a correponding yaml definition file and an SQL template in the filesystem.
testName is the name of the test to run, it must appear in the yaml definition file
params are BigQuery paramaters
The result will true iff the result is empty
*/
func (shot Shot) RunTest(query string, testName string, params []bigquery.QueryParameter) (bool, error) {
	metadata, err := getMetadata(shot.fsys, query)
	if err != nil {
		log.Errorf("failed to get metadata with %v", err)
		return false, err
	}

	test, ok := metadata.Tests[testName]
	if !ok {
		err := fmt.Errorf("test %q undefined", testName)
		log.Error(err)
		return false, err
	}

	testQuery, err := shot.getQuery(metadata.Name, testName, test.Args)
	if err != nil {
		return false, err
	}

	if test.Result.Source == "" {
		return false, fmt.Errorf("result source is missing")
	}

	resultQuery, err := shot.getQuery(test.Result.Source, testName, test.Result.Args)
	if err != nil {
		return false, err
	}

	queryIterator, err := runQuery(shot.bqClient, testQuery+"\nEXCEPT DISTINCT\n("+resultQuery+"\n)", params)
	if err != nil {
		return false, err
	}

	result, err := readResult(queryIterator, &map[string]bigquery.Value{})
	if err != nil {
		return false, err
	}

	return len(result) == 0, nil
}

/*
GetQuery returns the SQL query
query is the name of the query. There must be a correponding yaml definition file and an SQL template in the filesystem.
testName is the name of the test to run, it must appear in the yaml definition file
params are BigQuery paramaters
*/
func (shot Shot) GetQuery(query string, testName string, params []bigquery.QueryParameter) (string, error) {
	metadata, err := getMetadata(shot.fsys, query)
	if err != nil {
		log.Errorf("failed to get metadata with %v", err)
		return "", err
	}

	test, ok := metadata.Tests[testName]
	if !ok {
		err := fmt.Errorf("test %q undefined", testName)
		log.Error(err)
		return "", err
	}

	result, err := shot.getQuery(metadata.Name, testName, test.Args)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (shot Shot) loadAndRun(templateName string, testName string, args []argument, params []bigquery.QueryParameter, row interface{}) ([]interface{}, error) {
	query, err := shot.getQuery(templateName, testName, args)
	if err != nil {
		return nil, err
	}

	queryIterator, err := runQuery(shot.bqClient, query, params)
	if err != nil {
		return nil, err
	}

	result, err := readResult(queryIterator, row)
	if err != nil {
		return nil, err
	}

	return result, nil
}
