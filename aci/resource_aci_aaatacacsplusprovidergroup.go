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

func resourceAciTACACSPlusProviderGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciTACACSPlusProviderGroupCreate,
		UpdateContext: resourceAciTACACSPlusProviderGroupUpdate,
		ReadContext:   resourceAciTACACSPlusProviderGroupRead,
		DeleteContext: resourceAciTACACSPlusProviderGroupDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciTACACSPlusProviderGroupImport,
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

func getRemoteTACACSPlusProviderGroup(client *client.Client, dn string) (*models.TACACSPlusProviderGroup, error) {
	aaaTacacsPlusProviderGroupCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaTacacsPlusProviderGroup := models.TACACSPlusProviderGroupFromContainer(aaaTacacsPlusProviderGroupCont)
	if aaaTacacsPlusProviderGroup.DistinguishedName == "" {
		return nil, fmt.Errorf("TACACSPlusProviderGroup %s not found", aaaTacacsPlusProviderGroup.DistinguishedName)
	}
	return aaaTacacsPlusProviderGroup, nil
}

func setTACACSPlusProviderGroupAttributes(aaaTacacsPlusProviderGroup *models.TACACSPlusProviderGroup, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(aaaTacacsPlusProviderGroup.DistinguishedName)
	d.Set("description", aaaTacacsPlusProviderGroup.Description)
	aaaTacacsPlusProviderGroupMap, err := aaaTacacsPlusProviderGroup.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("annotation", aaaTacacsPlusProviderGroupMap["annotation"])
	d.Set("name", aaaTacacsPlusProviderGroupMap["name"])
	d.Set("name_alias", aaaTacacsPlusProviderGroupMap["nameAlias"])
	return d, nil
}

func resourceAciTACACSPlusProviderGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaTacacsPlusProviderGroup, err := getRemoteTACACSPlusProviderGroup(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setTACACSPlusProviderGroupAttributes(aaaTacacsPlusProviderGroup, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciTACACSPlusProviderGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] TACACSPlusProviderGroup: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	aaaTacacsPlusProviderGroupAttr := models.TACACSPlusProviderGroupAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaTacacsPlusProviderGroupAttr.Annotation = Annotation.(string)
	} else {
		aaaTacacsPlusProviderGroupAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		aaaTacacsPlusProviderGroupAttr.Name = Name.(string)
	}
	aaaTacacsPlusProviderGroup := models.NewTACACSPlusProviderGroup(fmt.Sprintf("userext/tacacsext/tacacsplusprovidergroup-%s", name), "uni", desc, nameAlias, aaaTacacsPlusProviderGroupAttr)
	err := aciClient.Save(aaaTacacsPlusProviderGroup)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(aaaTacacsPlusProviderGroup.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciTACACSPlusProviderGroupRead(ctx, d, m)
}

func resourceAciTACACSPlusProviderGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] TACACSPlusProviderGroup: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	aaaTacacsPlusProviderGroupAttr := models.TACACSPlusProviderGroupAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaTacacsPlusProviderGroupAttr.Annotation = Annotation.(string)
	} else {
		aaaTacacsPlusProviderGroupAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		aaaTacacsPlusProviderGroupAttr.Name = Name.(string)
	}
	aaaTacacsPlusProviderGroup := models.NewTACACSPlusProviderGroup(fmt.Sprintf("userext/tacacsext/tacacsplusprovidergroup-%s", name), "uni", desc, nameAlias, aaaTacacsPlusProviderGroupAttr)
	aaaTacacsPlusProviderGroup.Status = "modified"
	err := aciClient.Save(aaaTacacsPlusProviderGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(aaaTacacsPlusProviderGroup.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciTACACSPlusProviderGroupRead(ctx, d, m)
}

func resourceAciTACACSPlusProviderGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaTacacsPlusProviderGroup, err := getRemoteTACACSPlusProviderGroup(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setTACACSPlusProviderGroupAttributes(aaaTacacsPlusProviderGroup, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciTACACSPlusProviderGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "aaaTacacsPlusProviderGroup")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
