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

func resourceAciAccessGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciAccessAccessGroupCreate,
		UpdateContext: resourceAciAccessAccessGroupUpdate,
		ReadContext:   resourceAciAccessAccessGroupRead,
		DeleteContext: resourceAciAccessAccessGroupDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciAccessAccessGroupImport,
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"access_port_selector_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"fex_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tdn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DefaultFunc: func() (interface{}, error) {
					return "orchestrator:terraform", nil
				},
			},
		},
	}
}
func getRemoteAccessAccessGroup(client *client.Client, dn string) (*models.AccessAccessGroup, error) {
	infraRsAccBaseGrpCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	infraRsAccBaseGrp := models.AccessAccessGroupFromContainer(infraRsAccBaseGrpCont)

	if infraRsAccBaseGrp.DistinguishedName == "" {
		return nil, fmt.Errorf("AccessAccessGroup %s not found", infraRsAccBaseGrp.DistinguishedName)
	}

	return infraRsAccBaseGrp, nil
}

func setAccessAccessGroupAttributes(infraRsAccBaseGrp *models.AccessAccessGroup, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(infraRsAccBaseGrp.DistinguishedName)
	if dn != infraRsAccBaseGrp.DistinguishedName {
		d.Set("access_port_selector_dn", "")
	}
	infraRsAccBaseGrpMap, err := infraRsAccBaseGrp.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("access_port_selector_dn", GetParentDn(dn, fmt.Sprintf("/rsaccBaseGrp")))
	d.Set("annotation", infraRsAccBaseGrpMap["annotation"])
	d.Set("fex_id", infraRsAccBaseGrpMap["fexId"])
	d.Set("tdn", infraRsAccBaseGrpMap["tDn"])
	return d, nil
}

func resourceAciAccessAccessGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	infraRsAccBaseGrp, err := getRemoteAccessAccessGroup(aciClient, dn)

	if err != nil {
		return nil, err
	}
	pDN := GetParentDn(dn, fmt.Sprintf("/rsaccBaseGrp"))
	d.Set("access_port_selector_dn", pDN)
	schemaFilled, err := setAccessAccessGroupAttributes(infraRsAccBaseGrp, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciAccessAccessGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] AccessAccessGroup: Beginning Creation")
	aciClient := m.(*client.Client)
	AccessPortSelectorDn := d.Get("access_port_selector_dn").(string)

	infraRsAccBaseGrpAttr := models.AccessAccessGroupAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraRsAccBaseGrpAttr.Annotation = Annotation.(string)
	} else {
		infraRsAccBaseGrpAttr.Annotation = "{}"
	}
	if FexId, ok := d.GetOk("fex_id"); ok {
		infraRsAccBaseGrpAttr.FexId = FexId.(string)
	}
	if TDn, ok := d.GetOk("tdn"); ok {
		infraRsAccBaseGrpAttr.TDn = TDn.(string)
	}
	infraRsAccBaseGrp := models.NewAccessAccessGroup(fmt.Sprintf("rsaccBaseGrp"), AccessPortSelectorDn, infraRsAccBaseGrpAttr)

	err := aciClient.Save(infraRsAccBaseGrp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(infraRsAccBaseGrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciAccessAccessGroupRead(ctx, d, m)
}

func resourceAciAccessAccessGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] AccessAccessGroup: Beginning Update")

	aciClient := m.(*client.Client)

	AccessPortSelectorDn := d.Get("access_port_selector_dn").(string)

	infraRsAccBaseGrpAttr := models.AccessAccessGroupAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraRsAccBaseGrpAttr.Annotation = Annotation.(string)
	} else {
		infraRsAccBaseGrpAttr.Annotation = "{}"
	}
	if FexId, ok := d.GetOk("fex_id"); ok {
		infraRsAccBaseGrpAttr.FexId = FexId.(string)
	}
	if TDn, ok := d.GetOk("tdn"); ok {
		infraRsAccBaseGrpAttr.TDn = TDn.(string)
	}
	infraRsAccBaseGrp := models.NewAccessAccessGroup(fmt.Sprintf("rsaccBaseGrp"), AccessPortSelectorDn, infraRsAccBaseGrpAttr)

	infraRsAccBaseGrp.Status = "modified"

	err := aciClient.Save(infraRsAccBaseGrp)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(infraRsAccBaseGrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciAccessAccessGroupRead(ctx, d, m)

}

func resourceAciAccessAccessGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	infraRsAccBaseGrp, err := getRemoteAccessAccessGroup(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setAccessAccessGroupAttributes(infraRsAccBaseGrp, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciAccessAccessGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "infraRsAccBaseGrp")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
