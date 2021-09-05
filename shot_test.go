package espresso_test

import (
	"embed"
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/tufin/espresso"
	"github.com/tufin/espresso/env"
)

//go:embed queries/*
var sqlTemplates embed.FS

func TestFromEspresso(t *testing.T) {
	projectID := env.GetGCPProjectID()
	shot := espresso.NewShot(projectID, sqlTemplates)
	shot.RunTest(t, "queries", "report_summary", "Test1", []bigquery.QueryParameter{})
}
