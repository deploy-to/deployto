
variable "deployto_context" {
  type = map(string)
}

terraform {
  required_version = ">= 0.13"
  required_providers {
    yandex = {
      source = "yandex-cloud/yandex"
    }
  }
}

provider "yandex" {
  zone = var.deployto_context["yc_zone"]
  token     = var.deployto_context["yc_token"]
  cloud_id  = var.deployto_context["yc_cloud_id"]
  folder_id = var.deployto_context["yc_folder_id"]
}
