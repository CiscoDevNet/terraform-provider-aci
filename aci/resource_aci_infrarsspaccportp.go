package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAciInterfaceProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciInterfaceProfileCreate,
		UpdateContext: resourceAciInterfaceProfileUpdate,
		ReadContext:   resourceAciInterfaceProfileRead,
		DeleteContext: resourceAciInterfaceProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciInterfaceProfileImport,
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"spine_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				// Default:  "orchestrator:terraform",
				Computed: true,
				DefaultFunc: func() (interface{}, error) {
					return "orchestrator:terraform", nil
				},
			},

			"tdn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}
func getRemoteInterfaceProfile(client *client.Client, dn string) (*models.InterfaceProfile, error) {
	infraRsSpAccPortPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	infraRsSpAccPortP := models.InterfaceProfileFromContainer(infraRsSpAccPortPCont)

	if infraRsSpAccPortP.DistinguishedName == "" {
		return nil, fmt.Errorf("InterfaceProfile %s not found", infraRsSpAccPortP.DistinguishedName)
	}

	return infraRsSpAccPortP, nil
}

func setInterfaceProfileAttributes(infraRsSpAccPortP *models.InterfaceProfile, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(infraRsSpAccPortP.DistinguishedName)
	if dn != infraRsSpAccPortP.DistinguishedName {
		d.Set("spine_profile_dn", "")
	}

	infraRsSpAccPortPMap, err := infraRsSpAccPortP.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("tdn", infraRsSpAccPortPMap["tDn"])
	d.Set("annotation", infraRsSpAccPortPMap["annotation"])

	return d, nil
}

func resourceAciInterfaceProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	infraRsSpAccPortP, err := getRemoteInterfaceProfile(aciClient, dn)

	if err != nil {
		return nil, err
	}
	infraRsSpAccPortPMap, err := infraRsSpAccPortP.ToMap()
	if err != nil {
		return nil, err
	}
	tDn := infraRsSpAccPortPMap["tDn"]
	pDN := GetParentDn(dn, fmt.Sprintf("/rsspAccPortP-[%s]", tDn))
	d.Set("spine_profile_dn", pDN)
	schemaFilled, err := setInterfaceProfileAttributes(infraRsSpAccPortP, d)

	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciInterfaceProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] InterfaceProfile: Beginning Creation")
	aciClient := m.(*client.Client)

	tDn := d.Get("tdn").(string)

	SpineProfileDn := d.Get("spine_profile_dn").(string)

	infraRsSpAccPortPAttr := models.InterfaceProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraRsSpAccPortPAttr.Annotation = Annotation.(string)
	} else {
		infraRsSpAccPortPAttr.Annotation = "{}"
	}

	infraRsSpAccPortP := models.NewInterfaceProfile(fmt.Sprintf("rsspAccPortP-[%s]", tDn), SpineProfileDn, "", infraRsSpAccPortPAttr)

	err := aciClient.Save(infraRsSpAccPortP)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(infraRsSpAccPortP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciInterfaceProfileRead(ctx, d, m)
}

func resourceAciInterfaceProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] InterfaceProfile: Beginning Update")

	aciClient := m.(*client.Client)

	tDn := d.Get("tdn").(string)

	SpineProfileDn := d.Get("spine_profile_dn").(string)

	infraRsSpAccPortPAttr := models.InterfaceProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraRsSpAccPortPAttr.Annotation = Annotation.(string)
	} else {
		infraRsSpAccPortPAttr.Annotation = "{}"
	}

	infraRsSpAccPortP := models.NewInterfaceProfile(fmt.Sprintf("rsspAccPortP-[%s]", tDn), SpineProfileDn, "", infraRsSpAccPortPAttr)

	infraRsSpAccPortP.Status = "modified"

	err := aciClient.Save(infraRsSpAccPortP)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(infraRsSpAccPortP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciInterfaceProfileRead(ctx, d, m)

}

func resourceAciInterfaceProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	infraRsSpAccPortP, err := getRemoteInterfaceProfile(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setInterfaceProfileAttributes(infraRsSpAccPortP, d)

	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciInterfaceProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "infraRsSpAccPortP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
