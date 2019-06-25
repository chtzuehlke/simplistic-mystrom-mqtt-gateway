#!/bin/bash

#below: works on my machine ;-)
IP=$(ifconfig -a | grep inet\ | grep -v 127.0.0.1 | ruby -ne 'print $1 if /^.*inet (.+) netmask.*/')
HEX_NETMASK=$(ifconfig -a | grep inet\ | grep -v 127.0.0.1 | ruby -ne 'print $1 if /^.*netmask (.+) broadcast.*/')
CIDR=$(ruby -e "require 'ipaddr'; puts IPAddr.new($HEX_NETMASK,Socket::AF_INET).to_i.to_s(2).count('1')")

echo "Scanning: $IP/$CIDR"

SWITCH_IPS=$(nmap --open -p 80 $IP/$CIDR -oG - | ruby -ne 'print "#{$1}," if /^Host\: (\d+\.\d+\.\d+\.\d+) .*myStrom-Switch.*Ports\: 80.*/')

echo "Starting MQTT gateway for mystrom switches: $SWITCH_IPS"

./simplistic-mystrom-mqtt-gateway -conf conf.json -client-id testclient -mystrom-switch-ips $SWITCH_IPS
