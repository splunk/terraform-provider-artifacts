// Copyright 2021 Splunk, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
