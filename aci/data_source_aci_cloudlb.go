package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciCloudL4L7LoadBalancer() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciCloudL4L7LoadBalancerRead,
		SchemaVersion: 1,
		Schema: AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"active_active": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"allow_all": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_scaling": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"context_aware": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_resource_group": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"device_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"function_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_count": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_copy": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_instantiation": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_static_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"l4l7_device_application_security_group": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"l4l7_third_party_device": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"managed": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"max_instance_count": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"min_instance_count": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"mode": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"native_lb_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"package_model": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"promiscuous_mode": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"scheme": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"sku": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_mode": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"trunking": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloud_l4l7_load_balancer_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"relation_cloud_rs_ldev_to_cloud_subnet": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Set:      schema.HashString,
			},
			"aaa_domain_dn": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Set:      schema.HashString,
			},
			"static_ip_address": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
		}),
	}
}

func dataSourceAciCloudL4L7LoadBalancerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	TenantDn := d.Get("tenant_dn").(string)
	rn := fmt.Sprintf(RnCloudLB, name)
	dn := fmt.Sprintf("%s/%s", TenantDn, rn)

	_, err := getAndSetRemoteCloudL4L7LoadBalancerAttributes(aciClient, dn, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
