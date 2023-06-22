package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciCloudSubnet() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciCloudSubnetRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"cloud_cidr_pool_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"ip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"scope": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Use the sorted scope list to handle identical changes",
			},

			"usage": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Only applicable to the AWS vendor",
			},
			"relation_cloud_rs_subnet_to_flow_log": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"relation_cloud_rs_subnet_to_ctx": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_group_label": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Only applicable to the GCP vendor",
			},
		}),
	}
}

func dataSourceAciCloudSubnetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	ip := d.Get("ip").(string)

	rn := fmt.Sprintf("subnet-[%s]", ip)
	CloudCIDRPoolDn := d.Get("cloud_cidr_pool_dn").(string)

	dn := fmt.Sprintf("%s/%s", CloudCIDRPoolDn, rn)

	cloudSubnet, err := getRemoteCloudSubnet(aciClient, dn, d)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setCloudSubnetAttributes(cloudSubnet, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
