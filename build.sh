#!/bin/bash
env GOOS=linux   GOARCH=amd64 go build -o "mp3recover"
mv mp3recover ./linux_amd64

env GOOS=linux   GOARCH=386   go build -o "mp3recover"
mv mp3recover ./linux_386

env GOOS=windows GOARCH=amd64 go build -o "mp3recover"
mv mp3recover ./windows_amd64

env GOOS=windows GOARCH=386   go build -o "mp3recover"
mv mp3recover ./windows_386
