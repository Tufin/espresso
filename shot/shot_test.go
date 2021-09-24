package shot_test

import (
	"embed"
	"os"
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
	query, err := shot.NewShot(nil, "", "", os.DirFS("./queries/endpoints")).GetQuery("report_summary", "Test1", []bigquery.QueryParameter{})
	require.NoError(t, err)
	require.Equal(t,
		"\n\nWITH base AS (\n    SELECT\n        request_method,\n        path,\n        SUM(hit_count_yesterday) AS hit_count_yesterday,\n    FROM \n\n(\n    SELECT * FROM \n\n(\n    SELECT\n        \"GET\" AS request_method,\n        \"/api/rome/conf/master/tenants\" AS path,\n        200 AS status_code,\n        2 AS hit_count_yesterday,\n    UNION ALL\n    SELECT\n        \"GET\" AS request_method,\n        \"/api/rome/conf/master/tenants\" AS path,\n        500 AS status_code,\n        1 AS hit_count_yesterday,\n)\n\n\n)\n\n\n    WHERE status_code<400\n    GROUP BY \n        request_method,\n        path\n)\nSELECT\n    COUNT(*) AS total_endpoints,\n    COUNTIF(hit_count_yesterday>0) AS total_endpoints_yesterday,\n    SUM(hit_count_yesterday) AS hit_count_yesterday,\nFROM base\n\n",
		query)
}
