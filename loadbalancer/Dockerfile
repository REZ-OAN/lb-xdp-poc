FROM ubuntu:24.04

WORKDIR /home/xdp_lb

RUN apt-get update 

RUN apt-get install -y clang llvm libelf-dev libbpf-dev libpcap-dev gcc-multilib build-essential make linux-tools-common

RUN apt-get install -y linux-headers-$(uname -r) linux-tools-$(uname -r) linux-headers-generic linux-tools-generic

RUN apt-get install -y curl iproute2 iputils-ping nano dwarves tcpdump bind9-dnsutils

RUN apt-get install -y jq

RUN apt-get clean

COPY ./xdp/xdp_lb.c .

COPY ./attach_xdp.sh .

RUN clang -O2 -target bpf -g -c xdp_lb.c -o xdp_lb.o

RUN chmod +x attach_xdp.sh
ENTRYPOINT [ "./attach_xdp.sh" ]