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
  name = local.context["cluster"]
}

resource "yandex_mdb_postgresql_database" "default" {
  cluster_id = data.yandex_mdb_postgresql_cluster.default.id
  name       = local.context["name"]
  owner      = local.context["user"]
  lc_collate = "en_US.UTF-8"
  lc_type    = "en_US.UTF-8"
  extension {
    name = "uuid-ossp"
  }
  extension {
    name = "xml2"
  }
}

output "deployto_output" {
  value = yandex_mdb_postgresql_database.default
  sensitive = true
}
