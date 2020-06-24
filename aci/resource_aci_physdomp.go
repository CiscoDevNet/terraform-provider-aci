package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciPhysicalDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciPhysicalDomainCreate,
		Update: resourceAciPhysicalDomainUpdate,
		Read:   resourceAciPhysicalDomainRead,
		Delete: resourceAciPhysicalDomainDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciPhysicalDomainImport,
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
		}),
	}
}
func getRemotePhysicalDomain(client *client.Client, dn string) (*models.PhysicalDomain, error) {
	physDomPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	physDomP := models.PhysicalDomainFromContainer(physDomPCont)

	if physDomP.DistinguishedName == "" {
		return nil, fmt.Errorf("PhysicalDomain %s not found", physDomP.DistinguishedName)
	}

	return physDomP, nil
}

func setPhysicalDomainAttributes(physDomP *models.PhysicalDomain, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(physDomP.DistinguishedName)
	d.Set("description", physDomP.Description)
	physDomPMap, _ := physDomP.ToMap()

	d.Set("name", physDomPMap["name"])

	d.Set("annotation", physDomPMap["annotation"])
	d.Set("name_alias", physDomPMap["nameAlias"])
	return d
}

func resourceAciPhysicalDomainImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	physDomP, err := getRemotePhysicalDomain(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setPhysicalDomainAttributes(physDomP, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciPhysicalDomainCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] PhysicalDomain: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	physDomPAttr := models.PhysicalDomainAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		physDomPAttr.Annotation = Annotation.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		physDomPAttr.NameAlias = NameAlias.(string)
	}
	physDomP := models.NewPhysicalDomain(fmt.Sprintf("phys-%s", name), "uni", desc, physDomPAttr)

	err := aciClient.Save(physDomP)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if relationToinfraRsVlanNs, ok := d.GetOk("relation_infra_rs_vlan_ns"); ok {
		relationParam := relationToinfraRsVlanNs.(string)
		err = aciClient.CreateRelationinfraRsVlanNsFromPhysicalDomain(physDomP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_vlan_ns")
		d.Partial(false)

	}
	if relationToinfraRsVlanNsDef, ok := d.GetOk("relation_infra_rs_vlan_ns_def"); ok {
		relationParam := relationToinfraRsVlanNsDef.(string)
		err = aciClient.CreateRelationinfraRsVlanNsDefFromPhysicalDomain(physDomP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_vlan_ns_def")
		d.Partial(false)

	}
	if relationToinfraRsVipAddrNs, ok := d.GetOk("relation_infra_rs_vip_addr_ns"); ok {
		relationParam := relationToinfraRsVipAddrNs.(string)
		err = aciClient.CreateRelationinfraRsVipAddrNsFromPhysicalDomain(physDomP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_vip_addr_ns")
		d.Partial(false)

	}
	if relationToinfraRsDomVxlanNsDef, ok := d.GetOk("relation_infra_rs_dom_vxlan_ns_def"); ok {
		relationParam := relationToinfraRsDomVxlanNsDef.(string)
		err = aciClient.CreateRelationinfraRsDomVxlanNsDefFromPhysicalDomain(physDomP.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_dom_vxlan_ns_def")
		d.Partial(false)

	}

	d.SetId(physDomP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciPhysicalDomainRead(d, m)
}

func resourceAciPhysicalDomainUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] PhysicalDomain: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	physDomPAttr := models.PhysicalDomainAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		physDomPAttr.Annotation = Annotation.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		physDomPAttr.NameAlias = NameAlias.(string)
	}
	physDomP := models.NewPhysicalDomain(fmt.Sprintf("phys-%s", name), "uni", desc, physDomPAttr)

	physDomP.Status = "modified"

	err := aciClient.Save(physDomP)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if d.HasChange("relation_infra_rs_vlan_ns") {
		_, newRelParam := d.GetChange("relation_infra_rs_vlan_ns")
		err = aciClient.DeleteRelationinfraRsVlanNsFromPhysicalDomain(physDomP.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationinfraRsVlanNsFromPhysicalDomain(physDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_vlan_ns")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_vlan_ns_def") {
		_, newRelParam := d.GetChange("relation_infra_rs_vlan_ns_def")
		err = aciClient.CreateRelationinfraRsVlanNsDefFromPhysicalDomain(physDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_vlan_ns_def")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_vip_addr_ns") {
		_, newRelParam := d.GetChange("relation_infra_rs_vip_addr_ns")
		err = aciClient.DeleteRelationinfraRsVipAddrNsFromPhysicalDomain(physDomP.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationinfraRsVipAddrNsFromPhysicalDomain(physDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_vip_addr_ns")
		d.Partial(false)

	}
	if d.HasChange("relation_infra_rs_dom_vxlan_ns_def") {
		_, newRelParam := d.GetChange("relation_infra_rs_dom_vxlan_ns_def")
		err = aciClient.CreateRelationinfraRsDomVxlanNsDefFromPhysicalDomain(physDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_dom_vxlan_ns_def")
		d.Partial(false)

	}

	d.SetId(physDomP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciPhysicalDomainRead(d, m)

}

func resourceAciPhysicalDomainRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	physDomP, err := getRemotePhysicalDomain(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setPhysicalDomainAttributes(physDomP, d)

	infraRsVlanNsData, err := aciClient.ReadRelationinfraRsVlanNsFromPhysicalDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsVlanNs %v", err)

	} else {
		if _, ok := d.GetOk("relation_infra_rs_vlan_ns"); ok {
			tfName := d.Get("relation_infra_rs_vlan_ns").(string)
			if tfName != infraRsVlanNsData {
				d.Set("relation_infra_rs_vlan_ns", "")
			}
		}
	}

	infraRsVlanNsDefData, err := aciClient.ReadRelationinfraRsVlanNsDefFromPhysicalDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsVlanNsDef %v", err)

	} else {
		if _, ok := d.GetOk("relation_infra_rs_vlan_ns_def"); ok {
			tfName := d.Get("relation_infra_rs_vlan_ns_def").(string)
			if tfName != infraRsVlanNsDefData {
				d.Set("relation_infra_rs_vlan_ns_def", "")
			}
		}
	}

	infraRsVipAddrNsData, err := aciClient.ReadRelationinfraRsVipAddrNsFromPhysicalDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsVipAddrNs %v", err)

	} else {
		if _, ok := d.GetOk("relation_infra_rs_vip_addr_ns"); ok {
			tfName := d.Get("relation_infra_rs_vip_addr_ns").(string)
			if tfName != infraRsVipAddrNsData {
				d.Set("relation_infra_rs_vip_addr_ns", "")
			}
		}
	}

	infraRsDomVxlanNsDefData, err := aciClient.ReadRelationinfraRsDomVxlanNsDefFromPhysicalDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsDomVxlanNsDef %v", err)

	} else {
		if _, ok := d.GetOk("relation_infra_rs_dom_vxlan_ns_def"); ok {
			tfName := d.Get("relation_infra_rs_dom_vxlan_ns_def").(string)
			if tfName != infraRsDomVxlanNsDefData {
				d.Set("relation_infra_rs_dom_vxlan_ns_def", "")
			}
		}
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciPhysicalDomainDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "physDomP")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
