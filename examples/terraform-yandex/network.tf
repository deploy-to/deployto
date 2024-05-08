data "yandex_vpc_network" "default" {
  name = "default"
}

data "yandex_vpc_subnet" "default" {
  network_id     = data.yandex_vpc_network.default.id
}