FROM ubuntu:22.04

RUN apt-get update && \
    apt-get install -y iproute2 iputils-ping tcpdump curl bind9-dnsutils && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*