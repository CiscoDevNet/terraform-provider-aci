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

func resourceAciVXLANPool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciVXLANPoolCreate,
		UpdateContext: resourceAciVXLANPoolUpdate,
		ReadContext:   resourceAciVXLANPoolRead,
		DeleteContext: resourceAciVXLANPoolDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciVXLANPoolImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}
func getRemoteVXLANPool(client *client.Client, dn string) (*models.VXLANPool, error) {
	fvnsVxlanInstPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	fvnsVxlanInstP := models.VXLANPoolFromContainer(fvnsVxlanInstPCont)

	if fvnsVxlanInstP.DistinguishedName == "" {
		return nil, fmt.Errorf("VXLANPool %s not found", fvnsVxlanInstP.DistinguishedName)
	}

	return fvnsVxlanInstP, nil
}

func setVXLANPoolAttributes(fvnsVxlanInstP *models.VXLANPool, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(fvnsVxlanInstP.DistinguishedName)
	d.Set("description", fvnsVxlanInstP.Description)
	fvnsVxlanInstPMap, _ := fvnsVxlanInstP.ToMap()

	d.Set("name", fvnsVxlanInstPMap["name"])

	d.Set("annotation", fvnsVxlanInstPMap["annotation"])
	d.Set("name_alias", fvnsVxlanInstPMap["nameAlias"])
	return d
}

func resourceAciVXLANPoolImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	fvnsVxlanInstP, err := getRemoteVXLANPool(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setVXLANPoolAttributes(fvnsVxlanInstP, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciVXLANPoolCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] VXLANPool: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	fvnsVxlanInstPAttr := models.VXLANPoolAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvnsVxlanInstPAttr.Annotation = Annotation.(string)
	} else {
		fvnsVxlanInstPAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fvnsVxlanInstPAttr.NameAlias = NameAlias.(string)
	}
	fvnsVxlanInstP := models.NewVXLANPool(fmt.Sprintf("infra/vxlanns-%s", name), "uni", desc, fvnsVxlanInstPAttr)

	err := aciClient.Save(fvnsVxlanInstP)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fvnsVxlanInstP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciVXLANPoolRead(ctx, d, m)
}

func resourceAciVXLANPoolUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] VXLANPool: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	fvnsVxlanInstPAttr := models.VXLANPoolAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvnsVxlanInstPAttr.Annotation = Annotation.(string)
	} else {
		fvnsVxlanInstPAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fvnsVxlanInstPAttr.NameAlias = NameAlias.(string)
	}
	fvnsVxlanInstP := models.NewVXLANPool(fmt.Sprintf("infra/vxlanns-%s", name), "uni", desc, fvnsVxlanInstPAttr)

	fvnsVxlanInstP.Status = "modified"

	err := aciClient.Save(fvnsVxlanInstP)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fvnsVxlanInstP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciVXLANPoolRead(ctx, d, m)

}

func resourceAciVXLANPoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	fvnsVxlanInstP, err := getRemoteVXLANPool(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setVXLANPoolAttributes(fvnsVxlanInstP, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciVXLANPoolDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fvnsVxlanInstP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
