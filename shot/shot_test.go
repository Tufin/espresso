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

//go:embed queries
var sqlTemplates embed.FS

func TestEspressoShot_Embed(t *testing.T) {
	queryValues, resultValues, err := shot.NewShot(env.GetGCPProjectID(), sqlTemplates).RunTest("queries", "report_summary", "Test1", []bigquery.QueryParameter{})
	require.NoError(t, err)
	require.Equal(t, queryValues, resultValues)
}

func TestEspressoShot_Filesystem(t *testing.T) {
	fileSystem := os.DirFS(".")
	queryValues, resultValues, err := shot.NewShot(env.GetGCPProjectID(), fileSystem).RunTest("queries", "report_summary", "Test1", []bigquery.QueryParameter{})
	require.NoError(t, err)
	require.Equal(t, queryValues, resultValues)
}
