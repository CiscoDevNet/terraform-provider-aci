package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciCloudL4L7Device() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciCloudL4L7DeviceRead,
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
			"mode": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"package_model": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"promiscuous_mode": &schema.Schema{
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
			"relation_cloud_rs_ldev_to_ctx": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"aaa_domain_dn": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Set:      schema.HashString,
			},
			"interface_selectors": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: AppendAttrSchemas(map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"allow_all": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_point_selectors": &schema.Schema{
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: AppendAttrSchemas(map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"match_expression": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								}),
							},
						},
					}),
				},
			},
		}),
	}
}

func dataSourceAciCloudL4L7DeviceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	TenantDn := d.Get("tenant_dn").(string)
	rn := fmt.Sprintf(RnCloudLDev, name)
	dn := fmt.Sprintf("%s/%s", TenantDn, rn)

	_, err := getAndSetRemoteCloudL4L7DeviceAttributes(aciClient, dn, d)
	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}
	return nil
}
