package shot

import (
	"fmt"
	"io/fs"

	"cloud.google.com/go/bigquery"
	log "github.com/sirupsen/logrus"
	"github.com/tufin/espresso/bq"
	"github.com/tufin/espresso/shot/internal"
)

// Shot is used to load queries and run tests for them
type Shot struct {
	bqClient     bq.Client
	sqlTemplates fs.FS
	projectID    string
}

func NewShot(project string, fs fs.FS) Shot {
	return Shot{
		bqClient:     bq.NewClient(project),
		sqlTemplates: fs,
		projectID:    project,
	}
}

func (shot Shot) RunTest(testDefinitionPath string, templateName string, testName string, params []bigquery.QueryParameter) ([]map[string]bigquery.Value, []map[string]bigquery.Value, error) {

	metadata, err := internal.GetMetadata(shot.sqlTemplates, testDefinitionPath)
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

	client := bq.NewClient(shot.projectID)
	if err != nil {
		log.Errorf("failed to create client with %v", err)
		return nil, nil, err
	}

	queryValues, err := loadAndRun(client, shot.sqlTemplates, templateName, test.Args, params)
	if err != nil {
		return nil, nil, err
	}

	resultValues, err := loadAndRun(client, shot.sqlTemplates, test.Result.Source, test.Result.Args, params)
	if err != nil {
		return nil, nil, err
	}

	return queryValues, resultValues, nil
}

func loadAndRun(client bq.Client, fs fs.FS, templateName string, args []internal.Argument, params []bigquery.QueryParameter) ([]map[string]bigquery.Value, error) {
	query, err := internal.GetQuery(fs, templateName, args)
	if err != nil {
		return nil, err
	}

	queryIterator, err := internal.RunQuery(client, query, params)
	if err != nil {
		return nil, err
	}

	result, err := internal.ReadResult(queryIterator)
	if err != nil {
		return nil, err
	}

	return result, nil
}
