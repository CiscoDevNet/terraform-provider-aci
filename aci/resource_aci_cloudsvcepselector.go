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

func resourceAciCloudServiceEndpointSelector() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciCloudServiceEndpointSelectorCreate,
		UpdateContext: resourceAciCloudServiceEndpointSelectorUpdate,
		ReadContext:   resourceAciCloudServiceEndpointSelectorRead,
		DeleteContext: resourceAciCloudServiceEndpointSelectorDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciCloudServiceEndpointSelectorImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"cloud_service_epg_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"match_expression": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		})),
	}
}

func getRemoteCloudServiceEndpointSelector(client *client.Client, dn string) (*models.CloudServiceEndpointSelector, error) {
	cloudSvcEPSelectorCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	cloudSvcEPSelector := models.CloudServiceEndpointSelectorFromContainer(cloudSvcEPSelectorCont)
	if cloudSvcEPSelector.DistinguishedName == "" {
		return nil, fmt.Errorf("Cloud Service Endpoint Selector %s not found", dn)
	}
	return cloudSvcEPSelector, nil
}

func setCloudServiceEndpointSelectorAttributes(cloudSvcEPSelector *models.CloudServiceEndpointSelector, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(cloudSvcEPSelector.DistinguishedName)
	d.Set("description", cloudSvcEPSelector.Description)
	cloudSvcEPSelectorMap, err := cloudSvcEPSelector.ToMap()
	if err != nil {
		return d, err
	}
	dn := d.Id()
	if dn != cloudSvcEPSelector.DistinguishedName {
		d.Set("cloud_service_epg_dn", "")
	} else {
		d.Set("cloud_service_epg_dn", GetParentDn(cloudSvcEPSelector.DistinguishedName, fmt.Sprintf("/"+models.RnCloudSvcEPSelector, cloudSvcEPSelectorMap["name"])))
	}
	d.Set("annotation", cloudSvcEPSelectorMap["annotation"])
	d.Set("match_expression", cloudSvcEPSelectorMap["matchExpression"])
	d.Set("name", cloudSvcEPSelectorMap["name"])
	d.Set("name_alias", cloudSvcEPSelectorMap["nameAlias"])
	return d, nil
}

func resourceAciCloudServiceEndpointSelectorImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	cloudSvcEPSelector, err := getRemoteCloudServiceEndpointSelector(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setCloudServiceEndpointSelectorAttributes(cloudSvcEPSelector, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciCloudServiceEndpointSelectorCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Cloud Service Endpoint Selector: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	CloudServiceEPgDn := d.Get("cloud_service_epg_dn").(string)

	cloudSvcEPSelectorAttr := models.CloudServiceEndpointSelectorAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudSvcEPSelectorAttr.Annotation = Annotation.(string)
	} else {
		cloudSvcEPSelectorAttr.Annotation = "{}"
	}

	if MatchExpression, ok := d.GetOk("match_expression"); ok {
		cloudSvcEPSelectorAttr.MatchExpression = MatchExpression.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		cloudSvcEPSelectorAttr.Name = Name.(string)
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		cloudSvcEPSelectorAttr.NameAlias = NameAlias.(string)
	}
	cloudSvcEPSelector := models.NewCloudServiceEndpointSelector(fmt.Sprintf(models.RnCloudSvcEPSelector, name), CloudServiceEPgDn, desc, cloudSvcEPSelectorAttr)

	err := aciClient.Save(cloudSvcEPSelector)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cloudSvcEPSelector.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciCloudServiceEndpointSelectorRead(ctx, d, m)
}
func resourceAciCloudServiceEndpointSelectorUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Cloud Service Endpoint Selector: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	CloudServiceEPgDn := d.Get("cloud_service_epg_dn").(string)

	cloudSvcEPSelectorAttr := models.CloudServiceEndpointSelectorAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudSvcEPSelectorAttr.Annotation = Annotation.(string)
	} else {
		cloudSvcEPSelectorAttr.Annotation = "{}"
	}

	if MatchExpression, ok := d.GetOk("match_expression"); ok {
		cloudSvcEPSelectorAttr.MatchExpression = MatchExpression.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		cloudSvcEPSelectorAttr.Name = Name.(string)
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		cloudSvcEPSelectorAttr.NameAlias = NameAlias.(string)
	}
	cloudSvcEPSelector := models.NewCloudServiceEndpointSelector(fmt.Sprintf(models.RnCloudSvcEPSelector, name), CloudServiceEPgDn, desc, cloudSvcEPSelectorAttr)

	cloudSvcEPSelector.Status = "modified"

	err := aciClient.Save(cloudSvcEPSelector)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cloudSvcEPSelector.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciCloudServiceEndpointSelectorRead(ctx, d, m)
}

func resourceAciCloudServiceEndpointSelectorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	cloudSvcEPSelector, err := getRemoteCloudServiceEndpointSelector(aciClient, dn)
	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}

	_, err = setCloudServiceEndpointSelectorAttributes(cloudSvcEPSelector, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciCloudServiceEndpointSelectorDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "cloudSvcEPSelector")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
