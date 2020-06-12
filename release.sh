#!/bin/bash

go mod tidy

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/whg-linux-amd64.service main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o build/whg-darwin-amd64.service main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/whg-windows-amd64.service main.go
