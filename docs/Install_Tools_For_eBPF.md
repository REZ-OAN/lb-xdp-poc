# Install Necessary Tools to run and build eBPF codes

```
sudo apt-get update && \
    apt-get install -y clang llvm libelf-dev libbpf-dev libpcap-dev && \ 
    apt-get install -y gcc-multilib build-essential linux-tools-common && \
    apt-get install -y linux-headers-$(uname -r) linux-tools-$(uname -r) linux-headers-generic linux-tools-generic && \
    apt-get install -y dwarves
```