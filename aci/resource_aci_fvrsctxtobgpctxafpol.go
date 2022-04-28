package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciBGPAddressFamilyContextPolicyRelationship() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciBGPAddressFamilyContextPolicyRelationshipCreate,
		UpdateContext: resourceAciBGPAddressFamilyContextPolicyRelationshipUpdate,
		ReadContext:   resourceAciBGPAddressFamilyContextPolicyRelationshipRead,
		DeleteContext: resourceAciBGPAddressFamilyContextPolicyRelationshipDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciBGPAddressFamilyContextPolicyRelationshipImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"vrf_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"address_family": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ipv4-ucast",
					"ipv6-ucast",
				}, false),
			},
			"bgp_address_family_context_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		})),
	}
}

func getRemoteBGPAddressFamilyContextPolicyRelationship(client *client.Client, dn string) (*models.BGPAddressFamilyContextPolicyRelationship, error) {
	fvRsCtxToBgpCtxAfPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	fvRsCtxToBgpCtxAfPol := models.BGPAddressFamilyContextPolicyRelationshipFromContainer(fvRsCtxToBgpCtxAfPolCont)
	if fvRsCtxToBgpCtxAfPol.DistinguishedName == "" {
		return nil, fmt.Errorf("BGPAddressFamilyContextPolicyRelationship %s not found", fvRsCtxToBgpCtxAfPol.DistinguishedName)
	}
	return fvRsCtxToBgpCtxAfPol, nil
}

func setBGPAddressFamilyContextPolicyRelationshipAttributes(fvRsCtxToBgpCtxAfPol *models.BGPAddressFamilyContextPolicyRelationship, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(fvRsCtxToBgpCtxAfPol.DistinguishedName)
	if dn != fvRsCtxToBgpCtxAfPol.DistinguishedName {
		d.Set("vrf_dn", "")
	}

	fvRsCtxToBgpCtxAfPolMap, err := fvRsCtxToBgpCtxAfPol.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("address_family", fvRsCtxToBgpCtxAfPolMap["af"])
	d.Set("annotation", fvRsCtxToBgpCtxAfPolMap["annotation"])
	d.Set("bgp_address_family_context_dn", fvRsCtxToBgpCtxAfPolMap["tDn"])
	d.Set("vrf_dn", GetParentDn(dn, fmt.Sprintf("/"+models.RnfvRsCtxToBgpCtxAfPol, fvRsCtxToBgpCtxAfPolMap["tnBgpCtxAfPolName"], fvRsCtxToBgpCtxAfPolMap["af"])))
	return d, nil
}

func resourceAciBGPAddressFamilyContextPolicyRelationshipImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	fvRsCtxToBgpCtxAfPol, err := getRemoteBGPAddressFamilyContextPolicyRelationship(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setBGPAddressFamilyContextPolicyRelationshipAttributes(fvRsCtxToBgpCtxAfPol, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciBGPAddressFamilyContextPolicyRelationshipCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] BGPAddressFamilyContextPolicyRelationship: Beginning Creation")
	aciClient := m.(*client.Client)
	tnBgpCtxAfPolName := GetMOName(d.Get("bgp_address_family_context_dn").(string))
	af := d.Get("address_family").(string)
	VRFDn := d.Get("vrf_dn").(string)

	fvRsCtxToBgpCtxAfPolAttr := models.BGPAddressFamilyContextPolicyRelationshipAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		fvRsCtxToBgpCtxAfPolAttr.Annotation = Annotation.(string)
	} else {
		fvRsCtxToBgpCtxAfPolAttr.Annotation = "{}"
	}

	fvRsCtxToBgpCtxAfPolAttr.Af = af
	fvRsCtxToBgpCtxAfPolAttr.TnBgpCtxAfPolName = tnBgpCtxAfPolName

	fvRsCtxToBgpCtxAfPol := models.NewBGPAddressFamilyContextPolicyRelationship(fmt.Sprintf(models.RnfvRsCtxToBgpCtxAfPol, tnBgpCtxAfPolName, af), VRFDn, fvRsCtxToBgpCtxAfPolAttr)

	err := aciClient.Save(fvRsCtxToBgpCtxAfPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fvRsCtxToBgpCtxAfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciBGPAddressFamilyContextPolicyRelationshipRead(ctx, d, m)
}

func resourceAciBGPAddressFamilyContextPolicyRelationshipUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] BGPAddressFamilyContextPolicyRelationship: Beginning Update")
	aciClient := m.(*client.Client)
	tnBgpCtxAfPolName := GetMOName(d.Get("bgp_address_family_context_dn").(string))
	af := d.Get("address_family").(string)
	VRFDn := d.Get("vrf_dn").(string)

	fvRsCtxToBgpCtxAfPolAttr := models.BGPAddressFamilyContextPolicyRelationshipAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		fvRsCtxToBgpCtxAfPolAttr.Annotation = Annotation.(string)
	} else {
		fvRsCtxToBgpCtxAfPolAttr.Annotation = "{}"
	}

	fvRsCtxToBgpCtxAfPolAttr.Af = af
	fvRsCtxToBgpCtxAfPolAttr.TnBgpCtxAfPolName = tnBgpCtxAfPolName

	fvRsCtxToBgpCtxAfPol := models.NewBGPAddressFamilyContextPolicyRelationship(fmt.Sprintf(models.RnfvRsCtxToBgpCtxAfPol, tnBgpCtxAfPolName, af), VRFDn, fvRsCtxToBgpCtxAfPolAttr)

	fvRsCtxToBgpCtxAfPol.Status = "modified"

	err := aciClient.Save(fvRsCtxToBgpCtxAfPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fvRsCtxToBgpCtxAfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciBGPAddressFamilyContextPolicyRelationshipRead(ctx, d, m)
}

func resourceAciBGPAddressFamilyContextPolicyRelationshipRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	fvRsCtxToBgpCtxAfPol, err := getRemoteBGPAddressFamilyContextPolicyRelationship(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}

	_, err = setBGPAddressFamilyContextPolicyRelationshipAttributes(fvRsCtxToBgpCtxAfPol, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciBGPAddressFamilyContextPolicyRelationshipDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "fvRsCtxToBgpCtxAfPol")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
