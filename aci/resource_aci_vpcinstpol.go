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

func resourceAciVPCDomainPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciVPCDomainPolicyCreate,
		UpdateContext: resourceAciVPCDomainPolicyUpdate,
		ReadContext:   resourceAciVPCDomainPolicyRead,
		DeleteContext: resourceAciVPCDomainPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciVPCDomainPolicyImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{

			"dead_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		})),
	}
}

func getRemoteVPCDomainPolicy(client *client.Client, dn string) (*models.VPCDomainPolicy, error) {
	vpcInstPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	vpcInstPol := models.VPCDomainPolicyFromContainer(vpcInstPolCont)
	if vpcInstPol.DistinguishedName == "" {
		return nil, fmt.Errorf("VPCDomainPolicy %s not found", vpcInstPol.DistinguishedName)
	}
	return vpcInstPol, nil
}

func setVPCDomainPolicyAttributes(vpcInstPol *models.VPCDomainPolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(vpcInstPol.DistinguishedName)
	d.Set("description", vpcInstPol.Description)
	vpcInstPolMap, err := vpcInstPol.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("annotation", vpcInstPolMap["annotation"])
	d.Set("dead_intvl", vpcInstPolMap["deadIntvl"])
	d.Set("name", vpcInstPolMap["name"])
	d.Set("name_alias", vpcInstPolMap["nameAlias"])
	return d, nil
}

func resourceAciVPCDomainPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	vpcInstPol, err := getRemoteVPCDomainPolicy(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setVPCDomainPolicyAttributes(vpcInstPol, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciVPCDomainPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] VPCDomainPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	vpcInstPolAttr := models.VPCDomainPolicyAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vpcInstPolAttr.Annotation = Annotation.(string)
	} else {
		vpcInstPolAttr.Annotation = "{}"
	}

	if DeadIntvl, ok := d.GetOk("dead_intvl"); ok {
		vpcInstPolAttr.DeadIntvl = DeadIntvl.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		vpcInstPolAttr.Name = Name.(string)
	}
	vpcInstPol := models.NewVPCDomainPolicy(fmt.Sprintf("fabric/vpcInst-%s", name), "uni", desc, nameAlias, vpcInstPolAttr)
	err := aciClient.Save(vpcInstPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(vpcInstPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciVPCDomainPolicyRead(ctx, d, m)
}

func resourceAciVPCDomainPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] VPCDomainPolicy: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	vpcInstPolAttr := models.VPCDomainPolicyAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		vpcInstPolAttr.Annotation = Annotation.(string)
	} else {
		vpcInstPolAttr.Annotation = "{}"
	}

	if DeadIntvl, ok := d.GetOk("dead_intvl"); ok {
		vpcInstPolAttr.DeadIntvl = DeadIntvl.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		vpcInstPolAttr.Name = Name.(string)
	}
	vpcInstPol := models.NewVPCDomainPolicy(fmt.Sprintf("fabric/vpcInst-%s", name), "uni", desc, nameAlias, vpcInstPolAttr)
	vpcInstPol.Status = "modified"
	err := aciClient.Save(vpcInstPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(vpcInstPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciVPCDomainPolicyRead(ctx, d, m)
}

func resourceAciVPCDomainPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	vpcInstPol, err := getRemoteVPCDomainPolicy(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}
	setVPCDomainPolicyAttributes(vpcInstPol, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciVPCDomainPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "vpcInstPol")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
