# CUSTOM LOADBALANCER WITH XDP

## Table Of Contents
 - [Prerequisite](#prerequisite)
 - [Clone Repository](#step-1-clone-the-repository)
 - [Build Images](#step-2-build-necessary-docker-images)
 - [Generate Objectfile](#step-3-generate-necessary-files-using-bpf2go)
 - [Build **lb-backend**](#step-4-build-the-lb-backend-binary)
 - [Install Packages For Frontend](#step-5-install-the-necessary-packages)
 - [Hands On](#step-6-hands-on)
## Prerequisite

 - [Install Docker](https://docs.docker.com/engine/install/ubuntu/)
 - [Instal Go](https://go.dev/doc/install)
 - [Install Node and Npm](https://nodejs.org/en/download/package-manager)
 - [Install make](https://www.geeksforgeeks.org/how-to-install-make-on-ubuntu/)
 - [Install necessary tools to run eBPF code](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/docs/Install_Tools_For_eBPF.md)

## Step-1 (clone the repository)
```
 git clone https://github.com/REZ-OAN/lb-xdp-poc.git
```
## Step-2 (Build Necessary Docker Images)
Navigate to the `lb-xdp-poc` directory (**root** directory for the application)

 - For CLIENT-SERVER
```
make build_client
```
 - For LOADBALANCER-SERVER
 ```
 make build_lb
 ```
 [How i installed eBPF on docker image?](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/docs/How_To_eBPF_in_Docker.md)
 - For SERVER-BACKEND
```
make build_server
```
## Step-3 (Generate Necessary Files using bpf2go)
To interact with **bpf_maps** we have to convert the `bpf` code into go and object file.`github.com/cilium/ebpf/cmd/bpf2go` this module helps us to do this.Navigate to `lb-backend`.

To generate execute the following command : 
```
go generate
```
This will generate the necessary files for you.

You Will See these logs :

![go_generate_logs](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/generate.png)


## Step-4 (Build The lb-backend Binary)
Navigate to `lb-backend`. To **build** execute the following command :
```
go build
```
## Step-5 (Install The Necessary Packages)
Navigate to `lb-frontend`. To **install** necessary packages execute the following command:
```
npm i
```
## Step-6 (Hands On)
To see the hands on demonstration visit [HandsOn load_balancer_xdp](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/docs/HandsOn.md)
