package shot_test

import (
	"embed"
	"os"
	"reflect"
	"strings"
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/stretchr/testify/require"
	"github.com/tufin/espresso/env"
	"github.com/tufin/espresso/shot"
)

//go:embed queries/fruit
var templates embed.FS

func TestEspressoShot_Embed(t *testing.T) {
	queryValues, resultValues, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", templates).GetTestResults("fruit", "Test1", []bigquery.QueryParameter{}, &map[string]bigquery.Value{})
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}

func TestEspressoShot_Filesystem(t *testing.T) {
	queryValues, resultValues, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", os.DirFS("./queries/fruit")).GetTestResults("fruit", "Test1", []bigquery.QueryParameter{}, &map[string]bigquery.Value{})
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}

func TestEspressoShot_FilesystemWithDepth(t *testing.T) {
	queryValues, resultValues, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", os.DirFS("./")).GetTestResults("fruit", "Test1", []bigquery.QueryParameter{}, &map[string]bigquery.Value{})
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}

func TestEspressoShot_Hierarchical(t *testing.T) {
	queryValues, resultValues, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", os.DirFS("./queries/fruit")).GetTestResults("fruit", "Hierarchical", []bigquery.QueryParameter{}, &map[string]bigquery.Value{})
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}

func TestEspressoShot_HierarchicalInline(t *testing.T) {
	queryValues, resultValues, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", templates).GetTestResults("fruit", "HierarchicalInline", []bigquery.QueryParameter{}, &map[string]bigquery.Value{})
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}

func TestEspressoShot_StructValue(t *testing.T) {
	row := struct {
		Fruit string
	}{}

	queryValues, resultValues, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", os.DirFS("./queries/fruit")).GetTestResults("fruit", "Test1", []bigquery.QueryParameter{}, &row)
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}

func TestEspressoShot_InvalidRow(t *testing.T) {
	var row int
	_, _, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", templates).GetTestResults("fruit", "Test1", []bigquery.QueryParameter{}, &row)
	require.EqualError(t, err, "failed to iterate with bigquery: cannot convert *int to ValueLoader (need pointer to []Value, map[string]Value, or struct)")
}

func TestEspressoShot_Invalid(t *testing.T) {
	_, _, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", os.DirFS("./queries/invalid")).GetTestResults("invalid", "Test1", []bigquery.QueryParameter{}, &map[string]bigquery.Value{})
	require.EqualError(t, err, "invalid template \"invalid\" due to invalid arg \"Base\" lacks source, const and table definitions")
}

func TestEspressoShot_NoTemplate(t *testing.T) {
	_, _, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", os.DirFS("./")).GetTestResults("reuven", "Test1", []bigquery.QueryParameter{}, &map[string]bigquery.Value{})
	require.EqualError(t, err, "couldn't find definition file \"reuven.yaml\"")
}

func TestGetQuery(t *testing.T) {
	query, err := shot.NewShot(nil, "", "", os.DirFS("./queries/fruit")).GetQuery("fruit", "Test1", []bigquery.QueryParameter{})
	require.NoError(t, err)
	require.Equal(t,
		"\n\nWITH base AS (\n    \n\nSELECT\n    \"orange\" AS fruit\nUNION ALL\nSELECT\n    \"apple\"\n\n\n)\nSELECT\n    fruit\nFROM base\n\n",
		strings.ReplaceAll(query, "\r", ""))
}

func TestEspressoShot_Mismatch(t *testing.T) {
	queryValues, resultValues, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", templates).GetTestResults("fruit", "Mismatch", []bigquery.QueryParameter{}, &map[string]bigquery.Value{})
	require.NoError(t, err)
	require.Condition(t, func() bool {
		return !reflect.DeepEqual(queryValues, resultValues)
	})
}

func TestEspressoShot_NoResult(t *testing.T) {
	_, _, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", templates).GetTestResults("fruit", "NoResult", []bigquery.QueryParameter{}, &map[string]bigquery.Value{})
	require.EqualError(t, err, "result source is missing")
}

func TestEspressoShot_RunTest(t *testing.T) {
	env.Ophiuchus()
	empty, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", templates).RunTest("fruit", "Test1", []bigquery.QueryParameter{})
	require.NoError(t, err)
	require.True(t, empty)
}
