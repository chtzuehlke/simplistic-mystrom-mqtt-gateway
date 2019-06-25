#!/bin/bash

SWITCH_IPS=$1

./simplistic-mystrom-mqtt-gateway -conf conf.json -client-id testclient -mystrom-switch-ips $SWITCH_IPS
