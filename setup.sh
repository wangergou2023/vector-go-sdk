#!/bin/bash

echo "Installing OpenCV"

go get -u -d gocv.io/x/gocv
cd $GOPATH/src/gocv.io/x/gocv
make install
