package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciLeakInternalSubnet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciLeakInternalSubnetCreate,
		UpdateContext: resourceAciLeakInternalSubnetUpdate,
		ReadContext:   resourceAciLeakInternalSubnetRead,
		DeleteContext: resourceAciLeakInternalSubnetDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciLeakInternalSubnetImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"vrf_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vrf_scope": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"private",
					"public",
				}, false),
				Default: "private",
			},
			// Tenant and VRF Destinations
			"leak_to": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination_vrf_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"destination_vrf_scope": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"inherit",
								"private",
								"public",
							}, false),
							Default: "inherit",
						},
						"destination_tenant_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		})),
	}
}

func getRemoteLeakInternalSubnet(client *client.Client, dn string) (*models.LeakInternalSubnet, error) {
	leakInternalSubnetCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	leakInternalSubnet := models.LeakInternalSubnetFromContainer(leakInternalSubnetCont)
	if leakInternalSubnet.DistinguishedName == "" {
		return nil, fmt.Errorf("LeakInternalSubnet %s not found", dn)
	}
	return leakInternalSubnet, nil
}

func setLeakInternalSubnetAttributes(leakInternalSubnet *models.LeakInternalSubnet, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(leakInternalSubnet.DistinguishedName)
	d.Set("description", leakInternalSubnet.Description)
	if dn != leakInternalSubnet.DistinguishedName {
		d.Set("vrf_dn", "")
	}
	leakInternalSubnetMap, err := leakInternalSubnet.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("vrf_dn", GetParentDn(dn, fmt.Sprintf("/leakroutes/leakintsubnet-[%s]", leakInternalSubnetMap["ip"])))
	d.Set("annotation", leakInternalSubnetMap["annotation"])
	d.Set("ip", leakInternalSubnetMap["ip"])
	d.Set("name", leakInternalSubnetMap["name"])
	d.Set("vrf_scope", leakInternalSubnetMap["scope"])
	d.Set("name_alias", leakInternalSubnetMap["nameAlias"])
	return d, nil
}

func getAndSetLeakToObjects(client *client.Client, dn string, d *schema.ResourceData) (*schema.ResourceData, error) {
	leakToObjects, err := client.ListTenantandVRFdestinationforInterVRFLeakedRoutes(dn)
	if len(leakToObjects) >= 1 {
		newContent := make([]map[string]string, 0, 1)
		for _, leakToObject := range leakToObjects {
			newContent = append(newContent, map[string]string{
				"destination_tenant_name": leakToObject.DestinationTenantName,
				"destination_vrf_name":    leakToObject.DestinationCtxName,
				"destination_vrf_scope":   leakToObject.Scope,
			})
		}
		d.Set("leak_to", newContent)
		return d, nil
	} else if err != nil {
		d.Set("leak_to", nil)
		log.Printf("[DEBUG]: Could not find existing leakTo objects under the parent DN: %s", dn)
	}

	return d, nil
}

func resourceAciLeakInternalSubnetImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	leakInternalSubnet, err := getRemoteLeakInternalSubnet(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setLeakInternalSubnetAttributes(leakInternalSubnet, d)
	if err != nil {
		return nil, err
	}

	// leakTo - Beginning Import
	log.Printf("[DEBUG] %s: leakTo - Beginning Import with parent DN", dn)
	_, err = getAndSetLeakToObjects(aciClient, dn, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: leakTo - Import finished successfully with parent DN", dn)
	// leakTo - Import finished successfully

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciLeakInternalSubnetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LeakInternalSubnet: Beginning Creation")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	ip := d.Get("ip").(string)
	VRFDn := d.Get("vrf_dn").(string)

	leakInternalSubnetAttr := models.LeakInternalSubnetAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		leakInternalSubnetAttr.Annotation = Annotation.(string)
	} else {
		leakInternalSubnetAttr.Annotation = "{}"
	}

	if Ip, ok := d.GetOk("ip"); ok {
		leakInternalSubnetAttr.Ip = Ip.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		leakInternalSubnetAttr.Name = Name.(string)
	}

	if Scope, ok := d.GetOk("vrf_scope"); ok {
		leakInternalSubnetAttr.Scope = Scope.(string)
	}

	// Creation parent(leakRoutes) object
	leakRoutesAttr := models.InterVRFLeakedRoutesContainerAttributes{}
	leakRoutes := models.NewInterVRFLeakedRoutesContainer(fmt.Sprintf(models.RnleakRoutes), VRFDn, desc, nameAlias, leakRoutesAttr)

	err := aciClient.Save(leakRoutes)
	if err != nil {
		return diag.FromErr(err)
	}

	// Creation leakInternalSubnet object
	leakInternalSubnet := models.NewLeakInternalSubnet(fmt.Sprintf(models.RnleakInternalSubnet, ip), leakRoutes.DistinguishedName, desc, nameAlias, leakInternalSubnetAttr)

	err = aciClient.Save(leakInternalSubnet)
	if err != nil {
		return diag.FromErr(err)
	}

	// Creation leakTo objects under leakInternalSubnet
	leakInternalSubnet.Status = "modified"

	jsonPayload, _, err := aciClient.PrepareModel(leakInternalSubnet)
	if err != nil {
		return diag.FromErr(err)
	}
	jsonPayload.Array(leakInternalSubnet.ClassName, "children")

	// leakTo - Create Operations
	if leakToSet, ok := d.GetOk("leak_to"); ok {
		log.Printf("[DEBUG] leakTo: Beginning Creation")

		leakToList := leakToSet.(*schema.Set).List()

		for _, leakToObject := range leakToList {
			leakToObjectMap := leakToObject.(map[string]interface{})

			cidrJSON := []byte(fmt.Sprintf(`
			{
                "leakTo": {
                    "attributes": {
                        "tenantName": "%s",
                        "ctxName": "%s",
                        "status": "created",
                        "scope": "%s"
                    }
                }
            }`, leakToObjectMap["destination_tenant_name"].(string), leakToObjectMap["destination_vrf_name"].(string), leakToObjectMap["destination_vrf_scope"].(string)))

			cidrCon, err := container.ParseJSON(cidrJSON)
			if err != nil {
				return diag.FromErr(err)
			}

			jsonPayload.ArrayAppend(cidrCon.Data(), leakInternalSubnet.ClassName, "children")
		}

	}
	// leakTo - Create Operations

	log.Printf("[DEBUG] leakInternalSubnet object request data: %v", jsonPayload)

	leakToRequestData, err := aciClient.MakeRestRequest("POST", fmt.Sprintf("/api/node/mo/%s/%s.json", leakRoutes.DistinguishedName, fmt.Sprintf(models.RnleakInternalSubnet, ip)), jsonPayload, true)
	if err != nil {
		return diag.FromErr(err)
	}

	_, _, err = aciClient.Do(leakToRequestData)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(leakInternalSubnet.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciLeakInternalSubnetRead(ctx, d, m)
}

func resourceAciLeakInternalSubnetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LeakInternalSubnet: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	ip := d.Get("ip").(string)
	VRFDn := d.Get("vrf_dn").(string)

	leakInternalSubnetAttr := models.LeakInternalSubnetAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		leakInternalSubnetAttr.Annotation = Annotation.(string)
	} else {
		leakInternalSubnetAttr.Annotation = "{}"
	}

	if Ip, ok := d.GetOk("ip"); ok {
		leakInternalSubnetAttr.Ip = Ip.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		leakInternalSubnetAttr.Name = Name.(string)
	}

	if Scope, ok := d.GetOk("vrf_scope"); ok {
		leakInternalSubnetAttr.Scope = Scope.(string)
	}
	leakInternalSubnet := models.NewLeakInternalSubnet(fmt.Sprintf("leakintsubnet-[%s]", ip), VRFDn+"/leakroutes", desc, nameAlias, leakInternalSubnetAttr)

	leakInternalSubnet.Status = "modified"

	jsonPayload, _, err := aciClient.PrepareModel(leakInternalSubnet)
	if err != nil {
		return diag.FromErr(err)
	}
	jsonPayload.Array(leakInternalSubnet.ClassName, "children")

	// leakTo - Update Operations
	if d.HasChange("leak_to") {
		if leakToSet, ok := d.GetOk("leak_to"); ok {
			log.Printf("[DEBUG] leakTo: Beginning Creation")

			// Remove existing leakTo objects
			leakToObjects, err := aciClient.ListTenantandVRFdestinationforInterVRFLeakedRoutes(leakInternalSubnet.DistinguishedName)
			if len(leakToObjects) >= 1 {
				log.Printf("[DEBUG]: leakTo Existing Objects List: %v", leakToObjects)

				existingObjectsPayload, _, err := aciClient.PrepareModel(leakInternalSubnet)
				if err != nil {
					return diag.FromErr(err)
				}
				existingObjectsPayload.Array(leakInternalSubnet.ClassName, "children")

				for _, leakToObject := range leakToObjects {
					cidrJSON := []byte(fmt.Sprintf(`
					{
						"leakTo": {
							"attributes": {
								"dn": "%s",
								"status": "deleted"
							}
						}
					}`, leakToObject.DistinguishedName))

					cidrCon, err := container.ParseJSON(cidrJSON)
					if err != nil {
						return diag.FromErr(err)
					}
					existingObjectsPayload.ArrayAppend(cidrCon.Data(), leakInternalSubnet.ClassName, "children")
				}

				leakToRequestData, err := aciClient.MakeRestRequest("POST", fmt.Sprintf("/api/node/mo/%s/%s.json", VRFDn+"/leakroutes", fmt.Sprintf(models.RnleakInternalSubnet, ip)), existingObjectsPayload, true)
				if err != nil {
					return diag.FromErr(err)
				}

				_, _, err = aciClient.Do(leakToRequestData)
				if err != nil {
					return diag.FromErr(err)
				} else {
					log.Printf("[DEBUG]: leakTo existing objects destroy finished successfully")
				}
			} else if err != nil {
				log.Printf("[DEBUG]: Could not find existing leakTo objects under the parent DN: %s", leakInternalSubnet.DistinguishedName)
			}

			// Create new leakTo objects
			leakToList := leakToSet.(*schema.Set).List()

			for _, leakToObject := range leakToList {
				leakToObjectMap := leakToObject.(map[string]interface{})

				cidrJSON := []byte(fmt.Sprintf(`
				{
					"leakTo": {
						"attributes": {
							"tenantName": "%s",
							"ctxName": "%s",
							"status": "created",
							"scope": "%s"
						}
					}
				}`, leakToObjectMap["destination_tenant_name"].(string), leakToObjectMap["destination_vrf_name"].(string), leakToObjectMap["destination_vrf_scope"].(string)))

				cidrCon, err := container.ParseJSON(cidrJSON)
				if err != nil {
					return diag.FromErr(err)
				}

				jsonPayload.ArrayAppend(cidrCon.Data(), leakInternalSubnet.ClassName, "children")
			}

		} else {
			// Remove existing leakTo objects
			leakToObjects, err := aciClient.ListTenantandVRFdestinationforInterVRFLeakedRoutes(leakInternalSubnet.DistinguishedName)
			if len(leakToObjects) >= 1 {
				log.Printf("[DEBUG]: leakTo Existing Objects List: %v", leakToObjects)
				for _, leakToObject := range leakToObjects {
					cidrJSON := []byte(fmt.Sprintf(`
					{
						"leakTo": {
							"attributes": {
								"dn": "%s",
								"status": "deleted"
							}
						}
					}`, leakToObject.DistinguishedName))

					cidrCon, err := container.ParseJSON(cidrJSON)
					if err != nil {
						return diag.FromErr(err)
					}
					jsonPayload.ArrayAppend(cidrCon.Data(), leakInternalSubnet.ClassName, "children")
				}
			} else if err != nil {
				log.Printf("[DEBUG]: Could not find existing leakTo objects under the parent DN: %s", leakInternalSubnet.DistinguishedName)
			}
		}
	}
	// leakTo - Update Operations

	log.Printf("[DEBUG] leakInternalSubnet object request data: %v", jsonPayload)

	leakToRequestData, err := aciClient.MakeRestRequest("POST", fmt.Sprintf("/api/node/mo/%s/%s.json", VRFDn+"/leakroutes", fmt.Sprintf(models.RnleakInternalSubnet, ip)), jsonPayload, true)
	if err != nil {
		return diag.FromErr(err)
	}

	_, _, err = aciClient.Do(leakToRequestData)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(leakInternalSubnet.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciLeakInternalSubnetRead(ctx, d, m)
}

func resourceAciLeakInternalSubnetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	leakInternalSubnet, err := getRemoteLeakInternalSubnet(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	_, err = setLeakInternalSubnetAttributes(leakInternalSubnet, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	// leakTo - Beginning Read
	log.Printf("[DEBUG] %s: leakTo - Beginning Read with parent DN", dn)
	_, err = getAndSetLeakToObjects(aciClient, dn, d)
	if err != nil {
		return nil
	}
	log.Printf("[DEBUG] %s: leakTo - Read finished successfully with parent DN", dn)
	// leakTo - Read finished successfully

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciLeakInternalSubnetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "leakInternalSubnet")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
