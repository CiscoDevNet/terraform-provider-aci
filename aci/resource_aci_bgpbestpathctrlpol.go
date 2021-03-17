package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAciBgpBestPathPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciBgpBestPathPolicyCreate,
		Update: resourceAciBgpBestPathPolicyUpdate,
		Read:   resourceAciBgpBestPathPolicyRead,
		Delete: resourceAciBgpBestPathPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciBgpBestPathPolicyImport,
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

			"ctrl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"asPathMultipathRelax", "0",
				}, false),
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}
func getRemoteBgpBestPathPolicy(client *client.Client, dn string) (*models.BgpBestPathPolicy, error) {
	bgpBestPathCtrlPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	bgpBestPathCtrlPol := models.BgpBestPathPolicyFromContainer(bgpBestPathCtrlPolCont)

	if bgpBestPathCtrlPol.DistinguishedName == "" {
		return nil, fmt.Errorf("BgpBestPathPolicy %s not found", bgpBestPathCtrlPol.DistinguishedName)
	}

	return bgpBestPathCtrlPol, nil
}

func setBgpBestPathPolicyAttributes(bgpBestPathCtrlPol *models.BgpBestPathPolicy, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(bgpBestPathCtrlPol.DistinguishedName)
	d.Set("description", bgpBestPathCtrlPol.Description)
	dn := d.Id()
	if dn != bgpBestPathCtrlPol.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	bgpBestPathCtrlPolMap, _ := bgpBestPathCtrlPol.ToMap()

	d.Set("name", bgpBestPathCtrlPolMap["name"])

	d.Set("annotation", bgpBestPathCtrlPolMap["annotation"])
	if bgpBestPathCtrlPolMap["ctrl"] == "" {
		d.Set("ctrl", "0")
	} else {
		d.Set("ctrl", bgpBestPathCtrlPolMap["ctrl"])
	}
	d.Set("name_alias", bgpBestPathCtrlPolMap["nameAlias"])
	return d
}

func resourceAciBgpBestPathPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	bgpBestPathCtrlPol, err := getRemoteBgpBestPathPolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setBgpBestPathPolicyAttributes(bgpBestPathCtrlPol, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciBgpBestPathPolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] BgpBestPathPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	bgpBestPathCtrlPolAttr := models.BgpBestPathPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		bgpBestPathCtrlPolAttr.Annotation = Annotation.(string)
	} else {
		bgpBestPathCtrlPolAttr.Annotation = "{}"
	}
	if Ctrl, ok := d.GetOk("ctrl"); ok {
		bgpBestPathCtrlPolAttr.Ctrl = Ctrl.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		bgpBestPathCtrlPolAttr.NameAlias = NameAlias.(string)
	}
	bgpBestPathCtrlPol := models.NewBgpBestPathPolicy(fmt.Sprintf("bestpath-%s", name), TenantDn, desc, bgpBestPathCtrlPolAttr)

	err := aciClient.Save(bgpBestPathCtrlPol)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(bgpBestPathCtrlPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciBgpBestPathPolicyRead(d, m)
}

func resourceAciBgpBestPathPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] BgpBestPathPolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	bgpBestPathCtrlPolAttr := models.BgpBestPathPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		bgpBestPathCtrlPolAttr.Annotation = Annotation.(string)
	} else {
		bgpBestPathCtrlPolAttr.Annotation = "{}"
	}
	if Ctrl, ok := d.GetOk("ctrl"); ok {
		bgpBestPathCtrlPolAttr.Ctrl = Ctrl.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		bgpBestPathCtrlPolAttr.NameAlias = NameAlias.(string)
	}
	bgpBestPathCtrlPol := models.NewBgpBestPathPolicy(fmt.Sprintf("bestpath-%s", name), TenantDn, desc, bgpBestPathCtrlPolAttr)

	bgpBestPathCtrlPol.Status = "modified"

	err := aciClient.Save(bgpBestPathCtrlPol)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(bgpBestPathCtrlPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciBgpBestPathPolicyRead(d, m)

}

func resourceAciBgpBestPathPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	bgpBestPathCtrlPol, err := getRemoteBgpBestPathPolicy(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}

	setBgpBestPathPolicyAttributes(bgpBestPathCtrlPol, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciBgpBestPathPolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "bgpBestPathCtrlPol")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
