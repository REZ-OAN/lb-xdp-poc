# CUSTOM LOADBALANCER WITH XDP

## Prerequisite

 - [Install Docker]()
 - [Instal Go]()
 - [Install Node and Npm]()
 - [Install make]()
 - [Install necessary tools to run eBPF code]()

## Step-1 (Build Necessary Images)
Navigate to the `lb-xdp-poc` directory (**root** directory for the application)

 - For CLIENT-SERVER
```
make build_client
```
 - For LOADBALANCER-SERVER
 ```
 make build_lb
 ```
 - For SERVER-BACKEND
```
make build_server
```
## Step-2 (Generate Necessary Files using bpf2go)
To interact with **bpf_maps** we have to convert the `bpf` code into go and object file.`github.com/cilium/ebpf/cmd/bpf2go` this module helps us to do this.Navigate to `lb-backend`.

To generate execute the following command : 
```
go generate
```
This will generate the necessary files for you.

## Step-3 (Build The lb-backend Binary)
Navigate to `lb-backend`. To **build** execute the following command :
```
go build
```
## Step-4 (Install The Necessary Packages)
Navigate to `lb-frontend`. To **install** necessary packages execute the following command:
```
npm i
```
## Step-5 (Hands On)
 - [Hands On Demonstration]()