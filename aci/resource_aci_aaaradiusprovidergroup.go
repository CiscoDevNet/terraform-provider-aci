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

func resourceAciRADIUSProviderGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciRADIUSProviderGroupCreate,
		UpdateContext: resourceAciRADIUSProviderGroupUpdate,
		ReadContext:   resourceAciRADIUSProviderGroupRead,
		DeleteContext: resourceAciRADIUSProviderGroupDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciRADIUSProviderGroupImport,
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

func getRemoteRADIUSProviderGroup(client *client.Client, dn string) (*models.RADIUSProviderGroup, error) {
	aaaRadiusProviderGroupCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaRadiusProviderGroup := models.RADIUSProviderGroupFromContainer(aaaRadiusProviderGroupCont)
	if aaaRadiusProviderGroup.DistinguishedName == "" {
		return nil, fmt.Errorf("RADIUSProviderGroup %s not found", aaaRadiusProviderGroup.DistinguishedName)
	}
	return aaaRadiusProviderGroup, nil
}

func setRADIUSProviderGroupAttributes(aaaRadiusProviderGroup *models.RADIUSProviderGroup, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(aaaRadiusProviderGroup.DistinguishedName)
	d.Set("description", aaaRadiusProviderGroup.Description)
	aaaRadiusProviderGroupMap, err := aaaRadiusProviderGroup.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("annotation", aaaRadiusProviderGroupMap["annotation"])
	d.Set("name", aaaRadiusProviderGroupMap["name"])
	d.Set("name_alias", aaaRadiusProviderGroupMap["nameAlias"])
	return d, nil
}

func resourceAciRADIUSProviderGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaRadiusProviderGroup, err := getRemoteRADIUSProviderGroup(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setRADIUSProviderGroupAttributes(aaaRadiusProviderGroup, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciRADIUSProviderGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] RADIUSProviderGroup: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	aaaRadiusProviderGroupAttr := models.RADIUSProviderGroupAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaRadiusProviderGroupAttr.Annotation = Annotation.(string)
	} else {
		aaaRadiusProviderGroupAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		aaaRadiusProviderGroupAttr.Name = Name.(string)
	}
	aaaRadiusProviderGroup := models.NewRADIUSProviderGroup(fmt.Sprintf("userext/radiusext/radiusprovidergroup-%s", name), "uni", desc, nameAlias, aaaRadiusProviderGroupAttr)
	err := aciClient.Save(aaaRadiusProviderGroup)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(aaaRadiusProviderGroup.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciRADIUSProviderGroupRead(ctx, d, m)
}

func resourceAciRADIUSProviderGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] RADIUSProviderGroup: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	aaaRadiusProviderGroupAttr := models.RADIUSProviderGroupAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaRadiusProviderGroupAttr.Annotation = Annotation.(string)
	} else {
		aaaRadiusProviderGroupAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		aaaRadiusProviderGroupAttr.Name = Name.(string)
	}
	aaaRadiusProviderGroup := models.NewRADIUSProviderGroup(fmt.Sprintf("userext/radiusext/radiusprovidergroup-%s", name), "uni", desc, nameAlias, aaaRadiusProviderGroupAttr)
	aaaRadiusProviderGroup.Status = "modified"
	err := aciClient.Save(aaaRadiusProviderGroup)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(aaaRadiusProviderGroup.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciRADIUSProviderGroupRead(ctx, d, m)
}

func resourceAciRADIUSProviderGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaRadiusProviderGroup, err := getRemoteRADIUSProviderGroup(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setRADIUSProviderGroupAttributes(aaaRadiusProviderGroup, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciRADIUSProviderGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "aaaRadiusProviderGroup")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
