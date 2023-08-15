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

func resourceAciPrivateLinkLabelfortheserviceEPg() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciPrivateLinkLabelfortheserviceEPgCreate,
		UpdateContext: resourceAciPrivateLinkLabelfortheserviceEPgUpdate,
		ReadContext:   resourceAciPrivateLinkLabelfortheserviceEPgRead,
		DeleteContext: resourceAciPrivateLinkLabelfortheserviceEPgDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciPrivateLinkLabelfortheserviceEPgImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"cloud_service_epg_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		})),
	}
}

func getRemotePrivateLinkLabelfortheserviceEPg(client *client.Client, dn string) (*models.PrivateLinkLabelfortheserviceEPg, error) {
	cloudPrivateLinkLabelCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	cloudPrivateLinkLabel := models.PrivateLinkLabelfortheserviceEPgFromContainer(cloudPrivateLinkLabelCont)
	if cloudPrivateLinkLabel.DistinguishedName == "" {
		return nil, fmt.Errorf("PrivateLinkLabelfortheserviceEPg %s not found", dn)
	}
	return cloudPrivateLinkLabel, nil
}

func setPrivateLinkLabelfortheserviceEPgAttributes(cloudPrivateLinkLabel *models.PrivateLinkLabelfortheserviceEPg, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(cloudPrivateLinkLabel.DistinguishedName)
	d.Set("description", cloudPrivateLinkLabel.Description)
	cloudPrivateLinkLabelMap, err := cloudPrivateLinkLabel.ToMap()
	if err != nil {
		return d, err
	}
	dn := d.Id()
	if dn != cloudPrivateLinkLabel.DistinguishedName {
		d.Set("cloud_service_epg_dn", "")
	} else {
		d.Set("cloud_service_epg_dn", GetParentDn(cloudPrivateLinkLabel.DistinguishedName, fmt.Sprintf("/"+models.RnCloudPrivateLinkLabel, cloudPrivateLinkLabelMap["name"])))
	}
	d.Set("annotation", cloudPrivateLinkLabelMap["annotation"])
	d.Set("name", cloudPrivateLinkLabelMap["name"])
	d.Set("name_alias", cloudPrivateLinkLabelMap["nameAlias"])
	return d, nil
}

func resourceAciPrivateLinkLabelfortheserviceEPgImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	cloudPrivateLinkLabel, err := getRemotePrivateLinkLabelfortheserviceEPg(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setPrivateLinkLabelfortheserviceEPgAttributes(cloudPrivateLinkLabel, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciPrivateLinkLabelfortheserviceEPgCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] PrivateLinkLabelfortheserviceEPg: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	CloudServiceEPgDn := d.Get("cloud_service_epg_dn").(string)

	cloudPrivateLinkLabelAttr := models.PrivateLinkLabelfortheserviceEPgAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudPrivateLinkLabelAttr.Annotation = Annotation.(string)
	} else {
		cloudPrivateLinkLabelAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		cloudPrivateLinkLabelAttr.Name = Name.(string)
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		cloudPrivateLinkLabelAttr.NameAlias = NameAlias.(string)
	}
	cloudPrivateLinkLabel := models.NewPrivateLinkLabelfortheserviceEPg(fmt.Sprintf(models.RnCloudPrivateLinkLabel, name), CloudServiceEPgDn, desc, cloudPrivateLinkLabelAttr)

	err := aciClient.Save(cloudPrivateLinkLabel)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cloudPrivateLinkLabel.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciPrivateLinkLabelfortheserviceEPgRead(ctx, d, m)
}
func resourceAciPrivateLinkLabelfortheserviceEPgUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] PrivateLinkLabelfortheserviceEPg: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	CloudServiceEPgDn := d.Get("cloud_service_epg_dn").(string)

	cloudPrivateLinkLabelAttr := models.PrivateLinkLabelfortheserviceEPgAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudPrivateLinkLabelAttr.Annotation = Annotation.(string)
	} else {
		cloudPrivateLinkLabelAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		cloudPrivateLinkLabelAttr.Name = Name.(string)
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		cloudPrivateLinkLabelAttr.NameAlias = NameAlias.(string)
	}
	cloudPrivateLinkLabel := models.NewPrivateLinkLabelfortheserviceEPg(fmt.Sprintf(models.RnCloudPrivateLinkLabel, name), CloudServiceEPgDn, desc, cloudPrivateLinkLabelAttr)

	cloudPrivateLinkLabel.Status = "modified"

	err := aciClient.Save(cloudPrivateLinkLabel)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cloudPrivateLinkLabel.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciPrivateLinkLabelfortheserviceEPgRead(ctx, d, m)
}

func resourceAciPrivateLinkLabelfortheserviceEPgRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	cloudPrivateLinkLabel, err := getRemotePrivateLinkLabelfortheserviceEPg(aciClient, dn)
	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}

	_, err = setPrivateLinkLabelfortheserviceEPgAttributes(cloudPrivateLinkLabel, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciPrivateLinkLabelfortheserviceEPgDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "cloudPrivateLinkLabel")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
