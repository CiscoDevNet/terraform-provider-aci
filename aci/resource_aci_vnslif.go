package aci

import (
	"context"
	"fmt"
	"log"
	"sort"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAciLogicalInterface() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciLogicalInterfaceCreate,
		UpdateContext: resourceAciLogicalInterfaceUpdate,
		ReadContext:   resourceAciLogicalInterfaceRead,
		DeleteContext: resourceAciLogicalInterfaceDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciLogicalInterfaceImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"l4_l7_device_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"encap": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enhanced_lag_policy_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"relation_vns_rs_c_if_att_n": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Create relation to vnsCIf",
				Set:         schema.HashString,
			}})),
	}
}

func getRemoteLogicalInterface(client *client.Client, dn string) (*models.LogicalInterface, error) {
	vnsLIfCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	vnsLIf := models.LogicalInterfaceFromContainer(vnsLIfCont)
	if vnsLIf.DistinguishedName == "" {
		return nil, fmt.Errorf("Logical Interface %s not found", dn)
	}
	return vnsLIf, nil
}

func setLogicalInterfaceAttributes(vnsLIf *models.LogicalInterface, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(vnsLIf.DistinguishedName)
	vnsLIfMap, err := vnsLIf.ToMap()
	if err != nil {
		return d, err
	}
	dn := d.Id()
	if dn != vnsLIf.DistinguishedName {
		d.Set("l4_l7_device_dn", "")
	} else {
		d.Set("l4_l7_device_dn", GetParentDn(vnsLIf.DistinguishedName, fmt.Sprintf("/"+models.RnvnsLIf, vnsLIfMap["name"])))
	}
	d.Set("annotation", vnsLIfMap["annotation"])
	d.Set("encap", vnsLIfMap["encap"])
	d.Set("enhanced_lag_policy_name", vnsLIfMap["lagPolicyName"])
	d.Set("name", vnsLIfMap["name"])
	d.Set("name_alias", vnsLIfMap["nameAlias"])
	return d, nil
}

func resourceAciLogicalInterfaceImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	vnsLIf, err := getRemoteLogicalInterface(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setLogicalInterfaceAttributes(vnsLIf, d)
	if err != nil {
		return nil, err
	}
	vnsRsCIfAttNData, err := aciClient.ReadRelationvnsRsCIfAttN(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsCIfAttN %v", err)
		d.Set("relation_vns_rs_c_if_att_n", make([]string, 0, 1))
	} else {
		vnsRsCIfAttNDataList := toStringList(vnsRsCIfAttNData.(*schema.Set).List())
		sort.Strings(vnsRsCIfAttNDataList)
		d.Set("relation_vns_rs_c_if_att_n", vnsRsCIfAttNDataList)
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciLogicalInterfaceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LogicalInterface: Beginning Creation")
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	L4_L7DevicesDn := d.Get("l4_l7_device_dn").(string)

	vnsLIfAttr := models.LogicalInterfaceAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		vnsLIfAttr.Annotation = Annotation.(string)
	} else {
		vnsLIfAttr.Annotation = "{}"
	}

	if Encap, ok := d.GetOk("encap"); ok {
		vnsLIfAttr.Encap = Encap.(string)
	}

	if LagPolicyName, ok := d.GetOk("enhanced_lag_policy_name"); ok {
		vnsLIfAttr.LagPolicyName = LagPolicyName.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		vnsLIfAttr.Name = Name.(string)
	}
	vnsLIf := models.NewLogicalInterface(fmt.Sprintf(models.RnvnsLIf, name), L4_L7DevicesDn, nameAlias, vnsLIfAttr)

	err := aciClient.Save(vnsLIf)
	if err != nil {
		return diag.FromErr(err)
	}
	checkDns := make([]string, 0, 1)

	if relationTovnsRsCIfAttN, ok := d.GetOk("relation_vns_rs_c_if_att_n"); ok {
		relationParamList := toStringList(relationTovnsRsCIfAttN.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTovnsRsCIfAttN, ok := d.GetOk("relation_vns_rs_c_if_att_n"); ok {
		relationParamList := toStringList(relationTovnsRsCIfAttN.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationvnsRsCIfAttN(vnsLIf.DistinguishedName, vnsLIfAttr.Annotation, relationParam)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId(vnsLIf.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciLogicalInterfaceRead(ctx, d, m)
}
func resourceAciLogicalInterfaceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LogicalInterface: Beginning Update")
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	L4_L7DevicesDn := d.Get("l4_l7_device_dn").(string)

	vnsLIfAttr := models.LogicalInterfaceAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		vnsLIfAttr.Annotation = Annotation.(string)
	} else {
		vnsLIfAttr.Annotation = "{}"
	}

	if Encap, ok := d.GetOk("encap"); ok {
		vnsLIfAttr.Encap = Encap.(string)
	}

	if LagPolicyName, ok := d.GetOk("enhanced_lag_policy_name"); ok {
		vnsLIfAttr.LagPolicyName = LagPolicyName.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		vnsLIfAttr.Name = Name.(string)
	}
	vnsLIf := models.NewLogicalInterface(fmt.Sprintf(models.RnvnsLIf, name), L4_L7DevicesDn, nameAlias, vnsLIfAttr)

	vnsLIf.Status = "modified"

	err := aciClient.Save(vnsLIf)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_vns_rs_c_if_att_n") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_vns_rs_c_if_att_n")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())
		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_vns_rs_c_if_att_n") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_vns_rs_c_if_att_n")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationvnsRsCIfAttN(vnsLIf.DistinguishedName, relDn)

			if err != nil {
				return diag.FromErr(err)
			}
		}
		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationvnsRsCIfAttN(vnsLIf.DistinguishedName, vnsLIfAttr.Annotation, relDn)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId(vnsLIf.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciLogicalInterfaceRead(ctx, d, m)
}

func resourceAciLogicalInterfaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	vnsLIf, err := getRemoteLogicalInterface(aciClient, dn)
	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}

	_, err = setLogicalInterfaceAttributes(vnsLIf, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	vnsRsCIfAttNData, err := aciClient.ReadRelationvnsRsCIfAttN(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsCIfAttN %v", err)
		setRelationAttribute(d, "relation_vns_rs_c_if_att_n", make([]string, 0, 1))
	} else {
		vnsRsCIfAttNDataList := toStringList(vnsRsCIfAttNData.(*schema.Set).List())
		sort.Strings(vnsRsCIfAttNDataList)
		setRelationAttribute(d, "relation_vns_rs_c_if_att_n", vnsRsCIfAttNDataList)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciLogicalInterfaceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "vnsLIf")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
