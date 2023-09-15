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

const (
	CloudLBClassName                  = "cloudLB"
	CloudRsLDevToCloudSubnetClassName = "cloudRsLDevToCloudSubnet"
	RnCloudLB                         = "clb-%s"
)

func resourceAciCloudL4L7LoadBalancer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciCloudL4L7LoadBalancerCreate,
		UpdateContext: resourceAciCloudL4L7LoadBalancerUpdate,
		ReadContext:   resourceAciCloudL4L7LoadBalancerRead,
		DeleteContext: resourceAciCloudL4L7LoadBalancerDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciCloudL4L7LoadBalancerImport,
		},

		SchemaVersion: 1,
		Schema: AppendAttrSchemas(map[string]*schema.Schema{
			"tenant_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"version": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"active_active": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"allow_all": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"auto_scaling": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"context_aware": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"multi-Context",
					"single-Context",
				}, false),
			},
			"custom_resource_group": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"device_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"CLOUD",
					"PHYSICAL",
					"VIRTUAL",
				}, false),
			},
			"function_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"GoThrough",
					"GoTo",
					"L1",
					"L2",
					"None",
				}, false),
			},
			"instance_count": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_copy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"is_instantiation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"is_static_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"l4l7_device_application_security_group": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"l4l7_third_party_device": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"managed": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"max_instance_count": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"min_instance_count": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"mode": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"native_lb_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"package_model": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"promiscuous_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"scheme": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"internal",
					"internet",
				}, false),
			},
			"size": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"large",
					"medium",
					"small",
					"v2",
				}, false),
			},
			"sku": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"WAF",
					"WAF_v2",
					"standard",
					"standard_v2",
				}, false),
			},
			"service_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ADC",
					"COPY",
					"FW",
					"NATIVELB",
					"OTHERS",
				}, false),
			},
			"target_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"primary",
					"secondary",
					"unspecified",
				}, false),
			},
			"trunking": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"cloud_l4l7_load_balancer_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"application",
					"network",
				}, false),
			},
			"relation_cloud_rs_ldev_to_cloud_subnet": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"aaa_domain_dn": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
		}, GetNameAliasAttrSchema(), GetAnnotationAttrSchema()),
	}
}

func mapCloudL4L7LoadBalancerAttrs(status string, d *schema.ResourceData) map[string]string {
	return map[string]string{
		"activeActive":                       d.Get("active_active").(string),
		"allowAll":                           d.Get("allow_all").(string),
		"annotation":                         d.Get("annotation").(string),
		"autoScaling":                        d.Get("auto_scaling").(string),
		"contextAware":                       d.Get("context_aware").(string),
		"customRG":                           d.Get("custom_resource_group").(string),
		"devtype":                            d.Get("device_type").(string),
		"funcType":                           d.Get("function_type").(string),
		"instanceCount":                      d.Get("instance_count").(string),
		"isCopy":                             d.Get("is_copy").(string),
		"isInstantiation":                    d.Get("is_instantiation").(string),
		"isStaticIP":                         d.Get("is_static_ip").(string),
		"l4L7DeviceApplicationSecurityGroup": d.Get("l4l7_device_application_security_group").(string),
		"l4L7ThirdPartyDevice":               d.Get("l4l7_third_party_device").(string),
		"managed":                            d.Get("managed").(string),
		"maxInstanceCount":                   d.Get("max_instance_count").(string),
		"minInstanceCount":                   d.Get("min_instance_count").(string),
		"name":                               d.Get("name").(string),
		"nameAlias":                          d.Get("name_alias").(string),
		"nativeLBName":                       d.Get("native_lb_name").(string),
		"packageModel":                       d.Get("package_model").(string),
		"promMode":                           d.Get("promiscuous_mode").(string),
		"scheme":                             d.Get("scheme").(string),
		"size":                               d.Get("size").(string),
		"sku":                                d.Get("sku").(string),
		"svcType":                            d.Get("service_type").(string),
		"targetMode":                         d.Get("target_mode").(string),
		"trunking":                           d.Get("trunking").(string),
		"type":                               d.Get("cloud_l4l7_load_balancer_type").(string),
		"Version":                            d.Get("version").(string),
		"status":                             status,
	}
}

func mapCloudRsLDevToCloudSubnetAttrs(annotation, status string, subnetTargetDnList []string) []interface{} {
	cloudSubnetAttrsList := make([]interface{}, len(subnetTargetDnList))
	for index, endPointSelector := range subnetTargetDnList {
		cloudSubnetAttrsMap := map[string]interface{}{
			CloudRsLDevToCloudSubnetClassName: map[string]interface{}{
				"attributes": map[string]string{
					"annotation": annotation,
					"tDn":        endPointSelector,
					"status":     status,
				},
			},
		}
		cloudSubnetAttrsList[index] = cloudSubnetAttrsMap
	}
	return cloudSubnetAttrsList
}

func getAndSetRemoteCloudL4L7LoadBalancerAttributes(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	dnUrl := fmt.Sprintf("%s/%s.json?rsp-subtree=full", client.MOURL, dn)
	cloudLBCont, err := client.GetViaURL(dnUrl)
	if err != nil {
		return nil, err
	}

	cloudLBAttrs := cloudLBCont.S("imdata").Index(0).S(CloudLBClassName).S("attributes")

	cloudLBDistinguishedName := models.StripQuotes(cloudLBAttrs.S("dn").String())
	if cloudLBDistinguishedName == "" {
		return nil, fmt.Errorf("Cloud L4-L7 Load Balancer %s not found", dn)
	}
	d.Set("tenant_dn", GetParentDn(cloudLBDistinguishedName, fmt.Sprintf("/"+RnCloudLB, models.StripQuotes(cloudLBAttrs.S("name").String()))))
	d.Set("version", models.StripQuotes(cloudLBAttrs.S("Version").String()))
	d.Set("active_active", models.StripQuotes(cloudLBAttrs.S("activeActive").String()))
	d.Set("allow_all", models.StripQuotes(cloudLBAttrs.S("allowAll").String()))
	d.Set("annotation", models.StripQuotes(cloudLBAttrs.S("annotation").String()))
	d.Set("auto_scaling", models.StripQuotes(cloudLBAttrs.S("autoScaling").String()))
	d.Set("context_aware", models.StripQuotes(cloudLBAttrs.S("contextAware").String()))
	d.Set("custom_resource_group", models.StripQuotes(cloudLBAttrs.S("customRG").String()))
	d.Set("device_type", models.StripQuotes(cloudLBAttrs.S("devtype").String()))
	d.Set("function_type", models.StripQuotes(cloudLBAttrs.S("funcType").String()))
	d.Set("instance_count", models.StripQuotes(cloudLBAttrs.S("instanceCount").String()))
	d.Set("is_copy", models.StripQuotes(cloudLBAttrs.S("isCopy").String()))
	d.Set("is_instantiation", models.StripQuotes(cloudLBAttrs.S("isInstantiation").String()))
	d.Set("is_static_ip", models.StripQuotes(cloudLBAttrs.S("isStaticIP").String()))
	d.Set("l4l7_device_application_security_group", models.StripQuotes(cloudLBAttrs.S("l4L7DeviceApplicationSecurityGroup").String()))
	d.Set("l4l7_third_party_device", models.StripQuotes(cloudLBAttrs.S("l4L7ThirdPartyDevice").String()))
	d.Set("managed", models.StripQuotes(cloudLBAttrs.S("managed").String()))
	d.Set("max_instance_count", models.StripQuotes(cloudLBAttrs.S("maxInstanceCount").String()))
	d.Set("min_instance_count", models.StripQuotes(cloudLBAttrs.S("minInstanceCount").String()))
	d.Set("mode", models.StripQuotes(cloudLBAttrs.S("mode").String()))
	d.Set("name", models.StripQuotes(cloudLBAttrs.S("name").String()))
	d.Set("native_lb_name", models.StripQuotes(cloudLBAttrs.S("nativeLBName").String()))
	d.Set("package_model", models.StripQuotes(cloudLBAttrs.S("packageModel").String()))
	d.Set("promiscuous_mode", models.StripQuotes(cloudLBAttrs.S("promMode").String()))
	d.Set("scheme", models.StripQuotes(cloudLBAttrs.S("scheme").String()))
	d.Set("size", models.StripQuotes(cloudLBAttrs.S("size").String()))
	d.Set("sku", models.StripQuotes(cloudLBAttrs.S("sku").String()))
	d.Set("service_type", models.StripQuotes(cloudLBAttrs.S("svcType").String()))
	d.Set("target_mode", models.StripQuotes(cloudLBAttrs.S("targetMode").String()))
	d.Set("trunking", models.StripQuotes(cloudLBAttrs.S("trunking").String()))
	d.Set("cloud_l4l7_load_balancer_type", models.StripQuotes(cloudLBAttrs.S("type").String()))

	cloudLBChildCount, err := cloudLBCont.S("imdata").Index(0).S(CloudLBClassName).ArrayCount("children")
	if err != nil {
		return nil, err
	}

	cloudLBChildSubnetList := make([]string, 0)
	cloudLBChildAaaDomainList := make([]string, 0)

	cloudLBChild := cloudLBCont.S("imdata").Index(0).S(CloudLBClassName).S("children")
	for i := 0; i < cloudLBChildCount; i++ {
		cloudSubnetTDn := models.StripQuotes(cloudLBChild.Index(i).S(CloudRsLDevToCloudSubnetClassName).S("attributes").S("tDn").String())
		if cloudSubnetTDn != "{}" {
			cloudLBChildSubnetList = append(cloudLBChildSubnetList, cloudSubnetTDn)
			continue
		}
		aaaDomainName := models.StripQuotes(cloudLBChild.Index(i).S("aaaRbacAnnotation").S("attributes").S("domain").String())
		if aaaDomainName != "{}" {
			cloudLBChildAaaDomainList = append(cloudLBChildAaaDomainList, fmt.Sprintf("uni/userext/domain-%s", aaaDomainName))
			continue
		}
	}
	d.Set("relation_cloud_rs_ldev_to_cloud_subnet", cloudLBChildSubnetList)
	d.Set("aaa_domain_dn", cloudLBChildAaaDomainList)

	d.SetId(cloudLBDistinguishedName)
	return d, nil
}

func resourceAciCloudL4L7LoadBalancerImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	schemaFilled, err := getAndSetRemoteCloudL4L7LoadBalancerAttributes(aciClient, dn, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciCloudL4L7LoadBalancerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Cloud L4-L7 Load Balancer: Beginning Creation")
	aciClient := m.(*client.Client)

	annotation := d.Get("annotation").(string)
	checkDns := make([]string, 0, 1)
	if relationTocloudRsLDevToCloudSubnet, ok := d.GetOk("relation_cloud_rs_ldev_to_cloud_subnet"); ok {
		checkDns = append(checkDns, toStringList(relationTocloudRsLDevToCloudSubnet.(*schema.Set).List())...)
	}

	if aaaDomainDn, ok := d.GetOk("aaa_domain_dn"); ok {
		checkDns = append(checkDns, toStringList(aaaDomainDn.(*schema.Set).List())...)
	}

	d.Partial(true)
	err := checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	cloudLBChildList := make([]interface{}, 0)
	if relationTocloudRsLDevToCloudSubnet, ok := d.GetOk("relation_cloud_rs_ldev_to_cloud_subnet"); ok {
		cloudLBChildList = mapCloudRsLDevToCloudSubnetAttrs(annotation, "created", toStringList(relationTocloudRsLDevToCloudSubnet.(*schema.Set).List()))
	}

	if aaaDomainDn, ok := d.GetOk("aaa_domain_dn"); ok {
		cloudLBChildList = append(cloudLBChildList, mapListOfAaaDomainAttrs("created", toStringList(aaaDomainDn.(*schema.Set).List()))...)
	}

	cloudLBMapAttrs := mapCloudL4L7LoadBalancerAttrs("created", d)
	deleteEmptyValuesfromMap(cloudLBMapAttrs)
	cloudLBMap := map[string]interface{}{
		CloudLBClassName: map[string]interface{}{
			"attributes": cloudLBMapAttrs,
			"children":   cloudLBChildList,
		},
	}

	cloudLBDn := fmt.Sprintf("%s/%s", d.Get("tenant_dn").(string), fmt.Sprintf(RnCloudLB, d.Get("name").(string)))
	err = aciClient.PostObjectConfig(cloudLBDn, cloudLBMap)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(cloudLBDn)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciCloudL4L7LoadBalancerRead(ctx, d, m)
}

func resourceAciCloudL4L7LoadBalancerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Cloud L4-L7 Load Balancer: Beginning Update")
	aciClient := m.(*client.Client)

	annotation := d.Get("annotation").(string)

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_cloud_rs_ldev_to_cloud_subnet") {
		oldRel, newRel := d.GetChange("relation_cloud_rs_ldev_to_cloud_subnet")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		checkDns = append(checkDns, toStringList(newRelSet.Difference(oldRelSet).List())...)
	}

	if d.HasChange("aaa_domain_dn") {
		oldRel, newRel := d.GetChange("aaa_domain_dn")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		checkDns = append(checkDns, toStringList(newRelSet.Difference(oldRelSet).List())...)
	}

	d.Partial(true)
	err := checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	cloudLBChildList := make([]interface{}, 0)

	if d.HasChange("relation_cloud_rs_ldev_to_cloud_subnet") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_cloud_rs_ldev_to_cloud_subnet")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())
		cloudLBChildList = append(cloudLBChildList, mapCloudRsLDevToCloudSubnetAttrs(annotation, "deleted", relToDelete)...)
		cloudLBChildList = append(cloudLBChildList, mapCloudRsLDevToCloudSubnetAttrs(annotation, "created, modified", relToCreate)...)
	}

	if d.HasChange("aaa_domain_dn") {
		oldRel, newRel := d.GetChange("aaa_domain_dn")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())
		cloudLBChildList = append(cloudLBChildList, mapListOfAaaDomainAttrs("deleted", relToDelete)...)
		cloudLBChildList = append(cloudLBChildList, mapListOfAaaDomainAttrs("created, modified", relToCreate)...)
	}

	cloudLBMapAttrs := mapCloudL4L7LoadBalancerAttrs("modified", d)
	deleteEmptyValuesfromMap(cloudLBMapAttrs)
	cloudLBMap := map[string]interface{}{
		CloudLBClassName: map[string]interface{}{
			"attributes": cloudLBMapAttrs,
			"children":   cloudLBChildList,
		},
	}

	cloudLBDn := fmt.Sprintf("%s/%s", d.Get("tenant_dn").(string), fmt.Sprintf(RnCloudLB, d.Get("name").(string)))
	err = aciClient.PostObjectConfig(cloudLBDn, cloudLBMap)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(cloudLBDn)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciCloudL4L7LoadBalancerRead(ctx, d, m)
}

func resourceAciCloudL4L7LoadBalancerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	_, err := getAndSetRemoteCloudL4L7LoadBalancerAttributes(aciClient, dn, d)
	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciCloudL4L7LoadBalancerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, CloudLBClassName)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
