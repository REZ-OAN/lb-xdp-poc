#!/bin/bash

ip link set dev eth0 xdp off

ip link set dev eth0 xdp obj xdp_lb.o sec xdp

tail -f /dev/null

