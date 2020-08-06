package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciLeafProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciLeafProfileCreate,
		Update: resourceAciLeafProfileUpdate,
		Read:   resourceAciLeafProfileRead,
		Delete: resourceAciLeafProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciLeafProfileImport,
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

			"relation_infra_rs_acc_card_p": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_infra_rs_acc_port_p": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
		}),
	}
}
func getRemoteLeafProfile(client *client.Client, dn string) (*models.LeafProfile, error) {
	infraNodePCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	infraNodeP := models.LeafProfileFromContainer(infraNodePCont)

	if infraNodeP.DistinguishedName == "" {
		return nil, fmt.Errorf("LeafProfile %s not found", infraNodeP.DistinguishedName)
	}

	return infraNodeP, nil
}

func setLeafProfileAttributes(infraNodeP *models.LeafProfile, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(infraNodeP.DistinguishedName)
	d.Set("description", infraNodeP.Description)
	infraNodePMap, _ := infraNodeP.ToMap()

	d.Set("name", infraNodePMap["name"])

	d.Set("annotation", infraNodePMap["annotation"])
	d.Set("name_alias", infraNodePMap["nameAlias"])
	return d
}

func resourceAciLeafProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	infraNodeP, err := getRemoteLeafProfile(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setLeafProfileAttributes(infraNodeP, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciLeafProfileCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] LeafProfile: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	infraNodePAttr := models.LeafProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraNodePAttr.Annotation = Annotation.(string)
	} else {
		infraNodePAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraNodePAttr.NameAlias = NameAlias.(string)
	}
	infraNodeP := models.NewLeafProfile(fmt.Sprintf("infra/nprof-%s", name), "uni", desc, infraNodePAttr)

	err := aciClient.Save(infraNodeP)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if relationToinfraRsAccCardP, ok := d.GetOk("relation_infra_rs_acc_card_p"); ok {
		relationParamList := toStringList(relationToinfraRsAccCardP.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationinfraRsAccCardPFromLeafProfile(infraNodeP.DistinguishedName, relationParam)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_infra_rs_acc_card_p")
			d.Partial(false)
		}
	}
	if relationToinfraRsAccPortP, ok := d.GetOk("relation_infra_rs_acc_port_p"); ok {
		relationParamList := toStringList(relationToinfraRsAccPortP.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationinfraRsAccPortPFromLeafProfile(infraNodeP.DistinguishedName, relationParam)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_infra_rs_acc_port_p")
			d.Partial(false)
		}
	}

	d.SetId(infraNodeP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciLeafProfileRead(d, m)
}

func resourceAciLeafProfileUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] LeafProfile: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	infraNodePAttr := models.LeafProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraNodePAttr.Annotation = Annotation.(string)
	} else {
		infraNodePAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraNodePAttr.NameAlias = NameAlias.(string)
	}
	infraNodeP := models.NewLeafProfile(fmt.Sprintf("infra/nprof-%s", name), "uni", desc, infraNodePAttr)

	infraNodeP.Status = "modified"

	err := aciClient.Save(infraNodeP)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if d.HasChange("relation_infra_rs_acc_card_p") {
		oldRel, newRel := d.GetChange("relation_infra_rs_acc_card_p")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationinfraRsAccCardPFromLeafProfile(infraNodeP.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationinfraRsAccCardPFromLeafProfile(infraNodeP.DistinguishedName, relDn)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_infra_rs_acc_card_p")
			d.Partial(false)

		}

	}
	if d.HasChange("relation_infra_rs_acc_port_p") {
		oldRel, newRel := d.GetChange("relation_infra_rs_acc_port_p")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationinfraRsAccPortPFromLeafProfile(infraNodeP.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationinfraRsAccPortPFromLeafProfile(infraNodeP.DistinguishedName, relDn)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_infra_rs_acc_port_p")
			d.Partial(false)

		}

	}

	d.SetId(infraNodeP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciLeafProfileRead(d, m)

}

func resourceAciLeafProfileRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	infraNodeP, err := getRemoteLeafProfile(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setLeafProfileAttributes(infraNodeP, d)

	infraRsAccCardPData, err := aciClient.ReadRelationinfraRsAccCardPFromLeafProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsAccCardP %v", err)

	} else {
		d.Set("relation_infra_rs_acc_card_p", infraRsAccCardPData)
	}

	infraRsAccPortPData, err := aciClient.ReadRelationinfraRsAccPortPFromLeafProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsAccPortP %v", err)

	} else {
		d.Set("relation_infra_rs_acc_port_p", infraRsAccPortPData)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciLeafProfileDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "infraNodeP")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
