package main

import (
	"flag"
	"os"
	"reflect"

	"cloud.google.com/go/bigquery"
	log "github.com/sirupsen/logrus"
	"github.com/tufin/espresso/env"
	"github.com/tufin/espresso/shot"
)

var dir, query, test string

func init() {
	flag.StringVar(&dir, "dir", "", "base dir containing SQL files and definition file")
	flag.StringVar(&query, "query", "", "query name")
	flag.StringVar(&test, "test", "", "test name")
}

func main() {
	flag.Parse()

	queryValues, resultValues, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", os.DirFS(dir)).GetTestResults(query, test, []bigquery.QueryParameter{}, &map[string]bigquery.Value{})
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	if !reflect.DeepEqual(queryValues, resultValues) {
		println("test failed")
		os.Exit(1)
	}

	println("test passed")
}
