#!/usr/bin/env bash

[[ "travis" == $(echo "travis" | passcrux split --ratio 3/5 -e abc --stdin | tail -3 | passcrux combine -e abc --stdin) ]]
