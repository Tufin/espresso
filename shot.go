package espresso

import (
	"embed"
	"fmt"
	"testing"

	"cloud.google.com/go/bigquery"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/tufin/espresso/bq"
	"github.com/tufin/espresso/internal"
)

type Shot struct {
	bqClient     bq.Client
	sqlTemplates embed.FS
	projectID    string
}

func NewShot(project string, fs embed.FS) Shot {
	return Shot{
		bqClient:     bq.NewClient(project),
		sqlTemplates: fs,
		projectID:    project,
	}
}

func (shot Shot) RunTest(t *testing.T, path string, queryName string, testName string, params []bigquery.QueryParameter) error {

	metadata, err := internal.GetMetadata(shot.sqlTemplates, path, queryName)
	if err != nil {
		log.Errorf("failed to get metadata with '%v'", err)
		return err
	}

	test, ok := metadata.Tests[testName]
	if !ok {
		err := fmt.Errorf("test '%s' undefined", testName)
		log.Error(err)
		return err
	}

	client := bq.NewClient(shot.projectID)
	if err != nil {
		log.Errorf("failed to create client with '%v'", err)
		return err
	}

	queryValues, err := loadAndRun(client, shot.sqlTemplates, path, queryName, test.Args)
	if err != nil {
		log.Errorf("failed to run query with '%v'", err)
		return err
	}

	resultValues, err := loadAndRun(client, shot.sqlTemplates, path, test.Result, []internal.Argument{})
	if err != nil {
		log.Errorf("failed to run result query with '%v'", err)
		return err
	}

	require.Equal(t,
		queryValues,
		resultValues,
	)

	return nil
}

func loadAndRun(client bq.Client, fs embed.FS, path string, testName string, args []internal.Argument) ([]map[string]bigquery.Value, error) {
	query, err := internal.GetQuery(fs, path, testName, args)
	if err != nil {
		return nil, err
	}

	queryIterator, err := internal.RunQuery(client, query)
	if err != nil {
		return nil, err
	}

	return internal.ReadResult(queryIterator)
}
