package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciCloudContextProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciCloudContextProfileCreate,
		UpdateContext: resourceAciCloudContextProfileUpdate,
		ReadContext:   resourceAciCloudContextProfileRead,
		DeleteContext: resourceAciCloudContextProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciCloudContextProfileImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"tenant_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				ValidateFunc: validation.StringInSlice([]string{
					"regular",
					"shadow",
					"hosted",
					"container-overlay",
				}, false),
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

			"cloud_vendor": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the vendor",
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
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to cloudRegion",
			},
			"hub_network": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "hub network to enable transit gateway",
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

func setCloudContextProfileAttributes(cloudCtxProfile *models.CloudContextProfile, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(cloudCtxProfile.DistinguishedName)
	d.Set("description", cloudCtxProfile.Description)
	if dn != cloudCtxProfile.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	cloudCtxProfileMap, err := cloudCtxProfile.ToMap()

	if err != nil {
		return d, err
	}
	d.Set("tenant_dn", GetParentDn(dn, fmt.Sprintf("/ctxprofile-%s", GetMOName(cloudCtxProfile.DistinguishedName))))

	d.Set("name", GetMOName(cloudCtxProfile.DistinguishedName))

	d.Set("annotation", cloudCtxProfileMap["annotation"])
	d.Set("name_alias", cloudCtxProfileMap["nameAlias"])
	d.Set("type", cloudCtxProfileMap["type"])
	d.Set("primary_cidr", cloudCtxProfileMap["primary_cidr"])
	d.Set("region", cloudCtxProfileMap["region"])

	return d, nil
}

func resourceAciCloudContextProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	aciClient := m.(*client.Client)

	dn := d.Id()

	cloudCtxProfile, err := getRemoteCloudContextProfile(aciClient, dn)

	if err != nil {
		return nil, err
	}

	name := GetMOName(cloudCtxProfile.DistinguishedName)
	pDN := GetParentDn(dn, fmt.Sprintf("/ctxprofile-%s", name))
	d.Set("tenant_dn", pDN)
	schemaFilled, err := setCloudContextProfileAttributes(cloudCtxProfile, d)

	if err != nil {
		return nil, err
	}
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciCloudContextProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	cloudCtxProfileAttr := models.CloudContextProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudCtxProfileAttr.Annotation = Annotation.(string)
	} else {
		cloudCtxProfileAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		cloudCtxProfileAttr.NameAlias = NameAlias.(string)
	}
	if Type, ok := d.GetOk("type"); ok {
		cloudCtxProfileAttr.Type = Type.(string)
	}

	PrimaryCIDR := d.Get("primary_cidr").(string)

	Region := d.Get("region").(string)

	vendor := d.Get("cloud_vendor").(string)

	cloudCtxProfile := models.NewCloudContextProfile(fmt.Sprintf("ctxprofile-%s", name), TenantDn, desc, cloudCtxProfileAttr)

	checkDns := make([]string, 0, 1)

	if tempVar, ok := d.GetOk("relation_cloud_rs_to_ctx"); ok {
		checkDns = append(checkDns, tempVar.(string))
	}

	if relationTocloudRsCtxToFlowLog, ok := d.GetOk("relation_cloud_rs_ctx_to_flow_log"); ok {
		checkDns = append(checkDns, relationTocloudRsCtxToFlowLog.(string))
	}

	if relationTocloudRsCtxProfileToRegion, ok := d.GetOk("relation_cloud_rs_ctx_profile_to_region"); ok {
		checkDns = append(checkDns, relationTocloudRsCtxProfileToRegion.(string))
	}

	if temp, ok := d.GetOk("hub_network"); ok {
		checkDns = append(checkDns, temp.(string))
	}

	d.Partial(true)
	err := checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	var cloudRsCtx string
	if tempVar, ok := d.GetOk("relation_cloud_rs_to_ctx"); ok {
		cloudRsCtx = tempVar.(string)
		cloudRsCtx = GetMOName(cloudRsCtx)
	} else {
		cloudRsCtx = ""
	}
	cloudCtxProfile, err = aciClient.CreateCloudContextProfile(name, TenantDn, desc, cloudCtxProfileAttr, PrimaryCIDR, Region, vendor, cloudRsCtx)
	//err := aciClient.Save(cloudCtxProfile)
	if err != nil {
		return diag.FromErr(err)
	}

	if relationTocloudRsCtxToFlowLog, ok := d.GetOk("relation_cloud_rs_ctx_to_flow_log"); ok {
		relationParam := relationTocloudRsCtxToFlowLog.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationcloudRsCtxToFlowLogFromCloudContextProfile(cloudCtxProfile.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTocloudRsCtxProfileToRegion, ok := d.GetOk("relation_cloud_rs_ctx_profile_to_region"); ok {
		relationParam := relationTocloudRsCtxProfileToRegion.(string)
		err = aciClient.CreateRelationcloudRsCtxProfileToRegionFromCloudContextProfile(cloudCtxProfile.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if temp, ok := d.GetOk("hub_network"); ok {
		err := aciClient.CreateRelationcloudRsCtxProfileTocloudRsCtxProfileToGatewayRouterP(cloudCtxProfile.DistinguishedName, temp.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(cloudCtxProfile.DistinguishedName)
	return resourceAciCloudContextProfileRead(ctx, d, m)
}

func resourceAciCloudContextProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	cloudCtxProfileAttr := models.CloudContextProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudCtxProfileAttr.Annotation = Annotation.(string)
	} else {
		cloudCtxProfileAttr.Annotation = "{}"
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

	vendor := d.Get("cloud_vendor").(string)

	checkDns := make([]string, 0, 1)

	if tempVar, ok := d.GetOk("relation_cloud_rs_to_ctx"); ok {
		checkDns = append(checkDns, tempVar.(string))
	}

	if d.HasChange("relation_cloud_rs_ctx_to_flow_log") {
		_, newRelParam := d.GetChange("relation_cloud_rs_ctx_to_flow_log")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_cloud_rs_ctx_profile_to_region") {
		_, newRelParam := d.GetChange("relation_cloud_rs_ctx_profile_to_region")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("hub_network") {
		_, newRelParam := d.GetChange("hub_network")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err := checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	var cloudRsCtx string
	if tempVar, ok := d.GetOk("relation_cloud_rs_to_ctx"); ok {
		cloudRsCtx = tempVar.(string)
		cloudRsCtx = GetMOName(cloudRsCtx)
	} else {
		cloudRsCtx = ""
	}

	cloudCtxProfile, err = aciClient.UpdateCloudContextProfile(name, TenantDn, desc, cloudCtxProfileAttr, PrimaryCIDR, Region, vendor, cloudRsCtx)
	//err := aciClient.Save(cloudCtxProfile)

	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("relation_cloud_rs_ctx_to_flow_log") {
		_, newRelParam := d.GetChange("relation_cloud_rs_ctx_to_flow_log")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationcloudRsCtxToFlowLogFromCloudContextProfile(cloudCtxProfile.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationcloudRsCtxToFlowLogFromCloudContextProfile(cloudCtxProfile.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_cloud_rs_ctx_profile_to_region") {
		_, newRelParam := d.GetChange("relation_cloud_rs_ctx_profile_to_region")
		err = aciClient.DeleteRelationcloudRsCtxProfileToRegionFromCloudContextProfile(cloudCtxProfile.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationcloudRsCtxProfileToRegionFromCloudContextProfile(cloudCtxProfile.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("hub_network") {
		oldRelParam, newRelParam := d.GetChange("hub_network")
		err = aciClient.DeleteRelationcloudRsCtxProfileTocloudRsCtxProfileToGatewayRouterP(cloudCtxProfile.DistinguishedName, oldRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationcloudRsCtxProfileTocloudRsCtxProfileToGatewayRouterP(cloudCtxProfile.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(cloudCtxProfile.DistinguishedName)
	return resourceAciCloudContextProfileRead(ctx, d, m)

}

func resourceAciCloudContextProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	dn := d.Id()
	cloudCtxProfile, err := getRemoteCloudContextProfile(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	_, err = setCloudContextProfileAttributes(cloudCtxProfile, d)

	if err != nil {
		return diag.FromErr(err)
	}

	if hub, ok := d.GetOk("hub_network"); ok {
		dURL := fmt.Sprintf("%s/rsctxProfileToGatewayRouterP-[%s]", dn, hub.(string))
		_, err := aciClient.Get(dURL)
		if err != nil {
			d.Set("hub_network", "")
		}
	}
	return nil
}

func resourceAciCloudContextProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "cloudCtxProfile")
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diag.FromErr(err)
}
