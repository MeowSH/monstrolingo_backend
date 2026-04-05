variable "db_host" {
  type    = string
  default = "localhost"
}

variable "db_port" {
  type    = string
  default = "5435"
}

variable "db_user" {
  type    = string
  default = "monstrolingo"
}

variable "db_password" {
  type    = string
  default = "monstrolingo_local_dev"
}

variable "db_name" {
  type    = string
  default = "monstrolingo"
}

env "local" {
  url = "postgres://${var.db_user}:${var.db_password}@${var.db_host}:${var.db_port}/${var.db_name}?sslmode=disable"
  dev = "docker://postgres/16/dev?search_path=public"
  migration {
    dir = "file://db/migrations"
  }
}
