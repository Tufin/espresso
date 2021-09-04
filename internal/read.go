package internal

import (
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/tufin/espresso/bq"
	"google.golang.org/api/iterator"
)

func runSQL(bqClient bq.Client, query string) (bq.Iterator, error) {

	queryIterator, err := bqClient.QueryIterator(query, []bigquery.QueryParameter{})
	if err != nil {
		return nil, err
	}

	return queryIterator, nil
}

func readResult(queryIterator bq.Iterator) ([]map[string]bigquery.Value, error) {

	result := []map[string]bigquery.Value{}
	for {
		row := map[string]bigquery.Value{}
		err := queryIterator.Next(&row)
		if err != nil {
			if err == iterator.Done {
				break
			}
			err = fmt.Errorf("failed to iterate with '%v'", err)
			return nil, err
		}
		result = append(result, row)
	}
	return result, nil
}
