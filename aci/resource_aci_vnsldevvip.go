package aci

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciL4ToL7Devices() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciL4ToL7DevicesCreate,
		UpdateContext: resourceAciL4ToL7DevicesUpdate,
		ReadContext:   resourceAciL4ToL7DevicesRead,
		DeleteContext: resourceAciL4ToL7DevicesDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciL4ToL7DevicesImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"tenant_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"active": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"context_aware": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"multi-Context",
					"single-Context",
				}, false),
			},
			"device_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"CLOUD",
					"PHYSICAL",
					"VIRTUAL",
				}, false),
			},
			"function_type": {
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
			"is_copy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"managed": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "no",
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"legacy-Mode",
				}, false),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"promiscuous_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"service_type": {
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
			"trunking": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"relation_vns_rs_al_dev_to_dom_p": {
				Type:          schema.TypeSet,
				Optional:      true,
				Description:   "Create relation to vmmDomP",
				MaxItems:      1,
				ConflictsWith: []string{"relation_vns_rs_al_dev_to_phys_dom_p"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_dn": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"switching_mode": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Default:  "native",
							ValidateFunc: validation.StringInSlice([]string{
								"AVE",
								"native",
							}, false),
						},
					},
				},
			},
			"relation_vns_rs_al_dev_to_phys_dom_p": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to phys:DomP",
			}})),

		CustomizeDiff: func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
			// Plan time validation.
			device_type := diff.Get("device_type")
			_, virtual_ok := diff.GetOk("relation_vns_rs_al_dev_to_dom_p")
			_, physical_ok := diff.GetOk("relation_vns_rs_al_dev_to_phys_dom_p")
			if device_type.(string) == "VIRTUAL" && !virtual_ok {
				return errors.New(`"relation_vns_rs_al_dev_to_dom_p" is required when "device_type" is VIRTUAL`)
			} else if device_type.(string) == "PHYSICAL" && !physical_ok {
				return errors.New(`"relation_vns_rs_al_dev_to_phys_dom_p" is required when "device_type" is PHYSICAL`)
			}
			return nil
		},
	}

}

func getRemoteL4ToL7Devices(client *client.Client, dn string) (*models.L4ToL7Devices, error) {
	vnsLDevVipCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	vnsLDevVip := models.L4ToL7DevicesFromContainer(vnsLDevVipCont)
	if vnsLDevVip.DistinguishedName == "" {
		return nil, fmt.Errorf("L4ToL7Devices %s not found", dn)
	}
	return vnsLDevVip, nil
}

func setL4ToL7DevicesAttributes(vnsLDevVip *models.L4ToL7Devices, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(vnsLDevVip.DistinguishedName)
	vnsLDevVipMap, err := vnsLDevVip.ToMap()
	if err != nil {
		return d, err
	}
	dn := d.Id()
	if dn != vnsLDevVip.DistinguishedName {
		d.Set("tenant_dn", "")
	} else {
		d.Set("tenant_dn", GetParentDn(vnsLDevVip.DistinguishedName, fmt.Sprintf("/"+models.RnvnsLDevVip, vnsLDevVipMap["name"])))
	}
	d.Set("active", vnsLDevVipMap["activeActive"])
	d.Set("annotation", vnsLDevVipMap["annotation"])
	d.Set("context_aware", vnsLDevVipMap["contextAware"])
	d.Set("device_type", vnsLDevVipMap["devtype"])
	d.Set("function_type", vnsLDevVipMap["funcType"])
	d.Set("is_copy", vnsLDevVipMap["isCopy"])
	d.Set("managed", vnsLDevVipMap["managed"])
	d.Set("mode", vnsLDevVipMap["mode"])
	d.Set("name", vnsLDevVipMap["name"])
	d.Set("promiscuous_mode", vnsLDevVipMap["promMode"])
	d.Set("service_type", vnsLDevVipMap["svcType"])
	d.Set("trunking", vnsLDevVipMap["trunking"])
	d.Set("name_alias", vnsLDevVipMap["nameAlias"])
	return d, nil
}

func resourceAciL4ToL7DevicesImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	vnsLDevVip, err := getRemoteL4ToL7Devices(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setL4ToL7DevicesAttributes(vnsLDevVip, d)
	if err != nil {
		return nil, err
	}
	vnsRsALDevToDomPData, err := aciClient.ReadRelationvnsRsALDevToDomP(dn)
	relParams := make([]map[string]string, 0, 1)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsALDevToDomP %v", err)
		d.Set("relation_vns_rs_al_dev_to_dom_p", relParams)
	} else {
		relParamsList := vnsRsALDevToDomPData.([]map[string]string)
		for _, obj := range relParamsList {
			relParams = append(relParams, map[string]string{
				"switching_mode": obj["switchingMode"],
				"domain_dn":      obj["tDn"],
			})
		}
		d.Set("relation_vns_rs_al_dev_to_dom_p", relParams)
	}

	vnsRsALDevToPhysDomPData, err := aciClient.ReadRelationvnsRsALDevToPhysDomP(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsALDevToPhysDomP %v", err)
		d.Set("relation_vns_rs_al_dev_to_phys_dom_p", "")
	} else {
		d.Set("relation_vns_rs_al_dev_to_phys_dom_p", vnsRsALDevToPhysDomPData.(string))
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciL4ToL7DevicesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] L4ToL7Devices: Beginning Creation")
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	TenantDn := d.Get("tenant_dn").(string)

	vnsLDevVipAttr := models.L4ToL7DevicesAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		vnsLDevVipAttr.Annotation = Annotation.(string)
	} else {
		vnsLDevVipAttr.Annotation = "{}"
	}

	if Active, ok := d.GetOk("active"); ok {
		vnsLDevVipAttr.Active = Active.(string)
	}

	if ContextAware, ok := d.GetOk("context_aware"); ok {
		vnsLDevVipAttr.ContextAware = ContextAware.(string)
	}

	if DeviceType, ok := d.GetOk("device_type"); ok {
		vnsLDevVipAttr.Devtype = DeviceType.(string)
	}

	if FunctionType, ok := d.GetOk("function_type"); ok {
		vnsLDevVipAttr.FuncType = FunctionType.(string)
	}

	if IsCopy, ok := d.GetOk("is_copy"); ok {
		vnsLDevVipAttr.IsCopy = IsCopy.(string)
	}

	if Managed, ok := d.GetOk("managed"); ok {
		vnsLDevVipAttr.Managed = Managed.(string)
	}

	if Mode, ok := d.GetOk("mode"); ok {
		vnsLDevVipAttr.Mode = Mode.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		vnsLDevVipAttr.Name = Name.(string)
	}

	if PromiscuousMode, ok := d.GetOk("promiscuous_mode"); ok {
		vnsLDevVipAttr.PromMode = PromiscuousMode.(string)
	}

	if ServiceType, ok := d.GetOk("service_type"); ok {
		vnsLDevVipAttr.SvcType = ServiceType.(string)
	}

	if Trunking, ok := d.GetOk("trunking"); ok {
		vnsLDevVipAttr.Trunking = Trunking.(string)
	}
	vnsLDevVip := models.NewL4ToL7Devices(fmt.Sprintf(models.RnvnsLDevVip, name), TenantDn, nameAlias, vnsLDevVipAttr)

	err := aciClient.Save(vnsLDevVip)
	if err != nil {
		return diag.FromErr(err)
	}
	checkDns := make([]string, 0, 1)

	if relationTovnsRsALDevToDomP, ok := d.GetOk("relation_vns_rs_al_dev_to_dom_p"); ok {
		relationParamMap := relationTovnsRsALDevToDomP.(*schema.Set).List()
		for _, relationParam := range relationParamMap {
			innerMap := relationParam.(map[string]interface{})
			checkDns = append(checkDns, innerMap["domain_dn"].(string))
		}
	}

	if relationTovnsRsALDevToPhysDomP, ok := d.GetOk("relation_vns_rs_al_dev_to_phys_dom_p"); ok {
		relationParam := relationTovnsRsALDevToPhysDomP.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTovnsRsALDevToDomP, ok := d.GetOk("relation_vns_rs_al_dev_to_dom_p"); ok {
		relationParamMap := relationTovnsRsALDevToDomP.(*schema.Set).List()
		var switchingMode, domainDN string
		for _, relationParam := range relationParamMap {
			innerMap := relationParam.(map[string]interface{})
			switchingMode = innerMap["switching_mode"].(string)
			domainDN = innerMap["domain_dn"].(string)
		}
		err = aciClient.CreateRelationvnsRsALDevToDomP(vnsLDevVip.DistinguishedName, vnsLDevVipAttr.Annotation, switchingMode, domainDN)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if relationTovnsRsALDevToPhysDomP, ok := d.GetOk("relation_vns_rs_al_dev_to_phys_dom_p"); ok {
		relationParam := relationTovnsRsALDevToPhysDomP.(string)
		err = aciClient.CreateRelationvnsRsALDevToPhysDomP(vnsLDevVip.DistinguishedName, vnsLDevVipAttr.Annotation, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(vnsLDevVip.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciL4ToL7DevicesRead(ctx, d, m)
}
func resourceAciL4ToL7DevicesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] L4ToL7Devices: Beginning Update")
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	TenantDn := d.Get("tenant_dn").(string)

	vnsLDevVipAttr := models.L4ToL7DevicesAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		vnsLDevVipAttr.Annotation = Annotation.(string)
	} else {
		vnsLDevVipAttr.Annotation = "{}"
	}

	if Active, ok := d.GetOk("active"); ok {
		vnsLDevVipAttr.Active = Active.(string)
	}

	if ContextAware, ok := d.GetOk("context_aware"); ok {
		vnsLDevVipAttr.ContextAware = ContextAware.(string)
	}

	if DeviceType, ok := d.GetOk("device_type"); ok {
		vnsLDevVipAttr.Devtype = DeviceType.(string)
	}

	if FunctionType, ok := d.GetOk("function_type"); ok {
		vnsLDevVipAttr.FuncType = FunctionType.(string)
	}

	if IsCopy, ok := d.GetOk("is_copy"); ok {
		vnsLDevVipAttr.IsCopy = IsCopy.(string)
	}

	if Managed, ok := d.GetOk("managed"); ok {
		vnsLDevVipAttr.Managed = Managed.(string)
	}

	if Mode, ok := d.GetOk("mode"); ok {
		vnsLDevVipAttr.Mode = Mode.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		vnsLDevVipAttr.Name = Name.(string)
	}

	if PromiscuousMode, ok := d.GetOk("promiscuous_mode"); ok {
		vnsLDevVipAttr.PromMode = PromiscuousMode.(string)
	}

	if ServiceType, ok := d.GetOk("service_type"); ok {
		vnsLDevVipAttr.SvcType = ServiceType.(string)
	}

	if Trunking, ok := d.GetOk("trunking"); ok {
		vnsLDevVipAttr.Trunking = Trunking.(string)
	}
	vnsLDevVip := models.NewL4ToL7Devices(fmt.Sprintf(models.RnvnsLDevVip, name), TenantDn, nameAlias, vnsLDevVipAttr)

	vnsLDevVip.Status = "modified"

	err := aciClient.Save(vnsLDevVip)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_vns_rs_al_dev_to_dom_p") || d.HasChange("annotation") {
		_, newRel := d.GetChange("relation_vns_rs_al_dev_to_dom_p")
		newRelSet := newRel.(*schema.Set).List()
		for _, relationParam := range newRelSet {
			innerMap := relationParam.(map[string]interface{})
			checkDns = append(checkDns, innerMap["domain_dn"].(string))
		}
	}

	if d.HasChange("relation_vns_rs_al_dev_to_phys_dom_p") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_vns_rs_al_dev_to_phys_dom_p")
		checkDns = append(checkDns, newRelParam.(string))

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_vns_rs_al_dev_to_dom_p") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_vns_rs_al_dev_to_dom_p")
		newRelParamMap := newRelParam.(*schema.Set).List()
		var switchingMode, domainDN string
		for _, newRelationParam := range newRelParamMap {
			innerMap := newRelationParam.(map[string]interface{})
			switchingMode = innerMap["switching_mode"].(string)
			domainDN = innerMap["domain_dn"].(string)
		}
		err = aciClient.DeleteRelationvnsRsALDevToDomP(vnsLDevVip.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvnsRsALDevToDomP(vnsLDevVip.DistinguishedName, vnsLDevVipAttr.Annotation, switchingMode, domainDN)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("relation_vns_rs_al_dev_to_phys_dom_p") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_vns_rs_al_dev_to_phys_dom_p")
		err = aciClient.DeleteRelationvnsRsALDevToPhysDomP(vnsLDevVip.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvnsRsALDevToPhysDomP(vnsLDevVip.DistinguishedName, vnsLDevVipAttr.Annotation, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(vnsLDevVip.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciL4ToL7DevicesRead(ctx, d, m)
}

func resourceAciL4ToL7DevicesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	vnsLDevVip, err := getRemoteL4ToL7Devices(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}

	_, err = setL4ToL7DevicesAttributes(vnsLDevVip, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	vnsRsALDevToDomPData, err := aciClient.ReadRelationvnsRsALDevToDomP(dn)
	relParams := make([]map[string]string, 0, 1)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsALDevToDomP %v", err)
		setRelationAttribute(d, "relation_vns_rs_al_dev_to_dom_p", relParams)
	} else {
		relParamsList := vnsRsALDevToDomPData.([]map[string]string)
		for _, obj := range relParamsList {
			relParams = append(relParams, map[string]string{
				"switching_mode": obj["switchingMode"],
				"domain_dn":      obj["tDn"],
			})
		}
		setRelationAttribute(d, "relation_vns_rs_al_dev_to_dom_p", relParams)
	}

	vnsRsALDevToPhysDomPData, err := aciClient.ReadRelationvnsRsALDevToPhysDomP(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsALDevToPhysDomP %v", err)
		setRelationAttribute(d, "relation_vns_rs_al_dev_to_phys_dom_p", "")
	} else {
		setRelationAttribute(d, "relation_vns_rs_al_dev_to_phys_dom_p", vnsRsALDevToPhysDomPData.(string))
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciL4ToL7DevicesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "vnsLDevVip")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
