terraform {
  required_providers {
    yandex = {
      # source = "yandex-cloud/yandex"
      source = "terraform-registry.storage.yandexcloud.net/yandex-cloud/yandex" # Alternate link
    }
  }
  required_version = ">= 0.13"
}

