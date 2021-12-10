package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciOutofServiceFabricPath() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciOutofServiceFabricPathRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"annotation": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"pod_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"node_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"fex_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"interface": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		}),
	}
}

func dataSourceAciOutofServiceFabricPathRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	podId := fmt.Sprintf("%v", d.Get("pod_id"))
	nodeId := fmt.Sprintf("%v", d.Get("node_id"))
	interfaceName := d.Get("interface").(string)

	var interfaceDn string
	if FexId, ok := d.GetOk("fex_id"); ok {
		fexId := fmt.Sprintf("%v", FexId)
		interfaceDn = fmt.Sprintf("topology/pod-%s/paths-%s/extpaths-%s/pathep-[%s]", podId, nodeId, fexId, interfaceName)
	} else {
		interfaceDn = fmt.Sprintf("topology/pod-%s/paths-%s/pathep-[%s]", podId, nodeId, interfaceName)
	}

	dn := fmt.Sprintf(models.DnfabricRsOosPath, interfaceDn)
	fabricRsOosPath, err := getRemoteOutofServiceFabricPath(aciClient, dn)
	if err != nil {
		return err
	}
	d.SetId(dn)
	setOutofServiceFabricPathAttributes(fabricRsOosPath, d)

	return nil
}
