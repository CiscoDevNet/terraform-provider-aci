package aci

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	CloudLDevClassName        = "cloudLDev"
	CloudLIfClassName         = "cloudLIf"
	CloudEPSelectorClassName  = "cloudEPSelector"
	CloudRsLDevToCtxClassName = "cloudRsLDevToCtx"
	RnCloudLDev               = "cld-%s"
)

func resourceAciCloudL4L7Device() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciCloudL4L7DeviceCreate,
		UpdateContext: resourceAciCloudL4L7DeviceUpdate,
		ReadContext:   resourceAciCloudL4L7DeviceRead,
		DeleteContext: resourceAciCloudL4L7DeviceDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciCloudL4L7DeviceImport,
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
			"mode": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"relation_cloud_rs_ldev_to_ctx": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Create relation to fv:Ctx",
			},
			"aaa_domain_dn": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
				Set:      schema.HashString,
			},

			"interface_selectors": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: AppendAttrSchemas(map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
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
						"end_point_selectors": &schema.Schema{
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: AppendAttrSchemas(map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"match_expression": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
								}),
							},
						},
					}),
				},
			},
		}, GetNameAliasAttrSchema(), GetAnnotationAttrSchema()),
	}
}

func mapCloudL4L7DeviceAttrs(status string, d *schema.ResourceData) map[string]string {
	return map[string]string{
		"annotation":                         d.Get("annotation").(string),
		"Version":                            d.Get("version").(string),
		"activeActive":                       d.Get("active_active").(string),
		"contextAware":                       d.Get("context_aware").(string),
		"customRG":                           d.Get("custom_resource_group").(string),
		"devtype":                            d.Get("device_type").(string),
		"funcType":                           d.Get("function_type").(string),
		"instanceCount":                      d.Get("instance_count").(string),
		"isCopy":                             d.Get("is_copy").(string),
		"isInstantiation":                    d.Get("is_instantiation").(string),
		"l4L7DeviceApplicationSecurityGroup": d.Get("l4l7_device_application_security_group").(string),
		"l4L7ThirdPartyDevice":               d.Get("l4l7_third_party_device").(string),
		"managed":                            d.Get("managed").(string),
		"name":                               d.Get("name").(string),
		"nameAlias":                          d.Get("name_alias").(string),
		"packageModel":                       d.Get("package_model").(string),
		"promMode":                           d.Get("promiscuous_mode").(string),
		"svcType":                            d.Get("service_type").(string),
		"targetMode":                         d.Get("target_mode").(string),
		"trunking":                           d.Get("trunking").(string),
		"status":                             status,
	}
}

func mapCloudRsLDevToCtxAttrs(annotation, status string, vrf_dn string) map[string]interface{} {
	return map[string]interface{}{
		CloudRsLDevToCtxClassName: map[string]interface{}{
			"attributes": map[string]string{
				"tDn":    vrf_dn,
				"status": status,
			},
		},
	}
}

func mapCloudInterfaceSelectorAttrs(annotation, status string, interfaceSelectors []interface{}) []interface{} {
	cloudLIfAttrsMapList := make([]interface{}, 0)
	for _, interfaceSelector := range interfaceSelectors {
		interfaceSelectorAttrMap := interfaceSelector.(map[string]interface{})
		cloudLIfAttrsMap := map[string]string{
			"name":     interfaceSelectorAttrMap["name"].(string),
			"allowAll": interfaceSelectorAttrMap["allow_all"].(string),
			"status":   status,
		}

		cloudEPSelectorsMap := make([]map[string]interface{}, 0)
		for _, cloudEPSelector := range interfaceSelectorAttrMap["end_point_selectors"].(*schema.Set).List() {
			cloudEPSelectorMap := cloudEPSelector.(map[string]interface{})
			epStatus := ""
			if cloudEPSelectorMap["status"] != nil {
				epStatus = cloudEPSelectorMap["status"].(string)
			}
			cloudEPSelectorAttrsMap := map[string]interface{}{
				CloudEPSelectorClassName: map[string]interface{}{
					"attributes": map[string]string{
						"name":            cloudEPSelectorMap["name"].(string),
						"matchExpression": cloudEPSelectorMap["match_expression"].(string),
						"status":          epStatus,
					},
				},
			}
			cloudEPSelectorsMap = append(cloudEPSelectorsMap, cloudEPSelectorAttrsMap)
		}
		cloudLIfAttrsMapList = append(cloudLIfAttrsMapList, map[string]interface{}{
			CloudLIfClassName: map[string]interface{}{
				"attributes": cloudLIfAttrsMap,
				"children":   cloudEPSelectorsMap,
			},
		})
	}
	return cloudLIfAttrsMapList
}

func getAndSetRemoteCloudL4L7DeviceAttributes(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	dnUrl := fmt.Sprintf("%s/%s.json?rsp-subtree=full", client.MOURL, dn)
	cloudLDevCont, err := client.GetViaURL(dnUrl)
	if err != nil {
		return nil, err
	}

	cloudLDevAttrs := cloudLDevCont.S("imdata").Index(0).S(CloudLDevClassName).S("attributes")

	cloudLDevDistinguishedName := models.StripQuotes(cloudLDevAttrs.S("dn").String())
	if cloudLDevDistinguishedName == "" {
		return nil, fmt.Errorf("Cloud L4-L7 Third Party Device %s not found", dn)
	}
	d.Set("tenant_dn", GetParentDn(cloudLDevDistinguishedName, fmt.Sprintf("/"+RnCloudLDev, models.StripQuotes(cloudLDevAttrs.S("name").String()))))

	d.Set("annotation", models.StripQuotes(cloudLDevAttrs.S("annotation").String()))
	d.Set("active_active", models.StripQuotes(cloudLDevAttrs.S("activeActive").String()))
	d.Set("context_aware", models.StripQuotes(cloudLDevAttrs.S("contextAware").String()))
	d.Set("custom_resource_group", models.StripQuotes(cloudLDevAttrs.S("customRG").String()))
	d.Set("device_type", models.StripQuotes(cloudLDevAttrs.S("devtype").String()))
	d.Set("function_type", models.StripQuotes(cloudLDevAttrs.S("funcType").String()))
	d.Set("instance_count", models.StripQuotes(cloudLDevAttrs.S("instanceCount").String()))
	d.Set("is_copy", models.StripQuotes(cloudLDevAttrs.S("isCopy").String()))
	d.Set("is_instantiation", models.StripQuotes(cloudLDevAttrs.S("isInstantiation").String()))
	d.Set("l4l7_device_application_security_group", models.StripQuotes(cloudLDevAttrs.S("l4L7DeviceApplicationSecurityGroup").String()))
	d.Set("l4l7_third_party_device", models.StripQuotes(cloudLDevAttrs.S("l4L7ThirdPartyDevice").String()))
	d.Set("managed", models.StripQuotes(cloudLDevAttrs.S("managed").String()))
	d.Set("mode", models.StripQuotes(cloudLDevAttrs.S("mode").String()))
	d.Set("name", models.StripQuotes(cloudLDevAttrs.S("name").String()))
	d.Set("name_alias", models.StripQuotes(cloudLDevAttrs.S("nameAlias").String()))
	d.Set("package_model", models.StripQuotes(cloudLDevAttrs.S("packageModel").String()))
	d.Set("promiscuous_mode", models.StripQuotes(cloudLDevAttrs.S("promMode").String()))
	d.Set("service_type", models.StripQuotes(cloudLDevAttrs.S("svcType").String()))
	d.Set("target_mode", models.StripQuotes(cloudLDevAttrs.S("targetMode").String()))
	d.Set("trunking", models.StripQuotes(cloudLDevAttrs.S("trunking").String()))
	d.Set("version", models.StripQuotes(cloudLDevAttrs.S("Version").String()))

	cloudLDevChildCount, err := cloudLDevCont.S("imdata").Index(0).S(CloudLDevClassName).ArrayCount("children")
	if err != nil {
		return nil, err
	}

	cloudLDevChildAaaDomainList := make([]string, 0)
	cloudLIfChildList := make([]interface{}, 0)

	cloudLDevChild := cloudLDevCont.S("imdata").Index(0).S(CloudLDevClassName).S("children")
	for i := 0; i < cloudLDevChildCount; i++ {

		aaaDomainName := models.StripQuotes(cloudLDevChild.Index(i).S("aaaRbacAnnotation").S("attributes").S("domain").String())
		if aaaDomainName != "{}" {
			cloudLDevChildAaaDomainList = append(cloudLDevChildAaaDomainList, fmt.Sprintf("uni/userext/domain-%s", aaaDomainName))
			continue
		}

		cloudRsLDevToCtxTDn := models.StripQuotes(cloudLDevChild.Index(i).S(CloudRsLDevToCtxClassName).S("attributes").S("tDn").String())
		if cloudRsLDevToCtxTDn != "{}" {
			d.Set("relation_cloud_rs_ldev_to_ctx", cloudRsLDevToCtxTDn)
			continue
		}

		cloudLIfAttrs := cloudLDevChild.Index(i).S(CloudLIfClassName).S("attributes")
		cloudLIfName := models.StripQuotes(cloudLIfAttrs.S("name").String())

		if cloudLIfName != "{}" {
			logicalInterfaceAttrsMap := map[string]interface{}{
				"name":      cloudLIfName,
				"allow_all": models.StripQuotes(cloudLIfAttrs.S("allowAll").String()),
			}
			cloudLIfChildCount, err := cloudLDevChild.Index(i).S(CloudLIfClassName).ArrayCount("children")
			if err != nil {
				return nil, err
			}
			cloudEPSelectorAttrsList := make([]map[string]string, 0)
			for j := 0; j < cloudLIfChildCount; j++ {
				cloudEPSelectorAttrs := cloudLDevChild.Index(i).S(CloudLIfClassName).S("children").Index(j).S(CloudEPSelectorClassName).S("attributes")
				cloudEPSelectorAttrsMap := map[string]string{
					"name":             models.StripQuotes(cloudEPSelectorAttrs.S("name").String()),
					"match_expression": models.StripQuotes(cloudEPSelectorAttrs.S("matchExpression").String()),
				}
				cloudEPSelectorAttrsList = append(cloudEPSelectorAttrsList, cloudEPSelectorAttrsMap)
			}
			logicalInterfaceAttrsMap["end_point_selectors"] = cloudEPSelectorAttrsList
			cloudLIfChildList = append(cloudLIfChildList, logicalInterfaceAttrsMap)
		}
	}

	d.Set("interface_selectors", cloudLIfChildList)
	d.Set("aaa_domain_dn", cloudLDevChildAaaDomainList)
	d.SetId(cloudLDevDistinguishedName)
	return d, nil
}

func resourceAciCloudL4L7DeviceImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	schemaFilled, err := getAndSetRemoteCloudL4L7DeviceAttributes(aciClient, dn, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciCloudL4L7DeviceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] CloudL4L7Device: Beginning Creation")

	aciClient := m.(*client.Client)

	annotation := d.Get("annotation").(string)

	checkDns := make([]string, 0, 1)
	if relationCloudRsLDevToCtx, ok := d.GetOk("relation_cloud_rs_ldev_to_ctx"); ok {
		checkDns = append(checkDns, relationCloudRsLDevToCtx.(string))
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

	cloudLDevChildList := make([]interface{}, 0)
	// VRF Mapping
	cloudLDevChildList = append(cloudLDevChildList, mapCloudRsLDevToCtxAttrs(annotation, "created", d.Get("relation_cloud_rs_ldev_to_ctx").(string)))
	// AAA Domain Mapping
	if aaaDomainDn, ok := d.GetOk("aaa_domain_dn"); ok {
		cloudLDevChildList = append(cloudLDevChildList, mapListOfAaaDomainAttrs("created", toStringList(aaaDomainDn.(*schema.Set).List()))...)
	}
	// Interface Mapping
	if interfaceSelectors, ok := d.GetOk("interface_selectors"); ok {
		cloudLDevChildList = append(cloudLDevChildList, mapCloudInterfaceSelectorAttrs(annotation, "created", interfaceSelectors.(*schema.Set).List())...)
	}

	cloudLDevMapAttrs := mapCloudL4L7DeviceAttrs("created", d)
	deleteEmptyValuesfromMap(cloudLDevMapAttrs)
	cloudLDevMap := map[string]interface{}{
		CloudLDevClassName: map[string]interface{}{
			"attributes": cloudLDevMapAttrs,
			"children":   cloudLDevChildList,
		},
	}

	cloudLDevDn := fmt.Sprintf("%s/%s", d.Get("tenant_dn").(string), fmt.Sprintf(RnCloudLDev, d.Get("name").(string)))
	err = aciClient.PostObjectConfig(cloudLDevDn, cloudLDevMap)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(cloudLDevDn)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciCloudL4L7DeviceRead(ctx, d, m)
}
func resourceAciCloudL4L7DeviceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] CloudL4L7Device: Beginning Update")
	aciClient := m.(*client.Client)

	annotation := d.Get("annotation").(string)

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_cloud_rs_ldev_to_ctx") {
		if relationCloudRsLDevToCtx, ok := d.GetOk("relation_cloud_rs_ldev_to_ctx"); ok {
			checkDns = append(checkDns, relationCloudRsLDevToCtx.(string))
		}
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

	cloudLDevChildList := make([]interface{}, 0)

	// VRF Mapping
	if d.HasChange("relation_cloud_rs_ldev_to_ctx") || d.HasChange("annotation") {
		cloudLDevChildList = append(cloudLDevChildList, mapCloudRsLDevToCtxAttrs(annotation, "created,modified", d.Get("relation_cloud_rs_ldev_to_ctx").(string)))
	}

	// AAA Domain Mapping
	if d.HasChange("aaa_domain_dn") {
		oldRel, newRel := d.GetChange("aaa_domain_dn")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())
		cloudLDevChildList = append(cloudLDevChildList, mapListOfAaaDomainAttrs("deleted", relToDelete)...)
		cloudLDevChildList = append(cloudLDevChildList, mapListOfAaaDomainAttrs("created,modified", relToCreate)...)
	}

	// Interface Mapping
	if d.HasChange("interface_selectors") || d.HasChange("annotation") {
		oldInf, newInf := d.GetChange("interface_selectors")

		oldInfList := oldInf.(*schema.Set).List()
		newInfList := newInf.(*schema.Set).List()

		oldInfPayload := mapCloudInterfaceSelectorAttrs(annotation, "deleted", oldInfList)
		newInfPayload := mapCloudInterfaceSelectorAttrs(annotation, "created,modified", newInfList)

		newInfPayloadTemp := make([]interface{}, 0)
		newInfPayloadNames := make(map[string]bool, 0)

		for _, newInf := range newInfPayload {
			newInfPayloadNames[newInf.(map[string]interface{})[CloudLIfClassName].(map[string]interface{})["attributes"].(map[string]string)["name"]] = true
		}

		// Interface level check
		for _, oldInf := range oldInfPayload {
			oldInfName := oldInf.(map[string]interface{})[CloudLIfClassName].(map[string]interface{})["attributes"].(map[string]string)["name"]
			if newInfPayloadNames[oldInfName] {
				for _, newInf := range newInfPayload {
					if reflect.DeepEqual(oldInf, newInf) {
						break
					}

					newInfName := newInf.(map[string]interface{})[CloudLIfClassName].(map[string]interface{})["attributes"].(map[string]string)["name"]
					if oldInfName != newInfName {
						continue
					}

					if oldInfName == newInfName {
						oldInfEPList := oldInf.(map[string]interface{})[CloudLIfClassName].(map[string]interface{})["children"]
						newInfEPList := newInf.(map[string]interface{})[CloudLIfClassName].(map[string]interface{})["children"]
						newInfEPNames := make(map[string]bool, 0)

						// End Point Selectors level check
						for _, newInfEP := range newInfEPList.([]map[string]interface{}) {
							newInfEPNames[newInfEP[CloudEPSelectorClassName].(map[string]interface{})["attributes"].(map[string]string)["name"]] = true
						}

						for _, oldInfEP := range oldInfEPList.([]map[string]interface{}) {
							oldInfEPAttrs := oldInfEP[CloudEPSelectorClassName].(map[string]interface{})["attributes"].(map[string]string)
							if !newInfEPNames[oldInfEPAttrs["name"]] {
								oldInfEPAttrs["status"] = "deleted"
								newInf.(map[string]interface{})[CloudLIfClassName].(map[string]interface{})["children"] = append(newInf.(map[string]interface{})[CloudLIfClassName].(map[string]interface{})["children"].([]map[string]interface{}), oldInfEP)
							}
						}
					}
				}
			} else {
				newInfPayloadTemp = append(newInfPayloadTemp, oldInf)
			}
		}
		cloudLDevChildList = append(cloudLDevChildList, newInfPayloadTemp...)
		cloudLDevChildList = append(cloudLDevChildList, newInfPayload...)
	}

	cloudLDevMapAttrs := mapCloudL4L7DeviceAttrs("created,modified", d)
	deleteEmptyValuesfromMap(cloudLDevMapAttrs)

	cloudLDevMap := map[string]interface{}{
		CloudLDevClassName: map[string]interface{}{
			"attributes": cloudLDevMapAttrs,
			"children":   cloudLDevChildList,
		},
	}

	cloudLDevDn := fmt.Sprintf("%s/%s", d.Get("tenant_dn").(string), fmt.Sprintf(RnCloudLDev, d.Get("name").(string)))
	err = aciClient.PostObjectConfig(cloudLDevDn, cloudLDevMap)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(cloudLDevDn)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciCloudL4L7DeviceRead(ctx, d, m)
}

func resourceAciCloudL4L7DeviceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	_, err := getAndSetRemoteCloudL4L7DeviceAttributes(aciClient, dn, d)
	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciCloudL4L7DeviceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, CloudLDevClassName)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
