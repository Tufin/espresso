package shot

import (
	"fmt"
	"reflect"

	"cloud.google.com/go/bigquery"
	log "github.com/sirupsen/logrus"
	"github.com/tufin/espresso/shot/bq"
	"google.golang.org/api/iterator"
)

func runQuery(bqClient bq.Client, query string, params []bigquery.QueryParameter) (bq.Iterator, error) {

	queryIterator, err := bqClient.QueryIterator(query, params)
	if err != nil {
		log.Errorf("failed to run query with %v", err)
		fmt.Println("query:", query)
		return nil, err
	}

	return queryIterator, nil
}

func readResult(queryIterator bq.Iterator, row interface{}) ([]interface{}, error) {

	result := []interface{}{}

	for {
		rowCopy := reflect.New(reflect.ValueOf(row).Elem().Type()).Interface()

		err := queryIterator.Next(rowCopy)
		if err != nil {
			if err == iterator.Done {
				break
			}
			err = fmt.Errorf("failed to iterate with %v", err)
			return nil, err
		}
		result = append(result, rowCopy)
	}
	return result, nil
}
