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

func resourceAciMulticastAddressPool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciMulticastAddressPoolCreate,
		UpdateContext: resourceAciMulticastAddressPoolUpdate,
		ReadContext:   resourceAciMulticastAddressPoolRead,
		DeleteContext: resourceAciMulticastAddressPoolDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciMulticastAddressPoolImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		})),
	}
}

func getRemoteMulticastAddressPool(client *client.Client, dn string) (*models.MulticastAddressPool, error) {
	fvnsMcastAddrInstPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	fvnsMcastAddrInstP := models.MulticastAddressPoolFromContainer(fvnsMcastAddrInstPCont)
	if fvnsMcastAddrInstP.DistinguishedName == "" {
		return nil, fmt.Errorf("MulticastAddressPool %s not found", dn)
	}
	return fvnsMcastAddrInstP, nil
}

func setMulticastAddressPoolAttributes(fvnsMcastAddrInstP *models.MulticastAddressPool, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(fvnsMcastAddrInstP.DistinguishedName)
	d.Set("description", fvnsMcastAddrInstP.Description)
	fvnsMcastAddrInstPMap, err := fvnsMcastAddrInstP.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("annotation", fvnsMcastAddrInstPMap["annotation"])
	d.Set("name", fvnsMcastAddrInstPMap["name"])
	d.Set("name_alias", fvnsMcastAddrInstPMap["nameAlias"])
	return d, nil
}

func resourceAciMulticastAddressPoolImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	fvnsMcastAddrInstP, err := getRemoteMulticastAddressPool(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setMulticastAddressPoolAttributes(fvnsMcastAddrInstP, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciMulticastAddressPoolCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] MulticastAddressPool: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	fvnsMcastAddrInstPAttr := models.MulticastAddressPoolAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		fvnsMcastAddrInstPAttr.Annotation = Annotation.(string)
	} else {
		fvnsMcastAddrInstPAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		fvnsMcastAddrInstPAttr.Name = Name.(string)
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fvnsMcastAddrInstPAttr.NameAlias = NameAlias.(string)
	}
	fvnsMcastAddrInstP := models.NewMulticastAddressPool(fmt.Sprintf("infra/maddrns-%s", name), "uni", desc, fvnsMcastAddrInstPAttr)
	err := aciClient.Save(fvnsMcastAddrInstP)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fvnsMcastAddrInstP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciMulticastAddressPoolRead(ctx, d, m)
}
func resourceAciMulticastAddressPoolUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] MulticastAddressPool: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)

	fvnsMcastAddrInstPAttr := models.MulticastAddressPoolAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		fvnsMcastAddrInstPAttr.Annotation = Annotation.(string)
	} else {
		fvnsMcastAddrInstPAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		fvnsMcastAddrInstPAttr.Name = Name.(string)
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fvnsMcastAddrInstPAttr.NameAlias = NameAlias.(string)
	}
	fvnsMcastAddrInstP := models.NewMulticastAddressPool(fmt.Sprintf("infra/maddrns-%s", name), "uni", desc, fvnsMcastAddrInstPAttr)
	fvnsMcastAddrInstP.Status = "modified"

	err := aciClient.Save(fvnsMcastAddrInstP)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fvnsMcastAddrInstP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciMulticastAddressPoolRead(ctx, d, m)
}

func resourceAciMulticastAddressPoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	fvnsMcastAddrInstP, err := getRemoteMulticastAddressPool(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}

	_, err = setMulticastAddressPoolAttributes(fvnsMcastAddrInstP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciMulticastAddressPoolDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "fvnsMcastAddrInstP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
