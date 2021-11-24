package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciSwitchAssociation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciSwitchAssociationCreate,
		UpdateContext: resourceAciSwitchAssociationUpdate,
		ReadContext:   resourceAciSwitchAssociationRead,
		DeleteContext: resourceAciSwitchAssociationDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciSwitchAssociationImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"leaf_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"switch_association_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ALL",
					"range",
					"ALL_IN_POD",
				}, false),
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"relation_infra_rs_acc_node_p_grp": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
		}),
	}
}
func getRemoteSwitchAssociation(client *client.Client, dn string) (*models.SwitchAssociation, error) {
	infraLeafSCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	infraLeafS := models.SwitchAssociationFromContainer(infraLeafSCont)

	if infraLeafS.DistinguishedName == "" {
		return nil, fmt.Errorf("SwitchAssociation %s not found", infraLeafS.DistinguishedName)
	}

	return infraLeafS, nil
}

func setSwitchAssociationAttributes(infraLeafS *models.SwitchAssociation, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(infraLeafS.DistinguishedName)
	d.Set("description", infraLeafS.Description)

	if dn != infraLeafS.DistinguishedName {
		d.Set("leaf_profile_dn", "")
	}
	infraLeafSMap, _ := infraLeafS.ToMap()

	d.Set("name", infraLeafSMap["name"])
	d.Set("switch_association_type", infraLeafSMap["type"])
	d.Set("leaf_profile_dn", GetParentDn(infraLeafS.DistinguishedName, fmt.Sprintf("/leaves-%s-typ-%s", infraLeafSMap["name"], infraLeafSMap["type"])))
	d.Set("annotation", infraLeafSMap["annotation"])
	d.Set("name_alias", infraLeafSMap["nameAlias"])
	return d
}

func resourceAciSwitchAssociationImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	infraLeafS, err := getRemoteSwitchAssociation(aciClient, dn)

	if err != nil {
		return nil, err
	}
	infraLeafSMap, _ := infraLeafS.ToMap()
	name := infraLeafSMap["name"]
	ltype := infraLeafSMap["type"]
	pDN := GetParentDn(dn, fmt.Sprintf("/leaves-%s-typ-%s", name, ltype))
	d.Set("leaf_profile_dn", pDN)
	schemaFilled := setSwitchAssociationAttributes(infraLeafS, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciSwitchAssociationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SwitchAssociation: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	switch_association_type := d.Get("switch_association_type").(string)

	LeafProfileDn := d.Get("leaf_profile_dn").(string)

	infraLeafSAttr := models.SwitchAssociationAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraLeafSAttr.Annotation = Annotation.(string)
	} else {
		infraLeafSAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraLeafSAttr.NameAlias = NameAlias.(string)
	}
	if SwitchAssociation_type, ok := d.GetOk("switch_association_type"); ok {
		infraLeafSAttr.Switch_association_type = SwitchAssociation_type.(string)
	}
	infraLeafS := models.NewSwitchAssociation(fmt.Sprintf("leaves-%s-typ-%s", name, switch_association_type), LeafProfileDn, desc, infraLeafSAttr)

	err := aciClient.Save(infraLeafS)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationToinfraRsAccNodePGrp, ok := d.GetOk("relation_infra_rs_acc_node_p_grp"); ok {
		relationParam := relationToinfraRsAccNodePGrp.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationToinfraRsAccNodePGrp, ok := d.GetOk("relation_infra_rs_acc_node_p_grp"); ok {
		relationParam := relationToinfraRsAccNodePGrp.(string)
		err = aciClient.CreateRelationinfraRsAccNodePGrpFromSwitchAssociation(infraLeafS.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(infraLeafS.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciSwitchAssociationRead(ctx, d, m)
}

func resourceAciSwitchAssociationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SwitchAssociation: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	switch_association_type := d.Get("switch_association_type").(string)

	LeafProfileDn := d.Get("leaf_profile_dn").(string)

	infraLeafSAttr := models.SwitchAssociationAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraLeafSAttr.Annotation = Annotation.(string)
	} else {
		infraLeafSAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraLeafSAttr.NameAlias = NameAlias.(string)
	}
	if SwitchAssociation_type, ok := d.GetOk("switch_association_type"); ok {
		infraLeafSAttr.Switch_association_type = SwitchAssociation_type.(string)
	}
	infraLeafS := models.NewSwitchAssociation(fmt.Sprintf("leaves-%s-typ-%s", name, switch_association_type), LeafProfileDn, desc, infraLeafSAttr)

	infraLeafS.Status = "modified"

	err := aciClient.Save(infraLeafS)

	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_infra_rs_acc_node_p_grp") {
		_, newRelParam := d.GetChange("relation_infra_rs_acc_node_p_grp")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_infra_rs_acc_node_p_grp") {
		_, newRelParam := d.GetChange("relation_infra_rs_acc_node_p_grp")
		err = aciClient.DeleteRelationinfraRsAccNodePGrpFromSwitchAssociation(infraLeafS.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationinfraRsAccNodePGrpFromSwitchAssociation(infraLeafS.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(infraLeafS.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciSwitchAssociationRead(ctx, d, m)

}

func resourceAciSwitchAssociationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	infraLeafS, err := getRemoteSwitchAssociation(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setSwitchAssociationAttributes(infraLeafS, d)

	infraRsAccNodePGrpData, err := aciClient.ReadRelationinfraRsAccNodePGrpFromSwitchAssociation(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsAccNodePGrp %v", err)
		d.Set("relation_infra_rs_acc_node_p_grp", "")

	} else {
		d.Set("relation_infra_rs_acc_node_p_grp", infraRsAccNodePGrpData.(string))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciSwitchAssociationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "infraLeafS")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
