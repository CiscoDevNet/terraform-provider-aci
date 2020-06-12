package aci

import (
	"fmt"
	"log"
	"reflect"
	"sort"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciAny() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciAnyCreate,
		Update: resourceAciAnyUpdate,
		Read:   resourceAciAnyRead,
		Delete: resourceAciAnyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciAnyImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"vrf_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"match_t": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"pref_gr_memb": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"relation_vz_rs_any_to_cons": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_vz_rs_any_to_cons_if": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_vz_rs_any_to_prov": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
		}),
	}
}
func getRemoteAny(client *client.Client, dn string) (*models.Any, error) {
	vzAnyCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	vzAny := models.AnyFromContainer(vzAnyCont)

	if vzAny.DistinguishedName == "" {
		return nil, fmt.Errorf("Any %s not found", vzAny.DistinguishedName)
	}

	return vzAny, nil
}

func setAnyAttributes(vzAny *models.Any, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(vzAny.DistinguishedName)
	d.Set("description", vzAny.Description)
	// d.Set("vrf_dn", GetParentDn(vzAny.DistinguishedName))
	if dn != vzAny.DistinguishedName {
		d.Set("vrf_dn", "")
	}
	vzAnyMap, _ := vzAny.ToMap()

	d.Set("annotation", vzAnyMap["annotation"])
	d.Set("match_t", vzAnyMap["matchT"])
	d.Set("name_alias", vzAnyMap["nameAlias"])
	d.Set("pref_gr_memb", vzAnyMap["prefGrMemb"])
	return d
}

func resourceAciAnyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	vzAny, err := getRemoteAny(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setAnyAttributes(vzAny, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciAnyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Any: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	VRFDn := d.Get("vrf_dn").(string)

	vzAnyAttr := models.AnyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vzAnyAttr.Annotation = Annotation.(string)
	}
	if MatchT, ok := d.GetOk("match_t"); ok {
		vzAnyAttr.MatchT = MatchT.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vzAnyAttr.NameAlias = NameAlias.(string)
	}
	if PrefGrMemb, ok := d.GetOk("pref_gr_memb"); ok {
		vzAnyAttr.PrefGrMemb = PrefGrMemb.(string)
	}
	vzAny := models.NewAny(fmt.Sprintf("any"), VRFDn, desc, vzAnyAttr)

	err := aciClient.Save(vzAny)
	if err != nil {
		return err
	}
	d.Partial(true)
	d.Partial(false)

	if relationTovzRsAnyToCons, ok := d.GetOk("relation_vz_rs_any_to_cons"); ok {
		relationParamList := toStringList(relationTovzRsAnyToCons.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationvzRsAnyToConsFromAny(vzAny.DistinguishedName, relationParamName)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_vz_rs_any_to_cons")
			d.Partial(false)
		}
	}
	if relationTovzRsAnyToConsIf, ok := d.GetOk("relation_vz_rs_any_to_cons_if"); ok {
		relationParamList := toStringList(relationTovzRsAnyToConsIf.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationvzRsAnyToConsIfFromAny(vzAny.DistinguishedName, relationParamName)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_vz_rs_any_to_cons_if")
			d.Partial(false)
		}
	}
	if relationTovzRsAnyToProv, ok := d.GetOk("relation_vz_rs_any_to_prov"); ok {
		relationParamList := toStringList(relationTovzRsAnyToProv.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			relationParamName := GetMOName(relationParam)
			err = aciClient.CreateRelationvzRsAnyToProvFromAny(vzAny.DistinguishedName, relationParamName)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_vz_rs_any_to_prov")
			d.Partial(false)
		}
	}

	d.SetId(vzAny.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciAnyRead(d, m)
}

func resourceAciAnyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Any: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	VRFDn := d.Get("vrf_dn").(string)

	vzAnyAttr := models.AnyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vzAnyAttr.Annotation = Annotation.(string)
	}
	if MatchT, ok := d.GetOk("match_t"); ok {
		vzAnyAttr.MatchT = MatchT.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vzAnyAttr.NameAlias = NameAlias.(string)
	}
	if PrefGrMemb, ok := d.GetOk("pref_gr_memb"); ok {
		vzAnyAttr.PrefGrMemb = PrefGrMemb.(string)
	}
	vzAny := models.NewAny(fmt.Sprintf("any"), VRFDn, desc, vzAnyAttr)

	vzAny.Status = "modified"

	err := aciClient.Save(vzAny)

	if err != nil {
		return err
	}
	d.Partial(true)
	d.Partial(false)

	if d.HasChange("relation_vz_rs_any_to_cons") {
		oldRel, newRel := d.GetChange("relation_vz_rs_any_to_cons")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			relDnName := GetMOName(relDn)
			err = aciClient.DeleteRelationvzRsAnyToConsFromAny(vzAny.DistinguishedName, relDnName)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationvzRsAnyToConsFromAny(vzAny.DistinguishedName, relDnName)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_vz_rs_any_to_cons")
			d.Partial(false)

		}

	}
	if d.HasChange("relation_vz_rs_any_to_cons_if") {
		oldRel, newRel := d.GetChange("relation_vz_rs_any_to_cons_if")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			relDnName := GetMOName(relDn)
			err = aciClient.DeleteRelationvzRsAnyToConsIfFromAny(vzAny.DistinguishedName, relDnName)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationvzRsAnyToConsIfFromAny(vzAny.DistinguishedName, relDnName)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_vz_rs_any_to_cons_if")
			d.Partial(false)

		}

	}
	if d.HasChange("relation_vz_rs_any_to_prov") {
		oldRel, newRel := d.GetChange("relation_vz_rs_any_to_prov")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			relDnName := GetMOName(relDn)
			err = aciClient.DeleteRelationvzRsAnyToProvFromAny(vzAny.DistinguishedName, relDnName)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			relDnName := GetMOName(relDn)
			err = aciClient.CreateRelationvzRsAnyToProvFromAny(vzAny.DistinguishedName, relDnName)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_vz_rs_any_to_prov")
			d.Partial(false)

		}

	}

	d.SetId(vzAny.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciAnyRead(d, m)

}

func resourceAciAnyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	vzAny, err := getRemoteAny(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setAnyAttributes(vzAny, d)

	vzRsAnyToConsData, err := aciClient.ReadRelationvzRsAnyToConsFromAny(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vzRsAnyToCons %v", err)

	} else {
		if _, ok := d.GetOk("relation_vz_rs_any_to_cons"); ok {
			relationParamList := toStringList(d.Get("relation_vz_rs_any_to_cons").(*schema.Set).List())
			tfList := make([]string, 0, 1)
			for _, relationParam := range relationParamList {
				relationParamName := GetMOName(relationParam)
				tfList = append(tfList, relationParamName)
			}
			vzRsAnyToConsDataList := toStringList(vzRsAnyToConsData.(*schema.Set).List())
			sort.Strings(tfList)
			sort.Strings(vzRsAnyToConsDataList)

			if !reflect.DeepEqual(tfList, vzRsAnyToConsDataList) {
				d.Set("relation_vz_rs_any_to_cons", make([]string, 0, 1))
			}
		}
	}

	vzRsAnyToConsIfData, err := aciClient.ReadRelationvzRsAnyToConsIfFromAny(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vzRsAnyToConsIf %v", err)

	} else {
		if _, ok := d.GetOk("relation_vz_rs_any_to_cons_if"); ok {
			relationParamList := toStringList(d.Get("relation_vz_rs_any_to_cons_if").(*schema.Set).List())
			tfList := make([]string, 0, 1)
			for _, relationParam := range relationParamList {
				relationParamName := GetMOName(relationParam)
				tfList = append(tfList, relationParamName)
			}
			vzRsAnyToConsIfDataList := toStringList(vzRsAnyToConsIfData.(*schema.Set).List())
			sort.Strings(tfList)
			sort.Strings(vzRsAnyToConsIfDataList)

			if !reflect.DeepEqual(tfList, vzRsAnyToConsIfDataList) {
				d.Set("relation_vz_rs_any_to_cons_if", make([]string, 0, 1))
			}
		}
	}

	vzRsAnyToProvData, err := aciClient.ReadRelationvzRsAnyToProvFromAny(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vzRsAnyToProv %v", err)

	} else {
		if _, ok := d.GetOk("relation_vz_rs_any_to_prov"); ok {
			relationParamList := toStringList(d.Get("relation_vz_rs_any_to_prov").(*schema.Set).List())
			tfList := make([]string, 0, 1)
			for _, relationParam := range relationParamList {
				relationParamName := GetMOName(relationParam)
				tfList = append(tfList, relationParamName)
			}
			vzRsAnyToProvDataList := toStringList(vzRsAnyToProvData.(*schema.Set).List())
			sort.Strings(tfList)
			sort.Strings(vzRsAnyToProvDataList)

			if !reflect.DeepEqual(tfList, vzRsAnyToProvDataList) {
				d.Set("relation_vz_rs_any_to_prov", make([]string, 0, 1))
			}
		}
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciAnyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "vzAny")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
