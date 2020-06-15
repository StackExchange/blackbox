#!/bin/bash

rm -rf /tmp/bbhome-* && BLACKBOX_DEBUG=true go test -verbose "$@"
