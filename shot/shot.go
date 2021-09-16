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

func NewShot(project string, dataset string, fs fs.FS) Shot {

	return Shot{
		bqClient:  bq.NewClient(project),
		fsys:      fs,
		projectID: project,
		dataset:   dataset,
	}
}

func (shot Shot) RunQuery(query string, testName string, params []bigquery.QueryParameter, row interface{}) (interface{}, error) {
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

func (shot Shot) loadAndRun(templateName string, testName string, args []argument, params []bigquery.QueryParameter, row interface{}) (interface{}, error) {
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
