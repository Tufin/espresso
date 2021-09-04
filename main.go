package main

import (
	"embed"
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/stretchr/testify/require"
	"github.com/tufin/espresso/bq"
	"github.com/tufin/espresso/internal"
)

func NewShot(project string, fs embed.FS) Shot {
	return Shot{
		bqClient:     bq.NewClient(project),
		sqlTemplates: fs,
		projectID:    project,
	}
}

type Shot struct {
	bqClient     bq.Client
	sqlTemplates embed.FS
	projectID    string
}

func (shot Shot) Run(t *testing.T, query string, test string, params []bigquery.QueryParameter) error {

	query, err := internal.GetQuery(shot.sqlTemplates, query, test)
	if err != nil {
		return err
	}

	queryIterator, err := shot.bqClient.QueryIterator(
		query,
		params)
	require.NoError(t, err, query)

	resultIterator, err := internal.RunSQL(shot.projectID, shot.sqlTemplates, "report_summary_result")
	if err != nil {
		return err
	}

	require.Equal(t,
		internal.ReadResult(queryIterator),
		internal.ReadResult(resultIterator),
	)

	return nil
}
