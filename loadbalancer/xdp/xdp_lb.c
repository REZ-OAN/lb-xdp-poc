// added only the maps
#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/in.h>
#include <linux/ip.h>
#include <linux/tcp.h>
#include <bpf/bpf_endian.h>
#include <bpf/bpf_helpers.h>

#define ETH_ALEN 6

struct mac_addr {
    unsigned char addr[ETH_ALEN];
};

struct ip_mac {
    __u32 ip;
    struct mac_addr mac;
};


struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 256);
    __type(key, __u32);
    __type(value, struct mac_addr);
    __uint(pinning, LIBBPF_PIN_BY_NAME);
} client_mac_map SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 1);
    __type(key, __u32);
    __type(value, struct mac_addr);
    __uint(pinning, LIBBPF_PIN_BY_NAME);
} lb_mac_map SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_ARRAY);
    __uint(max_entries, 128);
    __type(key, unsigned int);
    __type(value, struct ip_mac);
    __uint(pinning, LIBBPF_PIN_BY_NAME);
} backend_server_map SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 1);
    __type(key, unsigned int);
    __type(value, unsigned int);
    __uint(pinning, LIBBPF_PIN_BY_NAME);
} bs_map_size_map SEC(".maps");

static __always_inline __u16 csum_fold_helper(__u64 csum)
{
#pragma unroll
    for (int i = 0; i < 4; i++) {
        if (csum >> 16) {
            csum = (csum & 0xffff) + (csum >> 16);
        }
    }
    return ~csum;
}

static __always_inline __u16 iph_csum(struct iphdr *iph)
{
    iph->check = 0;
    unsigned long long csum = bpf_csum_diff(0, 0, (__be32 *)iph, sizeof(struct iphdr), 0);
    return csum_fold_helper(csum);
}

unsigned int backend_server_index = 0;

SEC("xdp")
int xdp_load_balancer(struct xdp_md *ctx)
{
    // void *data_end = (void *)(long)ctx->data_end;
    // void *data = (void *)(long)ctx->data;
    // struct ethhdr *eth = data;
    
    // if ((void *)(eth + 1) > data_end) {
    //     return XDP_ABORTED;
    // }

    // if (eth->h_proto != bpf_htons(ETH_P_IP)) {
    //     return XDP_PASS;
    // }

    // struct iphdr *iph = data + sizeof(struct ethhdr);
    // if ((void *)(iph + 1) > data_end) {
    //     return XDP_ABORTED;
    // }

    // if (iph->protocol != IPPROTO_TCP) {
    //     return XDP_PASS;
    // }

    // struct tcphdr *tcph = (void *)iph + sizeof(*iph);
    // if ((void *)(tcph + 1) > data_end) {
    //     return XDP_ABORTED;
    // }

    // if (iph->saddr == client_ip) {
    //     // Client to backend
    //     bpf_printk("Client to Backend - Index: %d\n", backend_server_index);

    //     if (backend_server_index < sizeof(server_mac) / ETH_ALEN) {
    //         __builtin_memcpy(eth->h_dest, server_mac[backend_server_index], ETH_ALEN);
    //         iph->daddr = server_ip[backend_server_index];
    //         backend_server_index = (backend_server_index + 1) % (sizeof(server_mac) / ETH_ALEN);
    //     }
    // } else {
    //     // Backend to client
    //     bpf_printk("Backend to Client\n");

    //     __builtin_memcpy(eth->h_dest, client_mac, ETH_ALEN);
    //     iph->daddr = client_ip;
    // }

    // __builtin_memcpy(eth->h_source, load_balancer_mac, ETH_ALEN);
    // iph->saddr = load_balancer_ip;

    // // Recompute IP checksum
    // iph->check = iph_csum(iph);

    return XDP_PASS;
}

char _license[] SEC("license") = "GPL";
