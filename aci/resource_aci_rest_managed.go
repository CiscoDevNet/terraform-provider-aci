package provider

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAciRest() *schema.Resource {
	return &schema.Resource{
		Description: "Manages ACI Model Objects via REST API calls. This resource can only manage a single API object and its direct children. It is able to read the state and therefore reconcile configuration drift.",

		CreateContext: resourceAciRestCreate,
		UpdateContext: resourceAciRestUpdate,
		ReadContext:   resourceAciRestRead,
		DeleteContext: resourceAciRestDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAciRestImport,
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

func getAciRest(d *schema.ResourceData, c *container.Container) diag.Diagnostics {
	className := d.Get("class_name").(string)
	dn := d.Get("dn").(string)
	d.SetId(dn)

	content := d.Get("content")
	contentStrMap := toStrMap(content.(map[string]interface{}))
	newContent := make(map[string]interface{})
	restContent, ok := c.Search("imdata", className, "attributes").Index(0).Data().(map[string]interface{})

	if !ok {
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

func resourceAciRestReadHelper(ctx context.Context, d *schema.ResourceData, meta interface{}, expectObject bool) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	for attempts := 0; ; attempts++ {
		getChildren := false
		if len(d.Get("child").(*schema.Set).List()) > 0 {
			getChildren = true
		}
		cont, diags := ApicRest(d, meta, "GET", getChildren)
		if diags.HasError() {
			if ok := backoff(attempts, meta.(apiClient).Retries); !ok {
				return diags
			}
			log.Printf("[ERROR] Failed to read object: %s, retries: %v", diags[0].Summary, attempts)
			continue
		}

		// Check if we received an empty response without errors -> object has been deleted
		if cont == nil && diags == nil && !expectObject {
			d.SetId("")
			return nil
		}

		diags = getAciRest(d, cont)
		if !diags.HasError() {
			break
		}
		if ok := backoff(attempts, meta.(apiClient).Retries); !ok {
			return diags
		}
		log.Printf("[ERROR] Failed to decode response after reading object: %s, retries: %v", diags[0].Summary, attempts)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciRestCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if meta.(apiClient).IsMock {
		d.SetId(d.Get("dn").(string))
		return nil
	}

	log.Printf("[DEBUG] %s: Beginning Create", d.Id())

	for attempts := 0; ; attempts++ {
		_, diags := ApicRest(d, meta, "POST", false)
		if !diags.HasError() {
			break
		}
		if ok := backoff(attempts, meta.(apiClient).Retries); !ok {
			return diags
		}
		log.Printf("[ERROR] Failed to create object: %s, retries: %v", diags[0].Summary, attempts)
	}

	log.Printf("[DEBUG] %s: Create finished successfully", d.Id())
	return resourceAciRestReadHelper(ctx, d, meta, true)
}

func resourceAciRestUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if meta.(apiClient).IsMock {
		return nil
	}

	log.Printf("[DEBUG] %s: Beginning Update", d.Id())

	for attempts := 0; ; attempts++ {
		_, diags := ApicRest(d, meta, "POST", false)
		if !diags.HasError() {
			break
		}
		if ok := backoff(attempts, meta.(apiClient).Retries); !ok {
			return diags
		}
		log.Printf("[ERROR] Failed to update object: %s, retries: %v", diags[0].Summary, attempts)
	}

	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciRestReadHelper(ctx, d, meta, true)
}

func resourceAciRestRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if meta.(apiClient).IsMock {
		return nil
	}

	return resourceAciRestReadHelper(ctx, d, meta, false)
}

func resourceAciRestDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if meta.(apiClient).IsMock {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	for attempts := 0; ; attempts++ {
		_, diags := ApicRest(d, meta, "DELETE", false)
		if !diags.HasError() {
			break
		}
		if ok := backoff(attempts, meta.(apiClient).Retries); !ok {
			return diags
		}
		log.Printf("[ERROR] Failed to delete object: %s, retries: %v", diags[0].Summary, attempts)
	}

	d.SetId("")
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	return nil
}

func resourceAciRestImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())

	parts := strings.SplitN(d.Id(), ":", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("Unexpected format of ID (%s), expected class_name:dn", d.Id())
	}

	d.Set("dn", parts[1])
	d.Set("class_name", parts[0])
	d.SetId(parts[1])

	if diags := resourceAciRestReadHelper(ctx, d, meta, true); diags.HasError() {
		return nil, fmt.Errorf("Could not read object when importing: %s", diags[0].Summary)
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{d}, nil
}
