/*
 * Copyright 2021 Splunk, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"path/filepath"

	"terraform-provider-artifacts/internal/provider/internal/client"
)

func resourceUpload() *schema.Resource {
	return &schema.Resource{
		Description:   "Upload a file to Artifactory",
		CreateContext: resourceUploadCreate,
		ReadContext:   resourceUploadRead,
		UpdateContext: resourceUploadUpdate,
		DeleteContext: resourceUploadDelete,
		CustomizeDiff: resourceUploadDiff,
		Schema: map[string]*schema.Schema{
			uploadPathKey: {
				Description: "Path to upload to, relative to the provider's URL",
				Type:        schema.TypeString,
				Required:    true,
			},
			uploadFileKey: {
				Description: "File containing content to upload",
				Type:        schema.TypeString,
				Required:    true,
			},
			deleteOldPath: {
				Description: fmt.Sprintf("Set to false if the remote file should be orphaned on destruction of the resource or change of %s value. Defaults to true.", uploadPathKey),
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			sha1Key: {
				Description: "SHA1 of the uploaded file",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceUploadCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)

	d.SetId(artifactIDValue)

	uploadPath := d.Get(uploadPathKey).(string)
	filename := d.Get(uploadFileKey).(string)
	filePath, err := filepath.Abs(filename)
	if err != nil {
		return diag.Errorf("unable to determine absolute path for file %s", filePath)
	}

	if err := client.Upload(uploadPath, filePath); err != nil {
		return diag.Errorf("failure uploading file %s: %s", filePath, err)
	}

	return resourceUploadRead(ctx, d, meta)
}

func resourceUploadRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)

	uploadPath := d.Get(uploadPathKey).(string)

	checksums, err := client.Checksums(uploadPath)
	if err != nil {
		return diag.FromErr(err)
	}

	if checksums.SHA1 == "" {
		// missing SHA1 value indicates the resource wasn't found on the service, so mark this resource as missing
		d.SetId("")
	} else {
		if err := d.Set(sha1Key, checksums.SHA1); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceUploadUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)

	if d.HasChange(uploadPathKey) {
		if d.Get(deleteOldPath).(bool) {
			uploadPathInterfaceOld, _ := d.GetChange(uploadPathKey)
			uploadPathOld := uploadPathInterfaceOld.(string)
			if err := client.Delete(uploadPathOld); err != nil {
				return diag.Errorf("failure deleting old path %s: %s", uploadPathOld, err)
			}
		}
	}

	// after potentially deleting the old path, resourceUploadCreate does everything we need
	return resourceUploadCreate(ctx, d, meta)
}

func resourceUploadDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)

	if d.Get(deleteOldPath).(bool) {
		uploadPath := d.Get(uploadPathKey).(string)
		if err := client.Delete(uploadPath); err != nil {
			return diag.Errorf("error attempting delete: %s", err)
		}
	}

	return nil
}

func resourceUploadDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	client := meta.(*client.Client)

	uploadFile := d.Get(uploadFileKey).(string)

	sha1, err := client.SHA1(uploadFile)
	if err != nil {
		return err
	}

	if err := d.SetNew(sha1Key, sha1); err != nil {
		return err
	}

	return nil
}
