package aci

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciInterfaceConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciInterfaceConfigurationCreate,
		UpdateContext: resourceAciInterfaceConfigurationUpdate,
		ReadContext:   resourceAciInterfaceConfigurationRead,
		DeleteContext: resourceAciInterfaceConfigurationDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciInterfaceConfigurationImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"node": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(101, 4000),
			},
			"interface": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				StateFunc:        getInterfaceVal,
				ValidateDiagFunc: validateInterface,
			},
			"port_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"access",
					"fabric",
				}, false),
				Default: "access",
			},
			"role": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"leaf",
					"spine",
				}, false),
			},
			"policy_group": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"breakout": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"100g-2x",
					"100g-4x",
					"10g-4x",
					"25g-4x",
					"50g-8x",
					"none",
				}, false),
			},
			"admin_state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"up",
					"down",
				}, false),
			},
			"pc_member": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		})),
	}
}

func getInterfaceVal(interfaceVal interface{}) string {
	interfaceParts := strings.Split(interfaceVal.(string), "/")
	if len(interfaceParts) == 2 {
		interfaceParts = append(interfaceParts, "0")
		return strings.Join(interfaceParts, "/")
	}
	return interfaceVal.(string)
}

func getAdminState(adminState string) string {
	return map[string]string{
		"up":   "no",   // To handle user input
		"down": "yes",  // To handle user input
		"no":   "up",   // To handle APIC response
		"yes":  "down", // To handle APIC response
	}[adminState]
}

func validateInterface(value interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics
	errors := make([]map[string]string, 0)
	interfaceVal, validInterface := value.(string)

	invalidInterfaceError := map[string]string{
		"Summary": fmt.Sprintf("Interface: %v is invalid", interfaceVal),
		"Detail":  "The format must be either card/port/sub_port(1/1/1) or card/port(1/1).",
	}

	if validInterface {
		var err error
		var card, port, subPort int
		interfaceParts := strings.Split(interfaceVal, "/")

		if len(interfaceParts) == 2 || len(interfaceParts) == 3 {
			// Card ID validation
			if card, err = strconv.Atoi(interfaceParts[0]); err != nil || !InBetween(card, 1, 64) {
				errors = append(errors, map[string]string{
					"Summary": fmt.Sprintf("Card ID: %v is invalid", interfaceParts[0]),
					"Detail":  "Card ID must be in the range of 1 to 64.",
				})
			}
			// Port ID validation
			if port, err = strconv.Atoi(interfaceParts[1]); err != nil || !InBetween(port, 1, 128) {
				errors = append(errors, map[string]string{
					"Summary": fmt.Sprintf("Port ID: %v is invalid", interfaceParts[1]),
					"Detail":  "Port ID must be in the range of 1 to 128.",
				})
			}
			// Sub Port ID validation
			if len(interfaceParts) == 3 {
				if subPort, err = strconv.Atoi(interfaceParts[2]); err != nil || !InBetween(subPort, 0, 16) {
					errors = append(errors, map[string]string{
						"Summary": fmt.Sprintf("Sub Port ID: %v is invalid", interfaceParts[2]),
						"Detail":  "Sub Port ID must be in the range of 0 to 16.",
					})
				}
			}
		} else {
			// If the interfaceParts length is less than 2 or greater than 3, it returns error.
			errors = append(errors, invalidInterfaceError)
		}
	} else {
		// If not interface, it returns error.
		errors = append(errors, invalidInterfaceError)
	}
	for _, error := range errors {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       error["Summary"],
			Detail:        error["Detail"],
			AttributePath: path,
		})
	}
	return diags
}

func getAndSetRemoteInterfaceConfiguration(dn string, client *client.Client, d *schema.ResourceData) error {
	interfaceConfigCont, err := client.Get(dn)
	if err != nil {
		return err
	}

	var interfaceConfigMap map[string]string
	portType := d.Get("port_type").(string)

	if portType == "access" {
		interfaceConfig := models.InfraPortConfigurationFromContainer(interfaceConfigCont)
		interfaceConfigMap, err = interfaceConfig.ToMap()

		// To match the terraform plan stage changes with config file
		if interfaceConfigMap["brkoutMap"] == "none" && d.Get("breakout").(string) == "" {
			d.Set("breakout", nil)
		} else {
			d.Set("breakout", interfaceConfigMap["brkoutMap"])
		}

		d.Set("pc_member", interfaceConfigMap["pcMember"])
	} else if portType == "fabric" {
		interfaceConfig := models.FabricPortConfigurationFromContainer(interfaceConfigCont)
		interfaceConfigMap, err = interfaceConfig.ToMap()
	}

	if err != nil {
		return err
	}

	if interfaceConfigMap["dn"] == "" {
		d.Set("breakout", nil)
		d.Set("pc_member", nil)
		return fmt.Errorf("Interface Configuration: %s not found", dn)
	}

	if nodeId, err := strconv.Atoi(interfaceConfigMap["node"]); err == nil {
		d.Set("node", nodeId)
	} else {
		return err
	}

	d.Set("annotation", interfaceConfigMap["annotation"])
	d.Set("interface", fmt.Sprintf("%s/%s/%s", interfaceConfigMap["card"], interfaceConfigMap["port"], interfaceConfigMap["subPort"]))
	d.Set("role", interfaceConfigMap["role"])
	d.Set("policy_group", interfaceConfigMap["assocGrp"])
	d.Set("admin_state", getAdminState(interfaceConfigMap["shutdown"]))
	d.Set("description", interfaceConfigMap["description"])
	return nil
}

func resourceAciInterfaceConfigurationImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	if match, _ := regexp.MatchString("/infra/", dn); match {
		d.Set("port_type", "access")
	} else if match, _ := regexp.MatchString("/fabric/", dn); match {
		d.Set("port_type", "fabric")
	} else {
		return nil, fmt.Errorf("Interface Configuration DN: %s is not valid", dn)
	}
	err := getAndSetRemoteInterfaceConfiguration(dn, aciClient, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func resourceAciInterfaceConfigurationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Interface Configuration: Beginning Creation")
	aciClient := m.(*client.Client)

	node := strconv.Itoa(d.Get("node").(int))
	interfaceAttrMap := make(map[string]string)
	portType := d.Get("port_type").(string)
	parsedInterface := strings.Split(getInterfaceVal(d.Get("interface")), "/")

	if Annotation, ok := d.GetOk("annotation"); ok {
		interfaceAttrMap["Annotation"] = Annotation.(string)
	} else {
		interfaceAttrMap["Annotation"] = "{}"
	}

	interfaceAttrMap["Card"] = parsedInterface[0]
	interfaceAttrMap["Descr"] = d.Get("description").(string)
	interfaceAttrMap["Node"] = node
	interfaceAttrMap["Port"] = parsedInterface[1]
	interfaceAttrMap["Role"] = d.Get("role").(string)
	interfaceAttrMap["Shutdown"] = getAdminState(d.Get("admin_state").(string))
	interfaceAttrMap["SubPort"] = parsedInterface[2]

	var err error
	var interfaceDistinguishedName string

	if portType == "access" {

		breakout := d.Get("breakout").(string)
		policyGroup := d.Get("policy_group").(string)

		if breakout != "" && breakout != "none" && policyGroup != "" {
			return diag.FromErr(fmt.Errorf("Policy Group: %s and Breakout: %s cannot be configured togater.", policyGroup, breakout))
		} else if policyGroup != "" {
			interfaceAttrMap["AssocGrp"] = policyGroup
			interfaceAttrMap["BrkoutMap"] = "none"
		} else {
			interfaceAttrMap["BrkoutMap"] = breakout
		}

		interfaceAttrMap["PcMember"] = d.Get("pc_member").(string)
		accessInterfaceAttr := models.InfraPortConfigurationAttributes{}
		setModelAttributes(&accessInterfaceAttr, interfaceAttrMap)
		rn := fmt.Sprintf(models.RnInfraPortConfig, interfaceAttrMap["Node"], interfaceAttrMap["Card"], interfaceAttrMap["Port"], interfaceAttrMap["SubPort"])
		accessInterfaceConfig := models.NewInfraPortConfiguration(rn, models.ParentDnInfraPortConfig, interfaceAttrMap["Descr"], accessInterfaceAttr)
		err = aciClient.Save(accessInterfaceConfig)
		interfaceDistinguishedName = accessInterfaceConfig.DistinguishedName
	} else if portType == "fabric" {
		interfaceAttrMap["AssocGrp"] = d.Get("policy_group").(string)
		fabricInterfaceAttr := models.FabricPortConfigurationAttributes{}
		setModelAttributes(&fabricInterfaceAttr, interfaceAttrMap)
		rn := fmt.Sprintf(models.RnFabricPortConfig, interfaceAttrMap["Node"], interfaceAttrMap["Card"], interfaceAttrMap["Port"], interfaceAttrMap["SubPort"])
		fabricInterfaceConfig := models.NewFabricPortConfiguration(rn, models.ParentDnFabricPortConfig, interfaceAttrMap["Descr"], fabricInterfaceAttr)
		err = aciClient.Save(fabricInterfaceConfig)
		interfaceDistinguishedName = fabricInterfaceConfig.DistinguishedName
	}

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(interfaceDistinguishedName)

	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciInterfaceConfigurationRead(ctx, d, m)
}

func resourceAciInterfaceConfigurationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Interface Configuration: Beginning Update")
	aciClient := m.(*client.Client)

	node := strconv.Itoa(d.Get("node").(int))
	interfaceAttrMap := make(map[string]string)
	portType := d.Get("port_type").(string)
	parsedInterface := strings.Split(getInterfaceVal(d.Get("interface")), "/")

	if Annotation, ok := d.GetOk("annotation"); ok {
		interfaceAttrMap["Annotation"] = Annotation.(string)
	} else {
		interfaceAttrMap["Annotation"] = "{}"
	}

	interfaceAttrMap["Card"] = parsedInterface[0]
	interfaceAttrMap["Descr"] = d.Get("description").(string)
	interfaceAttrMap["Node"] = node
	interfaceAttrMap["Port"] = parsedInterface[1]
	interfaceAttrMap["Role"] = d.Get("role").(string)
	interfaceAttrMap["Shutdown"] = getAdminState(d.Get("admin_state").(string))
	interfaceAttrMap["SubPort"] = parsedInterface[2]

	var err error
	var interfaceDistinguishedName string

	if portType == "access" {

		breakout := d.Get("breakout").(string)
		policyGroup := d.Get("policy_group").(string)

		if breakout != "" && breakout != "none" && policyGroup != "" {
			return diag.FromErr(fmt.Errorf("Policy Group: %s and Breakout: %s cannot be configured togater.", policyGroup, breakout))
		} else if (breakout == "" || breakout == "none") && policyGroup == "" {
			interfaceAttrMap["AssocGrp"] = ""
			interfaceAttrMap["BrkoutMap"] = "none"
		} else if policyGroup != "" {
			interfaceAttrMap["AssocGrp"] = policyGroup
			interfaceAttrMap["BrkoutMap"] = "none"
		} else {
			interfaceAttrMap["BrkoutMap"] = breakout
		}

		interfaceAttrMap["PcMember"] = d.Get("pc_member").(string)
		accessInterfaceAttr := models.InfraPortConfigurationAttributes{}
		setModelAttributes(&accessInterfaceAttr, interfaceAttrMap)
		rn := fmt.Sprintf(models.RnInfraPortConfig, interfaceAttrMap["Node"], interfaceAttrMap["Card"], interfaceAttrMap["Port"], interfaceAttrMap["SubPort"])
		accessInterfaceConfig := models.NewInfraPortConfiguration(rn, models.ParentDnInfraPortConfig, interfaceAttrMap["Descr"], accessInterfaceAttr)
		accessInterfaceConfig.Status = "modified"
		err = aciClient.Save(accessInterfaceConfig)
		interfaceDistinguishedName = accessInterfaceConfig.DistinguishedName
	} else if portType == "fabric" {
		interfaceAttrMap["AssocGrp"] = d.Get("policy_group").(string)
		fabricInterfaceAttr := models.FabricPortConfigurationAttributes{}
		setModelAttributes(&fabricInterfaceAttr, interfaceAttrMap)
		rn := fmt.Sprintf(models.RnFabricPortConfig, interfaceAttrMap["Node"], interfaceAttrMap["Card"], interfaceAttrMap["Port"], interfaceAttrMap["SubPort"])
		fabricInterfaceConfig := models.NewFabricPortConfiguration(rn, models.ParentDnFabricPortConfig, interfaceAttrMap["Descr"], fabricInterfaceAttr)
		fabricInterfaceConfig.Status = "modified"
		err = aciClient.Save(fabricInterfaceConfig)
		interfaceDistinguishedName = fabricInterfaceConfig.DistinguishedName
	}

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(interfaceDistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciInterfaceConfigurationRead(ctx, d, m)
}

func resourceAciInterfaceConfigurationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := getAndSetRemoteInterfaceConfiguration(dn, aciClient, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciInterfaceConfigurationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	portType := d.Get("port_type").(string)
	var moName string
	if portType == "access" {
		moName = models.InfraPortConfigClassName
	} else if portType == "fabric" {
		moName = models.FabricPortConfigClassName
	}

	err := aciClient.DeleteByDn(dn, moName)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
