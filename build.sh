#!/usr/bin/env bash


env GOOS=darwin GOARCH=amd64 go build -o log_detector_mac_amd64 main.go
env GOOS=darwin GOARCH=arm64 go build -o log_detector_mac_arm64 main.go
env GOOS=windows GOARCH=amd64 go build -o log_detector.exe main.go
env GOOS=linux GOARCH=amd64 go build -o log_detector_linux_64 main.go