
variable "deployto_context" {
  type = object({
    name    = string
    address = string
  })
  sensitive = true
}

variable "YC_CLOUD_ID" {
  description = ""
  type        = string
}

variable "YC_FOLDER_ID" {
  description = ""
  type        = string
}

variable "YC_TOKEN" {
  description = "token"
  type        = string
  sensitive = true
}

variable "YC_ZONE" {
  description = "zone name"
  type        = string
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
  zone = var.YC_ZONE
  token     = var.YC_TOKEN
  cloud_id  = var.YC_CLOUD_ID
  folder_id = var.YC_FOLDER_ID
}
