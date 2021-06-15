package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceUpload(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceUploadCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("artifacts_upload.test", "upload_path", "sas-binary/terraform-provider-artifacts-test/test_file_1.txt"),
					resource.TestCheckResourceAttr("artifacts_upload.test", "sha1", "af3d968c42b3046f86296c7522b3b20dfdc58c59"),
				),
			},
			{
				Config: testResourceUploadUpdateURL,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("artifacts_upload.test", "upload_path", "sas-binary/terraform-provider-artifacts-test/test_file_2.txt"),
					resource.TestCheckResourceAttr("artifacts_upload.test", "sha1", "af3d968c42b3046f86296c7522b3b20dfdc58c59"),
				),
			},
			{
				Config: testResourceUploadUpdateFile,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("artifacts_upload.test", "upload_path", "sas-binary/terraform-provider-artifacts-test/test_file_2.txt"),
					resource.TestCheckResourceAttr("artifacts_upload.test", "sha1", "27f1703d965438b9f78d412d60d47816d878c9a5"),
				),
			},
		},
	})
}

const testResourceUploadCreateConfig = `
provider "artifacts" {
  url = "https://repo.splunk.com/artifactory"
}

resource "artifacts_upload" "test" {
  upload_path = "sas-binary/terraform-provider-artifacts-test/test_file_1.txt"
  upload_file = "test_files/source_file.txt"
}
`

const testResourceUploadUpdateURL = `
provider "artifacts" {
  url = "https://repo.splunk.com/artifactory"
}

resource "artifacts_upload" "test" {
  upload_path = "sas-binary/terraform-provider-artifacts-test/test_file_2.txt"
  upload_file = "test_files/source_file.txt"
}
`

const testResourceUploadUpdateFile = `
provider "artifacts" {
  url = "https://repo.splunk.com/artifactory"
}

resource "artifacts_upload" "test" {
  upload_path = "sas-binary/terraform-provider-artifacts-test/test_file_2.txt"
  upload_file = "test_files/source_file_update.txt"
}
`
