package aci

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// List of attributes to be not stored in state
var IgnoreAttr = []string{"extMngdBy", "lcOwn", "modTs", "monPolDn", "uid", "dn", "rn", "configQual", "configSt", "virtualIp", "annotation"}

// List of attributes to be only written to state from config
var WriteOnlyAttr = []string{"childAction"}

// List of classes where 'rsp-prop-include=config-only' does not return the desired objects/properties
var FullClasses = []string{"firmwareFwGrp", "maintMaintGrp", "maintMaintP", "firmwareFwP", "pkiExportEncryptionKey"}

// List of classes where an immediate GET following a POST might not reflect the created/updated object
var AllowEmptyReadClasses = []string{"firmwareFwGrp", "firmwareRsFwgrpp", "firmwareFwP", "fabricNodeBlk"}

// List of classes which do not support annotations
var NoAnnotationClasses = []string{"tagTag"}

func resourceAciRestManaged() *schema.Resource {
	return &schema.Resource{
		Description: "Manages ACI Model Objects via REST API calls. This resource can only manage a single API object and its direct children. It is able to read the state and therefore reconcile configuration drift.",

		CreateContext: resourceAciRestManagedCreate,
		UpdateContext: resourceAciRestManagedUpdate,
		ReadContext:   resourceAciRestManagedRead,
		DeleteContext: resourceAciRestManagedDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAciRestManagedImport,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The distinguished name of the object.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"dn": {
				Type:        schema.TypeString,
				Description: "Distinguished name of object being managed including its relative name, e.g. uni/tn-EXAMPLE_TENANT.",
				Required:    true,
				ForceNew:    true,
			},
			"class_name": {
				Type:        schema.TypeString,
				Description: "Which class object is being created. (Make sure there is no colon in the classname)",
				Required:    true,
				ForceNew:    true,
			},
			"content": {
				Type:        schema.TypeMap,
				Description: "Map of key-value pairs those needed to be passed to the Model object as parameters. Make sure the key name matches the name with the object parameter in ACI.",
				Optional:    true,
				Computed:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					content := d.Get("content")
					contentStrMap := toStrMap(content.(map[string]interface{}))
					key := strings.Split(k, ".")[1]
					if _, ok := contentStrMap[key]; ok {
						return false
					}
					return true
				},
			},
			"child": {
				Type:        schema.TypeSet,
				Description: "List of children.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rn": {
							Type:        schema.TypeString,
							Description: "The relative name of the child object.",
							Required:    true,
						},
						"class_name": {
							Type:        schema.TypeString,
							Description: "Class name of child object.",
							Optional:    true,
							Computed:    true,
						},
						"content": {
							Type:        schema.TypeMap,
							Description: "Map of key-value pairs which represents the attributes for the child object.",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func getAciRestManaged(d *schema.ResourceData, c *container.Container, expectObject bool) diag.Diagnostics {
	className := d.Get("class_name").(string)
	dn := d.Get("dn").(string)
	d.SetId(dn)

	content := d.Get("content")
	contentStrMap := toStrMap(content.(map[string]interface{}))
	newContent := make(map[string]interface{})
	restContent, ok := c.Search("imdata", className, "attributes").Index(0).Data().(map[string]interface{})

	if !ok {
		if expectObject && containsString(AllowEmptyReadClasses, className) {
			return nil
		}
		return diag.Errorf("Failed to retrieve REST payload for class: %s.", className)
	}

	for attr, value := range restContent {
		// Ignore certain attributes
		if !containsString(IgnoreAttr, attr) {
			newContent[attr] = value.(string)
		}
	}

	for attr, value := range contentStrMap {
		// Do not read/update write-only attributes, eg. 'childAction'
		if containsString(WriteOnlyAttr, attr) {
			newContent[attr] = value
		}
	}
	d.Set("content", newContent)

	newChildrenSet := make([]interface{}, 0, 1)
	for _, child := range d.Get("child").(*schema.Set).List() {
		newChildMap := make(map[string]interface{})
		childRn := child.(map[string]interface{})["rn"].(string)
		childClassName := child.(map[string]interface{})["class_name"].(string)
		childContent := child.(map[string]interface{})["content"]
		newChildMap["rn"] = childRn
		newChildMap["class_name"] = childClassName
		// Loop over retrieved children
		for _, rChild := range c.Search("imdata", className, "children").Index(0).Data().([]interface{}) {
			for rChildClassName, rChildObject := range rChild.(map[string]interface{}) {
				// Look for desired class
				if rChildClassName == childClassName {
					attrMap := rChildObject.(map[string]interface{})["attributes"].(map[string]interface{})
					for attr, value := range attrMap {
						// Find desired object by its rn
						if attr == "rn" && value.(string) == childRn {
							newChildContent := make(map[string]interface{})

							for key := range toStrMap(childContent.(map[string]interface{})) {
								newChildContent[key] = attrMap[key].(string)
							}
							newChildMap["content"] = newChildContent
						}
					}
				}
			}
		}
		newChildrenSet = append(newChildrenSet, newChildMap)
	}
	d.Set("child", newChildrenSet)

	return nil
}

func resourceAciRestManagedReadHelper(ctx context.Context, d *schema.ResourceData, m interface{}, expectObject bool) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	getChildren := false
	if len(d.Get("child").(*schema.Set).List()) > 0 {
		getChildren = true
	}
	cont, diags := MakeAciRestManagedQuery(d, m, "GET", getChildren)
	if diags.HasError() {
		return diags
	}

	// Check if we received an empty response without errors -> object has been deleted
	if cont == nil && diags == nil && !expectObject {
		d.SetId("")
		return nil
	}

	diags = getAciRestManaged(d, cont, expectObject)
	if diags.HasError() {
		return diags
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciRestManagedCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Create", d.Id())

	_, diags := MakeAciRestManagedQuery(d, m, "POST", false)
	if diags.HasError() {
		return diags
	}

	log.Printf("[DEBUG] %s: Create finished successfully", d.Id())
	return resourceAciRestManagedReadHelper(ctx, d, m, true)
}

func resourceAciRestManagedUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Update", d.Id())

	_, diags := MakeAciRestManagedQuery(d, m, "POST", false)
	if diags.HasError() {
		return diags
	}

	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciRestManagedReadHelper(ctx, d, m, true)
}

func resourceAciRestManagedRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceAciRestManagedReadHelper(ctx, d, m, false)
}

func resourceAciRestManagedDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	_, diags := MakeAciRestManagedQuery(d, m, "DELETE", false)
	if diags.HasError() {
		return diags
	}

	d.SetId("")
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	return nil
}

func resourceAciRestManagedImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())

	parts := strings.SplitN(d.Id(), ":", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("Unexpected format of ID (%s), expected class_name:dn", d.Id())
	}

	d.Set("dn", parts[1])
	d.Set("class_name", parts[0])
	d.SetId(parts[1])

	if diags := resourceAciRestManagedReadHelper(ctx, d, m, true); diags.HasError() {
		return nil, fmt.Errorf("Could not read object when importing: %s", diags[0].Summary)
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}

func MakeAciRestManagedQuery(d *schema.ResourceData, m interface{}, method string, children bool) (*container.Container, diag.Diagnostics) {
	aciClient := m.(*client.Client)
	path := "/api/mo/" + d.Get("dn").(string) + ".json"
	className := d.Get("class_name").(string)
	if method == "GET" {
		if children {
			path += "?rsp-subtree=children"
		} else if !containsString(FullClasses, className) {
			path += "?rsp-prop-include=config-only"
		}
	}
	var cont *container.Container = nil
	var err error

	if method == "POST" {
		content := d.Get("content")
		contentStrMap := make(map[string]string)
		if _, ok := contentStrMap["annotation"]; !ok {
			contentStrMap["annotation"] = "orchestrator:terraform"
		}
		for k, v := range toStrMap(content.(map[string]interface{})) {
			contentStrMap[k] = v
		}

		childrenSet := make([]interface{}, 0, 1)

		for _, child := range d.Get("child").(*schema.Set).List() {
			childMap := make(map[string]interface{})
			childClassName := child.(map[string]interface{})["class_name"]
			childContent := make(map[string]string)
			if _, ok := childContent["annotation"]; !ok {
				childContent["annotation"] = "orchestrator:terraform"
			}
			for k, v := range toStrMap(child.(map[string]interface{})["content"].(map[string]interface{})) {
				childContent[k] = v
			}
			childMap["class_name"] = childClassName.(string)
			childMap["content"] = childContent
			childrenSet = append(childrenSet, childMap)
		}

		cont, err = preparePayload(className, contentStrMap, childrenSet)
		if err != nil {
			return nil, diag.FromErr(err)
		}
	}

	req, err := aciClient.MakeRestRequest(method, path, cont, true)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	respCont, _, err := aciClient.Do(req)
	if err != nil {
		return respCont, diag.FromErr(err)
	}
	if respCont.S("imdata").Index(0).String() == "{}" {
		return nil, nil
	}
	err = client.CheckForErrors(respCont, method, false)
	if err != nil {
		return respCont, diag.FromErr(err)
	}
	if method == "POST" {
		return cont, nil
	} else {
		return respCont, nil
	}
}
