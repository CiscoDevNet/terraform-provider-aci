package aci

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciSubjectFilter() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciSubjectFilterCreate,
		UpdateContext: resourceAciSubjectFilterUpdate,
		ReadContext:   resourceAciSubjectFilterRead,
		DeleteContext: resourceAciSubjectFilterDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciSubjectFilterImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"contract_subject_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"deny",
					"permit",
				}, false),
			},
			"directives": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"log",
						"no_stats",
						"none",
					}, false),
				},
			},
			"priority_override": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"default",
					"level1",
					"level2",
					"level3",
				}, false),
			},
			"filter_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		}),
	}
}

func getRemoteSubjectFilter(client *client.Client, dn string) (*models.SubjectFilter, error) {
	vzRsSubjFiltAttCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	vzRsSubjFiltAtt := models.SubjectFilterFromContainer(vzRsSubjFiltAttCont)
	if vzRsSubjFiltAtt.DistinguishedName == "" {
		return nil, fmt.Errorf("Subject Filter %s not found", dn)
	}
	return vzRsSubjFiltAtt, nil
}

func setSubjectFilterAttributes(vzRsSubjFiltAtt *models.SubjectFilter, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(vzRsSubjFiltAtt.DistinguishedName)
	if dn != vzRsSubjFiltAtt.DistinguishedName {
		d.Set("contract_subject_dn", "")
	}
	vzRsSubjFiltAttMap, err := vzRsSubjFiltAtt.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("contract_subject_dn", GetParentDn(vzRsSubjFiltAtt.DistinguishedName, fmt.Sprintf("/"+models.RnvzRsSubjFiltAtt, vzRsSubjFiltAttMap["name"])))
	d.Set("annotation", vzRsSubjFiltAttMap["annotation"])
	d.Set("action", vzRsSubjFiltAttMap["action"])
	directivesGet := make([]string, 0, 1)
	for _, val := range strings.Split(vzRsSubjFiltAttMap["directives"], ",") {
		directivesGet = append(directivesGet, strings.Trim(val, " "))
	}
	d.Set("directives", directivesGet)
	d.Set("priority_override", vzRsSubjFiltAttMap["priorityOverride"])
	d.Set("filter_dn", vzRsSubjFiltAttMap["tDn"])
	return d, nil
}

func resourceAciSubjectFilterImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	vzRsSubjFiltAtt, err := getRemoteSubjectFilter(aciClient, dn)
	if err != nil {
		return nil, err
	}

	schemaFilled, err := setSubjectFilterAttributes(vzRsSubjFiltAtt, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciSubjectFilterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SubjectFilter: Beginning Creation")
	aciClient := m.(*client.Client)
	tnVzFilterName := GetMOName(d.Get("filter_dn").(string))
	ContractSubjectDn := d.Get("contract_subject_dn").(string)

	vzRsSubjFiltAttAttr := models.SubjectFilterAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		vzRsSubjFiltAttAttr.Annotation = Annotation.(string)
	} else {
		vzRsSubjFiltAttAttr.Annotation = "{}"
	}

	if Action, ok := d.GetOk("action"); ok {
		vzRsSubjFiltAttAttr.Action = Action.(string)
	}

	if Directives, ok := d.GetOk("directives"); ok {
		directivesList := make([]string, 0, 1)
		for _, val := range Directives.([]interface{}) {
			directivesList = append(directivesList, val.(string))
		}
		Directives := strings.Join(directivesList, ",")
		vzRsSubjFiltAttAttr.Directives = Directives
	}

	if PriorityOverride, ok := d.GetOk("priority_override"); ok {
		vzRsSubjFiltAttAttr.PriorityOverride = PriorityOverride.(string)
	}

	vzRsSubjFiltAttAttr.TnVzFilterName = tnVzFilterName

	vzRsSubjFiltAtt := models.NewSubjectFilter(fmt.Sprintf(models.RnvzRsSubjFiltAtt, tnVzFilterName), ContractSubjectDn, vzRsSubjFiltAttAttr)

	err := aciClient.Save(vzRsSubjFiltAtt)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(vzRsSubjFiltAtt.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciSubjectFilterRead(ctx, d, m)
}

func resourceAciSubjectFilterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SubjectFilter: Beginning Update")
	aciClient := m.(*client.Client)
	tnVzFilterName := GetMOName(d.Get("filter_dn").(string))
	ContractSubjectDn := d.Get("contract_subject_dn").(string)

	vzRsSubjFiltAttAttr := models.SubjectFilterAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		vzRsSubjFiltAttAttr.Annotation = Annotation.(string)
	} else {
		vzRsSubjFiltAttAttr.Annotation = "{}"
	}

	if Action, ok := d.GetOk("action"); ok {
		vzRsSubjFiltAttAttr.Action = Action.(string)
	}
	if Directives, ok := d.GetOk("directives"); ok {
		directivesList := make([]string, 0, 1)
		for _, val := range Directives.([]interface{}) {
			directivesList = append(directivesList, val.(string))
		}
		Directives := strings.Join(directivesList, ",")
		vzRsSubjFiltAttAttr.Directives = Directives
	}

	if PriorityOverride, ok := d.GetOk("priority_override"); ok {
		vzRsSubjFiltAttAttr.PriorityOverride = PriorityOverride.(string)
	}

	vzRsSubjFiltAttAttr.TnVzFilterName = tnVzFilterName

	vzRsSubjFiltAtt := models.NewSubjectFilter(fmt.Sprintf(models.RnvzRsSubjFiltAtt, tnVzFilterName), ContractSubjectDn, vzRsSubjFiltAttAttr)

	vzRsSubjFiltAtt.Status = "modified"

	err := aciClient.Save(vzRsSubjFiltAtt)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(vzRsSubjFiltAtt.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciSubjectFilterRead(ctx, d, m)
}

func resourceAciSubjectFilterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	vzRsSubjFiltAtt, err := getRemoteSubjectFilter(aciClient, dn)
	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}

	_, err = setSubjectFilterAttributes(vzRsSubjFiltAtt, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciSubjectFilterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "vzRsSubjFiltAtt")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
