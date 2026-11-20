FROM debian:12
RUN apt-get update && apt-get install -y --no-install-recommends jq ca-certificates curl && rm -rf /var/lib/apt/lists/*
COPY libwasmvm.x86_64.so /usr/lib/
COPY axmd /usr/local/bin/axmd
COPY init-chain.sh /usr/local/bin/init-chain.sh
COPY referrals.json /fixtures/referrals.json
RUN chmod +x /usr/local/bin/axmd /usr/local/bin/init-chain.sh
