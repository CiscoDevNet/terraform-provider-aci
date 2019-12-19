package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciCloudContextProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciCloudContextProfileCreate,
		Update: resourceAciCloudContextProfileUpdate,
		Read:   resourceAciCloudContextProfileRead,
		Delete: resourceAciCloudContextProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciCloudContextProfileImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"tenant_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"annotation": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Mo doc not defined in techpub!!!",
			},

			"name_alias": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Mo doc not defined in techpub!!!",
			},

			"type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "component type",
			},

			"primary_cidr": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Primary CIDR block",
			},

			"region": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "region",
			},

			"relation_cloud_rs_ctx_to_flow_log": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to cloudAwsFlowLogPol",
			},
			"relation_cloud_rs_to_ctx": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to fvCtx",
			},
			"relation_cloud_rs_ctx_profile_to_region": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to cloudRegion",
			},
		}),
	}
}

func getRemoteCloudContextProfile(client *client.Client, dn string) (*models.CloudContextProfile, error) {
	baseurlStr := "/api/node/mo"
	dnUrl := fmt.Sprintf("%s/%s.json?rsp-subtree=children", baseurlStr, dn)
	cloudCtxProfileCont, err := client.GetViaURL(dnUrl)
	if err != nil {
		return nil, err
	}

	cloudCtxProfile := models.CloudContextProfileFromContainer(cloudCtxProfileCont)

	if cloudCtxProfile.DistinguishedName == "" {
		return nil, fmt.Errorf("Bridge Domain %s not found", cloudCtxProfile.DistinguishedName)
	}

	return cloudCtxProfile, nil
}

func setCloudContextProfileAttributes(cloudCtxProfile *models.CloudContextProfile, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(cloudCtxProfile.DistinguishedName)
	d.Set("description", cloudCtxProfile.Description)
	d.Set("tenant_dn", GetParentDn(cloudCtxProfile.DistinguishedName))
	cloudCtxProfileMap, _ := cloudCtxProfile.ToMap()
	d.Set("name", GetMOName(cloudCtxProfile.DistinguishedName))

	d.Set("annotation", cloudCtxProfileMap["annotation"])
	d.Set("name_alias", cloudCtxProfileMap["nameAlias"])
	d.Set("type", cloudCtxProfileMap["type"])
	d.Set("primary_cidr", cloudCtxProfileMap["primary_cidr"])
	d.Set("region", cloudCtxProfileMap["region"])

	return d
}

func resourceAciCloudContextProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	aciClient := m.(*client.Client)

	dn := d.Id()

	cloudCtxProfile, err := getRemoteCloudContextProfile(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setCloudContextProfileAttributes(cloudCtxProfile, d)
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciCloudContextProfileCreate(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	cloudCtxProfileAttr := models.CloudContextProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudCtxProfileAttr.Annotation = Annotation.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		cloudCtxProfileAttr.NameAlias = NameAlias.(string)
	}
	if Type, ok := d.GetOk("type"); ok {
		cloudCtxProfileAttr.Type = Type.(string)
	}

	PrimaryCIDR := d.Get("primary_cidr").(string)

	Region := d.Get("region").(string)

	cloudCtxProfile := models.NewCloudContextProfile(fmt.Sprintf("ctxprofile-%s", name), TenantDn, desc, cloudCtxProfileAttr)

	cloudCtxProfile, err := aciClient.CreateCloudContextProfile(name, TenantDn, desc, models.CloudContextProfileAttributes{}, PrimaryCIDR, Region, d.Get("relation_cloud_rs_to_ctx").(string))
	//err := aciClient.Save(cloudCtxProfile)
	if err != nil {
		return err
	}

	if relationTocloudRsCtxToFlowLog, ok := d.GetOk("relation_cloud_rs_ctx_to_flow_log"); ok {
		relationParam := relationTocloudRsCtxToFlowLog.(string)
		err = aciClient.CreateRelationcloudRsCtxToFlowLogFromCloudContextProfile(cloudCtxProfile.DistinguishedName, relationParam)
		if err != nil {
			return err
		}

	}
	if relationTocloudRsCtxProfileToRegion, ok := d.GetOk("relation_cloud_rs_ctx_profile_to_region"); ok {
		relationParam := relationTocloudRsCtxProfileToRegion.(string)
		err = aciClient.CreateRelationcloudRsCtxProfileToRegionFromCloudContextProfile(cloudCtxProfile.DistinguishedName, relationParam)
		if err != nil {
			return err
		}

	}

	d.SetId(cloudCtxProfile.DistinguishedName)
	return resourceAciCloudContextProfileRead(d, m)
}

func resourceAciCloudContextProfileUpdate(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	cloudCtxProfileAttr := models.CloudContextProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudCtxProfileAttr.Annotation = Annotation.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		cloudCtxProfileAttr.NameAlias = NameAlias.(string)
	}
	if Type, ok := d.GetOk("type"); ok {
		cloudCtxProfileAttr.Type = Type.(string)
	}
	cloudCtxProfile := models.NewCloudContextProfile(fmt.Sprintf("ctxprofile-%s", name), TenantDn, desc, cloudCtxProfileAttr)

	cloudCtxProfile.Status = "modified"
	PrimaryCIDR := d.Get("primary_cidr").(string)

	Region := d.Get("region").(string)

	cloudCtxProfile, err := aciClient.CreateCloudContextProfile(name, TenantDn, desc, models.CloudContextProfileAttributes{}, PrimaryCIDR, Region, d.Get("relation_cloud_rs_to_ctx").(string))
	//err := aciClient.Save(cloudCtxProfile)

	if err != nil {
		return err
	}

	if d.HasChange("relation_cloud_rs_ctx_to_flow_log") {
		_, newRelParam := d.GetChange("relation_cloud_rs_ctx_to_flow_log")
		err = aciClient.DeleteRelationcloudRsCtxToFlowLogFromCloudContextProfile(cloudCtxProfile.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationcloudRsCtxToFlowLogFromCloudContextProfile(cloudCtxProfile.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}

	}
	if d.HasChange("relation_cloud_rs_ctx_profile_to_region") {
		_, newRelParam := d.GetChange("relation_cloud_rs_ctx_profile_to_region")
		err = aciClient.DeleteRelationcloudRsCtxProfileToRegionFromCloudContextProfile(cloudCtxProfile.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationcloudRsCtxProfileToRegionFromCloudContextProfile(cloudCtxProfile.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}

	}

	d.SetId(cloudCtxProfile.DistinguishedName)
	return resourceAciCloudContextProfileRead(d, m)

}

func resourceAciCloudContextProfileRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	dn := d.Id()
	cloudCtxProfile, err := getRemoteCloudContextProfile(aciClient, dn)

	if err != nil {
		return err
	}
	setCloudContextProfileAttributes(cloudCtxProfile, d)
	return nil
}

func resourceAciCloudContextProfileDelete(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "cloudCtxProfile")
	if err != nil {
		return err
	}

	d.SetId("")
	return err
}
