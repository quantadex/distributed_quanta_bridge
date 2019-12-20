#!/bin/sh

set -x
pwd

cd node
go build -ldflags "-X main.BuildStamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'`" # -X main.GitHash=`git rev-parse HEAD`"
cd ../cli/bitcoin
go build
cd ../../cli/ethereum
go build
cd ../../cli/litecoin
go build
cd ../../cli/bch
go build
#cd ../../cli/event_notifier
#go build
