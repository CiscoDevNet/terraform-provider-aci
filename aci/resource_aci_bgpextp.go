package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAciL3outBgpExternalPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciL3outBgpExternalPolicyCreate,
		Update: resourceAciL3outBgpExternalPolicyUpdate,
		Read:   resourceAciL3outBgpExternalPolicyRead,
		Delete: resourceAciL3outBgpExternalPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciL3outBgpExternalPolicyImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"l3_outside_dn": &schema.Schema{
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
func getRemoteL3outBgpExternalPolicy(client *client.Client, dn string) (*models.L3outBgpExternalPolicy, error) {
	bgpExtPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	bgpExtP := models.L3outBgpExternalPolicyFromContainer(bgpExtPCont)

	if bgpExtP.DistinguishedName == "" {
		return nil, fmt.Errorf("L3outBgpExternalPolicy %s not found", bgpExtP.DistinguishedName)
	}

	return bgpExtP, nil
}

func setL3outBgpExternalPolicyAttributes(bgpExtP *models.L3outBgpExternalPolicy, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(bgpExtP.DistinguishedName)
	d.Set("description", bgpExtP.Description)
	dn := d.Id()
	if dn != bgpExtP.DistinguishedName {
		d.Set("l3_outside_dn", "")
	}
	bgpExtPMap, _ := bgpExtP.ToMap()

	d.Set("annotation", bgpExtPMap["annotation"])
	d.Set("name_alias", bgpExtPMap["nameAlias"])
	return d
}

func resourceAciL3outBgpExternalPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	bgpExtP, err := getRemoteL3outBgpExternalPolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setL3outBgpExternalPolicyAttributes(bgpExtP, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciL3outBgpExternalPolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] L3outBgpExternalPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	L3OutsideDn := d.Get("l3_outside_dn").(string)

	bgpExtPAttr := models.L3outBgpExternalPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		bgpExtPAttr.Annotation = Annotation.(string)
	} else {
		bgpExtPAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		bgpExtPAttr.NameAlias = NameAlias.(string)
	}
	bgpExtP := models.NewL3outBgpExternalPolicy(fmt.Sprintf("bgpExtP"), L3OutsideDn, desc, bgpExtPAttr)

	err := aciClient.Save(bgpExtP)
	if err != nil {
		return err
	}
	d.Partial(true)
	d.Partial(false)

	checkDns := make([]string, 0, 1)

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return err
	}
	d.Partial(false)

	d.SetId(bgpExtP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciL3outBgpExternalPolicyRead(d, m)
}

func resourceAciL3outBgpExternalPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] L3outBgpExternalPolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	L3OutsideDn := d.Get("l3_outside_dn").(string)

	bgpExtPAttr := models.L3outBgpExternalPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		bgpExtPAttr.Annotation = Annotation.(string)
	} else {
		bgpExtPAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		bgpExtPAttr.NameAlias = NameAlias.(string)
	}
	bgpExtP := models.NewL3outBgpExternalPolicy(fmt.Sprintf("bgpExtP"), L3OutsideDn, desc, bgpExtPAttr)

	bgpExtP.Status = "modified"

	err := aciClient.Save(bgpExtP)

	if err != nil {
		return err
	}
	d.Partial(true)
	d.Partial(false)

	checkDns := make([]string, 0, 1)

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return err
	}
	d.Partial(false)

	d.SetId(bgpExtP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciL3outBgpExternalPolicyRead(d, m)

}

func resourceAciL3outBgpExternalPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	bgpExtP, err := getRemoteL3outBgpExternalPolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setL3outBgpExternalPolicyAttributes(bgpExtP, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciL3outBgpExternalPolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "bgpExtP")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
