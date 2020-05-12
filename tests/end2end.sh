#!/usr/bin/env bash

errcho() {
    (>&2 echo -e "\e[31m$1\e[0m")
}

FAILED=
if [[ "travis" != $(echo "travis" | $GOPATH/bin/passcrux split --ratio 3/5 | tail -3 | $GOPATH/bin/passcrux combine ) ]]; then
  FAILED=1
fi

for key in x travis GOLANG 'h@rdp#4ssW0rd'; do
  for enc in hex abc b32 b64; do
    for sepflag in '' '-s:' '-s-'; do
      if [[ "${key}" != $(echo "${key}" | $GOPATH/bin/passcrux split --ratio 3/5 -e "${enc}" "${sepflag}" | tail -3 \
          | $GOPATH/bin/passcrux combine -e "${enc}" "${sepflag}" ) ]]; then
        FAILED=1
        errcho "Failed encoding test: encoding: ${enc}, sep: [], key: [${key}]"
      fi
    done
  done
done

if [[ -n "$FAILED" ]]; then
  echo "Failed one or more tests"
  exit 1
fi

echo "PASS!"