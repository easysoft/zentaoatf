#!/usr/bin/env bash

rm -rf log/*
ztf run demo/lang demo/sample
ztf run demo/ztf/demo_test.sh