package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciX509Certificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciX509CertificateCreate,
		Update: resourceAciX509CertificateUpdate,
		Read:   resourceAciX509CertificateRead,
		Delete: resourceAciX509CertificateDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciX509CertificateImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"local_user_dn": &schema.Schema{
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

			"data": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}
func getRemoteX509Certificate(client *client.Client, dn string) (*models.X509Certificate, error) {
	aaaUserCertCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	aaaUserCert := models.X509CertificateFromContainer(aaaUserCertCont)

	if aaaUserCert.DistinguishedName == "" {
		return nil, fmt.Errorf("X509Certificate %s not found", aaaUserCert.DistinguishedName)
	}

	return aaaUserCert, nil
}

func setX509CertificateAttributes(aaaUserCert *models.X509Certificate, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(aaaUserCert.DistinguishedName)
	d.Set("description", aaaUserCert.Description)
	// d.Set("local_user_dn", GetParentDn(aaaUserCert.DistinguishedName))
	if dn != aaaUserCert.DistinguishedName {
		d.Set("local_user_dn", "")
	}
	aaaUserCertMap, _ := aaaUserCert.ToMap()

	d.Set("name", aaaUserCertMap["name"])

	d.Set("annotation", aaaUserCertMap["annotation"])
	d.Set("data", aaaUserCertMap["data"])
	d.Set("name_alias", aaaUserCertMap["nameAlias"])
	return d
}

func resourceAciX509CertificateImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	aaaUserCert, err := getRemoteX509Certificate(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setX509CertificateAttributes(aaaUserCert, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciX509CertificateCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] X509Certificate: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	LocalUserDn := d.Get("local_user_dn").(string)

	aaaUserCertAttr := models.X509CertificateAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaUserCertAttr.Annotation = Annotation.(string)
	}
	if Data, ok := d.GetOk("data"); ok {
		aaaUserCertAttr.Data = Data.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		aaaUserCertAttr.NameAlias = NameAlias.(string)
	}
	aaaUserCert := models.NewX509Certificate(fmt.Sprintf("usercert-%s", name), LocalUserDn, desc, aaaUserCertAttr)

	err := aciClient.Save(aaaUserCert)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(aaaUserCert.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciX509CertificateRead(d, m)
}

func resourceAciX509CertificateUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] X509Certificate: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	LocalUserDn := d.Get("local_user_dn").(string)

	aaaUserCertAttr := models.X509CertificateAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaUserCertAttr.Annotation = Annotation.(string)
	}
	if Data, ok := d.GetOk("data"); ok {
		aaaUserCertAttr.Data = Data.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		aaaUserCertAttr.NameAlias = NameAlias.(string)
	}
	aaaUserCert := models.NewX509Certificate(fmt.Sprintf("usercert-%s", name), LocalUserDn, desc, aaaUserCertAttr)

	aaaUserCert.Status = "modified"

	err := aciClient.Save(aaaUserCert)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(aaaUserCert.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciX509CertificateRead(d, m)

}

func resourceAciX509CertificateRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	aaaUserCert, err := getRemoteX509Certificate(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setX509CertificateAttributes(aaaUserCert, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciX509CertificateDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "aaaUserCert")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
