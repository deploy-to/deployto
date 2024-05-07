

variable "cloud_id" {
  description = ""
  type        = string
}

variable "folder_id" {
  description = ""
  type        = string
}

variable "token" {
  description = "token"
  type        = string
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
  token     = var.token
  cloud_id  = var.cloud_id
  folder_id = var.folder_id
}
