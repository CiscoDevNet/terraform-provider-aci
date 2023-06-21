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

func resourceAciCloudTemplateRegion() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciCloudTemplateRegionCreate,
		UpdateContext: resourceAciCloudTemplateRegionUpdate,
		ReadContext:   resourceAciCloudTemplateRegionRead,
		DeleteContext: resourceAciCloudTemplateRegionDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciCloudTemplateRegionImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"parent_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"hub_networking": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"disable",
					"enable",
				}, false),
			},
		})),
	}
}

func getActualHubNetworkingValue(hub_networking string) string {
	if hub_networking == "disable" {
		return "no"
	} else {
		return "yes"
	}
}

func setHubNetworkingValue(hub_networking string) string {
	if hub_networking == "no" {
		return "disable"
	} else {
		return "enable"
	}
}

func getRemoteCloudTemplateRegion(client *client.Client, dn string) (*models.CloudTemplateRegion, error) {
	cloudtemplateRegionDetailCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	cloudtemplateRegionDetail := models.CloudTemplateRegionFromContainer(cloudtemplateRegionDetailCont)
	if cloudtemplateRegionDetail.DistinguishedName == "" {
		return nil, fmt.Errorf("Cloud Template Region %s not found", dn)
	}
	return cloudtemplateRegionDetail, nil
}

func setCloudTemplateRegionAttributes(cloudtemplateRegionDetail *models.CloudTemplateRegion, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(cloudtemplateRegionDetail.DistinguishedName)
	cloudtemplateRegionDetailMap, err := cloudtemplateRegionDetail.ToMap()
	if err != nil {
		return d, err
	}
	dn := d.Id()
	if dn != cloudtemplateRegionDetail.DistinguishedName {
		d.Set("parent_dn", "")
	} else {
		d.Set("parent_dn", GetParentDn(cloudtemplateRegionDetail.DistinguishedName, fmt.Sprintf("/"+models.RnCloudtemplateRegionDetail)))
	}
	d.Set("annotation", cloudtemplateRegionDetailMap["annotation"])
	d.Set("hub_networking", setHubNetworkingValue(cloudtemplateRegionDetailMap["hubNetworkingEnabled"]))
	return d, nil
}

func resourceAciCloudTemplateRegionImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	cloudtemplateRegionDetail, err := getRemoteCloudTemplateRegion(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setCloudTemplateRegionAttributes(cloudtemplateRegionDetail, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciCloudTemplateRegionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] CloudTemplateRegion: Beginning Creation")
	aciClient := m.(*client.Client)
	CloudProviderandRegionNamesDn := d.Get("parent_dn").(string)

	cloudtemplateRegionDetailAttr := models.CloudTemplateRegionAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudtemplateRegionDetailAttr.Annotation = Annotation.(string)
	} else {
		cloudtemplateRegionDetailAttr.Annotation = "{}"
	}

	if HubNetworkingEnabled, ok := d.GetOk("hub_networking"); ok {
		cloudtemplateRegionDetailAttr.HubNetworkingEnabled = getActualHubNetworkingValue(HubNetworkingEnabled.(string))
	}
	cloudtemplateRegionDetail := models.NewCloudTemplateRegion(fmt.Sprintf(models.RnCloudtemplateRegionDetail), CloudProviderandRegionNamesDn, cloudtemplateRegionDetailAttr)

	err := aciClient.Save(cloudtemplateRegionDetail)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cloudtemplateRegionDetail.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciCloudTemplateRegionRead(ctx, d, m)
}
func resourceAciCloudTemplateRegionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] CloudTemplateRegion: Beginning Update")
	aciClient := m.(*client.Client)
	CloudProviderandRegionNamesDn := d.Get("parent_dn").(string)

	cloudtemplateRegionDetailAttr := models.CloudTemplateRegionAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudtemplateRegionDetailAttr.Annotation = Annotation.(string)
	} else {
		cloudtemplateRegionDetailAttr.Annotation = "{}"
	}

	if HubNetworkingEnabled, ok := d.GetOk("hub_networking"); ok {
		cloudtemplateRegionDetailAttr.HubNetworkingEnabled = getActualHubNetworkingValue(HubNetworkingEnabled.(string))
	}
	cloudtemplateRegionDetail := models.NewCloudTemplateRegion(fmt.Sprintf(models.RnCloudtemplateRegionDetail), CloudProviderandRegionNamesDn, cloudtemplateRegionDetailAttr)

	cloudtemplateRegionDetail.Status = "modified"

	err := aciClient.Save(cloudtemplateRegionDetail)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cloudtemplateRegionDetail.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciCloudTemplateRegionRead(ctx, d, m)
}

func resourceAciCloudTemplateRegionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	cloudtemplateRegionDetail, err := getRemoteCloudTemplateRegion(aciClient, dn)

	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}

	_, err = setCloudTemplateRegionAttributes(cloudtemplateRegionDetail, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciCloudTemplateRegionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "cloudtemplateRegionDetail")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
