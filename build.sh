#!/bin/sh
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=off go build -v -o gateway-device-service ./cmd