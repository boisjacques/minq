#!/bin/sh
sudo ifconfig enp3s0f0 10.0.1.10/24
sudo ifconfig enp3s0f1 10.0.3.10/24
sudo route add -net 10.0.0.0 netmask 255.255.255.0 gw 10.0.1.1
sudo route add -net 10.0.2.0 netmask 255.255.255.0 gw 10.0.3.1