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

data "yandex_vpc_network" "default" {
  name = local.context["network"]
}

resource "yandex_vpc_subnet" "default" {
  network_id  = data.yandex_vpc_network.default.id
  v4_cidr_blocks = ["10.5.0.0/24"]
  name        = local.context["name"]
  zone        = local.context["zone"]
}

output "deployto_output" {
  value = yandex_vpc_subnet.default
}