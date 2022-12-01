#!/bin/bash
cp /home/pi/vector-go-sdk/cmd/wirepod-ex/wirepod-ex.go /home/pi/wire-pod/chipper/plugins/
cd /home/pi/wire-pod/chipper
sudo /usr/local/go/bin/go get -u github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper
cd /home/pi/wire-pod/chipper/plugins/
sudo /usr/local/go/bin/go build -buildmode=plugin wirepod-ex.go
cd /home/pi/wire-pod/chipper
sudo ./start.sh
