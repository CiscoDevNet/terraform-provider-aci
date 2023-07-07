package aci

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciPIMInterfacePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciPIMInterfacePolicyCreate,
		UpdateContext: resourceAciPIMInterfacePolicyUpdate,
		ReadContext:   resourceAciPIMInterfacePolicyRead,
		DeleteContext: resourceAciPIMInterfacePolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciPIMInterfacePolicyImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"tenant_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"auth_key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"auth_t": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ah-md5",
					"none",
				}, false),
			},
			"ctrl": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"border",
						"passive",
						"strict-rfc-compliant",
					}, false),
				},
			},
			"dr_delay": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dr_prio": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"hello_itvl": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"jp_interval": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"secure_auth_key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func getRemotePIMInterfacePolicy(client *client.Client, dn string) (*models.PIMInterfacePolicy, error) {
	pimIfPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	pimIfPol := models.PIMInterfacePolicyFromContainer(pimIfPolCont)
	if pimIfPol.DistinguishedName == "" {
		return nil, fmt.Errorf("PIM Interface Policy %s not found", dn)
	}
	return pimIfPol, nil
}

func setPIMInterfacePolicyAttributes(pimIfPol *models.PIMInterfacePolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(pimIfPol.DistinguishedName)
	d.Set("description", pimIfPol.Description)
	pimIfPolMap, err := pimIfPol.ToMap()
	if err != nil {
		return d, err
	}
	dn := d.Id()
	if dn != pimIfPol.DistinguishedName {
		d.Set("tenant_dn", "")
	} else {
		d.Set("tenant_dn", GetParentDn(pimIfPol.DistinguishedName, fmt.Sprintf("/"+models.RnPimIfPol, pimIfPolMap["name"])))
	}
	d.Set("annotation", pimIfPolMap["annotation"])
	d.Set("auth_key", pimIfPolMap["authKey"])
	d.Set("auth_t", pimIfPolMap["authT"])
	ctrlGet := make([]string, 0, 1)
	for _, val := range strings.Split(pimIfPolMap["ctrl"], ",") {
		ctrlGet = append(ctrlGet, strings.Trim(val, " "))
	}
	sort.Strings(ctrlGet)
	if ctrlIntr, ok := d.GetOk("ctrl"); ok {
		ctrlAct := make([]string, 0, 1)
		for _, val := range ctrlIntr.([]interface{}) {
			ctrlAct = append(ctrlAct, val.(string))
		}
		sort.Strings(ctrlAct)
		if reflect.DeepEqual(ctrlAct, ctrlGet) {
			d.Set("ctrl", d.Get("ctrl").([]interface{}))
		} else {
			d.Set("ctrl", ctrlGet)
		}
	} else {
		d.Set("ctrl", ctrlGet)
	}
	d.Set("dr_delay", pimIfPolMap["drDelay"])
	d.Set("dr_prio", pimIfPolMap["drPrio"])
	d.Set("hello_itvl", pimIfPolMap["helloItvl"])
	d.Set("jp_interval", pimIfPolMap["jpInterval"])
	d.Set("name", pimIfPolMap["name"])
	d.Set("name_alias", pimIfPolMap["nameAlias"])
	d.Set("secure_auth_key", pimIfPolMap["secureAuthKey"])
	return d, nil
}

func resourceAciPIMInterfacePolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	pimIfPol, err := getRemotePIMInterfacePolicy(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setPIMInterfacePolicyAttributes(pimIfPol, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciPIMInterfacePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] PIMInterfacePolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	TenantDn := d.Get("tenant_dn").(string)

	pimIfPolAttr := models.PIMInterfacePolicyAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		pimIfPolAttr.Annotation = Annotation.(string)
	} else {
		pimIfPolAttr.Annotation = "{}"
	}

	if AuthKey, ok := d.GetOk("auth_key"); ok {
		pimIfPolAttr.AuthKey = AuthKey.(string)
	}

	if AuthT, ok := d.GetOk("auth_t"); ok {
		pimIfPolAttr.AuthT = AuthT.(string)
	}

	if Ctrl, ok := d.GetOk("ctrl"); ok {
		ctrlList := make([]string, 0, 1)
		for _, val := range Ctrl.([]interface{}) {
			ctrlList = append(ctrlList, val.(string))
		}
		Ctrl := strings.Join(ctrlList, ",")
		pimIfPolAttr.Ctrl = Ctrl
	}

	if DrDelay, ok := d.GetOk("dr_delay"); ok {
		pimIfPolAttr.DrDelay = DrDelay.(string)
	}

	if DrPrio, ok := d.GetOk("dr_prio"); ok {
		pimIfPolAttr.DrPrio = DrPrio.(string)
	}

	if HelloItvl, ok := d.GetOk("hello_itvl"); ok {
		pimIfPolAttr.HelloItvl = HelloItvl.(string)
	}

	if JpInterval, ok := d.GetOk("jp_interval"); ok {
		pimIfPolAttr.JpInterval = JpInterval.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		pimIfPolAttr.Name = Name.(string)
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		pimIfPolAttr.NameAlias = NameAlias.(string)
	}

	if SecureAuthKey, ok := d.GetOk("secure_auth_key"); ok {
		pimIfPolAttr.SecureAuthKey = SecureAuthKey.(string)
	}
	pimIfPol := models.NewPIMInterfacePolicy(fmt.Sprintf(models.RnPimIfPol, name), TenantDn, desc, pimIfPolAttr)

	err := aciClient.Save(pimIfPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pimIfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciPIMInterfacePolicyRead(ctx, d, m)
}
func resourceAciPIMInterfacePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] PIMInterfacePolicy: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	TenantDn := d.Get("tenant_dn").(string)

	pimIfPolAttr := models.PIMInterfacePolicyAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		pimIfPolAttr.Annotation = Annotation.(string)
	} else {
		pimIfPolAttr.Annotation = "{}"
	}

	if AuthKey, ok := d.GetOk("auth_key"); ok {
		pimIfPolAttr.AuthKey = AuthKey.(string)
	}

	if AuthT, ok := d.GetOk("auth_t"); ok {
		pimIfPolAttr.AuthT = AuthT.(string)
	}
	if Ctrl, ok := d.GetOk("ctrl"); ok {
		ctrlList := make([]string, 0, 1)
		for _, val := range Ctrl.([]interface{}) {
			ctrlList = append(ctrlList, val.(string))
		}
		Ctrl := strings.Join(ctrlList, ",")
		pimIfPolAttr.Ctrl = Ctrl
	}

	if DrDelay, ok := d.GetOk("dr_delay"); ok {
		pimIfPolAttr.DrDelay = DrDelay.(string)
	}

	if DrPrio, ok := d.GetOk("dr_prio"); ok {
		pimIfPolAttr.DrPrio = DrPrio.(string)
	}

	if HelloItvl, ok := d.GetOk("hello_itvl"); ok {
		pimIfPolAttr.HelloItvl = HelloItvl.(string)
	}

	if JpInterval, ok := d.GetOk("jp_interval"); ok {
		pimIfPolAttr.JpInterval = JpInterval.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		pimIfPolAttr.Name = Name.(string)
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		pimIfPolAttr.NameAlias = NameAlias.(string)
	}

	if SecureAuthKey, ok := d.GetOk("secure_auth_key"); ok {
		pimIfPolAttr.SecureAuthKey = SecureAuthKey.(string)
	}
	pimIfPol := models.NewPIMInterfacePolicy(fmt.Sprintf(models.RnPimIfPol, name), TenantDn, desc, pimIfPolAttr)

	pimIfPol.Status = "modified"

	err := aciClient.Save(pimIfPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pimIfPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciPIMInterfacePolicyRead(ctx, d, m)
}

func resourceAciPIMInterfacePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	pimIfPol, err := getRemotePIMInterfacePolicy(aciClient, dn)
	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}

	_, err = setPIMInterfacePolicyAttributes(pimIfPol, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciPIMInterfacePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, models.PimIfPolClassName)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
