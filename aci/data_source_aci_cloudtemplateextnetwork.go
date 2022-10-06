package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciCloudTemplateforExternalNetwork() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciCloudTemplateforExternalNetworkRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"infra_network_template_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"annotation": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"hub_network_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vrf_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func dataSourceAciCloudTemplateforExternalNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	CloudInfraNetworkTemplateDn := d.Get("infra_network_template_dn").(string)
	rn := fmt.Sprintf(models.RncloudtemplateExtNetwork, name)
	dn := fmt.Sprintf("%s/%s", CloudInfraNetworkTemplateDn, rn)

	cloudtemplateExtNetwork, err := getRemoteCloudTemplateforExternalNetwork(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)

	_, err = setCloudTemplateforExternalNetworkAttributes(cloudtemplateExtNetwork, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
