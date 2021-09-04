package internal

import (
	"fmt"

	"cloud.google.com/go/bigquery"
	log "github.com/sirupsen/logrus"
	"github.com/tufin/espresso/bq"
	"github.com/tufin/espresso/common"
	"google.golang.org/api/iterator"
)

func printResult(queryIterator bq.Iterator) {

	result := map[string]bigquery.Value{}
	for {
		err := queryIterator.Next(&result)
		if err != nil {
			if err == iterator.Done {
				return
			}
			err = fmt.Errorf("failed to iterate on with '%v'", err)
			log.Error(err)
			continue
		}
		common.PrintPretty(result)
	}
}
