# cloudyipam

This repository contains an extremely barebones IP address management (IPAM)
service. When deploying infrastructure in modern cloud environments such as
AWS, often there is a requirement to allocate and then subdivide IP address
ranges.

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
  `pkg/cloudyipam`
* a Dockerfile to build a PostgreSQL container with CloudyIPAM deployed
  in it, and the CLI tool also, for experimentation

## Terraform provider

This repository used to contain `terraform-provider-cloudyipam` also.
That code now lives in its [own repository](https://github.com/jsleeio/terraform-provider-cloudyipam).

## random notes/todo

* resource updates are not supported
* documentation is very thin
* the CLI client is not intended as a general purpose tool; it was mostly
  for me to exercise the API during development. But the framework is there
  to make a good tool, I guess?
