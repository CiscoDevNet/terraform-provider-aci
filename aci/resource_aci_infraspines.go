package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciSwitchSpineAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciSwitchSpineAssociationCreate,
		Update: resourceAciSwitchSpineAssociationUpdate,
		Read:   resourceAciSwitchSpineAssociationRead,
		Delete: resourceAciSwitchSpineAssociationDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciSwitchSpineAssociationImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"spine_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"spine_switch_association_type": &schema.Schema{
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

			"relation_infra_rs_spine_acc_node_p_grp": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		}),
	}
}
func getRemoteSwitchSpineAssociation(client *client.Client, dn string) (*models.SwitchSpineAssociation, error) {
	infraSpineSCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	infraSpineS := models.SwitchSpineAssociationFromContainer(infraSpineSCont)

	if infraSpineS.DistinguishedName == "" {
		return nil, fmt.Errorf("SwitchAssociation %s not found", infraSpineS.DistinguishedName)
	}

	return infraSpineS, nil
}

func setSwitchSpineAssociationAttributes(infraSpineS *models.SwitchSpineAssociation, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(infraSpineS.DistinguishedName)
	d.Set("description", infraSpineS.Description)
	if dn != infraSpineS.DistinguishedName {
		d.Set("spine_profile_dn", "")
	}
	infraSpineSMap, _ := infraSpineS.ToMap()

	d.Set("name", infraSpineSMap["name"])
	d.Set("spine_switch_association_type", infraSpineSMap["type"])
	d.Set("annotation", infraSpineSMap["annotation"])
	d.Set("name_alias", infraSpineSMap["nameAlias"])

	return d
}

func resourceAciSwitchSpineAssociationImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	infraSpineS, err := getRemoteSwitchSpineAssociation(aciClient, dn)

	if err != nil {
		return nil, err
	}
	infraSpineSMap, _ := infraSpineS.ToMap()
	name := infraSpineSMap["name"]
	satype := infraSpineSMap["type"]
	pDN := GetParentDn(dn, fmt.Sprintf("/spines-%s-typ-%s", name, satype))
	d.Set("spine_profile_dn", pDN)
	schemaFilled := setSwitchSpineAssociationAttributes(infraSpineS, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciSwitchSpineAssociationCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] SwitchAssociation: Beginning Creation")
	aciClient := m.(*client.Client)

	desc := d.Get("description").(string)
	name := d.Get("name").(string)

	infraSpineSAttr := models.SwitchSpineAssociationAttributes{}
	switchAssociationType := d.Get("spine_switch_association_type").(string)
	infraSpineSAttr.SwitchAssociationType = switchAssociationType

	SpineProfileDn := d.Get("spine_profile_dn").(string)

	if Annotation, ok := d.GetOk("annotation"); ok {
		infraSpineSAttr.Annotation = Annotation.(string)
	} else {
		infraSpineSAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraSpineSAttr.NameAlias = NameAlias.(string)
	}

	infraSpineS := models.NewSwitchSpineAssociation(fmt.Sprintf("spines-%s-typ-%s", name, switchAssociationType), SpineProfileDn, desc, infraSpineSAttr)

	err := aciClient.Save(infraSpineS)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.Partial(false)

	checkDns := make([]string, 0, 1)

	if relationToinfraRsSpineAccNodePGrp, ok := d.GetOk("relation_infra_rs_spine_acc_node_p_grp"); ok {
		relationParam := relationToinfraRsSpineAccNodePGrp.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return err
	}
	d.Partial(false)

	if relationToinfraRsSpineAccNodePGrp, ok := d.GetOk("relation_infra_rs_spine_acc_node_p_grp"); ok {
		relationParam := relationToinfraRsSpineAccNodePGrp.(string)
		err = aciClient.CreateRelationinfraRsSpineAccNodePGrpFromSwitchAssociation(infraSpineS.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.Partial(false)

	}

	d.SetId(infraSpineS.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciSwitchSpineAssociationRead(d, m)
}

func resourceAciSwitchSpineAssociationUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] SwitchAssociation: Beginning Update")

	aciClient := m.(*client.Client)

	desc := d.Get("description").(string)
	name := d.Get("name").(string)

	infraSpineSAttr := models.SwitchSpineAssociationAttributes{}

	switchAssociationType := d.Get("spine_switch_association_type").(string)
	infraSpineSAttr.SwitchAssociationType = switchAssociationType

	SpineProfileDn := d.Get("spine_profile_dn").(string)

	if Annotation, ok := d.GetOk("annotation"); ok {
		infraSpineSAttr.Annotation = Annotation.(string)
	} else {
		infraSpineSAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraSpineSAttr.NameAlias = NameAlias.(string)
	}

	infraSpineS := models.NewSwitchSpineAssociation(fmt.Sprintf("spines-%s-typ-%s", name, switchAssociationType), SpineProfileDn, desc, infraSpineSAttr)

	infraSpineS.Status = "modified"

	err := aciClient.Save(infraSpineS)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.Partial(false)

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_infra_rs_spine_acc_node_p_grp") {
		_, newRelParam := d.GetChange("relation_infra_rs_spine_acc_node_p_grp")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return err
	}
	d.Partial(false)

	if d.HasChange("relation_infra_rs_spine_acc_node_p_grp") {
		_, newRelParam := d.GetChange("relation_infra_rs_spine_acc_node_p_grp")
		err = aciClient.DeleteRelationinfraRsSpineAccNodePGrpFromSwitchAssociation(infraSpineS.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationinfraRsSpineAccNodePGrpFromSwitchAssociation(infraSpineS.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.Partial(false)

	}

	d.SetId(infraSpineS.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciSwitchSpineAssociationRead(d, m)

}

func resourceAciSwitchSpineAssociationRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	infraSpineS, err := getRemoteSwitchSpineAssociation(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setSwitchSpineAssociationAttributes(infraSpineS, d)

	infraRsSpineAccNodePGrpData, err := aciClient.ReadRelationinfraRsSpineAccNodePGrpFromSwitchAssociation(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsSpineAccNodePGrp %v", err)
		d.Set("relation_infra_rs_spine_acc_node_p_grp", "")

	} else {
		d.Set("relation_infra_rs_spine_acc_node_p_grp", infraRsSpineAccNodePGrpData)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciSwitchSpineAssociationDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "infraSpineS")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
