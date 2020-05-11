package aci

import (
	"errors"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/container"
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
				Required: true,
			},
			"dn": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAciRestCreate(d *schema.ResourceData, m interface{}) error {
	cont, err := PostAndSetStatus(d, m, "created, modified")
	if err != nil {
		return err
	}
	classNameIntf := d.Get("class_name")
	className := classNameIntf.(string)
	dn := cont.Search(className, "attributes", "dn").String()

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
	dn := cont.Search(className, "attributes", "dn").String()
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

	if content, ok := d.GetOk("content"); ok {
		contentStrMap := toStrMap(content.(map[string]interface{}))

		if classNameIntf, ok := d.GetOk("class_name"); ok {
			className := classNameIntf.(string)
			cont, err = preparePayload(className, contentStrMap)
			if err != nil {
				return nil, err
			}

		} else {
			return nil, errors.New("The className is required when content is provided explicitly")
		}

	}
	req, err := aciClient.MakeRestRequest("POST", path, cont, true)
	if err != nil {
		return nil, err
	}

	respCont, _, err := aciClient.Do(req)
	if err != nil {
		return nil, err
	}
	err = client.CheckForErrors(respCont, "POST")
	if err != nil {
		return nil, err
	}
	return cont, nil
}
