terraform {
  required_version = ">= 0.13"
  required_providers {
    yandex = {
      source = "yandex-cloud/yandex"
    }
  }
}

variable "deployto_context" {
  type = string
}

locals{
    context = yamldecode( var.deployto_context )
    provider_config = local.context["target"]["spec"]["terraform"]
}

provider "yandex" {
  zone      = local.provider_config["zone"]
  token     = local.provider_config["token"]
  cloud_id  = local.provider_config["cloud_id"]
  folder_id = local.provider_config["folder_id"]
}

resource "yandex_vpc_network" "default" {
  name = local.context["name"]
}

output "deployto_output" {
  value = yandex_vpc_network.default
}