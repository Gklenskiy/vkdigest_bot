FROM umputun/baseimage:buildgo-latest as build

ARG SKIP_TEST

ENV GOFLAGS="-mod=vendor"

COPY . /build/vkdigest
WORKDIR /build/vkdigest

RUN \
    if [ -z "$SKIP_TEST" ] ; then \
    go test -timeout=30s  ./... ;\
    else echo "skip tests" ; fi

RUN CGO_ENABLED=0 GOOS=linux go build -o vkdigest -ldflags "-s -w" ./app

FROM alpine:latest

RUN apk --no-cache add curl bash

COPY --from=build /build/vkdigest/vkdigest /srv/vkdigest
COPY --from=build /build/vkdigest/var /srv/var

WORKDIR /srv

RUN adduser -D user
USER user

CMD ["/srv/vkdigest", "bot"]