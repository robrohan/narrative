/*
---
title: TheRohans Infrastructure
author: "The Dev Squad"
date: 2022-04-08
rights: Copyright (c) Company, Pty Ltd
lang: en-GB
toc: true
papersize: A4
fontfamily: mathptmx
fontsize: 11pt
geometry:
- top=30mm
- left=20mm
- right=20mm
- bottom=30mm
abstract: |
	"This is an abstract"
numbersections: true
autoEqnLabels: true

theme: Berlin
colortheme: seahorse
---

# Introduction

We are using Terraform and AWS to host our applications. Static applications are simply hosted on S3, and for compute workloads we're using Kubernetes hosted on EKS.

## Provider Setup

We need to tell terraform we are using AWS. We do this using the provider call. We also set the region most of our things are going to be running in. The region value is sometimes overridden depending on the service we are using.

The region value could be a variable, but we are currently hard coding it.

*/
provider "aws" {
  region = "us-east-1"
}
/*

## Terraform State

Next, when we run terraform, to apply state we need to setup a place to store the state. The state file is used to compare our last known state with what we are asking terraform to do. There are several places to store state, but for our purposes, we are storing state in an S3 bucket.

---

**Note**: This bucket needs to have been created by hand before the state can be stored.

---

If you want to destroy everything, you first need to comment out the s3 backend, and then recreate the state locally (using init). Once you have the state locally, you can then destroy everything.

*/
terraform {
  backend "s3" {
    bucket  = "therohans.com-terraform-state"
    key     = "global/s3/terraform.tfstate"
    region  = "us-east-1"
    encrypt = true
  }
}
/*

We also make sure that our state file is encrypted using AES265. Additionally, we set the _prevent\_destory_ _lifecycle_ property so the bucket will not get deleted automatically.

*/
resource "aws_s3_bucket" "terraform_state" {
  bucket = "therohans.com-terraform-state"
  lifecycle {
    prevent_destroy = true
  }

  versioning {
    enabled = true
  }

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }
    }
  }
}

// Inlucde Buckets
module "therohans_s3_buckets" {
  source = "./modules/s3-buckets"
}

// Include EC2 Docker
module "therohans_ec2_docker" {
  source = "./modules/ec2-basic"
}

// Include EKS Setup
//module "therohans_eks_cluster" {
//	source = "./modules/eks-cluster"
//}

/*

[^Former2]: Tool to try to export infrasturcture created outside of terraform: https://former2.com/#section-outputs-tf

*/
