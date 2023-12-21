#!/bin/bash

CGO_LDFLAGS="-lpcre -static" CGO_ENABLED=1 go build  -tags 'netgo' --trimpath -o ../../bin/hellclientswitch ../