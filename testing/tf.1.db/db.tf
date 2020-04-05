data "terraform_remote_state" "zones" {
  backend = "etcdv3"
  config = {
    endpoints = [ "etcd:2379" ]
    prefix = "tf.0.zones/"
  }
}

resource "cloudyipam_subnet" "db" {
  for_each = data.terraform_remote_state.zones.outputs["zones"]
  usage    = "${each.key} (db)"
  zone_id  = each.value
}

output "subnet_cidr_blocks" {
  value = [ for subnet in cloudyipam_subnet.db: subnet.range ]
}
