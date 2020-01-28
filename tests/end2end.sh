#!/usr/bin/env bash

[[ "travis" == $(echo "travis" | $GOPATH/bin/passcrux split --ratio 3/5 -e abc --stdin | tail -3 | $GOPATH/bin/passcrux combine -e abc --stdin) ]]
