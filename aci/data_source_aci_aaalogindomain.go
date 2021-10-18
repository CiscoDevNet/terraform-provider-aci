package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciLoginDomain() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciLoginDomainReadContext,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"domain_auth_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"provider_group": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"realm": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"realm_sub_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func dataSourceAciLoginDomainReadContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)

	rn := fmt.Sprintf("userext/logindomain-%s", name)
	dn := fmt.Sprintf("uni/%s", rn)
	childDn := fmt.Sprintf("%s/domainauth", dn)
	aaaLoginDomain, err := getRemoteLoginDomain(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setLoginDomainAttributes(aaaLoginDomain, d)
	if err != nil {
		return diag.FromErr(err)
	}

	aaaDomainAuth, err := getRemoteAuthenticationMethodfortheDomain(aciClient, childDn)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setAuthenticationMethodfortheDomainAttributes(aaaDomainAuth, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
