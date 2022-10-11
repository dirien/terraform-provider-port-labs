package port

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccPortAction(t *testing.T) {
	identifier := genID()
	actionIdentifier := genID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port-labs_blueprint" "microservice" {
		title = "TF test microservice"
		icon = "Terraform"
		identifier = "%s"
		properties {
			identifier = "text"
			type = "string"
			title = "text"
		}
	}
	resource "port-labs_action" "restart_microservice" {
		title = "Restart service"
		icon = "Terraform"
		identifier = "%s"
		blueprint_identifier = port-labs_blueprint.microservice.identifier
		trigger = "DAY-2"
		invocation_method {
			type = "KAFKA"
		}
		user_properties {
			identifier = "clear_cache"
			type = "boolean"
			title = "Clear cache"
		}
	}
`, identifier, actionIdentifier)
	testAccActionConfigUpdate := fmt.Sprintf(`
	resource "port-labs_blueprint" "microservice" {
		title = "TF test microservice"
		icon = "Terraform"
		identifier = "%s"
		properties {
			identifier = "text"
			type = "string"
			title = "text"
		}
	}
	resource "port-labs_action" "restart_microservice" {
		title = "Restart service"
		icon = "Terraform"
		identifier = "%s"
		blueprint_identifier = port-labs_blueprint.microservice.identifier
		trigger = "DAY-2"
		invocation_method {
			type = "KAFKA"
		}
		user_properties {
			identifier = "clear_cache"
			type = "string"
			required = true
			title = "Clear cache"
			enum = ["yes", "no"]
		}
		user_properties {
			identifier = "submit_report"
			type = "boolean"
			title = "Submit report"
		}
	}
`, identifier, actionIdentifier)
	resource.Test(t, resource.TestCase{
		Providers: map[string]*schema.Provider{
			"port-labs": Provider(),
		},
		Steps: []resource.TestStep{
			{
				Config:  testAccActionConfigCreate,
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "title", "Restart service"),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "blueprint_identifier", identifier),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "invocation_method.0.type", "KAFKA"),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "trigger", "DAY-2"),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "user_properties.#", "1"),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "user_properties.0.identifier", "clear_cache"),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "user_properties.0.type", "boolean"),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "user_properties.0.title", "Clear cache"),
				),
			},
			{
				Config: testAccActionConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "title", "Restart service"),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "blueprint_identifier", identifier),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "invocation_method.0.type", "KAFKA"),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "trigger", "DAY-2"),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "user_properties.#", "2"),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "user_properties.0.identifier", "clear_cache"),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "user_properties.0.type", "string"),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "user_properties.0.title", "Clear cache"),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "user_properties.0.enum.0", "yes"),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "user_properties.0.required", "true"),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "user_properties.1.identifier", "submit_report"),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "user_properties.1.type", "boolean"),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "user_properties.1.title", "Submit report"),
				),
			},
		},
	})
}

func TestAccPortActionPropMeta(t *testing.T) {
	identifier := genID()
	actionIdentifier := genID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port-labs_blueprint" "microservice" {
		title = "TF test microservice"
		icon = "Terraform"
		identifier = "%s"
		properties {
			identifier = "text"
			type = "string"
			title = "text"
		}
	}
	resource "port-labs_action" "restart_microservice" {
		title = "Restart service"
		icon = "Terraform"
		identifier = "%s"
		blueprint_identifier = port-labs_blueprint.microservice.identifier
		trigger = "DAY-2"
		invocation_method {
			type = "KAFKA"
		}
		user_properties {
			identifier = "webhook_url"
			type = "string"
			title = "Webhook URL"
			description = "Webhook URL to send the request to"
			format = "url"
			default = "https://example.com"
			pattern = "^https://.*"
		}
	}
`, identifier, actionIdentifier)
	resource.Test(t, resource.TestCase{
		Providers: map[string]*schema.Provider{
			"port-labs": Provider(),
		},
		Steps: []resource.TestStep{
			{
				Config:  testAccActionConfigCreate,
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "title", "Restart service"),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "blueprint_identifier", identifier),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "trigger", "DAY-2"),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "invocation_method.0.type", "KAFKA"),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "user_properties.#", "1"),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "user_properties.0.identifier", "webhook_url"),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "user_properties.0.type", "string"),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "user_properties.0.title", "Webhook URL"),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "user_properties.0.description", "Webhook URL to send the request to"),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "user_properties.0.format", "url"),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "user_properties.0.default", "https://example.com"),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "user_properties.0.pattern", "^https://.*"),
				),
			},
		},
	})
}

func TestAccPortActionWebhookInvocation(t *testing.T) {
	identifier := genID()
	actionIdentifier := genID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port-labs_blueprint" "microservice" {
		title = "TF test microservice"
		icon = "Terraform"
		identifier = "%s"
		properties {
			identifier = "text"
			type = "string"
			title = "text"
		}
	}
	resource "port-labs_action" "restart_microservice" {
		title = "Restart service"
		icon = "Terraform"
		identifier = "%s"
		blueprint_identifier = port-labs_blueprint.microservice.identifier
		trigger = "DAY-2"
		user_properties {
			identifier = "multiselect"
			type = "string"
			title = "multiselect"
			description = "multiselect"
			format = "entity"
			blueprint = port-labs_blueprint.microservice.identifier
		}
		invocation_method {
			type = "KAFKA"
		}
	}
`, identifier, actionIdentifier)
	resource.Test(t, resource.TestCase{
		Providers: map[string]*schema.Provider{
			"port-labs": Provider(),
		},
		Steps: []resource.TestStep{
			{
				Config:  testAccActionConfigCreate,
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "user_properties.0.blueprint", identifier),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "user_properties.0.format", "entity"),
				),
			},
		},
	})
}

func TestAccPortActionEntityMultiselect(t *testing.T) {
	identifier := genID()
	actionIdentifier := genID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port-labs_blueprint" "microservice" {
		title = "TF test microservice"
		icon = "Terraform"
		identifier = "%s"
		properties {
			identifier = "text"
			type = "string"
			title = "text"
		}
	}
	resource "port-labs_action" "restart_microservice" {
		title = "Restart service"
		icon = "Terraform"
		identifier = "%s"
		blueprint_identifier = port-labs_blueprint.microservice.identifier
		trigger = "DAY-2"
		invocation_method {
			type = "WEBHOOK"
			url = "https://google.com"
		}
	}
`, identifier, actionIdentifier)
	resource.Test(t, resource.TestCase{
		Providers: map[string]*schema.Provider{
			"port-labs": Provider(),
		},
		Steps: []resource.TestStep{
			{
				Config:  testAccActionConfigCreate,
				Destroy: false,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "title", "Restart service"),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "blueprint_identifier", identifier),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "trigger", "DAY-2"),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "invocation_method.0.type", "WEBHOOK"),
					resource.TestCheckResourceAttr("port-labs_action.restart_microservice", "invocation_method.0.url", "https://google.com"),
				),
			},
		},
	})
}
