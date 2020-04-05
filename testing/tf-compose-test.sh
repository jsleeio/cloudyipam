#!/bin/sh

set -e -u
top="$(git rev-parse --show-toplevel)"
cd "$top/testing"

while ! nc -v "$CLOUDYIPAM_HOSTNAME" 5432 < /dev/null ; do
  sleep 2
done

for tf in $top/testing/tf.* ; do
  cd "$tf"
  component=$(basename "$tf")
  mkdir -p terraform.d/plugins/linux_amd64
  cp /tf/testing/terraform-provider-cloudyipam terraform.d/plugins/linux_amd64/
  terraform init -input=false -backend-config="prefix=$component/"
  terraform plan -var "component=$component"
  terraform apply -auto-approve -var="component=$component"
  rm -r terraform.d/plugins/linux_amd64
done
