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

func resourceAciServiceRedirectPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciServiceRedirectPolicyCreate,
		UpdateContext: resourceAciServiceRedirectPolicyUpdate,
		ReadContext:   resourceAciServiceRedirectPolicyRead,
		DeleteContext: resourceAciServiceRedirectPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciServiceRedirectPolicyImport,
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

			"anycast_enabled": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},

			"dest_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"L3",
					"L2",
					"L1",
				}, false),
			},

			"hashing_algorithm": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"sip-dip-prototype",
					"sip",
					"dip",
				}, false),
			},

			"max_threshold_percent": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"min_threshold_percent": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"program_local_pod_only": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},

			"resilient_hash_enabled": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},

			"threshold_down_action": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"permit",
					"deny",
					"bypass",
				}, false),
			},

			"threshold_enable": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},

			"relation_vns_rs_ipsla_monitoring_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
		}),
	}
}
func getRemoteServiceRedirectPolicy(client *client.Client, dn string) (*models.ServiceRedirectPolicy, error) {
	vnsSvcRedirectPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	vnsSvcRedirectPol := models.ServiceRedirectPolicyFromContainer(vnsSvcRedirectPolCont)

	if vnsSvcRedirectPol.DistinguishedName == "" {
		return nil, fmt.Errorf("ServiceRedirectPolicy %s not found", vnsSvcRedirectPol.DistinguishedName)
	}

	return vnsSvcRedirectPol, nil
}

func setServiceRedirectPolicyAttributes(vnsSvcRedirectPol *models.ServiceRedirectPolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(vnsSvcRedirectPol.DistinguishedName)
	d.Set("description", vnsSvcRedirectPol.Description)
	if dn != vnsSvcRedirectPol.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	vnsSvcRedirectPolMap, err := vnsSvcRedirectPol.ToMap()

	if err != nil {
		return d, err
	}

	d.Set("name", vnsSvcRedirectPolMap["name"])

	d.Set("anycast_enabled", vnsSvcRedirectPolMap["AnycastEnabled"])
	d.Set("annotation", vnsSvcRedirectPolMap["annotation"])
	d.Set("dest_type", vnsSvcRedirectPolMap["destType"])
	d.Set("hashing_algorithm", vnsSvcRedirectPolMap["hashingAlgorithm"])
	d.Set("max_threshold_percent", vnsSvcRedirectPolMap["maxThresholdPercent"])
	d.Set("min_threshold_percent", vnsSvcRedirectPolMap["minThresholdPercent"])
	d.Set("name_alias", vnsSvcRedirectPolMap["nameAlias"])
	d.Set("program_local_pod_only", vnsSvcRedirectPolMap["programLocalPodOnly"])
	d.Set("resilient_hash_enabled", vnsSvcRedirectPolMap["resilientHashEnabled"])
	d.Set("threshold_down_action", vnsSvcRedirectPolMap["thresholdDownAction"])
	d.Set("threshold_enable", vnsSvcRedirectPolMap["thresholdEnable"])
	return d, nil
}

func resourceAciServiceRedirectPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	vnsSvcRedirectPol, err := getRemoteServiceRedirectPolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	vnsSvcRedirectPolMap, err := vnsSvcRedirectPol.ToMap()

	if err != nil {
		return nil, err
	}
	name := vnsSvcRedirectPolMap["name"]
	pDN := GetParentDn(dn, fmt.Sprintf("/svcCont/svcRedirectPol-%s", name))
	d.Set("tenant_dn", pDN)
	schemaFilled, err := setServiceRedirectPolicyAttributes(vnsSvcRedirectPol, d)

	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciServiceRedirectPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ServiceRedirectPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	vnsSvcRedirectPolAttr := models.ServiceRedirectPolicyAttributes{}
	if AnycastEnabled, ok := d.GetOk("anycast_enabled"); ok {
		vnsSvcRedirectPolAttr.AnycastEnabled = AnycastEnabled.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vnsSvcRedirectPolAttr.Annotation = Annotation.(string)
	} else {
		vnsSvcRedirectPolAttr.Annotation = "{}"
	}
	if DestType, ok := d.GetOk("dest_type"); ok {
		vnsSvcRedirectPolAttr.DestType = DestType.(string)
	}
	if HashingAlgorithm, ok := d.GetOk("hashing_algorithm"); ok {
		vnsSvcRedirectPolAttr.HashingAlgorithm = HashingAlgorithm.(string)
	}
	if MaxThresholdPercent, ok := d.GetOk("max_threshold_percent"); ok {
		vnsSvcRedirectPolAttr.MaxThresholdPercent = MaxThresholdPercent.(string)
	}
	if MinThresholdPercent, ok := d.GetOk("min_threshold_percent"); ok {
		vnsSvcRedirectPolAttr.MinThresholdPercent = MinThresholdPercent.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vnsSvcRedirectPolAttr.NameAlias = NameAlias.(string)
	}
	if ProgramLocalPodOnly, ok := d.GetOk("program_local_pod_only"); ok {
		vnsSvcRedirectPolAttr.ProgramLocalPodOnly = ProgramLocalPodOnly.(string)
	}
	if ResilientHashEnabled, ok := d.GetOk("resilient_hash_enabled"); ok {
		vnsSvcRedirectPolAttr.ResilientHashEnabled = ResilientHashEnabled.(string)
	}
	if ThresholdDownAction, ok := d.GetOk("threshold_down_action"); ok {
		vnsSvcRedirectPolAttr.ThresholdDownAction = ThresholdDownAction.(string)
	}
	if ThresholdEnable, ok := d.GetOk("threshold_enable"); ok {
		vnsSvcRedirectPolAttr.ThresholdEnable = ThresholdEnable.(string)
	}
	vnsSvcRedirectPol := models.NewServiceRedirectPolicy(fmt.Sprintf("svcCont/svcRedirectPol-%s", name), TenantDn, desc, vnsSvcRedirectPolAttr)

	err := aciClient.Save(vnsSvcRedirectPol)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationTovnsRsIPSLAMonitoringPol, ok := d.GetOk("relation_vns_rs_ipsla_monitoring_pol"); ok {
		relationParam := relationTovnsRsIPSLAMonitoringPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTovnsRsIPSLAMonitoringPol, ok := d.GetOk("relation_vns_rs_ipsla_monitoring_pol"); ok {
		relationParam := relationTovnsRsIPSLAMonitoringPol.(string)
		err = aciClient.CreateRelationvnsRsIPSLAMonitoringPolFromServiceRedirectPolicy(vnsSvcRedirectPol.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(vnsSvcRedirectPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciServiceRedirectPolicyRead(ctx, d, m)
}

func resourceAciServiceRedirectPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ServiceRedirectPolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	vnsSvcRedirectPolAttr := models.ServiceRedirectPolicyAttributes{}
	if AnycastEnabled, ok := d.GetOk("anycast_enabled"); ok {
		vnsSvcRedirectPolAttr.AnycastEnabled = AnycastEnabled.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vnsSvcRedirectPolAttr.Annotation = Annotation.(string)
	} else {
		vnsSvcRedirectPolAttr.Annotation = "{}"
	}
	if DestType, ok := d.GetOk("dest_type"); ok {
		vnsSvcRedirectPolAttr.DestType = DestType.(string)
	}
	if HashingAlgorithm, ok := d.GetOk("hashing_algorithm"); ok {
		vnsSvcRedirectPolAttr.HashingAlgorithm = HashingAlgorithm.(string)
	}
	if MaxThresholdPercent, ok := d.GetOk("max_threshold_percent"); ok {
		vnsSvcRedirectPolAttr.MaxThresholdPercent = MaxThresholdPercent.(string)
	}
	if MinThresholdPercent, ok := d.GetOk("min_threshold_percent"); ok {
		vnsSvcRedirectPolAttr.MinThresholdPercent = MinThresholdPercent.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vnsSvcRedirectPolAttr.NameAlias = NameAlias.(string)
	}
	if ProgramLocalPodOnly, ok := d.GetOk("program_local_pod_only"); ok {
		vnsSvcRedirectPolAttr.ProgramLocalPodOnly = ProgramLocalPodOnly.(string)
	}
	if ResilientHashEnabled, ok := d.GetOk("resilient_hash_enabled"); ok {
		vnsSvcRedirectPolAttr.ResilientHashEnabled = ResilientHashEnabled.(string)
	}
	if ThresholdDownAction, ok := d.GetOk("threshold_down_action"); ok {
		vnsSvcRedirectPolAttr.ThresholdDownAction = ThresholdDownAction.(string)
	}
	if ThresholdEnable, ok := d.GetOk("threshold_enable"); ok {
		vnsSvcRedirectPolAttr.ThresholdEnable = ThresholdEnable.(string)
	}
	vnsSvcRedirectPol := models.NewServiceRedirectPolicy(fmt.Sprintf("svcCont/svcRedirectPol-%s", name), TenantDn, desc, vnsSvcRedirectPolAttr)

	vnsSvcRedirectPol.Status = "modified"

	err := aciClient.Save(vnsSvcRedirectPol)

	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_vns_rs_ipsla_monitoring_pol") {
		_, newRelParam := d.GetChange("relation_vns_rs_ipsla_monitoring_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_vns_rs_ipsla_monitoring_pol") {
		_, newRelParam := d.GetChange("relation_vns_rs_ipsla_monitoring_pol")
		err = aciClient.DeleteRelationvnsRsIPSLAMonitoringPolFromServiceRedirectPolicy(vnsSvcRedirectPol.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvnsRsIPSLAMonitoringPolFromServiceRedirectPolicy(vnsSvcRedirectPol.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(vnsSvcRedirectPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciServiceRedirectPolicyRead(ctx, d, m)

}

func resourceAciServiceRedirectPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	vnsSvcRedirectPol, err := getRemoteServiceRedirectPolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setServiceRedirectPolicyAttributes(vnsSvcRedirectPol, d)

	if err != nil {
		d.SetId("")
		return nil
	}

	vnsRsIPSLAMonitoringPolData, err := aciClient.ReadRelationvnsRsIPSLAMonitoringPolFromServiceRedirectPolicy(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsIPSLAMonitoringPol %v", err)
		d.Set("relation_vns_rs_ipsla_monitoring_pol", "")

	} else {
		d.Set("relation_vns_rs_ipsla_monitoring_pol", vnsRsIPSLAMonitoringPolData.(string))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciServiceRedirectPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "vnsSvcRedirectPol")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
