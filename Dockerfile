FROM golang:1.12.1-alpine3.9 as builder

USER root:root

RUN apk add --no-cache gcc musl-dev git
RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR $GOPATH/src/github.com/quantadex

ARG token
RUN git config --global url."https://$token@github.com/quantadex".insteadOf "https://github.com/quantadex"
RUN git clone --single-branch --branch graphene https://github.com/quantadex/distributed_quanta_bridge

WORKDIR $GOPATH/src/github.com/quantadex/distributed_quanta_bridge
RUN dep ensure --vendor-only
RUN git reset --hard
RUN ./build.sh

FROM alpine:3.7

EXPOSE 5000
EXPOSE 5100

ENV OUTDIR=/go/src/github.com/quantadex/distributed_quanta_bridge

COPY --from=builder $OUTDIR/node/node /usr/bin/quanta-bridge
COPY --from=builder $OUTDIR/cli/bitcoin/bitcoin /usr/bin/bitcoin_sync
COPY --from=builder $OUTDIR/cli/ethereum/ethereum /usr/bin/ethereum_sync
COPY --from=builder $OUTDIR/cli/litecoin/litecoin /usr/bin/litecoin_sync
COPY --from=builder $OUTDIR/cli/bch/bch /usr/bin/bch_sync

RUN ["chmod", "+x", "/usr/bin/quanta-bridge"]
RUN ["chmod", "+x", "/usr/bin/ethereum_sync"]
RUN ["chmod", "+x", "/usr/bin/bitcoin_sync"]
RUN ["chmod", "+x", "/usr/bin/litecoin_sync"]
RUN ["chmod", "+x", "/usr/bin/bch_sync"]

# ENTRYPOINT ["/usr/bin/quanta-bridge", "-config", "/data/crosschain.yml"]
