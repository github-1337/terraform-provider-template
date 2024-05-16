package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return &schema.Provider{
				DataSourcesMap: map[string]*schema.Resource{
					"idpfingerprint": dataSourceIDPFingerprint(),
				},
			}
		},
	})
}

func dataSourceIDPFingerprint() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIDPFingerprintRead,
		Schema: map[string]*schema.Schema{
			"idp_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
