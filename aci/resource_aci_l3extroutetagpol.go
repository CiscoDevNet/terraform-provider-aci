package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciL3outRouteTagPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciL3outRouteTagPolicyCreate,
		Update: resourceAciL3outRouteTagPolicyUpdate,
		Read:   resourceAciL3outRouteTagPolicyRead,
		Delete: resourceAciL3outRouteTagPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciL3outRouteTagPolicyImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"tenant_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

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

			"tag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}
func getRemoteL3outRouteTagPolicy(client *client.Client, dn string) (*models.L3outRouteTagPolicy, error) {
	l3extRouteTagPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	l3extRouteTagPol := models.L3outRouteTagPolicyFromContainer(l3extRouteTagPolCont)

	if l3extRouteTagPol.DistinguishedName == "" {
		return nil, fmt.Errorf("L3outRouteTagPolicy %s not found", l3extRouteTagPol.DistinguishedName)
	}

	return l3extRouteTagPol, nil
}

func setL3outRouteTagPolicyAttributes(l3extRouteTagPol *models.L3outRouteTagPolicy, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(l3extRouteTagPol.DistinguishedName)
	d.Set("description", l3extRouteTagPol.Description)
	dn := d.Id()
	if dn != l3extRouteTagPol.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	l3extRouteTagPolMap, _ := l3extRouteTagPol.ToMap()

	d.Set("name", l3extRouteTagPolMap["name"])

	d.Set("annotation", l3extRouteTagPolMap["annotation"])
	d.Set("name_alias", l3extRouteTagPolMap["nameAlias"])
	d.Set("tag", l3extRouteTagPolMap["tag"])
	return d
}

func resourceAciL3outRouteTagPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	l3extRouteTagPol, err := getRemoteL3outRouteTagPolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setL3outRouteTagPolicyAttributes(l3extRouteTagPol, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciL3outRouteTagPolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] L3outRouteTagPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	l3extRouteTagPolAttr := models.L3outRouteTagPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extRouteTagPolAttr.Annotation = Annotation.(string)
	} else {
		l3extRouteTagPolAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l3extRouteTagPolAttr.NameAlias = NameAlias.(string)
	}
	if Tag, ok := d.GetOk("tag"); ok {
		l3extRouteTagPolAttr.Tag = Tag.(string)
	}
	l3extRouteTagPol := models.NewL3outRouteTagPolicy(fmt.Sprintf("rttag-%s", name), TenantDn, desc, l3extRouteTagPolAttr)

	err := aciClient.Save(l3extRouteTagPol)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(l3extRouteTagPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciL3outRouteTagPolicyRead(d, m)
}

func resourceAciL3outRouteTagPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] L3outRouteTagPolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	l3extRouteTagPolAttr := models.L3outRouteTagPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extRouteTagPolAttr.Annotation = Annotation.(string)
	} else {
		l3extRouteTagPolAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l3extRouteTagPolAttr.NameAlias = NameAlias.(string)
	}
	if Tag, ok := d.GetOk("tag"); ok {
		l3extRouteTagPolAttr.Tag = Tag.(string)
	}
	l3extRouteTagPol := models.NewL3outRouteTagPolicy(fmt.Sprintf("rttag-%s", name), TenantDn, desc, l3extRouteTagPolAttr)

	l3extRouteTagPol.Status = "modified"

	err := aciClient.Save(l3extRouteTagPol)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(l3extRouteTagPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciL3outRouteTagPolicyRead(d, m)

}

func resourceAciL3outRouteTagPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	l3extRouteTagPol, err := getRemoteL3outRouteTagPolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setL3outRouteTagPolicyAttributes(l3extRouteTagPol, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciL3outRouteTagPolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "l3extRouteTagPol")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
