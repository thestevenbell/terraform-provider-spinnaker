package spinnaker

import (
	"strings"

	"github.com/armory-io/terraform-provider-spinnaker/spinnaker/api"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceApplication() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"application": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateApplicationName,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"permissions_read": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"permissions_write": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"permissions_execute": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Create: resourceApplicationCreate,
		Read:   resourceApplicationRead,
		Update: resourceApplicationUpdate,
		Delete: resourceApplicationDelete,
		Exists: resourceApplicationExists,
	}
}

// {
// 	"application": "applicationA",
// 	"cloudProviders": "kubernetes",
// 	"description": "Demo sponnet application",
// 	"email": "youremail@example.com",
// 	"spec": {
// 	   "email": "youremail@example.com",
// 	   "permissions": {
// 		  "EXECUTE": [
// 			 "com_sre_dev"
// 		  ],
// 		  "READ": [
// 			 "com_sre_dev"
// 		  ],
// 		  "WRITE": [
// 			 "com_sre_dev"
// 		  ]
// 	   }
// 	},
// 	"user": "youremail@example.com"
//  }

type applicationRead struct {
	Name       string `json:"name"`
	Attributes struct {
		Email string `json:"email"`
	} `json:"attributes"`
}

func resourceApplicationCreate(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	application := data.Get("application").(string)
	email := data.Get("email").(string)

	permissionsRead := data.Get("permissions_read").([]string)
	permissionsWrite := data.Get("permissions_write").([]string)
	permissionsExecute := data.Get("permissions_execute").([]string)

	if err := api.CreateApplication(client, application, email, permissionsRead, permissionsWrite, permissionsExecute); err != nil {
		return err
	}

	return resourceApplicationRead(data, meta)
}

func resourceApplicationRead(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	applicationName := data.Get("application").(string)
	var app applicationRead
	if err := api.GetApplication(client, applicationName, &app); err != nil {
		return err
	}

	return readApplication(data, app)
}

func resourceApplicationUpdate(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	applicationName := data.Get("application").(string)
	var app applicationRead
	if err := api.GetApplication(client, applicationName, &app); err != nil {
		return err
	}

	return readApplication(data, app)
}

func resourceApplicationDelete(data *schema.ResourceData, meta interface{}) error {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	applicationName := data.Get("application").(string)

	return api.DeleteAppliation(client, applicationName)
}

func resourceApplicationExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	clientConfig := meta.(gateConfig)
	client := clientConfig.client
	applicationName := data.Get("application").(string)

	var app applicationRead
	if err := api.GetApplication(client, applicationName, &app); err != nil {
		errmsg := err.Error()
		if strings.Contains(errmsg, "not found") {
			return false, nil
		}
		return false, err
	}

	if app.Name == "" {
		return false, nil
	}

	return true, nil
}

func readApplication(data *schema.ResourceData, application applicationRead) error {
	data.SetId(application.Name)
	return nil
}
