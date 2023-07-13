package aci

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciPIMInterfacePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciPIMInterfacePolicyCreate,
		UpdateContext: resourceAciPIMInterfacePolicyUpdate,
		ReadContext:   resourceAciPIMInterfacePolicyRead,
		DeleteContext: resourceAciPIMInterfacePolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciPIMInterfacePolicyImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"tenant_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"auth_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ah-md5",
					"none",
				}, false),
			},
			"control_state": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"border",
						"passive",
						"strict-rfc-compliant",
					}, false),
				},
				DiffSuppressFunc: suppressTypeListDiffFunc,
			},
			"designated_router_delay": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"designated_router_priority": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"hello_interval": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"join_prune_interval": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"inbound_join_prune_filter_policy": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"outbound_join_prune_filter_policy": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"neighbor_filter_policy": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		})),
	}
}

func getRemotePIMInterfacePolicy(client *client.Client, dn string) (*models.PIMInterfacePolicy, error) {
	pimIfPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	pimIfPol := models.PIMInterfacePolicyFromContainer(pimIfPolCont)
	if pimIfPol.DistinguishedName == "" {
		return nil, fmt.Errorf("PIM Interface Policy %s not found", dn)
	}
	return pimIfPol, nil
}

func setPIMInterfacePolicyAttributes(pimIfPol *models.PIMInterfacePolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(pimIfPol.DistinguishedName)
	d.Set("description", pimIfPol.Description)
	pimIfPolMap, err := pimIfPol.ToMap()
	if err != nil {
		return d, err
	}
	dn := d.Id()
	if dn != pimIfPol.DistinguishedName {
		d.Set("tenant_dn", "")
	} else {
		d.Set("tenant_dn", GetParentDn(pimIfPol.DistinguishedName, fmt.Sprintf("/"+models.RnPimIfPol, pimIfPolMap["name"])))
	}
	d.Set("annotation", pimIfPolMap["annotation"])
	d.Set("auth_type", pimIfPolMap["authT"])
	control_stateGet := make([]string, 0, 1)
	for _, val := range strings.Split(pimIfPolMap["ctrl"], ",") {
		control_stateGet = append(control_stateGet, strings.Trim(val, " "))
	}
	sort.Strings(control_stateGet)
	if control_stateIntr, ok := d.GetOk("control_state"); ok {
		control_stateAct := make([]string, 0, 1)
		for _, val := range control_stateIntr.([]interface{}) {
			control_stateAct = append(control_stateAct, val.(string))
		}
		sort.Strings(control_stateAct)
		if reflect.DeepEqual(control_stateAct, control_stateGet) {
			d.Set("control_state", d.Get("control_state").([]interface{}))
		} else {
			d.Set("control_state", control_stateGet)
		}
	} else {
		d.Set("control_state", control_stateGet)
	}
	d.Set("designated_router_delay", pimIfPolMap["drDelay"])
	d.Set("designated_router_priority", pimIfPolMap["drPrio"])
	d.Set("hello_interval", pimIfPolMap["helloItvl"])
	d.Set("join_prune_interval", pimIfPolMap["jpInterval"])
	d.Set("name", pimIfPolMap["name"])
	d.Set("name_alias", pimIfPolMap["nameAlias"])
	return d, nil
}

func getandSetPIMIfPolRelationshipAttributes(aciClient *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	jpInbFilterPolData, err := aciClient.ReadRelationPIMJPInbFilterPolrtdmcRsFilterToRtMapPol(fmt.Sprintf("%s/%s", dn, models.RnPimJPInbFilterPol))
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation jpInbFilterPolData %v", err)
		d.Set("inbound_join_prune_filter_policy", "")
	} else {
		d.Set("inbound_join_prune_filter_policy", jpInbFilterPolData.(string))
	}

	jpOutbFilterPolData, err := aciClient.ReadRelationPIMJPOutbFilterPolrtdmcRsFilterToRtMapPol(fmt.Sprintf("%s/%s", dn, models.RnPimJPOutbFilterPol))
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation jpOutbFilterPol %v", err)
		d.Set("outbound_join_prune_filter_policy", "")
	} else {
		d.Set("outbound_join_prune_filter_policy", jpOutbFilterPolData.(string))
	}

	neighborFilterPolData, err := aciClient.ReadRelationPIMNbrFilterPolrtdmcRsFilterToRtMapPol(fmt.Sprintf("%s/%s", dn, models.RnPimNbrFilterPol))
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation neighborFilterPol %v", err)
		d.Set("neighbor_filter_policy", "")
	} else {
		d.Set("neighbor_filter_policy", neighborFilterPolData.(string))
	}
	return d, nil
}

func resourceAciPIMInterfacePolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	pimIfPol, err := getRemotePIMInterfacePolicy(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setPIMInterfacePolicyAttributes(pimIfPol, d)
	if err != nil {
		return nil, err
	}

	_, err = getandSetPIMIfPolRelationshipAttributes(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] PimIfPol Relationship Attributes - Read finished successfully")
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciPIMInterfacePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] PIMInterfacePolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	TenantDn := d.Get("tenant_dn").(string)

	pimIfPolAttr := models.PIMInterfacePolicyAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		pimIfPolAttr.Annotation = Annotation.(string)
	} else {
		pimIfPolAttr.Annotation = "{}"
	}

	if AuthT, ok := d.GetOk("auth_type"); ok {
		pimIfPolAttr.AuthT = AuthT.(string)
	}

	if Control_state, ok := d.GetOk("control_state"); ok {
		control_stateList := make([]string, 0, 1)
		for _, val := range Control_state.([]interface{}) {
			control_stateList = append(control_stateList, val.(string))
		}
		Control_state := strings.Join(control_stateList, ",")
		pimIfPolAttr.Ctrl = Control_state
	}

	if DrDelay, ok := d.GetOk("designated_router_delay"); ok {
		pimIfPolAttr.DrDelay = DrDelay.(string)
	}

	if DrPrio, ok := d.GetOk("designated_router_priority"); ok {
		pimIfPolAttr.DrPrio = DrPrio.(string)
	}

	if HelloItvl, ok := d.GetOk("hello_interval"); ok {
		pimIfPolAttr.HelloItvl = HelloItvl.(string)
	}

	if JpInterval, ok := d.GetOk("join_prune_interval"); ok {
		pimIfPolAttr.JpInterval = JpInterval.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		pimIfPolAttr.Name = Name.(string)
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		pimIfPolAttr.NameAlias = NameAlias.(string)
	}

	pimIfPol := models.NewPIMInterfacePolicy(fmt.Sprintf(models.RnPimIfPol, name), TenantDn, desc, pimIfPolAttr)

	err := aciClient.Save(pimIfPol)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)
	if jpInboundFilterPolMap, ok := d.GetOk("inbound_join_prune_filter_policy"); ok {
		checkDns = append(checkDns, jpInboundFilterPolMap.(string))
	}

	if jpOutboundFilterPolMap, ok := d.GetOk("outbound_join_prune_filter_policy"); ok {
		checkDns = append(checkDns, jpOutboundFilterPolMap.(string))
	}

	if neighborFilterPolMap, ok := d.GetOk("neighbor_filter_policy"); ok {
		checkDns = append(checkDns, neighborFilterPolMap.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if jpInboundFilterPolMap, ok := d.GetOk("inbound_join_prune_filter_policy"); ok {
		jpInboundFilterPol := models.NewPIMJPInboundFilterPolicy(pimIfPol.DistinguishedName, "", models.PIMJPInboundFilterPolicyAttributes{})
		err := aciClient.Save(jpInboundFilterPol)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationPIMJPInbFilterPolrtdmcRsFilterToRtMapPol(jpInboundFilterPol.DistinguishedName, jpInboundFilterPolMap.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if jpOutboundFilterPolMap, ok := d.GetOk("outbound_join_prune_filter_policy"); ok {
		jpOutboundFilterPol := models.NewPIMJPOutboundFilterPolicy(pimIfPol.DistinguishedName, "", models.PIMJPOutboundFilterPolicyAttributes{})
		err := aciClient.Save(jpOutboundFilterPol)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationPIMJPOutbFilterPolrtdmcRsFilterToRtMapPol(jpOutboundFilterPol.DistinguishedName, jpOutboundFilterPolMap.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if neighborFilterPolMap, ok := d.GetOk("neighbor_filter_policy"); ok {
		neighborFilterPol := models.NewPIMNeighborFiterPolicy(pimIfPol.DistinguishedName, "", models.PIMNeighborFiterPolicyAttributes{})
		err := aciClient.Save(neighborFilterPol)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationPIMNbrFilterPolrtdmcRsFilterToRtMapPol(neighborFilterPol.DistinguishedName, neighborFilterPolMap.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(pimIfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciPIMInterfacePolicyRead(ctx, d, m)
}

func resourceAciPIMInterfacePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] PIMInterfacePolicy: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	TenantDn := d.Get("tenant_dn").(string)

	pimIfPolAttr := models.PIMInterfacePolicyAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		pimIfPolAttr.Annotation = Annotation.(string)
	} else {
		pimIfPolAttr.Annotation = "{}"
	}

	if AuthT, ok := d.GetOk("auth_type"); ok {
		pimIfPolAttr.AuthT = AuthT.(string)
	}
	if Control_state, ok := d.GetOk("control_state"); ok {
		control_stateList := make([]string, 0, 1)
		for _, val := range Control_state.([]interface{}) {
			control_stateList = append(control_stateList, val.(string))
		}
		Control_state := strings.Join(control_stateList, ",")
		pimIfPolAttr.Ctrl = Control_state
	}

	if DrDelay, ok := d.GetOk("designated_router_delay"); ok {
		pimIfPolAttr.DrDelay = DrDelay.(string)
	}

	if DrPrio, ok := d.GetOk("designated_router_priority"); ok {
		pimIfPolAttr.DrPrio = DrPrio.(string)
	}

	if HelloItvl, ok := d.GetOk("hello_interval"); ok {
		pimIfPolAttr.HelloItvl = HelloItvl.(string)
	}

	if JpInterval, ok := d.GetOk("join_prune_interval"); ok {
		pimIfPolAttr.JpInterval = JpInterval.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		pimIfPolAttr.Name = Name.(string)
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		pimIfPolAttr.NameAlias = NameAlias.(string)
	}

	pimIfPol := models.NewPIMInterfacePolicy(fmt.Sprintf(models.RnPimIfPol, name), TenantDn, desc, pimIfPolAttr)

	pimIfPol.Status = "modified"

	err := aciClient.Save(pimIfPol)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)
	if d.HasChange("inbound_join_prune_filter_policy") {
		_, newRelParam := d.GetChange("inbound_join_prune_filter_policy")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("outbound_join_prune_filter_policy") {
		_, newRelParam := d.GetChange("outbound_join_prune_filter_policy")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("neighbor_filter_policy") {
		_, newRelParam := d.GetChange("neighbor_filter_policy")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("inbound_join_prune_filter_policy") {
		_, newRelParam := d.GetChange("inbound_join_prune_filter_policy")
		jpInboundFilterPol := models.NewPIMJPInboundFilterPolicy(pimIfPol.DistinguishedName, "", models.PIMJPInboundFilterPolicyAttributes{})
		jpInboundFilterPol.Status = "created,modified"
		err := aciClient.Save(jpInboundFilterPol)
		if err != nil {
			return diag.FromErr(err)
		}

		err = aciClient.CreateRelationPIMJPInbFilterPolrtdmcRsFilterToRtMapPol(jpInboundFilterPol.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("outbound_join_prune_filter_policy") {
		_, newRelParam := d.GetChange("outbound_join_prune_filter_policy")
		jpOutboundFilterPol := models.NewPIMJPOutboundFilterPolicy(pimIfPol.DistinguishedName, "", models.PIMJPOutboundFilterPolicyAttributes{})
		jpOutboundFilterPol.Status = "created,modified"
		err := aciClient.Save(jpOutboundFilterPol)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationPIMJPOutbFilterPolrtdmcRsFilterToRtMapPol(jpOutboundFilterPol.DistinguishedName, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("neighbor_filter_policy") {
		_, newRelParam := d.GetChange("neighbor_filter_policy")
		neighborFilterPol := models.NewPIMNeighborFiterPolicy(pimIfPol.DistinguishedName, "", models.PIMNeighborFiterPolicyAttributes{})
		neighborFilterPol.Status = "created,modified"
		err := aciClient.Save(neighborFilterPol)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationPIMNbrFilterPolrtdmcRsFilterToRtMapPol(neighborFilterPol.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(pimIfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciPIMInterfacePolicyRead(ctx, d, m)
}

func resourceAciPIMInterfacePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	pimIfPol, err := getRemotePIMInterfacePolicy(aciClient, dn)
	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}

	_, err = setPIMInterfacePolicyAttributes(pimIfPol, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	_, err = getandSetPIMIfPolRelationshipAttributes(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] PimIfPol Relationship Attributes - Read finished successfully")
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciPIMInterfacePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, models.PimIfPolClassName)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
