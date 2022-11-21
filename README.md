# vector-go-sdk

A go SDK for Vector! This is an ***early alpha*** and not all features are present.  

## Usage

Currently. the SDK expects to find in your home directory the Python SDK data.
It consists of a .anki_vector folder where a sdk_config.ini file is placed. This file contains the configuration for each robot:

[ROBOT SERIAL NO]

cert = path to certificate (not needed by the go SDK)

ip = robot ip

name = robot name

guid = token for connecting to the robot

When you instantiate the robot you provide the serial number, then the SDK code will parse the .anki_vector/sdk_config.ini and get the data it needs.
I develop and test this version of the go SDK using a production Vector 1.0 and [wirepod](https://github.com/kercre123/wire-pod)

## Examples

Please see the cmd/examples directory for usage examples.

## Changelog 

***VERSION ALPHA_01***

I have taken the go SDK and did some changes.
- Integrated the functions of the great Wire's [vector-web-api-app](https://github.com/kercre123/vector-web-api-app) as sdk_wrapper functions
- Added many examples:
- The SDK supports features in these fields: 
  - Audio playback and volume settings
  - Custom eye colors
  - Animations
  - Images
  - Settings