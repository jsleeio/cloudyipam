# cloudyipam

This repository contains an extremely barebones IP address management (IPAM)
service and some related parts. When deploying infrastructure in modern cloud
environments such as AWS, often there is a requirement to allocate and then
subdivide IP address ranges.

`cloudyipam` attempts to provide this address range allocation such that it can
be used with extremely minimal pre-existing infrastructure.  It manages only
two kinds of resources: _zones_ and _subnets_. A _zone_ is an IP address range
that is intended to be divided into _subnets_ of equal size.

## what's included

* `database/cloudyipam.sql`: the IPAM service, implemented in PostgreSQL
  as a collection of stored procedures/functions
* `pkg/cloudyipam`: a thin "API client" library that interacts with the
  PostgreSQL components
* `cmd/cloudyipam`: a command-line client, built with `pkg/cloudyipam`
* `cmd/terraform-provider-cloudyipam`: a Terraform provider, built with
  `pkg/cloudyipam`
* `testing`: end-to-end test infrastructure using the Terraform provider
  and Docker Compose

## demo

Assuming you have Docker installed:

    $ git clone git@github.com:jsleeio/cloudyipam.git
    $ make
    $ testing/compose-up.sh

This demo uses the included Terraform provider to create some zones, ostensibly
one per availability zone in AWS's Oregon region, and then to allocate two sets
of subnets within those zones.

Once Terraform is done running the end-to-end test, you should be left with the
other containers still running.  One of them should have the `cloudyipam` CLI
tool available in it. This can be used to see the zones that Terraform created:

    $ docker exec -it testing_cloudyipam-cli_1 cloudyipam zone list
    Using config file: /root/.cloudyipam.yaml
    NAME		    ZONE ID					                      CIDR BLOCK	  SUBNET PREFIXLEN
    us-west-2a	feacde2e-5576-4e43-9dd3-f74cffe79052	10.0.0.0/10	  11
    us-west-2b	cd4062b1-8c7d-471d-a5d2-aa7d57de4a3a	10.64.0.0/10	11
    us-west-2c	0ba18fb9-2ef4-4c19-9126-32a1e4baf2ac	10.128.0.0/10	11
    us-west-2d	b1eef029-916f-4108-bd21-04ad82533e63	10.192.0.0/10	11

And the subnets:

    $ docker exec -it testing_cloudyipam-cli_1 cloudyipam subnet list
    Using config file: /root/.cloudyipam.yaml
    USAGE			        SUBNET ID				                      ZONE ID					                      CIDR BLOCK	  AVAILABLE
    us-west-2c (db)		cd64c8ca-9c7c-43e3-8da7-0aef9c53039f	0ba18fb9-2ef4-4c19-9126-32a1e4baf2ac	10.128.0.0/11	false
    us-west-2c (app)	abcccbee-01e3-4256-bded-d9abc3491b1d	0ba18fb9-2ef4-4c19-9126-32a1e4baf2ac	10.160.0.0/11	false
    us-west-2d (db)		4d939e2b-8fdf-48e8-b263-0944361c0642	b1eef029-916f-4108-bd21-04ad82533e63	10.192.0.0/11	false
    us-west-2d (app)	caf19003-61c6-4a18-b181-49eeed54e5a9	b1eef029-916f-4108-bd21-04ad82533e63	10.224.0.0/11	false
    us-west-2b (db)		4e79b1f2-5771-48e5-9ccf-065ab97d6fc1	cd4062b1-8c7d-471d-a5d2-aa7d57de4a3a	10.64.0.0/11	false
    us-west-2b (app)	9a51ba6f-06d6-4398-ab5e-b53f4d506626	cd4062b1-8c7d-471d-a5d2-aa7d57de4a3a	10.96.0.0/11	false
    us-west-2a (db)		1e6b9dce-decd-4ac0-8638-a5d5f6470551	feacde2e-5576-4e43-9dd3-f74cffe79052	10.0.0.0/11	  false
    us-west-2a (app)	c03f1eb1-b21b-4ef1-9e5e-71249e0600c3	feacde2e-5576-4e43-9dd3-f74cffe79052	10.32.0.0/11	false

To tear it all down:

    $ testing/compose-down.sh

## Terraform provider

Have a look at the examples in `testing/tf.*` and note in particular that the
zones and subnets are not allocated within the same Terraform state. Instead,
when the subnets are allocated using the `cloudyipam_subnet` resource,
Terraform remote state is interrogated to acquire the zone identifiers.

Anecdotally, this appears to mirror how people use Terraform in production, and
also acts as a workaround for a `terraform plan` behaviour when one resource
declaration (that is using the `for_each` directive) attempts to use the value
of a computed (ie. not known until apply) attribute of another resource.
Specifically, you need the zone ID to allocate a subnet within that zone.

## random notes/todo

* resource updates are not supported
* probably this repository should be split up
* documentation is very thin
* the CLI client is not intended as a general purpose tool; it was mostly
  for me to exercise the API during development. But the framework is there
  to make a good tool, I guess?
