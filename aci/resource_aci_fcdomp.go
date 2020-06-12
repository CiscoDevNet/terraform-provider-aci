package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciFCDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciFCDomainCreate,
		Update: resourceAciFCDomainUpdate,
		Read:   resourceAciFCDomainRead,
		Delete: resourceAciFCDomainDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciFCDomainImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"relation_infra_rs_vlan_ns": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_fc_rs_vsan_ns": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_fc_rs_vsan_attr": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_infra_rs_vlan_ns_def": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_infra_rs_vip_addr_ns": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_infra_rs_dom_vxlan_ns_def": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_fc_rs_vsan_attr_def": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_fc_rs_vsan_ns_def": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
		}),
	}
}
func getRemoteFCDomain(client *client.Client, dn string) (*models.FCDomain, error) {
	fcDomPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	fcDomP := models.FCDomainFromContainer(fcDomPCont)

	if fcDomP.DistinguishedName == "" {
		return nil, fmt.Errorf("FCDomain %s not found", fcDomP.DistinguishedName)
	}

	return fcDomP, nil
}

func setFCDomainAttributes(fcDomP *models.FCDomain, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(fcDomP.DistinguishedName)
	d.Set("description", fcDomP.Description)
	fcDomPMap, _ := fcDomP.ToMap()

	d.Set("name", fcDomPMap["name"])

	d.Set("annotation", fcDomPMap["annotation"])
	d.Set("name_alias", fcDomPMap["nameAlias"])
	return d
}

func resourceAciFCDomainImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	fcDomP, err := getRemoteFCDomain(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setFCDomainAttributes(fcDomP, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciFCDomainCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] FCDomain: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	fcDomPAttr := models.FCDomainAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fcDomPAttr.Annotation = Annotation.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fcDomPAttr.NameAlias = NameAlias.(string)
	}
	fcDomP := models.NewFCDomain(fmt.Sprintf("fc-%s", name), "uni", desc, fcDomPAttr)

	err := aciClient.Save(fcDomP)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if relationToinfraRsVlanNs, ok := d.GetOk("relation_infra_rs_vlan_ns"); ok {
		relationParam := relationToinfraRsVlanNs.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsVlanNsFromFCDomain(fcDomP.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_vlan_ns")
		d.Partial(false)

	}
	if relationTofcRsVsanNs, ok := d.GetOk("relation_fc_rs_vsan_ns"); ok {
		relationParam := relationTofcRsVsanNs.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfcRsVsanNsFromFCDomain(fcDomP.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fc_rs_vsan_ns")
		d.Partial(false)

	}
	if relationTofcRsVsanAttr, ok := d.GetOk("relation_fc_rs_vsan_attr"); ok {
		relationParam := relationTofcRsVsanAttr.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfcRsVsanAttrFromFCDomain(fcDomP.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fc_rs_vsan_attr")
		d.Partial(false)

	}
	if relationToinfraRsVlanNsDef, ok := d.GetOk("relation_infra_rs_vlan_ns_def"); ok {
		relationParam := relationToinfraRsVlanNsDef.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsVlanNsDefFromFCDomain(fcDomP.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_vlan_ns_def")
		d.Partial(false)

	}
	if relationToinfraRsVipAddrNs, ok := d.GetOk("relation_infra_rs_vip_addr_ns"); ok {
		relationParam := relationToinfraRsVipAddrNs.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsVipAddrNsFromFCDomain(fcDomP.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_vip_addr_ns")
		d.Partial(false)

	}
	if relationToinfraRsDomVxlanNsDef, ok := d.GetOk("relation_infra_rs_dom_vxlan_ns_def"); ok {
		relationParam := relationToinfraRsDomVxlanNsDef.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsDomVxlanNsDefFromFCDomain(fcDomP.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_dom_vxlan_ns_def")
		d.Partial(false)

	}
	if relationTofcRsVsanAttrDef, ok := d.GetOk("relation_fc_rs_vsan_attr_def"); ok {
		relationParam := relationTofcRsVsanAttrDef.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfcRsVsanAttrDefFromFCDomain(fcDomP.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fc_rs_vsan_attr_def")
		d.Partial(false)

	}
	if relationTofcRsVsanNsDef, ok := d.GetOk("relation_fc_rs_vsan_ns_def"); ok {
		relationParam := relationTofcRsVsanNsDef.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfcRsVsanNsDefFromFCDomain(fcDomP.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fc_rs_vsan_ns_def")
		d.Partial(false)

	}

	d.SetId(fcDomP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciFCDomainRead(d, m)
}

func resourceAciFCDomainUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] FCDomain: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	fcDomPAttr := models.FCDomainAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fcDomPAttr.Annotation = Annotation.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fcDomPAttr.NameAlias = NameAlias.(string)
	}
	fcDomP := models.NewFCDomain(fmt.Sprintf("fc-%s", name), "uni", desc, fcDomPAttr)

	fcDomP.Status = "modified"

	err := aciClient.Save(fcDomP)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if d.HasChange("relation_infra_rs_vlan_ns") {
		_, newRelParam := d.GetChange("relation_infra_rs_vlan_ns")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationinfraRsVlanNsFromFCDomain(fcDomP.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationinfraRsVlanNsFromFCDomain(fcDomP.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_vlan_ns")
		d.Partial(false)

	}
	if d.HasChange("relation_fc_rs_vsan_ns") {
		_, newRelParam := d.GetChange("relation_fc_rs_vsan_ns")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationfcRsVsanNsFromFCDomain(fcDomP.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationfcRsVsanNsFromFCDomain(fcDomP.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fc_rs_vsan_ns")
		d.Partial(false)

	}
	if d.HasChange("relation_fc_rs_vsan_attr") {
		_, newRelParam := d.GetChange("relation_fc_rs_vsan_attr")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationfcRsVsanAttrFromFCDomain(fcDomP.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationfcRsVsanAttrFromFCDomain(fcDomP.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fc_rs_vsan_attr")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_vlan_ns_def") {
		_, newRelParam := d.GetChange("relation_infra_rs_vlan_ns_def")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsVlanNsDefFromFCDomain(fcDomP.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_vlan_ns_def")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_vip_addr_ns") {
		_, newRelParam := d.GetChange("relation_infra_rs_vip_addr_ns")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationinfraRsVipAddrNsFromFCDomain(fcDomP.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationinfraRsVipAddrNsFromFCDomain(fcDomP.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_vip_addr_ns")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_dom_vxlan_ns_def") {
		_, newRelParam := d.GetChange("relation_infra_rs_dom_vxlan_ns_def")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsDomVxlanNsDefFromFCDomain(fcDomP.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_dom_vxlan_ns_def")
		d.Partial(false)

	}
	if d.HasChange("relation_fc_rs_vsan_attr_def") {
		_, newRelParam := d.GetChange("relation_fc_rs_vsan_attr_def")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationfcRsVsanAttrDefFromFCDomain(fcDomP.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fc_rs_vsan_attr_def")
		d.Partial(false)

	}
	if d.HasChange("relation_fc_rs_vsan_ns_def") {
		_, newRelParam := d.GetChange("relation_fc_rs_vsan_ns_def")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationfcRsVsanNsDefFromFCDomain(fcDomP.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_fc_rs_vsan_ns_def")
		d.Partial(false)

	}

	d.SetId(fcDomP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciFCDomainRead(d, m)

}

func resourceAciFCDomainRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	fcDomP, err := getRemoteFCDomain(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setFCDomainAttributes(fcDomP, d)

	infraRsVlanNsData, err := aciClient.ReadRelationinfraRsVlanNsFromFCDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsVlanNs %v", err)

	} else {
		if _, ok := d.GetOk("relation_infra_rs_vlan_ns"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_vlan_ns").(string))
			if tfName != infraRsVlanNsData {
				d.Set("relation_infra_rs_vlan_ns", "")
			}
		}
	}

	fcRsVsanNsData, err := aciClient.ReadRelationfcRsVsanNsFromFCDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fcRsVsanNs %v", err)

	} else {
		if _, ok := d.GetOk("relation_fc_rs_vsan_ns"); ok {
			tfName := GetMOName(d.Get("relation_fc_rs_vsan_ns").(string))
			if tfName != fcRsVsanNsData {
				d.Set("relation_fc_rs_vsan_ns", "")
			}
		}
	}

	fcRsVsanAttrData, err := aciClient.ReadRelationfcRsVsanAttrFromFCDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fcRsVsanAttr %v", err)

	} else {
		if _, ok := d.GetOk("relation_fc_rs_vsan_attr"); ok {
			tfName := GetMOName(d.Get("relation_fc_rs_vsan_attr").(string))
			if tfName != fcRsVsanAttrData {
				d.Set("relation_fc_rs_vsan_attr", "")
			}
		}
	}

	infraRsVlanNsDefData, err := aciClient.ReadRelationinfraRsVlanNsDefFromFCDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsVlanNsDef %v", err)

	} else {
		if _, ok := d.GetOk("relation_infra_rs_vlan_ns_def"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_vlan_ns_def").(string))
			if tfName != infraRsVlanNsDefData {
				d.Set("relation_infra_rs_vlan_ns_def", "")
			}
		}
	}

	infraRsVipAddrNsData, err := aciClient.ReadRelationinfraRsVipAddrNsFromFCDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsVipAddrNs %v", err)

	} else {
		if _, ok := d.GetOk("relation_infra_rs_vip_addr_ns"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_vip_addr_ns").(string))
			if tfName != infraRsVipAddrNsData {
				d.Set("relation_infra_rs_vip_addr_ns", "")
			}
		}
	}

	infraRsDomVxlanNsDefData, err := aciClient.ReadRelationinfraRsDomVxlanNsDefFromFCDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsDomVxlanNsDef %v", err)

	} else {
		if _, ok := d.GetOk("relation_infra_rs_dom_vxlan_ns_def"); ok {
			tfName := GetMOName(d.Get("relation_infra_rs_dom_vxlan_ns_def").(string))
			if tfName != infraRsDomVxlanNsDefData {
				d.Set("relation_infra_rs_dom_vxlan_ns_def", "")
			}
		}
	}

	fcRsVsanAttrDefData, err := aciClient.ReadRelationfcRsVsanAttrDefFromFCDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fcRsVsanAttrDef %v", err)

	} else {
		if _, ok := d.GetOk("relation_fc_rs_vsan_attr_def"); ok {
			tfName := GetMOName(d.Get("relation_fc_rs_vsan_attr_def").(string))
			if tfName != fcRsVsanAttrDefData {
				d.Set("relation_fc_rs_vsan_attr_def", "")
			}
		}
	}

	fcRsVsanNsDefData, err := aciClient.ReadRelationfcRsVsanNsDefFromFCDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fcRsVsanNsDef %v", err)

	} else {
		if _, ok := d.GetOk("relation_fc_rs_vsan_ns_def"); ok {
			tfName := GetMOName(d.Get("relation_fc_rs_vsan_ns_def").(string))
			if tfName != fcRsVsanNsDefData {
				d.Set("relation_fc_rs_vsan_ns_def", "")
			}
		}
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciFCDomainDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fcDomP")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
