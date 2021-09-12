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

var dir, def, test string

func init() {
	flag.StringVar(&dir, "dir", "", "base dir containing SQL files and definition file")
	flag.StringVar(&def, "def", "", "relative path of test definition file")
	flag.StringVar(&test, "test", "", "test name")
}

func main() {
	flag.Parse()

	fileSystem := os.DirFS(dir)
	queryValues, resultValues, err := shot.NewShot(env.GetGCPProjectID(), fileSystem).RunTest(def, test, []bigquery.QueryParameter{})
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	if !reflect.DeepEqual(queryValues, resultValues) {
		log.Info("test failed")
		os.Exit(1)
	}
}