package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciAttachableAccessEntityProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciAttachableAccessEntityProfileCreate,
		Update: resourceAciAttachableAccessEntityProfileUpdate,
		Read:   resourceAciAttachableAccessEntityProfileRead,
		Delete: resourceAciAttachableAccessEntityProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciAttachableAccessEntityProfileImport,
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

			"relation_infra_rs_dom_p": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
		}),
	}
}
func getRemoteAttachableAccessEntityProfile(client *client.Client, dn string) (*models.AttachableAccessEntityProfile, error) {
	infraAttEntityPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	infraAttEntityP := models.AttachableAccessEntityProfileFromContainer(infraAttEntityPCont)

	if infraAttEntityP.DistinguishedName == "" {
		return nil, fmt.Errorf("AttachableAccessEntityProfile %s not found", infraAttEntityP.DistinguishedName)
	}

	return infraAttEntityP, nil
}

func setAttachableAccessEntityProfileAttributes(infraAttEntityP *models.AttachableAccessEntityProfile, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(infraAttEntityP.DistinguishedName)
	d.Set("description", infraAttEntityP.Description)
	infraAttEntityPMap, _ := infraAttEntityP.ToMap()

	d.Set("name", infraAttEntityPMap["name"])

	d.Set("annotation", infraAttEntityPMap["annotation"])
	d.Set("name_alias", infraAttEntityPMap["nameAlias"])
	return d
}

func resourceAciAttachableAccessEntityProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	infraAttEntityP, err := getRemoteAttachableAccessEntityProfile(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setAttachableAccessEntityProfileAttributes(infraAttEntityP, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciAttachableAccessEntityProfileCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] AttachableAccessEntityProfile: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	infraAttEntityPAttr := models.AttachableAccessEntityProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraAttEntityPAttr.Annotation = Annotation.(string)
	} else {
		infraAttEntityPAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraAttEntityPAttr.NameAlias = NameAlias.(string)
	}
	infraAttEntityP := models.NewAttachableAccessEntityProfile(fmt.Sprintf("infra/attentp-%s", name), "uni", desc, infraAttEntityPAttr)

	err := aciClient.Save(infraAttEntityP)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	checkDns := make([]string, 0, 1)

	if relationToinfraRsDomP, ok := d.GetOk("relation_infra_rs_dom_p"); ok {
		relationParamList := toStringList(relationToinfraRsDomP.(*schema.Set).List())
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

	if relationToinfraRsDomP, ok := d.GetOk("relation_infra_rs_dom_p"); ok {
		relationParamList := toStringList(relationToinfraRsDomP.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationinfraRsDomPFromAttachableAccessEntityProfile(infraAttEntityP.DistinguishedName, relationParam)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_infra_rs_dom_p")
			d.Partial(false)
		}
	}

	d.SetId(infraAttEntityP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciAttachableAccessEntityProfileRead(d, m)
}

func resourceAciAttachableAccessEntityProfileUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] AttachableAccessEntityProfile: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	infraAttEntityPAttr := models.AttachableAccessEntityProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraAttEntityPAttr.Annotation = Annotation.(string)
	} else {
		infraAttEntityPAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraAttEntityPAttr.NameAlias = NameAlias.(string)
	}
	infraAttEntityP := models.NewAttachableAccessEntityProfile(fmt.Sprintf("infra/attentp-%s", name), "uni", desc, infraAttEntityPAttr)

	infraAttEntityP.Status = "modified"

	err := aciClient.Save(infraAttEntityP)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_infra_rs_dom_p") {
		oldRel, newRel := d.GetChange("relation_infra_rs_dom_p")
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

	if d.HasChange("relation_infra_rs_dom_p") {
		oldRel, newRel := d.GetChange("relation_infra_rs_dom_p")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationinfraRsDomPFromAttachableAccessEntityProfile(infraAttEntityP.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationinfraRsDomPFromAttachableAccessEntityProfile(infraAttEntityP.DistinguishedName, relDn)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_infra_rs_dom_p")
			d.Partial(false)

		}

	}

	d.SetId(infraAttEntityP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciAttachableAccessEntityProfileRead(d, m)

}

func resourceAciAttachableAccessEntityProfileRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	infraAttEntityP, err := getRemoteAttachableAccessEntityProfile(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setAttachableAccessEntityProfileAttributes(infraAttEntityP, d)

	infraRsDomPData, err := aciClient.ReadRelationinfraRsDomPFromAttachableAccessEntityProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsDomP %v", err)
		d.Set("relation_infra_rs_dom_p", make([]string, 0, 1))

	} else {
		d.Set("relation_infra_rs_dom_p", infraRsDomPData)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciAttachableAccessEntityProfileDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "infraAttEntityP")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
