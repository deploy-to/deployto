data "yandex_vpc_network" "default" {
  name = "default"
}

data "yandex_vpc_subnet" "default" {
  zone           = var.zone
  network_id     = data.yandex_vpc_network.default.id
}