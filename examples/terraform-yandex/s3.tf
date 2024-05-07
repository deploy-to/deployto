
resource "yandex_iam_service_account" "sa" {
  name = "deployto-s3-user"
}

// Назначение роли сервисному аккаунту
resource "yandex_resourcemanager_folder_iam_member" "sa-editor" {
  folder_id = var.folder_id
  role      = "storage.editor"
  member    = "serviceAccount:${yandex_iam_service_account.sa.id}"
}

// Создание статического ключа доступа
resource "yandex_iam_service_account_static_access_key" "sa-static-key" {
  service_account_id = yandex_iam_service_account.sa.id
  description        = "static access key for object storage"
}

// Создание бакета с использованием ключа
resource "yandex_storage_bucket" "test" {
  access_key = yandex_iam_service_account_static_access_key.sa-static-key.access_key
  secret_key = yandex_iam_service_account_static_access_key.sa-static-key.secret_key
  bucket     = "deployto"
}

locals {
  s3bucket = yandex_storage_bucket.test.bucket
  s3accesskey = yandex_iam_service_account_static_access_key.sa-static-key.access_key
  s3secretkey = yandex_iam_service_account_static_access_key.sa-static-key.secret_key
  s3uri = "s3://${yandex_storage_bucket.test.website_domain}/${yandex_storage_bucket.test.website_endpoint}"
}

output "s3bucket" {
  value = local.s3bucket
}

output "s3accesskey" {
  value = local.s3accesskey
  sensitive = true
}

output "s3secretkey" {
  value = local.s3secretkey
  sensitive = true
}

output "s3uri" {
  value = local.s3uri
}