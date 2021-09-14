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
	queryValues, resultValues, err := shot.NewShot(env.GetGCPProjectID(), "", endpointTemplates).RunTest("report_summary", "Test1", []bigquery.QueryParameter{})
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}

func TestEspressoShot_Filesystem(t *testing.T) {
	queryValues, resultValues, err := shot.NewShot(env.GetGCPProjectID(), "", os.DirFS("./queries/endpoints")).RunTest("report_summary", "Test1", []bigquery.QueryParameter{})
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}

func TestEspressoShot_FilesystemWithDepth(t *testing.T) {
	queryValues, resultValues, err := shot.NewShot(env.GetGCPProjectID(), "", os.DirFS("./")).RunTest("report_summary", "Test1", []bigquery.QueryParameter{})
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}

func TestEspressoShot_Const(t *testing.T) {
	queryValues, resultValues, err := shot.NewShot(env.GetGCPProjectID(), "", os.DirFS("./queries/fruit")).RunTest("fruit", "Test1", []bigquery.QueryParameter{})
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}

func TestEspressoShot_Invalid(t *testing.T) {
	_, _, err := shot.NewShot(env.GetGCPProjectID(), "", os.DirFS("./queries/invalid")).RunTest("invalid", "Test1", []bigquery.QueryParameter{})
	require.EqualError(t, err, "invalid template \"invalid\" due to invalid arg \"Base\" lacks source, const and table definitions")
}

func TestEspressoShot_NoTemplate(t *testing.T) {
	_, _, err := shot.NewShot(env.GetGCPProjectID(), "", os.DirFS("./")).RunTest("reuven", "Test1", []bigquery.QueryParameter{})
	require.EqualError(t, err, "couldn't find definition file \"reuven.yaml\"")
}
