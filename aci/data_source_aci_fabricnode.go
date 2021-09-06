package aci

import (
	"context"
	"fmt"
	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciFabricNodeOrg() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciFabricNodeReadOrg,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"fabric_pod_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"fabric_node_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"ad_st": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"apic_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"fabric_st": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"last_state_mod_ts": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciFabricNodeReadOrg(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	fabric_node_id := d.Get("fabric_node_id").(string)

	rn := fmt.Sprintf("node-%s", fabric_node_id)
	FabricPodDn := d.Get("fabric_pod_dn").(string)

	dn := fmt.Sprintf("%s/%s", FabricPodDn, rn)

	fabricNode, err := getRemoteFabricNodeOrg(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	setFabricNodeAttributesOrg(fabricNode, d)
	return nil
}

func getRemoteFabricNodeOrg(client *client.Client, dn string) (*models.OrgFabricNode, error) {
	fabricNodeCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	fabricNode := models.OrgFabricNodeFromContainer(fabricNodeCont)

	if fabricNode.DistinguishedName == "" {
		return nil, fmt.Errorf("FabricNode %s not found", fabricNode.DistinguishedName)
	}

	return fabricNode, nil
}

func setFabricNodeAttributesOrg(fabricNode *models.OrgFabricNode, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(fabricNode.DistinguishedName)
	d.Set("description", fabricNode.Description)
	dn := d.Id()
	if dn != fabricNode.DistinguishedName {
		d.Set("fabric_pod_dn", "")
	}
	fabricNodeMap, _ := fabricNode.ToMap()

	d.Set("fabric_pod_dn", GetParentDn(dn, fmt.Sprintf("/node-%s", fabricNodeMap["id"])))

	d.Set("fabric_node_id", fabricNodeMap["id"])

	d.Set("ad_st", fabricNodeMap["adSt"])
	d.Set("annotation", fabricNodeMap["annotation"])
	d.Set("apic_type", fabricNodeMap["apicType"])
	d.Set("fabric_st", fabricNodeMap["fabricSt"])
	d.Set("fabric_node_id", fabricNodeMap["id"])
	d.Set("last_state_mod_ts", fabricNodeMap["lastStateModTs"])
	d.Set("name_alias", fabricNodeMap["nameAlias"])
	return d
}
