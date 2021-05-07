package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAciLoopBackInterfaceProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciLoopBackInterfaceProfileCreate,
		Update: resourceAciLoopBackInterfaceProfileUpdate,
		Read:   resourceAciLoopBackInterfaceProfileRead,
		Delete: resourceAciLoopBackInterfaceProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciLoopBackInterfaceProfileImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"fabric_node_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"addr": &schema.Schema{
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

func getRemoteLoopBackInterfaceProfile(client *client.Client, dn string) (*models.LoopBackInterfaceProfile, error) {
	l3extLoopBackIfPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	l3extLoopBackIfP := models.LoopBackInterfaceProfileFromContainer(l3extLoopBackIfPCont)

	if l3extLoopBackIfP.DistinguishedName == "" {
		return nil, fmt.Errorf("LoopBackInterfaceProfile %s not found", l3extLoopBackIfP.DistinguishedName)
	}

	return l3extLoopBackIfP, nil
}

func setLoopBackInterfaceProfileAttributes(l3extLoopBackIfP *models.LoopBackInterfaceProfile, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()

	d.SetId(l3extLoopBackIfP.DistinguishedName)
	d.Set("description", l3extLoopBackIfP.Description)

	if dn != l3extLoopBackIfP.DistinguishedName {
		d.Set("fabric_node_dn", "")
	}

	l3extLoopBackIfPMap, _ := l3extLoopBackIfP.ToMap()
	d.Set("addr", l3extLoopBackIfPMap["addr"])
	d.Set("annotation", l3extLoopBackIfPMap["annotation"])
	d.Set("name_alias", l3extLoopBackIfPMap["nameAlias"])

	return d
}

func resourceAciLoopBackInterfaceProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()

	l3extLoopBackIfP, err := getRemoteLoopBackInterfaceProfile(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setLoopBackInterfaceProfileAttributes(l3extLoopBackIfP, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciLoopBackInterfaceProfileCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] LoopBackInterfaceProfile: Beginning Creation")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	addr := d.Get("addr").(string)

	FabricNodeDn := d.Get("fabric_node_dn").(string)

	l3extLoopBackIfPAttr := models.LoopBackInterfaceProfileAttributes{}
	if Addr, ok := d.GetOk("addr"); ok {
		l3extLoopBackIfPAttr.Addr = Addr.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extLoopBackIfPAttr.Annotation = Annotation.(string)
	} else {
		l3extLoopBackIfPAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l3extLoopBackIfPAttr.NameAlias = NameAlias.(string)
	}
	l3extLoopBackIfP := models.NewLoopBackInterfaceProfile(fmt.Sprintf("lbp-[%s]", addr), FabricNodeDn, desc, l3extLoopBackIfPAttr)

	err := aciClient.Save(l3extLoopBackIfP)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.Partial(false)

	d.SetId(l3extLoopBackIfP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciLoopBackInterfaceProfileRead(d, m)
}

func resourceAciLoopBackInterfaceProfileUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] LoopBackInterfaceProfile: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	addr := d.Get("addr").(string)

	FabricNodeDn := d.Get("fabric_node_dn").(string)

	l3extLoopBackIfPAttr := models.LoopBackInterfaceProfileAttributes{}
	if Addr, ok := d.GetOk("addr"); ok {
		l3extLoopBackIfPAttr.Addr = Addr.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extLoopBackIfPAttr.Annotation = Annotation.(string)
	} else {
		l3extLoopBackIfPAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l3extLoopBackIfPAttr.NameAlias = NameAlias.(string)
	}
	l3extLoopBackIfP := models.NewLoopBackInterfaceProfile(fmt.Sprintf("lbp-[%s]", addr), FabricNodeDn, desc, l3extLoopBackIfPAttr)

	l3extLoopBackIfP.Status = "modified"

	err := aciClient.Save(l3extLoopBackIfP)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.Partial(false)

	d.SetId(l3extLoopBackIfP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciLoopBackInterfaceProfileRead(d, m)

}

func resourceAciLoopBackInterfaceProfileRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	l3extLoopBackIfP, err := getRemoteLoopBackInterfaceProfile(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setLoopBackInterfaceProfileAttributes(l3extLoopBackIfP, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciLoopBackInterfaceProfileDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "l3extLoopBackIfP")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
