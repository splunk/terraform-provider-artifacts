package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"terraform-provider-artifacts/internal/provider/internal/client"
)

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		client := client.Client{
			URL:     d.Get(urlKey).(string),
			Context: ctx,
		}

		if usernameInterface, ok := d.GetOk(usernameKey); ok {
			client.Username = usernameInterface.(string)
			// we're counting on RequiredWith functionality to be valid, so set Password if Username is given
			client.Password = d.Get(passwordKey).(string)
		}

		return &client, nil
	}
}

// New returns a function that returns a pointer to a new schema.Provider for this provider.
func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				urlKey: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "URL of the artifacts service",
				},
				usernameKey: {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc(usernameEnvKey, nil),
					Description: "Username to use for authentication to the artifacts service",
				},
				passwordKey: {
					Type:         schema.TypeString,
					Optional:     true,
					DefaultFunc:  schema.EnvDefaultFunc(passwordEnvKey, nil),
					RequiredWith: []string{usernameKey},
					Description:  fmt.Sprintf("Password to use for authentication to the artifacts service. Must be set if %s is set.", usernameKey),
				},
			},
			ResourcesMap: map[string]*schema.Resource{
				uploadResourceKey: resourceUpload(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}