package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAciPhysicalDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciPhysicalDomainCreate,
		UpdateContext: resourceAciPhysicalDomainUpdate,
		ReadContext:   resourceAciPhysicalDomainRead,
		DeleteContext: resourceAciPhysicalDomainDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciPhysicalDomainImport,
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "orchestrator:terraform",
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"relation_infra_rs_vip_addr_ns": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_infra_rs_dom_vxlan_ns_def": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
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

func setPhysicalDomainAttributes(physDomP *models.PhysicalDomain, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(physDomP.DistinguishedName)
	//d.Set("description", physDomP.Description)
	physDomPMap, err := physDomP.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("name", physDomPMap["name"])

	d.Set("annotation", physDomPMap["annotation"])
	d.Set("name_alias", physDomPMap["nameAlias"])
	return d, nil
}

func getAndSetRelationinfraRsVlanNsFromPhysicalDomain(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsVlanNsData, err := client.ReadRelationinfraRsVlanNsFromPhysicalDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsVlanNs %v", err)
		d.Set("relation_infra_rs_vlan_ns", "")
		return d, err
	} else {
		d.Set("relation_infra_rs_vlan_ns", infraRsVlanNsData.(string))
	}
	return d, nil
}

func getAndSetRelationinfraRsVlanNsDefFromPhysicalDomain(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsVlanNsDefData, err := client.ReadRelationinfraRsVlanNsDefFromPhysicalDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsVlanNsDef %v", err)
		d.Set("relation_infra_rs_vlan_ns_def", "")
		return d, err
	} else {
		d.Set("relation_infra_rs_vlan_ns_def", infraRsVlanNsDefData.(string))
	}
	return d, nil
}

func getAndSetRelationinfraRsVipAddrNsFromPhysicalDomain(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsVipAddrNsData, err := client.ReadRelationinfraRsVipAddrNsFromPhysicalDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsVipAddrNs %v", err)
		d.Set("relation_infra_rs_vip_addr_ns", "")
		return d, err
	} else {
		d.Set("relation_infra_rs_vip_addr_ns", infraRsVipAddrNsData.(string))
	}
	return d, nil
}

func getAndSetRelationinfraRsDomVxlanNsDefFromPhysicalDomain(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	infraRsDomVxlanNsDefData, err := client.ReadRelationinfraRsDomVxlanNsDefFromPhysicalDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsDomVxlanNsDef %v", err)
		d.Set("relation_infra_rs_dom_vxlan_ns_def", "")
		return d, err
	} else {
		d.Set("relation_infra_rs_dom_vxlan_ns_def", infraRsDomVxlanNsDefData.(string))
	}
	return d, nil
}

func resourceAciPhysicalDomainImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	physDomP, err := getRemotePhysicalDomain(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setPhysicalDomainAttributes(physDomP, d)
	if err != nil {
		return nil, err
	}

	// infraRsVlanNs - Beginning Import
	log.Printf("[DEBUG] %s: infraRsVlanNs - Beginning Import with parent DN", dn)
	_, err = getAndSetRelationinfraRsVlanNsFromPhysicalDomain(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsVlanNs - Import finished successfully", d.Get("relation_infra_rs_vlan_ns"))
	}
	// infraRsVlanNs - Import finished successfully

	// infraRsVlanNsDef - Beginning Import
	log.Printf("[DEBUG] %s: infraRsVlanNsDef - Beginning Import with parent DN", dn)
	_, err = getAndSetRelationinfraRsVlanNsDefFromPhysicalDomain(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsVlanNsDef - Import finished successfully", d.Get("relation_infra_rs_vlan_ns_def"))
	}
	// infraRsVlanNsDef - Import finished successfully

	// infraRsVipAddrNs - Beginning Import
	log.Printf("[DEBUG] %s: infraRsVipAddrNs - Beginning Import with parent DN", dn)
	_, err = getAndSetRelationinfraRsVipAddrNsFromPhysicalDomain(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsVipAddrNs - Import finished successfully", d.Get("relation_infra_rs_vip_addr_ns"))
	}
	// infraRsVipAddrNs - Import finished successfully

	// infraRsDomVxlanNsDef - Beginning Import
	log.Printf("[DEBUG] %s: infraRsDomVxlanNsDef - Beginning Import with parent DN", dn)
	_, err = getAndSetRelationinfraRsDomVxlanNsDefFromPhysicalDomain(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsDomVxlanNsDef - Import finished successfully", d.Get("relation_infra_rs_dom_vxlan_ns_def"))
	}
	// infraRsDomVxlanNsDef - Import finished successfully

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciPhysicalDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] PhysicalDomain: Beginning Creation")
	aciClient := m.(*client.Client)
	//desc := d.Get("description").(string)

	name := d.Get("name").(string)

	physDomPAttr := models.PhysicalDomainAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		physDomPAttr.Annotation = Annotation.(string)
	} else {
		physDomPAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		physDomPAttr.NameAlias = NameAlias.(string)
	}
	physDomP := models.NewPhysicalDomain(fmt.Sprintf("phys-%s", name), "uni", "", physDomPAttr)

	err := aciClient.Save(physDomP)
	if err != nil {
		return diag.FromErr(err)
	}
	checkDns := make([]string, 0, 1)

	if relationToinfraRsVlanNs, ok := d.GetOk("relation_infra_rs_vlan_ns"); ok {
		relationParam := relationToinfraRsVlanNs.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationToinfraRsVlanNsDef, ok := d.GetOk("relation_infra_rs_vlan_ns_def"); ok {
		relationParam := relationToinfraRsVlanNsDef.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationToinfraRsVipAddrNs, ok := d.GetOk("relation_infra_rs_vip_addr_ns"); ok {
		relationParam := relationToinfraRsVipAddrNs.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationToinfraRsDomVxlanNsDef, ok := d.GetOk("relation_infra_rs_dom_vxlan_ns_def"); ok {
		relationParam := relationToinfraRsDomVxlanNsDef.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationToinfraRsVlanNs, ok := d.GetOk("relation_infra_rs_vlan_ns"); ok {
		relationParam := relationToinfraRsVlanNs.(string)
		err = aciClient.CreateRelationinfraRsVlanNsFromPhysicalDomain(physDomP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsVlanNsDef, ok := d.GetOk("relation_infra_rs_vlan_ns_def"); ok {
		relationParam := relationToinfraRsVlanNsDef.(string)
		err = aciClient.CreateRelationinfraRsVlanNsDefFromPhysicalDomain(physDomP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsVipAddrNs, ok := d.GetOk("relation_infra_rs_vip_addr_ns"); ok {
		relationParam := relationToinfraRsVipAddrNs.(string)
		err = aciClient.CreateRelationinfraRsVipAddrNsFromPhysicalDomain(physDomP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if relationToinfraRsDomVxlanNsDef, ok := d.GetOk("relation_infra_rs_dom_vxlan_ns_def"); ok {
		relationParam := relationToinfraRsDomVxlanNsDef.(string)
		err = aciClient.CreateRelationinfraRsDomVxlanNsDefFromPhysicalDomain(physDomP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(physDomP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciPhysicalDomainRead(ctx, d, m)
}

func resourceAciPhysicalDomainUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] PhysicalDomain: Beginning Update")

	aciClient := m.(*client.Client)
	//desc := d.Get("description").(string)
	name := d.Get("name").(string)

	physDomPAttr := models.PhysicalDomainAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		physDomPAttr.Annotation = Annotation.(string)
	} else {
		physDomPAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		physDomPAttr.NameAlias = NameAlias.(string)
	}
	physDomP := models.NewPhysicalDomain(fmt.Sprintf("phys-%s", name), "uni", "", physDomPAttr)

	physDomP.Status = "modified"

	err := aciClient.Save(physDomP)

	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_infra_rs_vlan_ns") {
		_, newRelParam := d.GetChange("relation_infra_rs_vlan_ns")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_infra_rs_vlan_ns_def") {
		_, newRelParam := d.GetChange("relation_infra_rs_vlan_ns_def")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_infra_rs_vip_addr_ns") {
		_, newRelParam := d.GetChange("relation_infra_rs_vip_addr_ns")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_infra_rs_dom_vxlan_ns_def") {
		_, newRelParam := d.GetChange("relation_infra_rs_dom_vxlan_ns_def")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_infra_rs_vlan_ns") {
		_, newRelParam := d.GetChange("relation_infra_rs_vlan_ns")
		err = aciClient.DeleteRelationinfraRsVlanNsFromPhysicalDomain(physDomP.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsVlanNsFromPhysicalDomain(physDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_vlan_ns_def") {
		_, newRelParam := d.GetChange("relation_infra_rs_vlan_ns_def")
		err = aciClient.CreateRelationinfraRsVlanNsDefFromPhysicalDomain(physDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_vip_addr_ns") {
		_, newRelParam := d.GetChange("relation_infra_rs_vip_addr_ns")
		err = aciClient.DeleteRelationinfraRsVipAddrNsFromPhysicalDomain(physDomP.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsVipAddrNsFromPhysicalDomain(physDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_dom_vxlan_ns_def") {
		_, newRelParam := d.GetChange("relation_infra_rs_dom_vxlan_ns_def")
		err = aciClient.CreateRelationinfraRsDomVxlanNsDefFromPhysicalDomain(physDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(physDomP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciPhysicalDomainRead(ctx, d, m)

}

func resourceAciPhysicalDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	physDomP, err := getRemotePhysicalDomain(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setPhysicalDomainAttributes(physDomP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	// infraRsVlanNs - Beginning Read
	log.Printf("[DEBUG] %s: infraRsVlanNs - Beginning Read with parent DN", dn)
	_, err = getAndSetRelationinfraRsVlanNsFromPhysicalDomain(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsVlanNs - Read finished successfully", d.Get("relation_infra_rs_vlan_ns"))
	}
	// infraRsVlanNs - Read finished successfully

	// infraRsVlanNsDef - Beginning Read
	log.Printf("[DEBUG] %s: infraRsVlanNsDef - Beginning Read with parent DN", dn)
	_, err = getAndSetRelationinfraRsVlanNsDefFromPhysicalDomain(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsVlanNsDef - Read finished successfully", d.Get("relation_infra_rs_vlan_ns_def"))
	}
	// infraRsVlanNsDef - Read finished successfully

	// infraRsVipAddrNs - Beginning Read
	log.Printf("[DEBUG] %s: infraRsVipAddrNs - Beginning Read with parent DN", dn)
	_, err = getAndSetRelationinfraRsVipAddrNsFromPhysicalDomain(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsVipAddrNs - Read finished successfully", d.Get("relation_infra_rs_vip_addr_ns"))
	}
	// infraRsVipAddrNs - Read finished successfully

	// infraRsDomVxlanNsDef - Beginning Read
	log.Printf("[DEBUG] %s: infraRsDomVxlanNsDef - Beginning Read with parent DN", dn)
	_, err = getAndSetRelationinfraRsDomVxlanNsDefFromPhysicalDomain(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsDomVxlanNsDef - Read finished successfully", d.Get("relation_infra_rs_dom_vxlan_ns_def"))
	}
	// infraRsDomVxlanNsDef - Read finished successfully

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciPhysicalDomainDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "physDomP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
