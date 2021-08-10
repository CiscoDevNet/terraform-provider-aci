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

func resourceAciL2Domain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciL2DomainCreate,
		UpdateContext: resourceAciL2DomainUpdate,
		ReadContext:   resourceAciL2DomainRead,
		DeleteContext: resourceAciL2DomainDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciL2DomainImport,
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
				// Default:  "orchestrator:terraform",
				Computed: true,
				DefaultFunc: func() (interface{}, error) {
					return "orchestrator:terraform", nil
				},
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
			"relation_extnw_rs_out": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_infra_rs_dom_vxlan_ns_def": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}
func getRemoteL2Domain(client *client.Client, dn string) (*models.L2Domain, error) {
	l2extDomPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	l2extDomP := models.L2DomainFromContainer(l2extDomPCont)

	if l2extDomP.DistinguishedName == "" {
		return nil, fmt.Errorf("L2Domain %s not found", l2extDomP.DistinguishedName)
	}

	return l2extDomP, nil
}

func setL2DomainAttributes(l2extDomP *models.L2Domain, d *schema.ResourceData) (*schema.ResourceData, error) {

	d.SetId(l2extDomP.DistinguishedName)
	l2extDomPMap, err := l2extDomP.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("name", l2extDomPMap["name"])

	d.Set("annotation", l2extDomPMap["annotation"])
	d.Set("name_alias", l2extDomPMap["nameAlias"])
	return d, nil
}

func resourceAciL2DomainImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	l2extDomP, err := getRemoteL2Domain(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setL2DomainAttributes(l2extDomP, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciL2DomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] L2Domain: Beginning Creation")
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	l2extDomPAttr := models.L2DomainAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l2extDomPAttr.Annotation = Annotation.(string)
	} else {
		l2extDomPAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l2extDomPAttr.NameAlias = NameAlias.(string)
	}
	l2extDomP := models.NewL2Domain(fmt.Sprintf("l2dom-%s", name), "uni", l2extDomPAttr)

	err := aciClient.Save(l2extDomP)
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
	if relationToextnwRsOut, ok := d.GetOk("relation_extnw_rs_out"); ok {
		relationParamList := toStringList(relationToextnwRsOut.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
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
		err = aciClient.CreateRelationinfraRsVlanNsFromL2Domain(l2extDomP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsVlanNsDef, ok := d.GetOk("relation_infra_rs_vlan_ns_def"); ok {
		relationParam := relationToinfraRsVlanNsDef.(string)
		err = aciClient.CreateRelationinfraRsVlanNsDefFromL2Domain(l2extDomP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToinfraRsVipAddrNs, ok := d.GetOk("relation_infra_rs_vip_addr_ns"); ok {
		relationParam := relationToinfraRsVipAddrNs.(string)
		err = aciClient.CreateRelationinfraRsVipAddrNsFromL2Domain(l2extDomP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToextnwRsOut, ok := d.GetOk("relation_extnw_rs_out"); ok {
		relationParamList := toStringList(relationToextnwRsOut.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationextnwRsOutFromL2Domain(l2extDomP.DistinguishedName, relationParam)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if relationToinfraRsDomVxlanNsDef, ok := d.GetOk("relation_infra_rs_dom_vxlan_ns_def"); ok {
		relationParam := relationToinfraRsDomVxlanNsDef.(string)
		err = aciClient.CreateRelationinfraRsDomVxlanNsDefFromL2Domain(l2extDomP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(l2extDomP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciL2DomainRead(ctx, d, m)
}

func resourceAciL2DomainUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] L2Domain: Beginning Update")

	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	l2extDomPAttr := models.L2DomainAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l2extDomPAttr.Annotation = Annotation.(string)
	} else {
		l2extDomPAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l2extDomPAttr.NameAlias = NameAlias.(string)
	}
	l2extDomP := models.NewL2Domain(fmt.Sprintf("l2dom-%s", name), "uni", l2extDomPAttr)

	l2extDomP.Status = "modified"

	err := aciClient.Save(l2extDomP)

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
	if d.HasChange("relation_extnw_rs_out") {
		oldRel, newRel := d.GetChange("relation_extnw_rs_out")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}

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
		err = aciClient.DeleteRelationinfraRsVlanNsFromL2Domain(l2extDomP.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsVlanNsFromL2Domain(l2extDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_vlan_ns_def") {
		_, newRelParam := d.GetChange("relation_infra_rs_vlan_ns_def")
		err = aciClient.CreateRelationinfraRsVlanNsDefFromL2Domain(l2extDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_infra_rs_vip_addr_ns") {
		_, newRelParam := d.GetChange("relation_infra_rs_vip_addr_ns")
		err = aciClient.DeleteRelationinfraRsVipAddrNsFromL2Domain(l2extDomP.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsVipAddrNsFromL2Domain(l2extDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_extnw_rs_out") {
		oldRel, newRel := d.GetChange("relation_extnw_rs_out")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationextnwRsOutFromL2Domain(l2extDomP.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}

		}

	}
	if d.HasChange("relation_infra_rs_dom_vxlan_ns_def") {
		_, newRelParam := d.GetChange("relation_infra_rs_dom_vxlan_ns_def")
		err = aciClient.CreateRelationinfraRsDomVxlanNsDefFromL2Domain(l2extDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(l2extDomP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciL2DomainRead(ctx, d, m)

}

func resourceAciL2DomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	l2extDomP, err := getRemoteL2Domain(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setL2DomainAttributes(l2extDomP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	infraRsVlanNsData, err := aciClient.ReadRelationinfraRsVlanNsFromL2Domain(dn)

	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsVlanNs %v", err)
		d.Set("relation_infra_rs_vlan_ns", "")

	} else {
		d.Set("relation_infra_rs_vlan_ns", infraRsVlanNsData.(string))
	}

	infraRsVlanNsDefData, err := aciClient.ReadRelationinfraRsVlanNsDefFromL2Domain(dn)

	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsVlanNsDef %v", err)
		d.Set("relation_infra_rs_vlan_ns_def", "")

	} else {
		d.Set("relation_infra_rs_vlan_ns_def", infraRsVlanNsDefData.(string))
	}

	infraRsVipAddrNsData, err := aciClient.ReadRelationinfraRsVipAddrNsFromL2Domain(dn)

	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsVipAddrNs %v", err)
		d.Set("relation_infra_rs_vip_addr_ns", "")

	} else {
		d.Set("relation_infra_rs_vip_addr_ns", infraRsVipAddrNsData.(string))
	}

	extnwRsOutData, err := aciClient.ReadRelationextnwRsOutFromL2Domain(dn)

	if err != nil {
		log.Printf("[DEBUG] Error while reading relation extnwRsOut %v", err)
		d.Set("relation_extnw_rs_out", make([]string, 0, 1))

	} else {
		d.Set("relation_extnw_rs_out", toStringList(extnwRsOutData.(*schema.Set).List()))
	}

	infraRsDomVxlanNsDefData, err := aciClient.ReadRelationinfraRsDomVxlanNsDefFromL2Domain(dn)

	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsDomVxlanNsDef %v", err)
		d.Set("relation_infra_rs_dom_vxlan_ns_def", "")

	} else {
		d.Set("relation_infra_rs_dom_vxlan_ns_def", infraRsDomVxlanNsDefData.(string))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciL2DomainDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "l2extDomP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
