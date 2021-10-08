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

func TestGetTestResults_Embed(t *testing.T) {
	queryValues, resultValues, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", templates).GetTestResults("fruit", "Test1", []bigquery.QueryParameter{}, &map[string]bigquery.Value{})
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}

func TestGetTestResults_Filesystem(t *testing.T) {
	queryValues, resultValues, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", os.DirFS("./queries/fruit")).GetTestResults("fruit", "Test1", []bigquery.QueryParameter{}, &map[string]bigquery.Value{})
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}

func TestGetTestResults_FilesystemWithDepth(t *testing.T) {
	queryValues, resultValues, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", os.DirFS("./")).GetTestResults("fruit", "Test1", []bigquery.QueryParameter{}, &map[string]bigquery.Value{})
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}

func TestGetTestResults_Hierarchical(t *testing.T) {
	queryValues, resultValues, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", os.DirFS("./queries/fruit")).GetTestResults("fruit", "Hierarchical", []bigquery.QueryParameter{}, &map[string]bigquery.Value{})
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}

func TestGetTestResults_HierarchicalInline(t *testing.T) {
	queryValues, resultValues, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", templates).GetTestResults("fruit", "HierarchicalInline", []bigquery.QueryParameter{}, &map[string]bigquery.Value{})
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}

func TestGetTestResults_StructValue(t *testing.T) {
	row := struct {
		Fruit string
	}{}

	queryValues, resultValues, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", os.DirFS("./queries/fruit")).GetTestResults("fruit", "Test1", []bigquery.QueryParameter{}, &row)
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}

func TestGetTestResults_InvalidRow(t *testing.T) {
	var row int
	_, _, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", templates).GetTestResults("fruit", "Test1", []bigquery.QueryParameter{}, &row)
	require.EqualError(t, err, "failed to iterate with bigquery: cannot convert *int to ValueLoader (need pointer to []Value, map[string]Value, or struct)")
}

func TestGetTestResults_Invalid(t *testing.T) {
	_, _, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", os.DirFS("./queries/invalid")).GetTestResults("invalid", "Test1", []bigquery.QueryParameter{}, &map[string]bigquery.Value{})
	require.EqualError(t, err, "invalid template \"invalid\" due to invalid arg \"Base\" lacks source, const and table definitions")
}

func TestGetTestResults_NoTemplate(t *testing.T) {
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

func TestGetTestResults_Mismatch(t *testing.T) {
	queryValues, resultValues, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", templates).GetTestResults("fruit", "Mismatch", []bigquery.QueryParameter{}, &map[string]bigquery.Value{})
	require.NoError(t, err)
	require.Condition(t, func() bool {
		return !reflect.DeepEqual(queryValues, resultValues)
	})
}

func TestGetTestResults_NoResult(t *testing.T) {
	_, _, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", templates).GetTestResults("fruit", "NoResult", []bigquery.QueryParameter{}, &map[string]bigquery.Value{})
	require.EqualError(t, err, "result source is missing")
}

func TestRunTest_Identical(t *testing.T) {
	empty, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", templates).RunTest("fruit", "Test1", []bigquery.QueryParameter{})
	require.NoError(t, err)
	require.True(t, empty)
}

func TestRunTest_Mismatch(t *testing.T) {
	empty, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", templates).RunTest("fruit", "Mismatch", []bigquery.QueryParameter{})
	require.NoError(t, err)
	require.False(t, empty)
}

func TestRunTest_Duplicates(t *testing.T) {
	empty, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", templates).RunTest("fruit", "Duplicates", []bigquery.QueryParameter{})
	require.NoError(t, err)
	// TODO: this test should fail because the results are not equal
	require.True(t, empty)
}

func TestGetTestResults_Duplicates(t *testing.T) {
	queryValues, resultValues, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", templates).GetTestResults("fruit", "Duplicates", []bigquery.QueryParameter{}, &map[string]bigquery.Value{})
	require.NoError(t, err)

	require.Condition(t, func() bool {
		return !reflect.DeepEqual(queryValues, resultValues)
	})
}
