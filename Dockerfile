FROM golang:alpine as build

RUN     apk update \
    &&  apk upgrade \
    &&  apk add --no-cache \
            git \
            make \
            bash

WORKDIR $GOPATH/src/github.com/xkortex/passcrux

COPY . ./

RUN go get

RUN make static

FROM build as inline_test

RUN ./tests/end2end.sh

FROM scratch

COPY --from=build /go/bin/passcrux /go/bin/passcrux

ENTRYPOINT ["/go/bin/passcrux"]
