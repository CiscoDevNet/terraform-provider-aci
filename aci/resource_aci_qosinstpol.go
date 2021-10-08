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

func resourceAciQOSInstancePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciQOSInstancePolicyCreate,
		UpdateContext: resourceAciQOSInstancePolicyUpdate,
		ReadContext:   resourceAciQOSInstancePolicyRead,
		DeleteContext: resourceAciQOSInstancePolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciQOSInstancePolicyImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{

			"etrap_age_timer": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"etrap_bw_thresh": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"etrap_byte_ct": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"etrap_st": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"fabric_flush_interval": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"fabric_flush_st": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"ctrl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"dot1p-preserve",
					"none",
				}, false),
			},
			"uburst_spine_queues": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"uburst_tor_queues": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func getRemoteQOSInstancePolicy(client *client.Client, dn string) (*models.QOSInstancePolicy, error) {
	qosInstPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	qosInstPol := models.QOSInstancePolicyFromContainer(qosInstPolCont)
	if qosInstPol.DistinguishedName == "" {
		return nil, fmt.Errorf("QOSInstancePolicy %s not found", qosInstPol.DistinguishedName)
	}
	return qosInstPol, nil
}

func setQOSInstancePolicyAttributes(qosInstPol *models.QOSInstancePolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(qosInstPol.DistinguishedName)
	d.Set("description", qosInstPol.Description)
	qosInstPolMap, err := qosInstPol.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("etrap_age_timer", qosInstPolMap["EtrapAgeTimer"])
	d.Set("etrap_bw_thresh", qosInstPolMap["EtrapBwThresh"])
	d.Set("etrap_byte_ct", qosInstPolMap["EtrapByteCt"])
	d.Set("etrap_st", qosInstPolMap["EtrapSt"])
	d.Set("fabric_flush_interval", qosInstPolMap["FabricFlushInterval"])
	d.Set("fabric_flush_st", qosInstPolMap["FabricFlushSt"])
	d.Set("annotation", qosInstPolMap["annotation"])
	if qosInstPolMap["ctrl"] == "" {
		d.Set("ctrl", "none")
	} else {
		d.Set("ctrl", qosInstPolMap["ctrl"])
	}
	d.Set("uburst_spine_queues", qosInstPolMap["uburstSpineQueues"])
	d.Set("uburst_tor_queues", qosInstPolMap["uburstTorQueues"])
	d.Set("name_alias", qosInstPolMap["nameAlias"])
	return d, nil
}

func resourceAciQOSInstancePolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	qosInstPol, err := getRemoteQOSInstancePolicy(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setQOSInstancePolicyAttributes(qosInstPol, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciQOSInstancePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] QOSInstancePolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := "default"
	qosInstPolAttr := models.QOSInstancePolicyAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		qosInstPolAttr.Annotation = Annotation.(string)
	} else {
		qosInstPolAttr.Annotation = "{}"
	}

	if EtrapAgeTimer, ok := d.GetOk("etrap_age_timer"); ok {
		qosInstPolAttr.EtrapAgeTimer = EtrapAgeTimer.(string)
	}

	if EtrapBwThresh, ok := d.GetOk("etrap_bw_thresh"); ok {
		qosInstPolAttr.EtrapBwThresh = EtrapBwThresh.(string)
	}

	if EtrapByteCt, ok := d.GetOk("etrap_byte_ct"); ok {
		qosInstPolAttr.EtrapByteCt = EtrapByteCt.(string)
	}

	if EtrapSt, ok := d.GetOk("etrap_st"); ok {
		qosInstPolAttr.EtrapSt = EtrapSt.(string)
	}

	if FabricFlushInterval, ok := d.GetOk("fabric_flush_interval"); ok {
		qosInstPolAttr.FabricFlushInterval = FabricFlushInterval.(string)
	}

	if FabricFlushSt, ok := d.GetOk("fabric_flush_st"); ok {
		qosInstPolAttr.FabricFlushSt = FabricFlushSt.(string)
	}

	if Ctrl, ok := d.GetOk("ctrl"); ok {
		qosInstPolAttr.Ctrl = Ctrl.(string)
	}

	qosInstPolAttr.Name = "default"

	if UburstSpineQueues, ok := d.GetOk("uburst_spine_queues"); ok {
		qosInstPolAttr.UburstSpineQueues = UburstSpineQueues.(string)
	}

	if UburstTorQueues, ok := d.GetOk("uburst_tor_queues"); ok {
		qosInstPolAttr.UburstTorQueues = UburstTorQueues.(string)
	}
	qosInstPol := models.NewQOSInstancePolicy(fmt.Sprintf("infra/qosinst-%s", name), "uni", desc, nameAlias, qosInstPolAttr)
	qosInstPol.Status = "modified"
	err := aciClient.Save(qosInstPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(qosInstPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciQOSInstancePolicyRead(ctx, d, m)
}

func resourceAciQOSInstancePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] QOSInstancePolicy: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := "default"
	qosInstPolAttr := models.QOSInstancePolicyAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		qosInstPolAttr.Annotation = Annotation.(string)
	} else {
		qosInstPolAttr.Annotation = "{}"
	}

	if EtrapAgeTimer, ok := d.GetOk("etrap_age_timer"); ok {
		qosInstPolAttr.EtrapAgeTimer = EtrapAgeTimer.(string)
	}

	if EtrapBwThresh, ok := d.GetOk("etrap_bw_thresh"); ok {
		qosInstPolAttr.EtrapBwThresh = EtrapBwThresh.(string)
	}

	if EtrapByteCt, ok := d.GetOk("etrap_byte_ct"); ok {
		qosInstPolAttr.EtrapByteCt = EtrapByteCt.(string)
	}

	if EtrapSt, ok := d.GetOk("etrap_st"); ok {
		qosInstPolAttr.EtrapSt = EtrapSt.(string)
	}

	if FabricFlushInterval, ok := d.GetOk("fabric_flush_interval"); ok {
		qosInstPolAttr.FabricFlushInterval = FabricFlushInterval.(string)
	}

	if FabricFlushSt, ok := d.GetOk("fabric_flush_st"); ok {
		qosInstPolAttr.FabricFlushSt = FabricFlushSt.(string)
	}

	if Ctrl, ok := d.GetOk("ctrl"); ok {
		qosInstPolAttr.Ctrl = Ctrl.(string)
	}

	qosInstPolAttr.Name = "default"

	if UburstSpineQueues, ok := d.GetOk("uburst_spine_queues"); ok {
		qosInstPolAttr.UburstSpineQueues = UburstSpineQueues.(string)
	}

	if UburstTorQueues, ok := d.GetOk("uburst_tor_queues"); ok {
		qosInstPolAttr.UburstTorQueues = UburstTorQueues.(string)
	}
	qosInstPol := models.NewQOSInstancePolicy(fmt.Sprintf("infra/qosinst-%s", name), "uni", desc, nameAlias, qosInstPolAttr)
	qosInstPol.Status = "modified"
	err := aciClient.Save(qosInstPol)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(qosInstPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciQOSInstancePolicyRead(ctx, d, m)
}

func resourceAciQOSInstancePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	qosInstPol, err := getRemoteQOSInstancePolicy(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	_, err = setQOSInstancePolicyAttributes(qosInstPol, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciQOSInstancePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	d.SetId("")
	var diags diag.Diagnostics
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Resource with class name qosInstPol cannot be deleted",
	})
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	return diags
}
