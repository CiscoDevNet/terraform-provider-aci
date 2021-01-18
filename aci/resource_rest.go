package aci

import (
	"errors"
	"fmt"
	"log"

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
		Create:        resourceAciRestCreate,
		Update:        resourceAciRestUpdate,
		Read:          resourceAciRestRead,
		Delete:        resourceAciRestDelete,
		CustomizeDiff: resourceAciRestCustomizeDiff,

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"path": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			// we set it automatically if file config is provided
			"class_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"content": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			"dn": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"payload": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}
func resourceAciRestCustomizeDiff(d *schema.ResourceDiff, m interface{}) error {

	if d.HasChange("path") || d.HasChange("class_name") {
		log.Printf("[DEBUG] path or class_name changed, resource %s will be recreated, no further check required",
			d.Id())
		return nil
	}

	if d.HasChange("content") {
		log.Printf("[DEBUG] content of resource %s has changed, checking if the dn is still valid \n", d.Id())
		var oldContent, newContent map[string]interface{}

		oldUntyped, newUntyped := d.GetChange("content")

		switch v := oldUntyped.(type) {
		case map[string]interface{}:
			oldContent = v
		default:
			oldContent = nil
		}

		switch v := newUntyped.(type) {
		case map[string]interface{}:
			newContent = v
		default:
			oldContent = nil
		}
		//New content or old content is empty, skipping diff compute
		if newContent == nil || oldContent == nil {
			return nil
		}

		if oldContent["dn"] != newContent["dn"] {
			d.SetNewComputed("dn")
		}

	} else if d.HasChange("payload") {
		//TODO : handle payload change
	}

	return nil
}

func resourceAciRestCreate(d *schema.ResourceData, m interface{}) error {
	cont, err := PostAndSetStatus(d, m, "created, modified")
	if err != nil {
		return err
	}
	classNameIntf := d.Get("class_name")
	className := classNameIntf.(string)
	dn := models.StripQuotes(models.StripSquareBrackets(cont.Search(className, "attributes", "dn").String()))

	if dn == "{}" {
		d.SetId(GetDN(d, m))

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
	classNameIntf := d.Get("class_name")
	className := classNameIntf.(string)
	dn := models.StripQuotes(models.StripSquareBrackets(cont.Search(className, "attributes", "dn").String()))
	if dn == "{}" {
		d.SetId(GetDN(d, m))

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

func GetDN(d *schema.ResourceData, m interface{}) string {
	aciClient := m.(*client.Client)
	path := d.Get("path").(string)
	className := d.Get("class_name").(string)
	cont, _ := aciClient.GetViaURL(path)
	dn := models.StripQuotes(models.StripSquareBrackets(cont.Search("imdata", className, "attributes", "dn").String()))
	return fmt.Sprintf("%s", dn)
}

// PostAndSetStatus is used to post schema and set the status
func PostAndSetStatus(d *schema.ResourceData, m interface{}, status string) (*container.Container, error) {
	aciClient := m.(*client.Client)
	path := d.Get("path").(string)
	var cont *container.Container
	var err error
	method := "POST"

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
		cont, err = payloadToContainer(payload)
		if err != nil {
			return nil, err
		}

		if status == "deleted" {
			method = "DELETE"
		}

	} else {
		return nil, fmt.Errorf("Either of payload or content is required")
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

func payloadToContainer(payload interface{}) (*container.Container, error) {
	var cont *container.Container
	payloadStr := payload.(string)
	if len(payloadStr) == 0 {
		return nil, fmt.Errorf("Payload cannot be empty string")
	}

	yamlJsonPayload, err := yaml.YAMLToJSON([]byte(payloadStr))

	if err != nil {
		// It may be possible that the payload is in JSON
		jsonPayload, err := container.ParseJSON([]byte(payloadStr))
		if err != nil {
			return nil, fmt.Errorf("Invalid format for yaml/JSON payload")
		}
		cont = jsonPayload
	} else {
		// we have valid yaml payload and we were able to convert it to json
		cont, err = container.ParseJSON(yamlJsonPayload)
		if err != nil {
			return nil, fmt.Errorf("Failed to convert YAML to JSON.")
		}
	}

	if err != nil {

		return nil, fmt.Errorf("Unable to parse the payload to JSON. Please check your payload")
	}
	return cont, nil
}
