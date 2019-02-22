FROM golang:1.10.4

EXPOSE 5000
EXPOSE 5100

ENV BITCOIN_VERSION 0.17.0
ENV BITCOIN_URL https://bitcoincore.org/bin/bitcoin-core-0.17.0/bitcoin-0.17.0-x86_64-linux-gnu.tar.gz
ENV BITCOIN_SHA256 9d6b472dc2aceedb1a974b93a3003a81b7e0265963bd2aa0acdcb17598215a4f
ENV BITCOIN_ASC_URL https://bitcoincore.org/bin/bitcoin-core-0.17.0/SHA256SUMS.asc
ENV BITCOIN_PGP_KEY 01EA5486DE18A882D4C2684590C8019E36C2E964

# install bitcoin binaries
RUN set -ex \
	&& cd /tmp \
	&& wget -qO bitcoin.tar.gz "$BITCOIN_URL" \
	&& echo "$BITCOIN_SHA256 bitcoin.tar.gz" | sha256sum -c - \
	&& gpg --batch --keyserver keyserver.ubuntu.com --recv-keys "$BITCOIN_PGP_KEY" \
	&& wget -qO bitcoin.asc "$BITCOIN_ASC_URL" \
	&& gpg --verify bitcoin.asc \
	&& tar -xzvf bitcoin.tar.gz -C /usr/local --strip-components=1 --exclude=*-qt \
	&& rm -rf /tmp/*

ADD node/node /usr/bin/quanta-bridge
ADD cli/bitcoin /usr/bin/bitcoin_sync
ADD cli/ethereum /usr/bin/ethereum_sync
RUN ["chmod", "+x", "/usr/bin/quanta-bridge","/usr/bin/ethereum_sync","/usr/bin/bitcoin_sync"]

# ENTRYPOINT ["/usr/bin/quanta-bridge", "-config", "/data/crosschain.yml"]
