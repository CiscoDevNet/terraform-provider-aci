package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciConcreteDevice() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciConcreteDeviceRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"l4_l7_devices_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vcenter_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vm_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"relation_vns_rs_c_dev_to_ctrlr_p": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to vmm:CtrlrP",
			},
		})),
	}
}

func dataSourceAciConcreteDeviceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	L4_L7DevicesDn := d.Get("l4_l7_devices_dn").(string)
	rn := fmt.Sprintf(models.RnvnsCDev, name)
	dn := fmt.Sprintf("%s/%s", L4_L7DevicesDn, rn)

	vnsCDev, err := getRemoteConcreteDevice(aciClient, dn)
	if err != nil {
		return nil
	}

	d.SetId(dn)

	_, err = setConcreteDeviceAttributes(vnsCDev, d)
	if err != nil {
		return nil
	}

	vnsRsCDevToCtrlrPData, err := aciClient.ReadRelationvnsRsCDevToCtrlrP(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsCDevToCtrlrP %v", err)
		d.Set("relation_vns_rs_c_dev_to_ctrlr_p", "")
	} else {
		d.Set("relation_vns_rs_c_dev_to_ctrlr_p", vnsRsCDevToCtrlrPData.(string))
	}

	return nil
}
