package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciFunctionNode() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciFunctionNodeRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"l4_l7_service_graph_template_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"func_template_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"func_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_copy": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"conn_copy_dn": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"managed": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"routing_mode": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"sequence_number": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"share_encap": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"relation_vns_rs_node_to_abs_func_prof": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"relation_vns_rs_node_to_l_dev": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"relation_vns_rs_node_to_m_func": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"relation_vns_rs_default_scope_to_term": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"relation_vns_rs_node_to_cloud_l_dev": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"l4_l7_device_interface_consumer_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"conn_consumer_dn": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"l4_l7_device_interface_consumer_connector_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"l4_l7_device_interface_consumer_attachment_notification": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"l4_l7_device_interface_provider_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"conn_provider_dn": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"l4_l7_device_interface_provider_connector_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"l4_l7_device_interface_provider_attachment_notification": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciFunctionNodeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("AbsNode-%s", name)
	l4l7ServiceGraphTemplateDn := d.Get("l4_l7_service_graph_template_dn").(string)

	dn := fmt.Sprintf("%s/%s", l4l7ServiceGraphTemplateDn, rn)

	vnsAbsNode, err := getRemoteFunctionNode(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setFunctionNodeAttributes(vnsAbsNode, d)
	if err != nil {
		return diag.FromErr(err)
	}

	err = getAndSetFunctionNodeRelationalAttributes(aciClient, dn, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
