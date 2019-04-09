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

# GPG keys required by litecoin
RUN set -ex \
      && for key in \
        B42F6819007F00F88E364FD4036A9C25BF357DD4 \
        FE3348877809386C \
      ; do \
        gpg --keyserver pgp.mit.edu --recv-keys "$key" || \
        gpg --keyserver keyserver.pgp.com --recv-keys "$key" || \
        gpg --keyserver ha.pool.sks-keyservers.net --recv-keys "$key" || \
        gpg --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys "$key" ; \
      done

ENV LITECOIN_VERSION=0.16.3
RUN curl -O https://download.litecoin.org/litecoin-${LITECOIN_VERSION}/linux/litecoin-${LITECOIN_VERSION}-x86_64-linux-gnu.tar.gz \
  && curl https://download.litecoin.org/litecoin-${LITECOIN_VERSION}/linux/litecoin-${LITECOIN_VERSION}-linux-signatures.asc | gpg --verify - \
  && tar --strip=2 -xzf *.tar.gz -C /usr/local/bin \
  && rm *.tar.gz

ADD node/node /usr/bin/quanta-bridge
ADD cli/bitcoin/bitcoin /usr/bin/bitcoin_sync
ADD cli/ethereum/ethereum /usr/bin/ethereum_sync
ADD cli/litecoin/litecoin /usr/bin/litecoin_sync

RUN ["chmod", "+x", "/usr/bin/quanta-bridge"]
RUN ["chmod", "+x", "/usr/bin/ethereum_sync"]
RUN ["chmod", "+x", "/usr/bin/bitcoin_sync"]
RUN ["chmod", "+x", "/usr/bin/litecoin_sync"]

# ENTRYPOINT ["/usr/bin/quanta-bridge", "-config", "/data/crosschain.yml"]
