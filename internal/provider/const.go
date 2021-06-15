package provider

const (
	urlKey = "url"
	// we don't track uploaded artifacts by an ID, just use any value for Id
	artifactIDValue   = "artifacts_id_value"
	usernameKey       = "username"
	usernameEnvKey    = "ARTIFACTORY_AUTH_USERNAME"
	passwordKey       = "password"
	passwordEnvKey    = "ARTIFACTORY_AUTH_PASSWORD"
	uploadResourceKey = "artifacts_upload"
	uploadPathKey     = "upload_path"
	uploadFileKey     = "upload_file"
	deleteOldPath     = "delete_old_path"
	sha1Key           = "sha1"
)
