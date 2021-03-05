package aci

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAciSpanningTreeInterfacePolicy() *schema.Resource {

	return &schema.Resource{
		Create: resourceAciSpanningTreeInterfacePolicyCreate,
		Update: resourceAciSpanningTreeInterfacePolicyUpdate,
		Read:   resourceAciSpanningTreeInterfacePolicyRead,
		Delete: resourceAciSpanningTreeInterfacePolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciSpanningTreeInterfacePolicyImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ctrl": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"bpdu-filter",
						"bpdu-guard",
					}, false),
				},
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
		return nil, fmt.Errorf("SpanningTreeInterfacePolicy %s not found", stpIfPol.DistinguishedName)
	}
	return stpIfPol, nil
}

func setSpanningTreeInterfacePolicyAttributes(stpIfPol *models.SpanningTreeInterfacePolicy, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(stpIfPol.DistinguishedName)
	d.Set("description", stpIfPol.Description)
	stpIfPolMap, _ := stpIfPol.ToMap()
	d.Set("name", stpIfPolMap["name"])
	d.Set("annotation", stpIfPolMap["annotation"])
	ctrlGet := make([]string, 0, 1)
	for _, val := range strings.Split(stpIfPolMap["ctrl"], ",") {
		ctrlGet = append(ctrlGet, strings.Trim(val, " "))
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
	d.Set("name_alias", stpIfPolMap["nameAlias"])
	return d
}

func resourceAciSpanningTreeInterfacePolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	stpIfPol, err := getRemoteSpanningTreeInterfacePolicy(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled := setSpanningTreeInterfacePolicyAttributes(stpIfPol, d)
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciSpanningTreeInterfacePolicyCreate(d *schema.ResourceData, m interface{}) error {
	return resourceAciSpanningTreeInterfacePolicyCreateOrUpdate(d, m, false)
}

func resourceAciSpanningTreeInterfacePolicyUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceAciSpanningTreeInterfacePolicyCreateOrUpdate(d, m, true)
}

func resourceAciSpanningTreeInterfacePolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	stpIfPol, err := getRemoteSpanningTreeInterfacePolicy(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}
	setSpanningTreeInterfacePolicyAttributes(stpIfPol, d)
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciSpanningTreeInterfacePolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "stpIfPol")
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return err
}
func resourceAciSpanningTreeInterfacePolicyCreateOrUpdate(d *schema.ResourceData, m interface{}, update bool) error {
	action := "Creation"
	if update == true {
		action = "Update"
	}
	log.Printf("[DEBUG] SpanningTreeInterfacePolicy: Beginning %s", action)
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	stpIfPolAttr := models.SpanningTreeInterfacePolicyAttributes{}
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
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	stpIfPol := models.NewSpanningTreeInterfacePolicy(fmt.Sprintf("infra/ifPol-%s", name), "uni", desc, nameAlias, stpIfPolAttr)
	if update == true {
		stpIfPol.Status = "modified"
	}
	err := aciClient.Save(stpIfPol)
	if err != nil {
		return err
	}
	d.Partial(true)
	d.SetPartial("name")
	d.Partial(false)
	d.SetId(stpIfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: %s finished successfully", d.Id(), action)
	return resourceAciSpanningTreeInterfacePolicyRead(d, m)
}
