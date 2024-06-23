# CUSTOM LOADBALANCER WITH XDP

## Prerequisite

 - [Install Docker]()
 - [Instal Go]()
 - [Install Node and Npm]()
 - [Install make]()
 - [Install necessary tools to run eBPF code]()

## Step-1 (Build Necessary Images)

 - For CLIENT-SERVER
```
# navigate to the lb-xdp-poc directory (root directory for the application)
make build_client
```