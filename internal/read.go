package internal

import (
	"embed"
	"fmt"

	"cloud.google.com/go/bigquery"
	log "github.com/sirupsen/logrus"
	"github.com/tufin/espresso/bq"
	"google.golang.org/api/iterator"
)

func RunSQL(projectID string, fs embed.FS, templateName string) (bq.Iterator, error) {
	bqClient := bq.NewClient(projectID)

	query, err := generateSQL(fs, templateName, map[string]string{})
	if err != nil {
		return nil, err
	}

	queryIterator, err := bqClient.QueryIterator(query, []bigquery.QueryParameter{})
	if err != nil {
		return nil, err
	}

	return queryIterator, nil
}

func ReadResult(queryIterator bq.Iterator) []map[string]bigquery.Value {

	result := []map[string]bigquery.Value{}
	for {
		row := map[string]bigquery.Value{}
		err := queryIterator.Next(&row)
		if err != nil {
			if err == iterator.Done {
				break
			}
			err = fmt.Errorf("failed to iterate on with '%v'", err)
			log.Error(err)
			continue
		}
		result = append(result, row)
	}
	return result
}
