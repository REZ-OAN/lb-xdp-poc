# Makefile

# Default CIDR if not provided
CIDR ?= 192.168.0.0/16

# Extract the network part from the CIDR to use in the network name
NETWORK_NAME = lb-net

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
# Phony target to ensure make does not confuse this with a file
.PHONY: create_network
