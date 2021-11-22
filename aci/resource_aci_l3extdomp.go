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

func resourceAciL3DomainProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciL3DomainProfileCreate,
		UpdateContext: resourceAciL3DomainProfileUpdate,
		ReadContext:   resourceAciL3DomainProfileRead,
		DeleteContext: resourceAciL3DomainProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciL3DomainProfileImport,
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				// Default:  "orchestrator:terraform",
				Computed: true,
				DefaultFunc: func() (interface{}, error) {
					return "orchestrator:terraform", nil
				},
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
				Type: schema.TypeString,

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
				Type: schema.TypeString,

				Optional: true,
			},
		},
	}
}
func getRemoteL3DomainProfile(client *client.Client, dn string) (*models.L3DomainProfile, error) {
	l3extDomPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	l3extDomP := models.L3DomainProfileFromContainer(l3extDomPCont)

	if l3extDomP.DistinguishedName == "" {
		return nil, fmt.Errorf("L3DomainProfile %s not found", l3extDomP.DistinguishedName)
	}

	return l3extDomP, nil
}

func setL3DomainProfileAttributes(l3extDomP *models.L3DomainProfile, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(l3extDomP.DistinguishedName)
	l3extDomPMap, err := l3extDomP.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("name", l3extDomPMap["name"])

	d.Set("annotation", l3extDomPMap["annotation"])
	d.Set("name_alias", l3extDomPMap["nameAlias"])
	return d, nil
}

func resourceAciL3DomainProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	l3extDomP, err := getRemoteL3DomainProfile(aciClient, dn)
	if err != nil {
		return nil, err
	}

	schemaFilled, err := setL3DomainProfileAttributes(l3extDomP, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciL3DomainProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] L3DomainProfile: Beginning Creation")
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	l3extDomPAttr := models.L3DomainProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extDomPAttr.Annotation = Annotation.(string)
	} else {
		l3extDomPAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l3extDomPAttr.NameAlias = NameAlias.(string)
	}
	l3extDomP := models.NewL3DomainProfile(fmt.Sprintf("l3dom-%s", name), "uni", "", l3extDomPAttr)

	err := aciClient.Save(l3extDomP)
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
		err = aciClient.CreateRelationinfraRsVlanNsFromL3DomainProfile(l3extDomP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if relationToinfraRsVlanNsDef, ok := d.GetOk("relation_infra_rs_vlan_ns_def"); ok {
		relationParam := relationToinfraRsVlanNsDef.(string)
		err = aciClient.CreateRelationinfraRsVlanNsDefFromL3DomainProfile(l3extDomP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if relationToinfraRsVipAddrNs, ok := d.GetOk("relation_infra_rs_vip_addr_ns"); ok {
		relationParam := relationToinfraRsVipAddrNs.(string)
		err = aciClient.CreateRelationinfraRsVipAddrNsFromL3DomainProfile(l3extDomP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationToextnwRsOut, ok := d.GetOk("relation_extnw_rs_out"); ok {
		relationParamList := toStringList(relationToextnwRsOut.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationextnwRsOutFromL3DomainProfile(l3extDomP.DistinguishedName, relationParam)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if relationToinfraRsDomVxlanNsDef, ok := d.GetOk("relation_infra_rs_dom_vxlan_ns_def"); ok {
		relationParam := relationToinfraRsDomVxlanNsDef.(string)
		err = aciClient.CreateRelationinfraRsDomVxlanNsDefFromL3DomainProfile(l3extDomP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(l3extDomP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciL3DomainProfileRead(ctx, d, m)
}

func resourceAciL3DomainProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] L3DomainProfile: Beginning Update")

	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	l3extDomPAttr := models.L3DomainProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extDomPAttr.Annotation = Annotation.(string)
	} else {
		l3extDomPAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l3extDomPAttr.NameAlias = NameAlias.(string)
	}
	l3extDomP := models.NewL3DomainProfile(fmt.Sprintf("l3dom-%s", name), "uni", "", l3extDomPAttr)

	l3extDomP.Status = "modified"

	err := aciClient.Save(l3extDomP)

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
		err = aciClient.DeleteRelationinfraRsVlanNsFromL3DomainProfile(l3extDomP.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsVlanNsFromL3DomainProfile(l3extDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("relation_infra_rs_vlan_ns_def") {
		_, newRelParam := d.GetChange("relation_infra_rs_vlan_ns_def")
		err = aciClient.CreateRelationinfraRsVlanNsDefFromL3DomainProfile(l3extDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("relation_infra_rs_vip_addr_ns") {
		_, newRelParam := d.GetChange("relation_infra_rs_vip_addr_ns")
		err = aciClient.DeleteRelationinfraRsVipAddrNsFromL3DomainProfile(l3extDomP.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsVipAddrNsFromL3DomainProfile(l3extDomP.DistinguishedName, newRelParam.(string))
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
			err = aciClient.CreateRelationextnwRsOutFromL3DomainProfile(l3extDomP.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}
		}

	}
	if d.HasChange("relation_infra_rs_dom_vxlan_ns_def") {
		_, newRelParam := d.GetChange("relation_infra_rs_dom_vxlan_ns_def")
		err = aciClient.CreateRelationinfraRsDomVxlanNsDefFromL3DomainProfile(l3extDomP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(l3extDomP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciL3DomainProfileRead(ctx, d, m)
}

func resourceAciL3DomainProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	l3extDomP, err := getRemoteL3DomainProfile(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}

	_, err = setL3DomainProfileAttributes(l3extDomP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	infraRsVlanNsData, err := aciClient.ReadRelationinfraRsVlanNsFromL3DomainProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsVlanNs %v", err)
		d.Set("relation_infra_rs_vlan_ns", "")

	} else {
		d.Set("relation_infra_rs_vlan_ns", infraRsVlanNsData.(string))
	}

	infraRsVlanNsDefData, err := aciClient.ReadRelationinfraRsVlanNsDefFromL3DomainProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsVlanNsDef %v", err)
		d.Set("relation_infra_rs_vlan_ns_def", "")

	} else {
		d.Set("relation_infra_rs_vlan_ns_def", infraRsVlanNsDefData.(string))
	}

	infraRsVipAddrNsData, err := aciClient.ReadRelationinfraRsVipAddrNsFromL3DomainProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsVipAddrNs %v", err)
		d.Set("relation_infra_rs_vip_addr_ns", "")

	} else {
		d.Set("relation_infra_rs_vip_addr_ns", infraRsVipAddrNsData.(string))
	}

	extnwRsOutData, err := aciClient.ReadRelationextnwRsOutFromL3DomainProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation extnwRsOut %v", err)
		d.Set("relation_extnw_rs_out", make([]string, 0, 1))

	} else {
		d.Set("relation_extnw_rs_out", toStringList(extnwRsOutData.(*schema.Set).List()))
	}

	infraRsDomVxlanNsDefData, err := aciClient.ReadRelationinfraRsDomVxlanNsDefFromL3DomainProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsDomVxlanNsDef %v", err)
		d.Set("relation_infra_rs_dom_vxlan_ns_def", "")

	} else {
		d.Set("relation_infra_rs_dom_vxlan_ns_def", infraRsDomVxlanNsDefData.(string))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciL3DomainProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "l3extDomP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
