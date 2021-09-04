package env

const (
	KeyGCPProjectID = "GCLOUD_PROJECT_ID"
)

func GetGCPProjectID() string {

	return GetEnvOrExit(KeyGCPProjectID)
}
