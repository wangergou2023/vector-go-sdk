#!/bin/bash

echo "Installing dependencies for Text To Speech..."
sudo apt-get install gcc libasound2 libasound2-dev -y
sudo apt install espeak -y
