package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciAccessPortSelector() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciAccessPortSelectorCreate,
		Update: resourceAciAccessPortSelectorUpdate,
		Read:   resourceAciAccessPortSelectorRead,
		Delete: resourceAciAccessPortSelectorDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciAccessPortSelectorImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"leaf_interface_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"access_port_selector_type": &schema.Schema{
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

			"relation_infra_rs_acc_base_grp": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
		}),
	}
}
func getRemoteAccessPortSelector(client *client.Client, dn string) (*models.AccessPortSelector, error) {
	infraHPortSCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	infraHPortS := models.AccessPortSelectorFromContainer(infraHPortSCont)

	if infraHPortS.DistinguishedName == "" {
		return nil, fmt.Errorf("AccessPortSelector %s not found", infraHPortS.DistinguishedName)
	}

	return infraHPortS, nil
}

func setAccessPortSelectorAttributes(infraHPortS *models.AccessPortSelector, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(infraHPortS.DistinguishedName)
	d.Set("description", infraHPortS.Description)
	// d.Set("leaf_interface_profile_dn", GetParentDn(infraHPortS.DistinguishedName))
	if dn != infraHPortS.DistinguishedName {
		d.Set("leaf_interface_profile_dn", "")
	}
	infraHPortSMap, _ := infraHPortS.ToMap()

	d.Set("name", infraHPortSMap["name"])

	d.Set("access_port_selector_type", infraHPortSMap["type"])

	d.Set("annotation", infraHPortSMap["annotation"])
	d.Set("name_alias", infraHPortSMap["nameAlias"])
	d.Set("access_port_selector_type", infraHPortSMap["type"])
	return d
}

func resourceAciAccessPortSelectorImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	infraHPortS, err := getRemoteAccessPortSelector(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setAccessPortSelectorAttributes(infraHPortS, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciAccessPortSelectorCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] AccessPortSelector: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	access_port_selector_type := d.Get("access_port_selector_type").(string)

	LeafInterfaceProfileDn := d.Get("leaf_interface_profile_dn").(string)

	infraHPortSAttr := models.AccessPortSelectorAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraHPortSAttr.Annotation = Annotation.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraHPortSAttr.NameAlias = NameAlias.(string)
	}
	if AccessPortSelector_type, ok := d.GetOk("access_port_selector_type"); ok {
		infraHPortSAttr.AccessPortSelector_type = AccessPortSelector_type.(string)
	}
	infraHPortS := models.NewAccessPortSelector(fmt.Sprintf("hports-%s-typ-%s", name, access_port_selector_type), LeafInterfaceProfileDn, desc, infraHPortSAttr)

	err := aciClient.Save(infraHPortS)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.SetPartial("access_port_selector_type")

	d.Partial(false)

	if relationToinfraRsAccBaseGrp, ok := d.GetOk("relation_infra_rs_acc_base_grp"); ok {
		relationParam := relationToinfraRsAccBaseGrp.(string)
		err = aciClient.CreateRelationinfraRsAccBaseGrpFromAccessPortSelector(infraHPortS.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_acc_base_grp")
		d.Partial(false)

	}

	d.SetId(infraHPortS.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciAccessPortSelectorRead(d, m)
}

func resourceAciAccessPortSelectorUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] AccessPortSelector: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	access_port_selector_type := d.Get("access_port_selector_type").(string)

	LeafInterfaceProfileDn := d.Get("leaf_interface_profile_dn").(string)

	infraHPortSAttr := models.AccessPortSelectorAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraHPortSAttr.Annotation = Annotation.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraHPortSAttr.NameAlias = NameAlias.(string)
	}
	if AccessPortSelector_type, ok := d.GetOk("access_port_selector_type"); ok {
		infraHPortSAttr.AccessPortSelector_type = AccessPortSelector_type.(string)
	}
	infraHPortS := models.NewAccessPortSelector(fmt.Sprintf("hports-%s-typ-%s", name, access_port_selector_type), LeafInterfaceProfileDn, desc, infraHPortSAttr)

	infraHPortS.Status = "modified"

	err := aciClient.Save(infraHPortS)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.SetPartial("access_port_selector_type")

	d.Partial(false)

	if d.HasChange("relation_infra_rs_acc_base_grp") {
		_, newRelParam := d.GetChange("relation_infra_rs_acc_base_grp")
		err = aciClient.DeleteRelationinfraRsAccBaseGrpFromAccessPortSelector(infraHPortS.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationinfraRsAccBaseGrpFromAccessPortSelector(infraHPortS.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_acc_base_grp")
		d.Partial(false)

	}

	d.SetId(infraHPortS.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciAccessPortSelectorRead(d, m)

}

func resourceAciAccessPortSelectorRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	infraHPortS, err := getRemoteAccessPortSelector(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setAccessPortSelectorAttributes(infraHPortS, d)

	infraRsAccBaseGrpData, err := aciClient.ReadRelationinfraRsAccBaseGrpFromAccessPortSelector(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsAccBaseGrp %v", err)

	} else {
		if _, ok := d.GetOk("relation_infra_rs_acc_base_grp"); ok {
			tfName := d.Get("relation_infra_rs_acc_base_grp").(string)
			if tfName != infraRsAccBaseGrpData {
				d.Set("relation_infra_rs_acc_base_grp", "")
			}
		}
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciAccessPortSelectorDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "infraHPortS")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
