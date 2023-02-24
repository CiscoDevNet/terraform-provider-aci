package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciL3InterfacePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciL3InterfacePolicyCreate,
		UpdateContext: resourceAciL3InterfacePolicyUpdate,
		ReadContext:   resourceAciL3InterfacePolicyRead,
		DeleteContext: resourceAciL3InterfacePolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciL3InterfacePolicyImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"bfd_isis": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"disabled",
					"enabled",
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
func getRemoteL3InterfacePolicy(client *client.Client, dn string) (*models.L3InterfacePolicy, error) {
	l3IfPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	l3IfPol := models.L3InterfacePolicyFromContainer(l3IfPolCont)

	if l3IfPol.DistinguishedName == "" {
		return nil, fmt.Errorf("L3 Interface Policy %s not found", dn)
	}

	return l3IfPol, nil
}

func setL3InterfacePolicyAttributes(l3IfPol *models.L3InterfacePolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(l3IfPol.DistinguishedName)
	d.Set("description", l3IfPol.Description)
	l3IfPolMap, err := l3IfPol.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("name", l3IfPolMap["name"])
	d.Set("annotation", l3IfPolMap["annotation"])
	d.Set("bfd_isis", l3IfPolMap["bfdIsis"])
	d.Set("name_alias", l3IfPolMap["nameAlias"])
	return d, nil
}

func resourceAciL3InterfacePolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	l3IfPol, err := getRemoteL3InterfacePolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setL3InterfacePolicyAttributes(l3IfPol, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciL3InterfacePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] L3InterfacePolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	l3IfPolAttr := models.L3InterfacePolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3IfPolAttr.Annotation = Annotation.(string)
	} else {
		l3IfPolAttr.Annotation = "{}"
	}
	if BfdIsis, ok := d.GetOk("bfd_isis"); ok {
		l3IfPolAttr.BfdIsis = BfdIsis.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l3IfPolAttr.NameAlias = NameAlias.(string)
	}
	l3IfPol := models.NewL3InterfacePolicy(fmt.Sprintf("fabric/l3IfP-%s", name), "uni", desc, l3IfPolAttr)

	err := aciClient.Save(l3IfPol)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(l3IfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciL3InterfacePolicyRead(ctx, d, m)
}

func resourceAciL3InterfacePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] L3InterfacePolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	l3IfPolAttr := models.L3InterfacePolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3IfPolAttr.Annotation = Annotation.(string)
	} else {
		l3IfPolAttr.Annotation = "{}"
	}
	if BfdIsis, ok := d.GetOk("bfd_isis"); ok {
		l3IfPolAttr.BfdIsis = BfdIsis.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l3IfPolAttr.NameAlias = NameAlias.(string)
	}
	l3IfPol := models.NewL3InterfacePolicy(fmt.Sprintf("fabric/l3IfP-%s", name), "uni", desc, l3IfPolAttr)

	l3IfPol.Status = "modified"

	err := aciClient.Save(l3IfPol)

	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(l3IfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciL3InterfacePolicyRead(ctx, d, m)

}

func resourceAciL3InterfacePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	l3IfPol, err := getRemoteL3InterfacePolicy(aciClient, dn)

	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}
	setL3InterfacePolicyAttributes(l3IfPol, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciL3InterfacePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "l3IfPol")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
