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

func resourceAciSpineInterfaceProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciSpineInterfaceProfileCreate,
		UpdateContext: resourceAciSpineInterfaceProfileUpdate,
		ReadContext:   resourceAciSpineInterfaceProfileRead,
		DeleteContext: resourceAciSpineInterfaceProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciSpineInterfaceProfileImport,
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
func getRemoteSpineInterfaceProfile(client *client.Client, dn string) (*models.SpineInterfaceProfile, error) {
	infraSpAccPortPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	infraSpAccPortP := models.SpineInterfaceProfileFromContainer(infraSpAccPortPCont)

	if infraSpAccPortP.DistinguishedName == "" {
		return nil, fmt.Errorf("SpineInterfaceProfile %s not found", infraSpAccPortP.DistinguishedName)
	}

	return infraSpAccPortP, nil
}

func setSpineInterfaceProfileAttributes(infraSpAccPortP *models.SpineInterfaceProfile, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(infraSpAccPortP.DistinguishedName)
	d.Set("description", infraSpAccPortP.Description)
	infraSpAccPortPMap, err := infraSpAccPortP.ToMap()

	if err != nil {
		return d, err
	}

	d.Set("name", infraSpAccPortPMap["name"])

	d.Set("annotation", infraSpAccPortPMap["annotation"])
	d.Set("name_alias", infraSpAccPortPMap["nameAlias"])
	return d, nil
}

func resourceAciSpineInterfaceProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	infraSpAccPortP, err := getRemoteSpineInterfaceProfile(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setSpineInterfaceProfileAttributes(infraSpAccPortP, d)

	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciSpineInterfaceProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SpineInterfaceProfile: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	infraSpAccPortPAttr := models.SpineInterfaceProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraSpAccPortPAttr.Annotation = Annotation.(string)
	} else {
		infraSpAccPortPAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraSpAccPortPAttr.NameAlias = NameAlias.(string)
	}
	infraSpAccPortP := models.NewSpineInterfaceProfile(fmt.Sprintf("infra/spaccportprof-%s", name), "uni", desc, infraSpAccPortPAttr)

	err := aciClient.Save(infraSpAccPortP)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(infraSpAccPortP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciSpineInterfaceProfileRead(ctx, d, m)
}

func resourceAciSpineInterfaceProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SpineInterfaceProfile: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	infraSpAccPortPAttr := models.SpineInterfaceProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraSpAccPortPAttr.Annotation = Annotation.(string)
	} else {
		infraSpAccPortPAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraSpAccPortPAttr.NameAlias = NameAlias.(string)
	}
	infraSpAccPortP := models.NewSpineInterfaceProfile(fmt.Sprintf("infra/spaccportprof-%s", name), "uni", desc, infraSpAccPortPAttr)

	infraSpAccPortP.Status = "modified"

	err := aciClient.Save(infraSpAccPortP)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(infraSpAccPortP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciSpineInterfaceProfileRead(ctx, d, m)

}

func resourceAciSpineInterfaceProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	infraSpAccPortP, err := getRemoteSpineInterfaceProfile(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setSpineInterfaceProfileAttributes(infraSpAccPortP, d)

	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciSpineInterfaceProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "infraSpAccPortP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
