resource "yandex_mdb_postgresql_cluster" "postgresql-single" {
  name        = "test"
  environment = "PRESTABLE"
  network_id  = yandex_vpc_network.default.id

  config {
    version = 12
    resources {
      resource_preset_id = "s2.micro"
      disk_type_id       = "network-ssd"
      disk_size          = 5
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
    zone      = var.zone
    subnet_id = yandex_vpc_subnet.k8s-mdb-subnet.id
  }
}

resource "yandex_mdb_postgresql_database" "foo" {
  cluster_id = yandex_mdb_postgresql_cluster.postgresql-single.id
  name       = "testdb"
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
  name       = "user_name"
  password   = "your_password"
}

locals {
  dbuser = yandex_mdb_postgresql_user.test.name
  dbpassword = yandex_mdb_postgresql_user.test.password
  dbhosts = yandex_mdb_postgresql_cluster.postgresql-single.host.*.fqdn
  dbname = yandex_mdb_postgresql_database.foo.name
  dburi = "postgresql://${local.dbuser}:${local.dbpassword}@:1/${local.dbname}"
}

output "dbuser" {
  value = local.dbuser
}

output "dbpassword" {
  value = local.dbpassword
  sensitive = true
}

output "dbhosts" {
  value = local.dbhosts[0]
}

output "dbname" {
  value = local.dbname
}

output "dburi" {
  value = local.dburi
  sensitive = true
}