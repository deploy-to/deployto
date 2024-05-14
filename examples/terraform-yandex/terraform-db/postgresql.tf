resource "yandex_mdb_postgresql_cluster" "postgresql-single" {
  name        = var.deployto_context["auth_hostname"]
  environment = "PRESTABLE"
  network_id  = data.yandex_vpc_network.default.id

  config {
    version = 12
    resources {
      resource_preset_id = "s2.micro"
      disk_type_id       = "network-ssd"
      disk_size          = 20
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
    zone      = var.deployto_context["yc_zone"]
    subnet_id = data.yandex_vpc_subnet.default.id
  }
}

resource "yandex_mdb_postgresql_database" "foo" {
  cluster_id = yandex_mdb_postgresql_cluster.postgresql-single.id
  name       = var.deployto_context["auth_database"]
  owner      = yandex_mdb_postgresql_user.test.name
  lc_collate = "en_US.UTF-8"
  lc_type    = "en_US.UTF-8"
  extension {
    name = "uuid-ossp"
  }
  extension {
    name = "xml2"
  }
}

resource "yandex_mdb_postgresql_user" "test" {
  cluster_id = yandex_mdb_postgresql_cluster.postgresql-single.id
  name       = var.deployto_context["auth_username"]
  password   = var.deployto_context["auth_password"]
}

output "dbuser" {
  value = yandex_mdb_postgresql_user.test.name
}

output "dbpassword" {
  value = yandex_mdb_postgresql_user.test.password
  sensitive = true
}

output "dbhosts" {
  value = yandex_mdb_postgresql_cluster.postgresql-single.host.*.fqdn[0]
}

output "dbname" {
  value = yandex_mdb_postgresql_database.foo.name
}

output "dburi" {
  value = "postgresql://${yandex_mdb_postgresql_user.test.name}:${yandex_mdb_postgresql_user.test.password}@${yandex_mdb_postgresql_cluster.postgresql-single.host.*.fqdn[0]}/${yandex_mdb_postgresql_database.foo.name}"
  sensitive = true
}

output "pgcluster" {
  value = yandex_mdb_postgresql_cluster.postgresql-single
  sensitive = true
}
