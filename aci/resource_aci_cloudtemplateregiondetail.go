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
					"disabled",
					"enabled",
				}, false),
			},
		})),
	}
}

func getRemoteCloudTemplateRegion(client *client.Client, dn string) (*models.CloudTemplateRegion, error) {
	cloudTemplateRegionDetailCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	cloudTemplateRegionDetail := models.CloudTemplateRegionFromContainer(cloudTemplateRegionDetailCont)
	if cloudTemplateRegionDetail.DistinguishedName == "" {
		return nil, fmt.Errorf("Cloud Template Region %s not found", dn)
	}
	return cloudTemplateRegionDetail, nil
}

func setCloudTemplateRegionAttributes(cloudTemplateRegionDetail *models.CloudTemplateRegion, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(cloudTemplateRegionDetail.DistinguishedName)
	cloudTemplateRegionDetailMap, err := cloudTemplateRegionDetail.ToMap()
	if err != nil {
		return d, err
	}
	if d.Id() != cloudTemplateRegionDetail.DistinguishedName {
		d.Set("parent_dn", "")
	} else {
		d.Set("parent_dn", GetParentDn(cloudTemplateRegionDetail.DistinguishedName, fmt.Sprintf("/"+models.RnCloudtemplateRegionDetail)))
	}
	d.Set("annotation", cloudTemplateRegionDetailMap["annotation"])
	d.Set("hub_networking", toggleOptions(cloudTemplateRegionDetailMap["hubNetworkingEnabled"]))
	return d, nil
}

func resourceAciCloudTemplateRegionImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	dn := d.Id()
	log.Printf("[DEBUG] %s: Beginning Import", dn)
	aciClient := m.(*client.Client)
	cloudTemplateRegionDetail, err := getRemoteCloudTemplateRegion(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setCloudTemplateRegionAttributes(cloudTemplateRegionDetail, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", dn)
	return []*schema.ResourceData{schemaFilled}, nil
}

func makeCloudTemplateRegion(ctx context.Context, d *schema.ResourceData, m interface{}, method string) diag.Diagnostics {
	aciClient := m.(*client.Client)
	cloudProviderandRegionNamesDn := d.Get("parent_dn").(string)

	cloudTemplateRegionDetailAttr := models.CloudTemplateRegionAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudTemplateRegionDetailAttr.Annotation = Annotation.(string)
	} else {
		cloudTemplateRegionDetailAttr.Annotation = "{}"
	}

	if HubNetworkingEnabled, ok := d.GetOk("hub_networking"); ok {
		cloudTemplateRegionDetailAttr.HubNetworkingEnabled = toggleOptions(HubNetworkingEnabled.(string))
	}
	cloudTemplateRegionDetail := models.NewCloudTemplateRegion(fmt.Sprintf(models.RnCloudtemplateRegionDetail), cloudProviderandRegionNamesDn, cloudTemplateRegionDetailAttr)

	if method == "update" {
		cloudTemplateRegionDetail.Status = "modified"
	}

	err := aciClient.Save(cloudTemplateRegionDetail)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(cloudTemplateRegionDetail.DistinguishedName)
	return nil
}

func resourceAciCloudTemplateRegionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] CloudTemplateRegion: Beginning Creation")
	makeCloudTemplateRegion(ctx, d, m, "create")
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciCloudTemplateRegionRead(ctx, d, m)
}
func resourceAciCloudTemplateRegionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] CloudTemplateRegion: Beginning Update")
	makeCloudTemplateRegion(ctx, d, m, "update")
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciCloudTemplateRegionRead(ctx, d, m)
}

func resourceAciCloudTemplateRegionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	dn := d.Id()
	log.Printf("[DEBUG] %s: Beginning Read", dn)
	aciClient := m.(*client.Client)

	cloudTemplateRegionDetail, err := getRemoteCloudTemplateRegion(aciClient, dn)

	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}

	_, err = setCloudTemplateRegionAttributes(cloudTemplateRegionDetail, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", dn)
	return nil
}

func resourceAciCloudTemplateRegionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	dn := d.Id()
	log.Printf("[DEBUG] %s: Beginning Destroy", dn)
	aciClient := m.(*client.Client)

	err := aciClient.DeleteByDn(dn, "cloudtemplateRegionDetail")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", dn)
	d.SetId("")
	return diag.FromErr(err)
}
