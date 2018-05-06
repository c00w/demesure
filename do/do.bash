#!/bin/bash

terraform apply --var="do_token=$DIGITAL_OCEAN_API_KEY"

IP=terraform state show digitalocean_droplet.nyc1 | grep ipv4_address | awk '{ print $3 }'
demesure "$IP:8080"
