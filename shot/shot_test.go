package shot_test

import (
	"embed"
	"os"
	"strings"
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/stretchr/testify/require"
	"github.com/tufin/espresso/env"
	"github.com/tufin/espresso/shot"
	"github.com/tufin/espresso/shot/bq"
)

//go:embed queries/endpoints
var endpointTemplates embed.FS

func TestEspressoShot_Embed(t *testing.T) {
	project := env.GetGCPProjectID()
	queryValues, resultValues, err := shot.NewShot(bq.NewClient(project), project, "", endpointTemplates).RunTest("report_summary", "Test1", []bigquery.QueryParameter{}, &map[string]bigquery.Value{})
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}

func TestEspressoShot_Filesystem(t *testing.T) {
	project := env.GetGCPProjectID()
	queryValues, resultValues, err := shot.NewShot(bq.NewClient(project), project, "", os.DirFS("./queries/endpoints")).RunTest("report_summary", "Test1", []bigquery.QueryParameter{}, &map[string]bigquery.Value{})
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}

func TestEspressoShot_FilesystemWithDepth(t *testing.T) {
	project := env.GetGCPProjectID()
	queryValues, resultValues, err := shot.NewShot(bq.NewClient(project), project, "", os.DirFS("./")).RunTest("report_summary", "Test1", []bigquery.QueryParameter{}, &map[string]bigquery.Value{})
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}

func TestEspressoShot_HierarchicalTemplates(t *testing.T) {
	project := env.GetGCPProjectID()
	queryValues, resultValues, err := shot.NewShot(bq.NewClient(project), project, "", endpointTemplates).RunTest("report_summary", "TestHierarchical", []bigquery.QueryParameter{}, &map[string]bigquery.Value{})
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}

func TestEspressoShot_Const(t *testing.T) {
	project := env.GetGCPProjectID()
	queryValues, resultValues, err := shot.NewShot(bq.NewClient(project), project, "", os.DirFS("./queries/fruit")).RunTest("fruit", "Test1", []bigquery.QueryParameter{}, &map[string]bigquery.Value{})
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}

func TestEspressoShot_StructValue(t *testing.T) {
	row := struct {
		Fruit string
	}{}

	project := env.GetGCPProjectID()
	queryValues, resultValues, err := shot.NewShot(bq.NewClient(project), project, "", os.DirFS("./queries/fruit")).RunTest("fruit", "Test1", []bigquery.QueryParameter{}, &row)
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}

func TestEspressoShot_Invalid(t *testing.T) {
	project := env.GetGCPProjectID()
	_, _, err := shot.NewShot(bq.NewClient(project), project, "", os.DirFS("./queries/invalid")).RunTest("invalid", "Test1", []bigquery.QueryParameter{}, &map[string]bigquery.Value{})
	require.EqualError(t, err, "invalid template \"invalid\" due to invalid arg \"Base\" lacks source, const and table definitions")
}

func TestEspressoShot_NoTemplate(t *testing.T) {
	project := env.GetGCPProjectID()
	_, _, err := shot.NewShot(bq.NewClient(project), project, "", os.DirFS("./")).RunTest("reuven", "Test1", []bigquery.QueryParameter{}, &map[string]bigquery.Value{})
	require.EqualError(t, err, "couldn't find definition file \"reuven.yaml\"")
}

func TestGetQuery(t *testing.T) {
	query, err := shot.NewShot(nil, "", "", os.DirFS("./queries/fruit")).GetQuery("fruit", "Test1", []bigquery.QueryParameter{})
	require.NoError(t, err)
	require.Equal(t,
		"\n\nWITH base AS (\n    SELECT\n        \"orange\" AS fruit\n    UNION ALL\n    SELECT\n        \"apple\"\n)\nSELECT\n    fruit\nFROM base\n\n",
		strings.ReplaceAll(query, "\r", ""))
}
