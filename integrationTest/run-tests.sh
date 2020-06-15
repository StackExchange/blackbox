#!/bin/bash

rm -rf /tmp/bbhome-* && go test "$@"
