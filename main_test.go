package main

import (
	"embed"
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/tufin/espresso/env"
)

//go:embed queries/*.sql queries/*.yaml
var sqlTemplates embed.FS

func TestFromEspresso(t *testing.T) {

	// env.Ophiuchus()
	projectID := env.GetGCPProjectID()
	shot := NewShot(projectID, sqlTemplates)
	shot.Run(t, "report_summary", "Test1", []bigquery.QueryParameter{})
}
