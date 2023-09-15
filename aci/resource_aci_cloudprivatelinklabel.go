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

func resourceAciCloudPrivateLinkLabel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciCloudPrivateLinkLabelCreate,
		UpdateContext: resourceAciCloudPrivateLinkLabelUpdate,
		ReadContext:   resourceAciCloudPrivateLinkLabelRead,
		DeleteContext: resourceAciCloudPrivateLinkLabelDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciCloudPrivateLinkLabelImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"parent_dn": {
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

func getRemoteCloudPrivateLinkLabel(client *client.Client, dn string) (*models.CloudPrivateLinkLabel, error) {
	cloudPrivateLinkLabelCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	cloudPrivateLinkLabel := models.CloudPrivateLinkLabelFromContainer(cloudPrivateLinkLabelCont)
	if cloudPrivateLinkLabel.DistinguishedName == "" {
		return nil, fmt.Errorf("Cloud Private Link Label %s not found", dn)
	}
	return cloudPrivateLinkLabel, nil
}

func setCloudPrivateLinkLabelAttributes(cloudPrivateLinkLabel *models.CloudPrivateLinkLabel, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(cloudPrivateLinkLabel.DistinguishedName)
	d.Set("description", cloudPrivateLinkLabel.Description)
	cloudPrivateLinkLabelMap, err := cloudPrivateLinkLabel.ToMap()
	if err != nil {
		return d, err
	}
	dn := d.Id()
	if dn != cloudPrivateLinkLabel.DistinguishedName {
		d.Set("parent_dn", "")
	} else {
		d.Set("parent_dn", GetParentDn(cloudPrivateLinkLabel.DistinguishedName, fmt.Sprintf("/"+models.RnCloudPrivateLinkLabel, cloudPrivateLinkLabelMap["name"])))
	}
	d.Set("annotation", cloudPrivateLinkLabelMap["annotation"])
	d.Set("name", cloudPrivateLinkLabelMap["name"])
	d.Set("name_alias", cloudPrivateLinkLabelMap["nameAlias"])
	return d, nil
}

func resourceAciCloudPrivateLinkLabelImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	cloudPrivateLinkLabel, err := getRemoteCloudPrivateLinkLabel(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setCloudPrivateLinkLabelAttributes(cloudPrivateLinkLabel, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciCloudPrivateLinkLabelCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Cloud Private Link Label: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	CloudServiceEPgDn := d.Get("parent_dn").(string)

	cloudPrivateLinkLabelAttr := models.CloudPrivateLinkLabelAttributes{}

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
	cloudPrivateLinkLabel := models.NewCloudPrivateLinkLabel(fmt.Sprintf(models.RnCloudPrivateLinkLabel, name), CloudServiceEPgDn, desc, cloudPrivateLinkLabelAttr)

	err := aciClient.Save(cloudPrivateLinkLabel)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cloudPrivateLinkLabel.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciCloudPrivateLinkLabelRead(ctx, d, m)
}
func resourceAciCloudPrivateLinkLabelUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Cloud Private Link Label: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	CloudServiceEPgDn := d.Get("parent_dn").(string)

	cloudPrivateLinkLabelAttr := models.CloudPrivateLinkLabelAttributes{}

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
	cloudPrivateLinkLabel := models.NewCloudPrivateLinkLabel(fmt.Sprintf(models.RnCloudPrivateLinkLabel, name), CloudServiceEPgDn, desc, cloudPrivateLinkLabelAttr)

	cloudPrivateLinkLabel.Status = "modified"

	err := aciClient.Save(cloudPrivateLinkLabel)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cloudPrivateLinkLabel.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciCloudPrivateLinkLabelRead(ctx, d, m)
}

func resourceAciCloudPrivateLinkLabelRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	cloudPrivateLinkLabel, err := getRemoteCloudPrivateLinkLabel(aciClient, dn)
	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}

	_, err = setCloudPrivateLinkLabelAttributes(cloudPrivateLinkLabel, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciCloudPrivateLinkLabelDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
