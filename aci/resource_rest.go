package aci

import (
	"errors"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/container"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/ghodss/yaml"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const Created = "created"
const Deleted = "deleted"

const ErrDistinguishedNameNotFound = "The Dn is not present in the content"

func resourceAciRest() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciRestCreate,
		Update: resourceAciRestUpdate,
		Read:   resourceAciRestRead,
		Delete: resourceAciRestDelete,

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"path": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			// we set it automatically if file config is provided
			"class_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"content": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			"dn": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"payload": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			// some dn lookups can take a long time, make this optional only checked if dn not supplied
			"lookup_dn": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceAciRestCreate(d *schema.ResourceData, m interface{}) error {
	cont, err := PostAndSetStatus(d, m, "created, modified")
	if err != nil {
		return err
	}
	className := GetClass(cont)

	dn := models.StripQuotes(models.StripSquareBrackets(cont.Search(className, "attributes", "dn").String()))
	lookupDn := d.Get("lookup_dn").(bool)

	if dn == "{}" && lookupDn {
		d.SetId(GetDN(d, cont, m))
	} else {
		d.SetId(dn)
	}

	return resourceAciRestRead(d, m)
}

func resourceAciRestUpdate(d *schema.ResourceData, m interface{}) error {
	cont, err := PostAndSetStatus(d, m, "modified")
	if err != nil {
		return err
	}
	className := GetClass(cont)

	dn := models.StripQuotes(models.StripSquareBrackets(cont.Search(className, "attributes", "dn").String()))
	lookupDn := d.Get("lookup_dn").(bool)

	if dn == "{}" && lookupDn {
		d.SetId(GetDN(d, cont, m))
	} else {
		d.SetId(dn)
	}

	return resourceAciRestRead(d, m)
}

func resourceAciRestRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceAciRestDelete(d *schema.ResourceData, m interface{}) error {
	_, err := PostAndSetStatus(d, m, Deleted)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

// GetDN performs a lookup on the APIC to find the real DN of the object created
func GetDN(d *schema.ResourceData, c *container.Container, m interface{}) string {
	aciClient := m.(*client.Client)
	path := d.Get("path").(string)
	class := GetClass(c)
	name := GetName(c)

	// Assume an object was defined directly to its path
	queryPath := fmt.Sprintf(`%v?query-target-filter=eq(%v.name,%v)`, path, class, name)

	cont, _ := aciClient.GetViaURL(queryPath)

	dn := models.StripQuotes(models.StripSquareBrackets(cont.Search("imdata", class, "attributes", "dn").String()))

	// Child object was created with parent as path
	if dn == "{}" {
		queryPath = fmt.Sprintf(`%v?query-target=children&query-target-filter=eq(%v.name,%v)`, path, class, name)

		childCont, _ := aciClient.GetViaURL(queryPath)

		dn = models.StripQuotes(models.StripSquareBrackets(childCont.Search("imdata", class, "attributes", "dn").String()))
	}

	return fmt.Sprintf("%s", dn)
}

// GetName returns the name of the object being created
func GetName(c *container.Container) string {
	var name string

	contMap, err := c.ChildrenMap()
	if err != nil {
		return name
	}

	for _, v := range contMap {
		name = v.Path("attributes.name").String()
		break
	}

	return name
}

// GetClass returns the normalized class (as first key of the container)
func GetClass(c *container.Container) string {
	var class string

	contMap, err := c.ChildrenMap()
	if err != nil {
		return class
	}

	for k, _ := range contMap {
		class = k
		break
	}

	return class
}

// PostAndSetStatus is used to post schema and set the status
func PostAndSetStatus(d *schema.ResourceData, m interface{}, status string) (*container.Container, error) {
	aciClient := m.(*client.Client)
	path := d.Get("path").(string)
	var cont *container.Container
	var err error
	method := "POST"

	cont, err = GetCont(d, status)
	if err != nil {
		return nil, err
	}

	req, err := aciClient.MakeRestRequest(method, path, cont, true)
	if err != nil {
		return nil, err
	}

	respCont, _, err := aciClient.Do(req)
	if err != nil {
		return nil, err
	}
	err = client.CheckForErrors(respCont, method, false)
	if err != nil {
		return nil, err
	}

	return cont, nil
}

// GetCont attempts to unify the different input methods for the provider by creating a single reusable parser
func GetCont(d *schema.ResourceData, status string) (*container.Container, error) {
	var cont *container.Container
	var err error

	if content, ok := d.GetOk("content"); ok {
		contentStrMap := toStrMap(content.(map[string]interface{}))

		if classNameIntf, ok := d.GetOk("class_name"); ok {
			className := classNameIntf.(string)
			cont, err = preparePayload(className, contentStrMap)
			if status == Deleted {
				cont.Set(status, className, "attributes", "status")
			}
			if err != nil {
				return nil, err
			}

		} else {
			return nil, errors.New("The className is required when content is provided explicitly")
		}

	} else if payload, ok := d.GetOk("payload"); ok {
		payloadStr := payload.(string)
		if len(payloadStr) == 0 {
			return nil, fmt.Errorf("Payload cannot be empty string")
		}

		yamlJSONPayload, err := yaml.YAMLToJSON([]byte(payloadStr))

		if err != nil {
			// It may be possible that the payload is in JSON
			jsonPayload, err := container.ParseJSON([]byte(payloadStr))
			if err != nil {
				return nil, fmt.Errorf("Invalid format for yaml/JSON payload")
			}
			cont = jsonPayload
		} else {
			// we have valid yaml payload and we were able to convert it to json
			cont, err = container.ParseJSON(yamlJSONPayload)
			if err != nil {
				return nil, fmt.Errorf("Failed to convert YAML to JSON")
			}
		}

		if err != nil {

			return nil, fmt.Errorf("Unable to parse the payload to JSON. Please check your payload")
		}

	} else {
		return nil, fmt.Errorf("Either of payload or content is required")
	}

	return cont, err
}
