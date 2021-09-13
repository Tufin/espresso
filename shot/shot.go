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
	bqClient     bq.Client
	sqlTemplates fs.FS
	projectID    string
	dataset      string
}

func NewShot(project string, dataset string, fs fs.FS) Shot {

	return Shot{
		bqClient:     bq.NewClient(project),
		sqlTemplates: fs,
		projectID:    project,
		dataset:      dataset,
	}
}

func (shot Shot) RunTest(testDefinitionPath string, testName string, params []bigquery.QueryParameter) ([]map[string]bigquery.Value, []map[string]bigquery.Value, error) {

	metadata, err := getMetadata(shot.sqlTemplates, testDefinitionPath)
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

	queryValues, err := shot.loadAndRun(metadata.Name, test.Args, params)
	if err != nil {
		return nil, nil, err
	}

	resultValues, err := shot.loadAndRun(test.Result.Source, test.Result.Args, params)
	if err != nil {
		return nil, nil, err
	}

	return queryValues, resultValues, nil
}

func (shot Shot) loadAndRun(templateName string, args []argument, params []bigquery.QueryParameter) ([]map[string]bigquery.Value, error) {
	query, err := shot.getQuery(templateName, args)
	if err != nil {
		return nil, err
	}

	queryIterator, err := runQuery(shot.bqClient, query, params)
	if err != nil {
		return nil, err
	}

	result, err := readResult(queryIterator)
	if err != nil {
		return nil, err
	}

	return result, nil
}
