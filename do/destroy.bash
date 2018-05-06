#!/bin/bash

terraform destroy --var="do_token=$DIGITAL_OCEAN_API_KEY"
