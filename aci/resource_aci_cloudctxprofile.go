package aci

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
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
				ValidateFunc: validation.StringInSlice([]string{
					"aws",
					"azure",
					"gcp",
				}, false),
			},

			"relation_cloud_rs_ctx_to_flow_log": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to cloudAwsFlowLogPol",
			},
			"relation_cloud_rs_to_ctx": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Create relation to fvCtx",
			},
			"hub_network": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "hub network to enable transit gateway",
			},
			"cloud_brownfield": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Import Brownfield Virtual Network",
				ForceNew:    true,
			},
			"access_policy_type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Ctx access policy type",
				ValidateFunc: validation.StringInSlice([]string{
					"read-only",
				}, false),
			},
		}),
	}
}

func getRemoteCloudContextProfile(client *client.Client, dn string, d *schema.ResourceData) (*models.CloudContextProfile, error) {
	baseurlStr := "/api/node/mo"
	dnUrl := fmt.Sprintf("%s/%s.json?rsp-subtree=children", baseurlStr, dn)
	cloudCtxProfileCont, err := client.GetViaURL(dnUrl)
	if err != nil {
		return nil, err
	}

	// To set relational attributes value
	setRelationalAttributes(cloudCtxProfileCont, d)

	cloudCtxProfile := models.CloudContextProfileFromContainer(cloudCtxProfileCont)

	if cloudCtxProfile.DistinguishedName == "" {
		return nil, fmt.Errorf("Cloud Context Profiles: %s not found", dn)
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
	d.Set("tenant_dn", GetParentDn(dn, fmt.Sprintf("/"+models.RncloudCtxProfile, GetMOName(cloudCtxProfile.DistinguishedName))))

	d.Set("name", GetMOName(cloudCtxProfile.DistinguishedName))
	d.Set("annotation", cloudCtxProfileMap["annotation"])
	d.Set("name_alias", cloudCtxProfileMap["nameAlias"])
	d.Set("type", cloudCtxProfileMap["type"])
	return d, nil
}

func setRelationalAttributes(cloudCtxProfileCont *container.Container, d *schema.ResourceData) {
	ChildContList, err := cloudCtxProfileCont.S("imdata").Index(0).S(models.CloudctxprofileClassName, "children").Children()
	if err != nil {
		log.Printf("[DEBUG]: Failed to set relational attributes : %v", err)
	}

	CloudVendorPattern := regexp.MustCompile(`uni/clouddomp/provp-(.+)/region-`)
	for _, childCont := range ChildContList {
		if childCont.Exists("cloudCidr") {
			d.Set("primary_cidr", G(childCont.S("cloudCidr", "attributes"), "addr"))
		} else if childCont.Exists("cloudRsCtxProfileToRegion") {
			d.Set("region", GetMOName(G(childCont.S("cloudRsCtxProfileToRegion", "attributes"), "tDn")))
			d.Set("cloud_vendor", CloudVendorPattern.FindStringSubmatch(G(childCont.S("cloudRsCtxProfileToRegion", "attributes"), "tDn"))[1])
		} else if childCont.Exists("cloudRsCtxProfileToGatewayRouterP") {
			d.Set("hub_network", G(childCont.S("cloudRsCtxProfileToGatewayRouterP", "attributes"), "tDn"))
		} else if childCont.Exists("cloudRsCtxToFlowLog") {
			d.Set("relation_cloud_rs_ctx_to_flow_log", G(childCont.S("cloudRsCtxToFlowLog", "attributes"), "tDn"))
		} else if childCont.Exists("cloudRsToCtx") {
			d.Set("relation_cloud_rs_to_ctx", G(childCont.S("cloudRsToCtx", "attributes"), "tDn"))
		} else if childCont.Exists("cloudBrownfield") {
			d.Set("cloud_brownfield", G(childCont.S("cloudBrownfield", "children", "cloudIDMapping", "attributes"), "cloudProviderId"))
		} else if childCont.Exists("cloudRsCtxProfileToAccessPolicy") {
			d.Set("access_policy_type", GetMOName(G(childCont.S("cloudRsCtxProfileToAccessPolicy", "attributes"), "tDn")))
		}
	}
}

func resourceAciCloudContextProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	cloudCtxProfile, err := getRemoteCloudContextProfile(aciClient, dn, d)

	if err != nil {
		return nil, err
	}

	name := GetMOName(cloudCtxProfile.DistinguishedName)
	parentDn := GetParentDn(dn, fmt.Sprintf("/"+models.RncloudCtxProfile, name))
	d.Set("tenant_dn", parentDn)
	schemaFilled, err := setCloudContextProfileAttributes(cloudCtxProfile, d)

	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciCloudContextProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] cloudCtxProfile: Beginning Creation")
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

	cloudBrownfield := ""
	if cloudBrownfield, ok := d.GetOk("cloud_brownfield"); ok {
		cloudBrownfield = cloudBrownfield.(string)
	}

	PrimaryCIDR := d.Get("primary_cidr").(string)

	Region := d.Get("region").(string)

	vendor := d.Get("cloud_vendor").(string)

	checkDns := make([]string, 0, 1)

	if tempVar, ok := d.GetOk("relation_cloud_rs_to_ctx"); ok {
		checkDns = append(checkDns, tempVar.(string))
	}

	if relationTocloudRsCtxToFlowLog, ok := d.GetOk("relation_cloud_rs_ctx_to_flow_log"); ok {
		checkDns = append(checkDns, relationTocloudRsCtxToFlowLog.(string))
	}

	if temp, ok := d.GetOk("hub_network"); ok {
		checkDns = append(checkDns, temp.(string))
	}

	accessPolicy := ""
	if accessPolicy, ok := d.GetOk("access_policy_type"); ok {
		accessDn := fmt.Sprintf("uni/tn-infra/accesspolicy-%s", accessPolicy.(string))
		checkDns = append(checkDns, accessDn)
	}

	d.Partial(true)
	err := checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	cloudRsCtx := ""
	if cloudRsCtxDn, ok := d.GetOk("relation_cloud_rs_to_ctx"); ok {
		cloudRsCtx = GetMOName(cloudRsCtxDn.(string))
	}

	if cloudRsCtx != "" && vendor == "gcp" {
		cloudCtxProfileAttr.VpcGroup = cloudRsCtx
	} else if cloudRsCtx == "" {
		return diag.FromErr(fmt.Errorf("Invalid Configuration relation_cloud_rs_to_ctx property cannot be empty for the Cloud APIC"))
	}

	cloudCtxProfile, err := aciClient.CreateCloudContextProfile(name, TenantDn, desc, cloudCtxProfileAttr, PrimaryCIDR, Region, vendor, cloudRsCtx, cloudBrownfield, accessPolicy)
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

	if temp, ok := d.GetOk("hub_network"); ok {
		err := aciClient.CreateRelationcloudRsCtxProfileTocloudRsCtxProfileToGatewayRouterP(cloudCtxProfile.DistinguishedName, temp.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(cloudCtxProfile.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciCloudContextProfileRead(ctx, d, m)
}

func resourceAciCloudContextProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] cloudCtxProfile: Beginning Update")
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

	checkDns := make([]string, 0, 1)

	if tempVar, ok := d.GetOk("relation_cloud_rs_to_ctx"); ok {
		checkDns = append(checkDns, tempVar.(string))
	}

	if d.HasChange("relation_cloud_rs_ctx_to_flow_log") {
		_, newRelParam := d.GetChange("relation_cloud_rs_ctx_to_flow_log")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("hub_network") {
		_, newRelParam := d.GetChange("hub_network")
		checkDns = append(checkDns, newRelParam.(string))
	}

	// if d.HasChange("cloud_brownfield") {
	// 	_, newRelParam := d.GetChange("cloud_brownfield")
	// 	checkDns = append(checkDns, newRelParam.(string))
	// }

	if d.HasChange("access_policy_type") {
		_, newRelParam := d.GetChange("access_policy_type")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err := checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	cloudRsCtx := ""
	if cloudRsCtxDn, ok := d.GetOk("relation_cloud_rs_to_ctx"); ok {
		cloudRsCtx = GetMOName(cloudRsCtxDn.(string))
	}

	if cloudRsCtx != "" && vendor == "gcp" {
		cloudCtxProfileAttr.VpcGroup = cloudRsCtx
	} else if cloudRsCtx == "" {
		return diag.FromErr(fmt.Errorf("Invalid Configuration relation_cloud_rs_to_ctx property cannot be empty for the Cloud APIC"))
	}

	accessPolicy := ""
	if accessPolicy, ok := d.GetOk("access_policy_type"); ok {
		accessPolicy = accessPolicy.(string)
	}

	cloudBrownfield := ""
	if cloudBrownfield, ok := d.GetOk("cloud_brownfield"); ok {
		cloudBrownfield = cloudBrownfield.(string)
	}

	cloudCtxProfile, err := aciClient.UpdateCloudContextProfile(name, TenantDn, desc, cloudCtxProfileAttr, PrimaryCIDR, Region, vendor, cloudRsCtx, cloudBrownfield, accessPolicy)

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
		if newRelParamName != "" {
			err = aciClient.CreateRelationcloudRsCtxToFlowLogFromCloudContextProfile(cloudCtxProfile.DistinguishedName, newRelParamName)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("hub_network") {
		oldRelParam, newRelParam := d.GetChange("hub_network")
		err = aciClient.DeleteRelationcloudRsCtxProfileTocloudRsCtxProfileToGatewayRouterP(cloudCtxProfile.DistinguishedName, oldRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		if newRelParam.(string) != "" {
			err = aciClient.CreateRelationcloudRsCtxProfileTocloudRsCtxProfileToGatewayRouterP(cloudCtxProfile.DistinguishedName, newRelParam.(string))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId(cloudCtxProfile.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciCloudContextProfileRead(ctx, d, m)
}

func resourceAciCloudContextProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()
	cloudCtxProfile, err := getRemoteCloudContextProfile(aciClient, dn, d)

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
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciCloudContextProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, models.CloudctxprofileClassName)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	return diag.FromErr(err)
}
