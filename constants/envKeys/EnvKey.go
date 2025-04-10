package envKeys

type EnvKey string

const (
	BucketName    EnvKey = "GCP_BUCKET_NAME"
	GCPAccessKey  EnvKey = "GCP_ACCESS_KEY"
	GCPPrivateKey EnvKey = "GCP_PRIVATE_KEY"
	GCPProjectID  EnvKey = "GCP_PROJECT_ID"
	StoragePath   EnvKey = "STORAGE_PATH"
	PORT          EnvKey = "PORT"
)
