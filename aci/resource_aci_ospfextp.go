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

func resourceAciL3outOspfExternalPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciL3outOspfExternalPolicyCreate,
		Update: resourceAciL3outOspfExternalPolicyUpdate,
		Read:   resourceAciL3outOspfExternalPolicyRead,
		Delete: resourceAciL3outOspfExternalPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciL3outOspfExternalPolicyImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"l3_outside_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"area_cost": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"area_ctrl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					oldList := strings.Split(old, ",")
					newList := strings.Split(new, ",")
					sort.Strings(oldList)
					sort.Strings(newList)

					return reflect.DeepEqual(oldList, newList)
				},
				ValidateFunc: schema.SchemaValidateFunc(validateCommaSeparatedStringInSlice([]string{
					"redistribute",
					"summary",
					"suppress-fa",
					"unspecified",
				}, false)),
			},

			"area_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"area_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"nssa",
					"regular",
					"stub",
				}, false),
			},

			"multipod_internal": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
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

func validateCommaSeparatedStringInSlice(valid []string, ignoreCase bool) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		// modified validation.StringInSlice function.
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}
		vals := strings.Split(v, ",")
		for _, val := range vals {
			match := false
			for _, str := range valid {
				if val == str || (ignoreCase && strings.ToLower(val) == strings.ToLower(str)) {
					match = true
				}
			}
			if !match {
				es = append(es, fmt.Errorf("expected %s to be one of %v, got %s", k, valid, val))
			}
		}
		return
	}
}

func getRemoteL3outOspfExternalPolicy(client *client.Client, dn string) (*models.L3outOspfExternalPolicy, error) {
	ospfExtPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	ospfExtP := models.L3outOspfExternalPolicyFromContainer(ospfExtPCont)

	if ospfExtP.DistinguishedName == "" {
		return nil, fmt.Errorf("L3outOspfExternalPolicy %s not found", ospfExtP.DistinguishedName)
	}

	return ospfExtP, nil
}

func setL3outOspfExternalPolicyAttributes(ospfExtP *models.L3outOspfExternalPolicy, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(ospfExtP.DistinguishedName)
	d.Set("description", ospfExtP.Description)
	dn := d.Id()
	if dn != ospfExtP.DistinguishedName {
		d.Set("l3_outside_dn", "")
	}
	ospfExtPMap, _ := ospfExtP.ToMap()

	d.Set("annotation", ospfExtPMap["annotation"])
	d.Set("area_cost", ospfExtPMap["areaCost"])
	d.Set("area_ctrl", ospfExtPMap["areaCtrl"])
	d.Set("area_id", ospfExtPMap["areaId"])
	d.Set("area_type", ospfExtPMap["areaType"])
	d.Set("multipod_internal", ospfExtPMap["multipodInternal"])
	d.Set("name_alias", ospfExtPMap["nameAlias"])
	return d
}

func resourceAciL3outOspfExternalPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	ospfExtP, err := getRemoteL3outOspfExternalPolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setL3outOspfExternalPolicyAttributes(ospfExtP, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciL3outOspfExternalPolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] L3outOspfExternalPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	L3OutsideDn := d.Get("l3_outside_dn").(string)

	ospfExtPAttr := models.L3outOspfExternalPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		ospfExtPAttr.Annotation = Annotation.(string)
	} else {
		ospfExtPAttr.Annotation = "{}"
	}
	if AreaCost, ok := d.GetOk("area_cost"); ok {
		ospfExtPAttr.AreaCost = AreaCost.(string)
	}
	if AreaCtrl, ok := d.GetOk("area_ctrl"); ok {
		ospfExtPAttr.AreaCtrl = AreaCtrl.(string)
	}
	if AreaId, ok := d.GetOk("area_id"); ok {
		ospfExtPAttr.AreaId = AreaId.(string)
	}
	if AreaType, ok := d.GetOk("area_type"); ok {
		ospfExtPAttr.AreaType = AreaType.(string)
	}
	if MultipodInternal, ok := d.GetOk("multipod_internal"); ok {
		ospfExtPAttr.MultipodInternal = MultipodInternal.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		ospfExtPAttr.NameAlias = NameAlias.(string)
	}
	ospfExtP := models.NewL3outOspfExternalPolicy(fmt.Sprintf("ospfExtP"), L3OutsideDn, desc, ospfExtPAttr)

	err := aciClient.Save(ospfExtP)
	if err != nil {
		return err
	}
	d.Partial(true)
	d.Partial(false)

	d.SetId(ospfExtP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciL3outOspfExternalPolicyRead(d, m)
}

func resourceAciL3outOspfExternalPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] L3outOspfExternalPolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	L3OutsideDn := d.Get("l3_outside_dn").(string)

	ospfExtPAttr := models.L3outOspfExternalPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		ospfExtPAttr.Annotation = Annotation.(string)
	} else {
		ospfExtPAttr.Annotation = "{}"
	}
	if AreaCost, ok := d.GetOk("area_cost"); ok {
		ospfExtPAttr.AreaCost = AreaCost.(string)
	}
	if AreaCtrl, ok := d.GetOk("area_ctrl"); ok {
		ospfExtPAttr.AreaCtrl = AreaCtrl.(string)
	}
	if AreaId, ok := d.GetOk("area_id"); ok {
		ospfExtPAttr.AreaId = AreaId.(string)
	}
	if AreaType, ok := d.GetOk("area_type"); ok {
		ospfExtPAttr.AreaType = AreaType.(string)
	}
	if MultipodInternal, ok := d.GetOk("multipod_internal"); ok {
		ospfExtPAttr.MultipodInternal = MultipodInternal.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		ospfExtPAttr.NameAlias = NameAlias.(string)
	}
	ospfExtP := models.NewL3outOspfExternalPolicy(fmt.Sprintf("ospfExtP"), L3OutsideDn, desc, ospfExtPAttr)

	ospfExtP.Status = "modified"

	err := aciClient.Save(ospfExtP)

	if err != nil {
		return err
	}
	d.Partial(true)
	d.Partial(false)

	d.SetId(ospfExtP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciL3outOspfExternalPolicyRead(d, m)

}

func resourceAciL3outOspfExternalPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	ospfExtP, err := getRemoteL3outOspfExternalPolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setL3outOspfExternalPolicyAttributes(ospfExtP, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciL3outOspfExternalPolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "ospfExtP")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
