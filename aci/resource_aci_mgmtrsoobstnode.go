package aci

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciMgmtStaticNode() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciMgmtStaticNodeCreate,
		UpdateContext: resourceAciMgmtStaticNodeUpdate,
		ReadContext:   resourceAciMgmtStaticNodeRead,
		DeleteContext: resourceAciMgmtStaticNodeDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciMgmtStaticNodeImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"management_epg_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"t_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"in_band",
					"out_of_band",
				}, false),
			},

			"addr": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"gw": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"v6_addr": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"v6_gw": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func getRemoteInbandStaticNode(client *client.Client, dn string) (*models.InbandStaticNode, error) {
	mgmtRsInBStNodeCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	mgmtRsInBStNode := models.InbandStaticNodeFromContainer(mgmtRsInBStNodeCont)

	if mgmtRsInBStNode.DistinguishedName == "" {
		return nil, fmt.Errorf("In Band Static Node %s not found", mgmtRsInBStNode.DistinguishedName)
	}

	return mgmtRsInBStNode, nil
}

func getRemoteOutofbandStaticNode(client *client.Client, dn string) (*models.OutofbandStaticNode, error) {
	mgmtRsOoBStNodeCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	mgmtRsOoBStNode := models.OutofbandStaticNodeFromContainer(mgmtRsOoBStNodeCont)

	if mgmtRsOoBStNode.DistinguishedName == "" {
		return nil, fmt.Errorf("Out of Band Static Node %s not found", mgmtRsOoBStNode.DistinguishedName)
	}

	return mgmtRsOoBStNode, nil
}

func setMgmtStaticNodeAttributes(mgmtRsOoBStNode *models.OutofbandStaticNode, mgmtRsInBStNode *models.InbandStaticNode, bandType string, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()

	if bandType == "in_band" {
		d.SetId(mgmtRsInBStNode.DistinguishedName)

		d.Set("description", mgmtRsInBStNode.Description)
		if dn != mgmtRsInBStNode.DistinguishedName {
			d.Set("management_epg_dn", "")
		}
	} else {
		d.SetId(mgmtRsOoBStNode.DistinguishedName)

		d.Set("description", mgmtRsOoBStNode.Description)
		if dn != mgmtRsOoBStNode.DistinguishedName {
			d.Set("management_epg_dn", "")
		}
	}

	if bandType == "in_band" {
		mgmtRsInBStNodeMap, err := mgmtRsInBStNode.ToMap()
		if err != nil {
			return d, err
		}

		d.Set("addr", mgmtRsInBStNodeMap["addr"])
		d.Set("annotation", mgmtRsInBStNodeMap["annotation"])
		d.Set("gw", mgmtRsInBStNodeMap["gw"])
		d.Set("t_dn", mgmtRsInBStNodeMap["tDn"])
		d.Set("v6_addr", mgmtRsInBStNodeMap["v6Addr"])
		d.Set("v6_gw", mgmtRsInBStNodeMap["v6Gw"])
		d.Set("management_epg_dn", GetParentDn(mgmtRsInBStNode.DistinguishedName, fmt.Sprintf("/rsinBStNode-[%s]", mgmtRsInBStNodeMap["tDn"])))

	} else {
		mgmtRsOoBStNodeMap, err := mgmtRsOoBStNode.ToMap()
		if err != nil {
			return d, err
		}

		d.Set("t_dn", mgmtRsOoBStNodeMap["tDn"])
		d.Set("addr", mgmtRsOoBStNodeMap["addr"])
		d.Set("annotation", mgmtRsOoBStNodeMap["annotation"])
		d.Set("gw", mgmtRsOoBStNodeMap["gw"])
		d.Set("v6_addr", mgmtRsOoBStNodeMap["v6Addr"])
		d.Set("v6_gw", mgmtRsOoBStNodeMap["v6Gw"])
		d.Set("management_epg_dn", GetParentDn(mgmtRsOoBStNode.DistinguishedName, fmt.Sprintf("/rsooBStNode-[%s]", mgmtRsOoBStNodeMap["tDn"])))

	}

	return d, nil
}

func resourceAciMgmtStaticNodeImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dns := strings.Split(d.Id(), ":")
	if len(dns) != 2 {
		return nil, fmt.Errorf("not getting enough arguments for the import operation")
	}

	dn := dns[0]
	bandType := dns[1]

	if bandType == "in_band" {
		mgmtRsInBStNode, err := getRemoteInbandStaticNode(aciClient, dn)

		if err != nil {
			return nil, err
		}
		schemaFilled, err := setMgmtStaticNodeAttributes(nil, mgmtRsInBStNode, "in_band", d)

		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

		return []*schema.ResourceData{schemaFilled}, nil

	} else {
		mgmtRsOoBStNode, err := getRemoteOutofbandStaticNode(aciClient, dn)

		if err != nil {
			return nil, err
		}
		schemaFilled, err := setMgmtStaticNodeAttributes(mgmtRsOoBStNode, nil, "out_of_band", d)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

		return []*schema.ResourceData{schemaFilled}, nil

	}

}

func resourceAciMgmtStaticNodeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Mgmt Static Node: Beginning Creation")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	tDn := d.Get("t_dn").(string)

	managementEPgDn := d.Get("management_epg_dn").(string)

	bandType := d.Get("type").(string)

	if bandType == "in_band" {
		mgmtRsInBStNodeAttr := models.InbandStaticNodeAttributes{}
		if Addr, ok := d.GetOk("addr"); ok {
			mgmtRsInBStNodeAttr.Addr = Addr.(string)
		}
		if Annotation, ok := d.GetOk("annotation"); ok {
			mgmtRsInBStNodeAttr.Annotation = Annotation.(string)
		} else {
			mgmtRsInBStNodeAttr.Annotation = "{}"
		}
		if Gw, ok := d.GetOk("gw"); ok {
			mgmtRsInBStNodeAttr.Gw = Gw.(string)
		}
		if V6Addr, ok := d.GetOk("v6_addr"); ok {
			mgmtRsInBStNodeAttr.V6Addr = V6Addr.(string)
		}
		if V6Gw, ok := d.GetOk("v6_gw"); ok {
			mgmtRsInBStNodeAttr.V6Gw = V6Gw.(string)
		}
		mgmtRsInBStNode := models.NewInbandStaticNode(fmt.Sprintf("rsinBStNode-[%s]", tDn), managementEPgDn, desc, mgmtRsInBStNodeAttr)

		err := aciClient.Save(mgmtRsInBStNode)
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(mgmtRsInBStNode.DistinguishedName)

	} else {
		mgmtRsOoBStNodeAttr := models.OutofbandStaticNodeAttributes{}

		if Addr, ok := d.GetOk("addr"); ok {
			mgmtRsOoBStNodeAttr.Addr = Addr.(string)
		}
		if Annotation, ok := d.GetOk("annotation"); ok {
			mgmtRsOoBStNodeAttr.Annotation = Annotation.(string)
		} else {
			mgmtRsOoBStNodeAttr.Annotation = "{}"
		}
		if Gw, ok := d.GetOk("gw"); ok {
			mgmtRsOoBStNodeAttr.Gw = Gw.(string)
		}
		if V6Addr, ok := d.GetOk("v6_addr"); ok {
			mgmtRsOoBStNodeAttr.V6Addr = V6Addr.(string)
		}
		if V6Gw, ok := d.GetOk("v6_gw"); ok {
			mgmtRsOoBStNodeAttr.V6Gw = V6Gw.(string)
		}

		mgmtRsOoBStNode := models.NewOutofbandStaticNode(fmt.Sprintf("rsooBStNode-[%s]", tDn), managementEPgDn, desc, mgmtRsOoBStNodeAttr)

		err := aciClient.Save(mgmtRsOoBStNode)
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(mgmtRsOoBStNode.DistinguishedName)
	}

	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciMgmtStaticNodeRead(ctx, d, m)
}

func resourceAciMgmtStaticNodeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Mgmt Static Node: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	tDn := d.Get("t_dn").(string)

	managementEPgDn := d.Get("management_epg_dn").(string)

	bandType := d.Get("type").(string)

	if bandType == "in_band" {
		mgmtRsInBStNodeAttr := models.InbandStaticNodeAttributes{}
		if Addr, ok := d.GetOk("addr"); ok {
			mgmtRsInBStNodeAttr.Addr = Addr.(string)
		}
		if Annotation, ok := d.GetOk("annotation"); ok {
			mgmtRsInBStNodeAttr.Annotation = Annotation.(string)
		} else {
			mgmtRsInBStNodeAttr.Annotation = "{}"
		}
		if Gw, ok := d.GetOk("gw"); ok {
			mgmtRsInBStNodeAttr.Gw = Gw.(string)
		}
		if V6Addr, ok := d.GetOk("v6_addr"); ok {
			mgmtRsInBStNodeAttr.V6Addr = V6Addr.(string)
		}
		if V6Gw, ok := d.GetOk("v6_gw"); ok {
			mgmtRsInBStNodeAttr.V6Gw = V6Gw.(string)
		}
		mgmtRsInBStNode := models.NewInbandStaticNode(fmt.Sprintf("rsinBStNode-[%s]", tDn), managementEPgDn, desc, mgmtRsInBStNodeAttr)

		mgmtRsInBStNode.Status = "modified"

		err := aciClient.Save(mgmtRsInBStNode)

		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(mgmtRsInBStNode.DistinguishedName)

	} else {
		mgmtRsOoBStNodeAttr := models.OutofbandStaticNodeAttributes{}
		if Addr, ok := d.GetOk("addr"); ok {
			mgmtRsOoBStNodeAttr.Addr = Addr.(string)
		}
		if Annotation, ok := d.GetOk("annotation"); ok {
			mgmtRsOoBStNodeAttr.Annotation = Annotation.(string)
		}
		if Gw, ok := d.GetOk("gw"); ok {
			mgmtRsOoBStNodeAttr.Gw = Gw.(string)
		}
		if V6Addr, ok := d.GetOk("v6_addr"); ok {
			mgmtRsOoBStNodeAttr.V6Addr = V6Addr.(string)
		}
		if V6Gw, ok := d.GetOk("v6_gw"); ok {
			mgmtRsOoBStNodeAttr.V6Gw = V6Gw.(string)
		}
		mgmtRsOoBStNode := models.NewOutofbandStaticNode(fmt.Sprintf("rsooBStNode-[%s]", tDn), managementEPgDn, desc, mgmtRsOoBStNodeAttr)

		mgmtRsOoBStNode.Status = "modified"

		err := aciClient.Save(mgmtRsOoBStNode)

		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(mgmtRsOoBStNode.DistinguishedName)
	}

	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciMgmtStaticNodeRead(ctx, d, m)

}

func resourceAciMgmtStaticNodeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()

	bandType := d.Get("type").(string)
	if bandType == "in_band" {
		mgmtRsInBStNode, err := getRemoteInbandStaticNode(aciClient, dn)

		if err != nil {
			d.SetId("")
			return nil
		}
		_, err = setMgmtStaticNodeAttributes(nil, mgmtRsInBStNode, "in_band", d)

		if err != nil {
			d.SetId("")
			return nil
		}

	} else {
		mgmtRsOoBStNode, err := getRemoteOutofbandStaticNode(aciClient, dn)

		if err != nil {
			d.SetId("")
			return nil
		}
		_, err = setMgmtStaticNodeAttributes(mgmtRsOoBStNode, nil, "out_of_band", d)

		if err != nil {
			d.SetId("")
			return nil
		}
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciMgmtStaticNodeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	bandType := d.Get("type").(string)

	if bandType == "in_band" {
		err := aciClient.DeleteByDn(dn, "mgmtRsInBStNode")
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		err := aciClient.DeleteByDn(dn, "mgmtRsOoBStNode")
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId("")

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	return nil
}
