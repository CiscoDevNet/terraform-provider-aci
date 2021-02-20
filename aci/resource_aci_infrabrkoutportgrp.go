package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAciLeafBreakoutPortGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciLeafBreakoutPortGroupCreate,
		Update: resourceAciLeafBreakoutPortGroupUpdate,
		Read:   resourceAciLeafBreakoutPortGroupRead,
		Delete: resourceAciLeafBreakoutPortGroupDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciLeafBreakoutPortGroupImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"brkout_map": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"100g-2x",
					"100g-4x",
					"10g-4x",
					"25g-4x",
					"50g-8x",
					"none",
				}, false),
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"relation_infra_rs_mon_brkout_infra_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
		}),
	}
}
func getRemoteLeafBreakoutPortGroup(client *client.Client, dn string) (*models.LeafBreakoutPortGroup, error) {
	infraBrkoutPortGrpCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	infraBrkoutPortGrp := models.LeafBreakoutPortGroupFromContainer(infraBrkoutPortGrpCont)

	if infraBrkoutPortGrp.DistinguishedName == "" {
		return nil, fmt.Errorf("LeafBreakoutPortGroup %s not found", infraBrkoutPortGrp.DistinguishedName)
	}

	return infraBrkoutPortGrp, nil
}

func setLeafBreakoutPortGroupAttributes(infraBrkoutPortGrp *models.LeafBreakoutPortGroup, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(infraBrkoutPortGrp.DistinguishedName)
	d.Set("description", infraBrkoutPortGrp.Description)
	infraBrkoutPortGrpMap, _ := infraBrkoutPortGrp.ToMap()

	d.Set("name", infraBrkoutPortGrpMap["name"])

	d.Set("annotation", infraBrkoutPortGrpMap["annotation"])
	d.Set("brkout_map", infraBrkoutPortGrpMap["brkoutMap"])
	d.Set("name_alias", infraBrkoutPortGrpMap["nameAlias"])
	return d
}

func resourceAciLeafBreakoutPortGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	infraBrkoutPortGrp, err := getRemoteLeafBreakoutPortGroup(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setLeafBreakoutPortGroupAttributes(infraBrkoutPortGrp, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciLeafBreakoutPortGroupCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] LeafBreakoutPortGroup: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	infraBrkoutPortGrpAttr := models.LeafBreakoutPortGroupAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraBrkoutPortGrpAttr.Annotation = Annotation.(string)
	} else {
		infraBrkoutPortGrpAttr.Annotation = "{}"
	}
	if BrkoutMap, ok := d.GetOk("brkout_map"); ok {
		infraBrkoutPortGrpAttr.BrkoutMap = BrkoutMap.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraBrkoutPortGrpAttr.NameAlias = NameAlias.(string)
	}
	infraBrkoutPortGrp := models.NewLeafBreakoutPortGroup(fmt.Sprintf("infra/funcprof/brkoutportgrp-%s", name), "uni", desc, infraBrkoutPortGrpAttr)

	err := aciClient.Save(infraBrkoutPortGrp)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if relationToinfraRsMonBrkoutInfraPol, ok := d.GetOk("relation_infra_rs_mon_brkout_infra_pol"); ok {
		relationParam := relationToinfraRsMonBrkoutInfraPol.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationinfraRsMonBrkoutInfraPolFromLeafBreakoutPortGroup(infraBrkoutPortGrp.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_mon_brkout_infra_pol")
		d.Partial(false)

	}

	d.SetId(infraBrkoutPortGrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciLeafBreakoutPortGroupRead(d, m)
}

func resourceAciLeafBreakoutPortGroupUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] LeafBreakoutPortGroup: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	infraBrkoutPortGrpAttr := models.LeafBreakoutPortGroupAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraBrkoutPortGrpAttr.Annotation = Annotation.(string)
	} else {
		infraBrkoutPortGrpAttr.Annotation = "{}"
	}
	if BrkoutMap, ok := d.GetOk("brkout_map"); ok {
		infraBrkoutPortGrpAttr.BrkoutMap = BrkoutMap.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraBrkoutPortGrpAttr.NameAlias = NameAlias.(string)
	}
	infraBrkoutPortGrp := models.NewLeafBreakoutPortGroup(fmt.Sprintf("infra/funcprof/brkoutportgrp-%s", name), "uni", desc, infraBrkoutPortGrpAttr)

	infraBrkoutPortGrp.Status = "modified"

	err := aciClient.Save(infraBrkoutPortGrp)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if d.HasChange("relation_infra_rs_mon_brkout_infra_pol") {
		_, newRelParam := d.GetChange("relation_infra_rs_mon_brkout_infra_pol")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationinfraRsMonBrkoutInfraPolFromLeafBreakoutPortGroup(infraBrkoutPortGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_mon_brkout_infra_pol")
		d.Partial(false)

	}

	d.SetId(infraBrkoutPortGrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciLeafBreakoutPortGroupRead(d, m)

}

func resourceAciLeafBreakoutPortGroupRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	infraBrkoutPortGrp, err := getRemoteLeafBreakoutPortGroup(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setLeafBreakoutPortGroupAttributes(infraBrkoutPortGrp, d)

	infraRsMonBrkoutInfraPolData, err := aciClient.ReadRelationinfraRsMonBrkoutInfraPolFromLeafBreakoutPortGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsMonBrkoutInfraPol %v", err)
		d.Set("relation_infra_rs_mon_brkout_infra_pol", "")

	} else {
		if relDn, ok := d.GetOk("relation_infra_rs_mon_brkout_infra_pol"); ok {
			tfName := GetMOName(relDn.(string))
			if tfName != infraRsMonBrkoutInfraPolData {
				d.Set("relation_infra_rs_mon_brkout_infra_pol", "")
			}
		}
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciLeafBreakoutPortGroupDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "infraBrkoutPortGrp")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
