resource "artifacts_upload" "latest_only" {
  upload_path = "uploaded/latest.txt"
  upload_file = "./artifact.txt"
}

resource "artifacts_upload" "retain_versions" {
  upload_path = "upload/version-1.0.0.txt"
  upload_file = "./artifact.txt"
  // the uploaded file at upload/version-1.0.0.txt will remain in place even if this resource is removed or
  // the upload_path changes, such as if the version is incremented. this can be useful if you wish to retain
  // previous versions even when a newer version is uploaded.
  delete_old_path = false
}
