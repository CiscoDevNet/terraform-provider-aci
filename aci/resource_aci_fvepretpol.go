package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAciEndPointRetentionPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciEndPointRetentionPolicyCreate,
		Update: resourceAciEndPointRetentionPolicyUpdate,
		Read:   resourceAciEndPointRetentionPolicyRead,
		Delete: resourceAciEndPointRetentionPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciEndPointRetentionPolicyImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"tenant_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"bounce_age_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"bounce_trig": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"hold_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"local_ep_age_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"move_freq": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"remote_ep_age_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}
func getRemoteEndPointRetentionPolicy(client *client.Client, dn string) (*models.EndPointRetentionPolicy, error) {
	fvEpRetPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	fvEpRetPol := models.EndPointRetentionPolicyFromContainer(fvEpRetPolCont)

	if fvEpRetPol.DistinguishedName == "" {
		return nil, fmt.Errorf("EndPointRetentionPolicy %s not found", fvEpRetPol.DistinguishedName)
	}

	return fvEpRetPol, nil
}

func setEndPointRetentionPolicyAttributes(fvEpRetPol *models.EndPointRetentionPolicy, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(fvEpRetPol.DistinguishedName)
	d.Set("description", fvEpRetPol.Description)
	d.Set("tenant_dn", GetParentDn(fvEpRetPol.DistinguishedName))
	fvEpRetPolMap, _ := fvEpRetPol.ToMap()

	d.Set("name", fvEpRetPolMap["name"])

	d.Set("annotation", fvEpRetPolMap["annotation"])
	d.Set("bounce_age_intvl", fvEpRetPolMap["bounceAgeIntvl"])
	d.Set("bounce_trig", fvEpRetPolMap["bounceTrig"])
	d.Set("hold_intvl", fvEpRetPolMap["holdIntvl"])
	d.Set("local_ep_age_intvl", fvEpRetPolMap["localEpAgeIntvl"])
	d.Set("move_freq", fvEpRetPolMap["moveFreq"])
	d.Set("name_alias", fvEpRetPolMap["nameAlias"])
	d.Set("remote_ep_age_intvl", fvEpRetPolMap["remoteEpAgeIntvl"])
	return d
}

func resourceAciEndPointRetentionPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	fvEpRetPol, err := getRemoteEndPointRetentionPolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setEndPointRetentionPolicyAttributes(fvEpRetPol, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciEndPointRetentionPolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] EndPointRetentionPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	fvEpRetPolAttr := models.EndPointRetentionPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvEpRetPolAttr.Annotation = Annotation.(string)
	}
	if BounceAgeIntvl, ok := d.GetOk("bounce_age_intvl"); ok {
		fvEpRetPolAttr.BounceAgeIntvl = BounceAgeIntvl.(string)
	}
	if BounceTrig, ok := d.GetOk("bounce_trig"); ok {
		fvEpRetPolAttr.BounceTrig = BounceTrig.(string)
	}
	if HoldIntvl, ok := d.GetOk("hold_intvl"); ok {
		fvEpRetPolAttr.HoldIntvl = HoldIntvl.(string)
	}
	if LocalEpAgeIntvl, ok := d.GetOk("local_ep_age_intvl"); ok {
		fvEpRetPolAttr.LocalEpAgeIntvl = LocalEpAgeIntvl.(string)
	}
	if MoveFreq, ok := d.GetOk("move_freq"); ok {
		fvEpRetPolAttr.MoveFreq = MoveFreq.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fvEpRetPolAttr.NameAlias = NameAlias.(string)
	}
	if RemoteEpAgeIntvl, ok := d.GetOk("remote_ep_age_intvl"); ok {
		fvEpRetPolAttr.RemoteEpAgeIntvl = RemoteEpAgeIntvl.(string)
	}
	fvEpRetPol := models.NewEndPointRetentionPolicy(fmt.Sprintf("epRPol-%s", name), TenantDn, desc, fvEpRetPolAttr)

	err := aciClient.Save(fvEpRetPol)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(fvEpRetPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciEndPointRetentionPolicyRead(d, m)
}

func resourceAciEndPointRetentionPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] EndPointRetentionPolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	fvEpRetPolAttr := models.EndPointRetentionPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvEpRetPolAttr.Annotation = Annotation.(string)
	}
	if BounceAgeIntvl, ok := d.GetOk("bounce_age_intvl"); ok {
		fvEpRetPolAttr.BounceAgeIntvl = BounceAgeIntvl.(string)
	}
	if BounceTrig, ok := d.GetOk("bounce_trig"); ok {
		fvEpRetPolAttr.BounceTrig = BounceTrig.(string)
	}
	if HoldIntvl, ok := d.GetOk("hold_intvl"); ok {
		fvEpRetPolAttr.HoldIntvl = HoldIntvl.(string)
	}
	if LocalEpAgeIntvl, ok := d.GetOk("local_ep_age_intvl"); ok {
		fvEpRetPolAttr.LocalEpAgeIntvl = LocalEpAgeIntvl.(string)
	}
	if MoveFreq, ok := d.GetOk("move_freq"); ok {
		fvEpRetPolAttr.MoveFreq = MoveFreq.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fvEpRetPolAttr.NameAlias = NameAlias.(string)
	}
	if RemoteEpAgeIntvl, ok := d.GetOk("remote_ep_age_intvl"); ok {
		fvEpRetPolAttr.RemoteEpAgeIntvl = RemoteEpAgeIntvl.(string)
	}
	fvEpRetPol := models.NewEndPointRetentionPolicy(fmt.Sprintf("epRPol-%s", name), TenantDn, desc, fvEpRetPolAttr)

	fvEpRetPol.Status = "modified"

	err := aciClient.Save(fvEpRetPol)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(fvEpRetPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciEndPointRetentionPolicyRead(d, m)

}

func resourceAciEndPointRetentionPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	fvEpRetPol, err := getRemoteEndPointRetentionPolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setEndPointRetentionPolicyAttributes(fvEpRetPol, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciEndPointRetentionPolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fvEpRetPol")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
