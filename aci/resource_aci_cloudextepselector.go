package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciCloudEndpointSelectorforExternalEPgs() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciCloudEndpointSelectorforExternalEPgsCreate,
		UpdateContext: resourceAciCloudEndpointSelectorforExternalEPgsUpdate,
		ReadContext:   resourceAciCloudEndpointSelectorforExternalEPgsRead,
		DeleteContext: resourceAciCloudEndpointSelectorforExternalEPgsDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciCloudEndpointSelectorforExternalEPgsImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"cloud_external_epg_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"is_shared": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},

			"match_expression": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"subnet": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		}),
	}
}
func getRemoteCloudEndpointSelectorforExternalEPgs(client *client.Client, dn string) (*models.CloudEndpointSelectorforExternalEPgs, error) {
	cloudExtEPSelectorCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudExtEPSelector := models.CloudEndpointSelectorforExternalEPgsFromContainer(cloudExtEPSelectorCont)

	if cloudExtEPSelector.DistinguishedName == "" {
		return nil, fmt.Errorf("CloudEndpointSelectorforExternalEPgs %s not found", cloudExtEPSelector.DistinguishedName)
	}

	return cloudExtEPSelector, nil
}

func setCloudEndpointSelectorforExternalEPgsAttributes(cloudExtEPSelector *models.CloudEndpointSelectorforExternalEPgs, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(cloudExtEPSelector.DistinguishedName)
	d.Set("description", cloudExtEPSelector.Description)

	if dn != cloudExtEPSelector.DistinguishedName {
		d.Set("cloud_external_epg_dn", "")
	}
	cloudExtEPSelectorMap, err := cloudExtEPSelector.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("cloud_external_epg_dn", GetParentDn(dn, fmt.Sprintf("/extepselector-[%s]", cloudExtEPSelectorMap["subnet"])))
	d.Set("name", cloudExtEPSelectorMap["name"])

	d.Set("annotation", cloudExtEPSelectorMap["annotation"])
	d.Set("is_shared", cloudExtEPSelectorMap["isShared"])
	d.Set("match_expression", cloudExtEPSelectorMap["matchExpression"])
	d.Set("name_alias", cloudExtEPSelectorMap["nameAlias"])
	d.Set("subnet", cloudExtEPSelectorMap["subnet"])
	return d, nil
}

func resourceAciCloudEndpointSelectorforExternalEPgsImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	cloudExtEPSelector, err := getRemoteCloudEndpointSelectorforExternalEPgs(aciClient, dn)

	if err != nil {
		return nil, err
	}
	cloudExtEPSelectorMap, err := cloudExtEPSelector.ToMap()
	if err != nil {
		return nil, err
	}
	subnet := cloudExtEPSelectorMap["subnet"]
	pDN := GetParentDn(dn, fmt.Sprintf("/extepselector-[%s]", subnet))
	d.Set("cloud_external_epg_dn", pDN)
	schemaFilled, err := setCloudEndpointSelectorforExternalEPgsAttributes(cloudExtEPSelector, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciCloudEndpointSelectorforExternalEPgsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] CloudEndpointSelectorforExternalEPgs: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	CloudExternalEPgDn := d.Get("cloud_external_epg_dn").(string)

	cloudExtEPSelectorAttr := models.CloudEndpointSelectorforExternalEPgsAttributes{}
	cloudExtEPSelectorAttr.Name = name
	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudExtEPSelectorAttr.Annotation = Annotation.(string)
	} else {
		cloudExtEPSelectorAttr.Annotation = "{}"
	}
	if IsShared, ok := d.GetOk("is_shared"); ok {
		cloudExtEPSelectorAttr.IsShared = IsShared.(string)
	}
	if MatchExpression, ok := d.GetOk("match_expression"); ok {
		cloudExtEPSelectorAttr.MatchExpression = MatchExpression.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		cloudExtEPSelectorAttr.NameAlias = NameAlias.(string)
	}
	if Subnet, ok := d.GetOk("subnet"); ok {
		cloudExtEPSelectorAttr.Subnet = Subnet.(string)
	}
	cloudExtEPSelector := models.NewCloudEndpointSelectorforExternalEPgs(fmt.Sprintf("extepselector-[%s]", cloudExtEPSelectorAttr.Subnet), CloudExternalEPgDn, desc, cloudExtEPSelectorAttr)

	err := aciClient.Save(cloudExtEPSelector)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cloudExtEPSelector.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciCloudEndpointSelectorforExternalEPgsRead(ctx, d, m)
}

func resourceAciCloudEndpointSelectorforExternalEPgsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] CloudEndpointSelectorforExternalEPgs: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	CloudExternalEPgDn := d.Get("cloud_external_epg_dn").(string)

	cloudExtEPSelectorAttr := models.CloudEndpointSelectorforExternalEPgsAttributes{}
	cloudExtEPSelectorAttr.Name = name

	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudExtEPSelectorAttr.Annotation = Annotation.(string)
	} else {
		cloudExtEPSelectorAttr.Annotation = "{}"
	}
	if IsShared, ok := d.GetOk("is_shared"); ok {
		cloudExtEPSelectorAttr.IsShared = IsShared.(string)
	}
	if MatchExpression, ok := d.GetOk("match_expression"); ok {
		cloudExtEPSelectorAttr.MatchExpression = MatchExpression.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		cloudExtEPSelectorAttr.NameAlias = NameAlias.(string)
	}
	if Subnet, ok := d.GetOk("subnet"); ok {
		cloudExtEPSelectorAttr.Subnet = Subnet.(string)
	}
	cloudExtEPSelector := models.NewCloudEndpointSelectorforExternalEPgs(fmt.Sprintf("extepselector-[%s]", cloudExtEPSelectorAttr.Subnet), CloudExternalEPgDn, desc, cloudExtEPSelectorAttr)

	cloudExtEPSelector.Status = "modified"

	err := aciClient.Save(cloudExtEPSelector)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cloudExtEPSelector.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciCloudEndpointSelectorforExternalEPgsRead(ctx, d, m)

}

func resourceAciCloudEndpointSelectorforExternalEPgsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	cloudExtEPSelector, err := getRemoteCloudEndpointSelectorforExternalEPgs(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setCloudEndpointSelectorforExternalEPgsAttributes(cloudExtEPSelector, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciCloudEndpointSelectorforExternalEPgsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "cloudExtEPSelector")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
