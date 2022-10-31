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

func resourceAciRecurringWindow() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciRecurringWindowCreate,
		UpdateContext: resourceAciRecurringWindowUpdate,
		ReadContext:   resourceAciRecurringWindowRead,
		DeleteContext: resourceAciRecurringWindowDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciRecurringWindowImport,
		},

		SchemaVersion: 1,
		Schema: AppendAttrSchemas(
			GetAnnotationAttrSchema(),
			GetNameAliasAttrSchema(),
			map[string]*schema.Schema{
				"scheduler_dn": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},
				"concur_cap": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},

				"day": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						"Friday",
						"Monday",
						"Saturday",
						"Sunday",
						"Thursday",
						"Tuesday",
						"Wednesday",
						"even-day",
						"every-day",
						"odd-day",
					}, false),
				},
				"hour": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"minute": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"name": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},
				"node_upg_interval": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"proc_break": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.Any(validation.StringInSlice([]string{
						"none",
					}, false), validateColonSeparatedTimeStamp()),
				},
				"proc_cap": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"time_cap": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.Any(validation.StringInSlice([]string{
						"unlimited",
					}, false), validateColonSeparatedTimeStamp()),
				},
			}),
	}
}

func getRemoteRecurringWindow(client *client.Client, dn string) (*models.RecurringWindow, error) {
	trigRecurrWindowPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	trigRecurrWindowP := models.RecurringWindowFromContainer(trigRecurrWindowPCont)
	if trigRecurrWindowP.DistinguishedName == "" {
		return nil, fmt.Errorf("RecurringWindow %s not found", trigRecurrWindowP.DistinguishedName)
	}
	return trigRecurrWindowP, nil
}

func setRecurringWindowAttributes(trigRecurrWindowP *models.RecurringWindow, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(trigRecurrWindowP.DistinguishedName)
	trigRecurrWindowPMap, err := trigRecurrWindowP.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("scheduler_dn", GetParentDn(trigRecurrWindowP.DistinguishedName, fmt.Sprintf("/recurrwinp-%s", trigRecurrWindowPMap["name"])))
	d.Set("concur_cap", trigRecurrWindowPMap["concurCap"])
	d.Set("annotation", trigRecurrWindowPMap["annotation"])
	d.Set("day", trigRecurrWindowPMap["day"])
	d.Set("hour", trigRecurrWindowPMap["hour"])
	d.Set("minute", trigRecurrWindowPMap["minute"])
	d.Set("name", trigRecurrWindowPMap["name"])
	d.Set("node_upg_interval", trigRecurrWindowPMap["nodeUpgInterval"])
	d.Set("proc_break", trigRecurrWindowPMap["procBreak"])
	d.Set("proc_cap", trigRecurrWindowPMap["procCap"])
	d.Set("time_cap", trigRecurrWindowPMap["timeCap"])
	d.Set("name_alias", trigRecurrWindowPMap["nameAlias"])
	return d, nil
}

func resourceAciRecurringWindowImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	trigRecurrWindowP, err := getRemoteRecurringWindow(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setRecurringWindowAttributes(trigRecurrWindowP, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciRecurringWindowCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] RecurringWindow: Beginning Creation")
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	SchedulerDn := d.Get("scheduler_dn").(string)

	trigRecurrWindowPAttr := models.RecurringWindowAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		trigRecurrWindowPAttr.Annotation = Annotation.(string)
	} else {
		trigRecurrWindowPAttr.Annotation = "{}"
	}

	if ConcurCap, ok := d.GetOk("concur_cap"); ok {
		trigRecurrWindowPAttr.ConcurCap = ConcurCap.(string)
	}

	if Day, ok := d.GetOk("day"); ok {
		trigRecurrWindowPAttr.Day = Day.(string)
	}

	if Hour, ok := d.GetOk("hour"); ok {
		trigRecurrWindowPAttr.Hour = Hour.(string)
	}

	if Minute, ok := d.GetOk("minute"); ok {
		trigRecurrWindowPAttr.Minute = Minute.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		trigRecurrWindowPAttr.Name = Name.(string)
	}

	if NodeUpgInterval, ok := d.GetOk("node_upg_interval"); ok {
		trigRecurrWindowPAttr.NodeUpgInterval = NodeUpgInterval.(string)
	}

	if ProcBreak, ok := d.GetOk("proc_break"); ok {
		trigRecurrWindowPAttr.ProcBreak = ProcBreak.(string)
	}

	if ProcCap, ok := d.GetOk("proc_cap"); ok {
		trigRecurrWindowPAttr.ProcCap = ProcCap.(string)
	}

	if TimeCap, ok := d.GetOk("time_cap"); ok {
		trigRecurrWindowPAttr.TimeCap = TimeCap.(string)
	}
	trigRecurrWindowP := models.NewRecurringWindow(fmt.Sprintf("recurrwinp-%s", name), SchedulerDn, nameAlias, trigRecurrWindowPAttr)

	err := aciClient.Save(trigRecurrWindowP)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(trigRecurrWindowP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciRecurringWindowRead(ctx, d, m)
}

func resourceAciRecurringWindowUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] RecurringWindow: Beginning Update")
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	SchedulerDn := d.Get("scheduler_dn").(string)
	trigRecurrWindowPAttr := models.RecurringWindowAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		trigRecurrWindowPAttr.Annotation = Annotation.(string)
	} else {
		trigRecurrWindowPAttr.Annotation = "{}"
	}

	if ConcurCap, ok := d.GetOk("concur_cap"); ok {
		trigRecurrWindowPAttr.ConcurCap = ConcurCap.(string)
	}

	if Day, ok := d.GetOk("day"); ok {
		trigRecurrWindowPAttr.Day = Day.(string)
	}

	if Hour, ok := d.GetOk("hour"); ok {
		trigRecurrWindowPAttr.Hour = Hour.(string)
	}

	if Minute, ok := d.GetOk("minute"); ok {
		trigRecurrWindowPAttr.Minute = Minute.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		trigRecurrWindowPAttr.Name = Name.(string)
	}

	if NodeUpgInterval, ok := d.GetOk("node_upg_interval"); ok {
		trigRecurrWindowPAttr.NodeUpgInterval = NodeUpgInterval.(string)
	}

	if ProcBreak, ok := d.GetOk("proc_break"); ok {
		trigRecurrWindowPAttr.ProcBreak = ProcBreak.(string)
	}

	if ProcCap, ok := d.GetOk("proc_cap"); ok {
		trigRecurrWindowPAttr.ProcCap = ProcCap.(string)
	}

	if TimeCap, ok := d.GetOk("time_cap"); ok {
		trigRecurrWindowPAttr.TimeCap = TimeCap.(string)
	}
	trigRecurrWindowP := models.NewRecurringWindow(fmt.Sprintf("recurrwinp-%s", name), SchedulerDn, nameAlias, trigRecurrWindowPAttr)

	trigRecurrWindowP.Status = "modified"
	err := aciClient.Save(trigRecurrWindowP)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(trigRecurrWindowP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciRecurringWindowRead(ctx, d, m)
}

func resourceAciRecurringWindowRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	trigRecurrWindowP, err := getRemoteRecurringWindow(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}
	setRecurringWindowAttributes(trigRecurrWindowP, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciRecurringWindowDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "trigRecurrWindowP")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
