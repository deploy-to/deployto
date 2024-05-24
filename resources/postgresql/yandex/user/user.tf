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

data "yandex_mdb_postgresql_cluster" "default" {
  name        = local.context["cluster"]
}

resource "yandex_mdb_postgresql_user" "default" {
  cluster_id = data.yandex_mdb_postgresql_cluster.default.id
  name       = local.context["name"]
  password   = local.context["password"]
}

output "deployto_output" {
  value = yandex_mdb_postgresql_user.default
  sensitive = true
}
