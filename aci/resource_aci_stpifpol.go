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

func resourceAciSpanningTreeInterfacePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciSpanningTreeInterfacePolicyCreate,
		UpdateContext: resourceAciSpanningTreeInterfacePolicyUpdate,
		ReadContext:   resourceAciSpanningTreeInterfacePolicyRead,
		DeleteContext: resourceAciSpanningTreeInterfacePolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciSpanningTreeInterfacePolicyImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{

			"ctrl": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"bpdu-filter",
						"bpdu-guard",
						"unspecified",
					}, false),
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		})),
	}
}

func getRemoteSpanningTreeInterfacePolicy(client *client.Client, dn string) (*models.SpanningTreeInterfacePolicy, error) {
	stpIfPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	stpIfPol := models.SpanningTreeInterfacePolicyFromContainer(stpIfPolCont)
	if stpIfPol.DistinguishedName == "" {
		return nil, fmt.Errorf("Spanning Tree Interface Policy %s not found", dn)
	}
	return stpIfPol, nil
}

func setSpanningTreeInterfacePolicyAttributes(stpIfPol *models.SpanningTreeInterfacePolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(stpIfPol.DistinguishedName)
	d.Set("description", stpIfPol.Description)
	stpIfPolMap, err := stpIfPol.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("annotation", stpIfPolMap["annotation"])
	ctrlGet := make([]string, 0, 1)
	for _, val := range strings.Split(stpIfPolMap["ctrl"], ",") {
		if val == "" {
			ctrlGet = append(ctrlGet, "unspecified")
		} else {
			ctrlGet = append(ctrlGet, strings.Trim(val, " "))
		}
	}
	sort.Strings(ctrlGet)
	if ctrlIntr, ok := d.GetOk("ctrl"); ok {
		ctrlAct := make([]string, 0, 1)
		for _, val := range ctrlIntr.([]interface{}) {
			ctrlAct = append(ctrlAct, val.(string))
		}
		sort.Strings(ctrlAct)
		if reflect.DeepEqual(ctrlAct, ctrlGet) {
			d.Set("ctrl", d.Get("ctrl").([]interface{}))
		} else {
			d.Set("ctrl", ctrlGet)
		}
	} else {
		d.Set("ctrl", ctrlGet)
	}
	d.Set("name", stpIfPolMap["name"])
	d.Set("name_alias", stpIfPolMap["nameAlias"])
	return d, nil
}

func resourceAciSpanningTreeInterfacePolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	stpIfPol, err := getRemoteSpanningTreeInterfacePolicy(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setSpanningTreeInterfacePolicyAttributes(stpIfPol, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciSpanningTreeInterfacePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SpanningTreeInterfacePolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	stpIfPolAttr := models.SpanningTreeInterfacePolicyAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		stpIfPolAttr.Annotation = Annotation.(string)
	} else {
		stpIfPolAttr.Annotation = "{}"
	}

	if Ctrl, ok := d.GetOk("ctrl"); ok {
		ctrlList := make([]string, 0, 1)
		for _, val := range Ctrl.([]interface{}) {
			ctrlList = append(ctrlList, val.(string))
		}
		Ctrl := strings.Join(ctrlList, ",")
		stpIfPolAttr.Ctrl = Ctrl
	}

	if Name, ok := d.GetOk("name"); ok {
		stpIfPolAttr.Name = Name.(string)
	}
	stpIfPol := models.NewSpanningTreeInterfacePolicy(fmt.Sprintf("infra/ifPol-%s", name), "uni", desc, nameAlias, stpIfPolAttr)
	err := aciClient.Save(stpIfPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(stpIfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciSpanningTreeInterfacePolicyRead(ctx, d, m)
}

func resourceAciSpanningTreeInterfacePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SpanningTreeInterfacePolicy: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	stpIfPolAttr := models.SpanningTreeInterfacePolicyAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		stpIfPolAttr.Annotation = Annotation.(string)
	} else {
		stpIfPolAttr.Annotation = "{}"
	}

	if Ctrl, ok := d.GetOk("ctrl"); ok {
		ctrlList := make([]string, 0, 1)
		for _, val := range Ctrl.([]interface{}) {
			ctrlList = append(ctrlList, val.(string))
		}
		Ctrl := strings.Join(ctrlList, ",")
		stpIfPolAttr.Ctrl = Ctrl
	}

	if Name, ok := d.GetOk("name"); ok {
		stpIfPolAttr.Name = Name.(string)
	}
	stpIfPol := models.NewSpanningTreeInterfacePolicy(fmt.Sprintf("infra/ifPol-%s", name), "uni", desc, nameAlias, stpIfPolAttr)
	stpIfPol.Status = "modified"
	err := aciClient.Save(stpIfPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(stpIfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciSpanningTreeInterfacePolicyRead(ctx, d, m)
}

func resourceAciSpanningTreeInterfacePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	stpIfPol, err := getRemoteSpanningTreeInterfacePolicy(aciClient, dn)
	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}
	_, err = setSpanningTreeInterfacePolicyAttributes(stpIfPol, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciSpanningTreeInterfacePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "stpIfPol")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
