package internal

import (
	"fmt"

	"cloud.google.com/go/bigquery"
	log "github.com/sirupsen/logrus"
	"github.com/tufin/espresso/bq"
	"google.golang.org/api/iterator"
)

func RunQuery(bqClient bq.Client, query string) (bq.Iterator, error) {

	queryIterator, err := bqClient.QueryIterator(query, []bigquery.QueryParameter{})
	if err != nil {
		log.Errorf("failed to run query with %v", err)
		fmt.Println("query:", query)
		return nil, err
	}

	return queryIterator, nil
}

func ReadResult(queryIterator bq.Iterator) ([]map[string]bigquery.Value, error) {

	result := []map[string]bigquery.Value{}
	for {
		row := map[string]bigquery.Value{}
		err := queryIterator.Next(&row)
		if err != nil {
			if err == iterator.Done {
				break
			}
			err = fmt.Errorf("failed to iterate with %v", err)
			return nil, err
		}
		result = append(result, row)
	}
	return result, nil
}
