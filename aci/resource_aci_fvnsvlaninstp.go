package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciVLANPool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciVLANPoolCreate,
		UpdateContext: resourceAciVLANPoolUpdate,
		ReadContext:   resourceAciVLANPoolRead,
		DeleteContext: resourceAciVLANPoolDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciVLANPoolImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"alloc_mode": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"dynamic",
					"static",
				}, false),
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}
func getRemoteVLANPool(client *client.Client, dn string) (*models.VLANPool, error) {
	fvnsVlanInstPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	fvnsVlanInstP := models.VLANPoolFromContainer(fvnsVlanInstPCont)

	if fvnsVlanInstP.DistinguishedName == "" {
		return nil, fmt.Errorf("VLAN Pool %s not found", dn)
	}

	return fvnsVlanInstP, nil
}

func setVLANPoolAttributes(fvnsVlanInstP *models.VLANPool, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(fvnsVlanInstP.DistinguishedName)
	d.Set("description", fvnsVlanInstP.Description)
	fvnsVlanInstPMap, _ := fvnsVlanInstP.ToMap()

	d.Set("name", fvnsVlanInstPMap["name"])

	d.Set("alloc_mode", fvnsVlanInstPMap["allocMode"])
	d.Set("annotation", fvnsVlanInstPMap["annotation"])
	d.Set("name_alias", fvnsVlanInstPMap["nameAlias"])
	return d
}

func resourceAciVLANPoolImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	fvnsVlanInstP, err := getRemoteVLANPool(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setVLANPoolAttributes(fvnsVlanInstP, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciVLANPoolCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] VLANPool: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	allocmode := d.Get("alloc_mode").(string)

	fvnsVlanInstPAttr := models.VLANPoolAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvnsVlanInstPAttr.Annotation = Annotation.(string)
	} else {
		fvnsVlanInstPAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fvnsVlanInstPAttr.NameAlias = NameAlias.(string)
	}
	fvnsVlanInstP := models.NewVLANPool(fmt.Sprintf("infra/vlanns-[%s]-%s", name, allocmode), "uni", desc, fvnsVlanInstPAttr)

	err := aciClient.Save(fvnsVlanInstP)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fvnsVlanInstP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciVLANPoolRead(ctx, d, m)
}

func resourceAciVLANPoolUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] VLANPool: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	allocmode := d.Get("alloc_mode").(string)

	fvnsVlanInstPAttr := models.VLANPoolAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvnsVlanInstPAttr.Annotation = Annotation.(string)
	} else {
		fvnsVlanInstPAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fvnsVlanInstPAttr.NameAlias = NameAlias.(string)
	}
	fvnsVlanInstP := models.NewVLANPool(fmt.Sprintf("infra/vlanns-[%s]-%s", name, allocmode), "uni", desc, fvnsVlanInstPAttr)

	fvnsVlanInstP.Status = "modified"

	err := aciClient.Save(fvnsVlanInstP)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fvnsVlanInstP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciVLANPoolRead(ctx, d, m)

}

func resourceAciVLANPoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	fvnsVlanInstP, err := getRemoteVLANPool(aciClient, dn)

	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}
	setVLANPoolAttributes(fvnsVlanInstP, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciVLANPoolDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fvnsVlanInstP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
