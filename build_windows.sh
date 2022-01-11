#!/bin/bash

go generate
export GOOS=windows
export GOARCH=amd64
go build -o gotit.exe