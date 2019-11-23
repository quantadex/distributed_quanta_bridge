package main

// go build -ldflags "-X main.BuildStamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.GitHash=`git rev-parse HEAD`"
var (
	Version = "1.0"
	BuildStamp = ""
	GitHash = ""
)
