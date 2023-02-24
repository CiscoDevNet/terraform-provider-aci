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

func resourceAciLACPPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciLACPPolicyCreate,
		UpdateContext: resourceAciLACPPolicyUpdate,
		ReadContext:   resourceAciLACPPolicyRead,
		DeleteContext: resourceAciLACPPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciLACPPolicyImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

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
						"symmetric-hash",
						"susp-individual",
						"graceful-conv",
						"load-defer",
						"fast-sel-hot-stdby",
					}, false),
				},
			},

			"max_links": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"min_links": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"off",
					"active",
					"passive",
					"mac-pin",
					"mac-pin-nicload",
				}, false),
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}
func getRemoteLACPPolicy(client *client.Client, dn string) (*models.LACPPolicy, error) {
	lacpLagPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	lacpLagPol := models.LACPPolicyFromContainer(lacpLagPolCont)

	if lacpLagPol.DistinguishedName == "" {
		return nil, fmt.Errorf("LACP Policy %s not found", dn)
	}

	return lacpLagPol, nil
}

func setLACPPolicyAttributes(lacpLagPol *models.LACPPolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(lacpLagPol.DistinguishedName)
	d.Set("description", lacpLagPol.Description)
	lacpLagPolMap, err := lacpLagPol.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("name", lacpLagPolMap["name"])

	d.Set("annotation", lacpLagPolMap["annotation"])
	ctrlsGet := make([]string, 0, 1)
	for _, val := range strings.Split(lacpLagPolMap["ctrl"], ",") {
		ctrlsGet = append(ctrlsGet, strings.Trim(val, " "))
	}
	sort.Strings(ctrlsGet)
	if _, ok := d.GetOk("ctrl"); ok {
		ctrlsAct := make([]string, 0, 1)
		for _, val := range d.Get("ctrl").([]interface{}) {
			ctrlsAct = append(ctrlsAct, val.(string))
		}
		sort.Strings(ctrlsAct)
		if reflect.DeepEqual(ctrlsAct, ctrlsGet) {
			d.Set("ctrl", d.Get("ctrl").([]interface{}))
		} else {
			d.Set("ctrl", ctrlsGet)
		}
	} else {
		d.Set("ctrl", ctrlsGet)
	}
	d.Set("max_links", lacpLagPolMap["maxLinks"])
	d.Set("min_links", lacpLagPolMap["minLinks"])
	d.Set("mode", lacpLagPolMap["mode"])
	d.Set("name_alias", lacpLagPolMap["nameAlias"])
	return d, nil
}

func resourceAciLACPPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	lacpLagPol, err := getRemoteLACPPolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setLACPPolicyAttributes(lacpLagPol, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciLACPPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LACPPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	lacpLagPolAttr := models.LACPPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		lacpLagPolAttr.Annotation = Annotation.(string)
	} else {
		lacpLagPolAttr.Annotation = "{}"
	}
	if Ctrl, ok := d.GetOk("ctrl"); ok {
		tp := Ctrl.([]interface{})
		ctrls := make([]string, 0, 1)
		for _, val := range tp {
			ctrls = append(ctrls, val.(string))
		}
		lacpLagPolAttr.Ctrl = strings.Join(ctrls, ",")
	}
	if MaxLinks, ok := d.GetOk("max_links"); ok {
		lacpLagPolAttr.MaxLinks = MaxLinks.(string)
	}
	if MinLinks, ok := d.GetOk("min_links"); ok {
		lacpLagPolAttr.MinLinks = MinLinks.(string)
	}
	if Mode, ok := d.GetOk("mode"); ok {
		lacpLagPolAttr.Mode = Mode.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		lacpLagPolAttr.NameAlias = NameAlias.(string)
	}
	lacpLagPol := models.NewLACPPolicy(fmt.Sprintf("infra/lacplagp-%s", name), "uni", desc, lacpLagPolAttr)

	err := aciClient.Save(lacpLagPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(lacpLagPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciLACPPolicyRead(ctx, d, m)
}

func resourceAciLACPPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LACPPolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	lacpLagPolAttr := models.LACPPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		lacpLagPolAttr.Annotation = Annotation.(string)
	} else {
		lacpLagPolAttr.Annotation = "{}"
	}
	if Ctrl, ok := d.GetOk("ctrl"); ok {
		tp := Ctrl.([]interface{})
		ctrls := make([]string, 0, 1)
		for _, val := range tp {
			ctrls = append(ctrls, val.(string))
		}
		lacpLagPolAttr.Ctrl = strings.Join(ctrls, ",")
	}
	if MaxLinks, ok := d.GetOk("max_links"); ok {
		lacpLagPolAttr.MaxLinks = MaxLinks.(string)
	}
	if MinLinks, ok := d.GetOk("min_links"); ok {
		lacpLagPolAttr.MinLinks = MinLinks.(string)
	}
	if Mode, ok := d.GetOk("mode"); ok {
		lacpLagPolAttr.Mode = Mode.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		lacpLagPolAttr.NameAlias = NameAlias.(string)
	}
	lacpLagPol := models.NewLACPPolicy(fmt.Sprintf("infra/lacplagp-%s", name), "uni", desc, lacpLagPolAttr)

	lacpLagPol.Status = "modified"

	err := aciClient.Save(lacpLagPol)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(lacpLagPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciLACPPolicyRead(ctx, d, m)

}

func resourceAciLACPPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	lacpLagPol, err := getRemoteLACPPolicy(aciClient, dn)

	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}
	_, err = setLACPPolicyAttributes(lacpLagPol, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciLACPPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "lacpLagPol")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
