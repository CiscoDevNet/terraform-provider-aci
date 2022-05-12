package aci

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciEPLoopProtectionPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciEPLoopProtectionPolicyCreate,
		UpdateContext: resourceAciEPLoopProtectionPolicyUpdate,
		ReadContext:   resourceAciEPLoopProtectionPolicyRead,
		DeleteContext: resourceAciEPLoopProtectionPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciEPLoopProtectionPolicyImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{

			"action": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"bd-learn-disable",
						"port-disable",
					}, false),
				},
			},
			"admin_st": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"disabled",
					"enabled",
				}, false),
			},
			"loop_detect_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"loop_detect_mult": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func getRemoteEPLoopProtectionPolicy(client *client.Client, dn string) (*models.EPLoopProtectionPolicy, error) {
	epLoopProtectPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	epLoopProtectP := models.EPLoopProtectionPolicyFromContainer(epLoopProtectPCont)
	if epLoopProtectP.DistinguishedName == "" {
		return nil, fmt.Errorf("EPLoopProtectionPolicy %s not found", epLoopProtectP.DistinguishedName)
	}
	return epLoopProtectP, nil
}

func setEPLoopProtectionPolicyAttributes(epLoopProtectP *models.EPLoopProtectionPolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(epLoopProtectP.DistinguishedName)
	d.Set("description", epLoopProtectP.Description)
	epLoopProtectPMap, err := epLoopProtectP.ToMap()
	if err != nil {
		return nil, err
	}
	actionGet := make([]string, 0, 1)
	if epLoopProtectPMap["action"] == "" {
		d.Set("action", actionGet)
	} else {
		for _, val := range strings.Split(epLoopProtectPMap["action"], ",") {
			actionGet = append(actionGet, strings.Trim(val, " "))
		}
		sort.Strings(actionGet)
		if actionIntr, ok := d.GetOk("action"); ok {
			actionAct := make([]string, 0, 1)
			for _, val := range actionIntr.([]interface{}) {
				actionAct = append(actionAct, val.(string))
			}
			sort.Strings(actionAct)
			if reflect.DeepEqual(actionAct, actionGet) {
				d.Set("action", d.Get("action").([]interface{}))
			} else {
				d.Set("action", actionGet)
			}
		} else {
			d.Set("action", actionGet)
		}
	}
	d.Set("admin_st", epLoopProtectPMap["adminSt"])
	d.Set("annotation", epLoopProtectPMap["annotation"])
	d.Set("loop_detect_intvl", epLoopProtectPMap["loopDetectIntvl"])
	d.Set("loop_detect_mult", epLoopProtectPMap["loopDetectMult"])
	d.Set("name_alias", epLoopProtectPMap["nameAlias"])
	return d, nil
}

func resourceAciEPLoopProtectionPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	epLoopProtectP, err := getRemoteEPLoopProtectionPolicy(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setEPLoopProtectionPolicyAttributes(epLoopProtectP, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciEPLoopProtectionPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] EPLoopProtectionPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := "default"
	epLoopProtectPAttr := models.EPLoopProtectionPolicyAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		epLoopProtectPAttr.Annotation = Annotation.(string)
	} else {
		epLoopProtectPAttr.Annotation = "{}"
	}

	if Action, ok := d.GetOk("action"); ok {
		actionList := make([]string, 0, 1)
		for _, val := range Action.([]interface{}) {
			actionList = append(actionList, val.(string))
		}
		Action := strings.Join(actionList, ",")
		epLoopProtectPAttr.Action = Action
	} else {
		epLoopProtectPAttr.Action = "{}"
	}

	if AdminSt, ok := d.GetOk("admin_st"); ok {
		epLoopProtectPAttr.AdminSt = AdminSt.(string)
	}

	if LoopDetectIntvl, ok := d.GetOk("loop_detect_intvl"); ok {
		epLoopProtectPAttr.LoopDetectIntvl = LoopDetectIntvl.(string)
	}

	if LoopDetectMult, ok := d.GetOk("loop_detect_mult"); ok {
		epLoopProtectPAttr.LoopDetectMult = LoopDetectMult.(string)
	}

	epLoopProtectPAttr.Name = "default"

	epLoopProtectP := models.NewEPLoopProtectionPolicy(fmt.Sprintf("infra/epLoopProtectP-%s", name), "uni", desc, nameAlias, epLoopProtectPAttr)
	epLoopProtectP.Status = "modified"
	err := aciClient.Save(epLoopProtectP)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(epLoopProtectP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciEPLoopProtectionPolicyRead(ctx, d, m)
}

func resourceAciEPLoopProtectionPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] EPLoopProtectionPolicy: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := "default"
	epLoopProtectPAttr := models.EPLoopProtectionPolicyAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		epLoopProtectPAttr.Annotation = Annotation.(string)
	} else {
		epLoopProtectPAttr.Annotation = "{}"
	}
	if Action, ok := d.GetOk("action"); ok {
		actionList := make([]string, 0, 1)
		for _, val := range Action.([]interface{}) {
			actionList = append(actionList, val.(string))
		}
		Action := strings.Join(actionList, ",")
		epLoopProtectPAttr.Action = Action
	} else {
		epLoopProtectPAttr.Action = "{}"
	}
	if AdminSt, ok := d.GetOk("admin_st"); ok {
		epLoopProtectPAttr.AdminSt = AdminSt.(string)
	}

	if LoopDetectIntvl, ok := d.GetOk("loop_detect_intvl"); ok {
		epLoopProtectPAttr.LoopDetectIntvl = LoopDetectIntvl.(string)
	}

	if LoopDetectMult, ok := d.GetOk("loop_detect_mult"); ok {
		epLoopProtectPAttr.LoopDetectMult = LoopDetectMult.(string)
	}

	epLoopProtectPAttr.Name = "default"

	epLoopProtectP := models.NewEPLoopProtectionPolicy(fmt.Sprintf("infra/epLoopProtectP-%s", name), "uni", desc, nameAlias, epLoopProtectPAttr)
	epLoopProtectP.Status = "modified"
	err := aciClient.Save(epLoopProtectP)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(epLoopProtectP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciEPLoopProtectionPolicyRead(ctx, d, m)
}

func resourceAciEPLoopProtectionPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	epLoopProtectP, err := getRemoteEPLoopProtectionPolicy(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setEPLoopProtectionPolicyAttributes(epLoopProtectP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciEPLoopProtectionPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	d.SetId("")
	var diags diag.Diagnostics
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Resource with class name epLoopProtectP cannot be deleted",
	})
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	return diags
}
