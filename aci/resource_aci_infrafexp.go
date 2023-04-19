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

func resourceAciFEXProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciFEXProfileCreate,
		UpdateContext: resourceAciFEXProfileUpdate,
		ReadContext:   resourceAciFEXProfileRead,
		DeleteContext: resourceAciFEXProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciFEXProfileImport,
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

func getRemoteFEXProfile(client *client.Client, dn string) (*models.FEXProfile, error) {
	infraFexPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	infraFexP := models.FEXProfileFromContainer(infraFexPCont)

	if infraFexP.DistinguishedName == "" {
		return nil, fmt.Errorf("FEX Profile %s not found", dn)
	}

	return infraFexP, nil
}

func setFEXProfileAttributes(infraFexP *models.FEXProfile, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(infraFexP.DistinguishedName)
	d.Set("description", infraFexP.Description)
	infraFexPMap, err := infraFexP.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("name", infraFexPMap["name"])

	d.Set("annotation", infraFexPMap["annotation"])
	d.Set("name_alias", infraFexPMap["nameAlias"])
	return d, nil
}

func resourceAciFEXProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	infraFexP, err := getRemoteFEXProfile(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setFEXProfileAttributes(infraFexP, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciFEXProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] FEXProfile: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	infraFexPAttr := models.FEXProfileAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		infraFexPAttr.Annotation = Annotation.(string)
	} else {
		infraFexPAttr.Annotation = "{}"
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraFexPAttr.NameAlias = NameAlias.(string)
	}
	infraFexP := models.NewFEXProfile(fmt.Sprintf("infra/fexprof-%s", name), "uni", desc, infraFexPAttr)

	err := aciClient.Save(infraFexP)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(infraFexP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciFEXProfileRead(ctx, d, m)
}

func resourceAciFEXProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] FEXProfile: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	infraFexPAttr := models.FEXProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraFexPAttr.Annotation = Annotation.(string)
	} else {
		infraFexPAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraFexPAttr.NameAlias = NameAlias.(string)
	}
	infraFexP := models.NewFEXProfile(fmt.Sprintf("infra/fexprof-%s", name), "uni", desc, infraFexPAttr)

	infraFexP.Status = "modified"

	err := aciClient.Save(infraFexP)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(infraFexP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciFEXProfileRead(ctx, d, m)

}

func resourceAciFEXProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	infraFexP, err := getRemoteFEXProfile(aciClient, dn)

	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}
	_, err = setFEXProfileAttributes(infraFexP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciFEXProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "infraFexP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
