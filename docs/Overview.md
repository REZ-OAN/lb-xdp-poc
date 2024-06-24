# Overview

## Application Folder Structure
```
----(root)lb-xdp-poc
            |
            |
            |----(subdir)client
            |               |
            |               |--Dockerfile # this is for the client Docker Image
            |
            |----(subdir)docs
            |               |
            |               |-- # this folder has the related documentations
            |
            |----(subdir)images
            |               |
            |               |-- # this folder has the images referenced inthe docs
            |
            |----(subdir)lb-backend
            |               |
            |               |-- # here are the codes for our Applications Backend
            |
            |----(subdir)lb-frontend
            |               |
            |               |-- # here are the codes for our Applications Frontend
            |
            |----(subdir)loadbalancer
            |               |
            |               |-- Dockerfile # this is for the load_balancing servers Docker Image 
            |               |
            |               |--(subdir)xdp # here are the eBPF (xdp) codes
            |
            |----(subdir)server
            |               |
            |               |-- Dockerfile # this is for the resquest_handling servers Docker Image
            |               |               serves the response to the client 
            |
            |---- Makefile # here we have wirte some script to automate our process a little bit
            |
            |---- README.md # a step wise explaination of whole setup

```
## Overview Of the System API's

For route `/api/create_net`
```
# payload
{
    subnet:<subnet(string)>
}
```
![/api/create_net overview](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/api-create-net.png)


For route `/api/launch_client`
```
# payload
{
    ip:<ip_addr(string)>
}
```
![/api/launch_client overview](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/api-launch-client.png)


For route `/api/launch_lb`
```
# payload
{
    ip:<ip_addr(string)>
}
```
![/api/launch_lb overview](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/api-launch-lb.png)


For route `/api/launch_server`
```
# payload
{
    ip:<ip_addr(string)>
}
```
![/api/launch_server overview](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/api-launch-server.png)


After launching the first server, launching rest of the servers only updates the map. Other process are same as shown in the figure.

## How it works

![working_overview](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/working.png)