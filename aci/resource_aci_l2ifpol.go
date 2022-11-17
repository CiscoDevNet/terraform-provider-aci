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

func resourceAciL2InterfacePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciL2InterfacePolicyCreate,
		UpdateContext: resourceAciL2InterfacePolicyUpdate,
		ReadContext:   resourceAciL2InterfacePolicyRead,
		DeleteContext: resourceAciL2InterfacePolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciL2InterfacePolicyImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

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

			"qinq": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"disabled",
					"edgePort",
					"corePort",
					"doubleQtagPort",
				}, false),
			},

			"vepa": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"disabled",
					"enabled",
				}, false),
			},

			"vlan_scope": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"global",
					"portlocal",
				}, false),
			},
		}),
	}
}
func getRemoteL2InterfacePolicy(client *client.Client, dn string) (*models.L2InterfacePolicy, error) {
	l2IfPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	l2IfPol := models.L2InterfacePolicyFromContainer(l2IfPolCont)

	if l2IfPol.DistinguishedName == "" {
		return nil, fmt.Errorf("L2InterfacePolicy %s not found", l2IfPol.DistinguishedName)
	}

	return l2IfPol, nil
}

func setL2InterfacePolicyAttributes(l2IfPol *models.L2InterfacePolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(l2IfPol.DistinguishedName)
	d.Set("description", l2IfPol.Description)
	l2IfPolMap, err := l2IfPol.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("name", l2IfPolMap["name"])

	d.Set("annotation", l2IfPolMap["annotation"])
	d.Set("name_alias", l2IfPolMap["nameAlias"])
	d.Set("qinq", l2IfPolMap["qinq"])
	d.Set("vepa", l2IfPolMap["vepa"])
	d.Set("vlan_scope", l2IfPolMap["vlanScope"])
	return d, nil
}

func resourceAciL2InterfacePolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	l2IfPol, err := getRemoteL2InterfacePolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setL2InterfacePolicyAttributes(l2IfPol, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciL2InterfacePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] L2InterfacePolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	l2IfPolAttr := models.L2InterfacePolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l2IfPolAttr.Annotation = Annotation.(string)
	} else {
		l2IfPolAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l2IfPolAttr.NameAlias = NameAlias.(string)
	}
	if Qinq, ok := d.GetOk("qinq"); ok {
		l2IfPolAttr.Qinq = Qinq.(string)
	}
	if Vepa, ok := d.GetOk("vepa"); ok {
		l2IfPolAttr.Vepa = Vepa.(string)
	}
	if VlanScope, ok := d.GetOk("vlan_scope"); ok {
		l2IfPolAttr.VlanScope = VlanScope.(string)
	}
	l2IfPol := models.NewL2InterfacePolicy(fmt.Sprintf("infra/l2IfP-%s", name), "uni", desc, l2IfPolAttr)

	err := aciClient.Save(l2IfPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(l2IfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciL2InterfacePolicyRead(ctx, d, m)
}

func resourceAciL2InterfacePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] L2InterfacePolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	l2IfPolAttr := models.L2InterfacePolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l2IfPolAttr.Annotation = Annotation.(string)
	} else {
		l2IfPolAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l2IfPolAttr.NameAlias = NameAlias.(string)
	}
	if Qinq, ok := d.GetOk("qinq"); ok {
		l2IfPolAttr.Qinq = Qinq.(string)
	}
	if Vepa, ok := d.GetOk("vepa"); ok {
		l2IfPolAttr.Vepa = Vepa.(string)
	}
	if VlanScope, ok := d.GetOk("vlan_scope"); ok {
		l2IfPolAttr.VlanScope = VlanScope.(string)
	}
	l2IfPol := models.NewL2InterfacePolicy(fmt.Sprintf("infra/l2IfP-%s", name), "uni", desc, l2IfPolAttr)

	l2IfPol.Status = "modified"

	err := aciClient.Save(l2IfPol)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(l2IfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciL2InterfacePolicyRead(ctx, d, m)

}

func resourceAciL2InterfacePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	l2IfPol, err := getRemoteL2InterfacePolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setL2InterfacePolicyAttributes(l2IfPol, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciL2InterfacePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "l2IfPol")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
