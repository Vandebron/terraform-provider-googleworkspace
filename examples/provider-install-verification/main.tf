terraform {
  required_providers {
    googleworkspace = {
      source = "hashicorp.com/vandebron/googleworkspace"
    }
  }
}

provider "googleworkspace" {
  credentials             = "../../credentials.json"
  impersonated_user_email = "admin@vandebron.nl"
}

data "googleworkspace_group" "example" {
  name = "everyone@vandebron.nl"
}


resource "googleworkspace_group" "test" {
  name        = "Test Terraform Group"
  email       = "test-terraform@vandebron.nl"
  description = "A group created via Terraform for testing purposes"
}
