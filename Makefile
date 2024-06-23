# Makefile

# Default CIDR if not provided
CIDR ?= 192.168.0.0/16

# Extract the network part from the CIDR to use in the network name
NETWORK_NAME = lb-net

# client ip default 192.168.0.3
CLIENT_IP ?= 192.168.0.3
CLIENT_NAME = client-server

# loadbalancer ip default 192.168.0.5
LB_IP ?= 192.168.0.5
LB_NAME = lb-server

## NETWORK
# Docker command to create the network
create_net:
	@echo "Creating Docker network with CIDR $(CIDR) and name $(NETWORK_NAME)"
	@docker network create --subnet=$(CIDR) $(NETWORK_NAME)

# get the bridge name which is created for the docker network
get_iface:
	@docker network ls --filter "name=$(NETWORK_NAME)" --format "{{.ID}}" | awk '{print "br-"$$1}'

# remove the docker network
remove_net:
	@echo "Removing Docker Network $(NETWORK_NAME)"
	@docker network rm -f $(NETWORK_NAME)

## CLIENT
# build client image
build_client:
	@echo "Building Client Image"
	@docker build -t client ./client

# exec to the container of client
exec_client :
	@docker exec -it $(CLIENT_NAME) bash

# run the client on specific network and ip
run_client :
	@echo "Launching Client On This IP $(CLIENT_IP) and Network $(NETWORK_NAME) with name $(CLIENT_NAME)"
	@docker run -d -it --net $(NETWORK_NAME) --ip $(CLIENT_IP) --name $(CLIENT_NAME) client

# get the client mac address
get_client_mac :
	@docker inspect $(CLIENT_NAME) | jq -r '.[0].NetworkSettings.Networks.[].MacAddress'
# remove client from the running container
remove_client:
	@echo "Removing $(CLIENT_NAME)"
	@docker stop $(CLIENT_NAME)
	@docker rm $(CLIENT_NAME)

## LOADBALANCER
# build load-balancer image
build_lb:
	@echo "Building Load-Balancer Image"
	@docker build -t lb ./loadbalancer

# exec to the container of load-balancer
exec_lb :
	@docker exec -it $(LB_NAME) bash

# run the load-balancer on specific network and ip
run_lb :
	@echo "Launching Load-Balancer On This IP $(LB_IP) and Network $(NETWORK_NAME) with name $(LB_NAME)"
	@docker run -d -it  --privileged --net $(NETWORK_NAME) --ip $(LB_IP) --name $(LB_NAME) lb

# get the load-balancer mac address
get_lb_mac :
	@docker inspect $(LB_NAME) | jq -r '.[0].NetworkSettings.Networks.[].MacAddress'

# remove load-balancer from the running container
remove_lb:
	@echo "Removing $(LB_NAME)"
	@docker stop $(LB_NAME)
	@docker rm $(LB_NAME)

## REMOVE ALL
# remove all configs
remove_all : remove_client remove_lb remove_net 


.PHONY: create_net get_iface remove_net build_client run_client exec_client remove_client build_lb run_lb exec_lb remove_lb remove_all
