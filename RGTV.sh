#!/bin/bash
progdir=$(dirname "$0")
cd $progdir
./RGTV > /tmp/RGTV.log 2>&1
