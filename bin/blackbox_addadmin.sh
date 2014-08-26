#!/usr/bin/env bash

set -e

cd keyrings/live
touch blackbox-admins.txt
sort  -fdu -o blackbox-admins.txt <(echo "$1") blackbox-admins.txt
