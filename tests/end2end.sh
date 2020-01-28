#!/usr/bin/env bash

[[ "travis" == $(echo "travis" | /tmp/passcrux split --ratio 3/5 -e abc --stdin | tail -3 | /tmp/passcrux combine -e abc --stdin) ]]
