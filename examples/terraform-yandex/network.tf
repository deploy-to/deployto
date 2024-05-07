data "yandex_vpc_network" "default" {
  name = "default"
}

resource "yandex_vpc_subnet" "k8s-mdb-subnet" {
  zone           = var.zone
  network_id     = data.yandex_vpc_network.default.id
  v4_cidr_blocks = ["10.5.0.0/24"]
}