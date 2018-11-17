#!/usr/bin/env bash

protoc ./pb/fingerprint.proto --go_out=plugins=grpc:.