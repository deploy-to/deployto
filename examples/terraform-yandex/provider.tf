

variable "yc_cloud_id" {
  description = ""
  type        = string
}

variable "yc_folder_id" {
  description = ""
  type        = string
}

variable "yc_token" {
  description = "token"
  type        = string
  sensitive = true
}

variable "zone" {
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
  zone = var.zone
  token     = var.yc_token
  cloud_id  = var.yc_cloud_id
  folder_id = var.yc_folder_id
}
