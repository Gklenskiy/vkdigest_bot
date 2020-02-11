FROM umputun/baseimage:buildgo-latest as build

ARG SKIP_TEST

ENV GOFLAGS="-mod=vendor"

COPY . /build/vkdigest_bot
WORKDIR /build/vkdigest_bot

RUN \
    if [ -z "$SKIP_TEST" ] ; then \
    go test -timeout=30s  ./... ;\
    else echo "skip tests" ; fi

RUN CGO_ENABLED=0 GOOS=linux go build -o vkdigest_bot -ldflags "-s -w" ./app

FROM alpine:latest

RUN apk --no-cache add curl bash

COPY --from=build /build/vkdigest_bot/vkdigest_bot /srv/vkdigest_bot

WORKDIR /srv

RUN adduser -D user
USER user

CMD ["/srv/vkdigest_bot", "bot"]