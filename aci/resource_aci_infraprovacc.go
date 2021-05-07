package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAciVlanEncapsulationforVxlanTraffic() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciVlanEncapsulationforVxlanTrafficCreate,
		Update: resourceAciVlanEncapsulationforVxlanTrafficUpdate,
		Read:   resourceAciVlanEncapsulationforVxlanTrafficRead,
		Delete: resourceAciVlanEncapsulationforVxlanTrafficDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciVlanEncapsulationforVxlanTrafficImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"attachable_access_entity_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}
func getRemoteVlanEncapsulationforVxlanTraffic(client *client.Client, dn string) (*models.VlanEncapsulationforVxlanTraffic, error) {
	infraProvAccCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	infraProvAcc := models.VlanEncapsulationforVxlanTrafficFromContainer(infraProvAccCont)

	if infraProvAcc.DistinguishedName == "" {
		return nil, fmt.Errorf("VlanEncapsulationforVxlanTraffic %s not found", infraProvAcc.DistinguishedName)
	}

	return infraProvAcc, nil
}

func setVlanEncapsulationforVxlanTrafficAttributes(infraProvAcc *models.VlanEncapsulationforVxlanTraffic, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(infraProvAcc.DistinguishedName)
	d.Set("description", infraProvAcc.Description)
	// d.Set("attachable_access_entity_profile_dn", GetParentDn(infraProvAcc.DistinguishedName))
	if dn != infraProvAcc.DistinguishedName {
		d.Set("attachable_access_entity_profile_dn", "")
	}
	infraProvAccMap, _ := infraProvAcc.ToMap()

	d.Set("annotation", infraProvAccMap["annotation"])
	d.Set("name_alias", infraProvAccMap["nameAlias"])
	return d
}

func resourceAciVlanEncapsulationforVxlanTrafficImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	infraProvAcc, err := getRemoteVlanEncapsulationforVxlanTraffic(aciClient, dn)

	if err != nil {
		return nil, err
	}
	pDN := GetParentDn(dn, fmt.Sprintf("/provacc"))
	d.Set("attachable_access_entity_profile_dn", pDN)
	schemaFilled := setVlanEncapsulationforVxlanTrafficAttributes(infraProvAcc, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciVlanEncapsulationforVxlanTrafficCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] VlanEncapsulationforVxlanTraffic: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	AttachableAccessEntityProfileDn := d.Get("attachable_access_entity_profile_dn").(string)

	infraProvAccAttr := models.VlanEncapsulationforVxlanTrafficAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraProvAccAttr.Annotation = Annotation.(string)
	} else {
		infraProvAccAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraProvAccAttr.NameAlias = NameAlias.(string)
	}
	infraProvAcc := models.NewVlanEncapsulationforVxlanTraffic(fmt.Sprintf("provacc"), AttachableAccessEntityProfileDn, desc, infraProvAccAttr)

	err := aciClient.Save(infraProvAcc)
	if err != nil {
		return err
	}
	d.Partial(true)
	d.Partial(false)

	d.SetId(infraProvAcc.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciVlanEncapsulationforVxlanTrafficRead(d, m)
}

func resourceAciVlanEncapsulationforVxlanTrafficUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] VlanEncapsulationforVxlanTraffic: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	AttachableAccessEntityProfileDn := d.Get("attachable_access_entity_profile_dn").(string)

	infraProvAccAttr := models.VlanEncapsulationforVxlanTrafficAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraProvAccAttr.Annotation = Annotation.(string)
	} else {
		infraProvAccAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraProvAccAttr.NameAlias = NameAlias.(string)
	}
	infraProvAcc := models.NewVlanEncapsulationforVxlanTraffic(fmt.Sprintf("provacc"), AttachableAccessEntityProfileDn, desc, infraProvAccAttr)

	infraProvAcc.Status = "modified"

	err := aciClient.Save(infraProvAcc)

	if err != nil {
		return err
	}
	d.Partial(true)
	d.Partial(false)

	d.SetId(infraProvAcc.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciVlanEncapsulationforVxlanTrafficRead(d, m)

}

func resourceAciVlanEncapsulationforVxlanTrafficRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	infraProvAcc, err := getRemoteVlanEncapsulationforVxlanTraffic(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setVlanEncapsulationforVxlanTrafficAttributes(infraProvAcc, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciVlanEncapsulationforVxlanTrafficDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "infraProvAcc")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
