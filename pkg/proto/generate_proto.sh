#!/usr/bin/env bash

protoc ./fingerprint.proto --go_out=plugins=grpc:.