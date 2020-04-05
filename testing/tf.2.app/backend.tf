terraform {
  backend "etcdv3" {
    endpoints = [ "etcd:2379" ]
  }
}

variable "component" {
  type = string
}
