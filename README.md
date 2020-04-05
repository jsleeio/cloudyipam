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
    NAME          ZONE ID                                 CIDR BLOCK       SUBNET PREFIXLEN
    us-west-2a    3df34b76-1575-4007-98b3-6be9b9ab7128    10.0.0.0/10      11
    us-west-2b    3b43631d-6fa0-4b8f-860c-c7a2f65a1b07    10.64.0.0/10     11
    us-west-2c    9d8ab9d7-e894-42bf-9a95-6f76ebf7aa64    10.128.0.0/10    11
    us-west-2d    1b21581a-5da6-414e-b210-1f1c779e2901    10.192.0.0/10    11

And the subnets:

    $ docker exec -it testing_cloudyipam-cli_1 cloudyipam subnet list
    Using config file: /root/.cloudyipam.yaml
    USAGE               SUBNET ID                               ZONE ID                                 CIDR BLOCK       AVAILABLE
    us-west-2d (db)     ff8e537a-5c70-44fd-a49f-41beec75e95f    1b21581a-5da6-414e-b210-1f1c779e2901    10.192.0.0/11    false
    us-west-2d (app)    85d114e5-7bf3-474d-8c23-2307358ab3c8    1b21581a-5da6-414e-b210-1f1c779e2901    10.224.0.0/11    false
    us-west-2b (db)     49073e7d-2c46-4763-8c9b-765105ea6e8d    3b43631d-6fa0-4b8f-860c-c7a2f65a1b07    10.64.0.0/11     false
    us-west-2b (app)    57c1930e-4b6c-4fc7-8b51-eabc6d8afd70    3b43631d-6fa0-4b8f-860c-c7a2f65a1b07    10.96.0.0/11     false
    us-west-2a (db)     7fd0f2df-e994-48f6-8a85-2472af6a0866    3df34b76-1575-4007-98b3-6be9b9ab7128    10.0.0.0/11      false
    us-west-2a (app)    c317fb4b-ace2-44c3-a3f3-a59070db095f    3df34b76-1575-4007-98b3-6be9b9ab7128    10.32.0.0/11     false
    us-west-2c (db)     c06a0f3d-cadc-416c-8930-cae152736d0c    9d8ab9d7-e894-42bf-9a95-6f76ebf7aa64    10.128.0.0/11    false
    us-west-2c (app)    a121b194-79cb-4217-9ecc-bdf55c93d787    9d8ab9d7-e894-42bf-9a95-6f76ebf7aa64    10.160.0.0/11    false

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
