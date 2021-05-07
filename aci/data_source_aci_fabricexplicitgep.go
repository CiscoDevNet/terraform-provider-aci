package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciVPCExplicitProtectionGroup() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciVPCExplicitProtectionGroupRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"vpc_explicit_protection_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciVPCExplicitProtectionGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("fabric/protpol/expgep-%s", name)

	dn := fmt.Sprintf("uni/%s", rn)

	fabricExplicitGEp, err := getRemoteVPCExplicitProtectionGroupDS(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	setVPCExplicitProtectionGroupAttributesDS(fabricExplicitGEp, d)
	return nil
}

func getRemoteVPCExplicitProtectionGroupDS(client *client.Client, dn string) (*models.VPCExplicitProtectionGroup, error) {
	baseurlStr := "/api/node/mo"
	dnUrl := fmt.Sprintf("%s/%s.json?rsp-subtree=children", baseurlStr, dn)
	fabricExplicitGEpCont, err := client.GetViaURL(dnUrl)
	if err != nil {
		return nil, err
	}

	fabricExplicitGEp := models.VPCExplicitProtectionGroupFromContainer(fabricExplicitGEpCont)

	if fabricExplicitGEp.DistinguishedName == "" {
		return nil, fmt.Errorf("VPCExplicitProtectionGroup %s not found", fabricExplicitGEp.DistinguishedName)
	}

	return fabricExplicitGEp, nil
}

func setVPCExplicitProtectionGroupAttributesDS(fabricExplicitGEp *models.VPCExplicitProtectionGroup, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(fabricExplicitGEp.DistinguishedName)
	d.Set("description", fabricExplicitGEp.Description)
	fabricExplicitGEpMap, _ := fabricExplicitGEp.ToMap()

	d.Set("name", fabricExplicitGEpMap["name"])

	d.Set("annotation", fabricExplicitGEpMap["annotation"])
	d.Set("vpc_explicit_protection_group_id", fabricExplicitGEpMap["id"])
	return d
}
