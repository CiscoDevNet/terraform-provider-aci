package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciAccessGeneric() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciAccessGenericCreate,
		Update: resourceAciAccessGenericUpdate,
		Read:   resourceAciAccessGenericRead,
		Delete: resourceAciAccessGenericDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciAccessGenericImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"attachable_access_entity_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"relation_infra_rs_func_to_epg": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
		}),
	}
}
func getRemoteAccessGeneric(client *client.Client, dn string) (*models.AccessGeneric, error) {
	infraGenericCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	infraGeneric := models.AccessGenericFromContainer(infraGenericCont)

	if infraGeneric.DistinguishedName == "" {
		return nil, fmt.Errorf("AccessGeneric %s not found", infraGeneric.DistinguishedName)
	}

	return infraGeneric, nil
}

func setAccessGenericAttributes(infraGeneric *models.AccessGeneric, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(infraGeneric.DistinguishedName)
	d.Set("description", infraGeneric.Description)
	if dn != infraGeneric.DistinguishedName {
		d.Set("attachable_access_entity_profile_dn", "")
	}
	infraGenericMap, _ := infraGeneric.ToMap()

	d.Set("name", infraGenericMap["name"])

	d.Set("annotation", infraGenericMap["annotation"])
	d.Set("name_alias", infraGenericMap["nameAlias"])
	return d
}

func resourceAciAccessGenericImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	infraGeneric, err := getRemoteAccessGeneric(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setAccessGenericAttributes(infraGeneric, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciAccessGenericCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] AccessGeneric: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	AttachableAccessEntityProfileDn := d.Get("attachable_access_entity_profile_dn").(string)

	infraGenericAttr := models.AccessGenericAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraGenericAttr.Annotation = Annotation.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraGenericAttr.NameAlias = NameAlias.(string)
	}
	infraGeneric := models.NewAccessGeneric(fmt.Sprintf("gen-%s", name), AttachableAccessEntityProfileDn, desc, infraGenericAttr)

	err := aciClient.Save(infraGeneric)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if relationToinfraRsFuncToEpg, ok := d.GetOk("relation_infra_rs_func_to_epg"); ok {
		relationParamList := toStringList(relationToinfraRsFuncToEpg.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationinfraRsFuncToEpgFromAccessGeneric(infraGeneric.DistinguishedName, relationParam)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_infra_rs_func_to_epg")
			d.Partial(false)
		}
	}

	d.SetId(infraGeneric.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciAccessGenericRead(d, m)
}

func resourceAciAccessGenericUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] AccessGeneric: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	AttachableAccessEntityProfileDn := d.Get("attachable_access_entity_profile_dn").(string)

	infraGenericAttr := models.AccessGenericAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraGenericAttr.Annotation = Annotation.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraGenericAttr.NameAlias = NameAlias.(string)
	}
	infraGeneric := models.NewAccessGeneric(fmt.Sprintf("gen-%s", name), AttachableAccessEntityProfileDn, desc, infraGenericAttr)

	infraGeneric.Status = "modified"

	err := aciClient.Save(infraGeneric)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if d.HasChange("relation_infra_rs_func_to_epg") {
		oldRel, newRel := d.GetChange("relation_infra_rs_func_to_epg")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationinfraRsFuncToEpgFromAccessGeneric(infraGeneric.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationinfraRsFuncToEpgFromAccessGeneric(infraGeneric.DistinguishedName, relDn)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.SetPartial("relation_infra_rs_func_to_epg")
			d.Partial(false)

		}

	}

	d.SetId(infraGeneric.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciAccessGenericRead(d, m)

}

func resourceAciAccessGenericRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	infraGeneric, err := getRemoteAccessGeneric(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setAccessGenericAttributes(infraGeneric, d)

	infraRsFuncToEpgData, err := aciClient.ReadRelationinfraRsFuncToEpgFromAccessGeneric(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsFuncToEpg %v", err)

	} else {
		d.Set("relation_infra_rs_func_to_epg", infraRsFuncToEpgData)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciAccessGenericDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "infraGeneric")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
