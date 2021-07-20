---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "Artifacts Provider"
subcategory: ""
description: |-
  Terraform provider to manage artifacts on an Artifactory service
---

# Artifacts Provider

## Description

Terraform provider to manage artifacts on an Artifactory service

## Example Usage

```terraform
provider "artifacts" {
  url = "https://example.com/artifactory"
  # username and password to be set by environment variables
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **url** (String) URL of the Artifactory service

### Optional

- **username** (String) Username used to authenticate to Artifactory. May be set via the `ARTIFACTORY_AUTH_USERNAME` environment variable instead.
- **password** (String) Password used to authenticate to Artifactory. Must be set if username is set. May be set via the `ARTIFACTORY_AUTH_PASSWORD` environment variable instead.