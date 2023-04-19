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

func resourceAciSAMLProviderGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciSAMLProviderGroupCreate,
		UpdateContext: resourceAciSAMLProviderGroupUpdate,
		ReadContext:   resourceAciSAMLProviderGroupRead,
		DeleteContext: resourceAciSAMLProviderGroupDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciSAMLProviderGroupImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		})),
	}
}

func getRemoteSAMLProviderGroup(client *client.Client, dn string) (*models.SAMLProviderGroup, error) {
	aaaSamlProviderGroupCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaSamlProviderGroup := models.SAMLProviderGroupFromContainer(aaaSamlProviderGroupCont)
	if aaaSamlProviderGroup.DistinguishedName == "" {
		return nil, fmt.Errorf("SAML Provider Group %s not found", dn)
	}
	return aaaSamlProviderGroup, nil
}

func setSAMLProviderGroupAttributes(aaaSamlProviderGroup *models.SAMLProviderGroup, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(aaaSamlProviderGroup.DistinguishedName)
	d.Set("description", aaaSamlProviderGroup.Description)
	aaaSamlProviderGroupMap, err := aaaSamlProviderGroup.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("annotation", aaaSamlProviderGroupMap["annotation"])
	d.Set("name", aaaSamlProviderGroupMap["name"])
	d.Set("name_alias", aaaSamlProviderGroupMap["nameAlias"])
	return d, nil
}

func resourceAciSAMLProviderGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaSamlProviderGroup, err := getRemoteSAMLProviderGroup(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setSAMLProviderGroupAttributes(aaaSamlProviderGroup, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciSAMLProviderGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SAMLProviderGroup: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	aaaSamlProviderGroupAttr := models.SAMLProviderGroupAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaSamlProviderGroupAttr.Annotation = Annotation.(string)
	} else {
		aaaSamlProviderGroupAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		aaaSamlProviderGroupAttr.Name = Name.(string)
	}
	aaaSamlProviderGroup := models.NewSAMLProviderGroup(fmt.Sprintf("userext/samlext/samlprovidergroup-%s", name), "uni", desc, nameAlias, aaaSamlProviderGroupAttr)
	err := aciClient.Save(aaaSamlProviderGroup)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(aaaSamlProviderGroup.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciSAMLProviderGroupRead(ctx, d, m)
}

func resourceAciSAMLProviderGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SAMLProviderGroup: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	aaaSamlProviderGroupAttr := models.SAMLProviderGroupAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaSamlProviderGroupAttr.Annotation = Annotation.(string)
	} else {
		aaaSamlProviderGroupAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		aaaSamlProviderGroupAttr.Name = Name.(string)
	}
	aaaSamlProviderGroup := models.NewSAMLProviderGroup(fmt.Sprintf("userext/samlext/samlprovidergroup-%s", name), "uni", desc, nameAlias, aaaSamlProviderGroupAttr)
	aaaSamlProviderGroup.Status = "modified"
	err := aciClient.Save(aaaSamlProviderGroup)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(aaaSamlProviderGroup.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciSAMLProviderGroupRead(ctx, d, m)
}

func resourceAciSAMLProviderGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaSamlProviderGroup, err := getRemoteSAMLProviderGroup(aciClient, dn)
	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}
	_, err = setSAMLProviderGroupAttributes(aaaSamlProviderGroup, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciSAMLProviderGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "aaaSamlProviderGroup")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
