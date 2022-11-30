#!/bin/bash
whoami
sudo -i -u pi bash << EOF
cd /home/pi/vector-go-sdk
/usr/local/go/bin/go run cmd/examples/rolldie/main.go --serial 005070ac
EOF
echo "Out"
whoami
