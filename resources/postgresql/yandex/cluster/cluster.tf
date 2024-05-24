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

data "yandex_vpc_subnet" "default" {
  name = local.context["subnet"]
}

resource "yandex_mdb_postgresql_cluster" "default" {
  name        = local.context["name"]
  environment = "PRESTABLE"
  network_id  = data.yandex_vpc_network.default.id

  config {
    version = 15
    resources {
      resource_preset_id = "s2.micro"
      disk_type_id       = "network-ssd"
      disk_size          = 10
    }
    postgresql_config = {
      max_connections                   = 395
      enable_parallel_hash              = true
      vacuum_cleanup_index_scale_factor = 0.2
      autovacuum_vacuum_scale_factor    = 0.34
      default_transaction_isolation     = "TRANSACTION_ISOLATION_READ_COMMITTED"
      shared_preload_libraries          = "SHARED_PRELOAD_LIBRARIES_AUTO_EXPLAIN,SHARED_PRELOAD_LIBRARIES_PG_HINT_PLAN"
    }
  }

  maintenance_window {
    type = "WEEKLY"
    day  = "SAT"
    hour = 12
  }

  host {
    zone      = local.context["zone"]
    subnet_id = data.yandex_vpc_subnet.default.id
    assign_public_ip = local.context["assign_public_ip"]
  }
}

output "deployto_output" {
  value = yandex_mdb_postgresql_cluster.default
  sensitive = true
}
