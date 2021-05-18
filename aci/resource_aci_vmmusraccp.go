package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciVMMCredential() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciVMMCredentialCreate,
		Update: resourceAciVMMCredentialUpdate,
		Read:   resourceAciVMMCredentialRead,
		Delete: resourceAciVMMCredentialDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciVMMCredentialImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"vmm_domain_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"pwd": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"usr": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func getRemoteVMMCredential(client *client.Client, dn string) (*models.VMMCredential, error) {
	vmmUsrAccPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	vmmUsrAccP := models.VMMCredentialFromContainer(vmmUsrAccPCont)
	if vmmUsrAccP.DistinguishedName == "" {
		return nil, fmt.Errorf("VMMCredential %s not found", vmmUsrAccP.DistinguishedName)
	}
	return vmmUsrAccP, nil
}

func setVMMCredentialAttributes(vmmUsrAccP *models.VMMCredential, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(vmmUsrAccP.DistinguishedName)
	d.Set("description", vmmUsrAccP.Description)
	vmmUsrAccPMap, _ := vmmUsrAccP.ToMap()
	d.Set("vmm_domain_dn", GetParentDn(vmmUsrAccP.DistinguishedName, fmt.Sprintf("/usracc-%s", vmmUsrAccPMap["name"])))
	d.Set("annotation", vmmUsrAccPMap["annotation"])
	d.Set("name", vmmUsrAccPMap["name"])
	d.Set("pwd", vmmUsrAccPMap["pwd"])
	d.Set("usr", vmmUsrAccPMap["usr"])
	d.Set("name_alias", vmmUsrAccPMap["nameAlias"])
	return d
}

func resourceAciVMMCredentialImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	vmmUsrAccP, err := getRemoteVMMCredential(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled := setVMMCredentialAttributes(vmmUsrAccP, d)
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciVMMCredentialCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] VMMCredential: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	VMMDomainDn := d.Get("vmm_domain_dn").(string)

	vmmUsrAccPAttr := models.VMMCredentialAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vmmUsrAccPAttr.Annotation = Annotation.(string)
	} else {
		vmmUsrAccPAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		vmmUsrAccPAttr.Name = Name.(string)
	}

	if Pwd, ok := d.GetOk("pwd"); ok {
		vmmUsrAccPAttr.Pwd = Pwd.(string)
	}

	if Usr, ok := d.GetOk("usr"); ok {
		vmmUsrAccPAttr.Usr = Usr.(string)
	}
	vmmUsrAccP := models.NewVMMCredential(fmt.Sprintf("usracc-%s", name), VMMDomainDn, desc, nameAlias, vmmUsrAccPAttr)

	err := aciClient.Save(vmmUsrAccP)
	if err != nil {
		return err
	}
	d.Partial(true)
	d.SetPartial("name")
	d.Partial(false)

	d.SetId(vmmUsrAccP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciVMMCredentialRead(d, m)
}

func resourceAciVMMCredentialUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] VMMCredential: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	VMMDomainDn := d.Get("vmm_domain_dn").(string)
	vmmUsrAccPAttr := models.VMMCredentialAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		vmmUsrAccPAttr.Annotation = Annotation.(string)
	} else {
		vmmUsrAccPAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		vmmUsrAccPAttr.Name = Name.(string)
	}

	if Pwd, ok := d.GetOk("pwd"); ok {
		vmmUsrAccPAttr.Pwd = Pwd.(string)
	}

	if Usr, ok := d.GetOk("usr"); ok {
		vmmUsrAccPAttr.Usr = Usr.(string)
	}
	vmmUsrAccP := models.NewVMMCredential(fmt.Sprintf("usracc-%s", name), VMMDomainDn, desc, nameAlias, vmmUsrAccPAttr)

	vmmUsrAccP.Status = "modified"
	err := aciClient.Save(vmmUsrAccP)
	if err != nil {
		return err
	}
	d.Partial(true)
	d.SetPartial("name")
	d.Partial(false)

	d.SetId(vmmUsrAccP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciVMMCredentialRead(d, m)
}

func resourceAciVMMCredentialRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	vmmUsrAccP, err := getRemoteVMMCredential(aciClient, dn)
	if err != nil {
		d.SetId("")
		return err
	}
	setVMMCredentialAttributes(vmmUsrAccP, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciVMMCredentialDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "vmmUsrAccP")
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return err
}
