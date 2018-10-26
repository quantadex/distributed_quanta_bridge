FROM golang:1.10.4

EXPOSE 5000
EXPOSE 5100

ADD node/node /usr/bin/quanta-bridge
RUN ["chmod", "+x", "/usr/bin/quanta-bridge"]

ENTRYPOINT ["/usr/bin/quanta-bridge", "-config", "/data/crosschain.yml"]
