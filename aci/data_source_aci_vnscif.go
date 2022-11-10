package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciConcreteInterface() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciConcreteInterfaceRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"concrete_device_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"encap": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vnic_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"relation_vns_rs_c_if_path_att": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to fabric:PathEp",
			},
		})),
	}
}

func dataSourceAciConcreteInterfaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	ConcreteDeviceDn := d.Get("concrete_device_dn").(string)
	rn := fmt.Sprintf(models.RnvnsCIf, name)
	dn := fmt.Sprintf("%s/%s", ConcreteDeviceDn, rn)

	vnsCIf, err := getRemoteConcreteInterface(aciClient, dn)
	if err != nil {
		return nil
	}

	d.SetId(dn)

	_, err = setConcreteInterfaceAttributes(vnsCIf, d)
	if err != nil {
		return nil
	}

	vnsRsCIfPathAttData, err := aciClient.ReadRelationvnsRsCIfPathAtt(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsCIfPathAtt %v", err)
		d.Set("relation_vns_rs_c_if_path_att", "")
	} else {
		d.Set("relation_vns_rs_c_if_path_att", vnsRsCIfPathAttData.(string))
	}

	return nil
}
