package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"terraform-provider-artifacts/internal/provider"
)

const (
	version = "dev"
)

func main() {
	opts := &plugin.ServeOpts{ProviderFunc: provider.New(version)}

	plugin.Serve(opts)
}
