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
	queryValues, resultValues, err := shot.NewShot(env.GetGCPProjectID(), "", endpointTemplates).RunTest("queries/endpoints/report_summary.yaml", "Test1", []bigquery.QueryParameter{})
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}

func TestEspressoShot_Filesystem(t *testing.T) {
	queryValues, resultValues, err := shot.NewShot(env.GetGCPProjectID(), "", os.DirFS("./queries/endpoints")).RunTest("report_summary.yaml", "Test1", []bigquery.QueryParameter{})
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}

func TestEspressoShot_FilesystemWithDepth(t *testing.T) {
	queryValues, resultValues, err := shot.NewShot(env.GetGCPProjectID(), "", os.DirFS("./")).RunTest("queries/endpoints/report_summary.yaml", "Test1", []bigquery.QueryParameter{})
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}

func TestEspressoShot_Const(t *testing.T) {
	queryValues, resultValues, err := shot.NewShot(env.GetGCPProjectID(), "", os.DirFS("./queries/fruit")).RunTest("fruit.yaml", "Test1", []bigquery.QueryParameter{})
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}
