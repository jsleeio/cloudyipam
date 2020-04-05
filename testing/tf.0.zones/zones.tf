variable "zones" {
  default = {
    "us-west-2a" = "10.0.0.0/10"
    "us-west-2b" = "10.64.0.0/10"
    "us-west-2c" = "10.128.0.0/10"
    "us-west-2d" = "10.192.0.0/10"
  }
}

resource "cloudyipam_zone" "zones" {
  for_each      = var.zones
  name          = each.key
  range         = each.value
  prefix_length = 11
}

output "zones" {
  value = {
    for zone in cloudyipam_zone.zones:
      zone.name => zone.id
  }
}
