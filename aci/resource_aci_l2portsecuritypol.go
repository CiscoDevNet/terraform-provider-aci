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

func resourceAciPortSecurityPolicy() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "The resource 'aci_port_security_policy' is deprecated, please refer to 'aci_port_security_interface_policy' instead. The resource will be removed in the next major version of the provider.",

		CreateContext: resourceAciPortSecurityPolicyCreate,
		UpdateContext: resourceAciPortSecurityPolicyUpdate,
		ReadContext:   resourceAciPortSecurityPolicyRead,
		DeleteContext: resourceAciPortSecurityPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciPortSecurityPolicyImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"maximum": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"timeout": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"violation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"protect",
				}, false),
			},
		}),
	}
}
func getRemotePortSecurityPolicy(client *client.Client, dn string) (*models.PortSecurityPolicy, error) {
	l2PortSecurityPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	l2PortSecurityPol := models.PortSecurityPolicyFromContainer(l2PortSecurityPolCont)

	if l2PortSecurityPol.DistinguishedName == "" {
		return nil, fmt.Errorf("Port Security Policy %s not found", dn)
	}

	return l2PortSecurityPol, nil
}

func setPortSecurityPolicyAttributes(l2PortSecurityPol *models.PortSecurityPolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(l2PortSecurityPol.DistinguishedName)
	d.Set("description", l2PortSecurityPol.Description)
	l2PortSecurityPolMap, err := l2PortSecurityPol.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("name", l2PortSecurityPolMap["name"])

	d.Set("annotation", l2PortSecurityPolMap["annotation"])
	d.Set("maximum", l2PortSecurityPolMap["maximum"])
	d.Set("name_alias", l2PortSecurityPolMap["nameAlias"])
	d.Set("timeout", l2PortSecurityPolMap["timeout"])
	d.Set("violation", l2PortSecurityPolMap["violation"])
	return d, nil
}

func resourceAciPortSecurityPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	l2PortSecurityPol, err := getRemotePortSecurityPolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setPortSecurityPolicyAttributes(l2PortSecurityPol, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciPortSecurityPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] PortSecurityPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	l2PortSecurityPolAttr := models.PortSecurityPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l2PortSecurityPolAttr.Annotation = Annotation.(string)
	} else {
		l2PortSecurityPolAttr.Annotation = "{}"
	}
	if Maximum, ok := d.GetOk("maximum"); ok {
		l2PortSecurityPolAttr.Maximum = Maximum.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l2PortSecurityPolAttr.NameAlias = NameAlias.(string)
	}
	if Timeout, ok := d.GetOk("timeout"); ok {
		l2PortSecurityPolAttr.Timeout = Timeout.(string)
	}
	if Violation, ok := d.GetOk("violation"); ok {
		l2PortSecurityPolAttr.Violation = Violation.(string)
	}
	l2PortSecurityPol := models.NewPortSecurityPolicy(fmt.Sprintf("infra/portsecurityP-%s", name), "uni", desc, l2PortSecurityPolAttr)

	err := aciClient.Save(l2PortSecurityPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(l2PortSecurityPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciPortSecurityPolicyRead(ctx, d, m)
}

func resourceAciPortSecurityPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] PortSecurityPolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	l2PortSecurityPolAttr := models.PortSecurityPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l2PortSecurityPolAttr.Annotation = Annotation.(string)
	} else {
		l2PortSecurityPolAttr.Annotation = "{}"
	}
	if Maximum, ok := d.GetOk("maximum"); ok {
		l2PortSecurityPolAttr.Maximum = Maximum.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		l2PortSecurityPolAttr.NameAlias = NameAlias.(string)
	}
	if Timeout, ok := d.GetOk("timeout"); ok {
		l2PortSecurityPolAttr.Timeout = Timeout.(string)
	}
	if Violation, ok := d.GetOk("violation"); ok {
		l2PortSecurityPolAttr.Violation = Violation.(string)
	}
	l2PortSecurityPol := models.NewPortSecurityPolicy(fmt.Sprintf("infra/portsecurityP-%s", name), "uni", desc, l2PortSecurityPolAttr)

	l2PortSecurityPol.Status = "modified"

	err := aciClient.Save(l2PortSecurityPol)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(l2PortSecurityPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciPortSecurityPolicyRead(ctx, d, m)

}

func resourceAciPortSecurityPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	l2PortSecurityPol, err := getRemotePortSecurityPolicy(aciClient, dn)

	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}
	_, err = setPortSecurityPolicyAttributes(l2PortSecurityPol, d)
	if err != nil {
		d.SetId("")
		return nil
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciPortSecurityPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "l2PortSecurityPol")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
