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

func resourceAciFilterRelationship() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciFilterRelationshipCreate,
		UpdateContext: resourceAciFilterRelationshipUpdate,
		ReadContext:   resourceAciFilterRelationshipRead,
		DeleteContext: resourceAciFilterRelationshipDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciFilterRelationshipImport,
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

func getRemoteFilterRelationship(client *client.Client, dn string) (*models.FilterRelationship, error) {
	vzRsFiltAttCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	vzRsFiltAtt := models.FilterRelationshipFromContainer(vzRsFiltAttCont)
	if vzRsFiltAtt.DistinguishedName == "" {
		return nil, fmt.Errorf("Filter %s not found", vzRsFiltAtt.DistinguishedName)
	}
	return vzRsFiltAtt, nil
}

func setFilterRelationshipAttributes(vzRsFiltAtt *models.FilterRelationship, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(vzRsFiltAtt.DistinguishedName)
	d.Set("description", vzRsFiltAtt.Description)
	if dn != vzRsFiltAtt.DistinguishedName {
		d.Set("contract_subject_dn", "")
	}

	vzRsFiltAttMap, err := vzRsFiltAtt.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("contract_subject_dn", GetParentDn(vzRsFiltAtt.DistinguishedName, fmt.Sprintf("/"+models.RnvzRsFiltAtt, vzRsFiltAttMap["name"])))
	d.Set("annotation", vzRsFiltAttMap["annotation"])
	d.Set("action", vzRsFiltAttMap["action"])
	directivesGet := make([]string, 0, 1)
	for _, val := range strings.Split(vzRsFiltAttMap["directives"], ",") {
		directivesGet = append(directivesGet, strings.Trim(val, " "))
	}
	d.Set("directives", directivesGet)
	d.Set("priority_override", vzRsFiltAttMap["priorityOverride"])
	d.Set("filter_dn", vzRsFiltAttMap["tDn"])
	return d, nil
}

func resourceAciFilterRelationshipImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	vzRsFiltAtt, err := getRemoteFilterRelationship(aciClient, dn)
	if err != nil {
		return nil, err
	}

	schemaFilled, err := setFilterRelationshipAttributes(vzRsFiltAtt, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciFilterRelationshipCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Filter: Beginning Creation")
	aciClient := m.(*client.Client)
	tnVzFilterName := GetMOName(d.Get("filter_dn").(string))
	ContractSubjectDn := d.Get("contract_subject_dn").(string)

	vzRsFiltAttAttr := models.FilterRelationshipAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		vzRsFiltAttAttr.Annotation = Annotation.(string)
	} else {
		vzRsFiltAttAttr.Annotation = "{}"
	}

	if Action, ok := d.GetOk("action"); ok {
		vzRsFiltAttAttr.Action = Action.(string)
	}

	if Directives, ok := d.GetOk("directives"); ok {
		directivesList := make([]string, 0, 1)
		for _, val := range Directives.([]interface{}) {
			directivesList = append(directivesList, val.(string))
		}
		Directives := strings.Join(directivesList, ",")
		vzRsFiltAttAttr.Directives = Directives
	}

	if PriorityOverride, ok := d.GetOk("priority_override"); ok {
		vzRsFiltAttAttr.PriorityOverride = PriorityOverride.(string)
	}

	vzRsFiltAttAttr.TnVzFilterName = tnVzFilterName

	vzRsFiltAtt := models.NewFilterRelationship(fmt.Sprintf(models.RnvzRsFiltAtt, tnVzFilterName), ContractSubjectDn, vzRsFiltAttAttr)

	err := aciClient.Save(vzRsFiltAtt)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(vzRsFiltAtt.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciFilterRelationshipRead(ctx, d, m)
}

func resourceAciFilterRelationshipUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Filter: Beginning Update")
	aciClient := m.(*client.Client)
	tnVzFilterName := GetMOName(d.Get("filter_dn").(string))
	ContractSubjectDn := d.Get("contract_subject_dn").(string)

	vzRsFiltAttAttr := models.FilterRelationshipAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		vzRsFiltAttAttr.Annotation = Annotation.(string)
	} else {
		vzRsFiltAttAttr.Annotation = "{}"
	}

	if Action, ok := d.GetOk("action"); ok {
		vzRsFiltAttAttr.Action = Action.(string)
	}
	if Directives, ok := d.GetOk("directives"); ok {
		directivesList := make([]string, 0, 1)
		for _, val := range Directives.([]interface{}) {
			directivesList = append(directivesList, val.(string))
		}
		Directives := strings.Join(directivesList, ",")
		vzRsFiltAttAttr.Directives = Directives
	}

	if PriorityOverride, ok := d.GetOk("priority_override"); ok {
		vzRsFiltAttAttr.PriorityOverride = PriorityOverride.(string)
	}

	vzRsFiltAttAttr.TnVzFilterName = tnVzFilterName

	vzRsFiltAtt := models.NewFilterRelationship(fmt.Sprintf(models.RnvzRsFiltAtt, tnVzFilterName), ContractSubjectDn, vzRsFiltAttAttr)

	vzRsFiltAtt.Status = "modified"

	err := aciClient.Save(vzRsFiltAtt)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(vzRsFiltAtt.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciFilterRelationshipRead(ctx, d, m)
}

func resourceAciFilterRelationshipRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	vzRsFiltAtt, err := getRemoteFilterRelationship(aciClient, dn)
	if err != nil {
		d.SetId("")
		// return diag.FromErr(err)
		return nil
	}

	_, err = setFilterRelationshipAttributes(vzRsFiltAtt, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciFilterRelationshipDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "vzRsFiltAtt")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
