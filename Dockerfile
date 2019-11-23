FROM golang:1.12.1-alpine3.9 as builder

USER root:root

RUN apk add --no-cache gcc musl-dev git git-lfs
RUN go get -u github.com/golang/dep/cmd/dep
RUN go get -u firebase.google.com/go
WORKDIR $GOPATH/src/github.com/quantadex

ARG token
RUN git config --global url."https://$token@github.com/quantadex".insteadOf "https://github.com/quantadex"
RUN git clone --depth=1 --single-branch --branch graphene https://github.com/quantadex/distributed_quanta_bridge

WORKDIR $GOPATH/src/github.com/quantadex/distributed_quanta_bridge
RUN tar xvf vendor.tar
RUN git reset --hard
RUN ./build.sh

FROM alpine:3.7

EXPOSE 5000
EXPOSE 5100

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
ENV OUTDIR=/go/src/github.com/quantadex/distributed_quanta_bridge

COPY --from=builder $OUTDIR/node/node /usr/bin/quanta-bridge
COPY --from=builder $OUTDIR/cli/bitcoin/bitcoin /usr/bin/bitcoin_sync
COPY --from=builder $OUTDIR/cli/ethereum/ethereum /usr/bin/ethereum_sync
COPY --from=builder $OUTDIR/cli/litecoin/litecoin /usr/bin/litecoin_sync
COPY --from=builder $OUTDIR/cli/bch/bch /usr/bin/bch_sync
COPY --from=builder $OUTDIR/cli/event_notifier/event_notifier /usr/bin/webhook_process

RUN ["chmod", "+x", "/usr/bin/quanta-bridge"]
RUN ["chmod", "+x", "/usr/bin/ethereum_sync"]
RUN ["chmod", "+x", "/usr/bin/bitcoin_sync"]
RUN ["chmod", "+x", "/usr/bin/litecoin_sync"]
RUN ["chmod", "+x", "/usr/bin/bch_sync"]
RUN ["chmod", "+x", "/usr/bin/webhook_process"]

# ENTRYPOINT ["/usr/bin/quanta-bridge", "-config", "/data/crosschain.yml"]
