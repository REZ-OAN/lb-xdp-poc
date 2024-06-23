// added all logic
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
    __uint(max_entries, 1);
    __type(key, __u32);
    __type(value, struct ip_mac);
    __uint(pinning, LIBBPF_PIN_BY_NAME);
} client_mac_map SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 1);
    __type(key, __u32);
    __type(value, struct ip_mac);
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
    __type(key, __u32);
    __type(value,__u32);
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

__u32 backend_server_index = 0;

SEC("xdp")
int xdp_load_balancer(struct xdp_md *ctx)
{

    void *data_end = (void *)(long)ctx->data_end;
    void *data = (void *)(long)ctx->data;
    struct ethhdr *eth = data;
    
    if ((void *)(eth + 1) > data_end) {
        bpf_printk("XDP_ABORTED: Invalid eth pointer\n");
        return XDP_ABORTED;
    }

    if (eth->h_proto != bpf_htons(ETH_P_IP)) {
        return XDP_PASS;
    }

    struct iphdr *iph = data + sizeof(struct ethhdr);
    if ((void *)(iph + 1) > data_end) {
        bpf_printk("XDP_ABORTED: Invalid iph pointer\n");
        return XDP_ABORTED;
    }

    if (iph->protocol != IPPROTO_TCP) {
        return XDP_PASS;
    }

    struct tcphdr *tcph = (void *)iph + sizeof(*iph);
    if ((void *)(tcph + 1) > data_end) {
        bpf_printk("XDP_ABORTED: Invalid tcph pointer\n");
        return XDP_ABORTED;
    }

    __u32 size_key = 0;
    __u32 *backend_map_size = bpf_map_lookup_elem(&bs_map_size_map, &size_key);

    if (!backend_map_size) {
        bpf_printk("XDP_ABORTED: backend_map_size is NULL\n");
        return XDP_ABORTED;
    }

    __u32 key = 0;
    struct ip_mac *client_ip_mac = bpf_map_lookup_elem(&client_mac_map, &key);

    if (!client_ip_mac) {
        bpf_printk("XDP_ABORTED: client_ip_mac is NULL\n");
        return XDP_ABORTED;
    }

    if (iph->saddr == client_ip_mac->ip) {
        // Client to backend
        bpf_printk("Client to Backend - Index: %d\n", backend_server_index);

        if (backend_server_index < *backend_map_size) {
            struct ip_mac *server_ip_mac = bpf_map_lookup_elem(&backend_server_map, &backend_server_index);
            if (server_ip_mac) {
                __builtin_memcpy(eth->h_dest, server_ip_mac->mac.addr, ETH_ALEN);
                iph->daddr = server_ip_mac->ip;
                backend_server_index = (backend_server_index + 1) % (*backend_map_size);
            } else {
                bpf_printk("XDP_ABORTED: server_ip_mac is NULL\n");
                return XDP_ABORTED;
            }
        } else {
            bpf_printk("XDP_ABORTED: Invalid backend_server_index\n");
            return XDP_ABORTED;
        }
    } else {
        bpf_printk("Backend to Client\n");
        struct ip_mac *client_ip_mac2 = bpf_map_lookup_elem(&client_mac_map, &key);
        if(client_ip_mac2){
        __builtin_memcpy(eth->h_dest, client_ip_mac->mac.addr, ETH_ALEN);
        iph->daddr = client_ip_mac->ip;
        }else {
            bpf_printk("client not found\n");
        }
    }

    struct ip_mac *lb_ip_mac = bpf_map_lookup_elem(&lb_mac_map, &key);
    if (lb_ip_mac) {
        __builtin_memcpy(eth->h_source, lb_ip_mac->mac.addr, ETH_ALEN);
        iph->saddr = lb_ip_mac->ip;
    } else {
        bpf_printk("XDP_ABORTED: lb_ip_mac is NULL\n");
        return XDP_ABORTED;
    }

    // Recompute IP checksum
    iph->check = iph_csum(iph);
    
    return XDP_TX;
}

char _license[] SEC("license") = "GPL";
