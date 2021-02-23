package aci

import (
	"fmt"
	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func resourceAciEndpointSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciEndpointSecurityGroupCreate,
		Update: resourceAciEndpointSecurityGroupUpdate,
		Read:   resourceAciEndpointSecurityGroupRead,
		Delete: resourceAciEndpointSecurityGroupDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciEndpointSecurityGroupImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"application_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"exception_tag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"flood_on_encap": &schema.Schema{
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

			"pc_enf_pref": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"pref_gr_memb": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"prio": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"userdom": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"relation_fv_rs_sec_inherited": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_fv_rs_prov": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_fv_rs_cons_if": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_fv_rs_cust_qos_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_fv_rs_cons": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_fv_rs_scope": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_fv_rs_prot_by": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_fv_rs_intra_epg": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
		}),
	}
}
func getRemoteEndpointSecurityGroup(client *client.Client, dn string) (*models.EndpointSecurityGroup, error) {
	fvESgCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	fvESg := models.EndpointSecurityGroupFromContainer(fvESgCont)

	if fvESg.DistinguishedName == "" {
		return nil, fmt.Errorf("EndpointSecurityGroup %s not found", fvESg.DistinguishedName)
	}

	return fvESg, nil
}

func setEndpointSecurityGroupAttributes(fvESg *models.EndpointSecurityGroup, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(fvESg.DistinguishedName)
	d.Set("description", fvESg.Description)
	fvESgMap, _ := fvESg.ToMap()

	d.Set("name", fvESgMap["name"])
	d.Set("application_profile_dn", GetParentDn(fvESg.DistinguishedName, fmt.Sprintf("/esg-%s", fvESgMap["name"])))

	d.Set("annotation", fvESgMap["annotation"])
	d.Set("exception_tag", fvESgMap["exceptionTag"])
	d.Set("flood_on_encap", fvESgMap["floodOnEncap"])
	d.Set("match_t", fvESgMap["matchT"])
	d.Set("name_alias", fvESgMap["nameAlias"])
	d.Set("pc_enf_pref", fvESgMap["pcEnfPref"])
	d.Set("pref_gr_memb", fvESgMap["prefGrMemb"])
	d.Set("prio", fvESgMap["prio"])
	d.Set("userdom", fvESgMap["userdom"])
	return d
}

func resourceAciEndpointSecurityGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	fvESg, err := getRemoteEndpointSecurityGroup(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setEndpointSecurityGroupAttributes(fvESg, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciEndpointSecurityGroupCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] EndpointSecurityGroup: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	ApplicationProfileDn := d.Get("application_profile_dn").(string)

	fvESgAttr := models.EndpointSecurityGroupAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvESgAttr.Annotation = Annotation.(string)
	}
	if ExceptionTag, ok := d.GetOk("exception_tag"); ok {
		fvESgAttr.ExceptionTag = ExceptionTag.(string)
	}
	if FloodOnEncap, ok := d.GetOk("flood_on_encap"); ok {
		fvESgAttr.FloodOnEncap = FloodOnEncap.(string)
	}
	if MatchT, ok := d.GetOk("match_t"); ok {
		fvESgAttr.MatchT = MatchT.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fvESgAttr.NameAlias = NameAlias.(string)
	}
	if PcEnfPref, ok := d.GetOk("pc_enf_pref"); ok {
		fvESgAttr.PcEnfPref = PcEnfPref.(string)
	}
	if PrefGrMemb, ok := d.GetOk("pref_gr_memb"); ok {
		fvESgAttr.PrefGrMemb = PrefGrMemb.(string)
	}
	if Prio, ok := d.GetOk("prio"); ok {
		fvESgAttr.Prio = Prio.(string)
	}
	if Userdom, ok := d.GetOk("userdom"); ok {
		fvESgAttr.Userdom = Userdom.(string)
	}
	fvESg := models.NewEndpointSecurityGroup(fmt.Sprintf("esg-%s", name), ApplicationProfileDn, desc, fvESgAttr)

	err := aciClient.Save(fvESg)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if relationTofvRsSecInherited, ok := d.GetOk("relation_fv_rs_sec_inherited"); ok {
		relationParamList := toStringList(relationTofvRsSecInherited.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsSecInheritedFromEndpointSecurityGroup(fvESg.DistinguishedName, relationParam)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_fv_rs_sec_inherited")
			d.Partial(false)
		}
	}
	if relationTofvRsProv, ok := d.GetOk("relation_fv_rs_prov"); ok {
		relationParamList := toStringList(relationTofvRsProv.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsProvFromEndpointSecurityGroup(fvESg.DistinguishedName, relationParam)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_fv_rs_prov")
			d.Partial(false)
		}
	}
	if relationTofvRsConsIf, ok := d.GetOk("relation_fv_rs_cons_if"); ok {
		relationParamList := toStringList(relationTofvRsConsIf.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsConsIfFromEndpointSecurityGroup(fvESg.DistinguishedName, relationParam)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_fv_rs_cons_if")
			d.Partial(false)
		}
	}
	if relationTofvRsCustQosPol, ok := d.GetOk("relation_fv_rs_cust_qos_pol"); ok {
		relationParam := relationTofvRsCustQosPol.(string)
		err = aciClient.CreateRelationfvRsCustQosPolFromEndpointSecurityGroup(fvESg.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fv_rs_cust_qos_pol")
		d.Partial(false)

	}
	if relationTofvRsCons, ok := d.GetOk("relation_fv_rs_cons"); ok {
		relationParamList := toStringList(relationTofvRsCons.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsConsFromEndpointSecurityGroup(fvESg.DistinguishedName, relationParam)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_fv_rs_cons")
			d.Partial(false)
		}
	}
	if relationTofvRsScope, ok := d.GetOk("relation_fv_rs_scope"); ok {
		relationParam := relationTofvRsScope.(string)
		err = aciClient.CreateRelationfvRsScopeFromEndpointSecurityGroup(fvESg.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fv_rs_scope")
		d.Partial(false)

	}
	if relationTofvRsProtBy, ok := d.GetOk("relation_fv_rs_prot_by"); ok {
		relationParamList := toStringList(relationTofvRsProtBy.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsProtByFromEndpointSecurityGroup(fvESg.DistinguishedName, relationParam)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_fv_rs_prot_by")
			d.Partial(false)
		}
	}
	if relationTofvRsIntraEpg, ok := d.GetOk("relation_fv_rs_intra_epg"); ok {
		relationParamList := toStringList(relationTofvRsIntraEpg.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsIntraEpgFromEndpointSecurityGroup(fvESg.DistinguishedName, relationParam)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_fv_rs_intra_epg")
			d.Partial(false)
		}
	}

	d.SetId(fvESg.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciEndpointSecurityGroupRead(d, m)
}

func resourceAciEndpointSecurityGroupUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] EndpointSecurityGroup: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	ApplicationProfileDn := d.Get("application_profile_dn").(string)

	fvESgAttr := models.EndpointSecurityGroupAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvESgAttr.Annotation = Annotation.(string)
	}
	if ExceptionTag, ok := d.GetOk("exception_tag"); ok {
		fvESgAttr.ExceptionTag = ExceptionTag.(string)
	}
	if FloodOnEncap, ok := d.GetOk("flood_on_encap"); ok {
		fvESgAttr.FloodOnEncap = FloodOnEncap.(string)
	}
	if MatchT, ok := d.GetOk("match_t"); ok {
		fvESgAttr.MatchT = MatchT.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fvESgAttr.NameAlias = NameAlias.(string)
	}
	if PcEnfPref, ok := d.GetOk("pc_enf_pref"); ok {
		fvESgAttr.PcEnfPref = PcEnfPref.(string)
	}
	if PrefGrMemb, ok := d.GetOk("pref_gr_memb"); ok {
		fvESgAttr.PrefGrMemb = PrefGrMemb.(string)
	}
	if Prio, ok := d.GetOk("prio"); ok {
		fvESgAttr.Prio = Prio.(string)
	}
	if Userdom, ok := d.GetOk("userdom"); ok {
		fvESgAttr.Userdom = Userdom.(string)
	}
	fvESg := models.NewEndpointSecurityGroup(fmt.Sprintf("esg-%s", name), ApplicationProfileDn, desc, fvESgAttr)

	fvESg.Status = "modified"

	err := aciClient.Save(fvESg)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if d.HasChange("relation_fv_rs_sec_inherited") {
		oldRel, newRel := d.GetChange("relation_fv_rs_sec_inherited")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationfvRsSecInheritedFromEndpointSecurityGroup(fvESg.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsSecInheritedFromEndpointSecurityGroup(fvESg.DistinguishedName, relDn)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_fv_rs_sec_inherited")
			d.Partial(false)

		}

	}
	if d.HasChange("relation_fv_rs_prov") {
		oldRel, newRel := d.GetChange("relation_fv_rs_prov")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationfvRsProvFromEndpointSecurityGroup(fvESg.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsProvFromEndpointSecurityGroup(fvESg.DistinguishedName, relDn)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_fv_rs_prov")
			d.Partial(false)

		}

	}
	if d.HasChange("relation_fv_rs_cons_if") {
		oldRel, newRel := d.GetChange("relation_fv_rs_cons_if")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationfvRsConsIfFromEndpointSecurityGroup(fvESg.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsConsIfFromEndpointSecurityGroup(fvESg.DistinguishedName, relDn)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_fv_rs_cons_if")
			d.Partial(false)

		}

	}
	if d.HasChange("relation_fv_rs_cust_qos_pol") {
		_, newRelParam := d.GetChange("relation_fv_rs_cust_qos_pol")
		err = aciClient.CreateRelationfvRsCustQosPolFromEndpointSecurityGroup(fvESg.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fv_rs_cust_qos_pol")
		d.Partial(false)

	}
	if d.HasChange("relation_fv_rs_cons") {
		oldRel, newRel := d.GetChange("relation_fv_rs_cons")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationfvRsConsFromEndpointSecurityGroup(fvESg.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsConsFromEndpointSecurityGroup(fvESg.DistinguishedName, relDn)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_fv_rs_cons")
			d.Partial(false)

		}

	}
	if d.HasChange("relation_fv_rs_scope") {
		_, newRelParam := d.GetChange("relation_fv_rs_scope")
		err = aciClient.CreateRelationfvRsScopeFromEndpointSecurityGroup(fvESg.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fv_rs_scope")
		d.Partial(false)

	}
	if d.HasChange("relation_fv_rs_prot_by") {
		oldRel, newRel := d.GetChange("relation_fv_rs_prot_by")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationfvRsProtByFromEndpointSecurityGroup(fvESg.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsProtByFromEndpointSecurityGroup(fvESg.DistinguishedName, relDn)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_fv_rs_prot_by")
			d.Partial(false)

		}

	}
	if d.HasChange("relation_fv_rs_intra_epg") {
		oldRel, newRel := d.GetChange("relation_fv_rs_intra_epg")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationfvRsIntraEpgFromEndpointSecurityGroup(fvESg.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsIntraEpgFromEndpointSecurityGroup(fvESg.DistinguishedName, relDn)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_fv_rs_intra_epg")
			d.Partial(false)

		}

	}

	d.SetId(fvESg.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciEndpointSecurityGroupRead(d, m)

}

func resourceAciEndpointSecurityGroupRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	fvESg, err := getRemoteEndpointSecurityGroup(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setEndpointSecurityGroupAttributes(fvESg, d)

	fvRsSecInheritedData, err := aciClient.ReadRelationfvRsSecInheritedFromEndpointSecurityGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsSecInherited %v", err)

	} else {
		d.Set("relation_fv_rs_sec_inherited", fvRsSecInheritedData)
	}

	fvRsProvData, err := aciClient.ReadRelationfvRsProvFromEndpointSecurityGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsProv %v", err)

	} else {
		d.Set("relation_fv_rs_prov", fvRsProvData)
	}

	fvRsConsIfData, err := aciClient.ReadRelationfvRsConsIfFromEndpointSecurityGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsConsIf %v", err)

	} else {
		d.Set("relation_fv_rs_cons_if", fvRsConsIfData)
	}

	fvRsCustQosPolData, err := aciClient.ReadRelationfvRsCustQosPolFromEndpointSecurityGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsCustQosPol %v", err)

	} else {
		d.Set("relation_fv_rs_cust_qos_pol", fvRsCustQosPolData)
	}

	fvRsConsData, err := aciClient.ReadRelationfvRsConsFromEndpointSecurityGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsCons %v", err)

	} else {
		d.Set("relation_fv_rs_cons", fvRsConsData)
	}

	fvRsScopeData, err := aciClient.ReadRelationfvRsScopeFromEndpointSecurityGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsScope %v", err)

	} else {
		d.Set("relation_fv_rs_scope", fvRsScopeData)
	}

	fvRsProtByData, err := aciClient.ReadRelationfvRsProtByFromEndpointSecurityGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsProtBy %v", err)

	} else {
		d.Set("relation_fv_rs_prot_by", fvRsProtByData)
	}

	fvRsIntraEpgData, err := aciClient.ReadRelationfvRsIntraEpgFromEndpointSecurityGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsIntraEpg %v", err)

	} else {
		d.Set("relation_fv_rs_intra_epg", fvRsIntraEpgData)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciEndpointSecurityGroupDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fvESg")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
