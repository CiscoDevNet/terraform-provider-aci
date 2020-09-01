package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciSpineProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciSpineProfileCreate,
		Update: resourceAciSpineProfileUpdate,
		Read:   resourceAciSpineProfileRead,
		Delete: resourceAciSpineProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciSpineProfileImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

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

			"relation_infra_rs_sp_acc_port_p": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
				Set:      schema.HashString,
			},
		}),
	}
}
func getRemoteSpineProfile(client *client.Client, dn string) (*models.SpineProfile, error) {
	infraSpinePCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	infraSpineP := models.SpineProfileFromContainer(infraSpinePCont)

	if infraSpineP.DistinguishedName == "" {
		return nil, fmt.Errorf("SpineProfile %s not found", infraSpineP.DistinguishedName)
	}

	return infraSpineP, nil
}

func setSpineProfileAttributes(infraSpineP *models.SpineProfile, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(infraSpineP.DistinguishedName)
	d.Set("description", infraSpineP.Description)
	infraSpinePMap, _ := infraSpineP.ToMap()

	d.Set("name", infraSpinePMap["name"])

	d.Set("annotation", infraSpinePMap["annotation"])
	d.Set("name_alias", infraSpinePMap["nameAlias"])
	return d
}

func resourceAciSpineProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	infraSpineP, err := getRemoteSpineProfile(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setSpineProfileAttributes(infraSpineP, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciSpineProfileCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] SpineProfile: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	infraSpinePAttr := models.SpineProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraSpinePAttr.Annotation = Annotation.(string)
	} else {
		infraSpinePAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraSpinePAttr.NameAlias = NameAlias.(string)
	}
	infraSpineP := models.NewSpineProfile(fmt.Sprintf("infra/spprof-%s", name), "uni", desc, infraSpinePAttr)

	err := aciClient.Save(infraSpineP)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	checkDns := make([]string, 0, 1)

	if relationToinfraRsSpAccPortP, ok := d.GetOk("relation_infra_rs_sp_acc_port_p"); ok {
		relationParamList := toStringList(relationToinfraRsSpAccPortP.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return err
	}
	d.Partial(false)

	if relationToinfraRsSpAccPortP, ok := d.GetOk("relation_infra_rs_sp_acc_port_p"); ok {
		relationParamList := toStringList(relationToinfraRsSpAccPortP.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationinfraRsSpAccPortPFromSpineProfile(infraSpineP.DistinguishedName, relationParam)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_infra_rs_sp_acc_port_p")
			d.Partial(false)
		}
	}

	d.SetId(infraSpineP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciSpineProfileRead(d, m)
}

func resourceAciSpineProfileUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] SpineProfile: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	infraSpinePAttr := models.SpineProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraSpinePAttr.Annotation = Annotation.(string)
	} else {
		infraSpinePAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraSpinePAttr.NameAlias = NameAlias.(string)
	}
	infraSpineP := models.NewSpineProfile(fmt.Sprintf("infra/spprof-%s", name), "uni", desc, infraSpinePAttr)

	infraSpineP.Status = "modified"

	err := aciClient.Save(infraSpineP)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_infra_rs_sp_acc_port_p") {
		oldRel, newRel := d.GetChange("relation_infra_rs_sp_acc_port_p")
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
		return err
	}
	d.Partial(false)

	if d.HasChange("relation_infra_rs_sp_acc_port_p") {
		oldRel, newRel := d.GetChange("relation_infra_rs_sp_acc_port_p")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationinfraRsSpAccPortPFromSpineProfile(infraSpineP.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationinfraRsSpAccPortPFromSpineProfile(infraSpineP.DistinguishedName, relDn)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_infra_rs_sp_acc_port_p")
			d.Partial(false)

		}
	}

	d.SetId(infraSpineP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciSpineProfileRead(d, m)

}

func resourceAciSpineProfileRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	infraSpineP, err := getRemoteSpineProfile(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setSpineProfileAttributes(infraSpineP, d)

	infraRsSpAccPortPData, err := aciClient.ReadRelationinfraRsSpAccPortPFromSpineProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsSpAccPortP %v", err)
		d.Set("relation_infra_rs_sp_acc_port_p", make([]interface{}, 0, 1))

	} else {
		d.Set("relation_infra_rs_sp_acc_port_p", infraRsSpAccPortPData)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciSpineProfileDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "infraSpineP")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
