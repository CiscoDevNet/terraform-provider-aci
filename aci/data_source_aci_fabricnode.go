package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciFabricNodeOrg() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciFabricNodeReadOrg,

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
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

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"node_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"role": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
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
	_, err = setFabricNodeAttributesOrg(fabricNode, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func getRemoteFabricNodeOrg(client *client.Client, dn string) (*models.TopologyFabricNode, error) {
	fabricNodeCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	fabricNode := models.TopologyFabricNodeFromContainer(fabricNodeCont)

	if fabricNode.DistinguishedName == "" {
		return nil, fmt.Errorf("FabricNode %s not found", fabricNode.DistinguishedName)
	}

	return fabricNode, nil
}

func setFabricNodeAttributesOrg(fabricNode *models.TopologyFabricNode, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(fabricNode.DistinguishedName)

	dn := d.Id()
	if dn != fabricNode.DistinguishedName {
		d.Set("fabric_pod_dn", "")
	}
	fabricNodeMap, err := fabricNode.ToMap()
	if err != nil {
		return nil, err
	}

	d.Set("fabric_pod_dn", GetParentDn(dn, fmt.Sprintf("/node-%s", fabricNodeMap["id"])))

	d.Set("fabric_node_id", fabricNodeMap["id"])

	d.Set("ad_st", fabricNodeMap["adSt"])
	d.Set("annotation", fabricNodeMap["annotation"])
	d.Set("apic_type", fabricNodeMap["apicType"])
	d.Set("fabric_st", fabricNodeMap["fabricSt"])
	d.Set("fabric_node_id", fabricNodeMap["id"])
	d.Set("name_alias", fabricNodeMap["nameAlias"])
	d.Set("address", fabricNodeMap["address"])
	d.Set("name", fabricNodeMap["name"])
	d.Set("node_type", fabricNodeMap["nodeType"])
	d.Set("role", fabricNodeMap["role"])
	return d, nil
}
