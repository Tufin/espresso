package shot_test

import (
	"embed"
	"os"
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/stretchr/testify/require"
	"github.com/tufin/espresso/env"
	"github.com/tufin/espresso/shot"
)

//go:embed queries/endpoints
var endpointTemplates embed.FS

func TestEspressoShot_Embed(t *testing.T) {
	queryValues, resultValues, err := shot.NewShot(env.GetGCPProjectID(), endpointTemplates).RunTest("queries/endpoints/report_summary.yaml", "report_summary", "Test1", []bigquery.QueryParameter{})
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}

func TestEspressoShot_Filesystem(t *testing.T) {
	queryValues, resultValues, err := shot.NewShot(env.GetGCPProjectID(), os.DirFS("./queries/endpoints")).RunTest("report_summary.yaml", "report_summary", "Test1", []bigquery.QueryParameter{})
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}

func TestEspressoShot_FilesystemWithDepth(t *testing.T) {
	queryValues, resultValues, err := shot.NewShot(env.GetGCPProjectID(), os.DirFS("./")).RunTest("queries/endpoints/report_summary.yaml", "report_summary", "Test1", []bigquery.QueryParameter{})
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}

func TestEspressoShot_TemplateNotFound(t *testing.T) {
	_, _, err := shot.NewShot(env.GetGCPProjectID(), os.DirFS("./queries/endpoints")).RunTest("report_summary.yaml", "xxx", "Test1", []bigquery.QueryParameter{})
	require.Error(t, err)
}

func TestEspressoShot_Const(t *testing.T) {
	queryValues, resultValues, err := shot.NewShot(env.GetGCPProjectID(), os.DirFS("./queries/fruit")).RunTest("fruit.yaml", "fruit", "Test1", []bigquery.QueryParameter{})
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}
