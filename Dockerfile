FROM golang:1.10.4

EXPOSE 5000
EXPOSE 5100

ADD node/node /usr/bin/quanta-bridge
ADD cli/bitcoin /usr/bin/bitcoin_sync
ADD cli/ethereum /usr/bin/ethereum_sync
RUN ["chmod", "+x", "/usr/bin/quanta-bridge","/usr/bin/*_sync"]

ENTRYPOINT ["/usr/bin/quanta-bridge", "-config", "/data/crosschain.yml"]
