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

func resourceAciIPSLAMonitoringPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciIPSLAMonitoringPolicyCreate,
		UpdateContext: resourceAciIPSLAMonitoringPolicyUpdate,
		ReadContext:   resourceAciIPSLAMonitoringPolicyRead,
		DeleteContext: resourceAciIPSLAMonitoringPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciIPSLAMonitoringPolicyImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"tenant_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"http_uri": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"http_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"HTTP/1.0",
					"HTTP/1.1",
				}, false),
			},
			"type_of_service": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"traffic_class_value": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"request_data_size": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sla_detect_multiplier": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sla_frequency": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sla_port": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sla_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"http",
					"icmp",
					"l2ping",
					"tcp",
				}, false),
			},
			"threshold": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"timeout": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func getRemoteIPSLAMonitoringPolicy(client *client.Client, dn string) (*models.IPSLAMonitoringPolicy, error) {
	fvIPSLAMonitoringPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	fvIPSLAMonitoringPol := models.IPSLAMonitoringPolicyFromContainer(fvIPSLAMonitoringPolCont)
	if fvIPSLAMonitoringPol.DistinguishedName == "" {
		return nil, fmt.Errorf("IP SLA Monitoring Policy %s not found", dn)
	}
	return fvIPSLAMonitoringPol, nil
}

func setIPSLAMonitoringPolicyAttributes(fvIPSLAMonitoringPol *models.IPSLAMonitoringPolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(fvIPSLAMonitoringPol.DistinguishedName)
	d.Set("description", fvIPSLAMonitoringPol.Description)
	fvIPSLAMonitoringPolMap, err := fvIPSLAMonitoringPol.ToMap()
	if err != nil {
		return d, err
	}
	dn := d.Id()
	if dn != fvIPSLAMonitoringPol.DistinguishedName {
		d.Set("tenant_dn", "")
	} else {
		d.Set("tenant_dn", GetParentDn(fvIPSLAMonitoringPol.DistinguishedName, fmt.Sprintf("/"+models.RnfvIPSLAMonitoringPol, fvIPSLAMonitoringPolMap["name"])))
	}
	d.Set("annotation", fvIPSLAMonitoringPolMap["annotation"])
	d.Set("http_uri", fvIPSLAMonitoringPolMap["httpUri"])
	if fvIPSLAMonitoringPolMap["httpVersion"] == "HTTP10" {
		d.Set("http_version", "HTTP/1.0")
	} else if fvIPSLAMonitoringPolMap["httpVersion"] == "HTTP11" {
		d.Set("http_version", "HTTP/1.1")
	}
	d.Set("type_of_service", fvIPSLAMonitoringPolMap["ipv4Tos"])
	d.Set("traffic_class_value", fvIPSLAMonitoringPolMap["ipv6TrfClass"])
	d.Set("name", fvIPSLAMonitoringPolMap["name"])
	d.Set("request_data_size", fvIPSLAMonitoringPolMap["reqDataSize"])
	d.Set("sla_detect_multiplier", fvIPSLAMonitoringPolMap["slaDetectMultiplier"])
	d.Set("sla_frequency", fvIPSLAMonitoringPolMap["slaFrequency"])
	d.Set("sla_port", fvIPSLAMonitoringPolMap["slaPort"])
	d.Set("sla_type", fvIPSLAMonitoringPolMap["slaType"])
	d.Set("threshold", fvIPSLAMonitoringPolMap["threshold"])
	d.Set("timeout", fvIPSLAMonitoringPolMap["timeout"])
	d.Set("name_alias", fvIPSLAMonitoringPolMap["nameAlias"])
	return d, nil
}

func resourceAciIPSLAMonitoringPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	fvIPSLAMonitoringPol, err := getRemoteIPSLAMonitoringPolicy(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setIPSLAMonitoringPolicyAttributes(fvIPSLAMonitoringPol, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciIPSLAMonitoringPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] IPSLAMonitoringPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	TenantDn := d.Get("tenant_dn").(string)

	fvIPSLAMonitoringPolAttr := models.IPSLAMonitoringPolicyAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		fvIPSLAMonitoringPolAttr.Annotation = Annotation.(string)
	} else {
		fvIPSLAMonitoringPolAttr.Annotation = "{}"
	}

	if HttpUri, ok := d.GetOk("http_uri"); ok {
		fvIPSLAMonitoringPolAttr.HttpUri = HttpUri.(string)
	}

	if HttpVersion, ok := d.GetOk("http_version"); ok {
		if HttpVersion.(string) == "HTTP/1.1" {
			fvIPSLAMonitoringPolAttr.HttpVersion = "HTTP11"
		} else {
			fvIPSLAMonitoringPolAttr.HttpVersion = "HTTP10"
		}
	}

	if Ipv4Tos, ok := d.GetOk("type_of_service"); ok {
		fvIPSLAMonitoringPolAttr.Ipv4Tos = Ipv4Tos.(string)
	}

	if Ipv6TrfClass, ok := d.GetOk("traffic_class_value"); ok {
		fvIPSLAMonitoringPolAttr.Ipv6TrfClass = Ipv6TrfClass.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		fvIPSLAMonitoringPolAttr.Name = Name.(string)
	}

	if ReqDataSize, ok := d.GetOk("request_data_size"); ok {
		fvIPSLAMonitoringPolAttr.ReqDataSize = ReqDataSize.(string)
	}

	if SlaDetectMultiplier, ok := d.GetOk("sla_detect_multiplier"); ok {
		fvIPSLAMonitoringPolAttr.SlaDetectMultiplier = SlaDetectMultiplier.(string)
	}

	if SlaFrequency, ok := d.GetOk("sla_frequency"); ok {
		fvIPSLAMonitoringPolAttr.SlaFrequency = SlaFrequency.(string)
	}

	if SlaPort, ok := d.GetOk("sla_port"); ok {
		fvIPSLAMonitoringPolAttr.SlaPort = SlaPort.(string)
	}

	if SlaType, ok := d.GetOk("sla_type"); ok {
		fvIPSLAMonitoringPolAttr.SlaType = SlaType.(string)
	}

	if Threshold, ok := d.GetOk("threshold"); ok {
		fvIPSLAMonitoringPolAttr.Threshold = Threshold.(string)
	}

	if Timeout, ok := d.GetOk("timeout"); ok {
		fvIPSLAMonitoringPolAttr.Timeout = Timeout.(string)
	}
	fvIPSLAMonitoringPol := models.NewIPSLAMonitoringPolicy(fmt.Sprintf(models.RnfvIPSLAMonitoringPol, name), TenantDn, desc, nameAlias, fvIPSLAMonitoringPolAttr)

	err := aciClient.Save(fvIPSLAMonitoringPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fvIPSLAMonitoringPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciIPSLAMonitoringPolicyRead(ctx, d, m)
}
func resourceAciIPSLAMonitoringPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] IP SLA Monitoring Policy: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	TenantDn := d.Get("tenant_dn").(string)

	fvIPSLAMonitoringPolAttr := models.IPSLAMonitoringPolicyAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		fvIPSLAMonitoringPolAttr.Annotation = Annotation.(string)
	} else {
		fvIPSLAMonitoringPolAttr.Annotation = "{}"
	}

	if HttpUri, ok := d.GetOk("http_uri"); ok {
		fvIPSLAMonitoringPolAttr.HttpUri = HttpUri.(string)
	}

	if HttpVersion, ok := d.GetOk("http_version"); ok {
		if HttpVersion.(string) == "HTTP/1.1" {
			fvIPSLAMonitoringPolAttr.HttpVersion = "HTTP11"
		} else {
			fvIPSLAMonitoringPolAttr.HttpVersion = "HTTP10"
		}
	}

	if Ipv4Tos, ok := d.GetOk("type_of_service"); ok {
		fvIPSLAMonitoringPolAttr.Ipv4Tos = Ipv4Tos.(string)
	}

	if Ipv6TrfClass, ok := d.GetOk("traffic_class_value"); ok {
		fvIPSLAMonitoringPolAttr.Ipv6TrfClass = Ipv6TrfClass.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		fvIPSLAMonitoringPolAttr.Name = Name.(string)
	}

	if ReqDataSize, ok := d.GetOk("request_data_size"); ok {
		fvIPSLAMonitoringPolAttr.ReqDataSize = ReqDataSize.(string)
	}

	if SlaDetectMultiplier, ok := d.GetOk("sla_detect_multiplier"); ok {
		fvIPSLAMonitoringPolAttr.SlaDetectMultiplier = SlaDetectMultiplier.(string)
	}

	if SlaFrequency, ok := d.GetOk("sla_frequency"); ok {
		fvIPSLAMonitoringPolAttr.SlaFrequency = SlaFrequency.(string)
	}

	if SlaPort, ok := d.GetOk("sla_port"); ok {
		fvIPSLAMonitoringPolAttr.SlaPort = SlaPort.(string)
	}

	if SlaType, ok := d.GetOk("sla_type"); ok {
		fvIPSLAMonitoringPolAttr.SlaType = SlaType.(string)
	}

	if Threshold, ok := d.GetOk("threshold"); ok {
		fvIPSLAMonitoringPolAttr.Threshold = Threshold.(string)
	}

	if Timeout, ok := d.GetOk("timeout"); ok {
		fvIPSLAMonitoringPolAttr.Timeout = Timeout.(string)
	}
	fvIPSLAMonitoringPol := models.NewIPSLAMonitoringPolicy(fmt.Sprintf(models.RnfvIPSLAMonitoringPol, name), TenantDn, desc, nameAlias, fvIPSLAMonitoringPolAttr)

	fvIPSLAMonitoringPol.Status = "modified"

	err := aciClient.Save(fvIPSLAMonitoringPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fvIPSLAMonitoringPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciIPSLAMonitoringPolicyRead(ctx, d, m)
}

func resourceAciIPSLAMonitoringPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	fvIPSLAMonitoringPol, err := getRemoteIPSLAMonitoringPolicy(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}

	_, err = setIPSLAMonitoringPolicyAttributes(fvIPSLAMonitoringPol, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciIPSLAMonitoringPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "fvIPSLAMonitoringPol")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
