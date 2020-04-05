data "terraform_remote_state" "zones" {
  backend = "etcdv3"
  config = {
    endpoints = [ "etcd:2379" ]
    prefix = "tf.0.zones/"
  }
}

resource "cloudyipam_subnet" "app" {
  for_each = data.terraform_remote_state.zones.outputs["zones"]
  usage    = "${each.key} (app)"
  zone_id  = each.value
}

output "subnet_cidr_blocks" {
  value = [ for subnet in cloudyipam_subnet.app: subnet.range ]
}
