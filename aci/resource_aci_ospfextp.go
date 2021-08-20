package aci

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciL3outOspfExternalPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciL3outOspfExternalPolicyCreate,
		UpdateContext: resourceAciL3outOspfExternalPolicyUpdate,
		ReadContext:   resourceAciL3outOspfExternalPolicyRead,
		DeleteContext: resourceAciL3outOspfExternalPolicyDelete,

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
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"redistribute",
						"summary",
						"suppress-fa",
						"unspecified",
					}, false),
				},
			},

			"area_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				StateFunc: func(val interface{}) string {
					if val.(string) == "backbone" {
						return "backbone"
					} else {
						numList := strings.Split(val.(string), ".")
						ip := []string{"0", "0", "0", "0"}
						if val.(string) != "" && len(numList) <= 4 {
							for i := 1; i <= len(numList); i++ {
								ip[4-i] = numList[len(numList)-i]
							}
						}
						return strings.Join(ip, ".")
					}
				},
				ValidateFunc: schema.SchemaValidateFunc(validateOspfIp()),
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

func validateOspfIp() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		// function to check if partial OSPF areaId is valid.
		val, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}
		if val == "backbone" {
			return
		}
		numList := strings.Split(val, ".")
		if len(numList) > 4 {
			es = append(es, fmt.Errorf("Invalid value for %s : %s", k, val))
			return
		}
		for _, v := range numList {
			intV, err := strconv.Atoi(v)
			if err != nil || intV > 255 {
				es = append(es, fmt.Errorf("Invalid value for %s : %s", k, val))
				return
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

func setL3outOspfExternalPolicyAttributes(ospfExtP *models.L3outOspfExternalPolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(ospfExtP.DistinguishedName)
	d.Set("description", ospfExtP.Description)
	dn := d.Id()
	if dn != ospfExtP.DistinguishedName {
		d.Set("l3_outside_dn", "")
	}
	ospfExtPMap, err := ospfExtP.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("l3_outside_dn", GetParentDn(dn, fmt.Sprintf("/ospfExtP")))

	d.Set("annotation", ospfExtPMap["annotation"])
	d.Set("area_cost", ospfExtPMap["areaCost"])
	d.Set("area_ctrl", ospfExtPMap["areaCtrl"])
	if ospfExtPMap["areaCtrl"] == "" {
		d.Set("area_ctrl", []string{
			"unspecified",
		})
	} else {
		areaCtrlGet := make([]string, 0, 1)
		for _, val := range strings.Split(ospfExtPMap["areaCtrl"], ",") {
			areaCtrlGet = append(areaCtrlGet, strings.Trim(val, " "))
		}
		sort.Strings(areaCtrlGet)
		if areaCtrlIntr, ok := d.GetOk("area_ctrl"); ok {
			areaCtrlAct := make([]string, 0, 1)
			for _, val := range areaCtrlIntr.([]interface{}) {
				areaCtrlAct = append(areaCtrlAct, val.(string))
			}
			sort.Strings(areaCtrlAct)
			if reflect.DeepEqual(areaCtrlAct, areaCtrlGet) {
				d.Set("area_ctrl", d.Get("area_ctrl").([]interface{}))
			} else {
				d.Set("area_ctrl", areaCtrlGet)
			}
		} else {
			d.Set("area_ctrl", areaCtrlGet)
		}
		d.Set("area_ctrl", ospfExtPMap["areaCtrl"])
	}
	d.Set("area_id", ospfExtPMap["areaId"])
	d.Set("area_type", ospfExtPMap["areaType"])
	d.Set("multipod_internal", ospfExtPMap["multipodInternal"])
	d.Set("name_alias", ospfExtPMap["nameAlias"])
	return d, nil
}

func resourceAciL3outOspfExternalPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	ospfExtP, err := getRemoteL3outOspfExternalPolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setL3outOspfExternalPolicyAttributes(ospfExtP, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciL3outOspfExternalPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		areaCtrlList := make([]string, 0, 1)
		for _, val := range AreaCtrl.([]interface{}) {
			areaCtrlList = append(areaCtrlList, val.(string))
		}
		AreaCtrl := strings.Join(areaCtrlList, ",")
		ospfExtPAttr.AreaCtrl = AreaCtrl
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
		return diag.FromErr(err)
	}

	d.SetId(ospfExtP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciL3outOspfExternalPolicyRead(ctx, d, m)
}

func resourceAciL3outOspfExternalPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		areaCtrlList := make([]string, 0, 1)
		for _, val := range AreaCtrl.([]interface{}) {
			areaCtrlList = append(areaCtrlList, val.(string))
		}
		AreaCtrl := strings.Join(areaCtrlList, ",")
		ospfExtPAttr.AreaCtrl = AreaCtrl
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
		return diag.FromErr(err)
	}

	d.SetId(ospfExtP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciL3outOspfExternalPolicyRead(ctx, d, m)

}

func resourceAciL3outOspfExternalPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	ospfExtP, err := getRemoteL3outOspfExternalPolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setL3outOspfExternalPolicyAttributes(ospfExtP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciL3outOspfExternalPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "ospfExtP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
