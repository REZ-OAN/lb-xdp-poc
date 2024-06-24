# Hands On of The Implemented Application

## Table Of Content

 - [Overview](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/docs/Overview.md)
 - [Start Both Backend and Frontend](#step-1-starting-backend-and-frontend)
 - [Create Network](#step-2-creating-docker-network)
 - [Launch Client](#step-3-launching-client)
 - [Launch LoadBalancer](#step-4-launching-load-balancer)
 - [Launch Servers](#step-5-launching-servers)
 - [Attach XDP](#step-6-attach-xdp)
 - [Testing Custom Load Balancer](#step-7-testing-custom-load-balancer)
 - [Clean The setup](#step-8-cleanup)

## Step-1 (Starting Backend and Frontend)
Navigate to `lb-xdp-poc/lb-backend` and then execute below command :
```
sudo ./lb-backend 
```
The `-port <port_no>` flag is for give the port number to server the backend.The `default` port is **8080**.

![start lb-backend](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/start_lb-backend.png)

**Note**: In the frontend code base, at present I have hard coded the url with `port` **8080** for api calling.

Now, navigate to `lb-xdp-poc/lb-frontend` and then execute the below command :
```
npm run dev
```
By default the frontend will launch at `port` **5173**.

![start lb-frontend](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/start_lb-frontend.png)

**Note**: In the backend code base, at present I have hard coded the url with `port` **5173** for cors origin problem.

## Step-2 (Creating Docker Network)
Open the Web-App by hit this url `https://localhost:5173` from you browser. You will see a window

like below :

![overview frontend](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/frontend_overview.png)

Then fillup the **SUBNET** input filed. And **click** the `Create Subnet` button. This will create a docker network on you host machine.

 - Create Subnet

![create_network](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/create_subnet.png)

 - After Creating Subnet Frontend Overview

![after creating subnet frontend overview](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/after_creating_subnet_frontend_log.png)

 - After Creating Subnet logs from the backend

![after creating subnet backend logs](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/after_creating_subnet_backend_log.png)

## Step-3 (Launching Client)
Staying on the web-app please fill the **CLIENT IP** field on the screen and **click** on the `Launch Client` button.

This will run a container with the specified ip and the client image you build before as per our instruction.

 - Launch Client

![launching client](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/launching_client.png)

 - After Launching Client Overview of the frontend

![overview of frontend after launching client](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/after_launching_client.png)

**Note**: When launching a client, the load balancer's IP and MAC are the same as the client's due to referencing the same address.

This setup won't cause issues as our backend processes values sequentially, ensuring only one is used at a time.

 - After Launching Client logs from the backend

![logs from the backend afte launching client](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/after_launching_client_backend_logs.png)

## Step-4 (Launching Load-Balancer)
On the web-app please fill the **LOAD-BALANCER IP** field on the screen and **click** on the `Launch LoadBalancer` button.

This will run a container with the specified ip and the lb image you build before as per our instruction.

 - Launch Load-Balancer

![launching load-balancer](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/launching_lb.png)

 - After Launching LoadBalancer Overview of the frontend

![overview of the frontend after launching load-balancer](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/after_launching_lb.png)

**Note**: Now you can see that the ip's and mac's are different. There is no issues. Though we will fix the glitch sooner.

 - After Launching LoadBalancer logs from the backend

![logs of the backend after launching load-balancer](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/after_launching_lb_backend_logs.png)

## Step-5 (Launching Servers)
Now quickly fill the **SERVER IP** field on the screen and **click** on the `Launch Server` button. This will run a container with the specified ip and the server image you build before as per our instruction.

 - Launching a server

![launch a server](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/launching_a_server.png)

 - After launching a server overview of the frontend

![overview of frontend after launching a server](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/after_launching_backend_Server_table_got_updated.png)

 - After launching a server logs from the backend

![logs from the backend after launching a server](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/after_launching_backend_server_backend_logs.png)

 - After launching second server logs from the backend

![logs from the backend after launching 2nd server](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/after_launching_2nd_Server.png)

 - After launching another two servers logs on the backend

![logs from the backend after launching two more servers](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/another_2_server_launching_backend_logs.png)

 - After launching the servers overview on the frontend

![overview of the frontend after launching the new servers](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/launched_4_servers.png)


## Step-6 (Attach XDP)
Just **click** on the `Attach XDP` button to attach the xdp program to the `bridge` interface which 

was selected from when creating the docker network.

 - Attaching Xdp overview

![overview of the frontend after attaching the xdp](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/attached_xdp.png)

 - Attaching Xdp logs from backend

![logs from backend after attaching the xdp](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/backend_logs_after_attaching_xdp.png)

## Step-7 (Testing Custom Load Balancer)
You must need to be on the root `lb-xdp-poc` directory to this.

 - Exec Into Client Container

![exec client container](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/exec_client1.png)

 - Try to request the load-balancer server

![requested for 1st time on the load-balancer server](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/first_curl_lb_Server.png)

 - Try to request to the load-balancer server for couple of times

### Second request

![requested for 2nd time on the load-balancer server](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/2nd_curl_lb_server.png)

### Third request
![requested for 3rd time on the load-balancer server](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/3rd_curl_lb_server.png)

### Fourth request
![requested for 4th time on the load-balancer server](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/4th_curl_lb_server.png)

### Fifth request
![requested for 5th time on the load-balancer server](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/5th_curl_lb_server.png)

 - Try to add a new server now

### Added a new server

![adding a new server](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/adding_new_server_after_attaching_xdp.png)

### Overview of the frontend after adding a new server

![adding a new server](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/table_Show_after_new_Server_adding.png)

 - Again sending request to the loadbalancing server
After adding the new server sending couple of request to the loadbalancing server

![after_added_new_server_req_res](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/after_Adding_new_server_then_curl.png)

## Step-8 (Cleanup)
First exit from the client container.

![exit from client container](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/exit_from_client.png)

If you don't on the root directory (`lb-xdp-poc` directory) then navigate to the root directory of the application. 

Then execute the command to clean up.
```
make remove_all
```
This will revome all of the pinned maps and the containers and also remove the docker network.

![clean_up logs](https://github.com/REZ-OAN/lb-xdp-poc/blob/main/images/remove_setup.png)