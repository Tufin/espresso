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
	queryValues, resultValues, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", templates).RunTest("fruit", "Test1", []bigquery.QueryParameter{}, &map[string]bigquery.Value{})
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}

func TestEspressoShot_Filesystem(t *testing.T) {
	queryValues, resultValues, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", os.DirFS("./queries/fruit")).RunTest("fruit", "Test1", []bigquery.QueryParameter{}, &map[string]bigquery.Value{})
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}

func TestEspressoShot_FilesystemWithDepth(t *testing.T) {
	queryValues, resultValues, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", os.DirFS("./")).RunTest("fruit", "Test1", []bigquery.QueryParameter{}, &map[string]bigquery.Value{})
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}

func TestEspressoShot_Hierarchical(t *testing.T) {
	queryValues, resultValues, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", os.DirFS("./queries/fruit")).RunTest("fruit", "Hierarchical", []bigquery.QueryParameter{}, &map[string]bigquery.Value{})
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}

func TestEspressoShot_HierarchicalInline(t *testing.T) {
	queryValues, resultValues, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", templates).RunTest("fruit", "HierarchicalInline", []bigquery.QueryParameter{}, &map[string]bigquery.Value{})
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}

func TestEspressoShot_StructValue(t *testing.T) {
	row := struct {
		Fruit string
	}{}

	queryValues, resultValues, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", os.DirFS("./queries/fruit")).RunTest("fruit", "Test1", []bigquery.QueryParameter{}, &row)
	require.NoError(t, err)
	require.ElementsMatch(t, queryValues, resultValues)
}

func TestEspressoShot_InvalidRow(t *testing.T) {
	var row int
	_, _, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", templates).RunTest("fruit", "Test1", []bigquery.QueryParameter{}, &row)
	require.EqualError(t, err, "failed to iterate with bigquery: cannot convert *int to ValueLoader (need pointer to []Value, map[string]Value, or struct)")
}

func TestEspressoShot_Invalid(t *testing.T) {
	_, _, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", os.DirFS("./queries/invalid")).RunTest("invalid", "Test1", []bigquery.QueryParameter{}, &map[string]bigquery.Value{})
	require.EqualError(t, err, "invalid template \"invalid\" due to invalid arg \"Base\" lacks source, const and table definitions")
}

func TestEspressoShot_NoTemplate(t *testing.T) {
	_, _, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", os.DirFS("./")).RunTest("reuven", "Test1", []bigquery.QueryParameter{}, &map[string]bigquery.Value{})
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
	queryValues, resultValues, err := shot.NewShotWithClient(env.GetGCPProjectID(), "", templates).RunTest("fruit", "Mismatch", []bigquery.QueryParameter{}, &map[string]bigquery.Value{})
	require.NoError(t, err)
	require.Condition(t, func() bool {
		return !reflect.DeepEqual(queryValues, resultValues)
	})
}
