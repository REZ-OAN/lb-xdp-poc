// added all logic and working and tested
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

__u32 serverIdx = (__u32)0;

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
  
    if(iph->saddr == client_ip_mac->ip) {
            if(serverIdx < *backend_map_size) {
                bpf_printk("client to backend\n");
                
                __u32 idx= (__u32)serverIdx;
                bpf_printk("recieved request from client IP : %d MAC : %x\n",client_ip_mac->ip,client_ip_mac->mac.addr);

                struct ip_mac *backend_ip_mac = bpf_map_lookup_elem(&backend_server_map,&idx);
                bpf_printk("serving idx %d\n",idx);
                if(!backend_ip_mac){
                    bpf_printk("server assignment failed\n");
                    return XDP_ABORTED;
                }
                bpf_printk("serving from backend-server IP  :%d, MAC : %x \n",backend_ip_mac->ip,backend_ip_mac->mac.addr);
                __builtin_memcpy(eth->h_dest, backend_ip_mac->mac.addr,ETH_ALEN);
                iph->daddr = backend_ip_mac->ip;
                serverIdx = (serverIdx + 1)%(*backend_map_size);
            } else {
                bpf_printk("idx value exceeded the max_length\n");
                return XDP_ABORTED;
            }
    } else {
        __u32 idx = (__u32)serverIdx;
        struct ip_mac *backend_ip_mac = bpf_map_lookup_elem(&backend_server_map,&idx);
        if(!backend_ip_mac){
            bpf_printk("not found any server with this idx %d\n",idx);
            return XDP_ABORTED;
        }
        bpf_printk("backend to client\n");
        __builtin_memcpy(eth->h_dest, client_ip_mac->mac.addr,ETH_ALEN);
        iph->daddr = client_ip_mac->ip;

    }
    struct ip_mac *lb_ip_mac = bpf_map_lookup_elem(&lb_mac_map,&key);
    if(!lb_ip_mac){
        bpf_printk("loadbalancing server not availabe\n");
        return XDP_ABORTED;
    }
    bpf_printk("changing source to lb-server %x\n",lb_ip_mac->mac.addr);
     __builtin_memcpy(eth->h_source, lb_ip_mac->mac.addr,ETH_ALEN);
    iph->saddr = lb_ip_mac->ip;
    // Recompute IP checksum
    iph->check = iph_csum(iph);
    
    return XDP_TX;
}

char _license[] SEC("license") = "GPL";
