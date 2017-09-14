#!/bin/bash

cd ../../dpdk
sudo ./pktgen -c 0xffffffff -n 4 -- -P -m "[1-2:3-4].0, [5-6:7-8].1, [9-10:11-12].2, [13-14:15-16].3" -T
reset

