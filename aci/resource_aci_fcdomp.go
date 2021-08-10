package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAciFCDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciFCDomainCreate,
		UpdateContext: resourceAciFCDomainUpdate,
		ReadContext:   resourceAciFCDomainRead,
		DeleteContext: resourceAciFCDomainDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciFCDomainImport,
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "orchestrator:terraform",
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
				Computed: true,
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
				Computed: true,
			},
			"relation_fc_rs_vsan_ns_def": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
				Computed: true,
			},
		},
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

func setFCDomainAttributes(fcDomP *models.FCDomain, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(fcDomP.DistinguishedName)

	fcDomPMap, err := fcDomP.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("name", fcDomPMap["name"])

	d.Set("annotation", fcDomPMap["annotation"])
	d.Set("name_alias", fcDomPMap["nameAlias"])
	return d, nil
}

func resourceAciFCDomainImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	fcDomP, err := getRemoteFCDomain(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setFCDomainAttributes(fcDomP, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciFCDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] FCDomain: Beginning Creation")
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	fcDomPAttr := models.FCDomainAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fcDomPAttr.Annotation = Annotation.(string)
	} else {
		fcDomPAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fcDomPAttr.NameAlias = NameAlias.(string)
	}
	fcDomP := models.NewFCDomain(fmt.Sprintf("fc-%s", name), "uni", fcDomPAttr)

	err := aciClient.Save(fcDomP)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationToinfraRsVlanNs, ok := d.GetOk("relation_infra_rs_vlan_ns"); ok {
		relationParam := relationToinfraRsVlanNs.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTofcRsVsanNs, ok := d.GetOk("relation_fc_rs_vsan_ns"); ok {
		relationParam := relationTofcRsVsanNs.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTofcRsVsanAttr, ok := d.GetOk("relation_fc_rs_vsan_attr"); ok {
		relationParam := relationTofcRsVsanAttr.(string)
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

	if relationTofcRsVsanAttrDef, ok := d.GetOk("relation_fc_rs_vsan_attr_def"); ok {
		relationParam := relationTofcRsVsanAttrDef.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTofcRsVsanNsDef, ok := d.GetOk("relation_fc_rs_vsan_ns_def"); ok {
		relationParam := relationTofcRsVsanNsDef.(string)
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
		err = aciClient.CreateRelationinfraRsVlanNsFromFCDomain(fcDomP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTofcRsVsanNs, ok := d.GetOk("relation_fc_rs_vsan_ns"); ok {
		relationParam := relationTofcRsVsanNs.(string)
		err = aciClient.CreateRelationfcRsVsanNsFromFCDomain(fcDomP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTofcRsVsanAttr, ok := d.GetOk("relation_fc_rs_vsan_attr"); ok {
		relationParam := relationTofcRsVsanAttr.(string)
		err = aciClient.CreateRelationfcRsVsanAttrFromFCDomain(fcDomP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsVlanNsDef, ok := d.GetOk("relation_infra_rs_vlan_ns_def"); ok {
		relationParam := relationToinfraRsVlanNsDef.(string)
		err = aciClient.CreateRelationinfraRsVlanNsDefFromFCDomain(fcDomP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsVipAddrNs, ok := d.GetOk("relation_infra_rs_vip_addr_ns"); ok {
		relationParam := relationToinfraRsVipAddrNs.(string)
		err = aciClient.CreateRelationinfraRsVipAddrNsFromFCDomain(fcDomP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsDomVxlanNsDef, ok := d.GetOk("relation_infra_rs_dom_vxlan_ns_def"); ok {
		relationParam := relationToinfraRsDomVxlanNsDef.(string)
		err = aciClient.CreateRelationinfraRsDomVxlanNsDefFromFCDomain(fcDomP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTofcRsVsanAttrDef, ok := d.GetOk("relation_fc_rs_vsan_attr_def"); ok {
		relationParam := relationTofcRsVsanAttrDef.(string)
		err = aciClient.CreateRelationfcRsVsanAttrDefFromFCDomain(fcDomP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTofcRsVsanNsDef, ok := d.GetOk("relation_fc_rs_vsan_ns_def"); ok {
		relationParam := relationTofcRsVsanNsDef.(string)
		err = aciClient.CreateRelationfcRsVsanNsDefFromFCDomain(fcDomP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(fcDomP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciFCDomainRead(ctx, d, m)
}

func resourceAciFCDomainUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] FCDomain: Beginning Update")

	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	fcDomPAttr := models.FCDomainAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fcDomPAttr.Annotation = Annotation.(string)
	} else {
		fcDomPAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fcDomPAttr.NameAlias = NameAlias.(string)
	}
	fcDomP := models.NewFCDomain(fmt.Sprintf("fc-%s", name), "uni", fcDomPAttr)

	fcDomP.Status = "modified"

	err := aciClient.Save(fcDomP)

	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_infra_rs_vlan_ns") {
		_, newRelParam := d.GetChange("relation_infra_rs_vlan_ns")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_fc_rs_vsan_ns") {
		_, newRelParam := d.GetChange("relation_fc_rs_vsan_ns")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_fc_rs_vsan_attr") {
		_, newRelParam := d.GetChange("relation_fc_rs_vsan_attr")
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

	if d.HasChange("relation_fc_rs_vsan_attr_def") {
		_, newRelParam := d.GetChange("relation_fc_rs_vsan_attr_def")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_fc_rs_vsan_ns_def") {
		_, newRelParam := d.GetChange("relation_fc_rs_vsan_ns_def")
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
		err = aciClient.DeleteRelationinfraRsVlanNsFromFCDomain(fcDomP.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsVlanNsFromFCDomain(fcDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_fc_rs_vsan_ns") {
		_, newRelParam := d.GetChange("relation_fc_rs_vsan_ns")
		err = aciClient.DeleteRelationfcRsVsanNsFromFCDomain(fcDomP.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationfcRsVsanNsFromFCDomain(fcDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_fc_rs_vsan_attr") {
		_, newRelParam := d.GetChange("relation_fc_rs_vsan_attr")
		err = aciClient.DeleteRelationfcRsVsanAttrFromFCDomain(fcDomP.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationfcRsVsanAttrFromFCDomain(fcDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_vlan_ns_def") {
		_, newRelParam := d.GetChange("relation_infra_rs_vlan_ns_def")
		err = aciClient.CreateRelationinfraRsVlanNsDefFromFCDomain(fcDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_vip_addr_ns") {
		_, newRelParam := d.GetChange("relation_infra_rs_vip_addr_ns")
		err = aciClient.DeleteRelationinfraRsVipAddrNsFromFCDomain(fcDomP.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsVipAddrNsFromFCDomain(fcDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_dom_vxlan_ns_def") {
		_, newRelParam := d.GetChange("relation_infra_rs_dom_vxlan_ns_def")
		err = aciClient.CreateRelationinfraRsDomVxlanNsDefFromFCDomain(fcDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_fc_rs_vsan_attr_def") {
		_, newRelParam := d.GetChange("relation_fc_rs_vsan_attr_def")
		err = aciClient.CreateRelationfcRsVsanAttrDefFromFCDomain(fcDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_fc_rs_vsan_ns_def") {
		_, newRelParam := d.GetChange("relation_fc_rs_vsan_ns_def")
		err = aciClient.CreateRelationfcRsVsanNsDefFromFCDomain(fcDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(fcDomP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciFCDomainRead(ctx, d, m)

}

func resourceAciFCDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	fcDomP, err := getRemoteFCDomain(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setFCDomainAttributes(fcDomP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	infraRsVlanNsData, err := aciClient.ReadRelationinfraRsVlanNsFromFCDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsVlanNs %v", err)
		d.Set("relation_infra_rs_vlan_ns", "")

	} else {
		d.Set("relation_infra_rs_vlan_ns", infraRsVlanNsData.(string))
	}

	fcRsVsanNsData, err := aciClient.ReadRelationfcRsVsanNsFromFCDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fcRsVsanNs %v", err)
		d.Set("relation_fc_rs_vsan_ns", "")

	} else {
		d.Set("relation_fc_rs_vsan_ns", fcRsVsanNsData.(string))
	}

	fcRsVsanAttrData, err := aciClient.ReadRelationfcRsVsanAttrFromFCDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fcRsVsanAttr %v", err)
		d.Set("relation_fc_rs_vsan_attr", "")

	} else {
		d.Set("relation_fc_rs_vsan_attr", fcRsVsanAttrData.(string))
	}

	infraRsVlanNsDefData, err := aciClient.ReadRelationinfraRsVlanNsDefFromFCDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsVlanNsDef %v", err)
		d.Set("relation_infra_rs_vlan_ns_def", "")

	} else {
		d.Set("relation_infra_rs_vlan_ns_def", infraRsVlanNsDefData.(string))
	}

	infraRsVipAddrNsData, err := aciClient.ReadRelationinfraRsVipAddrNsFromFCDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsVipAddrNs %v", err)
		d.Set("relation_infra_rs_vip_addr_ns", "")

	} else {
		d.Set("relation_infra_rs_vip_addr_ns", infraRsVipAddrNsData.(string))
	}

	infraRsDomVxlanNsDefData, err := aciClient.ReadRelationinfraRsDomVxlanNsDefFromFCDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsDomVxlanNsDef %v", err)
		d.Set("relation_infra_rs_dom_vxlan_ns_def", "")

	} else {
		d.Set("relation_infra_rs_dom_vxlan_ns_def", infraRsDomVxlanNsDefData.(string))
	}

	fcRsVsanAttrDefData, err := aciClient.ReadRelationfcRsVsanAttrDefFromFCDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fcRsVsanAttrDef %v", err)
		d.Set("relation_fc_rs_vsan_attr_def", "")

	} else {
		d.Set("relation_fc_rs_vsan_attr_def", fcRsVsanAttrDefData.(string))
	}

	fcRsVsanNsDefData, err := aciClient.ReadRelationfcRsVsanNsDefFromFCDomain(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fcRsVsanNsDef %v", err)
		d.Set("relation_fc_rs_vsan_ns_def", "")

	} else {
		d.Set("relation_fc_rs_vsan_ns_def", fcRsVsanNsDefData.(string))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciFCDomainDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fcDomP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
