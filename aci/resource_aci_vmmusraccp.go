package aci

import (
	"fmt"

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

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"vmm_domain_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"annotation": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Computed:    true,
				Description: "Mo doc not defined in techpub!!!",
			},

			"name_alias": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Computed:    true,
				Description: "Mo doc not defined in techpub!!!",
			},

			"pwd": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Computed:    true,
				Description: "user account profile password",
			},

			"usr": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Computed:    true,
				Description: "user name",
			},
		}),
	}
}

func getRemoteVMMCredential(client *client.Client, dn string) (*models.VMMCredential, error) {
	vmmUsrAccPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	vmmUsrAccP := models.VMMCredentialFromContainer(vmmUsrAccPCont)

	if vmmUsrAccP.DistinguishedName == "" {
		return nil, fmt.Errorf("Bridge Domain %s not found", vmmUsrAccP.DistinguishedName)
	}

	return vmmUsrAccP, nil
}

func setVMMCredentialAttributes(vmmUsrAccP *models.VMMCredential, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(vmmUsrAccP.DistinguishedName)
	d.Set("description", vmmUsrAccP.Description)
	// d.Set("vmm_domain_dn", GetParentDn(vmmUsrAccP.DistinguishedName))
	vmmUsrAccPMap, _ := vmmUsrAccP.ToMap()
	d.Set("name", vmmUsrAccPMap["name"])
	d.Set("vmm_domain_dn", GetParentDn(vmmUsrAccP.DistinguishedName, fmt.Sprintf("/usracc-%s", vmmUsrAccPMap["name"])))

	d.Set("annotation", vmmUsrAccPMap["annotation"])
	d.Set("name_alias", vmmUsrAccPMap["nameAlias"])
	d.Set("pwd", vmmUsrAccPMap["pwd"])
	d.Set("usr", vmmUsrAccPMap["usr"])
	return d
}

func resourceAciVMMCredentialImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	aciClient := m.(*client.Client)

	dn := d.Id()

	vmmUsrAccP, err := getRemoteVMMCredential(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setVMMCredentialAttributes(vmmUsrAccP, d)
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciVMMCredentialCreate(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	VMMDomainDn := d.Get("vmm_domain_dn").(string)

	vmmUsrAccPAttr := models.VMMCredentialAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vmmUsrAccPAttr.Annotation = Annotation.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vmmUsrAccPAttr.NameAlias = NameAlias.(string)
	}
	if Pwd, ok := d.GetOk("pwd"); ok {
		vmmUsrAccPAttr.Pwd = Pwd.(string)
	}
	if Usr, ok := d.GetOk("usr"); ok {
		vmmUsrAccPAttr.Usr = Usr.(string)
	}
	vmmUsrAccP := models.NewVMMCredential(fmt.Sprintf("usracc-%s", name), VMMDomainDn, desc, vmmUsrAccPAttr)

	err := aciClient.Save(vmmUsrAccP)
	if err != nil {
		return err
	}

	d.SetId(vmmUsrAccP.DistinguishedName)
	return resourceAciVMMCredentialRead(d, m)
}

func resourceAciVMMCredentialUpdate(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)
	VMMDomainDn := d.Get("vmm_domain_dn").(string)

	vmmUsrAccPAttr := models.VMMCredentialAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vmmUsrAccPAttr.Annotation = Annotation.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vmmUsrAccPAttr.NameAlias = NameAlias.(string)
	}
	if Pwd, ok := d.GetOk("pwd"); ok {
		vmmUsrAccPAttr.Pwd = Pwd.(string)
	}
	if Usr, ok := d.GetOk("usr"); ok {
		vmmUsrAccPAttr.Usr = Usr.(string)
	}
	vmmUsrAccP := models.NewVMMCredential(fmt.Sprintf("usracc-%s", name), VMMDomainDn, desc, vmmUsrAccPAttr)

	vmmUsrAccP.Status = "modified"

	err := aciClient.Save(vmmUsrAccP)

	if err != nil {
		return err
	}

	d.SetId(vmmUsrAccP.DistinguishedName)
	return resourceAciVMMCredentialRead(d, m)

}

func resourceAciVMMCredentialRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	dn := d.Id()
	vmmUsrAccP, err := getRemoteVMMCredential(aciClient, dn)

	if err != nil {
		return err
	}
	setVMMCredentialAttributes(vmmUsrAccP, d)
	return nil
}

func resourceAciVMMCredentialDelete(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "vmmUsrAccP")
	if err != nil {
		return err
	}

	d.SetId("")
	return err
}
