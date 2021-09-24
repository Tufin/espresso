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
query is the name of the query. There must be a correponding yaml definition file and a template in the filesystem.
testName is the name of the test to run, it must appear in the yaml definition file
params are BigQuery paramaters
row is the argument that will be passed to bigquery.RowIterator.Next
The result will be a slice of the same type of 'row' with the return values of the query
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
RunTest performs a BigQuery test by running two SQL queries: one for the test and another for the exptected result
query is the name of the query. There must be a correponding yaml definition file and a template in the filesystem.
testName is the name of the test to run, it must appear in the yaml definition file
params are BigQuery paramaters
row is the argument that will be passed to bigquery.RowIterator.Next
The results will be slices of the same type of 'row' with the return values of the query and corresponding result
*/
func (shot Shot) RunTest(query string, testName string, params []bigquery.QueryParameter, row interface{}) (interface{}, interface{}, error) {

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

	resultValues, err := shot.loadAndRun(test.Result.Source, testName, test.Result.Args, params, row)
	if err != nil {
		return nil, nil, err
	}

	return queryValues, resultValues, nil
}

/*
GetQuery returns the SQL query
query is the name of the query. There must be a correponding yaml definition file and a template in the filesystem.
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
